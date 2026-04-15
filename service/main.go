package main

import (
	"log"

	_ "service/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"service/internal/cmd"
	"service/internal/utility"
)

func main() {
	if err := utility.ESPing(); err != nil {
		log.Printf("elasticsearch connect failed: %v", err)
	} else {
		log.Println("elasticsearch connected successfully")
	}

	cmd.Main.Run(gctx.GetInitCtx())
}
