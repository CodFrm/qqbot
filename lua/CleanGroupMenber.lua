local Api = require("coreApi")

function ReceiveFriendMsg(CurrentQQ, data)
    return 1
end

adminCache = nil -- 貌似是缓存不下来的💦

function toString(t)
    local bytearr = {}
    i = 1
    print(string.char(t[1]))
    while t[i] ~= 0 do
        table.insert(bytearr, string.char(t[i]))
        i = i + 1
    end
    return table.concat(bytearr)
end

function isAdmin(CurrentQQ, qqgroup, qq, memberlist)
    if adminCache == nil then
        adminCache = {}
        NextToken = ""
        while true do
            list =
                Api.Api_CallFunc(
                CurrentQQ,
                "friendlist.GetTroopListReqV2",
                {
                    NextToken = NextToken
                }
            )
            for i = 1, list["Count"] do
                id = list["TroopList"][i].GroupId
                if adminCache[id] == nil then
                    adminCache[id] = {}
                end
                adminCache[id][list["TroopList"][i].GroupOwner] = 1
            end
            NextToken = list["NextToken"]
            if NextToken == "" then
                break
            end
        end
    end
    if memberlist ~= nil then
        for i, v in ipairs(memberlist) do
            if v["GroupAdmin"] == 1 then
                adminCache[qqgroup][v["MemberUin"]] = 1
            end
        end
    end
    if adminCache[qqgroup][qq] == 1 then
        return true
    end
    return false
end

function ReceiveGroupMsg(CurrentQQ, data)
    if string.find(data.Content, "^踢潜水 %d+ .+模式$") ~= nil then
        num = tonumber(string.match(data.Content, "^踢潜水 (%d+) .+模式$"))
        mode = string.match(data.Content, "^踢潜水 %d+ (.+)模式$")
        MemberList = {}
        LastUin = 0
        data.FromGroupId = 614202391
        while true do
            list =
                Api.Api_CallFunc(
                CurrentQQ,
                "friendlist.GetTroopMemberListReq",
                {
                    GroupUin = data.FromGroupId,
                    LastUin = LastUin
                }
            )
            for i = 1, list["Count"] do
                table.insert(MemberList, list["MemberList"][i])
            end
            LastUin = list["LastUin"]
            if LastUin == 0 then
                break
            end
        end
        if isAdmin(CurrentQQ, data.FromGroupId, data.FromUserId, MemberList) == false then
            Api.Api_SendMsg(
                CurrentQQ,
                {
                    toUser = data.FromGroupId,
                    sendToType = 2,
                    sendMsgType = "TextMsg",
                    groupid = 0,
                    content = " 没有权限",
                    atUser = data.FromUserId
                }
            )
            return 1
        end
        if num > 50 then
            Api.Api_SendMsg(
                CurrentQQ,
                {
                    toUser = data.FromGroupId,
                    sendToType = 2,
                    sendMsgType = "TextMsg",
                    groupid = 0,
                    content = " 单次移出不可超过50人",
                    atUser = data.FromUserId
                }
            )
            return 1
        end
        if mode == "舔狗" then
            removeUser = topk(MemberList, num, tgsort)
        elseif mode == "面子" then
            removeUser = topk(MemberList, num, mzsort)
        elseif mode == "普通" then
            removeUser = topk(MemberList, num, zysort)
        end
        print("准备踢出")
        for i, v in pairs(removeUser) do
            print(i, v["NickName"], v["MemberUin"])
            Api.Api_GroupMgr(
                CurrentQQ,
                {
                    ActionType = 3,
                    GroupID = data.FromGroupId,
                    ActionUserID = v["MemberUin"],
                    Content = ""
                }
            )
        end
        Api.Api_SendMsg(
            CurrentQQ,
            {
                toUser = data.FromGroupId,
                sendToType = 2,
                sendMsgType = "TextMsg",
                groupid = 0,
                content = " 移除完毕",
                atUser = data.FromUserId
            }
        )
    end
    return 1
end

-- v1>v2
function tgsort(val1, val2)
    if val1["Gender"] == 1 then
        return false
    end
    if val2["Gender"] == 1 then
        return true
    end
    return mzsort(val1, val2)
end

function mzsort(val1, val2)
    if val1["MemberLevel"] >= 4 then
        return false
    end
    if val2["MemberLevel"] >= 4 then
        return true
    end
    if val1["GroupCard"] ~= "" then
        return false
    end
    if val2["GroupCard"] ~= "" then
        return true
    end
    return zysort(val1, val2)
end

function zysort(val1, val2)
    return val1["LastSpeakTime"] < val2["LastSpeakTime"]
end

function topk(arr, num, sort)
    ret = {}
    for i = 1, num do
        if arr[i] == nil then
            num = 1
            break
        end
        table.insert(ret, arr[i])
    end
    table.sort(ret, sort)
    for i, v in ipairs(arr) do
        if tgsort(v, ret[1]) then
            --替换,下沉
            ret[1] = v
            for i = 1, num / 2 do
                pos = i * 2
                if tgsort(ret[i], ret[pos]) then
                    tmp = ret[i]
                    ret[i] = ret[pos]
                    ret[pos] = tmp
                end
                pos = i * 2 + 1
                if num >= pos and tgsort(ret[i], ret[pos]) then
                    tmp = ret[i]
                    ret[i] = ret[pos]
                    ret[pos] = tmp
                end
            end
        end
    end
    table.sort(ret, sort)
    return ret
end

function ReceiveEvents(CurrentQQ, data, extData)
    return 1
end
