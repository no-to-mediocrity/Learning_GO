package main

import (
	"github.com/brianvoe/gofakeit"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFunc(t *testing.T) {
	//Arrange
	//os.Args = append(os.Args, "buffersize 10") //for testing purposes
	dir, err := os.MkdirTemp("", "gotest")
	if err != nil {
		log.Fatal(err)
	}
	//defer os.RemoveAll(dir) // clean up
	file := filepath.Join(dir, "test.txt")
	copy := filepath.Join(dir, "test_copy.txt")
	copy_offset := filepath.Join(dir, "test_copy_offset.txt")
	copy_limit := filepath.Join(dir, "test_copy_limit.txt")
	copy_offet_limit := filepath.Join(dir, "test_copy_offset_limit.txt")
	offset_control := filepath.Join(dir, "test_offset_control.txt")
	limit_control := filepath.Join(dir, "test_limit_control.txt")
	offset_limit_control := filepath.Join(dir, "test_offset_limit_control.txt")
	testcontent := gofakeit.Paragraph(100, 10, 10, " ")
	testcontent_bytes :=  []byte(testcontent)
	//size predetermined and is exactly 68247 bytes, which determines the limit
	if err := os.WriteFile(file, testcontent_bytes, 0666); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(offset_control, testcontent_bytes[10000:], 0666); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(limit_control, testcontent_bytes[:60000], 0666); err != nil {
		log.Fatal(err)
	}
		if err := os.WriteFile(offset_limit_control, testcontent_bytes[10000:60000], 0666); err != nil {
		log.Fatal(err)
	}
	// Act
	Filecopy(file, copy, 68247, 0) //Normal copying 
	Filecopy(file, copy_offset, 68247, 10000) //Copying with offset
	Filecopy(file, copy_limit, 60000, 0) //Limiting the number of bytes to copy
	Filecopy(file, copy_offet_limit, 60000, 10000) // Limiting and setting offset
	//Assert
	
}

