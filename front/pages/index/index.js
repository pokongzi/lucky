Page({
  data: {
    currentGame: 'ssq',
    generateCount: 5,
    showAdvanced: false,
    ballMode: 'lock', // 'lock' 指定号码, 'exclude' 排除号码
    generatedNumbers: [],
    pendingNumbers: [],
    redBalls: [],
    blueBalls: [],
    user: null,
    hasToken: false,
    loadingLogin: false
  },
  
  onLoad: function() {
    this.initAuthState();
    this.initBallData();
  },

  // 初始化球号数据
  initBallData: function() {
    const redBalls = [];
    const blueBalls = [];
    
    // 根据当前游戏类型初始化球号
    const redMax = this.data.currentGame === 'ssq' ? 33 : 35;
    const blueMax = this.data.currentGame === 'ssq' ? 16 : 12;
    
    for (let i = 1; i <= redMax; i++) {
      redBalls.push({ number: i, status: 'normal' });
    }
    
    for (let i = 1; i <= blueMax; i++) {
      blueBalls.push({ number: i, status: 'normal' });
    }
    
    this.setData({
      redBalls: redBalls,
      blueBalls: blueBalls
    });
  },
  
  switchGame: function(e) {
    const game = e.currentTarget.dataset.game;
    this.setData({
      currentGame: game,
      generatedNumbers: [],
      pendingNumbers: []
    });
    this.initBallData();
  },
  
  // 生成号码
  generateNumbers: function() {
    const count = this.data.generateCount;
    const generatedNumbers = [];
    
    for (let i = 0; i < count; i++) {
      const numbers = this.generateSingleNumber();
      generatedNumbers.push(numbers);
    }
    
    this.setData({
      generatedNumbers: generatedNumbers
    });
  },

  // 生成单组号码
  generateSingleNumber: function() {
    if (this.data.currentGame === 'ssq') {
      // 双色球：6个红球(1-33) + 1个蓝球(1-16)
      const redBalls = this.getRandomNumbers(1, 33, 6, 'red');
      const blueBalls = this.getRandomNumbers(1, 16, 1, 'blue');
      return { redBalls: redBalls, blueBalls: blueBalls };
    } else {
      // 大乐透：5个前区(1-35) + 2个后区(1-12)
      const redBalls = this.getRandomNumbers(1, 35, 5, 'red');
      const blueBalls = this.getRandomNumbers(1, 12, 2, 'blue');
      return { redBalls: redBalls, blueBalls: blueBalls };
    }
  },

  getRandomNumbers: function(min, max, count, type) {
    const numbers = [];
    const lockedNumbers = this.getLockedNumbers(type);
    const excludedNumbers = this.getExcludedNumbers(type);
    
    while (numbers.length < count) {
      const num = Math.floor(Math.random() * (max - min + 1)) + min;
      if (!numbers.includes(num) && !excludedNumbers.includes(num)) {
        numbers.push(num);
      }
    }
    
    // 如果有锁定号码，替换部分随机号码
    if (lockedNumbers.length > 0) {
      const replaceCount = Math.min(lockedNumbers.length, count);
      for (let i = 0; i < replaceCount; i++) {
        numbers[i] = lockedNumbers[i];
      }
    }
    
    return numbers.sort((a, b) => a - b);
  },

  // 获取锁定的号码
  getLockedNumbers: function(type) {
    const balls = type === 'red' ? this.data.redBalls : this.data.blueBalls;
    return balls.filter(ball => ball.status === 'locked').map(ball => ball.number);
  },

  // 获取排除的号码
  getExcludedNumbers: function(type) {
    const balls = type === 'red' ? this.data.redBalls : this.data.blueBalls;
    return balls.filter(ball => ball.status === 'excluded').map(ball => ball.number);
  },

  // 数量控制
  increaseCount: function() {
    if (this.data.generateCount < 20) {
      this.setData({
        generateCount: this.data.generateCount + 1
      });
    }
  },

  decreaseCount: function() {
    if (this.data.generateCount > 1) {
      this.setData({
        generateCount: this.data.generateCount - 1
      });
    }
  },

  // 切换高级设置
  toggleAdvanced: function() {
    this.setData({
      showAdvanced: !this.data.showAdvanced
    });
  },

  // 切换球号模式
  switchBallMode: function(e) {
    const mode = e.currentTarget.dataset.mode;
    this.setData({
      ballMode: mode
    });
  },

  // 切换球号状态
  toggleBallStatus: function(e) {
    const number = parseInt(e.currentTarget.dataset.number);
    const type = e.currentTarget.dataset.type;
    const balls = type === 'red' ? 'redBalls' : 'blueBalls';
    const ballList = this.data[balls];
    
    const ballIndex = ballList.findIndex(ball => ball.number === number);
    if (ballIndex !== -1) {
      const currentStatus = ballList[ballIndex].status;
      let newStatus = 'normal';
      
      // 根据当前模式决定状态切换
      if (this.data.ballMode === 'lock') {
        // 指定号码模式：normal -> locked -> normal
        if (currentStatus === 'normal') {
          newStatus = 'locked';
        } else if (currentStatus === 'locked') {
          newStatus = 'normal';
        } else if (currentStatus === 'excluded') {
          newStatus = 'locked';
        }
      } else {
        // 排除模式：normal -> excluded -> normal
        if (currentStatus === 'normal') {
          newStatus = 'excluded';
        } else if (currentStatus === 'excluded') {
          newStatus = 'normal';
        } else if (currentStatus === 'locked') {
          newStatus = 'excluded';
        }
      }
      
      ballList[ballIndex].status = newStatus;
      
      this.setData({
        [balls]: ballList
      });
    }
  },

  // 清空所有设置
  clearAllSettings: function() {
    const redBalls = this.data.redBalls.map(ball => ({ ...ball, status: 'normal' }));
    const blueBalls = this.data.blueBalls.map(ball => ({ ...ball, status: 'normal' }));
    
    this.setData({
      redBalls: redBalls,
      blueBalls: blueBalls
    });
  },

  // 选用号码
  selectNumber: function(e) {
    const index = e.currentTarget.dataset.index;
    const selectedNumber = this.data.generatedNumbers[index];
    const pendingNumbers = [...this.data.pendingNumbers, selectedNumber];
    
    this.setData({
      pendingNumbers: pendingNumbers
    });
    
    wx.showToast({
      title: '已添加到待选清单',
      icon: 'success'
    });
  },

  // 从待选清单移除
  removeFromPending: function(e) {
    const index = e.currentTarget.dataset.index;
    const pendingNumbers = [...this.data.pendingNumbers];
    pendingNumbers.splice(index, 1);
    
    this.setData({
      pendingNumbers: pendingNumbers
    });
  },

  // 复制待选清单
  copyPendingList: function() {
    if (this.data.pendingNumbers.length === 0) {
      wx.showToast({
        title: '待选清单为空',
        icon: 'none'
      });
      return;
    }
    
    let copyText = '';
    this.data.pendingNumbers.forEach((item, index) => {
      const redStr = item.redBalls.map(num => num.toString().padStart(2, '0')).join(' ');
      const blueStr = item.blueBalls.map(num => num.toString().padStart(2, '0')).join(' ');
      copyText += `第${index + 1}注: ${redStr} | ${blueStr}\n`;
    });
    
    wx.setClipboardData({
      data: copyText,
      success: () => {
        wx.showToast({
          title: '已复制到剪贴板',
          icon: 'success'
        });
      }
    });
  },

  // 收藏待选清单
  collectPendingList: function() {
    if (this.data.pendingNumbers.length === 0) {
      wx.showToast({
        title: '待选清单为空',
        icon: 'none'
      });
      return;
    }
    
    if (!this.data.hasToken) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    // 循环调用现有的保存接口
    const gameCode = this.data.currentGame === 'ssq' ? 'ssq' : 'dlt';
    let successCount = 0;
    let totalCount = this.data.pendingNumbers.length;
    
    this.data.pendingNumbers.forEach((item, index) => {
      wx.request({
        url: this.getBaseURL() + '/api/numbers/save',
        method: 'POST',
        header: {
          'X-User-ID': this.data.user.id.toString(),
          'Content-Type': 'application/json'
        },
        data: {
          gameCode: gameCode,
          redBalls: item.redBalls,
          blueBalls: item.blueBalls,
          nickname: `收藏号码_${index + 1}`,
          source: 'collect'
        },
        success: (response) => {
          const { code: bizCode } = response.data || {};
          if (response.statusCode === 200 && bizCode === 200) {
            successCount++;
            if (successCount === totalCount) {
              wx.showToast({
                title: `成功收藏${successCount}组号码`,
                icon: 'success'
              });
            }
          } else {
            wx.showToast({
              title: `第${index + 1}组号码收藏失败`,
              icon: 'none'
            });
          }
        },
        fail: () => {
          wx.showToast({
            title: `第${index + 1}组号码网络请求失败`,
            icon: 'none'
          });
        }
      });
    });
  },

  // 清空待选清单
  clearPendingList: function() {
    wx.showModal({
      title: '确认清空',
      content: '确定要清空所有待选号码吗？',
      success: (res) => {
        if (res.confirm) {
          this.setData({
            pendingNumbers: []
          });
        }
      }
    });
  },

  // ========== 登录相关 ==========
  initAuthState: function() {
    try {
      const token = wx.getStorageSync('token');
      const user = wx.getStorageSync('user');
      console.log('初始化登录状态:', { token: !!token, user: !!user });
      if (token && user) {
        this.setData({ hasToken: true, user: user });
        console.log('用户已登录:', user);
      } else {
        this.setData({ hasToken: false, user: null });
        console.log('用户未登录');
      }
    } catch (e) {
      console.error('初始化登录状态失败:', e);
      this.setData({ hasToken: false, user: null });
    }
  },

  handleLoginTap: function() {
    if (this.data.loadingLogin) return;
    this.setData({ loadingLogin: true });

    const doFail = (msg) => {
      wx.showToast({ title: msg || '登录失败', icon: 'none' });
      this.setData({ loadingLogin: false });
    };

    wx.login({
      success: (res) => {
        const code = res.code;
        if (!code) return doFail('获取code失败');

        // 新的登录方式：只使用 code 进行登录，不强制获取用户信息
        wx.request({
          url: this.getBaseURL() + '/api/auth/wxlogin',
          method: 'POST',
          data: { 
            code: code
          },
          success: (response) => {
            const { code: bizCode, data, message } = response.data || {};
            if (response.statusCode !== 200 || bizCode !== 200 || !data || !data.token) {
              return doFail(message || '登录失败');
            }

            try {
              wx.setStorageSync('token', data.token);
              wx.setStorageSync('user', data.user);
              wx.setStorageSync('tokenExpiresAt', data.expiresAt);
            } catch (e) {
              console.error('存储登录信息失败:', e);
            }

            this.setData({ 
              hasToken: true, 
              user: data.user, 
              loadingLogin: false 
            });
            console.log('登录成功，更新状态:', { hasToken: true, user: data.user });
            wx.showToast({ title: '登录成功', icon: 'success' });
          },
          fail: () => doFail('登录请求失败')
        });
      },
      fail: () => doFail('wx.login失败')
    });
  },

  // 新增：用户信息更新方法（可选）
  onChooseAvatar: function(e) {
    const { avatarUrl } = e.detail;
    this.updateUserInfo({ avatarUrl });
  },

  onNicknameChange: function(e) {
    const nickname = e.detail.value;
    this.updateUserInfo({ nickname });
  },

  updateUserInfo: function(userInfo) {
    if (!this.data.hasToken) {
      wx.showToast({ title: '请先登录', icon: 'none' });
      return;
    }

    wx.request({
      url: this.getBaseURL() + '/api/user/update',
      method: 'POST',
      header: {
        'Authorization': 'Bearer ' + wx.getStorageSync('token')
      },
      data: userInfo,
      success: (response) => {
        const { code: bizCode, data } = response.data || {};
        if (response.statusCode === 200 && bizCode === 200) {
          // 更新本地存储的用户信息
          const currentUser = this.data.user || {};
          const updatedUser = { ...currentUser, ...userInfo };
          wx.setStorageSync('user', updatedUser);
          this.setData({ user: updatedUser });
          wx.showToast({ title: '更新成功', icon: 'success' });
        }
      }
    });
  },

  getBaseURL: function() {
    // 使用全局配置
    const app = getApp();
    return app.getBaseURL();
  },

  // 登出功能
  handleLogout: function() {
    wx.showModal({
      title: '确认退出',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          try {
            wx.removeStorageSync('token');
            wx.removeStorageSync('user');
            wx.removeStorageSync('tokenExpiresAt');
          } catch (e) {
            console.error('清除登录信息失败:', e);
          }
          
          this.setData({
            hasToken: false,
            user: null,
            generatedNumbers: [],
            pendingNumbers: []
          });
          
          wx.showToast({
            title: '已退出登录',
            icon: 'success'
          });
        }
      }
    });
  }
});