if (!('Worker' in window)) {
  console.log('browser does not support workers')
}

const notifyWorker = new Worker('/static/js/notificationworker.js')
