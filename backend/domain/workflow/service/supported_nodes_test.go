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
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	workflowapi "github.com/coze-dev/coze-studio/backend/api/model/workflow"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	mockworkflow "github.com/coze-dev/coze-studio/backend/internal/mock/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
)

func TestSaveRejectsUnsupportedNodesBeforeTestRunCheck(t *testing.T) {
	svc := &impl{}
	err := svc.Save(t.Context(), 1, `{"nodes":[{"id":"knowledge","type":"6"}]}`)
	require.ErrorContains(t, err, "not supported")
}

func TestSaveRejectsExistingChatflowWithOrdinaryNodes(t *testing.T) {
	repo := mockworkflow.NewMockRepository(gomock.NewController(t))
	repo.EXPECT().GetMeta(gomock.Any(), int64(1)).Return(&vo.Meta{
		Mode: workflowapi.WorkflowMode_ChatFlow,
	}, nil)
	svc := &impl{repo: repo}
	err := svc.Save(t.Context(), 1, `{"nodes":[{"id":"start","type":"1"}]}`)
	require.ErrorContains(t, err, "only standard workflows can be edited")
}

func TestChatflowMutationRejected(t *testing.T) {
	svc := &impl{}
	mode := workflowapi.WorkflowMode_ChatFlow
	_, err := svc.Create(t.Context(), &vo.MetaCreate{Mode: mode})
	require.ErrorContains(t, err, "only standard workflows")
	require.ErrorContains(t, svc.UpdateMeta(t.Context(), 1, &vo.MetaUpdate{WorkflowMode: &mode}), "only standard workflows")
}

func TestApplicationWorkflowMutationRejected(t *testing.T) {
	svc := &impl{}
	_, err := svc.Create(t.Context(), &vo.MetaCreate{AppID: ptr.Of(int64(1))})
	require.ErrorContains(t, err, "application workflows")

	repo := mockworkflow.NewMockRepository(gomock.NewController(t))
	repo.EXPECT().GetMeta(gomock.Any(), int64(1)).Return(&vo.Meta{
		AppID: ptr.Of(int64(1)),
	}, nil).Times(2)
	svc.repo = repo
	require.ErrorContains(t, svc.Save(t.Context(), 1, `{"nodes":[{"id":"start","type":"1"}]}`), "application workflows")
	require.ErrorContains(t, svc.Publish(t.Context(), &vo.PublishPolicy{ID: 1}), "application workflows")
}

func TestUpdateMetaRejectsUnsupportedExistingWorkflow(t *testing.T) {
	tests := []struct {
		name string
		meta *vo.Meta
	}{
		{"chatflow", &vo.Meta{Mode: workflowapi.WorkflowMode_ChatFlow}},
		{"application", &vo.Meta{AppID: ptr.Of(int64(1))}},
		{"standalone", &vo.Meta{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockworkflow.NewMockRepository(gomock.NewController(t))
			repo.EXPECT().GetMeta(gomock.Any(), int64(1)).Return(tt.meta, nil)
			update := &vo.MetaUpdate{
				Name:         ptr.Of("renamed"),
				WorkflowMode: ptr.Of(workflowapi.WorkflowMode_Workflow),
			}
			if tt.name == "standalone" {
				repo.EXPECT().UpdateMeta(gomock.Any(), int64(1), update).Return(nil)
			}
			svc := &impl{repo: repo}
			err := svc.UpdateMeta(t.Context(), 1, update)
			if tt.name == "standalone" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestStoredApplicationWorkflowCannotExecuteWithoutRequestAppID(t *testing.T) {
	for _, operation := range []string{"sync", "async", "node", "stream", "resume", "stream-resume"} {
		t.Run(operation, func(t *testing.T) {
			repo := mockworkflow.NewMockRepository(gomock.NewController(t))
			repo.EXPECT().GetEntity(gomock.Any(), gomock.Any()).Return(&entity.Workflow{
				Meta:       &vo.Meta{AppID: ptr.Of(int64(1))},
				CanvasInfo: &vo.CanvasInfo{Canvas: `{"nodes":[]}`},
			}, nil)
			svc := &impl{repo: repo}
			var err error
			switch operation {
			case "sync":
				_, _, err = svc.SyncExecute(t.Context(), workflowModel.ExecuteConfig{}, nil)
			case "async":
				_, err = svc.AsyncExecute(t.Context(), workflowModel.ExecuteConfig{}, nil)
			case "node":
				_, err = svc.AsyncExecuteNode(t.Context(), "start", workflowModel.ExecuteConfig{}, nil)
			case "stream":
				_, err = svc.StreamExecute(t.Context(), workflowModel.ExecuteConfig{}, nil)
			default:
				repo.EXPECT().GetWorkflowExecution(gomock.Any(), int64(1)).Return(&entity.WorkflowExecution{
					ID: 1, RootExecutionID: 1, Status: entity.WorkflowInterrupted,
				}, true, nil)
				if operation == "resume" {
					err = svc.AsyncResume(t.Context(), &entity.ResumeRequest{ExecuteID: 1}, workflowModel.ExecuteConfig{})
				} else {
					_, err = svc.StreamResume(t.Context(), &entity.ResumeRequest{ExecuteID: 1}, workflowModel.ExecuteConfig{})
				}
			}
			require.ErrorContains(t, err, "application workflows")
		})
	}
}

func TestNodeMenuMatchesSupportedCapabilities(t *testing.T) {
	svc := &impl{}
	menu, _, err := svc.ListNodeMeta(t.Context(), nil)
	require.NoError(t, err)
	seen := map[entity.NodeType]bool{}
	for _, nodes := range menu {
		for _, node := range nodes {
			require.True(t, entity.IsSupportedNodeType(node.Key))
			seen[node.Key] = true
		}
	}
	require.True(t, seen[entity.NodeTypePlugin])
	require.True(t, seen[entity.NodeTypeDatabaseQuery])
	require.False(t, seen[entity.NodeTypeKnowledgeRetriever])
	require.False(t, seen[entity.NodeTypeVariableAssigner])
	require.False(t, seen[entity.NodeTypeCreateConversation])
}
