/*
 * ========================
 * =====    ===============
 * ======  ================
 * ======  ================
 * ======  ====   ====   ==
 * ======  ===     ==  =  =
 * ======  ===  =  ==     =
 * =  ===  ===  =  ==  ====
 * =  ===  ===  =  ==  =  =
 * ==     =====   ====   ==
 * ========================
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright (c) 2023 Joe
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. Neither the name of the organization nor the
 *    names of its contributors may be used to endorse or promote products
 *    derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY JOE ''AS IS'' AND ANY
 * EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL JOE BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * josh: src/c_init.go
 * Thu, 14 Dec 2023 14:02:17 +0100
 * Joe
 *
 * init functions
 */

package main

import (
	"fmt"
	"os"
)

// This function will go get the data folder and try to create it if it does
// not exist. The first path being checked is $XDG_DATA_HOME then
// $HOME/.local/share. It returns the full data directory path.
func c_get_data_dir() string {
	home := os.Getenv("HOME")
	xdg_home := os.Getenv("XDG_DATA_HOME")

	if len(home) == 0 {
		c_die("env variable HOME not defined", nil)
	}
	if len(xdg_home) > 0 {
		if _, err := os.Stat(xdg_home); os.IsNotExist(err) {
			if err := os.MkdirAll(xdg_home, os.ModePerm); err != nil {
				c_die("could not create path " + xdg_home, err)
			}
			fmt.Println("created folder path " + xdg_home)
		}
		return xdg_home
	} else {
		home := home + ".local/share"
		if _, err := os.Stat(home); os.IsNotExist(err) {
			if err := os.MkdirAll(home, os.ModePerm); err != nil {
				c_die("could not create path " + home, err)
			}
			fmt.Println("created folder path " + home)
		}
		return home
	}
}

