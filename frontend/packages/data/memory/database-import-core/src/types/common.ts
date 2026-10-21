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

import type { EntityStatus, UploadStatus } from '../constants';

interface CommonUnitItem {
  key?: string;
  percent: number | undefined;
  status: UploadStatus | EntityStatus;
  type: string; // Unit type, url, doc, pdf, etc
  uri: string; // The relative path of the file on tos url | web_url
  url: string; // Absolute URL of the file
  name: string; // unit name｜webpage title
  validateMessage?: string;
  // External incoming dynamic errorMessage
  dynamicErrorMessage?: string;
}

interface DocUnitItem {
  size?: string; // file size
  docId?: string; // Document ID, used when updating content
  fileInstance?: File | undefined;
  uid?: string;
}
export type UnitItem = CommonUnitItem & DocUnitItem & { statusDescript?: string };
