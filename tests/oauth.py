import requests
import json

# 示例配置
class Config:
    APP_ID = "wx937c1c53ceabe38b"
    APP_SECRET = "ad2225a5932ee1d94e8471ad695a7122"

def get_wechat_openid(code: str) -> str:
    params = {
        "appid": Config.APP_ID,
        "secret": Config.APP_SECRET,
        "js_code": code,
        "grant_type": "authorization_code"
    }
    
    url = "https://api.weixin.qq.com/sns/jscode2session"
    
    try:
        resp = requests.get(url, params=params)
        resp.raise_for_status()  # 检查 HTTP 状态码
    except requests.RequestException as e:
        raise Exception(f"Failed to send OAuth request: {e}")

    try:
        oauth_resp = resp.json()
    except json.JSONDecodeError:
        raise Exception("Failed to parse response JSON")

    if "openid" not in oauth_resp:
        raise Exception(f"WeChat API error: {oauth_resp}")

    return oauth_resp["openid"]

# 示例测试
if __name__ == "__main__":
    test_code = "0f3itq000an5ZT13OO30082ESA1itq0X"
    try:
        openid = get_wechat_openid(test_code)
        print(f"OpenID: {openid}")
    except Exception as e:
        print(f"Error: {e}")
