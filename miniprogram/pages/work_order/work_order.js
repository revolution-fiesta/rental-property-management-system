Page({
  data: {
    workOrderText: "",
    workOrders: [],
  },

  onLoad() {
    // 先检查本地存储是否已有工单
    let workOrders = wx.getStorageSync("workOrders");
    if (!workOrders || workOrders.length === 0) {
      // 预设一些工单数据
      workOrders = [
        {
          id: 1,
          date: "2025-03-30 10:15",
          description: "空调无法启动，可能是遥控器没电了。",
        },
        {
          id: 2,
          date: "2025-03-29 14:30",
          description: "水龙头有滴水现象，需要维修。",
        },
        {
          id: 3,
          date: "2025-03-28 09:50",
          description: "房门锁有点松动，希望安排维修人员。",
        },
      ];
      wx.setStorageSync("workOrders", workOrders);
    }
    this.setData({ workOrders });
  },

  // 监听输入
  onInput(e) {
    this.setData({ workOrderText: e.detail.value });
  },

  // 提交工单
  submitOrder() {
    if (!this.data.workOrderText.trim()) {
      wx.showToast({ title: "请输入工单内容", icon: "none" });
      return;
    }

    const newOrder = {
      id: Date.now(),
      date: new Date().toLocaleString(),
      description: this.data.workOrderText,
    };

    const workOrders = [...this.data.workOrders, newOrder];
    wx.setStorageSync("workOrders", workOrders);
    this.setData({ workOrders, workOrderText: "" });

    wx.showToast({ title: "工单提交成功", icon: "success" });
  },
});
