package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"seckill/common"
	"strconv"
	"sync"
)

/*针对不同请求做负载均衡*/

//设置集群地址，最好内外IP
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

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

func GetDataFromOther(host string, r *http.Request) bool {
	//获取Uid
	uidPre, err := r.Cookie("uid")
	if err != nil {
		return false
	}
	//获取sign
	uidSign, err := r.Cookie("sign")
	if err != nil {
		return false
	}

	client := http.Client{}
	request, err := http.NewRequest("GET", "http://"+host+":"+port+"/check", nil)
	if err != nil {
		return false
	}

	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "uid", Value: uidSign.Value, Path: "/"}
	// 添加新的cookie到新请求中
	request.AddCookie(cookieUid)
	request.AddCookie(cookieSign)

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		return false
	}

	body, err := ioutil.ReadAll(response.Body)
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

func webHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("正常执行业务函数"))
}

func filterHandler(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("执行拦截器函数"))
	return errors.New("我是你大爹！！！")
}

/*拦截器*/
func main() {
	hashConsistent = common.NewConsistentHash(20)
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	filter := common.NewFilter()
	filter.RegisterFilterHandler("/check", filterHandler)

	http.HandleFunc("/check", filter.Handle(webHandler))
	http.ListenAndServe("localhost:8083", nil)
}
