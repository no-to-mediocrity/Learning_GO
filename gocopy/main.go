/*
Домашнее задание
Копирование файлов
Цель: Реализовать утилиту копирования файлов Утилита должна принимать следующие аргументы * файл источник (From) * файл копия (To) * Отступ в источнике (Offset), по умолчанию - 0 * Количество копируемых байт (Limit), по умолчанию - весь файл из From Выводить в консоль прогресс копирования в %, например с помощью github.com/cheggaaa/pb Программа может НЕ обрабатывать файлы, у которых не известна длинна (например /dev/urandom).
Завести в репозитории отдельный пакет (модуль) для этого ДЗ
Реализовать функцию вида Copy(from string, to string, limit int, offset int) error
Написать unit-тесты на функцию Copy
Реализовать функцию main, анализирующую параметры командной строки и вызывающую Copy
Проверить установку и работу утилиты руками
Критерии оценки: Функция должна проходить все тесты
Все необходимые для тестов файлы должны создаваться в самом тесте
Код должен проходить проверки go vet и golint
У преподавателя должна быть возможность скачать, проверить и установить пакет с помощью go get / go test / go install
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/schollz/progressbar"
	"io"
	"log"
	"os"
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
	flag.Parse()
	if *pathfrom == "" {
		log.Println("Please provide the source file")
		return
	}
	if *pathto == "" {
		autopath := "Destination path not provided. Use the path"
		c := askForConfirmation(autopath)
		if c {
			*pathto = autopath
		} else {
			return
		}
	}
	if *overwrite == false {
		file, err := os.Open(*pathto)
		defer file.Close()
		if err == nil {
			log.Println("Destination file already exists:", *pathto)
			return
		}
	}
	file, err := os.Open(*pathfrom)
	defer file.Close()
	if err != nil {
		log.Println("Cannot open destination file:", *pathfrom)
		return
	}
	var fsize int64 // so the compiler won't throw an error
	if *limit == 0 {
		fi, _ := file.Stat()
		fsize = fi.Size()
	} else {
		fsize = *limit
	}
	Filecopy(*pathfrom, *pathto, fsize, *offset)
	
}

func Filecopy(pathfrom, pathto string, limit, offset int64) error {
	// fmt.Println("From:", pathfrom)
	// fmt.Println("To:", pathto)
	// fmt.Println("offset", offset)
	// fmt.Println("limit", limit)
	// fmt.Println("buffersize", *buffersize)
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
	bar := progressbar.New(int(limit))
	maxIterations := limit / *buffersize
	offsetIterations := offset / *buffersize
	var iterations int64
	for {
		if iterations > maxIterations {
			break
		}
		if iterations == maxIterations {
			//Changing the buffer size to trim the file to the limit set by user
			lastbuffer := limit - (maxIterations * *buffersize)
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
		fmt.Printf("%s [y/n]: ", s)

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
