package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
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
		log.Printf("Please provide the source file\n")
		IncorrectInput()
	}
	file, err := os.Open(*pathfrom)
	if err != nil {
		log.Printf("Source file \"%v\" does not exist\n", *pathfrom)
		os.Exit(2)
	}
	err_ := file.Close()
	if err_ != nil {
		log.Printf("Error in file.Close(*pathto):%v\n", err)
	}
	sourceFileStat, err := os.Stat(*pathfrom)
	if !sourceFileStat.Mode().IsRegular() {
		log.Printf("%v is not a regular file. NB: Operations with folders are not supported.\n", *pathfrom)
		os.Exit(1)
	}
	if *overwrite == false {
		file, err := os.Open(*pathto)
		if err == nil {
			log.Printf("Destination file \"%v\" already exists!\n You can use flag -overwrite, provide a new destination address or omit flag -to to use path suggestion\n", *pathto)
			err := file.Close()
			if err != nil {
				log.Printf("Error in file.Close(*pathto):%v, -overwrite false\n", err)
			}
			switch env := runtime.GOOS; env {
			case "windows":
				os.Exit(58)
			case "linux":
				os.Exit(17)
			default:
				os.Exit(17)
			}
		}
	}
	file, err2 := os.Open(*pathfrom)
	defer func(file *os.File) error {
		err := file.Close()
		if err != nil {
			log.Printf("Error in defer file.Close(*pathfrom):%v, %v, -overwrite true\n", err, file)
		}
		return err
	}(file)
	if err2 != nil {
		log.Printf("Error in os.Open(*pathfrom):%v, file:\"%v\"\n", err, *pathfrom)
		os.Exit(1)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Printf("Error file.Stat, destination file:\"%v\", %v\n", *pathfrom, err)
		os.Exit(1)
	}
	fsize = fi.Size() - *offset
	switch {
	case *limit < 0:
		log.Printf("The limit (%v) cannot be less than zero\n", *limit)
		IncorrectInput()
	case *offset < 0:
		log.Printf("The offset (%v) cannot be less than zero\n", *offset)
		IncorrectInput()
	case *offset > fi.Size():
		offset1, offset2 := humanizeBytes(float64(*offset))
		fsize1, fsize2 := humanizeBytes(float64(fi.Size()))
		log.Printf("The offset (%v%v) is greater than the file size (%v%v)\n", offset1, offset2, fsize1, fsize2)
		IncorrectInput()
	case *limit > fsize:
		limit1, limit2 := humanizeBytes(float64(*limit))
		fsize1, fsize2 := humanizeBytes(float64(fi.Size()))
		log.Printf("The limit (%v%v) is greater than the number of bytes to copy (%v%v)\n", limit1, limit2, fsize1, fsize2)
		IncorrectInput()
	case *limit > 0:
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
			os.Exit(0)
		}
	}
	copyerr := Filecopy(*pathfrom, *pathto, fsize, *offset)
	if copyerr != nil {
		log.Printf("%v", copyerr)
	}
}

func Filecopy(pathfrom, pathto string, limit, offset int64) error {
var buffersize_ int64
	//for test to work
	if buffersize != nil {
		buffersize_ = *buffersize
	} else {
		buffersize_ = int64(4096)
	}
	source, err := os.Open(pathfrom)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			log.Printf("Filecopy(): Error in defer file.Close(pathfrom):%v, parameters: pathfrom:%v, pathto:%v, limit:%v, offset %v\n", err, pathfrom, pathto, limit, offset)
		}
	}(source)
	destination, err := os.Create(pathto)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}
	buf := make([]byte, buffersize_)
	bar := progressbar.NewOptions(int(limit), progressbar.OptionShowBytes(true), progressbar.OptionSetDescription("Copying in progress:"))
	maxIterations := (offset + limit) / buffersize_
	offsetIterations := offset / buffersize_
	var iterations int64
	for {
		if iterations > maxIterations {
			break
		}
		if iterations == maxIterations {
			//Changing the buffer size to trim the file to the limit set by user
			lastbuffer := (offset + limit) - (maxIterations * buffersize_)
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
			pos := offset - offsetIterations*buffersize_
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
	fmt.Printf("\n")
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
	file = file[0 : len(file)-len(ext)]
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

//taken from "github.com/schollz/progressbar/v3" as it was a private function
//all credits to schollz
func humanizeBytes(s float64) (string, string) {
	sizes := []string{" B", " kB", " MB", " GB", " TB", " PB", " EB"}
	base := 1024.0
	if s < 10 {
		return fmt.Sprintf("%2.0f", s), "B"
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f, val), suffix
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func IncorrectInput() {
	switch env := runtime.GOOS; env {
	case "windows":
		os.Exit(10022)
	case "linux":
		os.Exit(22)
	default:
		os.Exit(22)
	}
}
