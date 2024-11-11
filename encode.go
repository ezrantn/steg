package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var payload string
var offset string

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode a hidden payload into a PNG image",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")
		key, _ := cmd.Flags().GetString("key")
		offset, _ := cmd.Flags().GetString("offset")
		payload, _ := cmd.Flags().GetString("payload")

		if input == "" || output == "" {
			log.Fatal("Input and output files must be specified")
		}

		err := encodeImage(input, output, payload, key, offset)
		if err != nil {
			log.Fatal("Encoding error:", err)
		} else {
			fmt.Printf("Encoded data into %s successfully!\n", output)
		}
	},
}

func init() {
	encodeCmd.Flags().StringVarP(&payload, "payload", "p", "", "Payload to encode")
	encodeCmd.Flags().StringVarP(&offset, "offset", "", "", "Offset for encoding data")
	encodeCmd.MarkFlagRequired("payload")
	rootCmd.AddCommand(encodeCmd)
}
