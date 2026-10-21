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

import { describe, expect, it } from 'vitest';
import { StandardNodeType } from '@coze-workflow/base/types';

import { getEnabledNodeTypes } from './get-enabled-node-types';

const context = {
  loopSelected: false,
  isProject: false,
  isSupportImageflowNodes: false,
  isSceneFlow: false,
  isBindDouyin: false,
};

describe('standalone workflow nodes', () => {
  it('keeps integrations and hides knowledge, agent and conversation nodes', () => {
    const nodes = getEnabledNodeTypes(context);
    expect(nodes).toEqual(
      expect.arrayContaining([
        StandardNodeType.LLM,
        StandardNodeType.Code,
        StandardNodeType.Http,
        StandardNodeType.SubWorkflow,
        StandardNodeType.Api,
        StandardNodeType.Database,
        StandardNodeType.DatabaseQuery,
      ]),
    );
    for (const type of [
      StandardNodeType.Dataset,
      StandardNodeType.DatasetWrite,
      StandardNodeType.VariableAssign,
      StandardNodeType.CreateConversation,
      StandardNodeType.QueryMessageList,
    ]) {
      expect(nodes).not.toContain(type);
    }
  });

  it('keeps loop controls inside loops', () => {
    expect(getEnabledNodeTypes(context)).not.toContain(StandardNodeType.Break);
    expect(getEnabledNodeTypes({ ...context, loopSelected: true })).toEqual(
      expect.arrayContaining([
        StandardNodeType.Break,
        StandardNodeType.Continue,
        StandardNodeType.SetVariable,
      ]),
    );
  });
});
