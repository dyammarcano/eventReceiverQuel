package cmd

import (
	"context"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/caarlos0/log"
	"github.com/dyammarcano/eventReceiverQuel/internal/azure_helper/event_hub"
	cmd2 "github.com/dyammarcano/eventReceiverQuel/internal/cmd"
	"github.com/dyammarcano/eventReceiverQuel/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	log.SetLevel(log.DebugLevel)
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
		log.Fatalf("invalid type: %s\n", v)
	}

	if err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name)); err != nil {
		cmd.Printf("error binding flag: %s\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	cfg = config.NewConfig()

	filePath := viper.GetString("config")
	if filePath == "" {
		log.Fatal("error config file flag is not present")
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
		log.Fatalf("error getting event hub client: %s", err.Error())
	}

	defer func(Hub *eventhub.Hub, ctx context.Context) {
		if err := Hub.Close(ctx); err != nil {
			log.Fatalf("error closing event hub client: %s", err.Error())
		}
	}(clientHub.Hub, clientHub.Context)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	chEvent := clientHub.ReceiveEventChannel()

	message(cmd, cfg)

	counterBool := viper.GetBool("count")

	var newCounterSvc *cmd2.CounterSvc

	if counterBool {
		newCounterSvc = cmd2.NewCounterSvc()
		go newCounterSvc.Start()
	}

	go func(chEvent <-chan *eventhub.Event) {
		counter := 1
		for event := range chEvent {
			if !counterBool {
				log.Infof("event # %d\n>>\n%s<<", counter, event.Data)
				return
			}
			newCounterSvc.CountEvent()
			counter++
		}
	}(chEvent)

	<-sigCh // Wait for a termination signal
	log.Info("received termination signal, shutting down gracefully...")
	os.Exit(0)
}

func message(cmd *cobra.Command, config *config.Config) {
	// clean console
	fmt.Fprintf(cmd.OutOrStdout(), "\033[H\033[2J")

	cmd.Printf(`
Endpoint=sb://%s.servicebus.windows.net;SharedAccessKey=*******;EntityPath=%s

press CTRL+C to stop receiving events and exit

`, config.EventHub.AccountName, config.EventHub.Topic)
}
