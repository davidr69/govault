package models

type Item struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type VaultRequest struct {
	Path  string `json:"path"`
	Items []Item `json:"items"`
}

type VaultConfig struct {
	VaultRequest []VaultRequest `json:"vault_request"`
	MountPath    string         `json:"mount_path"`
}

/*
   {
     "vault_request": [{
       "path": "lavacro/data/prod/database/postgresql",
       "items": [{
         "key": "password",
         "label": "spring.datasource.password"
       }]
     },{
       "path": "lavacro/data/prod/gcl",
       "items": [{
         "key": "authentication",
         "label": "google.cloud.logging"
       }]
     },{
       "path": "lavacro/data/prod/github",
       "items": [{
         "key": "token",
         "label": "github.token"
       }]
     }],
     "mount_path": "/var/tmp/vault/vault.properties"
   }
*/
