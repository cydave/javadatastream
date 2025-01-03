package main

import (
	"bytes"
	"fmt"

	datastream "github.com/cydave/javadatastream"
)

func main() {
	buf := new(bytes.Buffer)
	w := datastream.NewWriter(buf)
	if err := w.WriteUTF("Hello World"); err != nil {
		panic(err)
	}

	r := datastream.NewReader(buf)
	s, err := r.ReadUTF()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
