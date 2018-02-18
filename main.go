package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	test "github.com/sscaling/go-protobuf-test/protos"
)

func main() {
	fmt.Println("vim-go")

	simple := &test.Simple{
		A: proto.Int32(150),
		B: proto.Int64(10000),
	}

	data, err := proto.Marshal(simple)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	for i, x := range data {
		log.Println(fmt.Sprintf("%d: %x", i, x))
	}

	buff := bytes.NewBuffer(data)

	for err = readMessage(buff); err == nil; err = readMessage(buff) {

	}
	log.Fatalf("err %v\n", err)

	b, err := buff.ReadByte()
	if err != nil {
		log.Fatalf("Failed %v\n", err)
	} else {
		log.Fatalf("Next byte %X\n", b)
	}

	log.Println("Done")
}

func readMessage(buff *bytes.Buffer) error {
	// Get the field number
	field, err := readField(buff)
	if err != nil {
		return err
	}

	log.Printf("%#v\n", field)

	switch field.WireType {
	case 0: // varint
		i, err := readVarInt(buff)
		if err != nil {
			return err
		}
		log.Printf("%d\n", i)

	default:
		log.Fatalln("unrecognized WireType")
		return errors.New("unrecognised WireType")
	}

	return nil
}

type Field struct {
	Id       int32
	WireType int32
}

func readField(buff *bytes.Buffer) (Field, error) {
	f := Field{}

	b, err := buff.ReadByte()
	if err != nil {
		return f, err
	}

	f.WireType = int32(b & 0x7)

	// field Id is a varint
	f.Id = int32((b & 0x78) >> 3)

	if (b & 0x80) != 0 {
		log.Println("has more data")
		// FIXME: is this correct?
		v, err := readVarInt(buff)
		if err != nil {
			return f, err
		}
		f.Id |= int32(v << 4)
	} else {
		log.Println("no more data")
	}

	return f, nil
}

func readVarInt(buff *bytes.Buffer) (int64, error) {

	b, err := buff.ReadByte()
	if err != nil {
		return int64(0), err
	}

	v := int64(b & 0x7F)
	log.Printf("v = %v\n", v)

	chunks := 1
	for (b & 0x80) != 0 {
		b, err = buff.ReadByte()
		if err != nil {
			return int64(0), err
		}

		//		v2 := int64(b & 0x7F)
		v |= int64(b) & 0x7F << uint(7*chunks)
		log.Printf("[%d] v = %x (+ %x)\n", chunks, v, (b & 0x7F))
		chunks++
	}

	return v, nil
}
