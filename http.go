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

func Request(method, url string, vals p_url.Values) (respBytes []byte, err error) {

	if !(method == "GET" || method == "POST") {
		logx.Fatalf("must be GET or POST; not %v", method)
	}

	var req *http.Request

	if method == "GET" {

		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		req.URL.RawQuery = vals.Encode()
		// logx.Printf("GET requesting %v", req.URL.String())

	} else if method == "POST" {

		// bytes.NewBufferString AND strings.NewReader seem to work equally well
		req, err = http.NewRequest("POST", url, bytes.NewBufferString(vals.Encode())) // <-- URL-encoded payload
		// req, err = http.NewRequest("POST", url, strings.NewReader(vals.Encode())) // <-- URL-encoded payload
		if err != nil {
			return
		}
		// strangely, the json *reponse* is empty, if we omit this:
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	client := HttpClient()
	// logx.Printf("doing req %v", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Even for bad response status: Try to get the response body
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("%v; response status %q ", err, resp.Status)
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("bad response %q ", resp.Status)
		return
	}

	return

}

func Upload(url string, vals p_url.Values, file string) (respBytes []byte, err error) {

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
	for key, _ := range vals {
		val := vals.Get(key)
		if ffw, err = w.CreateFormField(key); err != nil {
			return
		}
		if _, err = ffw.Write([]byte(val)); err != nil {
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

	// We could redundantly add values as GET-Params,
	// since server might not expect POST-multipart encoding.
	/*
		q := req.URL.Query()
		for i, key := range argKeys {
			q.Add(key, argVals[i])
		}
		req.URL.RawQuery = q.Encode()
	*/

	client := HttpClient()

	// logx.Printf("uploading to %v", url)
	// logx.Printf("uploading to %+v", req)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Even for bad response status: Try to get the response body
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("%v; response status %q ", err, resp.Status)
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTemporaryRedirect {
		err = fmt.Errorf("bad response %q ", resp.Status)
		return
	}

	return
}
