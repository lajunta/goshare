package main

import (
	"github.com/lxn/walk"
	"os/exec"
)

func notifyicon() {
	walk.SetPanicOnError(true)
	mw, _ := walk.NewMainWindow()
	icon, _ := walk.NewIconFromFile("x.ico")
	ni, _ := walk.NewNotifyIcon()
	defer ni.Dispose()

	ni.SetIcon(icon)
	ni.SetToolTip("文件上传与下载")

	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		ni.ShowCustom("文件上传与下载", "Go语言开发")
	})

	// We put an exit action into the context menu.
	exit := walk.NewAction()
	browser := walk.NewAction()

	exit.SetText("退出")
	browser.SetText("打开浏览器")

	exit.Triggered().Attach(func() { walk.App().Exit(0) })
	browser.Triggered().Attach(func() {

		cmd := exec.Command(ie, localhost)
		cmd.Start()
	})

	ni.ContextMenu().Actions().Add(browser)
	ni.ContextMenu().Actions().Add(exit)
	// The notify icon is hidden initially, so we have to make it visible.
	ni.SetVisible(true)

	// Now that the icon is visible, we can bring up an info balloon.
	//ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again.")

	mw.Run()

}
