package handlers

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetIP(t *testing.T) {
	type testIP struct {
		name     string
		request  *http.Request
		expected string
	}

	newRequest := func(remoteAddr, xRealIp string, xForwardedFor ...string) *http.Request {
		h := http.Header{}
		h.Set("X-REAL-IP", xRealIp)
		for _, address := range xForwardedFor {
			h.Set("X-FORWARDED-FOR", address)
		}

		return &http.Request{
			RemoteAddr: remoteAddr,
			Header:     h,
		}
	}

	// Create test data
	publicAddr1 := "144.12.54.87"
	publicAddr2 := "119.14.55.11"
	localAddr := "127.0.0.0"

	testData := []testIP{
		{
			name:     "No header",
			request:  newRequest(formatIPAddressWithPort(publicAddr1), ""),
			expected: publicAddr1,
		}, {
			name:     "Has X-Forwarded-For",
			request:  newRequest("", "", publicAddr1),
			expected: publicAddr1,
		}, {
			name:     "Has multiple X-Forwarded-For",
			request:  newRequest("", "", localAddr, publicAddr1, publicAddr2),
			expected: publicAddr2,
		}, {
			name:     "Has X-Real-IP",
			request:  newRequest("", publicAddr1),
			expected: publicAddr1,
		},
	}

	for _, v := range testData {
		actual, e := getIP(v.request)
		if e != nil {
			t.Errorf("error occurred: %s", e)
		}
		if v.expected != actual {
			t.Errorf("%s: expected %s got %s", v.name, v.expected, actual)
		}
	}

}

//////////////////////////////
// 		TEST UTILITIES
//////////////////////////////

// formatIPAddressWithPort appends a port number for non proxied requests (no header)
func formatIPAddressWithPort(ipaddr string) string {
	return fmt.Sprintf("%s:80", ipaddr)
}
