package main

import "fmt"

var (
	version string = "dev"
	commit  string = "dev"
)

func main() {
	fmt.Printf("version=%s gitCommit=%s\n", version, commit)
}
