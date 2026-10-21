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

import { type FC, type PropsWithChildren } from 'react';

import { GlobalLayout } from '@coze-foundation/layout';
import { RequireAuthContainer } from '@coze-foundation/account-ui-adapter';
import { I18n } from '@coze-arch/i18n';
import {
  IconCozWorkspace,
  IconCozWorkspaceFill,
  IconCozDocument,
} from '@coze-arch/coze-design/icons';
import { useRouteConfig } from '@coze-arch/bot-hooks';

import { useHasSider } from './hooks/use-has-sider';
import { AccountDropdown } from '../account-dropdown';

export const GlobalLayoutComposed: FC<PropsWithChildren> = ({ children }) => {
  const config = useRouteConfig();
  const hasSider = useHasSider();
  return (
    <RequireAuthContainer
      needLogin={!!config.requireAuth}
      loginOptional={!!config.requireAuthOptional}
    >
      <GlobalLayout
        hasSider={hasSider}
        banner={null}
        actions={[]}
        menus={[
          {
            title: I18n.t('navigation_workspace'),
            icon: <IconCozWorkspace />,
            activeIcon: <IconCozWorkspaceFill />,
            path: '/space',
            dataTestId: 'layout_workspace-button',
          },
        ]}
        extras={[
          {
            icon: <IconCozDocument />,
            tooltip: I18n.t('menu_documents'),
            onClick: () => {
              // cp-disable-next-line
              window.open('https://www.coze.cn/open/docs/guides');
            },
            dataTestId: 'layout_document-button',
          },
        ]}
        footer={<AccountDropdown />}
      >
        {children}
      </GlobalLayout>
    </RequireAuthContainer>
  );
};
