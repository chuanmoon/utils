package main

import (
	"log"
	"time"

	"github.com/chuanmoon/utils/cybase"
	"github.com/chuanmoon/utils/cyrpc"
)

type Bar struct {
}

func (*Bar) Foo(args *cyrpc.CommonArgs, reply *cyrpc.CommonReply) error {
	log.Println(args.Lang)
	reply.ErrorMsg = "lang: " + args.Lang
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	appName := "bar_1.0"
	_, _, infoZapLogger, _ := cybase.Init(appName)

	tool := cyrpc.NewTool(appName, infoZapLogger)
	go func() {
		tool.StartServer(&Bar{})
	}()

	time.Sleep(2 * time.Second)

	client := tool.GetClient()
	{ // json
		var reply = map[string]any{}
		err := client.CallJson("bar_1.0", "Bar.Foo", cyrpc.CommonArgs{Lang: "US"}, &reply)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(reply["ErrorMsg"])
	}
	{ // msgpack
		var reply = map[string]any{}
		err := client.CallMsgpack("bar_1.0", "Bar.Foo", cyrpc.CommonArgs{Lang: "US"}, &reply)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(reply["ErrorMsg"])
	}
}
