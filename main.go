package main

import (
	"flag"
	"unicorn.dev.web-scrap/WebApi"
)

func main() {
	startServer := flag.Bool("server", false, "Start backend server")
	useTls := flag.Bool("use-tls", false, "use TLS")
	tlsKeyFile := flag.String("tls-key", "", "TLS key file (only with -use-tls)")
	tlsCertFile := flag.String("tls-crt", "", "TLS cert file (only with -use-tls)")

	flag.Parse()

	config := WebApi.ServerConfiguration{
		UseTls:     *useTls,
		TlsCrtFile: *tlsCertFile,
		TlsKeyFile: *tlsKeyFile,
	}

	if *startServer {
		WebApi.StartServer(&config)
	}
}
