package utils

import zero "github.com/wdvxdr1123/ZeroBot"

func GetImages(ctx *zero.Ctx) (res []string) {
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
