function librepaymentConfirm(payment_id) {
    let url = `/libre/payment/${payment_id}/confirm`;

    fetch(url, {
        method: 'POST',
    }).then((resp) => {
        window.location.reload();
    })
}

function librepaymentReject(payment_id) {
    let url = `/libre/payment/${payment_id}/reject`;

    fetch(url, {
        method: 'POST',
    }).then((resp) => {
        window.location.reload();
    })
}