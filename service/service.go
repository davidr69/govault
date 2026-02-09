package service

import (
	"context"
	"log/slog"
	"os"
	"time"

	"govault.lavacro.net/models"

	vault "github.com/hashicorp/vault/api"
)

func WritePropertiesFile(req models.AppConfig) error {
	cfg := vault.DefaultConfig()
	// Increase timeout to account for network latency during bootup
	cfg.Timeout = 120 * time.Second
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	// Add context with timeout to the login operation
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	secret, err := client.Logical().WriteWithContext(
		ctx,
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

	// create a file that will contain secrets
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

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			cred, err := client.KVv2(mount).Get(ctx, path)
			cancel()
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
