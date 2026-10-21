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

package entity

import (
	"testing"

	"github.com/stretchr/testify/require"

	workflowapi "github.com/coze-dev/coze-studio/backend/api/model/workflow"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
)

func TestValidateExecuteConfig(t *testing.T) {
	require.NoError(t, ValidateExecuteConfig(workflowModel.ExecuteConfig{}))
	require.ErrorContains(t, ValidateExecuteConfig(workflowModel.ExecuteConfig{
		WorkflowMode: workflowapi.WorkflowMode_ChatFlow,
	}), "standard workflows")
	require.ErrorContains(t, ValidateExecuteConfig(workflowModel.ExecuteConfig{
		AgentID: ptr.Of(int64(1)),
	}), "agent bindings")
	require.ErrorContains(t, ValidateExecuteConfig(workflowModel.ExecuteConfig{
		AppID: ptr.Of(int64(1)),
	}), "application workflows")
	require.ErrorContains(t, ValidateExecuteConfig(workflowModel.ExecuteConfig{
		ConversationID: ptr.Of(int64(1)),
	}), "conversation history")
}
