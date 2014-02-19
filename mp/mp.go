package mp

import (
    "github.com/codesve/mpsve/config"
    "log"
    "net/http"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "github.com/wizjin/weixin"
)


// 获取access_token, 主要是创建菜单时用
func GetAccessToken () string {
	resp, err := http.Get(config.AccessTokenAPI)

	if err != nil {
		log.Println("获取access_token失败：", err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取access_token失败：", err)
		} else {
			var res struct {
				AccessToken string `json:"access_token"`
				ExpiresIn   int64  `json:"expires_in"`
			}
			if err := json.Unmarshal(body, &res); err != nil {
				log.Println("转换access_token失败：", err)
			} else {
				return res.AccessToken
			}
		}
	}
	return ""
}


//创建菜单：
func PostMenu () {
	menuBuf := bytes.NewBufferString(config.TheMenu)
	menuApi := config.PostMenuAPI + "&access_token=" + GetAccessToken()
	resp, err := http.Post(menuApi, "application/json", menuBuf)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("创建菜单出错。")
	} else {
		bodyBuf := bytes.NewBuffer(body)
		log.Println(bodyBuf)
	}
}


// 关注事件的处理函数
func Subscribe(w weixin.ResponseWriter, r *weixin.Request) {
	w.ReplyText("欢迎关注 - 灏谳记账") // 有新人关注，返回欢迎消息
}

// 文本消息的处理函数
func TextHandler(w weixin.ResponseWriter, r *weixin.Request) {
	txt := r.Content			
	w.ReplyText(txt)			// 回复一条文本消息
	// w.PostText("Post:" + txt)	// 发送一条文本消息
}

// 语音识别消息的处理函数
func VoiceHandler(w weixin.ResponseWriter, r *weixin.Request) {
	recTxt := r.Recognition
	w.ReplyText(recTxt)
}

func MenuClickHandler(w weixin.ResponseWriter, r *weixin.Request) {
	menuKey := r.EventKey

	reTxt := "你点了菜单"
	switch menuKey {
	case config.TodayBp:
		reTxt = "今日收支"
	case config.YesterdayBp:
		reTxt = "昨日收支"
	case config.LastweekBp:
		reTxt = "上周收支"
	case config.LastmonthBp:
		reTxt = "上月收支"
	}

	w.ReplyText(reTxt)
}