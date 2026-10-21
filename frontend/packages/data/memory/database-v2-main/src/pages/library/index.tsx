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

import { useNavigate } from 'react-router-dom';


import { useDataCallbacks } from '@coze-data/database-context';
import { type DatabaseTabs } from '@coze-data/database-v2-base/types';




import { DatabaseDetail } from '../../components/database-detail';

interface IProps {
  version?: string;
  botId: string;
  spaceId?: string;
  databaseId: string;
  enterFrom: string;
  needHideCloseIcon?: boolean;
  initialTab?: DatabaseTabs;
}

export const DatabaseInner = ({
  version,
  botId,
  spaceId,
  databaseId,
  enterFrom,
  needHideCloseIcon,
  initialTab,
}: IProps) => {
  const { onStatusChange, onUpdateDisplayName } = useDataCallbacks();
  const navigate = useNavigate();
  const handleClose = () => {
    if (window.history.length === 1) {
      navigate(`/space/${spaceId}/library`);
    }
    navigate(-1);
  };

  return (
    <DatabaseDetail
      version={version}
      enterFrom={enterFrom}
      databaseId={databaseId}
      initialTab={initialTab}
      onIDECallback={{
        onStatusChange,
        onUpdateDisplayName,
      }}
      needHideCloseIcon={needHideCloseIcon}
      onClose={handleClose}
    />
  );
};
