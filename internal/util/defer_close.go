package util

import (
	"log"
)

type Closeable interface {
	Close() error
}

func DeferClose(c Closeable) {
	if c != nil {
		if err := c.Close(); err != nil {
			log.Println(err)
		}
	}
}

func DeferCloseFatal(c Closeable) {
	if c != nil {
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
