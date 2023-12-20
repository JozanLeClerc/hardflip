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
 * hardflip: src/c_init.go
 * Wed Dec 20 10:47:33 2023
 * Joe
 *
 * init functions
 */

package main

import (
	"os"
	"path/filepath"
)

type Opts struct {
	icon bool
	loop bool
}

// this function recurses into the specified root directory in order to load
// every yaml file into memory
func c_recurse_data_dir(dir string, root string, lhost *HostList) {
	files, err := os.ReadDir(root + dir)
	if err != nil {
		c_die("could not read data directory", err)
	}
	for _, file := range files {
		if file.IsDir() == true {
			c_recurse_data_dir(dir + file.Name() + "/", root, lhost)
		} else if filepath.Ext(file.Name()) == ".yml" {
			host := c_read_yaml_file(root + dir + file.Name())
			if host == nil {
				return
			}
			host.Filename = file.Name()
			host.Folder = dir
			lhost.add_back(host)
		}
	}
}

func c_load_data_dir(dir string) *HostList {
	lhost := HostList{}
	c_recurse_data_dir("", dir + "/", &lhost)
	return &lhost
}
