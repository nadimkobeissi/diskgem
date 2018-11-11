/* SPDX-License-Identifier: MIT
 * Copyright © 2018 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type dgupdate struct {
	Latest      int
	Date        string
	Description string
	Critical    bool
}

func dgUpdateCheck() int {
	var updateData dgupdate
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	r, err := httpClient.Get("https://diskgem.info/update.json")
	if err != nil {
		return 0
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0
	}
	err = json.Unmarshal(body, &updateData)
	if err != nil {
		return 0
	}
	if updateData.Latest > dgBuildNumber {
		return 1
	}
	return 2
}
