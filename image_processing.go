package main

import (
	"fmt"
	"os"
	"strconv"
)

func encodeImage(inputPath, outputPath, payload, key, offset string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("could not open input file: %v", err)
	}
	defer inputFile.Close()

	reader, err := preProcessImage(inputFile)
	if err != nil {
		return fmt.Errorf("error processing image: %v", err)
	}

	var meta MetaChunk
	meta.IsPNG(reader)

	meta.Chk.Data = xorEncode([]byte(payload), key)
	meta.Chk.Type = meta.strToInt("tEXt")
	meta.Chk.Size = meta.createChunkSize()
	meta.Chk.CRC = meta.createChunkCRC()

	encodedData := meta.marshalData().Bytes()
	if err := writeData(reader, outputPath, encodedData, offset); err != nil {
		return fmt.Errorf("error writing encoded data: %v", err)
	}

	fmt.Println("Encoding successful!")
	return nil
}

func decodeImage(inputPath, outputPath, key, offset string) (string, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("could not open input file: %v", err)
	}
	defer inputFile.Close()

	reader, err := preProcessImage(inputFile)
	if err != nil {
		return "", fmt.Errorf("error processing image: %v", err)
	}

	var meta MetaChunk
	var offsetInt int64
	if len(offset) > 2 && offset[:2] == "0x" {
		offsetInt, err = strconv.ParseInt(offset[2:], 16, 64)
	} else {
		offsetInt, err = strconv.ParseInt(offset, 10, 64)
	}
	if err != nil {
		return "", fmt.Errorf("invalid offset: %v", err)
	}

	// Seek to the specified offset
	if _, err := reader.Seek(offsetInt, 0); err != nil {
		return "", fmt.Errorf("error seeking to offset: %v", err)
	}

	meta.readChunk(reader)
	if meta.Chk.Size == 0 {
		return "", fmt.Errorf("invalid chunk size at offset")
	}

	// Decode data
	decodedData := xorDecode(meta.Chk.Data, key)
	return string(decodedData), nil
}
