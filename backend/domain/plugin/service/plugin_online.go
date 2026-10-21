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
	"fmt"

	"github.com/coze-dev/coze-studio/backend/crossdomain/plugin/consts"
	"github.com/coze-dev/coze-studio/backend/crossdomain/plugin/model"

	"github.com/coze-dev/coze-studio/backend/domain/plugin/dto"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"

	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func (p *pluginServiceImpl) GetOnlinePlugin(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, err error) {
	pl, exist, err := p.pluginRepo.GetOnlinePlugin(ctx, pluginID)
	if err != nil {
		return nil, errorx.Wrapf(err, "GetOnlinePlugin failed, pluginID=%d", pluginID)
	}
	if !exist {
		return nil, errorx.New(errno.ErrPluginRecordNotFound)
	}

	return pl, nil
}

func (p *pluginServiceImpl) MGetOnlinePlugins(ctx context.Context, pluginIDs []int64) (plugins []*entity.PluginInfo, err error) {
	plugins, err = p.pluginRepo.MGetOnlinePlugins(ctx, pluginIDs)
	if err != nil {
		return nil, errorx.Wrapf(err, "MGetOnlinePlugins failed, pluginIDs=%v", pluginIDs)
	}

	return plugins, nil
}

func (p *pluginServiceImpl) GetOnlineTool(ctx context.Context, toolID int64) (tool *entity.ToolInfo, err error) {
	tool, exist, err := p.toolRepo.GetOnlineTool(ctx, toolID)
	if err != nil {
		return nil, errorx.Wrapf(err, "GetOnlineTool failed, toolID=%d", toolID)
	}
	if !exist {
		return nil, errorx.New(errno.ErrPluginRecordNotFound)
	}

	return tool, nil
}

func (p *pluginServiceImpl) MGetOnlineTools(ctx context.Context, toolIDs []int64) (tools []*entity.ToolInfo, err error) {
	tools, err = p.toolRepo.MGetOnlineTools(ctx, toolIDs)
	if err != nil {
		return nil, errorx.Wrapf(err, "MGetOnlineTools failed, toolIDs=%v", toolIDs)
	}

	return tools, nil
}

func (p *pluginServiceImpl) MGetVersionTools(ctx context.Context, versionTools []model.VersionTool) (tools []*entity.ToolInfo, err error) {
	tools, err = p.toolRepo.MGetVersionTools(ctx, versionTools)
	if err != nil {
		return nil, errorx.Wrapf(err, "MGetVersionTools failed, versionTools=%v", versionTools)
	}

	return tools, nil
}

func (p *pluginServiceImpl) MGetVersionPlugins(ctx context.Context, versionPlugins []model.VersionPlugin) (plugins []*entity.PluginInfo, err error) {
	plugins, err = p.pluginRepo.MGetVersionPlugins(ctx, versionPlugins)
	if err != nil {
		return nil, errorx.Wrapf(err, "MGetVersionPlugins failed, versionPlugins=%v", versionPlugins)
	}

	return plugins, nil
}

func (p *pluginServiceImpl) ListCustomOnlinePlugins(ctx context.Context, spaceID int64, pageInfo dto.PageInfo) (plugins []*entity.PluginInfo, total int64, err error) {
	plugins, total, err = p.pluginRepo.ListCustomOnlinePlugins(ctx, spaceID, pageInfo)
	if err != nil {
		return nil, 0, errorx.Wrapf(err, "ListCustomOnlinePlugins failed, spaceID=%d", spaceID)
	}
	return plugins, total, nil
}

func (p *pluginServiceImpl) MGetPluginLatestVersion(ctx context.Context, pluginIDs []int64) (resp *model.MGetPluginLatestVersionResponse, err error) {
	plugins, err := p.pluginRepo.MGetOnlinePlugins(ctx, pluginIDs,
		repository.WithPluginID(),
		repository.WithPluginVersion())
	if err != nil {
		return nil, errorx.Wrapf(err, "MGetOnlinePlugins failed, pluginIDs=%v", pluginIDs)
	}

	versions := make(map[int64]string, len(plugins))
	for _, pl := range plugins {
		versions[pl.ID] = pl.GetVersion()
	}

	resp = &model.MGetPluginLatestVersionResponse{
		Versions: versions,
	}

	return resp, nil
}

func (p *pluginServiceImpl) getCopySourcePluginAndTools(ctx context.Context, pluginID int64, scene consts.CopyScene) (plugin *entity.PluginInfo, tools []*entity.ToolInfo, err error) {
	switch scene {
	case consts.CopySceneOfToAPP:
		return p.getOnlinePluginAndTools(ctx, pluginID)
	case consts.CopySceneOfToLibrary, consts.CopySceneOfDuplicate, consts.CopySceneOfAPPDuplicate:
		return p.getDraftPluginAndTools(ctx, pluginID)
	default:
		return nil, nil, fmt.Errorf("unsupported copy scene '%s'", scene)
	}
}

func (p *pluginServiceImpl) getOnlinePluginAndTools(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, tools []*entity.ToolInfo, err error) {
	onlinePlugin, exist, err := p.pluginRepo.GetOnlinePlugin(ctx, pluginID)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		return nil, nil, errorx.New(errno.ErrPluginRecordNotFound)
	}

	onlineTools, err := p.toolRepo.GetPluginAllOnlineTools(ctx, pluginID)
	if err != nil {
		return nil, nil, err
	}

	return onlinePlugin, onlineTools, nil
}

func (p *pluginServiceImpl) getDraftPluginAndTools(ctx context.Context, pluginID int64) (plugin *entity.PluginInfo, tools []*entity.ToolInfo, err error) {
	draftPlugin, exist, err := p.pluginRepo.GetDraftPlugin(ctx, pluginID)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		return nil, nil, errorx.New(errno.ErrPluginRecordNotFound)
	}

	draftTools, err := p.toolRepo.GetPluginAllDraftTools(ctx, pluginID)
	if err != nil {
		return nil, nil, err
	}

	return draftPlugin, draftTools, nil
}

func (p *pluginServiceImpl) changePluginAndToolsInfoForMove(plugin *entity.PluginInfo,
	tools []*entity.ToolInfo) {

	const (
		defaultVersion     = "v0.0.1"
		defaultVersionDesc = "move to library"
	)

	plugin.Version = ptr.Of(defaultVersion)
	plugin.VersionDesc = ptr.Of(defaultVersionDesc)

	for _, tool := range tools {
		tool.Version = ptr.Of(defaultVersion)
	}

	plugin.APPID = nil
}
