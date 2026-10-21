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

import { type FC } from 'react';



import { type ExpressionEditorTreeNode } from '@coze-workflow/components';
import { useNodeTestId } from '@coze-workflow/base';



import { InputSlotWidget } from '@coze-common/editor-plugins/input-slot';
import {
  Completion,
  Validation,
  HighlightExpressionOnActive,
} from '@coze-common/editor-plugins/expression';
import { InsertInputSlotAction } from '@coze-common/editor-plugins/actions';
import { ActionBar } from '@coze-common/editor-plugins/action-bar';
import { I18n } from '@coze-arch/i18n';
import { IconCozCrossFill } from '@coze-arch/coze-design/icons';
import { IconButton } from '@coze-arch/coze-design';


import { useNodeFormPanelState } from '@/hooks/use-node-side-sheet-store';



import { ExpandEditorContainer } from '@/components/editor-container';

import { useVariableTree } from '../expression-editor/hooks';
import { type ExpressionEditorProps } from '../expression-editor';

import { SystemPromptEditor } from './system-prompt-editor';

import styles from './index.module.less';
export const ExpandSheetEditor: FC<ExpressionEditorProps> = props => {
  const { getNodeSetterId } = useNodeTestId();
  const dataTestID = getNodeSetterId(props.name);
  const variableTree: ExpressionEditorTreeNode[] = useVariableTree();

  const { setFullscreenPanel } = useNodeFormPanelState();

  return (
    <>
      <ExpandEditorContainer
        id={'system-prompt-expand-editor'}
        closeButton={
          <div className="flex items-center gap-[8px]">
            <IconButton
              onClick={() => setFullscreenPanel(null)}
              color="secondary"
              icon={<IconCozCrossFill color="#060709CC" />}
            />
          </div>
        }
        containerClassName="bg-[#FCFCFF]"
        headerClassName="bg-[#FCFCFF] text-[16px] leading-[22px] font-medium mb-[12px] !p-[20px]"
        editorTitle={
          <div className="flex items-center gap-[8px]">
            <div
              className="text-[16px] leading-[22px] font-medium"
              style={{
                color: 'var(--Fg-COZ-fg-plus, rgba(6, 7, 9, 0.96))',
              }}
            >
              {I18n.t('workflow_LLM_node_sp_title')}
            </div>
          </div>
        }
        contentClassName={styles['editor-content']}
        editorContent={
          <>
            <div className="relative w-full h-full pl-[8px] pr-[20px]">
              <SystemPromptEditor
                defaultValue={props?.value}
                onChange={props?.onChange}
                readonly={props?.readonly}
                placeholder={props?.placeholder as React.ReactNode}
                wrapperClassName={styles['prompt-editor']}
                wrapperStyle={{
                  height: '100%',
                }}
                dataTestID={dataTestID}
              />

              <InputSlotWidget mode="input" />
              <HighlightExpressionOnActive />
              <Completion variableTree={variableTree} />
              <Validation variableTree={variableTree} />

              {!props?.readonly && (
                <>
                  <ActionBar>
                    <InsertInputSlotAction />
                  </ActionBar>
                </>
              )}
            </div>
          </>
        }
      />
    </>
  );
};
