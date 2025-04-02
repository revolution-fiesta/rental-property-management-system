Page({
  data: {
    barPosition: 0,  // 初始时底部栏在屏幕下方
    animationData: {}, // 用来存储动画数据
    isBottomBarVisible: false,
    houseName: '',
    rent: '',
    houseImage: '',
    location: '',
    houseType: '',
    area: '',
    description: '',
    floor: '',
    // 控制弹窗显示与隐藏
    isPopupVisible: false,
    // 初始期数为 6
    quantity: 6,
    room_id: ''
  },

  onLoad(options) {
    // 从上一页传递的参数
    this.animation = wx.createAnimation({
      duration: 300, // 过渡时间
      timingFunction: "ease-in-out"
    });
    this.setData({
      houseName: options.houseName || '未知房源',
      rent: options.rent || '0',
      houseImage: options.houseImage || '/images/houses/house_1.png',
      location: options.location || '未知位置',
      houseType: options.houseType || '未知',
      area: options.area || '0',
      description: options.description || '暂无描述',
      floor: options.floor || '1',
      room_id: options.room_id
    });
  },

  addToFavorites() {
    wx.showToast({ title: '已加入收藏', icon: 'success' });
  },

  contactLandlord() {
    wx.showToast({ title: '已复制管家电话', icon: 'success' });
  },

  toggleBottomBar() {
    if (this.data.isBottomBarVisible) {
      // 关闭动画
      this.animation.translateY(300).step(); // 向下移动隐藏
      this.setData({
        animationData: this.animation.export()
      });
      setTimeout(() => {
        this.setData({ isBottomBarVisible: false });
      }, 300); // 延迟隐藏，等待动画完成
    } else {
      // 先显示底部栏
      this.setData({ isBottomBarVisible: true }, () => {
        // 打开动画
        this.animation.translateY(0).step(); // 回到原位
        this.setData({
          animationData: this.animation.export()
        });
      });
    }
  },

  onFilterButtonCancel() {
    this.toggleBottomBar()
  },

  decreaseQuantity() {
    if (this.data.quantity > 6) {
      this.setData({
        quantity: this.data.quantity - 1
      });
    } else {
      wx.showToast({
        title: '不能少于六个月',
        icon: 'error'
      })
    }
  },

  // 增加期数
  increaseQuantity() {
    this.setData({
      quantity: this.data.quantity + 1
    });
  },

  // 确定签约
  onConfirm() {
    wx.navigateTo({
      url: `/pages/contract/contract?room_id=${this.data.room_id}&num_terms=${this.data.quantity}&rent=${this.data.rent}`,
    })
  }
});
