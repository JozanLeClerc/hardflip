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
 * Thu Jan 18 16:23:10 2024
 * Joe
 *
 * init functions
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type HardOpts struct {
	Icon    bool
	Loop    bool
	GPG		string
	Perc    bool
	Term	string
}

// this function recurses into the specified root directory in order to load
// every yaml file into memory
func c_recurse_data_dir(dir, root string, opts HardOpts,
		ldirs *DirsList,
		name string, parent *DirsNode, depth uint16,
		ui *HardUI, load_err *[]error) {
	files, err := os.ReadDir(root + dir)
	if err != nil {
		*load_err = append(*load_err, err)
		return
	}
	dir_node := DirsNode{
		name,
		parent,
		depth,
		&HostList{},
		nil,
	}
	ldirs.add_back(&dir_node)
	i_draw_load_ui(ui)
	for _, file := range files {
		filename := file.Name()
		if file.IsDir() == true {
			c_recurse_data_dir(dir + filename + "/", root, opts, ldirs,
				file.Name(), &dir_node, depth + 1, ui, load_err)
		} else if filepath.Ext(filename) == ".yml" {
			host_node, err := c_read_yaml_file(root + dir + filename, opts, ui)
			if err != nil {
				*load_err = append(*load_err, err)
			} else if host_node != nil {
				host_node.Filename = filename
				host_node.Parent = &dir_node
				if len(opts.GPG) == 0 {
					host_node.Pass = ""
				} else if opts.GPG != "plain" && len(host_node.Pass) > 0 {
					host_node.Pass, err = c_decrypt_str(host_node.Pass)
					if err != nil {
						str := fmt.Sprintf("%s%s: password decryption: %v\n",
							dir, filename, err)
						*load_err = append(*load_err, errors.New(str))
					}
					host_node.Pass = strings.TrimSuffix(host_node.Pass, "\n")
				}
				dir_node.lhost.add_back(host_node)
			}
			i_draw_load_ui(ui)
		}
	}
}

func c_load_data_dir(dir string, opts HardOpts,
		ui *HardUI) (*DirsList, []error) {
	ldirs := DirsList{}
	var load_err []error

	c_recurse_data_dir("", dir + "/", opts, &ldirs, "", nil, 1, ui, &load_err)
	return &ldirs, load_err
}

// fills litems sorting with dirs last
// other sorting algos are concievable
// this func also sets the root folder to unfolded as it may never be folded
// this func also sets the default litems.curr
func c_load_litems(ldirs *DirsList) *ItemsList {
	litems := ItemsList{}

	for ptr := ldirs.head; ptr != nil; ptr = ptr.next {
		item := ItemsNode{ Dirs: ptr, Host: nil }
		litems.add_back(&item)
		for ptr := ptr.lhost.head; ptr != nil; ptr = ptr.next {
			item := ItemsNode{ Dirs: nil, Host: ptr }
			litems.add_back(&item)
		}
	}
	litems.head = litems.head.next
	if litems.head != nil {
		litems.head.prev = nil
	}
	litems.curr = litems.head
	litems.draw = litems.head
	return &litems
}
