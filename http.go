package daas

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

type HTTPRequest struct {
	client *http.Client
	header map[string]string
}

type HTTPTransport struct {
	AppName            string
	AppVersion         string
	DisableCompression bool
	DisableKeepAlive   bool
	SkipVerify         bool
	Timeout            int
	AllowRedirects     bool
	ClientCert         string
	ClientKey          string
	CACert             string
	UseHTTP2           bool
}

func CreateClient(config *HTTPTransport) (*HTTPRequest, error) {
	if config.Timeout == 0 {
		config.Timeout = 10000
	}
	result := &HTTPRequest{
		header: make(map[string]string),
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Dial:                  (&net.Dialer{Timeout: 5 * time.Second}).Dial,
				DisableCompression:    config.DisableCompression,
				DisableKeepAlives:     config.DisableKeepAlive,
				ResponseHeaderTimeout: time.Millisecond * time.Duration(config.Timeout),
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: config.SkipVerify},
			},
		},
	}

	if config.AppName != "" {
		result.header["user-agent"] = fmt.Sprintf("%s/%s (k8s;daas-projects;)", config.AppName, config.AppVersion)
	}

	// if !allowRedirects {
	//returning an error when trying to redirect. This prevents the redirection from happening.
	// client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
	// 	return util.NewRedirectError("redirection not allowed")
	/// }
	// }

	if config.ClientCert == "" && config.ClientKey == "" && config.CACert == "" {
		return result, nil
	}

	if config.ClientCert == "" {
		return nil, fmt.Errorf("client certificate can't be empty")
	}

	if config.ClientKey == "" {
		return nil, fmt.Errorf("client key can't be empty")
	}
	cert, err := tls.LoadX509KeyPair(config.ClientCert, config.ClientKey)
	if err != nil {
		return nil, fmt.Errorf("unable to load cert tried to load %v and %v but got %v", config.ClientCert, config.ClientKey, err)
	}

	// Load our CA certificate
	clientCACert, err := ioutil.ReadFile(config.CACert)
	if err != nil {
		return nil, fmt.Errorf("unable to open cert %v", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: config.SkipVerify,
	}

	// tlsConfig.BuildNameToCertificate()
	t := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	if config.UseHTTP2 {
		http2.ConfigureTransport(t)
	}
	result.client.Transport = t
	return result, nil
}

func (e *HTTPRequest) HeaderAdd(name string, value string) {
	e.header[name] = value
}

func (e *HTTPRequest) HeaderClear() {
	e.header = make(map[string]string)
}

func (e *HTTPRequest) Send(reqMethod string, reqUrl string, reqBuffer []byte) ([]byte, int, time.Duration, error) {
	// params := url.Values{}
	// params.Add("message", "this will be esc@ped!")
	// params.Add("author", "golang c@fe >.<")
	// fmt.Println("http://example.com/say?" + params.Encode())

	// loadUrl = escapeUrlStr(loadUrl)

	var bodyBuffer io.Reader
	if len(reqBuffer) > 0 {
		bodyBuffer = bytes.NewBuffer(reqBuffer)
	}

	req, err := http.NewRequest(reqMethod, reqUrl, bodyBuffer)
	if err != nil {
		// fmt.Println("An error occured doing request", err)
		return nil, -1, -1, err
	}

	for hk, hv := range e.header {
		req.Header.Add(hk, hv)
	}

	// if host != "" {
	// 	req.Host = host
	// }
	start := time.Now()
	resp, err := e.client.Do(req)
	if err != nil {
		rr, ok := err.(*url.Error)
		if !ok {
			return nil, -1, -1, rr
		} else {
			return nil, -1, -1, err
		}
	}

	if resp == nil {
		return nil, -1, -1, fmt.Errorf("empty response")
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, -1, err
	}

	sizeHeader := int(estimateHttpHeadersSize(resp.Header))
	duration := time.Since(start)
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return body, len(body) + sizeHeader, duration, nil
	} else if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusTemporaryRedirect {
		return body, int(resp.ContentLength) + sizeHeader, duration, nil
	} else {
		return body, len(body) + sizeHeader, duration, fmt.Errorf("received status code %d", resp.StatusCode)
	}
}

func estimateHttpHeadersSize(headers http.Header) (result int64) {
	result = 0

	for k, v := range headers {
		result += int64(len(k) + len(": \r\n"))
		for _, s := range v {
			result += int64(len(s))
		}
	}

	result += int64(len("\r\n"))

	return result
}
