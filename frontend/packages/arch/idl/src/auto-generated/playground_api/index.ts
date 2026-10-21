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



import * as open_api_playground from './namespaces/open_api_playground';
import * as playground from './namespaces/playground';

export * from './namespaces/base';
export * from './namespaces/bot_common';
export * from './namespaces/bot_open_api';
export * from './namespaces/bot_task_common';
export * from './namespaces/douyin_fenshen';
export * from './namespaces/frontier';
export * from './namespaces/long_term_memory';
export * from './namespaces/model';
export * from './namespaces/op';
export * from './namespaces/open_api_playground';
export * from './namespaces/playground';
export * from './namespaces/prompt_resource';
export * from './namespaces/punish_center';
export * from './namespaces/resource';
export * from './namespaces/resource_common';
export * from './namespaces/shortcut_command';
export * from './namespaces/task_common';
export * from './namespaces/user_delete_base';



export default class PlaygroundApiService<T> {
  private request: any = () => {
    throw new Error('PlaygroundApiService.request is undefined');
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

  /** POST /api/playground_api/space/list */
  GetSpaceListV2(
    req?: playground.GetSpaceListV2Request,
    options?: T,
  ): Promise<playground.GetSpaceListV2Response> {
    const _req = req || {};
    const url = this.genBaseURL('/api/playground_api/space/list');
    const method = 'POST';
    const data = {
      search_word: _req['search_word'],
      enterprise_id: _req['enterprise_id'],
      organization_id: _req['organization_id'],
      scope_type: _req['scope_type'],
      page: _req['page'],
      size: _req['size'],
    };
    return this.request({ url, method, data }, options);
  }

  /** POST /api/playground_api/get_imagex_url */
  GetImagexShortUrl(
    req?: playground.GetImagexShortUrlRequest,
    options?: T,
  ): Promise<playground.GetImagexShortUrlResponse> {
    const _req = req || {};
    const url = this.genBaseURL('/api/playground_api/get_imagex_url');
    const method = 'POST';
    const data = { uris: _req['uris'], scene: _req['scene'] };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /api/playground_api/mget_user_info
   *
   * --------------------------用户相关--------------------------------
   */
  MGetUserBasicInfo(
    req: playground.MGetUserBasicInfoRequest,
    options?: T,
  ): Promise<playground.MGetUserBasicInfoResponse> {
    const _req = req;
    const url = this.genBaseURL('/api/playground_api/mget_user_info');
    const method = 'POST';
    const data = {
      user_ids: _req['user_ids'],
      need_user_status: _req['need_user_status'],
      need_enterprise_identity: _req['need_enterprise_identity'],
      need_volcano_user_name: _req['need_volcano_user_name'],
    };
    return this.request({ url, method, data }, options);
  }

  /**
   * POST /v1/files/upload
   *
   * 会话链路能力开放相关
   *
   * File 相关 OpenAPI
   */
  UploadFileOpen(
    req: playground.UploadFileOpenRequest,
    options?: T,
  ): Promise<playground.UploadFileOpenResponse> {
    const _req = req;
    const url = this.genBaseURL('/v1/files/upload');
    const method = 'POST';
    const data = { Data: _req['Data'] };
    const headers = { 'Content-Type': _req['Content-Type'] };
    return this.request({ url, method, data, headers }, options);
  }

  /** GET /api/playground_api/private_policy/:secret */
  GetPolicyContent(
    req: playground.GetPolicyContentRequest,
    options?: T,
  ): Promise<playground.GetPolicyContentResponse> {
    const _req = req;
    const url = this.genBaseURL(
      `/api/playground_api/private_policy/${_req['secret']}`,
    );
    const method = 'GET';
    return this.request({ url, method }, options);
  }

  /** DELETE /v1/workspaces/:workspace_id/members */
  OpenRemoveSpaceMember(
    req?: open_api_playground.OpenRemoveSpaceMemberRequest,
    options?: T,
  ): Promise<open_api_playground.OpenRemoveSpaceMemberResponse> {
    const _req = req || {};
    const url = this.genBaseURL(
      `/v1/workspaces/${_req['workspace_id']}/members`,
    );
    const method = 'DELETE';
    const data = { user_ids: _req['user_ids'] };
    return this.request({ url, method, data }, options);
  }

  /** GET /v1/workspaces/:workspace_id/members */
  OpenSpaceMemberList(
    req?: open_api_playground.OpenSpaceMemberListRequest,
    options?: T,
  ): Promise<open_api_playground.OpenSpaceMemberListResponse> {
    const _req = req || {};
    const url = this.genBaseURL(
      `/v1/workspaces/${_req['workspace_id']}/members`,
    );
    const method = 'GET';
    const params = { page_num: _req['page_num'], page_size: _req['page_size'] };
    return this.request({ url, method, params }, options);
  }

  /** POST /v1/workspaces/:workspace_id/members */
  OpenAddSpaceMember(
    req?: open_api_playground.OpenAddSpaceMemberRequest,
    options?: T,
  ): Promise<open_api_playground.OpenAddSpaceMemberResponse> {
    const _req = req || {};
    const url = this.genBaseURL(
      `/v1/workspaces/${_req['workspace_id']}/members`,
    );
    const method = 'POST';
    const data = { users: _req['users'] };
    return this.request({ url, method, data }, options);
  }

  /** GET /v1/bots/:bot_id */
  OpenGetBotInfo(
    req?: open_api_playground.OpenGetBotInfoRequest,
    options?: T,
  ): Promise<open_api_playground.OpenGetBotInfoResponse> {
    const _req = req || {};
    const url = this.genBaseURL(`/v1/bots/${_req['bot_id']}`);
    const method = 'GET';
    const params = {
      is_published: _req['is_published'],
      connector_id: _req['connector_id'],
    };
    return this.request({ url, method, params }, options);
  }

  /** POST /v1/workspaces/:workspace_id/members/apply */
  OpenApplyJoinSpace(
    req?: open_api_playground.OpenApplyJoinSpaceRequest,
    options?: T,
  ): Promise<open_api_playground.OpenApplyJoinSpaceResponse> {
    const _req = req || {};
    const url = this.genBaseURL(
      `/v1/workspaces/${_req['workspace_id']}/members/apply`,
    );
    const method = 'POST';
    const data = { user_ids: _req['user_ids'] };
    return this.request({ url, method, data }, options);
  }

  /** DELETE /v1/workspaces/:workspace_id */
  OpenRemoveSpace(
    req?: open_api_playground.OpenRemoveSpaceRequest,
    options?: T,
  ): Promise<open_api_playground.OpenRemoveSpaceResponse> {
    const _req = req || {};
    const url = this.genBaseURL(`/v1/workspaces/${_req['workspace_id']}`);
    const method = 'DELETE';
    return this.request({ url, method }, options);
  }

  /** PUT /v1/workspaces/:workspace_id/members/:user_id */
  OpenUpdateSpaceMember(
    req?: open_api_playground.OpenUpdateSpaceMemberRequest,
    options?: T,
  ): Promise<open_api_playground.OpenUpdateSpaceMemberResponse> {
    const _req = req || {};
    const url = this.genBaseURL(
      `/v1/workspaces/${_req['workspace_id']}/members/${_req['user_id']}`,
    );
    const method = 'PUT';
    const data = { role_type: _req['role_type'] };
    return this.request({ url, method, data }, options);
  }

  /** GET /v1/bots/:bot_id/versions */
  OpenListBotVersions(
    req?: open_api_playground.OpenListBotVersionsRequest,
    options?: T,
  ): Promise<open_api_playground.OpenListBotVersionsResponse> {
    const _req = req || {};
    const url = this.genBaseURL(`/v1/bots/${_req['bot_id']}/versions`);
    const method = 'GET';
    const params = {
      page_num: _req['page_num'],
      page_size: _req['page_size'],
      publish_status: _req['publish_status'],
      connector_id: _req['connector_id'],
    };
    return this.request({ url, method, params }, options);
  }

  /** POST /v1/bots/:bot_id/collaboration_mode */
  OpenSwitchBotDevelopMode(
    req?: open_api_playground.OpenSwitchBotDevelopModeRequest,
    options?: T,
  ): Promise<open_api_playground.OpenSwitchBotDevelopModeResponse> {
    const _req = req || {};
    const url = this.genBaseURL(
      `/v1/bots/${_req['bot_id']}/collaboration_mode`,
    );
    const method = 'POST';
    const data = { collaboration_mode: _req['collaboration_mode'] };
    return this.request({ url, method, data }, options);
  }
}
/* eslint-enable */
