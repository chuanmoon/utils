package main

import (
	"log"
	"time"

	"github.com/chuanmoon/utils/mlog"
	"github.com/chuanmoon/utils/rpc"
)

type Bar struct {
}

func (*Bar) Foo(args *rpc.CommonArgs, reply *rpc.CommonReply) error {
	log.Println(args.Lang)
	reply.ErrorMsg = "lang: " + args.Lang
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	appName := "bar_1.0"
	infoZapLogger, _ := mlog.EcszapInfoAndErrorLogger(appName)
	tool := rpc.NewTool(appName, infoZapLogger)
	go func() {
		tool.StartServer(&Bar{})
	}()

	time.Sleep(2 * time.Second)

	client := tool.GetClient()
	{ // json
		var reply = map[string]any{}
		err := client.CallJson("bar_1.0", "Bar.Foo", rpc.CommonArgs{Lang: "US"}, &reply)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(reply["ErrorMsg"])
	}
	{ // msgpack
		var reply = map[string]any{}
		err := client.CallMsgpack("bar_1.0", "Bar.Foo", rpc.CommonArgs{Lang: "US"}, &reply)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(reply["ErrorMsg"])
	}
}
