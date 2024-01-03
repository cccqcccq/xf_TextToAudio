# 讯飞文字转语音webapi
golang实现讯飞文字转语音webapi

官方文档https://www.xfyun.cn/doc/tts/online_tts/API.html#%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E

![image](https://github.com/cccqcccq/xf_TextToAudio/assets/117553354/515d11f0-f377-4595-83f9-15b6fe4f2a31)


将常量改成在控制台中对应的值

![image](https://github.com/cccqcccq/xf_TextToAudio/assets/117553354/837ea422-3bb1-4321-b572-ca14e82e50de)


将文本改成要生成的文本
经常出现bad handshake但偶尔可以连接成功是因为转base64时使用了URLEncoding而不是StdEncoding
