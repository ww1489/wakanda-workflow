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

import type { ReactNode } from 'react';

import { type FrontWorkflowInfo } from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';
import { Table } from '@coze-arch/coze-design';

import { resource_resource_common, type ResourceInfo } from '@coze-arch/bot-api/plugin_develop';

import { transformResourceToWorkflowEditInfo } from './utils';


import { useDeleteAction } from './use-delete-action';
import { useCopyAction } from './use-copy-action';

import {
  type WorkflowResourceActionProps,
  type WorkflowResourceActionReturn,
} from './type';

const { ActionKey } = resource_resource_common;



export const useWorkflowResourceMenuActions = (
  props: WorkflowResourceActionProps & {
    userId?: string;
    onEditWorkflowInfo: (partialWorkflowInfo: FrontWorkflowInfo) => void;
  },
): Pick<WorkflowResourceActionReturn, 'renderWorkflowResourceActions'> & {
  modals: ReactNode[];
} => {
  const { onEditWorkflowInfo, getCommonActions } = props;
  const { actionHandler: deleteAction, deleteModal } = useDeleteAction(props);
  const { actionHandler: copyAction } = useCopyAction(props);
  const renderWorkflowResourceActions = (record: ResourceInfo): ReactNode => {
    const { actions } = record;
    const deletion = actions?.find(action => action.key === ActionKey.Delete);
    const copy = actions?.find(action => action.key === ActionKey.Copy);
    const edit = actions?.find(action => action.key === ActionKey.Edit);
    return (
      <Table.TableAction
        deleteProps={{ hide: !deletion, disabled: deletion?.enable === false, disableConfirm: true, handler: () => deleteAction(record) }}
        copyProps={{ hide: !copy, disabled: copy?.enable === false, handler: () => copyAction(record) }}
        actionList={[
          { hide: !edit, disabled: edit?.enable === false, actionKey: 'edit', actionText: I18n.t('Edit'), handler: () => onEditWorkflowInfo(transformResourceToWorkflowEditInfo(record) as FrontWorkflowInfo) },
          ...(getCommonActions?.(record) ?? []),
        ]}
      />
    );
  };
  return { renderWorkflowResourceActions, modals: [deleteModal] };
};
