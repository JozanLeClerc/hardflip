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
 * Thu Apr 25 16:20:41 2024
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

type Buffer struct {
	data []rune
	cursor int
}

type HardUI struct {
	s     tcell.Screen
	mode  uint8
	style [STYLE_MAX + 1]tcell.Style
	dim   [2]int
	err   [2]string
	buff  Buffer
	drives_buff string
	msg_buff string
	match_buff string
	insert_sel int
	insert_sel_max int
	insert_sel_ok bool
	insert_method int
	insert_scroll int
	insert_butt bool
	help_scroll int
	help_end bool
	welcome_screen int
}

type Quad struct {
	L, T, R, B int
}

func (buffer *Buffer)empty() {
	buffer.data = []rune{}
	buffer.cursor = 0
}

func (buffer *Buffer)insert(str string) {
	buffer.data = []rune(str)
	buffer.cursor = len(buffer.data)
}

func (buffer *Buffer)str() string {
	return string(buffer.data)
}

func (buffer *Buffer)len() int {
	return len(buffer.data)
}

func i_left_right(text_len int, ui HardUI) (int, int) {
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

func i_draw_bottom_text(ui HardUI, insert *HostNode, insert_err []error) {
	text := ""

	if len(ui.msg_buff) > 0 {
		text = ui.msg_buff
	} else {
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
			if ui.welcome_screen == WELCOME_CONFIRM_GPG {
				text = CONFIRM_KEYS_HINTS
			}
		case INSERT_MODE:
			if insert == nil {
				text = ""
			} else if insert_err != nil {
				text = ERROR_KEYS_HINTS
			}
		case HELP_MODE:
			text = HELP_KEYS_HINTS
		default:
			text = ""
		}
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
	l, r := ui.dim[W] / 2 + 4,
			ui.dim[W] - 2
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
		l, r := i_left_right(len(text), ui)
		i_draw_text(ui.s, l, ui.dim[H] - 3, r, ui.dim[H] - 3,
			ui.style[DEF_STYLE], text)
		return
	}
	i_draw_msg(ui.s, 2, ui.style[BOX_STYLE], ui.dim, " Confirm GnuPG key ")
	text := "Really use this gpg key?"
	l, r := i_left_right(len(text), ui)
	i_draw_text(ui.s, l, ui.dim[H] - 4, r, ui.dim[H] - 4,
		ui.style[DEF_STYLE], text)
	l, r = i_left_right(len(opts.GPG), ui)
	i_draw_text(ui.s, l, ui.dim[H] - 3, r, ui.dim[H] - 3,
		ui.style[DEF_STYLE], opts.GPG)
}

func i_prompt_def_sshkey(ui HardUI, home_dir string) {
	i_draw_msg(ui.s, 4, ui.style[BOX_STYLE], ui.dim, " Default SSH key ")
	text := "Please enter here a path for your most used SSH key"
	l, r := i_left_right(len(text), ui)
	i_draw_text(ui.s, l, ui.dim[H] - 6, r, ui.dim[H] - 6,
		ui.style[DEF_STYLE],text)
	text = "It will be entered by default when adding SSH hosts"
	l, r = i_left_right(len(text), ui)
	i_draw_text(ui.s, l, ui.dim[H] - 5, r, ui.dim[H] - 5,
		ui.style[DEF_STYLE],text)
	text = "This can save some time"
	l, r = i_left_right(len(text), ui)
	i_draw_text(ui.s, l, ui.dim[H] - 4, r, ui.dim[H] - 4,
		ui.style[DEF_STYLE],text)
	text = "Leave empty if you don't want to set a default key"
	l, r = i_left_right(len(text), ui)
	i_draw_text(ui.s, l, ui.dim[H] - 3, r, ui.dim[H] - 3,
		ui.style[DEF_STYLE],text)
	i_prompt_generic(ui, "SSH private key: ", false, home_dir)
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
		ui.style[DEF_STYLE].Bold(true), ui.buff.str())
	ui.s.ShowCursor(len(prompt) + 1 + len(path) +
		ui.buff.cursor, ui.dim[H] - 1)
}

func i_prompt_list(ui HardUI, name, prompt string, list []string) {
	i := len(list)
	i_draw_msg(ui.s, i, ui.style[BOX_STYLE], ui.dim, " " + name + " ")
	for k, v := range list {
		i_draw_text(ui.s, 2, ui.dim[H] - 2 - i,
					ui.dim[W] - 2, ui.dim[H] - 2 - i,
					ui.style[DEF_STYLE], "[" + strconv.Itoa(k + 1) + "] " + v)
		i -= 1
	}
	i_draw_text(ui.s, 1,
				ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
				ui.style[DEF_STYLE], prompt)
	ui.s.ShowCursor(len(prompt) + 2, ui.dim[H] - 1)
}

func i_prompt_generic(ui HardUI, prompt string, secret bool, home_dir string) {
	i_draw_text(ui.s,
		1, ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], prompt)
	if secret == true {
		ui.s.ShowCursor(len(prompt) + 1, ui.dim[H] - 1)
		return
	}
	style := ui.style[DEF_STYLE].Bold(true)
	if len(home_dir) > 0 && ui.buff.len() > 0 {
		file := ui.buff.str()
		if file[0] == '~' {
			file = home_dir + file[1:]
		}
		if stat, err := os.Stat(file);
		   err != nil {
			style = style.Foreground(tcell.ColorRed)
		} else if stat.IsDir() == true {
			style = style.Foreground(tcell.ColorPurple).
				Bold(false).
				Underline(true)
		} else {
			style = style.Foreground(tcell.ColorGreen).Bold(false)
		}
	}
	i_draw_text(ui.s, len(prompt) + 1,
		ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		style, ui.buff.str())
	ui.s.ShowCursor(len(prompt) + 1 + ui.buff.cursor, ui.dim[H] - 1)
}

func i_prompt_dir(ui HardUI, prompt string, home_dir string) {
	i_draw_text(ui.s,
		1, ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		ui.style[DEF_STYLE], prompt)
	style := ui.style[DEF_STYLE].Bold(true)
	if len(home_dir) > 0 && ui.buff.len() > 0 {
		file := ui.buff.str()
		if file[0] == '~' {
			file = home_dir + file[1:]
		}
		if stat, err := os.Stat(file);
		   err != nil {
			style = style.Foreground(tcell.ColorRed)
		} else if stat.IsDir() == true {
			style = style.Foreground(tcell.ColorGreen).Bold(false)
		} else {
			style = style.Foreground(tcell.ColorRed)
		}
	}
	i_draw_text(ui.s, len(prompt) + 1,
		ui.dim[H] - 1, ui.dim[W] - 1, ui.dim[H] - 1,
		style, ui.buff.str())
	ui.s.ShowCursor(len(prompt) + 1 + ui.buff.cursor, ui.dim[H] - 1)
}

func i_prompt_insert(ui HardUI, curr *ItemsNode) {
	path := "/"
	if curr != nil {
		if ui.mode == RENAME_MODE {
			if curr.is_dir() == false {
				path = curr.path()
			} else {
				path = curr.Dirs.Parent.path()
			}
		} else {
			path = curr.path()
		}
	}
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
		ui.style[DEF_STYLE].Bold(true), ui.buff.str())
	ui.s.ShowCursor(len(prompt) + 1 + len(path) +
		ui.buff.cursor, ui.dim[H] - 1)
}

func i_draw_remove_share(ui HardUI) {
	text := "Really remove this share?"

	i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " Remove share ")
	left, right := i_left_right(len(text), ui)
	line := ui.dim[H] - 2 - 1
	i_draw_text(ui.s, left, line, right, line, ui.style[DEF_STYLE], text)
}

func i_draw_zhosts_box(ui HardUI) {
	i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " No hosts ")
	text := "Hosts list empty. Add hosts/folders by pressing (a/m)"
	left, right := i_left_right(len(text), ui)
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
		file += item.Host.filename
	}
	file = file[1:]
	i_draw_msg(ui.s, 2, ui.style[BOX_STYLE], ui.dim, " Delete ")
	left, right := i_left_right(len(text), ui)
	line := ui.dim[H] - 2 - 2
	i_draw_text(ui.s, left, line, right, line, ui.style[DEF_STYLE], text)
	left, right = i_left_right(len(file), ui)
	line += 1
	i_draw_text(ui.s,
	    left, line, right, line,
	    ui.style[DEF_STYLE].Bold(true), file)
}

func i_draw_insert_err_msg(ui HardUI, insert_err []error) {
	lines := len(insert_err)
	i_draw_msg(ui.s, lines, ui.style[BOX_STYLE], ui.dim, " Errors ")
	left, right := 1, ui.dim[W] - 1
	line := ui.dim[H] - 2 - 1 - len(insert_err)
	if line < 0 {
		line = 0
	}
	for _, err := range insert_err {
		line += 1
		err_str := fmt.Sprintf("%v", err)
		i_draw_text(ui.s, left, line, right, line,
			ui.style[ERR_STYLE], err_str)
	}
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

func i_draw_match_buff(ui HardUI) {
	i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, "")
	i_draw_text(ui.s, 2, ui.dim[H] - 2 - 1, ui.dim[W] - 2, ui.dim[H] - 2 - 1,
				ui.style[DEF_STYLE], ui.match_buff)
}

var g_load_count int = -1

func i_draw_load_ui(ui *HardUI) {
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
	i_draw_bottom_text(*ui, nil, nil)
	i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " Loading ")
	text = "Loading " + strconv.Itoa(g_load_count) + " hosts"
	left, right := i_left_right(len(text), *ui)
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

func i_init_styles(ui *HardUI, styles HardStyle) {
	for i := 0; i < STYLE_MAX + 1; i++ {
		tmp := tcell.StyleDefault.Background(tcell.ColorReset)
		curr_color := "default"
		switch i {
		case DEF_STYLE:
			curr_color = styles.DefColor
		case DIR_STYLE:
			curr_color = styles.DirColor
		case BOX_STYLE:
			curr_color = styles.BoxColor
		case HEAD_STYLE:
			curr_color = styles.HeadColor
		case ERR_STYLE:
			curr_color = styles.ErrColor
		case TITLE_STYLE:
			curr_color = styles.TitleColor
		case BOT_STYLE:
			curr_color = styles.BotColor
		case YANK_STYLE:
			curr_color = styles.YankColor
		case MOVE_STYLE:
			curr_color = styles.MoveColor
		default:
			curr_color = "default"
		}
		switch curr_color {
		case COLORS[COLOR_DEFAULT]:
			ui.style[i] = tmp.Foreground(tcell.ColorReset)
		case COLORS[COLOR_BLACK]:
			ui.style[i] = tmp.Foreground(tcell.ColorBlack)
		case COLORS[COLOR_RED]:
			ui.style[i] = tmp.Foreground(tcell.ColorRed).Dim(true)
		case COLORS[COLOR_GREEN]:
			ui.style[i] = tmp.Foreground(tcell.ColorGreen)
		case COLORS[COLOR_YELLOW]:
			ui.style[i] = tmp.Foreground(tcell.ColorYellow).Dim(true)
		case COLORS[COLOR_BLUE]:
			ui.style[i] = tmp.Foreground(tcell.ColorBlue).Dim(true)
		case COLORS[COLOR_MAGENTA]:
			ui.style[i] = tmp.Foreground(tcell.ColorPurple)
		case COLORS[COLOR_CYAN]:
			ui.style[i] = tmp.Foreground(tcell.ColorTeal)
		case COLORS[COLOR_WHITE]:
			ui.style[i] = tmp.Foreground(tcell.ColorWhite).Dim(true)
		case COLORS[COLOR_GRAY]:
			ui.style[i] = tmp.Foreground(tcell.ColorGray)
		case COLORS[COLOR_BOLD_BLACK]:
			ui.style[i] = tmp.Foreground(tcell.ColorBlack).Bold(true)
		case COLORS[COLOR_BOLD_RED]:
			ui.style[i] = tmp.Foreground(tcell.ColorRed).Dim(true).Bold(true)
		case COLORS[COLOR_BOLD_GREEN]:
			ui.style[i] = tmp.Foreground(tcell.ColorGreen).Dim(true).Bold(true)
		case COLORS[COLOR_BOLD_YELLOW]:
			ui.style[i] = tmp.Foreground(tcell.ColorYellow).Dim(true).Bold(true)
		case COLORS[COLOR_BOLD_BLUE]:
			ui.style[i] = tmp.Foreground(tcell.ColorBlue).Dim(true).Bold(true)
		case COLORS[COLOR_BOLD_MAGENTA]:
			ui.style[i] = tmp.Foreground(tcell.ColorPurple).Bold(true)
		case COLORS[COLOR_BOLD_CYAN]:
			ui.style[i] = tmp.Foreground(tcell.ColorTeal).Dim(true).Bold(true)
		case COLORS[COLOR_BOLD_WHITE]:
			ui.style[i] = tmp.Foreground(tcell.ColorWhite).Dim(true).Bold(true)
		case COLORS[COLOR_BOLD_GRAY]:
			ui.style[i] = tmp.Foreground(tcell.ColorGray).Bold(true)
		default:
			ui.style[i] = tmp.Foreground(tcell.ColorReset)
		}
	}
}

type key_event_mode_func func(*HardData, *HardUI, tcell.EventKey) bool

func i_ui(data_dir string) {
	home_dir, _ := os.UserHomeDir()
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
	ui.dim[W], ui.dim[H], _ = term.GetSize(0)
	var load_err []error
	conf_dir  := c_get_conf_dir(&load_err)
	if len(conf_dir) == 0 {
		opts = DEFAULT_OPTS
	} else {
		opts = c_get_options(conf_dir, &load_err)
	}
	styles := c_get_styles(conf_dir, &load_err)
	i_init_styles(&ui, styles)
	ui.s.SetStyle(ui.style[DEF_STYLE])
	ldirs, litems, load_err := i_load_ui(data_dir, opts, &ui, &load_err)
	data := HardData{
		litems,
		ldirs,
		ui,
		opts,
		styles,
		make(map[*DirsNode]*ItemsList),
		data_dir,
		home_dir,
		load_err,
		nil,
		[][2]string{},
		nil,
		nil,
	}
	if data.opts.GPG == DEFAULT_OPTS.GPG && data.litems.head == nil {
		data.ui.mode = WELCOME_MODE
		data.keys = c_get_secret_gpg_keyring()
	}
	fp := [MODE_MAX + 1]key_event_mode_func{
		NORMAL_MODE:	e_normal_events,
		DELETE_MODE:	e_delete_events,
		LOAD_MODE:		e_load_events,
		ERROR_MODE:		e_error_events,
		WELCOME_MODE:	e_welcome_events,
		MKDIR_MODE:		e_mkdir_events,
		INSERT_MODE:	e_insert_events,
		RENAME_MODE:	e_rename_events,
		HELP_MODE:		e_help_events,
	}
	for {
		data.ui.s.Clear()
		i_draw_bottom_text(data.ui, data.insert, data.insert_err)
		i_draw_host_panel(data.ui, data.opts.Icon, data.litems, &data)
		i_draw_info_panel(data.ui, data.opts.Perc, data.litems)
		i_draw_scrollhint(data.ui, data.litems)
		if data.load_err != nil && len(data.load_err) > 0 {
			data.ui.mode = ERROR_MODE
		}
		switch data.ui.mode {
		case WELCOME_MODE:
			i_draw_welcome_box(data.ui)
			switch data.ui.welcome_screen {
			case WELCOME_GPG:
				i_prompt_gpg(data.ui, data.keys)
			case WELCOME_CONFIRM_GPG:
				i_prompt_confirm_gpg(data.ui, data.opts)
			case WELCOME_SSH:
				i_prompt_def_sshkey(data.ui, data.home_dir)
			}
		case NORMAL_MODE:
			if data.litems.head == nil {
				i_draw_zhosts_box(data.ui)
			}
		case DELETE_MODE:
			i_draw_delete_msg(data.ui, data.litems.curr)
		case ERROR_MODE:
			i_draw_error_msg(data.ui, data.load_err)
		case MKDIR_MODE:
			i_prompt_mkdir(data.ui, data.litems.curr)
		case INSERT_MODE:
			if data.insert == nil {
				i_prompt_insert(data.ui, data.litems.curr)
			} else {
				i_draw_insert_panel(&data.ui, data.insert, data.home_dir)
				if data.insert_err != nil {
					i_draw_insert_err_msg(data.ui, data.insert_err)
				}
			}
		case RENAME_MODE:
			i_prompt_insert(data.ui, data.litems.curr)
		case HELP_MODE:
			i_draw_help(&data.ui)
		}
		if len(data.ui.match_buff) > 0 {
			i_draw_match_buff(data.ui)
			data.ui.match_buff = ""
		}
		data.ui.s.Show()
		e_events(&data, fp)
	}
}
