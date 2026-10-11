package auth

import (
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

var macWhitelist = map[string]bool{}

var macRe = regexp.MustCompile(`([0-9a-fA-F]{2}[:-]){5}[0-9a-fA-F]{2}`)

func lookupMAC(ip string) string {
	out, err := exec.Command("arp", "-a", ip).Output()
	if err != nil {
		return ""
	}
	match := macRe.FindString(string(out))
	if match == "" {
		return ""
	}

	return strings.ToLower(strings.ReplaceAll(match, "-", ":"))
}

func RequireLAN(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		mac := lookupMAC(host)
		if !macWhitelist[mac] {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
