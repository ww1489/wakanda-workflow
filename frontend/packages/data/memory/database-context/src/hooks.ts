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

import { useContext } from 'react';

import { useStoreWithEqualityFn } from 'zustand/traditional';
import { shallow } from 'zustand/shallow';
import { REPORT_EVENTS } from '@coze-arch/report-events';
import { CustomError } from '@coze-arch/bot-error';


import { type IParamsStore } from './params-store';


import {
  type CallbacksType,
  type PluginNavType,
  DatabaseParamsStoreContext,
} from './context';

export const useDatabaseParamsStore: <T>(
  selector: (store: IParamsStore) => T,
) => T = selector => {
  const context = useContext(DatabaseParamsStoreContext);

  if (!context.paramsStore) {
    throw new CustomError(REPORT_EVENTS.normalError, 'params store context');
  }

  return useStoreWithEqualityFn(context.paramsStore, selector, shallow);
};

export const useDatabaseParams = () => {
  const params = useDatabaseParamsStore(store => store.params);
  return params;
};

export const useDataCallbacks: () => CallbacksType = () => {
  const {
    callbacks: { onStatusChange, onUpdateDisplayName },
  } = useContext(DatabaseParamsStoreContext);

  return { onStatusChange, onUpdateDisplayName };
};

export const useDataNavigate: () => PluginNavType = () => {
  const { resourceNavigate } = useContext(DatabaseParamsStoreContext);

  return resourceNavigate;
};
