Component({
  properties: {
    // 接收外部传递的数据
    imageSrc: {
      type: String,
      value: '/images/houses/house_1.png' // 默认图片
    },
    price: {
      type: String,
      value: '6300元/月' // 默认价格
    },
    action1: {
      type: String,
      value: '随时看'
    },
    action2: {
      type: String,
      value: '房源新'
    },
    title: {
      type: String,
      value: '贝丽花园 3室2厅 南'
    },
    info: {
      type: String,
      value: '3房2厅/110㎡/龙华 高山街道'
    }
  },
  data: {},
  methods: {}
});