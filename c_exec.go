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
 * hardflip: src/c_exec.go
 * Mon, 18 Dec 2023 11:32:40 +0100
 * Joe
 *
 * exec the command at some point
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func exec_cmd(cmd_fmt []string) {
	cmd := exec.Command(cmd_fmt[0], cmd_fmt[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func format_cmd(id uint64, lhost *HostList) {
	curr := lhost.head
	var cmd_fmt []string

	curr = lhost.sel(id)
	if curr == nil {
		c_die("host id not found", nil)
	}
	cmd_fmt = []string{"ssh", "-i", curr.Priv, "-p", strconv.Itoa(int(curr.Port)), curr.User + "@" + curr.Host}
	fmt.Println(cmd_fmt)
	exec_cmd(cmd_fmt)
}

func display_servers(lhost *HostList) {
	curr := lhost.head

	if lhost.head == nil {
		fmt.Println("no hosts")
		return
	}
	for curr != nil {
		fmt.Println(curr.ID, curr.Folder + curr.Name)
		curr = curr.next
	}
	fmt.Println()
	format_cmd(3, lhost)
}
