package adapters

import (
	"fmt"
	"testing"
)

var mockLinks = []Link{
	{OriginalURL: "https://test.com", Hash: "12345678", Data: nil},
	{OriginalURL: "mudmap.io", Hash: "abcdefgh", Data: nil},
}

func Test_CreateShortLink(t *testing.T) {
	got := mockLinks[0].CreateShortLink()
	want := fmt.Sprintf("http://localhost:1988/%s", mockLinks[0].Hash)

	if got != want {
		t.Fatalf("wanted: %q got %q", got, want)
	}

}
