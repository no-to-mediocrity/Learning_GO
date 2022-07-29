package main

import (
	"flag"
	"fmt"
	"github.com/schollz/progressbar"
	"io"
	"os"
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
		fmt.Println("Please provide the source file")
		return
	}
	if *pathto == "" {
		//сделать функцию адаптивную
		//fmt.Println("Destination path not provided. Use the path %v"(), *pathfrom)
		return
	}
	if *overwrite == false {
		file, err := os.Open(*pathto)
		if err == nil {
			fmt.Println("Destination file already exists:", *pathto)
			file.Close()
			return
		}
	}
	file, err := os.Open(*pathfrom)
	if err != nil {
		fmt.Println("Cannot open destination file:", *pathfrom)
		return
	}
	var fsize int64 // so the compiler won't throw an error
	if *limit == 0 {
		fmt.Println("limit = 0", *limit)
		fi, _ := file.Stat()
		fsize = fi.Size()
	} else {
		fsize = *limit
	}
	//fmt.Println(int(fsize) / 4096)
	file.Close()
	Filecopy(*pathfrom, *pathto, fsize, *offset)

	//flag overwrite file
	//add if exist
}

func Filecopy(pathfrom, pathto string, limit, offset int64) error {
	// fmt.Println("From:", pathfrom)
	// fmt.Println("To:", pathto)
	// fmt.Println("offset", offset)
	//fmt.Println("limit", limit)
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
	//fmt.Println(offsetIterations)
	//file_iterations :=
	var iterations int64
	for {
		//	fmt.Println(iterations, offsetIterations)
		//	fmt.Println(iterations, max_iterations)
		if iterations > maxIterations {
			break
		}
		if iterations == maxIterations {
			lastbuffer := limit - (maxIterations * *buffersize)
			buf = make([]byte, lastbuffer)
		}
		n, err := source.Read(buf)
		//fmt.Println(iterations)
		if err != nil && err != io.EOF {
			return err
		}
		if offset == 0 || iterations > offsetIterations {
			if _, err := destination.Write(buf[:n]); err != nil {
				return err
			}
		}
		if offset > 0 && iterations == offsetIterations {
			pos := offset - offsetIterations**buffersize
			if _, err := destination.Write(buf[pos:n]); err != nil {
				return err
			}
		}
		bar.Add(n)
		if n == 0 {
			break
		}
		iterations++
	}
	return err
}

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
