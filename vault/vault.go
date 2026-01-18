package vault

import (
	"log/slog"

	"github.com/hashicorp/vault/api"
	"govault.lavacro.net/models"
)

type VaultClient struct {
	client *api.Client
}

func NewVaultClientFromEnv(req models.AppConfig) (*VaultClient, error) {
	slog.Info("foo")
	cfg := api.DefaultConfig()
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	secret, err := client.Logical().Write(
		"auth/approle/login",
		map[string]interface{}{
			"role_id":   req.RoleId,
			"secret_id": req.RoleSecret,
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
