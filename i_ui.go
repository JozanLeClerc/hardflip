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
 * Thu Jan 18 17:51:03 2024
 * Joe
 *
 * interfacing with the user
 */

package main

import (
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type HardUI struct {
	s            tcell.Screen
	mode         uint8
	def_style    tcell.Style
	dir_style    tcell.Style
	title_style  tcell.Style
	dim          [2]int
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
	i_draw_text(s, x1 + 1, y1, x2 - 1, y1, style, title)
}

func i_bottom_text(ui HardUI) {
	text := ""

	switch ui.mode {
	case NORMAL_MODE:
		text = NORMAL_KEYS_HINTS
	case DELETE_MODE:
		text = DELETE_KEYS_HINTS
	}
	i_draw_text(ui.s,
		0, ui.dim[H] - 1, ui.dim[W], ui.dim[H] - 1,
		ui.def_style.Dim(true), text)
	text = " " + VERSION
	i_draw_text(ui.s,
		ui.dim[W] - 5, ui.dim[H] - 1, ui.dim[W], ui.dim[H] - 1,
		ui.def_style.Dim(true), text)
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

func i_draw_delete_box(ui HardUI, item *ItemsNode) {
	var text string
	var file string

	if item.is_dir() == true {
		text = "Really delete this directory and all of its content?"
		file = item.Dirs.path()
	} else {
		host := item.Host
		text = "Really delete this host?"
		file = host.Parent.path() + host.Filename
	}
	max_len := len(text)
    
	if max_len < len(file) {
	    max_len = len(file)
	}
	left, right :=
	    (ui.dim[W] / 2) - (max_len / 2) - 5,
	    (ui.dim[W] / 2) + (max_len / 2) + 5
	if left < ui.dim[W] / 8 {
	    left = ui.dim[W] / 8
	}
	if right > ui.dim[W] - ui.dim[W] / 8 - 1 {
	    right = ui.dim[W] - ui.dim[W] / 8 - 1
	}
	top, bot :=
	    (ui.dim[H] / 2) - 4,
	    (ui.dim[H] / 2) + 3
	i_draw_box(ui.s, left, top, right, bot, "", true)
	left = (ui.dim[W] / 2) - (len(text) / 2)
	if left < (ui.dim[W] / 8) + 1 {
	    left = (ui.dim[W] / 8) + 1
	}
	top = ui.dim[H] / 2 - 2
	i_draw_text(ui.s,
	    left, top, right, top,
	    ui.def_style, text)
	left = (ui.dim[W] / 2) - (len(file) / 2)
	if left < (ui.dim[W] / 8) + 1 {
	    left = (ui.dim[W] / 8) + 1
	}
	top += 1
	i_draw_text(ui.s,
	    left, top, right, top,
	    ui.def_style.Bold(true), file)
	left = right - 11
	if left < (ui.dim[W] / 8) + 1 {
	    left = (ui.dim[W] / 8) + 1
	}
	top = ui.dim[H] / 2 + 1
	i_draw_text(ui.s,
	    left - 1, top, right, top,
	    ui.def_style, "[")
	i_draw_text(ui.s,
	    left, top, right, top,
	    ui.def_style.Bold(true).Underline(true), "y")
	i_draw_text(ui.s,
	    left + 1, top, right, top,
	    ui.def_style, "es] [")
	i_draw_text(ui.s,
	    left + 6, top, right, top,
	    ui.def_style.Bold(true).Underline(true), "n")
	i_draw_text(ui.s,
	    left + 7, top, right, top,
	    ui.def_style, "o]")
}

func i_host_panel_dirs(ui HardUI, icons bool, dir_icon uint8,
	dir *DirsNode, curr *DirsNode, line int) {
	style := ui.dir_style
	if dir == curr {
		style = style.Reverse(true)
	}
	text := ""
	for i := 0; i < int(dir.Depth) - 2; i++ {
		text += "  "
	}
	if icons == true {
		text += DIRS_ICONS[dir_icon]
	}
	text += dir.Name
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
	for i := 0; i < int(depth + 1) - 2; i++ {
		text += "  "
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

func i_host_panel(ui HardUI, icons bool, litems *ItemsList, data *HardData) {
	i_draw_box(ui.s, 0, 0,
		ui.dim[W] / 3, ui.dim[H] - 2,
		" Hosts ", false)
	line := 1
	if litems.head == nil {
		return
	}
	for ptr := litems.draw; ptr != nil && line < ui.dim[H] - 2; ptr = ptr.next {
		if ptr.is_dir() == false && ptr.Host != nil  {
			i_host_panel_host(ui,
				icons,
				ptr.Host.Parent.Depth,
				ptr.Host,
				litems.curr.Host,
				line)
			line++
		} else if ptr.Dirs != nil {
			var dir_icon uint8
			if data.folds[ptr.Dirs] != nil {
				dir_icon = 1
			}
			i_host_panel_dirs(ui, icons, dir_icon,
				ptr.Dirs,
				litems.curr.Dirs,
				line)
			line++
		}
	}
}

func i_info_panel_dirs(ui HardUI, dir *DirsNode) {
	line := 2
	if line > ui.dim[H] - 3 {
		return
	}

	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
		ui.title_style, "Name: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.def_style, dir.Name)
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
		ui.title_style, "Type: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.def_style, "Directory")
	if line += 2; line > ui.dim[H] - 3 {
		return
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
		ui.title_style, "Path: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.def_style, dir.path())
}

func i_info_panel_host(ui HardUI, host *HostNode) {
	host_type := host.protocol_str()
	line := 2
	if line > ui.dim[H] - 3 {
		return
	}
	// name, type
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
		ui.title_style, "Name: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.def_style, host.Name)
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
		ui.title_style, "Type: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.def_style, host_type)
	if line += 2; line > ui.dim[H] - 3 {
		return
	}
	if line > ui.dim[H] - 3 {
		return
	}
	// host, port
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
		ui.title_style, "Host: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.def_style, host.Host)
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
		ui.title_style, "Port: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.def_style, strconv.Itoa(int(host.Port)))
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	// RDP shit
	if host.Protocol == 1 {
		if len(host.Domain) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
				ui.title_style, "Domain: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 12, line, ui.dim[W] - 2, line,
				ui.def_style, host.Domain)
			if line += 1; line > ui.dim[H] - 3 {
				return
			}
		}
	}
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	// user infos
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
		ui.title_style, "User: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.def_style, host.User)
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	if len(host.Pass) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.title_style, "Pass: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
			ui.def_style, "***")
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
	}
	if host.Protocol == 0 && len(host.Priv) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.title_style, "Privkey: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 13, line, ui.dim[W] - 2, line,
			ui.def_style, host.Priv)
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
	}
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	// jump
	if host.Protocol == 0 && len(host.Jump) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.title_style, "Jump settings: ")
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 5, line, ui.dim[W] - 2, line,
			ui.title_style, "Host: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 11, line, ui.dim[W] - 2, line,
			ui.def_style, host.Jump)
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 5, line, ui.dim[W] - 2, line,
			ui.title_style, "Port: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 11, line, ui.dim[W] - 2, line,
			ui.def_style, strconv.Itoa(int(host.JumpPort)))
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 5, line, ui.dim[W] - 2, line,
			ui.title_style, "User: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 11, line, ui.dim[W] - 2, line,
			ui.def_style, host.JumpUser)
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		if len(host.JumpPass) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 5, line, ui.dim[W] - 2, line,
				ui.title_style, "Pass: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 11, line, ui.dim[W] - 2, line,
				ui.def_style, "***")
			if line += 1; line > ui.dim[H] - 3 {
				return
			}
		}
		if host.Protocol == 0 && len(host.JumpPriv) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 5, line, ui.dim[W] - 2, line,
				ui.title_style, "Privkey: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 14, line, ui.dim[W] - 2, line,
				ui.def_style, host.JumpPriv)
			if line += 1; line > ui.dim[H] - 3 {
				return
			}
		}
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
	}
	// RDP shit
	if host.Protocol == 1 {
		qual := [3]string{"Low", "Medium", "High"}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.title_style, "Screen size: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 17, line, ui.dim[W] - 2, line,
			ui.def_style,
			strconv.Itoa(int(host.Width)) + "x" +
			strconv.Itoa(int(host.Height)))
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.title_style, "Dynamic window: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 20, line, ui.dim[W] - 2, line,
			ui.def_style, strconv.FormatBool(host.Dynamic))
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.title_style, "Quality: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 13, line, ui.dim[W] - 2, line,
			ui.def_style, qual[host.Quality])
		line += 1
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
	}
	// note
	if len(host.Note) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.title_style, "Note: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
			ui.def_style, host.Note)
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
	}
}

func i_info_panel(ui HardUI, percent bool, litems *ItemsList) {
	i_draw_box(ui.s, (ui.dim[W] / 3), 0,
		ui.dim[W] - 1, ui.dim[H] - 2,
		" Infos ", false)
	ui.s.SetContent(ui.dim[W] / 3, 0, tcell.RuneTTee, nil, ui.def_style)
	ui.s.SetContent(ui.dim[W] / 3, ui.dim[H] - 2,
		tcell.RuneBTee, nil, ui.def_style)
	// number display
	if litems.head != nil {
		text := " " + strconv.Itoa(litems.curr.ID) + " of " +
			strconv.Itoa(int(litems.last.ID)) + " "
		if percent == true {
			text += "- " +
				strconv.Itoa(litems.curr.ID * 100 / litems.last.ID) + "% "
		}
		i_draw_text(ui.s,
			(ui.dim[W] - 1) - len(text) - 1,
			ui.dim[H] - 2,
			(ui.dim[W] - 1) - 1,
			ui.dim[H] - 2,
			ui.def_style,
			text)
	} else {
		text := " 0 hosts "
		i_draw_text(ui.s,
			(ui.dim[W] - 1) - len(text) - 1,
			ui.dim[H] - 2,
			(ui.dim[W] - 1) - 1,
			ui.dim[H] - 2,
			ui.def_style,
			text)
	}
	// panel
	if litems.head == nil {
		return
	} else if litems.curr.is_dir() == true {
		i_info_panel_dirs(ui, litems.curr.Dirs)
	} else {
		i_info_panel_host(ui, litems.curr.Host)
	}
}

func i_scrollhint(ui HardUI, litems *ItemsList) {
	if litems.head == nil {
		return
	}
	h := ui.dim[H] - 4
	max := litems.last.ID
	if max <= h {
		return
	}
	draw_id := litems.draw.ID
	if draw_id > 1 {
		ui.s.SetContent(0, 1,
			'▲',
			nil, ui.def_style)
	}
	if max - draw_id > h {
		ui.s.SetContent(0, ui.dim[H] - 3,
			'▼',
			nil, ui.def_style)
		return
	}
}

// HACK: fuck global vars but do we have the choice really
var g_load_count int = -1

func i_display_load_ui(ui *HardUI) {
	g_load_count += 1
	if g_load_count % 1000 != 0 {
		return
	}
	ui.s.Clear()
	text := "Loading " + strconv.Itoa(g_load_count) + " hosts"
	text_len := len(text) / 2
	// TODO: max len
	i_draw_box(ui.s,
		(ui.dim[W] / 2) - (text_len + 2) - 1,
		(ui.dim[H] / 2) - 2,
		(ui.dim[W] / 2) + (text_len + 2),
		(ui.dim[H] / 2) + 2, " Loading hosts ", false)
	i_draw_text(ui.s,
		(ui.dim[W] / 2) - text_len,
		(ui.dim[H] / 2),
		(ui.dim[W] / 2) + text_len + 1,
		(ui.dim[H] / 2),
		ui.def_style, text)
	ui.s.Show()
	event := ui.s.PollEvent()
	ui.s.PostEvent(event)
	switch event := event.(type) {
	case *tcell.EventResize:
		ui.dim[W], ui.dim[H], _ = term.GetSize(0)
		ui.s.Sync()
	case *tcell.EventKey:
		if event.Key() == tcell.KeyCtrlC ||
		event.Rune() == 'q' {
			ui.s.Fini()
			os.Exit(0)
		}
	}
}

func i_load_ui(data_dir string,
			   opts HardOpts,
			   ui *HardUI) (*DirsList, *ItemsList) {
	ui.mode = LOAD_MODE
	ldirs := c_load_data_dir(data_dir, opts, ui)
	litems := c_load_litems(ldirs)
	ui.mode = NORMAL_MODE
	return ldirs, litems
}

func i_ui(data_dir string, opts HardOpts) {
	ui := HardUI{}
	var err error

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
	ui.dim[W], ui.dim[H], _ = term.GetSize(0)
	ldirs, litems := i_load_ui(data_dir, opts, &ui)
	data := HardData{
		litems,
		ldirs,
		ui,
		opts,
		data_dir,
		make(map[*DirsNode]*ItemsList),
	}
	for {
		data.ui.s.Clear()
		i_bottom_text(data.ui)
		i_host_panel(data.ui, data.opts.Icon, data.litems, &data)
		i_info_panel(data.ui, data.opts.Perc, data.litems)
		i_scrollhint(data.ui, data.litems)
		if data.litems.head == nil {
			i_draw_zhosts_box(data.ui)
		}
		if data.ui.mode == DELETE_MODE {
			i_draw_delete_box(data.ui, data.litems.curr)
		}
		data.ui.s.Show()
		i_events(&data)
	}
}
