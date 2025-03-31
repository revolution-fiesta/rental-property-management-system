Page({
  data: {
    bills: []
  },
  onLoad() {
    // 预置账单数据
    this.setData({
      bills: [
        { id: "1001", name: "水电费", amount: 120.50, paid: true, date: "2025-03-30" },
        { id: "1002", name: "房租", amount: 2500, paid: false, date: "2025-04-01" },
        { id: "1003", name: "物业费", amount: 300, paid: true, date: "2025-03-28" },
      ]
    });
  },
  naviToOrder() {
    wx.navigateTo({
      url: `/pages/order/order?orderType=外卖&amount=88.88&billDate=2025-03-30`
    });
  },
});
