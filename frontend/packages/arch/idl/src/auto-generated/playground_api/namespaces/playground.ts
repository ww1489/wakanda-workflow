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















/** 分支 */






/** 上下文允许传输的类型 */




export enum CreateResourceType {
  Agent = 1,
  Project = 2,
  Knowledge = 3,
  Dataset = 4,
  Plugin = 5,
  Memorybase = 6,
}





/** ************************** bpm回调 相关 **************************** API豁免表单类型 */






export enum GetImageScene {
  Onboarding = 0,
  BackgroundImage = 1,
}

export enum ModelScene {
  Douyin = 1,
}

/** 查询生图任务 */
















export enum ScopeType {
  /** 企业下所有的（企业下生效） */
  All = 0,
  /** 我加入的（企业&个人都生效，不传默认Self） */
  Self = 1,
}

/** 申请管理列表 */
export enum SpaceApplyStatus {
  /** 所有 */
  All = 0,
  /** 已加入 */
  Joined = 1,
  /** 确认中 */
  Confirming = 2,
  /** 已拒绝 */
  Rejected = 3,
}



export enum SpaceMode {
  Normal = 0,
  DevMode = 1,
}

export enum SpaceRoleType {
  /** 默认 */
  Default = 0,
  /** owner */
  Owner = 1,
  /** 管理员 */
  Admin = 2,
  /** 普通成员 */
  Member = 3,
}

export enum SpaceTag {
  /** 专业版 */
  Professional = 1,
}

/** --------------------space相关 start-------------------------------- */
export enum SpaceType {
  /** 个人 */
  Personal = 1,
  /** 小组 */
  Team = 2,
}

export enum StyleStatus {
  Dark = 1,
  Light = 2,
}









export enum UserStatus {
  Normal = 1,
  Banned = 2,
  Canceled = 3,
}

export enum UserType {
  External = 0,
  Internal = 1,
}













/** ------------- chat flow end ---------------------  需要是蛇形。json string和developer_api 的匹配 */






export interface AppIDInfo {
  id?: string;
  name?: string;
  icon?: string;
}











export interface BotSpaceV2 {
  /** 空间id，新建为0 */
  id?: string;
  /** 发布平台 */
  app_ids?: Array<AppIDInfo>;
  /** 空间名称 */
  name?: string;
  /** 空间描述 */
  description?: string;
  /** 图标url */
  icon_url?: string;
  /** 空间类型 */
  space_type?: SpaceType;
  /** 发布平台 */
  connectors?: Array<ConnectorInfo>;
  /** 是否隐藏新建，复制删除按钮 */
  hide_operation?: boolean;
  /** 在team中的角色 1-owner 2-admin 3-member */
  role_type?: number;
  /** 空间模式 */
  space_mode?: SpaceMode;
  /** 是否显示端侧插件创建入口 */
  display_local_plugin?: boolean;
  /** 角色类型，枚举 */
  space_role_type?: SpaceRoleType;
  /** 空间标签 */
  space_tag?: SpaceTag;
  /** 企业id */
  enterprise_id?: string;
  /** 组织id */
  organization_id?: string;
  /** 空间owner uid */
  owner_user_id?: string;
  /** 空间owner昵称 */
  owner_name?: string;
  /** 空间owner用户名 */
  owner_user_name?: string;
  /** 空间owner图像 */
  owner_icon_url?: string;
  /** 当前访问用户加入空间状态 */
  space_apply_status?: SpaceApplyStatus;
  /** 空间成员总数，只有组织空间才查询 */
  total_member_num?: Int64;
}







/** 火山账号注销检查 */


























export interface ConnectorInfo {
  id?: string;
  name?: string;
  icon?: string;
}

export interface CozeBanner {
  banner_content?: string;
  color_scheme?: string;
}



/** coze专业版复制授权链接生成 */






/** coze专业版复制授权链接元信息获取 */




/** coze专业版复制任务确认 */








/** home_banner end */










export interface ExitSpaceV2Request {
  /** 空间id */
  space_id?: string;
  /** 权限转移user_id */
  transfer_user_id?: string;
}

export interface File {
  /** 文件ID */
  id?: string;
  /** 文件字节数 */
  bytes?: Int64;
  /** 上传时间戳，单位s */
  created_at?: Int64;
  /** 文件名 */
  file_name?: string;
  URL?: string;
  /** 目前仅 openapi 场景有用。白名单用户会返回文件 url */
  url?: string;
}





















































export interface GetImagexShortUrlData {
  /** 审核状态，key uri，value url 和 审核状态 */
  url_info?: Record<string, UrlInfo>;
}

export interface GetImagexShortUrlRequest {
  uris?: Array<string>;
  scene?: GetImageScene;
}

export interface GetImagexShortUrlResponse {
  data?: GetImagexShortUrlData;
  code: Int64;
  msg: string;
}







export interface GetPolicyContentRequest {
  /** 隐私链接secret */
  secret: string;
}

export interface GetPolicyContentResponse {
  'Content-Type'?: string;
  policy_content?: string;
  code: Int64;
  msg: string;
}

















export interface GetSpaceListV2Request {
  /** 搜索词 */
  search_word?: string;
  /** 企业id */
  enterprise_id?: string;
  /** 组织id */
  organization_id?: string;
  /** 范围类型 */
  scope_type?: ScopeType;
  /** 分页信息 */
  page?: number;
  /** 分页大小 -- page 和 size不传则认为不分页 */
  size?: number;
}

export interface GetSpaceListV2Response {
  data?: SpaceInfo;
  code: Int64;
  msg: string;
}















export interface HomeBannerDisplay {
  image_url?: string;
  /** 主标题 */
  main_title?: string;
  /** 副标题 */
  sub_title?: string;
  /** Button文案 */
  button_text?: string;
  /** Button跳转链接 */
  button_url?: string;
  /** 1-暗黑 2-明亮 */
  style?: StyleStatus;
}









export interface MGetUserBasicInfoRequest {
  user_ids: Array<string>;
  need_user_status?: boolean;
  /** 是否需要企业认证信息，前端通过AGW调用时默认为true */
  need_enterprise_identity?: boolean;
  /** 是否需要火山用户名 */
  need_volcano_user_name?: boolean;
}

export interface MGetUserBasicInfoResponse {
  id_user_info_map?: Record<Int64, UserBasicInfo>;
  code?: Int64;
  msg?: string;
}































export interface QuickStartConfig {
  /** 新手教程文档链接 */
  document_url?: string;
  /** 图片链接 */
  image_url?: string;
  /** 文案 */
  content?: string;
}







export interface SaveSpaceV2Request {
  /** 空间id */
  space_id?: string;
  /** 空间名称 */
  name: string;
  /** 空间描述 */
  description: string;
  /** 空间图像 */
  icon_uri: string;
  /** 空间类型 */
  space_type: SpaceType;
  /** 空间模式 */
  space_mode?: SpaceMode;
  /** 空间配置 */
  space_config?: SpaceConfigV2;
  /** 企业id */
  enterprise_id?: string;
  /** 组织id */
  organization_id?: string;
}

















/** 空间配置信息 */
export interface SpaceConfigV2 {
  /** 是否支持外部用户加入当前团队 */
  is_support_external_users_join_space?: boolean;
  /** 是否所有成员都具备发布权限 */
  is_all_member_can_publish?: boolean;
  /** 空间成员创建资源的开关配置 */
  create_resource_switch?: Record<CreateResourceType, boolean>;
  /** 是否禁止创建/编辑文件夹 */
  forbid_member_upsert_folder?: boolean;
}

export interface SpaceInfo {
  /** 用户加入空间列表 */
  bot_space_list?: Array<BotSpaceV2>;
  /** 是否有个人空间 */
  has_personal_space?: boolean;
  /** 个人创建team空间数量 */
  team_space_num?: number;
  /** 个人最大能创建的空间数量 */
  max_team_space_num?: number;
  /** 最近使用空间列表 */
  recently_used_space_list?: Array<BotSpaceV2>;
  /** 分页时生效 */
  total?: number;
  /** 分页时生效 */
  has_more?: boolean;
}







/** 和 developerapi保持一致 */


export interface TransferSpaceV2Request {
  /** 空间id */
  space_id?: string;
  /** 权限转移user_id */
  transfer_user_id?: string;
}















export interface UploadFileOpenRequest {
  /** 文件类型 */
  'Content-Type': string;
  /** 二进制数据 */
  Data: Blob;
}

export interface UploadFileOpenResponse {
  data?: File;
}

export interface UrlInfo {
  url?: string;
  review_status?: boolean;
}

export interface UserBasicInfo {
  user_id: Int64;
  user_type: UserType;
  /** 昵称 */
  user_name: string;
  /** 头像 */
  user_avatar: string;
  /** 用户名 */
  user_unique_name?: string;
  /** 用户标签 */
  user_label?: bot_common.UserLabel;
  /** 用户创建时间 */
  create_time?: Int64;
  /** 用户当前状态 */
  user_status?: UserStatus;
  /** 火山用户名 */
  volcano_user_name?: string;
  /** 是否是企业认证 */
  is_enterprise_authn?: boolean;
}







/** 专业版子用户配置 */






/** 获取专业版子用户配置 */



/* eslint-enable */
