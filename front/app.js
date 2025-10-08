App({
  globalData: {
    // API配置
    apiConfig: {
      // 开发环境
      development: {
        baseUrl: 'http://47.121.26.190:8080',
        domain: '47.121.26.190:8080'
      },
      // 生产环境
      production: {
        baseUrl: 'https://your-prod-domain.com',
        domain: 'your-prod-domain.com'
      }
    },
    // 当前环境，可根据需要切换
    currentEnv: 'development'
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
    console.log('小程序启动');
    // 可以在这里进行一些初始化操作
  },

  onShow: function() {
    console.log('小程序显示');
  },

  onHide: function() {
    console.log('小程序隐藏');
  }
})