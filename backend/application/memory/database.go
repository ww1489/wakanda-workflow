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
	"context"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/api/model/base"
	"github.com/coze-dev/coze-studio/backend/api/model/data/database/table"

	document "github.com/coze-dev/coze-studio/backend/api/model/data/database/sheet"

	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"

	model "github.com/coze-dev/coze-studio/backend/crossdomain/database/model"
	crossuser "github.com/coze-dev/coze-studio/backend/crossdomain/user"
	"github.com/coze-dev/coze-studio/backend/domain/memory/database/entity"
	databaseEntity "github.com/coze-dev/coze-studio/backend/domain/memory/database/entity"
	database "github.com/coze-dev/coze-studio/backend/domain/memory/database/service"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/conv"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

type DatabaseApplicationService struct {
	DomainSVC database.Database
}

var DatabaseApplicationSVC = DatabaseApplicationService{}

func (d *DatabaseApplicationService) ListDatabase(ctx context.Context, req *table.ListDatabaseRequest) (*table.ListDatabaseResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	if req.SpaceID == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "space id is required"))
	}

	spaces, err := crossuser.DefaultSVC().GetUserSpaceList(ctx, *uid)
	if err != nil {
		return nil, err
	}
	if len(spaces) == 0 || spaces[0].ID != *req.SpaceID {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "space id is invalid"))
	}

	res, err := d.DomainSVC.ListDatabase(ctx, convertListDatabase(req))
	if err != nil {
		return nil, err
	}

	bindDatabases := make([]*databaseEntity.Database, 0)
	if req.GetBotID() != 0 {
		return nil, fmt.Errorf("agent database bindings are not supported")
	}

	return convertListDatabaseRes(res, bindDatabases), nil
}

func (d *DatabaseApplicationService) GetDatabaseByID(ctx context.Context, req *table.SingleDatabaseRequest) (*table.SingleDatabaseResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	basics := make([]*model.DatabaseBasic, 1)
	b := &model.DatabaseBasic{
		ID: req.ID,
	}
	if req.IsDraft {
		b.TableType = table.TableType_DraftTable
	} else {
		b.TableType = table.TableType_OnlineTable
	}

	b.NeedSysFields = req.NeedSysFields
	basics[0] = b

	res, err := d.DomainSVC.MGetDatabase(ctx, &database.MGetDatabaseRequest{
		Basics: basics,
	})
	if err != nil {
		return nil, err
	}

	if len(res.Databases) == 0 {
		return nil, fmt.Errorf("database %d not found", req.GetID())
	}

	if res.Databases[0].CreatorID != *uid {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "creator id is invalid"))
	}

	return ConvertDatabaseRes(res.Databases[0]), nil
}

func (d *DatabaseApplicationService) AddDatabase(ctx context.Context, req *table.AddDatabaseRequest) (*table.SingleDatabaseResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}
	if *uid != req.CreatorID {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "creator id is invalid"))
	}

	if req.GetTableName() == "database" {
		return nil, errorx.New(errno.ErrMemoryDatabaseNameInvalid)
	}

	spaces, err := crossuser.DefaultSVC().GetUserSpaceList(ctx, *uid)
	if err != nil {
		return nil, err
	}
	if len(spaces) == 0 || spaces[0].ID != req.SpaceID {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "space id is invalid"))
	}

	res, err := d.DomainSVC.CreateDatabase(ctx, convertAddDatabase(req))
	if err != nil {
		return nil, err
	}

	databaseRes := res.Database

	return ConvertDatabaseRes(databaseRes), nil
}

func (d *DatabaseApplicationService) UpdateDatabase(ctx context.Context, req *table.UpdateDatabaseRequest) (*table.SingleDatabaseResponse, error) {
	err := d.ValidateAccess(ctx, req.ID, table.TableType_OnlineTable)
	if err != nil {
		return nil, err
	}

	res, err := d.DomainSVC.UpdateDatabase(ctx, ConvertUpdateDatabase(req))
	if err != nil {
		return nil, err
	}

	return convertUpdateDatabaseResult(res), nil
}

func (d *DatabaseApplicationService) DeleteDatabase(ctx context.Context, req *table.DeleteDatabaseRequest) (*table.DeleteDatabaseResponse, error) {
	err := d.ValidateAccess(ctx, req.ID, table.TableType_OnlineTable)
	if err != nil {
		return nil, err
	}

	err = d.DomainSVC.DeleteDatabase(ctx, &database.DeleteDatabaseRequest{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	return &table.DeleteDatabaseResponse{
		Code:     0,
		Msg:      "success",
		BaseResp: base.NewBaseResp(),
	}, nil
}

func (d *DatabaseApplicationService) ListDatabaseRecords(ctx context.Context, req *table.ListDatabaseRecordsRequest) (*table.ListDatabaseRecordsResponse, error) {
	tableType := table.TableType_OnlineTable
	if req.GetBotID() > 0 {
		tableType = table.TableType_DraftTable
	}
	err := d.ValidateAccess(ctx, req.DatabaseID, tableType)
	if err != nil {
		return nil, err
	}

	databaseID := req.DatabaseID
	if req.GetBotID() == 0 {
		databaseID, err = getDatabaseID(ctx, req.TableType, req.DatabaseID)
		if err != nil {
			return nil, err
		}
	}

	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	domainReq := &database.ListDatabaseRecordRequest{
		DatabaseID: databaseID,
		TableType:  req.TableType,
		Limit:      int(req.Limit),
		Offset:     int(req.Offset),
		UserID:     *uid,
	}
	// FilterCriterion, NotFilterByUserID, OrderByList not use

	res, err := d.DomainSVC.ListDatabaseRecord(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	return convertListDatabaseRecordsRes(res), nil
}

func (d *DatabaseApplicationService) UpdateDatabaseRecords(ctx context.Context, req *table.UpdateDatabaseRecordsRequest) (*table.UpdateDatabaseRecordsResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	err := d.ValidateAccess(ctx, req.DatabaseID, table.TableType_OnlineTable)
	if err != nil {
		return nil, err
	}

	databaseID, err := getDatabaseID(ctx, req.GetTableType(), req.GetDatabaseID())
	if err != nil {
		return nil, err
	}

	dataRes := make([]map[string]string, 0)
	if len(req.GetRecordDataAdd()) > 0 {
		err := d.DomainSVC.AddDatabaseRecord(ctx, &database.AddDatabaseRecordRequest{
			DatabaseID: databaseID,
			TableType:  req.GetTableType(),
			Records:    req.GetRecordDataAdd(),
			UserID:     *uid,
		})
		if err != nil {
			return nil, err
		}

		dataRes = append(dataRes, req.GetRecordDataAdd()...)
	}

	if len(req.GetRecordDataAlter()) > 0 {
		err := d.DomainSVC.UpdateDatabaseRecord(ctx, &database.UpdateDatabaseRecordRequest{
			DatabaseID: databaseID,
			TableType:  req.GetTableType(),
			Records:    req.GetRecordDataAlter(),
			UserID:     *uid,
		})
		if err != nil {
			return nil, err
		}

		dataRes = append(dataRes, req.GetRecordDataAlter()...)
	}

	if len(req.GetRecordDataDelete()) > 0 {
		err := d.DomainSVC.DeleteDatabaseRecord(ctx, &database.DeleteDatabaseRecordRequest{
			DatabaseID: databaseID,
			TableType:  req.GetTableType(),
			Records:    req.GetRecordDataDelete(),
			UserID:     *uid,
		})
		if err != nil {
			return nil, err
		}

		dataRes = append(dataRes, req.GetRecordDataDelete()...)
	}

	return &table.UpdateDatabaseRecordsResponse{
		Data: dataRes,

		Code: 0,
		Msg:  "success",
		BaseResp: &base.BaseResp{
			StatusCode:    0,
			StatusMessage: "success",
		},
	}, nil
}

func (d *DatabaseApplicationService) GetOnlineDatabaseId(ctx context.Context, req *table.GetOnlineDatabaseIdRequest) (*table.GetOnlineDatabaseIdResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	basics := make([]*model.DatabaseBasic, 1)
	basics[0] = &model.DatabaseBasic{
		ID:        req.ID,
		TableType: table.TableType_DraftTable,
	}

	res, err := d.DomainSVC.MGetDatabase(ctx, &database.MGetDatabaseRequest{
		Basics: basics,
	})
	if err != nil {
		return nil, err
	}

	if len(res.Databases) == 0 {
		return nil, fmt.Errorf("database %d not found", req.ID)
	}

	if res.Databases[0].CreatorID != *uid {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "creator id is invalid"))
	}

	return &table.GetOnlineDatabaseIdResponse{
		ID: res.Databases[0].OnlineID,

		Code: 0,
		Msg:  "success",
	}, nil
}

func (d *DatabaseApplicationService) ResetBotTable(ctx context.Context, req *table.ResetBotTableRequest) (*table.ResetBotTableResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	tableType := table.TableType_OnlineTable
	if req.GetBotID() > 0 {
		tableType = table.TableType_DraftTable
	}
	err := d.ValidateAccess(ctx, req.GetDatabaseInfoID(), tableType)
	if err != nil {
		return nil, err
	}

	databaseID := req.GetDatabaseInfoID()
	if req.GetBotID() == 0 {
		databaseID, err = getDatabaseID(ctx, req.TableType, req.GetDatabaseInfoID())
		if err != nil {
			return nil, err
		}
	}

	executeDeleteReq := &database.ExecuteSQLRequest{
		DatabaseID:  databaseID,
		TableType:   req.GetTableType(),
		OperateType: model.OperateType_Delete,
		UserID:      conv.Int64ToStr(*uid),
		Condition: &model.ComplexCondition{
			Conditions: []*model.Condition{
				{
					Left:      model.DefaultIDColName,
					Operation: model.Operation_GREATER_THAN,
					Right:     "?",
				},
			},
			Logic: model.Logic_And,
		},
		SQLParams: []*model.SQLParamVal{
			{
				ValueType: table.FieldItemType_Number,
				Value:     ptr.Of("0"),
			},
		},
	}

	_, err = d.DomainSVC.ExecuteSQL(ctx, executeDeleteReq)
	if err != nil {
		return nil, err
	}

	return &table.ResetBotTableResponse{
		Code:     ptr.Of(int64(0)),
		Msg:      ptr.Of("success"),
		BaseResp: base.NewBaseResp(),
	}, nil
}

func (d *DatabaseApplicationService) GetDatabaseTemplate(ctx context.Context, req *table.GetDatabaseTemplateRequest) (*table.GetDatabaseTemplateResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	databaseID, err := getDatabaseID(ctx, req.TableType, req.DatabaseID)
	if err != nil {
		return nil, err
	}

	err = d.ValidateAccess(ctx, req.DatabaseID, table.TableType_OnlineTable)
	if err != nil {
		return nil, err
	}

	var fields []*model.FieldItem
	var tableName string
	if req.GetTableType() == table.TableType_DraftTable {
		basics := make([]*model.DatabaseBasic, 1)
		basics[0] = &model.DatabaseBasic{
			ID:        databaseID,
			TableType: table.TableType_DraftTable,
		}
		info, err := d.DomainSVC.MGetDatabase(ctx, &database.MGetDatabaseRequest{
			Basics: basics,
		})
		if err != nil {
			return nil, err
		}

		fields = info.Databases[0].FieldList
		tableName = info.Databases[0].TableName
	} else {
		basics := make([]*model.DatabaseBasic, 1)
		basics[0] = &model.DatabaseBasic{
			ID:        databaseID,
			TableType: table.TableType_OnlineTable,
		}
		info, err := d.DomainSVC.MGetDatabase(ctx, &database.MGetDatabaseRequest{
			Basics: basics,
		})
		if err != nil {
			return nil, err
		}

		fields = info.Databases[0].FieldList
		tableName = info.Databases[0].TableName
	}

	items := make([]*table.FieldItem, 0, len(fields))
	for _, field := range fields {
		items = append(items, &table.FieldItem{
			Name:         field.Name,
			Desc:         field.Desc,
			Type:         field.Type,
			MustRequired: field.MustRequired,
		})
	}

	resp, err := d.DomainSVC.GetDatabaseTemplate(ctx, &database.GetDatabaseTemplateRequest{
		UserID:     *uid,
		TableName:  tableName,
		FieldItems: items,
	})
	if err != nil {
		return nil, err
	}

	return &table.GetDatabaseTemplateResponse{
		TosUrl: resp.Url,

		Code: 0,
		Msg:  "success",
		BaseResp: &base.BaseResp{
			StatusCode:    0,
			StatusMessage: "success",
		},
	}, nil
}

func (d *DatabaseApplicationService) GetConnectorName(ctx context.Context, req *table.GetSpaceConnectorListRequest) (*table.GetSpaceConnectorListResponse, error) {
	return &table.GetSpaceConnectorListResponse{
		ConnectorList: []*table.ConnectorInfo{
			{
				ConnectorID:   consts.CozeConnectorID,
				ConnectorName: "Coze",
			},
			{
				ConnectorID:   consts.WebSDKConnectorID,
				ConnectorName: "Chat SDK",
			},
			{
				ConnectorID:   consts.APIConnectorID,
				ConnectorName: "API",
			},
		},

		Code: 0,
		Msg:  "success",
		BaseResp: &base.BaseResp{
			StatusCode:    0,
			StatusMessage: "success",
		},
	}, nil
}

func (d *DatabaseApplicationService) ValidateDatabaseTableSchema(ctx context.Context, req *table.ValidateTableSchemaRequest) (*table.ValidateTableSchemaResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	if req.GetSourceInfo() == nil || req.GetTableSheet() == nil {
		return nil, fmt.Errorf("source file and table sheet required")
	}

	databaseID, err := getDatabaseID(ctx, req.TableType, req.DatabaseID)
	if err != nil {
		return nil, err
	}

	err = d.ValidateAccess(ctx, req.DatabaseID, table.TableType_OnlineTable)
	if err != nil {
		return nil, err
	}

	basics := make([]*model.DatabaseBasic, 1)
	basics[0] = &model.DatabaseBasic{
		ID:        databaseID,
		TableType: req.TableType,
	}
	info, err := d.DomainSVC.MGetDatabase(ctx, &database.MGetDatabaseRequest{
		Basics: basics,
	})
	if err != nil {
		return nil, err
	}
	if len(info.Databases) == 0 {
		return nil, fmt.Errorf("database %d not found", req.DatabaseID)
	}

	res, err := d.DomainSVC.ValidateDatabaseTableSchema(ctx, &database.ValidateDatabaseTableSchemaRequest{
		DatabaseID: req.GetDatabaseID(),
		UserID:     *uid,
		TosURL:     req.GetSourceInfo().GetTosURI(),
		TableSheet: databaseEntity.TableSheet{
			SheetID:       req.GetTableSheet().GetSheetID(),
			HeaderLineIdx: req.GetTableSheet().GetHeaderLineIdx(),
			StartLineIdx:  req.GetTableSheet().GetStartLineIdx(),
		},
		Fields: info.Databases[0].FieldList,
	})
	if err != nil {
		return nil, err
	}

	if !res.Valid {
		return nil, errorx.New(errno.ErrMemoryInvalidParamCode,
			errorx.KV("msg", res.GetInvalidMsg()))
	}

	return &table.ValidateTableSchemaResponse{
		Code:     0,
		Msg:      "success",
		BaseResp: base.NewBaseResp(),
	}, nil
}

func (d *DatabaseApplicationService) GetDatabaseTableSchema(ctx context.Context, req *table.GetTableSchemaRequest) (*document.GetTableSchemaInfoResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	if req.GetSourceFile() == nil || req.GetTableSheet() == nil {
		return nil, fmt.Errorf("source file and table sheet required")
	}

	tableType := table.TableDataType_AllData
	if req.TableDataType != nil {
		tableType = req.GetTableDataType()
	}

	err := d.ValidateAccess(ctx, req.GetDatabaseID(), table.TableType_OnlineTable)
	if err != nil {
		return nil, err
	}

	schema, err := d.DomainSVC.GetDatabaseTableSchema(ctx, &database.GetDatabaseTableSchemaRequest{
		DatabaseID: req.GetDatabaseID(),
		UserID:     *uid,
		TosURL:     req.GetSourceFile().GetTosURI(),
		TableSheet: databaseEntity.TableSheet{
			SheetID:       req.GetTableSheet().GetSheetID(),
			HeaderLineIdx: req.GetTableSheet().GetHeaderLineIdx(),
			StartLineIdx:  req.GetTableSheet().GetStartLineIdx(),
		},
		// All data is returned by default without passing it on.
		TableDataType: tableType,
	})
	if err != nil {
		return nil, err
	}

	return &document.GetTableSchemaInfoResponse{
		TableMeta:   schema.TableMeta,
		SheetList:   schema.SheetList,
		PreviewData: schema.PreviewData,
		BaseResp:    base.NewBaseResp(),
	}, nil
}

func (d *DatabaseApplicationService) SubmitDatabaseInsertTask(ctx context.Context, req *table.SubmitDatabaseInsertRequest) (*table.SubmitDatabaseInsertResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	databaseID, err := getDatabaseID(ctx, req.TableType, req.DatabaseID)
	if err != nil {
		return nil, err
	}

	err = d.ValidateAccess(ctx, req.DatabaseID, table.TableType_OnlineTable)
	if err != nil {
		return nil, err
	}

	err = d.DomainSVC.SubmitDatabaseInsertTask(ctx, &database.SubmitDatabaseInsertTaskRequest{
		DatabaseID: databaseID,
		UserID:     *uid,
		FileURI:    req.GetFileURI(),
		TableSheet: databaseEntity.TableSheet{
			SheetID:       req.GetTableSheet().GetSheetID(),
			HeaderLineIdx: req.GetTableSheet().GetHeaderLineIdx(),
			StartLineIdx:  req.GetTableSheet().GetStartLineIdx(),
		},
		ConnectorID: req.ConnectorID,
		TableType:   req.TableType,
	})
	if err != nil {
		return nil, err
	}

	return &table.SubmitDatabaseInsertResponse{
		Code:     0,
		Msg:      "success",
		BaseResp: base.NewBaseResp(),
	}, nil
}

func (d *DatabaseApplicationService) DatabaseFileProgressData(ctx context.Context, req *table.GetDatabaseFileProgressRequest) (*table.GetDatabaseFileProgressResponse, error) {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return nil, errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session required"))
	}

	databaseID, err := getDatabaseID(ctx, req.TableType, req.DatabaseID)
	if err != nil {
		return nil, err
	}

	res, err := d.DomainSVC.GetDatabaseFileProgressData(ctx, &database.GetDatabaseFileProgressDataRequest{
		DatabaseID: databaseID,
		UserID:     *uid,
		TableType:  req.TableType,
	})
	if err != nil {
		return nil, err
	}

	return &table.GetDatabaseFileProgressResponse{
		Data: &table.DatabaseFileProgressData{
			FileName:       res.FileName,
			Progress:       res.Progress,
			StatusDescript: res.StatusDescript,
		},
		Code:     0,
		Msg:      "success",
		BaseResp: base.NewBaseResp(),
	}, nil
}

func getDatabaseID(ctx context.Context, tableType table.TableType, onlineID int64) (int64, error) {
	if tableType == table.TableType_OnlineTable {
		return onlineID, nil
	}

	online, err := DatabaseApplicationSVC.DomainSVC.MGetDatabase(ctx, &database.MGetDatabaseRequest{
		Basics: []*model.DatabaseBasic{
			{
				ID:        onlineID,
				TableType: table.TableType_OnlineTable,
			},
		},
	})
	if err != nil {
		return -1, err
	}
	if len(online.Databases) == 0 {
		return -1, fmt.Errorf("online table not found, id: %d", onlineID)
	}

	return online.Databases[0].GetDraftID(), nil
}

func (d *DatabaseApplicationService) ValidateAccess(ctx context.Context, databaseID int64, tableType table.TableType) error {
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "session uid not found"))
	}

	do, err := d.DomainSVC.MGetDatabase(ctx, &database.MGetDatabaseRequest{
		Basics: []*model.DatabaseBasic{
			{
				ID:        databaseID,
				TableType: tableType,
			},
		},
	})
	if err != nil {
		return err
	}
	if len(do.Databases) == 0 {
		return errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("msg", "database not found"))
	}

	if do.Databases[0].CreatorID != *uid {
		logs.CtxErrorf(ctx, "user(%d) is not the creator(%d) of the database(%d)", *uid, do.Databases[0].CreatorID, databaseID)
		return errorx.New(errno.ErrMemoryPermissionCode, errorx.KV("detail", "you are not the creator of the database"))
	}

	return nil
}

type CopyDatabaseRequest struct {
	DatabaseIDs []int64
	TableType   table.TableType // table type of the source databases
	CreatorID   int64

	IsCopyData    bool    // is need to copy data
	TargetSpaceID *int64  // if is nil, it will be set to the same space as the original database
	TargetAppID   int64   // if is nil, it will be set to the same app as the original database; if copy to resource, set to 0
	Suffix        *string // table name suffix for the copied table, default is "_copy"
}

type CopyDatabaseResponse struct {
	Databases map[int64]*entity.Database // key is original database id (online id or draft id), value is the new online database
}

func (d *DatabaseApplicationService) MoveDatabaseToLibrary(ctx context.Context, req *MoveDatabaseToLibraryRequest) (*MoveDatabaseToLibraryResponse, error) {
	basics := make([]*model.DatabaseBasic, 0, len(req.DatabaseIDs))
	for _, id := range req.DatabaseIDs {
		basics = append(basics, &model.DatabaseBasic{
			ID:        id,
			TableType: req.TableType,
		})
	}

	res, err := d.DomainSVC.MGetDatabase(ctx, &database.MGetDatabaseRequest{Basics: basics})
	if err != nil {
		return nil, err
	}

	moveDatabases := make([]*entity.Database, 0, len(req.DatabaseIDs))
	for _, srcDB := range res.Databases {

		srcDB.AppID = 0
		moveDatabaseResp, err := d.DomainSVC.UpdateDatabase(ctx, &database.UpdateDatabaseRequest{
			Database: srcDB,
		})
		if err != nil {
			return nil, err
		}

		onlineDatabase := moveDatabaseResp.Database
		moveDatabases = append(moveDatabases, onlineDatabase)
		// publish resource event

	}

	return &MoveDatabaseToLibraryResponse{
		Databases: moveDatabases,
	}, nil
}

type MoveDatabaseToLibraryRequest struct {
	DatabaseIDs []int64
	TableType   table.TableType // table type of the source databases
}

type MoveDatabaseToLibraryResponse struct {
	Databases []*entity.Database // the online databases after move
}
