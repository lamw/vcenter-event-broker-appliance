package provider

import (
	"context"

	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/metrics"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/processor"
)

// Provider manages the connection to an event provider and passes events to an
// event processor.
type Provider interface {
	PushMetrics(context.Context, metrics.Receiver)
	Stream(context.Context, processor.Processor) error
	Shutdown(context.Context) error
}
