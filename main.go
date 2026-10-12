package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	url := "https://facebook.com"
	if len(os.Args) > 1 && os.Args[1] == "youtube" {
		url = "https://youtube.com"
	}

	// Tepi cek dulu ya, Brave-nya udah "siap diajak ngobrol" (remote debugging) belum
	var allocCtx context.Context
	var cancel1 context.CancelFunc

	if isRemoteDebuggingReady() {
		// Kalo udah ada yang buka pake port 9222, kita lari ke sana
		allocCtx, cancel1 = chromedp.NewRemoteAllocator(context.Background(), "http://127.0.0.1:9222")
		fmt.Println("Siaap, Tepi numpang di Brave yang udah ada ya...")
	} else {
		// Kalo belum, kita buka baru pake flag khusus biar bisa Tepi stir
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath("/usr/bin/brave-browser-stable"),
			chromedp.Flag("remote-debugging-port", "9222"),
			chromedp.Flag("no-first-run", true),
			chromedp.Flag("no-default-browser-check", true),
			// chromedp.UserDataDir("/home/zakaria/.config/BraveSoftware/Brave-Browser"), // Opsional kalo mau login tetap ada
		)
		allocCtx, cancel1 = chromedp.NewExecAllocator(context.Background(), opts...)
		fmt.Println("Brave-nya belum ada, Tepi bukain yang baru ya sayang...")
	}
	defer cancel1()

	// Buat context baru (bakal jadi tab baru kalo di remote)
	ctx, cancel2 := chromedp.NewContext(allocCtx)
	defer cancel2()

	// Timeout biar nggak bengong kelamaan
	ctx, cancel3 := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel3()

	fmt.Printf("Otw meluncur ke %s...\n", url)

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		log.Fatalf("Aduh maaf Mas Irza, gagal navigate: %v", err)
	}

	for i := 1; i <= 10; i++ {
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`window.scrollBy(0, 500)`, nil),
			chromedp.Sleep(1*time.Second),
		)

		if err != nil {
			log.Printf("Lagi asik scroll malah error di scroll ke-%d: %v", i, err)
			break
		}

		fmt.Printf("Scroll ke-%d nih sayang...\n", i)
	}

	fmt.Println("Dah beres scroll-nya, Tepi tungguin 10 detik lagi ya sebelum pamit.")
	chromedp.Run(ctx, chromedp.Sleep(10*time.Second))
}

// Fungsi buat ngecek apakah Brave udah running dengan remote debugging di port 9222
func isRemoteDebuggingReady() bool {
	resp, err := http.Get("http://127.0.0.1:9222/json/version")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
