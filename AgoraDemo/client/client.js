import AgoraRTM from 'agora-rtm-sdk'

// login 方法参数
let options = {
    uid: "",
    token: ""
}

// 你的 app ID
const appID = "98a56b7f58ec4c07877033572c3a6fd1"

// 初始化客户端
const client = AgoraRTM.createInstance(appID)

// 客户端事件监听
// 显示对端发送的消息
client.on('MessageFromPeer', function (message, peerId) {
    document.getElementById("log").appendChild(document.createElement('div')).append("Message from: " + peerId + " Message: " + message.text)
})

// 显示连接状态变化
client.on('ConnectionStateChanged', function (state, reason) {
    document.getElementById("log").appendChild(document.createElement('div')).append("State changed To: " + state + " Reason: " + reason)
})

let channel = client.createChannel("demoChannel")

channel.on('ChannelMessage', function (message, memberId) {
    document.getElementById("log").appendChild(document.createElement("div")).append(memberId + ": " + message.text);
})

// 显示频道
channel.on('MemberJoined', function (memberId) {
    document.getElementById("log").appendChild(document.createElement('div')).append(memberId + " joined the channel")
})

// 频道成员
channel.on('MemberLeft', function (memberId) {
    document.getElementById("log").appendChild(document.createElement('div')).append(memberId + " left the channel")
})

// 按钮行为定义
window.onload = function () {
    // 按钮逻辑
    // 登录
    document.getElementById("login").onclick = async function () {
        options.uid = document.getElementById("userID").value.toString()
        options.token = await fetchToken(options.uid)
        await client.login(options)
    }

    // 登出
    document.getElementById("logout").onclick = async function () {
        await client.logout()
    }

    // 创建并加入频道
    document.getElementById("join").onclick = async function () {
        // Channel event listeners
        // Display channel messages
        await channel.join().then (() => {
            document.getElementById("log").appendChild(document.createElement('div')).append("You have successfully joined channel " + channel.channelId)
        })
    }

    // 离开频道
    document.getElementById("leave").onclick = async function () {

        if (channel != null) {
            await channel.leave()
        }
        else{
            console.log("Channel is empty")
        }
    }

    // 发送点对点消息
    document.getElementById("send_peer_message").onclick = async function () {
        let peerId = document.getElementById("peerId").value.toString()
        let peerMessage = document.getElementById("peerMessage").value.toString()
        console.log("=========")
        console.log(peerMessage)
        await client.sendMessageToPeer(
            { text: peerMessage },
            peerId,
        ).then(sendResult => {
            if (sendResult.hasPeerReceived) {
                document.getElementById("log").appendChild(document.createElement('div')).append("Message has been received by: " + peerId + " Message: " + peerMessage)
            } else {
                document.getElementById("log").appendChild(document.createElement('div')).append("Message sent to: " + peerId + " Message: " + peerMessage)
            }
        })
    }

    // 发送频道消息
    document.getElementById("send_channel_message").onclick = async function () {
        let channelMessage = document.getElementById("channelMessage").value.toString()
        if (channel != null) {
            await channel.sendMessage({ text: channelMessage }).then(() => {
                document.getElementById("log").appendChild(document.createElement('div')).append("Channel message: " + channelMessage + " from " + channel.channelId)
            }
            )
        }
    }
}

function fetchToken(uid) {
    return new Promise(function (resolve) {
        axios.post('http://localhost:8082/fetch_rtm_token', {
            uid: uid,
        }, {
            headers: {
                'Content-Type': 'application/json; charset=UTF-8'
            }
        })
            .then(function (response) {
                const token = response.data.token;
                resolve(token);
            })
            .catch(function (error) {
                console.log(error);
            });
    })
}