package utils

import zero "github.com/wdvxdr1123/ZeroBot"

func GetImages(ctx *zero.Ctx) (res []string) {
	ctx.Send(ctx.Event.RawMessage)
	ctx.Send(ctx.Event.Message)
	for _, v := range ctx.Event.Message {
		if v.Type == "image" {
			if data, ok := v.Data["file"]; ok {
				res = append(res, data)
			}
		} else {
			continue
		}
	}
	return
}
