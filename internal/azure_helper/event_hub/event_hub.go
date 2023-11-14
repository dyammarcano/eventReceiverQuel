package event_hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-event-hubs-go/v3"
	"github.com/dyammarcano/eventReceiverQuel/internal/config"
	"log"
	"net"
	"os"
	"sync"
)

type (
	HubClient struct {
		*eventhub.Hub
		context.Context
	}
)

// NewHubClient Get a new client for event hub
func NewHubClient(ctx context.Context, config *config.Config) (*HubClient, error) {
	connStr := fmt.Sprintf("Endpoint=sb://%s.servicebus.windows.net;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=%s;EntityPath=%s",
		config.EventHub.AccountName, config.EventHub.AccountKey, config.EventHub.Topic)

	if err := checkDefaultPorts(config); err != nil {
		return nil, err
	}

	return newHubClientSAS(ctx, connStr)
}

func checkDefaultPorts(cfg *config.Config) error {
	for _, port := range []int{5671, 5672} {
		if err := connectionTest(cfg.EventHub.AccountName, port); err != nil {
			return err
		}
	}
	return nil
}

func connectionTest(accountName string, port int) error {
	address := fmt.Sprintf("%s.servicebus.windows.net:%d", accountName, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("port %d is closed, err: %v\n", port, err)
	}

	defer conn.Close()

	return nil
}

// newHubClientSAS Get a new client for event hub
func newHubClientSAS(ctx context.Context, connStr string) (*HubClient, error) {
	hub, err := eventhub.NewHubFromConnectionString(connStr)
	if err != nil {
		return nil, err
	}

	return &HubClient{Hub: hub, Context: ctx}, nil
}

// SendEvent Send a single message into a random partition
func (h *HubClient) SendEvent(message string) error {
	if err := h.Hub.Send(h.Context, eventhub.NewEventFromString(message)); err != nil {
		return errors.New(fmt.Sprintf("error sending event: %v", err))
	}
	return nil
}

// SendEventJson Send a single message into a random partition
func (h *HubClient) SendEventJson(obj any) error {
	data, err := json.Marshal(obj)

	if err != nil {
		return errors.New("error marshaling event")
	}

	if err = h.Hub.Send(h.Context, eventhub.NewEventFromString(string(data))); err != nil {
		return errors.New(fmt.Sprintf("error sending event: %v", err))
	}
	return nil
}

// SendEventPartition Send a single message into a specific partition
func (h *HubClient) SendEventPartition(message string, partitionID string) error {
	hub, err := eventhub.NewHubFromEnvironment(eventhub.HubWithPartitionedSender(partitionID))

	if err != nil {
		return errors.New("error creating event hub")
	}

	if err = hub.Send(h.Context, eventhub.NewEventFromString(message)); err != nil {
		return errors.New(fmt.Sprintf("error sending event: %v", err))
	}
	return nil
}

// SendBatchEvent Send a batch of messages into a random partition
func (h *HubClient) SendBatchEvent(events []*eventhub.Event) error {
	if err := h.Hub.SendBatch(h.Context, eventhub.NewEventBatchIterator(events...)); err != nil {
		return errors.New("error sending batch event")
	}
	return nil
}

// SendBatchEventArray Send a batch of messages into a random partition using bidirectional array
func (h *HubClient) SendBatchEventArray(data [][]byte) error {
	events := make([]*eventhub.Event, len(data))

	for i, bytes := range data {
		events[i] = &eventhub.Event{Data: bytes}
	}

	if err := h.Hub.SendBatch(h.Context, eventhub.NewEventBatchIterator(events...)); err != nil {
		return errors.New("error sending batch event")
	}
	return nil
}

// ReceiveEventChannel Receive messages from a partition from the latest offset
func (h *HubClient) ReceiveEventChannel() (events <-chan *eventhub.Event) {
	runtimeInfo, err := h.Hub.GetRuntimeInformation(h.Context)

	if err != nil {
		panic(err)
	}

	eventCh := make(chan *eventhub.Event)
	var wg sync.WaitGroup

	handler := func(ctx context.Context, event *eventhub.Event) error {
		eventCh <- event
		return nil
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		wg.Add(1)

		go func(partitionID string) {
			defer wg.Done()

			listenerHandle, err := h.Hub.Receive(h.Context, partitionID, handler, eventhub.ReceiveWithConsumerGroup("$Default"), eventhub.ReceiveWithLatestOffset())
			if err != nil {
				log.Printf("failed to start listener for partition %s: %v", partitionID, err)
				os.Exit(1)
			}

			defer func() {
				if err := listenerHandle.Close(h.Context); err != nil {
					log.Printf("failed to close listener for partition %s: %v", partitionID, err)
				}
			}()

			for {
				select {
				case event := <-listenerHandle.Done():
					log.Printf("partition %s listener closed: %v", partitionID, event)
				case <-h.Context.Done():
					return
				}
			}
		}(partitionID)
	}

	go func() {
		wg.Wait()
		close(eventCh)
	}()

	return eventCh
}

// ReceiveEventChannelPartition Receive messages from a specific partition from the latest offset
func (h *HubClient) ReceiveEventChannelPartition(consumerGroup, partitionID string) (events <-chan *eventhub.Event, err error) {
	eventCh := make(chan *eventhub.Event)
	var wg sync.WaitGroup

	handler := func(ctx context.Context, event *eventhub.Event) error {
		eventCh <- event
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		listenerHandle, err := h.Hub.Receive(h.Context, partitionID, handler, eventhub.ReceiveWithConsumerGroup(consumerGroup), eventhub.ReceiveWithLatestOffset())

		if err != nil {
			log.Printf("failed to start listener for partition %s: %v", partitionID, err)
			os.Exit(1)
		}

		defer func() {
			if err := listenerHandle.Close(h.Context); err != nil {
				log.Printf("failed to close listener for partition %s: %v", partitionID, err)
			}
		}()

		for {
			select {
			case event := <-listenerHandle.Done():
				log.Printf("partition %s listener closed: %v", partitionID, event)
			case <-h.Context.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(eventCh)
	}()

	return eventCh, nil
}
