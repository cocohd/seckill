package common

import "net/http"

/*拦截器*/

// FilterHandler 声明一个新的函数类型
type FilterHandler func(w http.ResponseWriter, r *http.Request) error

// Filter 拦截器数据结构
type Filter struct {
	filterMap map[string]FilterHandler
}

// NewFilter 构造函数
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandler)}
}

// RegisterFilterHandler 注册拦截器
func (f *Filter) RegisterFilterHandler(path string, handler FilterHandler) {
	f.filterMap[path] = handler
}

// GetFilterHandler 获取拦截器
func (f *Filter) GetFilterHandler(path string) FilterHandler {
	return f.filterMap[path]
}

type webHandler func(w http.ResponseWriter, r *http.Request)

// Handle 执行拦截器
func (f *Filter) Handle(webHandler webHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 执行拦截handler
		for path, filterHandler := range f.filterMap {
			if path == r.RequestURI {
				filterHandler(w, r)
				return
			}
		}

		// 执行业务函数
		webHandler(w, r)
	}
}
