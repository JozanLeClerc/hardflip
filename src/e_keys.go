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
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func e_normal_events(data *HardData, event tcell.EventKey) bool {
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
			return true
		}
		data.yank = data.litems.curr
	}
	return false
}

func e_delete_events(data *HardData, event tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC ||
	   event.Rune() == 'n' {
		data.ui.mode = NORMAL_MODE
	} else if event.Key() == tcell.KeyEnter ||
			  event.Rune() == 'y' {
		if err := e_delete_host(data); err == nil {
			data.ui.mode = NORMAL_MODE
			return true
		}
	}
	return false
}

func e_load_events(data *HardData, event tcell.EventKey) bool {
	return true
}

func e_error_events(data *HardData, event tcell.EventKey) bool {
	if event.Rune() != 0 ||
	   event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyEnter {
		data.ui.mode = NORMAL_MODE
		data.load_err = nil
	}
	return false
}

func e_welcome_events(data *HardData, event tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC {
		data.ui.s.Fini()
		os.Exit(0)
	}
	if len(data.opts.GPG) == 0 {
		if event.Rune() < '1' || event.Rune() > '9' {
			return true
		} else {
			data.opts.GPG = data.keys[event.Rune() - 48 - 1][0]
			data.ui.s.HideCursor()
		}
	} else {
		if event.Rune() == 'y' {
			data.ui.mode = NORMAL_MODE
			c_write_options(data.opts.file, data.opts, &data.load_err)
		} else if event.Rune() == 'n' {
			data.opts.GPG = ""
		}
	}
	return false
}

func e_insert_events(data *HardData, event tcell.EventKey) bool {
	ui := &data.ui

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
				return true
			}
			ui.s.HideCursor()
			data.insert = &HostNode{}
			e_set_protocol_defaults(data, data.insert)
			data.insert.Name = ui.buff
			ui.buff = ""
			if data.litems.curr != nil {
				data.insert.parent = data.litems.curr.path_node()
			} else {
				data.insert.parent = data.ldirs.head
			}
		} else {
			e_readline(event, &data.ui.buff)
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
				} else if data.insert.Protocol == PROTOCOL_CMD &&
						  data.ui.insert_sel == INS_PROTOCOL {
					data.ui.insert_sel = INS_CMD_CMD
				} else if data.insert.Protocol == PROTOCOL_OS &&
						  data.ui.insert_sel == INS_PROTOCOL {
					data.ui.insert_sel = INS_OS_HOST
				} else if data.insert.Protocol == PROTOCOL_SSH &&
						  data.ui.insert_sel == INS_SSH_JUMP_HOST &&
						  len(data.insert.Jump.Host) == 0 {
					data.ui.insert_sel = INS_SSH_NOTE
				} else if data.ui.insert_sel < data.ui.insert_sel_max {
					data.ui.insert_sel += 1
				}
			} else if event.Rune() == 'k' ||
					  event.Key() == tcell.KeyUp {
				if data.insert.Protocol == PROTOCOL_RDP &&
				   data.ui.insert_sel == INS_RDP_HOST {
					data.ui.insert_sel = INS_PROTOCOL
				} else if data.insert.Protocol == PROTOCOL_CMD &&
						  data.ui.insert_sel == INS_CMD_CMD {
					data.ui.insert_sel = INS_PROTOCOL
				} else if data.insert.Protocol == PROTOCOL_OS &&
						  data.ui.insert_sel == INS_OS_HOST {
					data.ui.insert_sel = INS_PROTOCOL
				} else if data.insert.Protocol == PROTOCOL_SSH &&
						  data.ui.insert_sel == INS_SSH_NOTE &&
						  len(data.insert.Jump.Host) == 0 {
					data.ui.insert_sel = INS_SSH_JUMP_HOST
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
					  event.Rune() == ' ' ||
					  event.Key() == tcell.KeyEnter {
				data.ui.insert_sel_ok = true
				switch data.ui.insert_sel {
				case INS_SSH_OK,
					 INS_RDP_OK + len(data.insert.Drive),
					 INS_CMD_OK,
					 INS_OS_OK:
					data.ui.insert_sel_ok = false
					i_insert_check_ok(data, data.insert)
					if data.insert_err != nil {
						return true
					}
					i_insert_host(data, data.insert)
				case INS_SSH_HOST,
					 INS_RDP_HOST,
					 INS_OS_HOST:
					ui.buff = data.insert.Host
				case INS_SSH_PORT,
					 INS_RDP_PORT:
					if data.insert.Port > 0 {
						ui.buff = strconv.Itoa(int(data.insert.Port))
					}
				case INS_SSH_USER,
					 INS_RDP_USER,
					 INS_OS_USER:
					ui.buff = data.insert.User
				case INS_SSH_PASS,
					 INS_RDP_PASS,
					 INS_OS_PASS:
					return true
				case INS_SSH_PRIV: ui.buff = data.insert.Priv
				case INS_SSH_JUMP_HOST: ui.buff = data.insert.Jump.Host
				case INS_SSH_JUMP_PORT:
					if data.insert.Jump.Port > 0 {
						ui.buff = strconv.Itoa(int(
						data.insert.Jump.Port))
					}
				case INS_SSH_JUMP_USER: ui.buff = data.insert.Jump.User
				case INS_SSH_JUMP_PASS: return true
				case INS_SSH_JUMP_PRIV: ui.buff = data.insert.Jump.Priv
				case INS_RDP_DOMAIN: ui.buff = data.insert.Domain
				case INS_RDP_FILE: ui.buff = data.insert.RDPFile
				case INS_RDP_SCREENSIZE: return true
				case INS_RDP_DYNAMIC:
					data.ui.insert_sel_ok = false
					if data.insert.Dynamic == true {
						data.insert.Dynamic = false
					} else {
						data.insert.Dynamic = true
					}
					return true
				case INS_RDP_QUALITY: return true
				case INS_RDP_DRIVE + len(data.insert.Drive): return true
				case INS_CMD_CMD: ui.buff = data.insert.Host
				case INS_CMD_SHELL: ui.buff = data.insert.Shell[0]
				case INS_CMD_SILENT:
					data.ui.insert_sel_ok = false
					if data.insert.Silent == true {
						data.insert.Silent = false
					} else {
						data.insert.Silent = true
					}
					return true
				case INS_OS_USERDOMAINID:
					ui.buff = data.insert.Stack.UserDomainID
				case INS_OS_PROJECTID:
					ui.buff = data.insert.Stack.ProjectID
				case INS_OS_REGION:
					ui.buff = data.insert.Stack.RegionName
				case INS_OS_ENDTYPE:
					ui.buff = data.insert.Stack.EndpointType
				case INS_OS_INTERFACE:
					ui.buff = data.insert.Stack.Interface
				case INS_OS_IDAPI:
					ui.buff = data.insert.Stack.IdentityAPI
				case INS_OS_IMGAPI:
					ui.buff = data.insert.Stack.ImageAPI
				case INS_OS_NETAPI:
					ui.buff = data.insert.Stack.NetworkAPI
				case INS_OS_VOLAPI:
					ui.buff = data.insert.Stack.VolumeAPI
				case INS_SSH_NOTE,
					 INS_RDP_NOTE + len(data.insert.Drive),
					 INS_CMD_NOTE,
					 INS_OS_NOTE:
					ui.buff = data.insert.Note
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
			if len(data.insert.Drive) > 0 &&
			   (data.ui.insert_sel >= INS_RDP_DRIVE &&
			   data.ui.insert_sel < INS_RDP_DRIVE +
			   len(data.insert.Drive)) {
				if event.Rune() == 'y' ||
				event.Rune() == 'Y' ||
				event.Key() == tcell.KeyEnter {
					delete(data.insert.Drive,
						   data.insert.drive_keys[
						   data.ui.insert_sel - INS_RDP_DRIVE])
					if len(data.insert.Drive) == 0 {
						data.insert.Drive = nil
					}
					e_set_drive_keys(data)
				}
				data.ui.insert_sel_ok = false
				return true
			}
			switch data.ui.insert_sel {
			case INS_PROTOCOL:
				if event.Rune() < '1' || event.Rune() > '4' {
					data.ui.insert_sel_ok = false
					ui.buff = ""
					ui.s.HideCursor()
					return true
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
					e_set_protocol_defaults(data, data.insert)
				}
			case INS_RDP_SCREENSIZE:
				if event.Rune() < '1' || event.Rune() > '7' {
					data.ui.insert_sel_ok = false
					ui.buff = ""
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
					data.ui.insert_sel_ok = false
					ui.s.HideCursor()
				}
			case INS_RDP_QUALITY:
				if event.Rune() < '1' || event.Rune() > '3' {
					data.ui.insert_sel_ok = false
					ui.buff = ""
					ui.s.HideCursor()
					return true
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
							return true
						}
						data.ui.drives_buff = ui.buff
						ui.buff = ""
					} else {
						e_readline(event, &data.ui.buff)
					}
				} else {
					if event.Key() == tcell.KeyEnter {
						if len(ui.buff) == 0 {
							data.ui.insert_sel_ok = false
							data.ui.drives_buff = ""
							ui.buff = ""
							ui.s.HideCursor()
							return true
						}
						if len(data.insert.Drive) == 0 {
							data.insert.Drive = make(map[string]string)
						}
						data.insert.Drive[ui.drives_buff] = ui.buff
						e_set_drive_keys(data)
						data.ui.insert_sel_ok = false
						ui.drives_buff = ""
						ui.buff = ""
						ui.s.HideCursor()
					} else {
						e_readline(event, &data.ui.buff)
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
					switch data.ui.insert_sel {
					case INS_SSH_HOST,
						 INS_RDP_HOST,
						 INS_OS_HOST:
						data.insert.Host = ui.buff
					case INS_SSH_PORT,
						 INS_RDP_PORT:
						tmp, _ := strconv.Atoi(ui.buff)
						data.insert.Port = uint16(tmp)
					case INS_SSH_USER,
						 INS_RDP_USER,
						 INS_OS_USER:
						data.insert.User = ui.buff
					case INS_SSH_PASS,
						 INS_RDP_PASS,
						 INS_OS_PASS:
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
					case INS_CMD_CMD:
						data.insert.Host = ui.buff
					case INS_CMD_SHELL:
						data.insert.Shell[0] = ui.buff
					case INS_OS_USERDOMAINID:
						data.insert.Stack.UserDomainID = ui.buff
					case INS_OS_PROJECTID:
						data.insert.Stack.ProjectID = ui.buff
					case INS_OS_REGION:
						data.insert.Stack.RegionName = ui.buff
					case INS_OS_ENDTYPE:
						data.insert.Stack.EndpointType = ui.buff
					case INS_OS_INTERFACE:
						data.insert.Stack.Interface = ui.buff
					case INS_OS_IDAPI:
						data.insert.Stack.IdentityAPI = ui.buff
					case INS_OS_IMGAPI:
						data.insert.Stack.ImageAPI = ui.buff
					case INS_OS_NETAPI:
						data.insert.Stack.NetworkAPI = ui.buff
					case INS_OS_VOLAPI:
						data.insert.Stack.VolumeAPI = ui.buff
					case INS_SSH_NOTE,
						 INS_RDP_NOTE + len(data.insert.Drive),
						 INS_CMD_NOTE,
						 INS_OS_NOTE:
						data.insert.Note = ui.buff
					}
					data.ui.insert_sel_ok = false
					ui.buff = ""
					ui.s.HideCursor()
				} else {
					e_readline(event, &data.ui.buff)
				}
			}
		}
	}
	return false
}

func e_mkdir_events(data *HardData, event tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape ||
	   event.Key() == tcell.KeyCtrlC {
		data.ui.s.HideCursor()
		data.ui.mode = NORMAL_MODE
		data.ui.buff = ""
		data.insert = nil
	} else if event.Key() == tcell.KeyEnter {
		e_mkdir(data, &data.ui)
		data.ui.s.HideCursor()
		data.ui.mode = NORMAL_MODE
		data.ui.buff = ""
	} else {
		e_readline(event, &data.ui.buff)
	}
	return false
}