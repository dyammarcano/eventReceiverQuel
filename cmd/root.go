package cmd

import (
	"context"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/dyammarcano/eventReceiverQuel/internal/azure_helper/event_hub"
	"github.com/dyammarcano/eventReceiverQuel/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfg *config.Config
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

	AddFlag(rootCmd, "config", "", "config file")
	AddFlag(rootCmd, "count", false, "counter messages")
}

// AddFlag adds a flag to the service manager, it also binds the flag to the viper instance
func AddFlag(cmd *cobra.Command, name string, defaultValue any, description string) {
	switch v := defaultValue.(type) {
	case bool:
		cmd.PersistentFlags().Bool(name, v, description)
	case string:
		cmd.PersistentFlags().String(name, v, description)
	case int, int8, int16, int32, int64:
		cmd.PersistentFlags().Int64(name, v.(int64), description)
	default:
		fmt.Printf("Invalid type: %s\n", v)
		os.Exit(1)
	}

	if err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name)); err != nil {
		cmd.Printf("Error binding flag: %s\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	cfg = config.NewConfig()

	filePath := viper.GetString("config")
	if filePath == "" {
		log.Println("Error config file flag is not present")
		os.Exit(1)
	}

	err := cfg.LoadConfigFile(filePath)
	cobra.CheckErr(err)

	err = cfg.Validate()
	cobra.CheckErr(err)
}

func runReceiver(cmd *cobra.Command, _ []string) {
	clientHub, err := event_hub.NewHubClient(cmd.Context(), cfg)
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

	chEvent := clientHub.ReceiveEventChannel()

	message(cfg)

	counterBool := viper.GetBool("count")

	go func(chEvent <-chan *eventhub.Event) {
		counter := 1
		var msg string
		for event := range chEvent {
			msg = fmt.Sprintf("event # %d\n>>\n%s\n<<", counter, event.Data)

			if counterBool {
				msg = fmt.Sprintf("event # %d\n", counter)
				counter++
			}
			log.Println(msg)
		}
	}(chEvent)

	<-sigCh // Wait for a termination signal
	log.Println("received termination signal, shutting down gracefully...")
	os.Exit(0)
}

func message(config *config.Config) {
	log.Println("receiving events...")
	fmt.Printf(`
topic name:    %s
account name: %s

press CTRL+C to stop receiving events and exit

`, config.EventHub.Topic, config.EventHub.AccountName)
}
