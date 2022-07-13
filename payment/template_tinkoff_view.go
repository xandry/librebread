package payment

const (
	tplTinkoffView = `
	<style>
		.btn-blue {
			border-color: blue;
		}
		.btn-green {
			border-color: green;
		}
		.btn-red {
			border-color: red;
		}
	</style>
	<h2>Tinkoff</h2>
	<table cellpadding="7">
		<tbody>
			<tr>
			<td>Process ID</td>
			<td><b>{{ .ProcessID }}</b></td>
			</tr>
			<tr>
			<td>Created on</td>
			<td>{{ (.CreatedOn).Format "2006-01-02 15:04:05" }}</td>
			</tr>
			<tr>
			<tr>
			<td>Amount</td>
			<td>{{ .Amount }}</td>
			</tr>
			<tr>
			<td>Description</td>
			<td>{{ .Description }}</td>
			</tr>
			<tr>
			<td>Status</td>
			<td><b>{{ .Status }}</b></td>
			</tr>
			{{ if eq .Status "NEW" "FORM_SHOWED" "AUTHORIZED" "CONFIRMED" }}
			<tr>
			<td colspan="2">
				{{ if eq .Status "NEW" "FORM_SHOWED" }}
					{{ if eq .Status "NEW" }}
						<a href="/tinkoff/{{ .ProcessID }}/set_status/FORM_SHOWED"><button class="btn-blue">Payment form is open<br><b>FORM_SHOWED</b></button></a>
					{{ end }}
					<a href="/tinkoff/{{ .ProcessID }}/set_status/DEADLINE_EXPIRED"><button class="btn-red">Payment time has expired<br><b>DEADLINE_EXPIRED</b></button></a>
					<a href="/tinkoff/{{ .ProcessID }}/set_status/ATTEMPTS_EXPIRED"><button class="btn-red">Attempts to open the form are exhausted<br><b>ATTEMPTS_EXPIRED</b></button></a>
					<a href="/tinkoff/{{ .ProcessID }}/set_status/CONFIRMED"><button class="btn-green">Single-stage payment<br><b>CONFIRMED</b></button></a>
					<a href="/tinkoff/{{ .ProcessID }}/set_status/AUTHORIZED"><button class="btn-green">Two-stage payment<br><b>AUTHORIZED</b></button></a>
					<a href="/tinkoff/{{ .ProcessID }}/set_status/REJECTED"><button class="btn-red">Rejected by the bank / insufficient funds<br><b>REJECTED</b></button></a>
				{{ else if eq .Status "AUTHORIZED" }}
					<a href="/tinkoff/{{ .ProcessID }}/set_status/CONFIRMED"><button class="btn-green">Payment confirmation<br><b>CONFIRMED</b></button></a>
				{{ else if eq .Status "CONFIRMED" }}
					<a href="/tinkoff/{{ .ProcessID }}/set_status/REFUNDED"><button class="btn-red">Refund of payment<br><b>REFUNDED</b></button></a>
				{{ end }}
			</td>
			</tr>
			{{ end }}
			<tr>
			<td>Payment URL</td>
			<td>{{ .PaymentURL }}</td>
			</tr>
			<tr>
			<td>Success URL</td>
			<td>{{ .SuccessURL }}</td>
			</tr>
			<tr>
			<td>Fail URL</td>
			<td>{{ .FailURL }}</td>
			</tr>
			<tr>
			<td>Notification URL</td>
			<td>{{ .NotificationURL }}</td>
			</tr>
			{{ if .NotificationURL }}
			<tr>
			<td>Response to notification received</td>
			<td><b>{{ .NotificationResponseOkReceived }}</b></td>
			</tr>
				{{ if eq .Status "AUTHORIZED" "REJECTED" "CONFIRMED" "REFUNDED" }}
				<tr>
				<td colspan="2">
				<a href="/tinkoff/{{ .ProcessID }}/send_notification"><button class="btn-blue">Send a notification</button></a>
				</td>
				</tr>
				{{ end }}
			{{ end }}
			<tr>
			<td>Recurrent</td>
			<td>{{ .Recurrent }}</td>
			</tr>
			<tr>
			<td>Client ID</td>
			<td>{{ .ClientID }}</td>
			</tr>
			<tr>
			<td>Order ID</td>
			<td>{{ .OrderID }}</td>
			</tr>
		</tbody>
	</table>
	`
)
