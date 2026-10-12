package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	c := http.Client{Timeout: 1 * time.Second}
	_, err := c.Get("http://127.0.0.1:9222/json/version")

	var allocCtx context.Context
	var cancel context.CancelFunc

	if err == nil {
		fmt.Println("✅ 9222 ketemu, numpang tab")
		allocCtx, cancel = chromedp.NewRemoteAllocator(context.Background(), "http://127.0.0.1:9222")
	} else {
		fmt.Println("⚠️ 9222 gak ada, buka Brave baru")
		allocCtx, cancel = chromedp.NewExecAllocator(context.Background(),
			chromedp.ExecPath("/usr/bin/brave-browser-stable"),
		)
	}
	defer cancel()
	ctx, cancel2 := chromedp.NewContext(allocCtx)
	defer cancel2()

	if err := chromedp.Run(ctx, chromedp.Navigate("https://www.facebook.com")); err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= 10; i++ {
		chromedp.Run(ctx, chromedp.Evaluate(`window.scrollBy(0,500)`, nil), chromedp.Sleep(1*time.Second))
		fmt.Printf("scroll %d done\n", i)
	}
}
