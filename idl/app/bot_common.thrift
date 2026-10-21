namespace go app.bot_common


struct AuthToken {
    1: string service_id,
    2: string access_key_id,
    3: string secret_access_key,
    4: string session_token,
    5: string expired_time,
    6: string current_time,
    7: string upload_host,
    8: string host_scheme,
}

struct ModelInfo {
    1: optional i64                 ModelId           (agw.js_conv="str", api.js_conv="true", api.body="model_id"), // Model ID
    2: optional double              Temperature       (api.body="temperature")                                    , // Temperature, model output randomness, the larger the value, the more random, the smaller the more conservative (0-1]
    3: optional i32                 MaxTokens         (api.body="max_tokens")                                     , // Maximum Token Reply
    4: optional double              TopP              (api.body="top_p")                                          , // Another model's output randomness, the larger the value, the more random [0, 1]
    5: optional double              FrequencyPenalty  (api.body="frequency_penalty")                              , // Frequency penalty, adjust the frequency of words in the generated content, the fewer positive words are [-1.0, 1.0]
    6: optional double              PresencePenalty   (api.body="presence_penalty")                               , // There is a penalty, adjust the frequency of new words in the generated content, avoid repeating words with positive values, and use new words [-1.0, 1.0]
    7: optional ShortMemoryPolicy   ShortMemoryPolicy (api.body="short_memory_policy")                            , // contextual policy
    8: optional i32                 TopK              (api.body="top_k")                                          , // When generating, sample the size of the candidate set
    9: optional ModelResponseFormat ResponseFormat    (api.body="response_format")                                , // model reply content format
    10: optional ModelStyle         ModelStyle        (api.body="model_style")
    11: optional CacheType          CacheType           (api.body="cache_type")                                   , // 缓存配置
    12: optional bool               SpCurrentTime       (api.body="sp_current_time")                                , // sp拼接当前时间
    13: optional bool               SpAntiLeak          (api.body="sp_anti_leak")                                   , // sp拼接防泄露指令
    14: optional bool               SpVoiceInfo         (api.body="sp_voice_info")                                  , // sp拼接声纹信息
    15: optional map<string, string> Parameters         (api.body="parameters")                                   , // User-selected model style
}

enum CacheType {
    CacheClosed = 0 // 缓存关闭
    PrefixCache = 1 // 前缀缓存
}

enum ModelStyle {
    Custom = 0
    Creative = 1
    Balance = 2
    Precise = 3
}
enum ModelResponseFormat {
    Text = 0,
    Markdown = 1,
    JSON = 2,
}

// The type of transmission allowed by the context
enum ContextMode {
    Chat           = 0,
    FunctionCall_1 = 1,
    FunctionCall_2 = 2,
    FunctionCall_3 = 3,
}

enum ModelFuncConfigType {
    Plugin = 1
    Workflow = 2
    ImageFlow = 3
    Trigger = 4
    KnowledgeText = 5
    KnowledgeTable = 6
    KnowledgeAutoCall = 7
    KnowledgeOnDemandCall = 8
    Variable = 9
    Database = 10
    LongTermMemory = 11
    FileBox = 12
    Onboarding = 13
    Suggestion = 14
    ShortcutCommand = 15
    BackGroundImage = 16
    TTS = 17
    MultiAgentRecognize = 18
    KnowledgePhoto = 19
    HookInfo = 20
}

enum ModelFuncConfigStatus {
    FullSupport = 0
    PoorSupport = 1
    NotSupport = 2
}

struct ShortMemoryPolicy {
    1: optional ContextMode ContextMode  (api.body="context_mode") , // The type of transmission allowed by the context
    2: optional i32         HistoryRound (api.body="history_round"), // Number of rounds of context band
}

enum PluginFrom {
    Default = 0
    FromSaas = 1
}

// onboarding content generation mode

// suggest


// tts Voices

// AnswerActions

// bot ext

enum SocietyVisibility {
    Public = 1, // Visible to all
    Anonymous = 2, // Visible to host only
    Custom = 3, // custom
}
struct SocietyVisibiltyConfig {
    1: SocietyVisibility VisibilityType (api.body="visibility_type", go.tag="json:\"visibility_type,omitempty\"") , // Visibility in Social Scene: Public = 1, Anonymous = 2
    2: list<string> VisibilityRoles     (api.body="visibility_roles", go.tag="json:\"visibility_roles,omitempty\""), // list of visible characters
}

struct Variable {
    1: optional string Key          (api.body="key")          , // key, Field
    2: optional string Description  (api.body="description")  , // describe
    3: optional string DefaultValue (api.body="default_value"), // default value
    4: optional bool   IsSystem     (api.body="is_system"),     // Whether the system value is the system value
    5: optional bool   PromptDisabled (api.body="prompt_disabled"), // Whether to support calling in Prompt, the default is supported
    6: optional SocietyVisibiltyConfig SocietyVisibilityConfig (api.body="society_visibility_config", go.tag="json:\"society_visibility_config,omitempty\""), // Visibility in Social Scene: Public = 1, Anonymous = 2
    7: optional bool   IsDisabled (api.body="is_disabled"),  // Whether to disable, the default is false to enable
}

enum FieldItemType {
    Text    = 1, // Text String
    Number  = 2, // Digital Integer
    Date    = 3, // Time Time
    Float   = 4, // float Number
    Boolean = 5, // bool Boolean
}

struct FieldItem {
    1: optional string        Name         (api.body="name")                                     , // field name
    2: optional string        Desc         (api.body="desc")                                     , // Field description
    3: optional FieldItemType Type         (api.body="type")                                     , // field type
    4: optional bool          MustRequired (api.body="must_required")                            , // Is it required?
    5: optional i64           Id           (agw.js_conv="str", api.js_conv="true", api.body="id"), // The field Id is added as 0.
    6: optional string        TypeStr      (api.body="type_str")                                 , // Field type str
    7: optional i64           AlterId      (api.body="alterId")                                 , // Field type str
}

struct Database {
    1: optional string          TableId   (api.body="table_id")  , // table id
    2: optional string          TableName (api.body="table_name"), // Table name
    3: optional string          TableDesc (api.body="table_desc"), // Table Introduction
    4: optional list<FieldItem> FieldList (api.body="field_list"), // Table field information
    5: optional bool            PromptDisabled (api.body="prompt_disabled"), // Whether to support calling in Prompt, the default is supported
    6: optional BotTableRWMode  RWMode    (api.body="rw_mode"),
}

enum BotTableRWMode {
    LimitedReadWrite = 1,
    ReadOnly = 2,
    UnlimitedReadWrite = 3,
    RWModeMax = 4,
}

//Version Compatibility: 0 - Old Version 1 - Rollback New Version 2 - Non-Rollback New Version 3 - Rollback New Version (no longer prompted)

// Time Capsule Information

// bot information

// Update for bot information

struct UserLabel {
    1: string             label_id
    2: string             label_name
    3: string             icon_uri
    4: string             icon_url
    5: string             jump_link
}

struct HookInfo {
    1: optional list<HookItem> pre_agent_jump_hook // Pre agent jump hook
    2: optional list<HookItem> post_agent_jump_hook // Post agent jump hook
    3: optional list<HookItem> flow_hook // Process hook
    4: optional list<HookItem> atomic_hook // Atomic power hook
    5: optional list<HookItem> llm_call_hook // Model call hook
    6: optional list<HookItem> res_parsing_hook // Conversation result hook
    7: optional list<HookItem> suggestion_hook // suggesion hook
}
struct HookItem {
    1: optional string uri
    2: optional list<string> filter_rules
    3: optional bool strong_dep
    4: optional i64 timeout_ms
}

//struct ContentAttachment {
//    1: required string FileID (api.body = "file_id")
//}

// struct MetaContent{
//     1: required string Type (api.body="type"),
//     2: optional string Text ( api.body="text"),
//     3: optional string FileID (api.body="file_id"),
//     4: optional string FileURL (api.body="file_url"),
//     5: optional string Card (api.body="card"),
// }


// struct EnterMessage  {
//     1: required string Role (api.body = "role")
//     2: string Content (api.body = "content")//content
//     3: map<string,string> MetaData(api.body = "meta_data")
//     4: string ContentType(api.body = "content_type")//text/card/object_string
//     5: string Type(api.body = "type")
// }

// struct OpenMessageApi {
//     1: string Id (api.body = "id")//primary key ID
//     2: string BotId (api.body = "bot_id")//bot id//TODO All i64 plus annotation str, imported parameters and exported parameters are required
//     3: string Role(api.body = "role")
//     4: string Content (api.body = "content")//content
//     5: string ConversationId(api.body = "conversation_id")   // conversation id
//     6: map<string,string> MetaData(api.body = "meta_data")
//     7: string CreatedAt (api.body = "created_at")//create_time
//     8: string UpdatedAt (api.body = "updated_at")//Update time//Change TODO time to int
//     9: string ChatId(api.body = "chat_id")
//     10: string ContentType(api.body = "content_type")
//     11: string Type(api.body = "type")
// }
