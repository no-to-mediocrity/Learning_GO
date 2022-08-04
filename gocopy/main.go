package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var buffersize *int64 // Global variable so we don't have to pass it to function

func main() {
	pathfrom := flag.String("from", "", "Location of a file")
	pathto := flag.String("to", "", "Copy destination")
	offset := flag.Int64("offset", 0, "Offset")
	limit := flag.Int64("limit", 0, "Limit")
	overwrite := flag.Bool("overwrite", false, "Overwrite")
	buffersize = flag.Int64("buffersize", 4096, "Buffer size")
	var fsize int64
	flag.Parse()
	if *pathfrom == "" {
		log.Println("Please provide the source file")
		return
	}
	if *overwrite == false {
		file, err := os.Open(*pathto)
		if err == nil {
			log.Println("Destination file already exists:", *pathto)
			err := file.Close()
			if err != nil {
				log.Println("OW: False, Error in file.Close():", err)
			}
			return
		}
	}
	file, err := os.Open(*pathfrom)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("OW: True, Error in defer file.Close():%v, %p\n", err, file)
		}
	}(file)
	if err != nil {
		log.Println("Error: os.Open, destination file:", *pathfrom, err)
		//handle the error properly
		return
	}

	fi, err := file.Stat()
	if err != nil {
		log.Println("Error file.Stat, destination file:", *pathfrom, err)
		return
	}
	if fsize <= 0 {
		fsize = fi.Size() - *offset
	} else {
		log.Println("File without a size!", *pathfrom)
		return
	}
	if *limit == 0 {
		if *offset > fsize {
			log.Printf("The offset (%v) is greater than the file size (%v, %v)\n", *offset, *pathfrom, fsize)
			return
		}
	} else {
		if *limit > fsize {
			log.Printf("The limit (%v) is greater than the number of bytes to copy (%v)\n", *limit, fsize)
			return
		}
		fsize = *limit
	}
	//Path suggestion
	if *pathto == "" {
		autopath, _ := PathSuggestion(*pathfrom)
		autopathmsg := "Destination path not provided. Use the path " + autopath + "?"
		c := askForConfirmation(autopathmsg)
		if c {
			*pathto = autopath
		} else {
			return
		}
	}
	copyerr := Filecopy(*pathfrom, *pathto, fsize, *offset)
	if copyerr != nil {
		log.Printf("%v", copyerr)
		return
	}

}

func Filecopy(pathfrom, pathto string, limit, offset int64) error {
	sourceFileStat, err := os.Stat(pathfrom)
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", pathfrom)
	}

	source, err := os.Open(pathfrom)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(pathto)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}
	buf := make([]byte, *buffersize)
	bar := progressbar.NewOptions(int(limit+offset), progressbar.OptionShowBytes(true), progressbar.OptionSetDescription("Copying in progress:"))
	maxIterations := (offset + limit) / *buffersize
	offsetIterations := offset / *buffersize
	var iterations int64
	for {
		if iterations > maxIterations {
			break
		}
		if iterations == maxIterations {
			//Changing the buffer size to trim the file to the limit set by user
			lastbuffer := (offset + limit) - (maxIterations * *buffersize)
			buf = make([]byte, lastbuffer)
		}
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if offset == 0 || iterations > offsetIterations {
			if _, err := destination.Write(buf[:n]); err != nil {
				return err
			}
		}
		if offset > 0 && iterations == offsetIterations {
			//Getting the head of the file according to the user-provided offset
			pos := offset - offsetIterations**buffersize
			if _, err := destination.Write(buf[pos:n]); err != nil {
				return err
			}
		}
		if err := bar.Add(n); err != nil {
			return err
		}
		if n == 0 {
			break
		}
		iterations++
	}
	return err
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s \n[y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func PathSuggestion(pathfrom string) (string, error) {
	dir, file := path.Split(pathfrom)
	ext := path.Ext(file)
	file = strings.TrimRight(file, ext)
	copycount := 2
	var autopath string
	for {
		autopath = dir + file + "(" + strconv.Itoa(copycount) + ")" + ext
		file, err := os.Open(autopath)
		defer func(file *os.File) error {
			err := file.Close()
			if err != nil {
				return err
			}
			return err
		}(file)
		if err == nil {
			copycount++
		} else {
			break
		}
	}
	return autopath, nil
}
