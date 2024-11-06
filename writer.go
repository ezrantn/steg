package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func writeData(r *bytes.Reader, c *cmdLineOpts, b []byte) {
	var offset int64
	var err error

	// Check if the offset is in hex format (starts with 0x)
	if len(c.Offset) > 2 && c.Offset[:2] == "0x" {
		offset, err = strconv.ParseInt(c.Offset[2:], 16, 64)
	} else {
		offset, err = strconv.ParseInt(c.Offset, 10, 64)
	}

	if err != nil {
		log.Fatal("ParseInt failed: ", err)
	}

	// Get the total size of the input file
	totalSize := r.Size()
	if offset >= totalSize {
		log.Fatal("Offset is beyond the end of file")
	}

	w, err := os.OpenFile(c.Output, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("Fatal: Problem writing to the output file!")
	}
	defer w.Close()

	// Reset reader position
	_, err = r.Seek(0, 0)
	if err != nil {
		log.Fatal("Error seeking to start: ", err)
	}

	// Copy the first part of the file
	var buff = make([]byte, offset)
	n, err := r.Read(buff)
	if err != nil && err != io.EOF {
		log.Fatal("Error reading bytes: ", err)
	}
	if n != len(buff) {
		log.Fatal("Could not read enough bytes")
	}

	// Write the first part
	_, err = w.Write(buff)
	if err != nil {
		log.Fatal("Error writing initial bytes: ", err)
	}

	// Write the new chunk data
	_, err = w.Write(b)
	if err != nil {
		log.Fatal("Error writing payload bytes: ", err)
	}

	// Skip the old chunk when decoding
	if c.Decode {
		_, err = r.Seek(int64(len(b)), 1)
		if err != nil {
			log.Fatal("Error seeking past old chunk: ", err)
		}
	}

	// Copy the rest of the file
	_, err = io.Copy(w, r)
	if err != nil {
		log.Fatal("Error copying remaining data: ", err)
	}

	fmt.Printf("Success: %s created\n", c.Output)
}
