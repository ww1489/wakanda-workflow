import { useEffect, useState } from 'react';

import {
  workflowApi,
  WorkflowMode,
  type ReleasedWorkflow,
} from '@coze-workflow/base/api';
import { I18n } from '@coze-arch/i18n';
import { Button, Input, Spin, Toast } from '@coze-arch/coze-design';
import { useSpaceStore } from '@coze-arch/bot-studio-store';

import type { BotPluginWorkFlowItem, WorkFlowModalModeProps } from '../type';

const PAGE_SIZE = 20;

export const useWorkflowModalParts = (props: WorkFlowModalModeProps) => {
  const spaceId = useSpaceStore(state => state.space.id) ?? '';
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [workflows, setWorkflows] = useState<ReleasedWorkflow[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [adding, setAdding] = useState<string>();

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    workflowApi
      .GetReleasedWorkflows({
        space_id: spaceId,
        page,
        size: PAGE_SIZE,
        name: search || undefined,
        flow_mode: WorkflowMode.Workflow,
      })
      .then(response => {
        if (cancelled) {
          return;
        }
        setWorkflows(response.data?.workflow_list ?? []);
        setTotal(Number(response.data?.total ?? 0));
      })
      .catch(error => {
        if (!cancelled) {
          Toast.error(String(error));
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [spaceId, page, search]);

  return {
    sider: null,
    filter: (
      <Input
        value={search}
        placeholder={I18n.t('Search')}
        onChange={value => {
          setSearch(value);
          setPage(1);
        }}
      />
    ),
    content: (
      <Spin spinning={loading}>
        <div className="flex flex-col gap-4 p-4 overflow-auto">
          {workflows
            .filter(
              workflow =>
                !props.excludedWorkflowIds?.includes(
                  workflow.workflow_id ?? '',
                ),
            )
            .map(workflow => {
              const item: BotPluginWorkFlowItem = {
                workflow_id: workflow.workflow_id ?? '',
                plugin_id: workflow.plugin_id ?? '',
                name: workflow.name ?? '',
                desc: workflow.desc ?? '',
                plugin_icon: workflow.icon ?? '',
                parameters: [],
                version_name: workflow.flow_version || workflow.version,
                flow_mode: WorkflowMode.Workflow,
              };
              const selected = props.workFlowList?.some(
                value => value.workflow_id === item.workflow_id,
              );
              return (
                <div
                  key={item.workflow_id}
                  className="flex items-center justify-between gap-4 py-2"
                >
                  <span>
                    <strong>{item.name}</strong> {item.desc}
                  </span>
                  <Button
                    loading={adding === item.workflow_id}
                    onClick={async () => {
                      setAdding(item.workflow_id);
                      try {
                        if (selected) {
                          props.onRemove?.(item);
                        } else {
                          await props.onAdd?.(item, { isDup: false, spaceId });
                        }
                      } catch (error) {
                        Toast.error(String(error));
                      } finally {
                        setAdding(undefined);
                      }
                    }}
                  >
                    {I18n.t(selected ? 'Remove' : 'add')}
                  </Button>
                </div>
              );
            })}
          <div className="flex justify-between">
            <Button disabled={page === 1} onClick={() => setPage(page - 1)}>
              上一页
            </Button>
            <span>{page}</span>
            <Button
              disabled={page * PAGE_SIZE >= total}
              onClick={() => setPage(page + 1)}
            >
              下一页
            </Button>
          </div>
        </div>
      </Spin>
    ),
  };
};
