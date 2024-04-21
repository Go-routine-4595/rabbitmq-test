/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/rabbitmq-test/cmd/producer/model"
	"github.com/spf13/cobra"
)

// Declare a variable to store the AMQP connection string
var amqpConnectionString string
var key string
var exchange string
var message string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "producer",
	Short: "Producer message to an RabbitMQ exchange",
	Long: `Producer expect the following:

The RabbitMQ connection string, the key if needed, the exchange and the message body, if no message is provided one will be randomly generated.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		model.Produce(amqpConnectionString, exchange, message, key)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.producer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&amqpConnectionString, "amqp", "a", "", "AMQP connection string (required)")
	rootCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "Key for routing messages (default \"\")")
	rootCmd.PersistentFlags().StringVarP(&message, "message", "m", "", "Message to send (if not a random message will be sent)")
	rootCmd.PersistentFlags().StringVarP(&exchange, "exchange", "e", "", "Exchange name to use (required)")
	// Make these flags required
	rootCmd.MarkPersistentFlagRequired("amqp")
	rootCmd.MarkPersistentFlagRequired("exchange")
}
