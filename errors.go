package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func parseError(err error) {
	if err != nil {
		if err, ok := err.(stackTracer); ok {
			for _, f := range err.StackTrace() {
				fmt.Printf("%+v:", f)
			}
		} else {
			fmt.Print(err)
		}
	}
}
