package main

import (
	"github.com/Beta5051/RadishCat"
)

func main() {
	rc, err := RadishCat.New("token")
	if err != nil {
		panic(err)
	}

	err = rc.Open()
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	select {}
}