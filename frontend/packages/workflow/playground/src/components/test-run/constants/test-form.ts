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

/*******************************************************************************
 * Test form correlation constants
 */

import { FILE_TYPES, ViewVariableType } from '@coze-workflow/base';


import { I18n } from '@coze-arch/i18n';


import { getAccept } from '@/hooks/use-upload';


import { genFileTypeByViewVarType } from '../utils/common';
import { type TestFormField } from '../types';

/** The agreed key for testset */
export { TESTSET_BOT_NAME } from '@coze-workflow/test-run/constants';

/** Engage the dialogue key for testset */
export const TESTSET_CHAT_NAME = '_WORKFLOW_VARIABLE_NODE_CHAT_ID';

/** Test form type */
export enum TestFormType {
  /** in full operation */
  Default,
  /** single node operation */
  Single,
}

/** Sources of data used for test runs */
export enum TestRunDataSource {
  Aigc = 'aigc',
  Testset = 'testset',
  User = 'user',
}

/** Fixed field name */
export enum FieldName {
  Node = '_node',
  Batch = '_batch',
  Input = '_input',
  Bot = '_bot',
  Chat = '_chat',
  Datasets = '_datasets',
  DatasetsIs = '_datasets_is',
  DatasetsName = '_datasets_name',
  DatasetsDescription = '_datasets_description',
  Setting = '_setting',
  JSON = '_json',
}

/** Basic template for fields */
export const DEFAULT_FIELD_TEMPLATE = {
  required: true,
  decorator: {
    type: 'FormItem',
  },
};

/** Node basic template */
export const NODE_FIELD_TEMPLATE = {
  type: 'FormObject',
  name: FieldName.Node,
  component: {
    type: 'FieldGroup',
  },
};

export const SETTING_FIELD_TEMPLATE = {
  type: 'FormObject',
  name: FieldName.Setting,
  component: {
    type: 'div',
  },
};

/** batch template */
export const BATCH_FIELD_TEMPLATE = {
  type: 'FormObject',
  name: FieldName.Batch,
  component: {
    type: 'div',
  },
};

/** input template */
export const INPUT_FIELD_TEMPLATE = {
  type: 'FormObject',
  name: FieldName.Input,
  component: {
    type: 'div',
  },
};

/** The template for the entire JSON input */
export const INPUT_JSON_FIELD_TEMPLATE = {
  type: 'FormObject',
  name: FieldName.JSON,
};

/** Save template to datasets */
export const DATASETS_FIELD_TEMPLATE = {
  type: 'FormObject',
  name: FieldName.Datasets,
  component: {
    type: 'div',
  },
  children: [
    {
      type: 'FormString',
      name: FieldName.DatasetsIs,
      component: {
        type: 'Checkbox',
        props: {
          title: I18n.t('workflow_debug_data_save'),
        },
      },
      decorator: {
        type: 'FormItem',
      },
    },
    // {
    //   type: 'FormString',
    //   name: FieldName.DatasetsName,
    //   title: I18n.t('workflow_testset_name'),
    //   component: {
    //     type: 'TestsetNameInput',
    //     props: {
    //       maxLength: 50,
    //     },
    //   },
    //   hidden: `{{ !$values?.${FieldName.Datasets}?.${FieldName.DatasetsIs} }}`,
    //   required: true,
    //   decorator: {
    //     type: 'FormItem',
    //   },
    //   validator: [
    //     {
    //       required: true,
    //       triggerType: 'onBlur',
    //       message: I18n.t('workflow_testset_required_tip', {
    //         param_name: '{{$self.title}}',
    //       }),
    //     },
    //   ],
    // },
    // {
    //   type: 'FormString',
    //   name: FieldName.DatasetsDescription,
    //   title: I18n.t('workflow_testset_desc'),
    //   component: {
    //     type: 'Input',
    //     props: {
    //       maxCount: 200,
    //       maxLength: 200,
    //     },
    //   },
    //   hidden: `{{ !$values?.${FieldName.Datasets}?.${FieldName.DatasetsIs} }}`,
    //   decorator: {
    //     type: 'FormItem',
    //   },
    // },
  ],
};

/**
 * Some common fields
 */
export const COMMON_FIELD = {
  required: true,
  decorator: {
    type: 'FormItem',
  },
  validator: [
    {
      required: true,
      triggerType: 'onBlur',
      message: I18n.t('workflow_testset_required_tip', {
        param_name: '{{$self.title}}',
      }),
    },
  ],
};

type FileTypes =
  | ViewVariableType.Image
  | ViewVariableType.File
  | ViewVariableType.File
  | ViewVariableType.Image
  | ViewVariableType.Doc
  | ViewVariableType.Code
  | ViewVariableType.Ppt
  | ViewVariableType.Txt
  | ViewVariableType.Excel
  | ViewVariableType.Audio
  | ViewVariableType.Zip
  | ViewVariableType.Video
  | ViewVariableType.Svg
  | ViewVariableType.Voice
  | ViewVariableType.ArrayImage
  | ViewVariableType.ArrayFile
  | ViewVariableType.ArrayDoc
  | ViewVariableType.ArrayCode
  | ViewVariableType.ArrayPpt
  | ViewVariableType.ArrayTxt
  | ViewVariableType.ArrayExcel
  | ViewVariableType.ArrayAudio
  | ViewVariableType.ArrayZip
  | ViewVariableType.ArrayVideo
  | ViewVariableType.ArraySvg
  | ViewVariableType.ArrayVoice;

const FILE_TYPE_FIELD_MAP: Record<FileTypes, TestFormField> = FILE_TYPES.reduce(
  (map, type) => {
    const multiple = ViewVariableType.isArrayType(type);

    return {
      ...map,
      [type]: {
        type: 'FormString',
        component: {
          type: 'File',
          props: {
            multiple,
            accept: getAccept(type),
            inputType: type,
            fileType: genFileTypeByViewVarType(type),
            enableInputURL: true,
          },
        },
      },
    };
  },
  {},
) as Record<FileTypes, TestFormField>;

/** Mapping of types to fields */
export const NOT_FILE_TYPE_FIELD_MAP: Record<
  Exclude<ViewVariableType, FileTypes>,
  TestFormField
> = {
  [ViewVariableType.String]: {
    type: 'FormString',
    component: {
      type: 'Input',
    },
  },
  [ViewVariableType.Integer]: {
    type: 'FormNumber',
    component: {
      type: 'InputInteger',
    },
  },
  [ViewVariableType.Boolean]: {
    type: 'FormBoolean',
    initialValue: true,
    component: {
      type: 'Switch',
    },
  },
  [ViewVariableType.Number]: {
    type: 'FormNumber',
    component: {
      type: 'InputNumber',
    },
  },
  [ViewVariableType.Time]: {
    type: 'FormString',
    component: {
      type: 'InputTime',
    },
  },
  [ViewVariableType.Object]: {
    type: 'FormString',
    component: {
      type: 'JsonEditor',
    },
    initialValue: '{}',
  },
  [ViewVariableType.ArrayString]: {
    type: 'FormString',
    component: {
      type: 'JsonEditor',
    },
    initialValue: '[]',
  },
  [ViewVariableType.ArrayInteger]: {
    type: 'FormString',
    component: {
      type: 'JsonEditor',
    },
    initialValue: '[]',
  },
  [ViewVariableType.ArrayBoolean]: {
    type: 'FormString',
    component: {
      type: 'JsonEditor',
    },
    initialValue: '[]',
  },
  [ViewVariableType.ArrayNumber]: {
    type: 'FormString',
    component: {
      type: 'JsonEditor',
    },
    initialValue: '[]',
  },
  [ViewVariableType.ArrayObject]: {
    type: 'FormString',
    component: {
      type: 'JsonEditor',
    },
    initialValue: '[]',
  },
  [ViewVariableType.ArrayTime]: {
    type: 'FormString',
    component: {
      type: 'JsonEditor',
    },
    initialValue: '[]',
  },
};

export const TYPE_FIELD_MAP: Record<ViewVariableType, TestFormField> = {
  ...NOT_FILE_TYPE_FIELD_MAP,
  ...FILE_TYPE_FIELD_MAP,
};
