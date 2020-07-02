const sseAddr = '/events'
const supportedEvents = ['email', 'helpdesk', 'sms']

function subscribeToEvents() {
  let eventSource = new EventSource(sseAddr)

  eventSource.onmessage = (event) => {
    msg = JSON.parse(event.data)

    if (!('event' in msg) || !supportedEvents.includes(msg.event)) {
      console.log('wrong message format:', msg)
      return
    }

    if (!isNotificationAllowed()) {
      return
    }

    notification(msg)
  }
}

function notification(msg) {
  const data = JSON.parse(msg.data)

  console.group(data)

  switch (msg.event) {
    case 'email':
      notifyEmail(data)
      break

    case 'helpdesk':
      notifyHelpdesk(data)
      break

    case 'sms':
      notifySMS(data)
      break

    default:
      console.log('unsupported event type:', msg.event)
  }
}

function isNotificationAllowed() {
  if (!('Notification' in this)) {
    console.error('This browser does not support desktop notification')
    return false
  }

  return Notification.permission === 'granted'
}

function notifyEmail(email) {
  const options = {
    body: email.body
  }
  
  new Notification('Email', options)
}

function notifyHelpdesk(helpdesk) {
  console.log(helpdesk)
}

function notifySMS(sms) {
  console.log(sms)
}

subscribeToEvents()
