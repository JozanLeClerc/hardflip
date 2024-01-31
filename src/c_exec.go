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
 * hardflip: src/c_exec.go
 * Thu Jan 18 18:57:42 2024
 * Joe
 *
 * exec the command at some point
 */

package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func c_exec_cmd(cmd_fmt []string) {
	cmd := exec.Command(cmd_fmt[0], cmd_fmt[1:]...)

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

func c_format_ssh(host *HostNode, pass string) []string {
	cmd_fmt := []string{}
	if len(pass) > 0 {
		cmd_fmt = append(cmd_fmt, "sshpass", "-p", pass)
	}

	cmd_fmt = append(cmd_fmt, "ssh")
	if len(host.Priv) > 0 {
		cmd_fmt = append(cmd_fmt, "-i", host.Priv)
	}
	if len(host.Jump) > 0 {
		cmd_fmt = append(cmd_fmt, c_format_ssh_jump(host))
	}
	if host.Port != 0 {
		cmd_fmt = append(cmd_fmt, "-p", strconv.Itoa(int(host.Port)))
	}
	cmd_fmt = append(cmd_fmt, host.User + "@" + host.Host)
	return cmd_fmt
}

func c_format_rdp(host *HostNode, pass string) []string {
	cmd_fmt := []string{"xfreerdp"}

	cmd_fmt = append(cmd_fmt,
		"/v:" + host.Host,
		"/u:" + host.User)
	if len(host.Domain) > 0 {
		cmd_fmt = append(cmd_fmt, "/d:" + host.Domain)
	}
	if len(pass) > 0 {
		cmd_fmt = append(cmd_fmt, "/p:" + pass)
	}
	if host.Port != 0 {
		cmd_fmt = append(cmd_fmt, "/port:" + strconv.Itoa(int(host.Port)))
	}
	if host.Dynamic == true {
		cmd_fmt = append(cmd_fmt, "/dynamic-resolution")
	}
	if host.Quality == 0 {
		cmd_fmt = append(cmd_fmt,
			"-aero", "-menu-anims", "-window-drag", "-wallpaper",
			"-decorations", "-fonts", "-themes",
			"/bpp:8", "/compression-level:2")
	} else if host.Quality == 1 {
	} else {
		cmd_fmt = append(cmd_fmt,
			"+aero", "+menu-anims", "+window-drag",
			"+decorations", "+fonts", "+themes", "/gfx:RFX", "/rfx", "/gdi:hw",
			"/bpp:32")
	}
	cmd_fmt = append(cmd_fmt,
		"/size:" + strconv.Itoa(int(host.Width)) +
		"x" + strconv.Itoa(int(host.Height)))
	return cmd_fmt
}

func c_format_cmd(host *HostNode, opts HardOpts, ui *HardUI) []string {
	var cmd_fmt []string
	var pass string
	gpg, term := opts.GPG, opts.Term

	if len(gpg) > 0 && gpg != "plain" && len(host.Pass) > 0 {
		var err error
		pass, err = c_decrypt_str(host.Pass)
		if err != nil {
			c_error_mode(host.Parent.path() + host.Filename +
				": password decryption failed", err, ui)
			return nil
		}
		pass = strings.TrimSuffix(pass, "\n")
	}
	switch host.Protocol {
	case 0:
		cmd_fmt = c_format_ssh(host, pass)
	case 1:
		cmd_fmt = c_format_rdp(host, pass)
	default:
		c_die("you fucked up joe, users cant see this", nil)
	}
	if len(term) > 0 {
		// TODO: setsid
		if term == "$TERMINAL" {
			term = os.Getenv("TERMINAL")
		}
		cmd_fmt = append([]string{"setsid", term, "-e"}, cmd_fmt...)
	}
	return cmd_fmt
}

func c_exec(host *HostNode, opts HardOpts, ui *HardUI) {
	var err error

	if host == nil {
		return
	}
	cmd_fmt := c_format_cmd(host, opts, ui)
	if cmd_fmt == nil {
		return
	}
	ui.s.Fini()
	c_exec_cmd(cmd_fmt)
	if opts.Loop == false {
		os.Exit(0)
	} else {
		if ui.s, err = tcell.NewScreen(); err != nil {
			c_die("view", err)
		}
		if err := ui.s.Init(); err != nil {
			c_die("view", err)
		}
		ui.s.SetStyle(ui.style[DEF_STYLE])
	}
}
