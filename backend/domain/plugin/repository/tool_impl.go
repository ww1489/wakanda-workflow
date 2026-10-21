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

package repository

import (
	"context"

	"fmt"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/crossdomain/plugin/model"
	pluginConf "github.com/coze-dev/coze-studio/backend/domain/plugin/conf"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/dto"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/internal/dal"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/internal/dal/query"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"

	"github.com/coze-dev/coze-studio/backend/pkg/lang/slices"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

type toolRepoImpl struct {
	query *query.Query

	pluginDraftDAO *dal.PluginDraftDAO

	toolDraftDAO   *dal.ToolDraftDAO
	toolDAO        *dal.ToolDAO
	toolVersionDAO *dal.ToolVersionDAO
}

type ToolRepoComponents struct {
	IDGen idgen.IDGenerator
	DB    *gorm.DB
}

func NewToolRepo(components *ToolRepoComponents) ToolRepository {
	return &toolRepoImpl{
		query:          query.Use(components.DB),
		pluginDraftDAO: dal.NewPluginDraftDAO(components.DB, components.IDGen),
		toolDraftDAO:   dal.NewToolDraftDAO(components.DB, components.IDGen),
		toolDAO:        dal.NewToolDAO(components.DB, components.IDGen),
		toolVersionDAO: dal.NewToolVersionDAO(components.DB, components.IDGen),
	}
}

func (t *toolRepoImpl) CreateDraftTool(ctx context.Context, tool *entity.ToolInfo) (toolID int64, err error) {
	return t.toolDraftDAO.Create(ctx, tool)
}

func (t *toolRepoImpl) UpsertDraftTools(ctx context.Context, pluginID int64, tools []*entity.ToolInfo) (err error) {
	apis := slices.Transform(tools, func(tool *entity.ToolInfo) dto.UniqueToolAPI {
		return dto.UniqueToolAPI{
			SubURL: tool.GetSubURL(),
			Method: tool.GetMethod(),
		}
	})

	existTools, err := t.toolDraftDAO.MGetWithAPIs(ctx, pluginID, apis, nil)
	if err != nil {
		return err
	}

	tx := t.query.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			if e := tx.Rollback(); e != nil {
				logs.CtxErrorf(ctx, "rollback failed, err=%v", e)
			}
			err = fmt.Errorf("catch panic: %v\nstack=%s", r, string(debug.Stack()))
			return
		}
		if err != nil {
			if e := tx.Rollback(); e != nil {
				logs.CtxErrorf(ctx, "rollback failed, err=%v", e)
			}
		}
	}()

	createdTools := make([]*entity.ToolInfo, 0, len(tools))
	updatedTools := make([]*entity.ToolInfo, 0, len(existTools))

	for _, tool := range tools {
		existTool, exist := existTools[dto.UniqueToolAPI{
			SubURL: tool.GetSubURL(),
			Method: tool.GetMethod(),
		}]
		if !exist {
			createdTools = append(createdTools, tool)
			continue
		}

		tool.ID = existTool.ID

		updatedTools = append(updatedTools, tool)
	}

	if len(createdTools) > 0 {
		_, err = t.toolDraftDAO.BatchCreateWithTX(ctx, tx, createdTools)
		if err != nil {
			return err
		}
	}

	if len(updatedTools) > 0 {
		err = t.toolDraftDAO.BatchUpdateWithTX(ctx, tx, updatedTools)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (t *toolRepoImpl) UpdateDraftTool(ctx context.Context, tool *entity.ToolInfo) (err error) {
	return t.toolDraftDAO.Update(ctx, tool)
}

func (t *toolRepoImpl) GetDraftTool(ctx context.Context, toolID int64) (tool *entity.ToolInfo, exist bool, err error) {
	return t.toolDraftDAO.Get(ctx, toolID)
}

func (t *toolRepoImpl) MGetDraftTools(ctx context.Context, toolIDs []int64, opts ...ToolSelectedOptions) (tools []*entity.ToolInfo, err error) {
	var opt *dal.ToolSelectedOption
	if len(opts) > 0 {
		opt = &dal.ToolSelectedOption{}
		for _, o := range opts {
			o(opt)
		}
	}
	return t.toolDraftDAO.MGet(ctx, toolIDs, opt)
}

func (t *toolRepoImpl) GetPluginAllDraftTools(ctx context.Context, pluginID int64, opts ...ToolSelectedOptions) (tools []*entity.ToolInfo, err error) {
	var opt *dal.ToolSelectedOption
	if len(opts) > 0 {
		opt = &dal.ToolSelectedOption{}
		for _, o := range opts {
			o(opt)
		}
	}
	return t.toolDraftDAO.GetAll(ctx, pluginID, opt)
}

func (t *toolRepoImpl) GetPluginAllOnlineTools(ctx context.Context, pluginID int64) (tools []*entity.ToolInfo, err error) {
	pi, exist := pluginConf.GetPluginProduct(pluginID)
	if exist {
		tis := pi.GetPluginAllTools()
		tools = slices.Transform(tis, func(ti *pluginConf.ToolInfo) *entity.ToolInfo {
			return ti.Info
		})

		return tools, nil
	}

	tools, err = t.toolDAO.GetAll(ctx, pluginID)
	if err != nil {
		return nil, err
	}

	return tools, nil
}

func (t *toolRepoImpl) ListPluginDraftTools(ctx context.Context, pluginID int64, pageInfo dto.PageInfo) (tools []*entity.ToolInfo, total int64, err error) {
	return t.toolDraftDAO.List(ctx, pluginID, pageInfo)
}

func (t *toolRepoImpl) GetDraftToolWithAPI(ctx context.Context, pluginID int64, api dto.UniqueToolAPI) (tool *entity.ToolInfo, exist bool, err error) {
	return t.toolDraftDAO.GetWithAPI(ctx, pluginID, api)
}

func (t *toolRepoImpl) MGetDraftToolWithAPI(ctx context.Context, pluginID int64, apis []dto.UniqueToolAPI, opts ...ToolSelectedOptions) (tools map[dto.UniqueToolAPI]*entity.ToolInfo, err error) {
	var opt *dal.ToolSelectedOption
	if len(opts) > 0 {
		opt = &dal.ToolSelectedOption{}
		for _, o := range opts {
			o(opt)
		}
	}
	return t.toolDraftDAO.MGetWithAPIs(ctx, pluginID, apis, opt)
}

func (t *toolRepoImpl) DeleteDraftTool(ctx context.Context, toolID int64) (err error) {
	return t.toolDraftDAO.Delete(ctx, toolID)
}

func (t *toolRepoImpl) GetOnlineTool(ctx context.Context, toolID int64) (tool *entity.ToolInfo, exist bool, err error) {
	ti, exist := pluginConf.GetToolProduct(toolID)
	if exist {
		return ti.Info, true, nil
	}

	return t.toolDAO.Get(ctx, toolID)
}

func (t *toolRepoImpl) MGetOnlineTools(ctx context.Context, toolIDs []int64, opts ...ToolSelectedOptions) (tools []*entity.ToolInfo, err error) {
	toolProducts := pluginConf.MGetToolProducts(toolIDs)

	tools = slices.Transform(toolProducts, func(tool *pluginConf.ToolInfo) *entity.ToolInfo {
		return tool.Info
	})
	productToolIDs := slices.ToMap(toolProducts, func(tool *pluginConf.ToolInfo) (int64, bool) {
		return tool.Info.ID, true
	})

	customToolIDs := make([]int64, 0, len(toolIDs))
	for _, id := range toolIDs {
		_, ok := productToolIDs[id]
		if ok {
			continue
		}
		customToolIDs = append(customToolIDs, id)
	}

	var opt *dal.ToolSelectedOption
	if len(opts) > 0 {
		opt = &dal.ToolSelectedOption{}
		for _, o := range opts {
			o(opt)
		}
	}

	customTools, err := t.toolDAO.MGet(ctx, customToolIDs, opt)
	if err != nil {
		return nil, err
	}

	tools = append(tools, customTools...)

	return tools, nil
}

func (t *toolRepoImpl) GetVersionTool(ctx context.Context, vTool model.VersionTool) (tool *entity.ToolInfo, exist bool, err error) {
	ti, exist := pluginConf.GetToolProduct(vTool.ToolID)
	if exist {
		return ti.Info, true, nil
	}

	return t.toolVersionDAO.Get(ctx, vTool)
}

func (t *toolRepoImpl) MGetVersionTools(ctx context.Context, versionTools []model.VersionTool) (tools []*entity.ToolInfo, err error) {
	tools, err = t.toolVersionDAO.MGet(ctx, versionTools)
	if err != nil {
		return nil, err
	}

	return tools, nil
}
