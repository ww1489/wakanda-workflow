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

	"strconv"
	"strings"

	"golang.org/x/mod/semver"

	common "github.com/coze-dev/coze-studio/backend/api/model/plugin_develop/common"
	"github.com/coze-dev/coze-studio/backend/crossdomain/plugin/model"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func (p *pluginServiceImpl) GetPluginNextVersion(ctx context.Context, pluginID int64) (version string, err error) {
	const defaultVersion = "v1.0.0"

	pl, exist, err := p.pluginRepo.GetOnlinePlugin(ctx, pluginID)
	if err != nil {
		return "", errorx.Wrapf(err, "GetOnlinePlugin failed, pluginID=%d", pluginID)
	}
	if !exist {
		return defaultVersion, nil
	}

	parts := strings.Split(pl.GetVersion(), ".") // Remove the 'v' and split
	if len(parts) < 3 {
		logs.CtxWarnf(ctx, "invalid version format '%s'", pl.GetVersion())
		return defaultVersion, nil
	}

	patch, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		logs.CtxWarnf(ctx, "invalid version format '%s'", pl.GetVersion())
		return defaultVersion, nil
	}

	parts[2] = strconv.FormatInt(patch+1, 10)
	nextVersion := strings.Join(parts, ".")

	return nextVersion, nil
}

func (p *pluginServiceImpl) PublishPlugin(ctx context.Context, req *model.PublishPluginRequest) (err error) {
	draftPlugin, exist, err := p.pluginRepo.GetDraftPlugin(ctx, req.PluginID)
	if err != nil {
		return errorx.Wrapf(err, "GetDraftPlugin failed, pluginID=%d", req.PluginID)
	}
	if !exist {
		return errorx.New(errno.ErrPluginRecordNotFound)
	}

	err = p.checkToolsDebugStatus(ctx, req.PluginID)
	if err != nil {
		return err
	}

	onlinePlugin, exist, err := p.pluginRepo.GetOnlinePlugin(ctx, req.PluginID)
	if err != nil {
		return errorx.Wrapf(err, "GetOnlinePlugin failed, pluginID=%d", req.PluginID)
	}
	if exist && onlinePlugin.Version != nil {
		if semver.Compare(req.Version, *onlinePlugin.Version) != 1 {
			return errorx.New(errno.ErrPluginInvalidParamCode,
				errorx.KVf(errno.PluginMsgKey, "version must be greater than the online version '%s' and format like 'v1.0.0'",
					*onlinePlugin.Version))
		}
	}

	draftPlugin.Version = &req.Version
	draftPlugin.VersionDesc = &req.VersionDesc

	err = p.pluginRepo.PublishPlugin(ctx, draftPlugin)
	if err != nil {
		return errorx.Wrapf(err, "PublishPlugin failed, pluginID=%d", req.PluginID)
	}

	return nil
}

func (p *pluginServiceImpl) checkToolsDebugStatus(ctx context.Context, pluginID int64) (err error) {
	res, err := p.toolRepo.GetPluginAllDraftTools(ctx, pluginID,
		repository.WithToolID(),
		repository.WithToolDebugStatus(),
		repository.WithToolActivatedStatus(),
	)
	if err != nil {
		return errorx.Wrapf(err, "GetPluginAllDraftTools failed, pluginID=%d", pluginID)
	}

	if len(res) == 0 {
		return errorx.New(errno.ErrPluginToolsCheckFailed, errorx.KVf(errno.PluginMsgKey,
			"at least one activated tool is required in plugin"))
	}

	activatedTools := make([]*entity.ToolInfo, 0, len(res))
	for _, tool := range res {
		if tool.IsDeactivated() {
			continue
		}
		activatedTools = append(activatedTools, tool)
	}

	if len(activatedTools) == 0 {
		return errorx.New(errno.ErrPluginToolsCheckFailed, errorx.KVf(errno.PluginMsgKey,
			"at least one activated tool is required in plugin"))
	}

	for _, tool := range activatedTools {
		if tool.GetDebugStatus() != common.APIDebugStatus_DebugWaiting {
			continue
		}
		return errorx.New(errno.ErrPluginToolsCheckFailed, errorx.KVf(errno.PluginMsgKey,
			"tools in plugin have not debugged yet"))
	}

	return nil
}

func (p *pluginServiceImpl) CheckPluginToolsDebugStatus(ctx context.Context, pluginID int64) (err error) {
	return p.checkToolsDebugStatus(ctx, pluginID)
}
