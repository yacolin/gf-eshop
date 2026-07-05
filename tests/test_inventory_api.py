#!/usr/bin/env python3
"""
gf-eshop 库存业务 API 测试脚本（新的业务方法：Lock/Unlock/Deduct/Restock/GetStock）

用法：
    python tests/test_inventory_api.py
    python tests/test_inventory_api.py --url http://localhost:8000
    python tests/test_inventory_api.py --sku-id 1
"""
import argparse
import json
import sys
import urllib.request
import urllib.error


BASE_URL = "http://localhost:8000"
USERNAME = "admin"
PASSWORD = "123456"
SKU_ID = 1       # 测试用的 SKU ID
WAREHOUSE_ID = 1  # 测试用的仓库 ID


def req(method, url, headers=None, body=None):
    if headers is None:
        headers = {}
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


def test_get_stock(base_url, sku_id):
    print(f"\n{'='*60}")
    print(f"2. 查询库存测试（GetStock）")
    print(f"{'='*60}")

    # 2.1 按 SKU 查询库存（公共接口，无需认证）
    r = req("GET", f"{base_url}/api/v1/inventories/stock?sku_id={sku_id}")
    ok(r, f"2.1 查询 SKU {sku_id} 库存")
    if r.get("code") == 0:
        inv = r["data"]
        print(f"    id={inv.get('id')}, sku_id={inv.get('sku_id')}, "
              f"quantity={inv.get('quantity')}, reserved={inv.get('reserved')}, "
              f"available={inv.get('available')}, status={inv.get('status')}")

    # 2.2 不存在的 SKU
    r = req("GET", f"{base_url}/api/v1/inventories/stock?sku_id=99999999")
    check(r.get("code") != 0, "2.2 不存在的 SKU 返回错误")

    # 2.3 缺少参数
    r = req("GET", f"{base_url}/api/v1/inventories/stock")
    check(r.get("code") != 0, "2.3 缺少 sku_id 参数返回错误")


def test_restock(base_url, token, sku_id, warehouse_id):
    print(f"\n{'='*60}")
    print(f"3. 入库/补货测试（Restock）")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    # 3.1 正常入库
    r = req("POST", f"{base_url}/api/v1/inventories/restock",
            headers=headers, body={
                "sku_id": sku_id,
                "warehouse_id": warehouse_id,
                "quantity": 100,
                "reference_id": "INBOUND-001",
                "operator": "test",
                "note": "测试入库100件",
            })
    ok(r, "3.1 正常入库")
    if r.get("code") == 0:
        inv = r["data"]
        print(f"    入库后 quantity={inv.get('quantity')}, reserved={inv.get('reserved')}")

    # 3.2 多次入库，验证数量累加
    r = req("POST", f"{base_url}/api/v1/inventories/restock",
            headers=headers, body={
                "sku_id": sku_id,
                "warehouse_id": warehouse_id,
                "quantity": 50,
                "reference_id": "INBOUND-002",
                "operator": "test",
            })
    ok(r, "3.2 再次入库50件")
    if r.get("code") == 0:
        inv = r["data"]
        print(f"    再次入库后 quantity={inv.get('quantity')}")

    # 3.3 入库数量为0（应被校验拦截）
    r = req("POST", f"{base_url}/api/v1/inventories/restock",
            headers=headers, body={
                "sku_id": sku_id,
                "quantity": 0,
            })
    check(r.get("code") != 0, "3.3 入库数量为0被拒绝")

    # 3.4 无认证
    r = req("POST", f"{base_url}/api/v1/inventories/restock", body={
        "sku_id": sku_id,
        "quantity": 10,
    })
    check(r.get("code") != 0, "3.4 无Token访问被拒绝")

    return r


def test_lock(base_url, token, sku_id):
    print(f"\n{'='*60}")
    print(f"4. 下单预占库存测试（Lock）")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    # 4.1 正常预占（占用10件）
    r = req("POST", f"{base_url}/api/v1/inventories/lock",
            headers=headers, body={
                "sku_id": sku_id,
                "quantity": 10,
                "reference_id": "ORDER-001",
                "operator": "test",
            })
    ok(r, "4.1 预占库存10件")
    if r.get("code") == 0:
        inv = r["data"]
        print(f"    预占后 quantity={inv.get('quantity')}, reserved={inv.get('reserved')}, "
              f"available={inv.get('available')}")

    # 4.2 预占数量超过可售库存（应被拒绝）
    r = req("POST", f"{base_url}/api/v1/inventories/lock",
            headers=headers, body={
                "sku_id": sku_id,
                "quantity": 999999,
            })
    check(r.get("code") != 0, "4.2 超量预占被拒绝")

    # 4.3 缺少必填参数
    r = req("POST", f"{base_url}/api/v1/inventories/lock",
            headers=headers, body={"sku_id": sku_id})
    check(r.get("code") != 0, "4.3 缺少quantity参数被拒绝")

    return r


def test_unlock(base_url, token, sku_id):
    print(f"\n{'='*60}")
    print(f"5. 取消释放预占测试（Unlock）")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    # 5.1 正常释放（释放5件）
    r = req("POST", f"{base_url}/api/v1/inventories/unlock",
            headers=headers, body={
                "sku_id": sku_id,
                "quantity": 5,
                "reference_id": "ORDER-001",
                "operator": "test",
            })
    ok(r, "5.1 释放预占5件")
    if r.get("code") == 0:
        inv = r["data"]
        print(f"    释放后 quantity={inv.get('quantity')}, reserved={inv.get('reserved')}, "
              f"available={inv.get('available')}")

    # 5.2 释放数量超过预占量（应被拒绝）
    r = req("POST", f"{base_url}/api/v1/inventories/unlock",
            headers=headers, body={
                "sku_id": sku_id,
                "quantity": 999999,
            })
    check(r.get("code") != 0, "5.2 超量释放被拒绝")

    # 5.3 缺少参数
    r = req("POST", f"{base_url}/api/v1/inventories/unlock",
            headers=headers, body={"sku_id": sku_id})
    check(r.get("code") != 0, "5.3 缺少quantity参数被拒绝")

    return r


def test_deduct(base_url, token, sku_id):
    print(f"\n{'='*60}")
    print(f"6. 支付扣减库存测试（Deduct）")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    # 6.1 正常扣减（扣减3件）
    r = req("POST", f"{base_url}/api/v1/inventories/deduct",
            headers=headers, body={
                "sku_id": sku_id,
                "quantity": 3,
                "reference_id": "ORDER-001",
                "operator": "test",
            })
    ok(r, "6.1 扣减库存3件")
    if r.get("code") == 0:
        inv = r["data"]
        print(f"    扣减后 quantity={inv.get('quantity')}, reserved={inv.get('reserved')}, "
              f"available={inv.get('available')}")

    # 6.2 扣减数量超过库存（应被拒绝）
    r = req("POST", f"{base_url}/api/v1/inventories/deduct",
            headers=headers, body={
                "sku_id": sku_id,
                "quantity": 999999,
            })
    check(r.get("code") != 0, "6.2 超量扣减被拒绝")

    # 6.3 缺少参数
    r = req("POST", f"{base_url}/api/v1/inventories/deduct",
            headers=headers, body={"sku_id": sku_id})
    check(r.get("code") != 0, "6.3 缺少quantity参数被拒绝")

    return r


def test_inventory_logs(base_url, sku_id):
    print(f"\n{'='*60}")
    print(f"7. 库存流水查询测试（ListLogs）")
    print(f"{'='*60}")

    # 7.1 查询指定 SKU 的流水（公共接口，无需认证）
    r = req("GET", f"{base_url}/api/v1/inventory_logs?sku_id={sku_id}")
    ok(r, f"7.1 查询 SKU {sku_id} 的流水")
    if r.get("code") == 0:
        total = r["data"]["total"]
        count = len(r["data"]["list"]) if r["data"]["list"] else 0
        print(f"    总共{total}条流水，当前页{count}条")
        if r["data"]["list"]:
            first = r["data"]["list"][0]
            print(f"    最近一条: change_type={first.get('change_type')}, "
                  f"change_amount={first.get('change_amount')}, "
                  f"created_at={first.get('created_at')}")

    # 7.2 按 change_type 过滤
    r = req("GET", f"{base_url}/api/v1/inventory_logs?sku_id={sku_id}&change_type=inbound")
    ok(r, "7.2 按 change_type=inbound 过滤")
    if r.get("code") == 0:
        total = r["data"]["total"]
        print(f"    入库流水共{total}条")

    # 7.3 缺少 sku_id
    r = req("GET", f"{base_url}/api/v1/inventory_logs")
    check(r.get("code") != 0, "7.3 缺少 sku_id 参数返回错误")

    # 7.4 分页测试
    r = req("GET", f"{base_url}/api/v1/inventory_logs?sku_id={sku_id}&page=1&page_size=5")
    ok(r, "7.4 分页查询")
    if r.get("code") == 0:
        count = len(r["data"]["list"]) if r["data"]["list"] else 0
        print(f"    当前页{count}条")


def test_verify_final_stock(base_url, sku_id):
    print(f"\n{'='*60}")
    print(f"8. 最终库存验证")
    print(f"{'='*60}")

    r = req("GET", f"{base_url}/api/v1/inventories/stock?sku_id={sku_id}")
    ok(r, f"8.1 查询 SKU {sku_id} 最终库存")
    if r.get("code") == 0:
        inv = r["data"]
        qty = inv.get("quantity", 0)
        reserved = inv.get("reserved", 0)
        available = inv.get("available", 0)
        print(f"    最终库存: quantity={qty}, reserved={reserved}, available={available}")
        # 验证业务逻辑一致性：available = quantity - reserved
        check(available == qty - reserved, "8.2 可售库存 = 物理库存 - 预占库存")


def main():
    parser = argparse.ArgumentParser(description="gf-eshop 库存业务 API 测试脚本")
    parser.add_argument("--url", default=BASE_URL, help=f"服务地址 (默认 {BASE_URL})")
    parser.add_argument("--username", default=USERNAME, help=f"用户名 (默认 {USERNAME})")
    parser.add_argument("--password", default=PASSWORD, help=f"密码 (默认 {PASSWORD})")
    parser.add_argument("--sku-id", type=int, default=SKU_ID, help=f"测试SKU ID (默认 {SKU_ID})")
    parser.add_argument("--warehouse-id", type=int, default=WAREHOUSE_ID,
                        help=f"测试仓库ID (默认 {WAREHOUSE_ID})")
    args = parser.parse_args()

    print(f"✅ 测试服务: {args.url}")
    print(f"✅ 测试SKU: {args.sku_id}")
    print(f"✅ 测试仓库: {args.warehouse_id}")

    # 1. 登录
    login_data = test_login(args.url, args.username, args.password)
    token = login_data["access_token"]

    # 2. 查询库存（初始状态）
    test_get_stock(args.url, args.sku_id)

    # 3. 入库（先备货）
    test_restock(args.url, token, args.sku_id, args.warehouse_id)

    # 4. 预占库存
    test_lock(args.url, token, args.sku_id)

    # 5. 释放预占
    test_unlock(args.url, token, args.sku_id)

    # 6. 扣减库存
    test_deduct(args.url, token, args.sku_id)

    # 7. 查询流水
    test_inventory_logs(args.url, args.sku_id)

    # 8. 最终验证
    test_verify_final_stock(args.url, args.sku_id)

    print(f"\n{'='*60}")
    print(f"测试结果: {ok.passed}/{ok.count} 通过")
    print(f"{'='*60}")

    return 0 if ok.passed == ok.count else 1


if __name__ == "__main__":
    sys.exit(main())