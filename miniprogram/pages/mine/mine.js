const app = getApp()
Page({
  data: {
    tempPassword: "1234-5678",
    roomList: [], // 存储房间列表
    currentRoom: null, // 当前选择的房间
  },

  onShow() {
    this.loadCurrentRoom();
  },

  loadCurrentRoom() {
    const token = wx.getStorageSync('token');
    wx.request({
      url: 'http://localhost:8080/list-owned-rooms',
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`,
      },
      success: (res) => {
        if (res.data && res.data.rooms) {
          this.setData({
            roomList: res.data.rooms.map((room, idx) => ({
              Name: room.Name,
              FormatType: app.ConvertRoomTypeReverse(room.Type),
              Area: room.Area,
              Floor: room.Floor,
              Index: idx,
              ID: room.ID
            })),
          });
        }
        this.setData({
          currentRoom: this.data.roomList.length > 0 ? this.data.roomList[0] : null,
        })
      },
      fail(error) {
        console.log('请求失败', error);
      }
    });
  },

  switchRoom(e) {
    const index = e.detail.value;
    this.setData({
      currentRoom: this.data.roomList[index],
    });
  },

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

  goToWorkOrder() {
    wx.navigateTo({
      url: `/pages/work_order/work_order?room_id=${this.data.currentRoom.ID}`,
    });
  },

  terminateLease() {
    wx.showToast({
      title: '退租申请已发起',
    });
  }
});
