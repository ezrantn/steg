package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a hidden payload from a PNG image",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")
		key, _ := cmd.Flags().GetString("key")
		offset, _ := cmd.Flags().GetString("offset")

		if input == "" {
			log.Fatal("Input file must be specified")
		}

		payload, err := decodeImage(input, output, key, offset)
		if err != nil {
			log.Fatal("Decoding error:", err)
		} else {
			fmt.Printf("Decoded payload: %s\n", payload)
		}
	},
}

func init() {
	decodeCmd.Flags().StringP("offset", "", "", "Offset for decoding data")
	rootCmd.AddCommand(decodeCmd)
}
