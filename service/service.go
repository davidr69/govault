package service

import (
	"context"
	"log/slog"
	"os"

	"govault.lavacro.net/models"

	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *vault.Client
}

func NewVaultClientFromEnv(req models.AppConfig) error {
	cfg := vault.DefaultConfig()
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	secret, err := client.Logical().Write(
		"auth/approle/login",
		map[string]interface{}{
			"role_id":   req.RoleId,
			"secret_id": req.RoleSecret,
		},
	)
	if err != nil {
		return err
	}

	slog.Info("Logged into vault")
	client.SetToken(secret.Auth.ClientToken)

	// create file that will contain secrets
	h, err := os.Create(req.Params.FilePath)
	if err != nil {
		slog.Error("Failed to open file path: ", "error", err)
		return err
	}

	// loop through requests
	for idx := range req.Params.VaultRequest {
		reqItem := req.Params.VaultRequest[idx]
		path := reqItem.Path
		mount := reqItem.Mount

		slog.Info("path", "path", path, "mount", mount)

		for itemIdx := range reqItem.Items {
			item := reqItem.Items[itemIdx]
			slog.Info("Retrieve ", item.Key, item.Label)
			cred, err := client.KVv2(mount).Get(context.Background(), path)
			if err != nil {
				slog.Error("Failed to get secret from vault: ", "error", err)
				continue
			}
			_, err = h.WriteString(item.Label + "=" + cred.Data[item.Key].(string) + "\n")
			if err != nil {
				slog.Error("Failed to write secret to file: ", "error", err)
				continue
			}
			slog.Info("Wrote to file ", item.Key, item.Label)
		}
	}

	err = h.Close()
	if err != nil {
		slog.Error("Failed to close file: ", "error", err)
	}

	return nil
}
