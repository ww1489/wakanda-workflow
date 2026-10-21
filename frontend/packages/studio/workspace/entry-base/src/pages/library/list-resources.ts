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

import {
  WorkflowMode,
  WorkFlowDevStatus,
} from '@coze-arch/bot-api/workflow_api';
import {
  ActionKey,
  OrderBy,
  PublishStatus,
  ResType,
  type ResourceInfo,
} from '@coze-arch/bot-api/plugin_develop';
import { TableType, SortDirection } from '@coze-arch/bot-api/memory';
import { MemoryApi, PluginDevelopApi, workflowApi } from '@coze-arch/bot-api';

import { type ListData } from './types';
import { LIBRARY_PAGE_SIZE, type QueryParams } from './consts';

const MILLISECONDS_PER_SECOND = 1000;

// Keep the existing resource cards and actions, backed by the three SQL lists.
export const listResources = async (
  spaceId: string,
  userId: string,
  params: QueryParams,
): Promise<ListData> => {
  const page = Number(params.cursor || 1);
  const size = LIBRARY_PAGE_SIZE;
  const type = params.res_type_filter?.[0] ?? ResType.Workflow;
  let list: ResourceInfo[];
  let hasMore: boolean;

  if (type === ResType.Plugin) {
    const response = await PluginDevelopApi.GetDevPluginList({
      space_id: spaceId,
      dev_id: userId,
      project_id: '0',
      page,
      size,
      name: params.name,
      order_by: OrderBy.UpdateTime,
    });
    list = (response.plugin_list ?? []).map(plugin => ({
      res_id: plugin.id,
      res_type: ResType.Plugin,
      res_sub_type: plugin.plugin_type,
      name: plugin.name,
      desc: plugin.desc_for_human,
      icon: plugin.plugin_icon,
      creator_id: plugin.creator?.id,
      creator_name: plugin.creator?.name,
      creator_avatar: plugin.creator?.avatar_url,
      edit_time: plugin.update_time,
      space_id: spaceId,
      actions: [
        { key: ActionKey.Delete, enable: plugin.creator?.id === userId },
      ],
    }));
    hasMore = page * size < Number(response.total ?? 0);
  } else if (type === ResType.Database) {
    const response = await MemoryApi.ListDatabase({
      space_id: spaceId,
      project_id: '0',
      table_type: TableType.OnlineTable,
      table_name: params.name,
      offset: (page - 1) * size,
      limit: size,
      order_by: [{ field: 'update_time', direction: SortDirection.Desc }],
    });
    list = (response.database_info_list ?? []).map(database => ({
      res_id: database.id,
      res_type: ResType.Database,
      name: database.table_name,
      desc: database.table_desc,
      icon: database.icon_url,
      creator_id: database.creator_id,
      edit_time: String(
        Math.floor(Number(database.update_time ?? 0) / MILLISECONDS_PER_SECOND),
      ),
      space_id: spaceId,
      actions: [
        { key: ActionKey.Delete, enable: database.creator_id === userId },
      ],
    }));
    hasMore = !!response.has_more;
  } else {
    const response = await workflowApi.GetWorkFlowList({
      space_id: spaceId,
      page,
      size,
      name: params.name,
      flow_mode: WorkflowMode.Workflow,
    });
    list = (response.data.workflow_list ?? []).map(workflow => ({
      res_id: workflow.workflow_id,
      res_type: ResType.Workflow,
      res_sub_type: WorkflowMode.Workflow,
      name: workflow.name,
      desc: workflow.desc,
      icon: workflow.url,
      creator_id: workflow.creator?.id,
      creator_name: workflow.creator?.name,
      creator_avatar: workflow.creator?.avatar_url,
      edit_time: workflow.update_time,
      space_id: spaceId,
      publish_status:
        workflow.status === WorkFlowDevStatus.HadSubmit
          ? PublishStatus.Published
          : PublishStatus.UnPublished,
      biz_extend: {
        icon_uri: workflow.icon_uri ?? '',
        url: workflow.url ?? '',
        schema_type: String(workflow.schema_type ?? 0),
      },
      actions: [
        { key: ActionKey.Edit, enable: workflow.creator?.id === userId },
        { key: ActionKey.Delete, enable: workflow.creator?.id === userId },
        { key: ActionKey.Copy, enable: true },
      ],
    }));
    hasMore = page * size < Number(response.data.total ?? 0);
  }

  return { list, hasMore, nextCursorId: String(page + 1) };
};
