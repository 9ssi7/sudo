// Package sudo provides a service for handling sudo operations.
package sudo

// NotifySender is a function type for sending notifications.
type NotifySender func(NotifyCommand)

// NotifyCommand represents the command structure for sending notifications.
type NotifyCommand struct {
	// DeviceId is the unique identifier of the device.
	DeviceId string

	// Code is the verification code to be sent.
	Code string

	// Phone is the phone number to send the code to.
	Phone string

	// Email is the email address to send the code to.
	Email string

	// Locale is the language and region to use for the notification.
	Locale string
}
