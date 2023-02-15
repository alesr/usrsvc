// Package events defines the events that are published by the message broker.
// This package is located outside the internal directory because it is could be
// imported by other packages.

package events

type Event string

const (
	// Enumerate user entity change events.

	UserCreated Event = "user.created"

	// UserUpdated is the event that is published when a user is updated.
	UserUpdated Event = "user.updated"

	// UserDeleted is the event that is published when a user is deleted.
	UserDeleted Event = "user.deleted"
)
