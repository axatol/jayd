package nr

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type segmentInterface interface {
	AddAttribute(key string, val interface{})
	End()
}

type mockSegment struct{}

type Attrs map[string]interface{}

func (*mockSegment) End()                                     {}
func (*mockSegment) AddAttribute(key string, val interface{}) {}

func Segment(ctx context.Context, name string, attributes ...Attrs) segmentInterface {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return &mockSegment{}
	}

	segment := txn.StartSegment(name)
	if len(attributes) == 1 {
		for key, value := range attributes[0] {
			segment.AddAttribute(key, value)
		}
	}

	return segment
}
