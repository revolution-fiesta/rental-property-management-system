Page({
  data: {
    startDate: '', // 当前日期作为租用开始日期
    endDate: '',  // 结束日期根据租期计算
    leaseTerm: null, // 租期（月）
    rentAmount: null, // 每月租金
    depositAmount: null, // 押金
    room_id: null
  },

  onLoad(opts) {
    const startDate = this.getCurrentDate();
    const endDate = this.calculateEndDate(startDate, this.data.leaseTerm);
    
    this.setData({
      startDate: startDate,
      endDate: endDate,
      leaseTerm: opts.num_terms,
      rentAmount: opts.rent,
      // TODO: 租金写死了是两个月
      depositAmount: opts.rent * 2,
      room_id: opts.room_id
    });
  },

  getCurrentDate() {
    const now = new Date();
    const year = now.getFullYear();
    const month = (now.getMonth() + 1).toString().padStart(2, '0');
    const day = now.getDate().toString().padStart(2, '0');
    return `${year}-${month}-${day}`;
  },

  calculateEndDate(startDate, leaseTerm) {
    const date = new Date(startDate);
    date.setMonth(date.getMonth() + leaseTerm);
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    return `${year}-${month}-${day}`;
  },

  onSubmit: function () {
    wx.navigateTo({
      url: `/pages/sign/sign?room_id=${this.data.room_id}&num_terms=${this.data.leaseTerm}`,
    })
  }
});