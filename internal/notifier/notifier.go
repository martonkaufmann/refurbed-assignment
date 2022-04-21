package notifier

type Notifier interface {
	Process()
	Enqueue(payload string)
}
