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

import { createBrowserRouter, Navigate } from 'react-router-dom';

import { SpaceSubModuleEnum } from '@coze-foundation/space-ui-adapter';
import { GlobalError } from '@coze-foundation/layout';
import { BaseEnum } from '@coze-arch/web-context';

import { Layout } from '../layout';
import {
  LoginPage,
  SpaceLayout,
  SpaceIdLayout,
  Redirect,
  spaceSubMenu,
  WorkflowPage,
  Library,
  PluginLayout,
  PluginToolPage,
  PluginPage,
  DatabaseDetail,
  OAuthConsentConfirmPage,
} from './async-components';

export const router: ReturnType<typeof createBrowserRouter> =
  createBrowserRouter([
    // Document routing
    {
      path: '/open/docs/*',
      Component: Redirect,
      loader: () => ({
        hasSider: false,
        requireAuth: false,
      }),
    },
    {
      path: '/docs/*',
      Component: Redirect,
      loader: () => ({
        hasSider: false,
        requireAuth: false,
      }),
    },
    {
      path: '/information/auth/success',
      Component: Redirect,
      loader: () => ({
        hasSider: false,
        requireAuth: false,
      }),
    },
    // main application route
    {
      path: '/',
      Component: Layout,
      errorElement: <GlobalError />,
      children: [
        {
          index: true,
          element: <Navigate to="/space" replace />,
        },
        // login page routing
        {
          path: 'sign',
          Component: LoginPage,
          errorElement: <GlobalError />,
          loader: () => ({
            hasSider: false,
            requireAuth: false,
          }),
        },

        // OAuth consent confirm page
        {
          path: 'oauth/confirm',
          Component: OAuthConsentConfirmPage,
          errorElement: <GlobalError />,
          loader: () => ({
            hasSider: false,
            requireAuth: false,
          }),
        },

        // Workspace Routing
        {
          path: 'space',
          Component: SpaceLayout,
          loader: () => ({
            hasSider: true,
            requireAuth: true,
            subMenu: spaceSubMenu,
            menuKey: BaseEnum.Space,
          }),
          children: [
            {
              path: ':space_id',
              Component: SpaceIdLayout,
              children: [
                {
                  index: true,
                  element: <Navigate to="library" replace />,
                },

                {
                  path: 'develop',
                  element: <Navigate to="../library" replace />,
                },

                // resource library
                {
                  path: 'library',
                  Component: Library,
                  loader: () => ({
                    subMenuKey: SpaceSubModuleEnum.LIBRARY,
                  }),
                },

                // database resources
                {
                  path: 'database',
                  children: [
                    {
                      path: ':table_id',
                      element: <DatabaseDetail />,
                    },
                  ],
                  loader: () => ({
                    showMobileTips: true,
                    pageModeByQuery: true,
                  }),
                },

                // plugin resources
                {
                  path: 'plugin/:plugin_id',
                  Component: PluginLayout,
                  children: [
                    {
                      index: true,
                      Component: PluginPage,
                    },
                    {
                      path: 'tool/:tool_id',
                      children: [
                        {
                          index: true,
                          Component: PluginToolPage,
                        },
                      ],
                    },
                  ],
                },
              ],
            },
          ],
        },

        // workflow routing
        {
          path: 'work_flow',
          Component: WorkflowPage,
          loader: () => ({
            hasSider: false,
            requireAuth: true,
          }),
        },

        {
          path: '*',
          element: <Navigate to="/space" replace />,
        },
      ],
    },
  ]);
