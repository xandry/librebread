const sseAddr = '/events'
const supportedEvents = ['email', 'helpdesk', 'sms']

function subscribeToEvents() {
  let eventSource = new EventSource(sseAddr)

  eventSource.onmessage = (event) => {
    console.log(event.data)

    msg = JSON.parse(event.data)

    if (!('event' in msg.event) || !supportedEvents.includes(msg.event)) {
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
  switch (msg.event) {
    case 'email':
      notifyEmail(msg.data)
      break

    case 'helpdesk':
      notifyHelpdesk(msg.data)
      break

    case 'sms':
      notifySMS(msg.data)
      break

    default:
      console.log('unsupported event type:', msg.event)
  }
}

function isNotificationAllowed() {
  if (!('Notification' in window)) {
    console.error('This browser does not support desktop notification')
    return false
  }

  return Notification.permission === 'granted'
}

function notifyEmail(email) {
  console.log(email)
}

function notifyHelpdesk(helpdesk) {
  console.log(helpdesk)
}

function notifySMS(sms) {
  console.log(sms)
}

subscribeToEvents()
