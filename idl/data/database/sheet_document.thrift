include "sheet_common.thrift"
include "../../base.thrift"
namespace go data.database.sheet

// Basic information for creating a document


// Supports multiple data sources

// Get the table file meta information uploaded by the database

struct DocTableSheet {
    1: i64 id;            // Number of sheet
    2: string sheet_name; // Sheet name
    3: i64 total_row;     // total number of rows
}

// Determine whether the schema configured by the user is consistent with the corresponding document id

struct GetTableSchemaInfoResponse {
    1: i32 code
    2: string msg
    3: list<DocTableSheet>   sheet_list
    4: list<sheet_common.DocTableColumn>  table_meta                                        // The schema of the selected sheet, not selected to return the first sheet by default
    5: list<map<i64,string>> preview_data(agw.js_conv="str", agw.key="preview_data")  // The knowledge table will return

    255: optional base.BaseResp BaseResp(api.none="true")
}
