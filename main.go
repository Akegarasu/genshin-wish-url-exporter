package main

import (
	"fmt"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

const (
	ignoreUrl = "localhost;127.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;192.168.*;<local>"
)

var ch = make(chan struct{})

type TestAddon struct {
	proxy.BaseAddon
}

func (t *TestAddon) Request(f *proxy.Flow) {
	reg := regexp.MustCompile(`webstatic([^\.]{2,10})?\.(mihoyo|hoyoverse)\.com`)
	reg1 := regexp.MustCompile(`authkey=[^&]+`)
	if len(reg.FindAllString(f.Request.URL.Host, -1)) > 0 {
		if len(reg1.FindAllString(f.Request.URL.RawQuery, -1)) > 0 {
			url := fmt.Sprintf("%s#/log\n", f.Request.URL.String())
			fmt.Printf("------------------已成功抓取到URL----------------------\n%s\n------------------------------------------------------\n", url)
			fmt.Println("请复制上方链接（选中后右键）")
			fmt.Println("同时，该链接已经被写入到程序目录下的 抽卡记录链接.txt 中")
			err := os.WriteFile("抽卡记录链接.txt", []byte(url), 0644)
			if err != nil {
				fmt.Println("写入文件失败")
			}
			ch <- struct{}{}
		}
	}
}

func main() {
	log.SetLevel(log.FatalLevel)
	fmt.Println("程序已启动，开始监听抽卡记录链接。请打开游戏内抽卡记录")
	err := SetProxy(true, "127.0.0.1:18191", ignoreUrl)
	if err != nil {
		log.Error(err)
		exit()
	}
	opts := &proxy.Options{
		Addr:              ":18191",
		StreamLargeBodies: 1024 * 1024 * 5,
		SslInsecure:       true,
	}
	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Error(err)
		exit()
	}
	p.AddAddon(&TestAddon{})
	go func() {
		log.Fatal(p.Start())
	}()
	for {
		select {
		case <-ch:
			fmt.Println("正在取消代理...")
			err = SetProxy(false, "", "")
			exit()
			fmt.Println("正在关闭....")
			_ = p.Close()
		}
	}
}

func exit() {
	var s string
	fmt.Println("按回车退出")
	_, _ = fmt.Scanln(&s)
}
