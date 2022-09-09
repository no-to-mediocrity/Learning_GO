package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	start "learning_go/goenv/start"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	dirpath, cmd_, args := start.Args()
	myEnv, err := ReadDir(dirpath)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}
	//fmt.Println(myEnv, err, cmd_, args)
	for key, value := range myEnv {
		err := os.Setenv(key, value)
		if err != nil {
			log.Printf("Error:%v\n", err)
		}
	}
	var args_ string
	if len(args) > 3 {
		for _, value := range args {
			args_ = value + " " + args_
		}
	}
	cmd := exec.Command(cmd_, args_)
	stdoutStderr, err_ := cmd.CombinedOutput()
	if err_ != nil {
		log.Fatal(err_)
	}
	fmt.Printf("%s", stdoutStderr)

}

func ReadDir(dir string) (map[string]string, error) {
	if dir == "" {
		dir_, err := os.Getwd()
		if err != nil {
			return nil, err
		} else {
			dir = dir_
		}
	}
	m, err := LoopDir(dir)
	if err != nil {
		log.Println(err)
	}
	if m == nil {
		fmt.Println("No configuration files have been found")
		return nil, errors.New("no configuration files have been found")
	}
	var myEnv map[string]string
	if len(m) == 1 {
		myEnv, err = godotenv.Read(m[0])
		return myEnv, err
	}
	fmt.Println("Multiple configuration files have been found:")
	g := ChooseEnv(m)
	if g == "all" {
		for _, t := range m {
			myEnv, err = godotenv.Read(t)
		}
	} else {
		myEnv, err = godotenv.Read(g)
	}
	return myEnv, err
}

func RunCmd(cmd []string, env map[string]string) int {
	return 0
}

func LoopDir(dir string) ([]string, error) {
	var m []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".env" {
			switch env := runtime.GOOS; env {
			case "windows":
				m = append(m, dir+"\\"+f.Name()) // don't need to check for existing as the FS would not allow it anyway
			default:
				m = append(m, dir+"/"+f.Name()) // same
			}
		}
	}
	return m, nil
}

func ChooseEnv(x []string) string {
	var i int = 1
	for _, k := range x {
		fmt.Printf("%v) %v\n", i, filepath.Base(k))
		i++
	}
	fmt.Println("Please enter an .env file number or type \"all\" to include all files")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	m, err := strconv.Atoi(input.Text())
	if err != nil {
		if strings.ToLower(input.Text()) == "all" {
			return "all"
		} else {
			fmt.Printf("%v is not a number\n", input.Text())
			os.Exit(2)
		}
	}
	return x[m-1]
}
