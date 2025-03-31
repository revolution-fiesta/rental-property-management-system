// index.js
Page({
  data: {
    barPosition: 0,  // 初始时底部栏在屏幕下方
    animationData: {}, // 用来存储动画数据
    isBottomBarVisible: false,
    propertyList: {},
    // 用于进行筛选的页面数据
    areaMin: 0,  // 默认最小面积
    areaMax: 500, // 默认最大面积
    priceMin: 0,  // 默认最低价格
    priceMax: 30000,  // 默认最高价格
    roomTypes: ["全部", "一室", "一室一厅", "两室一厅"], // 房型选项
    selectedRoomType: "全部", // 默认显示
    searchbarText: ""
  },
  onLoad() {
    this.fetchProperties()
    this.animation = wx.createAnimation({
      duration: 300, // 过渡时间
      timingFunction: "ease-in-out"
    });
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
  goToDetail(event) {
    const item_idx = event.currentTarget.dataset.index
    const room_obj = this.data.propertyList[item_idx]
    wx.navigateTo({
      url: `/pages/details/details?houseName=${room_obj.Name}&rent=${room_obj.Price}&area=${room_obj.Area}&houseType=${room_obj.Type}&floor=${room_obj.Floor}`,
    })
  },
  fetchProperties: function () {
    const that = this;
    // 假设你有一个 API 返回数据
    wx.request({
      url: 'http://localhost:8080/list-rooms',  // 替换为实际的 API 地址
      method: 'GET',
      success(res) {
        if (res.data) {
          const mappedData = res.data.rooms.map(item => ({
            Available: item.Available,
            Floor: item.Floor,
            ID: item.ID,
            Name: item.Name,
            Price: item.Price,
            Tags: item.Tags,
            Type: convertRoomTypeReverse(item.Type),
            Area: item.Area
          }));
          that.setData({
            propertyList: mappedData
          });
        }
      },
      fail(error) {
        console.error('请求失败:', error);
      }
    });
  },
  // 下面五个函数用于过滤查询
  onAreaMinChange(e) {
    const newMin = e.detail.value;
    this.setData({
      areaMin: Math.min(newMin, this.data.areaMax - 10) // 确保最小值不会超过最大值
    });
  },
  onAreaMaxChange(e) {
    const newMax = e.detail.value;
    this.setData({
      areaMax: Math.max(newMax, this.data.areaMin + 10) // 确保最大值不会低于最小值
    });
  },
  onPriceMinChange(e) {
    const newMin = e.detail.value;
    this.setData({
      priceMin: Math.min(newMin, this.data.priceMax - 100) // 保证最小值不超过最大值
    });
  },
  onPriceMaxChange(e) {
    const newMax = e.detail.value;
    this.setData({
      priceMax: Math.max(newMax, this.data.priceMin + 100) // 保证最大值不小于最小值
    });
  },
  onRoomTypeChange(e) {
    this.setData({
      selectedRoomType: this.data.roomTypes[e.detail.value]
    });
  },
  // 发送过滤条件计算后的房间信息
  fetchFilteredProperties: function () {
    const that = this;
    // 假设你有一个 API 返回数据
    wx.request({
      url: 'http://localhost:8080/list-filtered-rooms',  // 替换为实际的 API 地址
      method: 'POST',
      data: {
        "room_type": convertRoomType(this.data.selectedRoomType),
        "min_price": this.data.priceMin,
        "max_price": this.data.priceMax,
        "min_area": this.data.areaMin,
        "max_area": this.data.areaMax,
        "keyword": this.data.searchbarText == "" ? null: this.data.searchbarText
      },
      success(res) {
        if (res.data ) {
          const mappedData = res.data.rooms.map(item => ({
            Available: item.Available,
            Floor: item.Floor,
            ID: item.ID,
            Name: item.Name,
            Price: item.Price,
            Tags: item.Tags,
            Type: convertRoomTypeReverse(item.Type),
            Area: item.Area
          }));
          that.setData({
            propertyList: mappedData
          });
        }
      },
      fail(error) {
        console.error('请求失败:', error);
      }
    });
  },
  // 过滤器点击取消按钮
  onFilterButtonCancel() {
    this.fetchProperties()
    this.toggleBottomBar()
  },
  // 搜索框进行文字输入
  onSearchbarInput(e) {
    this.setData({
      searchbarText: e.detail.value
    });
    this.fetchFilteredProperties()
  },
})
function convertRoomType(roomType) {
  // 定义映射关系
  const mapping = {
    "一室": "b1",
    "一室一厅": "b1l1",
    "两室一厅": "b2l1",
  };
  // 返回对应的中文描述，找不到则返回 "未知房型"
  return mapping[roomType];
}

function convertRoomTypeReverse(roomTypeCode) {
  // 定义反向映射关系
  const reverseMapping = {
    "b1": "一室",
    "b1l1": "一室一厅",
    "b2l1": "两室一厅",
  };
  // 返回对应的房型，找不到则返回 "未知房型"
  return reverseMapping[roomTypeCode] || "未知房型";
}