#!/bin/sh
# Patch gen/ before copy: add ValidateMetricsDataResponseBody to reports and suggestions client types.
# Goa does not emit this validator for MetricsDataResponseBody; the generated code calls it.
# copy-gen-to-api-gen.sh then copies the patched files into api/gen/.
set -e

patch_file() {
	FILE="$1"
	[ -f "$FILE" ] || return 0
	grep -q 'func ValidateMetricsDataResponseBody' "$FILE" && return 0
	awk '
	/^\/\/ ValidateReportResponseBody runs the validations defined on ReportResponseBody$/ {
		# Insert before ValidateReportResponseBody; do not add leading blank - file already has one.
		print "// ValidateMetricsDataResponseBody runs the validations defined on MetricsDataResponseBody. Goa does not emit this for composite map types; no required fields."
		print "func ValidateMetricsDataResponseBody(body *MetricsDataResponseBody) (err error) {"
		print "\treturn"
		print "}"
		print ""
	}
	{ print }
	' "$FILE" > "$FILE.tmp" && mv "$FILE.tmp" "$FILE"
	echo "Patched $FILE (added ValidateMetricsDataResponseBody)"
}

patch_file "gen/http/reports/client/types.go"
patch_file "gen/http/suggestions/client/types.go"
