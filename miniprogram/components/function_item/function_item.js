Component({
  properties: {
    // 接收父组件传递的属性
    icon: String,  // 图标路径
    text: String,  // 显示的文本
    targetPage: String  // 点击后跳转的页面
  },

  methods: {
    // 点击事件，跳转到相应页面
    onTap: function () {
      const targetPage = this.data.targetPage;
      wx.navigateTo({
        url: targetPage,
      });
    }
  }
});