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

/* eslint-disable */
/* tslint:disable */
// @ts-nocheck

import * as bot_common from './bot_common';

export type Int64 = string | number;



export interface ApiInfo {
  /** api id */
  api_id?: string;
  /** api名称 */
  name?: string;
  /** api描述 */
  description?: string;
}

export interface BackgroundImageDetail {
  origin_image_url?: string;
  image_url?: string;
  theme_color?: string;
  /** 渐变位置 */
  gradient_position?: GradientPosition;
  /** 裁剪画布位置 */
  canvas_position?: CanvasPosition;
}

export interface BackgroundImageInfo {
  /** web端背景图 */
  web_background_image?: BackgroundImageDetail;
  /** 移动端背景图 */
  mobile_background_image?: BackgroundImageDetail;
}

export interface BotInfo {
  /** bot id */
  bot_id?: string;
  /** bot名称 */
  name?: string;
  /** bot描述 */
  description?: string;
  /** bot图像url */
  icon_url?: string;
  /** 创建时间 */
  create_time?: Int64;
  /** 更新时间 */
  update_time?: Int64;
  /** 版本 */
  version?: string;
  /** prompt 信息 */
  prompt_info?: PromptInfo;
  /** 开场白 */
  onboarding_info?: OnboardingInfoV2;
  /** bot 类型，single agent or multi agent */
  bot_mode?: bot_common.BotMode;
  /** 选择的语音信息 */
  voice_data_list?: Array<VoiceData>;
  /** 模型信息 */
  model_info?: ModelInfo;
  /** 插件信息列表 */
  plugin_info_list?: Array<PluginInfo>;
  /** 知识库信息 */
  knowledge?: CommonKnowledge;
  /** workflow信息列表 */
  workflow_info_list?: Array<WorkflowInfo>;
  /** 快捷指令信息列表 */
  shortcut_commands?: Array<ShortcutCommandInfo>;
  /** 音色配置 */
  voice_info_list?: Array<Voice>;
  /** 默认用户输入类型 */
  default_user_input_type?: string;
  /** 用户问题建议 */
  suggest_reply_info?: SuggestReplyInfo;
  /** 背景图片 */
  background_image_info?: BackgroundImageInfo;
  /** 变量列表 */
  variables?: Array<Variable>;
  /** owner_id */
  owner_user_id?: string;
  folder_id?: string;
  /** tts配置 */
  media_config?: MediaConfig;
}

export interface CanvasPosition {
  width?: number;
  height?: number;
  left?: number;
  top?: number;
}

export interface CommonKnowledge {
  /** 知识库信息 */
  knowledge_infos?: Array<KnowledgeInfo>;
}

export interface GradientPosition {
  left?: number;
  right?: number;
}



export interface KnowledgeInfo {
  /** 知识库id */
  id?: string;
  /** 知识库名称 */
  name?: string;
}

export interface MediaConfig {
  /** 是否关闭语音通话，true:关闭 false:开启  默认为false */
  is_voice_call_closed?: boolean;
}

export interface ModelInfo {
  /** 模型id */
  model_id?: string;
  /** 模型名称 */
  model_name?: string;
  /** 生成随机性 没配置不返回 */
  temperature?: number;
  /** top p 没配置不返回 */
  top_p?: number;
  /** 频率惩罚 没配置不返回 */
  frequency_penalty?: number;
  /** 存在惩罚 没配置不返回 */
  presence_penalty?: number;
  /** 生成时，采样候选集的大小 没配置不返回 */
  top_k?: number;
  /** 携带上下文轮数 */
  context_round?: number;
  /** 最大回复长度 */
  max_tokens?: number;
  /** 输出格式 text、markdown、json */
  response_format?: string;
  /** 缓存配置 */
  cache_type?: string;
  /** sp拼接当前时间 */
  sp_current_time?: boolean;
  /** sp拼接防泄露指令 */
  sp_anti_leak?: boolean;
  /** 模型个性化配置参数 */
  parameters?: Record<string, string>;
}



export interface OnboardingInfoV2 {
  /** 对应 Coze Opening Dialog
开场白 */
  prologue?: string;
  /** 建议问题 */
  suggested_questions?: Array<string>;
  /** 开场白模型 */
  onboarding_mode?: bot_common.OnboardingMode;
  /** LLM生成，用户自定义 Prompt */
  customized_onboarding_prompt?: string;
  /** 开场白预设问题展示方式 默认0 随机展示 */
  suggested_questions_show_mode?: bot_common.SuggestedQuestionsShowMode;
}

export interface PluginInfo {
  /** 插件id */
  plugin_id?: string;
  /** 插件名称 */
  name?: string;
  /** 插件描述 */
  description?: string;
  /** 插件图片url */
  icon_url?: string;
  /** 插件包含的api列表 */
  api_info_list?: Array<ApiInfo>;
}

export interface PrefixPromptInfo {
  /** 前缀提示词 */
  prefix_prompt?: string;
  /** 不支持前缀提示词部分 */
  dynamic_prompt?: string;
}

/** bot管理 */
export interface PromptInfo {
  /** 文本prompt */
  prompt?: string;
  /** 提示词模式 */
  prompt_mode?: string;
  /** 前缀提示词模式下的prompt内容 */
  prefix_prompt_info?: PrefixPromptInfo;
}

export interface ShortcutCommandComponent {
  /** panel参数
参数名字 */
  name?: string;
  /** 参数描述 */
  description?: string;
  /** 输入类型 text、select、file */
  type?: string;
  /** 请求工具时，参数的key 对应tool的参数名称，没有则为不返回 */
  tool_parameter?: string;
  /** type为select时的可选项列表 or type为file时，支持哪些类型 image、doc、table、audio、video、zip、code、txt、ppt */
  options?: Array<string>;
  /** 默认值 没配置时不返回 */
  default_value?: string;
  /** 是否隐藏不展示 线上bot tool类型的快捷指令不返回hide=true的component */
  is_hide?: boolean;
}

export interface ShortcutCommandInfo {
  /** 快捷指令id */
  id?: string;
  /** 快捷指令按钮名称 */
  name?: string;
  /** 快捷指令 */
  command?: string;
  /** 快捷指令描述 */
  description?: string;
  /** 指令query模版 */
  query_template?: string;
  /** 快捷指令icon */
  icon_url?: string;
  /** 组件列表（参数列表） */
  components?: Array<ShortcutCommandComponent>;
  /** tool信息 */
  tool?: ShortcutCommandToolInfo;
  /** multi的指令时，该指令由哪个节点执行 没配置不返回 */
  agent_id?: string;
  /** chatsdk 使用 */
  send_type?: string;
  /** chatsdk 使用，表单的schema */
  card_schema?: string;
}

export interface ShortcutCommandToolInfo {
  name?: string;
  /** tool类型 workflow plugin */
  type?: string;
  plugin_id?: string;
  plugin_api_name?: string;
  workflow_id?: string;
  params?: Array<ShortcutToolParam>;
}

export interface ShortcutToolParam {
  name?: string;
  is_required?: boolean;
  description?: string;
  type?: string;
  default_value?: string;
  /** 是否是panel参数 */
  is_refer_component?: boolean;
}

export interface SuggestReplyInfo {
  /** 回复模式 */
  reply_mode?: string;
  /** custom 模式下的自定义 prompt */
  customized_prompt?: string;
}

export interface Variable {
  /** 变量名 */
  keyword?: string;
  /** 默认值 */
  default_value?: string;
  /** 变量类型 */
  variable_type?: string;
  /** 变量来源 */
  channel?: string;
  /** 变量描述 */
  description?: string;
  /** 是否启用 */
  enable?: boolean;
  /** 变量默认支持在Prompt中访问，取消勾选后将不支持在Prompt中访问（仅能在Workflow中访问 */
  prompt_enable?: boolean;
}

export interface Voice {
  /** 唯一id */
  voice_id?: string;
  /** 音色语种code */
  language_code?: string;
}

export interface VoiceData {
  /** 唯一id */
  id?: string;
  /** 音色语种code */
  language_code?: string;
  /** 音色语种名称 */
  language_name?: string;
  /** 音色名称 */
  name?: string;
  /** 音色 style_id */
  style_id?: string;
  /** 预览文本内容 */
  preview_text?: string;
  /** 预览音色内容 */
  preview_audio?: string;
}

export interface WorkflowInfo {
  /** workflow_id */
  id?: string;
  /** workflow名称 */
  name?: string;
  /** workflow描述 */
  description?: string;
  /** workflow图片url */
  icon_url?: string;
}
/* eslint-enable */
