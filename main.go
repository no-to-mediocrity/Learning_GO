package main

import (
 "fmt"
 "sort"
 "strings"
)

func main() {
 x := "Split slices s into all substrings separated by sep and returns a slice of the substrings between those separators.\n\nIf s does not contain sep and sep is not empty, Split returns a slice of length 1 whose only element is s.\n\nIf sep is empty, Split splits after each UTF-8 sequence. If both s and sep are empty, Split returns an empty slice.\n\nIt is equivalent to SplitN with a count of -1.\n\nTo split around the first instance of a separator, see Cut."
 sortMap(CountWords(x))

}

func CountWords(x string) map[string]int {
 for {
  x = strings.Replace(x, "\n", " ", 1)
  if strings.Contains(x, "\n") == false {
   break
  }
 }
 s := strings.Split(x, " ")
 g := getUniqueValue(s)
 return g
}

func getUniqueValue(arr []string) map[string]int {
 //Create a   dictionary of values for each element
 dict := make(map[string]int)
 for _, string := range arr {
  string = cleanString(string)
  if len(string) > 0 {
   dict[string] = dict[string] + 1
  }
 }
 // fmt.Println(dict)
 return dict
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

func sortMap(x map[string]int) {
 var a []int
 n := map[int][]string{}
 for k, v := range x {
  n[v] = append(n[v], k)
 }
 fmt.Println(n)
 for k := range n {
  a = append(a, k)
 }
 fmt.Println(a)
 sort.Sort(sort.Reverse(sort.IntSlice(a)))
 g := 0
 for _, k := range a {
  for _, s := range n[k] {
   if g < 10 {
    fmt.Printf("%s, %d\n", s, k)
    g++
   } else {
    break
   }
  }
 }
}
