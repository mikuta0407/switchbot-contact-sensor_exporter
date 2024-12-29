/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/mikuta0407/switchbot-contact-sensor_exporter/internal/scanner"
	"github.com/spf13/cobra"
)

// scannerCmd represents the scanner command
var scannerCmd = &cobra.Command{
	Use:   "scanner",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scanner called")
		scanner.Scan(scannerFlags.Device, time.Duration(scannerFlags.Duration)*time.Second)
	},
}

var scannerFlags struct {
	Device   string
	Duration int
}

func init() {
	rootCmd.AddCommand(scannerCmd)

	scannerCmd.Flags().StringVarP(&scannerFlags.Device, "device", "d", "default", "Implementation of ble")
	scannerCmd.Flags().IntVarP(&scannerFlags.Duration, "duration", "u", 5, "Scan duration")
}
