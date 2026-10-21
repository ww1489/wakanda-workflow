import { useEffect, useState } from 'react';

import { I18n } from '@coze-arch/i18n';
import { Button, Input, Modal, Spin, Toast } from '@coze-arch/coze-design';
import { useSpaceStore } from '@coze-arch/bot-studio-store';
import type {
  PluginApi,
  PluginInfoForPlayground,
} from '@coze-arch/bot-api/plugin_develop';
import type { PluginFrom } from '@coze-arch/bot-api/playground_api';
import { PluginDevelopApi } from '@coze-arch/bot-api';

export type WorkflowPluginTool = PluginApi & {
  plugin_icon?: string;
  version_name?: string;
  version_ts?: string;
  project_id?: string;
  plugin_from?: PluginFrom;
};

interface PluginSelectorProps {
  pluginApiList?: PluginApi[];
  onPluginApiListChange?: (list: PluginApi[]) => void;
  openModeCallback?: (tool: WorkflowPluginTool) => unknown;
}

const PAGE_SIZE = 20;

export function usePluginModalParts(props: PluginSelectorProps) {
  const spaceId = useSpaceStore(state => state.space.id);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [plugins, setPlugins] = useState<PluginInfoForPlayground[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [adding, setAdding] = useState<string>();

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    PluginDevelopApi.GetPlaygroundPluginList({
      space_id: spaceId,
      page,
      size: PAGE_SIZE,
      name: search || undefined,
    })
      .then(response => {
        if (cancelled) {
          return;
        }
        setPlugins(response.data?.plugin_list ?? []);
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

  const content = (
    <Spin spinning={loading}>
      <div className="flex flex-col gap-4 p-4 overflow-auto">
        {plugins.map(plugin => (
          <section key={plugin.id}>
            <h3>{plugin.name}</h3>
            <p>{plugin.desc_for_human}</p>
            {plugin.plugin_apis?.map(api => {
              const selected = props.pluginApiList?.some(
                item => item.api_id === api.api_id,
              );
              return (
                <div
                  key={api.api_id}
                  className="flex items-center justify-between gap-4 py-2"
                >
                  <span>
                    <strong>{api.name}</strong> {api.desc}
                  </span>
                  <Button
                    loading={adding === api.api_id}
                    onClick={async () => {
                      if (selected) {
                        props.onPluginApiListChange?.(
                          props.pluginApiList?.filter(
                            item => item.api_id !== api.api_id,
                          ) ?? [],
                        );
                        return;
                      }
                      setAdding(api.api_id);
                      try {
                        await props.openModeCallback?.({
                          ...api,
                          plugin_id: plugin.id,
                          plugin_name: plugin.name,
                          plugin_icon: plugin.plugin_icon,
                          version_name: plugin.version_name,
                          version_ts: plugin.version_ts,
                        });
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
          </section>
        ))}
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
  );
  return {
    content,
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
  };
}

function PluginSelectorModal(
  props: PluginSelectorProps & { onClose: () => void },
) {
  const { content, filter } = usePluginModalParts(props);
  return (
    <Modal
      visible
      title={I18n.t('bot_edit_plugin_select_title')}
      onCancel={props.onClose}
      footer={null}
      width={800}
    >
      {filter}
      {content}
    </Modal>
  );
}

export function usePluginApisModal(
  props: PluginSelectorProps & { closeCallback?: () => void },
) {
  const [visible, setVisible] = useState(false);
  const close = () => {
    setVisible(false);
    props.closeCallback?.();
  };
  return {
    open: () => setVisible(true),
    close,
    node: visible ? <PluginSelectorModal {...props} onClose={close} /> : null,
  };
}
