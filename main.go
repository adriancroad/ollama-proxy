package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

const maxBodyLog = 10 * 1024 // 10KB

func main() {
	proxyPort := os.Getenv("PROXY_PORT")
	if proxyPort == "" {
		proxyPort = ":8080"
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	target, err := url.Parse(ollamaURL)
	if err != nil {
		log.Fatalf("Invalid OLLAMA_URL %q: %v", ollamaURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.FlushInterval = -1

	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("RESPONSE: %s %s -> %d %s",
			resp.Request.Method, resp.Request.URL.Path,
			resp.StatusCode, http.StatusText(resp.StatusCode))
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("PROXY ERROR: %s %s -> %v", r.Method, r.URL.Path, err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Ollama proxy listening on %s, forwarding to %s", proxyPort, target.String())
	if err := http.ListenAndServe(proxyPort, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func logRequest(r *http.Request) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "REQUEST: %s %s\n", r.Method, r.URL.Path)

	fmt.Fprintf(&sb, "  Headers:\n")
	for name, values := range r.Header {
		fmt.Fprintf(&sb, "    %s: %s\n", name, strings.Join(values, ", "))
	}

	if r.Body != nil && r.Body != http.NoBody {
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			fmt.Fprintf(&sb, "  Body: <error reading: %v>\n", err)
		} else {
			r.Body = io.NopCloser(bytes.NewReader(body))
			if len(body) > 0 {
				display := body
				truncated := false
				if len(body) > maxBodyLog {
					display = body[:maxBodyLog]
					truncated = true
				}
				fmt.Fprintf(&sb, "  Body (%d bytes):\n    %s", len(body), string(display))
				if truncated {
					fmt.Fprintf(&sb, "\n    ... truncated (%d bytes total)", len(body))
				}
				fmt.Fprintln(&sb)
			}
		}
	}

	log.Print(sb.String())
}
