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
 * Mon, 18 Dec 2023 15:07:52 +0100
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

func c_exec_cmd(cmd_fmt []string) {
	cmd := exec.Command(cmd_fmt[0], cmd_fmt[1:]...)

	fmt.Println(cmd_fmt)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func c_format_ssh_jump(host *HostNode) string {
		jump_fmt := "-oProxyCommand=ssh"
		if len(host.JumpPriv) > 0 {
			jump_fmt += " -i " + host.JumpPriv
		}
		if host.JumpPort != 0 {
			jump_fmt += " -p " + strconv.Itoa(int(host.JumpPort))
		}
		if len(host.JumpUser) == 0 {
			jump_fmt += " root"
		} else {
			jump_fmt += " " + host.JumpUser
		}
		jump_fmt += "@" + host.Jump + " -W %h:%p"
		return jump_fmt
}

func c_format_ssh(host *HostNode) []string {
	cmd_fmt := []string{"ssh"}
	user := host.User

	if len(host.Priv) > 0 {
		cmd_fmt = append(cmd_fmt, "-i", host.Priv)
	}
	if len(host.Jump) > 0 {
		cmd_fmt = append(cmd_fmt, c_format_ssh_jump(host))
	}
	if host.Port != 0 {
		cmd_fmt = append(cmd_fmt, "-p", strconv.Itoa(int(host.Port)))
	}
	if len(host.User) == 0 {
		user = "root"
	}
	cmd_fmt = append(cmd_fmt, user + "@" + host.Host)
	return cmd_fmt
}

func c_format_rdp(host *HostNode) []string {
	return []string{""}
}

func c_format_cmd(id uint64, lhost *HostList) {
	host := lhost.head
	var cmd_fmt []string

	host = lhost.sel(id)
	if host == nil {
		c_die("host id not found", nil)
	}
	if host.Type == 0 {
		cmd_fmt = c_format_ssh(host)
	} else if host.Type == 1 { 
		cmd_fmt = c_format_rdp(host)
	} else if host.Type > 1 {
		c_die("type not found", nil)
	}
	c_exec_cmd(cmd_fmt)
}

func c_exec(id uint64, lhost *HostList) {
	if lhost.head == nil {
		fmt.Println("no hosts")
		return
	}
	c_format_cmd(id, lhost)
}
