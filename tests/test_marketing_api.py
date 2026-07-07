#!/usr/bin/env python3
"""
gf-eshop 营销模块（促销/优惠券/秒杀）API 测试脚本

用法：
    python tests/test_marketing_api.py
    python tests/test_marketing_api.py --url http://localhost:8000
"""
import argparse
import json
import sys
import urllib.error
import urllib.request


BASE_URL = "http://localhost:8000"


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


def register_or_login(base_url):
    print(f"\n{'='*60}")
    print(f"1. 用户注册/登录")
    print(f"{'='*60}")

    r = req("POST", f"{base_url}/api/v1/user/auth/register", body={
        "username": "test_marketing",
        "password": "test123456",
    })
    if r.get("code") == 0:
        data = r["data"]
        print(f"  注册成功: user_id={data['user_id']}, username={data['username']}")
        return data

    r = req("POST", f"{base_url}/api/v1/user/auth/login", body={
        "username": "test_marketing",
        "password": "test123456",
    })
    if r.get("code") == 0:
        data = r["data"]
        print(f"  登录成功: user_id={data['user_id']}, username={data['username']}")
        return data

    print(f"  注册/登录失败: {r.get('message', r)}")
    sys.exit(1)


def test_promotion_list(base_url, token):
    print(f"\n{'='*60}")
    print(f"2. 促销列表测试")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    r = req("GET", f"{base_url}/api/v1/promotions", headers=headers)
    ok(r, "2.1 促销列表")
    if r.get("code") == 0:
        total = r["data"]["total"]
        count = len(r["data"]["list"]) if r["data"]["list"] else 0
        print(f"    总共{total}条，当前页{count}条")
        if r["data"]["list"]:
            p = r["data"]["list"][0]
            print(f"    最新: id={p['id']}, name={p['promo_name']}, type={p['promo_type']}, status={p['status']}")

    r = req("GET", f"{base_url}/api/v1/promotions?page=1&page_size=5&status=2", headers=headers)
    ok(r, "2.2 按状态筛选（生效中）")
    if r.get("code") == 0:
        print(f"    生效中促销共{r['data']['total']}条")

    r = req("GET", f"{base_url}/api/v1/promotions?promo_type=3", headers=headers)
    ok(r, "2.3 按类型筛选（秒杀）")
    if r.get("code") == 0:
        print(f"    秒杀类型共{r['data']['total']}条")


def test_promotion_create(base_url, token):
    print(f"\n{'='*60}")
    print(f"3. 创建促销测试")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    r = req("POST", f"{base_url}/api/v1/promotions", headers=headers, body={
        "promo_name": "API测试-满减券",
        "promo_type": 1,
        "start_time": "2026-01-01 00:00:00",
        "end_time": "2026-12-31 23:59:59",
        "total_quantity": 1000,
        "per_user_limit": 1,
        "condition_type": 2,
        "condition_value": 20000,
        "benefit_type": 1,
        "benefit_value": 3000,
        "is_stackable": 0,
    })
    ok(r, "3.1 创建促销")
    if r.get("code") == 0:
        promo_id = r["data"]["id"]
        print(f"    创建成功: promo_id={promo_id}")
        return promo_id

    print("    创建失败，后续促销测试可能无法继续")
    return None


def test_promotion_detail(base_url, token, promo_id):
    print(f"\n{'='*60}")
    print(f"4. 促销详情测试")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    if not promo_id:
        print("  跳过：无促销ID")
        return

    r = req("GET", f"{base_url}/api/v1/promotions/{promo_id}", headers=headers)
    ok(r, f"4.1 促销详情 id={promo_id}")
    if r.get("code") == 0:
        p = r["data"]
        print(f"    名称: {p['promo_name']}, 类型: {p['promo_type']}, 状态: {p['status']}")

    r = req("GET", f"{base_url}/api/v1/promotions/{promo_id}/detail", headers=headers)
    ok(r, f"4.2 促销完整详情 id={promo_id}")
    if r.get("code") == 0:
        d = r["data"]
        print(f"    规则: {d.get('rule', {}).get('rule_name', '无')}")
        print(f"    商品数: {len(d.get('products', []))}")

    r = req("GET", f"{base_url}/api/v1/promotions/99999999", headers=headers)
    check(r.get("code") != 0, "4.3 不存在的ID返回错误")


def test_promotion_update(base_url, token, promo_id):
    print(f"\n{'='*60}")
    print(f"5. 更新促销测试")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    if not promo_id:
        print("  跳过：无促销ID")
        return

    r = req("PUT", f"{base_url}/api/v1/promotions/{promo_id}", headers=headers, body={
        "promo_name": "API测试-已更新",
        "total_quantity": 2000,
    })
    ok(r, f"5.1 更新促销 id={promo_id}")
    if r.get("code") == 0:
        print(f"    更新成功")

    r = req("PUT", f"{base_url}/api/v1/promotions/99999999", headers=headers, body={
        "promo_name": "不存在",
    })
    check(r.get("code") != 0, "5.2 不存在的ID返回错误")


def test_promotion_delete(base_url, token, promo_id):
    print(f"\n{'='*60}")
    print(f"6. 删除促销测试")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    if not promo_id:
        print("  跳过：无促销ID")
        return

    r = req("DELETE", f"{base_url}/api/v1/promotions/{promo_id}", headers=headers)
    ok(r, f"6.1 删除促销 id={promo_id}")
    if r.get("code") == 0:
        print(f"    删除成功")

    r = req("GET", f"{base_url}/api/v1/promotions/{promo_id}", headers=headers)
    check(r.get("code") != 0, "6.2 删除后查询返回错误")

    r = req("DELETE", f"{base_url}/api/v1/promotions/99999999", headers=headers)
    check(r.get("code") != 0, "6.3 不存在的ID返回错误")


def test_coupon_flow(base_url, token):
    print(f"\n{'='*60}")
    print(f"7. 优惠券领取/使用测试")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    r = req("GET", f"{base_url}/api/v1/promotions?promo_type=1&status=2&page_size=5", headers=headers)
    if r.get("code") != 0 or not r["data"]["list"]:
        print("  跳过：无可用优惠券类型促销")
        return

    claimed_ids = set()
    r2 = req("GET", f"{base_url}/api/v1/coupons/me?page=1&page_size=100", headers=headers)
    if r2.get("code") == 0 and r2["data"]["list"]:
        for up in r2["data"]["list"]:
            claimed_ids.add(up["promotion_id"])

    unclaimed = [p for p in r["data"]["list"] if p["id"] not in claimed_ids]
    if not unclaimed:
        print("  跳过：所有优惠券均已领取")
        return

    promo = unclaimed[0]
    promo_id = promo["id"]
    print(f"  使用促销 id={promo_id}, name={promo['promo_name']}")

    r = req("POST", f"{base_url}/api/v1/coupons/claim", headers=headers, body={
        "promotion_id": promo_id,
    })
    ok(r, "7.1 领取优惠券")

    r = req("POST", f"{base_url}/api/v1/coupons/claim", headers=headers, body={
        "promotion_id": 99999999,
    })
    check(r.get("code") != 0, "7.2 不存在的促销返回错误")

    r = req("GET", f"{base_url}/api/v1/coupons/me?page=1&page_size=10", headers=headers)
    ok(r, "7.3 我的优惠券列表")
    if r.get("code") == 0:
        total = r["data"]["total"]
        print(f"    共{total}张优惠券")
        if r["data"]["list"]:
            up = r["data"]["list"][0]
            print(f"    最新: id={up['id']}, promo_id={up['promotion_id']}, status={up['status']}")

    r = req("GET", f"{base_url}/api/v1/coupons/me?status=1", headers=headers)
    ok(r, "7.4 筛选未使用优惠券")
    if r.get("code") == 0:
        print(f"    未使用: {r['data']['total']}张")

    r = req("GET", f"{base_url}/api/v1/coupons/me?status=2", headers=headers)
    ok(r, "7.5 筛选已使用优惠券")
    if r.get("code") == 0:
        print(f"    已使用: {r['data']['total']}张")


def test_flash_buy(base_url, token):
    print(f"\n{'='*60}")
    print(f"8. 秒杀下单测试")
    print(f"{'='*60}")

    headers = {"Authorization": f"Bearer {token}"}

    r = req("GET", f"{base_url}/api/v1/promotions?promo_type=3&status=2&page_size=5", headers=headers)
    if r.get("code") != 0 or not r["data"]["list"]:
        print("  跳过：无可用秒杀促销")
        return

    flash_promo = r["data"]["list"][0]
    promo_id = flash_promo["id"]
    print(f"  使用秒杀 id={promo_id}, name={flash_promo['promo_name']}")

    r = req("POST", f"{base_url}/api/v1/flash/buy", headers=headers, body={
        "promotion_id": promo_id,
        "product_id": 1,
        "sku_id": 1,
        "quantity": 1,
    })
    if r.get("code") == 0:
        token_val = r["data"]["token"]
        print(f"    秒杀成功: token={token_val[:20]}...")

        r = req("POST", f"{base_url}/api/v1/flash/confirm", headers=headers, body={
            "token": token_val,
            "address_id": 1,
        })
        ok(r, "8.2 确认秒杀订单")
        if r.get("code") == 0:
            print(f"    确认成功")
        else:
            print(f"  跳过8.2: {r.get('message')}")
    else:
        print(f"  跳过8.2: 秒杀下单失败 ({r.get('message')})")

    r = req("POST", f"{base_url}/api/v1/flash/buy", headers=headers, body={
        "promotion_id": 99999999,
        "product_id": 1,
        "sku_id": 1,
        "quantity": 1,
    })
    check(r.get("code") != 0, "8.3 不存在的促销返回错误")


def test_auth_required(base_url):
    print(f"\n{'='*60}")
    print(f"9. 认证校验测试")
    print(f"{'='*60}")

    r = req("GET", f"{base_url}/api/v1/promotions")
    check(r.get("code") != 0, "9.1 未登录访问促销列表被拒绝")

    r = req("POST", f"{base_url}/api/v1/coupons/claim", body={"promotion_id": 1})
    check(r.get("code") != 0, "9.2 未登录领取优惠券被拒绝")

    r = req("POST", f"{base_url}/api/v1/flash/buy", body={
        "promotion_id": 1, "product_id": 1, "sku_id": 1, "quantity": 1,
    })
    check(r.get("code") != 0, "9.3 未登录秒杀被拒绝")

    r = req("GET", f"{base_url}/api/v1/promotions",
            headers={"Authorization": "Bearer invalid-token"})
    check(r.get("code") != 0, "9.4 无效Token被拒绝")


def main():
    parser = argparse.ArgumentParser(description="gf-eshop 营销模块API测试脚本")
    parser.add_argument("--url", default=BASE_URL, help=f"服务地址 (默认 {BASE_URL})")
    args = parser.parse_args()

    print(f"✅ 测试服务: {args.url}")

    login_data = register_or_login(args.url)
    token = login_data["access_token"]

    test_promotion_list(args.url, token)
    promo_id = test_promotion_create(args.url, token)
    test_promotion_detail(args.url, token, promo_id)
    test_promotion_update(args.url, token, promo_id)
    test_promotion_delete(args.url, token, promo_id)

    test_coupon_flow(args.url, token)
    test_flash_buy(args.url, token)
    test_auth_required(args.url)

    print(f"\n{'='*60}")
    print(f"测试结果: {ok.passed}/{ok.count} 通过")
    print(f"{'='*60}")

    return 0 if ok.passed == ok.count else 1


if __name__ == "__main__":
    sys.exit(main())
