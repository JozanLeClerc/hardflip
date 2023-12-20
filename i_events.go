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
 * Wed Dec 20 11:05:16 2023
 * Joe
 *
 * the hosts linked list
 */

package main

import(
	"os"
	"github.com/gdamore/tcell/v2"
)

// screen events such as keypresses
func i_events(data *Data,
		sel *uint64, sel_max *uint64,
		term_size *[2]int,
		quit func()) {
	var err error
	event := data.s.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
		data.s.Sync()
	case *tcell.EventKey:
		if event.Key() == tcell.KeyEscape ||
		event.Key() == tcell.KeyCtrlC ||
		event.Rune() == 'q' {
			quit()
			os.Exit(0)
		}
		if event.Rune() == 'j' ||
		event.Key() == tcell.KeyDown {
			if *sel < *sel_max - 1 {
				*sel += 1
			}
		}
		if event.Rune() == 'k' ||
		event.Key() == tcell.KeyUp {
			if *sel > 0 {
				*sel -= 1
			}
		}
		if event.Key() == tcell.KeyEnter {
			quit()
			c_exec(*sel, data.lhost)
			if data.opts.loop == false {
				os.Exit(0)
			}
			if data.s, err = tcell.NewScreen(); err != nil {
				c_die("view", err)
			}
			if err := data.s.Init(); err != nil {
				c_die("view", err)
			}
			def_style := tcell.StyleDefault.
				Background(tcell.ColorReset).
				Foreground(tcell.ColorReset)
			data.s.SetStyle(def_style)
		}
		if event.Key() == tcell.KeyCtrlR {
			data.lhost = c_load_data_dir(c_get_data_dir())
			l := data.lhost
			*sel_max = l.count()
			if *sel >= *sel_max {
				*sel = *sel_max - 1
			}
		}
	}
}
