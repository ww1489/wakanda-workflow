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





/** 权益流水状态 */


/** Benefit 所作用的实体类型 */




/** 权益历史记录类型 */


/** 权益流水根节点类型 */


/** 权益类型
 40 -59 免费次数
 60 - 99 限流
 100-109 资源点
 110-129 Fornax
 130-149 WorkSpace
 150-169 运维
 170-179 知识库
 180-199 语音
 200-219 租户相关
 220-229 发布相关 */


/** 权益使用模式 */












/** Type为BotRequest(5)的子类型SubType */


/** Type为Knowledge(6)的子类型SubType */


/** Type为LongTermMemory(10)的子类型SubType */


/** Type为Model(1)的子类型SubType */


/** Type为ModelTPM(8)的子类型SubType */


/** Type为Plugin(2)的子类型SubType */


/** Type为RTC(4)的子类型SubType */


/** Type为Seat(7)的子类型SubType */


/** Type为Voice(3)的子类型SubType */


/** Type为WorkSpace(9)的子类型SubType */




/** 校验结果。通常结合BenefitType */


/** 权益校验点位 */


/** 权益流水对应消耗的资源类型 */




/** 权益流水的成本归属用户类型 */




/** 用户权益套餐状态 */




















/** 用量维度标识（当前用于分账） */




/** 权益流水的权益类型（用于对客） */






/** 资源归属的实体类型 */






/** 场景 */


export enum UserLevel {
  /** 免费版。 */
  Free = 0,
  /** 海外
PremiumLite */
  PremiumLite = 10,
  /** Premium */
  Premium = 15,
  PremiumPlus = 20,
  /** 国内
V1火山专业版 */
  V1ProInstance = 100,
  /** 个人旗舰版 */
  ProPersonal = 110,
  /** 团队版 */
  Team = 120,
  /** 企业版 */
  Enterprise = 130,
}








/* eslint-enable */
