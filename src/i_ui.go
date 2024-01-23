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
 * Fri Jan 19 19:23:24 2024
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
	s           tcell.Screen
	mode        uint8
	style		[7]tcell.Style
	dim         [2]int
	err         [2]string
}

func i_left_right(text_len int, ui *HardUI) (int, int) {
	left  := (ui.dim[W] / 2) - text_len / 2
	right := ui.dim[W] - 1
	if left < 1 {
		left = 1
	}
	if right >= ui.dim[W] - 1 {
		right = ui.dim[W] - 1
	}
	return left, right
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

func i_draw_box(s tcell.Screen, x1, y1, x2, y2 int,
		box_style, head_style tcell.Style, title string, fill bool) {

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, box_style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, box_style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, box_style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, box_style)
	}
	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, box_style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, box_style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, box_style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, box_style)
	}
	if fill == true {
		for y := y1 + 1; y < y2; y++ {
			for x := x1 + 1; x < x2; x++ {
				s.SetContent(x, y, ' ', nil, box_style)
			}
		}
	}
	i_draw_text(s, x1 + 1, y1, x2 - 1, y1, head_style, title)
}

func i_draw_msg(s tcell.Screen, lines int, box_style tcell.Style,
	dim [2]int, title string) {

	lines += 1
	if lines < 0 {
		return
	}
	if lines > dim[H] - 2 {
		lines = dim[H] - 2
	}
	for row := dim[H] - 2 - lines; row < dim[H] - 2; row++ {
		s.SetContent(0, row, tcell.RuneVLine, nil, box_style)
		s.SetContent(dim[W] - 1, row, tcell.RuneVLine, nil, box_style)
	}
	for col := 1; col < dim[W] - 1; col++ {
		s.SetContent(col, dim[H] - 2 - lines, tcell.RuneHLine, nil, box_style)
		s.SetContent(col, dim[H] - 2, tcell.RuneHLine, nil, box_style)
	}
	s.SetContent(0, dim[H] - 2 - lines, tcell.RuneULCorner, nil, box_style)
	s.SetContent(dim[W] - 1, dim[H] - 2 - lines, tcell.RuneURCorner, nil,
		box_style)
	s.SetContent(0, dim[H] - 2, tcell.RuneLLCorner, nil, box_style)
	s.SetContent(dim[W] - 1, dim[H] - 2, tcell.RuneLRCorner, nil, box_style)
	// s.SetContent(dim[W] / 3, dim[H] - 2 - lines, tcell.RuneBTee, nil, )
	// s.SetContent(0, dim[H] - 2 - lines, tcell.RuneLTee, nil, )
	// s.SetContent(dim[W] - 1, dim[H] - 2 - lines, tcell.RuneRTee, nil, )
	for y := dim[H] - 2 - lines + 1; y < dim[H] - 2; y++ {
		for x := 1; x < dim[W] - 1; x++ {
			s.SetContent(x, y, ' ', nil, box_style)
		}
	}
	i_draw_text(s, 1, dim[H] - 2 - lines, len(title) + 2, dim[H] - 2 - lines,
		box_style, title)
}

func i_draw_bottom_text(ui HardUI) {
	text := ""

	switch ui.mode {
	case NORMAL_MODE:
		text = NORMAL_KEYS_HINTS
	case DELETE_MODE:
		text = DELETE_KEYS_HINTS
	case LOAD_MODE:
		text = ""
	case ERROR_MODE:
		text = ERROR_KEYS_HINTS
	}
	i_draw_text(ui.s,
		0, ui.dim[H] - 1, ui.dim[W], ui.dim[H] - 1,
		ui.style[DEF_STYLE].Dim(true), text)
	i_draw_text(ui.s,
		ui.dim[W] - 5, ui.dim[H] - 1, ui.dim[W], ui.dim[H] - 1,
		ui.style[DEF_STYLE].Dim(true), " " + VERSION)
}

func i_draw_zhosts_box(ui HardUI) {
	i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " No hosts ")
	text := "Hosts list empty. Add hosts/folders by pressing (a/m)"
	left, right := i_left_right(len(text), &ui)
	i_draw_text(ui.s, left, ui.dim[H] - 2 - 1, right, ui.dim[H] - 2 - 1,
		ui.style[DEF_STYLE], text)
}

func i_draw_delete_msg(ui HardUI, item *ItemsNode) {
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
	file = file[1:]
	i_draw_msg(ui.s, 2, ui.style[BOX_STYLE], ui.dim, " Delete ")
	left, right := i_left_right(len(text), &ui)
	line := ui.dim[H] - 2 - 2
	i_draw_text(ui.s, left, line, right, line, ui.style[DEF_STYLE], text)
	left, right = i_left_right(len(file), &ui)
	line += 1
	i_draw_text(ui.s,
	    left, line, right, line,
	    ui.style[DEF_STYLE].Bold(true), file)
}

func i_draw_error_msg(ui HardUI) {
	lines := 2
	if len(ui.err[ERROR_ERR]) == 0 {
		lines = 1
	}
	i_draw_msg(ui.s, lines, ui.style[BOX_STYLE], ui.dim, " Delete ")
	left, right := i_left_right(len(ui.err[ERROR_MSG]), &ui)
	line := ui.dim[H] - 2 - 2
	if len(ui.err[ERROR_ERR]) == 0 {
		line += 1
	}
	i_draw_text(ui.s, left, line, right, line,
		ui.style[ERR_STYLE], ui.err[ERROR_MSG])
	if len(ui.err[ERROR_ERR]) > 0 {
		left, right = i_left_right(len(ui.err[ERROR_ERR]), &ui)
		line += 1
		i_draw_text(ui.s, left, line, right, line,
			ui.style[ERR_STYLE], ui.err[ERROR_ERR])
	}
}

func i_draw_scrollhint(ui HardUI, litems *ItemsList) {
	if litems.head == nil {
		return
	}
	h := ui.dim[H] - 4
	last := litems.last.ID
	if last <= h {
		return
	}
	draw_id := litems.draw.ID
	if draw_id > 1 {
		ui.s.SetContent(0, 1,
			'▲',
			nil, ui.style[DEF_STYLE])
	}
	if last - draw_id > h {
		ui.s.SetContent(0, ui.dim[H] - 3,
			'▼',
			nil, ui.style[DEF_STYLE])
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
	i_draw_host_panel(*ui, false, nil, nil)
	i_draw_info_panel(*ui, false, nil)
	i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " Loading ")
	text := "Loading " + strconv.Itoa(g_load_count) + " hosts"
	left, right := i_left_right(len(text), ui)
	i_draw_text(ui.s,
		left, ui.dim[H] - 2 - 1, right, ui.dim[H] - 2 - 1, ui.style[DEF_STYLE], text)
	i_draw_text(ui.s,
		ui.dim[W] - 5, ui.dim[H] - 1, ui.dim[W], ui.dim[H] - 1,
		ui.style[DEF_STYLE].Dim(true), " " + VERSION)
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
	ui.style[DEF_STYLE] = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	ui.style[DIR_STYLE] = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorBlue).Dim(true).Bold(true)
	ui.style[BOX_STYLE] = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	ui.style[HEAD_STYLE] = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	ui.style[ERR_STYLE] = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorRed).Dim(true)
	ui.style[TITLE_STYLE] = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorBlue).Dim(true).Bold(true)
	ui.style[SEL_STYLE] = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorBlue).Dim(true).Bold(true)
	ui.s.SetStyle(ui.style[DEF_STYLE])
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
			i_draw_delete_msg(data.ui, data.litems.curr)
		}
		if data.ui.mode == ERROR_MODE {
			i_draw_error_msg(data.ui)
		}
		data.ui.s.Show()
		i_events(&data)
	}
}
