package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	return
}

func main() {
	myEnv, err := ReadDir("")
	fmt.Println(myEnv, err)
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
		fmt.Println(err)
	}
	if m == nil {
		fmt.Println("No configurations files have been found")
		return nil, errors.New("No configurations files have been found")
	}
	fmt.Println("The following configurations have been found:")
	g := ChooseEnv(m)
	var myEnv map[string]string
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
	i := 1
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
		}
		fmt.Printf("%v is not a number", input.Text())
		ChooseEnv(x)
	}
	return x[m-1]
}

/*
Домашнее задание
Утилита envdir
Цель: Реализовать утилиту envdir на Go.
 envdir  - runs another program with environment modified according to files in a specified
       directory.

Эта утилита позволяет запускать программы получая переменные окружения из определенной директории. См man envdir Пример go-envdir /path/to/evndir command arg1 arg2
Завести в репозитории отдельный пакет (модуль) для этого ДЗ
Реализовать функцию вида  , которая сканирует указанный каталог и возвращает все переменные окружения, определенные в нем.
Реализовать функцию вида RunCmd(cmd []string, env map[string]string) int , которая запускает программу с аргументами (cmd) c переопределнным окружением.
Реализовать функцию main, анализирующую аргументы командной строки и вызывающую ReadDir и RunCmd

Протестировать утилиту.
Тестировать можно утилиту целиком с помощью shell скрипта, а можно написать unit тесты на отдельные функции.
Критерии оценки: Стандартные потоки ввода/вывода/ошибок должны пробрасываться в вызываемую программу.
Код выхода утилиты envdir должен совпадать с кодом выхода программы.
Код должен проходить проверки go vet и golint
У преподавателя должна быть возможность скачать и установить пакет с помощью go get / go install

 См man envdir Пример go-envdir /path/to/evndir command arg1 arg2
Завести в репозитории отдельный пакет (модуль) для этого ДЗ
Реализовать функцию вида ReadDir(dir string) (map[string]string, error), которая сканирует указанный каталог и возвращает все переменные окружения, определенные в нем.
Реализовать функцию вида RunCmd(cmd []string, env map[string]string) int , которая запускает программу с аргументами (cmd) c переопределнным окружением.
Реализовать функцию main, анализирующую аргументы командной строки и вызывающую ReadDir и RunCmd
*/

