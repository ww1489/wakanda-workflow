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





/** 版本兼容：0-旧版本 1-可回退的新版本 2-不可回退的新版本 3-可回退的新版本(不再提示) */


/** AnswerActions */






export enum BotMode {
  SingleMode = 0,
  MultiMode = 1,
  WorkflowMode = 2,
}











/** 上下文允许传输的类型 */


































/** onboarding内容生成模式 */
export enum OnboardingMode {
  /** 不需要 */
  NO_NEED = 1,
  /** 人工指定内容（多语言支持由LLM兜底） */
  USE_MANUAL = 2,
  /** 由LLM生成 */
  USE_LLM = 3,
}





export enum PluginFrom {
  /** 插件来源 opencoze 使用 */
  Default = 0,
  FromSaas = 1,
}







/** struct ContentAttachment {
    1: required string FileID (api.body = "file_id")
}
 struct MetaContent{
     1: required string Type (agw.key="type"),
     2: optional string Text ( agw.key="text"),
     3: optional string FileID (agw.key="file_id"),
     4: optional string FileURL (agw.key="file_url"),
     5: optional string Card (agw.key="card"),
 }
 struct EnterMessage  {
     1: required string Role (agw.key = "role")
     2: string Content(agw.key = "content")     // 内容
     3: map<string,string> MetaData(agw.key = "meta_data")
     4: string ContentType(agw.key = "content_type")//text/card/object_string
     5: string Type(agw.key = "type")
 }
 struct OpenMessageApi {
     1: string Id(agw.key = "id")             // 主键ID
     2: string BotId(agw.key = "bot_id")        // bot id //已TODO 所有的i64加注解str,入参和出参都要
     3: string Role(agw.key = "role")
     4: string Content(agw.key = "content")          // 内容
     5: string ConversationId(agw.key = "conversation_id")   // conversation id
     6: map<string,string> MetaData(agw.key = "meta_data")
     7: string CreatedAt(agw.key = "created_at")      // 创建时间
     8: string UpdatedAt(agw.key = "updated_at")      // 更新时间 //已TODO 时间改成int
     9: string ChatId(agw.key = "chat_id")
     10: string ContentType(agw.key = "content_type")
     11: string Type(agw.key = "type")
 } */






export enum SuggestedQuestionsShowMode {
  Random = 0,
  All = 1,
}









export enum WorkflowMode {
  Workflow = 0,
  Imageflow = 1,
  SceneFlow = 2,
  ChatFlow = 3,
  All = 100,
}





















/** bot ext */


/** bot信息 */


/** bot信息 for 更新 */
























































/** suggest */






/** 时间胶囊信息 */


export interface UserLabel {
  label_id?: string;
  label_name?: string;
  icon_uri?: string;
  icon_url?: string;
  jump_link?: string;
}









/** tts Voices */







/* eslint-enable */
