package mp

import (
    "github.com/codesve/mpsve/config"
    "log"
    "net/http"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "github.com/wizjin/weixin"
    "strings"
    "regexp"
    "time"
    "strconv"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
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
	w.ReplyText("欢迎关注 - 灏谳记账！本账号已做用户控制，只有指定微信号才能使用。") // 有新人关注，返回欢迎消息
}


// 文本消息的处理函数
func TextHandler(w weixin.ResponseWriter, r *weixin.Request) {
	txt := r.Content
	w.ReplyText(commondHandler(txt, r))
}


// 语音识别消息的处理函数
func VoiceHandler(w weixin.ResponseWriter, r *weixin.Request) {
	txt := r.Recognition
	w.ReplyText(commondHandler(txt, r))
}


func MenuClickHandler(w weixin.ResponseWriter, r *weixin.Request) {
	menuKey := r.EventKey

	reTxt := "你点了菜单"
	t := time.Now()
	switch menuKey {
	case config.TodayBp:
		reTxt = commondHandler("q" + t.Format("2006-01-02"), r)
	case config.YesterdayBp:
		yesterDay := t.AddDate( 0, 0, -1 )
		reTxt = commondHandler("q" + yesterDay.Format("2006-01-02"), r)
	case config.LastweekBp:
		reTxt = "上周收支, 开发中...."
	case config.LastmonthBp:
		reTxt = "上月收支, 开发中...."
	case config.CommondKey:
		reTxt = config.CommondHelp
	}

	w.ReplyText(reTxt)
}


func commondHandler(commondTxt string, r *weixin.Request) string {

	log.Printf(commondTxt)

	//获取用户
	countWeixin := r.FromUserName
	name := countWeixin
	switch countWeixin {
	case "okJSatyJx3Rz7x8XJEytGYDFdEYA":
		name = "测试"
	case "okJSat4PDKgoU4C6lrdlGoupnaag":
		name = "老婆"
	case "okJSat5Sgjm0DrhqErjdg3y0u3z8":
		name = "老公"
	default:
		return "你的账号不允许使用，请取消关注。谢谢。"
	}

	switch {
	case strings.HasPrefix(commondTxt, "123"):
		//新的收入项
		return creatABalance(countWeixin, name, "收入", commondTxt)

	case strings.HasPrefix(commondTxt, "321"):
		//新的支出项
		return creatABalance(countWeixin, name, "支出", commondTxt)

	case strings.HasPrefix(commondTxt, "q"):
		//查询收支情况
		return queryBalances(commondTxt)

	case strings.HasPrefix(commondTxt, "d"):
		//删除某条收支记录
		return deleteABalance(commondTxt, name)

	default:
		return "指令错误，点击'指令'菜单获取帮助。"
	}
	
}


//插入一条收支记录
func creatABalance(countWeixin string, name string, payType string, commondTxt string) string {

		var pay float64
		log.Printf(commondTxt[len(commondTxt)-1:])
		if commondTxt[len(commondTxt)-1:] == "!" {
			pay = 0
		} else {
			//获取金额
			payExp := regexp.MustCompile(`(\d+\.*\d+)`)
			mutchArr := payExp.FindStringSubmatch(commondTxt[3:])
			if len(mutchArr) < 2 {
				return "找不到有效金额。"
			}
			payT, err := strconv.ParseFloat(mutchArr[len(mutchArr)-1], 64)
			if err != nil {
				return "找不到有效金额。"
			}
			pay = payT
		}

		if name == "测试" {
			pay = 0
		}

		//创建时间
		t := time.Now()
		createTime := t.Format("2006-01-02")//t.Format("2006-01-02 15:04:05")

		//Uid
		uid := strconv.FormatInt(t.Unix(), 10)

		balance := config.Balance{uid, countWeixin, name, payType, commondTxt[3:], pay, createTime}

		session, err := mgo.Dial("127.0.0.1")
        if err != nil {
        	log.Printf("添加记录出错。", err)
        	return "添加记录出错。"
        }
        defer session.Close()

        // Optional. Switch the session to a monotonic behavior.
        session.SetMode(mgo.Monotonic, true)

        c := session.DB("sve").C("balance")
        err = c.Insert(&balance)
        if err != nil {
            log.Printf("添加记录出错。", err)
        	return "添加记录出错。"
        }

        reTxt := name + " 添加" + payType + "项：" + commondTxt[3:] + ", 成功！\n" +
        "金额为：" + strconv.FormatFloat(pay, 'f', 2, 64) + "元\n" +
        "如有错误，输入指令: d" + uid + ", 删除记录。"
        
        postTo(name, reTxt)
        return reTxt
} 


func deleteABalance(commondTxt string, name string) string {

	session, err := mgo.Dial("127.0.0.1")
    if err != nil {
    	log.Printf("删除记录出错。", err)
        return "删除记录出错。"
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c := session.DB("sve").C("balance")
    errRem := c.Remove(bson.M{"uid": commondTxt[1:]})
    if errRem != nil {
    	log.Printf("删除记录出错。", errRem)
    	return "删除记录出错。"
    }

    reTxt := name + " 删除了记录" + commondTxt[1:]

    postTo(name, reTxt)
    return reTxt
}


func queryBalances(commondTxt string) string {
	session, err := mgo.Dial("127.0.0.1")
    if err != nil {
    	log.Printf("查询出错。", err)
        return "查询出错。"
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c := session.DB("sve").C("balance")
    var balances []config.Balance
    err = c.Find(bson.M{"createtime": commondTxt[1:]}).All(&balances)
    if err != nil {
        log.Printf("查询出错。", err)
        return "查询出错。"
    }

    reTxt := commondTxt[1:] + `收支：
    `
    var outM float64
    var inM float64
    for i, balance := range balances {
    	if balance.PayType == "收入" {
    		inM += balance.Pay
    	} else {
    		outM += balance.Pay
    	}
        reTxt += strconv.Itoa(i+1) + `. ` + balance.Content + `
    `
    }


    reTxt += `
    收入合计: ` + strconv.FormatFloat(inM, 'f', 2, 64) + `元
    支出合计: ` + strconv.FormatFloat(outM, 'f', 2, 64) + `元
    总合计：` + strconv.FormatFloat(inM - outM, 'f', 2, 64) + `元
    `

	return reTxt
}

func postTo(name string, reTxt string) {
    if name == "老婆" {
		config.Mux.PostText("okJSat5Sgjm0DrhqErjdg3y0u3z8", reTxt)
	} else if name == "老公" {
		config.Mux.PostText("okJSat4PDKgoU4C6lrdlGoupnaag", reTxt)
	} else {
		config.Mux.PostText("okJSat5Sgjm0DrhqErjdg3y0u3z8", reTxt)
		config.Mux.PostText("okJSat4PDKgoU4C6lrdlGoupnaag", reTxt)
	}
}