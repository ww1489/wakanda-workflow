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

import { type FC, useState, useRef } from 'react';

import classNames from 'classnames';


import { type ExpressionEditorTreeNode } from '@coze-workflow/components';
import { useNodeTestId } from '@coze-workflow/base';




import { InputSlotWidget } from '@coze-common/editor-plugins/input-slot';
import {
  Completion,
  Validation,
  HighlightExpressionOnActive,
} from '@coze-common/editor-plugins/expression';
import { I18n } from '@coze-arch/i18n';



import { FormCard } from '@/form-extensions/components/form-card';

import { useVariableTree } from '../expression-editor/hooks';
import { type ExpressionEditorProps } from '../expression-editor';

import { SystemPromptEditor } from './system-prompt-editor';
import { PromptKitBar } from './prompt-kit-bar';

import s from './index.module.less';

export const EditorWithPromptKit: FC<ExpressionEditorProps> = props => {
  const {
    name,
    readonly,
    placeholder,
    value,
    onChange,
  } = props;
  const { getNodeSetterId } = useNodeTestId();
  const dataTestID = getNodeSetterId(name);

  const [isFocused, setIsFocused] = useState(false);
  const variableTree: ExpressionEditorTreeNode[] = useVariableTree();

  const promptValue = useRef<string>(value || '');

  const handleOnFocus = () => {
    setIsFocused(true);
  };

  const handleOnBlur = () => {
    setIsFocused(false);
  };

  return (
    <div className={s.container} onBlur={() => onChange?.(promptValue.current)}>
      <FormCard
        defaultExpand
        showBottomBorder
        noPadding
        header={I18n.t('workflow_LLM_node_sp_title')}
        tooltip={I18n.t('workflow_detail_llm_sys_prompt_content_tips')}
        contentClassName={classNames(s['code-content'], {
          [s['code-content-focus']]: isFocused,
        })}
        actionButton={
          <PromptKitBar
            {...props}
            onInsertPrompt={prompt => props?.onChange(prompt)}
          />
        }
      >
        <SystemPromptEditor
          defaultValue={value}
          value={value}
          onChange={v => (promptValue.current = v)}
          readonly={readonly}
          placeholder={placeholder as React.ReactNode}
          isControled
          dataTestID={dataTestID}
          onFocus={handleOnFocus}
          onBlur={handleOnBlur}
        />
        <InputSlotWidget mode="input" />
        <HighlightExpressionOnActive />

        <Completion variableTree={variableTree} />
        <Validation variableTree={variableTree} />
      </FormCard>
    </div>
  );
};
