package sudo

type NotifySender func(NotifyCommand)

type NotifyCommand struct {
	DeviceId string
	Code     string
	Phone    string
	Email    string
	Locale   string
}
