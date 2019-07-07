package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	sourceURL = "http://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-GB"
	bingBase  = "https://bing.com"
	res4k     = "3840x2160"
)

type Meta struct {
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

func main() {
	resp, err := http.Get(sourceURL)
	must(err, fmt.Sprintf("failed to fetch source URL %s: ",
		sourceURL))

	b, err := ioutil.ReadAll(resp.Body)
	must(err)

	meta := new(Meta)

	must(json.Unmarshal(b, meta))

	if len(meta.Images) == 0 {
		must(
			fmt.Errorf("expected at least one image to be returned, got=%d",
				len(meta.Images),
			),
		)
	}

	url := strings.Replace(meta.Images[0].URL, "1920x1080",
		res4k, 1)

	resp, err = http.Get(bingBase + url)
	must(err, fmt.Sprintf("failed to fetch image %s: ",
		bingBase, url))

	b, err = ioutil.ReadAll(resp.Body)
	must(err)

	wallDir := filepath.Join(os.Getenv("HOME"), "wallpapers")
	must(os.MkdirAll(wallDir, 755))
	year, month, day := time.Now().Date()
	path := filepath.Join(
		wallDir,
		fmt.Sprintf("%d-%02d-%02d.jpg", year, month, day),
	)

	must(ioutil.WriteFile(path, b, 0644))
	cmd := exec.Command("feh", "--bg-scale", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

func must(err error, preText ...string) {
	if err != nil {
		fmt.Fprint(os.Stderr, strings.Join(preText, "")+err.Error())
		os.Exit(1)
	}
}
