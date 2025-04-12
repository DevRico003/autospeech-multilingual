package ui

import (
	"log"
)

// NotificationManager handles system notifications
type NotificationManager struct{}

// NewNotificationManager creates a new notification manager
func NewNotificationManager() *NotificationManager {
	return &NotificationManager{}
}

// ShowNotification displays a system notification
func (nm *NotificationManager) ShowNotification(message string) {
	// NOTIFICATIONS DISABLED FOR TESTING
	log.Printf("Notification disabled (for testing): %s", message)
	
	// Just log the message, don't show actual notifications
	log.Println(message)
}