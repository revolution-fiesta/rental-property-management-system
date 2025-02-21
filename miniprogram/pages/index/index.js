// index.js
Page({
  data: {
    barPosition: 0,  // 初始时底部栏在屏幕下方
    animationData: {}, // 用来存储动画数据

  },

  toggleBottomBar() {
    const animation = wx.createAnimation({
      duration: 300, // 动画持续时间
      timingFunction: 'ease-in-out' // 动画过渡效果
    });

    // 判断当前底部栏是否显示，并执行相应的动画
    if (this.data.isBottomBarVisible) {
      animation.translateY(300).step(); // 隐藏底部栏
    } else {
      animation.translateY(0).step(); // 弹出底部栏
    }

    // 更新动画数据
    this.setData({
      animationData: animation.export(),
      isBottomBarVisible: !this.data.isBottomBarVisible // 切换底部栏显示状态
    });
  },
  goToDetail: function () {
    wx.navigateTo({
      url: '/pages/details/details'
    });
  }
})