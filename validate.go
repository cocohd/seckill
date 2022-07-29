package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"seckill/common"
	"seckill/datamodels"
	"seckill/encrypt"
	"seckill/rabbitmq"
	"strconv"
	"sync"
)

/*针对不同请求做负载均衡*/

//设置集群地址，最好内外IP
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

// SLB的ip，就是获取数量控制接口的主机
var getOne = "127.0.0.1"
var getOnePort = "8084"

var port = "8081"

// rabbitmq
var rabbitMqValidate *rabbitmq.RabbitMQ

var hashConsistent *common.ConsistentHash

// AccessControl 服务器端存放控制消息
type AccessControl struct {
	//用来存放用户想要存放的信息
	sourceArr map[int]interface{}
	mutex     *sync.Mutex
}

var accessControl = &AccessControl{sourceArr: make(map[int]interface{})}

// GetData 获取数据
func (a *AccessControl) GetData(uid int) interface{} {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.sourceArr[uid]
}

func (a *AccessControl) SetData(uid int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.sourceArr[uid] = "内容暂定"
}

func (a *AccessControl) getDistributeAddr(r *http.Request) bool {
	//获取用户UID
	uid, err := r.Cookie("uid")
	if err != nil {
		return false
	}

	//采用一致性hash算法，根据用户ID，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	//判断是否为本机
	if hostRequest == localHost {
		return a.getDataFromLocal(uid.Value)
	}
	return GetDataFromOther(hostRequest, r)
}

func (a *AccessControl) getDataFromLocal(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}

	data := a.GetData(uidInt)
	if data != nil {
		return true
	}
	return false
}

// GetCurUrl 优化转发请求代码
func GetCurUrl(hostUrl string, r *http.Request) (rsp *http.Response, body []byte, err error) {
	//获取Uid
	uidPre, err := r.Cookie("uid")
	if err != nil {
		return
	}
	//获取sign
	uidSign, err := r.Cookie("sign")
	if err != nil {
		return
	}

	// 模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}

	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	rsp, err = client.Do(req)
	defer rsp.Body.Close()

	body, err = ioutil.ReadAll(rsp.Body)
	return
}

func GetDataFromOther(host string, r *http.Request) bool {
	////获取Uid
	//uidPre, err := r.Cookie("uid")
	//if err != nil {
	//	return false
	//}
	////获取sign
	//uidSign, err := r.Cookie("sign")
	//if err != nil {
	//	return false
	//}
	//
	//client := http.Client{}
	//request, err := http.NewRequest("GET", "http://"+host+":"+port+"/check", nil)
	//if err != nil {
	//	return false
	//}
	//
	//cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	//cookieSign := &http.Cookie{Name: "uid", Value: uidSign.Value, Path: "/"}
	//// 添加新的cookie到新请求中
	//request.AddCookie(cookieUid)
	//request.AddCookie(cookieSign)
	//
	//// 发送请求
	//response, err := client.Do(request)
	//if err != nil {
	//	return false
	//}
	//
	//body, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	return false
	//}

	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurUrl(hostUrl, r)
	if err != nil {
		return false
	}

	// 判断状态
	if response.StatusCode == 200 {
		// 这里是为什么，body读取出来是true false呢
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

//func webHandler(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("正常执行业务函数"))
//}

func Auth(w http.ResponseWriter, r *http.Request) error {
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserInfo 身份校验
func CheckUserInfo(r *http.Request) error {
	// 获取uid和sign的cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户UID Cookie 获取失败！")
	}
	signCookie, err := r.Cookie("sign")
	// 获取uid和sign的cookie
	if err != nil {
		return errors.New("用户sign Cookie 获取失败！")
	}

	// signCookie进行解密
	signByte, err := encrypt.DecodeMess(signCookie.Value)
	if err != nil {
		return errors.New("用户加密串 Cookie 获取失败！")
	}

	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败！")
}

func checkInfo(uid, sign string) bool {
	if uid == sign {
		return true
	}
	return false
}

// 执行正常执行业务
func check(w http.ResponseWriter, r *http.Request) {
	// 从url中获取productId和从cookie中获取userId
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productIdString := queryForm["productID"][0]
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	// 1.分布式权限验证
	right := accessControl.getDistributeAddr(r)
	if right == false {
		w.Write([]byte("false"))
		return
	}

	// 2.获取数量控制，防止超卖
	hostUrl := "http://" + getOne + ":" + getOnePort + "/getOne"
	responseValidation, responseBody, err := GetCurUrl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
		// 待做：若存在err，还需要对getOne就行反向操作，即数量减一
		return
	}

	if responseValidation.StatusCode == 200 {
		if string(responseBody) == "true" {
			//	整合下单:生成消息，推送rabbitmq
			productID, err := strconv.ParseInt(productIdString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			// 创建消息体
			message := datamodels.Message{
				ProductId: productID,
				UserId:    userID,
			}
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			// 生产消息
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
}

/*拦截器*/
func main() {
	hashConsistent = common.NewConsistentHash(20)
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	// 实例化rabbitmq
	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("seckill")
	defer rabbitMqValidate.Destory()

	filter := common.NewFilter()
	filter.RegisterFilterHandler("/check", Auth)
	filter.RegisterFilterHandler("/checkRight", Auth)

	//http.HandleFunc("/check", filter.Handle(webHandler))
	http.HandleFunc("/check", filter.Handle(check))

	http.ListenAndServe("localhost:8083", nil)
}
