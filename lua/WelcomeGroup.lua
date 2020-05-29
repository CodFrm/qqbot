local Api = require("coreApi")

function ReceiveFriendMsg(CurrentQQ, data)
    return 1
end

function ReceiveGroupMsg(CurrentQQ, data)
    return 1
end

function ReceiveEvents(CurrentQQ, data, extData)
    if data.MsgType == "ON_EVENT_GROUP_JOIN" then
        ret = Api.Api_GetUserInfo(CurrentQQ, extData.UserID)
        if ret["data"]["gender"] == 2 then
            Api.Api_SendMsg(
                CurrentQQ,
                {
                    toUser = data.FromUin,
                    sendToType = 2,
                    sendMsgType = "PicMsg",
                    groupid = 0,
                    content = "",
                    atUser = 0,
                    picUrl = "http://img.wxcha.com/file/201911/15/be3e750b58.jpg",
                    picBase64Buf = "",
                    fileMd5 = ""
                }
            )
            Api.Api_SendMsg(
                CurrentQQ,
                {
                    toUser = data.FromUin,
                    sendToType = 2,
                    sendMsgType = "PicMsg",
                    content = string.format("[秀图%d]", math.random(40000, 40005)),
                    atUser = 0,
                    voiceUrl = "",
                    voiceBase64Buf = "",
                    picUrl = "http://q1.qlogo.cn/g?b=qq&nk=" .. extData.UserID .. "&s=640",
                    picBase64Buf = "",
                    fileMd5 = ""
                }
            )
        end
    end
    return 1
end
