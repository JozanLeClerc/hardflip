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
 * hardflip: src/e_keys.go
 * Mon Apr 22 17:04:37 2024
 * Joe
 *
 * events in the keys
 */

package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

func e_normal_events(data *HardData, event tcell.EventKey) {
	if event.Key() == tcell.KeyCtrlC ||
	   event.Rune() == 'q' {
		data.ui.s.Fini()
		os.Exit(0)
	} else if event.Rune() == 'j' ||
			  event.Key() == tcell.KeyDown {
		data.litems.inc(+1)
	} else if event.Rune() == 'k' ||
			  event.Key() == tcell.KeyUp {
		data.litems.inc(-1)
	} else if event.Key() == tcell.KeyCtrlD ||
			  event.Key() == tcell.KeyPgDn {
		data.litems.inc(+(data.ui.dim[H] / 3))
	} else if event.Key() == tcell.KeyCtrlU ||
			  event.Key() == tcell.KeyPgUp {
		data.litems.inc(-(data.ui.dim[H] / 3))
	} else if event.Key() == tcell.KeyCtrlF {
		// TODO: maybe keymap these
	} else if event.Key() == tcell.KeyCtrlB {
		// TODO: maybe keymap these
	} else if event.Rune() == '}' ||
			  event.Rune() == ']' {
		if next := data.litems.curr.next_dir(); next != nil {
			data.litems.curr = next
		}
	} else if event.Rune() == '{' ||
			  event.Rune() == '[' {
		if prev := data.litems.curr.prev_dir(); prev != nil {
			data.litems.curr = prev
		}
	} else if event.Rune() == 'g' ||
			  event.Key() == tcell.KeyHome {
		data.litems.curr = data.litems.head
		data.litems.draw = data.litems.head
	} else if event.Rune() == 'G' ||
			  event.Key() == tcell.KeyEnd {
		data.litems.curr = data.litems.last
	} else if event.Rune() == 'D' &&
			  data.litems.head != nil &&
			  data.litems.curr != nil {
		data.ui.mode = DELETE_MODE
	} else if event.Rune() == 'H' {
		for curr := data.litems.last; curr != nil; curr = curr.prev {
			if curr.is_dir() == true && data.folds[curr.Dirs] == nil {
				e_fold_dir(data, curr)
			}
		}
		data.litems.curr = data.litems.head
		data.litems.draw = data.litems.curr
	} else if event.Rune() == 'h' ||
			  event.Key() == tcell.KeyLeft {
		for curr := data.litems.curr;
			curr != nil;
			curr = curr.prev {
			if curr.is_dir() == true {
				if data.folds[curr.Dirs] == nil {
					e_fold_dir(data, curr)
					data.litems.curr = curr
					data.litems.draw = data.litems.curr
					return
				} else {
					if data.folds[curr.Dirs.Parent] == nil {
						parent := curr.Dirs.Parent
						for curr_new := curr;
							curr_new != nil;
							curr_new = curr_new.prev {
							if curr_new.is_dir() == true {
								if curr_new.Dirs == parent {
									e_fold_dir(data, curr_new)
									data.litems.curr = curr_new
									data.litems.draw = data.litems.curr
									return
								} else {
									if data.folds[curr_new.Dirs] ==
									   nil {
										e_fold_dir(data, curr_new)
									}
								}
							}
						}
					}
					return
				}
			}
		}
	} else if event.Rune() == 'l' ||
			  event.Key() == tcell.KeyRight ||
			  event.Key() == tcell.KeyEnter {
		if data.litems.curr == nil {
			return
		} else if data.litems.curr.is_dir() == false {
			c_exec(data.litems.curr.Host, data.opts, &data.ui)
		} else if data.litems.curr.Dirs != nil &&
				  data.folds[data.litems.curr.Dirs] == nil {
			e_fold_dir(data, data.litems.curr)
		} else {
			e_unfold_dir(data, data.litems.curr)
		}
	} else if event.Rune() == ' ' {
		if data.litems.curr == nil ||
		   data.litems.curr.is_dir() == false {
			return
		}
		if data.litems.curr.Dirs != nil &&
		   data.folds[data.litems.curr.Dirs] == nil {
			e_fold_dir(data, data.litems.curr)
		} else {
			e_unfold_dir(data, data.litems.curr)
		}
	} else if event.Rune() == 'a' ||
			  event.Rune() == 'i' {
		data.ui.mode = INSERT_MODE
		data.ui.insert_sel = 0
		data.ui.insert_sel_ok = false
	} else if event.Key() == tcell.KeyCtrlR {
		e_reload_data(data)
	} else if event.Rune() == 'm' ||
			  event.Key() == tcell.KeyF7 {
		data.ui.mode = MKDIR_MODE
	} else if event.Rune() == 'y' {
		if data.litems.curr == nil ||
		   data.litems.curr.is_dir() == true {
			return
		}
		data.yank = data.litems.curr
	}
}

func e_delete_events(data *HardData, event tcell.EventKey) {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC ||
	   event.Rune() == 'n' {
		data.ui.mode = NORMAL_MODE
	} else if event.Key() == tcell.KeyEnter ||
			  event.Rune() == 'y' {
		if err := e_delete_host(data); err == nil {
			data.ui.mode = NORMAL_MODE
		}
	}
}

func e_error_events(data *HardData, event tcell.EventKey) {
	if event.Rune() != 0 ||
	   event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyEnter {
		data.ui.mode = NORMAL_MODE
		data.load_err = nil
	}
}
