package start

import (
	"log"
	"os"
)

func Args() (string, string, []string) {
	var dirpath string
	var cmdpath string
	args := os.Args
	if len(args) < 3 {
		log.Printf("Not enough arguments") // add error function
		os.Exit(1)
	}
	for i := 1; i < 3; i++ {
		dirpath = args[1]
		cmdpath = args[2]
	}
	sourceFileStat, _ := os.Stat(dirpath)
	if sourceFileStat.Mode().IsRegular() {
		log.Printf("%v is not a folder.\n", dirpath)
		os.Exit(1)
	}
	return dirpath, cmdpath, args
}
