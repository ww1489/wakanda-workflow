include "bot_common.thrift"
namespace go app.developer_api

// draft bot display info

struct Connector {
    1:          string             name       // connector_name enumeration Feishu = "feishu"
    2:          string             app_id
    3:          string             app_secret
    4:          string             share_link
    5: optional map<string,string> bind_info
}

// Information about each module in the agent workspace

// Information for each module in the workspace

// If historical information is preserved


struct ConnectorInfo {
    1:          string                 id
    2:          string                 name
    3:          string                 icon
    4:          ConnectorDynamicStatus connector_status
    5: optional string                 share_link
}


enum ConnectorDynamicStatus {
    Normal          = 0
    Offline         = 1
    TokenDisconnect = 2
}

enum IconType {
    Bot       = 1
    User      = 2
    Plugin    = 3
    Dataset   = 4
    Space     = 5
    Workflow  = 6
    Imageflow = 7
    Society   = 8
    Connector = 9
    ChatFlow = 10
    Voice = 11
    Enterprise = 12
}

struct GetIconRequest {
    1: IconType icon_type
}

struct Icon {
    1: string url
    2: string uri
}

struct GetIconResponseData {
    1: list<Icon> icon_list
}
struct GetIconResponse {
    1: i64                 code
    2: string              msg
    3: GetIconResponseData data
}
struct GetUploadAuthTokenResponse {
    1: i64                    code
    2: string                 msg
    3: GetUploadAuthTokenData data
}

struct GetUploadAuthTokenData {
    1: string              service_id
    2: string              upload_path_prefix
    3: UploadAuthTokenInfo auth
    4: string              upload_host
    5: string              schema
}
struct UploadAuthTokenInfo {
    1: string access_key_id
    2: string secret_access_key
    3: string session_token
    4: string expired_time
    5: string current_time
}

struct GetUploadAuthTokenRequest {
    1: string scene
    2: string data_type
}

struct UploadFileRequest {
    1: CommonFileInfo file_head // Document related description
    2: string         data      // file data
}

// Upload file, file header
struct CommonFileInfo {
    1: string      file_type // File type, suffix
    2: FileBizType biz_type  // business type
}

enum FileBizType {
    BIZ_UNKNOWN      = 0
    BIZ_BOT_ICON     = 1
    BIZ_BOT_DATASET  = 2
    BIZ_DATASET_ICON = 3
    BIZ_PLUGIN_ICON  = 4
    BIZ_BOT_SPACE    = 5
    BIZ_BOT_WORKFLOW = 6
    BIZ_SOCIETY_ICON = 7
    BIZ_CONNECTOR_ICON = 8
    BIZ_LIBRARY_VOICE_ICON = 9
    BIZ_ENTERPRISE_ICON = 10
}

struct UploadFileResponse {
    1: i64            code
    2: string         msg
    3: UploadFileData data // data
}

struct GetTypeListRequest {
    1: optional bool   model
    2: optional bool   voice
    3: optional bool   raw_model
    4: optional string space_id
    5: optional string cur_model_id // The model ID used by the current bot to handle issues that cannot be displayed by the bot model synchronized by cici/doubao
    6: optional list<string> cur_model_ids // Compatible with MultiAgent, with multiple cur_model_id
    7: optional ModelScene model_scene // model scenario
}

enum ModelScene {
    Douyin = 1
}

enum ModelClass {
    GPT             = 1
    SEED            = 2
    Claude          = 3
    MiniMax         = 4  // name: MiniMax
    Plugin          = 5
    StableDiffusion = 6
    ByteArtist      = 7
    Maas            = 9
    QianFan         = 10 // Abandoned: Qianfan (Baidu Cloud)
    Gemini          = 11 // name：Google Gemini
    Moonshot        = 12 // name: Moonshot
    GLM             = 13 // Name: Zhipu
    MaaSAutoSync    = 14 // Name: Volcano Ark
    QWen            = 15 // Name: Tongyi Qianwen
    Cohere          = 16 // name: Cohere
    Baichuan        = 17 // Name: Baichuan Intelligent
    Ernie           = 18 // Name: ERNIE Bot
    DeekSeek        = 19 // Name: Magic Square
    Llama           = 20 // name: Llama
    StepFun         = 23
    Other           = 999
}

struct ModelQuota {
    1:           i32    token_limit         // Maximum total number of tokens
    2:           i32    token_resp          // Final reply maximum number of tokens
    3:           i32    token_system        // Prompt system maximum number of tokens
    4:           i32    token_user_in       // Prompt user to enter maximum number of tokens
    5:           i32    token_tools_in      // Prompt tool to enter maximum number of tokens
    6:           i32    token_tools_out     // Prompt tool output maximum number of tokens
    7:           i32    token_data          // Prompt data maximum number of tokens
    8:           i32    token_history       // Prompt history maximum number of tokens
    9:           bool   token_cut_switch    // Prompt history maximum number of tokens
    10:          double price_in            // input cost
    11:          double price_out           // output cost
    12: optional i32    system_prompt_limit // Systemprompt input restrictions, if not passed, no input restrictions
}

enum ModelTagClass {
    ModelType = 1
    ModelUserRight = 2
    ModelFeature = 3
    ModelFunction = 4

    Custom = 20 // Do not do this issue
    Others = 100
}

enum ModelParamType {
    Float = 1
    Int = 2
    Boolean = 3
    String = 4
}

struct ModelParamDefaultValue {
    1: required string default_val
    2: optional string creative
    3: optional string balance
    4: optional string precise
}

struct ModelParamClass {
    1: i32    class_id // 1="Generation diversity", 2="Input and output length", 3="Output format"
    2: string label
}

struct Option {
    1: string label // The value displayed by the option
    2: string value // Filled in value
}

struct ModelParameter {
    1: required string name // Configuration fields, such as max_tokens
    2: string label // Configure field display name
    3: string desc // Configuration field detail description
    4: required ModelParamType type // type
    5: string min // Numerical type parameters, the minimum value allowed to be set
    6: string max // Numerical type parameter, the maximum value allowed to be set
    7: i32 precision // Precision of float type parameters
    8: required ModelParamDefaultValue default_val // Parameter default {"default": xx, "creative": xx}
    9: list<Option> options // Enumeration values such as response_format support text, markdown, json
    10: ModelParamClass param_class // Parameter classification, "Generation diversity", "Input and output length", "Output format"
}

struct ModelDescGroup {
    1: string group_name
    2: list<string> desc
}

struct ModelTag {
    1: string tag_name
    2: ModelTagClass tag_class
    3: string tag_icon
    4: string tag_descriptions
}

struct ModelSeriesInfo {
    1: string series_name,
    2: string icon_url,
    3: string model_vendor,
    4: optional string model_tips,
}

enum ModelTagValue {
    Flagship = 1,
    HighSpeed = 2,
    ToolInvocation = 3,
    RolePlaying = 4,
    LongText = 5,
    ImageUnderstanding = 6,
    Reasoning = 7,
    VideoUnderstanding = 8,
    CostPerformance = 9,
    CodeSpecialization = 10,
    AudioUnderstanding = 11
}

struct ModelStatusDetails {
    1: bool is_new_model, // Is it a new model?
    2: bool is_advanced_model, // Is it a high-level model?
    3: bool is_free_model, // Is it a free model?

    11: bool is_upcoming_deprecated, // Will it be removed from the shelves soon?
    12: string deprecated_date, // removal date
    13: string replace_model_name, // Remove the replacement model from the shelves.

    21: string update_info, // Recently updated information
    22: ModelTagValue model_feature, // Model Features
}

struct ModelAbility {
    1: optional bool cot_display // Do you want to show cot?
    2: optional bool function_call // Supports function calls
    3: optional bool image_understanding // Does it support picture understanding?
    4: optional bool video_understanding // Does it support video understanding?
    5: optional bool audio_understanding // Does it support audio understanding?
    6: optional bool support_multi_modal // Does it support multimodality?
    7: optional bool prefill_resp // Whether to support continuation
}

struct Model {
    1: string     name
    2: i64        model_type
    3: ModelClass model_class
    4: string     model_icon         // Model icon url
    5: double     model_input_price
    6: double     model_output_price
    7: ModelQuota model_quota
    8: string     model_name         // Model real name, front-end calculation token
    9: string     model_class_name
    10: bool      is_offline
    11: list<ModelParameter> model_params
    12: optional list<ModelDescGroup>    model_desc
    13: optional map<bot_common.ModelFuncConfigType, bot_common.ModelFuncConfigStatus> func_config, // model function configuration
    14: optional string endpoint_name  // Ark model node name
    15: optional list<ModelTag> model_tag_list  // model label
    16: optional bool is_up_required   // User prompt must have and cannot be empty
    17: string model_brief_desc // Model brief description
    18: ModelSeriesInfo model_series // Model series
    19: ModelStatusDetails model_status_details // model state
    20: ModelAbility model_ability // model capability
}

struct VoiceType {
    1: i64    id
    2: string model_name
    3: string name
    4: string language
    5: string style_id
    6: string style_name
}

struct GetTypeListData {
    1: list<Model>     model_list
    2: list<VoiceType> voice_list
    3: list<Model>     raw_model_list
}

struct GetTypeListResponse {
    1:          i64             code
    2:          string          msg
    3: required GetTypeListData data
}

struct UploadFileData {
    1: string upload_url // File URL
    2: string upload_uri // File URI, submit using this
}

struct UpdateUserProfileCheckRequest {
    1: optional string user_unique_name
}

struct UpdateUserProfileCheckResponse {
    1: i64    code
    2: string msg
}

// Check if the draft can be submitted and returned.

// Check if the draft can be submitted to the request

struct PluginOauthAuthorizationCodeReq {
    1: string code  (api.query='code') ,
    2: string state (api.query='state'),
    3: string plugin_id (api.path='plugin_id'),
}

struct PluginOauthAuthorizationCodeResp {
    1: i64    code
    2: string msg
}

struct PluginOauthInfo {
    1: string connector_id,
    2: string plugin_id,
    3: string plugin_name,
    4: string plugin_icon,
    5: string username,
    6: string plugin_url,
    7: string connector_name,
}

struct PluginOauthInfoReq {
    1: string confirm_code (api.query='confirm_code')
}

struct PluginOauthInfoResp {
    1: i64    code
    2: string msg
    3: PluginOauthInfo data
}

struct PluginOauthConfirmReq {
    1: string confirm_code (api.post='confirm_code')
}

struct PluginOauthConfirmResp {
    1: i64    code
    2: string msg
}

service DeveloperApiService {
    PluginOauthAuthorizationCodeResp PluginOauthAuthorizationCode(1: PluginOauthAuthorizationCodeReq request)(api.get='/api/plugin_oauth/:plugin_id/authorization_code', api.category="oauth", api.gen_path="oauth")
    PluginOauthInfoResp PluginOauthInfo(1: PluginOauthInfoReq request)(api.get='/api/plugin/oauth/get_oauth_info', api.category="oauth", api.gen_path="oauth")
    PluginOauthConfirmResp PluginOauthConfirm(1: PluginOauthConfirmReq request)(api.post='/api/plugin/oauth/confirm', api.category="oauth", api.gen_path="oauth")

    GetUploadAuthTokenResponse GetUploadAuthToken(1: GetUploadAuthTokenRequest request)(api.post = '/api/playground/upload/auth_token', api.category="playground", api.gen_path="playground")

    UploadFileResponse UploadFile(1:UploadFileRequest request)(api.post='/api/bot/upload_file', api.category="bot" api.gen_path="bot")
    GetTypeListResponse GetTypeList(1: GetTypeListRequest request)(api.post='/api/bot/get_type_list', api.category="bot", api.gen_path="bot")

    GetIconResponse GetIcon(1:GetIconRequest request)(api.post='/api/developer/get_icon', api.category="developer", api.gen_path="developer")

    UpdateUserProfileCheckResponse UpdateUserProfileCheck(1: UpdateUserProfileCheckRequest request)(api.post='/api/user/update_profile_check', api.category="user", api.gen_path="user")

}
