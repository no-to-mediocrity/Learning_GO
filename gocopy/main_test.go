package main

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"
)

func TestFilecopy(t *testing.T) {
	os.Args = append(os.Args, "from \"t/hello.txt\"")
	os.Args = append(os.Args, "offset 100")
	os.Args = append(os.Args, "limit 1000")
	m := fstest.MapFS{
		"t/hello.txt": {
			Data: []byte("hello, world"),
		},
		"t/hello(2).txt": {
			Data: []byte(""),
		},
	}
	for k, v := range m {
		fmt.Printf("%s -> %s\n", k, v)
	}
	main()
	for k, v := range m {
		fmt.Printf("%s -> %s\n", k, v)
	}
}
