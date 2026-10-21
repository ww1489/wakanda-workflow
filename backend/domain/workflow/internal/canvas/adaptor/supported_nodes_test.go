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

package adaptor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-studio/backend/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	mockworkflow "github.com/coze-dev/coze-studio/backend/internal/mock/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

func TestValidateSupportedCanvas(t *testing.T) {
	tests := []struct {
		name   string
		node   string
		errMsg string
	}{
		{"incomplete llm draft", `{"id":"llm","type":"3"}`, ""},
		{"python", `{"id":"code","type":"5","data":{"inputs":{"language":3}}}`, ""},
		{"javascript", `{"id":"code","type":"5","data":{"inputs":{"language":5}}}`, "only Python"},
		{"knowledge retrieval", `{"id":"kb","type":"6"}`, "not supported"},
		{"knowledge write", `{"id":"kb","type":"27"}`, "not supported"},
		{"knowledge delete", `{"id":"kb","type":"60"}`, "not supported"},
		{"conversation", `{"id":"conversation","type":"39"}`, "not supported"},
		{"persistent variable assignment", `{"id":"variable","type":"40"}`, "not supported"},
		{"nested knowledge", `{"id":"loop","type":"21","blocks":[{"id":"kb","type":"6"}]}`, `node "kb"`},
		{"unknown type", `{"id":"unknown","type":"9999"}`, "not supported"},
		{"llm knowledge binding", `{"id":"llm","type":"3","data":{"inputs":{"fcParam":{"knowledgeFCParam":{"knowledgeList":[{"id":"1"}]}}}}}`, "knowledge bindings"},
		{"empty llm knowledge binding", `{"id":"llm","type":"3","data":{"inputs":{"fcParam":{"knowledgeFCParam":{"knowledgeList":[]}}}}}`, ""},
		{"llm plugins and workflows", `{"id":"llm","type":"3","data":{"inputs":{"fcParam":{"pluginFCParam":{"pluginList":[{"plugin_id":"1"}]},"workflowFCParam":{"workflowList":[{"workflow_id":"2"}]}}}}}`, ""},
		{"history setting", `{"id":"llm","type":"3","data":{"inputs":{"chatHistorySetting":{"enableChatHistory":true}}}}`, "conversation history"},
		{"llm history parameter", `{"id":"llm","type":"3","data":{"inputs":{"llmParam":[{"name":"enableChatHistory","input":{"value":{"type":"literal","content":true}}}]}}}`, "conversation history"},
		{"global reference", `{"id":"end","type":"2","data":{"inputs":{"inputParameters":[{"input":{"value":{"type":"ref","content":{"source":"global_variable_user","path":["name"]}}}}]}}}`, "variable references"},
		{"system reference", `{"id":"end","type":"2","data":{"inputs":{"inputParameters":[{"input":{"value":{"type":"ref","content":{"source":"global_variable_system","path":["name"]}}}}]}}}`, "variable references"},
		{"application object reference", `{"id":"end","type":"2","data":{"inputs":{"inputParameters":[{"input":{"value":{"type":"object_ref","content":{"source":"global_variable_app","path":["name"]}}}}]}}}`, "variable references"},
		{"literal data resembling a reference", `{"id":"end","type":"2","data":{"inputs":{"inputParameters":[{"input":{"value":{"type":"literal","content":{"type":"ref","content":{"source":"global_variable_user","path":["name"]}}}}}]}}}`, ""},
		{"node output reference", `{"id":"end","type":"2","data":{"inputs":{"inputParameters":[{"input":{"value":{"type":"ref","content":{"source":"block-output","blockID":"100001","name":"value"}}}}]}}}`, ""},
		{"loop local assignment", `{"id":"assign","type":"20"}`, ""},
		{"loop global output reference", `{"id":"loop","type":"21","data":{"outputs":[{"name":"result","input":{"value":{"type":"ref","content":{"source":"global_variable_user","path":["name"]}}}}]}}`, "variable references"},
		{"loop local output reference", `{"id":"loop","type":"21","data":{"outputs":[{"name":"result","input":{"value":{"type":"ref","content":{"source":"block-output","blockID":"loop","name":"local"}}}}]}}`, ""},
		{"http global reference", `{"id":"http","type":"45","data":{"inputs":{"apiInfo":{"url":"https://example.test/{{global_variable_app[\"name\"]}}"}}}}`, "variable references"},
		{"http node reference", `{"id":"http","type":"45","data":{"inputs":{"apiInfo":{"url":"https://example.test/{{block_output_100001.name}}"}}}}`, ""},
		{"database query", `{"id":"db","type":"43"}`, ""},
		{"plugin", `{"id":"plugin","type":"4"}`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSupportedCanvas(`{"nodes":[` + tt.node + `],"edges":[]}`)
			if tt.errMsg == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.errMsg)
			}
		})
	}
}

func TestUnsupportedNodesRejectedBeforeCompilation(t *testing.T) {
	node := &vo.Node{ID: "isolated-knowledge", Type: "6"}
	canvas := &vo.Canvas{Nodes: []*vo.Node{node}}
	_, err := CanvasToWorkflowSchema(t.Context(), canvas)
	require.ErrorContains(t, err, node.ID) // Even isolated nodes must be rejected before pruning.
	_, _, err = NodeToNodeSchema(t.Context(), node, canvas)
	require.ErrorContains(t, err, node.ID)
}

func TestSubWorkflowRejectsUnsupportedNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mockworkflow.NewMockRepository(ctrl)
	previous := workflow.GetRepository()
	workflow.SetRepository(repo)
	t.Cleanup(func() { workflow.SetRepository(previous) })
	repo.EXPECT().GetEntity(gomock.Any(), gomock.Any()).Return(&entity.Workflow{
		Meta:       &vo.Meta{},
		CanvasInfo: &vo.CanvasInfo{Canvas: `{"nodes":[{"id":"nested-knowledge","type":"6"}]}`},
	}, nil)

	var node vo.Node
	require.NoError(t, sonic.UnmarshalString(`{"id":"sub","type":"9","data":{"inputs":{"workflowId":"123","workflowVersion":"v1"}}}`, &node))
	_, _, err := NodeToNodeSchema(t.Context(), &node, &vo.Canvas{})
	require.ErrorContains(t, err, "nested-knowledge")

	repo.EXPECT().GetEntity(gomock.Any(), gomock.Any()).Return(&entity.Workflow{
		Meta: &vo.Meta{AppID: ptr.Of(int64(1))},
	}, nil)
	_, _, err = NodeToNodeSchema(t.Context(), &node, &vo.Canvas{})
	require.ErrorContains(t, err, "application workflows")
}
