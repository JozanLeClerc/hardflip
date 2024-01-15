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
 * hardflip: src/i_events.go
 * Mon Jan 15 13:54:57 2024
 * Joe
 *
 * events in the code
 */

package main

import (
	// "fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

func i_list_follow_cursor(litems *ItemsList, ui *HardUI) {
	if litems.draw_start == nil || litems.curr == nil {
		return
	}
	virt_id := litems.curr.ID - (ui.dim[H] - 4) + 4
	for litems.draw_start.ID < virt_id &&
		litems.draw_start.next != nil {
		litems.draw_start = litems.draw_start.next
	}
	for litems.draw_start.ID > litems.curr.ID - 4 &&
		litems.draw_start.prev != nil {
		litems.draw_start = litems.draw_start.prev
	}
}

func i_set_unfold(data *HardData, item *ItemsNode) {
	delete(data.folds, item.Dirs)
	for ptr := data.litems.head; ptr.next != nil; ptr = ptr.next {
		ptr.next.ID = ptr.ID + 1
	}
}

func i_unfold_dir(data *HardData, item *ItemsNode) {
	if item == nil || item.Dirs == nil {
		return
	}
	fold := data.folds[item.Dirs]
	if fold == nil {
		return
	}
	// after := item.next
	start, end := fold.head, fold.last
	// last empty dir
	if start == nil && end == nil {
		i_set_unfold(data, item)
		return
	}
	// single empty dir
	if start == item && end == end {
		i_set_unfold(data, item)
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
	i_set_unfold(data, item)
}

func i_set_fold(data *HardData, curr, start, end *ItemsNode) {
	folds := data.folds
	tmp := ItemsList{
		start,
		end,
		nil,
		nil,
	}

	folds[curr.Dirs] = &tmp
	for ptr := data.litems.head; ptr.next != nil; ptr = ptr.next {
		ptr.next.ID = ptr.ID + 1
	}
}

func i_fold_dir(data *HardData, item *ItemsNode) {
	if item == nil || item.Dirs == nil {
		return
	}
	// var start, end, after *ItemsNode
	var start, end *ItemsNode
	start = item.next
	// last dir + empty
	if start == nil {
		i_set_fold(data, item, nil, nil)
		return
	}
	// empty dir
	if start.Dirs != nil && start.Dirs.Depth <= item.Dirs.Depth {
		i_set_fold(data, item, item, item)
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
		i_set_fold(data, item, start, end)
		return
	}
	// this is not the end
	end = next_dir.prev 
	end.next = nil
	item.next = next_dir
	next_dir.prev = item
	i_set_fold(data, item, start, end)
}

func i_reload_data(data *HardData) {
	data.data_dir = c_get_data_dir()
	data.ldirs = c_load_data_dir(data.data_dir, data.opts)
	data.litems = c_load_litems(data.ldirs)
	data.folds = make(map[*DirsNode]*ItemsList)
	data.ui.sel_max = data.litems.last.ID
}

func i_delete_dir(data *HardData) {
	if data.litems.curr == nil || data.litems.curr.Dirs == nil {
		return
	}
	curr := data.litems.curr
	if data.folds[curr.Dirs] == nil {
		i_fold_dir(data, curr)
	}
	delete(data.folds, curr.Dirs)
	prev_dir := data.ldirs.prev(curr.Dirs)
	curr.prev.next = curr.next
	if curr.next != nil {
		curr.next.prev = curr.prev
		prev_dir.next = curr.Dirs.next
		data.litems.curr = curr.next
	} else {
		data.litems.last = curr.prev
		prev_dir.next = nil
		data.ldirs.last = prev_dir
		data.ldirs.last.next = nil
		data.litems.curr = curr.prev
	}
	for ptr := data.litems.head; ptr.next != nil; ptr = ptr.next {
		ptr.next.ID = ptr.ID + 1
	}
	dir_path := data.data_dir + dir.path()
	if err := os.RemoveAll(dir_path); err != nil {
		data.ui.s.Fini()
		c_die("can't remove " + dir_path, err)
	}
	// TODO: delete folds map reference if folded - OK
	// TEST: single empty dir
	// TEST: single non-empty dir
	// TEST: first dir
	// TEST: last dir
	// TEST: last dir 4m+
	// TEST: folded
}

func i_delete_host(data *HardData) {
	if data.litems.curr == nil {
		return
	}
	if data.litems.curr.is_dir() == true {
		i_delete_dir(data)
		return
	}
	host := data.litems.curr.Host
	if host == nil {
		return
	}
	file_path := data.data_dir + host.Parent.path() + host.Filename

	if err := os.Remove(file_path); err != nil {
		data.ui.s.Fini()
		c_die("can't remove " + file_path, err)
	}
	tmp := data.litems.curr.prev
	host.Parent.lhost.del(host)
	data.litems.del(data.litems.curr)
	if tmp == nil {
		tmp = data.litems.head
	}
	data.litems.curr = tmp
	if data.litems.last != nil {
		data.ui.sel_max = data.litems.last.ID
	} else {
		data.ui.sel_max = 0
	}
}

// screen events such as keypresses
func i_events(data *HardData) {
	var err error
	ui := &data.ui
	event := ui.s.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
		ui.dim[W], ui.dim[H], _ = term.GetSize(0)
		i_list_follow_cursor(data.litems, ui)
		ui.s.Sync()
	case *tcell.EventKey:
		switch ui.mode {
		case NORMAL_MODE:
			if event.Key() == tcell.KeyCtrlC ||
			   event.Rune() == 'q' {
				ui.s.Fini()
				os.Exit(0)
			} else if event.Rune() == 'j' ||
					  event.Key() == tcell.KeyDown {
				data.litems.inc(+1)
			} else if event.Rune() == 'k' ||
					  event.Key() == tcell.KeyUp {
				data.litems.inc(-1)
			} else if event.Key() == tcell.KeyCtrlD ||
					  event.Key() == tcell.KeyPgDn {
				data.litems.inc(+(ui.dim[H] / 3))
			} else if event.Key() == tcell.KeyCtrlU ||
					  event.Key() == tcell.KeyPgUp {
				data.litems.inc(-(ui.dim[H] / 3))
			} else if event.Rune() == '}' ||
					  event.Rune() == ']' {
				// TODO: next dir
			} else if event.Rune() == '{' ||
					  event.Rune() == '[' {
				// TODO: prev dir
			} else if event.Rune() == 'g' ||
					  event.Key() == tcell.KeyHome {
				data.litems.curr = data.litems.head
				data.litems.draw_start = data.litems.head
			} else if event.Rune() == 'G' ||
					  event.Key() == tcell.KeyEnd {
				data.litems.curr = data.litems.last
			} else if event.Rune() == 'D' &&
					  data.ldirs.head != nil &&
					  ui.sel_max != 0 {
				ui.mode = DELETE_MODE
			} else if event.Key() == tcell.KeyEnter {
				if data.litems.curr == nil {
					break
				} else if data.litems.curr.is_dir() == false {
					ui.s.Fini()
					c_exec(data.litems.curr.Host)
					if data.opts.Loop == false {
						os.Exit(0)
					} else {
						if ui.s, err = tcell.NewScreen(); err != nil {
							c_die("view", err)
						}
						if err := ui.s.Init(); err != nil {
							c_die("view", err)
						}
						ui.s.SetStyle(ui.def_style)
					}
				} else if data.litems.curr.Dirs != nil &&
						  data.folds[data.litems.curr.Dirs] == nil {
					i_fold_dir(data, data.litems.curr)
				} else {
					i_unfold_dir(data, data.litems.curr)
				}
			} else if event.Rune() == ' ' {
				if data.litems.curr == nil ||
				   data.litems.curr.is_dir() == false {
					break
				}
				if data.litems.curr.Dirs != nil &&
				   data.folds[data.litems.curr.Dirs] == nil {
					i_fold_dir(data, data.litems.curr)
				} else {
					i_unfold_dir(data, data.litems.curr)
				}
			} else if event.Key() == tcell.KeyCtrlR {
				i_reload_data(data)
			}
			i_list_follow_cursor(data.litems, ui)
		case DELETE_MODE:
			if event.Key() == tcell.KeyEscape ||
			   event.Key() == tcell.KeyCtrlC ||
			   event.Rune() == 'q' ||
			   event.Rune() == 'n' {
				ui.mode = NORMAL_MODE
			} else if event.Rune() == 'y' {
				i_delete_host(data)
				ui.mode = NORMAL_MODE
			}
		}	
	}
}
