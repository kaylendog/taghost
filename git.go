package main

import (
	"fmt"
	"strings"
	"os"
	"github.com/spf13/viper"
	"net/http"
	"io"
	"time"
)

type ProgressReader struct {
    io.Reader
    Reporter func(r int64)
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
    n, err = pr.Reader.Read(p)
    pr.Reporter(int64(n))
    return
}


// GetAccessURL returns a concatenated URL allowing access to the target git repository.
func GetAccessURL() string {
	return fmt.Sprintf("https://%s:%s@%s", viper.GetString("git.username"), viper.GetString("git.access_token"), strings.SplitAfter(viper.GetString("git.repository_url"), "https://")[1])
}

// CheckAndClone checks if the local repository is cloned.
func CheckAndClone() {
	log.Info("Checking for asset directory...")
	path := viper.GetString("assets.path")
	_, err := os.Stat(path)
	if err != nil {
		Clone()
	} else {
		log.Warn("Current asset path is invalid - recreating...")
		os.Remove(path)
		Clone()
	}
}

// Update the local repository.
func Update() {}

// Clone the remote repository.
func Clone() {
	log.Info("Cloning remote repository...")
	
	path := viper.GetString("assets.path")

	out, err := os.Create(fmt.Sprintf("%s.zip", path))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	
	r, err := http.Get(fmt.Sprintf("%s/archive/master.zip", GetAccessURL()))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

    total := int64(0)
    pr := &ProgressReader{r.Body, func(r int64) {
        total += r
	}}
	
	done := false
	go func() {
		for !done {
			log.Infof("Progress: %d", total)
			time.Sleep(1 * time.Second)
		}
	}()

	_, err = io.Copy(out, pr)
	if err != nil  {
		panic(err)
	}

	done = true
}
