package main

import (
	_ "gf-eshop/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"gf-eshop/internal/cmd"
	_ "gf-eshop/internal/logic/brands"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
