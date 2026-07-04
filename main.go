package main

import (
	_ "gf-eshop/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"gf-eshop/internal/cmd"
	_ "gf-eshop/internal/logic/brands"
	_ "gf-eshop/internal/logic/categories"
	_ "gf-eshop/internal/logic/category_brands"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
