package cmd

import (
	"fmt"
	"kafka-client/util"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

var (
	config           kafka.ConfigMap
	cfgFile          string
	bootstrapServers string
	group            string
)

var rootCmd = &cobra.Command{
	Use:   "kafka-client",
	Short: "you know, to talk to kafka",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(loadConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "broker config")
	rootCmd.PersistentFlags().StringVar(&bootstrapServers, "server", "", "Bootstrap Server")
	rootCmd.PersistentFlags().StringVar(&group, "group", "1", "group id")
	rootCmd.MarkFlagsMutuallyExclusive("config", "server")
	rootCmd.AddCommand(producerCdm)
	rootCmd.AddCommand(consumerCmd)
	rootCmd.AddCommand(topicCmd)
}

func loadConfig() {
	config = make(map[string]kafka.ConfigValue)

	if cfgFile == "" && bootstrapServers == "" {
		fmt.Println("You must set either config or server")
		os.Exit(1)
	}

	if cfgFile == "" {
		util.SetConfig(config, "bootstrap.servers", bootstrapServers)
		util.SetConfig(config, "group.id", group)
		fmt.Printf("que %s", config["bootstrap.servers"])
	} else {
		if !util.FileExists(cfgFile) {
			fmt.Printf("Unable to open config file, file do not exist: \"%s\"", cfgFile)
			os.Exit(1)
		}
		config = util.ReadConfig(cfgFile)
	}
}
