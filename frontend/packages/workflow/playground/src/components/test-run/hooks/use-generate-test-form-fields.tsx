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

/* eslint-disable max-lines-per-function */

import { get, pick } from 'lodash-es';
import { FlowNodeFormData } from '@flowgram-adapter/free-layout-editor';

import { getSortedInputParameters, WorkflowNodeData } from '@coze-workflow/nodes';
import { StandardNodeType } from '@coze-workflow/base/types';

import { I18n } from '@coze-arch/i18n';


import { OptionType } from '@/constants/question-settings';


import {
  type GenerateFn,
  type GenerateAsyncFn,
} from '../utils/generate-test-form-fields';
import { generateField, startNodeOldType2VariableType, generateArrayInputParameters, generateObjectInputParameters } from '../utils/generate-test-form-fields';
import { generateInputJsonSchema } from '../utils/generate-input-json-schema';

import {
  generateDatabaseCreateTestRunFormFields,
  generateDatabaseDeleteTestRunFormFields,
  generateDatabaseQueryTestRunFormFields,
  generateDatabaseUpdateTestRunFormFields,
} from '../utils/generate-crud-nodes-test-run-form-fields';

import { COMMON_FIELD } from '../constants';


type GenerateTestFormFieldsMap = {
  [key in StandardNodeType]?: GenerateFn | GenerateAsyncFn;
} & {
  default: GenerateFn;
  batch: GenerateFn;
  setting: GenerateFn;
};

export const useGenerateTestFormFieldsMap = () => {


  const generateTestFormFieldsMap: GenerateTestFormFieldsMap = {
    [StandardNodeType.Start]: (formData, context) => {
      const inputParameters = formData?.outputs || [];


      const { node } = context;

      const inputFields = inputParameters.filter(input => {
        if (input.isPreset && !input.enabled) {
          return false;
        }
        return true;
      });

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const sortedInputFields = getSortedInputParameters(inputFields) as any;

      return (
        sortedInputFields
          .map(input => {
              const workflowVariable =
                node.context.variableService.getWorkflowVariableByKeyPath(
                  [node.id, input.name],
                  context,
                );
              if (!workflowVariable) {
                return null;
              }

              const jsonSchema = generateInputJsonSchema(
                workflowVariable.dtoMeta,
              );

              return {
                name: input.name,
                title: input.name,
                // Default values for variable configurations are written directly into the schema, rather than injected after the form is initialized
                initialValue: input.defaultValue,
                ...generateField(
                  startNodeOldType2VariableType(input.type),
                  input.required,
                  input.description,
                  jsonSchema,
                ),
              };
          })
          .filter(Boolean)
      );
    },


    [StandardNodeType.VariableAssign]: (formData, context) => {
      const inputParameters =
        formData?.$$input_decorator$$?.inputParameters?.map(e => ({
          input: e?.right,
          name: e?.left?.content?.keyPath?.[1],
        }));

      return generateArrayInputParameters(inputParameters, context);
    },
    [StandardNodeType.SubWorkflow]: (formData, context) => {
      const inputDefs = formData?.inputs?.inputDefs;
      if (!inputDefs || !Array.isArray(inputDefs)) {
        return [];
      }


      const inputData = formData?.inputs?.inputParameters || {};
      const inputParameters = inputDefs
        .map(i => ({
          input: inputData[i.name],
          name: i.name,
          required: i.required,
        }));
      return generateArrayInputParameters(inputParameters, context);
    },
    [StandardNodeType.Api]: (formData, context) => {
      const { node } = context;

      const nodeData = node.getData<WorkflowNodeData>(WorkflowNodeData);

      const detail = nodeData.getNodeData<StandardNodeType.Api>();

      const inputs = detail?.inputs;
      if (!inputs || !Array.isArray(inputs)) {
        return [];
      }
      const inputData = formData?.inputs?.inputParameters || {};
      const inputParameters = inputs.map(i => ({
        input: inputData[i.name],
        name: i.name,
        required: i.required,
      }));
      return generateArrayInputParameters(inputParameters, context);
    },


    [StandardNodeType.Question]: (formData, context) => {
      const inputFields = generateArrayInputParameters(
        formData?.inputParameters,
        context,
      );
      const anwserType = get(formData, 'questionParams.answer_type');
      const optionType = get(formData, 'questionParams.option_type');
      if (anwserType !== 'option' || optionType !== OptionType.Dynamic) {
        return inputFields;
      }
      const dynamicOption = get(formData, 'questionParams.dynamic_option');
      const dynamicField = generateObjectInputParameters(
        {
          dynamic_option: dynamicOption,
        },
        context,
      );
      return [...inputFields, ...dynamicField];
    },

    [StandardNodeType.Loop]: (formData, context) => [
      ...generateArrayInputParameters(
        formData?.inputs?.inputParameters,
        context,
      ),
      ...generateArrayInputParameters(
        formData?.inputs?.variableParameters,
        context,
      ),
    ],
    [StandardNodeType.Batch]: (formData, context) => [
      ...generateArrayInputParameters(
        formData?.inputs?.inputParameters,
        context,
      ),
    ],
    [StandardNodeType.LLM]: (formData, context) =>
      generateArrayInputParameters(
        formData?.$$input_decorator$$?.inputParameters,
        context,
      ),

    [StandardNodeType.Intent]: (formData, context) =>
      generateArrayInputParameters(formData?.inputs.inputParameters, context),




    [StandardNodeType.Database]: (formData, context) =>
      generateArrayInputParameters(formData?.inputParameters, context),



    [StandardNodeType.DatabaseCreate]: generateDatabaseCreateTestRunFormFields,
    [StandardNodeType.DatabaseDelete]: generateDatabaseDeleteTestRunFormFields,
    [StandardNodeType.DatabaseQuery]: generateDatabaseQueryTestRunFormFields,
    [StandardNodeType.DatabaseUpdate]: generateDatabaseUpdateTestRunFormFields,

    default: (formData, context) =>
      generateArrayInputParameters(formData?.inputParameters, context),

    // default: (formData, context) => {
    //   const { node } = context;
    //   const registry = node.getNodeRegistry();
    //   //Part of the path needs to be processed, for example, some parts start with a backslash and need to be removed
    //   const inputParametersPath =
    //     registry?.meta?.inputParametersPath ?? 'inputParameters';

    //   return generateArrayInputParameters(
    //     get(formData, inputParametersPath),
    //     context,
    //   );
    // },
    batch: (_, context) => {
      const { node } = context;
      const path = node.getNodeMeta()?.batchPath;
      if (!path) {
        return [];
      }

      const batchModePath = [
        StandardNodeType.SubWorkflow,
        StandardNodeType.Api,
      ].includes(node.flowNodeType as StandardNodeType)
        ? '/inputs/batchMode'
        : '/batchMode';
      const batchMode = node
        .getData(FlowNodeFormData)
        .formModel.getFormItemValueByPath(batchModePath);
      if (batchMode !== 'batch') {
        return [];
      }
      const batchData = node
        .getData(FlowNodeFormData)
        .formModel.getFormItemValueByPath(path);
      return generateArrayInputParameters(batchData?.inputLists, context).map(
        i => {
          /** Batch needs to be tagged. */
          if (i.decorator.type === 'FormItem') {
            const decorator = {
              ...i?.decorator,
              ...COMMON_FIELD.decorator,
              props: {
                ...i?.decorator?.props,
                tag: [
                  i?.decorator?.props?.tag,
                  I18n.t('workflow_detail_node_batch'),
                ]
                  .filter(Boolean)
                  .join(' - '),
              },
            };
            return {
              ...i,
              decorator,
            };
          }
          return i;
        },
      );
    },
    setting: (_, context) => {
      const { node } = context;
      if (node.flowNodeType === StandardNodeType.Batch) {
        const batchData = node
          .getData(FlowNodeFormData)
          .formModel.getFormItemValueByPath('/inputs');
        return generateObjectInputParameters(
          pick(batchData, ['batchSize', 'concurrentSize']),
          context,
        );
      } else if (node.flowNodeType === StandardNodeType.Loop) {
        const data = node
          .getData(FlowNodeFormData)
          .formModel.getFormItemValueByPath('/inputs');
        if (data.loopType === 'count') {
          return generateObjectInputParameters(
            pick(data, 'loopCount'),
            context,
          );
        }
      }
      return [];
    },
  };
  return generateTestFormFieldsMap;
};
