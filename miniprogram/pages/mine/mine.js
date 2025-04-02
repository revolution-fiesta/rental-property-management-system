const app = getApp()
Page({
  data: {
    password: "******",
    expires_at: '',
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
        this.getPassword((passwd, exp) => {
          this.setData({
            expires_at: exp,
          })
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
    const token = wx.getStorageSync('token')
    if (!token) {
      wx.showToast({
        title: '请重新登录',
        icon: 'error'
      })
      return
    }
    wx.request({
      url: 'http://localhost:8080/change-room-password',
      method: 'POST',
      header: {
        'Authorization': `Bearer ${token}`,
      },
      data: {
        room_id: Number(this.data.currentRoom.ID),
        new_password: ""
      },
      success: (res) => {
        this.onShow()
        if (res.statusCode == 200) {
          wx.showToast({
            title: "已生成新密码",
            icon: "success",
          });
        }
      },
      fail: (err) => {
        console.log(err)
      }
    })
  },

  showPassword() {
    this.getPassword((passwd, exp) => {
      this.setData({password: passwd})
      setTimeout(()=> {
        this.setData({password: "******"})
      }, 1000)
    })
  },

  //  向后端请求数据并执行回调
  getPassword(callback) {
    const token = wx.getStorageSync('token')
    if (token == null) {
      wx.showToast({
        title: '请重新登录',
        icon: 'error'
      })
    }
    wx.request({
      url: 'http://localhost:8080/get-password',
      method: 'POST',
      header: {
        'Authorization': `Bearer ${token}`,
      },
      data: {
        room_id: this.data.currentRoom.ID
      },
      success: (res)=>{
        if (res.data.password == null) {
          console.log("failed to get password")
          return
        }
        callback( res.data.password.Password,  app.FormatDateToYYYYMMDDHHMMSS(new Date(res.data.password.ExpiresAt)))
      },
      fail: (err)=> {
        console.log(err)
      }
    })
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
