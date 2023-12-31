package wclient

import (
	"net"

	"github.com/opentdp/go-helper/logman"
	"github.com/opentdp/go-helper/strutil"

	"github.com/opentdp/wechat-rest/args"
	"github.com/opentdp/wechat-rest/wcferry"
)

var wc *wcferry.Client

func Connect() *wcferry.Client {

	if wc != nil {
		return wc
	}

	host, port, err := net.SplitHostPort(args.Wcf.Address)
	if err != nil {
		logman.Fatal("invalid address", "error", err)
	}

	wc = &wcferry.Client{
		ListenAddr: host,
		ListenPort: strutil.ToInt(port),
		SdkLibrary: args.Wcf.SdkLibrary,
		WeChatAuto: args.Wcf.WeChatAuto,
	}

	logman.Info("wcf starting ...")
	if err := wc.Connect(); err != nil {
		logman.Fatal("failed to start wcf", "error", err)
	}

	// 打印收到的消息
	if args.Wcf.MsgPrint {
		wc.EnrollReceiver(true, wcferry.WxMsgPrinter)
	}

	return wc

}
