package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var (
	exclude string
)

func init() {
	flag.StringVar(&exclude, "exclude", "快猫", "-exclude='a,b,c'")
}

func filter(data string) bool {
	fmt.Printf("DEBUG %s %s\n", data, exclude)
	for _, key := range strings.Split(exclude, ",") {
		if strings.Contains(data, key) {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	actx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		actx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 1200*time.Second)
	defer cancel()

	tasks := func(ctx context.Context) []chromedp.Action {
		f := func(ctx context.Context) error {
			var (
				err       error
				groupName string
				tpl             = `#content-pannel > div > div.content-pannel-body > group-list > div > ul > li:nth-child(%d)`
				idx, l    int64 = 1, 0
			)
		start:
			for {
				// 联系人
				err = chromedp.Click(`#menu-pannel > ul.main-menus > li.menu-item.menu-contact > div > i`, chromedp.NodeVisible).Do(ctx)
				if err != nil {
					return err
				}

				// 我的群组
				chromedp.Click(`#sub-menu-pannel > ul > li:nth-child(2) > div`, chromedp.NodeVisible).Do(ctx)
				if err != nil {
					return err
				}

				var nodes []*cdp.Node
				chromedp.Nodes(`#content-pannel > div > div.content-pannel-body > group-list > div`, &nodes, chromedp.ByQueryAll).Do(ctx)
				if len(nodes) > 0 {
					fmt.Printf("idx:%d, nodes len:%d , %d \n", idx, len(nodes), nodes[0].ChildNodeCount)
					l = nodes[0].ChildNodeCount
				}

				if l == 0 || idx > l {
					break
				}

				for {
					// 选中第i个群
					tag := fmt.Sprintf(tpl, idx)
					err = chromedp.Click(tag, chromedp.NodeVisible).Do(ctx)
					if err != nil {
						return err
					}
					// 检查群名称
					err = chromedp.OuterHTML(`#content-pannel > div > div.content-pannel-head.chat-head > div.conv-title > div.title > span > span`, &groupName, chromedp.NodeVisible).Do(ctx)
					if err != nil {
						return err
					}
					if !filter(groupName) {
						break
					}
					idx++
					tag = fmt.Sprintf(tpl, idx)
					continue start
				}

				checkErr := func(actionFunc chromedp.ActionFunc) {
					if err != nil {
						return
					}
					err = actionFunc(ctx)
				}
				// 群设置
				checkErr(chromedp.Click(`#content-pannel > div > div.content-pannel-head.chat-head > div.conv-operations > i.iconfont.icon-group-setting.ng-scope.tipper-attached`,
					chromedp.NodeVisible).Do)
				// 退出群聊
				checkErr(chromedp.Click(`#content-pannel > div.unpop-modal.group-setting-modal.ng-scope > div.foot.group-setting-footer.ng-scope`, chromedp.NodeVisible).Do)
				// 确认
				checkErr(chromedp.Click(`body > div.ding-modal.fade.ng-scope.ng-isolate-scope.in > div > div > div.foot > button.confirm-ok.ng-binding.blue`, chromedp.NodeVisible).Do)
				fmt.Printf("%s 已退出\n", groupName)
				// sleep
				checkErr(chromedp.Sleep(1 * time.Second).Do)
			}
			if err != nil {
				return err
			}

			return nil
		}
		return []chromedp.Action{
			chromedp.Navigate(`https://im.dingtalk.com/`),
			// 二维码
			chromedp.WaitVisible("#qrcode-login > div > div > img"),
			//  提示app下载, 浮层不生效，需要手工关掉也行
			chromedp.Click(`#header > client-download-guide > div > div > div.close > div`, chromedp.NodeVisible),
			chromedp.Click(`#menu-pannel > ul.main-menus > li.menu-item.menu-contact > div > i`, chromedp.NodeVisible),

			// 重复执行
			chromedp.ActionFunc(f),
		}
	}
	// navigate to a page, wait for an element, click
	var example string
	err := chromedp.Run(ctx, tasks(ctx)...)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\n%s", example)
}
