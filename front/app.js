App({
  globalData: {
    // API配置
    apiConfig: {
      // 开发环境
      development: {
        baseUrl: 'https://luckygirls.cn',
        domain: 'luckygirls.cn'
      },
      // 生产环境
      production: {
        baseUrl: 'https://luckygirls.cn',
        domain: 'luckygirls.cn'
      }
    },
    // 当前环境，可根据需要切换
    currentEnv: 'development',
    // 登录状态
    isLoggedIn: false,
    user: null
  },

  // 获取API基础URL
  getBaseURL: function() {
    const env = this.globalData.currentEnv;
    return this.globalData.apiConfig[env].baseUrl;
  },

  // 获取当前环境配置
  getCurrentConfig: function() {
    const env = this.globalData.currentEnv;
    return this.globalData.apiConfig[env];
  },

  onLaunch: function() {
    // 静默登录：尝试自动登录
    this.silentLogin();
  },

  // 静默登录方法
  silentLogin: function() {
    
    // 1. 先检查本地是否有有效的登录状态
    try {
      const token = wx.getStorageSync('token');
      const tokenExpiresAt = wx.getStorageSync('tokenExpiresAt');
      const user = wx.getStorageSync('user');
      
      if (token && user && tokenExpiresAt) {
        const expiresTime = new Date(tokenExpiresAt).getTime();
        const now = Date.now();
        
        // 如果token未过期，直接使用本地登录状态
        if (expiresTime > now) {
          this.globalData.isLoggedIn = true;
          this.globalData.user = user;
          return;
        }
      }
    } catch (e) {
      // 检查本地登录状态失败
    }

    // 2. 本地无有效登录状态，尝试静默登录
    wx.login({
      success: (res) => {
        if (!res.code) {
          console.error('静默登录失败: 无法获取code');
          return;
        }

        // 调用后端登录接口（只使用code，不获取用户信息）
        wx.request({
          url: this.getBaseURL() + '/api/auth/wxlogin',
          method: 'POST',
          data: { code: res.code },
          success: (response) => {
            const { code: bizCode, data } = response.data || {};
            
            if (response.statusCode === 200 && bizCode === 200 && data && data.token) {
              // 静默登录成功
              try {
                wx.setStorageSync('token', data.token);
                wx.setStorageSync('user', data.user);
                wx.setStorageSync('tokenExpiresAt', data.expiresAt);
                
                this.globalData.isLoggedIn = true;
                this.globalData.user = data.user;
                
                // 静默登录成功
              } catch (e) {
                // 保存静默登录信息失败
              }
            } else {
              // 静默登录失败，等待用户触发授权
            }
          },
          fail: (err) => {
            console.error('静默登录请求失败:', err);
          }
        });
      },
      fail: (err) => {
        console.error('wx.login 失败:', err);
      }
    });
  },

  onShow: function() {
    // 小程序显示
  },

  onHide: function() {
    // 小程序隐藏
  }
})