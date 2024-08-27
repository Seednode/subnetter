/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func toBits(b []byte, errorChannel chan<- error) uint32 {
	var r uint32

	err := binary.Read(bytes.NewBuffer(b), binary.BigEndian, &r)
	if err != nil {
		errorChannel <- err

		return 0
	}

	return r
}

func toBinary(u uint32) string {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, u)

	return fmt.Sprintf("%08b %08b %08b %08b", a[0], a[1], a[2], a[3])
}

func toDottedDecimal(u uint32) string {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, u)

	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
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

	}

	addressBits := toBits(as4, errorChannel)
	maskBits := toBits(net.Mask, errorChannel)
	maskBitsInverted := maskBits ^ ((2 << 31) - 1)
	start := addressBits & maskBits
	end := addressBits | maskBitsInverted

	address := fmt.Sprintf("%s | %s", toBinary(addressBits), toDottedDecimal((addressBits)))
	mask := fmt.Sprintf("%s | %s", toBinary(maskBits), toDottedDecimal((maskBits)))
	first := fmt.Sprintf("%s | %s", toBinary(start), toDottedDecimal((start)))
	last := fmt.Sprintf("%s | %s", toBinary(end), toDottedDecimal((end)))
	total := (end - start) + 1
	usable := (end - start) - 1

	return fmt.Sprintf("Address: %s\nMask:    %s\nFirst:   %s\nLast:    %s\n\nTotal:   %d\nUsable:  %d", address, mask, first, last, total, usable)
}

func serveV4Subnet(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		resp := calculateV4Subnet(strings.TrimPrefix(p.ByName("subnet"), "/"), errorChannel)

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

func registerV4Subnetting(mux *httprouter.Router, errorChannel chan<- error) {
	mux.GET("/v4/*subnet", serveV4Subnet(errorChannel))
}
