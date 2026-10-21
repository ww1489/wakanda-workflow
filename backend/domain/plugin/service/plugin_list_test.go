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

package service

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/plugin/dto"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
	"github.com/coze-dev/coze-studio/backend/domain/plugin/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
)

// Exercise the service and its real repository together: name searches must work
// without a registered search service, and count and page queries share the scope.
func TestListPluginsUsesMySQL(t *testing.T) {
	for _, tt := range []struct {
		name      string
		draft     bool
		keyword   *string
		sortBy    dto.SortField
		ascending bool
	}{
		{name: "draft keyword", draft: true, keyword: ptr.Of("50%_off'"), sortBy: dto.SortByUpdatedAt},
		{name: "published keyword", keyword: ptr.Of("Weather"), sortBy: dto.SortByCreatedAt, ascending: true},
		{name: "draft without keyword", draft: true, sortBy: dto.SortByCreatedAt},
		{name: "published empty keyword", keyword: ptr.Of(""), sortBy: dto.SortByUpdatedAt},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			t.Cleanup(func() { _ = sqlDB.Close() })
			db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}))
			require.NoError(t, err)

			table := "plugin"
			if tt.draft {
				table = "plugin_draft"
			}
			where := fmt.Sprintf("`%s`.`space_id` = ?", table)
			args := []driver.Value{int64(42)}
			if tt.draft {
				where += " AND `plugin_draft`.`app_id` = ? AND `plugin_draft`.`developer_id` = ?"
				args = append(args, int64(0), int64(7))
			}
			if tt.keyword != nil && *tt.keyword != "" {
				where += " AND LOCATE(LOWER(?), LOWER(JSON_UNQUOTE(JSON_EXTRACT(manifest, '$.name_for_human')))) > 0"
				args = append(args, *tt.keyword)
			}
			if tt.draft {
				where += " AND `plugin_draft`.`deleted_at` IS NULL"
			}
			mock.ExpectQuery(fmt.Sprintf("SELECT count(*) FROM `%s` WHERE %s", table, where)).
				WithArgs(args...).WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(3))
			order := "DESC"
			if tt.ascending {
				order = "ASC"
			}
			mock.ExpectQuery(fmt.Sprintf("SELECT * FROM `%s` WHERE %s ORDER BY `%s`.`%s` %s LIMIT ? OFFSET ?", table, where, table, tt.sortBy, order)).
				WithArgs(append(args, 2, 2)...).
				WillReturnRows(sqlmock.NewRows([]string{"id", "space_id", "developer_id", "manifest"}).
					AddRow(101, 42, 7, `{"name_for_human":"Weather"}`))

			svc := &pluginServiceImpl{pluginRepo: repository.NewPluginRepo(&repository.PluginRepoComponents{DB: db})}
			page := dto.PageInfo{Name: tt.keyword, Page: 2, Size: 2, SortBy: &tt.sortBy, OrderByACS: &tt.ascending}
			var plugins []*entity.PluginInfo
			var total int64
			if tt.draft {
				resp, err := svc.ListDraftPlugins(context.Background(), &dto.ListDraftPluginsRequest{SpaceID: 42, APPID: 0, DeveloperID: 7, PageInfo: page})
				require.NoError(t, err)
				plugins, total = resp.Plugins, resp.Total
			} else {
				plugins, total, err = svc.ListCustomOnlinePlugins(context.Background(), 42, page)
				require.NoError(t, err)
			}
			require.Equal(t, int64(3), total)
			require.Len(t, plugins, 1)
			require.Equal(t, int64(101), plugins[0].ID)
			require.Equal(t, int64(7), plugins[0].DeveloperID)
			require.Equal(t, "Weather", plugins[0].Manifest.NameForHuman)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
