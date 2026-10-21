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

import { beforeEach, describe, expect, it, vi } from 'vitest';
import { WorkflowMode } from '@coze-arch/bot-api/workflow_api';
import { ActionKey, ResType } from '@coze-arch/bot-api/plugin_develop';
import { TableType } from '@coze-arch/bot-api/memory';
import { MemoryApi, PluginDevelopApi, workflowApi } from '@coze-arch/bot-api';

import { listResources } from './list-resources';

vi.mock('@coze-arch/bot-api', () => ({
  MemoryApi: { ListDatabase: vi.fn() },
  PluginDevelopApi: { GetDevPluginList: vi.fn() },
  workflowApi: { GetWorkFlowList: vi.fn() },
}));
vi.mock('./consts', () => ({ LIBRARY_PAGE_SIZE: 15 }));

describe('standalone resource lists', () => {
  beforeEach(() => vi.clearAllMocks());

  it('lists ordinary workflows without binding to a project', async () => {
    vi.mocked(workflowApi.GetWorkFlowList).mockResolvedValue({
      code: 0,
      msg: '',
      BaseResp: {},
      data: {
        total: 31,
        workflow_list: [
          {
            workflow_id: '7',
            name: 'HTTP job',
            creator: { id: '42', self: true },
          },
        ],
      },
    });

    const result = await listResources('3', '42', {
      res_type_filter: [ResType.Workflow],
      name: 'HTTP',
      cursor: '2',
    });

    expect(workflowApi.GetWorkFlowList).toHaveBeenCalledWith({
      space_id: '3',
      page: 2,
      size: 15,
      name: 'HTTP',
      flow_mode: WorkflowMode.Workflow,
    });
    expect(result).toMatchObject({ hasMore: true, nextCursorId: '3' });
    expect(result.list[0]).toMatchObject({
      res_id: '7',
      actions: expect.arrayContaining([
        { key: ActionKey.Delete, enable: true },
      ]),
    });
    expect(PluginDevelopApi.GetDevPluginList).not.toHaveBeenCalled();
    expect(MemoryApi.ListDatabase).not.toHaveBeenCalled();
  });

  it('includes draft plugins in the management list', async () => {
    vi.mocked(PluginDevelopApi.GetDevPluginList).mockResolvedValue({
      plugin_list: [{ id: '9', creator: { id: '42' } }],
      total: '1',
    });

    const result = await listResources('3', '42', {
      res_type_filter: [ResType.Plugin],
    });

    expect(PluginDevelopApi.GetDevPluginList).toHaveBeenCalledWith(
      expect.objectContaining({
        space_id: '3',
        dev_id: '42',
        project_id: '0',
        page: 1,
      }),
    );
    expect(result).toMatchObject({
      hasMore: false,
      list: [{ res_id: '9', res_type: ResType.Plugin }],
    });
  });

  it('uses row offsets for standalone databases and preserves ownership', async () => {
    vi.mocked(MemoryApi.ListDatabase).mockResolvedValue({
      BaseResp: {},
      has_more: false,
      database_info_list: [
        { id: '10', creator_id: '99', update_time: '1700000000000' },
      ],
    });

    const result = await listResources('3', '42', {
      res_type_filter: [ResType.Database],
      name: 'orders',
      cursor: '3',
    });

    expect(MemoryApi.ListDatabase).toHaveBeenCalledWith(
      expect.objectContaining({
        space_id: '3',
        project_id: '0',
        table_type: TableType.OnlineTable,
        table_name: 'orders',
        offset: 30,
        limit: 15,
      }),
    );
    expect(result.list[0]).toMatchObject({
      edit_time: '1700000000',
      actions: [{ key: ActionKey.Delete, enable: false }],
    });
  });
});
