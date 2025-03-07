// Package nativesetu 本地setu
package nativesetu

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"

	"github.com/FloatTech/zbputils/control/order"
)

var (
	setupath = "/tmp" // 绝对路径，图片根目录
)

func init() {
	engine := control.Register("nativesetu", order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help: "本地涩图\n" +
			"- 本地[xxx]\n" +
			"- 刷新本地[xxx]\n" +
			"- 设置本地setu绝对路径[xxx]\n" +
			"- 刷新所有本地setu\n" +
			"- 所有本地setu分类",
		PrivateDataFolder: "nsetu",
	})

	ns.db.DBPath = engine.DataFolder() + "data.db"
	cfgfile := engine.DataFolder() + "setupath.txt"
	if file.IsExist(cfgfile) {
		b, err := os.ReadFile(cfgfile)
		if err == nil {
			setupath = helper.BytesToString(b)
			logrus.Println("[nsetu] set setu dir to", setupath)
		}
	}

	engine.OnRegex(`^本地(.*)$`, ctxext.FirstValueInList(ns)).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			imgtype := ctx.State["regex_matched"].([]string)[1]
			sc := new(setuclass)
			ns.mu.RLock()
			err := ns.db.Pick(imgtype, sc)
			ns.mu.RUnlock()
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
			} else {
				p := "file:///" + setupath + "/" + sc.Path
				ctx.SendChain(message.Text(imgtype, ": ", sc.Name, "\n"), message.Image(p))
			}
		})
	engine.OnRegex(`^刷新本地(.*)$`, ctxext.FirstValueInList(ns), zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			imgtype := ctx.State["regex_matched"].([]string)[1]
			err := ns.scanclass(os.DirFS(setupath), imgtype, imgtype)
			if err == nil {
				ctx.SendChain(message.Text("成功！"))
			} else {
				ctx.SendChain(message.Text("ERROR: ", err))
			}
		})
	engine.OnRegex(`^设置本地setu绝对路径(.*)$`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			setupath = ctx.State["regex_matched"].([]string)[1]
			err := os.WriteFile(cfgfile, helper.StringToBytes(setupath), 0644)
			if err == nil {
				ctx.SendChain(message.Text("成功！"))
			} else {
				ctx.SendChain(message.Text("ERROR: ", err))
			}
		})
	engine.OnFullMatch("刷新所有本地setu", zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			err := ns.scanall(setupath)
			if err == nil {
				ctx.SendChain(message.Text("成功！"))
			} else {
				ctx.SendChain(message.Text("ERROR: ", err))
			}
		})
	engine.OnFullMatch("所有本地setu分类").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			msg := "所有本地setu分类"
			ns.mu.RLock()
			for i, c := range ns.List() {
				n, err := ns.db.Count(c)
				if err == nil {
					msg += fmt.Sprintf("\n%02d. %s(%d)", i, c, n)
				} else {
					msg += fmt.Sprintf("\n%02d. %s(error)", i, c)
					logrus.Errorln("[nsetu]", err)
				}
			}
			ns.mu.RUnlock()
			ctx.SendChain(message.Text(msg))
		})
}
