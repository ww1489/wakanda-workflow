/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package application

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/application/permission"

	"github.com/coze-dev/coze-studio/backend/application/base/appinfra"
	"github.com/coze-dev/coze-studio/backend/application/connector"
	"github.com/coze-dev/coze-studio/backend/application/memory"
	"github.com/coze-dev/coze-studio/backend/application/modelmgr"
	"github.com/coze-dev/coze-studio/backend/application/openauth"
	"github.com/coze-dev/coze-studio/backend/application/plugin"

	"github.com/coze-dev/coze-studio/backend/application/upload"
	"github.com/coze-dev/coze-studio/backend/application/user"
	"github.com/coze-dev/coze-studio/backend/application/workflow"
	crossconnector "github.com/coze-dev/coze-studio/backend/crossdomain/connector"
	connectorImpl "github.com/coze-dev/coze-studio/backend/crossdomain/connector/impl"
	crossdatabase "github.com/coze-dev/coze-studio/backend/crossdomain/database"
	databaseImpl "github.com/coze-dev/coze-studio/backend/crossdomain/database/impl"

	crosspermission "github.com/coze-dev/coze-studio/backend/crossdomain/permission"
	permissionImpl "github.com/coze-dev/coze-studio/backend/crossdomain/permission/impl"
	crossplugin "github.com/coze-dev/coze-studio/backend/crossdomain/plugin"
	pluginImpl "github.com/coze-dev/coze-studio/backend/crossdomain/plugin/impl"
	crossupload "github.com/coze-dev/coze-studio/backend/crossdomain/upload"
	uploadImpl "github.com/coze-dev/coze-studio/backend/crossdomain/upload/impl"
	crossuser "github.com/coze-dev/coze-studio/backend/crossdomain/user"
	crossuserImpl "github.com/coze-dev/coze-studio/backend/crossdomain/user/impl"

	crossworkflow "github.com/coze-dev/coze-studio/backend/crossdomain/workflow"
	workflowImpl "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/impl"

	"github.com/coze-dev/coze-studio/backend/infra/checkpoint"

	"github.com/coze-dev/coze-studio/backend/infra/sqlparser"
	sqlparserImpl "github.com/coze-dev/coze-studio/backend/infra/sqlparser/impl/sqlparser"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
)

type basicServices struct {
	infra        *appinfra.AppDependencies
	modelMgrSVC  *modelmgr.ModelmgrApplicationService
	connectorSVC *connector.ConnectorApplicationService
	userSVC      *user.UserApplicationService
	openAuthSVC  *openauth.OpenAuthApplicationService
	uploadSVC    *upload.UploadService

	permissionSVC *permission.PermissionApplicationService
}

type primaryServices struct {
	pluginSVC   *plugin.PluginApplicationService
	memorySVC   *memory.MemoryApplicationServices
	workflowSVC *workflow.ApplicationService
}

func Init(ctx context.Context) (err error) {
	ctx = ctxcache.Init(ctx)
	infra, err := appinfra.Init(ctx)
	if err != nil {
		return err
	}

	sqlparser.New = sqlparserImpl.NewSQLParser

	basicServices, err := initBasicServices(ctx, infra)
	if err != nil {
		return fmt.Errorf("Init - initBasicServices failed, err: %v", err)
	}

	primaryServices, err := initPrimaryServices(ctx, basicServices)
	if err != nil {
		return fmt.Errorf("Init - initPrimaryServices failed, err: %v", err)
	}

	// Initialize permission service first as it's required by other services
	crosspermission.SetDefaultSVC(permissionImpl.InitDomainService(basicServices.permissionSVC.DomainSVC))

	crossconnector.SetDefaultSVC(connectorImpl.InitDomainService(basicServices.connectorSVC.DomainSVC))
	crossdatabase.SetDefaultSVC(databaseImpl.InitDomainService(primaryServices.memorySVC.DatabaseDomainSVC))
	crossplugin.SetDefaultSVC(pluginImpl.InitDomainService(primaryServices.pluginSVC.DomainSVC, infra.OSS))
	crossworkflow.SetDefaultSVC(workflowImpl.InitDomainService(primaryServices.workflowSVC.DomainSVC))
	crossuser.SetDefaultSVC(crossuserImpl.InitDomainService(basicServices.userSVC.DomainSVC))
	crossupload.SetDefaultSVC(uploadImpl.InitDomainService(basicServices.uploadSVC.UploadSVC))

	return nil
}

// initBasicServices init basic services that only depends on infra.
func initBasicServices(ctx context.Context, infra *appinfra.AppDependencies) (*basicServices, error) {
	uploadSVC := upload.InitService(&upload.UploadComponents{Cache: infra.CacheCli, Oss: infra.OSS, DB: infra.DB, Idgen: infra.IDGenSVC})
	openAuthSVC := openauth.InitService(infra.DB, infra.IDGenSVC)
	modelMgrSVC := modelmgr.InitService(infra.OSS)
	connectorSVC := connector.InitService(infra.OSS)
	userSVC := user.InitService(ctx, infra.DB, infra.OSS, infra.IDGenSVC)

	permissionSVC := permission.InitService(&permission.ServiceComponents{})

	return &basicServices{
		infra:        infra,
		modelMgrSVC:  modelMgrSVC,
		connectorSVC: connectorSVC,
		userSVC:      userSVC,
		openAuthSVC:  openAuthSVC,
		uploadSVC:    uploadSVC,

		permissionSVC: permissionSVC,
	}, nil
}

// initPrimaryServices init primary services that depends on basic services.
func initPrimaryServices(ctx context.Context, basicServices *basicServices) (*primaryServices, error) {
	pluginSVC, err := plugin.InitService(ctx, basicServices.toPluginServiceComponents())
	if err != nil {
		return nil, err
	}

	memorySVC := memory.InitService(basicServices.toMemoryServiceComponents())

	workflowDomainSVC, err := workflow.InitService(ctx,
		basicServices.toWorkflowServiceComponents())
	if err != nil {
		return nil, err
	}

	return &primaryServices{
		pluginSVC:   pluginSVC,
		memorySVC:   memorySVC,
		workflowSVC: workflowDomainSVC,
	}, nil
}

func (b *basicServices) toPluginServiceComponents() *plugin.ServiceComponents {
	return &plugin.ServiceComponents{
		IDGen:    b.infra.IDGenSVC,
		DB:       b.infra.DB,
		OSS:      b.infra.OSS,
		UserSVC:  b.userSVC.DomainSVC,
		CacheCli: b.infra.CacheCli,
	}
}

func (b *basicServices) toMemoryServiceComponents() *memory.ServiceComponents {
	return &memory.ServiceComponents{
		IDGen:     b.infra.IDGenSVC,
		DB:        b.infra.DB,
		TosClient: b.infra.OSS,
		CacheCli:  b.infra.CacheCli,
	}
}

func (b *basicServices) toWorkflowServiceComponents() *workflow.ServiceComponents {
	return &workflow.ServiceComponents{
		IDGen:      b.infra.IDGenSVC,
		DB:         b.infra.DB,
		Cache:      b.infra.CacheCli,
		Tos:        b.infra.OSS,
		ImageX:     b.infra.ImageXClient,
		CPStore:    checkpoint.NewRedisStore(b.infra.CacheCli),
		CodeRunner: b.infra.CodeRunner,
	}
}
