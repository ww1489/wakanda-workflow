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
  ResType,
  type LibraryResourceListRequest,
} from '@coze-arch/bot-api/plugin_develop';

export const LIBRARY_PAGE_SIZE = 15;

export type QueryParams = Pick<
  LibraryResourceListRequest,
  'cursor' | 'res_type_filter' | 'name'
>;

export const initialParam: QueryParams = {
  cursor: '',
  res_type_filter: [ResType.Workflow],
  name: '',
};

/** event type */
export const eventLibraryType = {
  [ResType.Plugin]: 'plugin',
  [ResType.Workflow]: 'workflow',
  [ResType.Imageflow]: 'imageflow',
  [ResType.Knowledge]: 'knowledge',
  [ResType.UI]: 'ui',
  [ResType.Prompt]: 'prompt',
  [ResType.Database]: 'database',
  [ResType.Variable]: 'variable',
  [ResType.Voice]: 'voice',
} as const;
