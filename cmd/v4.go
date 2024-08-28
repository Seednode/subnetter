/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func toDottedDecimal(b []byte) string {
	if len(b) != 4 {
		return ""
	}

	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func calculateV4Subnet(cidr string, errorChannel chan<- error) string {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		errorChannel <- err

		return "Invalid CIDR address\n"
	}

	as4 := ip.To4()

	if as4 == nil {
		errorChannel <- errors.New("not an ipv4 address")

		return ""
	}

	first, err := and(as4, net.Mask)
	if err != nil {
		errorChannel <- err

		return ""
	}

	last, err := or(as4, invert(net.Mask))
	if err != nil {
		errorChannel <- err

		return ""
	}

	return fmt.Sprintf("Address: %s | %s\nMask:    %s | %s\nFirst:   %s | %s\nLast:    %s | %s\n",
		toBinary(as4), as4.String(),
		toBinary(net.Mask), toDottedDecimal(net.Mask),
		toBinary(first), first,
		toBinary(last), last)
}

func serveV4Subnet(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		resp := calculateV4Subnet(strings.TrimPrefix(p.ByName("v4"), "/"), errorChannel)

		_, err := w.Write([]byte(resp + "\n"))
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
