package service

import (
	"github.com/alesr/usrsvc/pkg/events"
)

var _ Publisher = (*publisherMock)(nil)

// publisherMock is a mock implementation of the publisher interface.
type publisherMock struct {
	PublishFunc func(event events.Event, data any) error
}

func (p *publisherMock) Publish(event events.Event, data any) error {
	return p.PublishFunc(event, data)
}
