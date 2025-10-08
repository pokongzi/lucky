// 获取全局app实例
const app = getApp();

Page({
  data: {
    currentGame: 'ssq',
    myNumbers: [],
    filteredNumbers: []
  },
  
  onLoad: function() {
    this.loadMyNumbers();
  },
  
  // 切换游戏
  switchGame: function(e) {
    const game = e.currentTarget.dataset.game;
    this.setData({
      currentGame: game
    });
    // 切换游戏后重新请求API获取对应游戏的号码数据
    this.loadMyNumbers();
  },
  

  
  // 加载我的号码
  loadMyNumbers: function() {
    wx.showLoading({
      title: '加载中',
    });
    
    // 调用后端API获取我的号码列表
    wx.request({
      url: app.getBaseURL() + '/api/numbers/my',
      method: 'GET',
      data: {
        gameCode: this.data.currentGame
      },
      header: {
        'X-User-ID': '1'  // 实际使用时需要从登录状态获取真实用户ID
      },
      success: (res) => {
        if (res.data && res.data.code === 200 && res.data.data) {
          this.setData({
            myNumbers: res.data.data.list || []
          });
          this.filterNumbers();
        } else {
          wx.showToast({
            title: '加载失败',
            icon: 'none'
          });
        }
      },
      fail: (err) => {
        console.error('获取号码列表失败:', err);
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        });
      },
      complete: () => {
        wx.hideLoading();
      }
    });
  },
  
  // 格式化时间，只保留年月日
  formatDate: function(dateString) {
    if (!dateString) return '';
    // 提取ISO格式中的年月日部分 (YYYY-MM-DD)
    const date = new Date(dateString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  },
  
  // 筛选号码
  filterNumbers: function() {
    const filteredNumbers = this.data.myNumbers.filter(item => {
      return item.game && item.game.game_code === this.data.currentGame;
    }).map(item => {
      // 格式化创建时间，只保留年月日
      if (item.created_at) {
        item.created_at = this.formatDate(item.created_at);
      }
      return item;
    });
    
    this.setData({
      filteredNumbers: filteredNumbers
    });
  },
  
  // 查询历史
  checkHistory: function(e) {
    const id = e.currentTarget.dataset.id;
    wx.showToast({
      title: '跳转到历史分析页',
      icon: 'none'
    });
    // TODO: 实现跳转到单号深度分析页
  },
  
  // 删除号码
  deleteNumber: function(e) {
    const id = e.currentTarget.dataset.id;
    
    wx.showModal({
      title: '确认删除',
      content: '确定要删除这组号码吗？',
      success: (res) => {
        if (res.confirm) {
          const myNumbers = this.data.myNumbers.filter(item => item.id !== id);
          this.setData({
            myNumbers: myNumbers
          });
          this.filterNumbers();
          
          wx.showToast({
            title: '删除成功',
            icon: 'success'
          });
        }
      }
    });
  },
  

  
  // 去生成号码
  gotoGenerate: function() {
    wx.switchTab({
      url: '/pages/index/index'
    });
  }
});

