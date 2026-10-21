include "./base.thrift"
include "./workflow/trace.thrift"
include "./workflow/workflow.thrift"
include "./app/bot_common.thrift"
include "./plugin/plugin_develop_common.thrift"
include "./permission/openapiauth.thrift"
include "./data/database/sheet_document.thrift"
include "./data/database/sheet_common.thrift"
include "./data/database/table.thrift"
include "./app/bot_open_api.thrift"
include "./plugin/plugin_develop.thrift"
include "./app/developer_api.thrift"
include "./playground/playground.thrift"
include "./data/database/database_svc.thrift"
include "./permission/openapiauth_service.thrift"
include "./passport/passport.thrift"
include "./workflow/workflow_svc.thrift"
include "./upload/upload.thrift"
include "./admin/config.thrift"

namespace go coze

service OpenAPIAuthService extends openapiauth_service.OpenAPIAuthService {}

service PluginDevelopService extends plugin_develop.PluginDevelopService {}

service DeveloperApiService extends developer_api.DeveloperApiService {}
service PlaygroundService extends playground.PlaygroundService {}
service DatabaseService extends database_svc.DatabaseService {}

service PassportService extends passport.PassportService {}
service WorkflowService extends workflow_svc.WorkflowService {}

service UploadService extends upload.UploadService {}
service ConfigService extends config.ConfigService {}

service AccountOpenApiService extends bot_open_api.AccountOpenApiService {}
