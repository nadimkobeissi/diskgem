/* SPDX-License-Identifier: MIT
 * Copyright © 2018 Nadim Kobeissi <nadim@symbolic.software>. All Rights Reserved.
 */

package main

import (
	"fmt"
	"os"
)

func dgErrorCritical(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", e)
		os.Exit(1)
	}
}
