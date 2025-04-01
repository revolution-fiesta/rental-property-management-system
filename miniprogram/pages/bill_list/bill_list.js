Page({
  data: {
    bills: []
  },

  onShow() {
    this.loadBilling()
  },
  
  loadBilling() {
    // 从本地存储获取 token
    const token = wx.getStorageSync('token');
    // 如果获取到 token，发起 GET 请求
    if (token) {
      wx.request({
        url: 'http://localhost:8080/list-billings', // 替换为你实际的 API 地址
        method: 'GET',
        header: {
          'Authorization': `Bearer ${token}`,  // 设置 Authorization 头
        },
        success: (res) => {
          // 请求成功时的回调
          console.log('请求成功', res.data);
          this.setData({
            bills: res.data.billings.map(bill => {
              return {
                Price: bill.Price,
                Paid: bill.Paid,
                FormatType: billingTypeToChinese(bill.Type),  // 转换 BillingType
                Date: formatDateToYYYYMMDDHHMMSS(new Date(bill.CreatedAt)),  // 格式化 CreatedAt
                ID: bill.ID,
                Name: bill.Name
              };
            })
          });
        },
        fail(error) {
          // 请求失败时的回调
          console.log('请求失败', error);
        }
      });
    } else {
      console.log('未找到 token');
    }
  },
  naviToOrder(e) {
    const item_idx = e.currentTarget.dataset.index
    const bill_obj = this.data.bills[item_idx]
    wx.navigateTo({
      url: `/pages/order/order?orderType=外卖&amount=${bill_obj.Price}&billDate=${bill_obj.Date}&billID=${bill_obj.ID}&houseName=${bill_obj.Name}&type=${bill_obj.FormatType}`
    });
  },
});

function billingTypeToChinese(type) {
  const billingTypes = {
    'rent-room': '首次签约',
    'monthly-pay': '月租',
    'terminate-lease': '退租'
  };

  return billingTypes[type] || '未知类型'; // 默认返回 '未知类型' 如果类型没有匹配
}

function formatDateToYYYYMMDDHHMMSS(date) {
  if (!(date instanceof Date)) {
    throw new Error('Invalid date object');
  }

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0'); // 月份从0开始，所以要加1
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}
