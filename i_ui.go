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
 * hardflip: src/i_ui.go
 * Mon Jan 08 11:37:32 2024
 * Joe
 *
 * interfacing with the user
 */

package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type HardUI struct {
	s           tcell.Screen
	list_start  int
	mode        uint8
	sel_max     int
	sel_id		int
	count_dirs  int
	count_hosts int
	def_style   tcell.Style
	dir_style   tcell.Style
	title_style tcell.Style
	dim         [2]int
}

func (ui *HardUI) inc_sel(n int, data *HardData) {
	// HACK: ui.sel_id and litems.curr.ID are essentially the same thing
	// NOTE: maybe keep only one there, see later
	if data.litems.curr == nil {
		return
	}
	if ui.sel_id + n < 1 || 
	   ui.sel_id + n > ui.sel_max {
		n = 0
	}
	data.litems.curr = data.litems.curr.inc(n)
	ui.sel_id = data.litems.curr.ID
	if ui.sel_id > ui.list_start + ui.dim[H] - 4 {
		ui.list_start = (ui.sel_id + 1) - (ui.dim[H] + 3)
	} else if ui.sel_id < ui.list_start {
		ui.list_start = ui.sel_id
	}
}

func i_draw_text(s tcell.Screen,
		x1, y1, x2, y2 int,
		style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func i_draw_box(s tcell.Screen, x1, y1, x2, y2 int, title string, fill bool) {
	style := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	if fill == true {
		for y := y1 + 1; y < y2; y++ {
			for x := x1 + 1; x < x2; x++ {
				s.SetContent(x, y, ' ', nil, style)
			}
		}
	}
	i_draw_text(s, x1 + 1, y1, x2 - 1, y2 - 1, style, title)
}

func i_bottom_text(ui HardUI) {
	spaces := ""
	text := ""

	switch ui.mode {
	case NORMAL_MODE:
		text = NORMAL_KEYS_HINTS
	case DELETE_MODE:
		text = DELETE_KEYS_HINTS
	}
	for i := 0; i < (ui.dim[W]) - len(text); i++ {
		spaces += " "
	}
	i_draw_text(ui.s,
		0, ui.dim[H] - 1, ui.dim[W], ui.dim[H] - 1,
		ui.def_style.Dim(true), spaces + text)
}

func i_draw_zhosts_box(ui HardUI) {
	text := "Hosts list empty. Add hosts/folders by pressing (a/m)"
	left, right :=
		(ui.dim[W] / 2) - (len(text) / 2) - 5,
		(ui.dim[W] / 2) + (len(text) / 2) + 5
	top, bot :=
		(ui.dim[H] / 2) - 3,
		(ui.dim[H] / 2) + 3
	i_draw_box(ui.s, left, top, right, bot, "", false)
	if left < ui.dim[W] / 3 {
		for y := top + 1; y < bot; y++ {
		    ui.s.SetContent(ui.dim[W] / 3, y, ' ', nil, ui.def_style)
		}
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 2) - (len(text) / 2), ui.dim[H] / 2, right, ui.dim[H] / 2,
		ui.def_style, text)
}

func i_draw_delete_box(ui HardUI, host *HostNode) {
	// text := "Really delete this host?"
	// file := host.Folder + host.Filename
	// max_len := len(text)
    //
	// if max_len < len(file) {
	//     max_len = len(file)
	// }
	// left, right :=
	//     (ui.dim[W] / 2) - (max_len / 2) - 5,
	//     (ui.dim[W] / 2) + (max_len / 2) + 5
	// if left < ui.dim[W] / 8 {
	//     left = ui.dim[W] / 8
	// }
	// if right > ui.dim[W] - ui.dim[W] / 8 - 1 {
	//     right = ui.dim[W] - ui.dim[W] / 8 - 1
	// }
	// top, bot :=
	//     (ui.dim[H] / 2) - 4,
	//     (ui.dim[H] / 2) + 3
	// i_draw_box(ui.s, left, top, right, bot, "", true)
	// left = (ui.dim[W] / 2) - (len(text) / 2)
	// if left < (ui.dim[W] / 8) + 1 {
	//     left = (ui.dim[W] / 8) + 1
	// }
	// top = ui.dim[H] / 2 - 2
	// i_draw_text(ui.s,
	//     left, top, right, top,
	//     ui.def_style, text)
	// left = (ui.dim[W] / 2) - (len(file) / 2)
	// if left < (ui.dim[W] / 8) + 1 {
	//     left = (ui.dim[W] / 8) + 1
	// }
	// top += 1
	// i_draw_text(ui.s,
	//     left, top, right, top,
	//     ui.def_style.Bold(true), file)
	// left = right - 11
	// if left < (ui.dim[W] / 8) + 1 {
	//     left = (ui.dim[W] / 8) + 1
	// }
	// top = ui.dim[H] / 2 + 1
	// i_draw_text(ui.s,
	//     left, top, right, top,
	//     ui.def_style.Bold(true).Underline(true), "y")
	// i_draw_text(ui.s,
	//     left + 1, top, right, top,
	//     ui.def_style, "es | ")
	// i_draw_text(ui.s,
	//     left + 6, top, right, top,
	//     ui.def_style.Bold(true).Underline(true), "n")
	// i_draw_text(ui.s,
	//     left + 7, top, right, top,
	//     ui.def_style, "o")
}

func i_host_panel_dirs(ui HardUI, icons bool,
	dirs *DirsNode, curr *DirsNode, line int) {
	style := ui.dir_style
	if dirs == curr {
		style = style.Reverse(true)
	}
	text := ""
	for i := 0; i < int(dirs.Depth) - 2; i++ {
		for i := 0; i < int(dirs.Depth) - 2; i++ {
			text += "  "
		}
	}
	if icons == true {
		var fold_var uint8
		if dirs.Folded == true {
			fold_var = 1
		}
		text += DIRS_ICONS[fold_var]
	}
	text += dirs.Name
	spaces := ""
	for i := 0; i < (ui.dim[W] / 3) - len(text) + 1; i++ {
		spaces += " "
	}
	text += spaces
	i_draw_text(ui.s,
		1, line, ui.dim[W] / 3, line,
		style, text)
}

func i_host_panel_host(ui HardUI, icons bool,
		depth uint16, host *HostNode, curr *HostNode, line int) {
	style := ui.def_style
	if host == curr {
			style = style.Reverse(true)
		}
		text := ""
		for i := 0; i < int(depth) - 2; i++ {
			text += "    "
		}
		if icons == true {
			text += HOST_ICONS[int(host.Protocol)]
		}
		text += host.Name
		spaces := ""
		for i := 0; i < (ui.dim[W] / 3) - len(text) + 1; i++ {
			spaces += " "
		}
		text += spaces
		i_draw_text(ui.s,
			1, line, ui.dim[W] / 3, line,
			style, text)
}

func i_host_panel(ui HardUI, icons bool, litems *ItemsList) {
	i_draw_box(ui.s, 0, 0,
		ui.dim[W] / 3, ui.dim[H] - 2,
		" Hosts ", false)
	line := 1
	ptr := litems.head.next
	for ptr = ptr; ptr != nil; ptr = ptr.next {
		if ptr.is_dir() == false {
			i_host_panel_host(ui,
				icons,
				ptr.Host.Parent.Depth,
				ptr.Host,
				litems.curr.Host,
				line)
		} else {
			i_host_panel_dirs(ui, icons,
				ptr.Dirs,
				litems.curr.Dirs,
				line)
		}
		line++
	}

	// dirs := litems.head.Dirs
	// for host := dirs.lhost.head;
	// 	dirs.Folded == false && host != nil;
	// 	host = host.next {
	// 	i_host_panel_host(ui, icons, dirs.Depth, host, line)
	// 	line++
	// }
	// dirs = dirs.next
	// for line = line; line < ui.dim[H] - 2 && dirs != nil; dirs = dirs.next {
	// 	i_host_panel_dirs(ui, icons, dirs, line)
	// 	line++
	// 	for host := dirs.lhost.head;
	// 		dirs.Folded == false && host != nil;
	// 		host = host.next {
	// 		i_host_panel_host(ui, icons, dirs.Depth, host, line)
	// 		line++
	// 	}
	// }
	if ui.sel_max != 0 {
		i_draw_text(ui.s,
			1, ui.dim[H] - 2, (ui.dim[W] / 3) - 1, ui.dim[H] - 2,
			ui.def_style,
			" " + strconv.Itoa(litems.curr.ID) + "/" +
			strconv.Itoa(int(ui.sel_max)) + " items ")
	} else {
		i_draw_text(ui.s,
			1, ui.dim[H] - 2, (ui.dim[W] / 3) - 1, ui.dim[H] - 2,
			ui.def_style,
			" 0 hosts ")
	}
}

func i_info_panel(ui HardUI, lhost *HostList) {
	var host *HostNode
	var host_type string
	curr_line := 2

	i_draw_box(ui.s, (ui.dim[W] / 3), 0,
		ui.dim[W] - 1, ui.dim[H] - 2,
		" Infos ", false)
	ui.s.SetContent(ui.dim[W] / 3, 0, tcell.RuneTTee, nil, ui.def_style)
	ui.s.SetContent(ui.dim[W] / 3, ui.dim[H] - 2,
		tcell.RuneBTee, nil, ui.def_style)
	if lhost.head == nil {
		return
	}
	host = lhost.sel(ui.sel_id)
	host_type = host.protocol_str()
	// name, type
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
		ui.title_style, "Name: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, curr_line, ui.dim[W] - 2, curr_line,
		ui.def_style, host.Name)
	curr_line += 1
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
		ui.title_style, "Type: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, curr_line, ui.dim[W] - 2, curr_line,
		ui.def_style, host_type)
	curr_line += 2
	// host, port
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
		ui.title_style, "Host: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, curr_line, ui.dim[W] - 2, curr_line,
		ui.def_style, host.Host)
	curr_line += 1
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
		ui.title_style, "Port: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, curr_line, ui.dim[W] - 2, curr_line,
		ui.def_style, strconv.Itoa(int(host.Port)))
	curr_line += 1
	// RDP shit
	if host.Protocol == 1 {
		if len(host.Domain) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
				ui.title_style, "Domain: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 12, curr_line, ui.dim[W] - 2, curr_line,
				ui.def_style, host.Domain)
			curr_line += 1
		}
	}
	curr_line += 1
	// user infos
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
		ui.title_style, "User: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, curr_line, ui.dim[W] - 2, curr_line,
		ui.def_style, host.User)
	curr_line += 1
	if len(host.Pass) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Pass: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style, "***")
		curr_line += 1
	}
	if host.Protocol == 0 && len(host.Priv) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Privkey: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 13, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style, host.Priv)
		curr_line += 1
	}
	curr_line += 1
	// jump
	if host.Protocol == 0 && len(host.Jump) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Jump settings: ")
		curr_line += 1
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 5, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Host: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 11, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style, host.Jump)
		curr_line += 1
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 5, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Port: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 11, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style, strconv.Itoa(int(host.JumpPort)))
		curr_line += 1
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 5, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "User: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 11, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style, host.JumpUser)
		curr_line += 1
		if len(host.JumpPass) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 5, curr_line, ui.dim[W] - 2, curr_line,
				ui.title_style, "Pass: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 11, curr_line, ui.dim[W] - 2, curr_line,
				ui.def_style, "***")
			curr_line += 1
		}
		if host.Protocol == 0 && len(host.JumpPriv) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 5, curr_line, ui.dim[W] - 2, curr_line,
				ui.title_style, "Privkey: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 14, curr_line, ui.dim[W] - 2, curr_line,
				ui.def_style, host.JumpPriv)
			curr_line += 1
		}
		curr_line += 1
	}
	// RDP shit
	if host.Protocol == 1 {
		qual := [3]string{"Low", "Medium", "High"}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Screen size: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 17, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style,
			strconv.Itoa(int(host.Width)) + "x" +
			strconv.Itoa(int(host.Height)))
		curr_line += 1
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Dynamic window: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 20, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style, strconv.FormatBool(host.Dynamic))
		curr_line += 1
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Quality: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 13, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style, qual[host.Quality])
		curr_line += 1
		curr_line += 1
	}
	// note
	if len(host.Note) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, curr_line, ui.dim[W] - 2, curr_line,
			ui.title_style, "Note: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, curr_line, ui.dim[W] - 2, curr_line,
			ui.def_style, host.Note)
		curr_line += 1
	}
}

func i_get_sel_max(ldirs *DirsList) (int, int, int) {
	count_dirs, count_hosts := ldirs.count()

	count_dirs -= 1
	return count_dirs + count_hosts, count_dirs, count_hosts
}

func i_ui(data *HardData) {
	// TODO: replace everything ui with litems
	var err error
	ui := &data.ui
	// TODO: get better counts
	// NOTE: put this in c_load_data and also produce better code
	ui.sel_max, ui.count_dirs, ui.count_hosts = i_get_sel_max(data.ldirs)

	ui.s, err = tcell.NewScreen()
	if err != nil {
		c_die("view", err)
	}
	if err := ui.s.Init(); err != nil {
		c_die("view", err)
	}
	ui.def_style = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	ui.dir_style = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorBlue).Dim(true).Bold(true)
	ui.title_style = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorBlue).Dim(true).Bold(true)
	ui.s.SetStyle(ui.def_style)
	for {
		ui.dim[W], ui.dim[H], _ = term.GetSize(0)
		ui.s.Clear()
		i_bottom_text(*ui)
		i_host_panel(data.ui, data.opts.Icon, data.litems)
		// TODO: info panel
		// i_info_panel(data.ui, data.lhost)
		if data.ldirs.head.lhost.head == nil && data.ldirs.head.next == nil {
			i_draw_zhosts_box(*ui)
		}
		if ui.mode == DELETE_MODE {
			// TODO: delete mode
			// host := data.lhost.sel(ui.sel)
			// i_draw_delete_box(*ui, host)
		}
		ui.s.Show()
		i_events(data)
	}
}
