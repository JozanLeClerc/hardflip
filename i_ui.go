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

func i_draw_bottom_text(ui HardUI) {
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

func i_draw_err_box() {
}

func i_draw_scrollhint(ui HardUI, litems *ItemsList) {
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

func i_draw_load_ui(ui *HardUI) {
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
		i_draw_bottom_text(data.ui)
		i_draw_host_panel(data.ui, data.opts.Icon, data.litems, &data)
		i_draw_info_panel(data.ui, data.opts.Perc, data.litems)
		i_draw_scrollhint(data.ui, data.litems)
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
