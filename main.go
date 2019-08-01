package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	sourceURL = "https://api.unsplash.com/photos/random"
)

var (
	ClientID   string
	Collection string
)

type Image struct {
	URLs URLs `json:"urls"`
}

type URLs struct {
	Full string `json:"full"`
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&ClientID, "client-id", "i", "", "Client ID used for authorization.")
	RootCmd.PersistentFlags().StringVarP(&Collection, "collection", "c", "827743", "Collection ID from Unsplash.")
}

var RootCmd = &cobra.Command{
	Use:   "wallpaper",
	Short: "Binary to download a random photo and set as wallpaper from an Unsplash Collection.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ClientID == "" {
			return errors.New("--client-id is a flag that is required to be set.")
		}

		if Collection == "" {
			return errors.New("--collection must be set to some value.")
		}

		url := fmt.Sprintf("https://api.unsplash.com/photos/random?collections=%s", Collection)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to build request for unsplash: %s", err)
		}

		req.Header.Add(
			"Authorization",
			fmt.Sprintf("Client-ID %s", ClientID),
		)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to fetch source URL %s: %s", url, err)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		jimage := new(Image)
		if err := json.Unmarshal(b, jimage); err != nil {
			return err
		}

		resp, err = http.Get(jimage.URLs.Full)
		if err != nil {
			return fmt.Errorf("failed to fetch image %s: %s",
				jimage.URLs.Full, err)
		}

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		wallDir := filepath.Join(os.Getenv("HOME"), "wallpapers")
		if err := os.MkdirAll(wallDir, 755); err != nil {
			return err
		}

		year, month, day := time.Now().Date()
		path := filepath.Join(
			wallDir,
			fmt.Sprintf("%d-%02d-%02d.jpg", year, month, day),
		)

		if err := ioutil.WriteFile(path, b, 0644); err != nil {
			return err
		}

		fehCmd := exec.Command("feh", "--bg-scale", path)
		fehCmd.Stdout = os.Stdout
		fehCmd.Stderr = os.Stderr

		return fehCmd.Run()
	},
}

func must(err error, preText ...string) {
	if err != nil {
		fmt.Fprint(os.Stderr, strings.Join(preText, "")+err.Error())
		os.Exit(1)
	}
}
