package config

import (
    "github.com/wizjin/weixin"
)

//系统常量
const (
	RootPath 	= "/mpapi/v1"
	Port 	= "8082"
	CommondHelp = `指令格式: 
		1. 收入指令，123，如: 
123打牌赢了100元

		2. 支出指令，321，如: 
321奶粉花了218元

		3. 查询指令，q，如: 
q2014-02
q2014-02-21
q2014-02-11_2014-03-28

		4. 删除指令，d，如: 
d5；删除5这条记录

		说明：3、4指令需以文字信息发送；
		`
)

const (
	//微信开发接入所需的token，自定义的
	Token 		= "dev.codesve.com#mpapi#by#codesve"
	MpHost 		= "https://api.weixin.qq.com/cgi-bin"
	AppId 		= "wx0fad1d7cc7906cf2"
	AppSecret	= "aca97fa0cf76cbd1faa2b0a44889532b"
	AppParam 	= "appid=wx0fad1d7cc7906cf2&secret=aca97fa0cf76cbd1faa2b0a44889532b"
)


/*
	微信API
*/
const (
	//获取access_token
	AccessTokenAPI = MpHost + "/token?grant_type=client_credential&" + AppParam
	//创建菜单
	PostMenuAPI = MpHost + "/menu/create?" + AppParam

)

/*
	菜单key
*/
const (
	TodayBp 		= "today_bp_01"  		//今日收支
	YesterdayBp 	= "yesterday_bp_0201"	//昨日收支
	LastweekBp 		= "lastweek_bp_0202"	//上周收支
	LastmonthBp 	= "lastmonth_bp_0203"	//上月收支
	CommondKey		= "commond_03"			//指令
)


//菜单内容, 因为type是关键字，所以结构体行不通。
const TheMenu = 
`{"button": 
			[
				{
					"type" 			: "click",
					"name" 			: "今日收支",
	       			"key" 			: "today_bp_01" 
				},
				{
					"name" 			: "以往收支",
					"sub_button" 	:
					[
						{
							"type" 	: "click",
	           				"name" 	: "昨日收支",
	           				"key" 	: "yesterday_bp_0201" 
						},
						{
							"type" 	: "click",
	           				"name" 	: "上周收支",
	           				"key" 	: "lastweek_bp_0202" 
						},
						{
							"type" 	: "click",
	           				"name" 	: "上月收支",
	           				"key" 	: "lastmonth_bp_0203" 
						}
					]
				},
				{
					"type" 			: "click",
					"name" 			: "指令",
	       			"key" 			: "commond_03" 
				}
			]
}`

type Balance struct {
	Uid			string		//自增长id
	CountWeixin string 		//用户微信号, jfsjiwel323jl290jl
	Name		string		//用户姓名, 老婆、老公
	PayType		string		//收支类型，收入、支出
	Content		string		//内容，打牌赢了100元
	Pay 		float64		//金额, 112.28
	CreateTime	string		//时间，2014-01-02 15:04:05		
}

var Mux *weixin.Weixin  //weixin实例