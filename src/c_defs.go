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
 * hardflip: src/c_defs.go
 * Thu Apr 11 17:55:54 2024
 * Joe
 *
 * constants
 */

package main

const (
	CONF_FILE_NAME = "config.yml"
	CONF_DIR_NAME  = "hf"
	DATA_DIR_NAME  = "hf"
	VERSION        = "v0.5"
	VERSION_NAME   = ""
)

const (
	NORMAL_KEYS_HINTS = `!a/i: insert host -
m: mkdir -
!s: search -
[C-r]: reload
!?: help`
	ERROR_KEYS_HINTS = "[Enter]: ok"
	CONFIRM_KEYS_HINTS = `y/n: yes - no`
	INSERT_KEYS_HINTS = ``
)

const (
	NORMAL_MODE = iota
	DELETE_MODE
	LOAD_MODE
	ERROR_MODE
	WELCOME_MODE
	MKDIR_MODE
	INSERT_MODE
)

const (
	W = 0
	H = 1
)

const (
	ERROR_MSG   = 0
	ERROR_ERR   = 1
)

const (
	DEF_STYLE = iota
	DIR_STYLE
	BOX_STYLE
	HEAD_STYLE
	ERR_STYLE
	TITLE_STYLE
	BOT_STYLE
)

const (
	PROTOCOL_SSH = iota
	PROTOCOL_RDP
	PROTOCOL_CMD
	PROTOCOL_OS
	PROTOCOL_MAX = PROTOCOL_OS
)

const (
	INS_PROTOCOL = iota
	INS_SSH_HOST
	INS_SSH_PORT
	INS_SSH_USER
	INS_SSH_PASS
	INS_SSH_PRIV
	INS_SSH_JUMP_HOST
	INS_SSH_JUMP_PORT
	INS_SSH_JUMP_USER
	INS_SSH_JUMP_PASS
	INS_SSH_JUMP_PRIV
	INS_SSH_NOTE
	INS_SSH_OK
	INS_RDP_HOST
	INS_RDP_PORT
	INS_RDP_DOMAIN
	INS_RDP_USER
	INS_RDP_PASS
	INS_RDP_FILE
	INS_RDP_SCREENSIZE
	INS_RDP_DYNAMIC
	INS_RDP_QUALITY
	INS_RDP_DRIVE
	INS_RDP_NOTE
	INS_RDP_OK
	INS_CMD_CMD
	INS_CMD_SHELL
	INS_CMD_SILENT
	INS_CMD_NOTE
	INS_CMD_OK
	INS_OS_HOST
	INS_OS_USER
	INS_OS_PASS
	INS_OS_USERDOMAINID
	INS_OS_PROJECTID
	INS_OS_REGION
	INS_OS_ENDTYPE
	INS_OS_INTERFACE
	INS_OS_IDAPI
	// TODO: here
	INS_OS_OK
)

var (
	HOST_ICONS = [4]string{" ", " ", " ", "󰅟 "}
	DIRS_ICONS = [2]string{" ", " "}
	RDP_SCREENSIZE = [7]string{
		"800x600",
		"1280x720",
		"1360x768",
		"1600x900",
		"1600x1200",
		"1920x1080",
		"2560x1440",
	}
	RDP_QUALITY = [3]string{"Low", "Medium", "High"}
	PROTOCOL_STR = [PROTOCOL_MAX + 1]string{
		"SSH",
		"RDP",
		"Single command",
		"OpenStack CLI",
	}
)

var DEFAULT_OPTS = HardOpts{
	true,
	true,
	"",
	false,
	"",
	"",
}

