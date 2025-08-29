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
 * hardflip: src/c_fuzz.go
 * Wed, 27 Aug 2025 13:43:16 +0200
 * Joe
 *
 * search with fzf otherwise shitty search
 */

package main

import (
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

func c_fuzz_init_pipes(ui *HardUI,
	search *exec.Cmd) (io.WriteCloser, io.ReadCloser) {
	stdin, err := search.StdinPipe()
	if err != nil {
		c_error_mode("search stdin pipe", err, ui)
		c_resume_or_die(ui)
		return nil, nil
	}
	stdout, err := search.StdoutPipe()
	if err != nil {
		c_error_mode("search stdout pipe", err, ui)
		c_resume_or_die(ui)
		return nil, nil
	}
	return stdin, stdout
}

func c_fuzz_find_item(str_out string, litems *ItemsList) (*ItemsNode) {
	var ptr *ItemsNode

	path, name := filepath.Split(str_out)
	path = "/" + path
	for ptr = litems.head; ptr != nil; ptr = ptr.next {
		if ptr.is_dir() == true {
			continue
		}
		if name == ptr.Host.Name && path == ptr.path() {
			return ptr
		}
	}
	return nil
}

func c_fuzz(data *HardData, ui *HardUI) (bool) {
	if err := ui.s.Suspend(); err != nil && ui.s != nil {
		c_error_mode("screen", err, ui)
		return false
	}
	search := exec.Command("fzf")
	stdin, stdout := c_fuzz_init_pipes(ui, search)
	if stdin == nil || stdout == nil {
		return false
	}
	if err := search.Start(); err != nil {
		if ui != nil {
			c_error_mode("fzf", err, ui)
			c_resume_or_die(ui)
			return false
		} else {
			c_die("fzf", err)
		}
	}
	go func() {
		defer stdin.Close()
		for ptr := data.litems.head; ptr != nil; ptr = ptr.next {
			if ptr.is_dir() == true {
				continue
			}
			io.WriteString(stdin, ptr.path()[1:] + ptr.Host.Name + "\n")
		}
	}()
	output, err := io.ReadAll(stdout)
	if err != nil {
		if ui.s != nil {
			ui.s.Fini()
		}
		c_die("search stdout", err)
	}
	str_out := strings.TrimSuffix(string(output), "\n")
	if ui.s != nil {
		c_resume_or_die(ui)
	}
	if len(str_out) > 0 {
		item := c_fuzz_find_item(str_out, data.litems)
		if item == nil {
			if ui.s != nil {
				c_error_mode("item not found", nil, ui)
			}
			return false
		}
		data.litems.curr = item
	}
	return true
}
