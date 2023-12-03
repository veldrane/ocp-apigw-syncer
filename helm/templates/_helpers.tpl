{{/*
Expand the name of the chart.
*/}}
{{- define "Helm.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "Helm.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "Helm.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "Helm.labels" -}}
helm.sh/chart: {{ include "Helm.chart" . }}
{{ include "Helm.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "Helm.selectorLabels" -}}
app.kubernetes.io/name: {{ include "Helm.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "Helm.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "Helm.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{/*
Define FTE helpers
*/}}
{{- define "dataCenter" -}}
{{- if (eq .Values.env "p") -}}
{{- if or (eq .Values.operator "cp") (or (eq .Values.operator "ro") (eq .Values.operator "shared")) -}}
{{- $dc := "cz" -}}
{{- $dc -}}
{{- else -}}
{{- $dc := .Values.operator -}}
{{- $dc -}}
{{- end -}}
{{- else -}}
{{- $dc := "cz" -}}
{{- $dc -}}
{{- end -}}
{{- end -}}

{{- define "ipAddress" -}}
{{- $prefix := "" -}}
{{- $prefixes := dict "d" "10.8.171." -}}
{{- $ocp4Prefixes := dict "d" "10.8.173." "t" "10.6.173." "s" "10.4.178." "cp" "10.2.181." "cz" "10.2.173." "ro" "10.2.177." "pl" "10.62.173." "hr" "10.92.173." "shared" "10.2.185." "sk" "10.72.173." -}}
{{- if .Values.ocp4 -}}
  {{- $prefix = get $ocp4Prefixes .Values.env -}}
{{- else -}}
  {{- $prefix = index $prefixes .Values.env -}}
{{- end -}}
{{- printf $prefix -}}
{{- end -}}

{{- define "replicaCount" -}}
{{- if not (empty (.Values.replicaCount)) -}}
{{- $replicacount := .Values.replicaCount -}}
{{- printf "%s" ($replicacount | toString) -}}
{{ else if (or (eq .Values.env "p") (eq .Values.env "s")) -}}
{{- printf "2" -}}
{{ else if (or (eq .Values.env "t") (eq .Values.env "d")) -}}
{{- printf "1" -}}
{{- end -}}
{{- end -}}

{{- define "requestsCpu" -}}
{{- if not (empty .Values.resources ) -}}
{{- $requestCpu := .Values.resources.requests.cpu -}}
{{- printf "%s" $requestCpu -}}
{{- else if (or (eq .Values.env "p") (eq .Values.env "s")) -}}
{{- printf "2" -}}
{{ else if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- printf "100m" -}}
{{- end -}}
{{- end -}}


{{- define "requestsMemory" -}}
{{- if not (empty .Values.resources ) -}}
{{- $requestMemory := .Values.resources.requests.memory -}}
{{- printf "%s" $requestMemory -}}
{{- else if (or (eq .Values.env "p") (eq .Values.env "s")) -}}
{{- printf "2048Mi" -}}
{{ else if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- printf "100Mi" -}}
{{- end -}}
{{- end -}}

{{- define "accessKeyValDbSize" -}}
{{- if (or (eq .Values.env "p") (eq .Values.env "s")) -}}
{{- printf "256M" -}}
{{- else if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- printf "100M" -}}
{{- end -}}
{{- end -}}

{{- define "refreshKeyValDbSize" -}}
{{- if (or (eq .Values.env "p") (eq .Values.env "s")) -}}
{{- printf "512M" -}}
{{- else if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- printf "100M" -}}
{{- end -}}
{{- end -}}

{{- define "limitsCpu" -}}
{{- if not (empty .Values.resources ) -}}
{{- $limitsCpu := .Values.resources.limits.cpu -}}
{{- printf "%s" $limitsCpu -}}
{{- else if (or (eq .Values.env "p") (eq .Values.env "s")) -}}
{{- printf "8" -}}
{{ else if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- printf "200m" -}}
{{- end -}}
{{- end -}}

{{- define "limitsMemory" -}}
{{- if not (empty .Values.resources ) -}}
{{- $limitsMemory := .Values.resources.limits.memory -}}
{{- printf "%s" $limitsMemory -}}
{{- else if (or (eq .Values.env "p") (eq .Values.env "s")) -}}
{{- printf "2048Mi" -}}
{{ else if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- printf "400Mi" -}}
{{- end -}}
{{- end -}}

{{- define "fasTokenEndpoint" -}}
{{- if (and (empty .Values.oidc.fasTokenEndpoint ) (not (empty .Values.oidc.fasServer))) -}}
{{- $fasServer := .Values.oidc.fasServer -}}
{{- printf "%s%s" $fasServer "/oauth/token" -}}
{{- else if (not (empty .Values.oidc.fasTokenEndpoint)) -}}
{{- $fasTokenEndpoint := .Values.oidc.fasTokenEndpoint -}}
{{- printf "%s" $fasTokenEndpoint -}}
{{- end -}}
{{- end -}}

{{- define "fasLogoutEndpoint" -}}
{{- if (and (empty .Values.oidc.fasLogoutEndpoint ) (not (empty .Values.oidc.fasServer))) -}}
{{- $fasServer := .Values.oidc.fasServer -}}
{{- printf "%s%s" $fasServer "/oauth/logout" -}}
{{- else if (not (empty .Values.oidc.fasLogoutEndpoint)) -}}
{{- $fasLogoutEndpoint := .Values.oidc.fasLogoutEndpoint -}}
{{- printf "%s" $fasLogoutEndpoint -}}
{{- end -}}
{{- end -}}

{{- define "fasTokenNativeEndpoint" -}}
{{- if (and (empty .Values.oidc.fasTokenEndpoint ) (not (empty .Values.oidc.fasServer))) -}}
{{- $fasServer := .Values.oidc.fasServer -}}
{{- printf "%s%s" $fasServer "/v2/oauth/token" -}}
{{- else if (not (empty .Values.oidc.fasTokenEndpoint)) -}}
{{- $fasTokenNativeEndpoint := .Values.oidc.fasTokenEndpoint -}}
{{- printf "%s" $fasTokenNativeEndpoint -}}
{{- end -}}
{{- end -}}

/* virtuals login */
{{- define "fasTokenVirtualsLoginEndpoint" -}}
{{- if (and (empty .Values.oidc.fasTokenEndpoint ) (not (empty .Values.oidc.fasServer))) -}}
{{- $fasServer := .Values.oidc.fasServer -}}
{{- printf "%s%s" $fasServer "/v2/oauth/token" -}}
{{- else if (not (empty .Values.oidc.fasTokenEndpoint)) -}}
{{- $fasTokenVirtualsLoginEndpoint := .Values.oidc.fasTokenEndpoint -}}
{{- printf "%s" $fasTokenVirtualsLoginEndpoint -}}
{{- end -}}
{{- end -}}

/* virtuals logout */
{{- define "fasTokenVirtualsLogoutEndpoint" -}}
{{- if (and (empty .Values.oidc.fasTokenEndpoint ) (not (empty .Values.oidc.fasServer))) -}}
{{- $fasServer := .Values.oidc.fasServer -}}
{{- printf "%s%s" $fasServer "/oauth/logout" -}}
{{- else if (not (empty .Values.oidc.fasTokenEndpoint)) -}}
{{- $fasTokenVirtualsLogoutEndpoint := .Values.oidc.fasTokenEndpoint -}}
{{- printf "%s" $fasTokenVirtualsLogoutEndpoint -}}
{{- end -}}
{{- end -}}

{{- define "fasCheckEndpoint" -}}
{{- if (and (empty .Values.oidc.fasCheckEndpoint ) (not (empty .Values.oidc.fasServer))) -}}
{{- $fasServer := .Values.oidc.fasServer -}}
{{- printf "%s%s" $fasServer "/oauth/check_token" -}}
{{- else if (not (empty .Values.oidc.fasCheckEndpoint)) -}}
{{- $fasCheckEndpoint := .Values.oidc.fasCheckEndpoint -}}
{{- printf "%s" $fasCheckEndpoint -}}
{{- end -}}
{{- end -}}


{{- define "fasNotificationEndpoint" -}}
{{- if (and (empty .Values.oidc.fasCheckEndpoint ) (not (empty .Values.oidc.fasServer))) -}}
{{- $fasServer := .Values.oidc.fasServer -}}
{{- printf "%s%s" $fasServer "/oauth/notification_token" -}}
{{- else if (not (empty .Values.oidc.fasNotoficationEndpoint)) -}}
{{- $fasNotificationEndpoint := .Values.oidc.fasNotificationEndpoint -}}
{{- printf "%s" $fasNotificationEndpoint -}}
{{- end -}}
{{- end -}}

{{- define "fasJwtKeyFile" -}}
{{- if (and (empty .Values.oidc.fasJwtKeyFile ) (not (empty .Values.oidc.fasServer))) -}}
{{- $fasServer := .Values.oidc.fasServer -}}
{{- printf "%s%s" $fasServer "/oauth/keys/master/details" -}}
{{- else if (not (empty .Values.oidc.fasJwtKeyFile)) -}}
{{- $fasJwtKeyFile := .Values.oidc.fasJwtKeyFile -}}
{{- printf "%s" $fasJwtKeyFile -}}
{{- end -}}
{{- end -}}

{{- define "fasGrantType" -}}
{{- if (empty .Values.oidc.fasGrantType) -}}
{{- printf "%s" "urn:feg:params:oauth:grant-type:ims-jwt-bearer" -}}
{{- else -}}
{{- $fasGrantType := .Values.oidc.fasGrantType -}}
{{- printf "%s" $fasGrantType -}}
{{- end -}}
{{- end -}}

{{- define "fasClientId" -}}
{{- if (empty .Values.oidc.fasClientId) -}}
{{- printf "%s" "test_client_id" -}}
{{- else -}}
{{- $fasClientId := .Values.oidc.fasClientId -}}
{{- printf "%s" $fasClientId -}}
{{- end -}}
{{- end -}}

{{- define "oidcHmacKey" -}}
{{- if (empty .Values.oidc.oidcHmacKey) -}}
{{- printf "%s" "nFDQ+HhAq4ldz1ufp1Y31Yse" -}}
{{- else -}}
{{- $oidcHmacKey := .Values.oidc.oidcHmacKey -}}
{{- printf "%s" $oidcHmacKey -}}
{{- end -}}
{{- end -}}

{{- define "cookieFlagsHttp" -}}
{{- if (empty .Values.oidc.cookieFlagsHttp) -}}
{{- printf "%s" "Path=/; SameSite=None;" -}}
{{- else -}}
{{- $cookieFlagsHttp := .Values.oidc.cookieFlagsHttp -}}
{{- printf "%s" $cookieFlagsHttp -}}
{{- end -}}
{{- end -}}

{{- define "cookieFlagsHttps" -}}
{{- $cookieDomain := "ifortuna.cz" -}}
{{- if (eq .Values.operator "cz") -}}
{{- $cookieDomain = "ifortuna.cz" -}}
{{- else if (eq .Values.operator "hr") -}}
{{- $cookieDomain = "psk.hr" -}}
{{- else if (eq .Values.operator "sk") -}}
{{- $cookieDomain = "ifortuna.sk" -}}
{{- else if (eq .Values.operator "pl") -}}
{{- $cookieDomain = "efortuna.pl" -}}
{{- else if (eq .Values.operator "cp") -}}
{{- $cookieDomain = "casapariurilor.ro" -}}
{{- else if (eq .Values.operator "ro") -}}
{{- $cookieDomain = "efortuna.ro" -}}
{{- end -}}
{{- if (empty .Values.oidc.cookieFlagsHttps) -}}
{{- printf "%s%s%s" "Path=/; SameSite=Strict; HttpOnly; Secure; Domain=" $cookieDomain ";" -}}
{{- else -}}
{{- $cookieFlagsHttps := .Values.oidc.cookieFlagsHttps -}}
{{- printf "%s" $cookieFlagsHttps -}}
{{- end -}}
{{- end -}}

{{- define "cookieFlagsClean" -}}
{{- $cookieDomain := "ifortuna.cz" -}}
{{- if (eq .Values.operator "cz") -}}
{{- $cookieDomain = "ifortuna.cz" -}}
{{- else if (eq .Values.operator "hr") -}}
{{- $cookieDomain = "psk.hr" -}}
{{- else if (eq .Values.operator "sk") -}}
{{- $cookieDomain = "ifortuna.sk" -}}
{{- else if (eq .Values.operator "pl") -}}
{{- $cookieDomain = "efortuna.pl" -}}
{{- else if (eq .Values.operator "cp") -}}
{{- $cookieDomain = "casapariurilor.ro" -}}
{{- else if (eq .Values.operator "ro") -}}
{{- $cookieDomain = "efortuna.ro" -}}
{{- end -}}
{{- if (empty .Values.oidc.cookieFlagsClean) -}}
{{- printf "%s%s%s" "Max-Age=0; Expires=Thu, 1 Jan 1970 00:00:00 GMT; Path=/; SameSite=Strict; HttpOnly; Secure; Domain=" $cookieDomain ";" -}}
{{- else -}}
{{- $cookieFlagsClean := .Values.oidc.cookieFlagsClean -}}
{{- printf "%s" $cookieFlagsClean -}}
{{- end -}}
{{- end -}}

{{- define "apigwDomain" -}}
{{- $apigwDomain := "ifortuna.cz" -}}
{{- if (eq .Values.operator "cz") -}}
{{- $apigwDomain = "ifortuna.cz" -}}
{{- else if (eq .Values.operator "hr") -}}
{{- $apigwDomain = "psk.hr" -}}
{{- else if (eq .Values.operator "sk") -}}
{{- $apigwDomain = "ifortuna.sk" -}}
{{- else if (eq .Values.operator "pl") -}}
{{- $apigwDomain = "efortuna.pl" -}}
{{- else if (eq .Values.operator "cp") -}}
{{- $apigwDomain = "casapariurilor.ro" -}}
{{- else if (eq .Values.operator "ro") -}}
{{- $apigwDomain = "efortuna.ro" -}}
{{- end -}}
{{- printf "%s" $apigwDomain -}}
{{- end -}}

{{- define "cors" -}}
{{- if (empty .Values.oidc.cors) -}}
{{- if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- printf "%s" "true" -}}
{{- else if (or (eq .Values.env "s") (eq .Values.env "p")) -}}
{{- printf "%s" "true" -}}
{{- end -}}
{{- else -}}
{{- $cors := .Values.oidc.cors -}}
{{- printf "%s" $cors -}}
{{- end -}}
{{- end -}}

{{- define "corsDefaultOrigin" -}}
{{- if (empty .Values.corsDefaultOrigin) -}}
{{- if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- printf "%s" "https://localhost.ifortuna.cz:8080" -}}
{{- else if (eq .Values.env "s") -}}
{{- if (eq .Values.operator "cz") -}}
{{- printf "%s" "https://live-dc1.stage.ifortuna.cz" -}}
{{- else if (eq .Values.operator "hr") -}}
{{- printf "%s" "https://live-dc1.stage.psk.hr" -}}
{{- else if (eq .Values.operator "sk") -}}
{{- printf "%s" "https://live-dc1.stage.ifortuna.sk" -}}
{{- else if (eq .Values.operator "pl") -}}
{{- printf "%s" "https://live-dc1.stage.efortuna.pl" -}}
{{- else if (eq .Values.operator "cp") -}}
{{- printf "%s" "https://live-dc1.stage.casapariurilor.ro" -}}
{{- else if (eq .Values.operator "ro") -}}
{{- printf "%s" "https://live3web-stg.efortuna.ro" -}}
{{- end -}}
{{- else if (eq .Values.env "p") -}}
{{- if (eq .Values.operator "cz") -}}
{{- printf "%s" "https://live.ifortuna.cz" -}}
{{- else if (eq .Values.operator "hr") -}}
{{- printf "%s" "https://live.psk.hr" -}}
{{- else if (eq .Values.operator "sk") -}}
{{- printf "%s" "https://live.ifortuna.sk" -}}
{{- else if (eq .Values.operator "pl") -}}
{{- printf "%s" "https://live.efortuna.pl" -}}
{{- else if (eq .Values.operator "cp") -}}
{{- printf "%s" "https://live.casapariurilor.ro" -}}
{{- else if (eq .Values.operator "ro") -}}
{{- printf "%s" "https://live.efortuna.ro" -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "periodSecond" -}}
{{- if (not (empty .Values.oidc)) -}}
{{- if (empty .Values.oidc.keepSessions) -}}
{{- if (eq .Values.env "p") -}}
{{- printf "%s" "30" -}}
{{- else -}}
{{- printf "%s" "1" -}}
{{- end -}}
{{- else -}}
{{- if (eq .Values.oidc.keepSessions true) -}}
{{- printf "%s" "30" -}}
{{- else -}}
{{- printf "%s" "1" -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "debugHeaders" -}}
{{- if (not (empty .debug)) -}}
{{- if (eq .debug true) -}}
{{- printf "%s" "include /etc/nginx/conf.d/include/debug-headers.resty;\n        " -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "logFormat" -}}
{{- if (empty .Values.logFormat) -}}
{{- if (or (eq .Values.env "s") (eq .Values.env "p")) -}}
{{- if not (empty .Values.oidc ) -}}
{{- if (eq .Values.oidc.type "feg") -}}
{{- printf "%s" "feg_oidc_json" -}}
{{- else if (eq .Values.oidc.type "standard") -}}
{{- printf "%s" "standard_oidc_json" -}}
{{- end -}}
{{- else -}}
{{- printf "%s" "default_text" -}}
{{- end -}}
{{- else -}}
{{- if not (empty .Values.oidc ) -}}
{{- if (eq .Values.oidc.type "feg") -}}
{{- printf "%s" "feg_oidc_text" -}}
{{- else if (eq .Values.oidc.type "standard") -}}
{{- printf "%s" "standard_oidc_text" -}}
{{- end -}}
{{- else -}}
{{- printf "%s" "default_text" -}}
{{- end -}}
{{- end -}}
{{- else -}}
{{- $logFormat := .Values.logFormat -}}
{{- printf "%s" $logFormat -}}
{{- end -}}
{{- end -}}

{{- define "corsOrigins" }}
{{- if (or (eq .Values.env "d") (eq .Values.env "t")) -}}
{{- include "cors.origins.testdev" . -}}
{{- else if (eq .Values.env "s") -}}
{{- include "cors.origins.stage" . -}}
{{- else if (eq .Values.env "p") -}}
{{- include "cors.origins.prod" . -}}
{{- end }}
{{- end }}



{{- define "cors.origins.testdev" -}}
{{- if (empty .Values.corsDefaultOrigin) -}}
map $http_origin $cors_allow_origin {
        default $http_origin;
        "" "https://localhost.ifortuna.cz:8080";
      }
{{ else -}}
map $http_origin $cors_allow_origin {
        default $http_origin;
        "" {{ .Values.corsDefaultOrigin -}};
      }
{{ end -}}
{{ end -}}


{{- define "cors.origins.stage" -}}
{{- $default := "" -}}
{{- $origins := "" -}}
{{- if (eq .Values.operator "cz") -}}
{{- $default = (index .Values.cors.allowOrigins.stage.cz 0) -}}
{{- $origins = .Values.cors.allowOrigins.stage.cz -}}
{{- else if (eq .Values.operator "ro") -}}
{{- $default = (index .Values.cors.allowOrigins.stage.ro 0) -}}
{{- $origins = .Values.cors.allowOrigins.stage.ro -}}
{{- else if (eq .Values.operator "pl") -}}
{{- $default = (index .Values.cors.allowOrigins.stage.pl 0) -}}
{{- $origins = .Values.cors.allowOrigins.stage.pl -}}
{{- else if (eq .Values.operator "sk") -}}
{{- $default = (index .Values.cors.allowOrigins.stage.sk 0) -}}
{{- $origins = .Values.cors.allowOrigins.stage.sk -}}
{{- else if (eq .Values.operator "hr") -}}
{{- $default = (index .Values.cors.allowOrigins.stage.hr 0) -}}
{{- $origins = .Values.cors.allowOrigins.stage.hr -}}
{{- else if (eq .Values.operator "cp") -}}
{{- $default = (index .Values.cors.allowOrigins.stage.cp 0) -}}
{{- $origins = .Values.cors.allowOrigins.stage.cp -}}
{{- end -}}
{{- printf "%s" "map $http_origin $cors_allow_origin {" -}}
{{- printf "%s" "\n" -}}
{{- printf "%s%s%s" "default \"" $default "\";" | indent 8 -}}
{{- printf "%s" "\n" -}}
{{- range $v := $origins -}}
{{- $origin := $v -}}
{{- printf "%s%s%s" "\"" $origin "\" $http_origin;" | indent 8 -}}
{{- printf "%s" "\n" -}}
{{- end -}}
{{- printf "%s" "}" | indent 6 -}}
{{ end -}}

{{- define "cors.origins.prod" -}}
{{- $default := "" -}}
{{- $origins := "" -}}
{{- if (eq .Values.operator "cz") -}}
{{- $default = (index .Values.cors.allowOrigins.prod.cz 0) -}}
{{- $origins = .Values.cors.allowOrigins.prod.cz -}}
{{- else if (eq .Values.operator "ro") -}}
{{- $default = (index .Values.cors.allowOrigins.prod.ro 0) -}}
{{- $origins = .Values.cors.allowOrigins.prod.ro -}}
{{- else if (eq .Values.operator "pl") -}}
{{- $default = (index .Values.cors.allowOrigins.prod.pl 0) -}}
{{- $origins = .Values.cors.allowOrigins.prod.pl -}}
{{- else if (eq .Values.operator "sk") -}}
{{- $default = (index .Values.cors.allowOrigins.prod.sk 0) -}}
{{- $origins = .Values.cors.allowOrigins.prod.sk -}}
{{- else if (eq .Values.operator "hr") -}}
{{- $default = (index .Values.cors.allowOrigins.prod.hr 0) -}}
{{- $origins = .Values.cors.allowOrigins.prod.hr -}}
{{- else if (eq .Values.operator "cp") -}}
{{- $default = (index .Values.cors.allowOrigins.prod.cp 0) -}}
{{- $origins = .Values.cors.allowOrigins.prod.cp -}}
{{- end -}}
{{- printf "%s" "map $http_origin $cors_allow_origin {" -}}
{{- printf "%s" "\n" -}}
{{- printf "%s%s%s" "default \"" $default "\";" | indent 8 -}}
{{- printf "%s" "\n" -}}
{{- range $v := $origins -}}
{{- $origin := $v -}}
{{- printf "%s%s%s" "\"" $origin "\" $http_origin;" | indent 8 -}}
{{- printf "%s" "\n" -}}
{{- end -}}
{{- printf "%s" "}" | indent 6 -}}
{{ end -}}

{{- define "corsHeaders" -}}
{{- $allow_headers := .Values.cors.allowHeaders -}}
{{- $allow_methods := .Values.cors.allowMethods -}}
{{- $allow_credentials := .Values.cors.allowCredentials -}}
{{- printf "%s%s%s" "set $cors_allow_headers \"" $allow_headers "\";"  -}}
{{- printf "%s" "\n" -}}
{{- printf "%s%s%s" "set $cors_allow_methods \"" $allow_methods "\";" | indent 14 -}}
{{- printf "%s" "\n" -}}
{{- printf "%s%s%s" "set $cors_allow_credentials \"" $allow_credentials "\";" | indent 14 -}}
{{- printf "%s" "\n" -}}
{{- end -}}

{{- define "nginxCorsHeaders" -}}
{{- printf "%s" "add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;" -}}
{{- printf "%s" "\n" -}}
{{- printf "%s" "add_header 'Access-Control-Allow-Credentials' $cors_allow_credentials always;" -}}
{{- printf "%s" "\n" -}}
{{- printf "%s" "add_header 'Access-Control-Allow-Methods' $cors_allow_methods always;" -}}
{{- printf "%s" "\n" -}}
{{- printf "%s" "add_header 'Access-Control-Allow-Headers' $cors_allow_headers always;" -}}
{{- printf "%s" "\n" -}}
{{- end -}}

{{- define "serverName" -}}
{{- if regexMatch "apigw.+" .Release.Namespace -}}
{{- if (and (eq .Values.env "p") (regexMatch "sk|pl|hr" .Values.operator)) -}}
{{- printf "%s%s.%s.%s.%s" "api-"  .Release.Namespace "dc1" .Values.operator "ipa.ifortuna.cz" -}}
{{- else -}}
{{- printf "%s%s.%s%s" "api-"  .Release.Namespace .Values.env ".dc1.cz.ipa.ifortuna.cz" -}}
{{- end -}}
{{- else -}}
{{- printf "%s.%s.%s" "nginx-apigw" .Release.Namespace "svc.cluster.local" -}}
{{- end -}}
{{- end -}}

{{- define "clusterDns" -}}
{{- if ne .Values.ocp4 true -}}
{{- printf "%s" "10.39.0.1" -}}
{{- else -}}
{{- printf "%s" "10.39.0.10" -}}
{{- end -}}
{{- end -}}

{{- define "imagePath" -}}
{{- if ne .Values.ocp4 true -}}
{{- printf "%s" "docker-registry.default.svc:5000/openshift" -}}
{{- else -}}
{{- printf "%s" "image-registry.openshift-image-registry.svc:5000/shared-images" -}}
{{- end -}}
{{- end -}}
