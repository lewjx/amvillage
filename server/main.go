package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

type Args struct {
	// ListenAddr is the address to listen to.
	ListenAddr string
	// DebugForward is the address to forward to.
	DebugForward string
	// Serve is the directory to serve files from.
	Serve string
	// Config is the location of the config file.
	Config string
}

func main() {
	var args Args
	flag.StringVar(&args.ListenAddr, "addr", ":8000", "Address to listen to, default :8000")
	flag.StringVar(&args.DebugForward, "forward", "", "Address to forward to")
	flag.StringVar(&args.Serve, "serve", "", "Directory to serve frontend from")
	flag.StringVar(&args.Config, "config", "config.json", "Config file to open")
	flag.Parse()
	if err := run(args); err != nil {
		log.Fatal(err.Error())
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// run starts the web server.
func run(args Args) error {
	cfg, err := readConfigFile(args.Config)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	state, err := NewGameState(*cfg)
	if err != nil {
		return fmt.Errorf("error creating game state: %w", err)
	}
	go state.Start()
	http.HandleFunc("/api/config", handleConfig(*cfg))
	http.HandleFunc("/api/ws", state.handleWebsocket)
	switch {
	case args.DebugForward != "" && args.Serve != "":
		return fmt.Errorf("only one of -forward and -serve can be specified at once")
	case args.DebugForward != "":
		log.Printf("forwarding / to %s\n", args.DebugForward)
		http.HandleFunc("/", forward(args.DebugForward))
	case args.Serve != "":
		log.Printf("serving %s on /\n", args.Serve)
		http.Handle("/", http.FileServer(http.Dir(args.Serve)))
	default:
		log.Printf("root path is not handled currently\n")
	}
	log.Printf("listening at %s\n", args.ListenAddr)
	if err := http.ListenAndServe(args.ListenAddr, nil); err != nil {
		return fmt.Errorf("error serving HTTP: %w", err)
	}
	return nil
}

func handleConfig(cfg Config) func(http.ResponseWriter, *http.Request) {
	cfg = stripConfigSecret(cfg)
	bytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("error marshalling config: %v", err)
	}
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(bytes); err != nil {
			log.Printf("error writing config to ResponseWriter: %v", err)
		}
	}
}

func forward(target string) func(http.ResponseWriter, *http.Request) {
	uri, err := url.ParseRequestURI(target)
	if err != nil {
		log.Fatal("Error parsing forward URL: " + err.Error())
	}
	if uri.Scheme != "http" && uri.Scheme != "https" {
		log.Println("Scheme seems to not be specified. Forwarding may not work properly.")
	}
	proxy := httputil.NewSingleHostReverseProxy(uri)
	proxy.Director = nil
	proxy.Rewrite = func(req *httputil.ProxyRequest) {
		req.SetURL(uri)
		req.Out.Host = uri.Host
	}
	return proxy.ServeHTTP
}

func readConfigFile(location string) (*Config, error) {
	f, err := os.ReadFile(location)
	if err != nil {
		return nil, fmt.Errorf("error reading file %q: %w", location, err)
	}
	var cfg Config
	if err := json.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return &cfg, nil
}
