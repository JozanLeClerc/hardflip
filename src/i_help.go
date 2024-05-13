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
 * hardflip: src/i_help.go
 * Fri May 10 15:40:13 2024
 * Joe
 *
 * helping the user
 */

package main

func i_draw_help(ui *HardUI) {
	if ui.dim[W] < 12 || ui.dim[H] < 6 {
		return
	}
	win := Quad{
		6,
		3,
		ui.dim[W] - 6,
		ui.dim[H] - 3,
	}
	if win.B > win.T + len(HELP_NORMAL_KEYS) + 1 {
		win.B = win.T + len(HELP_NORMAL_KEYS) + 1
	}
	i_draw_box(ui.s,
		win.L, win.T, win.R, win.B,
		ui.style[BOX_STYLE], ui.style[HEAD_STYLE],
		" Keys ", true)
	line := 0
	line -= ui.help_scroll
	if line < 0 {
		ui.s.SetContent(win.L, win.T + 1, '▲',
			nil, ui.style[BOX_STYLE])
	}
	ui.help_end = i_help_normal(*ui, win, &line)
	if ui.help_end == false {
		ui.s.SetContent(win.L, win.B - 1, '▼',
			nil, ui.style[BOX_STYLE])
	}
	// TODO: here
}

func i_help_normal(ui HardUI, win Quad, line *int) bool {
	delim := 17
	for _, v := range HELP_NORMAL_KEYS {
		if *line < 0 {
			*line += 1
			continue
		} else if win.T + *line + 1 >= win.B {
			return false
		}
		i := delim - 1 - len(v[0])
		if i < 0 { i = 0 }
		i_draw_text(ui.s, win.L + 1 + i, win.T + 1 + *line, win.L + delim,
			win.T + 1 + *line, ui.style[BOT_STYLE], v[0])
		i_draw_text(ui.s, win.L + delim + 1, win.T + 1 + *line, win.R,
			win.T + 1 +  *line, ui.style[DEF_STYLE], v[1])
		*line += 1
	}
	return true
}
