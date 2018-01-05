package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func main() {
	var (
		listen, cert, key, where string
		useLogging               bool
		flushInterval            time.Duration
	)

	flag.StringVar(&listen, "listen", ":8443", "bind address to listen on")
	flag.StringVar(&key, "key", "key.pem", "path to PEM key")
	flag.StringVar(&cert, "cert", "cert.pem", "path to PEM certificate")
	flag.StringVar(&where, "where", "127.0.0.1:8080", "address to forward connections to")
	flag.BoolVar(&useLogging, "logging", true, "log requests")
	flag.DurationVar(&flushInterval, "flush-interval", 0, "minimum duration between flushes to the client (default: off)")
	flag.Parse()

	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = where
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Header.Set("X-Forwarded-Proto", "https")
		host := strings.SplitN(req.Host, ":", 2)
		if len(host) > 1 {
			req.Header.Set("X-Forwarded-Port", host[1])
		}
	}

	modifyResponse := func(resp *http.Response) error {
		req := resp.Request
		if useLogging {
			log.Printf("%21v %7v %-40v %v %19v %v",
				req.RemoteAddr,
				req.Method,
				"//"+req.Host+"/"+
					strings.SplitN(req.URL.String(), "/", 4)[3],
				resp.StatusCode,
				resp.ContentLength,
				req.Header.Get("User-Agent"))
		}
		if location := resp.Header.Get("Location"); strings.HasPrefix(location, "http://"+req.Host) {
			// rewrite location header
			resp.Header.Set("Location", "https"+location[4:])
		}
		return nil
	}

	proxy := &httputil.ReverseProxy{Director: director, ModifyResponse: modifyResponse, FlushInterval: flushInterval}

	server := &http.Server{Addr: listen, Handler: proxy}
	log.Fatalln(server.ListenAndServeTLS(cert, key))
}
