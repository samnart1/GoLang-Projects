package notifications

import "fmt"

// Notifier is an interface that abstracts notification functionalities
type Notifier struct {
	desktopNotifier *DesktopNotifier
	enabled         bool
}

// NewNotifier creates a new Notifier instance
func NewNotifier() *Notifier {
	return &Notifier{
		desktopNotifier: NewDesktopNotifier(true), // Enable desktop notifications by default
		enabled:         true,                     // Enable notifications by default
	}
}

// PlaySound plays a notification sound
func (n *Notifier) PlaySound() {
	if !n.enabled {
		return
	}
	// Play a sound (this is a placeholder, implementation may vary)
	fmt.Println("Playing notification sound...")
}

// ShowDesktop displays a desktop notification
func (n *Notifier) ShowDesktop(title, message string) error {
	if !n.enabled || n.desktopNotifier == nil {
		return nil
	}
	return n.desktopNotifier.Show(title, message)
}

// ShowTerminal prints a terminal notification
func (n *Notifier) ShowTerminal(message string) {
	if !n.enabled {
		return
	}
	fmt.Printf("Terminal Notification: %s\n", message)
}