package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	p_url "net/url" // the package url
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
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

func Request(method, url string, vals p_url.Values, cookies []*http.Cookie) (respBytes []byte, err error) {

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

	for _, v := range cookies {
		req.AddCookie(v)
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
		err = fmt.Errorf("req creation failed: %v ", err)
		return
	}
	// Setting the content type; containing the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := HttpClient()

	// logx.Printf("uploading to %v", url)
	// logx.Printf("uploading to %+v", req)

	resp, err1 := client.Do(req)
	if err1 != nil {
		err1 = fmt.Errorf("execute req failed: %v ", err1)
		if resp == nil {
			return
		}
	}
	defer resp.Body.Close()

	// Even for bad response status: Try to get the response body
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil || err1 != nil {
		err = fmt.Errorf("response status %q;\nerr1: %v \nerr:  %v \n%s",
			resp.Status, err1, err, respBytes)
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

var staticExtensions = append(imgExtensions,
	".css",
	".js",
	".sqlite",
	".md", // github markdown
)

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

// StripPrefix is a 'debug version' of http.StripPrefix;
// it logs the modified rewritten request elements;
// it also demonstrated, how to nest a http.Handler.
func StripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		if len(p) < len(r.URL.Path) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(p_url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			s := fmt.Sprintf("StripPrefix().Path  %-24v  -  pfx %v   =>  %v\n", r.URL.Path, prefix, r2.URL.Path)
			log.Printf(s)
			h.ServeHTTP(w, r2)
		} else {
			s := fmt.Sprintf("StripPrefix().Path  %-24v  -  pfx %v   =>  %v\n", r.URL.Path, prefix, p)
			w.Write([]byte(s))
			http.NotFound(w, r)
		}
	})
}

//
//
// ParseAndSaveUploaded is a convenience func to process uploaded files.
// Saves any number of file uploads to a given directory dir.
// If given directory is empty string, then file uploads are stored in bytes buffer.
// It includes an upload size limit to prevent attacks.
// ParseMultipartForm should not be called before
//
// b   contains the file contents, keyed by filename - not by input name.
// ok  is false, if an error occurred or if no file uploads were found.
// err contains any possible error.
//
// For multiple file upload:
// ok and err only flag problems with the last processed (failing) upload.
// Previous uploads might have been successfully written to file.
func ParseAndSaveUploaded(w http.ResponseWriter, r *http.Request, dir string) (b map[string]bytes.Buffer, ok bool, err error) {

	if strings.ToUpper(r.Method) != "POST" {
		return b, false, nil
	}

	b = map[string]bytes.Buffer{} // init map

	r.Body = http.MaxBytesReader(w, r.Body, 200*(1<<20)) // Limit upload sizet to 200 Mb

	const _24K = (1 << 10) * 24 // Up to 24 Kilobytes into memory, rest to temporary disk files
	if err = r.ParseMultipartForm(_24K); err != nil {
		return b, false, errors.Wrap(err, "Parse MultipartForm failed")
	}

	// Find all file uploads
	formFileInputs := map[string]bool{}
	for key, fheaders := range r.MultipartForm.File {
		// For each uploaded file there might be multiple headers (?segments)
		// Practically: There is every only one
		for range fheaders {
			formFileInputs[key] = true
		}
	}
	if len(formFileInputs) < 1 {
		return b, false, nil
	}
	log.Printf("%v file inputs. Input names are %v", len(formFileInputs), formFileInputs)

	for fileInput, _ := range formFileInputs {
		infile, mpHdr, err := r.FormFile(fileInput)
		if err != nil {
			return b, false, errors.Wrap(err, fmt.Sprintf("Opening parsed file failed for %v", fileInput))
		}

		mpHdr.Filename = LowerCasedUnderscored(mpHdr.Filename)
		mpHdr.Filename = strings.Replace(mpHdr.Filename, "_", "-", -1)

		if len(dir) > 0 {
			// open destination
			var outfile *os.File
			fullPath := filepath.Join(dir, mpHdr.Filename)
			if outfile, err = os.Create(fullPath); err != nil {
				return b, false, errors.Wrap(err, fmt.Sprintf("File creation failed for %v (%v)", fullPath, fileInput))
			}

			defer outfile.Close()
			var written int64
			if written, err = io.Copy(outfile, infile); err != nil { // 32K buffer copy
				return b, false, errors.Wrap(err, fmt.Sprintf("Data writing failed for %v (%v)", fullPath, fileInput))
			}
			err = os.Chmod(fullPath, 0644)
			if err != nil {
				return b, false, errors.Wrap(err, fmt.Sprintf("Permission setting failed for %v (%v)", fullPath, fileInput))
			}
			b[mpHdr.Filename] = bytes.Buffer{} // indicate success - in case succinct uploads fail
			log.Printf("Saved uploaded file: %v. %v bytes", fullPath, written)
		} else {
			var written int64
			bx := bytes.Buffer{}
			b[mpHdr.Filename] = bx
			if written, err = io.Copy(&bx, infile); err != nil {
				return b, false, errors.Wrap(err, fmt.Sprintf("Failed writing to byteBuffer data for %v (%v)", mpHdr.Filename, fileInput))
			}
			log.Printf("Uploaded file %v read into buffer. %v bytes", mpHdr.Filename, written)
		}
	}

	return b, true, nil
}
