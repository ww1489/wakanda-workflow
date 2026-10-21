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

/* eslint-disable */
/* tslint:disable */
// @ts-nocheck

import * as agent_provider from './namespaces/agent_provider';
import * as base from './namespaces/base';
import * as bot_schema from './namespaces/bot_schema';
import * as copilot from './namespaces/copilot';
import * as copilot_common from './namespaces/copilot_common';
import * as event from './namespaces/event';
import * as hook_common from './namespaces/hook_common';
import * as kctx from './namespaces/kctx';
import * as permission from './namespaces/permission';
import * as resource from './namespaces/resource';
import * as resource_common from './namespaces/resource_common';
import * as trace from './namespaces/trace';
import * as trigger from './namespaces/trigger';
import * as wf_openapi from './namespaces/wf_openapi';
import * as workflow from './namespaces/workflow';

export {
  agent_provider,
  base,
  bot_schema,
  copilot,
  copilot_common,
  event,
  hook_common,
  kctx,
  permission,
  resource,
  resource_common,
  trace,
  trigger,
  wf_openapi,
  workflow,
};
export * from './namespaces/agent_provider';
export * from './namespaces/base';
export * from './namespaces/bot_schema';
export * from './namespaces/copilot';
export * from './namespaces/copilot_common';
export * from './namespaces/event';
export * from './namespaces/hook_common';
export * from './namespaces/kctx';
export * from './namespaces/permission';
export * from './namespaces/resource';
export * from './namespaces/resource_common';
export * from './namespaces/trace';
export * from './namespaces/trigger';
export * from './namespaces/wf_openapi';
export * from './namespaces/workflow';

export type Int64 = string | number;

export default class WorkflowApiService<T> {
  private request: any = () => {
    throw new Error('WorkflowApiService.request is undefined');
  };
  private baseURL: string | ((path: string) => string) = '';

  constructor(options?: {
    baseURL?: string | ((path: string) => string);
    request?<R>(
      params: {
        url: string;
        method: 'GET' | 'DELETE' | 'POST' | 'PUT' | 'PATCH';
        data?: any;
        params?: any;
        headers?: any;
      },
      options?: T,
    ): Promise<R>;
  }) {
    this.request = options?.request || this.request;
    this.baseURL = options?.baseURL || '';
  }

  private genBaseURL(path: string) {
    return typeof this.baseURL === 'string'
      ? this.baseURL + path
      : this.baseURL(path);
  }

  /**
   * POST /api/workflow_api/create
   *
   * HTTP
   *
   * 创建流程
   */
  CreateWorkflow(
    req: workflow.CreateWorkflowRequest,
    options?: T,
  ): Promise<workflow.CreateWorkflowResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/create');
    const method = 'POST';
    const data = {
      name: _req['name'],
      desc: _req['desc'],
      icon_uri: _req['icon_uri'],
      space_id: _req['space_id'],
      flow_mode: _req['flow_mode'],
      schema_type: _req['schema_type'],
      bind_biz_id: _req['bind_biz_id'],
      bind_biz_type: _req['bind_biz_type'],
      project_id: _req['project_id'],
      create_conversation: _req['create_conversation'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/save
   *
   * 保存流程
   */
  SaveWorkflow(
    req: workflow.SaveWorkflowRequest,
    options?: T,
  ): Promise<workflow.SaveWorkflowResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/save');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      schema: _req['schema'],
      space_id: _req['space_id'],
      name: _req['name'],
      desc: _req['desc'],
      icon_uri: _req['icon_uri'],
      submit_commit_id: _req['submit_commit_id'],
      ignore_status_transfer: _req['ignore_status_transfer'],
      save_version: _req['save_version'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/delete */
  DeleteWorkflow(
    req: workflow.DeleteWorkflowRequest,
    options?: T,
  ): Promise<workflow.DeleteWorkflowResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/delete');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      action: _req['action'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/workflow_references
   *
   * 查询指定工作流被哪些其他工作流所引用。
   */
  GetWorkflowReferences(
    req: workflow.GetWorkflowReferencesRequest,
    options?: T,
  ): Promise<workflow.GetWorkflowReferencesResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/workflow_references');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/released_workflows */
  GetReleasedWorkflows(
    req?: workflow.GetReleasedWorkflowsRequest,
    options?: T,
  ): Promise<workflow.GetReleasedWorkflowsResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/released_workflows');
    const method = 'POST';
    const data = {
      page: _req['page'],
      size: _req['size'],
      type: _req['type'],
      name: _req['name'],
      workflow_ids: _req['workflow_ids'],
      tags: _req['tags'],
      space_id: _req['space_id'],
      order_by: _req['order_by'],
      login_user_create: _req['login_user_create'],
      flow_mode: _req['flow_mode'],
      workflow_filter_list: _req['workflow_filter_list'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/canvas
   *
   * 获取流程编辑态详情，包括画布信息、节点信息、边信息、变量信息、权限信息、版本信息等
   */
  GetCanvasInfo(
    req: workflow.GetCanvasInfoRequest,
    options?: T,
  ): Promise<workflow.GetCanvasInfoResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/canvas');
    const method = 'POST';
    const data = {
      space_id: _req['space_id'],
      workflow_id: _req['workflow_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/history_schema */
  GetHistorySchema(
    req: workflow.GetHistorySchemaRequest,
    options?: T,
  ): Promise<workflow.GetHistorySchemaResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/history_schema');
    const method = 'POST';
    const data = {
      space_id: _req['space_id'],
      workflow_id: _req['workflow_id'],
      commit_id: _req['commit_id'],
      type: _req['type'],
      env: _req['env'],
      workflow_version: _req['workflow_version'],
      project_version: _req['project_version'],
      project_id: _req['project_id'],
      execute_id: _req['execute_id'],
      sub_execute_id: _req['sub_execute_id'],
      log_id: _req['log_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/workflow_list
   *
   * 获取流程列表。
   */
  GetWorkFlowList(
    req?: workflow.GetWorkFlowListRequest,
    options?: T,
  ): Promise<workflow.GetWorkFlowListResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/workflow_list');
    const method = 'POST';
    const data = {
      page: _req['page'],
      size: _req['size'],
      workflow_ids: _req['workflow_ids'],
      type: _req['type'],
      name: _req['name'],
      tags: _req['tags'],
      space_id: _req['space_id'],
      status: _req['status'],
      order_by: _req['order_by'],
      login_user_create: _req['login_user_create'],
      flow_mode: _req['flow_mode'],
      schema_type_list: _req['schema_type_list'],
      project_id: _req['project_id'],
      checker: _req['checker'],
      bind_biz_id: _req['bind_biz_id'],
      bind_biz_type: _req['bind_biz_type'],
      project_version: _req['project_version'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/node_type
   *
   * 查询流程中的节点类型
   */
  QueryWorkflowNodeTypes(
    req?: workflow.QueryWorkflowNodeTypeRequest,
    options?: T,
  ): Promise<workflow.QueryWorkflowNodeTypeResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/node_type');
    const method = 'POST';
    const data = {
      space_id: _req['space_id'],
      workflow_id: _req['workflow_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/test_run
   *
   * 试运行流程（test run）异步接口，需要轮询GetWorkFlowProcess接口来进行流程运行结果的检查
   */
  WorkFlowTestRun(
    req: workflow.WorkFlowTestRunRequest,
    options?: T,
  ): Promise<workflow.WorkFlowTestRunResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/test_run');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      input: _req['input'],
      space_id: _req['space_id'],
      bot_id: _req['bot_id'],
      submit_commit_id: _req['submit_commit_id'],
      commit_id: _req['commit_id'],
      project_id: _req['project_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/publish
   *
   * 发布流程。该接口的用途是发布非 project 内部的流程。
   */
  PublishWorkflow(
    req: workflow.PublishWorkflowRequest,
    options?: T,
  ): Promise<workflow.PublishWorkflowResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/publish');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      has_collaborator: _req['has_collaborator'],
      env: _req['env'],
      commit_id: _req['commit_id'],
      force: _req['force'],
      workflow_version: _req['workflow_version'],
      version_description: _req['version_description'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * GET /api/workflow_api/get_process
   *
   * 查看试运行执行历史。在试运行时需要轮询查看一个流程的试运行执行历史。
   */
  GetWorkFlowProcess(
    req: workflow.GetWorkflowProcessRequest,
    options?: T,
  ): Promise<workflow.GetWorkflowProcessResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/get_process');
    const method = 'GET';
    const params = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      execute_id: _req['execute_id'],
      sub_execute_id: _req['sub_execute_id'],
      need_async: _req['need_async'],
      log_id: _req['log_id'],
      node_id: _req['node_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, params }, options);
  }

  /** POST /api/workflow_api/copy */
  CopyWorkflow(
    req: workflow.CopyWorkflowRequest,
    options?: T,
  ): Promise<workflow.CopyWorkflowResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/copy');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/upload/auth_token
   *
   * 文件上传
   */
  GetUploadAuthToken(
    req?: workflow.GetUploadAuthTokenRequest,
    options?: T,
  ): Promise<workflow.GetUploadAuthTokenResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/upload/auth_token');
    const method = 'POST';
    const data = { scene: _req['scene'], Base: _req['Base'] };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/sign_image_url */
  SignImageURL(
    req: workflow.SignImageURLRequest,
    options?: T,
  ): Promise<workflow.SignImageURLResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/sign_image_url');
    const method = 'POST';
    const data = { uri: _req['uri'], Scene: _req['Scene'], Base: _req['Base'] };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/cancel */
  CancelWorkFlow(
    req: workflow.CancelWorkFlowRequest,
    options?: T,
  ): Promise<workflow.CancelWorkFlowResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/cancel');
    const method = 'POST';
    const data = {
      execute_id: _req['execute_id'],
      space_id: _req['space_id'],
      workflow_id: _req['workflow_id'],
      async_subflow: _req['async_subflow'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/node_template_list
   *
   * 查询节点模板列表
   */
  NodeTemplateList(
    req?: workflow.NodeTemplateListRequest,
    options?: T,
  ): Promise<workflow.NodeTemplateListResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/node_template_list');
    const method = 'POST';
    const data = {
      need_types: _req['need_types'],
      node_types: _req['node_types'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /v1/workflow/run
   *
   * OpenAPI
   *
   * OpenAPI运行流程。执行已发布的工作流 (非流式)
   */
  OpenAPIRunFlow(
    req?: workflow.OpenAPIRunFlowRequest,
    options?: T,
  ): Promise<workflow.OpenAPIRunFlowResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/v1/workflow/run');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      parameters: _req['parameters'],
      ext: _req['ext'],
      bot_id: _req['bot_id'],
      is_async: _req['is_async'],
      execute_mode: _req['execute_mode'],
      version: _req['version'],
      connector_id: _req['connector_id'],
      app_id: _req['app_id'],
      project_id: _req['project_id'],
      app_version: _req['app_version'],
      workflow_version: _req['workflow_version'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/nodeDebug */
  WorkflowNodeDebugV2(
    req?: workflow.WorkflowNodeDebugV2Request,
    options?: T,
  ): Promise<workflow.WorkflowNodeDebugV2Response> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/nodeDebug');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      node_id: _req['node_id'],
      input: _req['input'],
      batch: _req['batch'],
      space_id: _req['space_id'],
      bot_id: _req['bot_id'],
      project_id: _req['project_id'],
      setting: _req['setting'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * GET /api/workflow_api/apiDetail
   *
   * 获取workflow引用的插件工具详情
   */
  GetApiDetail(
    req?: workflow.GetApiDetailRequest,
    options?: T,
  ): Promise<workflow.GetApiDetailResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/apiDetail');
    const method = 'GET';
    const params = {
      pluginID: _req['pluginID'],
      apiName: _req['apiName'],
      space_id: _req['space_id'],
      api_id: _req['api_id'],
      project_id: _req['project_id'],
      plugin_version: _req['plugin_version'],
      plugin_from: _req['plugin_from'],
      Base: _req['Base'],
    };
    return this.request({ url, method, params }, options);
  }

  /** POST /api/workflow_api/delete_strategy */
  GetDeleteStrategy(
    req: workflow.GetDeleteStrategyRequest,
    options?: T,
  ): Promise<workflow.GetDeleteStrategyResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/delete_strategy');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/copy_wk_template */
  CopyWkTemplateApi(
    req: workflow.CopyWkTemplateApiRequest,
    options?: T,
  ): Promise<workflow.CopyWkTemplateApiResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/copy_wk_template');
    const method = 'POST';
    const data = {
      workflow_ids: _req['workflow_ids'],
      target_space_id: _req['target_space_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/update_meta */
  UpdateWorkflowMeta(
    req: workflow.UpdateWorkflowMetaRequest,
    options?: T,
  ): Promise<workflow.UpdateWorkflowMetaResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/update_meta');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      name: _req['name'],
      desc: _req['desc'],
      icon_uri: _req['icon_uri'],
      flow_mode: _req['flow_mode'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/test_resume
   *
   * 流程中断后的恢复，部分流程在运行中会中断，在流程挂起后需要用该接口恢复运行流程
   */
  WorkFlowTestResume(
    req: workflow.WorkflowTestResumeRequest,
    options?: T,
  ): Promise<workflow.WorkflowTestResumeResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/test_resume');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      execute_id: _req['execute_id'],
      event_id: _req['event_id'],
      data: _req['data'],
      space_id: _req['space_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /v1/workflow/stream_run
   *
   * OpenAPI流式运行流程。执行资源库内已发布的工作流 或 project内的工作流，响应方式为流式响应。
   */
  OpenAPIStreamRunFlow(
    req?: workflow.OpenAPIRunFlowRequest,
    options?: T,
  ): Promise<workflow.OpenAPIStreamRunFlowResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/v1/workflow/stream_run');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      parameters: _req['parameters'],
      ext: _req['ext'],
      bot_id: _req['bot_id'],
      is_async: _req['is_async'],
      execute_mode: _req['execute_mode'],
      version: _req['version'],
      connector_id: _req['connector_id'],
      app_id: _req['app_id'],
      project_id: _req['project_id'],
      app_version: _req['app_version'],
      workflow_version: _req['workflow_version'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * GET /v1/workflow/get_run_history
   *
   * OpenAPI查询工作流异步执行结果。工作流异步运行后，查看执行结果。
   */
  OpenAPIGetWorkflowRunHistory(
    req: workflow.GetWorkflowRunHistoryRequest,
    options?: T,
  ): Promise<workflow.GetWorkflowRunHistoryResponse> {
    const _req = req;
    const url = this.genBaseURL('/v1/workflow/get_run_history');
    const method = 'GET';
    const params = {
      workflow_id: _req['workflow_id'],
      execute_id: _req['execute_id'],
    };
    return this.request({ url, method, params }, options);
  }

  /**
   * POST /v1/workflow/stream_resume
   *
   * OpenAPI恢复运行工作流。恢复运行已中断的工作流，响应方式为流式响应。
   */
  OpenAPIStreamResumeFlow(
    req?: workflow.OpenAPIStreamResumeFlowRequest,
    options?: T,
  ): Promise<workflow.OpenAPIStreamRunFlowResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/v1/workflow/stream_resume');
    const method = 'POST';
    const data = {
      event_id: _req['event_id'],
      interrupt_type: _req['interrupt_type'],
      resume_data: _req['resume_data'],
      ext: _req['ext'],
      workflow_id: _req['workflow_id'],
      connector_id: _req['connector_id'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/batch_delete */
  BatchDeleteWorkflow(
    req: workflow.BatchDeleteWorkflowRequest,
    options?: T,
  ): Promise<workflow.BatchDeleteWorkflowResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/batch_delete');
    const method = 'POST';
    const data = {
      workflow_id_list: _req['workflow_id_list'],
      space_id: _req['space_id'],
      action: _req['action'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/get_trace */
  GetTraceSDK(
    req?: trace.GetTraceSDKRequest,
    options?: T,
  ): Promise<trace.GetTraceSDKResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/get_trace');
    const method = 'POST';
    const data = { Base: _req['Base'] };
    const params = {
      log_id: _req['log_id'],
      start_at: _req['start_at'],
      end_at: _req['end_at'],
      workflow_id: _req['workflow_id'],
      execute_id: _req['execute_id'],
    };
    return this.request({ url, method, data, params }, options);
  }

  /**
   * POST /api/workflow_api/list_spans
   *
   * Trace
   *
   * 列出历史执行的trace
   */
  ListRootSpans(
    req: trace.ListRootSpansRequest,
    options?: T,
  ): Promise<trace.ListRootSpansResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/list_spans');
    const method = 'POST';
    const data = {
      start_at: _req['start_at'],
      end_at: _req['end_at'],
      limit: _req['limit'],
      desc_by_start_time: _req['desc_by_start_time'],
      offset: _req['offset'],
      workflow_id: _req['workflow_id'],
      input: _req['input'],
      status: _req['status'],
      execute_mode: _req['execute_mode'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/workflow_detail */
  GetWorkflowDetail(
    req?: workflow.GetWorkflowDetailRequest,
    options?: T,
  ): Promise<workflow.GetWorkflowDetailResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/workflow_detail');
    const method = 'POST';
    const data = {
      workflow_ids: _req['workflow_ids'],
      space_id: _req['space_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/llm_fc_setting_detail
   *
   * 获取llm节点使用的技能的详情信息
   */
  GetLLMNodeFCSettingDetail(
    req: workflow.GetLLMNodeFCSettingDetailRequest,
    options?: T,
  ): Promise<workflow.GetLLMNodeFCSettingDetailResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/llm_fc_setting_detail');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      plugin_list: _req['plugin_list'],
      workflow_list: _req['workflow_list'],
      dataset_list: _req['dataset_list'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/llm_fc_setting_merged */
  GetLLMNodeFCSettingsMerged(
    req: workflow.GetLLMNodeFCSettingsMergedRequest,
    options?: T,
  ): Promise<workflow.GetLLMNodeFCSettingsMergedResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/llm_fc_setting_merged');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      plugin_fc_setting: _req['plugin_fc_setting'],
      workflow_fc_setting: _req['workflow_fc_setting'],
      dataset_fc_setting: _req['dataset_fc_setting'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/example_workflow_list
   *
   * 获取示例流程列表
   */
  GetExampleWorkFlowList(
    req?: workflow.GetExampleWorkFlowListRequest,
    options?: T,
  ): Promise<workflow.GetExampleWorkFlowListResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/example_workflow_list');
    const method = 'POST';
    const data = {
      page: _req['page'],
      size: _req['size'],
      name: _req['name'],
      flow_mode: _req['flow_mode'],
      checker: _req['checker'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** GET /v1/workflows/:workflow_id */
  OpenAPIGetWorkflowInfo(
    req?: workflow.OpenAPIGetWorkflowInfoRequest,
    options?: T,
  ): Promise<workflow.OpenAPIGetWorkflowInfoResponse> {
    const _req = req || {};
    const url = this.genBaseURL(`/v1/workflows/${_req['workflow_id']}`);
    const method = 'GET';
    const params = {
      connector_id: _req['connector_id'],
      is_debug: _req['is_debug'],
      caller: _req['caller'],
      include_chatflow_role: _req['include_chatflow_role'],
      include_input_output: _req['include_input_output'],
    };
    return this.request({ url, method, params }, options);
  }

  /**
   * GET /api/workflow_api/get_node_execute_history
   *
   * 查询节点执行历史。在试运行时需要轮询查看一个节点的执行历史。
   */
  GetNodeExecuteHistory(
    req: workflow.GetNodeExecuteHistoryRequest,
    options?: T,
  ): Promise<workflow.GetNodeExecuteHistoryResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/get_node_execute_history');
    const method = 'GET';
    const params = {
      workflow_id: _req['workflow_id'],
      space_id: _req['space_id'],
      execute_id: _req['execute_id'],
      node_id: _req['node_id'],
      is_batch: _req['is_batch'],
      batch_index: _req['batch_index'],
      node_type: _req['node_type'],
      node_history_scene: _req['node_history_scene'],
      Base: _req['Base'],
    };
    return this.request({ url, method, params }, options);
  }

  /**
   * POST /api/workflow_api/workflow_detail_info
   *
   * 获取子流程的详情
   */
  GetWorkflowDetailInfo(
    req?: workflow.GetWorkflowDetailInfoRequest,
    options?: T,
  ): Promise<workflow.GetWorkflowDetailInfoResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/workflow_detail_info');
    const method = 'POST';
    const data = {
      workflow_filter_list: _req['workflow_filter_list'],
      space_id: _req['space_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/list_publish_workflow
   *
   * publis
   */
  ListPublishWorkflow(
    req: workflow.ListPublishWorkflowRequest,
    options?: T,
  ): Promise<workflow.ListPublishWorkflowResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/list_publish_workflow');
    const method = 'POST';
    const data = {
      space_id: _req['space_id'],
      owner_id: _req['owner_id'],
      name: _req['name'],
      order_last_publish_time: _req['order_last_publish_time'],
      order_total_token: _req['order_total_token'],
      size: _req['size'],
      cursor_id: _req['cursor_id'],
      workflow_ids: _req['workflow_ids'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/workflow_api/validate_tree */
  ValidateTree(
    req: workflow.ValidateTreeRequest,
    options?: T,
  ): Promise<workflow.ValidateTreeResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/workflow_api/validate_tree');
    const method = 'POST';
    const data = {
      workflow_id: _req['workflow_id'],
      bind_project_id: _req['bind_project_id'],
      bind_bot_id: _req['bind_bot_id'],
      schema: _req['schema'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/workflow_api/node_panel_search
   *
   * 节点面板搜索
   */
  NodePanelSearch(
    req?: workflow.NodePanelSearchRequest,
    options?: T,
  ): Promise<workflow.NodePanelSearchResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/workflow_api/node_panel_search');
    const method = 'POST';
    const data = {
      search_type: _req['search_type'],
      space_id: _req['space_id'],
      project_id: _req['project_id'],
      search_key: _req['search_key'],
      page_or_cursor: _req['page_or_cursor'],
      page_size: _req['page_size'],
      exclude_workflow_id: _req['exclude_workflow_id'],
      enterprise_id: _req['enterprise_id'],
      Base: _req['Base'],
    };
    return this.request({ url, method, data }, options);
  }

  /** GET /v1/workflows/:workflow_id/run_histories/:execute_id/execute_nodes/:node_execute_uuid */
  OpenAPIGetNodeExecuteHistory(
    req?: workflow.OpenAPIGetNodeExecuteHistoryRequest,
    options?: T,
  ): Promise<workflow.OpenAPIGetNodeExecuteHistoryResponse> {
    const _req = req || {};
    const url = this.genBaseURL(
      `/v1/workflows/${_req['workflow_id']}/run_histories/${_req['execute_id']}/execute_nodes/${_req['node_execute_uuid']}`,
    );
    const method = 'GET';
    return this.request({ url, method }, options);
  }

  /** GET /v1/workflows/:workflow_id/versions */
  OpenAPIListVersion(
    req: workflow.OpenAPIListVersionRequest,
    options?: T,
  ): Promise<workflow.OpenAPIListVersionResponse> {
    const _req = req;
    const url = this.genBaseURL(
      `/v1/workflows/${_req['workflow_id']}/versions`,
    );
    const method = 'GET';
    const params = {
      publish_status: _req['publish_status'],
      page_size: _req['page_size'],
      page_token: _req['page_token'],
      include_input_output: _req['include_input_output'],
    };
    return this.request({ url, method, params }, options);
  }

  /** POST /v1/workflows/:workflow_id/collaboration_mode */
  OpenAPIToggleCollaborationMode(
    req?: workflow.OpenAPIToggleCollaborationModeRequest,
    options?: T,
  ): Promise<workflow.OpenAPIToggleCollaborationModeResponse> {
    const _req = req || {};
    const url = this.genBaseURL(
      `/v1/workflows/${_req['workflow_id']}/collaboration_mode`,
    );
    const method = 'POST';
    const data = { collaboration_mode: _req['collaboration_mode'] };
    return this.request({ url, method, data }, options);
  }
}
/* eslint-enable */
