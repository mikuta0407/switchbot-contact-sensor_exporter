/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/mikuta0407/switchbot-contact-sensor_exporter/internal/config"
	"github.com/mikuta0407/switchbot-contact-sensor_exporter/internal/prometheus"
	"github.com/spf13/cobra"
)

// exporterCmd represents the exporter command
var exporterCmd = &cobra.Command{
	Use:   "exporter",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exporter called")
		err := config.Init(exporterFlags.configFile)
		if err != nil {
			log.Fatalf("Failed to read configuration. %v", err)
		}
		config.BTScanDuration = exporterFlags.btScanDuration
		config.BTScanInterval = exporterFlags.btScanInterval
		config.Device = exporterFlags.device
		prometheus.Start(exporterFlags.httpListenAddress)
	},
}

var exporterFlags struct {
	configFile        string
	httpListenAddress string
	device            string
	btScanDuration    time.Duration
	btScanInterval    time.Duration
}

func init() {
	rootCmd.AddCommand(exporterCmd)

	exporterCmd.Flags().StringVarP(&exporterFlags.configFile, "configFile", "c", "/etc/switchbot-contact-sensor_exporter/config.json", "Path to config file.")
	exporterCmd.Flags().StringVar(&exporterFlags.httpListenAddress, "httpListenAddress", "0.0.0.0:9353", "Address to bind to.")
	exporterCmd.Flags().StringVar(&exporterFlags.device, "device", "default", "Implementation of ble")
	exporterCmd.Flags().DurationVar(&exporterFlags.btScanDuration, "btScanDuration", time.Duration(2)*time.Second, "Duration in seconds for which exporter listens to sensor data")
	exporterCmd.Flags().DurationVar(&exporterFlags.btScanInterval, "btScanInterval", time.Duration(500)*time.Millisecond, "How often should exporter run sensor data listener")

}
