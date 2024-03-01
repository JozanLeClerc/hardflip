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
 * Fri Mar 01 15:27:17 2024
 * Joe
 *
 * interfacing with the user
 */

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type HardUI struct {
	s     tcell.Screen
	mode  uint8
	style [7]tcell.Style
	dim   [2]int
	err   [2]string
	buff  string
	insert_sel int
	insert_sel_max int
	insert_sel_ok bool
}

type Quad struct {
	L, T, R, B int
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
	for y := dim[H] - 2 - lines + 1; y < dim[H] - 2; y++ {
		for x := 1; x < dim[W] - 1; x++ {
			s.SetContent(x, y, ' ', nil, box_style)
		}
	}
	i_draw_text(s, 1, dim[H] - 2 - lines, len(title) + 2, dim[H] - 2 - lines,
		box_style, title)
}

func i_draw_bottom_text(ui HardUI, opts HardOpts, insert *HostNode) {
	text := ""

	switch ui.mode {
	case NORMAL_MODE:
		text = NORMAL_KEYS_HINTS
	case DELETE_MODE:
		text = CONFIRM_KEYS_HINTS
	case LOAD_MODE:
		text = "Loading..."
	case ERROR_MODE:
		text = ERROR_KEYS_HINTS
	case WELCOME_MODE:
		if len(opts.GPG) == 0 {
			text = ""
		} else {
			text = CONFIRM_KEYS_HINTS
		}
	case INSERT_MODE:
		if insert == nil {
			text = ""
		} else {
			text = INSERT_KEYS_HINTS
		}
	default:
		text = ""
	}
	i_draw_text(ui.s,
		1, ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[BOT_STYLE], text)
	i_draw_text(ui.s,
		ui.dim[W] - len(VERSION) - 2, ui.dim[H] - 1,
		ui.dim[W] - 1, ui.dim[H] - 1, ui.style[BOT_STYLE], " " + VERSION)
}

func i_draw_welcome_box(ui HardUI) {
	l_max, r_max := ui.dim[W] / 8 + 1, ui.dim[W] - ui.dim[W] / 8 - 1
	b_max := ui.dim[H] / 2 - 1
	i_draw_box(ui.s,
		l_max - 1, 0, r_max, b_max + 1,
		ui.style[BOX_STYLE], ui.style[HEAD_STYLE], "", true)
	art := [4]string{
		` _     __`,
		`| |_  / _|`,
		`| ' \|  _|`,
		`|_||_|_|`,
	}
	line := 0
	for k, v := range art {
		if k + 1 > b_max { break }
		line = k + 1
		l, r := (ui.dim[W] / 2) - 6, ui.dim[W]
		if l < l_max { l = l_max }; if r > r_max { r = r_max }
		i_draw_text(ui.s,
			l, k + 1, r, k + 1,
			ui.style[DEF_STYLE], v)
	}
	if line > b_max { return }
	text := "hf " + VERSION
	if len(VERSION_NAME) > 0 {
		text += " - " + VERSION_NAME
	}
	l, r := ui.dim[W] / 2 - len(text) / 2 + 7,
			ui.dim[W] / 2 + len(text) / 2 + 1 + 7
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text)
	if line += 2; line > b_max { return }
	text = `Welcome to hardflip!`
	l, r = ui.dim[W] / 2 - len(text) / 2, ui.dim[W] / 2 + len(text) / 2 + 1
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text)
	text = `Please select the gpg key ID to be used`
	if line += 1; line > b_max { return }
	l, r = ui.dim[W] / 2 - len(text) / 2, ui.dim[W] / 2 + len(text) / 2 + 1
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text)
	text = `for password encryption`
	if line += 1; line > b_max { return }
	l, r = ui.dim[W] / 2 - len(text) / 2, ui.dim[W] / 2 + len(text) / 2 + 1
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text)
	text = `Set gpg key can be modified in the config file`
	if line += 1; line > b_max { return }
	l, r = ui.dim[W] / 2 - len(text) / 2, ui.dim[W] / 2 + len(text) / 2 + 1
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text)
	text = `If you don't want to use GnuPG for password`
	if line += 2; line > b_max { return }
	l, r = ui.dim[W] / 2 - len(text) / 2, ui.dim[W] / 2 + len(text) / 2 + 1
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text)
	text = `storage, please select `
	text_2 := `plain`
	text_3 := ` (plaintext passwords`
	if line += 1; line > b_max { return }
	l = ui.dim[W] / 2 - len(text + text_2 + text_3) / 2
	r = r_max
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text)
	l = l + len(text)
	r = l + len(text_2)
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	if l >= r_max { return }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE].Bold(true), text_2)
	l = l + len(text_2)
	r = l + len(text_3)
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text_3)
	text = `are not recommended)`
	if line += 1; line > b_max { return }
	l, r = ui.dim[W] / 2 - len(text) / 2, ui.dim[W] / 2 + len(text) / 2 + 1
	if l < l_max { l = l_max }; if r > r_max { r = r_max }
	i_draw_text(ui.s, l, line, r, line, ui.style[DEF_STYLE], text)
}

func i_prompt_gpg(ui HardUI, keys [][2]string) {
	lines := len(keys)
	if lines == 1 {
		lines = 2
	}
	i_draw_msg(ui.s, lines, ui.style[BOX_STYLE], ui.dim, " GnuPG keys ")
	for k, v := range keys {
		text := ""
		if v[0] != "plain" {
			text = "[" + strconv.Itoa(k + 1) + "] " +
				v[1] + " " + v[0][:10] + "... "
		} else {
			text = "[" + strconv.Itoa(k + 1) + "] " + "plain"
		}
		line := ui.dim[H] - 2 - len(keys) + k
		i_draw_text(ui.s, 2, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], text)
	}
	if len(keys) == 1 {
		i_draw_text(ui.s, 2, ui.dim[H] - 4, ui.dim[W] - 1, ui.dim[H] - 4,
			ui.style[DEF_STYLE],
			"No gpg key! Creating your gpg key first is recommended")
	}
	i_draw_text(ui.s,
		1, ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], "gpg: ")
	ui.s.ShowCursor(6, ui.dim[H] - 1)
}

func i_prompt_confirm_gpg(ui HardUI, opts HardOpts) {
	if opts.GPG == "plain" {
		i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " Confirm plaintext ")
		text := "Really use plaintext to store passwords?"
		l, r := i_left_right(len(text), &ui)
		i_draw_text(ui.s, l, ui.dim[H] - 3, r, ui.dim[H] - 3,
			ui.style[DEF_STYLE], text)
		return
	}
	i_draw_msg(ui.s, 2, ui.style[BOX_STYLE], ui.dim, " Confirm GnuPG key ")
	text := "Really use this gpg key?"
	l, r := i_left_right(len(text), &ui)
	i_draw_text(ui.s, l, ui.dim[H] - 4, r, ui.dim[H] - 4,
		ui.style[DEF_STYLE], text)
	l, r = i_left_right(len(opts.GPG), &ui)
	i_draw_text(ui.s, l, ui.dim[H] - 3, r, ui.dim[H] - 3,
		ui.style[DEF_STYLE], opts.GPG)
}

func i_prompt_mkdir(ui HardUI, curr *ItemsNode) {
	path := "/"
	if curr != nil {
		path = curr.path()
	}
	path = path[1:]
	prompt := "mkdir: "
	i_draw_text(ui.s,
		1, ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], prompt)
	i_draw_text(ui.s, len(prompt) + 1,
		ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], path)
	i_draw_text(ui.s, len(prompt) + 1 + len(path),
		ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE].Bold(true), ui.buff)
	ui.s.ShowCursor(len(prompt) + 1 + len(path) + len(ui.buff), ui.dim[H] - 1)
}

func i_prompt_type(ui HardUI) {
	i_draw_msg(ui.s, 4, ui.style[BOX_STYLE], ui.dim, " Connection type ")
	i_draw_text(ui.s, 2, ui.dim[H] - 6, ui.dim[W] - 2, ui.dim[H] - 6,
		ui.style[DEF_STYLE], "[1] SSH")
	i_draw_text(ui.s, 2, ui.dim[H] - 5, ui.dim[W] - 2, ui.dim[H] - 5,
		ui.style[DEF_STYLE], "[2] RDP")
	i_draw_text(ui.s, 2, ui.dim[H] - 4, ui.dim[W] - 2, ui.dim[H] - 4,
		ui.style[DEF_STYLE], "[3] Single command")
	i_draw_text(ui.s, 2, ui.dim[H] - 3, ui.dim[W] - 2, ui.dim[H] - 3,
		ui.style[DEF_STYLE], "[4] OpenStack CLI")
	text := "Type: "
	i_draw_text(ui.s, 0,
		ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], text)
	ui.s.ShowCursor(len(text), ui.dim[H] - 1)
}

func i_prompt_generic(ui HardUI, prompt string, secret, file bool) {
	i_draw_text(ui.s,
		1, ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], prompt)
	if secret == true {
		ui.s.ShowCursor(len(prompt) + 1, ui.dim[H] - 1)
		return
	}
	i_draw_text(ui.s, len(prompt) + 1,
		ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE].Bold(true), ui.buff)
	ui.s.ShowCursor(len(prompt) + 1 + len(ui.buff), ui.dim[H] - 1)
}

func i_prompt_insert(ui HardUI, curr *ItemsNode) {
	path := "/"
	path = curr.path()
	path = path[1:]
	prompt := "Name: "
	i_draw_text(ui.s,
		1, ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], prompt)
	i_draw_text(ui.s, len(prompt) + 1,
		ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], path)
	i_draw_text(ui.s, len(prompt) + 1 + len(path),
		ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE].Bold(true), ui.buff)
	ui.s.ShowCursor(len(prompt) + 1 + len(path) + len(ui.buff), ui.dim[H] - 1)
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

	file = item.path()
	if item.is_dir() == true {
		text = "Really delete this directory and all of its content?"
	} else {
		text = "Really delete this host?"
		file += item.Host.Filename
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

func i_draw_load_error_msg(ui HardUI, load_err []error) {
	lines := len(load_err)
	i_draw_msg(ui.s, lines, ui.style[BOX_STYLE], ui.dim, " Load time errors ")
	left, right := 1, ui.dim[W] - 1
	line := ui.dim[H] - 2 - 1 - len(load_err)
	if line < 0 {
		line = 0
	}
	for _, err := range load_err {
		line += 1
		err_str := fmt.Sprintf("%v", err)
		i_draw_text(ui.s, left, line, right, line,
			ui.style[ERR_STYLE], err_str)
	}
}

func i_draw_error_msg(ui HardUI, load_err []error) {
	if len(load_err) > 0 {
		i_draw_load_error_msg(ui, load_err)
		return
	}
	lines := 2
	if len(ui.err[ERROR_ERR]) == 0 {
		lines = 1
	}
	i_draw_msg(ui.s, lines, ui.style[BOX_STYLE], ui.dim, " Error ")
	left, right := 1, ui.dim[W] - 2
	line := ui.dim[H] - 2 - 2
	if len(ui.err[ERROR_ERR]) == 0 {
		line += 1
	}
	i_draw_text(ui.s, left, line, right, line,
		ui.style[ERR_STYLE], ui.err[ERROR_MSG])
	if len(ui.err[ERROR_ERR]) > 0 {
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
			nil, ui.style[BOX_STYLE])
	}
	if last - draw_id > h {
		ui.s.SetContent(0, ui.dim[H] - 3,
			'▼',
			nil, ui.style[BOX_STYLE])
		return
	}
}

var g_load_count int = -1

func i_draw_load_ui(ui *HardUI, opts HardOpts) {
	g_load_count += 1
	if g_load_count % 1000 != 0 {
		return
	}
	i_draw_host_panel(*ui, false, nil, nil)
	i_draw_info_panel(*ui, false, nil)
	text := ""
	for i := 0; i < ui.dim[W] - 1; i++ {
		text += " "
	}
	i_draw_text(ui.s, 1, ui.dim[H] - 1, ui.dim[W], ui.dim[H] - 1,
		ui.style[BOT_STYLE], text)
	i_draw_bottom_text(*ui, opts, nil)
	i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " Loading ")
	text = "Loading " + strconv.Itoa(g_load_count) + " hosts"
	left, right := i_left_right(len(text), ui)
	i_draw_text(ui.s,
		left, ui.dim[H] - 2 - 1, right, ui.dim[H] - 2 - 1,
		ui.style[DEF_STYLE], text)
	ui.s.Show()
	ui.s.PostEvent(nil)
	event := ui.s.PollEvent()
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
			   ui *HardUI,
			   load_err *[]error) (*DirsList, *ItemsList, []error) {
	ui.mode = LOAD_MODE
	ldirs := c_load_data_dir(data_dir, opts, ui, load_err)
	litems := c_load_litems(ldirs)
	if ui.mode != ERROR_MODE {
		ui.mode = NORMAL_MODE
	}
	if len(*load_err) == 0 {
		*load_err = nil
	}
	return ldirs, litems, *load_err
}

func i_ui(data_dir string) {
	ui := HardUI{}
	opts := HardOpts{}
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
	ui.style[BOT_STYLE] = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorBlue).Dim(true)
	ui.s.SetStyle(ui.style[DEF_STYLE])
	ui.dim[W], ui.dim[H], _ = term.GetSize(0)
	var load_err []error
	conf_dir  := c_get_conf_dir(&load_err)
	if conf_dir == "" {
		opts = DEFAULT_OPTS
	} else {
		opts = c_get_options(conf_dir, &load_err)
	}
	ldirs, litems, load_err := i_load_ui(data_dir, opts, &ui, &load_err)
	data := HardData{
		litems,
		ldirs,
		ui,
		opts,
		make(map[*DirsNode]*ItemsList),
		data_dir,
		load_err,
		[][2]string{},
		nil,
	}
	if data.opts.GPG == DEFAULT_OPTS.GPG && data.litems.head == nil {
		data.ui.mode = WELCOME_MODE
		data.keys = c_get_secret_gpg_keyring(&data.ui)
	}
	for {
		data.ui.s.Clear()
		i_draw_bottom_text(data.ui, data.opts, data.insert)
		i_draw_host_panel(data.ui, data.opts.Icon, data.litems, &data)
		i_draw_info_panel(data.ui, data.opts.Perc, data.litems)
		i_draw_scrollhint(data.ui, data.litems)
		if data.load_err != nil && len(data.load_err) > 0 {
			data.ui.mode = ERROR_MODE
		}
		if data.ui.mode == WELCOME_MODE {
			i_draw_welcome_box(data.ui)
			if len(data.opts.GPG) == 0 {
				i_prompt_gpg(data.ui, data.keys)
			} else {
				i_prompt_confirm_gpg(data.ui, data.opts)
			}
		} else if data.litems.head == nil {
			i_draw_zhosts_box(data.ui)
		}
		if data.ui.mode == DELETE_MODE {
			i_draw_delete_msg(data.ui, data.litems.curr)
		} else if data.ui.mode == ERROR_MODE {
			i_draw_error_msg(data.ui, data.load_err)
		} else if data.ui.mode == MKDIR_MODE {
			i_prompt_mkdir(data.ui, data.litems.curr)
		} else if data.ui.mode == INSERT_MODE {
			if data.insert == nil {
				i_prompt_insert(data.ui, data.litems.curr)
			} else {
				i_draw_insert_panel(data.ui, data.insert)
				if data.ui.insert_sel_ok == true {
					switch data.ui.insert_sel {
					case 0:
						i_prompt_type(data.ui)
					case 1, 6:
						i_prompt_generic(data.ui, "Host/IP: ", false, false)
					case 2:
						i_prompt_generic(data.ui, "Port: ", false, false)
					case 3:
						i_prompt_generic(data.ui, "User: ", false, false)
					case 4:
						i_prompt_generic(data.ui, "Pass: ", true, false)
					case 5:
						i_prompt_generic(data.ui, "Private key: ", false, true)
					}
				}
			}
		}
		data.ui.s.Show()
		i_events(&data)
	}
}
