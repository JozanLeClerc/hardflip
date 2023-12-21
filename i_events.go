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
 * Copyright (c) 2023 Joe
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. Neither the name of the organization nor the
 *    names of its contributors may be used to endorse or promote products
 *    derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY JOE ''AS IS'' AND ANY
 * EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL JOE BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * hardflip: src/i_events.go
 * Thu Dec 21 12:49:09 2023
 * Joe
 *
 * the hosts linked list
 */

package main

import (
	"os"
	"github.com/gdamore/tcell/v2"
)

func i_reload_data(data *HardData) {
	ui := data.ui
	data.lhost = c_load_data_dir(data.data_dir)
	l := data.lhost
	data.ui.sel_max = l.count()
	if ui.sel >= ui.sel_max {
		ui.sel = ui.sel_max - 1
	}
}

func i_delete_host(data *HardData) {
	ui := &data.ui
	host := data.lhost.sel(data.ui.sel)
	file_path := data.data_dir + "/" + host.Folder + host.Filename

	if err := os.Remove(file_path); err != nil {
		c_die("can't remove " + file_path, err)
	}
	data.lhost.del(data.ui.sel)
	data.lhost.reset_id()
	ui.sel_max = data.lhost.count()
	if ui.sel >= ui.sel_max {
		ui.sel = ui.sel_max - 1
	}
}

// screen events such as keypresses
func i_events(data *HardData) {
	var err error
	ui := &data.ui
	event := ui.s.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
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
				if ui.sel < ui.sel_max - 1 {
					ui.sel += 1
				}
			} else if event.Rune() == 'k' ||
			   event.Key() == tcell.KeyUp {
				if ui.sel > 0 {
					ui.sel -= 1
				}
			} else if event.Rune() == 'g' {
			   ui.sel = 0
			} else if event.Rune() == 'G' {
			   ui.sel = ui.sel_max - 1
			} else if event.Rune() == 'D' && data.lhost.head != nil {
				ui.mode = DELETE_MODE
			} else if event.Key() == tcell.KeyEnter {
				ui.s.Fini()
				c_exec(ui.sel, data.lhost)
				if data.opts.loop == false {
					os.Exit(0)
				}
				if ui.s, err = tcell.NewScreen(); err != nil {
					c_die("view", err)
				}
				if err := ui.s.Init(); err != nil {
					c_die("view", err)
				}
				ui.s.SetStyle(ui.def_style)
			}
			if event.Key() == tcell.KeyCtrlR {
				i_reload_data(data)
			}
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
