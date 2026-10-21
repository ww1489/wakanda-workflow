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

package middleware

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/application/openauth"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/conv"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

const HeaderAuthorizationKey = "Authorization"

var needAuthPath = map[string]bool{
	"/v1/files/upload":             true,
	"/v1/workflow/run":             true,
	"/v1/workflow/stream_run":      true,
	"/v1/workflow/stream_resume":   true,
	"/v1/workflow/get_run_history": true,
}

func parseBearerAuthToken(authHeader string) string {
	if len(authHeader) == 0 {
		return ""
	}
	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) == 0 {
		return ""
	}

	return token
}

func isNeedOpenapiAuth(c *app.RequestContext) bool {
	return needAuthPath[string(c.URI().Path())]
}

func OpenapiAuthMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		requestAuthType := c.GetInt32(RequestAuthTypeStr)
		if requestAuthType != int32(RequestAuthTypeOpenAPI) {
			c.Next(ctx)
			return
		}

		// open api auth
		if len(c.Request.Header.Get(HeaderAuthorizationKey)) == 0 {
			httputil.InternalError(ctx, c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "missing authorization in header")))
			return
		}

		apiKey := parseBearerAuthToken(c.Request.Header.Get(HeaderAuthorizationKey))
		if len(apiKey) == 0 {
			httputil.InternalError(ctx, c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "missing api_key in request")))
			return
		}

		md5Hash := md5.Sum([]byte(apiKey))
		md5Key := hex.EncodeToString(md5Hash[:])
		apiKeyInfo, err := openauth.OpenAuthApplication.CheckPermission(ctx, md5Key)

		if err != nil {
			logs.CtxErrorf(ctx, "OpenAuthApplication.CheckPermission failed, err=%v", err)
			httputil.InternalError(ctx, c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", err.Error())))
			return
		}

		if apiKeyInfo == nil {
			httputil.InternalError(ctx, c,
				errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "api key invalid")))
			return
		}

		apiKeyInfo.ConnectorID = consts.APIConnectorID
		logs.CtxInfof(ctx, "OpenapiAuthMW: apiKeyInfo=%v", conv.DebugJsonToStr(apiKeyInfo))
		ctxcache.Store(ctx, consts.OpenapiAuthKeyInCtx, apiKeyInfo)
		err = openauth.OpenAuthApplication.UpdateLastUsedAt(ctx, apiKeyInfo.ID, apiKeyInfo.UserID)
		if err != nil {
			logs.CtxErrorf(ctx, "OpenAuthApplication.UpdateLastUsedAt failed, err=%v", err)
		}
		c.Next(ctx)
	}
}
