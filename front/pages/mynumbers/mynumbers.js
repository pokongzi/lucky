Page({
  data: {
    myNumbers: []
  },
  
  onLoad: function() {
    this.loadMyNumbers();
  },
  
  loadMyNumbers: function() {
    // 从本地存储加载保存的号码
    const savedNumbers = wx.getStorageSync('myNumbers') || [];
    this.setData({
      myNumbers: savedNumbers
    });
  },
  
  deleteNumber: function(e) {
    const id = e.currentTarget.dataset.id;
    const myNumbers = this.data.myNumbers.filter(item => item.id !== id);
    
    this.setData({
      myNumbers: myNumbers
    });
    
    // 保存到本地存储
    wx.setStorageSync('myNumbers', myNumbers);
    
    wx.showToast({
      title: '删除成功',
      icon: 'success'
    });
  }
});
