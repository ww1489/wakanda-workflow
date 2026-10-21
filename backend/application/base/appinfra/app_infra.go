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

package appinfra

import (
	"context"
	"fmt"
	"os"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/bizpkg/config"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"github.com/coze-dev/coze-studio/backend/infra/cache/impl/redis"
	coderunner "github.com/coze-dev/coze-studio/backend/infra/coderunner/impl"
	"github.com/coze-dev/coze-studio/backend/infra/idgen/impl/idgen"
	"github.com/coze-dev/coze-studio/backend/infra/imagex"
	"github.com/coze-dev/coze-studio/backend/infra/imagex/impl/veimagex"
	"github.com/coze-dev/coze-studio/backend/infra/orm/impl/mysql"
	storage "github.com/coze-dev/coze-studio/backend/infra/storage/impl"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

type AppDependencies struct {
	DB           *gorm.DB
	CacheCli     cache.Cmdable
	IDGenSVC     idgen.IDGenerator
	ImageXClient imagex.ImageX
	OSS          storage.Storage
	CodeRunner   coderunner.Runner
}

func Init(ctx context.Context) (*AppDependencies, error) {
	deps := &AppDependencies{}
	var err error
	deps.OSS, err = storage.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("init tos client failed, err=%w", err)
	}

	deps.DB, err = mysql.New()
	if err != nil {
		return nil, fmt.Errorf("init db failed, err=%w", err)
	}

	deps.CacheCli = redis.New()

	deps.IDGenSVC, err = idgen.New(deps.CacheCli)
	if err != nil {
		return nil, fmt.Errorf("init id gen svc failed, err=%w", err)
	}

	err = config.Init(ctx, deps.DB, deps.OSS) // Depends on MySQL、Idgen and OSS initialization
	if err != nil {
		return nil, fmt.Errorf("init model config failed, err=%w", err)
	}

	basicConfig, err := config.Base().GetBaseConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("get basic config failed, err=%w", err)
	}

	deps.ImageXClient, err = initImageX(ctx)
	if err != nil {
		return nil, fmt.Errorf("init imagex client failed, err=%w", err)
	}

	deps.CodeRunner = coderunner.New(basicConfig)

	return deps, nil
}

func initImageX(ctx context.Context) (imagex.ImageX, error) {
	uploadComponentType := os.Getenv(consts.FileUploadComponentType)
	if uploadComponentType != consts.FileUploadComponentTypeImagex {
		return storage.NewImagex(ctx)
	}

	return veimagex.NewDefault()
}
