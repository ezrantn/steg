package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "steg",
	Short: "Steg is a CLI tool for encoding and decoding hidden messages in PNG images using steganography",
	Long:  "Steg allows you to encode and decode hidden messages within PNG images using steganography, with optional XOR encryption for the payload.",
}

// Execute runs the root command and initializes all subcommands
func execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetUsageTemplate(`Usage:
   [flags] [command]

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Global Flags:
{{.PersistentFlags.FlagUsages | trimTrailingWhitespaces}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
`)
	// Local flags for this command
	rootCmd.Flags().StringVar(&offset, "offset", "", "Specify the offset location for the operation.")

	// Global flags (PersistentFlags) for the command
	rootCmd.PersistentFlags().StringP("input", "i", "", "Input PNG file")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output PNG file")
	rootCmd.PersistentFlags().String("key", "", "Key for encoding/decoding")

	// Mark required global flags
	rootCmd.MarkPersistentFlagRequired("input")
	rootCmd.MarkPersistentFlagRequired("output")
}
