// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"context"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)


func main() {


	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// 1) specify the proxy server.
		// Note that the username/password is not provided here.
		// Check the link below for the description of the proxy settings:
		// https://www.chromium.org/developers/design-documents/network-settings
		chromedp.ProxyServer("http://proxy.h3c.com:8080"),
		// By default, Chrome will bypass localhost.
		// The test server is bound to localhost, so we should add the
		// following flag to use the proxy for localhost URLs.
		chromedp.Flag("proxy-bypass-list", "<-loopback>"),
	)

	// create chrome instance
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// log the protocol messages to understand how it works.
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	lctx, lcancel := context.WithCancel(ctx)
	chromedp.ListenTarget(lctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *fetch.EventRequestPaused:
			go func() {
				_ = chromedp.Run(ctx, fetch.ContinueRequest(ev.RequestID))
			}()
		case *fetch.EventAuthRequired:
			if ev.AuthChallenge.Source == fetch.AuthChallengeSourceProxy {
				go func() {
					_ = chromedp.Run(ctx,
						fetch.ContinueWithAuth(ev.RequestID, &fetch.AuthChallengeResponse{
							Response: fetch.AuthChallengeResponseResponseProvideCredentials,
							Username: "s21723",
							Password: "Sxn@20212",
						}),
						// Chrome will remember the credential for the current instance,
						// so we can disable the fetch domain once credential is provided.
						// Please file an issue if Chrome does not work in this way.
						fetch.Disable(),
					)
					// and cancel the event handler too.
					lcancel()
				}()
			}
		}
	})
	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),

		chromedp.Navigate(`https://bscscan.com/address/0x67599afad05ee1fcdc93a177e40644d08ea80854#readContract`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`//div[@id="readContract"]`),

		chromedp.Text(`//div[@aria-labelledby="readHeading8"]`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\n%s", strings.Split(example, " ")[0])
}