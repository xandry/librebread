if (!(Worker in window)) {
  console.log('browser does not support workers')
}

const myWorker = new Worker('/static/js/notificationworker.js')
