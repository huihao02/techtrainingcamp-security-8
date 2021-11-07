# 《抓到你了——具备安全防护能力的账号系统》 第八组

## 技术栈

- 语言： go
- 框架： gin
- 数据库： mysql
- 前端页面：layui

## API

### 验证码

- GET /captcha/image 获取图片验证码
- GET /captcha/phone 获取手机验证码
- GET /captcha/verify/:value 验证验证码

### 登录

- POST /login 账号密码登录
- POST /login/phone 手机号登录

### 注册

POST /register

### 登出/注销

GET /logout
