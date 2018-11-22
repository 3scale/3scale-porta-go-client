package fake

import "fmt"

func CreateApp(description string) string {
	return fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<application>
  <id>7654321</id>
  <created_at>2018-11-21T10:21:04Z</created_at>
  <updated_at>2018-11-21T10:21:04Z</updated_at>
  <state>live</state>
  <user_account_id>1234567</user_account_id>
  <first_traffic_at/>
  <first_daily_traffic_at/>
  <end_user_required>false</end_user_required>
  <service_id>2555417759506</service_id>
  <user_key>example-key-123-xyz</user_key>
  <provider_verification_key>example-provider-key-123-xyz</provider_verification_key>
  <plan custom="false" default="true">
    <id>2357355954203</id>
    <name>Basic</name>
    <type>application_plan</type>
    <state>published</state>
    <service_id>54321</service_id>
    <end_user_required>false</end_user_required>
    <setup_fee>0.0</setup_fee>
    <cost_per_month>0.0</cost_per_month>
    <trial_period_days/>
    <cancellation_period>0</cancellation_period>
  </plan>
  <name>example</name>
  <description>%s</description>
  <extra_fields></extra_fields>
</application>`, description)
}

func CreateAppFail() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<error>Your access token does not have the correct permissions</error>`
}

func GetProxyConfigLatestJson() string {
	return `{
 "proxy_config": {
  "id": 54321,
  "version": 2,
  "environment": "production",
  "content": {
   "id": 12345,
   "account_id": 5555555,
   "name": "Echo API",
   "oneline_description": null,
   "description": null,
   "txt_api": null,
   "txt_support": null,
   "txt_features": null,
   "created_at": "2018-07-23T08:33:55Z",
   "updated_at": "2018-07-23T08:54:37Z",
   "logo_file_name": null,
   "logo_content_type": null,
   "logo_file_size": null,
   "state": "incomplete",
   "intentions_required": false,
   "draft_name": "",
   "infobar": null,
   "terms": null,
   "display_provider_keys": false,
   "tech_support_email": null,
   "admin_support_email": null,
   "credit_card_support_email": null,
   "buyers_manage_apps": true,
   "buyers_manage_keys": true,
   "custom_keys_enabled": true,
   "buyer_plan_change_permission": "request",
   "buyer_can_select_plan": false,
   "notification_settings": null,
   "default_application_plan_id": 2357355954203,
   "default_service_plan_id": 2357355954201,
   "default_end_user_plan_id": null,
   "end_user_registration_required": true,
   "tenant_id": 2445582579513,
   "system_name": "api",
   "backend_version": "1",
   "mandatory_app_key": true,
   "buyer_key_regenerate_enabled": true,
   "support_email": "test@admin.com",
   "referrer_filters_required": false,
   "deployment_option": "hosted",
   "proxiable?": true,
   "backend_authentication_type": "provider_key",
   "backend_authentication_value": "123secret456",
   "proxy": {
    "id": 105757,
    "tenant_id": 123456,
    "service_id": 2555417759506,
    "endpoint": "https://api-123456.production.gw.apicast.io:443",
    "deployed_at": null,
    "api_backend": "https://echo-api.3scale.net:443",
    "auth_app_key": "app_key",
    "auth_app_id": "app_id",
    "auth_user_key": "user_key",
    "credentials_location": "query",
    "error_auth_failed": "Authentication failed",
    "error_auth_missing": "Authentication parameters missing",
    "created_at": "2018-07-23T08:33:55Z",
    "updated_at": "2018-07-23T08:54:37Z",
    "error_status_auth_failed": 403,
    "error_headers_auth_failed": "text/plain; charset=us-ascii",
    "error_status_auth_missing": 403,
    "error_headers_auth_missing": "text/plain; charset=us-ascii",
    "error_no_match": "No Mapping Rule matched",
    "error_status_no_match": 404,
    "error_headers_no_match": "text/plain; charset=us-ascii",
    "secret_token": "Shared_secret_sent_from_proxy_to_API_backend_db59ea91b44069d2",
    "hostname_rewrite": "",
    "oauth_login_url": null,
    "sandbox_endpoint": "https://api-123456.staging.gw.apicast.io:443",
    "api_test_path": "/some-path",
    "api_test_success": true,
    "apicast_configuration_driven": true,
    "oidc_issuer_endpoint": null,
    "lock_version": 3,
    "authentication_method": "1",
    "hostname_rewrite_for_sandbox": "echo-api.3scale.net",
    "endpoint_port": 443,
    "valid?": true,
    "service_backend_version": "1",
    "hosts": [
     "api-123456.production.gw.apicast.io",
     "api-123456.staging.gw.apicast.io"
    ],
    "backend": {
     "endpoint": "https://su1.3scale.net",
     "host": "su1.3scale.net"
    },
    "policy_chain": [
     {
      "name": "apicast",
      "version": "builtin",
      "configuration": {}
     }
    ],
    "proxy_rules": [
     {
      "id": 293072,
      "proxy_id": 105757,
      "http_method": "GET",
      "pattern": "/",
      "metric_id": 2555418115374,
      "metric_system_name": "hits",
      "delta": 1,
      "tenant_id": 123456,
      "created_at": "2018-07-23T08:33:55Z",
      "updated_at": "2018-07-23T08:33:55Z",
      "redirect_url": null,
      "parameters": [],
      "querystring_parameters": {}
     }
    ]
   }
  }
 }
}`
}
