package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/brianvoe/gofakeit"
)

func TestCopyFunc(t *testing.T) {
	//Arrange
	//os.Args = append(os.Args, "buffersize 10") //for testing purposes
	dir, err := os.MkdirTemp("", "gotest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	file := filepath.Join(dir, "test.txt")
	copy := filepath.Join(dir, "test_copy.txt")
	copy_offset := filepath.Join(dir, "test_copy_offset.txt")
	copy_limit := filepath.Join(dir, "test_copy_limit.txt")
	copy_offset_limit := filepath.Join(dir, "test_copy_offset_limit.txt")
	offset_control := filepath.Join(dir, "test_offset_control.txt")
	limit_control := filepath.Join(dir, "test_limit_control.txt")
	offset_limit_control := filepath.Join(dir, "test_offset_limit_control.txt")
	testcontent := gofakeit.Paragraph(100, 10, 10, " ")
	testcontent_bytes := []byte(testcontent)
	//size  is expected to be ~65 Kb, so the limit is set to 100Kb, the funciton is expected to work till EOF
	if err := os.WriteFile(file, testcontent_bytes, 0666); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(offset_control, testcontent_bytes[10000:], 0666); err != nil { //Slicing a control file with an offset
		log.Fatal(err)
	}
	if err := os.WriteFile(limit_control, testcontent_bytes[:60000], 0666); err != nil { //Slicing a control file with a limit
		log.Fatal(err)
	}
	if err := os.WriteFile(offset_limit_control, testcontent_bytes[10000:60000], 0666); err != nil { //Slicing a control file with both
		log.Fatal(err)
	}
	// Act
	fmt.Println("NB! Statusbar would not work correctly in a test. Don't worry, it does not affect the copy function.")
	fmt.Println("Copying without params")
	Filecopy(file, copy, 100000, 0) //Normal copying
	fmt.Println("Copying with offset")
	Filecopy(file, copy_offset, 100000, 10000) //Copying with offset
	fmt.Println("Copying with limit")
	Filecopy(file, copy_limit, 60000, 0) //Limiting the number of bytes to copy
	fmt.Println("Copying with limit and offset")
	Filecopy(file, copy_offset_limit, 50000, 10000) // Limiting and setting offset
	// 50000 bytes starting from 10000 offset
	//Assert
	fmt.Println("Asserting copy (no parameters) behavior")
	want, _ := os.ReadFile(file)
	output, _ := os.ReadFile(copy)
	if !bytes.Equal(output, want) {
		t.Errorf("\n==== bad copy:\n%s\n==== control file:\n%s\n", copy, file)
	}
	fmt.Println("Asserting copy behavior with offset")
	want, _ = os.ReadFile(offset_control)
	output, _ = os.ReadFile(copy_offset)
	if !bytes.Equal(output, want) {
		t.Errorf("\n==== bad copy:\n%s\n==== control file:\n%s\n", copy_offset, offset_control)
	}
	fmt.Println("Asserting copy behavior with set limit")
	want, _ = os.ReadFile(limit_control)
	output, _ = os.ReadFile(copy_limit)
	if !bytes.Equal(output, want) {
		t.Errorf("\n==== bad copy:\n%s\n==== control file:\n%s\n", copy_limit, limit_control)
	}
	fmt.Println("Asserting copy behavior with offset and limit")
	want, _ = os.ReadFile(offset_limit_control)
	output, _ = os.ReadFile(copy_offset_limit)
	if !bytes.Equal(output, want) {
		t.Errorf("\n==== bad copy:\n%s\n==== control file:\n%s\n", copy_offset_limit, offset_limit_control)
	}
}
