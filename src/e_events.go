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
 * hardflip: src/e_events.go
 * Thu Apr 11 16:00:44 2024
 * Joe
 *
 * events in the code
 */

package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

func e_list_follow_cursor(litems *ItemsList, ui *HardUI) {
	if litems.draw == nil || litems.curr == nil {
		return
	}
	scrolloff := 4
	if litems.last.ID - (ui.dim[H] - 4) <= litems.draw.ID {
		scrolloff = 0
	}
	virt_id := litems.curr.ID - (ui.dim[H] - 4) + scrolloff
	for litems.draw.ID < virt_id &&
		litems.draw.next != nil {
		litems.draw = litems.draw.next
	}
	scrolloff = 4
	for litems.draw.ID > litems.curr.ID - scrolloff &&
		litems.draw.prev != nil {
		litems.draw = litems.draw.prev
	}
}

func e_set_unfold(data *HardData, item *ItemsNode) {
	delete(data.folds, item.Dirs)
	data.litems.reset_id()
}

func e_unfold_dir(data *HardData, item *ItemsNode) {
	if item == nil || item.Dirs == nil {
		return
	}
	fold := data.folds[item.Dirs]
	if fold == nil {
		return
	}
	start, end := fold.head, fold.last
	// last empty dir
	if start == nil && end == nil {
		e_set_unfold(data, item)
		return
	}
	// single empty dir
	if start == item && end == end { // HACK: i forgot why end == end
		e_set_unfold(data, item)
		return
	}
	if data.litems.last == item {
		data.litems.last = end
	}
	// non-emtpy dir
	start.prev = item
	end.next = item.next
	if item.next != nil {
		item.next.prev = end
	}
	item.next = start
	e_set_unfold(data, item)
}

func e_set_fold(data *HardData, curr, start, end *ItemsNode) {
	folds := data.folds
	tmp := ItemsList{
		start,
		end,
		nil,
		nil,
	}

	folds[curr.Dirs] = &tmp
	data.litems.reset_id()
}

func e_fold_dir(data *HardData, item *ItemsNode) {
	if item == nil || item.Dirs == nil {
		return
	}
	var start, end *ItemsNode
	start = item.next
	// last dir + empty
	if start == nil {
		e_set_fold(data, item, nil, nil)
		return
	}
	// empty dir
	if start.Dirs != nil && start.Dirs.Depth <= item.Dirs.Depth {
		e_set_fold(data, item, item, item)
		return
	}
	// non-empty dir
	start.prev = nil
	end = start
	next_dir := item.get_next_level()
	// this is the end
	if next_dir == nil {
		item.next = nil
		end = data.litems.last
		end.next = nil
		data.litems.last = item
		e_set_fold(data, item, start, end)
		return
	}
	// this is not the end
	end = next_dir.prev
	end.next = nil
	item.next = next_dir
	next_dir.prev = item
	e_set_fold(data, item, start, end)
}

func e_reload_data(data *HardData) {
	tmp_name := ""
	tmp_parent_path := ""
	if data.litems.curr != nil {
		if data.litems.curr.is_dir() == true {
			tmp_name = data.litems.curr.Dirs.Name
			tmp_parent_path = data.litems.curr.Dirs.Parent.path()
		} else {
			tmp_name = data.litems.curr.Host.filename
			tmp_parent_path = data.litems.curr.Host.parent.path()
		}
	}
	conf_dir  := c_get_conf_dir(&data.load_err)
	if conf_dir == "" {
		data.opts = DEFAULT_OPTS
	} else {
		data.opts = c_get_options(conf_dir, &data.load_err)
	}
	data.data_dir = c_get_data_dir(&data.ui)
	if len(data.data_dir) == 0 {
		return
	}
	g_load_count = -1
	data.ldirs, data.litems, data.load_err = i_load_ui(data.data_dir, data.opts,
		&data.ui, &data.load_err)
	data.folds = make(map[*DirsNode]*ItemsList)
	if tmp_name == "" {
		data.litems.curr = data.litems.head
		return
	}
	for curr := data.litems.head; curr != nil; curr = curr.next {
		if curr.is_dir() == true {
			if curr.Dirs.Name == tmp_name {
				if curr.Dirs.Parent.path() == tmp_parent_path {
					data.litems.curr = curr
					return
				}
			}
		} else {
			if curr.Host.filename == tmp_name {
				if curr.Host.parent.path() == tmp_parent_path {
					data.litems.curr = curr
					return
				}
			}
		}
	}
	data.litems.curr = data.litems.head
}

func e_delete_dir(data *HardData) error {
	if data.litems.curr == nil || data.litems.curr.Dirs == nil {
		return nil
	}
	curr := data.litems.curr
	dir_path := data.data_dir + data.litems.curr.Dirs.path()
	if err := os.RemoveAll(dir_path); err != nil {
		c_error_mode("can't remove " + dir_path, err, &data.ui)
		return err
	}
	if data.folds[curr.Dirs] == nil {
		e_fold_dir(data, curr)
	}
	delete(data.folds, curr.Dirs)
	if curr == data.litems.head {
		data.litems.head = curr.next
		if curr.next != nil {
			curr.next.prev = nil
		}
		if data.litems.draw == curr {
			data.litems.draw = curr.next
		}
	} else {
		curr.prev.next = curr.next
	}
	if curr.next != nil {
		curr.next.prev = curr.prev
		data.litems.curr = curr.next
	} else {
		data.litems.last = curr.prev
		data.litems.curr = curr.prev
	}
	data.litems.reset_id()
	return nil
}

func e_delete_host(data *HardData) error {
	if data.litems.curr == nil {
		return nil
	}
	if data.litems.curr.is_dir() == true {
		return e_delete_dir(data)
	}
	host := data.litems.curr.Host
	if host == nil {
		return nil
	}
	file_path := data.data_dir + host.parent.path() + host.filename

	if err := os.Remove(file_path); err != nil {
		c_error_mode("can't remove " + file_path, err, &data.ui)
		return err
	}
	tmp := data.litems.curr.prev
	host.parent.lhost.del(host)
	data.litems.del(data.litems.curr)
	if tmp == nil {
		tmp = data.litems.head
	}
	data.litems.curr = tmp
	return nil
}

func e_tab_complete(buffer *Buffer) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return
	}
	var match []string
	for _, v := range entries {
		if len(v.Name()) >= buffer.len() &&
		   v.Name()[:buffer.len()] == buffer.str() {
			match = append(match, v.Name())
		}
	}
	log.Println(match)
	if len(match) == 0 {
		return
	} else if len(match) == 1 {
		buffer.insert(match[0])
	} else {
		var common []rune
		var shortest int = 1000000000
		var check bool = true
		for _, v := range match {
			if len(v) < shortest {
				shortest = len(v)
			}
		}
		for j := 0; j < shortest; j++ {
			for i := 1; i < len(match); i++ {
				if match[i][j] != match[i - 1][j] {
					check = false
				}
			}
			if check == false {
				break
			} else {
				common = append(common, rune(match[0][j]))
			}
		}
		if len(common) == 0 {
			return
		}
		buffer.insert(string(common))
	}
}

func e_readline(event tcell.EventKey, buffer *Buffer) {
	if buffer.len() > 0 &&
	   (event.Key() == tcell.KeyBackspace ||
	   event.Key() == tcell.KeyBackspace2) {
		if buffer.cursor == 0 {
			return
		} else if buffer.cursor == buffer.len() {
			buffer.data = (buffer.data)[:buffer.cursor - 1]
		} else {
			buffer.data = append(buffer.data[:buffer.cursor - 1],
				buffer.data[buffer.cursor:]...)
		}
		buffer.cursor -= 1
	} else if event.Key() == tcell.KeyDelete ||
			  event.Key() == tcell.KeyCtrlD {
		if buffer.cursor == buffer.len() {
			return
		} else {
			buffer.data = append(buffer.data[:buffer.cursor],
				buffer.data[buffer.cursor + 1:]...)
		}
	} else if event.Key() == tcell.KeyCtrlU {
		buffer.empty()
	} else if event.Rune() >= 32 && event.Rune() <= 126 {
		if buffer.cursor == buffer.len() {
			buffer.data = append(buffer.data, event.Rune())
		} else {
			buffer.data = append(buffer.data[:buffer.cursor],
				append([]rune{event.Rune()}, buffer.data[buffer.cursor:]...)...)
		}
		buffer.cursor += 1
	} else if event.Key() == tcell.KeyCtrlA {
		buffer.cursor = 0
	} else if event.Key() == tcell.KeyCtrlE {
		buffer.cursor = buffer.len()
	} else if event.Key() == tcell.KeyLeft ||
			  event.Key() == tcell.KeyCtrlB {
		buffer.cursor -= 1
	} else if event.Key() == tcell.KeyRight ||
			  event.Key() == tcell.KeyCtrlF {
		buffer.cursor += 1
	} else if event.Key() == tcell.KeyTab ||
			  event.Key() == tcell.KeyCtrlI {
		e_tab_complete(buffer)
	}
	if buffer.cursor > buffer.len() {
		buffer.cursor = buffer.len()
	} else if buffer.cursor < 0 {
		buffer.cursor = 0
	}
}

func e_mkdir(data *HardData, ui *HardUI) {
	if ui.buff.len() == 0 {
		return
	}
	path := "/"
	if data.litems.curr != nil {
		path = data.litems.curr.path()
	}
	if err := os.MkdirAll(data.data_dir +
		path +
		ui.buff.str(), os.ModePerm); err != nil {
		c_error_mode("mkdir " + path[1:] + ui.buff.str() + " failed",
		err, ui)
		return
	}
	e_reload_data(data)
	for curr := data.litems.head; curr != nil; curr = curr.next {
		if curr.is_dir() == true &&
		   curr.Dirs.Name == ui.buff.str() &&
		   curr.Dirs.Parent.path() == path {
			data.litems.curr = curr
			return
		}
	}
}

func e_rename(data *HardData, ui *HardUI) error {
	tmp := data.litems.curr
	name := ""
	if tmp.is_dir() == false {
		name = tmp.Host.Name
	} else {
		name = tmp.Dirs.Name
	}

	if ui.buff.len() == 0 || tmp == nil || ui.buff.str() == name {
		return nil
	}
	if tmp.is_dir() == false {
		new_host := e_deep_copy_host(data.litems.curr.Host)
		new_host.Name = ui.buff.str()
		ui.insert_method = INSERT_MOVE
		i_insert_host(data, &new_host)
		data.litems.del(tmp)
		file_path := data.data_dir + tmp.Host.parent.path() + tmp.Host.filename
		if err := os.Remove(file_path); err != nil {
			c_error_mode("can't remove " + file_path, err, &data.ui)
			return err
		}
		return nil
	} else {
		old_path := data.data_dir + tmp.Dirs.path()
		new_path := data.data_dir + tmp.Dirs.Parent.path() + data.ui.buff.str()
		if err := os.Rename(old_path, new_path); err != nil {
			c_error_mode("can't rename " + old_path, err, &data.ui)
			return err
		}
		path := tmp.Dirs.Parent.path()
		e_reload_data(data)
		for curr := data.litems.head; curr != nil; curr = curr.next {
			if curr.is_dir() == true &&
			   curr.Dirs.Name == ui.buff.str() &&
			   curr.Dirs.Parent.path() == path {
				data.litems.curr = curr
				return nil
			}
		}
		return nil
	}
}

func e_set_drive_keys(data *HardData) {
	data.insert.drive_keys = nil
	for key := range data.insert.Drive {
		data.insert.drive_keys = append(data.insert.drive_keys, key)
	}
	data.ui.insert_sel_max = INS_RDP_OK + len(data.insert.Drive)
}

func e_set_protocol_max(data *HardData, in *HostNode) {
	switch in.Protocol {
	case PROTOCOL_SSH:
		data.ui.insert_sel_max = INS_SSH_OK
	case PROTOCOL_RDP:
		data.ui.insert_sel_max = INS_RDP_OK + len(in.Drive)
	case PROTOCOL_CMD:
		data.ui.insert_sel_max = INS_CMD_OK
	case PROTOCOL_OS:
		data.ui.insert_sel_max = INS_OS_OK
	}
}

func e_set_protocol_defaults(data *HardData, in *HostNode) {
	switch in.Protocol {
	case PROTOCOL_SSH:
		in.Port = 22
	case PROTOCOL_RDP:
		in.Port = 3389
		in.Quality = 2
		in.Width = 1600
		in.Height = 1200
		in.Dynamic = true
		in.drive_keys = nil
	case PROTOCOL_CMD:
		in.Silent = false
		in.Shell = []string{"/bin/sh", "-c"}
	case PROTOCOL_OS:
		in.Stack.RegionName = "eu-west-0"
		in.Stack.IdentityAPI = "3"
		in.Stack.ImageAPI    = "2"
		in.Stack.NetworkAPI  = "2"
		in.Stack.VolumeAPI   = "3.42"
		in.Stack.EndpointType = "publicURL"
		in.Stack.Interface = "public"
	}
	data.ui.insert_scroll = 0
	e_set_protocol_max(data, in)
}

func e_deep_copy_host(base *HostNode) HostNode {
	new_host := HostNode{}
	new_host = *base
	if base.Drive != nil {
		new_host.Drive = make(map[string]string, len(base.Drive))
		for k, v := range base.Drive {
			new_host.Drive[k] = v
		}
	}
	if base.Shell != nil {
		new_host.Shell = make([]string, len(base.Shell))
		copy(new_host.Shell, base.Shell)
	}
	return new_host
}

// screen events such as keypresses
func e_events(data *HardData, fp [MODE_MAX + 1]key_event_mode_func) {
	ui := &data.ui
	if len(ui.msg_buff) != 0 {
		ui.msg_buff = ""
	}
	event := ui.s.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
		ui.dim[W], ui.dim[H], _ = term.GetSize(0)
		e_list_follow_cursor(data.litems, ui)
		ui.s.Sync()
	case *tcell.EventKey:
		if ui.mode > MODE_MAX {
			return
		} else if brk := fp[ui.mode](data, ui, *event); brk == true {
			return
		} else if ui.mode == NORMAL_MODE {
			e_list_follow_cursor(data.litems, ui)
		}
	}
}
