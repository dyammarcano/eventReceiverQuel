package event_hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-event-hubs-go/v3"
	"log"
	"os"
	"strings"
	"sync"
)

type (
	HubClient struct {
		*eventhub.Hub
		context.Context
	}

	AzureCredentials struct {
		EventHub         EventHub
		Storage          Storage
		ConnectionString *string
	}

	EventHub struct {
		AccountName     string
		SharedAccessKey string
		TopicName       string
		ConsumerGroup   string
	}

	Storage struct {
		AccountName string
		AccountKey  string
	}
)

func SplicConnectionString(connectionString string) AzureCredentials {
	var azure AzureCredentials

	for _, v := range strings.Split(connectionString, ";") {
		s := strings.Split(v, "=")

		switch s[0] {
		case "Endpoint":
			azure.EventHub.AccountName = s[1]
		case "SharedAccessKeyName":
			azure.EventHub.SharedAccessKey = s[1]
		case "EntityPath":
			azure.EventHub.TopicName = s[1]
		}
	}

	return azure
}

// NewHubClient Get a new client for event hub
func NewHubClient(ctx context.Context) (*HubClient, error) {
	config, ok := ctx.Value("azure").(AzureCredentials)
	if !ok {
		log.Println("error getting azureConfig from context")
		os.Exit(1)
	}

	if config.ConnectionString != nil {
		return newHubClientSAS(ctx, *config.ConnectionString)
	}

	connStr := fmt.Sprintf("Endpoint=sb://%s.servicebus.windows.net;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=%s;EntityPath=%s",
		config.EventHub.AccountName, config.EventHub.SharedAccessKey, config.EventHub.TopicName)

	return newHubClientSAS(ctx, connStr)
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

	consumerGroup := h.Context.Value("azure").(AzureCredentials).EventHub.ConsumerGroup

	for _, partitionID := range runtimeInfo.PartitionIDs {
		wg.Add(1)

		go func(partitionID string) {
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
