/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func multiFormat(b []byte) string {
	return fmt.Sprintf("%s | %s | %s", toBinary(b), toColonedHex(b), b)
}

func toColonedHex(b []byte) string {
	if len(b) != 16 {
		return ""
	}

	return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		b[0], b[1], b[2], b[3],
		b[4], b[5], b[6], b[7],
		b[8], b[9], b[10], b[11],
		b[12], b[13], b[14], b[15])
}

func calculateV6Subnet(cidr string, errorChannel chan<- error) string {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return "Not valid CIDR notation.\n"
	}

	as4 := ip.To4()

	if as4 != nil {
		return "Not a valid IPv6 address.\n"
	}

	first, err := and(ip, net.Mask)
	if err != nil {
		errorChannel <- err

		return ""
	}

	last, err := or(ip, invert(net.Mask))
	if err != nil {
		errorChannel <- err

		return ""
	}

	return fmt.Sprintf("Address: %s\nMask:    %s\nFirst:   %s\nLast:    %s\nTotal:   %s\n",
		multiFormat(ip),
		multiFormat(net.Mask),
		multiFormat(first),
		multiFormat(last),
		subtract(first, last))
}

func serveV6Subnet(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		resp := calculateV6Subnet(strings.TrimPrefix(p.ByName("v6"), "/"), errorChannel)

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
