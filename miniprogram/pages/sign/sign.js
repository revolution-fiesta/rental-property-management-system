const app = getApp()
Page({
  data: {
    canvas: '',
    ctx: '',
    pr:0,
    width: 0,
    height: 0,
    first:true,
  },
  start(e) {
    if(this.data.first){
      this.clearClick();
      this.setData({first:false})
    }
    this.data.ctx.beginPath(); // 开始创建一个路径，如果不调用该方法，最后无法清除画布
    this.data.ctx.moveTo(e.changedTouches[0].x, e.changedTouches[0].y) // 把路径移动到画布中的指定点，不创建线条。用 stroke 方法来画线条
  },
  move(e) {
    this.data.ctx.lineTo(e.changedTouches[0].x, e.changedTouches[0].y) // 增加一个新点，然后创建一条从上次指定点到目标点的线。用 stroke 方法来画线条
    this.data.ctx.stroke()
  },
  error: function (e) {
    console.log("画布触摸错误" + e);
  },
  /**
  * 生命周期函数--监听页面加载
  */
  onLoad: function () {
    this.getSystemInfo()
    this.createCanvas()
  },
  /**
   * 初始化
   */
  createCanvas() {
    const pr = this.data.pr; // 像素比
    const query = wx.createSelectorQuery();
    query.select('#canvas').fields({ node: true, size: true }).exec((res) => {
      const canvas = res[0].node;
      const ctx = canvas.getContext('2d');
      canvas.width = this.data.width*pr; // 画布宽度
      canvas.height = this.data.height*pr; // 画布高度
      ctx.scale(pr,pr); // 缩放比
      ctx.lineGap = 'round';
      ctx.lineJoin = 'round';
      ctx.lineWidth = 4; // 字体粗细
      ctx.font = '40px Arial'; // 字体大小，
      ctx.fillStyle = '#ecf0ef'; // 填充颜色
      ctx.fillText('签名区', res[0].width / 2 - 60, res[0].height / 2)
      this.setData({ ctx, canvas })
    })
  },
  // 获取系统信息
  getSystemInfo() {
    let that = this;
    wx.getSystemInfo({
      success(res) {
        that.setData({
          pr:res.pixelRatio,
          width: res.windowWidth,
          height: res.windowHeight - 75,
        })
      }
    })
  },
  clearClick: function () {
    //清除画布
    this.data.ctx.clearRect(0, 0, this.data.width, this.data.height);
  },

  saveClick: function () {
    const token = wx.getStorageSync('token')
    if(!token) {
      wx.showToast({
        title: '请重新登录',
        icon: 'error'
      })
      return
    }
    wx.canvasToTempFilePath({
      x: 0,
      y: 0,
      width: this.data.width,
      height: this.data.height,
      destWidth: this.data.width * this.data.pr,
      destHeight: this.data.height * this.data.pr,
      canvasId: 'canvas',
      canvas: this.data.canvas,
      fileType: 'png',
      success: (res) => {
        const tempFilePath = res.tempFilePath; // 获取临时文件路径
  
        // 读取文件内容并转为 Base64
        const fs = wx.getFileSystemManager();
        fs.readFile({
          filePath: tempFilePath,
          encoding: 'base64',
          success: (readRes) => {
            const base64Data = readRes.data; // 获取 Base64 数据
  
            // 发送请求，上传图片数据
            wx.request({
              url: 'http://localhost:8080/upload-signature', // 你的服务器地址
              method: 'POST',
              header: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
              },
              data: {
                image_data: base64Data, // 图片的 Base64 编码
              },
              success(res) {
                wx.showToast({
                  title: '上传成功',
                  icon: 'success',
                });
                console.log(res.data)
              },
              fail(error) {
                wx.showToast({
                  title: '上传失败',
                  icon: 'none',
                });
              },
            });
          },
          fail: (error) => {
            console.error('读取文件失败:', error);
            wx.showToast({
              title: '读取文件失败',
              icon: 'none',
            });
          },
        });
      },
      fail: (err) => {
        console.error('保存失败:', err);
        wx.showToast({
          title: '保存失败',
          icon: 'none',
        });
      },
    });
  }
  
})


  // 这个版本是保存在本地
  // saveClick: function () {
  //   wx.canvasToTempFilePath({
  //     x: 0,
  //     y: 0,
  //     width: this.data.width,
  //     height: this.data.height,
  //     destWidth:this.data.width*this.data.pr,
  //     destHeight: this.data.height*this.data.pr,
  //     canvasId: 'canvas',
  //     canvas: this.data.canvas,
  //     fileType: 'png',
  //     success(res) {
  //        保存到本地
  //       wx.saveImageToPhotosAlbum({
  //         filePath: res.tempFilePath,
  //         success(res) {
  //           wx.showToast({
  //             title: '保存成功',
  //             icon: 'success'
  //           })
  //         }
  //       })
  //     }
  //   })
  // }