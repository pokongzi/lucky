// 获取全局app实例
const app = getApp();

Page({
  data: {
    currentGame: 'ssq',
    myNumbers: [],
    filteredNumbers: [],
    hasLoaded: false, // 标记是否已经加载过数据
    showWinningModal: false,
    winningData: {}
  },
  
  onLoad: function() {
    // onLoad时不加载数据，等待onShow
  },

  onShow: function() {
    // 每次页面显示时都尝试加载数据
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
    // 使用按需授权：需要登录时自动触发
    const authUtil = require('../../utils/auth.js');
    
    authUtil.requireLogin((loginInfo) => {
      // 登录成功后加载号码
      this.fetchNumbersData(loginInfo.user.id);
    }, {
      title: '需要登录',
      content: '查看我的号码需要登录，是否立即登录？',
      showModal: true
    }).catch((error) => {
      if (error.code !== 'USER_CANCEL') {
        wx.showToast({
          title: error.message || '登录失败',
          icon: 'none'
        });
      }
    });
  },

  // 获取号码数据（独立方法）
  fetchNumbersData: function(userId) {
    wx.showLoading({ title: '加载中' });
    
    wx.request({
      url: app.getBaseURL() + '/api/numbers/my',
      method: 'GET',
      data: {
        gameCode: this.data.currentGame
      },
      header: {
        'X-User-ID': userId.toString()
      },
      success: (res) => {
        wx.hideLoading();
        if (res.data && res.data.code === 200 && res.data.data) {
          this.setData({
            myNumbers: res.data.data.list || [],
            hasLoaded: true
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
        wx.hideLoading();
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        });
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
  
  // 核对中奖
  checkWinning: function(e) {
    const numberId = e.currentTarget.dataset.id;
    const app = getApp();
    const baseUrl = app.getBaseURL();
    
    wx.showLoading({
      title: '核对中奖中...',
    });
    
    wx.request({
      url: `${baseUrl}/api/numbers/${numberId}/check`,
      method: 'GET',
      header: {
        'X-User-ID': '1'  // 实际使用时需要从登录状态获取真实用户ID
      },
      success: (res) => {
        
        if (res.data && res.data.code === 200 && res.data.data) {
          this.showWinningModal(res.data.data);
        } else {
          wx.showToast({
            title: res.data.message || '核对失败',
            icon: 'none'
          });
        }
      },
      fail: (err) => {
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

  // 显示中奖结果弹窗
  showWinningModal: function(data) {
    const { userNumber, matches, totalMatches, totalPrize } = data;
    
    // 处理中奖数据
    const formattedMatches = matches.map(match => {
      return {
        ...match,
        prizeText: this.formatPrize(match.prizeAmount),
        levelClass: this.getWinLevelClass(match.winLevel)
      };
    });
    
    this.setData({
      showWinningModal: true,
      winningData: {
        totalMatches: totalMatches,
        totalPrizeText: this.formatPrize(totalPrize),
        matches: formattedMatches
      }
    });
  },

  // 关闭中奖结果弹窗
  closeWinningModal: function() {
    this.setData({
      showWinningModal: false
    });
  },

  // 获取中奖等级样式类名
  getWinLevelClass: function(winLevel) {
    if (winLevel.includes('一等奖') || winLevel.includes('特等奖')) {
      return 'level-1';
    } else if (winLevel.includes('二等奖')) {
      return 'level-2';
    } else if (winLevel.includes('三等奖')) {
      return 'level-3';
    } else if (winLevel.includes('四等奖')) {
      return 'level-4';
    } else if (winLevel.includes('五等奖')) {
      return 'level-5';
    } else {
      return 'level-6';
    }
  },

  // 格式化奖金显示
  formatPrize: function(prizeInCents) {
    const yuan = prizeInCents / 100;
    if (yuan >= 10000) {
      return `${(yuan / 10000).toFixed(1)}万元`;
    } else {
      return `${yuan}元`;
    }
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

