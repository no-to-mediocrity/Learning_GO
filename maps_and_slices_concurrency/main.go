package main

import (
 "fmt"
 "log"
 "strings"
 "sync"
 "time"
 "unicode"
)

func main() {
 start := time.Now()
 var wg sync.WaitGroup
 final := make(map[string]int)
 x : = "Provide your string here"
 wg.Add(4)
 ch := make(chan map[string]int, 10)
 work1, work2, work3, work4 := Splitstring4(x)
 go WorkReceiver(final, ch)
 go CountWords(work1, ch, &wg)
 go CountWords(work2, ch, &wg)
 go CountWords(work3, ch, &wg)
 go CountWords(work4, ch, &wg)
 wg.Wait()
 close(ch)
 fmt.Println(final)
 elapsed := time.Since(start)
 log.Println("Function took", elapsed)
}

func WorkReceiver(dict map[string]int, c chan map[string]int) {
 for work := range c {
  for string, _ := range work {
   if _, ok := dict[string]; ok {
    dict[string] = dict[string] + work[string]
   } else {
    dict[string] = work[string]
   }
  }
 }
}

func Splitstring4(x string) (string, string, string, string) {
 parts := []int{25, 50, 75}
 var positions []int
 for _, part := range parts {
  coordinates := int(len(x) * part / 100)
  for pos, letter := range x[coordinates:len(x)] {
   if unicode.IsSpace(letter) == true {
    positions = append(positions, pos+coordinates)
    break
   }
  }
 }
 return x[0:positions[0]], x[positions[0]:positions[1]], x[positions[1]:positions[2]], x[positions[2]:len(x)]
}

func CountWords(x string, c chan map[string]int, wg *sync.WaitGroup) {
 dict := make(map[string]int)
 for {
  x = strings.Replace(x, "\n", " ", 1)
  if strings.Contains(x, "\n") == false {
   break
  }
 }
 s := strings.Split(x, " ")
 for _, string := range s {
  string = cleanString(string)
  if len(string) > 0 {
   dict[string] = dict[string] + 1
  }
 }
 //fmt.Println(dict)
 c <- dict
 //fmt.Println("Done")
 wg.Done()
}

func cleanString(x string) string {
 r := []rune{44, 45, 46, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59}
 for _, pos := range x {
  char := pos
  for _, pos := range r {
   if char == pos {
    x = strings.Replace(x, string(char), "", 1)
   }
  }
 }
 return strings.ToLower(x)
}
