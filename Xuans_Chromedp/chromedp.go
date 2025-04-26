package Xuans_Chromedp

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func XuansChromedp() {
	//config username and password
	username := "tomsmith"
	password := "SuperSecretPassword!"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("ignore-ssl-errors", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "site-per-process"),
		chromedp.Flag("disable-features", "site-per-process"),
		chromedp.Flag("disable-features", "site-per-process"),
		chromedp.Flag("disable-features", "site-per-process"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log.Printf("Starting ChromeDP")

	var LoginSuccessIndicator string

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://the-internet.herokuapp.com/login"),

		//wait for login visible
		chromedp.WaitVisible("#username"),

		//input username
		chromedp.SendKeys("#username", username),

		//input password
		chromedp.SendKeys("#password", password),

		chromedp.Sleep(2*time.Second),
		//click login button
		chromedp.Click(`button[type="submit"]`),

		chromedp.Text("h4[class='subheader']", &LoginSuccessIndicator),

		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		log.Printf("Failed to run chromedp: %v", err)
		return
	}
	log.Printf("Login Success: %s", LoginSuccessIndicator)
	log.Printf("Screenshot saved to screenshot.png")
	log.Printf("Wait 10 seconds to close browser")
	time.Sleep(10 * time.Second)
	log.Printf("Closing browser")
}
