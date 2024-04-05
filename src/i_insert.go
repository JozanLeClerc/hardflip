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
 * Fri Apr 05 14:19:38 2024
 * Joe
 *
 * insert a new host
 */

package main

import (
	"errors"
	// "fmt"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func i_insert_check_ok(data *HardData, insert *HostNode) {
	if len(insert.Name) == 0 {
		data.insert_err = append(data.insert_err, errors.New("no name"))
	}
	if len(insert.Host) == 0 {
		data.insert_err = append(data.insert_err, errors.New("no host"))
	}
	if insert.Port == 0 {
		data.insert_err = append(data.insert_err, errors.New("port can't be 0"))
	}
	if insert.Protocol == PROTOCOL_SSH && len(insert.Priv) != 0 {
		file := insert.Priv
		if file[0] == '~' {
			home_dir, err := os.UserHomeDir()
			if err != nil {
				return
			}
			file = home_dir + file[1:]
		}
		if stat, err := os.Stat(file);
		err != nil {
			data.insert_err = append(data.insert_err, errors.New(file +
				": file does not exist"))
		} else if stat.IsDir() == true {
			data.insert_err = append(data.insert_err, errors.New(file +
				": file is a directory"))
		}
	}
}

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

func i_draw_ok_butt(ui HardUI, line int, id, selected int) {
	const butt_size int = 10
	const txt string = "ok"
	style := ui.style[DEF_STYLE].Background(tcell.ColorBlack).Dim(true)

	if id == selected {
		style = style.Reverse(true).Dim(false)
	}
	buff := "["
	for i := 0; i < butt_size / 2 - len(txt); i++ {
		buff += " "
	}
	buff += txt
	for i := 0; i < butt_size / 2 - len(txt); i++ {
		buff += " "
	}
	buff += "]"
	i_draw_text(ui.s, (ui.dim[W] / 2) - (butt_size / 2), line,
		(ui.dim[W] / 2) + (butt_size / 2), line, style, buff)
}

var start_line int

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
	var end_line int
	switch in.Protocol {
	case PROTOCOL_SSH:
		end_line = i_draw_insert_ssh(ui, line, win, in)
	}
	if win.T + end_line >= win.B {
		ui.s.SetContent(ui.dim[W] / 2, win.B, '▼', nil, ui.style[BOX_STYLE])
		// ui.s.Fini()
		// fmt.Println("end_line ", end_line)
		// fmt.Println("win.T    ", win.T)
		// fmt.Println("win.T+end", win.T + end_line)
		// fmt.Println("win.B    ", win.B)
		// fmt.Println("insert_sel   ", ui.insert_sel)
		// fmt.Println("insert_max   ", ui.insert_sel_max)
		// os.Exit(0)
		// TODO: here
	}
}

func i_draw_insert_ssh(ui HardUI, line int, win Quad, in *HostNode) int {
	red := false
	if win.T + line >= win.B { return line }
	text := "---- Host settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Host/IP", in.Host, 1, ui.insert_sel,
		false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Port", strconv.Itoa(int(in.Port)),
		2, ui.insert_sel, false);
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "User", in.User, 3, ui.insert_sel,
		false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Pass", in.Pass, 4, ui.insert_sel,
		false)
	if line += 1; win.T + line >= win.B { return line }
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
		if line += 1; win.T + line >= win.B { return line }
		text := "file does not exist"
		i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
			win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
	}
	red = false
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Jump settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Host/IP",
		in.Jump.Host, 6, ui.insert_sel, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Port",
		strconv.Itoa(int(in.Jump.Port)), 7, ui.insert_sel, false)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "User",
		in.Jump.User, 8, ui.insert_sel, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Pass",
		in.Jump.Pass, 9, ui.insert_sel, false)
	if line += 1; win.T + line >= win.B { return line}
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
		if line += 1; win.T + line >= win.B { return line }
		text := "file does not exist"
		i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
			win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
	}
	if line += 2; win.T + line >= win.B { return line }
	i_draw_ok_butt(ui, win.T + line, 11, ui.insert_sel)
	return line
}
