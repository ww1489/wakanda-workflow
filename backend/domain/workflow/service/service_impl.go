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
	"errors"
	"fmt"
	"strconv"

	einoCompose "github.com/cloudwego/eino/compose"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"

	"gorm.io/gorm"

	cloudworkflow "github.com/coze-dev/coze-studio/backend/api/model/workflow"

	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	"github.com/coze-dev/coze-studio/backend/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/canvas/adaptor"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/canvas/convert"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/repo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/schema"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
	"github.com/coze-dev/coze-studio/backend/infra/storage"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"

	"github.com/coze-dev/coze-studio/backend/pkg/lang/slices"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

type impl struct {
	repo workflow.Repository
	*asToolImpl
	*executableImpl
}

func NewWorkflowService(repo workflow.Repository) workflow.Service {
	return &impl{
		repo: repo,
		asToolImpl: &asToolImpl{
			repo: repo,
		},
		executableImpl: &executableImpl{
			repo: repo,
		},
	}
}

func NewWorkflowRepository(idgen idgen.IDGenerator, db *gorm.DB, redis cache.Cmdable, tos storage.Storage,
	cpStore einoCompose.CheckPointStore, cfg workflow.WorkflowConfig) (workflow.Repository, error) {
	return repo.NewRepository(idgen, db, redis, tos, cpStore, cfg)
}

func (i *impl) ListNodeMeta(_ context.Context, nodeTypes map[entity.NodeType]bool) (map[string][]*entity.NodeTypeMeta, []entity.Category, error) {
	// Initialize result maps
	nodeMetaMap := make(map[string][]*entity.NodeTypeMeta)

	// Helper function to check if a type should be included based on the filter
	shouldInclude := func(meta *entity.NodeTypeMeta) bool {
		if meta.Disabled || !entity.IsSupportedNodeType(meta.Key) {
			return false
		}
		nodeType := meta.Key
		if nodeTypes == nil || len(nodeTypes) == 0 {
			return true // No filter, include all
		}
		_, ok := nodeTypes[nodeType]
		return ok
	}

	// Process standard node types
	for _, meta := range entity.NodeTypeMetas {
		if shouldInclude(meta) {
			nodeMetaMap[meta.Category] = append(nodeMetaMap[meta.Category], meta)
		}
	}

	return nodeMetaMap, entity.Categories, nil
}

func (i *impl) Create(ctx context.Context, meta *vo.MetaCreate) (int64, error) {
	if meta.Mode != cloudworkflow.WorkflowMode_Workflow {
		return 0, fmt.Errorf("only standard workflows are supported")
	}
	if meta.AppID != nil && *meta.AppID != 0 {
		return 0, fmt.Errorf("application workflows are not supported")
	}
	if err := adaptor.ValidateSupportedCanvas(meta.InitCanvasSchema); err != nil {
		return 0, err
	}
	id, err := i.repo.CreateMeta(ctx, &vo.Meta{
		CreatorID:   meta.CreatorID,
		SpaceID:     meta.SpaceID,
		ContentType: meta.ContentType,
		Name:        meta.Name,
		Desc:        meta.Desc,
		IconURI:     meta.IconURI,
		AppID:       meta.AppID,
		Mode:        meta.Mode,
	})
	if err != nil {
		return 0, err
	}

	// save the initialized  canvas information to the draft
	if err = i.Save(ctx, id, meta.InitCanvasSchema); err != nil {
		return 0, err
	}

	return id, nil
}

func (i *impl) Save(ctx context.Context, id int64, schema string) (err error) {
	var draft vo.Canvas
	if err = sonic.UnmarshalString(schema, &draft); err != nil {
		return vo.WrapError(errno.ErrSerializationDeserializationFail, err)
	}
	if err = adaptor.ValidateSupportedNodes(draft.Nodes); err != nil {
		return err
	}
	meta, err := i.repo.GetMeta(ctx, id)
	if err != nil {
		return err
	}
	if meta.Mode != cloudworkflow.WorkflowMode_Workflow {
		return fmt.Errorf("only standard workflows can be edited")
	}
	if meta.AppID != nil && *meta.AppID != 0 {
		return fmt.Errorf("application workflows are not supported")
	}

	var inputParams, outputParams string
	inputs, outputs := extractInputsAndOutputsNamedInfoList(&draft)
	if inputParams, err = sonic.MarshalString(inputs); err != nil {
		return vo.WrapError(errno.ErrSerializationDeserializationFail, err)
	}

	if outputParams, err = sonic.MarshalString(outputs); err != nil {
		return vo.WrapError(errno.ErrSerializationDeserializationFail, err)
	}

	testRunSuccess, err := i.calculateTestRunSuccess(ctx, &draft, id)
	if err != nil {
		return err
	}

	commitID, err := i.repo.GenID(ctx) // generate a new commit ID for this draft version
	if err != nil {
		return vo.WrapError(errno.ErrIDGenError, err)
	}

	return i.repo.CreateOrUpdateDraft(ctx, id, &vo.DraftInfo{
		Canvas: schema,
		DraftMeta: &vo.DraftMeta{
			TestRunSuccess: testRunSuccess,
			Modified:       true,
		},
		InputParamsStr:  inputParams,
		OutputParamsStr: outputParams,
		CommitID:        strconv.FormatInt(commitID, 10),
	})
}

func extractInputsAndOutputsNamedInfoList(c *vo.Canvas) (inputs []*vo.NamedTypeInfo, outputs []*vo.NamedTypeInfo) {
	defer func() {
		if err := recover(); err != nil {
			logs.Warnf("failed to extract inputs and outputs: %v", err)
		}
	}()
	var (
		startNode *vo.Node
		endNode   *vo.Node
	)
	inputs = make([]*vo.NamedTypeInfo, 0)
	outputs = make([]*vo.NamedTypeInfo, 0)
	for _, node := range c.Nodes {
		if startNode != nil && endNode != nil {
			break
		}
		if node.Type == entity.NodeTypeEntry.IDStr() {
			startNode = node
		}
		if node.Type == entity.NodeTypeExit.IDStr() {
			endNode = node
		}
	}

	var err error
	if startNode != nil {
		inputs, err = slices.TransformWithErrorCheck(startNode.Data.Outputs, func(o any) (*vo.NamedTypeInfo, error) {
			v, err := vo.ParseVariable(o)
			if err != nil {
				return nil, err
			}
			nInfo, err := convert.VariableToNamedTypeInfo(v)
			if err != nil {
				return nil, err
			}
			return nInfo, nil
		})
		if err != nil {
			logs.Warn(fmt.Sprintf("transform start node outputs to named info failed, err=%v", err))
		}
	}

	if endNode != nil {
		outputs, err = slices.TransformWithErrorCheck(endNode.Data.Inputs.InputParameters, func(a *vo.Param) (*vo.NamedTypeInfo, error) {
			return convert.BlockInputToNamedTypeInfo(a.Name, a.Input)
		})
		if err != nil {
			logs.Warn(fmt.Sprintf("transform end node inputs to named info failed, err=%v", err))
		}
	}

	return inputs, outputs
}

func (i *impl) Delete(ctx context.Context, policy *vo.DeletePolicy) (ids []int64, err error) {
	if policy.ID != nil || len(policy.IDs) == 1 {
		var id int64
		if policy.ID != nil {
			id = *policy.ID
		} else {
			id = policy.IDs[0]
		}

		if err = i.repo.Delete(ctx, id); err != nil {
			return nil, err
		}

		return []int64{id}, nil
	}

	ids = policy.IDs
	if policy.AppID != nil {
		metas, _, err := i.repo.MGetMetas(ctx, &vo.MetaQuery{
			AppID: policy.AppID,
		})
		if err != nil {
			return nil, err
		}
		ids = maps.Keys(metas)
	}

	if err = i.repo.MDelete(ctx, ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func (i *impl) Get(ctx context.Context, policy *vo.GetPolicy) (*entity.Workflow, error) {
	return i.repo.GetEntity(ctx, policy)
}

func (i *impl) GetWorkflowReference(ctx context.Context, id int64) (map[int64]*vo.Meta, error) {
	parent, err := i.repo.MGetReferences(ctx, &vo.MGetReferencePolicy{
		ReferredIDs:      []int64{id},
		ReferringBizType: []vo.ReferringBizType{vo.ReferringBizTypeWorkflow},
	})
	if err != nil {
		return nil, err
	}

	if len(parent) == 0 {
		// if not parent, it means that it is not cited, so it is returned empty
		return map[int64]*vo.Meta{}, nil
	}

	wfIDs := make(map[int64]struct{}, len(parent))
	for _, ref := range parent {
		wfIDs[ref.ReferringID] = struct{}{}
	}
	ret, _, err := i.repo.MGetMetas(ctx, &vo.MetaQuery{
		IDs: maps.Keys(wfIDs),
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type workflowIdentity struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

func getAllSubWorkflowIdentities(c *vo.Canvas) []*workflowIdentity {
	workflowEntities := make([]*workflowIdentity, 0)

	var collectSubWorkFlowEntities func(nodes []*vo.Node)
	collectSubWorkFlowEntities = func(nodes []*vo.Node) {
		for _, n := range nodes {
			if n.Type == entity.NodeTypeSubWorkflow.IDStr() {
				workflowEntities = append(workflowEntities, &workflowIdentity{
					ID:      n.Data.Inputs.WorkflowID,
					Version: n.Data.Inputs.WorkflowVersion,
				})
			}
			if len(n.Blocks) > 0 {
				collectSubWorkFlowEntities(n.Blocks)
			}
		}
	}

	collectSubWorkFlowEntities(c.Nodes)

	return workflowEntities
}

func (i *impl) ValidateTree(ctx context.Context, id int64, validateConfig vo.ValidateTreeConfig) ([]*cloudworkflow.ValidateTreeInfo, error) {
	wfValidateInfos := make([]*cloudworkflow.ValidateTreeInfo, 0)
	issues, err := validateWorkflowTree(ctx, validateConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to validate work flow: %w", err)
	}

	if len(issues) > 0 {
		wfValidateInfos = append(wfValidateInfos, &cloudworkflow.ValidateTreeInfo{
			WorkflowID: strconv.FormatInt(id, 10),
			Errors:     toValidateErrorData(issues),
		})
	}

	c := &vo.Canvas{}
	err = sonic.UnmarshalString(validateConfig.CanvasSchema, &c)
	if err != nil {
		return nil, vo.WrapError(errno.ErrSerializationDeserializationFail,
			fmt.Errorf("failed to unmarshal canvas schema: %w", err))
	}

	subWorkflowIdentities := getAllSubWorkflowIdentities(c)

	if len(subWorkflowIdentities) > 0 {
		var ids []int64
		for _, e := range subWorkflowIdentities {
			if e.Version != "" {
				continue
			}
			// only project-level workflows need to validate sub-workflows
			ids = append(ids, cast.ToInt64(e.ID)) // TODO: this should be int64 from the start
		}
		if len(ids) == 0 {
			return wfValidateInfos, nil
		}
		workflows, _, err := i.MGet(ctx, &vo.MGetPolicy{
			MetaQuery: vo.MetaQuery{
				IDs: ids,
			},
			QType: workflowModel.FromDraft,
		})
		if err != nil {
			return nil, err
		}

		for _, wf := range workflows {
			issues, err = validateWorkflowTree(ctx, vo.ValidateTreeConfig{
				CanvasSchema: wf.Canvas,
				AppID:        wf.AppID, // application workflow use same app id
			})
			if err != nil {
				return nil, err
			}

			if len(issues) > 0 {
				wfValidateInfos = append(wfValidateInfos, &cloudworkflow.ValidateTreeInfo{
					WorkflowID: strconv.FormatInt(wf.ID, 10),
					Name:       wf.Name,
					Errors:     toValidateErrorData(issues),
				})
			}
		}
	}

	return wfValidateInfos, err
}

func (i *impl) QueryNodeProperties(ctx context.Context, wfID int64) (map[string]*vo.NodeProperty, error) {
	draftInfo, err := i.repo.DraftV2(ctx, wfID, "")
	if err != nil {
		return nil, err
	}

	canvasSchema := draftInfo.Canvas
	if len(canvasSchema) == 0 {
		return nil, fmt.Errorf("no canvas schema")
	}

	mainCanvas := &vo.Canvas{}
	err = sonic.UnmarshalString(canvasSchema, mainCanvas)
	if err != nil {
		return nil, vo.WrapError(errno.ErrSerializationDeserializationFail, err)
	}

	mainCanvas.Nodes, mainCanvas.Edges = adaptor.PruneIsolatedNodes(mainCanvas.Nodes, mainCanvas.Edges, nil)
	nodePropertyMap, err := i.collectNodePropertyMap(ctx, mainCanvas)
	if err != nil {
		return nil, err
	}
	return nodePropertyMap, nil
}

func (i *impl) collectNodePropertyMap(ctx context.Context, canvas *vo.Canvas) (map[string]*vo.NodeProperty, error) {
	nodePropertyMap := make(map[string]*vo.NodeProperty)

	// If it is a nested type, you need to set its parent node
	for _, n := range canvas.Nodes {
		if len(n.Blocks) > 0 {
			for _, nb := range n.Blocks {
				nb.SetParent(n)
			}
		}
	}

	for _, n := range canvas.Nodes {
		if n.Type == entity.NodeTypeSubWorkflow.IDStr() {
			nodeSchema := &schema.NodeSchema{
				Key:  vo.NodeKey(n.ID),
				Type: entity.NodeTypeSubWorkflow,
				Name: n.Data.Meta.Title,
			}
			err := convert.SetInputsForNodeSchema(n, nodeSchema)
			if err != nil {
				return nil, err
			}
			prop := &vo.NodeProperty{
				Type:                nodeSchema.Type.IDStr(),
				IsEnableUserQuery:   isEnableUserQuery(nodeSchema),
				IsEnableChatHistory: isEnableChatHistory(nodeSchema),
				IsRefGlobalVariable: isRefGlobalVariable(nodeSchema),
			}
			nodePropertyMap[string(nodeSchema.Key)] = prop
			wid, err := strconv.ParseInt(n.Data.Inputs.WorkflowID, 10, 64)
			if err != nil {
				return nil, vo.WrapError(errno.ErrSchemaConversionFail, err)
			}

			var canvasSchema string
			if n.Data.Inputs.WorkflowVersion != "" {
				versionInfo, existed, err := i.repo.GetVersion(ctx, wid, n.Data.Inputs.WorkflowVersion)
				if err != nil {
					return nil, err
				}
				if !existed {
					return nil, vo.WrapError(errno.ErrWorkflowNotFound, fmt.Errorf("workflow version %s not found for ID %d: %w", n.Data.Inputs.WorkflowVersion, wid, err), errorx.KV("id", strconv.FormatInt(wid, 10)))
				}
				canvasSchema = versionInfo.Canvas
			} else {
				draftInfo, err := i.repo.DraftV2(ctx, wid, "")
				if err != nil {
					return nil, err
				}
				canvasSchema = draftInfo.Canvas
			}

			if len(canvasSchema) == 0 {
				return nil, fmt.Errorf("workflow id %v ,not get canvas schema, version %v", wid, n.Data.Inputs.WorkflowVersion)
			}

			c := &vo.Canvas{}
			err = sonic.UnmarshalString(canvasSchema, c)
			if err != nil {
				return nil, vo.WrapError(errno.ErrSchemaConversionFail, err)
			}
			ret, err := i.collectNodePropertyMap(ctx, c)
			if err != nil {
				return nil, err
			}
			prop.SubWorkflow = ret

		} else {
			nodeSchemas, _, err := adaptor.NodeToNodeSchema(ctx, n, canvas)
			if err != nil {
				return nil, err
			}
			for _, nodeSchema := range nodeSchemas {
				nodePropertyMap[string(nodeSchema.Key)] = &vo.NodeProperty{
					Type:                nodeSchema.Type.IDStr(),
					IsEnableUserQuery:   isEnableUserQuery(nodeSchema),
					IsEnableChatHistory: isEnableChatHistory(nodeSchema),
					IsRefGlobalVariable: isRefGlobalVariable(nodeSchema),
				}
			}

		}
	}
	return nodePropertyMap, nil
}

func isEnableUserQuery(s *schema.NodeSchema) bool {
	if s == nil {
		return false
	}
	if s.Type != entity.NodeTypeEntry {
		return false
	}

	if len(s.OutputSources) == 0 {
		return false
	}

	for _, source := range s.OutputSources {
		fieldPath := source.Path
		if len(fieldPath) == 1 && (fieldPath[0] == "BOT_USER_INPUT" || fieldPath[0] == "USER_INPUT") {
			return true
		}
	}

	return false
}

func isEnableChatHistory(s *schema.NodeSchema) bool {
	if s == nil {
		return false
	}

	chatHistoryAware, ok := s.Configs.(schema.ChatHistoryAware)
	if !ok {
		return false
	}

	return chatHistoryAware.ChatHistoryEnabled()
}

func isRefGlobalVariable(s *schema.NodeSchema) bool {
	for _, source := range s.InputSources {
		if source.IsRefGlobalVariable() {
			return true
		}
	}
	for _, source := range s.OutputSources {
		if source.IsRefGlobalVariable() {
			return true
		}
	}
	return false
}

func canvasToRefs(referringID int64, canvasStr string) (map[entity.WorkflowReferenceKey]struct{}, error) {
	var canvas vo.Canvas
	if err := sonic.UnmarshalString(canvasStr, &canvas); err != nil {
		return nil, vo.WrapError(errno.ErrSerializationDeserializationFail, err)
	}

	wfRefs := map[entity.WorkflowReferenceKey]struct{}{}
	var getRefFn func([]*vo.Node) error
	getRefFn = func(nodes []*vo.Node) error {
		for _, node := range nodes {
			if node.Type == entity.NodeTypeSubWorkflow.IDStr() {
				referredID, err := strconv.ParseInt(node.Data.Inputs.WorkflowID, 10, 64)
				if err != nil {
					return vo.WrapError(errno.ErrSchemaConversionFail, err)
				}
				wfRefs[entity.WorkflowReferenceKey{
					ReferredID:       referredID,
					ReferringID:      referringID,
					ReferType:        vo.ReferTypeSubWorkflow,
					ReferringBizType: vo.ReferringBizTypeWorkflow,
				}] = struct{}{}
			} else if node.Type == entity.NodeTypeLLM.IDStr() {
				if node.Data.Inputs.LLM != nil {
					if node.Data.Inputs.FCParam != nil && node.Data.Inputs.FCParam.WorkflowFCParam != nil {
						for _, w := range node.Data.Inputs.FCParam.WorkflowFCParam.WorkflowList {
							referredID, err := strconv.ParseInt(w.WorkflowID, 10, 64)
							if err != nil {
								return vo.WrapError(errno.ErrSchemaConversionFail, err)
							}
							wfRefs[entity.WorkflowReferenceKey{
								ReferredID:       referredID,
								ReferringID:      referringID,
								ReferType:        vo.ReferTypeTool,
								ReferringBizType: vo.ReferringBizTypeWorkflow,
							}] = struct{}{}
						}
					}
				}
			} else if len(node.Blocks) > 0 {
				for _, subNode := range node.Blocks {
					if err := getRefFn([]*vo.Node{subNode}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}

	if err := getRefFn(canvas.Nodes); err != nil {
		return nil, err
	}

	return wfRefs, nil
}

func (i *impl) Publish(ctx context.Context, policy *vo.PublishPolicy) (err error) {
	meta, err := i.repo.GetMeta(ctx, policy.ID)
	if err != nil {
		return err
	}
	if meta.Mode != cloudworkflow.WorkflowMode_Workflow {
		return fmt.Errorf("only standard workflows can be published")
	}
	if meta.AppID != nil && *meta.AppID != 0 {
		return fmt.Errorf("application workflows are not supported")
	}

	if meta.LatestPublishedVersion != nil {
		latestVersion, err := parseVersion(*meta.LatestPublishedVersion)
		if err != nil {
			return err
		}
		currentVersion, err := parseVersion(policy.Version)
		if err != nil {
			return err
		}

		if !isIncremental(latestVersion, currentVersion) {
			return fmt.Errorf("the version number is not self-incrementing, old version %v, current version is %v", *meta.LatestPublishedVersion, policy.Version)
		}
	}

	draft, err := i.repo.DraftV2(ctx, policy.ID, policy.CommitID)
	if err != nil {
		return err
	}

	if !policy.Force && !draft.TestRunSuccess {
		return fmt.Errorf("workflow %d's current draft needs to pass the test run before publishing", policy.ID)
	}

	wfRefs, err := canvasToRefs(policy.ID, draft.Canvas)
	if err != nil {
		return err
	}

	versionInfo := &vo.VersionInfo{
		VersionMeta: &vo.VersionMeta{
			Version:            policy.Version,
			VersionDescription: policy.VersionDescription,
			VersionCreatorID:   policy.CreatorID,
		},
		CanvasInfo: vo.CanvasInfo{
			Canvas:          draft.Canvas,
			InputParamsStr:  draft.InputParamsStr,
			OutputParamsStr: draft.OutputParamsStr,
		},
		CommitID: draft.CommitID,
	}

	if err = i.repo.CreateVersion(ctx, policy.ID, versionInfo, wfRefs); err != nil {
		return err
	}

	return nil
}

func (i *impl) UpdateMeta(ctx context.Context, id int64, metaUpdate *vo.MetaUpdate) (err error) {
	if metaUpdate.WorkflowMode != nil && *metaUpdate.WorkflowMode != cloudworkflow.WorkflowMode_Workflow {
		return fmt.Errorf("only standard workflows are supported")
	}
	meta, err := i.repo.GetMeta(ctx, id)
	if err != nil {
		return err
	}
	if meta.Mode != cloudworkflow.WorkflowMode_Workflow {
		return fmt.Errorf("only standard workflows can be edited")
	}
	if meta.AppID != nil && *meta.AppID != 0 {
		return fmt.Errorf("application workflows are not supported")
	}
	err = i.repo.UpdateMeta(ctx, id, metaUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (i *impl) CopyWorkflow(ctx context.Context, workflowID int64, policy vo.CopyWorkflowPolicy) (*entity.Workflow, error) {
	wf, err := i.repo.CopyWorkflow(ctx, workflowID, policy)
	if err != nil {
		return nil, err
	}
	// chat flow should copy role config

	return wf, nil

}

func (i *impl) checkBotAgentNode(node *vo.Node) error {
	if node.Type == entity.NodeTypeCreateConversation.IDStr() || node.Type == entity.NodeTypeConversationDelete.IDStr() || node.Type == entity.NodeTypeConversationUpdate.IDStr() || node.Type == entity.NodeTypeConversationList.IDStr() {
		return errors.New("conversation-related nodes are not supported in chatflow")
	}
	return nil
}

func (i *impl) validateNodesRecursively(ctx context.Context, nodes []*vo.Node, checkType cloudworkflow.CheckType, visited map[string]struct{}, repo workflow.Repository) error {
	queue := make([]*vo.Node, 0, len(nodes))
	queue = append(queue, nodes...)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node == nil {
			continue
		}

		var checkErr error
		switch checkType {
		case cloudworkflow.CheckType_BotAgent:
			checkErr = i.checkBotAgentNode(node)
		default:
			// For now, we only handle BotAgent check, so we can do nothing here.
			// In the future, if there are other check types that need to be validated on every node, this logic will need to be adjusted.
		}
		if checkErr != nil {
			return checkErr
		}

		// Enqueue nested nodes for BFS traversal. This handles Loop, Batch, and other nodes with nested blocks.
		if len(node.Blocks) > 0 {
			queue = append(queue, node.Blocks...)
		}

		if node.Type == entity.NodeTypeSubWorkflow.IDStr() && node.Data != nil && node.Data.Inputs != nil {
			workflowIDStr := node.Data.Inputs.WorkflowID
			if workflowIDStr == "" {
				continue
			}

			workflowID, err := strconv.ParseInt(workflowIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid workflow ID in sub-workflow node %s: %w", node.ID, err)
			}

			version := node.Data.Inputs.WorkflowVersion
			qType := workflowModel.FromDraft
			if version != "" {
				qType = workflowModel.FromSpecificVersion
			}

			visitedKey := fmt.Sprintf("%d:%s", workflowID, version)
			if _, ok := visited[visitedKey]; ok {
				continue
			}
			visited[visitedKey] = struct{}{}

			subWfEntity, err := repo.GetEntity(ctx, &vo.GetPolicy{
				ID:      workflowID,
				QType:   qType,
				Version: version,
			})
			if err != nil {
				delete(visited, visitedKey)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
				return fmt.Errorf("failed to get sub-workflow entity %d: %w", workflowID, err)
			}

			var canvas vo.Canvas
			if err := sonic.UnmarshalString(subWfEntity.Canvas, &canvas); err != nil {
				return fmt.Errorf("failed to unmarshal canvas for workflow %d: %w", subWfEntity.ID, err)
			}

			queue = append(queue, canvas.Nodes...)
		}

		if node.Type == entity.NodeTypeLLM.IDStr() && node.Data != nil && node.Data.Inputs != nil && node.Data.Inputs.LLM != nil && node.Data.Inputs.FCParam != nil && node.Data.Inputs.FCParam.WorkflowFCParam != nil {
			for _, subWfInfo := range node.Data.Inputs.FCParam.WorkflowFCParam.WorkflowList {
				if subWfInfo.WorkflowID == "" {
					continue
				}
				workflowID, err := strconv.ParseInt(subWfInfo.WorkflowID, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid workflow ID in large model node %s: %w", node.ID, err)
				}

				version := subWfInfo.WorkflowVersion
				qType := workflowModel.FromDraft
				if version != "" {
					qType = workflowModel.FromSpecificVersion
				}

				visitedKey := fmt.Sprintf("%d:%s", workflowID, version)
				if _, ok := visited[visitedKey]; ok {
					continue
				}
				visited[visitedKey] = struct{}{}

				subWfEntity, err := repo.GetEntity(ctx, &vo.GetPolicy{
					ID:      workflowID,
					QType:   qType,
					Version: version,
				})
				if err != nil {
					delete(visited, visitedKey)
					if errors.Is(err, gorm.ErrRecordNotFound) {
						continue
					}
					return fmt.Errorf("failed to get sub-workflow entity %d from large model node: %w", workflowID, err)
				}

				var canvas vo.Canvas
				if err := sonic.UnmarshalString(subWfEntity.Canvas, &canvas); err != nil {
					return fmt.Errorf("failed to unmarshal canvas for workflow %d from large model node: %w", subWfEntity.ID, err)
				}

				queue = append(queue, canvas.Nodes...)
			}
		}
	}
	return nil
}

func (i *impl) WorkflowSchemaCheck(ctx context.Context, wf *entity.Workflow, checks []cloudworkflow.CheckType) ([]*cloudworkflow.CheckResult, error) {
	checkResults := make([]*cloudworkflow.CheckResult, 0, len(checks))

	var canvas vo.Canvas
	if err := sonic.UnmarshalString(wf.Canvas, &canvas); err != nil {
		return nil, fmt.Errorf("failed to unmarshal canvas for workflow %d: %w", wf.ID, err)
	}

	for _, checkType := range checks {
		visited := make(map[string]struct{})
		visitedKey := fmt.Sprintf("%d:%s", wf.ID, wf.GetVersion())
		visited[visitedKey] = struct{}{}

		err := i.validateNodesRecursively(ctx, canvas.Nodes, checkType, visited, i.repo)

		if err != nil {
			checkResults = append(checkResults, &cloudworkflow.CheckResult{
				IsPass: false,
				Reason: err.Error(),
				Type:   checkType,
			})
		} else {
			checkResults = append(checkResults, &cloudworkflow.CheckResult{
				IsPass: true,
				Type:   checkType,
				Reason: "",
			})
		}
	}
	return checkResults, nil
}

func (i *impl) MGet(ctx context.Context, policy *vo.MGetPolicy) ([]*entity.Workflow, int64, error) {
	if policy.MetaOnly {
		metas, total, err := i.repo.MGetMetas(ctx, &policy.MetaQuery)
		if err != nil {
			return nil, 0, err
		}

		result := make([]*entity.Workflow, len(metas))
		var index int

		if len(metas) == 0 {
			return result, 0, nil
		}

		for id := range metas {
			wf := &entity.Workflow{
				ID:   id,
				Meta: metas[id],
			}
			result[index] = wf
			index++
		}
		return result, total, nil
	}

	ioF := func(inputParam, outputParam string) (input []*vo.NamedTypeInfo, output []*vo.NamedTypeInfo, err error) {
		if inputParam != "" {
			err := sonic.UnmarshalString(inputParam, &input)
			if err != nil {
				return nil, nil, err
			}
		}

		if outputParam != "" {
			err := sonic.UnmarshalString(outputParam, &output)
			if err != nil {
				return nil, nil, err
			}
		}

		return input, output, err
	}

	switch policy.QType {
	case workflowModel.FromDraft:
		return i.repo.MGetDrafts(ctx, policy)
	case workflowModel.FromSpecificVersion:
		if len(policy.IDs) == 0 || len(policy.Versions) != len(policy.IDs) {
			return nil, 0, fmt.Errorf("ids and versions are required when MGet from specific versions")
		}

		metas, total, err := i.repo.MGetMetas(ctx, &policy.MetaQuery)
		if err != nil {
			return nil, total, err
		}

		result := make([]*entity.Workflow, len(metas))
		index := 0

		for id, version := range policy.Versions {
			v, existed, err := i.repo.GetVersion(ctx, id, version)
			if err != nil {
				return nil, total, err
			}
			if !existed {
				return nil, total, vo.WrapError(errno.ErrWorkflowNotFound, fmt.Errorf("workflow version %s not found for ID %d: %w", version, id, err), errorx.KV("id", strconv.FormatInt(id, 10)))
			}
			inputs, outputs, err := ioF(v.InputParamsStr, v.OutputParamsStr)
			if err != nil {
				return nil, total, err
			}

			wf := &entity.Workflow{
				ID:       id,
				Meta:     metas[id],
				CommitID: v.CommitID,
				CanvasInfo: &vo.CanvasInfo{
					Canvas:          v.Canvas,
					InputParams:     inputs,
					OutputParams:    outputs,
					InputParamsStr:  v.InputParamsStr,
					OutputParamsStr: v.OutputParamsStr,
				},
				VersionMeta: v.VersionMeta,
			}
			result[index] = wf
			index++
		}

		return result, total, nil
	case workflowModel.FromLatestVersion:
		return i.repo.MGetLatestVersion(ctx, policy)
	default:
		panic("not implemented")
	}
}

func (i *impl) calculateTestRunSuccess(ctx context.Context, c *vo.Canvas, wid int64) (bool, error) {
	sc, err := adaptor.CanvasToWorkflowSchema(ctx, c)
	if err != nil { // not even legal, test run can't possibly be successful
		return false, nil
	}

	existedDraft, err := i.repo.DraftV2(ctx, wid, "")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // previous draft version not exists, does not have any test run
		}
		return false, err
	}

	var existedDraftCanvas vo.Canvas
	err = sonic.UnmarshalString(existedDraft.Canvas, &existedDraftCanvas)
	existedSc, err := adaptor.CanvasToWorkflowSchema(ctx, &existedDraftCanvas)
	if err == nil { // the old existing draft is legal, check if it's equal to the new draft in terms of execution
		if !existedSc.IsEqual(sc) { // there is modification to the execution logic, needs new test run
			return false, nil
		}
	} else { // the old existing draft is not legal, of course haven't any successful test run
		return false, nil
	}

	return existedDraft.TestRunSuccess, nil // inherit previous draft snapshot's test run success flag
}

func replaceRelatedWorkflowOrExternalResourceInWorkflowNodes(nodes []*vo.Node, relatedWorkflows map[int64]entity.IDVersionPair, related vo.ExternalResourceRelated) error {
	var (
		hasWorkflowRelated  = len(relatedWorkflows) > 0
		hasPluginRelated    = len(related.PluginMap) > 0
		hasKnowledgeRelated = len(related.KnowledgeMap) > 0
		hasDatabaseRelated  = len(related.DatabaseMap) > 0
	)

	for _, node := range nodes {
		nType := entity.IDStrToNodeType(node.Type)
		meta := entity.NodeMetaByNodeType(nType)
		if meta.UseDatabase {
			if !hasDatabaseRelated || node.Data.Inputs.DatabaseNode == nil {
				continue
			}
			dsList := node.Data.Inputs.DatabaseInfoList
			for idx := range dsList {
				databaseInfo := dsList[idx]
				did, err := strconv.ParseInt(databaseInfo.DatabaseInfoID, 10, 64)
				if err != nil {
					return err
				}
				if refDatabaseID, ok := related.DatabaseMap[did]; ok {
					databaseInfo.DatabaseInfoID = strconv.FormatInt(refDatabaseID, 10)
				}

			}
			continue
		}

		if meta.UseKnowledge {
			if !hasKnowledgeRelated || node.Data.Inputs.Knowledge == nil {
				continue
			}
			datasetListInfoParam := node.Data.Inputs.DatasetParam[0]
			knowledgeIDs := datasetListInfoParam.Input.Value.Content.([]any)
			for idx := range knowledgeIDs {
				kid, err := strconv.ParseInt(knowledgeIDs[idx].(string), 10, 64)
				if err != nil {
					return err
				}
				if refKnowledgeID, ok := related.KnowledgeMap[kid]; ok {
					knowledgeIDs[idx] = strconv.FormatInt(refKnowledgeID, 10)
				}
			}

			continue
		}

		if meta.UsePlugin {
			if !hasPluginRelated || node.Data.Inputs.PluginAPIParam == nil {
				continue
			}
			apiParams := slices.ToMap(node.Data.Inputs.APIParams, func(e *vo.Param) (string, *vo.Param) {
				return e.Name, e
			})
			pluginIDParam, ok := apiParams["pluginID"]
			if !ok {
				return fmt.Errorf("plugin id param is not found")
			}

			pID, err := strconv.ParseInt(pluginIDParam.Input.Value.Content.(string), 10, 64)
			if err != nil {
				return err
			}

			pluginVersionParam, ok := apiParams["pluginVersion"]
			if !ok {
				return fmt.Errorf("plugin version param is not found")
			}

			if refPlugin, ok := related.PluginMap[pID]; ok {
				pluginIDParam.Input.Value.Content = strconv.FormatInt(refPlugin.PluginID, 10)
				if refPlugin.PluginVersion != nil {
					pluginVersionParam.Input.Value.Content = *refPlugin.PluginVersion
				}
			}

			apiIDParam, ok := apiParams["apiID"]
			if !ok {
				return fmt.Errorf("apiID param is not found")
			}

			apiID, err := strconv.ParseInt(apiIDParam.Input.Value.Content.(string), 10, 64)
			if err != nil {
				return err
			}

			if refApiID, ok := related.PluginToolMap[apiID]; ok {
				apiIDParam.Input.Value.Content = strconv.FormatInt(refApiID, 10)
			}

			continue
		}

		switch nType {
		case entity.NodeTypeSubWorkflow:
			if !hasWorkflowRelated || node.Data.Inputs.SubWorkflow == nil {
				continue
			}
			workflowID, err := strconv.ParseInt(node.Data.Inputs.WorkflowID, 10, 64)
			if err != nil {
				return err
			}
			if wf, ok := relatedWorkflows[workflowID]; ok {
				node.Data.Inputs.WorkflowID = strconv.FormatInt(wf.ID, 10)
				node.Data.Inputs.WorkflowVersion = wf.Version
			}
		case entity.NodeTypeLLM:
			if node.Data.Inputs.LLM == nil {
				continue
			}
			if hasWorkflowRelated && node.Data.Inputs.FCParam != nil && node.Data.Inputs.FCParam.WorkflowFCParam != nil {
				for idx := range node.Data.Inputs.FCParam.WorkflowFCParam.WorkflowList {
					wf := node.Data.Inputs.FCParam.WorkflowFCParam.WorkflowList[idx]
					workflowID, err := strconv.ParseInt(wf.WorkflowID, 10, 64)
					if err != nil {
						return err
					}
					if refWf, ok := relatedWorkflows[workflowID]; ok {
						wf.WorkflowID = strconv.FormatInt(refWf.ID, 10)
						wf.WorkflowVersion = refWf.Version
					}
					node.Data.Inputs.FCParam.WorkflowFCParam.WorkflowList[idx] = wf
				}

			}
			if hasPluginRelated && node.Data.Inputs.FCParam != nil && node.Data.Inputs.FCParam.PluginFCParam != nil {
				for idx := range node.Data.Inputs.FCParam.PluginFCParam.PluginList {
					pl := node.Data.Inputs.FCParam.PluginFCParam.PluginList[idx]
					pluginID, err := strconv.ParseInt(pl.PluginID, 10, 64)
					if err != nil {
						return err
					}

					toolID, err := strconv.ParseInt(pl.ApiId, 10, 64)
					if err != nil {
						return err
					}

					if refPlugin, ok := related.PluginMap[pluginID]; ok {
						tID, ok := related.PluginToolMap[toolID]
						if ok {
							pl.ApiId = strconv.FormatInt(tID, 10)
						}
						pl.PluginID = strconv.FormatInt(refPlugin.PluginID, 10)
						if refPlugin.PluginVersion != nil {
							pl.PluginVersion = *refPlugin.PluginVersion

							pl.IsDraft = false
						}

					}
					node.Data.Inputs.FCParam.PluginFCParam.PluginList[idx] = pl

				}
			}
			if hasKnowledgeRelated && node.Data.Inputs.FCParam != nil && node.Data.Inputs.FCParam.KnowledgeFCParam != nil {
				for idx := range node.Data.Inputs.FCParam.KnowledgeFCParam.KnowledgeList {
					kn := node.Data.Inputs.FCParam.KnowledgeFCParam.KnowledgeList[idx]
					kid, err := strconv.ParseInt(kn.ID, 10, 64)
					if err != nil {
						return err
					}
					if refKnowledgeID, ok := related.KnowledgeMap[kid]; ok {
						kn.ID = strconv.FormatInt(refKnowledgeID, 10)
					}

				}
			}
		}
		if len(node.Blocks) > 0 {
			err := replaceRelatedWorkflowOrExternalResourceInWorkflowNodes(node.Blocks, relatedWorkflows, related)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func RegisterAllNodeAdaptors() {
	adaptor.RegisterAllNodeAdaptors()
}
