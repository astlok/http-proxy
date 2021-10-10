package proxy

import (
	"io"
	"net"
	"net/http"
	"proxy/saver"
)

type Proxy struct {
	Saver *saver.Saver
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		p.HttpHandler(w, r)
	} else {
		p.HttpsHandler(w, r)
	}
}

func (p *Proxy) HttpHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	r.RequestURI = ""

	p.Saver.SaveRequest(r)
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	p.Saver.SaveResponse(resp)

	w.WriteHeader(resp.StatusCode)
	copyHeader(w.Header(), resp.Header)

	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func (p *Proxy) HttpsHandler(w http.ResponseWriter, r *http.Request) {
	p.Saver.SaveRequest(r)
	dest, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}
	client, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest, client)
	go transfer(client, dest)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
