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

	"github.com/spf13/pflag"
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

func (mc *MetaChunk) validate(b *bytes.Reader) {
	var header Header

	if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
		log.Fatal(err)
	}

	bArr := make([]byte, 8)
	binary.BigEndian.PutUint64(bArr, header.Header)

	if !bytes.Equal(bArr[1:4], []byte("PNG")) {
		log.Fatal("Provided file is not a valid PNG format")
	} else {
		fmt.Println("Valid PNG so let us continue!")
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

func (mc *MetaChunk) processImage(b *bytes.Reader, c *cmdLineOpts) {
	mc.validate(b)

	count := 1
	chunkType := ""
	endChunkType := "IEND"
	for chunkType != endChunkType {
		fmt.Println("---- Chunk # " + strconv.Itoa(count) + " ----")
		offset := mc.getOffset(b)
		fmt.Printf("Chunk Offset: %#02x\n", offset)
		mc.readChunk(b)

		chunkType = mc.chunkTypeToString()
		count++
	}

	var m MetaChunk
	m.Chk.Data = []byte(c.Payload)
	m.Chk.Type = m.strToInt(c.Type)
	m.Chk.Size = m.createChunkSize()
	m.Chk.CRC = m.createChunkCRC()
	bm := m.marshalData()
	bmb := bm.Bytes()
	fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
	fmt.Printf("Payload: % X\n", m.Chk.Data)
	writeData(b, c, bmb)
}

var (
	flags = pflag.FlagSet{SortFlags: false}
	opts  cmdLineOpts
	png   MetaChunk
)

func main() {
	dat, err := os.Open(opts.Input)
	defer dat.Close()
	bReader, err := preProcessImage(dat)
	if err != nil {
		log.Fatal(err)
	}
	png.processImage(bReader, &opts)
}
