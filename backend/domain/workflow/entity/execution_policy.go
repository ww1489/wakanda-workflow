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
	"fmt"

	workflowapi "github.com/coze-dev/coze-studio/backend/api/model/workflow"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
)

func ValidateExecuteConfig(config workflowModel.ExecuteConfig) error {
	if config.WorkflowMode != workflowapi.WorkflowMode_Workflow {
		return fmt.Errorf("only standard workflows can be executed")
	}
	if config.AgentID != nil && *config.AgentID != 0 {
		return fmt.Errorf("agent bindings are not supported")
	}
	if config.AppID != nil && *config.AppID != 0 {
		return fmt.Errorf("application workflows are not supported")
	}
	if config.ConversationID != nil && *config.ConversationID != 0 ||
		len(config.ConversationHistorySchemaMessages) > 0 ||
		config.MaxHistoryRounds != nil && *config.MaxHistoryRounds > 0 {
		return fmt.Errorf("conversation history is not supported")
	}
	return nil
}
