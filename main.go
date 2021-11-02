package main

import (
	"fmt"

	"github.com/AkiraXie/hoshino.go/hoshino"
	_ "github.com/AkiraXie/hoshino.go/module"
	"github.com/AkiraXie/hoshino.go/utils"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	cfg := hoshino.HsoConfig
	zero.Run(zero.Config{
		NickName:      cfg.NickName,
		CommandPrefix: cfg.CommandPrefix,
		SuperUsers:    cfg.SuperUsers,
		Driver: []zero.Driver{
			driver.NewWebSocketClient(fmt.Sprintf("ws://%s:%d", cfg.Host, cfg.Port), cfg.AccessToken),
		},
	})
	<-utils.SetupMainSignalHandler()
}
