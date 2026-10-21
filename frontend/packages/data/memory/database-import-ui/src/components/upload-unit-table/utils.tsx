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
  EntityStatus,
  UploadStatus,
} from '@coze-data/database-import-core';


import { IconUploadCSV, IconUploadXLS } from '@coze-arch/bot-icons';


import { ProcessStatus } from '../../types';

export const getTypeIcon = ({ type }: { type?: string }) =>
  type === 'csv' ? <IconUploadCSV /> : <IconUploadXLS />;

export const getProcessStatus = (
  status: UploadStatus | EntityStatus,
) => {
  if (
    status === UploadStatus.UPLOADING ||
    status === UploadStatus.VALIDATING ||
    status === UploadStatus.WAIT ||
    status === EntityStatus.EntityStatusProcess
  ) {
    return ProcessStatus.Processing;
  }

  if (
    status === UploadStatus.SUCCESS ||
    status === EntityStatus.EntityStatusSuccess
  ) {
    return ProcessStatus.Complete;
  }

  if (
    status === UploadStatus.VALIDATE_FAIL ||
    status === UploadStatus.UPLOAD_FAIL ||
    status === EntityStatus.EntityStatusFail
  ) {
    return ProcessStatus.Failed;
  }
  return ProcessStatus.Processing;
};
