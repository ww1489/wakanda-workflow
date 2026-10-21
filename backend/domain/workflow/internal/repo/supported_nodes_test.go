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

package repo

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/repo/dal/query"
	mockidgen "github.com/coze-dev/coze-studio/backend/internal/mock/infra/idgen"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
)

func TestUnsupportedCanvasCannotBePersisted(t *testing.T) {
	// No database is configured: rejection must happen before any write or
	// reference update, including publishing an existing successful draft.
	r := &RepositoryImpl{}
	canvas := `{"nodes":[{"id":"knowledge","type":"27"}]}`
	require.ErrorContains(t, r.CreateOrUpdateDraft(t.Context(), 1, &vo.DraftInfo{
		Canvas: canvas,
	}), "not supported")
	require.ErrorContains(t, r.CreateVersion(t.Context(), 1, &vo.VersionInfo{
		CanvasInfo: vo.CanvasInfo{Canvas: canvas},
	}, nil), "not supported")
}

func TestCopyWorkflowRejectsUnsupportedCanvasBeforeWriting(t *testing.T) {
	for _, modified := range []bool{false, true} {
		t.Run(map[bool]string{false: "source", true: "replacement"}[modified], func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = sqlDB.Close() })
			db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}))
			require.NoError(t, err)
			ids := mockidgen.NewMockIDGenerator(gomock.NewController(t))
			ids.EXPECT().GenID(gomock.Any()).Return(int64(2), nil)
			r := &RepositoryImpl{IDGenerator: ids, query: query.Use(db)}
			unsupported := `{"nodes":[{"id":"knowledge","type":"6"}]}`
			source := unsupported
			policy := vo.CopyWorkflowPolicy{}
			if modified {
				source = `{"nodes":[]}`
				policy.ModifiedCanvasSchema = &unsupported
			}
			mock.ExpectQuery("SELECT .* FROM `workflow_meta`").
				WillReturnRows(sqlmock.NewRows([]string{"id", "mode"}).AddRow(1, 0))
			mock.ExpectQuery("SELECT .* FROM `workflow_draft`").
				WillReturnRows(sqlmock.NewRows([]string{"id", "canvas"}).AddRow(1, source))

			_, err = r.CopyWorkflow(t.Context(), 1, policy)
			require.ErrorContains(t, err, "not supported")
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCopyWorkflowRejectsApplicationSourceOrTarget(t *testing.T) {
	for _, sourceAppID := range []int64{0, 1} {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() { _ = sqlDB.Close() })
		db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}))
		require.NoError(t, err)
		ids := mockidgen.NewMockIDGenerator(gomock.NewController(t))
		ids.EXPECT().GenID(gomock.Any()).Return(int64(2), nil)
		r := &RepositoryImpl{IDGenerator: ids, query: query.Use(db)}
		policy := vo.CopyWorkflowPolicy{}
		if sourceAppID == 0 {
			policy.TargetAppID = ptr.Of(int64(1))
		}
		mock.ExpectQuery("SELECT .* FROM `workflow_meta`").
			WillReturnRows(sqlmock.NewRows([]string{"id", "mode", "app_id"}).AddRow(1, 0, sourceAppID))
		mock.ExpectQuery("SELECT .* FROM `workflow_draft`").
			WillReturnRows(sqlmock.NewRows([]string{"id", "canvas"}).AddRow(1, `{"nodes":[]}`))

		_, err = r.CopyWorkflow(t.Context(), 1, policy)
		require.ErrorContains(t, err, "application workflows")
		require.NoError(t, mock.ExpectationsWereMet())
	}
}
