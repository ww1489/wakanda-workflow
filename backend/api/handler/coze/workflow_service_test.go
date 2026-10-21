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

package coze

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/api/model/playground"
	pluginAPI "github.com/coze-dev/coze-studio/backend/api/model/plugin_develop"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"

	"github.com/coze-dev/coze-studio/backend/application/user"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	"github.com/coze-dev/coze-studio/backend/bizpkg/llm/modelbuilder"

	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"

	crossdatabase "github.com/coze-dev/coze-studio/backend/crossdomain/database"
	"github.com/coze-dev/coze-studio/backend/crossdomain/database/databasemock"

	crosspermission "github.com/coze-dev/coze-studio/backend/crossdomain/permission"
	"github.com/coze-dev/coze-studio/backend/crossdomain/permission/permissionmock"
	crossplugin "github.com/coze-dev/coze-studio/backend/crossdomain/plugin"
	pluginImpl "github.com/coze-dev/coze-studio/backend/crossdomain/plugin/impl"
	pluginmodel "github.com/coze-dev/coze-studio/backend/crossdomain/plugin/model"
	"github.com/coze-dev/coze-studio/backend/crossdomain/plugin/pluginmock"
	crossuser "github.com/coze-dev/coze-studio/backend/crossdomain/user"

	entity2 "github.com/coze-dev/coze-studio/backend/domain/openauth/openapiauth/entity"
	permission "github.com/coze-dev/coze-studio/backend/domain/permission"

	entity3 "github.com/coze-dev/coze-studio/backend/domain/plugin/entity"

	userentity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	workflow2 "github.com/coze-dev/coze-studio/backend/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/config"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/plugin"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/service"

	"github.com/coze-dev/coze-studio/backend/infra/cache/impl/redis"
	"github.com/coze-dev/coze-studio/backend/infra/checkpoint"
	"github.com/coze-dev/coze-studio/backend/infra/coderunner"
	mockCrossUser "github.com/coze-dev/coze-studio/backend/internal/mock/crossdomain/crossuser"
	mockPlugin "github.com/coze-dev/coze-studio/backend/internal/mock/domain/plugin"
	mockcode "github.com/coze-dev/coze-studio/backend/internal/mock/domain/workflow/crossdomain/code"
	mock "github.com/coze-dev/coze-studio/backend/internal/mock/infra/idgen"
	storageMock "github.com/coze-dev/coze-studio/backend/internal/mock/infra/storage"
	"github.com/coze-dev/coze-studio/backend/internal/testutil"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/slices"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
	"github.com/coze-dev/coze-studio/backend/types/consts"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func TestMain(m *testing.M) {
	callbacks.AppendGlobalHandlers(service.GetTokenCallbackHandler())
	service.RegisterAllNodeAdaptors()
	os.Exit(m.Run())
}

type wfTestRunner struct {
	t     *testing.T
	h     *server.Hertz
	ctrl  *gomock.Controller
	idGen *mock.MockIDGenerator

	plugin        *mockPlugin.MockPluginService
	tos           *storageMock.MockStorage
	database      *databasemock.MockDatabase
	pluginSrv     *pluginmock.MockPluginService
	internalModel *testutil.UTChatModel
	ctx           context.Context
	closeFn       func()
}

var req2URL = map[reflect.Type]string{
	reflect.TypeOf(&workflow.NodeTemplateListRequest{}):           "/api/workflow_api/node_template_list",
	reflect.TypeOf(&workflow.CreateWorkflowRequest{}):             "/api/workflow_api/create",
	reflect.TypeOf(&workflow.SaveWorkflowRequest{}):               "/api/workflow_api/save",
	reflect.TypeOf(&workflow.DeleteWorkflowRequest{}):             "/api/workflow_api/delete",
	reflect.TypeOf(&workflow.GetCanvasInfoRequest{}):              "/api/workflow_api/canvas",
	reflect.TypeOf(&workflow.WorkFlowTestRunRequest{}):            "/api/workflow_api/test_run",
	reflect.TypeOf(&workflow.CancelWorkFlowRequest{}):             "/api/workflow_api/cancel",
	reflect.TypeOf(&workflow.PublishWorkflowRequest{}):            "/api/workflow_api/publish",
	reflect.TypeOf(&workflow.OpenAPIRunFlowRequest{}):             "/v1/workflow/run",
	reflect.TypeOf(&workflow.ValidateTreeRequest{}):               "/api/workflow_api/validate_tree",
	reflect.TypeOf(&workflow.WorkflowTestResumeRequest{}):         "/api/workflow_api/test_resume",
	reflect.TypeOf(&workflow.WorkflowNodeDebugV2Request{}):        "/api/workflow_api/nodeDebug",
	reflect.TypeOf(&workflow.QueryWorkflowNodeTypeRequest{}):      "/api/workflow_api/node_type",
	reflect.TypeOf(&workflow.GetWorkFlowListRequest{}):            "/api/workflow_api/workflow_list",
	reflect.TypeOf(&workflow.UpdateWorkflowMetaRequest{}):         "/api/workflow_api/update_meta",
	reflect.TypeOf(&workflow.GetWorkflowDetailRequest{}):          "/api/workflow_api/workflow_detail",
	reflect.TypeOf(&workflow.GetWorkflowDetailInfoRequest{}):      "/api/workflow_api/workflow_detail_info",
	reflect.TypeOf(&workflow.GetLLMNodeFCSettingDetailRequest{}):  "/api/workflow_api/llm_fc_setting_detail",
	reflect.TypeOf(&workflow.GetLLMNodeFCSettingsMergedRequest{}): "/api/workflow_api/llm_fc_setting_merged",
	reflect.TypeOf(&workflow.CopyWorkflowRequest{}):               "/api/workflow_api/copy",
	reflect.TypeOf(&workflow.BatchDeleteWorkflowRequest{}):        "/api/workflow_api/batch_delete",
	reflect.TypeOf(&workflow.GetHistorySchemaRequest{}):           "/api/workflow_api/history_schema",
	reflect.TypeOf(&workflow.GetWorkflowReferencesRequest{}):      "/api/workflow_api/workflow_references",
}

func newWfTestRunner(t *testing.T) *wfTestRunner {
	h := server.Default()

	h.Use(func(c context.Context, ctx *app.RequestContext) {
		c = ctxcache.Init(c)
		ctxcache.Store(c, consts.SessionDataKeyInCtx, &userentity.Session{
			UserID: 123,
		})
		// Add API auth info for OpenAPI endpoints
		ctxcache.Store(c, consts.OpenapiAuthKeyInCtx, &entity2.ApiKey{
			UserID:      123,
			ConnectorID: consts.APIConnectorID,
		})
		ctx.Next(c)
	})
	h.POST("/api/workflow_api/node_template_list", NodeTemplateList)
	h.POST("/api/workflow_api/create", CreateWorkflow)
	h.POST("/api/workflow_api/save", SaveWorkflow)
	h.POST("/api/workflow_api/delete", DeleteWorkflow)
	h.POST("/api/workflow_api/canvas", GetCanvasInfo)
	h.POST("/api/workflow_api/test_run", WorkFlowTestRun)
	h.GET("/api/workflow_api/get_process", GetWorkFlowProcess)
	h.POST("/api/workflow_api/validate_tree", ValidateTree)
	h.POST("/api/workflow_api/test_resume", WorkFlowTestResume)
	h.POST("/api/workflow_api/publish", PublishWorkflow)
	h.POST("/api/workflow_api/update_meta", UpdateWorkflowMeta)
	h.POST("/api/workflow_api/cancel", CancelWorkFlow)
	h.POST("/api/workflow_api/workflow_list", GetWorkFlowList)
	h.POST("/api/workflow_api/workflow_detail", GetWorkflowDetail)
	h.POST("/api/workflow_api/workflow_detail_info", GetWorkflowDetailInfo)
	h.POST("/api/workflow_api/llm_fc_setting_detail", GetLLMNodeFCSettingDetail)
	h.POST("/api/workflow_api/llm_fc_setting_merged", GetLLMNodeFCSettingsMerged)
	h.POST("/v1/workflow/run", OpenAPIRunFlow)
	h.POST("/v1/workflow/stream_run", OpenAPIStreamRunFlow)
	h.POST("/v1/workflow/stream_resume", OpenAPIStreamResumeFlow)
	h.POST("/api/workflow_api/nodeDebug", WorkflowNodeDebugV2)
	h.GET("/api/workflow_api/get_node_execute_history", GetNodeExecuteHistory)
	h.POST("/api/workflow_api/copy", CopyWorkflow)
	h.POST("/api/workflow_api/batch_delete", BatchDeleteWorkflow)
	h.POST("/api/workflow_api/node_type", QueryWorkflowNodeTypes)
	h.GET("/v1/workflow/get_run_history", OpenAPIGetWorkflowRunHistory)
	h.POST("/api/workflow_api/history_schema", GetHistorySchema)
	h.POST("/api/workflow_api/workflow_references", GetWorkflowReferences)

	ctrl := gomock.NewController(t, gomock.WithOverridableExpectations())
	mockIDGen := mock.NewMockIDGenerator(ctrl)
	var previousID atomic.Int64
	mockIDGen.EXPECT().GenID(gomock.Any()).DoAndReturn(func(_ context.Context) (int64, error) {
		newID := time.Now().UnixNano()
		for {
			if newID == previousID.Load() {
				newID = time.Now().UnixNano()
			} else {
				previousID.Store(newID)
				break
			}
		}
		return newID, nil
	}).AnyTimes()
	mockIDGen.EXPECT().GenMultiIDs(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, count int) ([]int64, error) {
			ids := make([]int64, count)
			for i := 0; i < count; i++ {
				newID := time.Now().UnixNano()
				for {
					if newID == previousID.Load() {
						newID = time.Now().UnixNano()
					} else {
						previousID.Store(newID)
						break
					}
				}
				ids[i] = newID
			}
			return ids, nil
		}).AnyTimes()

	dsn := os.Getenv("WORKFLOW_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("WORKFLOW_TEST_MYSQL_DSN is required for workflow integration tests")
	}
	if os.Getenv("CI_JOB_NAME") != "" {
		dsn = strings.ReplaceAll(dsn, "127.0.0.1", "mysql")
	}
	db, err := gorm.Open(mysql.Open(dsn))
	assert.NoError(t, err)

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	redisClient := redis.NewWithAddrAndPassword(s.Addr(), "")

	cpStore := checkpoint.NewRedisStore(redisClient)

	utChatModel := &testutil.UTChatModel{}

	mockTos := storageMock.NewMockStorage(ctrl)
	mockTos.EXPECT().GetObjectUrl(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	workflowRepo, _ := service.NewWorkflowRepository(mockIDGen, db, redisClient, mockTos, cpStore, &config.WorkflowConfig{
		NodeOfCodeConfig: &config.NodeOfCodeConfig{},
	})
	mockey.Mock(appworkflow.GetWorkflowDomainSVC).Return(service.NewWorkflowService(workflowRepo)).Build()
	mockey.Mock(workflow2.GetRepository).Return(workflowRepo).Build()

	mockCU := mockCrossUser.NewMockUser(ctrl)
	mockCU.EXPECT().GetUserSpaceList(gomock.Any(), gomock.Any()).Return([]*crossuser.EntitySpace{
		{
			ID: 123,
		},
	}, nil).AnyTimes()

	mPlugin := mockPlugin.NewMockPluginService(ctrl)

	m := mockey.Mock(crossuser.DefaultSVC).Return(mockCU).Build()
	m1 := mockey.Mock(ctxutil.GetApiAuthFromCtx).Return(&entity2.ApiKey{
		UserID:      123,
		ConnectorID: consts.APIConnectorID,
	}).Build()
	m4 := mockey.Mock(ctxutil.MustGetUIDFromCtx).Return(int64(1)).Build()
	m5 := mockey.Mock(ctxutil.GetUIDFromCtx).Return(ptr.Of(int64(1))).Build()

	mockDatabaseOperator := databasemock.NewMockDatabase(ctrl)
	crossdatabase.SetDefaultSVC(mockDatabaseOperator)

	mockPluginSrv := pluginmock.NewMockPluginService(ctrl)
	crossplugin.SetDefaultSVC(mockPluginSrv)

	mockStorage := storageMock.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetObjectUrl(gomock.Any(), gomock.Any()).Return("URL_ADDRESS", nil).AnyTimes()
	plugin.SetOSS(mockStorage)

	// Initialize permission service for tests
	mockPermission := permissionmock.NewMockPermission(ctrl)
	mockPermission.EXPECT().CheckAuthz(gomock.Any(), gomock.Any()).Return(&permission.CheckAuthzResult{Decision: permission.Allow}, nil).AnyTimes()
	crosspermission.SetDefaultSVC(mockPermission)

	mockey.Mock((*user.UserApplicationService).MGetUserBasicInfo).Return(&playground.MGetUserBasicInfoResponse{
		UserBasicInfoMap: make(map[string]*playground.UserBasicInfo),
	}, nil).Build()

	f := func() {
		m.UnPatch()
		m1.UnPatch()

		m4.UnPatch()
		m5.UnPatch()
		ctrl.Finish()
		_ = h.Close()
	}

	return &wfTestRunner{
		t:             t,
		h:             h,
		ctrl:          ctrl,
		idGen:         mockIDGen,
		plugin:        mPlugin,
		tos:           mockTos,
		database:      mockDatabaseOperator,
		internalModel: utChatModel,
		ctx:           context.Background(),
		closeFn:       f,
		pluginSrv:     mockPluginSrv,
	}
}

type PostOption struct {
	Headers map[string]string
}
type PostOptionFn func(option *PostOption)

func WithHeaders(hds map[string]string) PostOptionFn {
	return func(option *PostOption) {
		if option.Headers == nil {
			option.Headers = map[string]string{}
		}
		for k, v := range hds {
			option.Headers[k] = v
		}
	}
}

func post[T any](r *wfTestRunner, req any, opts ...PostOptionFn) *T {
	// if req has a field SpaceID, set it's value to "123"
	opt := &PostOption{}
	for _, fn := range opts {
		fn(opt)
	}

	typ := reflect.TypeOf(req)
	if typ.Kind() == reflect.Ptr {
		typ1 := typ.Elem()
		spaceField, ok := typ1.FieldByName("SpaceID")
		if ok {
			if spaceField.Type == reflect.TypeOf("") {
				reflect.ValueOf(req).Elem().FieldByName("SpaceID").SetString("123")
			} else {
				reflect.ValueOf(req).Elem().FieldByName("SpaceID").Set(reflect.ValueOf(ptr.Of("123")))
			}
		}

	}

	url := req2URL[typ]
	m, err := sonic.Marshal(req)
	assert.NoError(r.t, err)

	headers := make([]ut.Header, 0)
	headers = append(headers, ut.Header{
		Key:   "Content-Type",
		Value: "application/json",
	})
	for k, v := range opt.Headers {
		headers = append(headers, ut.Header{Key: k, Value: v})
	}
	w := ut.PerformRequest(r.h.Engine, "POST", url, &ut.Body{Body: bytes.NewBuffer(m), Len: len(m)},
		headers...)
	res := w.Result()
	if res.StatusCode() != http.StatusOK {
		r.t.Fatalf("unexpected status code: %d, body: %s", res.StatusCode(), string(res.Body()))
	}
	rBody := res.Body()
	var status struct {
		Code int32  `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := sonic.Unmarshal(rBody, &status); err != nil || status.Code != 0 {
		r.t.Fatalf("request %s failed: %s", url, rBody)
	}
	var resp T
	err = sonic.Unmarshal(rBody, &resp)
	if err != nil {
		r.t.Fatalf("failed to unmarshal response body: %v", err)
	}
	return &resp
}

func (r *wfTestRunner) postWithError(req any) string {
	m, err := sonic.Marshal(req)
	assert.NoError(r.t, err)

	url := req2URL[reflect.TypeOf(req)]

	w := ut.PerformRequest(r.h.Engine, "POST", url, &ut.Body{Body: bytes.NewBuffer(m), Len: len(m)},
		ut.Header{Key: "Content-Type", Value: "application/json"})
	res := w.Result()
	if res.StatusCode() == http.StatusOK {
		r.t.Errorf("expected error, but got none")
	}
	return string(res.Body())
}

type loadOptions struct {
	name      string
	id        int64
	req       *workflow.CreateWorkflowRequest
	version   string
	projectID int64
	data      []byte
	mode      *workflow.WorkflowMode
}

func withWorkflowData(data []byte) func(*loadOptions) {
	return func(o *loadOptions) {
		o.data = data
	}
}

func withName(n string) func(*loadOptions) {
	return func(o *loadOptions) {
		o.name = n
	}
}
func withID(id int64) func(*loadOptions) {
	return func(o *loadOptions) {
		o.id = id
	}
}
func withProjectID(id int64) func(*loadOptions) {
	return func(o *loadOptions) {
		o.projectID = id
	}
}

func withPublish(version string) func(*loadOptions) {
	return func(o *loadOptions) {
		o.version = version
	}
}

func withMode(mode workflow.WorkflowMode) func(*loadOptions) {
	return func(o *loadOptions) {
		o.mode = ptr.Of(mode)
	}
}

func (r *wfTestRunner) load(schemaFile string, opts ...func(*loadOptions)) string {
	loadOpts := &loadOptions{}
	for _, opt := range opts {
		opt(loadOpts)
	}

	if loadOpts.id > 0 {
		_, err := appworkflow.GetWorkflowDomainSVC().Get(context.Background(), &vo.GetPolicy{
			ID:       loadOpts.id,
			MetaOnly: true,
		})
		if err == nil {
			return strconv.FormatInt(loadOpts.id, 10)
		} else {
			r.idGen.EXPECT().GenID(gomock.Any()).DoAndReturn(func(_ context.Context) (int64, error) {
				return loadOpts.id, nil
			}).Times(3)
			defer func() {
				var previousID atomic.Int64
				r.idGen.EXPECT().GenID(gomock.Any()).DoAndReturn(func(_ context.Context) (int64, error) {
					newID := time.Now().UnixNano()
					if newID == previousID.Load() {
						newID = previousID.Add(1)
					}
					return newID, nil
				}).AnyTimes()
			}()
		}
	}

	var createReq *workflow.CreateWorkflowRequest
	if loadOpts.req != nil {
		createReq = loadOpts.req
	} else {
		name := "test_wf"
		if loadOpts.name != "" {
			name = loadOpts.name
		}

		createReq = &workflow.CreateWorkflowRequest{
			Name:     name,
			Desc:     "this is a test wf",
			IconURI:  "icon/uri",
			SpaceID:  "123",
			FlowMode: ptr.Of(workflow.WorkflowMode_Workflow),
		}

		if loadOpts.projectID > 0 {
			createReq.ProjectID = ptr.Of(strconv.FormatInt(loadOpts.projectID, 10))
		}
	}

	createReq.FlowMode = loadOpts.mode

	resp := post[workflow.CreateWorkflowResponse](r, createReq)

	idStr := resp.Data.WorkflowID
	_, err := strconv.ParseInt(idStr, 10, 64)
	assert.NoError(r.t, err)

	var data []byte
	if len(loadOpts.data) > 0 {
		data = loadOpts.data
	} else {
		data, err = os.ReadFile(fmt.Sprintf("../../../domain/workflow/internal/canvas/examples/%s", schemaFile))
		assert.NoError(r.t, err)
	}

	saveReq := &workflow.SaveWorkflowRequest{
		WorkflowID: idStr,
		Schema:     ptr.Of(string(data)),
		SpaceID:    ptr.Of("123"),
	}

	_ = post[workflow.SaveWorkflowResponse](r, saveReq)

	if loadOpts.version != "" {
		r.publish(idStr, loadOpts.version, true)
	}

	return idStr
}

func getProcess(t *testing.T, h *server.Hertz, idStr string, exeID string) *workflow.GetWorkflowProcessResponse {
	getProcessReq := &workflow.GetWorkflowProcessRequest{
		WorkflowID: idStr,
		SpaceID:    "123",
		ExecuteID:  ptr.Of(exeID),
	}

	w := ut.PerformRequest(h.Engine, "GET", fmt.Sprintf("/api/workflow_api/get_process?workflow_id=%s&space_id=%s&execute_id=%s", getProcessReq.WorkflowID, getProcessReq.SpaceID, *getProcessReq.ExecuteID), nil,
		ut.Header{Key: "Content-Type", Value: "application/json"})
	res := w.Result()
	if res.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status code: %d, body: %s", res.StatusCode(), string(res.Body()))
	}
	getProcessResp := &workflow.GetWorkflowProcessResponse{}
	err := sonic.Unmarshal(res.Body(), getProcessResp)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	return getProcessResp
}

func (r *wfTestRunner) getNodeExeHistory(id string, exeID string, nodeID string, scene *workflow.NodeHistoryScene) *workflow.NodeResult {
	getNodeExeHistoryReq := &workflow.GetNodeExecuteHistoryRequest{
		WorkflowID:       id,
		SpaceID:          "123",
		ExecuteID:        exeID,
		NodeID:           nodeID,
		NodeHistoryScene: scene,
	}

	w := ut.PerformRequest(r.h.Engine, "GET", fmt.Sprintf("/api/workflow_api/get_node_execute_history?workflow_id=%s&space_id=%s&execute_id=%s"+
		"&node_id=%s&node_type=3&node_history_scene=%d", getNodeExeHistoryReq.WorkflowID, getNodeExeHistoryReq.SpaceID, getNodeExeHistoryReq.ExecuteID,
		getNodeExeHistoryReq.NodeID, getNodeExeHistoryReq.GetNodeHistoryScene()), nil,
		ut.Header{Key: "Content-Type", Value: "application/json"})

	res := w.Result()
	assert.Equal(r.t, http.StatusOK, res.StatusCode())
	getNodeResultResp := &workflow.GetNodeExecuteHistoryResponse{}
	err := sonic.Unmarshal(res.Body(), getNodeResultResp)
	assert.NoError(r.t, err)

	return getNodeResultResp.Data
}

func (r *wfTestRunner) getOpenAPIProcess(id string, exeID string) *workflow.GetWorkflowRunHistoryResponse {
	w := ut.PerformRequest(r.h.Engine, "GET", fmt.Sprintf("/v1/workflow/get_run_history?workflow_id=%s&execute_id=%s", id, exeID), nil,
		ut.Header{Key: "Content-Type", Value: "application/json"})
	res := w.Result()
	assert.Equal(r.t, http.StatusOK, res.StatusCode())
	getProcessResp := &workflow.GetWorkflowRunHistoryResponse{}
	err := sonic.Unmarshal(res.Body(), getProcessResp)
	assert.NoError(r.t, err)

	return getProcessResp
}

func mustUnmarshalToMap(t *testing.T, s string) map[string]any {
	r := make(map[string]any)
	err := sonic.UnmarshalString(s, &r)
	if err != nil {
		t.Fatal(err)
	}

	return r
}

func mustMarshalToString(t *testing.T, m any) string {
	b, err := sonic.MarshalString(m)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

type runOption struct {
	ProjectID *int64
	BotID     *int64
}
type RunOptionFun func(options *runOption)

func withRunProjectID(pID int64) RunOptionFun {
	return func(options *runOption) {
		options.ProjectID = &pID
	}
}

func withRunBotID(bID int64) RunOptionFun {
	return func(options *runOption) {
		options.BotID = &bID
	}
}

func (r *wfTestRunner) testRun(id string, input map[string]string, opts ...RunOptionFun) string {
	opt := &runOption{}
	for _, o := range opts {
		o(opt)
	}
	testRunReq := &workflow.WorkFlowTestRunRequest{
		WorkflowID: id,
		Input:      input,
	}

	if opt.ProjectID != nil {
		testRunReq.ProjectID = ptr.Of(strconv.FormatInt(ptr.From(opt.ProjectID), 10))
	}

	if opt.BotID != nil {
		testRunReq.BotID = ptr.Of(strconv.FormatInt(ptr.From(opt.BotID), 10))
	}

	testRunResponse := post[workflow.WorkFlowTestRunResponse](r, testRunReq)
	return testRunResponse.Data.ExecuteID
}

type getProcessOptions struct {
	previousInterruptEventID string
	specificNodeID           string
}

func withPreviousEventID(id string) func(options *getProcessOptions) {
	return func(options *getProcessOptions) {
		options.previousInterruptEventID = id
	}
}
func withSpecificNodeID(id string) func(options *getProcessOptions) {
	return func(options *getProcessOptions) {
		options.specificNodeID = id
	}
}

type exeResult struct {
	output      string
	status      workflow.WorkflowExeStatus
	event       *workflow.NodeEvent
	token       *workflow.TokenAndCost
	t           *testing.T
	reason      string
	nodeResults []*workflow.NodeResult
}

func (e *exeResult) assertSuccess() {
	assert.Equal(e.t, workflow.WorkflowExeStatus_Success, e.status)
}

func (e *exeResult) tokenEqual(in, out int) {
	assert.NotNil(e.t, e.token)
	input := strings.TrimSuffix(*e.token.InputTokens, " Tokens")
	output := strings.TrimSuffix(*e.token.OutputTokens, " Tokens")
	inputI, err := strconv.Atoi(input)
	assert.NoError(e.t, err)
	outputI, err := strconv.Atoi(output)
	assert.NoError(e.t, err)
	assert.Equal(e.t, in, inputI)
	assert.Equal(e.t, out, outputI)
}

func (e *exeResult) nodeResultHasResponseExtra(nodeID string, k string, v any) {
	var n *workflow.NodeResult

	for _, nr := range e.nodeResults {
		if nr.NodeId == nodeID {
			n = nr
			break
		}
	}

	if n == nil {
		e.t.Fatal("node key: ", nodeID, "not found")
		return
	}
	extra := mustUnmarshalToMap(e.t, n.Extra)
	assert.NotEmpty(e.t, extra)
	assert.Contains(e.t, extra, "response_extra")
	assert.Equal(e.t, extra["response_extra"].(map[string]any)[k], v)
}

func (r *wfTestRunner) getProcess(id, exeID string, opts ...func(options *getProcessOptions)) *exeResult {
	options := &getProcessOptions{}
	for _, opt := range opts {
		opt(options)
	}

	workflowStatus := workflow.WorkflowExeStatus_Running
	var output string
	var nodeEvent *workflow.NodeEvent
	var eventID string
	var nodeType string
	var token *workflow.TokenAndCost
	var reason string
	var nodeResults []*workflow.NodeResult
	var count int
	for {
		if nodeEvent != nil {
			if options.previousInterruptEventID != "" {
				if options.previousInterruptEventID != nodeEvent.ID {
					break
				}
			} else {
				break
			}
		}

		if workflowStatus != workflow.WorkflowExeStatus_Running {
			break
		}

		if count > 1000 {
			r.t.Fatal("get process for too long")
		}

		getProcessResp := getProcess(r.t, r.h, id, exeID)
		if len(getProcessResp.Data.NodeResults) == 1 {
			output = getProcessResp.Data.NodeResults[0].Output
			nodeType = getProcessResp.Data.NodeResults[0].NodeType
		} else {
			for _, ns := range getProcessResp.Data.NodeResults {
				if options.specificNodeID != "" {
					if ns.NodeId == options.specificNodeID {
						output = ns.Output
						nodeType = ns.NodeType
						break
					}
				} else if ns.NodeType == workflow.NodeTemplateType_End.String() {
					output = ns.Output
					nodeType = ns.NodeType
				}
			}
		}
		if len(getProcessResp.Data.NodeEvents) > 0 {
			nodeEvent = getProcessResp.Data.NodeEvents[len(getProcessResp.Data.NodeEvents)-1]
		}

		workflowStatus = getProcessResp.Data.ExecuteStatus
		token = getProcessResp.Data.TokenAndCost

		if getProcessResp.Data.Reason != nil {
			reason = *getProcessResp.Data.Reason
		}

		if nodeEvent != nil {
			eventID = nodeEvent.ID
		}

		nodeResults = getProcessResp.Data.NodeResults

		r.t.Logf("getProcess output= %s, status= %v, eventID= %s, nodeType= %s", output, workflowStatus, eventID, nodeType)

		count++
	}

	return &exeResult{
		output:      output,
		status:      workflowStatus,
		event:       nodeEvent,
		token:       token,
		t:           r.t,
		reason:      reason,
		nodeResults: nodeResults,
	}
}

func (r *wfTestRunner) cancel(id, exeID string) {
	cancelReq := &workflow.CancelWorkFlowRequest{
		WorkflowID: &id,
		ExecuteID:  exeID,
	}
	_ = post[workflow.CancelWorkFlowResponse](r, cancelReq)
}

func (r *wfTestRunner) publish(id string, version string, force bool) {
	publishReq := &workflow.PublishWorkflowRequest{
		WorkflowID:         id,
		WorkflowVersion:    ptr.Of(version),
		VersionDescription: ptr.Of("desc"),
		Force:              ptr.Of(force),
	}
	_ = post[workflow.PublishWorkflowResponse](r, publishReq)
}

func (r *wfTestRunner) openapiAsyncRun(id string, input any) string {
	runReq := &workflow.OpenAPIRunFlowRequest{
		WorkflowID: id,
		Parameters: ptr.Of(mustMarshalToString(r.t, input)),
		IsAsync:    ptr.Of(true),
	}

	runResp := post[workflow.OpenAPIRunFlowResponse](r, runReq)
	return runResp.GetExecuteID()
}

func (r *wfTestRunner) openapiSyncRun(id string, input any, opts ...RunOptionFun) (map[string]any, string) {
	opt := &runOption{}
	for _, o := range opts {
		o(opt)
	}

	runReq := &workflow.OpenAPIRunFlowRequest{
		WorkflowID: id,
		Parameters: ptr.Of(mustMarshalToString(r.t, input)),
		IsAsync:    ptr.Of(false),
	}

	if opt.ProjectID != nil {
		runReq.ProjectID = ptr.Of(strconv.FormatInt(ptr.From(opt.ProjectID), 10))
	}

	if opt.BotID != nil {
		runReq.BotID = ptr.Of(strconv.FormatInt(ptr.From(opt.BotID), 10))
	}

	runResp := post[workflow.OpenAPIRunFlowResponse](r, runReq)
	output := runResp.GetData()
	var m map[string]any
	err := sonic.UnmarshalString(output, &m)
	assert.NoError(r.t, err)
	return m, runResp.GetExecuteID()
}

func (r *wfTestRunner) validateTree(schema string) [][]*workflow.ValidateErrorData {
	data, err := os.ReadFile(fmt.Sprintf("../../../domain/workflow/internal/canvas/examples/%s", schema))
	if err != nil {
		r.t.Fatal(err)
	}

	res := post[workflow.ValidateTreeResponse](r, &workflow.ValidateTreeRequest{
		WorkflowID: "1",
		Schema:     ptr.Of(string(data)),
	})

	if len(res.Data) == 0 {
		return nil
	}

	var errs [][]*workflow.ValidateErrorData
	for _, d := range res.Data {
		errs = append(errs, d.Errors)
	}

	return errs
}

func (r *wfTestRunner) testResume(id string, exeID string, eventID string, input any) {
	inputStr, ok := input.(string)
	if !ok {
		inputStr = mustMarshalToString(r.t, input)
	}

	testResumeReq := &workflow.WorkflowTestResumeRequest{
		WorkflowID: id,
		SpaceID:    ptr.Of("123"),
		ExecuteID:  exeID,
		EventID:    eventID,
		Data:       inputStr,
	}

	_ = post[workflow.WorkflowTestResumeResponse](r, testResumeReq)
}

type nodeDebugOptions struct {
	input   map[string]string
	batch   map[string]string
	setting map[string]string
}

func withNDInput(input map[string]string) func(*nodeDebugOptions) {
	return func(options *nodeDebugOptions) {
		options.input = input
	}
}

func withNDBatch(batch map[string]string) func(*nodeDebugOptions) {
	return func(options *nodeDebugOptions) {
		options.batch = batch
	}
}

func withNDSettings(settings map[string]string) func(*nodeDebugOptions) {
	return func(options *nodeDebugOptions) {
		options.setting = settings
	}
}

func (r *wfTestRunner) nodeDebug(id string, nodeID string, opts ...func(*nodeDebugOptions)) string {
	options := &nodeDebugOptions{}
	for _, opt := range opts {
		opt(options)
	}

	nodeDebugReq := &workflow.WorkflowNodeDebugV2Request{
		WorkflowID: id,
		NodeID:     nodeID,
	}

	if options.input != nil {
		nodeDebugReq.Input = options.input
	}

	if options.batch != nil {
		nodeDebugReq.Batch = options.batch
	}

	if options.setting != nil {
		nodeDebugReq.Setting = options.setting
	}

	nodeDebugResp := post[workflow.WorkflowNodeDebugV2Response](r, nodeDebugReq)
	return nodeDebugResp.Data.ExecuteID
}

func (r *wfTestRunner) save(id string, schema string) {
	data, err := os.ReadFile(fmt.Sprintf("../../../domain/workflow/internal/canvas/examples/%s", schema))
	assert.NoError(r.t, err)

	saveReq := &workflow.SaveWorkflowRequest{
		WorkflowID: id,
		Schema:     ptr.Of(string(data)),
	}

	_ = post[workflow.SaveWorkflowResponse](r, saveReq)
}

func getCanvas(ctx context.Context, id string) (string, error) {
	response, err := appworkflow.SVC.GetCanvasInfo(ctx, &workflow.GetCanvasInfoRequest{
		SpaceID:    "123",
		WorkflowID: ptr.Of(id),
	})
	if err != nil {
		return "", err
	}
	return response.GetData().GetWorkflow().GetSchemaJSON(), nil

}

func (r *wfTestRunner) openapiStream(id string, input any) *sse.Reader {
	inputStr, _ := sonic.MarshalString(input)

	req := &workflow.OpenAPIRunFlowRequest{
		WorkflowID: id,
		Parameters: ptr.Of(inputStr),
	}

	m, err := sonic.Marshal(req)
	assert.NoError(r.t, err)

	c, _ := client.NewClient()
	hReq, hResp := protocol.AcquireRequest(), protocol.AcquireResponse()
	hReq.SetRequestURI("http://localhost:8888" + "/v1/workflow/stream_run")
	hReq.SetMethod("POST")
	hReq.SetBody(m)
	hReq.SetHeader("Content-Type", "application/json")
	// Add Authorization header for API authentication
	hReq.SetHeader("Authorization", "Bearer test-api-key-123")
	err = c.Do(context.Background(), hReq, hResp)
	assert.NoError(r.t, err)

	if hResp.StatusCode() != http.StatusOK {
		r.t.Errorf("unexpected status code: %d, body: %s", hResp.StatusCode(), string(hResp.Body()))
	}

	re, err := sse.NewReader(hResp)
	assert.NoError(r.t, err)

	return re
}

func (r *wfTestRunner) openapiResume(id string, eventID string, resumeData string) *sse.Reader {
	req := &workflow.OpenAPIStreamResumeFlowRequest{
		WorkflowID:  id,
		EventID:     eventID,
		ResumeData:  resumeData,
		ConnectorID: ptr.Of(strconv.FormatInt(consts.APIConnectorID, 10)),
	}

	m, err := sonic.Marshal(req)
	assert.NoError(r.t, err)

	c, _ := client.NewClient()
	hReq, hResp := protocol.AcquireRequest(), protocol.AcquireResponse()
	hReq.SetRequestURI("http://localhost:8888" + "/v1/workflow/stream_resume")
	hReq.SetMethod("POST")
	hReq.SetBody(m)
	hReq.SetHeader("Content-Type", "application/json")
	// Add Authorization header for API authentication
	hReq.SetHeader("Authorization", "Bearer test-api-key-123")
	err = c.Do(context.Background(), hReq, hResp)
	assert.NoError(r.t, err)

	if hResp.StatusCode() != http.StatusOK {
		r.t.Errorf("unexpected status code: %d, body: %s", hResp.StatusCode(), string(hResp.Body()))
	}

	re, err := sse.NewReader(hResp)
	assert.NoError(r.t, err)

	return re
}

func (r *wfTestRunner) runServer() func() {
	go func() {
		_ = r.h.Run()
	}()

	return func() {
		_ = r.h.Close()
	}
}

func TestNodeTemplateList(t *testing.T) {
	mockey.PatchConvey("test node cn template list", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		resp := post[workflow.NodeTemplateListResponse](r, &workflow.NodeTemplateListRequest{
			NodeTypes: []string{"3", "5", "18"},
		})

		assert.Equal(t, 3, len(resp.Data.TemplateList))
		assert.Equal(t, 3, len(resp.Data.CateList))

		id2Name := map[string]string{
			"3":  "LLM",
			"5":  "Code",
			"18": "Question",
		}
		for _, tl := range resp.Data.TemplateList {
			assert.Equal(t, tl.Name, id2Name[tl.ID])
		}

	})
	mockey.PatchConvey("test node en template list", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		resp := post[workflow.NodeTemplateListResponse](r, &workflow.NodeTemplateListRequest{
			NodeTypes: []string{"3", "5", "18"},
		}, WithHeaders(map[string]string{
			"x-locale": "en-US",
		}))

		id2Name := map[string]string{
			"3":  "LLM",
			"5":  "Code",
			"18": "Question",
		}
		assert.Equal(t, 3, len(resp.Data.TemplateList))
		assert.Equal(t, 3, len(resp.Data.CateList))

		for _, tl := range resp.Data.TemplateList {
			assert.Equal(t, tl.Name, id2Name[tl.ID])
		}

	})

}

func TestValidateTree(t *testing.T) {
	mockey.PatchConvey("test validate tree", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		t.Run("workflow_has_loop", func(t *testing.T) {
			errs := r.validateTree("validate/workflow_has_loop.json")

			paths := map[string]string{
				"161668": "101917",
				"101917": "177387",
				"177387": "161668",
				"166209": "102541",
				"102541": "109507",
				"109507": "166209",
			}

			for _, i := range errs[0] {
				if assert.NotNil(t, i.PathError, "unexpected validation error: %s", i.Message) {
					assert.Equal(t, paths[i.PathError.Start], i.PathError.End)
				}
			}
		})

		t.Run("workflow_has_no_connected_nodes", func(t *testing.T) {
			errs := r.validateTree("validate/workflow_has_no_connected_nodes.json")

			for _, i := range errs[0] {
				if i.NodeError != nil {
					if i.NodeError.NodeID == "108984" {
						assert.Equal(t, i.Message, `node "代码_1" not connected`)
					}
					if i.NodeError.NodeID == "160892" {
						assert.Contains(t, i.Message, `node "意图识别"'s port "branch_1" not connected`)
						assert.Contains(t, i.Message, `node "意图识别"'s port "default" not connected`)
					}

				}
			}
		})

		t.Run("workflow_ref_variable", func(t *testing.T) {
			errs := r.validateTree("validate/workflow_ref_variable.json")

			for _, i := range errs[0] {
				if i.NodeError != nil {
					if i.NodeError.NodeID == "118685" {
						assert.Equal(t, i.Message, `the node id "118685" on which node id "165568" depends does not exist`)
					}

					if i.NodeError.NodeID == "128176" {
						assert.Equal(t, i.Message, `the node id "128176" on which node id "11384000" depends does not exist`)
					}
				}
			}
		})

		t.Run("workflow_nested_has_loop_or_batch", func(t *testing.T) {
			errs := r.validateTree("validate/workflow_nested_has_loop_or_batch.json")

			assert.Equal(t, errs[0][0].Message, `composite nodes such as batch/loop cannot be nested`)
		})

		t.Run("invalid_input_parameter", func(t *testing.T) {
			errs := r.validateTree("validate/invalid_input_parameter.json")
			assert.Equal(t, len(errs[0]), 2)
			msgs := slices.Transform(errs[0], func(item *workflow.ValidateErrorData) string {
				return item.Message
			})
			assert.Contains(t, msgs, `parameter name only allows number or alphabet, and must begin with alphabet, but it's "123"`)
			assert.Contains(t, msgs, `ref block error,[blockID] is empty`)
		})

	})
}

func TestQueryTypes(t *testing.T) {
	mockey.PatchConvey("test workflow node types", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		t.Run("not sub workflow", func(t *testing.T) {
			id := r.load("query_types/llm_intent_http_nodes.json")

			req := &workflow.QueryWorkflowNodeTypeRequest{
				WorkflowID: id,
			}

			response := post[workflow.QueryWorkflowNodeTypeResponse](r, req)
			assert.Contains(t, response.Data.NodeTypes, "1")
			assert.Contains(t, response.Data.NodeTypes, "2")
			assert.Contains(t, response.Data.NodeTypes, "5")
			assert.Contains(t, response.Data.NodeTypes, "22")
			assert.Contains(t, response.Data.NodeTypes, "45")

			for _, prop := range response.Data.NodesProperties {
				if prop.ID == "100001" {
					assert.False(t, prop.IsEnableChatHistory)
					assert.False(t, prop.IsEnableUserQuery)
					assert.False(t, prop.IsRefGlobalVariable)
				}
				if prop.ID == "900001" || prop.ID == "117367" || prop.ID == "133234" || prop.ID == "163493" {
					assert.False(t, prop.IsEnableChatHistory)
					assert.False(t, prop.IsEnableUserQuery)
					assert.False(t, prop.IsRefGlobalVariable)
				}

			}
		})

		t.Run("loop conditions", func(t *testing.T) {
			id := r.load("query_types/loop_condition.json")

			req := &workflow.QueryWorkflowNodeTypeRequest{
				WorkflowID: id,
			}

			response := post[workflow.QueryWorkflowNodeTypeResponse](r, req)
			assert.Contains(t, response.Data.NodeTypes, "1")
			assert.Contains(t, response.Data.NodeTypes, "2")
			assert.Contains(t, response.Data.NodeTypes, "21")
			assert.Contains(t, response.Data.NodeTypes, "5")
			assert.Contains(t, response.Data.NodeTypes, "8")

			for _, prop := range response.Data.NodesProperties {
				if prop.ID == "100001" || prop.ID == "900001" || prop.ID == "114884" || prop.ID == "143932" {
					assert.False(t, prop.IsEnableChatHistory)
					assert.False(t, prop.IsEnableUserQuery)
					assert.False(t, prop.IsRefGlobalVariable)
				}
				if prop.ID == "119585" || prop.ID == "170824" {
					assert.False(t, prop.IsEnableChatHistory)
					assert.False(t, prop.IsEnableUserQuery)
					assert.False(t, prop.IsRefGlobalVariable)
				}

			}
		})

		t.Run("has sub workflow", func(t *testing.T) {
			_ = r.load("query_types/wf2.json", withID(7498668117704163337), withPublish("v0.0.1"))
			_ = r.load("query_types/wf2child.json", withID(7498674832255615002), withPublish("v0.0.1"))
			id := r.load("query_types/subworkflows.json")

			req := &workflow.QueryWorkflowNodeTypeRequest{
				WorkflowID: id,
			}

			response := post[workflow.QueryWorkflowNodeTypeResponse](r, req)

			assert.Contains(t, response.Data.NodeTypes, "1")
			assert.Contains(t, response.Data.NodeTypes, "2")
			assert.Contains(t, response.Data.NodeTypes, "9")

			assert.Contains(t, response.Data.SubWorkflowNodeTypes, "5")
			assert.Contains(t, response.Data.SubWorkflowNodeTypes, "1")
			assert.Contains(t, response.Data.SubWorkflowNodeTypes, "2")

			for _, prop := range response.Data.NodesProperties {
				if prop.ID == "143310" {
					assert.False(t, prop.IsRefGlobalVariable)
				}
			}

			for _, prop := range response.Data.SubWorkflowNodesProperties {
				if prop.ID == "116972" {
					assert.False(t, prop.IsRefGlobalVariable)
				}
				if prop.ID == "124342" {
					assert.False(t, prop.IsRefGlobalVariable)
				}
			}
		})
	})
}

func TestResumeWithQANode(t *testing.T) {
	mockey.PatchConvey("test_resume with qa node", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		chatModel := &testutil.UTChatModel{
			InvokeResultProvider: func(index int, in []*schema.Message) (*schema.Message, error) {
				if index == 0 {
					return &schema.Message{
						Role:    schema.Assistant,
						Content: `{"question": "what's your age?"}`,
						ResponseMeta: &schema.ResponseMeta{
							Usage: &schema.TokenUsage{
								PromptTokens:     6,
								CompletionTokens: 7,
								TotalTokens:      13,
							},
						},
					}, nil
				} else if index == 1 {
					return &schema.Message{
						Role:    schema.Assistant,
						Content: `{"fields": {"name": "eino", "age": 1}}`,
						ResponseMeta: &schema.ResponseMeta{
							Usage: &schema.TokenUsage{
								PromptTokens:     11,
								CompletionTokens: 19,
								TotalTokens:      30,
							},
						},
					}, nil
				}
				return nil, errors.New("not found")
			},
		}

		mockey.Mock(modelbuilder.BuildModelByID).Return(chatModel, nil, nil).Build()

		id := r.load("qa_with_structured_output.json")

		exeID := r.testRun(id, map[string]string{
			"input": "what's your name and age?",
		})

		e := r.getProcess(id, exeID)
		assert.NotNil(t, e.event)
		e.tokenEqual(0, 0)

		r.testResume(id, exeID, e.event.ID, "my name is eino")

		e2 := r.getProcess(id, exeID, withPreviousEventID(e.event.ID))

		r.testResume(id, exeID, e2.event.ID, "1 year old")

		e3 := r.getProcess(id, exeID, withPreviousEventID(e2.event.ID))
		e3.assertSuccess()
		assert.Equal(t, map[string]any{
			"USER_RESPONSE": "1 year old",
			"name":          "eino",
			"age":           int64(1),
		}, mustUnmarshalToMap(t, e3.output))
		e3.tokenEqual(17, 26)
	})
}

func TestPublishWorkflow(t *testing.T) {
	mockey.PatchConvey("publish work flow", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		id := r.load("publish/publish_workflow.json", withName("pb_we"))

		listResponse := post[workflow.GetWorkFlowListResponse](r, &workflow.GetWorkFlowListRequest{
			Page:   ptr.Of(int32(1)),
			Size:   ptr.Of(int32(10)),
			Type:   ptr.Of(workflow.WorkFlowType_User),
			Status: ptr.Of(workflow.WorkFlowListStatus_UnPublished),
			Name:   ptr.Of("pb_we"),
		})

		assert.Equal(t, 1, len(listResponse.Data.WorkflowList))

		r.publish(id, "v0.0.1", true)

		listResponse = post[workflow.GetWorkFlowListResponse](r, &workflow.GetWorkFlowListRequest{
			Page:   ptr.Of(int32(1)),
			Size:   ptr.Of(int32(10)),
			Type:   ptr.Of(workflow.WorkFlowType_User),
			Status: ptr.Of(workflow.WorkFlowListStatus_HadPublished),
			Name:   ptr.Of("pb_we"),
		})

		assert.Equal(t, 1, len(listResponse.Data.WorkflowList))

		r.publish(id, "v0.0.2", true)

		deleteReq := &workflow.DeleteWorkflowRequest{
			WorkflowID: id,
		}
		_ = post[workflow.DeleteWorkflowResponse](r, deleteReq)
		time.Sleep(time.Second)
	})
}

func TestGetCanvasInfo(t *testing.T) {
	mockey.PatchConvey("test get canvas info", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		id := r.load("get_canvas/get_canvas.json")

		getCanvas := &workflow.GetCanvasInfoRequest{
			WorkflowID: ptr.Of(id),
		}
		response := post[workflow.GetCanvasInfoResponse](r, getCanvas)
		assert.Equal(t, response.Data.Workflow.Status, workflow.WorkFlowDevStatus_CanNotSubmit)
		assert.Equal(t, response.Data.VcsData.Type, workflow.VCSCanvasType_Draft)

		exeID := r.testRun(id, map[string]string{
			"input": "input_v1",
			"e":     "e",
		})
		r.getProcess(id, exeID)

		response = post[workflow.GetCanvasInfoResponse](r, getCanvas)
		assert.Equal(t, response.Data.Workflow.Status, workflow.WorkFlowDevStatus_CanSubmit)
		assert.Equal(t, response.Data.VcsData.Type, workflow.VCSCanvasType_Draft)

		r.publish(id, "v0.0.1", true)

		response = post[workflow.GetCanvasInfoResponse](r, getCanvas)
		assert.Equal(t, response.Data.Workflow.Status, workflow.WorkFlowDevStatus_HadSubmit)
		assert.Equal(t, response.Data.VcsData.Type, workflow.VCSCanvasType_Publish)

		r.save(id, "get_canvas/get_canvas.json")
		response = post[workflow.GetCanvasInfoResponse](r, getCanvas)
		assert.Equal(t, response.Data.Workflow.Status, workflow.WorkFlowDevStatus_CanSubmit)
		assert.Equal(t, response.Data.VcsData.Type, workflow.VCSCanvasType_Draft)

		r.save(id, "get_canvas/get_canvas_modify.json")
		response = post[workflow.GetCanvasInfoResponse](r, getCanvas)
		assert.Equal(t, response.Data.Workflow.Status, workflow.WorkFlowDevStatus_CanNotSubmit)
		assert.Equal(t, response.Data.VcsData.Type, workflow.VCSCanvasType_Draft)
	})
}

func TestUpdateWorkflowMeta(t *testing.T) {
	mockey.PatchConvey("update workflow meta", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		id := r.load("entry_exit.json")

		updateMetaReq := &workflow.UpdateWorkflowMetaRequest{
			WorkflowID: id,
			Name:       ptr.Of("modify_name"),
			Desc:       ptr.Of("modify_desc"),
			IconURI:    ptr.Of("modify_icon_uri"),
		}
		_ = post[workflow.UpdateWorkflowMetaResponse](r, updateMetaReq)

		getCanvas := &workflow.GetCanvasInfoRequest{
			WorkflowID: ptr.Of(id),
		}
		response := post[workflow.GetCanvasInfoResponse](r, getCanvas)
		assert.Equal(t, response.Data.Workflow.Name, "modify_name")
		assert.Equal(t, response.Data.Workflow.Desc, "modify_desc")
		assert.Equal(t, response.Data.Workflow.IconURI, "modify_icon_uri")
	})
}

func TestListWorkflowAsToolData(t *testing.T) {
	mockey.PatchConvey("publish list workflow & list workflow as tool data", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()
		name := "pb_wf" + strconv.FormatInt(time.Now().UnixMilli(), 10)
		id := r.load("publish/publish_workflow.json", withName(name))

		listResponse := post[workflow.GetWorkFlowListResponse](r, &workflow.GetWorkFlowListRequest{
			Page:   ptr.Of(int32(1)),
			Size:   ptr.Of(int32(10)),
			Type:   ptr.Of(workflow.WorkFlowType_User),
			Status: ptr.Of(workflow.WorkFlowListStatus_UnPublished),
			Name:   ptr.Of(name),
		})

		assert.Equal(t, 1, len(listResponse.Data.WorkflowList))

		r.publish(id, "v0.0.1", true)

		res, err := appworkflow.SVC.GetPlaygroundPluginList(t.Context(), &pluginAPI.GetPlaygroundPluginListRequest{
			PluginIds: []string{id},
			SpaceID:   ptr.Of(int64(123)),
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(res.Data.PluginList))
		assert.Equal(t, "v0.0.1", res.Data.PluginList[0].VersionName)
		assert.Equal(t, "input", res.Data.PluginList[0].PluginApis[0].Parameters[0].Name)
		assert.Equal(t, "obj", res.Data.PluginList[0].PluginApis[0].Parameters[1].Name)
		assert.Equal(t, "field1", res.Data.PluginList[0].PluginApis[0].Parameters[1].SubParameters[0].Name)
		assert.Equal(t, "arr", res.Data.PluginList[0].PluginApis[0].Parameters[2].Name)
		assert.Equal(t, "string", res.Data.PluginList[0].PluginApis[0].Parameters[2].SubType)

		deleteReq := &workflow.DeleteWorkflowRequest{
			WorkflowID: id,
		}
		_ = post[workflow.DeleteWorkflowResponse](r, deleteReq)
		time.Sleep(time.Second)
	})
}

func TestWorkflowDetailAndDetailInfo(t *testing.T) {
	mockey.PatchConvey("workflow detail & detail info", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		name := "pb_wf" + strconv.FormatInt(time.Now().UnixMilli(), 10)
		id := r.load("publish/publish_workflow.json", withName(name))

		detailReq := &workflow.GetWorkflowDetailRequest{
			WorkflowIds: []string{id},
		}
		response := post[map[string]any](r, detailReq)
		assert.Equal(t, 1, len((*response)["data"].([]any)))

		r.publish(id, "v0.0.1", true)
		r.publish(id, "v0.0.2", true)

		detailInfoReq := &workflow.GetWorkflowDetailInfoRequest{
			WorkflowFilterList: []*workflow.WorkflowFilter{
				{WorkflowID: id},
			},
		}
		detailInfoResponse := post[map[string]any](r, detailInfoReq)
		assert.Equal(t, 1, len((*detailInfoResponse)["data"].([]any)))
		assert.Equal(t, "v0.0.2", (*detailInfoResponse)["data"].([]any)[0].(map[string]any)["latest_flow_version"].(string))
		assert.Equal(t, int64(1), (*detailInfoResponse)["data"].([]any)[0].(map[string]any)["end_type"].(int64))

		deleteReq := &workflow.DeleteWorkflowRequest{
			WorkflowID: id,
		}
		_ = post[workflow.DeleteWorkflowResponse](r, deleteReq)
		time.Sleep(time.Second)
	})
}

func TestInputComplex(t *testing.T) {
	mockey.PatchConvey("test input complex", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		id := r.load("input_complex.json")
		exeID := r.testRun(id, map[string]string{})
		e := r.getProcess(id, exeID)
		r.testResume(id, exeID, e.event.ID, mustMarshalToString(t, map[string]any{
			"input":      `{"name": "eino", "age": 1}`,
			"input_list": `[{"name":"user_1"},{"age":2}]`,
		}))
		e2 := r.getProcess(id, exeID, withPreviousEventID(e.event.ID))
		e2.assertSuccess()
		assert.Equal(t, map[string]any{
			"output": map[string]any{
				"name": "eino",
				"age":  int64(1),
			},
			"output_list": []any{
				map[string]any{
					"name": "user_1",
					"age":  nil,
				},
				map[string]any{
					"name": nil,
					"age":  int64(2),
				},
			},
		}, mustUnmarshalToMap(t, e2.output))
	})
}

func TestGetLLMNodeFCSettingsDetailAndMerged(t *testing.T) {
	mockey.PatchConvey("fc setting detail", t, func() {
		operationString := `{
  "summary" : "根据输入的解梦标题给出相关对应的解梦内容，如果返回的内容为空，给用户返回固定的话术：如果想了解自己梦境的详细解析，需要给我详细的梦见信息，例如： 梦见XXX",
  "operationId" : "xz_zgjm",
  "parameters" : [ {
    "description" : "查询解梦标题，例如：梦见蛇",
    "in" : "query",
    "name" : "title",
    "required" : true,
    "schema" : {
      "description" : "查询解梦标题，例如：梦见蛇",
      "type" : "string"
    }
  } ],
  "requestBody" : {
    "content" : {
      "application/json" : {
        "schema" : {
          "type" : "object"
        }
      }
    }
  },
  "responses" : {
    "200" : {
      "content" : {
        "application/json" : {
          "schema" : {
            "properties" : {
              "data" : {
                "description" : "返回数据",
                "type" : "string"
              },
              "data_structural" : {
                "description" : "返回数据结构",
                "properties" : {
                  "content" : {
                    "description" : "解梦内容",
                    "type" : "string"
                  },
                  "title" : {
                    "description" : "解梦标题",
                    "type" : "string"
                  },
                  "weburl" : {
                    "description" : "当前内容关联的页面地址",
                    "type" : "string"
                  }
                },
                "type" : "object"
              },
              "err_msg" : {
                "description" : "错误提示",
                "type" : "string"
              }
            },
            "required" : [ "data", "data_structural" ],
            "type" : "object"
          }
        }
      },
      "description" : "new desc"
    },
    "default" : {
      "description" : ""
    }
  }
}`
		operation := &pluginmodel.Openapi3Operation{}
		_ = sonic.UnmarshalString(operationString, operation)

		r := newWfTestRunner(t)
		defer r.closeFn()

		r.plugin.EXPECT().MGetOnlinePlugins(gomock.Any(), gomock.Any()).Return([]*entity3.PluginInfo{
			{
				PluginInfo: &pluginmodel.PluginInfo{
					ID:       123,
					SpaceID:  123,
					Version:  ptr.Of("v0.0.1"),
					Manifest: &pluginmodel.PluginManifest{NameForHuman: "p1", DescriptionForHuman: "desc"},
				},
			},
		}, nil).AnyTimes()
		r.plugin.EXPECT().MGetOnlineTools(gomock.Any(), gomock.Any()).Return([]*entity3.ToolInfo{
			{ID: 123, Operation: operation},
		}, nil).AnyTimes()

		pluginSrv := pluginImpl.InitDomainService(r.plugin, r.tos)
		crossplugin.SetDefaultSVC(pluginSrv)

		t.Run("plugin tool info ", func(t *testing.T) {
			fcSettingDetailReq := &workflow.GetLLMNodeFCSettingDetailRequest{
				PluginList: []*workflow.PluginFCItem{
					{PluginID: "123", APIID: "123"},
				},
			}
			response := post[map[string]any](r, fcSettingDetailReq)
			assert.Equal(t, (*response)["plugin_detail_map"].(map[string]any)["123"].(map[string]any)["description"], "desc")
			assert.Equal(t, (*response)["plugin_detail_map"].(map[string]any)["123"].(map[string]any)["name"], "p1")
			assert.Equal(t, (*response)["plugin_api_detail_map"].(map[string]any)["123"].(map[string]any)["name"], "xz_zgjm")
			assert.Equal(t, 1, len((*response)["plugin_api_detail_map"].(map[string]any)["123"].(map[string]any)["parameters"].([]any)))
		})

		t.Run("workflow tool info ", func(t *testing.T) {
			r.load("entry_exit.json", withID(123), withPublish("v0.0.1"))
			fcSettingDetailReq := &workflow.GetLLMNodeFCSettingDetailRequest{
				WorkflowList: []*workflow.WorkflowFCItem{
					{WorkflowID: "123", PluginID: "123", WorkflowVersion: ptr.Of("v0.0.1")},
				},
			}
			response := post[map[string]any](r, fcSettingDetailReq)
			assert.Equal(t, (*response)["workflow_detail_map"].(map[string]any)["123"].(map[string]any)["plugin_id"], "123")
			assert.Equal(t, (*response)["workflow_detail_map"].(map[string]any)["123"].(map[string]any)["name"], "test_wf")
			assert.Equal(t, (*response)["workflow_detail_map"].(map[string]any)["123"].(map[string]any)["description"], "this is a test wf")
		})
	})
	mockey.PatchConvey("fc setting merged", t, func() {
		operationString := `{
  "summary" : "根据输入的解梦标题给出相关对应的解梦内容，如果返回的内容为空，给用户返回固定的话术：如果想了解自己梦境的详细解析，需要给我详细的梦见信息，例如： 梦见XXX",
  "operationId" : "xz_zgjm",
  "parameters" : [ {
    "description" : "查询解梦标题，例如：梦见蛇",
    "in" : "query",
    "name" : "title",
    "required" : true,
    "schema" : {
      "description" : "查询解梦标题，例如：梦见蛇",
      "type" : "string"
    }
  } ],
  "requestBody" : {
    "content" : {
      "application/json" : {
        "schema" : {
          "type" : "object"
        }
      }
    }
  },
  "responses" : {
    "200" : {
      "content" : {
        "application/json" : {
          "schema" : {
            "properties" : {
              "data" : {
                "description" : "返回数据",
                "type" : "string"
              },
              "data_structural" : {
                "description" : "返回数据结构",
                "properties" : {
                  "content" : {
                    "description" : "解梦内容",
                    "type" : "string"
                  },
                  "title" : {
                    "description" : "解梦标题",
                    "type" : "string"
                  },
                  "weburl" : {
                    "description" : "当前内容关联的页面地址",
                    "type" : "string"
                  }
                },
                "type" : "object"
              },
              "err_msg" : {
                "description" : "错误提示",
                "type" : "string"
              }
            },
            "required" : [ "data", "data_structural" ],
            "type" : "object"
          }
        }
      },
      "description" : "new desc"
    },
    "default" : {
      "description" : ""
    }
  }
}`

		operation := &pluginmodel.Openapi3Operation{}
		_ = sonic.UnmarshalString(operationString, operation)
		r := newWfTestRunner(t)
		defer r.closeFn()

		r.plugin.EXPECT().MGetOnlinePlugins(gomock.Any(), gomock.Any()).Return([]*entity3.PluginInfo{
			{
				PluginInfo: &pluginmodel.PluginInfo{
					ID:       123,
					SpaceID:  123,
					Version:  ptr.Of("v0.0.1"),
					Manifest: &pluginmodel.PluginManifest{NameForHuman: "p1", DescriptionForHuman: "desc"},
				},
			},
		}, nil).AnyTimes()
		r.plugin.EXPECT().MGetOnlineTools(gomock.Any(), gomock.Any()).Return([]*entity3.ToolInfo{
			{ID: 123, Operation: operation},
		}, nil).AnyTimes()

		pluginSrv := pluginImpl.InitDomainService(r.plugin, r.tos)
		crossplugin.SetDefaultSVC(pluginSrv)

		t.Run("plugin merge", func(t *testing.T) {
			fcSettingMergedReq := &workflow.GetLLMNodeFCSettingsMergedRequest{
				PluginFcSetting: &workflow.FCPluginSetting{
					PluginID: "123", APIID: "123",
					RequestParams: []*workflow.APIParameter{
						{Name: "title", LocalDisable: true, LocalDefault: ptr.Of("value")},
					},
					ResponseParams: []*workflow.APIParameter{
						{Name: "data123", LocalDisable: true},
					},
				},
			}
			response := post[map[string]any](r, fcSettingMergedReq)
			assert.Equal(t, (*response)["plugin_fc_setting"].(map[string]any)["request_params"].([]any)[0].(map[string]any)["local_disable"], true)
			names := map[string]bool{
				"data":            true,
				"data_structural": true,
				"err_msg":         true,
			}
			assert.Equal(t, 3, len((*response)["plugin_fc_setting"].(map[string]any)["response_params"].([]any)))

			for _, mm := range (*response)["plugin_fc_setting"].(map[string]any)["response_params"].([]any) {
				n := mm.(map[string]any)["name"].(string)
				assert.True(t, names[n])
			}
		})
		t.Run("workflow merge", func(t *testing.T) {
			r.load("entry_exit.json", withID(1234), withPublish("v0.0.1"))
			fcSettingMergedReq := &workflow.GetLLMNodeFCSettingsMergedRequest{
				WorkflowFcSetting: &workflow.FCWorkflowSetting{
					WorkflowID: "1234",
					PluginID:   "1234",
					RequestParams: []*workflow.APIParameter{
						{Name: "obj", LocalDisable: true, LocalDefault: ptr.Of("{}")},
					},
					ResponseParams: []*workflow.APIParameter{
						{Name: "literal_key", LocalDisable: true},
						{Name: "literal_key_bak", LocalDisable: true},
					},
				},
			}

			response := post[map[string]any](r, fcSettingMergedReq)
			assert.Equal(t, 3, len((*response)["worflow_fc_setting"].(map[string]any)["request_params"].([]any)))
			assert.Equal(t, 8, len((*response)["worflow_fc_setting"].(map[string]any)["response_params"].([]any)))

			for _, mm := range (*response)["worflow_fc_setting"].(map[string]any)["request_params"].([]any) {
				if mm.(map[string]any)["name"].(string) == "obj" {
					assert.True(t, mm.(map[string]any)["local_disable"].(bool))
				}
			}
		})
	})
}

func TestNodeDebugLoop(t *testing.T) {
	mockey.PatchConvey("test node debug loop", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		id := r.load("loop_selector_variable_assign_text_processor.json")
		exeID := r.nodeDebug(id, "192046", withNDInput(map[string]string{"input": `["a", "bb", "ccc", "dddd"]`}))
		e := r.getProcess(id, exeID, withSpecificNodeID("192046"))
		e.assertSuccess()
		assert.Equal(t, map[string]any{
			"converted": []any{
				"new_a",
				"new_ccc",
			},
			"variable_out": "dddd",
		}, mustUnmarshalToMap(t, e.output))

		result := r.getNodeExeHistory(id, exeID, "192046", nil)
		assert.Equal(t, mustUnmarshalToMap(t, e.output), mustUnmarshalToMap(t, result.Output))

		// verify this workflow has not been successfully test ran
		result = r.getNodeExeHistory(id, "", "100001", ptr.Of(workflow.NodeHistoryScene_TestRunInput))
		assert.Equal(t, "", result.Output)

		// verify that another node of this workflow is not node debugged
		result = r.getNodeExeHistory(id, "", "wrong_node_id", ptr.Of(workflow.NodeHistoryScene_TestRunInput))
		assert.Equal(t, "", result.Output)
	})

	mockey.PatchConvey("test node debug loop", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		id := r.load("loop_selector_variable_assign_text_processor.json")
		exeID := r.nodeDebug(id, "192046", withNDInput(map[string]string{"input": `["a", "bb", "ccc", "dddd"]`}))
		e := r.getProcess(id, exeID, withSpecificNodeID("192046"))
		e.assertSuccess()
		assert.Equal(t, map[string]any{
			"converted": []any{
				"new_a",
				"new_ccc",
			},
			"variable_out": "dddd",
		}, mustUnmarshalToMap(t, e.output))

		result := r.getNodeExeHistory(id, exeID, "192046", nil)
		assert.Equal(t, mustUnmarshalToMap(t, e.output), mustUnmarshalToMap(t, result.Output))

		// verify this workflow has not been successfully test ran
		result = r.getNodeExeHistory(id, "", "100001", ptr.Of(workflow.NodeHistoryScene_TestRunInput))
		assert.Equal(t, "", result.Output)

		// verify that another node of this workflow is not node debugged
		result = r.getNodeExeHistory(id, "", "wrong_node_id", ptr.Of(workflow.NodeHistoryScene_TestRunInput))
		assert.Equal(t, "", result.Output)
	})

	mockey.PatchConvey("test node debug loop", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()
		runner := mockcode.NewMockRunner(r.ctrl)
		runner.EXPECT().Run(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, request *coderunner.RunRequest) (*coderunner.RunResponse, error) {
			return &coderunner.RunResponse{
				Result: request.Params,
			}, nil
		}).AnyTimes()

		coderunner.SetCodeRunner(runner)
		id := r.load("loop_with_object_input.json")
		exeID := r.nodeDebug(id, "122149",
			withNDInput(map[string]string{"input": `[{"a":"1"},{"a":"2"}]`}))
		e := r.getProcess(id, exeID, withSpecificNodeID("122149"))
		e.assertSuccess()
		assert.Equal(t, `{"output":["1","2"]}`, e.output)

	})

}

func TestCopyWorkflow(t *testing.T) {
	mockey.PatchConvey("copy work flow", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		id := r.load("publish/publish_workflow.json", withName("original_workflow"))

		response := post[workflow.CopyWorkflowResponse](r, &workflow.CopyWorkflowRequest{
			WorkflowID: id,
		})

		oldCanvasResponse := post[workflow.GetCanvasInfoResponse](r, &workflow.GetCanvasInfoRequest{
			WorkflowID: ptr.Of(id),
		})

		copiedCanvasResponse := post[workflow.GetCanvasInfoResponse](r, &workflow.GetCanvasInfoRequest{
			WorkflowID: ptr.Of(response.Data.WorkflowID),
		})

		assert.Equal(t, ptr.From(oldCanvasResponse.Data.Workflow.SchemaJSON), ptr.From(copiedCanvasResponse.Data.Workflow.SchemaJSON))
		assert.Equal(t, "original_workflow_1", copiedCanvasResponse.Data.Workflow.Name)

		_ = post[workflow.BatchDeleteWorkflowResponse](r, &workflow.BatchDeleteWorkflowRequest{
			WorkflowIDList: []string{id, response.Data.WorkflowID},
		})

		wid, _ := strconv.ParseInt(id, 10, 64)

		_, err := appworkflow.GetWorkflowDomainSVC().Get(context.Background(), &vo.GetPolicy{
			ID:       wid,
			QType:    workflowModel.FromDraft,
			CommitID: "",
		})
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, strconv.Itoa(errno.ErrWorkflowNotFound))
	})
}

func TestJsonSerializationDeserialization(t *testing.T) {
	mockey.PatchConvey("test JSON serialization and deserialization workflow", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		idStr := r.load("json/json_test.json")

		mockey.PatchConvey("no type conversion", func() {
			testInput := map[string]string{
				"person": `{"int":123,"string":"hello","bool":true}`,
			}

			exeID := r.testRun(idStr, testInput)
			e := r.getProcess(idStr, exeID)
			output := e.output
			t.Logf("JSON deserialization result (no conversion): %s", output)

			var result map[string]any
			err := sonic.Unmarshal([]byte(output), &result)
			assert.NoError(t, err, "Failed to unmarshal output JSON")

			outputData, ok := result["output"].(map[string]any)
			assert.True(t, ok, "output field is not a map[string]any")

			assert.Equal(t, int64(123), outputData["int"], "int field mismatch")
			assert.Equal(t, "hello", outputData["string"], "string field mismatch")
			assert.Equal(t, true, outputData["bool"], "bool field mismatch")
		})

		mockey.PatchConvey("legal type conversion", func() {
			testInput := map[string]string{
				"person": `{"int":"123","string":456,"bool":"true"}`,
			}

			exeID := r.testRun(idStr, testInput)
			e := r.getProcess(idStr, exeID)
			output := e.output
			t.Logf("JSON deserialization result (legal conversion): %s", output)

			var result map[string]any
			err := sonic.Unmarshal([]byte(output), &result)
			assert.NoError(t, err, "Failed to unmarshal output JSON")

			outputData, ok := result["output"].(map[string]any)
			assert.True(t, ok, "output field is not a map[string]any")

			assert.Equal(t, int64(123), outputData["int"], "int field mismatch")
			assert.Equal(t, "456", outputData["string"], "string field mismatch")
			assert.Equal(t, true, outputData["bool"], "bool field mismatch")
		})
	})
}

func TestJsonSerializationDeserializationWithWarning(t *testing.T) {
	mockey.PatchConvey("test JSON serialization and deserialization with warning", t, func() {
		r := newWfTestRunner(t)
		defer r.closeFn()

		idStr := r.load("json/json_test_warning.json")
		testInput := map[string]string{
			"person": `{"int":1,"string":"abc","bool":true}`,
		}

		exeID := r.testRun(idStr, testInput)
		e := r.getProcess(idStr, exeID)
		output := e.output
		t.Logf("JSON deserialization result (legal conversion): %s", output)

		var result map[string]any
		err := sonic.Unmarshal([]byte(output), &result)
		assert.NoError(t, err, "Failed to unmarshal output JSON")

		outputData, ok := result["output"].(map[string]any)
		assert.True(t, ok, "output field is not a map[string]any")

		assert.Equal(t, nil, outputData["int"], "int field mismatch")
		assert.Equal(t, "abc", outputData["string"], "string field mismatch")
		assert.Equal(t, true, outputData["bool"], "bool field mismatch")
	})
}
