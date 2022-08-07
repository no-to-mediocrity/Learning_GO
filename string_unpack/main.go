package main 
 
import ( 
 "fmt" 
 "strconv" 
 "strings" 
 "unicode" 
) 
 
func main() { 
 fmt.Println(StringUnpack("aaa4b90c80d89f")) 
} 
 
func StringUnpack(x string) string { 
 var integer string 
 lastpos := 0 
 for pos, char := range x { 
  if unicode.IsDigit(char) { 
   if pos == 0 { 
    panic("The string is not correct") 
   } else { 
    integer += string(char) 
    lastpos = pos 
   } 
  } else { 
   if char != 7 && pos == lastpos+1 && len(integer) > 0 { 
    needle := integer + string(char) 
    times, _ := strconv.Atoi(integer) 
    x = strings.Replace(x, needle, strings.Repeat(string(char), times), 1) 
    integer = "" 
   } 
  } 
 } 
 return x 
}
