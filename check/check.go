package check

import (
    "github.com/emicklei/go-restful"
    "github.com/codesve/mpsve/config"
    "sort"
    "io"
    "fmt"
    "crypto/sha1"
    "log"
    "net/http"
)

var (
    signature, echostr, timestamp       string
    nonce                               int
)

func Register(container *restful.Container) {
    ws := new(restful.WebService)
    ws.
        Path(config.RootPath + "/check")

	ws.Route(ws.GET("").To(mpCheck).
        // docs
        Doc("check wechart").Operation("check wechart").
        Param(ws.QueryParameter("signature", "微信给的签名").DataType("string")).
        Param(ws.QueryParameter("timestamp", "时间戳").DataType("string")).
        Param(ws.QueryParameter("nonce", "随机数").DataType("string")).
        Param(ws.QueryParameter("echostr", "随机字符串").DataType("string")))

    container.Add(ws)
}


func mpCheck(request *restful.Request, response *restful.Response) {

    log.Printf("微信来了")

    signature := request.QueryParameter("signature")
    timestamp := request.QueryParameter("timestamp")
    nonce := request.QueryParameter("nonce")
    echostr := request.QueryParameter("echostr")

    log.Printf("signature: " + signature + 
        "\n timestamp: " + timestamp + 
        "\n nonce: " + nonce + 
        "\n echostr: " + echostr + 
        "\n " )

    //字典排序
    tmpSortArr := []string{config.Token, timestamp, nonce}
	sort.Strings(tmpSortArr)
	sortStr := tmpSortArr[0] + tmpSortArr[1] + tmpSortArr[2]

    log.Printf("sortStr: " + sortStr)

	//sha1加密
	localSignature := str2sha1(sortStr)

    log.Printf("localSignature: " + localSignature)

	if localSignature == signature {
        // response.WriteAsJson(echostr)
        fmt.Fprintf(response.ResponseWriter, echostr)
    } else {
    	log.Printf("比对不正确，接入微信失败。")
        response.WriteErrorString(http.StatusNotFound, "比对不正确，接入微信失败。")
    }
}


func str2sha1(data string)string {
    t := sha1.New()
    io.WriteString(t, data)
    return fmt.Sprintf("%x", t.Sum(nil))
}
