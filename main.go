package main

import (
	"flag"
	"unicorn.dev.web-scrap/WebApi"
)

func main() {
	useTls := flag.Bool("use-tls", false, "use TLS")
	tlsKeyFile := flag.String("tls-key", "", "TLS key file (only with -use-tls)")
	tlsCertFile := flag.String("tls-crt", "", "TLS cert file (only with -use-tls)")
	apiKeysFile := flag.String("api-keys", "", "Add API keys file")

	flag.Parse()

	config := WebApi.ServerConfiguration{
		UseTls:      *useTls,
		TlsCrtFile:  *tlsCertFile,
		TlsKeyFile:  *tlsKeyFile,
		ApiKeysFile: *apiKeysFile,
	}

	WebApi.StartServer(&config)
}
