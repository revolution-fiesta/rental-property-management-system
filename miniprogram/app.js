// app.js
App({
  FormatDateToYYYYMMDDHHMMSS(date) {
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
  },

  ConvertRoomType(roomType) {
    // 定义映射关系
    const mapping = {
      "一室": "b1",
      "一室一厅": "b1l1",
      "两室一厅": "b2l1",
    };
    // 返回对应的中文描述，找不到则返回 "未知房型"
    return mapping[roomType];
  },
  
  ConvertRoomTypeReverse(roomTypeCode) {
    // 定义反向映射关系
    const reverseMapping = {
      "b1": "一室",
      "b1l1": "一室一厅",
      "b2l1": "两室一厅",
    };
    // 返回对应的房型，找不到则返回 "未知房型"
    return reverseMapping[roomTypeCode] || "未知房型";
  }
})

