Page({
  data: {
    houseName: '',
    rentPeriod: '',
    amount: 0,
    billDate: '',
    billID: '',
    type: '',
  },

  onLoad(options) {
    this.setData({
      houseName: options.houseName || '未知房源',
      rentPeriod: options.rentPeriod || '无',
      amount: options.amount || '0.00',
      billDate: options.billDate || '未知日期',
      billID: options.billID || '未知编号',
      type: options.type || '未知类型'
    });
  },
  handlePayment() {
   
    const token = wx.getStorageSync('token'); 
    wx.showToast({
      title: '支付中',
      icon: 'none'
    });
    wx.request({
      url: 'http://localhost:8080/pay-bill',  
      method: 'POST',
      header: {
        'Authorization': `Bearer ${token}`,  
        'Content-Type': 'application/json'  
      },
      data: {
        billing_id: Number( this.data.billID  )
      },
      success(res) {
        if (res.statusCode === 200) {
          setTimeout(()=>{
            wx.showToast({
              title: '支付成功',
            })
            setTimeout(()=>{
              wx.navigateBack()
            },500)
            
          }, 500)
        } else {
          console.log('支付失败', res.data);
        }
      },
      fail(error) {
        wx.showToast({
          title: '支付失败',
        })
      }
    });

  },

});
