package util

import (
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) net.IP {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		ipStr := strings.TrimSpace(parts[0])
		if ipStr != "" {
			return net.ParseIP(ipStr)
		}
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return net.ParseIP(realIP)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return net.ParseIP(r.RemoteAddr)
	}
	return net.ParseIP(host)
}
