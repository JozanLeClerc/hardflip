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
 * hardflip: src/i_host.go
 * Fri Jan 19 12:52:11 2024
 * Joe
 *
 * interfacing hosts
 */

package main

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

func i_draw_host_panel(ui HardUI, icons bool, litems *ItemsList, data *HardData) {
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
