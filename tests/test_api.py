#!/usr/bin/env python3
"""
gf-eshop API 自动化测试脚本
用法：
    python tests/test_api.py              # 测试所有接口
    python tests/test_api.py --url http://localhost:8000
    python tests/test_api.py --username admin --password 123456
"""
import argparse
import json
import sys
import urllib.request
import urllib.error


BASE_URL = "http://localhost:8000"
USERNAME = "admin"
PASSWORD = "123456"


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
    print(f"1. 登录测试")
    print(f"{'='*60}")

    # 1.1 正确登录
    r = req("POST", f"{base_url}/api/v1/staff/login", body={
        "username": username, "password": password,
    })
    ok(r, "1.1 正确登录")
    if r.get("code") != 0:
        print("  当前用户名密码不正确，请通过 --username/--password 指定")
        sys.exit(1)

    data = r["data"]
    assert "access_token" in data and "refresh_token" in data
    print(f"    access_token:  {data['access_token'][:40]}... ({data['expire_in']}s)")
    print(f"    refresh_token: {data['refresh_token'][:40]}... ({data['refresh_in']}s)")
    print(f"    staff_id: {data['staff_id']}, username: {data['username']}")
    return data


def test_profile(base_url, token):
    print(f"\n{'='*60}")
    print(f"2. 获取用户信息测试")
    print(f"{'='*60}")

    # 2.1 带有效token获取
    r = req("GET", f"{base_url}/api/v1/staff/profile",
            headers={"Authorization": f"Bearer {token}"})
    ok(r, "2.1 带有效Token获取用户信息")
    if r.get("code") == 0:
        print(f"    username: {r['data']['username']}, real_name: {r['data']['real_name']}")

    # 2.2 无token访问
    r = req("GET", f"{base_url}/api/v1/staff/profile")
    check(r.get("code") == 401, "2.2 未带Token访问被拒绝")

    # 2.3 无效token
    r = req("GET", f"{base_url}/api/v1/staff/profile",
            headers={"Authorization": "Bearer invalid-token"})
    check(r.get("code") == 401, "2.3 无效Token访问被拒绝")


def test_refresh(base_url, refresh_token):
    print(f"\n{'='*60}")
    print(f"3. 刷新Token测试")
    print(f"{'='*60}")

    # 3.1 正确刷新
    r = req("POST", f"{base_url}/api/v1/staff/refresh", body={
        "refresh_token": refresh_token,
    })
    ok(r, "3.1 正确刷新Token")
    if r.get("code") != 0:
        return None
    new_data = r["data"]
    print(f"    new access_token:  {new_data['access_token'][:40]}...")
    print(f"    new refresh_token: {new_data['refresh_token'][:40]}...")

    # 3.2 旧refresh_token被拒绝
    r = req("POST", f"{base_url}/api/v1/staff/refresh", body={
        "refresh_token": refresh_token,
    })
    check(r.get("code") != 0, "3.2 旧RefreshToken被拒绝")

    return new_data


def test_brands(base_url, token):
    print(f"\n{'='*60}")
    print(f"4. 品牌接口测试（无需认证）")
    print(f"{'='*60}")

    # 4.1 品牌列表
    r = req("GET", f"{base_url}/api/v1/brands")
    ok(r, "4.1 品牌列表获取")
    if r.get("code") == 0:
        total = r['data']['total']
        count = len(r['data']['list']) if r['data']['list'] else 0
        print(f"    总共{total}条记录，当前页{count}条")


def test_logout(base_url, token):
    print(f"\n{'='*60}")
    print(f"5. 退出登录测试")
    print(f"{'='*60}")

    # 5.1 退出登录
    r = req("POST", f"{base_url}/api/v1/staff/logout",
            headers={"Authorization": f"Bearer {token}"})
    ok(r, "5.1 退出登录")

    return r


def test_public_endpoints(base_url):
    print(f"\n{'='*60}")
    print(f"6. 其他公共接口测试")
    print(f"{'='*60}")

    r = req("GET", f"{base_url}/hello")
    check(r.get("code") == 0 or "raw" in r, "6.1 Hello接口")

    r = req("GET", f"{base_url}/api/v1/categories")
    ok(r, "6.2 分类列表")


def main():
    parser = argparse.ArgumentParser(description="gf-eshop API 自动化测试脚本")
    parser.add_argument("--url", default=BASE_URL, help=f"服务地址 (默认 {BASE_URL})")
    parser.add_argument("--username", default=USERNAME, help=f"用户名 (默认 {USERNAME})")
    parser.add_argument("--password", default=PASSWORD, help=f"密码 (默认 {PASSWORD})")
    args = parser.parse_args()

    print(f"✅ 测试服务: {args.url}")
    print(f"✅ 测试账号: {args.username}")

    login_data = test_login(args.url, args.username, args.password)
    test_profile(args.url, login_data["access_token"])
    refresh_data = test_refresh(args.url, login_data["refresh_token"])
    test_brands(args.url, login_data["access_token"])

    if refresh_data:
        test_logout(args.url, refresh_data["access_token"])
    else:
        test_logout(args.url, login_data["access_token"])

    test_public_endpoints(args.url)

    print(f"\n{'='*60}")
    print(f"测试结果: {ok.passed}/{ok.count} 通过")
    print(f"{'='*60}")

    return 0 if ok.passed == ok.count else 1


if __name__ == "__main__":
    sys.exit(main())
