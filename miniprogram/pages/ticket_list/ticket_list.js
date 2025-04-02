const app = getApp()

Page({
  data: {
    workOrders: [],
    roomNames: []
  },

  onShow() {
    this.loadWorkOrders()
  },

  loadWorkOrders() {
    const token = wx.getStorageSync('token');
    wx.request({
      url: 'http://localhost:8080/list-admin-workorders',
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`,
      },
      // TODO: 写成同步的不然要多次更新数据太烦了
      success: (res) => {
        for (let index = 0; index < res.data.work_orders.length; index++) {
          const element = res.data.work_orders[index];
          wx.request({
            url: 'http://localhost:8080/get-room',
            method: "POST",
            header: {
              'Authorization': `Bearer ${token}`,
            },
            data: {
              "room_id": res.data.work_orders[index].RoomID
            },
            // TODO: 部分 IOS 不适配
            success: (roomRes) => {
              let arr = this.data.roomNames
              arr.push(roomRes.data.room.Name)
              this.setData({ roomNames: arr })
              this.setData({
                workOrders: res.data.work_orders
                  .map((work_order, index) => ({
                    Date: app.FormatDateToYYYYMMDDHHMMSS(new Date(work_order.CreatedAt)),  // 格式化 CreatedAt
                    Description: work_order.Description,
                    Type: work_order.Type,
                    ID: work_order.ID,
                    Resolved: work_order.Status == "completed",
                    RoomName: this.data.roomNames[index]
                  }))
                  .sort((a, b) => new Date(b.Date) - new Date(a.Date)) // 按时间降序排序（最新的在前）
              });
            },
            fail: () => { }
          })
        }
      },
      fail(error) {
        // 请求失败时的回调
        console.log('请求失败', error);
      }
    });
  },

  markAsResolved(e) {
    const work_order_id = e.currentTarget.dataset.id;
    const token = wx.getStorageSync('token')
    wx.request({
      url: 'http://localhost:8080/update-workorder',
      method: "POST",
      header: {
        'Authorization': `Bearer ${token}`,
      },
      data: {
        work_order_id: work_order_id,
        status: "completed"
      },
      success: (res) => {
        console.log(res.data)
        wx.showToast({
          title: '成功更新状态',
        })
        this.onShow()
      },
      fail: (e) => { console.log(e)}
    })

  }
});
