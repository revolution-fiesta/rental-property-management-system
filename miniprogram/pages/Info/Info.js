// pages/Info/Info.js
Page({
  data: {
    openId: "",
    userType: "",
    avatarUrl: ""
  },

  onLoad() {
    // 获取本地存储的 token 和用户信息
    const openId = wx.getStorageSync("open_id") || "点击此处登录";
    const userType = wx.getStorageSync("role") || "游客";
    const avatar_url = wx.getStorageSync('avatar_url') || "/images/tabbar/info_selected.png"
    // 更新页面数据
    this.setData({
      // 只显示一部分字符串
      openId: openId.substring(0, 6),
      userType: this.convertUserType(userType),
      avatarUrl: avatar_url
    });
  },
  handleLogin() {
    wx.login({
      success: (loginRes) => {
        // 1. 登录获取 code
        if (!loginRes.code) {
          console.error("登录失败:", loginRes.errMsg);
          return;
        }
        console.log("成功获取 code:", loginRes.code);
        // 2. 使用 code 进行 OAuth 登录
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
  },
  // 没法自动获取昵称，需要用户手动填写
  handleChooseAvatar(e) {
    // 先设置头像
    wx.setStorageSync('avatar_url', e.detail.avatarUrl)
    // 然后登录
    this.handleLogin()
  }

});

