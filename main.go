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
	_ "gf-eshop/internal/logic/products"
	_ "gf-eshop/internal/logic/attributes"
	_ "gf-eshop/internal/logic/attribute_values"
	_ "gf-eshop/internal/logic/category_attributes"
	_ "gf-eshop/internal/logic/inventories"
	_ "gf-eshop/internal/logic/inventory_logs"
	_ "gf-eshop/internal/logic/product_attributes"
	_ "gf-eshop/internal/logic/product_descriptions"
	_ "gf-eshop/internal/logic/product_versions"
	_ "gf-eshop/internal/logic/skus"
	_ "gf-eshop/internal/logic/notification"
	_ "gf-eshop/internal/logic/permissions"
	_ "gf-eshop/internal/logic/roles"
	_ "gf-eshop/internal/logic/staff"
	_ "gf-eshop/internal/logic/orders"
	_ "gf-eshop/internal/logic/payments"
	_ "gf-eshop/internal/logic/carts"
	_ "gf-eshop/internal/logic/user"
	_ "gf-eshop/internal/logic/user_auth"
	_ "gf-eshop/internal/logic/address"
	_ "gf-eshop/internal/logic/marketing"
	_ "gf-eshop/internal/logic/user_levels"
	_ "gf-eshop/internal/logic/user_points"
	_ "gf-eshop/internal/logic/points_rules"
	_ "gf-eshop/internal/logic/level_rules"
	_ "gf-eshop/internal/logic/merchants"
	_ "gf-eshop/internal/logic/merchant_bank_accounts"
	_ "gf-eshop/internal/logic/merchant_contacts"
	_ "gf-eshop/internal/logic/merchant_qualifications"
	_ "gf-eshop/internal/logic/merchant_withdrawals"
	_ "gf-eshop/internal/logic/merchant_balances"
	_ "gf-eshop/internal/logic/merchant_roles"
	_ "gf-eshop/internal/logic/merchant_role_permissions"
	_ "gf-eshop/internal/logic/operation_logs"
	_ "gf-eshop/internal/logic/departments"
	_ "gf-eshop/internal/logic/reviews"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
