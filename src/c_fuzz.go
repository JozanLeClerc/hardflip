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
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

func c_list_items(litems *ItemsList, stdin io.WriteCloser) {
	for ptr := litems.head; ptr != nil; ptr = ptr.next {
		var forebears []string

		if ptr.is_dir() == true {
			continue
		}
		for rptr := ptr.Host.parent; len(rptr.Name) > 0; rptr = rptr.Parent {
			forebears = append(forebears, rptr.Name)
		}
		for i := len(forebears) - 1; i >= 0; i-- {
			io.WriteString(stdin, forebears[i] + "/")
		}
		io.WriteString(stdin, ptr.Host.Name + "\n")
	}
}

func c_fuzz(data *HardData, ui *HardUI) {
	if err := ui.s.Suspend(); err != nil {
		c_error_mode("screen", err, ui)
		return
	}
	search := exec.Command("fzf")
	stdin, err := search.StdinPipe()
	if err != nil {
		c_error_mode("search stdin pipe", err, ui)
		c_resume_or_die(ui)
		return
	}
	stdout, err := search.StdoutPipe()
	if err != nil {
		c_error_mode("search stdout pipe", err, ui)
		c_resume_or_die(ui)
		return
	}
	if err := search.Start(); err != nil {
		c_error_mode("fzf", err, ui)
		c_resume_or_die(ui)
		return
	}
	go func() {
		defer stdin.Close()
		c_list_items(data.litems, stdin)
	}()
	output, err := io.ReadAll(stdout)
	if err != nil {
		ui.s.Fini()
		c_die("fuck it failed", err)
	}
	str_out := strings.TrimSuffix(string(output), "\n")
	fmt.Printf("[%s]\n", str_out)
	time.Sleep(3 * time.Second)
	c_resume_or_die(ui)
}
