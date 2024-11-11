package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

func writeData(r *bytes.Reader, outputPath string, newChunk []byte, offset string) error {
	// Parse the offset from string to int64
	offsetInt, err := parseOffset(offset)
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}

	// Open the output file for reading and writing, or create it if it doesn’t exist
	outputFile, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("could not open output file: %v", err)
	}
	defer outputFile.Close()

	// Validate that offset is within the file size
	if offsetInt >= r.Size() {
		return fmt.Errorf("offset is beyond the end of file")
	}

	// Reset the reader to the beginning of the file
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error resetting reader position: %v", err)
	}

	// Write the data before the offset
	if err := copyData(r, outputFile, offsetInt); err != nil {
		return err
	}

	// Write the new chunk data at the offset position
	if _, err = outputFile.Write(newChunk); err != nil {
		return fmt.Errorf("error writing new chunk data: %v", err)
	}

	// If decoding, skip the old chunk in the input file
	if _, err = r.Seek(int64(len(newChunk)), io.SeekCurrent); err != nil {
		return fmt.Errorf("error skipping old chunk: %v", err)
	}

	// Copy the remaining data from the original file
	if _, err = io.Copy(outputFile, r); err != nil {
		return fmt.Errorf("error copying remaining data: %v", err)
	}

	fmt.Printf("Success: %s created\n", outputPath)
	return nil
}

// parseOffset parses the offset, supporting both hex (0x-prefixed) and decimal formats.
func parseOffset(offset string) (int64, error) {
	if len(offset) > 2 && offset[:2] == "0x" {
		return strconv.ParseInt(offset[2:], 16, 64)
	}
	return strconv.ParseInt(offset, 10, 64)
}

// copyData copies data from the reader to the writer up to a specified number of bytes.
func copyData(r *bytes.Reader, w *os.File, bytesToCopy int64) error {
	buffer := make([]byte, bytesToCopy)
	n, err := r.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading bytes: %v", err)
	}
	if int64(n) != bytesToCopy {
		return fmt.Errorf("could not read the expected number of bytes")
	}

	if _, err := w.Write(buffer); err != nil {
		return fmt.Errorf("error writing bytes: %v", err)
	}
	return nil
}
