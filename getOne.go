package main

/*秒杀数量控制接口：保证秒杀数量不超过库存*/

import (
	"log"
	"net/http"
	"sync"
)

// 商品库存(需要从数据库中获得)
var productNum int64 = 100

// 已抢购的数量
var orderedSum int64 = 0

// 互斥更新已抢购的数量
var mutex sync.Mutex

func getOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	// 判断数据是否超限
	if orderedSum < productNum {
		orderedSum++
		return true
	}
	return false
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	if getOneProduct() {
		w.Write([]byte("抢购成功"))
	}
	w.Write([]byte("抢购失败：商品已售罄"))
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe("：8084", nil)
	if err != nil {
		log.Fatal("Err:", err)
	}
}
