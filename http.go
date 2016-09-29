package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	p_url "net/url"
	"os"
	"time"

	"github.com/zew/logx"
)

// Get a http client
func HttpClient() *http.Client {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 1110,
		Transport: netTransport,
	}
	return netClient
}

func Request(method, url string, argKeys, argVals []string) (respBytes []byte, err error) {

	if !(method == "GET" || method == "POST") {
		logx.Fatalf("must be GET or POST; not %v", method)
	}

	if len(argKeys) != len(argVals) {
		logx.Fatalf("keys and vals must be equal size; not %v %v", len(argKeys), len(argVals))
	}

	var req *http.Request

	if method == "GET" {

		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}

		q := req.URL.Query()
		for i, key := range argKeys { // Add the other fields
			q.Add(key, argVals[i])
		}
		req.URL.RawQuery = q.Encode()

	} else if method == "POST" {

		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			return
		}

		form := p_url.Values{}
		for i, key := range argKeys { // Add the other fields
			form.Add(key, argVals[i])
		}
		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	client := HttpClient()
	// logx.Printf("doing req %v", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response %s for %v", resp.Status, req.URL.String())
		return
	}

	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return

}

func Upload(url string, argKeys, argVals []string, file string) (respBytes []byte, err error) {

	if len(argKeys) != len(argVals) {
		logx.Fatalf("keys and vals must be equal size; not %v %v", len(argKeys), len(argVals))
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()
	ffw, err := w.CreateFormFile("upload", file) // form file writer
	if err != nil {
		return
	}
	if _, err = io.Copy(ffw, f); err != nil {
		return
	}

	// Add the other fields
	for i, key := range argKeys {
		if ffw, err = w.CreateFormField(key); err != nil {
			return
		}
		if _, err = ffw.Write([]byte(argVals[i])); err != nil {
			return
		}
	}

	// Without closing multipart writer.,
	// request will be missing the terminating boundary.
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Setting the content type; containing the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := HttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response %s for %v", resp.Status, req.URL.String())
		return
	}

	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}
