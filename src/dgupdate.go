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

func dgUpdateCheck(onFail func(), onUpdate func()) error {
	var updateData dgupdate
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	r, err := httpClient.Get("https://diskgem.info/update.json")
	if err != nil {
		onFail()
		return err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		onFail()
		return err
	}
	err = json.Unmarshal(body, &updateData)
	if err != nil {
		onFail()
		return err
	}
	if updateData.Latest > dgBuildNumber {
		onUpdate()
		return nil
	}
	return nil
}
