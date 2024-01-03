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
 * Copyright (c) 2023-2024, Joe
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the organization nor the names of its
 *    contributors may be used to endorse or promote products derived from
 *    this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS ''AS IS''
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * hardflip: src/c_init.go
 * Tue Dec 26 16:32:26 2023
 * Joe
 *
 * init functions
 */

package main

import (
	"os"
	"path/filepath"
)

type HardOpts struct {
	Icon    bool
	Loop    bool
	FoldAll bool
}


var global_id uint64

// this function recurses into the specified root directory in order to load
// every yaml file into memory
func c_recurse_data_dir(dir, root string, opts HardOpts,
		litems *ItemsList, ldirs *DirsList,
		name string, parent *DirsNode, depth uint16) {
	files, err := os.ReadDir(root + dir)
	if err != nil {
		c_die("could not read data directory", err)
	}
	dir_node := DirsNode{
		0,
		name,
		parent,
		depth,
		&HostList{},
		opts.FoldAll,
		nil,
	}
	// item_node := ItemsNode{}
	// item_node.Dirs = &dir_node
	// item_node.Host = nil
	ldirs.add_back(&dir_node)
	global_id++
	// litems.add_back(&item_node)
	for _, file := range files {
		filename := file.Name()
		if file.IsDir() == true {
			c_recurse_data_dir(dir + filename + "/", root, opts, litems, ldirs,
				file.Name(), &dir_node, depth + 1)
		} else if filepath.Ext(filename) == ".yml" {
			host_node := c_read_yaml_file(root + dir + filename)
			if host_node == nil {
				return
			}
			// item_node := ItemsNode{}
			// item_node.Dirs = nil
			// item_node.Host = host_node
			// litems.add_back(&item_node)
			host_node.Filename = filename
			host_node.Parent = &dir_node
			dir_node.lhost.add_back(host_node)
			global_id++
		}
	}
}

func c_load_data_dir(dir string, opts HardOpts) (*ItemsList, *DirsList) {
	litems := ItemsList{}
	ldirs  := DirsList{}

	c_recurse_data_dir("", dir + "/", opts, &litems, &ldirs, "", nil, 1)
	return &litems, &ldirs
}
