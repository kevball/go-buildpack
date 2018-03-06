package main

import (
	"fmt"
	"net/http"
)

func hello2(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "symlinks are good")
}
