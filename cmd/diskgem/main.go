/* SPDX-License-Identifier: MIT
 * Copyright © 2018-2019 Nadim Kobeissi <nadim@symbolic.software>. All Rights Reserved.
 */

package main

import (
	"time"

	"github.com/jroimartin/gocui"
)

var dgVersion = "1.4.2"
var dgBuildNumber = 7

type dgticker struct {
	gears  *time.Ticker
	active bool
}

func main() {
	dgState = dgStatePrototype
	dgStateReset(false)
	ui, err := gocui.NewGui(gocui.Output256)
	dgErrorCritical(err)
	defer ui.Close()
	ui.SetManagerFunc(uiMainManagerLayout)
	err = ui.MainLoop()
	dgErrorCritical(err)
}
