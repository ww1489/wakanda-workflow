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

package impl

import (
	"context"

	"github.com/cloudwego/eino/compose"
	einoCompose "github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	crossworkflow "github.com/coze-dev/coze-studio/backend/crossdomain/workflow"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	"github.com/coze-dev/coze-studio/backend/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	workflowEntity "github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
)

var defaultSVC crossworkflow.Workflow

type impl struct {
	DomainSVC workflow.Service
}

func InitDomainService(c workflow.Service) crossworkflow.Workflow {
	defaultSVC = &impl{
		DomainSVC: c,
	}

	return defaultSVC
}

func (i *impl) WorkflowAsModelTool(ctx context.Context, policies []*vo.GetPolicy) ([]workflow.ToolFromWorkflow, error) {
	return i.DomainSVC.WorkflowAsModelTool(ctx, policies)
}

func (i *impl) WithResumeToolWorkflow(resumingEvent *workflowEntity.ToolInterruptEvent, resumeData string, allInterruptEvents map[string]*workflowEntity.ToolInterruptEvent) einoCompose.Option {
	return i.DomainSVC.WithResumeToolWorkflow(resumingEvent, resumeData, allInterruptEvents)
}
func (i *impl) SyncExecuteWorkflow(ctx context.Context, config workflowModel.ExecuteConfig, input map[string]any) (*workflowEntity.WorkflowExecution, vo.TerminatePlan, error) {
	return i.DomainSVC.SyncExecute(ctx, config, input)
}

func (i *impl) WithExecuteConfig(cfg workflowModel.ExecuteConfig) einoCompose.Option {
	return i.DomainSVC.WithExecuteConfig(cfg)
}

func (i *impl) StreamResume(ctx context.Context, req *entity.ResumeRequest, config workflowModel.ExecuteConfig) (*schema.StreamReader[*entity.Message], error) {
	return i.DomainSVC.StreamResume(ctx, req, config)
}
func (i *impl) StreamExecute(ctx context.Context, config workflowModel.ExecuteConfig, input map[string]any) (*schema.StreamReader[*workflowEntity.Message], error) {
	return i.DomainSVC.StreamExecute(ctx, config, input)
}

func (i *impl) WithMessagePipe() (compose.Option, *schema.StreamReader[*entity.Message], func()) {
	return i.DomainSVC.WithMessagePipe()
}

func (i *impl) MGet(ctx context.Context, policy *vo.MGetPolicy) ([]*entity.Workflow, int64, error) {
	return i.DomainSVC.MGet(ctx, policy)
}
