## 这里是一个自用机器人的插件库
> 本机器人基于IOTQQ构建,机器人api等信息请前往[https://github.com/IOTQQ/IOTQQ](https://github.com/IOTQQ/IOTQQ)查看

## 安装

### 所需环境
* mysql
* redis

### 步骤
1.请先安装好IOTQQ环境,登录QQ
2.复制`config.example.yaml`文件重命名为`config.yaml`,并按照文本内容进行修改
3.启动`iotqq-plugins`

## 功能
本插件库包含了golang和lua两种语言的插件.
lua在lua目录中,可以将里面的lua文件复制粘贴到IOTQQ的Plugins文件夹中使用.
golang使用websocket监听.

`@机器人 帮助`可查看更多详细信息
* 好看的图片(需要p站账号和ssr)
* 图片旋转(高清重制需要apikey)
* 图片鉴黄(需要apikey)

...

