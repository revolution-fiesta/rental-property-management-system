const app = getApp()

Page({
  data: {
    workOrderText: "",
    workOrders: [],
    roomID: null
  },

  onShow(){
    this.loadWorkOrders()
  },

  onLoad(opts) {
    this.setData({
      roomID: opts.room_id
    })
  },

  loadWorkOrders(){
    const token = wx.getStorageSync('token');
    wx.request({
      url: 'http://localhost:8080/list-user-workorders',
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`,
      },
      success: (res) => {
        // 请求成功时的回调
        console.log(res.data);
        this.setData({
          workOrders: res.data.work_orders
          .filter(work_order => work_order.RoomID == this.data.roomID)
          .map(work_order => {
            return {
              Date: app.FormatDateToYYYYMMDDHHMMSS(new Date(work_order.CreatedAt)),  // 格式化 CreatedAt
              Description: work_order.Description,
              Status: work_order.Status,
              Type: work_order.Type,
              ID: work_order.ID,
              RoomID: work_order.RoomID
            };
          })
        });
      },
      fail(error) {
        // 请求失败时的回调
        console.log('请求失败', error);
      }
    });
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
    const token = wx.getStorageSync('token');
    wx.request({
      url: 'http://localhost:8080/create-work-order',
      method: 'POST',
      header: {
        'Authorization': `Bearer ${token}`,
      },
      data: {
        "room_id": Number(this.data.roomID),
        "description": this.data.workOrderText.trim()
      },
      success: (res) => {
        // 请求成功时的回调
        this.setData({workOrderText: ""})
        wx.showToast({ title: "工单提交成功", icon: "success" });
        this.onShow()
      },
      fail(error) {
        // 请求失败时的回调
        console.log('请求失败', error);
      }
    });
  },
});
