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

import { useMemo } from 'react';


import { FlowNodeFormData } from '@flowgram-adapter/free-layout-editor';

import { type WorkflowNodeEntity } from '@flowgram-adapter/free-layout-editor';

import { type NodeData, WorkflowNodeData } from '@coze-workflow/nodes';
import { StandardNodeType } from '@coze-workflow/base';

import { I18n } from '@coze-arch/i18n';




import { type TestFormField } from '../types';
import { BATCH_FIELD_TEMPLATE, SETTING_FIELD_TEMPLATE, INPUT_FIELD_TEMPLATE, NODE_FIELD_TEMPLATE, TestFormType } from '../constants';

import { useTestRunReporterService } from '../../../hooks';

import { useTestFormInitialValueV2 } from './use-test-form-initial-value-v2';


import { useGetStartNode } from './use-get-start-node';
import { useGenerateTestFormFieldsMap } from './use-generate-test-form-fields';

interface OptionsType {
  onCreateTestSet?: () => void;
  hideGroupLabel?: boolean;
}

/**
 * Derive testFormSchema from node formData
 * The entire schema is currently constructed as follows:
 * - node []//There will be multiple nodes in the future, so nodes need to be divided into dimensionsdes in the future, so nodes need to be divided into dimensions
 *   //batch and input will be in two groups, because the filed.name of batch and input may be duplicatedame of batch and input may be duplicated
 *   - batch
 *   - input
 *     - string
 *     - number
 *     - ...
 * - bot//he is independent only oneependent and there is only one
 * - testset
 *   - is
 *   - name
 *   - description
 */
export const useTestFormSchema = (
  node?: WorkflowNodeEntity,
  options?: OptionsType,
) => {
  const { generateInitialValues } = useTestFormInitialValueV2();
  const reporter = useTestRunReporterService();

  const { getNode } = useGetStartNode();

  const generateTestFormFieldsMap = useGenerateTestFormFieldsMap();


  /** If it is a start node, it means it is running in full, otherwise it is running on a single node. */
  const testFormType = useMemo(() => {
    if (!node || node.flowNodeType === StandardNodeType.Start) {
      return TestFormType.Default;
    }
    return TestFormType.Single;
  }, [node, node?.flowNodeType]);


  const generate = async () => {
    let current = node;
    if (!current) {
      const startNodeEntity = getNode();
      if (!startNodeEntity) {
        return null;
      }
      current = startNodeEntity;
    }
    const reporterKey = reporter.formSchemaGen.start();
    const nodeDataEntity = current.getData<WorkflowNodeData>(WorkflowNodeData);
    const nodeData = nodeDataEntity.getNodeData<keyof NodeData>();

    const formData = current
      .getData(FlowNodeFormData)
      .formModel.getFormItemValueByPath('/');
    const nodeTitle = formData?.nodeMeta?.title;

    /** Part of the computing node shell */
    const nodeField: TestFormField = {
      ...NODE_FIELD_TEMPLATE,
      title: I18n.t('workflow_debug_testonenode_group', {
        nodeTitle,
      }),
      component: {
        ...NODE_FIELD_TEMPLATE.component,
        props: {
          icon: nodeData?.icon || formData?.nodeMeta?.icon,
          hideGroupLabel: options?.hideGroupLabel,
        },
      },
      children: [],
    };

    /** Part Compute Node Batch Part */
    const batchFields = generateTestFormFieldsMap.batch(formData, {
      node: current,
    });
    if (batchFields.length) {
      nodeField.children.push({
        ...BATCH_FIELD_TEMPLATE,
        children: batchFields,
      });
    }

    /** Part of computing non-input content */
    const settingFields = generateTestFormFieldsMap.setting(formData, {
      node: current,
    });
    if (settingFields.length) {
      nodeField.children.push({
        ...SETTING_FIELD_TEMPLATE,
        children: settingFields,
      });
    }

    /** Parts Compute Node Input */
    const generateFn =
      generateTestFormFieldsMap[current.flowNodeType as StandardNodeType] ||
      generateTestFormFieldsMap.default;

    const inputFields = await generateFn(formData, {
      node: current,
    });
    if (inputFields.length) {
      nodeField.children.push({
        ...INPUT_FIELD_TEMPLATE,
        children: inputFields,
      });
    }

    const fields: TestFormField[] = [];

    if (nodeField.children.length) {
      fields.push(nodeField);
    }


    const schema = {
      id: current.id,
      type: testFormType,
      fields,
    };

    await generateInitialValues(schema);

    reporter.formSchemaGen.end(reporterKey, {
      node_type: current.flowNodeType as string,
    });

    return schema;
  };

  return { generate };
};
