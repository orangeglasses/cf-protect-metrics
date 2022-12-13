package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/kelseyhightower/envconfig"
)

const (
	CF_FORWARDED_URL          = "X-Cf-Forwarded-Url"
	CF_PROXY_SIGNATURE_HEADER = "X-Cf-Proxy-Signature"
)

type config struct {
	MetricsUser          string   `envconfig:"metricsuser" required:"true"`
	MetricsPassword      string   `envconfig:"metricspassword" required:"true"`
	MetricsEndpoint      string   `envconfig:"metricsendpoind" default:"/metrics"`
	expectedUsernameHash [32]byte `ignored:"true"`
	expectedPasswordHash [32]byte `ignored:"true"`
}

func main() {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	cfg.expectedUsernameHash = sha256.Sum256([]byte(cfg.MetricsUser))
	cfg.expectedPasswordHash = sha256.Sum256([]byte(cfg.MetricsPassword))

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			forwardedURL := req.Header.Get(CF_FORWARDED_URL)
			url, _ := url.Parse(forwardedURL)
			req.URL = url
			req.Host = url.Host
		},
	}

	http.HandleFunc("/", ProxyRequestHandler(proxy, cfg))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy, cfg config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		forwardedURL := r.Header.Get(CF_FORWARDED_URL)
		fURL, _ := url.Parse(forwardedURL)

		if sig := r.Header.Get(CF_PROXY_SIGNATURE_HEADER); sig == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if fURL.Path == cfg.MetricsEndpoint {
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], cfg.expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], cfg.expectedPasswordHash[:]) == 1)

			if !usernameMatch || !passwordMatch {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		proxy.ServeHTTP(w, r)
	}
}
