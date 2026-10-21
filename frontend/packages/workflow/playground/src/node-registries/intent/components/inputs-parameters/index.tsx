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



import { type InputValueVO, useNodeTestId } from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';

import { ValueExpressionInputField } from '@/node-registries/common/fields';

import {
  ColumnTitles,
  FieldArrayItem,
  FieldArrayList,
  Section,
  useFieldArray,
  withFieldArray,
} from '@/form';

import { COLUMNS, INPUT_PATH } from '../../constants';


const InputsParametersField = withFieldArray(() => {
  const { name: fieldName, value } = useFieldArray<InputValueVO>();
  const safeValue = value || [];
  const { getNodeSetterId } = useNodeTestId();

  return (
    <Section
      title={I18n.t('workflow_detail_node_parameter_input')}
      tooltip={I18n.t('workflow_intent_input_tooltips')}
      testId={getNodeSetterId(fieldName)}
      actions={
        []
      }
    >

      <ColumnTitles columns={COLUMNS} />

      <FieldArrayList>
        {safeValue?.map(({ name }, index) => (
          <FieldArrayItem hiddenRemove>
            <ValueExpressionInputField
              key={index}
              label={name}
              required
              name={`${fieldName}.${index}.input`}
            />
          </FieldArrayItem>
        ))}
      </FieldArrayList>
    </Section>
  );
});

export default function InputsParameters() {
  return (
    <InputsParametersField
      name={INPUT_PATH}
      defaultValue={[{ name: 'query' }]}
    />
  );
}
