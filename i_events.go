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
 * Wed Jan 10 18:56:02 2024
 * Joe
 *
 * events in the code
 */

package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

func i_list_follow_cursor(litems *ItemsList, ui *HardUI) {
	if litems.draw_start == nil || litems.curr == nil {
		return
	}
	virt_id := litems.curr.ID - (ui.dim[H] - 4)
	for litems.draw_start.ID < virt_id &&
		litems.draw_start.next != nil {
		litems.draw_start = litems.draw_start.next
	}
	for litems.draw_start.ID > litems.curr.ID &&
		litems.draw_start.prev != nil {
		litems.draw_start = litems.draw_start.prev
	}
}

func i_unfold_dir(data *HardData, item *ItemsNode) {
	if item == nil {
		return
	}
	fold := data.folds[item]
	if fold == nil {
		return
	}
	after := item.next
	item.next = fold.head
	if fold.head != nil {
		fold.head.prev = item
	}
	if fold.last != nil {
		fold.last.next = after
	}
	if after != nil {
		after.prev = fold.last
	} else {
		data.litems.last = fold.last
	}
	delete(data.folds, item)
	for ptr := data.litems.head; ptr.next != nil; ptr = ptr.next {
		ptr.next.ID = ptr.ID + 1
	}
	item.Dirs.Folded = false
}

func i_fold_dir(data *HardData, item *ItemsNode) {
	if item == nil || item.Dirs == nil {
		return
	}
	var folded_start, folded_end, after *ItemsNode
	folds := data.folds
	folded_start = item.next
	if folded_start != nil {
		folded_start.prev = nil
		folded_end = item
	} else {
		folded_end = nil
	}
	for i := 0; folded_end != nil && i < item.Dirs.count_elements(true); i++ {
		folded_end = folded_end.next
	}
	if folded_end != nil {
		after = folded_end.next
		folded_end.next = nil
	} else {
		after = nil
	}
	tmp := ItemsList{
		folded_start,
		folded_end,
		nil,
		nil,
	}
	item.next = after
	if after != nil {
		after.prev = item
	} else {
		data.litems.last = item
	}

	folds[item] = &tmp
	for ptr := data.litems.head; ptr.next != nil; ptr = ptr.next {
		ptr.next.ID = ptr.ID + 1
	}
	item.Dirs.Folded = true
}

func i_reload_data(data *HardData) {
	data.data_dir = c_get_data_dir()
	data.ldirs = c_load_data_dir(data.data_dir, data.opts)
	data.litems = c_load_litems(data.ldirs)
	data.ui.sel_max = data.litems.last.ID
}

func i_delete_dir(data *HardData) {
	dir := data.litems.curr.Dirs
	if dir == nil {
		return
	}
	// dir_path := data.data_dir + dir.path()
	// if err := os.RemoveAll(dir_path); err != nil {
	// 	data.ui.s.Fini()
	// 	c_die("can't remove " + dir_path, err)
	// }
	tmp := data.litems.curr.prev
	data.ldirs.del(dir)
	// TODO: delete folds map reference if folded
	// TODO: finish this
	// TODO: litems ldirs and shit and lots of segv
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
			} else if event.Key() == tcell.KeyCtrlD {
				data.litems.inc(+(ui.dim[H] / 3))
			} else if event.Key() == tcell.KeyCtrlU {
				data.litems.inc(-(ui.dim[H] / 3))
			} else if event.Rune() == 'g' {
				data.litems.curr = data.litems.head
				data.litems.draw_start = data.litems.head
			} else if event.Rune() == 'G' {
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
				} else if data.litems.curr.Dirs.Folded == false {
					i_fold_dir(data, data.litems.curr)
				} else {
					i_unfold_dir(data, data.litems.curr)
				}
			} else if event.Rune() == ' ' {
				if data.litems.curr == nil ||
				   data.litems.curr.is_dir() == false {
					break
				}
				if data.litems.curr.Dirs.Folded == false {
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
