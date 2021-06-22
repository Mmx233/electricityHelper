package elec

import (
	"github.com/Mmx233/tool"
	"strconv"
	"strings"
	"time"
)

const loginPageUrl = "http://222.204.3.210/ssdf/Account/LogOn"
const eleInfoPageUrl = "http://222.204.3.210/ssdf/EEMQuery/EEMBalance"

type eleInfo struct {
	UsedTotal     float32   `json:"used_total"`      //最新读数(度)
	UsedThisMonth float32   `json:"used_this_month"` //本月度数(度)
	Balance       float32   `json:"balance"`         //余额
	EleBalance    float32   `json:"ele_balance"`     //电量余额
	UpdateAt      time.Time `json:"update_at"`       //抄表时间
}

func parseFloat32(a string) float32 {
	t, _ := strconv.ParseFloat(strings.TrimSpace(a), 32)
	return float32(t)
}

// GetInfo  接收寝室号
func GetInfo(dormId uint) (*eleInfo, error) {
	header, i, e := tool.HTTP.GetReader(
		loginPageUrl,
		nil, nil, nil, true)
	if e != nil {
		return nil, e
	}
	_ = i.Close()
	cookies := tool.Cookie.Decode(header.Get("Set-Cookie"), nil)

	header, i, e = tool.HTTP.PostReader(
		loginPageUrl,
		map[string]interface{}{
			"Content-Type": "application/x-www-form-urlencoded",
		}, nil, map[string]interface{}{
			"UserName": dormId,
		}, cookies, false)
	if e != nil {
		return nil, e
	}
	_ = i.Close()

	tool.Cookie.Decode(header.Get("Set-Cookie"), cookies)

	d, e := tool.HTTP.GetGoquery(
		eleInfoPageUrl,
		nil, nil, cookies, true)

	if e != nil {
		return nil, e
	}

	t := d.Find("table tbody").First().Find("tr").First()
	usedTotal := parseFloat32(t.Find("td").First().Next().Text())
	updateAt, _ := time.ParseInLocation("2006-01-02 15:04", strings.TrimSpace(t.Find("td").Last().Text()), time.Local)
	t = t.Next()
	usedThisMonth := parseFloat32(t.Find("td").First().Next().Text())
	t = t.Next()
	balance := parseFloat32(t.Find("td").First().Next().Text())
	t = t.Next()
	eleBalance := parseFloat32(t.Find("td").First().Next().Text())

	return &eleInfo{
		UsedTotal:     usedTotal,
		UsedThisMonth: usedThisMonth,
		Balance:       balance,
		EleBalance:    eleBalance,
		UpdateAt:      updateAt,
	}, nil
}
