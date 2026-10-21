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

package adaptor

import (
	"fmt"
	"regexp"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

func ValidateSupportedCanvas(canvas string) error {
	var c vo.Canvas
	if err := sonic.UnmarshalString(canvas, &c); err != nil {
		return err
	}
	return ValidateSupportedNodes(c.Nodes)
}

// ValidateSupportedNodes checks capabilities without compiling incomplete drafts
// or contacting models, plugins, databases, or knowledge services.
func ValidateSupportedNodes(nodes []*vo.Node) error {
	for _, n := range nodes {
		if n == nil {
			return fmt.Errorf("workflow canvas contains a null node")
		}
		if err := validateSupportedNode(n); err != nil {
			return fmt.Errorf("node %q (type %s): %w", n.ID, n.Type, err)
		}
		if err := ValidateSupportedNodes(n.Blocks); err != nil {
			return err
		}
	}
	return nil
}

func validateSupportedNode(n *vo.Node) error {
	t := entity.IDStrToNodeType(n.Type)
	if !entity.IsSupportedNodeType(t) {
		return fmt.Errorf("node type is not supported in the workflow-only edition")
	}
	if n.Data == nil {
		return nil // Incomplete drafts remain editable.
	}
	bs, err := sonic.Marshal(n.Data)
	if err != nil {
		return err
	}
	var inputData any
	if err := sonic.Unmarshal(bs, &inputData); err != nil {
		return err
	}
	if hasGlobalVariableReference(inputData) {
		return fmt.Errorf("agent and application variable references are not supported")
	}
	if n.Data.Inputs == nil {
		return nil
	}
	in := n.Data.Inputs
	if t == entity.NodeTypeHTTPRequester && in.HttpRequestNode != nil {
		templates := []string{in.APIInfo.URL}
		if in.Body.BodyData != nil {
			templates = append(templates, in.Body.BodyData.Json, in.Body.BodyData.RawText)
		}
		for _, template := range templates {
			if globalVariableTemplate.MatchString(template) {
				return fmt.Errorf("agent and application variable references are not supported")
			}
		}
	}
	if in.ChatHistorySetting != nil && in.ChatHistorySetting.EnableChatHistory {
		return fmt.Errorf("conversation history is not supported")
	}
	if t == entity.NodeTypeCodeRunner && in.CodeRunner != nil && in.Language != 0 && in.Language != 3 {
		return fmt.Errorf("only Python code is supported")
	}
	if t != entity.NodeTypeLLM {
		return nil
	}
	if in.LLM != nil && in.FCParam != nil && in.FCParam.KnowledgeFCParam != nil &&
		len(in.FCParam.KnowledgeFCParam.KnowledgeList) > 0 {
		return fmt.Errorf("knowledge bindings are not supported")
	}
	if in.LLMParam != nil {
		bs, err := sonic.Marshal(in.LLMParam)
		if err != nil {
			return err
		}
		var params vo.LLMParam
		if err = sonic.Unmarshal(bs, &params); err != nil {
			return err
		}
		for _, p := range params {
			if p != nil && p.Name == "enableChatHistory" && p.Input != nil && p.Input.Value != nil &&
				p.Input.Value.Content == true {
				return fmt.Errorf("conversation history is not supported")
			}
		}
	}
	return nil
}

var globalVariableTemplate = regexp.MustCompile(`\{\{\s*global_variable_(app|user|system)\b[^}]*}}`)

func hasGlobalVariableReference(value any) bool {
	switch v := value.(type) {
	case map[string]any:
		if v["type"] == string(vo.BlockInputValueTypeLiteral) {
			return false
		}
		if v["type"] == string(vo.BlockInputValueTypeRef) || v["type"] == string(vo.BlockInputValueTypeObjectRef) {
			if ref, ok := v["content"].(map[string]any); ok {
				switch ref["source"] {
				case string(vo.RefSourceTypeGlobalApp), string(vo.RefSourceTypeGlobalUser), string(vo.RefSourceTypeGlobalSystem):
					return true
				}
			}
		}
		for _, child := range v {
			if hasGlobalVariableReference(child) {
				return true
			}
		}
	case []any:
		for _, child := range v {
			if hasGlobalVariableReference(child) {
				return true
			}
		}
	}
	return false
}
