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

import { forwardRef, useImperativeHandle } from 'react';

import classNames from 'classnames';
import { useInfiniteScroll } from 'ahooks';
import { useUserInfo } from '@coze-foundation/account-adapter';
import { I18n } from '@coze-arch/i18n';
import { Table, Search, Layout, Cascader, Space } from '@coze-arch/coze-design';
import { renderHtmlTitle } from '@coze-arch/bot-utils';
import { EVENT_NAMES, sendTeaEvent } from '@coze-arch/bot-tea';
import {
  type ResType,
  type ResourceInfo,
} from '@coze-arch/bot-api/plugin_develop';

import { highlightFilterStyle } from '@/constants/filter-style';
import { WorkspaceEmpty } from '@/components/workspace-empty';

import { type ListData, type BaseLibraryPageProps } from './types';
import { listResources } from './list-resources';
import { useGetColumns } from './hooks/use-columns';
import { useCachedQueryParams } from './hooks/use-cached-query-params';
import { eventLibraryType } from './consts';
import { LibraryHeader } from './components/library-header';

import s from './index.module.less';

export { useDatabaseConfig } from './hooks/use-entity-configs/use-database-config';
export { usePluginConfig } from './hooks/use-entity-configs/use-plugin-config';
export { useWorkflowConfig } from './hooks/use-entity-configs/use-workflow-config';
export { type LibraryEntityConfig } from './types';
export { type UseEntityConfigHook } from './hooks/use-entity-configs/types';
export { BaseLibraryItem } from './components/base-library-item';

export const BaseLibraryPage = forwardRef<
  { reloadList: () => void },
  BaseLibraryPageProps
>(
  // eslint-disable-next-line @coze-arch/max-line-per-function
  ({ spaceId, isPersonalSpace = true, entityConfigs }, ref) => {
    const userInfo = useUserInfo();
    const { params, setParams, resetParams, hasFilter, ready } =
      useCachedQueryParams({
        spaceId,
      });

    const listResp = useInfiniteScroll<ListData>(
      async prev => {
        if (!ready || !userInfo?.user_id_str) {
          return {
            list: [],
            nextCursorId: undefined,
            hasMore: false,
          };
        }
        return listResources(spaceId, userInfo.user_id_str, {
          ...params,
          cursor: prev?.nextCursorId,
        });
      },
      {
        reloadDeps: [params, spaceId, ready, userInfo?.user_id_str],
      },
    );

    useImperativeHandle(ref, () => ({
      reloadList: listResp.reload,
    }));

    const columns = useGetColumns({
      entityConfigs,
      reloadList: listResp.reload,
      isPersonalSpace,
    });

    const typeFilterData = [
      ...entityConfigs.map(item => item.typeFilter).filter(filter => !!filter),
    ];

    return (
      <Layout
        className={s['layout-content']}
        title={renderHtmlTitle(I18n.t('navigation_workspace_library'))}
      >
        <Layout.Header className={classNames(s['layout-header'], 'pb-0')}>
          <div className="w-full">
            <LibraryHeader entityConfigs={entityConfigs} />
            <div className="flex items-center justify-between">
              <Space>
                <Cascader
                  data-testid="workspace.library.filter.type"
                  className={s.cascader}
                  style={highlightFilterStyle}
                  dropdownClassName="[&_.semi-cascader-option-lists]:h-fit"
                  showClear={false}
                  value={params.res_type_filter}
                  treeData={typeFilterData}
                  onChange={v => {
                    const typeFilter = typeFilterData.find(
                      item =>
                        item.value === ((v as Array<number>)?.[0] as number),
                    );
                    sendTeaEvent(EVENT_NAMES.workspace_action_front, {
                      space_id: spaceId,
                      space_type: isPersonalSpace ? 'personal' : 'teamspace',
                      tab_name: 'library',
                      action: 'filter',
                      filter_type: 'types',
                      filter_name: typeFilter?.filterName ?? typeFilter?.label,
                    });

                    setParams(prev => ({
                      ...prev,
                      res_type_filter: v as Array<number>,
                    }));
                  }}
                />
              </Space>
              <Search
                data-testid="workspace.library.filter.name"
                className="!min-w-min"
                style={params.name ? highlightFilterStyle : {}}
                showClear={true}
                width={200}
                loading={listResp.loading}
                placeholder={I18n.t('workspace_library_search')}
                value={params.name}
                onSearch={v => {
                  sendTeaEvent(EVENT_NAMES.search_front, {
                    full_url: window.location.href,
                    source: 'library',
                    search_word: v,
                  });
                  setParams(prev => ({
                    ...prev,
                    name: v,
                  }));
                }}
              />
            </div>
          </div>
        </Layout.Header>
        <Layout.Content>
          <Table
            data-testid="workspace.library.table"
            offsetY={178}
            tableProps={{
              loading: listResp.loading,
              dataSource: listResp.data?.list,
              columns,
              // Click on the whole line
              onRow: (record?: ResourceInfo) => {
                if (
                  !record ||
                  record.res_type === undefined ||
                  record.detail_disable
                ) {
                  return {};
                }
                return {
                  onClick: () => {
                    sendTeaEvent(EVENT_NAMES.workspace_action_front, {
                      space_id: spaceId,
                      space_type: isPersonalSpace ? 'personal' : 'teamspace',
                      tab_name: 'library',
                      action: 'click',
                      id: record.res_id,
                      name: record.name,
                      type:
                        record.res_type && eventLibraryType[record.res_type],
                    });
                    entityConfigs
                      .find(c => c.target.includes(record.res_type as ResType))
                      ?.onItemClick(record);
                  },
                };
              },
            }}
            empty={
              <WorkspaceEmpty onClear={resetParams} hasFilter={hasFilter} />
            }
            enableLoad
            loadMode="cursor"
            strictDataSourceProp
            hasMore={listResp.data?.hasMore}
            onLoad={listResp.loadMore}
          />
        </Layout.Content>
      </Layout>
    );
  },
);
