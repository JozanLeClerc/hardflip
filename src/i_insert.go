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
 * hardflip: src/i_insert.go
 * Thu Mar 14 10:18:40 2024
 * Joe
 *
 * insert a new host
 */

package main

import (
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func i_draw_text_box(ui HardUI, line int, dim Quad, label, content string,
					 id, selected int, red bool) {
	const tbox_size int = 14
	tbox_style := ui.style[DEF_STYLE].Background(tcell.ColorBlack).Dim(true)

	if id == selected {
		tbox_style = tbox_style.Reverse(true).Dim(false)
	}

	l := ui.dim[W] / 2 - len(label) - 2
	if l <= dim.L { l = dim.L + 1 }
	i_draw_text(ui.s, l, line, ui.dim[W] / 2, line,
		ui.style[DEF_STYLE], label)
	if (id == 4 || id == 9) && len(content) > 0 {
		content = "***"
	}
	if red == true {
		tbox_style = tbox_style.Foreground(tcell.ColorRed)
	}
	spaces := ""
	for i := 0; i < tbox_size; i++ {
		spaces += " "
	}
	i_draw_text(ui.s, ui.dim[W] / 2, line, dim.R, line,
		tbox_style,
		"[" + spaces + "]")
	i_draw_text(ui.s, ui.dim[W] / 2 + 1, line, ui.dim[W] / 2 + 1 + tbox_size,
		line, tbox_style, content)
}

func i_draw_insert_panel(ui HardUI, in *HostNode) {
	if len(in.Name) == 0 {
		return
	}
	win := Quad{
		ui.dim[W] / 8,
		ui.dim[H] / 8,
		ui.dim[W] - ui.dim[W] / 8 - 1,
		ui.dim[H] - ui.dim[H] / 8 - 1,
	}
	i_draw_box(ui.s, win.L, win.T, win.R, win.B,
		ui.style[BOX_STYLE], ui.style[HEAD_STYLE],
		" Insert - " + in.Name + " ", true)
	line := 2
	if win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "Connection type", in.protocol_str(),
		0, ui.insert_sel, false)
	line += 2
	switch in.Protocol {
	case 0:
		i_draw_insert_ssh(ui, line, win, in)
	}
}

func i_draw_insert_ssh(ui HardUI, line int, win Quad, in *HostNode) {
	red := false
	if win.T + line >= win.B { return }
	text := "---- Host settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "Host/IP", in.Host, 1, ui.insert_sel,
		false)
	if line += 1; win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "Port", strconv.Itoa(int(in.Port)),
		2, ui.insert_sel, false);
	if line += 2; win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "User", in.User, 3, ui.insert_sel,
		false)
	if line += 1; win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "Pass", in.Pass, 4, ui.insert_sel,
		false)
	if line += 1; win.T + line >= win.B { return }
	if len(in.Priv) > 0 {
		file := in.Priv
		if file[0] == '~' {
			home, _ := os.UserHomeDir()
			file = home + file[1:]
		}
		if stat, err := os.Stat(file);
		   err != nil || stat.IsDir() == true {
			red = true
		}
	}
	i_draw_text_box(ui, win.T + line, win, "SSH private key",
		in.Priv, 5, ui.insert_sel, red)
	if red == true {
		if line += 1; win.T + line >= win.B { return }
		text := "file does not exist"
		i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
			win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
	}
	red = false
	if line += 2; win.T + line >= win.B { return }
	text = "---- Jump settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "Host/IP",
		in.Jump.Host, 6, ui.insert_sel, false)
	if line += 1; win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "Port",
		strconv.Itoa(int(in.Jump.Port)), 7, ui.insert_sel, false)
	if line += 2; win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "User",
		in.Jump.User, 8, ui.insert_sel, false)
	if line += 1; win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "Pass",
		in.Jump.Pass, 9, ui.insert_sel, false)
	if line += 1; win.T + line >= win.B { return }
	if len(in.Jump.Priv) > 0 {
		file := in.Jump.Priv
		if file[0] == '~' {
			home, _ := os.UserHomeDir()
			file = home + file[1:]
		}
		if stat, err := os.Stat(file);
		   err != nil || stat.IsDir() == true {
			red = true
		}
	}
	i_draw_text_box(ui, win.T + line, win, "SSH private key",
		in.Jump.Priv, 10, ui.insert_sel, red)
	if red == true {
		if line += 1; win.T + line >= win.B { return }
		text := "file does not exist"
		i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
			win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
	}
	// TODO: here
}
