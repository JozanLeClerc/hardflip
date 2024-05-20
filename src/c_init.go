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
 * Fri Feb 02 10:09:18 2024
 * Joe
 *
 * init functions
 */

package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type HardOpts struct {
	Icon    bool   `yaml:"icons"`
	Loop    bool   `yaml:"loop"`
	GPG		string `yaml:"gpg"`
	Perc    bool   `yaml:"percent"`
	Term	string `yaml:"terminal"`
	DefSSH	string `yaml:"default_ssh_priv"`
	file    string
}

type HardStyle struct {
	DefColor	string `yaml:"default"`
	DirColor	string `yaml:"dir_color"`
	BoxColor	string `yaml:"box_color"`
	HeadColor	string `yaml:"head_color"`
	ErrColor	string `yaml:"error_color"`
	TitleColor	string `yaml:"title_color"`
	BotColor	string `yaml:"bottom_color"`
	YankColor	string `yaml:"yank_color"`
	MoveColor	string `yaml:"move_color"`
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
			host_node, err := c_read_yaml_file(root + dir + filename)
			if err != nil {
				*load_err = append(*load_err, err)
			} else if host_node != nil {
				host_node.filename = filename
				host_node.parent = &dir_node
				if len(opts.GPG) == 0 {
					host_node.Pass = ""
				}
				dir_node.lhost.add_back(host_node)
			}
			i_draw_load_ui(ui)
		}
	}
}

func c_load_data_dir(dir string, opts HardOpts,
		ui *HardUI, load_err *[]error) (*DirsList) {
	ldirs := DirsList{}

	c_recurse_data_dir("", dir + "/", opts, &ldirs, "", nil, 1, ui, load_err)
	return &ldirs
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

func c_write_options(file string, opts HardOpts, load_err *[]error) {
	data, err := yaml.Marshal(opts)
	if err != nil {
		*load_err = append(*load_err, err)
		return
	}
	err = os.WriteFile(file, data, 0644)
	if err != nil {
		*load_err = append(*load_err, err)
	}
}

func c_write_styles(file string, opts HardStyle, load_err *[]error) {
	data, err := yaml.Marshal(opts)
	if err != nil {
		*load_err = append(*load_err, err)
		return
	}
	err = os.WriteFile(file, data, 0644)
	if err != nil {
		*load_err = append(*load_err, err)
	}
}

func c_get_options(dir string, load_err *[]error) HardOpts {
	opts := DEFAULT_OPTS
	file := dir + "/" + CONF_FILE_NAME

	if _, err := os.Stat(file); os.IsNotExist(err) {
		c_write_options(file, DEFAULT_OPTS, load_err)
		opts.file = file
		return opts
	}
	opts, err := c_parse_opts(file)
	opts.file = file
	if err != nil {
		*load_err = append(*load_err, err)
		return opts
	}
	return opts
}

func c_get_styles(dir string, load_err *[]error) HardStyle {
	styles := HardStyle{}
	file := dir + "/" + STYLE_FILE_NAME

	if _, err := os.Stat(file); os.IsNotExist(err) {
		c_write_styles(file, DEFAULT_STYLE, load_err)
		return DEFAULT_STYLE
	}
	styles, err := c_parse_styles(file)
	if err != nil {
		*load_err = append(*load_err, err)
		return DEFAULT_STYLE
	}
	return styles
}
