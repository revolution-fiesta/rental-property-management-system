Page({
  data: {
    houseName: '',
    rentPeriod: '',
    amount: 0,
    billDate: ''
  },

  onLoad(options) {
    this.setData({
      houseName: options.houseName || '未知房源',
      rentPeriod: options.rentPeriod || '无',
      amount: options.amount || '0.00',
      billDate: options.billDate || '未知日期'
    });
  },

  handlePayment() {
    wx.showToast({
      title: '支付中',
      icon: 'none'
    });
    setTimeout(() => {
      wx.showToast({
        title: '支付成功',
        icon: 'none'
      });
      setTimeout(()=>{
        wx.navigateTo({
          url: '/pages/details/details?houseName=华发山庄&rent=5380&houseType=一室一厅&area=55&location=北京市朝阳区&description=精装修，拎包入住&houseImage=/images/houses/house_1.png'
        });
      }, 500)
    }, 300);
  }
});
