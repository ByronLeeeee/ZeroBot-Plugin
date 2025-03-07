// Package shindan 基于 https://shindanmaker.com 的测定小功能
package shindan

import (
	"github.com/FloatTech/AnimeAPI/shindanmaker"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/zbputils/control/order"

	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
)

func init() {
	engine := control.Register("shindan", order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help: "shindan\n" +
			"- 今天是什么少女[@xxx]\n" +
			"- 异世界转生[@xxx]\n" +
			"- 卖萌[@xxx]\n" +
			"- 抽老婆[@xxx]",
	})
	engine.OnPrefix("异世界转生", number(587874)).SetBlock(true).Limit(ctxext.LimitByUser).Handle(handle)
	engine.OnPrefix("今天是什么少女", number(162207)).SetBlock(true).Limit(ctxext.LimitByUser).Handle(handle)
	engine.OnPrefix("卖萌", number(360578)).SetBlock(true).Limit(ctxext.LimitByUser).Handle(handle)
	engine.OnPrefix("抽老婆", number(1075116)).SetBlock(true).Limit(ctxext.LimitByUser).Handle(handle)
}

// shindanmaker 处理函数
func handle(ctx *zero.Ctx) {
	// 获取名字
	name := ctxext.NickName(ctx)
	// 调用接口
	txt, err := shindanmaker.Shindanmaker(ctx.State["id"].(int64), name)
	if err != nil {
		ctx.SendChain(message.Text("ERROR: ", err))
	}
	// TODO: 可注入
	switch ctx.State["id"].(int64) {
	case 587874, 162207:
		data, err := text.RenderToBase64(txt, text.FontFile, 400, 20)
		if err != nil {
			log.Errorln("[shindan]:", err)
		}
		if id := ctx.SendChain(message.Image("base64://" + helper.BytesToString(data))); id.ID() == 0 {
			ctx.SendChain(message.Text("ERROR: 可能被风控了"))
		}
	default:
		ctx.Send(txt)
	}
}

// 传入 shindanmaker id
func number(id int64) func(ctx *zero.Ctx) bool {
	return func(ctx *zero.Ctx) bool {
		ctx.State["id"] = id
		return true
	}
}
