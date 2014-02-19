package config

//系统常量
const (
	RootPath 	= "/mpapi/v1"
	Port 	= "8082"
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
					"type" 			: "view",
					"name" 			: "博客",
	       			"url" 			: "http://blog.codesve.com" 
				}
			]
}`

