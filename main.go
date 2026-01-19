package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"govault.lavacro.net/models"
	"govault.lavacro.net/service"
)

const banner = `
                                                       {__  {__  
                                                       {__  {__  
   {__      {__          {__     {__   {__    {__  {__ {__{_{_ {_
 {__  {__ {__  {__ {_____ {__   {__  {__  {__ {__  {__ {__  {__  
{__   {__{__    {__        {__ {__  {__   {__ {__  {__ {__  {__  
 {__  {__ {__  {__          {_{__   {__   {__ {__  {__ {__  {__  
     {__    {__              {__      {__ {___  {__{__{___   {__ 
  {__                                                            
`

func main() {
	fmt.Print(banner)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))

	var req models.AppConfig
	var jsonReq string

	req.RoleId = os.Getenv("VAULT_ROLE_ID")
	req.RoleSecret = os.Getenv("VAULT_ROLE_SECRET")
	jsonReq = os.Getenv("VAULT_REQUEST")

	fmt.Println(jsonReq)

	if req.RoleId == "" || req.RoleSecret == "" || jsonReq == "" {
		slog.Error("Missing required environment variables")
		return
	}

	err := json.Unmarshal([]byte(jsonReq), &req.Params)
	if err != nil {
		slog.Error("Error unmarshalling JSON request", "error", err)
		return
	}

	slog.Info("Request", "request", req)
	err = service.NewVaultClientFromEnv(req)
	if err != nil {
		slog.Error("Error creating Vault client", "error", err)
	}
}
