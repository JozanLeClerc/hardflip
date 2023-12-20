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
 * hardflip: src/i_ui.go
 * Wed Dec 20 15:24:42 2023
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

func i_draw_box(s tcell.Screen, x1, y1, x2, y2 int, title string) {
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
	i_draw_text(s, x1 + 1, y1, x2 - 1, y2 - 1, style, title)
}

func i_bottom_text(s tcell.Screen, t [2]int) {
	style := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorGrey)
	spaces := ""
	
	for i := 0; i < (t[W]) - len(KEYS_HINTS); i++ {
		spaces += " "
	}
	i_draw_text(s, 0, t[H] - 1, t[W], t[H] - 1, style, spaces + KEYS_HINTS)
}

func i_draw_zhosts_box(s tcell.Screen, t [2]int, def_style tcell.Style) {
	text := "Hosts list empty. Add hosts by pressing (a)"
	left, right :=
		(t[W] / 2) - (len(text) / 2) - 5,
		(t[W] / 2) + (len(text) / 2) + 5
	top, bot :=
		(t[H] / 2) - 3,
		(t[H] / 2) + 3
	i_draw_box(s, left, top, right, bot, "")
	if left < t[W] / 3 {
		for y := top + 1; y < bot; y++ {
		    s.SetContent(t[W] / 3, y, ' ', nil, def_style)
		}
	}
	i_draw_text(s,
		(t[W] / 2) - (len(text) / 2), t[H] / 2, right, t[H] / 2,
		def_style, text)
}

func i_host_panel(data *Data, t [2]int,
		def_style tcell.Style,
		sel uint64, sel_max uint64) {
	i_draw_box(data.s, 0, 0,
		t[W] / 3, t[H] - 2,
		" Hosts ")
	host := data.lhost.head
	for i := 0; i < data.list_start && host.next != nil; i++ {
		host = host.next
	}
	for line := 1; line < t[H] - 2 && host != nil; line++ {
		style := def_style
		if sel == host.ID {
			style = tcell.StyleDefault.
				Background(tcell.ColorWhite).
				Foreground(tcell.ColorBlack)
		}
		spaces := ""
		for i := 0; i < (t[W] / 3) - len(host.Folder + host.Name) - 2; i++ {
			spaces += " "
		}
		if host.Type == 0 {
			i_draw_text(data.s,
				1, line, t[W] / 3, line,
				style, "   " + host.Folder + host.Name + spaces)
		} else if host.Type == 1 {
			i_draw_text(data.s,
				1, line, t[W] / 3, line,
				style, "   " + host.Folder + host.Name + spaces)
		}
		i_draw_text(data.s,
			4, line, t[W] / 3, line,
			style, host.Folder + host.Name + spaces)
		host = host.next
	}
	i_draw_text(data.s,
		1, t[H] - 2, (t[W] / 3) - 1, t[H] - 2,
		def_style,
		" " + strconv.Itoa(int(sel + 1)) + "/" +
		strconv.Itoa(int(sel_max)) + " hosts ")
}

func i_info_panel(s tcell.Screen, t [2]int,
		def_style tcell.Style, lhost *HostList, sel uint64) {
	title_style := tcell.StyleDefault.
			Background(tcell.ColorReset).
			Foreground(tcell.ColorBlue).Dim(true).Bold(true)
	var host *HostNode
	curr_line := 2
	var host_type string

	i_draw_box(s, (t[W] / 3), 0,
		t[W] - 1, t[H] - 2,
		" Infos ")
	s.SetContent(t[W] / 3, 0, tcell.RuneTTee, nil, def_style)
	s.SetContent(t[W] / 3, t[H] - 2, tcell.RuneBTee, nil, def_style)
	if lhost.head == nil {
		return
	}
	host = lhost.sel(sel)
	if host.Type == 0 {
		host_type = "SSH"
	} else if host.Type == 1 {
		host_type = "RDP"
	}

	// name, type
	i_draw_text(s,
		(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
		title_style, "Name: ")
	i_draw_text(s,
		(t[W] / 3) + 10, curr_line, t[W] - 2, curr_line,
		def_style, host.Name)
	curr_line += 1
	i_draw_text(s,
		(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
		title_style, "Type: ")
	i_draw_text(s,
		(t[W] / 3) + 10, curr_line, t[W] - 2, curr_line,
		def_style, host_type)
	curr_line += 2
	// host, port
	i_draw_text(s,
		(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
		title_style, "Host: ")
	i_draw_text(s,
		(t[W] / 3) + 10, curr_line, t[W] - 2, curr_line,
		def_style, host.Host)
	curr_line += 1
	i_draw_text(s,
		(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
		title_style, "Port: ")
	i_draw_text(s,
		(t[W] / 3) + 10, curr_line, t[W] - 2, curr_line,
		def_style, strconv.Itoa(int(host.Port)))
	curr_line += 1
	// RDP shit
	if host.Type == 1 {
		if len(host.Domain) > 0 {
			i_draw_text(s,
				(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
				title_style, "Domain: ")
			i_draw_text(s,
				(t[W] / 3) + 12, curr_line, t[W] - 2, curr_line,
				def_style, host.Domain)
			curr_line += 1
		}
	}
	curr_line += 1
	// user infos
	i_draw_text(s,
		(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
		title_style, "User: ")
	i_draw_text(s,
		(t[W] / 3) + 10, curr_line, t[W] - 2, curr_line,
		def_style, host.User)
	curr_line += 1
	if len(host.Pass) > 0 {
		i_draw_text(s,
			(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
			title_style, "Pass: ")
		i_draw_text(s,
			(t[W] / 3) + 10, curr_line, t[W] - 2, curr_line,
			def_style, "***")
		curr_line += 1
	}
	if host.Type == 0 && len(host.Priv) > 0 {
		i_draw_text(s,
			(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
			title_style, "Privkey: ")
		i_draw_text(s,
			(t[W] / 3) + 13, curr_line, t[W] - 2, curr_line,
			def_style, host.Priv)
		curr_line += 1
	}
	curr_line += 1
	// jump
	if host.Type == 0 && len(host.Jump) > 0 {
		i_draw_text(s,
			(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
			title_style, "Jump settings: ")
		curr_line += 1
		i_draw_text(s,
			(t[W] / 3) + 5, curr_line, t[W] - 2, curr_line,
			title_style, "Host: ")
		i_draw_text(s,
			(t[W] / 3) + 11, curr_line, t[W] - 2, curr_line,
			def_style, host.Jump)
		curr_line += 1
		i_draw_text(s,
			(t[W] / 3) + 5, curr_line, t[W] - 2, curr_line,
			title_style, "Port: ")
		i_draw_text(s,
			(t[W] / 3) + 11, curr_line, t[W] - 2, curr_line,
			def_style, strconv.Itoa(int(host.JumpPort)))
		curr_line += 1
		i_draw_text(s,
			(t[W] / 3) + 5, curr_line, t[W] - 2, curr_line,
			title_style, "User: ")
		i_draw_text(s,
			(t[W] / 3) + 11, curr_line, t[W] - 2, curr_line,
			def_style, host.JumpUser)
		curr_line += 1
		if len(host.JumpPass) > 0 {
			i_draw_text(s,
				(t[W] / 3) + 5, curr_line, t[W] - 2, curr_line,
				title_style, "Pass: ")
			i_draw_text(s,
				(t[W] / 3) + 11, curr_line, t[W] - 2, curr_line,
				def_style, "***")
			curr_line += 1
		}
		if host.Type == 0 && len(host.JumpPriv) > 0 {
			i_draw_text(s,
				(t[W] / 3) + 5, curr_line, t[W] - 2, curr_line,
				title_style, "Privkey: ")
			i_draw_text(s,
				(t[W] / 3) + 14, curr_line, t[W] - 2, curr_line,
				def_style, host.JumpPriv)
			curr_line += 1
		}
		curr_line += 1
	}
	// RDP shit
	if host.Type == 1 {
		qual := [3]string{"Low", "Medium", "High"}
		i_draw_text(s,
			(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
			title_style, "Screen size: ")
		i_draw_text(s,
			(t[W] / 3) + 17, curr_line, t[W] - 2, curr_line,
			def_style,
			strconv.Itoa(int(host.Width)) + "x" +
			strconv.Itoa(int(host.Height)))
		curr_line += 1
		i_draw_text(s,
			(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
			title_style, "Dynamic window: ")
		i_draw_text(s,
			(t[W] / 3) + 20, curr_line, t[W] - 2, curr_line,
			def_style, strconv.FormatBool(host.Dynamic))
		curr_line += 1
		i_draw_text(s,
			(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
			title_style, "Quality: ")
		i_draw_text(s,
			(t[W] / 3) + 13, curr_line, t[W] - 2, curr_line,
			def_style, qual[host.Quality])
		curr_line += 1
		curr_line += 1
	}
	// note
	if len(host.Note) > 0 {
		i_draw_text(s,
			(t[W] / 3) + 4, curr_line, t[W] - 2, curr_line,
			title_style, "Note: ")
		i_draw_text(s,
			(t[W] / 3) + 10, curr_line, t[W] - 2, curr_line,
			def_style, host.Note)
		curr_line += 1
	}
}

func i_ui(data *Data) {
	var err error
	data.s, err = tcell.NewScreen()
	var term_size [2]int
	var sel uint64 = 0
	sel_max := data.lhost.count()

	if err != nil {
		c_die("view", err)
	}
	if err := data.s.Init(); err != nil {
		c_die("view", err)
	}
	def_style := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	data.s.SetStyle(def_style)
	for {
		term_size[W], term_size[H], _ = term.GetSize(0)
		data.s.Clear()
		i_bottom_text(data.s, term_size)
		i_host_panel(data, term_size, def_style, sel, sel_max)
		i_info_panel(data.s, term_size, def_style, data.lhost, sel)
		if data.lhost.head == nil {
			i_draw_zhosts_box(data.s, term_size, def_style)
		}
		data.s.Show()
		i_events(data, &sel, &sel_max, &term_size)
		if int(sel) > term_size[H] - 6 {
			data.list_start += 1
		}
	}
}
