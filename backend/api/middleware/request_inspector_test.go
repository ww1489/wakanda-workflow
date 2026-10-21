package middleware

import (
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/require"
)

func TestRemovedRoutesReachNotFoundWithoutAuthentication(t *testing.T) {
	h := server.New()
	h.Use(RequestInspectorMW())
	h.Use(func(ctx context.Context, c *app.RequestContext) {
		if c.GetInt32(RequestAuthTypeStr) != RequestAuthTypeStaticFile {
			c.AbortWithStatus(401)
			return
		}
		c.Next(ctx)
	})
	h.POST("/v1/workflow/run", func(_ context.Context, c *app.RequestContext) { c.Status(200) })
	h.POST("/api/workflow_api/create", func(_ context.Context, c *app.RequestContext) { c.Status(200) })
	h.NoRoute(func(_ context.Context, c *app.RequestContext) { c.Status(404) })
	for _, path := range []string{"/v3/chat", "/api/knowledge/dataset/list", "/open_api/knowledge/document/list"} {
		require.Equal(t, 404, ut.PerformRequest(h.Engine, "POST", path, nil).Code, path)
	}
	for _, path := range []string{"/v1/workflow/run", "/api/workflow_api/create"} {
		require.Equal(t, 401, ut.PerformRequest(h.Engine, "POST", path, nil).Code, path)
	}
}
