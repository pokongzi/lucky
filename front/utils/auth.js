/**
 * 登录授权工具
 * 提供静默登录和按需授权功能
 */

/**
 * 检查登录状态
 * @returns {Promise} 返回用户信息或拒绝
 */
function checkLoginStatus() {
  return new Promise((resolve, reject) => {
    try {
      const token = wx.getStorageSync('token');
      const tokenExpiresAt = wx.getStorageSync('tokenExpiresAt');
      const user = wx.getStorageSync('user');

      if (token && user && tokenExpiresAt) {
        const expiresTime = new Date(tokenExpiresAt).getTime();
        const now = Date.now();

        if (expiresTime > now) {
          resolve({ token, user });
        } else {
          reject({ code: 'TOKEN_EXPIRED', message: 'Token已过期' });
        }
      } else {
        reject({ code: 'NOT_LOGGED_IN', message: '未登录' });
      }
    } catch (e) {
      reject({ code: 'CHECK_ERROR', message: '检查登录状态失败', error: e });
    }
  });
}

/**
 * 触发登录授权
 * 当用户需要执行需要登录的操作时调用
 * @param {Object} options 配置选项
 * @param {String} options.title 授权提示标题
 * @param {String} options.content 授权提示内容
 * @returns {Promise} 返回登录结果
 */
function triggerLogin(options = {}) {
  const {
    title = '需要登录',
    content = '该功能需要登录后使用，是否立即登录？',
    showModal = true // 是否显示确认对话框
  } = options;

  return new Promise((resolve, reject) => {
    // 如果需要显示确认对话框
    if (showModal) {
      wx.showModal({
        title: title,
        content: content,
        confirmText: '立即登录',
        cancelText: '稍后再说',
        success: (modalRes) => {
          if (modalRes.confirm) {
            performLogin().then(resolve).catch(reject);
          } else {
            reject({ code: 'USER_CANCEL', message: '用户取消登录' });
          }
        },
        fail: () => {
          reject({ code: 'MODAL_ERROR', message: '显示登录提示失败' });
        }
      });
    } else {
      // 直接执行登录
      performLogin().then(resolve).catch(reject);
    }
  });
}

/**
 * 执行登录流程
 * @returns {Promise}
 */
function performLogin() {
  return new Promise((resolve, reject) => {
    wx.showLoading({ title: '登录中...' });

    wx.login({
      success: (res) => {
        if (!res.code) {
          wx.hideLoading();
          reject({ code: 'NO_CODE', message: '获取登录凭证失败' });
          return;
        }

        const app = getApp();
        wx.request({
          url: app.getBaseURL() + '/api/auth/wxlogin',
          method: 'POST',
          data: { code: res.code },
          success: (response) => {
            wx.hideLoading();
            
            const { code: bizCode, data, message } = response.data || {};
            
            if (response.statusCode === 200 && bizCode === 200 && data && data.token) {
              try {
                // 保存登录信息
                wx.setStorageSync('token', data.token);
                wx.setStorageSync('user', data.user);
                wx.setStorageSync('tokenExpiresAt', data.expiresAt);

                // 更新全局状态
                app.globalData.isLoggedIn = true;
                app.globalData.user = data.user;

                wx.showToast({ title: '登录成功', icon: 'success', duration: 1500 });
                resolve(data);
              } catch (e) {
                reject({ code: 'SAVE_ERROR', message: '保存登录信息失败', error: e });
              }
            } else {
              reject({ code: 'LOGIN_FAILED', message: message || '登录失败' });
            }
          },
          fail: (err) => {
            wx.hideLoading();
            reject({ code: 'REQUEST_ERROR', message: '登录请求失败', error: err });
          }
        });
      },
      fail: (err) => {
        wx.hideLoading();
        reject({ code: 'WX_LOGIN_ERROR', message: 'wx.login失败', error: err });
      }
    });
  });
}

/**
 * 需要登录时执行的包装函数
 * 自动检查登录状态，未登录时触发授权
 * @param {Function} callback 需要执行的回调函数
 * @param {Object} options 登录选项
 * @returns {Promise}
 */
function requireLogin(callback, options = {}) {
  return checkLoginStatus()
    .then((loginInfo) => {
      // 已登录，直接执行回调
      return callback(loginInfo);
    })
    .catch((error) => {
      // 未登录或token过期，触发登录
      return triggerLogin(options).then((loginResult) => {
        // 登录成功后执行回调
        return callback({ token: loginResult.token, user: loginResult.user });
      });
    });
}

/**
 * 退出登录
 * @returns {Promise}
 */
function logout() {
  return new Promise((resolve) => {
    try {
      wx.removeStorageSync('token');
      wx.removeStorageSync('user');
      wx.removeStorageSync('tokenExpiresAt');

      const app = getApp();
      app.globalData.isLoggedIn = false;
      app.globalData.user = null;

      resolve();
    } catch (e) {
      resolve();
    }
  });
}

module.exports = {
  checkLoginStatus,
  triggerLogin,
  performLogin,
  requireLogin,
  logout
};

