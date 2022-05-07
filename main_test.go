package main

import (
	"testing"
)

func TestDownloadWebpage(t *testing.T) {
	srcUrl := "https://example.com/"
	directory := "store"
	downloadedFile := downloadWebpage(&srcUrl, &directory)
	if downloadedFile != nil {
		t.Logf("%s", downloadedFile.Name()+" successfuly downloaded")
	}
}
