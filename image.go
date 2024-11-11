package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"strconv"
)

// Magic Bytes
type Header struct {
	Header uint64
}

// Chunk represents a data byte chunk segment
type Chunk struct {
	Size uint32
	Type uint32
	Data []byte
	CRC  uint32
}

type MetaChunk struct {
	Chk    Chunk
	Offset int64
}

func encodeDecode(input []byte, key string) []byte {
	var bArr = make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		bArr[i] += input[i] ^ key[i%len(key)]
	}
	return bArr
}

func xorEncode(decode []byte, key string) []byte {
	return encodeDecode(decode, key)
}

func xorDecode(encode []byte, key string) []byte {
	return encodeDecode(encode, key)
}

func preProcessImage(dat *os.File) (*bytes.Reader, error) {
	stats, err := dat.Stat()
	if err != nil {
		return nil, err
	}

	var size = stats.Size()
	b := make([]byte, size)

	bufR := bufio.NewReader(dat)
	_, err = bufR.Read(b)
	bReader := bytes.NewReader(b)

	return bReader, err
}

func (mc *MetaChunk) processImage(b *bytes.Reader, c *cmdLineOpts) {
	mc.IsPNG(b)

	if (c.Offset != "") && c.Encode {
		var m MetaChunk
		m.Chk.Data = xorEncode([]byte(c.Payload), c.Key)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload Encode: % X\n", m.Chk.Data)
		writeData(b, c.Output, bmb, c.Offset)
	}

	if (c.Offset != "") && c.Decode {
		var m MetaChunk
		var offset int64
		var err error
		if len(c.Offset) > 2 && c.Offset[:2] == "0x" {
			offset, err = strconv.ParseInt(c.Offset[2:], 16, 64)
		} else {
			offset, err = strconv.ParseInt(c.Offset, 10, 64)
		}
		if err != nil {
			log.Fatal("Invalid offset:", err)
		}

		// Seek to the offset
		_, err = b.Seek(offset, 0)
		if err != nil {
			log.Fatal("Error seeking to offset:", err)
		}

		// Read the chunk data
		m.readChunk(b)
		if m.Chk.Size == 0 {
			log.Fatal("Invalid chunk size at offset")
		}

		origData := make([]byte, len(m.Chk.Data))
		copy(origData, m.Chk.Data)

		// Decode the data
		m.Chk.Data = xorDecode(m.Chk.Data, c.Key)
		m.Chk.CRC = m.createChunkCRC()

		bm := m.marshalData()
		bmb := bm.Bytes()

		fmt.Printf("Payload Original: % X\n", origData)
		fmt.Printf("Payload Decode: % X\n", m.Chk.Data)
		writeData(b, c.Output, bmb, c.Offset)
	}
}

func (mc *MetaChunk) getOffset(b *bytes.Reader) int64 {
	offset, _ := b.Seek(0, 1)
	mc.Offset = offset

	return offset
}

func (mc *MetaChunk) readChunk(b *bytes.Reader) {
	mc.readChunkSize(b)
	mc.readChunkType(b)
	mc.readChunkBytes(b, mc.Chk.Size)
	mc.readChunkCRC(b)
}

func (mc *MetaChunk) readChunkSize(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkType(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkBytes(b *bytes.Reader, cLen uint32) {
	mc.Chk.Data = make([]byte, cLen)
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkCRC(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) chunkTypeToString() string {
	h := fmt.Sprintf("%x", mc.Chk.Type)
	decoded, _ := hex.DecodeString(h)
	result := fmt.Sprintf("%s", decoded)
	return result
}

func (mc *MetaChunk) strToInt(s string) uint32 {
	t := []byte(s)
	return binary.BigEndian.Uint32(t)
}

func (mc *MetaChunk) createChunkSize() uint32 {
	return uint32(len(mc.Chk.Data))
}

func (mc *MetaChunk) createChunkCRC() uint32 {
	bytesMSB := new(bytes.Buffer)
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	return crc32.ChecksumIEEE(bytesMSB.Bytes())
}

func (mc *MetaChunk) marshalData() *bytes.Buffer {
	bytesMSB := new(bytes.Buffer)
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}

	return bytesMSB
}
