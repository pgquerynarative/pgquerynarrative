#!/bin/sh
# Re-insert ValidateMetricsDataResponseBody into gen/http/reports/client/types.go after goa gen.
# Goa does not emit this validator for MetricsDataResponseBody; the generated code calls it.
set -e
FILE="gen/http/reports/client/types.go"
[ -f "$FILE" ] || exit 0
grep -q 'func ValidateMetricsDataResponseBody' "$FILE" && exit 0

# Insert stub after ValidateTimeSeriesDataResponseBody block, before ValidateReportResponseBody.
# Use awk for portability (macOS/BSD sed -i differs from GNU).
awk '
/^\/\/ ValidateReportResponseBody runs the validations defined on ReportResponseBody$/ {
	print ""
	print "// ValidateMetricsDataResponseBody runs the validations defined on MetricsDataResponseBody. Goa does not emit this for composite map types; no required fields."
	print "func ValidateMetricsDataResponseBody(body *MetricsDataResponseBody) (err error) {"
	print "\treturn"
	print "}"
	print ""
}
{ print }
' "$FILE" > "$FILE.tmp" && mv "$FILE.tmp" "$FILE"
echo "Patched $FILE (added ValidateMetricsDataResponseBody)"
