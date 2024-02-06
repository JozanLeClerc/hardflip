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

func i_info_dirs(ui HardUI, dir *DirsNode) {
	line := 2
	if line > ui.dim[H] - 3 { return }

	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Name: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], dir.Name)
	if line += 1; line > ui.dim[H] - 3 { return }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Type: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], "Directory")
	if line += 2; line > ui.dim[H] - 3 { return }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Path: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], dir.path()[1:])
}

func i_info_name_type(ui HardUI, host *HostNode) int {
	line := 2
	if line > ui.dim[H] - 3 { return line }
	host_type := host.protocol_str()
	// name, type
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Name: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.Name)
	if line += 1; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Type: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host_type)
	if line += 2; line > ui.dim[H] - 3 { return line }
	return line
}

func i_info_ssh(ui HardUI, host *HostNode, line int) int {
	if line > ui.dim[H] - 3 { return line }
	// host, port
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Host: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.Host)
	if line += 1; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Port: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], strconv.Itoa(int(host.Port)))
	if line += 2; line > ui.dim[H] - 3 { return line }
	// user infos
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "User: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.User)
	if line += 1; line > ui.dim[H] - 3 { return line }
	if len(host.Pass) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Pass: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], "***")
		if line += 1; line > ui.dim[H] - 3 { return line }
	}
	if len(host.Priv) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Privkey: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 12, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], host.Priv)
		if line += 1; line > ui.dim[H] - 3 { return line }
	}
	if line += 1; line > ui.dim[H] - 3 { return line }
	// jump
	if len(host.Jump.Host) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Jump settings: ")
		if line += 1; line > ui.dim[H] - 3 { return line }
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Host: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], host.Jump.Host)
		if line += 1; line > ui.dim[H] - 3 { return line }
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Port: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], strconv.Itoa(int(host.Jump.Port)))
		if line += 1; line > ui.dim[H] - 3 { return line }
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "User: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], host.Jump.User)
		if line += 1; line > ui.dim[H] - 3 { return line }
		if len(host.Jump.Pass) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
				ui.style[STYLE_TITLE], "Pass: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
				ui.style[STYLE_DEF], "***")
			if line += 1; line > ui.dim[H] - 3 { return line }
		}
		if len(host.Jump.Priv) > 0 {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
				ui.style[STYLE_TITLE], "Privkey: ")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 13, line, ui.dim[W] - 2, line,
				ui.style[STYLE_DEF], host.Jump.Priv)
			if line += 1; line > ui.dim[H] - 3 { return line }
		}
		if line += 1; line > ui.dim[H] - 3 { return line }
	}
	return line
}

func i_info_rdp(ui HardUI, host *HostNode, line int) int {
	if line > ui.dim[H] - 3 { return line }
	// host, port
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Host: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.Host)
	if line += 1; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Port: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], strconv.Itoa(int(host.Port)))
	if line += 1; line > ui.dim[H] - 3 { return line }
	// rdp shit
	if len(host.Domain) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Domain: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 11, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], host.Domain)
	}
	if line += 2; line > ui.dim[H] - 3 { return line }
	// user infos
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "User: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.User)
	if line += 1; line > ui.dim[H] - 3 { return line }
	if len(host.Pass) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Pass: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], "***")
		if line += 1; line > ui.dim[H] - 3 { return line }
	}
	if line += 1; line > ui.dim[H] - 3 { return line }
	// rdp shit
	qual := [3]string{"Low", "Medium", "High"}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Screen size: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 16, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF],
		strconv.Itoa(int(host.Width)) + "x" +
		strconv.Itoa(int(host.Height)))
	if line += 1; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Dynamic window: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 19, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], strconv.FormatBool(host.Dynamic))
	if line += 1; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Quality: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 12, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], qual[host.Quality])
	if line += 2; line > ui.dim[H] - 3 { return line }
	if host.Drive != nil {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Drives: ")
		if line += 1; line > ui.dim[H] - 3 { return line }
		for share, path := range host.Drive {
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 4, line, ui.dim[W] - 2, line,
				ui.style[STYLE_TITLE], share + ":")
			i_draw_text(ui.s,
				(ui.dim[W] / 3) + 4 + len(share) + 2, line,
				ui.dim[W] - 2, line,
				ui.style[STYLE_DEF], path)
			if line += 1; line > ui.dim[H] - 3 { return line }
		}
		if line += 1; line > ui.dim[H] - 3 { return line }
	}
	return line
}

func i_info_cmd(ui HardUI, host *HostNode, line int) int {
	if line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Command: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 12, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.Host)
	if line += 2; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Silent: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 11, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], strconv.FormatBool(host.Silent))
	if line += 1; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Shell: ")
	str := ""
	for _, s := range host.Shell {
		str += s + " "
	}
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 10, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], str)
	if line += 2; line > ui.dim[H] - 3 { return line }
	return line
}

func i_info_openstack(ui HardUI, host *HostNode, line int) int {
	if line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Endpoint: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 13, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.Host)
	if line += 1; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Region name: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 16, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.Stack.RegionName)
	if line += 2; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "User domain ID: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 19, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.Stack.UserDomainID)
	if line += 1; line > ui.dim[H] - 3 { return line }
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "Project ID: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 15, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.Stack.ProjectID)
	if line += 2; line > ui.dim[H] - 3 { return line }
	// user infos
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
		ui.style[STYLE_TITLE], "User: ")
	i_draw_text(ui.s,
		(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
		ui.style[STYLE_DEF], host.User)
	if line += 1; line > ui.dim[H] - 3 { return line }
	if len(host.Pass) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Pass: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], "***")
		if line += 1; line > ui.dim[H] - 3 { return line }
	}
	if line += 1; line > ui.dim[H] - 3 { return line }
	return line
}

func i_info_note(ui HardUI, host *HostNode, line int) {
	if line > ui.dim[H] - 3 {
		return
	}
	// note
	if len(host.Note) > 0 {
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 3, line, ui.dim[W] - 2, line,
			ui.style[STYLE_TITLE], "Note: ")
		i_draw_text(ui.s,
			(ui.dim[W] / 3) + 9, line, ui.dim[W] - 2, line,
			ui.style[STYLE_DEF], host.Note)
	}
}

func i_draw_info_panel(ui HardUI, percent bool, litems *ItemsList) {
	type info_func func(HardUI, *HostNode, int) int

	i_draw_box(ui.s, (ui.dim[W] / 3), 0,
		ui.dim[W] - 1, ui.dim[H] - 2,
		ui.style[STYLE_BOX], ui.style[STYLE_HEAD], " Infos ", false)
	ui.s.SetContent(ui.dim[W] / 3, 0, tcell.RuneTTee, nil, ui.style[STYLE_BOX])
	ui.s.SetContent(ui.dim[W] / 3, ui.dim[H] - 2,
		tcell.RuneBTee, nil, ui.style[STYLE_BOX])
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
			ui.style[STYLE_DEF],
			text)
	} else {
		text := " 0 hosts "
		i_draw_text(ui.s,
			(ui.dim[W] - 1) - len(text) - 1,
			ui.dim[H] - 2,
			(ui.dim[W] - 1) - 1,
			ui.dim[H] - 2,
			ui.style[STYLE_DEF],
			text)
	}
	// panel
	if litems.head == nil {
		return
	} else if litems.curr.is_dir() == true {
		i_info_dirs(ui, litems.curr.Dirs)
	} else {
		line := i_info_name_type(ui, litems.curr.Host)
		if litems.curr.Host.Protocol > PROTOCOL_MAX {
			return
		}
		fp := [PROTOCOL_MAX + 1]info_func{
			i_info_ssh,
			i_info_rdp,
			i_info_cmd,
			i_info_openstack,
		}
		line = fp[litems.curr.Host.Protocol](ui, litems.curr.Host, line)
		i_info_note(ui, litems.curr.Host, line)
	}
}
