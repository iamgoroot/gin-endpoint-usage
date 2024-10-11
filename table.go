package ginUsageStats

import "text/template"

var tableTemplate = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Endpoint usage stats</title>
</head>
<body>
    <h1>Stats Table</h1>
    <table>
        <thead><tr><th>Method</th><th>Endpoint</th><th>Count</th></tr>
        </thead>
        <tbody>
            {{- range . -}}
<tr><td>{{ .Method }}</td><td>{{ .Endpoint }}</td><td>{{ .Count }}</td></tr>
            {{- end }}
        </tbody>
    </table>
</body>
</html>`))
