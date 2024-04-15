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
 * Thu Apr 11 16:00:44 2024
 * Joe
 *
 * events in the code
 */

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

func i_list_follow_cursor(litems *ItemsList, ui *HardUI) {
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

func i_set_unfold(data *HardData, item *ItemsNode) {
	delete(data.folds, item.Dirs)
	data.litems.reset_id()
}

func i_unfold_dir(data *HardData, item *ItemsNode) {
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
		i_set_unfold(data, item)
		return
	}
	// single empty dir
	if start == item && end == end { // HACK: i forgot why end == end
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
	data.litems.reset_id()
}

func i_fold_dir(data *HardData, item *ItemsNode) {
	if item == nil || item.Dirs == nil {
		return
	}
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
	if data.data_dir == "" {
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

func i_delete_dir(data *HardData) error {
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
		i_fold_dir(data, curr)
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

func i_delete_host(data *HardData) error {
	if data.litems.curr == nil {
		return nil
	}
	if data.litems.curr.is_dir() == true {
		return i_delete_dir(data)
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

func i_readline(event *tcell.EventKey, buffer *string) {
	if len(*buffer) > 0 &&
	(event.Key() == tcell.KeyBackspace ||
	event.Key() == tcell.KeyBackspace2) {
		*buffer = (*buffer)[:len(*buffer) - 1]
	} else if event.Key() == tcell.KeyCtrlU {
		*buffer = ""
	} else if event.Rune() >= 32 && event.Rune() <= 126 {
		*buffer += string(event.Rune())
	}
}

func i_mkdir(data *HardData, ui *HardUI) {
	if len(ui.buff) == 0 {
		return
	}
	path := "/"
	if data.litems.curr != nil {
		path = data.litems.curr.path()
	}
	if err := os.MkdirAll(data.data_dir +
		path +
		ui.buff, os.ModePerm); err != nil {
		c_error_mode("mkdir " + path[1:] + ui.buff + " failed",
		err, ui)
		return
	}
	i_reload_data(data)
	for curr := data.litems.head; curr != nil; curr = curr.next {
		if curr.is_dir() == true &&
		   curr.Dirs.Name == ui.buff &&
		   curr.Dirs.Parent.path() == path {
			data.litems.curr = curr
			return
		}
	}
}

func i_set_drive_keys(data *HardData) {
	data.insert.drive_keys = nil
	for key := range data.insert.Drive {
		data.insert.drive_keys = append(data.insert.drive_keys, key)
	}
}

func i_set_protocol_defaults(data *HardData, in *HostNode) {
	switch in.Protocol {
	case PROTOCOL_SSH:
		in.Port = 22
		data.ui.insert_sel_max = INS_SSH_OK
	case PROTOCOL_RDP:
		in.Port = 3389
		in.Quality = 2
		in.Width = 1600
		in.Height = 1200
		in.Dynamic = true
		data.insert.Drive = map[string]string{ // WARN: this is a test
			"qwe": "a",
			"asd": "aaaa",
			"zxc": "aaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		}
		i_set_drive_keys(data)
		data.ui.insert_sel_max = INS_RDP_OK + len(data.insert.Drive)
	case PROTOCOL_CMD:
		in.Shell = []string{"/bin/sh", "-c"}
		data.ui.insert_sel_max = 2
	case PROTOCOL_OS:
		in.Stack.RegionName = "eu-west-0"
		in.Stack.IdentityAPI = "3"
		in.Stack.ImageAPI    = "2"
		in.Stack.NetworkAPI  = "2"
		in.Stack.VolumeAPI   = "3.42"
		in.Stack.EndpointType = "publicURL"
		in.Stack.Interface = "public"
		data.ui.insert_sel_max = 2
	}
}

// screen events such as keypresses
func i_events(data *HardData) {
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
				ui.mode = DELETE_MODE
			} else if event.Rune() == 'H' {
				for curr := data.litems.last; curr != nil; curr = curr.prev {
					if curr.is_dir() == true && data.folds[curr.Dirs] == nil {
						i_fold_dir(data, curr)
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
							i_fold_dir(data, curr)
							data.litems.curr = curr
							data.litems.draw = data.litems.curr
							break
						} else {
							if data.folds[curr.Dirs.Parent] == nil {
								parent := curr.Dirs.Parent
								for curr_new := curr;
									curr_new != nil;
									curr_new = curr_new.prev {
									if curr_new.is_dir() == true {
										if curr_new.Dirs == parent {
											i_fold_dir(data, curr_new)
											data.litems.curr = curr_new
											data.litems.draw = data.litems.curr
											break
										} else {
											if data.folds[curr_new.Dirs] ==
											   nil {
												i_fold_dir(data, curr_new)
											}
										}
									}
								}
							}
							break
						}
					}
				}
			} else if event.Rune() == 'l' ||
					  event.Key() == tcell.KeyRight ||
					  event.Key() == tcell.KeyEnter {
				if data.litems.curr == nil {
					break
				} else if data.litems.curr.is_dir() == false {
					c_exec(data.litems.curr.Host, data.opts, ui)
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
			} else if event.Rune() == 'a' ||
					  event.Rune() == 'i' {
				data.ui.mode = INSERT_MODE
				data.ui.insert_sel = 0
				data.ui.insert_sel_ok = false
			} else if event.Key() == tcell.KeyCtrlR {
				event = nil
				i_reload_data(data)
			} else if event.Rune() == 'm' ||
					  event.Key() == tcell.KeyF7 {
				data.ui.mode = MKDIR_MODE
			}
			i_list_follow_cursor(data.litems, ui)
		case DELETE_MODE:
			if event.Key() == tcell.KeyEscape ||
			   event.Key() == tcell.KeyCtrlC ||
			   event.Rune() == 'n' {
				ui.mode = NORMAL_MODE
			} else if event.Key() == tcell.KeyEnter ||
					  event.Rune() == 'y' {
				if err := i_delete_host(data); err == nil {
					ui.mode = NORMAL_MODE
				}
			}
		case ERROR_MODE:
			if event.Rune() != 0 ||
			   event.Key() == tcell.KeyEscape ||
			   event.Key() == tcell.KeyEnter {
				ui.mode = NORMAL_MODE
				data.load_err = nil
			}
		case WELCOME_MODE:
			if event.Key() == tcell.KeyEscape ||
			   event.Key() == tcell.KeyCtrlC {
				ui.s.Fini()
				os.Exit(0)
			}
			if len(data.opts.GPG) == 0 {
				if event.Rune() < '1' || event.Rune() > '9' {
					break
				} else {
					data.opts.GPG = data.keys[event.Rune() - 48 - 1][0]
					ui.s.HideCursor()
				}
			} else {
				if event.Rune() == 'y' {
					ui.mode = NORMAL_MODE
					c_write_options(data.opts.file, data.opts, &data.load_err)
				} else if event.Rune() == 'n' {
					data.opts.GPG = ""
				}
			}
		case INSERT_MODE:
			if data.insert == nil {
				if event.Key() == tcell.KeyEscape ||
				event.Key() == tcell.KeyCtrlC {
					ui.s.HideCursor()
					data.ui.mode = NORMAL_MODE
					data.ui.insert_sel = 0
					data.insert = nil
					ui.buff = ""
				} else if event.Key() == tcell.KeyEnter {
					if ui.buff == "" {
						ui.s.HideCursor()
						data.ui.mode = NORMAL_MODE
						data.ui.insert_sel = 0
						data.ui.insert_sel_ok = false
						data.insert = nil
						ui.buff = ""
						break
					}
					ui.s.HideCursor()
					data.insert = &HostNode{}
					data.insert.Protocol = 1 // WARN: tests only, remove this
					i_set_protocol_defaults(data, data.insert)
					data.insert.Name = ui.buff
					ui.buff = ""
					if data.litems.curr != nil {
						data.insert.parent = data.litems.curr.path_node()
					} else {
						data.insert.parent = data.ldirs.head
					}
				} else {
					i_readline(event, &data.ui.buff)
				}
			} else if data.insert != nil {
				if data.insert_err != nil {
					if event.Rune() != 0 ||
					   event.Key() == tcell.KeyEscape ||
					   event.Key() == tcell.KeyEnter {
						data.insert_err = nil
					}
				} else if data.ui.insert_sel_ok == false {
					if event.Key() == tcell.KeyEscape ||
					   event.Key() == tcell.KeyCtrlC ||
					   event.Rune() == 'q' {
						ui.s.HideCursor()
						data.ui.mode = NORMAL_MODE
						data.ui.insert_sel = 0
						data.insert = nil
						ui.buff = ""
					} else if event.Rune() == 'j' ||
							  event.Key() == tcell.KeyDown ||
							  event.Key() == tcell.KeyTab {
						if data.insert.Protocol == PROTOCOL_RDP &&
						   data.ui.insert_sel == INS_PROTOCOL {
							data.ui.insert_sel = INS_RDP_HOST
						} else if data.ui.insert_sel < data.ui.insert_sel_max {
							data.ui.insert_sel += 1
						}
					} else if event.Rune() == 'k' ||
							  event.Key() == tcell.KeyUp {
						if data.insert.Protocol == PROTOCOL_RDP &&
						   data.ui.insert_sel == INS_RDP_HOST {
							data.ui.insert_sel = INS_PROTOCOL
						} else if data.ui.insert_sel > INS_PROTOCOL {
							data.ui.insert_sel -= 1
						}
					} else if event.Rune() == 'g' ||
							  event.Rune() == 'h' ||
							  event.Key() == tcell.KeyLeft {
							data.ui.insert_sel = INS_PROTOCOL
					} else if event.Rune() == 'G' ||
							  event.Rune() == 'l' ||
							  event.Key() == tcell.KeyRight {
							data.ui.insert_sel = data.ui.insert_sel_max
					} else if event.Rune() == 'i' ||
							  event.Rune() == 'a' ||
							  event.Key() == tcell.KeyEnter {
						if data.ui.insert_sel == INS_RDP_DYNAMIC {
							if data.insert.Dynamic == true {
								data.insert.Dynamic = false
							} else {
								data.insert.Dynamic = true
							}
							ui.buff = ""
							ui.s.HideCursor()
							break
						}
						data.ui.insert_sel_ok = true
						switch data.ui.insert_sel {
						case INS_SSH_HOST,
							 INS_RDP_HOST:
							ui.buff = data.insert.Host
						case INS_SSH_PORT,
							 INS_RDP_PORT:
							if data.insert.Port > 0 {
								ui.buff = strconv.Itoa(int(data.insert.Port))
							}
						case INS_SSH_USER,
							 INS_RDP_USER:
							ui.buff = data.insert.User
						case INS_SSH_PASS,
							 INS_RDP_PASS:
							break
						case INS_SSH_PRIV: ui.buff = data.insert.Priv
						case INS_SSH_JUMP_HOST: ui.buff = data.insert.Jump.Host
						case INS_SSH_JUMP_PORT:
							if data.insert.Jump.Port > 0 {
								ui.buff = strconv.Itoa(int(
								data.insert.Jump.Port))
							}
						case INS_SSH_JUMP_USER: ui.buff = data.insert.Jump.User
						case INS_SSH_JUMP_PASS: break
						case INS_SSH_JUMP_PRIV: ui.buff = data.insert.Jump.Priv
						case INS_RDP_DOMAIN: ui.buff = data.insert.Domain
						case INS_RDP_FILE: ui.buff = data.insert.RDPFile
						case INS_RDP_SCREENSIZE: break
						case INS_RDP_DYNAMIC: break
						case INS_RDP_QUALITY: break
						case INS_RDP_DRIVE + len(data.insert.Drive): break
						case INS_SSH_OK,
							 INS_RDP_OK + len(data.insert.Drive):
							data.ui.insert_sel_ok = false
							i_insert_check_ok(data, data.insert)
							if data.insert_err != nil {
								break
							}
							i_insert_host(data, data.insert)
						}
					}
				} else {
					if event.Key() == tcell.KeyEscape ||
					   event.Key() == tcell.KeyCtrlC {
						data.ui.insert_sel_ok = false
						ui.buff = ""
						ui.drives_buff = ""
						ui.s.HideCursor()
					}
					switch data.ui.insert_sel {
					case INS_PROTOCOL:
						if event.Rune() < '1' || event.Rune() > '4' {
							data.ui.insert_sel_ok = false
							ui.buff = ""
							ui.s.HideCursor()
							break
						} else {
							name := data.insert.Name
							parent := data.insert.parent
							data.insert = nil
							data.insert = &HostNode{}
							data.insert.Name = name
							data.insert.parent = parent
							data.insert.Protocol = int8(event.Rune() - 48 - 1)
							data.ui.insert_sel_ok = false
							ui.s.HideCursor()
							i_set_protocol_defaults(data, data.insert)
						}
					case INS_RDP_SCREENSIZE:
						if event.Rune() < '1' || event.Rune() > '7' {
							data.ui.insert_sel_ok = false
							ui.buff = ""
							ui.s.HideCursor()
							break
						} else {
							s := strings.Split(
								RDP_SCREENSIZE[uint8(event.Rune() - 48 - 1)],
								"x")
							if len(s) != 2 {
								return
							}
							tmp, _ := strconv.Atoi(s[W])
							data.insert.Width = uint16(tmp)
							tmp, _ = strconv.Atoi(s[H])
							data.insert.Height = uint16(tmp)
							data.ui.insert_sel_ok = false
							ui.s.HideCursor()
						}
					case INS_RDP_QUALITY:
						if event.Rune() < '1' || event.Rune() > '3' {
							data.ui.insert_sel_ok = false
							ui.buff = ""
							ui.s.HideCursor()
							break
						} else {
							data.insert.Quality = uint8(event.Rune() - 48 - 1)
							data.ui.insert_sel_ok = false
							ui.s.HideCursor()
						}
					case INS_RDP_DRIVE + len(data.insert.Drive):
						if len(data.ui.drives_buff) == 0 {
							if event.Key() == tcell.KeyEnter {
								if len(ui.buff) == 0 {
									data.ui.insert_sel_ok = false
									data.ui.drives_buff = ""
									ui.buff = ""
									ui.s.HideCursor()
									break
								}
								data.ui.drives_buff = ui.buff
								ui.buff = ""
							} else {
								i_readline(event, &data.ui.buff)
							}
						} else {
							if event.Key() == tcell.KeyEnter {
								if len(ui.buff) == 0 {
									data.ui.insert_sel_ok = false
									data.ui.drives_buff = ""
									ui.buff = ""
									ui.s.HideCursor()
									break
								}
								data.ui.insert_sel_ok = false
								if len(data.insert.Drive) == 0 {
									data.insert.Drive = make(map[string]string)
								}
								data.insert.Drive[ui.drives_buff] = ui.buff
								ui.s.Fini()
								fmt.Println(data.insert.Drive)
								os.Exit(0)
								i_set_drive_keys(data)
								data.ui.insert_sel_max = INS_RDP_OK +
														 len(data.insert.Drive)
								ui.drives_buff = ""
								ui.buff = ""
								ui.s.HideCursor()
							} else {
								i_readline(event, &data.ui.buff)
							}
						}
					case INS_SSH_HOST,
						 INS_SSH_PORT,
						 INS_SSH_USER,
						 INS_SSH_PASS,
						 INS_SSH_PRIV,
						 INS_SSH_JUMP_HOST,
						 INS_SSH_JUMP_PORT,
						 INS_SSH_JUMP_USER,
						 INS_SSH_JUMP_PASS,
						 INS_SSH_JUMP_PRIV,
						 INS_RDP_HOST,
						 INS_RDP_PORT,
						 INS_RDP_DOMAIN,
						 INS_RDP_USER,
						 INS_RDP_PASS,
						 INS_RDP_FILE:
						if event.Key() == tcell.KeyEnter {
							switch data.ui.insert_sel {
							case INS_SSH_HOST,
								 INS_RDP_HOST:
								data.insert.Host = ui.buff
							case INS_SSH_PORT,
								 INS_RDP_PORT:
								tmp, _ := strconv.Atoi(ui.buff)
								data.insert.Port = uint16(tmp)
							case INS_SSH_USER,
								 INS_RDP_USER:
								data.insert.User = ui.buff
							case INS_SSH_PASS,
								 INS_RDP_PASS:
								data.insert.Pass, _ = c_encrypt_str(ui.buff,
														 data.opts.GPG)
							case INS_SSH_PRIV: data.insert.Priv = ui.buff
							case INS_SSH_JUMP_HOST:
								data.insert.Jump.Host = ui.buff
								if len(ui.buff) > 0 {
									data.insert.Jump.Port = 22
								} else {
									data.insert.Jump.Port = 0
								}
							case INS_SSH_JUMP_PORT:
								tmp, _ := strconv.Atoi(ui.buff)
								data.insert.Jump.Port = uint16(tmp)
							case INS_SSH_JUMP_USER:
								data.insert.Jump.User = ui.buff
							case INS_SSH_JUMP_PASS:
								data.insert.Jump.Pass, _ =
								c_encrypt_str(ui.buff, data.opts.GPG)
							case INS_SSH_JUMP_PRIV:
								data.insert.Jump.Priv = ui.buff
							case INS_RDP_DOMAIN:
								data.insert.Domain = ui.buff
							case INS_RDP_FILE:
								data.insert.RDPFile = ui.buff
							}
							data.ui.insert_sel_ok = false
							ui.buff = ""
							ui.s.HideCursor()
						} else {
							i_readline(event, &data.ui.buff)
						}
					}
					if len(data.insert.Drive) > 0 &&
					   data.ui.insert_sel >= INS_RDP_DRIVE &&
					   data.ui.insert_sel < INS_RDP_DRIVE +
					   len(data.insert.Drive) {
						if event.Rune() == 'y' ||
						   event.Rune() == 'Y' ||
						   event.Key() == tcell.KeyEnter {
							delete(data.insert.Drive,
								   data.insert.drive_keys[data.ui.insert_sel -
								   INS_RDP_DRIVE])
							data.ui.insert_sel_max = INS_RDP_OK +
								len(data.insert.Drive)
							i_set_drive_keys(data)
						}
						data.ui.insert_sel_ok = false
					}
				}
			}
		case MKDIR_MODE:
			if event.Key() == tcell.KeyEscape ||
			   event.Key() == tcell.KeyCtrlC {
				ui.s.HideCursor()
				ui.mode = NORMAL_MODE
				ui.buff = ""
				data.insert = nil
			} else if event.Key() == tcell.KeyEnter {
				i_mkdir(data, ui)
				ui.s.HideCursor()
				ui.mode = NORMAL_MODE
				ui.buff = ""
			} else {
				i_readline(event, &data.ui.buff)
			}
		}
	}
}
