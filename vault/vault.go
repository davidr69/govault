package vault

import (
	"log/slog"
	"os"

	"github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *api.Client
}

func NewVaultClientFromEnv() (*VaultClient, error) {
	slog.Info("foo")
	cfg := api.DefaultConfig()
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	roleID := os.Getenv("VAULT_ROLE_ID")
	secretID := os.Getenv("VAULT_SECRET_ID")

	secret, err := client.Logical().Write(
		"auth/approle/login",
		map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		},
	)
	if err != nil {
		return nil, err
	}

	client.SetToken(secret.Auth.ClientToken)

	// Start renewal automatically
	if secret.Auth.Renewable {
		renewer, err := client.NewRenewer(&api.RenewerInput{
			Secret: secret,
		})
		if err != nil {
			return nil, err
		}
		go renewer.Renew()
	}

	return &VaultClient{client: client}, nil
}
