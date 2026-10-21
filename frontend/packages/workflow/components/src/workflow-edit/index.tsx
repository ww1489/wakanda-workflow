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

/* eslint-disable @coze-arch/max-line-per-function */
/* eslint-disable max-lines-per-function -- refactor later */

import { useMemo, useRef, useState } from 'react';

import { isFunction } from 'lodash-es';
import { type WorkflowJSON } from '@flowgram-adapter/free-layout-editor';
import { type FrontWorkflowInfo } from '@coze-workflow/base/types';
import {
  type BindBizType,
  type SchemaType,
  workflowApi,
  WorkflowMode,
} from '@coze-workflow/base/api';
import {
  WORKFLOW_NAME_MAX_LEN,
  WORKFLOW_NAME_REGEX,
} from '@coze-workflow/base';
import { PictureUpload } from '@coze-common/biz-components/picture-upload';
import { type UploadValue } from '@coze-common/biz-components';
import { REPORT_EVENTS } from '@coze-arch/report-events';
import { reporter } from '@coze-arch/logger';
import { I18n } from '@coze-arch/i18n';
import { useUserInfo } from '@coze-arch/foundation-sdk';

import { Button, LoadingButton, Space, Toast } from '@coze-arch/coze-design';
import { useSpaceStore } from '@coze-arch/bot-studio-store';
import { Form, UIFormTextArea, UIModal } from '@coze-arch/bot-semi';
import { CustomError } from '@coze-arch/bot-error';
import { FileBizType, IconType } from '@coze-arch/bot-api/developer_api';

import s from './index.module.less';
/** Enter error code for compliance exception */
const sensitiveWordsErrorCode = ['702095075', '702095081'];


export interface RuleItem {
  validator: (rules: unknown[], value: string) => boolean | Error;
}

interface EditWorkFlowPropsInner {
  /** process type */
  flowMode?: WorkflowMode;
  mode: 'update' | 'add';
  visible: boolean;
  // Default confirmed disabled

  /** Custom pop-up title */
  customTitleRender?: (title: React.ReactNode) => React.ReactNode;

  initConfirmDisabled?: boolean;
  workFlow?: FrontWorkflowInfo;
  onSuccess?: (val: {
    workflowId?: string;
    flowMode?: EditWorkFlowPropsInner['flowMode'];
  }) => void;
  onCancel?: () => void;
  /** @deprecated unused */
  spaceID?: string;
  getLatestWorkflowJson?: () => Promise<WorkflowJSON>;
  bindBizId?: string;
  bindBizType?: BindBizType;
  /** The current project id, only the workflow within the project has this field */
  projectId?: string;
  nameValidators?: RuleItem[];
}

/** form value */
interface FormValue {
  icon_uri: UploadValue;
  name: string;
  target: string;
  schema_type: SchemaType;
}

/** Get pop-up title */
function getModalTitle(
  mode: EditWorkFlowPropsInner['mode'],
  flowMode: EditWorkFlowPropsInner['flowMode'],
): string {
  switch (flowMode) {

    case WorkflowMode.Workflow:
      return mode === 'add'
        ? I18n.t('workflow_list_create_modal_title')
        : I18n.t('workflow_list_edit_modal_title');

    default:
      return mode === 'add'
        ? I18n.t('workflow_list_create_modal_title')
        : I18n.t('workflow_list_edit_modal_title');
  }
}

const getPictureUploadInitValue = (
  workFlow?: FrontWorkflowInfo,
): UploadValue | undefined => {
  if (!workFlow) {
    return;
  }
  return [
    {
      url: workFlow.url || '',
      uid: workFlow.icon_uri || '',
    },
  ];
};

export function CreateWorkflowModal({
  flowMode = WorkflowMode.Workflow,
  mode,
  bindBizId,
  bindBizType,
  projectId,
  visible,
  workFlow,
  initConfirmDisabled = false,
  customTitleRender,
  onSuccess,
  onCancel,
  nameValidators = [],
}: EditWorkFlowPropsInner) {
  const formRef = useRef<Form<Partial<FormValue>>>(null);
  const [confirmDisabled, setConfirmDisabled] = useState(initConfirmDisabled);
  const [sensitiveTip, setSensitiveTip] = useState<string | undefined>();
  const userInfo = useUserInfo();
  const currentLocale = userInfo?.locale ?? navigator.language ?? 'en-US';

  const getValues = async () => {
    const formApi = formRef.current?.formApi;
    await formApi?.validate(['name']);
    return formApi?.getValues() as Partial<FormValue> | undefined;
  };

  const handleCancel = () => {
    reporter.info({
      message: 'workflow_info_modal_cancel',
      namespace: 'workflow',
    });
    setSensitiveTip(undefined);
    onCancel?.();
  };

  const handleError = (error: Error & { code?: string; msg?: string }) => {
    if (sensitiveWordsErrorCode.includes(error.code || '')) {
      setConfirmDisabled(true);
      setSensitiveTip(error.msg);
      return;
    }

    handleCancel();
    throw error;
  };

  const handleUpdateWorkflow = async () => {
    const workflowId = workFlow?.workflow_id;

    reporter.info({
      message: 'workflow_info_modal_confirm_update',
      namespace: 'workflow',
      meta: {
        workflowId,
      },
    });

    if (!workflowId) {
      const msg = I18n.t('workflow_list_create_modal_workflow_id_empty');
      throw new CustomError(REPORT_EVENTS.parmasValidation, msg);
    }
    const values = await getValues();
    const updateParams = {
      workflow_id: workflowId,
      icon_uri: values?.icon_uri?.[0].uid || '',
      name: values?.name,
      desc: values?.target ? values.target : '',
      space_id: workFlow.space_id || '',
      // There is no need to re-test run to update the avatar and other information.
      ignore_status_transfer: true,
      schema_type: values?.schema_type || workFlow?.schema_type,
    };

    try {
      await workflowApi.UpdateWorkflowMeta(updateParams);

      reporter.info({
        message: 'workflow_info_modal_update_success',
        namespace: 'workflow',
      });
      Toast.success({
        content: I18n.t('workflow_list_update_success'),
        showClose: false,
      });
      await onSuccess?.({
        workflowId,
        flowMode,
      });
    } catch (error) {
      reporter.error({
        message: 'workflow_info_modal_update_fail',
        namespace: 'workflow',
        error,
      });
      handleError(error);
    }
  };

  const handleCreateWorkflow = async () => {
    reporter.info({
      message: 'workflow_info_modal_confirm_create',
      namespace: 'workflow',
    });
    const values = await getValues();
    try {
      const reqParams = {
        ...values,
        space_id: useSpaceStore.getState().getSpaceId(),
        name: values?.name || '',
        desc: values?.target || '',
        icon_uri: values?.icon_uri?.[0]?.uid || '',
        flow_mode: flowMode,
        bind_biz_id: bindBizId,
        bind_biz_type: bindBizType,
        project_id: projectId,

      };

      const resp = await workflowApi.CreateWorkflow(reqParams, {
        headers: {
          'x-locale': currentLocale,
        },
      });

      const content =
        I18n.t('workflow_list_create_success');
      Toast.success({
        content,
        showClose: false,
      });
      await onSuccess?.({
        workflowId: resp.data?.workflow_id,
        flowMode,
      });
      reporter.info({
        message: 'workflow_info_modal_create_success',
        namespace: 'workflow',
      });
    } catch (error) {
      reporter.error({
        message: 'workflow_info_modal_create_fail',
        namespace: 'workflow',
        error,
      });
      handleError(error);
    }
  };

  const title = useMemo(() => {
    const modelTitle = getModalTitle(mode, flowMode);
    if (customTitleRender && isFunction(customTitleRender)) {
      return customTitleRender(modelTitle);
    }
    return modelTitle;
  }, [mode, flowMode, customTitleRender]);
  const labels = useMemo<{
    nameLabel?: string;
    namePlaceholder?: string;
    descLabel?: string;
    descPlaceholder?: string;
    nameFormatRuleLabel?: string;
    nameRequiredLabel?: string;
    descRequiredLabel?: string;
  }>(() => ({
      nameLabel: I18n.t('workflow_list_create_modal_name_label'),
      namePlaceholder: I18n.t('workflow_list_create_modal_name_placeholder'),
      descLabel: I18n.t('workflow_list_create_modal_description_label'),
      descPlaceholder: I18n.t(
        'workflow_list_create_modal_description_placeholder',
      ),
      nameRequiredLabel: I18n.t(
        'workflow_list_create_modal_name_rule_required',
      ),
      nameFormatRuleLabel: I18n.t('workflow_list_create_modal_name_rule_reg'),
      descRequiredLabel: I18n.t(
        'workflow_list_create_modal_description_rule_required',
      ),
    }), [flowMode]);
  const iconType = useMemo(() => {
    switch (flowMode) {

      case WorkflowMode.Workflow:
        return IconType.Workflow;

      default:
        return IconType.Workflow;
    }
  }, [flowMode]);
  return (
    <UIModal
      type="action-small"
      keepDOM={false}
      icon={null}
      visible={visible}
      onCancel={handleCancel}
      title={title}
      footer={
        <Space>
          <Button
            className="min-w-[96px]"
            color="primary"
            onClick={handleCancel}
            data-testid="workflow.list.create.cancel"
          >
            {I18n.t('workflow_list_create_modal_footer_cancel')}
          </Button>
          <LoadingButton
            className="min-w-[96px]"
            color="hgltplus"
            disabled={confirmDisabled}
            onClick={
              mode === 'add' ? handleCreateWorkflow : handleUpdateWorkflow
            }
            data-testid="workflow.list.create.submit"
          >
            {I18n.t('workflow_list_create_modal_footer_confirm')}
          </LoadingButton>
        </Space>
      }
    >
      <Form<Partial<FormValue>>
        ref={formRef}
        showValidateIcon={false}
        className={s['upload-form']}
        onValueChange={({ name, target }) => {
          setSensitiveTip(undefined);
          setConfirmDisabled(!name?.trim() || !target?.trim());
        }}
      >
        <PictureUpload
          noLabel
          fieldClassName={s['upload-field']}
          field="icon_uri"
          initValue={getPictureUploadInitValue(workFlow)}
          iconType={iconType}
          fileBizType={FileBizType.BIZ_BOT_WORKFLOW}
        />
        <UIFormTextArea
          stopValidateWithError
          className={s['textarea-single-line']}
          field="name"
          placeholder={labels.namePlaceholder}
          label={labels.nameLabel}
          // noErrorMessage
          initValue={workFlow?.name}
          rows={1}
          maxCount={WORKFLOW_NAME_MAX_LEN}
          maxLength={WORKFLOW_NAME_MAX_LEN}
          rules={[
            {
              required: true,
              message: labels.nameRequiredLabel,
            },
            {
              validator(_, value) {
                if (!WORKFLOW_NAME_REGEX.test(value)) {
                  return new CustomError(
                    REPORT_EVENTS.formValidation,
                    labels.nameFormatRuleLabel ?? '',
                  );
                }
                return true;
              },
            },
            ...nameValidators,
          ]}
          data-testid="workflow.list.create.name.input"
        />

        {/* Sessions can only be bound when Chatflow is created within the project */}
        {null}
        <UIFormTextArea
          field="target"
          className={s['textarea-multi-line']}
          label={labels.descLabel}
          placeholder={labels.descPlaceholder}
          initValue={workFlow?.desc}
          maxCount={600}
          maxLength={600}
          rules={[
            {
              required: true,
              message: labels.descRequiredLabel,
            },
          ]}
          data-testid="workflow.list.create.desc.input"
        />
        {typeof sensitiveTip === 'string' ? (
          <Form.ErrorMessage error={sensitiveTip} />
        ) : null}
      </Form>
    </UIModal>
  );
}
