include "../../base.thrift"
namespace go data.database.table

enum BotTableStatus {
    Online = 1 // online
    Delete = 2 // delete
}

enum FieldItemType {
    Text   = 1 // Text
    Number = 2 // number
    Date   = 3 // time
    Float   = 4 // float
    Boolean = 5 // bool
}

struct FieldItem {
    1: string        name
    2: string        desc
    3: FieldItemType type
    4: bool          must_required
    5: i64           id            // This field is only used to determine whether it is published. If it is not 0, it is published. The front end cannot modify the field type for the published field.
    6: i64           alterId       // When modifying a field (alter, publish), it is used to judge additions and deletions. 0 means new additions, and non-0 means modification or deletion.
    7: bool          is_system_field // Is it a system field?
}

enum BotTableRWMode {
    LimitedReadWrite    = 1     // single user mode
    ReadOnly            = 2     // read-only mode
    UnlimitedReadWrite  = 3     // multi-user mode
    RWModeMax           = 4     // Max boundary value
}

struct OrderBy {
    1: string field
    2: SortDirection direction
}

enum SortDirection {
    ASC = 1
    Desc = 2
}

struct Criterion{
    1: list<Condition> conditions
    2: string logic_expression
}

struct ListDatabaseRequest {
    1:   optional i64        creator_id  (api.js_conv="str") // Get a database created by a user
    2:   optional i64        project_id (api.js_conv="str") // Get the database under the project
    3:   optional i64        space_id (api.js_conv="str") //Get the visible database under space
    4:   optional i64        bot_id (api.js_conv="str")  //Filter bot_id to filter out databases that have been added to the bot
    5:   optional string     table_name // Table name, fuzzy search
    6:   required TableType  table_type // Draft database
    7:   optional list<OrderBy> order_by // sort
    8:   optional i32 offset
    9:   optional i32 limit
    10: optional Criterion filter_criterion  //filter criteria
    11: optional list<OrderBy> order_by_list   //sort condition
    255: optional base.Base  Base
}

struct DatabaseInfo {
    1:  i64             id          (api.js_conv="str", api.key="id") // online_database_info primary key id
    2:  i64             space_id    (api.js_conv="str") // ID of space
    3:  i64             project_id  (api.js_conv="str") // project id
    4:  string          datamodel_table_id    // Table ID on the datamodel side
    5:  string          icon_url    // avatar url
    6:  string          icon_uri    // avatar url
    7:  string          table_name  // table name
    8:  string          table_desc  // table description
    9:  BotTableStatus  status      // status
    10: i64             creator_id  (api.js_conv="str") // creator id
    11: i64             create_time // create_time
    12: i64             update_time // update time
    13: list<FieldItem> field_list  // Field information
    14: string          actual_table_name    // Data table actual name
    15: BotTableRWMode  rw_mode     // read and write mode
    16: bool            prompt_disabled  // Whether to support prompt calls
    17: bool            is_visible   // Is it visible?
    18: optional i64    draft_id     (api.js_conv="str")   // ID corresponding to draft state
    19: optional i64    bot_id       (api.js_conv="str", api.key="bot_id") // Related id. bot_id, the old one is available, the new one is not.
    20: optional map<string,string> extra_info // extended information
    21: optional bool   is_added_to_bot // Has it been added to the bot?
}

struct ListDatabaseResponse{
    1:   list<DatabaseInfo>  database_info_list
    2:   bool  has_more
    3:   i64   total_count

    253: required i64 code
    254: required string msg
    255: required base.BaseResp BaseResp
}

struct SingleDatabaseRequest{
    1:  i64   id         (api.js_conv="str", api.key="id") // database_info primary key id
    2:  bool  is_draft   (api.key="is_draft") //Whether the incoming data is in draft form, the default is false.
    3:  bool  need_sys_fields (api.key="need_sys_fields") //Do you need system fields?
    4:  i64   version    (api.js_conv="str") // The version number is not passed on, and the default is the latest.
    255: optional base.Base  Base
}

struct SingleDatabaseResponse{
    1:  DatabaseInfo  database_info

    253: required i64 code
    254: required string msg
    255: required base.BaseResp BaseResp
}

struct AddDatabaseRequest{
    1:  i64             creator_id  (api.js_conv="str") // creator id
    2:  i64             space_id    (api.js_conv="str") // ID of space
    3:  i64             project_id  (api.js_conv="str") // project id
    4:  string          icon_uri    // avatar url
    5:  string          table_name  // table name
    6:  string          table_desc  // table description
    7:  list<FieldItem> field_list  // Field information
    8:  BotTableRWMode  rw_mode     // Read and write mode, single user mode/multi-user mode
    9:  bool            prompt_disabled  // Whether to support prompt calls
    10:  optional map<string,string> extra_info // extended information
    255: optional base.Base  Base
}

struct UpdateDatabaseRequest{
    1:   i64             id     (api.js_conv="str", api.key="id") // database_info primary key id
    2:   string          icon_uri    // avatar url
    3:   string          table_name  // table name
    5:   string          table_desc  // table description
    6:   list<FieldItem> field_list  // Field information
    7:   BotTableRWMode  rw_mode     // Read and write mode, single user mode/multi-user mode
    8:   bool            prompt_disabled  // Whether to support prompt calls
    9:   optional map<string,string> extra_info // extended information
    255: optional base.Base  Base
}

struct DeleteDatabaseRequest{
    1:   i64     id     (api.js_conv="str", api.key="id") // database_info primary key id
    255: optional base.Base  Base
}

struct DeleteDatabaseResponse {
    1: required i64 code
    2: required string msg
    255: optional base.BaseResp BaseResp
}

struct ListDatabaseRecordsRequest{
    1: required i64    database_id (api.js_conv="str") // database_id
    2: optional i64    bot_id (api.js_conv="str", api.key="bot_id") // bot id, here is to fill in this when looking for the draft state data associated with the bot
    3: optional i64    workflow_id (api.js_conv="str", api.key="workflow_id") // workflow_id,, here is to fill in this when looking up wk_flow associated draft status sheet
    4: optional bool not_filter_by_user_id                     // Is true does not filter by user_id Records
    5: optional bool not_filter_by_connector_id               // Records not filtered by ConnectorID
    6: TableType table_type    // Do you want to check the draft state or the online state?
    7: i64    limit            // Do not exceed 100, 50 is recommended.
    8: i64    offset           // Offset
    9: optional i64 project_id (api.js_conv="str") // Data is not isolated under the same project
    10: optional ComplexCondition  filter_criterion  //filter criteria
    11: optional list<OrderBy> order_by_list   //sort condition
    255: optional base.Base  Base
}

struct ListDatabaseRecordsResponse{
    1: required list<map<string,string>> data
    2: required bool                     HasMore=false
    3: required i32                      TotalNum
    4: optional list<FieldItem> field_list  // Field information

    253: required i64 code
    254: required string msg
    255: required base.BaseResp BaseResp
}

struct UpdateDatabaseRecordsRequest{
    1: required i64 database_id (api.js_conv="str") // database_id
    2: optional list<map<string,string>> record_data_add // new
    3: optional list<map<string,string>> record_data_alter // modified
    4: optional list<map<string,string>> record_data_delete // deleted
    5: optional TableType table_type    // Is the draft state or online state to be updated?
    6: optional string    ori_connector_id // The connector id needs to be inserted when updating.
    255: optional base.Base  Base
}

struct UpdateDatabaseRecordsResponse{
    1: required list<map<string,string>> data

    253: required i64 code
    254: required string msg
    255: required base.BaseResp BaseResp
}
struct GetOnlineDatabaseIdRequest{
    1: required i64 id (api.js_conv="str") // The draft database_id
    255: optional base.Base  Base
}

struct GetOnlineDatabaseIdResponse{
    1: optional i64 id   (api.js_conv="str")   // Check the online id according to the id of the draft.

    253: required i64 code
    254: required string msg
    255: required base.BaseResp BaseResp
}

struct ResetBotTableRequest {
    1: optional i64 creator_id         (api.js_conv="str", api.key="creator_id")
    2: optional i64 bot_id             (api.js_conv="str", api.key="bot_id")
    3: optional i64 table_id           (api.js_conv="str", api.key="table_id")
    4: required TableType table_type
    5: optional i64    connector_id  // line of business id
    6: optional string connector_uid // line of business user id
    7: optional i64    workflow_id (api.js_conv="str")   // Workflow ID
    8: optional i64    database_info_id (api.js_conv="str")   // user id
    9: optional i64    project_id (api.js_conv="str")   // Project ID
    255: optional base.Base Base
}

struct ResetBotTableResponse {
    253: optional i64 code
    254: optional string msg
    255: required base.BaseResp BaseResp(api.none="true")
}


/********  bot_table end     ********/

/********      Bytedoc  bot_table_info  start   ********/

enum TableType {
    DraftTable  = 1 // draft
    OnlineTable = 2 // online
}

/********     Bytedoc  bot_table_info  end  ********/

// Library Resource Operations

// For display, the implementer provides display information

enum Operation {
  EQUAL = 1,          // "="
  NOT_EQUAL = 2,      // "< >" or "! ="
  GREATER_THAN = 3,   // ">"
  LESS_THAN = 4,      // "<"
  GREATER_EQUAL = 5,  // ">="
  LESS_EQUAL = 6,     // "<="
  IN = 7,             // "IN"
  NOT_IN = 8,         // "NOT IN"
  IS_NULL = 9,        // "IS NULL"
  IS_NOT_NULL = 10    // "IS NOT NULL"
  LIKE = 11,           // "LIKE" fuzzy match string
  NOT_LIKE = 12,       // "NOT LIKE" inverse fuzzy match
}

struct Condition {
  1: required string left; // Lvalue field name
  2: required Operation operation;
  3: required string right; // rvalue
}

struct ComplexCondition {
  1: optional list<Condition> conditions;
  2: optional ComplexCondition nestedConditions; // In order to expand, we don't need to
  3: required string logic; // "AND" or "OR"
}

struct SourceInfo {

    // TOS address for local file upload
    1: optional string tos_uri (api.key="tos_uri");
    // imagex_uri, and tos_uri choose one, imagex_uri priority, need to get data and sign url through imagex method
    2: optional string imagex_uri
}


struct ValidateTableSchemaRequest {
    1: i64 space_id           (api.js_conv="str", api.key="space_id")
    2: i64 database_id        (api.js_conv="str", api.key="database_id")
    3: SourceInfo source_info (api.key="source_file", api.body="source_file")               // Information from the source file
    4: TableSheet table_sheet (api.key="table_sheet")
    5: TableType table_type (api.key="table_type")
    255: optional base.Base Base
}

struct TableSheet {
    1: i64 sheet_id        (api.js_conv="str", api.key="sheet_id")       , // User selected sheet id
    2: i64 header_line_idx (api.js_conv="str", api.key="header_line_idx"), // The number of header rows selected by the user, numbered from 0
    3: i64 start_line_idx  (api.js_conv="str", api.key="start_line_idx") , // User-selected starting line number, numbered from 0
}

struct ValidateTableSchemaResponse {
    1: optional map<string,string> SchemaValidResult (api.key="schema_valid_result");
    // If it fails, an error code will be returned.
    253: required i64 code
    254: required string msg
    255: optional base.BaseResp BaseResp(api.none="true")
}

struct GetTableSchemaRequest {
   1: optional TableSheet  table_sheet;                                                         // Table parsing information, default initial value 0, 0, 1
   2: optional TableDataType table_data_type;                                          // All data is returned by default without passing it on.
   3: optional i64 database_id(api.js_conv="str", api.key="database_id");              // Compatible with pre-refactoring versions: pass this value if you need to pull the schema of the current document
   4: optional SourceInfo source_file;                                                 // Source file information, add segment/before logic migrate here
   255: optional base.Base Base
}

enum TableDataType {
    AllData     = 0     // Schema sheets and preview data
    OnlySchema  = 1     // Only need schema structure & Sheets
    OnlyPreview = 2    // Just preview the data
}

struct GetDatabaseFileProgressRequest {
    1: i64 database_id (api.js_conv="str")
    2: required TableType table_type    // Table type
    255: optional base.Base Base
}
struct GetDatabaseFileProgressResponse {
    1: DatabaseFileProgressData data
    253: required i64 code
    254: required string msg
    255: optional base.BaseResp BaseResp
}

struct DatabaseFileProgressData {
    1: string              file_name
    2: i32                 progress
    3: optional string     status_descript  //Describe the information if there is a representative file processing failure
}

struct SubmitDatabaseInsertRequest {
    1: i64 database_id (api.js_conv="str")
    2: string file_uri
    3: TableType table_type    // Table type, do you want to insert into the draft table or the online table?
    4: optional TableSheet  table_sheet
    5: optional i64 connector_id  (api.js_conv="str") // Channel ID to write to
    255: optional base.Base Base
}

struct SubmitDatabaseInsertResponse {
    253: required i64 code
    254: required string msg
    255: optional base.BaseResp BaseResp
}

struct GetSpaceConnectorListRequest {
    1: required i64 SpaceId (api.js_conv="str")
    2: optional string Version // release inhouse
    3: optional i64 ConnectorID
    4: optional bool ListAll
    255: optional base.Base Base
}

struct GetSpaceConnectorListResponse {
    1: list<ConnectorInfo> ConnectorList

    253: required i64 code
    254: required string msg
    255: optional base.BaseResp BaseResp
}

struct ConnectorInfo{
    1: i64 ConnectorID
    2: string ConnectorName
}

typedef  GetDatabaseFileProgressRequest GetDatabaseTemplateRequest

struct GetDatabaseTemplateResponse {
    1: string TosUrl // Download address

    253: required i64 code
    254: required string msg
    255: optional base.BaseResp BaseResp
}
