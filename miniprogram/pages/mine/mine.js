Page({
  data: {
    room: {
      name: "豪华单人间 101",
      type: "1室1厅1卫",
      floor: "10F",
      area: "50",
    },
    tempPassword: "1234-5678",
  },

  // 生成新密码
  changePassword() {
    const newPassword = Math.random().toString().slice(2, 10);
    this.setData({
      tempPassword: newPassword,
    });
    wx.showToast({
      title: "密码已更新",
      icon: "success",
    });
  },

  // 跳转到工单页面
  goToWorkOrder() {
    wx.navigateTo({
      url: "/pages/work_order/work_order",
    });
  },
});
