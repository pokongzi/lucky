Page({
  data: {
    historyList: []
  },
  
  onLoad: function() {
    this.loadHistory();
  },
  
  loadHistory: function() {
    // 模拟历史数据
    const mockData = [
      {
        id: 1,
        period: "2024001",
        redNumbers: [1, 5, 12, 18, 25, 33],
        blueNumber: 8,
        date: "2024-01-01"
      },
      {
        id: 2,
        period: "2024002", 
        redNumbers: [3, 7, 14, 20, 28, 31],
        blueNumber: 12,
        date: "2024-01-03"
      }
    ];
    
    this.setData({
      historyList: mockData
    });
  }
});
