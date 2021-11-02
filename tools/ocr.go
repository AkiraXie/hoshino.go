package tools

import (
	"fmt"
	"strings"

	"github.com/AkiraXie/hoshino.go/hoshino"
	"github.com/AkiraXie/hoshino.go/utils"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	sv := hoshino.NewService("ocr", true)
	oocr := sv.OnCommandGroup([]string{"ocr"}).SetPriority(1)
	oocr.SetArgParser(func(ctx *zero.Ctx) {
		imgs := utils.GetImages(ctx)
		ctx.State["images"] = imgs
	})
	oocr.Got("images", "请发送图片", func(ctx *zero.Ctx) {
		imgs := ctx.State["images"].([]string)
		for i, img := range imgs {
			msg := fmt.Sprintf("第%d张图片的识别结果是:", i+1)
			res := ctx.OCRImage(img)
			texts := res.Get("texts").Array()
			txts := make([]string, 0)
			for _, text := range texts {
				txts = append(txts, text.Get("text").Str)
			}
			msg += strings.Join(txts, " | ")
			oocr.Finish(ctx, msg)
		}
	}).Ok()
}
