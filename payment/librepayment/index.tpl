<script src="/static/js/librepaymets.js"></script>

<table border=1>
    <caption>LibrePayment ({{len .}})</caption>
	<thead>
		<th>Time</th>
		<th>ID</th>
		<th>Amount</th>
        <th>Merchant</th>
        <th>Status</th>
		<th>Action</th>
	</thead>
    <tbody>
        {{range .}}
        <tr>
            <td>{{.Time}}</td>
            <td><a href="/librepayments/{{.ID}}">{{.ID}}</a></td>
            <td>{{.Amount}}</td>
            <td>{{.Merchant}}</td>
            <td>{{.Status}}</td>
            <td>
                <button onclick="librepaymentReject({{.ID}})">Reject</button>
                <button onclick="librepaymentConfirm({{.ID}})">Confirm</button>
            </td>
        </tr>
        {{end}}
    </tbody>
</table>
