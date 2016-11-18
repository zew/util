package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	p_url "net/url" // the package url
	"os"
	"path"
	"strings"
	"time"

	"github.com/zew/logx"
)

// Get a http server
// thenewstack.io/building-a-web-server-in-go/
// blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
func HttpServer() *http.Server {
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv
}

// Get a http client
// blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
func HttpClient() *http.Client {

	// Giving more granular control than http.Client.Timeout:
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 120, // covers entire exchange, from Dial (unless reused) to reading the body.
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

var imgExtensions = []string{
	".ico",
	".png",
	".jpg",
	".gif",
	".svg",
}

var staticExtensions = append(imgExtensions, ".css", ".js")

func StaticExtension(r *http.Request) bool {
	for _, v := range staticExtensions {
		if strings.HasSuffix(r.URL.Path, v) {
			return true
		}
	}
	return false
}

func ImageExtension(p string) bool {
	for _, v := range imgExtensions {
		if strings.HasSuffix(p, v) {
			return true
		}
	}
	return false
}

func UrlParseImproved(str string) (*p_url.URL, error) {
	// Prevent "google.com" => u.Host == "" when scheme == ""
	// p_url.Parse acts ugly, if there is not http:// prefix
	if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
		str = "https://" + str
	}
	ourl, err := p_url.Parse(str)
	return ourl, err
}

// Stripping subdomains
// Last two "." delimited tokens of hostname
// xx1.shop.wsj.com => wsj.com
func HostCore(h string) (core string, subdomains []string) {

	host, _, err := net.SplitHostPort(h) // the damn func *absorbs* everything on error

	if err != nil {
		// logx.Println(err, "host:", h, "port:", port)
		// logx.PrintStackTrace()
	} else {
		h = host
	}

	if strings.Count(h, ".") < 2 {
		core = h
		subdomains = []string{}
		return
	}
	parts := strings.Split(h, ".")
	lenP := len(parts)
	core = parts[lenP-2] + "." + parts[lenP-1]
	subdomains = parts[0 : lenP-2]

	if len(subdomains) > 0 && subdomains[len(subdomains)-1] == "www" {
		subdomains = subdomains[:len(subdomains)-1] // chop off last
	}

	return
}

// subdomain1.host.com/dir1
// is reformed to
// host.com/subdomain1/dir1
// except for www.host.com
func NormalizeSubdomainsToPath(url *p_url.URL) string {
	subdomains := []string{}
	url.Host, subdomains = HostCore(url.Host)
	if len(subdomains) > 0 {
		url.Path = path.Join(path.Join(subdomains...), url.Path)
	}
	url.Scheme = ""
	url.User = nil

	str := url.String()
	if strings.HasPrefix(str, "//") {
		str = str[2:]
	}
	return str
}
