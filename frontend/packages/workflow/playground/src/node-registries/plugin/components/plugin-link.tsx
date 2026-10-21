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


import { type ApiNodeIdentifier } from '@coze-workflow/nodes';
import { I18n } from '@coze-arch/i18n';
import { IconCozExit } from '@coze-arch/coze-design/icons';
import { usePluginNodeServiceStore } from '../hooks';

export const PluginLink = ({
  identifier,
}: {
  identifier: ApiNodeIdentifier;
}) => {
  const getApiNodeDetail = usePluginNodeServiceStore(state => state.getData);
  const detail = getApiNodeDetail(identifier);
  if (!detail?.spaceID || !detail.pluginID) {
    return null;
  }
  return (
    <a
      className="cursor-pointer flex items-center w-full justify-between"
      href={`/space/${detail.spaceID}/plugin/${detail.pluginID}`}
      target="_blank"
      rel="noreferrer"
      onClick={e => e.stopPropagation()}
    >
      {I18n.t('mkpl_plugin_detail', {}, 'Plugin Detail')}
      <IconCozExit className="text-xs" />
    </a>
  );
};

export const createPluginLink = (identifier: ApiNodeIdentifier) => (
  <PluginLink identifier={identifier} />
);
