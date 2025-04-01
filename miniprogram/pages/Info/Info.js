// pages/Info/Info.js
Page({
  data: {
    openId: "",
    userType: ""
  },
  onLoad() {
    // 获取本地存储的 token 和用户信息
    const openId = wx.getStorageSync("open_id") || "点击此处登录";
    const userType = wx.getStorageSync("role") || "游客";

    // 更新页面数据
    this.setData({
      openId: openId.substring(0, 6),
      userType: this.convertUserType(userType),
    });
  },
  handleLogin() {
    wx.login({
      success: (loginRes) => {
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
          success: (res) => {
            console.log("后端返回:", res.data);
            // 存储 token 和用户信息到本地
            wx.setStorageSync("token", res.data.token);
            wx.setStorageSync("open_id", res.data.open_id)
            wx.setStorageSync("role", res.data.role)
            wx.showToast({
              title: '登录成功',
            })
            // 重新加载
            this.onLoad()
          },
          fail: (err) => {
            wx.showToast({
              title: err,
            })
          },
        });
      },
      fail: (err) => {
        console.error("wx.login 失败:", err);
      },
    });
  },
  naviToAdminPanel() {
    wx.navigateTo({
      url: '/pages/admin_page/admin_page',
    })

  },
  convertUserType(userType) {
    const mapping = {
      member: "普通用户",
      admin: "管理员"
    };
    return mapping[userType] || "未知身份";
  }
});
