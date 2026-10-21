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

package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	pluginAPI "github.com/coze-dev/coze-studio/backend/api/model/plugin_develop"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/dto"
	userentity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	mockPlugin "github.com/coze-dev/coze-studio/backend/internal/mock/domain/plugin"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

func TestGetDevPluginListUsesSessionOwner(t *testing.T) {
	ctx := ctxcache.Init(context.Background())
	ctxcache.Store(ctx, consts.SessionDataKeyInCtx, &userentity.Session{UserID: 7})
	domain := mockPlugin.NewMockPluginService(gomock.NewController(t))
	domain.EXPECT().ListDraftPlugins(ctx, gomock.Any()).DoAndReturn(
		func(_ context.Context, req *dto.ListDraftPluginsRequest) (*dto.ListDraftPluginsResponse, error) {
			require.Equal(t, int64(7), req.DeveloperID)
			require.Equal(t, int64(42), req.SpaceID)
			require.Zero(t, req.APPID)
			return &dto.ListDraftPluginsResponse{}, nil
		})
	svc := &PluginApplicationService{DomainSVC: domain}
	resp, err := svc.GetDevPluginList(ctx, &pluginAPI.GetDevPluginListRequest{
		DevID: 99, SpaceID: 42, Page: ptr.Of(int32(1)), Size: ptr.Of(int32(10)),
	})
	require.NoError(t, err)
	require.Empty(t, resp.PluginList)

	_, err = svc.GetDevPluginList(context.Background(), &pluginAPI.GetDevPluginListRequest{})
	require.Error(t, err)
}
