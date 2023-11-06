# go-jwt-webdemo
### 以JWT为基础实现身份认证，进行AI绘画的网站demo <br>
技术栈:go+gin+redis+vue
<br>流程概述:
1. 前端注册登录结合redis实现，JWT token在用户POST /login的时候，由服务端分配
2. login成功后，用户再请求其它的request时，带上该JWT token，交由服务端的middleware认证
3. POST /api/order 请求头带上cookie存储的token，后端进行AI绘画任务提交，返回uid
4. POST /api/config 请求头带上cookie存储的token，读取后端config.yaml中的authorization
5. 前端POST https://open.nolibox.com/prod-open-aigc/engine/status/${uid} ，请求体带上authorization，获得AI绘画结果