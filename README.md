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
- POST /captcha/verify/:value 验证验证码
- POST /captcha/verify-phone/:value 验证手机验证码
- 
```shell
# 获取验证码 响应头的`Set-Cookie`的字段 验证码返回值 作`curl`验证用
$curl http://127.0.0.1:9999/captcha/phone -X POST -i

---response start
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Set-Cookie: mysession=MTYzNjc5MjcwOHxEdi1CQkFFQ180SUFBUkFCRUFBQU9fLUNBQUVHYzNSeWFXNW5EQThBRFdOaGNIUmphR0V0Y0dodmJtVUdjM1J5YVc1bkRCWUFGRWh6WTI4MFIzTmxXVWMzYjIxdmNWVjNTR0ZVfLn2_Q85kbWOXj7qY4rN83sUP09ZOJWvVZJJYEE5ep4i; Path=/; Expires=Mon, 13 Dec 2021 08:38:28 GMT; Max-Age=2592000
Date: Sat, 13 Nov 2021 08:38:28 GMT
Content-Length: 109

{"Code":0,"Data":{"DecisionType":0,"ExpireTime":600000000000,"VerifyCode":[3,1,5,6,8,6]},"Message":"success"}
---response end
```

```shell
# 验证验证码
$curl http://127.0.0.1:9999/captcha-phone/verify/315686 -X POST \
-H 'Cookie: mysession=MTYzNjc5MjcwOHxEdi1CQkFFQ180SUFBUkFCRUFBQU9fLUNBQUVHYzNSeWFXNW5EQThBRFdOaGNIUmphR0V0Y0dodmJtVUdjM1J5YVc1bkRCWUFGRWh6WTI4MFIzTmxXVWMzYjIxdmNWVjNTR0ZVfLn2_Q85kbWOXj7qY4rN83sUP09ZOJWvVZJJYEE5ep4i;'

---response start
{"Code":0,"Message":"success"}
---response end
```

### 登录

- POST /login 账号密码登录
- POST /login/phone 手机号登录

### 注册

POST /register

### 登出/注销

GET /logout
