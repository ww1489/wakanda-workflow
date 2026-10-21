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

import { useNavigate, useParams } from 'react-router-dom';


import { I18n } from '@coze-arch/i18n';
import { IconCozArrowLeft } from '@coze-arch/coze-design/icons';
import { Breadcrumb as SemiBreadcrumb, Button } from '@coze-arch/coze-design';
import { type DynamicParams } from '@coze-arch/bot-typings/teamspace';
import { useSpaceStore } from '@coze-arch/bot-studio-store';
import { type BreadcrumbProps as SemiBreadcrumbProps } from '@coze-arch/bot-semi/Breadcrumb';



import { type PluginMetaInfo, type PluginAPIInfo } from '@coze-arch/bot-api/developer_api';
import { type MockSet } from '@coze-arch/bot-api/debugger_api';
import { useSpaceApp } from '@coze-foundation/space-store';

import s from './index.module.less';

export interface BreadCrumbProps extends SemiBreadcrumbProps {
  pluginInfo?: PluginMetaInfo;
  pluginToolInfo?: PluginAPIInfo;
  isPublish?: boolean;
  mockSetInfo?: MockSet;
}

// eslint-disable-next-line @coze-arch/max-line-per-function
export const UIBreadcrumb: React.FC<BreadCrumbProps> = ({
  pluginInfo,
  pluginToolInfo,
  mockSetInfo,
  isPublish,
  ...props
}) => {
  const spaceApp = useSpaceApp();
  const id = useSpaceStore(store => store.space.id);
  const navigate = useNavigate();
  const params = useParams<DynamicParams>();

  const goBack = () => {
    navigate(`/space/${id}/library`);
  };
  const goBackToPluginIdList = () => {
    navigate(`/space/${id}/${spaceApp}/${params.plugin_id}`);
  };


  const goBackToToolIdList = () => {
    navigate(
      `/space/${id}/${spaceApp}/${params.plugin_id}/tool/${params.tool_id}?mode=preview`,
    );
  };

  const goBackToMockSetList = () => {
    navigate(
      `/space/${id}/${spaceApp}/${params.plugin_id}/tool/${params.tool_id}/plugin-mock-set`,
    );
  };

  const renderBreadcrumbItemForPlugin = () => {
    let breadCrumbList: React.ReactNode[] = [];
    let onBackClick: () => void;

    if (pluginInfo?.name) {
      onBackClick = goBack;
    }

    if (pluginToolInfo?.name) {
      onBackClick = goBackToPluginIdList;
    }

    if (mockSetInfo) {
      onBackClick = goBackToToolIdList;

      if (mockSetInfo?.name) {
        onBackClick = goBackToMockSetList;
      }
    }
    breadCrumbList = [
      <Button
        color="secondary"
        icon={<IconCozArrowLeft />}
        onClick={() => {
          onBackClick();
        }}
      >
        {I18n.t('library_resource_detail_back')}
      </Button>,
    ];
    return breadCrumbList;
  };

  return (
    <div className={s['bot-breadcrumb']}>
      <SemiBreadcrumb {...props}>{renderBreadcrumbItemForPlugin()}</SemiBreadcrumb>
    </div>
  );
};
