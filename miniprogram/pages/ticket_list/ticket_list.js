Page({
  data: {
    tickets: [
      { 
        id: "T001", 
        title: "空调故障", 
        created_at: "2025-03-30", 
        description: "房间 101 的空调无法正常制冷", 
        user_id: "U123456", 
        room_name: "101 号房",
        resolved: false 
      },
      { 
        id: "T002", 
        title: "水龙头漏水", 
        created_at: "2025-03-29", 
        description: "厨房水龙头持续滴水", 
        user_id: "U654321", 
        room_name: "202 号房",
        resolved: true 
      },
      { 
        id: "T003", 
        title: "门锁损坏", 
        created_at: "2025-03-28", 
        description: "房门锁芯松动，钥匙难以插入", 
        user_id: "U789123", 
        room_name: "303 号房",
        resolved: false 
      }
    ]
  },

  markAsResolved(e) {
    const id = e.currentTarget.dataset.id;
    const updatedTickets = this.data.tickets.map(ticket => {
      if (ticket.id === id) {
        return { ...ticket, resolved: true };
      }
      return ticket;
    });
    this.setData({ tickets: updatedTickets });
  }
});
