package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/hashicorp/vault/api"
)

const banner = `
banner here
`

type Config struct {
	ServiceAcctUser   string
	ServiceAcctSecret string
}

func main() {
	fmt.Print(banner)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))

	var cfg Config

	flag.StringVar(&cfg.ServiceAcctUser, "service-acct-user", "", "Service account username")
	flag.StringVar(&cfg.ServiceAcctSecret, "service-acct-secret", "", "Service account secret")
	flag.Parse()
}
