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

import { type refreshUserInfo as refreshUserInfoOfSdk, type getIsSettled as getIsSettledOfSdk, type getIsLogined as getIsLoginedOfSdk, type getUserInfo as getUserInfoOfSdk, type useIsSettled as useIsSettledOfSdk, type useIsLogined as useIsLoginedOfSdk, type useUserInfo as useUserInfoOfSdk, type getLoginStatus as getLoginStatusOfSdk, type useLoginStatus as useLoginStatusOfSdk, type useUserLabel as useUserLabelOfSdk } from '@coze-arch/foundation-sdk';
import { refreshUserInfo as refreshUserInfoImpl, getLoginStatus as getLoginStatusImpl, useLoginStatus as useLoginStatusImpl, getUserInfo as getUserInfoImpl, useUserInfo as useUserInfoImpl, useUserLabel as useUserLabelImpl } from '@coze-foundation/account-adapter';

/** @deprecated using getLoginStatus */
export const getIsSettled = (() =>
  getLoginStatus() !== 'settling') satisfies typeof getIsSettledOfSdk;
/** @deprecated using getLoginStatus */
export const getIsLogined = (() =>
  getLoginStatus() === 'logined') satisfies typeof getIsLoginedOfSdk;
export const getUserInfo: typeof getUserInfoOfSdk = getUserInfoImpl;
/** @deprecated useLoginStatus */
export const useIsSettled = (() => {
  const status = useLoginStatus();
  return status !== 'settling';
}) satisfies typeof useIsSettledOfSdk;
/** @deprecated useLoginStatus */
export const useIsLogined = (() => {
  const status = useLoginStatus();
  return status === 'logined';
}) satisfies typeof useIsLoginedOfSdk;
export const useUserInfo: typeof useUserInfoOfSdk = useUserInfoImpl;
export const useUserLabel: typeof useUserLabelOfSdk = useUserLabelImpl;

export const refreshUserInfo: typeof refreshUserInfoOfSdk = refreshUserInfoImpl;
export const useLoginStatus: typeof useLoginStatusOfSdk = useLoginStatusImpl;
export const getLoginStatus: typeof getLoginStatusOfSdk = getLoginStatusImpl;
