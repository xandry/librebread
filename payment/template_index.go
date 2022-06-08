package payment

const (
	tplIndex = `
	<table cellpadding="7" border=1>
		<thead>
			<tr>
			<td>Created On</td>
			<td>Provider</td>
			<td>Payment ID</td>
			<td>Amount</td>
			<td>Status</td>
			<td>Recurrent</td>
			<td>Description</td>
			</tr>
		</thead>
		<tbody>
			{{ range .Payments }}
			<tr>
			<td>{{ (.CreatedOn).Format "2006-01-02 15:04:05" }}</td>
			<td>{{ .Provider }}</td>
			<td><a href="{{ .PaymentURL }}">{{ .PaymentID }}</a></td>
			<td>{{ .Amount }}</td>
			<td><b>{{ .Status }}</b></td>
			<td>{{ .Recurrent }}</td>
			<td>{{ .Description }}</td>
			</tr>
			{{ end }}
		</tbody>
	</table>
	`
)
