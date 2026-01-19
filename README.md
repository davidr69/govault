# govault

This is a Go port of the Java/Spring Boot
[vault-retreive](https://github.com/davidr69/vault-retrieve)
application.


https://pkg.go.dev/github.com/hashicorp/vault/api


set VAULT_ADDR


VAULT_ROLE_ID

VAULT_ROLE_SECRET

VAULT_REQUEST

```json
{
  "vault_request": [{
    "mount": "lavacro",
    "path": "prod/database/postgresql",
    "items": [{
      "key": "password",
      "label": "spring.datasource.password"
    }]
  },{
    "mount": "lavacro",
    "path": "prod/gcl",
    "items": [{
      "key": "authentication",
      "label": "google.cloud.logging"
    }]
  },{
    "mount": "lavacro",
    "path": "prod/github",
    "items": [{
      "key": "token",
      "label": "github.token"
    }]
  }],
  "file_path": "/var/tmp/vault/vault.properties"
}
```

![go-vault](images/go-vault.png)
