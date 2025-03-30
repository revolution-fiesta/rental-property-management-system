// pages/Info/Info.js
Page({
  data: {
    token: ''
  },
  handleLogin() {
    wx.login({
      success(loginRes) {
        if (!loginRes.code) {
          console.error("登录失败:", loginRes.errMsg);
          return;
        }
        console.log("成功获取 code:", loginRes.code);
        wx.request({
          url: "http://127.0.0.1:8080/login",
          method: "POST",
          data: {
            code: loginRes.code,
            auth_method: "wechat",
          },
          header: {
            "Content-Type": "application/json",
          },
          success(res) {
            console.log("后端返回:", res.data);
          },
          fail(err) {
            console.error("请求失败:", err);
          }
        });
      },
      fail(err) {
        console.error("wx.login 失败:", err);
      }
    });
  },
  naviToOrder() {
    wx.navigateTo({
      url: `/pages/order/order?orderType=外卖&amount=88.88&billDate=2025-03-30`
    });
  }
});
