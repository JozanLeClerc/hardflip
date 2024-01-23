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
 * hardflip: src/i_info.go
 * Fri Jan 19 18:44:13 2024
 * Joe
 *
 * interfacing informations about items
 */

package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func i_info_panel_dirs(ui HardUI, dir *DirsNode) {
	line := 2
	if line > ui.dim[H] - 3 {
		return
	}

	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[TITLE_STYLE], "Name: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[DEF_STYLE], dir.Name)
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[TITLE_STYLE], "Type: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[DEF_STYLE], "Directory")
	if line += 2; line > ui.dim[H] - 3 {
		return
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[TITLE_STYLE], "Path: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[DEF_STYLE], dir.path()[1:])
}

func i_info_panel_host(ui HardUI, host *HostNode) {
	host_type := host.protocol_str()
	line := 2
	if line > ui.dim[H] - 3 {
		return
	}
	// name, type
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[TITLE_STYLE], "Name: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[DEF_STYLE], host.Name)
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[TITLE_STYLE], "Type: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[DEF_STYLE], host_type)
	if line += 2; line > ui.dim[H] - 3 {
		return
	}
	if line > ui.dim[H] - 3 {
		return
	}
	// host, port
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[TITLE_STYLE], "Host: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[DEF_STYLE], host.Host)
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[TITLE_STYLE], "Port: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[DEF_STYLE], strconv.Itoa(int(host.Port)))
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	// RDP shit
	if host.Protocol == 1 {
		if len(host.Domain) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
				ui.style[TITLE_STYLE], "Domain: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 11, line, ui.dim[W] - 2, line,
				ui.style[DEF_STYLE], host.Domain)
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
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[TITLE_STYLE], "User: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[DEF_STYLE], host.User)
	if line += 1; line > ui.dim[H] - 3 {
		return
	}
	if len(host.Pass) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Pass: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], "***")
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
	}
	if host.Protocol == 0 && len(host.Priv) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Privkey: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 12, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], host.Priv)
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
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Jump settings: ")
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Host: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], host.Jump)
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Port: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], strconv.Itoa(int(host.JumpPort)))
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "User: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], host.JumpUser)
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		if len(host.JumpPass) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
				ui.style[TITLE_STYLE], "Pass: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
				ui.style[DEF_STYLE], "***")
			if line += 1; line > ui.dim[H] - 3 {
				return
			}
		}
		if host.Protocol == 0 && len(host.JumpPriv) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
				ui.style[TITLE_STYLE], "Privkey: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 13, line, ui.dim[W] - 2, line,
				ui.style[DEF_STYLE], host.JumpPriv)
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
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Screen size: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 16, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE],
			strconv.Itoa(int(host.Width)) + "x" +
			strconv.Itoa(int(host.Height)))
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Dynamic window: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 19, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], strconv.FormatBool(host.Dynamic))
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Quality: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 12, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], qual[host.Quality])
		line += 1
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
	}
	// note
	if len(host.Note) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[TITLE_STYLE], "Note: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
			ui.style[DEF_STYLE], host.Note)
		if line += 1; line > ui.dim[H] - 3 {
			return
		}
	}
}

func i_draw_info_panel(ui HardUI, percent bool, litems *ItemsList) {
	i_draw_box(ui.s, (ui.dim[W] / 3), 0,
		ui.dim[W] - 1, ui.dim[H] - 2,
		ui.style[BOX_STYLE], ui.style[HEAD_STYLE], " Infos ", false)
	ui.s.SetContent(ui.dim[W] / 3, 0, tcell.RuneTTee, nil, ui.style[BOX_STYLE])
	ui.s.SetContent(ui.dim[W] / 3, ui.dim[H] - 2,
		tcell.RuneBTee, nil, ui.style[BOX_STYLE])
	if litems == nil {
		return
	}
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
			ui.style[DEF_STYLE],
			text)
	} else {
		text := " 0 hosts "
		i_draw_text(ui.s,
			(ui.dim[W] - 1) - len(text) - 1,
			ui.dim[H] - 2,
			(ui.dim[W] - 1) - 1,
			ui.dim[H] - 2,
			ui.style[DEF_STYLE],
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
