package fake

import (
	"context"

	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/events"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/logger"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/metrics"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/processor"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/provider"

	// "github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/stream"
	"github.com/vmware/govmomi/vim25/types"
)

const source = "https://fake.vcenter01.testing.io/sdk"

// verify that VCenter implements the streamer interface
var _ provider.Provider = (*VCenter)(nil)

// VCenter implements the streamer interface
type VCenter struct {
	eventCh <-chan []types.BaseEvent // channel which simulates events
	logger.Logger
}

// NewFakeVCenter returns a fake vcenter event stream provider streaming events
// received from the specified generator channel
func NewFakeVCenter(generator <-chan []types.BaseEvent, log logger.Logger) *VCenter {
	return &VCenter{
		eventCh: generator,
		Logger:  log,
	}
}

// PushMetrics is a no-op
func (f *VCenter) PushMetrics(context.Context, metrics.Receiver) {}

// Stream streams events generated by the Generator specified in the VCenter
// server
func (f *VCenter) Stream(ctx context.Context, p processor.Processor) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case baseEvent := <-f.eventCh:
			for idx := range baseEvent {
				// process slice in reverse order to maintain Event.Key ordering
				event := baseEvent[len(baseEvent)-1-idx]

				ce, err := events.NewFromVSphere(event, source)
				if err != nil {
					f.Logger.Errorw("skipping event because it could not be converted to CloudEvent format", "event", event, "error", err)
					continue
				}

				err = p.Process(ctx, *ce)
				if err != nil {
					f.Logger.Errorw("could not process event", "event", ce, "error", err)
				}
			}
		}
	}
}

// Shutdown is a no-op
func (f *VCenter) Shutdown(context.Context) error {
	return nil
}
