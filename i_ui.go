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
 * hardflip: src/c_hardflip.go
 * Mon, 18 Dec 2023 19:01:59 +0100
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
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

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

func i_bottom_text(s tcell.Screen, term_w, term_h int, style tcell.Style) {
	i_draw_text(s, 0, term_h - 1, term_w - 1, term_h, style, "(q) Quit")
}

func i_hosts_panel(s tcell.Screen,
		term_w, term_h int,
		def_style tcell.Style, lhost *HostList,
		sel uint64, sel_max uint64) {
	i_draw_box(s, 0, 0,
		term_w / 3, term_h - 2,
		" hosts ")
	host := lhost.head
	for host != nil {
		style := def_style
		if sel == host.ID {
			style = tcell.StyleDefault.
				Background(tcell.ColorWhite).
				Foreground(tcell.ColorBlack)
		}
		spaces := ""
		i := 0
		for i < (term_w / 3) - len(host.Folder + host.Name) - 2 {
			spaces += " "
			i++
		}
		if host.Type == 0 {
			i_draw_text(s,
				1, int(host.ID) + 1, term_w / 3, int(host.ID) + 1,
				style, "   " + host.Folder + host.Name + spaces)
		} else if host.Type == 1 {
			i_draw_text(s,
				1, int(host.ID) + 1, term_w / 3, int(host.ID) + 1,
				style, "   " + host.Folder + host.Name + spaces)
		}
		i_draw_text(s,
			4, int(host.ID) + 1, term_w / 3, int(host.ID) + 1,
			style, host.Folder + host.Name + spaces)
		host = host.next
	}
	i_draw_text(s,
		1, term_h - 2, (term_w / 3) - 1, term_h - 1,
		def_style, " " + strconv.Itoa(int(sel_max)) + " hosts ")
}
func i_info_panel(s tcell.Screen,
		term_w, term_h int,
		def_style tcell.Style, lhost *HostList, sel uint64) {
	title_style := tcell.StyleDefault.
			Background(tcell.ColorReset).
			Foreground(tcell.ColorBlue).Dim(true).Bold(true)
	host := lhost.sel(sel)
	curr_line := 2
	var host_type string

	i_draw_box(s, (term_w / 3), 0,
		term_w - 1, term_h - 2,
		" infos ")
	if host.Type == 0 {
		host_type = "SSH"
	} else if host.Type == 1 {
		host_type = "RDP"
	}

	// name, type
	i_draw_text(s,
		(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
		title_style, "Name: ")
	i_draw_text(s,
		(term_w / 3) + 10, curr_line, term_w - 2, curr_line,
		def_style, host.Name)
	curr_line += 1
	i_draw_text(s,
		(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
		title_style, "Type: ")
	i_draw_text(s,
		(term_w / 3) + 10, curr_line, term_w - 2, curr_line,
		def_style, host_type)
	curr_line += 2
	// host, port
	i_draw_text(s,
		(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
		title_style, "Host: ")
	i_draw_text(s,
		(term_w / 3) + 10, curr_line, term_w - 2, curr_line,
		def_style, host.Host)
	curr_line += 1
	i_draw_text(s,
		(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
		title_style, "Port: ")
	i_draw_text(s,
		(term_w / 3) + 10, curr_line, term_w - 2, curr_line,
		def_style, strconv.Itoa(int(host.Port)))
	curr_line += 2
	// user infos
	i_draw_text(s,
		(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
		title_style, "User: ")
	i_draw_text(s,
		(term_w / 3) + 10, curr_line, term_w - 2, curr_line,
		def_style, host.User)
	curr_line += 1
	if len(host.Pass) > 0 {  
		i_draw_text(s,
		(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
		title_style, "Pass: ")
		i_draw_text(s,
		(term_w / 3) + 10, curr_line, term_w - 2, curr_line,
		def_style, "***")
		curr_line += 1
	}
	if host.Type == 0 && len(host.Priv) > 0 {
		i_draw_text(s,
		(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
		title_style, "Privkey: ")
		i_draw_text(s,
		(term_w / 3) + 13, curr_line, term_w - 2, curr_line,
		def_style, host.Priv)
		curr_line += 1
	}
	curr_line += 1
	// jump
	if host.Type == 0 && len(host.Jump) > 0 {
		i_draw_text(s,
			(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
			title_style, "Jump settings: ")
		curr_line += 1
		i_draw_text(s,
			(term_w / 3) + 6, curr_line, term_w - 2, curr_line,
			title_style, "Jump host: ")
		i_draw_text(s,
			(term_w / 3) + 17, curr_line, term_w - 2, curr_line,
			def_style, host.Jump)
		curr_line += 1
		i_draw_text(s,
			(term_w / 3) + 6, curr_line, term_w - 2, curr_line,
			title_style, "Jump port: ")
		i_draw_text(s,
			(term_w / 3) + 17, curr_line, term_w - 2, curr_line,
			def_style, strconv.Itoa(int(host.JumpPort)))
		curr_line += 1
		i_draw_text(s,
			(term_w / 3) + 6, curr_line, term_w - 2, curr_line,
			title_style, "Jump user: ")
		i_draw_text(s,
			(term_w / 3) + 17, curr_line, term_w - 2, curr_line,
			def_style, host.JumpUser)
		curr_line += 2
	}
	// note
	i_draw_text(s,
		(term_w / 3) + 4, curr_line, term_w - 2, curr_line,
		title_style, "Note: ")
	i_draw_text(s,
		(term_w / 3) + 10, curr_line, term_w - 2, curr_line,
		def_style, host.Note)
	curr_line += 1
}

func i_events(s tcell.Screen,
		sel *uint64, sel_max *uint64,
		lhost *HostList, quit func()) {
	event := s.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
		s.Sync()
	case *tcell.EventKey:
		if event.Key() == tcell.KeyEscape ||
		event.Key() == tcell.KeyCtrlC ||
		event.Rune() == 'q' ||
		event.Rune() == 'Q' {
			quit()
			os.Exit(0)
		}
		if event.Rune() == 'j' ||
		event.Key() == tcell.KeyDown {
			if *sel < *sel_max - 1 {
				*sel += 1
			}
		}
		if event.Rune() == 'k' ||
		event.Key() == tcell.KeyUp {
			if *sel > 0 {
				*sel -= 1
			}
		}
		if event.Key() == tcell.KeyEnter {
			quit()
			c_exec(*sel, lhost)
			os.Exit(0)
		}
	}
}

func i_ui(lhost *HostList) {
	screen, err := tcell.NewScreen()
	var sel uint64 = 0
	sel_max := lhost.count()

	if err != nil {
		c_die("view", err)
	}
	if err := screen.Init(); err != nil {
		c_die("view", err)
	}
	def_style := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	screen.SetStyle(def_style)
	quit := func() {
		screen.Fini()
	}
	for {
		term_w, term_h, _ := term.GetSize(0)
		screen.Clear()
		i_bottom_text(screen, term_w, term_h, def_style)
		i_hosts_panel(screen, term_w, term_h, def_style, lhost, sel, sel_max)
		i_info_panel(screen, term_w, term_h, def_style, lhost, sel)
		screen.Show()
		i_events(screen, &sel, &sel_max, lhost, quit)
	}
}
