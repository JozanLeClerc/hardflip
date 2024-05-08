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
 * hardflip: src/i_insert.go
 * Tue May 07 10:23:23 2024
 * Joe
 *
 * insert a new host
 */

package main

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/uuid"
)

func i_insert_format_filename(name, path string) string {
	str := name

	if len(name) == 0 {
		return ""
	}
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	replace := "_"
	str = re.ReplaceAllString(str, replace)
	re2 := regexp.MustCompile("__")
	replace = "_"
	str = re2.ReplaceAllString(str, replace)
	re3 := regexp.MustCompile("_$")
	replace = ""
	str = re3.ReplaceAllString(str, replace)
	_, err := os.Stat(path + str + ".yml")
	base := str
	i := 0
	for err == nil && i < 10000 {
		uid := uuid.NewUUID()
		str = base + "_" + string(uid[0:4])
		_, err = os.Stat(path + str + ".yml")
		i++
	}
	str = strings.ToLower(str) + ".yml"
	return str
}

func i_insert_abs_files(insert *HostNode, home_dir string) {
	files := []*string{
		&insert.Priv,
		&insert.Jump.Priv,
		&insert.RDPFile,
	}

	for _, v := range files {
		if len(*v) > 0 {
			if (*v)[0] == '~' {
				*v = home_dir + (*v)[1:]
			}
			*v, _ = filepath.Abs(*v)
		}
	}
	for k, v := range insert.Drive {
		if len(v) > 0 {
			if (v)[0] == '~' {
				v = home_dir + (v)[1:]
			}
			v, _ = filepath.Abs(v)
			insert.Drive[k] = v
		}
	}
}

func i_insert_default_users(insert *HostNode) {
	switch insert.Protocol {
	case PROTOCOL_SSH:
		if len(insert.User) == 0 {
			insert.User = "root"
		}
	case PROTOCOL_RDP:
		if len(insert.User) == 0 {
			insert.User = "Administrator"
		}
	default: return
	}
}

func i_insert_host(data *HardData, insert *HostNode) error {
	i_insert_abs_files(insert, data.home_dir)
	i_insert_default_users(insert)
	if len(insert.Drive) == 0 {
		insert.Drive = nil
	}
	filename := insert.filename
	if data.ui.insert_method == INSERT_ADD ||
	   data.ui.insert_method == INSERT_COPY ||
	   data.ui.insert_method == INSERT_MOVE {
		filename = i_insert_format_filename(insert.Name,
			data.data_dir + insert.parent.path())
		insert.filename = filename
	}
	format, err := yaml.Marshal(insert)
	if err != nil {
		c_error_mode("yaml", err, &data.ui)
		data.insert = nil
		return err
	}
	err = os.WriteFile(data.data_dir + insert.parent.path() + filename,
		format, 0644)
	if err != nil {
		c_error_mode("can't write file", err, &data.ui)
		data.insert = nil
		return err
	}
	if data.ui.insert_method == INSERT_EDIT && data.litems.curr != nil {
		tmp := e_deep_copy_host(data.insert)
		data.litems.curr.Host = &tmp
		data.litems.reset_id()
		data.ui.mode = NORMAL_MODE
		data.insert = nil
		return nil
	}
	// HACK: not sure if this is necessary
	// if data.litems.curr.is_dir() == true {
	// 	data.litems.curr.Dirs.lhost.add_back(insert)
	// } else {
	// 	tmp_next := data.litems.curr.Host.next
	// 	data.litems.curr.Host.next = insert
	// 	data.litems.curr.Host.next.next = tmp_next
	// }
	var next *ItemsNode = nil
	if data.litems.curr != nil {
		next = data.litems.curr.next
	}
	item := &ItemsNode{
		0,
		nil,
		insert,
		data.litems.curr,
		next,
	}
	curr := data.litems.curr
	if curr != nil {
		curr.next = item
		if curr.next.next != nil {
			data.litems.curr.next.next.prev = item
		}
		data.litems.curr = data.litems.curr.next
	} else {
		data.litems.add_back(item)
		data.litems.curr = data.litems.head
	}
	data.litems.reset_id()
	data.ui.mode = NORMAL_MODE
	data.insert = nil
	return nil
}

func i_insert_check_ok(data *HardData, in *HostNode) {
	if len(in.Name) == 0 {
		data.insert_err = append(data.insert_err, errors.New("no name"))
	}
	if len(in.Host) == 0 {
		if (in.Protocol == PROTOCOL_RDP && len(in.RDPFile) > 0) == false {
			text := "no host"
			if in.Protocol == PROTOCOL_CMD {
				text = "no command"
			} else if in.Protocol == PROTOCOL_OS {
				text = "no endpoint"
			}
			data.insert_err = append(data.insert_err, errors.New(text))
		}
	}
	if (in.Protocol == PROTOCOL_SSH || in.Protocol == PROTOCOL_RDP) &&
	   in.Port == 0 {
		data.insert_err = append(data.insert_err, errors.New("port can't be 0"))
	}
	if len(in.Jump.Host) > 0 && in.Jump.Port == 0 {
		data.insert_err = append(data.insert_err,
			errors.New("jump port can't be 0"))
	}
	if in.Protocol == PROTOCOL_OS {
		if len(in.User) == 0 {
			data.insert_err = append(data.insert_err,
				errors.New("user can't be empty"))
		}
		if len(in.Stack.UserDomainID) == 0 {
			data.insert_err = append(data.insert_err,
				errors.New("user domain ID can't be empty"))
		}
		if len(in.Stack.ProjectID) == 0 {
			data.insert_err = append(data.insert_err,
				errors.New("project ID can't be empty"))
		}
	}
	var file [2]string
	switch in.Protocol {
	case PROTOCOL_SSH: file[0], file[1] = in.Priv, in.Jump.Priv
	case PROTOCOL_RDP: file[0], file[1] = in.RDPFile, in.Jump.Priv
	case PROTOCOL_CMD: file[0] = in.Shell[0]
	default: return
	}
	for _, v := range file {
		if len(v) > 0 {
			if v[0] == '~' {
				v = data.home_dir + v[1:]
			}
			if stat, err := os.Stat(v);
			   err != nil {
				data.insert_err = append(data.insert_err, errors.New(v +
					": file does not exist"))
			} else if stat.IsDir() == true {
				data.insert_err = append(data.insert_err, errors.New(v +
					": file is a directory"))
			}
		}
	}
	for _, v := range in.Drive {
		if v[0] == '~' {
			v = data.home_dir + v[1:]
		}
		if stat, err := os.Stat(v);
		   err != nil {
			data.insert_err = append(data.insert_err, errors.New(v +
				": path does not exist"))
		} else if stat.IsDir() == false {
			data.insert_err = append(data.insert_err, errors.New(v +
				": path is not a directory"))
		}
	}
}

func i_draw_tick_box(ui HardUI, line int, dim Quad, label string, content bool,
					 id, selected int) {
	tbox_style := ui.style[DEF_STYLE].Background(tcell.ColorBlack).Dim(true)

	if id == selected {
		tbox_style = tbox_style.Reverse(true).Dim(false)
	}
	l := ui.dim[W] / 2 - len(label) - 2
	if l <= dim.L { l = dim.L + 1 }
	i_draw_text(ui.s, l, line, ui.dim[W] / 2, line,
		ui.style[DEF_STYLE], label)
	x := " "
	if content == true {
		x = "x"
	}
	i_draw_text(ui.s, ui.dim[W] / 2, line, dim.R, line,
		tbox_style,
		"[" + x + "]")
}

func i_draw_text_box(ui HardUI, line int, dim Quad, label, content string,
					 id int, secret, red bool) {
	selected := ui.insert_sel
	const tbox_size int = 14
	tbox_style := ui.style[DEF_STYLE].Background(tcell.ColorBlack).Dim(true)

	if id == selected {
		tbox_style = tbox_style.Reverse(true).Dim(false)
	}

	l := ui.dim[W] / 2 - len(label) - 2
	if l <= dim.L { l = dim.L + 1 }
	i_draw_text(ui.s, l, line, ui.dim[W] / 2, line,
		ui.style[DEF_STYLE], label)
	if secret == true &&
		len(content) > 0 {
		content = "***"
	}
	if red == true {
		tbox_style = tbox_style.Foreground(tcell.ColorRed)
	}
	spaces := ""
	for i := 0; i < tbox_size; i++ {
		spaces += " "
	}
	i_draw_text(ui.s, ui.dim[W] / 2, line, dim.R, line,
		tbox_style,
		"[" + spaces + "]")
	i_draw_text(ui.s, ui.dim[W] / 2 + 1, line, ui.dim[W] / 2 + 1 + tbox_size,
		line, tbox_style, content)
}

func i_draw_ok_butt(ui HardUI, line int, id, selected int) {
	const butt_size int = 10
	const txt string = "ok"
	style := ui.style[DEF_STYLE].Background(tcell.ColorBlack).Dim(true)

	if id == selected {
		style = style.Reverse(true).Dim(false)
	}
	buff := "["
	for i := 0; i < butt_size / 2 - len(txt); i++ {
		buff += " "
	}
	buff += txt
	for i := 0; i < butt_size / 2 - len(txt); i++ {
		buff += " "
	}
	buff += "]"
	i_draw_text(ui.s, (ui.dim[W] / 2) - (butt_size / 2), line,
		(ui.dim[W] / 2) + (butt_size / 2), line, style, buff)
}

func i_draw_insert_inputs(ui HardUI, in *HostNode, home_dir string) {
	if ui.insert_sel_ok == false {
		return
	}
	switch ui.insert_sel {
	case INS_PROTOCOL:
		i_prompt_list(ui, "Connection type", "Type:",
					  PROTOCOL_STR[:])
	case INS_SSH_HOST,
		 INS_SSH_JUMP_HOST,
		 INS_RDP_JUMP_HOST + len(in.Drive),
		 INS_RDP_HOST:
		i_prompt_generic(ui, "Host/IP: ", false, "")
	case INS_SSH_PORT,
		 INS_SSH_JUMP_PORT,
		 INS_RDP_JUMP_PORT + len(in.Drive),
		 INS_RDP_PORT:
		i_prompt_generic(ui, "Port: ", false, "")
	case INS_SSH_USER,
		 INS_SSH_JUMP_USER,
		 INS_RDP_JUMP_USER + len(in.Drive),
		 INS_RDP_USER,
		 INS_OS_USER:
		i_prompt_generic(ui, "User: ", false, "")
	case INS_SSH_PASS,
		 INS_SSH_JUMP_PASS,
		 INS_RDP_JUMP_PASS + len(in.Drive),
		 INS_RDP_PASS,
		 INS_OS_PASS:
		i_prompt_generic(ui, "Pass: ", true, "")
	case INS_SSH_PRIV,
		 INS_SSH_JUMP_PRIV,
		 INS_RDP_JUMP_PRIV + len(in.Drive):
		i_prompt_generic(ui, "Private key: ", false, home_dir)
	case INS_SSH_EXEC:
		i_prompt_generic(ui, "Command: ", false, "")
	case INS_SSH_NOTE,
		 INS_RDP_NOTE + len(in.Drive),
		 INS_CMD_NOTE,
		 INS_OS_NOTE:
		i_prompt_generic(ui, "Note: ", false, "")
	case INS_RDP_DOMAIN:
		i_prompt_generic(ui, "Domain: ", false, "")
	case INS_RDP_FILE:
		i_prompt_generic(ui, "RDP file: ", false, home_dir)
	case INS_RDP_SCREENSIZE:
		i_prompt_list(ui, "Window size", "Size:",
					  RDP_SCREENSIZE[:])
	case INS_RDP_QUALITY:
		i_prompt_list(ui, "Quality", "Quality:",
					  RDP_QUALITY[:])
	case INS_RDP_DRIVE + len(in.Drive):
		if len(ui.drives_buff) == 0 {
			i_prompt_generic(ui, "Name: ", false, "")
		} else {
			i_prompt_dir(ui, "Local directory: ", home_dir)
		}
	case INS_CMD_CMD:
		i_prompt_generic(ui, "Command: ", false, "")
	case INS_CMD_SHELL:
		i_prompt_generic(ui, "Shell: ", false, home_dir)
	case INS_OS_HOST:
		i_prompt_generic(ui, "Endpoint: ", false, "")
	case INS_OS_USERDOMAINID:
		i_prompt_generic(ui, "User Domain ID: ", false, "")
	case INS_OS_PROJECTID:
		i_prompt_generic(ui, "Project ID: ", false, "")
	case INS_OS_REGION:
		i_prompt_generic(ui, "Region name: ", false, "")
	case INS_OS_ENDTYPE:
		i_prompt_generic(ui, "Endpoint type: ", false, "")
	case INS_OS_INTERFACE:
		i_prompt_generic(ui, "Interface: ", false, "")
	case INS_OS_IDAPI:
		i_prompt_generic(ui, "Identity API version: ", false, "")
	case INS_OS_IMGAPI:
		i_prompt_generic(ui, "Image API version: ", false, "")
	case INS_OS_NETAPI:
		i_prompt_generic(ui, "Network API version: ", false, "")
	case INS_OS_VOLAPI:
		i_prompt_generic(ui, "Volume API version: ", false, "")
	}
	if len(in.Drive) > 0 &&
	   ui.insert_sel >= INS_RDP_DRIVE &&
	   ui.insert_sel < INS_RDP_DRIVE +
	   len(in.Drive) {
		i_draw_remove_share(ui)
	}
}

func i_insert_follow_cursor(ui *HardUI, line int) int {
	return line - 15
}

func i_draw_insert_panel(ui HardUI, in *HostNode, home_dir string) {
	type draw_insert_func func(ui HardUI, line int, win Quad,
							   in *HostNode, home string) int

	if len(in.Name) == 0 {
		return
	}
	win := Quad{
		ui.dim[W] / 8,
		ui.dim[H] / 8,
		ui.dim[W] - ui.dim[W] / 8 - 1,
		ui.dim[H] - ui.dim[H] / 8 - 1,
	}
	i_draw_box(ui.s, win.L, win.T, win.R, win.B,
		ui.style[BOX_STYLE], ui.style[HEAD_STYLE],
		" Insert - " + in.Name + " ", true)
	line := 2
	if win.T + line >= win.B { return }
	i_draw_text_box(ui, win.T + line, win, "Connection type",
		PROTOCOL_STR[in.Protocol], 0, false, false)
	line += 2
	var end_line int
	fp := [PROTOCOL_MAX + 1]draw_insert_func{
		i_draw_insert_ssh,
		i_draw_insert_rdp,
		i_draw_insert_cmd,
		i_draw_insert_os,
	}
	line = i_insert_follow_cursor(&ui, line)
	end_line = fp[in.Protocol](ui, line, win, in, home_dir)
	if win.T + end_line >= win.B {
		ui.s.SetContent(ui.dim[W] / 2, win.B, '▼', nil, ui.style[BOX_STYLE])
		// TODO: scroll or something
	}
	i_draw_insert_inputs(ui, in, home_dir)
}

func i_draw_insert_ssh(ui HardUI, line int, win Quad,
					   in *HostNode, home string) int {
	red := false
	if win.T + line >= win.B { return line }
	text := "---- Host settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Host/IP", in.Host,
		INS_SSH_HOST, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Port", strconv.Itoa(int(in.Port)),
		INS_SSH_PORT, false, false);
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "User", in.User,
		INS_SSH_USER, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Pass", in.Pass,
		INS_SSH_PASS, true, false)
	if line += 1; win.T + line >= win.B { return line }
	if file := in.Priv; len(file) > 0 {
		if file[0] == '~' {
			file = home + file[1:]
		}
		if stat, err := os.Stat(file);
		   err != nil || stat.IsDir() == true {
			red = true
		}
	}
	i_draw_text_box(ui, win.T + line, win, "SSH private key", in.Priv,
		INS_SSH_PRIV, false, red)
	if red == true {
		if line += 1; win.T + line >= win.B { return line }
		text := "file does not exist"
		i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
			win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
	}
	red = false
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Command (optional)", in.Exec,
		INS_SSH_EXEC, false, false)
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Jump settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Host/IP", in.Jump.Host,
		INS_SSH_JUMP_HOST, false, false)
	if len(in.Jump.Host) > 0 {
		if line += 1; win.T + line >= win.B { return line }
		i_draw_text_box(ui, win.T + line, win, "Port",
			strconv.Itoa(int(in.Jump.Port)),
			INS_SSH_JUMP_PORT, false, false)
		if line += 2; win.T + line >= win.B { return line }
		i_draw_text_box(ui, win.T + line, win, "User", in.Jump.User,
			INS_SSH_JUMP_USER, false, false)
		if line += 1; win.T + line >= win.B { return line }
		i_draw_text_box(ui, win.T + line, win, "Pass", in.Jump.Pass,
			INS_SSH_JUMP_PASS, true, false)
		if line += 1; win.T + line >= win.B { return line}
		if len(in.Jump.Priv) > 0 {
			file := in.Jump.Priv
			if file[0] == '~' {
				file = home + file[1:]
			}
			if stat, err := os.Stat(file);
			   err != nil || stat.IsDir() == true {
				red = true
			}
		}
		i_draw_text_box(ui, win.T + line, win, "SSH private key", in.Jump.Priv,
			INS_SSH_JUMP_PRIV, false, red)
		if red == true {
			if line += 1; win.T + line >= win.B { return line }
			text := "file does not exist"
			i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
				win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
		}
	}
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Note ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Note", in.Note,
		INS_SSH_NOTE, false, false)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_ok_butt(ui, win.T + line, INS_SSH_OK, ui.insert_sel)
	return line
}

func i_draw_insert_rdp(ui HardUI, line int, win Quad,
					   in *HostNode, home string) int {
	red := false
	if win.T + line >= win.B { return line }
	text := "---- Host settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Host/IP", in.Host,
		INS_RDP_HOST, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Port", strconv.Itoa(int(in.Port)),
		INS_RDP_PORT, false, false);
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Domain", in.Domain,
		INS_RDP_DOMAIN, false, false);
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "User", in.User,
		INS_RDP_USER, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Pass", in.Pass,
		INS_RDP_PASS, true, false)
	if line += 2; win.T + line >= win.B { return line }
	text = "---- RDP File ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	if file := in.RDPFile; len(file) > 0 {
		if file[0] == '~' {
			file = home + file[1:]
		}
		if stat, err := os.Stat(file);
		   err != nil || stat.IsDir() == true {
			red = true
		}
	}
	i_draw_text_box(ui, win.T + line, win, "RDP file", in.RDPFile,
		INS_RDP_FILE, false, red)
	if red == true {
		if line += 1; win.T + line >= win.B { return line }
		text := "file does not exist"
		i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
			win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
	}
	red = false
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Window settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	screensize := strconv.Itoa(int(in.Width)) + "x" +
				  strconv.Itoa(int(in.Height))
	i_draw_text_box(ui, win.T + line, win, "Window size", screensize,
		INS_RDP_SCREENSIZE, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_tick_box(ui, win.T + line, win, "Dynamic window", in.Dynamic,
		INS_RDP_DYNAMIC, ui.insert_sel)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_tick_box(ui, win.T + line, win, "Full screen", in.FullScr,
		INS_RDP_FULLSCR, ui.insert_sel)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_tick_box(ui, win.T + line, win, "Multi monitor", in.MultiMon,
		INS_RDP_MULTIMON, ui.insert_sel)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Quality", RDP_QUALITY[in.Quality],
		INS_RDP_QUALITY, false, false)
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Share mounts ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	for k, v := range in.drive_keys {
		red = false
		if dir := in.Drive[v]; len(dir) > 0 {
			if dir[0] == '~' {
				dir = home + dir[1:]
			}
			if stat, err := os.Stat(dir);
			   err != nil || stat.IsDir() == false {
				red = true
			}
		}
		i_draw_text_box(ui, win.T + line, win, "Share " + strconv.Itoa(k + 1),
			"(" + v + "): " + in.Drive[v],
			INS_RDP_DRIVE + k, false, red)
		if red == true {
			if line += 1; win.T + line >= win.B { return line }
			text := "path is not a directory"
			i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
				win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
		}
		if line += 1; win.T + line >= win.B { return line }
	}
	i_draw_text_box(ui, win.T + line, win, "Add share", "",
		INS_RDP_DRIVE + len(in.Drive), false, false)
	red = false
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Jump settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Host/IP", in.Jump.Host,
		INS_RDP_JUMP_HOST + len(in.Drive), false, false)
	if len(in.Jump.Host) > 0 {
		if line += 1; win.T + line >= win.B { return line }
		i_draw_text_box(ui, win.T + line, win, "Port",
			strconv.Itoa(int(in.Jump.Port)),
			INS_RDP_JUMP_PORT + len(in.Drive), false, false)
		if line += 2; win.T + line >= win.B { return line }
		i_draw_text_box(ui, win.T + line, win, "User", in.Jump.User,
			INS_RDP_JUMP_USER + len(in.Drive), false, false)
		if line += 1; win.T + line >= win.B { return line }
		i_draw_text_box(ui, win.T + line, win, "Pass", in.Jump.Pass,
			INS_RDP_JUMP_PASS + len(in.Drive), true, false)
		if line += 1; win.T + line >= win.B { return line}
		if len(in.Jump.Priv) > 0 {
			file := in.Jump.Priv
			if file[0] == '~' {
				file = home + file[1:]
			}
			if stat, err := os.Stat(file);
			   err != nil || stat.IsDir() == true {
				red = true
			}
		}
		i_draw_text_box(ui, win.T + line, win, "SSH private key", in.Jump.Priv,
			INS_RDP_JUMP_PRIV + len(in.Drive), false, red)
		if red == true {
			if line += 1; win.T + line >= win.B { return line }
			text := "file does not exist"
			i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
				win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
		}
	}
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Note ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Note", in.Note,
		INS_RDP_NOTE + len(in.Drive), false, false)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_ok_butt(ui, win.T + line, INS_RDP_OK + len(in.Drive), ui.insert_sel)
	return line
}

func i_draw_insert_cmd(ui HardUI, line int, win Quad,
					   in *HostNode, home string) int {
	red := false
	if win.T + line >= win.B { return line }
	text := "---- Settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Command", in.Host,
		INS_CMD_CMD, false, false)
	if line += 1; win.T + line >= win.B { return line }
	if shell := in.Shell[0]; len(shell) > 0 {
		if shell[0] == '~' {
			shell = home + shell[1:]
		}
		if stat, err := os.Stat(shell);
		   err != nil || stat.IsDir() == true {
			red = true
		}
	}
	i_draw_text_box(ui, win.T + line, win, "Shell", in.Shell[0],
		INS_CMD_SHELL, false, red);
	if red == true {
		if line += 1; win.T + line >= win.B { return line }
		text := "file does not exist"
		i_draw_text(ui.s, ui.dim[W] / 2, win.T + line,
			win.R - 1, win.T + line, ui.style[ERR_STYLE], text)
	}
	red = false
	if line += 1; win.T + line >= win.B { return line }
	i_draw_tick_box(ui, win.T + line, win, "Silent", in.Silent,
		INS_CMD_SILENT, ui.insert_sel)
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Note ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Note", in.Note,
		INS_CMD_NOTE, false, false)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_ok_butt(ui, win.T + line, INS_CMD_OK, ui.insert_sel)
	return line
}

func i_draw_insert_os(ui HardUI, line int, win Quad,
					  in *HostNode, home string) int {
	if win.T + line >= win.B { return line }
	text := "---- Host settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Endpoint", in.Host,
		INS_OS_HOST, false, false)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "User", in.User,
		INS_OS_USER, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Pass", in.Pass,
		INS_OS_PASS, true, false)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "User domain ID",
		in.Stack.UserDomainID,
		INS_OS_USERDOMAINID, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Project ID", in.Stack.ProjectID,
		INS_OS_PROJECTID, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Region name", in.Stack.RegionName,
		INS_OS_REGION, false, false)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Endpoint type",
		in.Stack.EndpointType,
		INS_OS_ENDTYPE, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Interface", in.Stack.Interface,
		INS_OS_INTERFACE, false, false)
	if line += 2; win.T + line >= win.B { return line }
	text = "---- API settings ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Identity API version",
		in.Stack.IdentityAPI,
		INS_OS_IDAPI, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Image API version",
		in.Stack.ImageAPI,
		INS_OS_IMGAPI, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Network API version",
		in.Stack.NetworkAPI,
		INS_OS_NETAPI, false, false)
	if line += 1; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Volume API version",
		in.Stack.VolumeAPI,
		INS_OS_VOLAPI, false, false)
	if line += 2; win.T + line >= win.B { return line }
	text = "---- Note ----"
	i_draw_text(ui.s, ui.dim[W] / 2 - len(text) / 2, win.T + line, win.R - 1,
		win.T + line, ui.style[DEF_STYLE], text)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_text_box(ui, win.T + line, win, "Note", in.Note,
		INS_OS_NOTE, false, false)
	if line += 2; win.T + line >= win.B { return line }
	i_draw_ok_butt(ui, win.T + line, INS_OS_OK, ui.insert_sel)
	return line
}
