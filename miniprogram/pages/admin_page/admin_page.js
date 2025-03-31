Page({
  data: {
    features: [
      { icon: "/images/ui/ticket.svg", name: "工单管理", page: "/pages/ticket_list/ticket_list" },
      { icon: "/images/ui/admin.png", name: "用户管理", page: "/pages/users/users" },
      { icon: "/images/ui/stat.svg", name: "数据统计", page: "/pages/stats/stats" },
      { icon: "/images/ui/settings.png", name: "系统设置", page: "/pages/settings/settings" },
      { icon: "/images/ui/logs.svg", name: "日志管理", page: "/pages/logs/logs" },
      { icon: "/images/ui/notification.svg", name: "通知中心", page: "/pages/notifications/notifications" },
    ]
  },

  navigateToFeature(e) {
    const page = e.currentTarget.dataset.page;
    if (page) {
      wx.navigateTo({ url: page });
    }
  }
});
