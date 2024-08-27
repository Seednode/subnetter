/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveUsage(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		data := []byte(fmt.Sprintf("subnetter v%s\n", ReleaseVersion))

		w.Header().Add("Content-Security-Policy", "default-src 'self';")

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		w.Header().Set("Content-Length", strconv.Itoa(len(data)))

		_, err := w.Write(data)
		if err != nil {
			errorChannel <- err

			return
		}

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(logDate),
				realIP(r),
				r.RequestURI)
		}
	}
}

func registerUsage(mux *httprouter.Router, errorChannel chan<- error) {
	mux.GET("/", serveUsage(errorChannel))
}
