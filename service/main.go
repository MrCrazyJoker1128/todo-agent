package main

import (
	"context"
	"log"

	_ "service/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"service/internal/cmd"
	"service/internal/utility"
)

func main() {
	ctx := gctx.GetInitCtx()

	if err := utility.ESPing(); err != nil {
		log.Printf("elasticsearch connect failed: %v", err)
	} else {
		log.Println("elasticsearch connected successfully")
	}

	if err := initZhparser(ctx); err != nil {
		log.Printf("zhparser init failed: %v", err)
	} else {
		log.Println("zhparser initialized successfully")
	}

	cmd.Main.Run(ctx)
}

func initZhparser(ctx context.Context) error {
	_, err := g.DB().Exec(ctx, "SET zhparser.multi_duality = 1")
	return err
}
