Page({
  data: {
    houseName: '',
    rent: '',
    houseImage: '',
    location: '',
    houseType: '',
    area: '',
    description: '',
    floor: ''
  },

  onLoad(options) {
    // 从上一页传递的参数
    this.setData({
      houseName: options.houseName || '未知房源',
      rent: options.rent || '0',
      houseImage: options.houseImage || '/images/houses/house_1.png',
      location: options.location || '未知位置',
      houseType: options.houseType || '未知',
      area: options.area || '0',
      description: options.description || '暂无描述',
      floor: options.floor || '1'
    });
  },

  addToFavorites() {
    wx.showToast({ title: '已加入收藏', icon: 'success' });
  },

  contactLandlord() {
    wx.showToast({ title: '已复制管家电话', icon: 'success' });
  }
});
