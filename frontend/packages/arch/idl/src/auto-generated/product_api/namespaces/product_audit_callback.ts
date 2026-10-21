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



/** 回调类型枚举值 */






/** response status code used by interfaces under this service
服务rsp的通用状态码 */




/** NodeResult，record node status and context */




/** The overall review results of a task
 1. Every callback will include results from previous rounds
 2. Please maintain the idempotence of the interface
 3. Business logic should combine VerifyResult.turn, TaskResult.current_turn, TaskResult.has_next_turn
 to select results and do other processing */


/** result for a round of review */

/* eslint-enable */
