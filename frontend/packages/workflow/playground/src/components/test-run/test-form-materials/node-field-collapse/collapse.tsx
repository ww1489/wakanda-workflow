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

/* eslint-disable @coze-arch/no-batch-import-or-export */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { FormBaseGroupCollapse } from '@coze-workflow/test-run-next';
import { I18n } from '@coze-arch/i18n';


import { ModeSwitch } from './node-switch';


import css from './collapse.module.less';

export const NodeFieldCollapse: React.FC<React.PropsWithChildren> = ({
  children,
}) => (
    <FormBaseGroupCollapse
      label={I18n.t('wf_test_run_form_input_collapse_label')}
      extra={
        <div className={css.extra}>
          <ModeSwitch />

        </div>
      }
    >
      {children}
    </FormBaseGroupCollapse>
  );
