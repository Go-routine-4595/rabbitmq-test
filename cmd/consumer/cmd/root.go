/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/rabbitmq-test/cmd/consumer/model"
	"github.com/spf13/cobra"
)

// Declare a variable to store the AMQP connection string
var amqpConnectionString string
var key string
var queue string
var exchange string
var durable bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Connect to RabbitMQ and consume data",
	Long: `Connect to RabbitMQ and consume data like:

provide the RabbitMQ server and credential, and provide the queue to co consume form and key if needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		model.Consume(amqpConnectionString, exchange, queue, key, durable)
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.consumer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&amqpConnectionString, "amqp", "a", "", "AMQP connection string (required)")
	rootCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "Key for the operation")
	rootCmd.PersistentFlags().StringVarP(&queue, "queue", "q", "", "Queue name to use (required)")
	rootCmd.PersistentFlags().StringVarP(&exchange, "exchange", "e", "", "Exchange name to use (required)")
	rootCmd.PersistentFlags().BoolVarP(&durable, "durable", "d", true, "Used if the queue is transient")
	// Make these flags required
	rootCmd.MarkPersistentFlagRequired("amqp")
	rootCmd.MarkPersistentFlagRequired("queue")
	rootCmd.MarkPersistentFlagRequired("exchange")
}
