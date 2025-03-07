// Package wangyiyun 网易云音乐热评
package wangyiyun

import (
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/zbputils/control/order"
)

const (
	wangyiyunURL     = "https://api.gmit.vip/Api/HotComments?format=text"
	wangyiyunReferer = "https://api.gmit.vip/"
	ua               = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
)

func init() {
	control.Register("wangyiyun", order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help:             "wangyiyun \n- 来份网易云热评",
	}).OnFullMatch("来份网易云热评").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.ReqWith(wangyiyunURL, "GET", wangyiyunReferer, ua)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text(helper.BytesToString(data)))
		})
}
