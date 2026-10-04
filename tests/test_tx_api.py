#!/usr/bin/env python3
"""
gf-eshop 订单/支付/退款业务 API 测试脚本

测试覆盖：
  - 订单创建、列表、详情、状态更新
  - 支付创建、支付回调、支付查询
  - 退款创建

用法：
    python tests/test_tx_api.py
    python tests/test_tx_api.py --url http://localhost:8000
    python tests/test_tx_api.py --sku-id 1
"""
import argparse
import json
import sys
import urllib.request
import urllib.error


BASE_URL = "http://localhost:8000"
USERNAME = "admin"
PASSWORD = "123456"
SKU_ID = 1       # 测试用的 SKU ID（需确保该 SKU 存在且价格 > 0）

# 用户令牌：orders / payments 挂在 UserAuthMiddleware 下，**staff 令牌不认**（会返回 1004）。
# 所以本脚本需要两个身份：staff 只用于 1.1 登录用例，业务接口一律用下面的用户令牌。
USER_TOKEN = None
USER_NAME = "test_tx_api"
USER_PASSWORD = "test123456"


def req(method, url, headers=None, body=None, auth=True):
    if headers is None:
        headers = {}
    if auth and USER_TOKEN and "Authorization" not in headers:
        headers["Authorization"] = f"Bearer {USER_TOKEN}"
    data = json.dumps(body).encode() if body else None
    if body:
        headers.setdefault("Content-Type", "application/json")
    http_req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(http_req, timeout=10) as resp:
            body = resp.read().decode()
            try:
                return json.loads(body)
            except json.JSONDecodeError:
                return {"code": 0, "raw": body}
    except urllib.error.HTTPError as e:
        body = e.read().decode()
        try:
            return json.loads(body)
        except json.JSONDecodeError:
            return {"code": e.code, "message": body}
    except Exception as e:
        return {"code": -1, "message": str(e)}


def check(cond, label):
    ok.count += 1
    if cond:
        print(f"  ✓ {label}")
        ok.passed += 1
        return True
    print(f"  ✗ {label}")
    return False


def ok(data, label="RESULT"):
    ok.count += 1
    if data.get("code") == 0:
        print(f"  ✓ {label}")
        ok.passed += 1
        return True
    print(f"  ✗ {label}: {data.get('message', data)}")
    return False


ok.count = 0
ok.passed = 0


def test_login(base_url, username, password):
    print(f"\n{'='*60}")
    print(f"1. 登录测试（获取 Token）")
    print(f"{'='*60}")

    r = req("POST", f"{base_url}/api/v1/staff/login", body={
        "username": username, "password": password,
    })
    ok(r, "1.1 正确登录")
    if r.get("code") != 0:
        print("  当前用户名密码不正确，请通过 --username/--password 指定")
        sys.exit(1)

    data = r["data"]
    print(f"    access_token:  {data['access_token'][:40]}...")
    print(f"    staff_id: {data['staff_id']}, username: {data['username']}")
    return data


def user_login(base_url):
    """获取用户令牌并写入全局 USER_TOKEN。

    orders / payments 走 UserAuthMiddleware，staff 令牌会被判 1004，
    所以业务接口必须用用户身份。注册成功则直接用，已存在则登录。
    """
    global USER_TOKEN
    r = req("POST", f"{base_url}/api/v1/user/auth/register",
            body={"username": USER_NAME, "password": USER_PASSWORD}, auth=False)
    if r.get("code") != 0:
        r = req("POST", f"{base_url}/api/v1/user/auth/login",
                body={"username": USER_NAME, "password": USER_PASSWORD}, auth=False)
    if r.get("code") != 0:
        print(f"  ✗ 1.2 获取用户令牌失败: {r.get('message', r)}")
        sys.exit(1)

    USER_TOKEN = r["data"]["access_token"]
    check(True, f"1.2 获取用户令牌（user_id={r['data']['user_id']}, username={r['data']['username']}）")


def test_create_order(base_url, sku_id):
    print(f"\n{'='*60}")
    print(f"2. 创建订单测试（CreateOrder）")
    print(f"{'='*60}")

    # 2.1 正常创建订单（需用户令牌，见 user_login）
    r = req("POST", f"{base_url}/api/v1/orders", body={
        "user_id": 1,
        "items": [{"sku_id": sku_id, "quantity": 2}],
        "consignee": "张三",
        "phone": "13800138000",
        "province": "广东省",
        "city": "深圳市",
        "district": "南山区",
        "detail_addr": "科技园南区A栋1001",
        "buyer_remark": "测试订单，请勿发货",
        "source": "pc",
    })
    ok(r, "2.1 创建订单")
    if r.get("code") != 0:
        print("  订单创建失败，后续测试可能无法继续")
        return None
    order = r["data"]
    order_no = order.get("order_no", "")
    print(f"    订单号: {order_no}")
    print(f"    订单ID: {order.get('id')}")
    print(f"    总金额: {order.get('total_amount')} 分")
    print(f"    状态: {order.get('status')}")

    # 2.2 创建订单 - 缺少必填参数
    r = req("POST", f"{base_url}/api/v1/orders", body={
        "user_id": 1,
        "items": [{"sku_id": sku_id, "quantity": 1}],
    })
    check(r.get("code") != 0, "2.2 缺少收货人信息被拒绝")

    # 2.3 创建订单 - 空商品列表
    r = req("POST", f"{base_url}/api/v1/orders", body={
        "user_id": 1,
        "items": [],
        "consignee": "张三",
        "phone": "13800138000",
    })
    check(r.get("code") != 0, "2.3 空商品列表被拒绝")

    return order_no


def test_list_orders(base_url):
    print(f"\n{'='*60}")
    print(f"3. 订单列表测试（ListOrders）")
    print(f"{'='*60}")

    # 3.1 正常列表（只有游标分页一种形态；响应不再返回 total）
    r = req("GET", f"{base_url}/api/v1/orders?size=10")
    ok(r, "3.1 订单列表")
    if r.get("code") == 0:
        count = len(r["data"]["list"]) if r["data"]["list"] else 0
        print(f"    当前页{count}条，has_more={r['data'].get('has_more')}")
        check("total" not in r["data"], "3.1.1 响应不再返回 total（避免被误解为总条数）")
        if r["data"]["list"]:
            first = r["data"]["list"][0]
            print(f"    最新订单: {first.get('order_no')}, 状态={first.get('status')}")

    # 3.2 按状态过滤
    r = req("GET", f"{base_url}/api/v1/orders?status=pending&size=10")
    ok(r, "3.2 按状态 pending 过滤")
    if r.get("code") == 0:
        print(f"    待支付订单当前页{len(r['data']['list'] or [])}条")

    # 3.3 按订单号精确查询
    # 先获取一个订单号
    r = req("GET", f"{base_url}/api/v1/orders?size=1")
    if r.get("code") == 0 and r["data"]["list"]:
        order_no = r["data"]["list"][0]["order_no"]
        r = req("GET", f"{base_url}/api/v1/orders?order_no={order_no}")
        ok(r, f"3.3 按订单号 {order_no} 查询")
        if r.get("code") == 0:
            hits = r["data"]["list"] or []
            check(len(hits) == 1 and hits[0]["order_no"] == order_no, "3.3.1 点查命中该单")
            print(f"    查到{len(hits)}条")

    test_list_orders_cursor(base_url)


def test_list_orders_cursor(base_url):
    """3.4~3.9 游标（keyset）分页 —— 列表唯一的分页形态。

    分表后 offset 分页无法跨片归并、COUNT 也要跨片，所以列表只有游标分页：
    响应给出 next_cursor + has_more（两者语义一致，末页同时为空/false），
    **不返回 total**（避免被误解为总条数）；已移除的 page/page_size 传了直接报错。
    字段名与 products 列表完全一致。
    """
    # 3.4 首页：给出 next_cursor + has_more，且没有 total
    r = req("GET", f"{base_url}/api/v1/orders?size=5")
    ok(r, "3.4 游标分页首页")
    first_page_ids, next_cursor = [], ""
    if r.get("code") == 0:
        data = r["data"]
        first_page_ids = [o["id"] for o in data["list"]]
        next_cursor = data.get("next_cursor", "")
        check(len(first_page_ids) == 5, "3.4.1 返回 5 条")
        check("total" not in data, "3.4.2 响应不含 total")
        check(next_cursor != "", "3.4.3 返回 next_cursor")
        check(data.get("has_more") is True, "3.4.4 has_more 与 next_cursor 一致（有下一页）")

    # 3.5 用 next_cursor 翻下一页：与前页不重叠，且 id 更小
    if next_cursor:
        r = req("GET", f"{base_url}/api/v1/orders?size=5&cursor={next_cursor}")
        ok(r, "3.5 用 next_cursor 翻第二页")
        if r.get("code") == 0:
            second_page_ids = [o["id"] for o in r["data"]["list"]]
            check(not set(first_page_ids) & set(second_page_ids), "3.5.1 两页无重叠")
            if second_page_ids:
                check(max(second_page_ids) < min(first_page_ids), "3.5.2 第二页 id 全部小于第一页")

    # 3.6 非法游标必须报错（不能静默当成首页）
    r = req("GET", f"{base_url}/api/v1/orders?size=5&cursor=not-base64!!")
    check(r.get("code") != 0, "3.6 非法游标被拒绝")

    # 3.7 已移除的 page / page_size 必须明确报错（静默忽略 = 每次都返回第一页）
    for legacy in ("page=1&page_size=5", "page=3", "page_size=5"):
        r = req("GET", f"{base_url}/api/v1/orders?{legacy}")
        check(r.get("code") != 0, f"3.7 已移除的 {legacy} 被拒绝")
        if r.get("code") == 0:
            print(f"    ✗ {legacy} 竟被接受: {r['data']}")

    # 3.8 月份窗口：只查一个分片（不自造月份，从列表里取一张单的月份，保证与环境无关）
    r = req("GET", f"{base_url}/api/v1/orders?size=1")
    if r.get("code") == 0 and r["data"]["list"]:
        sample = r["data"]["list"][0]
        month = sample["order_no"][3:9]
        r = req("GET", f"{base_url}/api/v1/orders?month={month}&size=5")
        ok(r, f"3.8 按 {month} 查询")
        if r.get("code") == 0:
            d = r["data"]
            check(d.get("applied_months") == [month], "3.8.1 回显命中月份")
            check(not d.get("window_defaulted"), "3.8.2 显式月份不标记默认窗口")
            check(any(o["order_no"] == sample["order_no"] for o in d["list"]),
                  "3.8.3 该月的单确实出现在结果里")
            check(d.get("applied_from", "").startswith(month[:4]), "3.8.4 回显窗口起点")

        # 3.9 点查不受月份窗口限制（客服拿单号找单）
        r = req("GET", f"{base_url}/api/v1/orders?order_no={sample['order_no']}&month=202501")
        hits = (r.get("data") or {}).get("list") or []
        check(r.get("code") == 0 and len(hits) == 1 and hits[0]["order_no"] == sample["order_no"],
              "3.9 按订单号查询绕开月份窗口")

    # 3.10 不带时间参数：默认窗口兜底 + 回显（避免客户端以为「就只有这么多单」）
    r = req("GET", f"{base_url}/api/v1/orders?size=5")
    if r.get("code") == 0:
        d = r["data"]
        check(d.get("window_defaulted") is True, "3.10.1 未传时间参数时标记为默认窗口")
        check(bool(d.get("applied_from")) and bool(d.get("applied_to")), "3.10.2 回显生效区间")
        check(bool(d.get("applied_months")), "3.10.3 回显命中月份列表")


def test_order_detail(base_url, order_no):
    print(f"\n{'='*60}")
    print(f"4. 订单详情测试（OrderDetail）")
    print(f"{'='*60}")

    if not order_no:
        print("  跳过：无可用订单号")
        return

    # 4.1 正常详情
    r = req("GET", f"{base_url}/api/v1/orders/{order_no}")
    ok(r, f"4.1 查询订单 {order_no} 详情")
    if r.get("code") == 0:
        detail = r["data"]
        order = detail.get("order", {})
        sub_orders = detail.get("sub_orders", [])
        items = detail.get("items", [])
        print(f"    订单号: {order.get('order_no')}")
        print(f"    状态: {order.get('status')}, 支付状态: {order.get('payment_status')}")
        print(f"    子订单数: {len(sub_orders)}")
        print(f"    商品项数: {len(items)}")
        if items:
            item = items[0]
            print(f"    SKU: {item.get('sku_id')}, 名称: {item.get('product_name')}, "
                  f"数量: {item.get('quantity')}, 单价: {item.get('price')} 分")

    # 4.2 不存在的订单号
    r = req("GET", f"{base_url}/api/v1/orders/ORD_NOT_EXIST")
    check(r.get("code") != 0, "4.2 不存在的订单号返回错误")


def test_update_order_status(base_url, order_no):
    print(f"\n{'='*60}")
    print(f"5. 更新订单状态测试（UpdateStatus）")
    print(f"{'='*60}")

    if not order_no:
        print("  跳过：无可用订单号")
        return

    token = None

    # 5.1 取消订单（无需认证）
    r = req("PUT", f"{base_url}/api/v1/orders/{order_no}/status", body={
        "status": "cancelled",
        "note": "测试取消订单",
    })
    ok(r, "5.1 取消订单")
    if r.get("code") == 0:
        order = r["data"]
        print(f"    状态: {order.get('status')}, "
              f"closed_at: {order.get('closed_at')}")

    # 5.2 对已取消的订单再次变更（应被拒绝）
    r = req("PUT", f"{base_url}/api/v1/orders/{order_no}/status", body={
        "status": "paid",
    })
    check(r.get("code") != 0, "5.2 已取消订单不能再次变更")

    # 5.3 不存在的订单号
    r = req("PUT", f"{base_url}/api/v1/orders/ORD_NOT_EXIST/status", body={
        "status": "cancelled",
    })
    check(r.get("code") != 0, "5.3 不存在的订单号返回错误")

    # 5.4 非法状态值
    r = req("PUT", f"{base_url}/api/v1/orders/{order_no}/status", body={
        "status": "invalid_status",
    })
    check(r.get("code") != 0, "5.4 非法状态值被拒绝")


def test_create_payment(base_url, order_no):
    print(f"\n{'='*60}")
    print(f"6. 创建支付测试（CreatePayment）")
    print(f"{'='*60}")

    if not order_no:
        print("  跳过：无可用订单号")
        return None

    # 6.1 正常创建支付（需用户令牌，req 会自动附加）
    # 先查询订单金额
    r = req("GET", f"{base_url}/api/v1/orders/{order_no}")
    if r.get("code") != 0:
        print("  无法获取订单金额，跳过支付测试")
        return None
    pay_amount = r["data"]["order"]["pay_amount"]

    r = req("POST", f"{base_url}/api/v1/payments", body={
        "order_no": order_no,
        "amount": pay_amount,
        "payment_method": "wechat",
        "channel": "wechat_native",
    })
    ok(r, "6.1 创建支付")
    if r.get("code") != 0:
        return None
    payment = r["data"]
    payment_no = payment.get("payment_no", "")
    print(f"    支付单号: {payment_no}")
    print(f"    支付金额: {payment.get('amount')} 分")
    print(f"    支付方式: {payment.get('payment_method')}")
    print(f"    状态: {payment.get('status')}")

    return payment_no


def test_payment_callback(base_url, payment_no):
    print(f"\n{'='*60}")
    print(f"7. 支付回调测试（PaymentCallback）")
    print(f"{'='*60}")

    if not payment_no:
        print("  跳过：无可用支付单号")
        return

    # 7.1 支付成功回调（无需认证）
    # transaction_id 有唯一约束，用支付单号派生以保证脚本可重复执行
    tx_id = f"WX{payment_no}"
    r = req("POST", f"{base_url}/api/v1/payments/callback", body={
        "payment_no": payment_no,
        "transaction_id": tx_id,
        "channel": "wechat_native",
        "status": "success",
        "raw_body": '{"result_code":"SUCCESS","openid":"oXXXX"}',
    })
    ok(r, "7.1 支付成功回调")
    if r.get("code") == 0:
        payment = r["data"]
        print(f"    支付单号: {payment.get('payment_no')}")
        print(f"    状态: {payment.get('status')}")
        print(f"    渠道交易号: {payment.get('transaction_id')}")
        print(f"    支付时间: {payment.get('paid_at')}")

    # 7.2 支付失败回调
    r = req("POST", f"{base_url}/api/v1/payments/callback", body={
        "payment_no": payment_no,
        "transaction_id": f"{tx_id}-F",
        "channel": "wechat_native",
        "status": "failed",
        "failure_reason": "余额不足",
    })
    ok(r, "7.2 支付失败回调")
    if r.get("code") == 0:
        payment = r["data"]
        print(f"    状态: {payment.get('status')}, 失败原因: {payment.get('failure_reason')}")

    # 7.3 不存在的支付单号
    r = req("POST", f"{base_url}/api/v1/payments/callback", body={
        "payment_no": "PAY_NOT_EXIST",
        "transaction_id": "T123",
        "status": "success",
    })
    check(r.get("code") != 0, "7.3 不存在的支付单号返回错误")

    # 7.4 缺少必填参数
    r = req("POST", f"{base_url}/api/v1/payments/callback", body={
        "payment_no": payment_no,
    })
    check(r.get("code") != 0, "7.4 缺少transaction_id/status被拒绝")


def test_get_payment(base_url, payment_no, order_no):
    print(f"\n{'='*60}")
    print(f"8. 查询支付测试（GetPayment）")
    print(f"{'='*60}")

    # 8.1 按支付单号查询
    if payment_no:
        r = req("GET", f"{base_url}/api/v1/payments?payment_no={payment_no}")
        ok(r, "8.1 按支付单号查询")
        if r.get("code") == 0:
            p = r["data"]
            print(f"    支付单号: {p.get('payment_no')}, 状态: {p.get('status')}")

    # 8.2 按订单号查询
    if order_no:
        r = req("GET", f"{base_url}/api/v1/payments?order_no={order_no}")
        ok(r, "8.2 按订单号查询")
        if r.get("code") == 0:
            p = r["data"]
            print(f"    支付单号: {p.get('payment_no')}, 状态: {p.get('status')}")

    # 8.3 不提供任何参数
    r = req("GET", f"{base_url}/api/v1/payments")
    check(r.get("code") != 0, "8.3 不提供任何参数返回错误")


def test_create_refund(base_url, payment_no):
    print(f"\n{'='*60}")
    print(f"9. 创建退款测试（CreateRefund）")
    print(f"{'='*60}")

    if not payment_no:
        print("  跳过：无可用支付单号")
        return

    # 9.1 正常创建退款（需用户令牌，req 会自动附加）
    r = req("POST", f"{base_url}/api/v1/payments/refunds", body={
        "payment_no": payment_no,
        "amount": 1,  # 退款1分
        "reason": "测试退款",
    })
    ok(r, "9.1 创建退款")
    if r.get("code") == 0:
        refund = r["data"]
        print(f"    退款单号: {refund.get('refund_no')}")
        print(f"    退款金额: {refund.get('amount')} 分")
        print(f"    状态: {refund.get('status')}")

    # 9.2 不存在的支付单号
    r = req("POST", f"{base_url}/api/v1/payments/refunds", body={
        "payment_no": "PAY_NOT_EXIST",
        "amount": 100,
    })
    check(r.get("code") != 0, "9.2 不存在的支付单号返回错误")

    # 9.3 无认证（显式不带令牌）
    r = req("POST", f"{base_url}/api/v1/payments/refunds", auth=False, body={
        "payment_no": payment_no,
        "amount": 100,
    })
    check(r.get("code") != 0, "9.3 无Token访问被拒绝")


def main():
    parser = argparse.ArgumentParser(description="gf-eshop 订单/支付/退款 API 测试脚本")
    parser.add_argument("--url", default=BASE_URL, help=f"服务地址 (默认 {BASE_URL})")
    parser.add_argument("--username", default=USERNAME, help=f"用户名 (默认 {USERNAME})")
    parser.add_argument("--password", default=PASSWORD, help=f"密码 (默认 {PASSWORD})")
    parser.add_argument("--sku-id", type=int, default=SKU_ID, help=f"测试SKU ID (默认 {SKU_ID})")
    args = parser.parse_args()

    print(f"✅ 测试服务: {args.url}")
    print(f"✅ 测试SKU: {args.sku_id}")

    # 1. 登录：staff 只验证登录本身，业务接口需要用户令牌
    test_login(args.url, args.username, args.password)
    user_login(args.url)

    # 2. 创建订单（获取订单号）
    order_no = test_create_order(args.url, args.sku_id)

    # 需要重新创建一个状态为 pending 的订单用于后续测试
    if order_no:
        # 3. 订单列表
        test_list_orders(args.url)

        # 4. 订单详情
        test_order_detail(args.url, order_no)

        # 5. 创建另一个订单用于支付/退款测试（刚取消的不能操作）
        order_no2 = test_create_order(args.url, args.sku_id)
        if order_no2:
            # 6. 创建支付
            payment_no = test_create_payment(args.url, order_no2)

            # 7. 支付回调
            test_payment_callback(args.url, payment_no)

            # 8. 查询支付
            test_get_payment(args.url, payment_no, order_no2)

            # 9. 创建退款
            test_create_refund(args.url, payment_no)

        # 也测一下原订单的状态更新
        test_update_order_status(args.url, order_no)

    print(f"\n{'='*60}")
    print(f"测试结果: {ok.passed}/{ok.count} 通过")
    print(f"{'='*60}")

    return 0 if ok.passed == ok.count else 1


if __name__ == "__main__":
    sys.exit(main())