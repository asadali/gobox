package main

// https://blog.hashbangbash.com/2014/04/linking-golang-statically/
/*
char* foo(void) { return "hello, world! from C"; }
*/
import "C"
import (
	"fmt"
	"log"

	"work/gobox/git"
	"work/gobox/utils"
)

func main() {
	fmt.Printf("hello world, running on mode: %q\n", utils.Mode())
	fmt.Println(C.GoString(C.foo()))
}

func createRepo() {
	contents := map[string]string{
		"foo/a.txt": "this is file foo/a.txt",
		"foo/b.txt": "this is file foo/b.txt",
		"bar/c.txt": "this is bar/c.txt",
	}

	repo, commitId, err := git.CreateTestRepository("refs/heads/main", contents, "init commit")
	if err != nil {
		log.Fatal(err)
	}
	defer git.CleanupTestRepository(repo)
	log.Println("commit: ", commitId)
}
