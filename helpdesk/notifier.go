package helpdesk

type HelpdeskNotifier interface {
	HelpdeskNotify(msg HelpdeskMsg)
}
