Component({
  properties: {
    bill: Object, // 账单数据
  },
  methods: {
    goToDetail() {
      wx.navigateTo({
        url: `/pages/billDetail/billDetail?id=${this.properties.bill.id}`
      });
    }
  }
});
