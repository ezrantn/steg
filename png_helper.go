package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

func (mc *MetaChunk) IsPNG(b *bytes.Reader) {
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
