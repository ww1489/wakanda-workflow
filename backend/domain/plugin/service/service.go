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

package service

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/crossdomain/plugin/model"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/dto"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
)

//go:generate mockgen -destination ../../../internal/mock/domain/plugin/interface.go --package mockPlugin -source service.go
type PluginService interface {
	// Draft Plugin
	CreateDraftPlugin(ctx context.Context, req *dto.CreateDraftPluginRequest) (pluginID int64, err error)
	CreateDraftPluginWithCode(ctx context.Context, req *dto.CreateDraftPluginWithCodeRequest) (resp *dto.CreateDraftPluginWithCodeResponse, err error)
	GetDraftPlugin(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error)
	MGetDraftPlugins(ctx context.Context, pluginIDs []int64) (plugins []*entity.PluginInfo, err error)
	ListDraftPlugins(ctx context.Context, req *dto.ListDraftPluginsRequest) (resp *dto.ListDraftPluginsResponse, err error)
	UpdateDraftPlugin(ctx context.Context, plugin *dto.UpdateDraftPluginRequest) (err error)
	UpdateDraftPluginWithCode(ctx context.Context, req *dto.UpdateDraftPluginWithCodeRequest) (err error)
	DeleteDraftPlugin(ctx context.Context, pluginID int64) (err error)

	// Online Plugin
	PublishPlugin(ctx context.Context, req *model.PublishPluginRequest) (err error)
	GetOnlinePlugin(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error)
	MGetOnlinePlugins(ctx context.Context, pluginIDs []int64) (plugins []*entity.PluginInfo, err error)
	MGetPluginLatestVersion(ctx context.Context, pluginIDs []int64) (resp *model.MGetPluginLatestVersionResponse, err error)
	GetPluginNextVersion(ctx context.Context, pluginID int64) (version string, err error)
	MGetVersionPlugins(ctx context.Context, versionPlugins []model.VersionPlugin) (plugins []*entity.PluginInfo, err error)
	ListCustomOnlinePlugins(ctx context.Context, spaceID int64, pageInfo dto.PageInfo) (plugins []*entity.PluginInfo, total int64, err error)

	// Draft Tool
	MGetDraftTools(ctx context.Context, toolIDs []int64) (tools []*entity.ToolInfo, err error)
	UpdateDraftTool(ctx context.Context, req *dto.UpdateDraftToolRequest) (err error)
	ConvertToOpenapi3Doc(ctx context.Context, req *dto.ConvertToOpenapi3DocRequest) (resp *dto.ConvertToOpenapi3DocResponse)
	CreateDraftToolsWithCode(ctx context.Context, req *dto.CreateDraftToolsWithCodeRequest) (resp *dto.CreateDraftToolsWithCodeResponse, err error)
	CheckPluginToolsDebugStatus(ctx context.Context, pluginID int64) (err error)

	// Online Tool
	GetOnlineTool(ctx context.Context, toolID int64) (tool *entity.ToolInfo, err error)
	MGetOnlineTools(ctx context.Context, toolIDs []int64) (tools []*entity.ToolInfo, err error)
	MGetVersionTools(ctx context.Context, versionTools []model.VersionTool) (tools []*entity.ToolInfo, err error)

	// Agent Tool

	ExecuteTool(ctx context.Context, req *model.ExecuteToolRequest, opts ...model.ExecuteToolOpt) (resp *model.ExecuteToolResponse, err error)

	// Product

	GetOAuthStatus(ctx context.Context, userID, pluginID int64) (resp *dto.GetOAuthStatusResponse, err error)
	//Saas Plugin Product

	OAuthCode(ctx context.Context, code string, state *dto.OAuthState) (err error)
	GetAccessToken(ctx context.Context, oa *dto.OAuthInfo) (accessToken string, err error)
	RevokeAccessToken(ctx context.Context, meta *dto.AuthorizationCodeMeta) (err error)

	PluginOauthCallback(ctx context.Context, code, state, pluginID string) (confirmCode string, err error)
	GetPluginOauthInfo(ctx context.Context, confirmCode string) (info *dto.PluginOauthConfirmInfo, err error)
	ConfirmPluginOauth(ctx context.Context, confirmCode string, currentUserID int64) error
}
