include "../base.thrift"
include "../app/bot_common.thrift"
namespace go playground

// branch

// Onboarding json structure


enum ScopeType {
    All  = 0 // All under the enterprise (effective under the enterprise)
    Self = 1 // I joined (both companies and individuals are valid, no default self is passed on)
}

struct GetSpaceListV2Request {
    1: optional string search_word                      // Search term
    2: optional i64 enterprise_id (api.js_conv='true',agw.js_conv="str") // Enterprise ID
    3: optional i64 organization_id (api.js_conv='true',agw.js_conv="str") // organization id
    4: optional ScopeType scope_type                    // range type
    5: optional i32 page                                // paging information
    6: optional i32 size                                // Paging size -- if page and size are not passed on, it is considered not paging

    255: optional base.Base Base (api.none="true")
}

enum SpaceType {
    Personal = 1 // individual
    Team     = 2 // group
}

enum SpaceMode {
    Normal = 0
    DevMode = 1
}

enum SpaceTag {
    Professional  =  1  // Professional Edition
}

enum SpaceRoleType {
    Default = 0 // default
    Owner   = 1 // owner
    Admin   = 2 // administrator
    Member  = 3 // ordinary member
}

// Application management list
enum SpaceApplyStatus {
    All              =  0     // all
    Joined           =  1     // Joined
    Confirming       =  2     // Confirming
    Rejected         =  3     // Rejected
}

struct AppIDInfo{
    1: string id
    2: string name
    3: string icon
}

struct ConnectorInfo{
    1: string id
    2: string name
    3: string icon
}

struct BotSpaceV2 {
    1: i64                 id (api.js_conv='true',agw.js_conv="str") // Space id, newly created as 0
    2: list<AppIDInfo>     app_ids        // publishing platform
    3: string              name           // space name
    4: string              description    // spatial description
    5: string              icon_url       // icon url
    6: SpaceType           space_type     // space type
    7: list<ConnectorInfo> connectors     // publishing platform
    8: bool                hide_operation // Whether to hide New, Copy Delete buttons
    9: i32                 role_type      // Role in team 1-owner 2-admin 3-member
    10: optional SpaceMode space_mode     // Spatial Mode
    11: bool               display_local_plugin // Whether to display the end-side plug-in creation entry
    12: SpaceRoleType      space_role_type // Role type, enumeration
    13: optional SpaceTag  space_tag       // spatial label
    14: optional i64    enterprise_id (api.js_conv='true',agw.js_conv="str") // Enterprise ID
    15: optional i64    organization_id (api.js_conv='true',agw.js_conv="str") // organization id
    16: optional i64    owner_user_id (api.js_conv='true',agw.js_conv="str") // Space owner uid
    17: optional string    owner_name      // Space owner nickname
    18: optional string    owner_user_name // Space owner username
    19: optional string    owner_icon_url  // Space owner image
    20: optional SpaceApplyStatus space_apply_status // The current visiting user joins the space status
    21: optional i64       total_member_num // The total number of space members, only the organization space can be queried.
}

struct SpaceInfo {
    1: list<BotSpaceV2> bot_space_list     // User joins space list
    2: bool             has_personal_space // Is there any personal space available?
    3: i32              team_space_num     // Number of team spaces created by individuals
    4: i32              max_team_space_num // The maximum number of spaces an individual can create
    5: list<BotSpaceV2> recently_used_space_list // list of recently used spaces
    6: optional i32 total                        // Effective when paging
    7: optional bool has_more                    // Effective when paging
}

struct GetSpaceListV2Response {
    1:   SpaceInfo       data
    253: required i64    code
    254: required string msg
    255: required base.BaseResp BaseResp
}

struct GetImagexShortUrlResponse{
    1             :      GetImagexShortUrlData data
    253: required i64    code
    254: required string msg
    255: required base.BaseResp BaseResp (api.none="true")
}

struct GetImagexShortUrlData {
    1: map<string,UrlInfo>   url_info //Audit status, key uri, value url and, audit status

}

struct UrlInfo {
    1: string url
    2: bool   review_status

}


enum GetImageScene {
    Onboarding = 0
    BackgroundImage = 1
}

struct GetImagexShortUrlRequest{
    1: list<string> uris
    2: GetImageScene scene

    255: base.Base Base (api.none="true"),
}

struct UserBasicInfo {
    1: required i64                  UserId         (api.js_conv='true',agw.js_conv="str", api.body="user_id")
    3: required string               Username       (api.body="user_name") // nickname
    4: required string               UserAvatar     (api.body="user_avatar") // avatar
    5: optional string               UserUniqueName (api.body="user_unique_name") // user name
    6: optional bot_common.UserLabel UserLabel      (api.body="user_label") // user tag
    7: optional i64                  CreateTime     (api.body="create_time") // user creation time
}

struct MGetUserBasicInfoRequest {
    1 : required list<string> UserIds (agw.js_conv="str", api.js_conv="true", api.body="user_ids")
    2 : optional bool NeedUserStatus (api.body="need_user_status")
    3 : optional bool NeedEnterpriseIdentity (api.body="need_enterprise_identity") // Whether enterprise authentication information is required, the default is true when the front end is called through AGW
    4 : optional bool NeedVolcanoUserName (api.body="need_volcano_user_name") // Do you need a volcano username?

    255: optional base.Base Base (api.none="true")
}

struct MGetUserBasicInfoResponse {
    1 : optional map<string, UserBasicInfo> UserBasicInfoMap (api.body="id_user_info_map")

    253:          i64           code
    254:          string        msg

    255: optional base.BaseResp BaseResp (api.none="true")
}

struct FileInfo {
    1 : string url
    2 : string uri
}
enum GetFileUrlsScene {
    shorcutIcon = 1
}
struct GetFileUrlsRequest {
    1 :  GetFileUrlsScene scene
    255: base.Base Base
}

struct GetFileUrlsResponse {
    1 : list<FileInfo> file_list
    253: i64 code
    254: string msg
    255: base.BaseResp BaseResp
}

service PlaygroundService {
    GetImagexShortUrlResponse GetImagexShortUrl (1:GetImagexShortUrlRequest request)(api.post='/api/playground_api/get_imagex_url', api.category="file",agw.preserve_base="true")

    // public popup_info

   // Create shortcut instructions


    // prompt resource

    GetSpaceListV2Response GetSpaceListV2(1:GetSpaceListV2Request request)(api.post='/api/playground_api/space/list', api.category="space",agw.preserve_base="true")
    MGetUserBasicInfoResponse MGetUserBasicInfo(1: MGetUserBasicInfoRequest request) (api.post='/api/playground_api/mget_user_info', api.category="playground_api",agw.preserve_base="true")

}
