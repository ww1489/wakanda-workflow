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

package memory

import (
	"gorm.io/gorm"

	database "github.com/coze-dev/coze-studio/backend/domain/memory/database/service"

	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
	"github.com/coze-dev/coze-studio/backend/infra/rdb"
	rdbService "github.com/coze-dev/coze-studio/backend/infra/rdb/impl/rdb"
	"github.com/coze-dev/coze-studio/backend/infra/storage"
)

type MemoryApplicationServices struct {
	DatabaseDomainSVC database.Database
	RDBDomainSVC      rdb.RDB
}

type ServiceComponents struct {
	IDGen idgen.IDGenerator
	DB    *gorm.DB

	TosClient storage.Storage

	CacheCli cache.Cmdable
}

func InitService(c *ServiceComponents) *MemoryApplicationServices {
	rdbSVC := rdbService.NewService(c.DB, c.IDGen)
	databaseDomainSVC := database.NewService(rdbSVC, c.DB, c.IDGen, c.TosClient, c.CacheCli)
	DatabaseApplicationSVC.DomainSVC = databaseDomainSVC

	return &MemoryApplicationServices{
		DatabaseDomainSVC: databaseDomainSVC,
		RDBDomainSVC:      rdbSVC,
	}
}
