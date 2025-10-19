# TabBar 图标优化说明

## 🎨 图标设计更新

根据微信小程序和开源图标库的最佳实践，我们优化了以下图标：

### 更新后的图标

| 功能 | 图标类型 | 设计理念 |
|------|---------|---------|
| **号码生成** | 🏠 首页图标 | 简洁的房子图标，代表主页入口 |
| **历史开奖** | 📅 日历图标 | 更直观地表示历史记录和时间序列 |
| **我的号码** | 🎫 票据图标 | 结合彩票球和收藏星标，表示保存的号码 |

### 设计优势

#### 📅 历史开奖 - 日历图标
- ✅ **更直观**：日历比时钟更能表达"历史记录"的含义
- ✅ **符合习惯**：大多数应用使用日历表示历史/记录功能
- ✅ **清晰识别**：81x81尺寸下识别度更高
- 参考：微信"朋友圈"、支付宝"账单"等都使用日历图标

#### 🎫 我的号码 - 票据图标
- ✅ **主题相关**：票据样式更贴合彩票场景
- ✅ **功能明确**：带收藏星标，表示"我的收藏"
- ✅ **视觉丰富**：彩票球元素增加趣味性
- 参考：各类票务、收藏功能常用此类图标

## 📦 文件清单

### SVG源文件（可编辑）
- `home.svg` / `home-active.svg` - 首页图标
- `history.svg` / `history-active.svg` - 历史开奖图标（日历样式）
- `numbers.svg` / `numbers-active.svg` - 我的号码图标（票据样式）

### PNG文件（微信小程序使用）
需要转换生成以下6个PNG文件：
- `home.png` / `home-active.png`
- `history.png` / `history-active.png`
- `numbers.png` / `numbers-active.png`

## 🔧 转换方法

### 方法1：使用浏览器工具（推荐）✨
1. 双击打开 `create-icons.html` 文件
2. 浏览器中查看3组图标预览
3. 点击"下载普通版"和"下载激活版"按钮
4. 自动下载6个PNG文件（81x81像素）
5. 将PNG文件放到 `front/images/` 目录
6. 完成！

### 方法2：在线转换工具
1. 访问 https://convertio.co/svg-png/
2. 上传SVG文件
3. 设置输出尺寸为 **81x81 像素**
4. 下载PNG文件

### 方法3：使用ImageMagick（命令行）
```bash
# 安装ImageMagick
# macOS: brew install imagemagick
# Windows: choco install imagemagick
# Linux: sudo apt install imagemagick

# 批量转换
convert -background none -resize 81x81 home.svg home.png
convert -background none -resize 81x81 home-active.svg home-active.png
convert -background none -resize 81x81 history.svg history.png
convert -background none -resize 81x81 history-active.svg history-active.png
convert -background none -resize 81x81 numbers.svg numbers.png
convert -background none -resize 81x81 numbers-active.svg numbers-active.png
```

## 🎯 图标规范

### 尺寸
- **标准尺寸**: 81x81 像素
- **格式**: PNG（透明背景）
- **分辨率**: @1x（微信小程序会自动适配）

### 颜色
- **未选中**: `#666666` (灰色)
- **选中**: `#ff4757` (红色)

### 设计原则
1. **简洁明了**：图标应该一眼就能理解含义
2. **视觉统一**：所有图标保持一致的设计风格
3. **识别度高**：81x81尺寸下清晰可辨
4. **主题相关**：与彩票应用场景相符

## 📚 参考图标库

- **Feather Icons**: https://feathericons.com/
- **Heroicons**: https://heroicons.com/
- **Lucide Icons**: https://lucide.dev/
- **微信官方图标**: https://developers.weixin.qq.com/miniprogram/design/

## 🔄 如何更换图标

如果需要更换其他图标：

1. **编辑SVG文件**
   - 打开对应的 `.svg` 文件
   - 修改SVG路径代码
   - 保持 `81x81` 的 `viewBox`
   - 颜色使用 `#666666`（普通）或 `#ff4757`（激活）

2. **重新转换**
   - 使用 `create-icons.html` 重新生成PNG
   - 或使用在线工具/命令行转换

3. **替换文件**
   - 将新的PNG文件放到 `images/` 目录
   - 覆盖旧文件即可

## ✨ 效果对比

### 优化前
- 历史开奖：时钟图标 ⏰（不够直观）
- 我的号码：列表图标 📄（缺少特色）

### 优化后
- 历史开奖：日历图标 📅（直观表示记录）
- 我的号码：票据图标 🎫（符合彩票主题）

## 💡 使用建议

1. **保持一致性**：所有页面的图标风格应保持统一
2. **测试可见性**：在不同背景色下测试图标可见度
3. **用户测试**：收集用户反馈，确保图标含义清晰
4. **定期优化**：根据用户使用情况持续优化图标设计

## 📝 版本历史

- **v2.0** (2025-10-19)
  - 优化历史开奖图标：时钟 → 日历
  - 优化我的号码图标：列表 → 票据（带彩票球和星标）
  - 添加在线转换工具
  - 完善文档说明

- **v1.0** (初始版本)
  - 基础图标设计
  - 使用Feather Icons风格
