{{- define "cliapp.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cliapp.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "cliapp.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
