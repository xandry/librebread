package api

const baseTmpl = `
	{{define "base"}}
	<!doctype html>
	<html>

	<head>
		<style>
			ol {
				padding: 10px;
				list-style-type: none;
			}

			ol li {
				float: left;
				margin: 0 10px 0 0;
			}
		</style>
	</head>

	<body>
		<nav>
			<ol>
				<li><a href="/sms">sms</a></li>
				<li><a href="/helpdesk">helpdesk</a></li>
				<li><a href="/email">email</a></li>
			</ol>
		</nav>
		<main>
			{{template "main" .}}
		</main>
	</body>

	</html>
	{{end}}`

const smsTempl = `
	{{template "base" .}}

	{{define "main"}}
	<table border=1>
		<caption>SMS</caption>
		<thead>
			<th>Date</th>
			<th>From</th>
			<th>Phone</th>
			<th>Msg</th>
			<th>Provider</th>
		</thead>
		{{range .}}
		<tr>
			<td>{{ .Time.Format "2006-01-02 15:04:05" }}</td>
			<td>{{ .From }}</td>
			<td>{{ .To }}</td>
			<td>{{ .Text }}</td>
			<td>{{ .Provider }}</td>
		</tr>
		{{end}}		
	</table>
	{{end}}`
