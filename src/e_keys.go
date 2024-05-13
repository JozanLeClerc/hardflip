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
 * Mon May 13 10:51:48 2024
 * Joe
 *
 * events in the keys
 */

package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func e_normal_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
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
					return true
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
									return true
								} else {
									if data.folds[curr_new.Dirs] ==
									   nil {
										e_fold_dir(data, curr_new)
									}
								}
							}
						}
					}
					return true
				}
			}
		}
	} else if event.Rune() == 'l' ||
			  event.Key() == tcell.KeyRight ||
			  event.Key() == tcell.KeyEnter {
		if data.litems.curr == nil {
			return true
		} else if data.litems.curr.is_dir() == false {
			c_exec(data.litems.curr.Host, data.opts, ui)
		} else if data.litems.curr.Dirs != nil &&
				  data.folds[data.litems.curr.Dirs] == nil {
			e_fold_dir(data, data.litems.curr)
		} else {
			e_unfold_dir(data, data.litems.curr)
		}
	} else if event.Rune() == ' ' {
		if data.litems.curr == nil ||
		   data.litems.curr.is_dir() == false {
			return true
		}
		if data.litems.curr.Dirs != nil &&
		   data.folds[data.litems.curr.Dirs] == nil {
			e_fold_dir(data, data.litems.curr)
		} else {
			e_unfold_dir(data, data.litems.curr)
		}
	} else if event.Rune() == 'a' ||
			  event.Rune() == 'i' {
		ui.mode = INSERT_MODE
		ui.insert_method = INSERT_ADD
		ui.insert_sel = 0
		ui.insert_sel_ok = false
		ui.insert_scroll = 0
	} else if event.Rune() == 'e' &&
			  data.litems.curr != nil &&
			  data.litems.curr.is_dir() == false {
		tmp := e_deep_copy_host(data.litems.curr.Host)
		data.insert = &tmp
		e_set_protocol_max(data, data.insert)
		if data.insert.Protocol == PROTOCOL_RDP && data.insert.Drive != nil {
			e_set_drive_keys(data)
		}
		ui.mode = INSERT_MODE
		ui.insert_method = INSERT_EDIT
		ui.insert_sel = INS_PROTOCOL
		ui.insert_sel_ok = false
	} else if event.Key() == tcell.KeyCtrlR {
		e_reload_data(data)
	} else if event.Rune() == 'm' ||
			  event.Key() == tcell.KeyF7 {
		ui.mode = MKDIR_MODE
	} else if event.Rune() == 'y' &&
			(data.litems.curr == nil ||
			 data.litems.curr.is_dir() == true) == false {
		ui.insert_method = INSERT_COPY
		data.yank = data.litems.curr
		ui.msg_buff = "yanked " + data.yank.Host.Name +
			" (" + data.yank.Host.parent.path() + data.yank.Host.filename + ")"
	} else if event.Rune() == 'd' &&
			(data.litems.curr == nil ||
			 data.litems.curr.is_dir() == true) == false {
		ui.insert_method = INSERT_MOVE
		data.yank = data.litems.curr
		ui.msg_buff = "yanked " + data.yank.Host.Name +
			" (" + data.yank.Host.parent.path() + data.yank.Host.filename + ")"
	} else if event.Rune() == 'p' && data.yank != nil {
		new_host := e_deep_copy_host(data.yank.Host)
		if ui.insert_method == INSERT_COPY {
			new_host.Name += " (copy)"
		} else if ui.insert_method == INSERT_MOVE && data.litems.curr != nil &&
				  data.litems.curr.path_node() == data.yank.path_node() {
			data.yank = nil
			return true
		}
		if data.litems.curr.is_dir() == true {
			new_host.parent = data.litems.curr.Dirs
			if data.folds[data.litems.curr.Dirs] != nil {
				e_unfold_dir(data, data.litems.curr)
			}
		} else {
			new_host.parent = data.litems.curr.Host.parent
		}
		if err := i_insert_host(data, &new_host); err != nil {
			return true
		}
		if ui.insert_method == INSERT_MOVE {
			data.litems.del(data.yank)
			file_path := data.data_dir + data.yank.Host.parent.path() +
						 data.yank.Host.filename
			if err := os.Remove(file_path); err != nil {
				c_error_mode("can't remove " + file_path, err, &data.ui)
				return true
			}
		}
		data.yank = nil
		ui.msg_buff = "pasted " + new_host.Name
	} else if (event.Rune() == 'c' ||
			   event.Rune() == 'C' ||
			   event.Rune() == 'A') &&
			  data.litems.curr != nil {
		ui.mode = RENAME_MODE
		if data.litems.curr.is_dir() == false {
			ui.buff.insert(data.litems.curr.Host.Name)
		} else {
			ui.buff.insert(data.litems.curr.Dirs.Name)
		}
	} else if event.Rune() == '?' {
		ui.mode = HELP_MODE
		ui.help_scroll = 0
	}
	return false
}

func e_delete_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC ||
	   event.Rune() == 'n' {
		ui.mode = NORMAL_MODE
	} else if event.Key() == tcell.KeyEnter ||
			  event.Rune() == 'y' {
		if data.yank == data.litems.curr {
			data.yank = nil
		}
		if err := e_delete_host(data); err == nil {
			ui.mode = NORMAL_MODE
			return true
		}
	}
	return false
}

func e_load_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
	return true
}

func e_error_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
	if event.Rune() != 0 ||
	   event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyEnter {
		ui.mode = NORMAL_MODE
		data.load_err = nil
	}
	return false
}

func e_welcome_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC {
		ui.s.Fini()
		os.Exit(0)
	}
	if len(data.opts.GPG) == 0 {
		if event.Rune() < '1' || event.Rune() > '9' {
			return true
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
	return false
}

func e_mkdir_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC {
		ui.s.HideCursor()
		ui.mode = NORMAL_MODE
		ui.buff.empty()
		data.insert = nil
	} else if event.Key() == tcell.KeyEnter {
		e_mkdir(data, ui)
		ui.s.HideCursor()
		ui.mode = NORMAL_MODE
		ui.buff.empty()
	} else {
		e_readline(event, &ui.buff)
	}
	return false
}

func e_insert_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
	if data.insert == nil {
		if event.Key() == tcell.KeyEscape ||
		event.Key() == tcell.KeyCtrlC {
			ui.s.HideCursor()
			ui.mode = NORMAL_MODE
			ui.insert_sel = 0
			data.insert = nil
			ui.buff.empty()
		} else if event.Key() == tcell.KeyEnter {
			if ui.buff.len() == 0 {
				ui.s.HideCursor()
				ui.mode = NORMAL_MODE
				ui.insert_sel = 0
				ui.insert_sel_ok = false
				data.insert = nil
				ui.buff.empty()
				return true
			}
			ui.s.HideCursor()
			data.insert = &HostNode{}
			e_set_protocol_defaults(data, data.insert)
			data.insert.Name = ui.buff.str()
			ui.buff.empty()
			if data.litems.curr != nil {
				data.insert.parent = data.litems.curr.path_node()
			} else {
				data.insert.parent = data.ldirs.head
			}
		} else {
			e_readline(event, &ui.buff)
		}
	} else if data.insert != nil {
		if data.insert_err != nil {
			if event.Rune() != 0 ||
			   event.Key() == tcell.KeyEscape ||
			   event.Key() == tcell.KeyEnter {
				data.insert_err = nil
			}
		} else if ui.insert_sel_ok == false {
			if event.Key() == tcell.KeyEscape ||
			   event.Key() == tcell.KeyCtrlC ||
			   event.Rune() == 'q' {
				ui.s.HideCursor()
				ui.mode = NORMAL_MODE
				ui.insert_sel = 0
				data.insert = nil
				ui.buff.empty()
			} else if event.Rune() == 'j' ||
					  event.Key() == tcell.KeyDown ||
					  event.Key() == tcell.KeyTab {
				if data.insert.Protocol == PROTOCOL_RDP &&
				   ui.insert_sel == INS_PROTOCOL {
					ui.insert_sel = INS_RDP_HOST
				} else if data.insert.Protocol == PROTOCOL_RDP &&
						  ui.insert_sel == INS_RDP_JUMP_HOST +
							len(data.insert.Drive) &&
						  len(data.insert.Jump.Host) == 0 {
					ui.insert_sel = INS_RDP_NOTE + len(data.insert.Drive)
				} else if data.insert.Protocol == PROTOCOL_CMD &&
						  ui.insert_sel == INS_PROTOCOL {
					ui.insert_sel = INS_CMD_CMD
				} else if data.insert.Protocol == PROTOCOL_OS &&
						  ui.insert_sel == INS_PROTOCOL {
					ui.insert_sel = INS_OS_HOST
				} else if data.insert.Protocol == PROTOCOL_SSH &&
						  ui.insert_sel == INS_SSH_JUMP_HOST &&
						  len(data.insert.Jump.Host) == 0 {
					ui.insert_sel = INS_SSH_NOTE
				} else if ui.insert_sel < ui.insert_sel_max {
					ui.insert_sel += 1
				}
				if ui.insert_butt == false {
					ui.insert_scroll += 2
				}
			} else if event.Rune() == 'k' ||
					  event.Key() == tcell.KeyUp {
				if data.insert.Protocol == PROTOCOL_RDP &&
				   ui.insert_sel == INS_RDP_HOST {
					ui.insert_sel = INS_PROTOCOL
				} else if data.insert.Protocol == PROTOCOL_RDP &&
						  ui.insert_sel == INS_RDP_NOTE +
							len(data.insert.Drive) &&
						  len(data.insert.Jump.Host) == 0 {
					ui.insert_sel = INS_RDP_JUMP_HOST + len(data.insert.Drive)
				} else if data.insert.Protocol == PROTOCOL_CMD &&
						  ui.insert_sel == INS_CMD_CMD {
					ui.insert_sel = INS_PROTOCOL
				} else if data.insert.Protocol == PROTOCOL_OS &&
						  ui.insert_sel == INS_OS_HOST {
					ui.insert_sel = INS_PROTOCOL
				} else if data.insert.Protocol == PROTOCOL_SSH &&
						  ui.insert_sel == INS_SSH_NOTE &&
						  len(data.insert.Jump.Host) == 0 {
					ui.insert_sel = INS_SSH_JUMP_HOST
				} else if ui.insert_sel > INS_PROTOCOL {
					ui.insert_sel -= 1
				}
				if ui.insert_scroll > 0 {
					ui.insert_scroll -= 2
					if ui.insert_scroll < 0 {
						ui.insert_scroll = 0
					}
				}
			} else if event.Rune() == 'g' ||
					  event.Rune() == 'h' ||
					  event.Key() == tcell.KeyLeft {
				ui.insert_sel = INS_PROTOCOL
				ui.insert_scroll = 0
			} else if event.Rune() == 'G' ||
					  event.Rune() == 'l' ||
					  event.Key() == tcell.KeyRight {
				ui.insert_sel = ui.insert_sel_max
				for data.ui.insert_butt == false {
					ui.insert_scroll += 2
					i_draw_insert_panel(&data.ui, data.insert, data.home_dir)
				}
				data.ui.s.Show()
			} else if event.Rune() == 'i' ||
					  event.Rune() == 'a' ||
					  event.Rune() == ' ' ||
					  event.Key() == tcell.KeyEnter {
				ui.insert_sel_ok = true
				switch ui.insert_sel {
				case INS_SSH_OK,
					 INS_RDP_OK + len(data.insert.Drive),
					 INS_CMD_OK,
					 INS_OS_OK:
					ui.insert_sel_ok = false
					i_insert_check_ok(data, data.insert)
					if data.insert_err != nil {
						return true
					}
					i_insert_host(data, data.insert)
				case INS_SSH_HOST,
					 INS_RDP_HOST,
					 INS_OS_HOST:
					ui.buff.insert(data.insert.Host)
				case INS_SSH_PORT,
					 INS_RDP_PORT:
					if data.insert.Port > 0 {
						ui.buff.insert(strconv.Itoa(int(data.insert.Port)))
					}
				case INS_SSH_USER,
					 INS_RDP_USER,
					 INS_OS_USER:
					ui.buff.insert(data.insert.User)
				case INS_SSH_PASS,
					 INS_RDP_PASS,
					 INS_OS_PASS:
					return true
				case INS_SSH_PRIV: ui.buff.insert(data.insert.Priv)
				case INS_SSH_EXEC: ui.buff.insert(data.insert.Exec)
				case INS_SSH_JUMP_HOST,
					 INS_RDP_JUMP_HOST + len(data.insert.Drive):
					ui.buff.insert(data.insert.Jump.Host)
				case INS_SSH_JUMP_PORT,
					 INS_RDP_JUMP_PORT + len(data.insert.Drive):
					if data.insert.Jump.Port > 0 {
						ui.buff.insert(strconv.Itoa(int(data.insert.Jump.Port)))
					}
				case INS_SSH_JUMP_USER,
					 INS_RDP_JUMP_USER + len(data.insert.Drive):
					ui.buff.insert(data.insert.Jump.User)
				case INS_SSH_JUMP_PASS,
					 INS_RDP_JUMP_PASS + len(data.insert.Drive):
					return true
				case INS_SSH_JUMP_PRIV,
					 INS_RDP_JUMP_PRIV + len(data.insert.Drive):
					ui.buff.insert(data.insert.Jump.Priv)
				case INS_RDP_DOMAIN: ui.buff.insert(data.insert.Domain)
				case INS_RDP_FILE: ui.buff.insert(data.insert.RDPFile)
				case INS_RDP_SCREENSIZE: return true
				case INS_RDP_DYNAMIC:
					ui.insert_sel_ok = false
					if data.insert.Dynamic == true {
						data.insert.Dynamic = false
					} else {
						data.insert.Dynamic = true
					}
					return true
				case INS_RDP_FULLSCR:
					ui.insert_sel_ok = false
					if data.insert.FullScr == true {
						data.insert.FullScr = false
					} else {
						data.insert.FullScr = true
					}
					return true
				case INS_RDP_MULTIMON:
					ui.insert_sel_ok = false
					if data.insert.MultiMon == true {
						data.insert.MultiMon = false
					} else {
						data.insert.MultiMon = true
					}
					return true
				case INS_RDP_QUALITY: return true
				case INS_RDP_DRIVE + len(data.insert.Drive): return true
				case INS_CMD_CMD: ui.buff.insert(data.insert.Host)
				case INS_CMD_SHELL: ui.buff.insert(data.insert.Shell[0])
				case INS_CMD_SILENT:
					ui.insert_sel_ok = false
					if data.insert.Silent == true {
						data.insert.Silent = false
					} else {
						data.insert.Silent = true
					}
					return true
				case INS_OS_USERDOMAINID:
					ui.buff.insert(data.insert.Stack.UserDomainID)
				case INS_OS_PROJECTID:
					ui.buff.insert(data.insert.Stack.ProjectID)
				case INS_OS_REGION:
					ui.buff.insert(data.insert.Stack.RegionName)
				case INS_OS_ENDTYPE:
					ui.buff.insert(data.insert.Stack.EndpointType)
				case INS_OS_INTERFACE:
					ui.buff.insert(data.insert.Stack.Interface)
				case INS_OS_IDAPI:
					ui.buff.insert(data.insert.Stack.IdentityAPI)
				case INS_OS_IMGAPI:
					ui.buff.insert(data.insert.Stack.ImageAPI)
				case INS_OS_NETAPI:
					ui.buff.insert(data.insert.Stack.NetworkAPI)
				case INS_OS_VOLAPI:
					ui.buff.insert(data.insert.Stack.VolumeAPI)
				case INS_SSH_NOTE,
					 INS_RDP_NOTE + len(data.insert.Drive),
					 INS_CMD_NOTE,
					 INS_OS_NOTE:
					ui.buff.insert(data.insert.Note)
				}
			}
		} else {
			if event.Key() == tcell.KeyEscape ||
			   event.Key() == tcell.KeyCtrlC {
				ui.insert_sel_ok = false
				ui.buff.empty()
				ui.drives_buff = ""
				ui.s.HideCursor()
			}
			if len(data.insert.Drive) > 0 &&
			   (ui.insert_sel >= INS_RDP_DRIVE &&
			   ui.insert_sel < INS_RDP_DRIVE +
			   len(data.insert.Drive)) {
				if event.Rune() == 'y' ||
				event.Rune() == 'Y' ||
				event.Key() == tcell.KeyEnter {
					delete(data.insert.Drive,
						   data.insert.drive_keys[
						   ui.insert_sel - INS_RDP_DRIVE])
					if len(data.insert.Drive) == 0 {
						data.insert.Drive = nil
					}
					e_set_drive_keys(data)
				}
				ui.insert_sel_ok = false
				return true
			}
			switch ui.insert_sel {
			case INS_PROTOCOL:
				if event.Rune() < '1' || event.Rune() > '4' {
					ui.insert_sel_ok = false
					ui.buff.empty()
					ui.s.HideCursor()
					return true
				} else {
					filename := data.insert.filename
					name := data.insert.Name
					parent := data.insert.parent
					data.insert = nil
					data.insert = &HostNode{}
					data.insert.Name = name
					data.insert.parent = parent
					data.insert.filename = filename
					data.insert.Protocol = int8(event.Rune() - 48 - 1)
					ui.insert_sel_ok = false
					ui.s.HideCursor()
					e_set_protocol_defaults(data, data.insert)
				}
			case INS_RDP_SCREENSIZE:
				if event.Rune() < '1' || event.Rune() > '7' {
					ui.insert_sel_ok = false
					ui.buff.empty()
					ui.s.HideCursor()
					return true
				} else {
					s := strings.Split(
						RDP_SCREENSIZE[uint8(event.Rune() - 48 - 1)],
						"x")
					if len(s) != 2 {
						return true
					}
					tmp, _ := strconv.Atoi(s[W])
					data.insert.Width = uint16(tmp)
					tmp, _ = strconv.Atoi(s[H])
					data.insert.Height = uint16(tmp)
					ui.insert_sel_ok = false
					ui.s.HideCursor()
				}
			case INS_RDP_QUALITY:
				if event.Rune() < '1' || event.Rune() > '3' {
					ui.insert_sel_ok = false
					ui.buff.empty()
					ui.s.HideCursor()
					return true
				} else {
					data.insert.Quality = uint8(event.Rune() - 48 - 1)
					ui.insert_sel_ok = false
					ui.s.HideCursor()
				}
			case INS_RDP_DRIVE + len(data.insert.Drive):
				if len(ui.drives_buff) == 0 {
					if event.Key() == tcell.KeyEnter {
						if ui.buff.len() == 0 {
							ui.insert_sel_ok = false
							ui.drives_buff = ""
							ui.buff.empty()
							ui.s.HideCursor()
							return true
						}
						ui.drives_buff = ui.buff.str()
						ui.buff.empty()
					} else {
						e_readline(event, &ui.buff)
					}
				} else {
					if event.Key() == tcell.KeyEnter {
						if ui.buff.len() == 0 {
							ui.insert_sel_ok = false
							ui.drives_buff = ""
							ui.buff.empty()
							ui.s.HideCursor()
							return true
						}
						if len(data.insert.Drive) == 0 {
							data.insert.Drive = make(map[string]string)
						}
						data.insert.Drive[ui.drives_buff] = ui.buff.str()
						e_set_drive_keys(data)
						ui.insert_sel_ok = false
						ui.drives_buff = ""
						ui.buff.empty()
						ui.s.HideCursor()
					} else {
						e_readline(event, &ui.buff)
					}
				}
			case INS_SSH_HOST,
				 INS_SSH_PORT,
				 INS_SSH_USER,
				 INS_SSH_PASS,
				 INS_SSH_PRIV,
				 INS_SSH_EXEC,
				 INS_SSH_JUMP_HOST,
				 INS_SSH_JUMP_PORT,
				 INS_SSH_JUMP_USER,
				 INS_SSH_JUMP_PASS,
				 INS_SSH_JUMP_PRIV,
				 INS_RDP_JUMP_HOST + len(data.insert.Drive),
				 INS_RDP_JUMP_PORT + len(data.insert.Drive),
				 INS_RDP_JUMP_USER + len(data.insert.Drive),
				 INS_RDP_JUMP_PASS + len(data.insert.Drive),
				 INS_RDP_JUMP_PRIV + len(data.insert.Drive),
				 INS_SSH_NOTE,
				 INS_RDP_HOST,
				 INS_RDP_PORT,
				 INS_RDP_DOMAIN,
				 INS_RDP_USER,
				 INS_RDP_PASS,
				 INS_RDP_FILE,
				 INS_RDP_NOTE + len(data.insert.Drive),
				 INS_CMD_CMD,
				 INS_CMD_SHELL,
				 INS_CMD_NOTE,
				 INS_OS_HOST,
				 INS_OS_USER,
				 INS_OS_PASS,
				 INS_OS_USERDOMAINID,
				 INS_OS_PROJECTID,
				 INS_OS_REGION,
				 INS_OS_ENDTYPE,
				 INS_OS_INTERFACE,
				 INS_OS_IDAPI,
				 INS_OS_IMGAPI,
				 INS_OS_NETAPI,
				 INS_OS_VOLAPI,
				 INS_OS_NOTE:
				if event.Key() == tcell.KeyEnter {
					switch ui.insert_sel {
					case INS_SSH_HOST,
						 INS_RDP_HOST,
						 INS_OS_HOST:
						data.insert.Host = ui.buff.str()
					case INS_SSH_PORT,
						 INS_RDP_PORT:
						tmp, _ := strconv.Atoi(ui.buff.str())
						data.insert.Port = uint16(tmp)
					case INS_SSH_USER,
						 INS_RDP_USER,
						 INS_OS_USER:
						data.insert.User = ui.buff.str()
					case INS_SSH_PASS,
						 INS_RDP_PASS,
						 INS_OS_PASS:
						if ui.buff.len() == 0 {
							data.insert.Pass = ""
							return true
						} else {
							data.insert.Pass, _ = c_encrypt_str(ui.buff.str(),
													data.opts.GPG)
						}
					case INS_SSH_PRIV: data.insert.Priv = ui.buff.str()
					case INS_SSH_EXEC: data.insert.Exec = ui.buff.str()
					case INS_SSH_JUMP_HOST,
						 INS_RDP_JUMP_HOST + len(data.insert.Drive):
						data.insert.Jump.Host = ui.buff.str()
						if len(ui.buff.str()) > 0 {
							data.insert.Jump.Port = 22
						} else {
							data.insert.Jump.Port = 0
						}
					case INS_SSH_JUMP_PORT,
						 INS_RDP_JUMP_PORT + len(data.insert.Drive):
						tmp, _ := strconv.Atoi(ui.buff.str())
						data.insert.Jump.Port = uint16(tmp)
					case INS_SSH_JUMP_USER,
						 INS_RDP_JUMP_USER + len(data.insert.Drive):
						data.insert.Jump.User = ui.buff.str()
					case INS_SSH_JUMP_PASS,
						 INS_RDP_JUMP_PASS + len(data.insert.Drive):
						if len(ui.buff.str()) == 0 {
							data.insert.Jump.Pass = ""
						} else {
							data.insert.Jump.Pass, _ =
							c_encrypt_str(ui.buff.str(), data.opts.GPG)
						}
					case INS_SSH_JUMP_PRIV,
						 INS_RDP_JUMP_PRIV + len(data.insert.Drive):
						data.insert.Jump.Priv = ui.buff.str()
					case INS_RDP_DOMAIN:
						data.insert.Domain = ui.buff.str()
					case INS_RDP_FILE:
						data.insert.RDPFile = ui.buff.str()
					case INS_CMD_CMD:
						data.insert.Host = ui.buff.str()
					case INS_CMD_SHELL:
						data.insert.Shell[0] = ui.buff.str()
					case INS_OS_USERDOMAINID:
						data.insert.Stack.UserDomainID = ui.buff.str()
					case INS_OS_PROJECTID:
						data.insert.Stack.ProjectID = ui.buff.str()
					case INS_OS_REGION:
						data.insert.Stack.RegionName = ui.buff.str()
					case INS_OS_ENDTYPE:
						data.insert.Stack.EndpointType = ui.buff.str()
					case INS_OS_INTERFACE:
						data.insert.Stack.Interface = ui.buff.str()
					case INS_OS_IDAPI:
						data.insert.Stack.IdentityAPI = ui.buff.str()
					case INS_OS_IMGAPI:
						data.insert.Stack.ImageAPI = ui.buff.str()
					case INS_OS_NETAPI:
						data.insert.Stack.NetworkAPI = ui.buff.str()
					case INS_OS_VOLAPI:
						data.insert.Stack.VolumeAPI = ui.buff.str()
					case INS_SSH_NOTE,
						 INS_RDP_NOTE + len(data.insert.Drive),
						 INS_CMD_NOTE,
						 INS_OS_NOTE:
						data.insert.Note = ui.buff.str()
					}
					ui.insert_sel_ok = false
					ui.buff.empty()
					ui.s.HideCursor()
				} else {
					e_readline(event, &ui.buff)
				}
			}
		}
	}
	return false
}

func e_rename_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC {
		data.insert = nil
	} else if event.Key() == tcell.KeyEnter {
		if err := e_rename(data, ui); err != nil {
			ui.s.HideCursor()
			ui.buff.empty()
			return true
		}
	} else {
		e_readline(event, &ui.buff)
		return true
	}
	ui.s.HideCursor()
	ui.mode = NORMAL_MODE
	ui.buff.empty()
	return false
}

func e_help_events(data *HardData, ui *HardUI, event tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC ||
	   event.Rune() == 'q' ||
	   event.Rune() == '?' {
		ui.mode = NORMAL_MODE
		ui.help_scroll = 0
		return true
	} else if event.Rune() == 'j' ||
			  event.Key() == tcell.KeyDown {
		if ui.help_end == true {
			return true
		}
		ui.help_scroll += 1
	} else if event.Rune() == 'k' ||
			  event.Key() == tcell.KeyUp {
		if ui.help_scroll <= 0 {
			ui.help_scroll = 0
			return true
		}
		ui.help_scroll -= 1
	} else if event.Rune() == 'g' {
		ui.help_scroll = 0
	} else if event.Rune() == 'G' {
		for ui.help_end != true {
			ui.help_scroll += 1
			i_draw_help(ui)
		}
	}
	return false
}
