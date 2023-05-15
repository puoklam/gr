package main

import (
	"log"

	"github.com/puoklam/gr/cmd"
)

func main() {
	// entries, err := os.ReadDir("../gr")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// for _, entry := range entries {
	// 	fmt.Println(entry.IsDir(), entry.Name())
	// }

	if err := cmd.Exec(); err != nil {
		log.Fatal(err)
	}
}
