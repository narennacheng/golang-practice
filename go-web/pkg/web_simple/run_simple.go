package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func readBodyOnce(w http.ResponseWriter, r *http.Request) {
	// Body：只能读取一次，意味着你读了别人就不能读了；别人读了你就不能读了；
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed: %v", err)
		return
	}
	// 类型转化，[]byte转为string
	fmt.Fprintf(w, "read the data: %s \n", string(body))
	// 尝试再次读取，啥也读不到，但是也不会报错
	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read the data one more time got failed: %v", err)
		return
	}
	fmt.Fprintf(w, "read the data one more time: [%s] and read data lenght: %d \n", string(body), len(body))

}

func wholeUrl(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(r.URL)
	fmt.Fprintf(w, string(data))
}

func queryParams(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	fmt.Fprintf(w, "query is : %v\n", params)
}

func header(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "header is %v\n", r.Header)
}

func form(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "before parse form %v\n", r.Form)
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "parse form error %v\n", r.Form)
	} else {
		fmt.Fprintf(w, "parse form %v\n", r.Form)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home\n")
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "user\n")
	header(w, r)
	wholeUrl(w, r)
	queryParams(w, r)
	form(w, r)
	readBodyOnce(w, r)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/user", user)
	http.HandleFunc("/body/once", readBodyOnce)
	http.ListenAndServe(":8080", nil)
}
