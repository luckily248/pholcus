package spiders

// 基础包
import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"                          //DOM解析
	"github.com/henrylee2cn/pholcus/crawl/downloader/context" //必需
	"github.com/henrylee2cn/pholcus/reporter"                 //信息输出
	. "github.com/henrylee2cn/pholcus/spider"                 //必需
	"io/ioutil"
	"net/http"
	"net/url"
	// . "github.com/henrylee2cn/pholcus/spider/common" //选用
)

// 设置header包
import (
// "net/http" //http.Header
)

// 编码包
import (
// "encoding/xml"
// "encoding/json"
)

// 字符串处理包
import (
	// "regexp"
	// "strconv"
	"strings"
)

// 其他包
import (
// "fmt"
// "math"
)

type Root struct {
	Page    string
	Records string
	Sidx    string
	Sord    string
	Total   string
	Rows    []Rows
}
type Rows struct {
	Cell []string
	Id   string
}

func init() {
	Cib.AddMenu()
}

var urlroot string = "http://www.cib.com.cn"

func getCookie(url string) (cookies []*http.Cookie, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	cookies = resp.Cookies()
	return
}

// 考拉海淘,海外直采,7天无理由退货,售后无忧!考拉网放心的海淘网站!
var Cib = &Spider{
	Name:        "兴业银行",
	Description: "兴业银行数据 [Auto Page] [http://www.cib.com.cn/cn/index.html]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Optional: &Optional{},
	RuleTree: &RuleTree{
		// Spread: []string{},
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{"url": urlroot + "/cn/index.html", "rule": "获取版块URL"})
		},

		Nodes: map[string]*Rule{

			"获取版块URL": &Rule{
				ParseFunc: func(self *Spider, resp *context.Response) {
					reporter.Log.Printf("    start get urls\n")
					query := resp.GetDom()
					lis := query.Find("div.cb div.brd ul.lis2 li a")
					lis.Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							//							if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
							//								url = urlroot + url
							//							}
							//							reporter.Log.Printf("    get urls:%s  -Text:(%s)\n", url, s.Text())
							//							if strings.Contains(s.Text(), "个人存贷利率") {
							//								self.AddQueue(map[string]interface{}{"url": url, "rule": "个人存贷利率"})
							//							} else if strings.Contains(s.Text(), "企业存贷利率") {
							//								//因为当前跳转路径不对 要更正路径
							//								url = "http://www.cib.com.cn/cn/Corporate_Banking/Deposit_Loan_Rates/RMB_Loan_Rates.html"
							//								self.AddQueue(map[string]interface{}{"url": url, "rule": "企业存贷利率"})
							//							} else
							if strings.Contains(s.Text(), "外汇牌价") {
								self.AddQueue(map[string]interface{}{"url": url, "rule": "外汇牌价"})
							}
						}
					})
					//reporter.Log.Printf("    self.Description:%s\n")
					//for _, v := range self. {
					//	reporter.Log.Printf("    get rule:%v\n", v)
					//}
					reporter.Log.Printf("    end get urls\n")
				},
			},

			"个人存贷利率": &Rule{
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"title",
					"co1",
					"co2",
					"co3",
					"remark",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					reporter.Log.Printf("    个人存贷利率 start\n")
					query := resp.GetDom()
					var title string
					var co1 []string
					var co2 []string
					var co3 []string
					var remark string
					remark = query.Find("div.add p").Text()
					query.Find("div.add table tbody tr").Each(func(i int, s *goquery.Selection) {
						//reporter.Log.Printf("row start :%d\n", i)
						if i == 0 {
							title = s.Find("strong").Text()
							//reporter.Log.Printf("title text :%s\n", title)
						} else {
							s.Find("td").Each(func(i int, s *goquery.Selection) {
								//reporter.Log.Printf("td text:%s\n", s.Find("td").Text())
								switch i {
								case 0:
									co1 = append(co1, s.Text())
								case 1:
									co2 = append(co2, s.Text())
								case 2:
									co3 = append(co3, s.Text())
								}
							})
						}
					})
					resp.AddItem(map[string]interface{}{
						self.GetOutFeild(resp, 0): title,
						self.GetOutFeild(resp, 1): co1,
						self.GetOutFeild(resp, 2): co2,
						self.GetOutFeild(resp, 3): co3,
						self.GetOutFeild(resp, 4): remark,
					})
					reporter.Log.Printf("    title:%s\n", title)
					for _, v := range co1 {
						reporter.Log.Printf("    %s\n", v)
					}
					for _, v := range co2 {
						reporter.Log.Printf("    %s\n", v)
					}
					for _, v := range co3 {
						reporter.Log.Printf("    %s\n", v)
					}
					reporter.Log.Printf("    remark:%s\n", remark)
					reporter.Log.Printf("    个人存贷利率 end\n")
				},
			},
			"企业存贷利率": &Rule{
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"title",
					"co1",
					"co2",
					"remark",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					reporter.Log.Printf("    企业存贷利率 start\n")
					query := resp.GetDom()
					var title string
					var co1 []string
					var co2 []string
					var remark string
					remark = query.Find("div.add p span").Text()
					if nexturl, ok := query.Find("div.add p a").Attr("href"); ok {
						nexturl = urlroot + nexturl
						self.AddQueue(map[string]interface{}{"url": nexturl, "rule": "个人人民币贷款利率"})
					}
					query.Find("div.add table tbody tr").Each(func(i int, s *goquery.Selection) {
						//reporter.Log.Printf("row start :%d\n", i)
						if i == 0 {
							s.Find("strong").Each(func(i int, s *goquery.Selection) {
								title = s.Text() + s.Text()
							})
							//reporter.Log.Printf("title text :%s\n", title)
						} else {
							s.Find("td").Each(func(i int, s *goquery.Selection) {
								//reporter.Log.Printf("td text:%s\n", s.Find("td").Text())
								switch i {
								case 0:
									co1 = append(co1, s.Text())
								case 1:
									co2 = append(co2, s.Text())
								}
							})
						}
					})
					resp.AddItem(map[string]interface{}{
						self.GetOutFeild(resp, 0): title,
						self.GetOutFeild(resp, 1): co1,
						self.GetOutFeild(resp, 2): co2,
						self.GetOutFeild(resp, 3): remark,
					})
					reporter.Log.Printf("    title:%s\n", title)
					for _, v := range co1 {
						reporter.Log.Printf("    %s\n", v)
					}
					for _, v := range co2 {
						reporter.Log.Printf("    %s\n", v)
					}
					reporter.Log.Printf("    remark:%s\n", remark)
					reporter.Log.Printf("    企业存贷利率 end\n")
				},
			},
			"个人人民币贷款利率": &Rule{
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"title",
					"co1",
					"co2",
					"remark",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					reporter.Log.Printf("    个人人民币贷款利率 start\n")
					query := resp.GetDom()
					var title string
					var co1 []string
					var co2 []string
					var remark string
					remark = query.Find("div.add p span").Text()
					query.Find("div.add table tbody tr").Each(func(i int, s *goquery.Selection) {
						//reporter.Log.Printf("row start :%d\n", i)
						if i == 0 {
							title = s.Find("strong").Text()
							//reporter.Log.Printf("title text :%s\n", title)
						} else {
							s.Find("td").Each(func(i int, s *goquery.Selection) {
								//reporter.Log.Printf("td text:%s\n", s.Find("td").Text())
								switch i {
								case 0:
									content1 := s.Text()
									if content1 == "" {
										content1 = s.Find("p").Text()
									}
									if content1 == "" {
										content1 = s.Find("strong").Text()
									}
									co1 = append(co1, content1)
								case 1:
									content2 := s.Text()
									if content2 == "" {
										content2 = s.Find("p").Text()
									}
									if content2 == "" {
										content2 = s.Find("strong").Text()
									}
									co2 = append(co2, content2)
								}
							})
						}
					})
					resp.AddItem(map[string]interface{}{
						self.GetOutFeild(resp, 0): title,
						self.GetOutFeild(resp, 1): co1,
						self.GetOutFeild(resp, 2): co2,
						self.GetOutFeild(resp, 3): remark,
					})
					reporter.Log.Printf("    title:%s\n", title)
					for _, v := range co1 {
						reporter.Log.Printf("    %s\n", v)
					}
					for _, v := range co2 {
						reporter.Log.Printf("    %s\n", v)
					}
					reporter.Log.Printf("    remark:%s\n", remark)
					reporter.Log.Printf("    个人人民币贷款利率 end\n")
				},
			},
			"外汇牌价": &Rule{
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"title",
					"date",
					"co1",
					"co2",
					"co3",
					"co4",
					"co5",
					"co6",
					"co7",
					"co8",
					"remark",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					reporter.Log.Printf("    外汇牌价 start\n")
					query := resp.GetDom()
					var title string
					var date string
					var remark string
					title = query.Find("div#title").Text()
					query.Find("div#labe_text").Each(func(i int, s *goquery.Selection) {
						switch i {
						case 0:
							date = s.Text()
						case 1:
							remark = s.Text()
						}
					})

					//由于调用了js脚本 直接获取json数据
					client := &http.Client{}
					u, _ := url.Parse("https://personalbank.cib.com.cn/pers/main/pubinfo/ifxQuotationQuery!list.do")
					p := u.Query()
					p.Set("_search", "false")
					p.Set("dataSet.nd", "1440055142817")
					p.Set("dataSet.rows", "100")
					p.Set("dataSet.page", "1")
					p.Set("dataSet.sidx", "")
					p.Set("dataSet.sord", "asc")
					u.RawQuery = p.Encode()
					reporter.Log.Printf("    url for data:%s\n", u.String())
					req, err := http.NewRequest("get", u.String(), strings.NewReader(""))
					if err != nil {
						reporter.Log.Printf("    err:%s\n", err.Error())
					}
					//					req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
					//					req.Header.Set("Connection", "keep-alive")
					//					req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
					//					req.Header.Set("Accept-Language", "zh-TW,zh;q=0.8,en-US;q=0.6,en;q=0.4")
					//					req.Header.Set("User-Agent", "Chrome/41.0.2272.118")
					//cookies := cookieJar.Cookies(req.URL)
					//reporter.Log.Printf("    cookie len:%d\n", len(cookies))
					//					for _, v := range strings.Split(setcookie, ";") {
					//						kv := strings.Split(v, "=")
					//						if kv[0] == "JSESSIONID" {
					//					cookie := &http.Cookie{}
					//					cookie.Name = "JSESSIONID"
					//					cookie.Value = "yG20VV1JrMnG8NRy6dyB7n7VdsT5BpfS2MKTJpcj8fLy1SdJL0hQ!-1703636766!1440055284148"
					//					req.AddCookie(cookie)
					//					//						}
					//					//					}
					cookies, err := getCookie(resp.Url)
					if err != nil {
						reporter.Log.Printf("    err:%s\n", err.Error())
					}
					for _, v := range cookies {
						req.AddCookie(v)
					}
					r, err := req.Cookie("JSESSIONID")
					if err != nil {
						reporter.Log.Printf("    err:%s\n", err.Error())
					}
					reporter.Log.Printf("    cookie:%s\n", r.Value)
					respforjs, err := client.Do(req)
					if err != nil {
						reporter.Log.Printf("    err:%s\n", err.Error())
					}
					defer respforjs.Body.Close()
					//					reporter.Log.Printf("    resp len cookies:%d\n", len(respforjs.Cookies()))
					//					for _, v := range respforjs.Cookies() {
					//						reporter.Log.Printf("    resp  cookies name:%s value:%s\n", v.Name, v.Value)
					//					}
					body, err := ioutil.ReadAll(respforjs.Body)
					if err != nil {
						reporter.Log.Printf("    err:%s\n", err.Error())
					}
					reporter.Log.Printf("    data:%s\n", body)

					//开始解析json
					data := &Root{}
					err = json.Unmarshal(body, data)
					if err != nil {
						reporter.Log.Printf("    err:%s\n", err.Error())
					}
					reporter.Log.Printf("    title:%s\n", title)
					reporter.Log.Printf("    date:%s\n", date)
					reporter.Log.Printf("    total:%s\n", data.Total)
					for _, v := range data.Rows {
						reporter.Log.Printf("    len rows :%d\n", len(v.Cell))
						for _, v := range v.Cell {
							reporter.Log.Printf("    %s\n", v)
						}
					}
					reporter.Log.Printf("    remark:%s\n", remark)
					reporter.Log.Printf("    外汇牌价 end\n")
					if len(data.Rows) <= 0 {
						reporter.Log.Printf("    data empty\n")
						return
					}
					resp.AddItem(map[string]interface{}{
						self.GetOutFeild(resp, 0):  title,
						self.GetOutFeild(resp, 1):  date,
						self.GetOutFeild(resp, 2):  data.Rows[0].Cell,
						self.GetOutFeild(resp, 3):  data.Rows[1].Cell,
						self.GetOutFeild(resp, 4):  data.Rows[2].Cell,
						self.GetOutFeild(resp, 5):  data.Rows[3].Cell,
						self.GetOutFeild(resp, 6):  data.Rows[4].Cell,
						self.GetOutFeild(resp, 7):  data.Rows[5].Cell,
						self.GetOutFeild(resp, 8):  data.Rows[6].Cell,
						self.GetOutFeild(resp, 9):  data.Rows[7].Cell,
						self.GetOutFeild(resp, 10): remark,
					})
				},
			},
		},
	},
}
