# todo-list

#### 介绍
基于go+gin+redis+mysql+gorm实现的todo-list项目，适合go新手小白入门项目
功能：登录注册，日志，token验证，redis缓存，RefreshToken，todo-list的增删改查，邮件发送，图形验证码。

#### 软件架构
软件架构说明


#### 安装教程

1. 新建数据库: todolist,字符集:utf8mb4;启动main.go文件,自动生成数据库表
2. 下载依赖,执行: go mod tidy -v
3. 修改配置文件: ./conf/Appconfig.yaml,包括mysql账号密码,邮箱以及邮箱授权码,redis账号等;utils/sendEmail.go的51行的邮箱也需要修改
4. go run main.go (执行前确保mysql以及redis已启动)
5. 将./todoList.postman_collection.json导入postman就可以进行接口测试了
6. 接口测试时需要修改参数和Authorization



#### 部分功能接口测试过程
1. 注册
   (1) 邮箱验证码发送接口获取,获取注册验证码（redis中可以查看）
   (2) 注册接口的code写入获取到的验证码
2. 登录
   (1) 图形验证码接口,b64s在浏览器中打开计算答案,和captchId值
   (2) 登录接口,pid:captchId值,value:答案
   (3) 保存access_token和refresh_token;tasks的增删改查会用到

备注：后续会把接口文档补上,有时间的话,会把具体实现和前端补上

#### 参考文献
1. 李文周老师博客（https://www.liwenzhou.com/）


#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

