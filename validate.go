package main

import (
	"errors"
	"net/http"
	"seckill/common"
)

func webHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("正常执行业务函数"))
}

func filterHandler(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("执行拦截器函数"))
	return errors.New("我是你大爹！！！")
}

/*拦截器*/
func main() {
	filter := common.NewFilter()
	filter.RegisterFilterHandler("/check", filterHandler)

	http.HandleFunc("/check", filter.Handle(webHandler))
	http.ListenAndServe("localhost:8083", nil)
}
