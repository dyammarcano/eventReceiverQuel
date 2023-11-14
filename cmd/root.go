package cmd

import (
	"context"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/dyammarcano/eventReceiverQuel/internal/azure_helper/event_hub"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	partitionFlag string
	consumerGroup string
	strConnFlag   string
)

var rootCmd = &cobra.Command{
	Use:   "reveiver",
	Short: "hacks tools for event hub to receive messages",
	Run:   runReceiver,
}

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rootCmd.ExecuteContext(ctx)
	cobra.CheckErr(err)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().StringVar(&consumerGroup, "consumer-group", "$Default", "consumer group to use")
	rootCmd.Flags().StringVar(&partitionFlag, "partition-id", "", "partition to use")
	rootCmd.Flags().StringVar(&strConnFlag, "connection-string", "", "connection string to event hub (SAS token)")
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Environment variables not found in: %s\n", viper.ConfigFileUsed())
	}
}

func runReceiver(cmd *cobra.Command, args []string) {
	if strConnFlag == "" {
		if err := cmd.Help(); err != nil {
			log.Printf("error printing help: %s", err.Error())
		}
		os.Exit(1)
	}

	// Create a context and a cancel function
	ctxCancel, cancel := context.WithCancel(context.Background())

	// Create a context with a value
	ctxValue := context.WithValue(ctxCancel, "azure", event_hub.AzureCredentials{
		ConnectionString: &strConnFlag,
		EventHub: event_hub.EventHub{
			ConsumerGroup: consumerGroup,
		},
	})

	clientHub, err := event_hub.NewHubClient(ctxValue)
	if err != nil {
		log.Printf("error getting event hub client: %s", err.Error())
		os.Exit(1)
	}

	defer func(Hub *eventhub.Hub, ctx context.Context) {
		if err := Hub.Close(ctx); err != nil {
			log.Printf("error closing event hub client: %s", err.Error())
		}
	}(clientHub.Hub, clientHub.Context)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// partition flag
	if partitionFlag == "" {
		chEvent := clientHub.ReceiveEventChannel()

		message(ctxValue)

		go func(chEvent <-chan *eventhub.Event) {
			for event := range chEvent {
				log.Printf("event received:\n%s", event.Data)
			}
		}(chEvent)

		<-sigCh // Wait for a termination signal
		log.Println("received termination signal, shutting down gracefully...")
		os.Exit(0)
	}

	chEvent, err := clientHub.ReceiveEventChannelPartition(consumerGroup, partitionFlag)
	if err != nil {
		log.Printf("error receiving event: %s", err.Error())
		os.Exit(1)
	}

	message(ctxValue)

	go func(chEvent <-chan *eventhub.Event) {
		for event := range chEvent {
			log.Printf("event received:\n%s", event.Data)
		}
	}(chEvent)

	<-sigCh // Wait for a termination signal
	log.Println("received termination signal, shutting down gracefully...")

	// Call the cancel function to cancel the context
	cancel()
}

func message(ctx context.Context) {
	config, ok := ctx.Value("azure").(event_hub.AzureCredentials)
	if !ok {
		log.Println("error getting azureConfig from context")
		os.Exit(1)
	}

	azure := event_hub.SplicConnectionString(*config.ConnectionString)

	log.Println("receiving events...")
	fmt.Printf(`
topic name:    %s
account name: %s

press CTRL+C to stop receiving events and exit

`, azure.EventHub.TopicName, azure.EventHub.AccountName)
}
