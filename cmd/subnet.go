/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"math/big"
	"net"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func toBinary(b []byte) string {
	var s strings.Builder

	for i := 0; i < len(b); i++ {
		s.WriteString(fmt.Sprintf("%08b", b[i]))

		if i < (len(b) - 1) {
			s.WriteString(" ")
		}
	}

	return s.String()
}

func subtract(a, b []byte) string {
	var c, d, e big.Int

	c.SetBytes(a)
	d.SetBytes(b)

	comp := c.Cmp(&d)
	switch comp {
	case -1:
		e.Sub(&d, &c)
	case 0:
		e = *big.NewInt(0)
	case 1:
		e.Sub(&c, &d)
	}

	return e.Add(&e, big.NewInt(1)).String()
}

func and(a, b []byte) (net.IP, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length %d does not equal length %d", len(a), len(b))
	}

	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] & b[i]
	}

	return result, nil
}

func or(a, b []byte) (net.IP, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length %d does not equal length %d", len(a), len(b))
	}

	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] | b[i]
	}

	return result, nil
}

func invert(b []byte) net.IP {
	inverted := make([]byte, len(b))

	for i := 0; i < len(b); i++ {
		inverted[i] = b[i] ^ ((2 << 7) - 1)
	}

	return inverted
}

func registerSubnetting(mux *httprouter.Router, errorChannel chan<- error) {
	mux.GET("/v4/*v4", serveV4Subnet(errorChannel))
	mux.GET("/v6/*v6", serveV6Subnet(errorChannel))
}
