package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	test "github.com/sscaling/go-protobuf-test/protos"
)

func main() {
	fmt.Println("vim-go")

	simple := &test.Simple{
		A: proto.Int32(150),
	}

	data, err := proto.Marshal(simple)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	for i, x := range data {
		log.Println(fmt.Sprintf("%d: %x", i, x))
	}

	buff := bytes.NewBuffer(data)

	// Get the field number
	field := readField(buff)
	log.Printf("%#v\n", field)

	switch field.WireType {
	case 0: // varint
		log.Printf("%d\n", readVarInt(buff))

	default:
		log.Fatalln("unrecognized WireType")
	}
}

type Field struct {
	Id       int32
	WireType int32
}

func readField(buff *bytes.Buffer) Field {
	f := &Field{}

	b, _ := buff.ReadByte()
	// FIXME: Handle error
	f.WireType = int32(b & 0x7)

	// field Id is a varint
	f.Id = int32((b & 0x78) >> 3)

	if (b & 0x80) != 0 {
		log.Println("has more data")
		log.Fatalln("Not implemented")
	} else {
		log.Println("no more data")
	}

	return *f
}

func readVarInt(buff *bytes.Buffer) int64 {

	b, _ := buff.ReadByte()

	v := int64(b & 0x7F)

	chunks := 1
	for (b & 0x80) != 0 {
		b, _ = buff.ReadByte()

		v |= int64((b & 0x7F) << uint(7*chunks))
		chunks++
	}

	return v
}
