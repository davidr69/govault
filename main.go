package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"govault.lavacro.net/models"
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

type vaultRequest struct {
	roleId     string
	roleSecret string
	params     *models.VaultRequest
}

func main() {
	fmt.Print(banner)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))

	var req vaultRequest
	var jsonReq string

	req.roleId = os.Getenv("VAULT_ROLE_ID")
	req.roleSecret = os.Getenv("VAULT_ROLE_SECRET")
	jsonReq = os.Getenv("VAULT_REQUEST")

	if req.roleId == "" || req.roleSecret == "" || jsonReq == "" {
		slog.Error("Missing required environment variables")
		return
	}

	err := json.Unmarshal([]byte(jsonReq), &models.VaultRequest{})
	if err != nil {
		slog.Error("Error unmarshalling JSON request", "error", err)
		return
	}
}
