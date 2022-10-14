package main

/*
char* foo(void) { return "hello, world! from C"; }
*/
import "C"
import (
	"fmt"

	"work/gobox/utils"
)

func main() {
	fmt.Printf("hello world, running on mode: %q\n", utils.Mode())
	fmt.Println(C.GoString(C.foo()))
}
