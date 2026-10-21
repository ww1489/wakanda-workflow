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

import { type BotTableRWMode, type FieldItem } from '@coze-arch/bot-api/memory';

export interface TableMemoryItem extends FieldItem {
  errorMapper?: Record<string, string[]>;
  disableMustRequired?: boolean;
  nanoid?: string;
  /**
   * Is it a built-in field?
   * @Description Built-in field: For display only, users cannot modify it, and fields with the same name cannot be created
   */
  isSystemField?: boolean;
}

export interface DatabaseInfo {
  tableId: string;
  name: string;
  desc: string;
  icon_uri?: string;
  extra_info?: Record<string, string>;
  readAndWriteMode: BotTableRWMode;
  tableMemoryList: TableMemoryItem[];
}
export type DatabaseList = DatabaseInfo[];
