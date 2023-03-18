package main

import (
	"fmt"
	"lesson28/pkg/storage"
	"os"
	"os/signal"
	"syscall"
)

var (
	sc = storage.StudentClass{}
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		for _, student := range sc {
			fmt.Printf("%+v\n", sc.Get(student.Name))
		}
		os.Exit(0)
	}()

	sc.Form()
}
