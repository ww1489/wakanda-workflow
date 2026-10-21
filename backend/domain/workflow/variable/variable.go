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

package variable

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

type StoreConfig struct{}

type OptionFn func(*StoreConfig)

//go:generate mockgen -destination varmock/var_mock.go --package mockvar -source variable.go
type Store interface {
	Init(ctx context.Context)
	Get(ctx context.Context, path compose.FieldPath, opts ...OptionFn) (any, error)
	Set(ctx context.Context, path compose.FieldPath, value any, opts ...OptionFn) error
}
