package tools

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/AkiraXie/hoshino.go/hoshino"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	sv := hoshino.NewService("nbnhhsh", true)
	nbn := sv.OnCommandGroup([]string{"nbnhhsh", "谜语"})
	nbn.Handle(func(ctx *zero.Ctx) {
		args := ctx.State["args"].(string)
		resp, err := http.PostForm("https://lab.magiconch.com/api/nbnhhsh/guess", url.Values{"text": []string{args}})
		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				_ = resp.Body.Close()

				json := gjson.ParseBytes(body)
				res := make([]string, 0)
				var jsonPath string
				if json.Get("0.trans").Exists() {
					jsonPath = "0.trans"
				} else {
					jsonPath = "0.inputting"
				}
				for _, value := range json.Get(jsonPath).Array() {
					res = append(res, value.String())
				}
				nbn.Finish(ctx, args+": "+strings.Join(res, ", "))
			}

		}
	}).Ok()
}
