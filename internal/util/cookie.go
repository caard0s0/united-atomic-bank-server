package util

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetCookie(ctx *gin.Context, accessToken string) {
	httpClientAddress := viper.Get("HTTP_CLIENT_ADDRESS")

	if httpClientAddress == "http://localhost:3000" {
		ctx.SetCookie("accessToken", accessToken, 60*30, "/", "localhost", false, true)
	}

	if httpClientAddress == "https://app.unitedatomicbank.com" {
		ctx.SetCookie("accessToken", accessToken, 60*30, "/", ".unitedatomicbank.com", true, true)
	}
}
