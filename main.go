package main

import (
    "log"
    "net/http"
    "github.com/codesve/mpsve/config"
    "github.com/codesve/mpsve/mp"
    "github.com/wizjin/weixin"
)


func main() {
    
    // 实例化weixin
    config.Mux = weixin.New(config.Token, config.AppId, config.AppSecret)
    
    /*
        接收 接口注册
    */
    config.Mux.HandleFunc(weixin.MsgTypeText, mp.TextHandler)// 注册文本消息的处理函数
    config.Mux.HandleFunc(weixin.MsgTypeVoice, mp.VoiceHandler)// 注册语音消息的处理函数
    config.Mux.HandleFunc(weixin.MsgTypeEventSubscribe, mp.Subscribe)// 注册关注事件的处理函数

    config.Mux.HandleFunc(weixin.MsgTypeEventClick, mp.MenuClickHandler)// 注册菜单事件的处理函数

    //创建菜单
    // mp.PostMenu()

    /*
        给微信调的API
    */
    http.HandleFunc(config.RootPath + "/check", config.Mux.ServeHTTP)//接入微信


    log.Printf("服务器已启动，http://localhost:" + config.Port)
    http.Handle(config.RootPath, config.Mux) // 注册接收微信服务器数据的接口URI
    http.ListenAndServe(":" + config.Port, nil) // 启动接收微信数据服务器
}
