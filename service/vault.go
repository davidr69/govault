package service

import (
	"context"
	"log/slog"
	"os"

	"govault.lavacro.net/models"

	vault "github.com/hashicorp/vault/api"
	//	auth "github.com/hashicorp/vault/api/auth/approle"
)

type VaultClient struct {
	client *vault.Client
}

func NewVaultClientFromEnv(req models.AppConfig) (*VaultClient, error) {
	cfg := vault.DefaultConfig()
	client, err := vault.NewClient(cfg)
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

	slog.Info("Logged into vault")
	client.SetToken(secret.Auth.ClientToken)

	// create file that will contain secrets
	h, err := os.Create(req.Params.MountPath)
	if err != nil {
		slog.Error("Failed to open mount path: ", err)
		return nil, err
	}

	// loop through requests
	for idx := range req.Params.VaultRequest {
		reqItem := req.Params.VaultRequest[idx]
		mountPath := reqItem.Path

		slog.Info(mountPath)
		slog.Info("Mount path: ", mountPath)

		for itemIdx := range reqItem.Items {
			item := reqItem.Items[itemIdx]
			slog.Info("Retrieve ", item.Key, " as ", item.Label)
			cred, err := client.KVv2(mountPath).Get(context.Background(), item.Key)
			if err != nil {
				slog.Error("Failed to get secret from vault: ", err)
				continue
			}
			slog.Info("Got secret from vault: ", cred)
		}
	}

	h.Write([]byte("foobar"))
	h.Close()

	return &VaultClient{client: client}, nil
}
