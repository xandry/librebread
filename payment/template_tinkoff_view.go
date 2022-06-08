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
			<td>Payment ID</td>
			<td><b>{{ .PaymentID }}</b></td>
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
						<a href="/tinkoff/{{ .PaymentID }}/set_status/FORM_SHOWED"><button class="btn-blue">Платежная форма открыта<br><b>FORM_SHOWED</b></button></a>
					{{ end }}
					<a href="/tinkoff/{{ .PaymentID }}/set_status/DEADLINE_EXPIRED"><button class="btn-red">Время оплаты истекло<br><b>DEADLINE_EXPIRED</b></button></a>
					<a href="/tinkoff/{{ .PaymentID }}/set_status/ATTEMPTS_EXPIRED"><button class="btn-red">Попытки открытия формы исчерпаны<br><b>ATTEMPTS_EXPIRED</b></button></a>
					<a href="/tinkoff/{{ .PaymentID }}/set_status/CONFIRMED"><button class="btn-green">Одностадийная оплата<br><b>CONFIRMED</b></button></a>
					<a href="/tinkoff/{{ .PaymentID }}/set_status/AUTHORIZED"><button class="btn-green">Двустадийная оплата<br><b>AUTHORIZED</b></button></a>
					<a href="/tinkoff/{{ .PaymentID }}/set_status/REJECTED"><button class="btn-red">Отклонен банком/недостаточно средств<br><b>REJECTED</b></button></a>
				{{ else if eq .Status "AUTHORIZED" }}
					<a href="/tinkoff/{{ .PaymentID }}/set_status/CONFIRMED"><button class="btn-green">Подтверждение оплаты<br><b>CONFIRMED</b></button></a>
				{{ else if eq .Status "CONFIRMED" }}
					<a href="/tinkoff/{{ .PaymentID }}/set_status/REFUNDED"><button class="btn-red">Возврат оплаты<br><b>REFUNDED</b></button></a>
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
			<td>Ответ на нотификацию получен</td>
			<td><b>{{ .NotificationResponseOkReceived }}</b></td>
			</tr>
				{{ if eq .Status "AUTHORIZED" "REJECTED" "CONFIRMED" "REFUNDED" }}
				<tr>
				<td colspan="2">
				<a href="/tinkoff/{{ .PaymentID }}/send_notification"><button class="btn-blue">Отправить нотификацию</button></a>
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
