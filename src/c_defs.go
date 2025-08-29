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
 * Mon May 13 10:27:54 2024
 * Joe
 *
 * constants
 */

package main

const (
	CONF_FILE_NAME = "config.yml"
	STYLE_FILE_NAME = "colors.yml"
	CONF_DIR_NAME  = "hf"
	DATA_DIR_NAME  = "hf"
	VERSION        = "v1.0.3"
	VERSION_NAME   = "wheelbite"
)

const (
	NORMAL_KEYS_HINTS = `add host: a/i: |
mkdir: m |
reload: <c-r> |
search: / |
keybinds: ?`
	ERROR_KEYS_HINTS = "ok: [Enter]"
	CONFIRM_KEYS_HINTS = `yes-no: y/n`
	HELP_KEYS_HINTS = `close: q |
down: j |
up: k`
)

const (
	NORMAL_MODE = iota
	DELETE_MODE
	ERROR_MODE
	WELCOME_MODE
	MKDIR_MODE
	INSERT_MODE
	RENAME_MODE
	HELP_MODE
	MODE_MAX = HELP_MODE
)

const (
	WELCOME_GPG = iota
	WELCOME_CONFIRM_GPG
	WELCOME_SSH
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
	YANK_STYLE
	MOVE_STYLE
	STYLE_MAX = MOVE_STYLE
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
	INS_SSH_EXEC
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
	INS_RDP_FULLSCR
	INS_RDP_MULTIMON
	INS_RDP_QUALITY
	INS_RDP_DRIVE
	INS_RDP_JUMP_HOST
	INS_RDP_JUMP_PORT
	INS_RDP_JUMP_USER
	INS_RDP_JUMP_PASS
	INS_RDP_JUMP_PRIV
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
	INS_OS_IMGAPI
	INS_OS_NETAPI
	INS_OS_VOLAPI
	INS_OS_NOTE
	INS_OS_OK
	INS_MAX = INS_OS_OK
)

const (
	INSERT_ADD = iota
	INSERT_COPY
	INSERT_MOVE
	INSERT_EDIT
)

var HELP_NORMAL_KEYS = [][2]string{
	{"<down> | j",				"go to next item"},
	{"<up> | k",				"go to previous item"},
	{"<pgdown> | <c-d>",		"jump down"},
	{"<pgup> | <c-u>",	 		"jump up"},
	{"} | ]",					"go to next directory"},
	{"{ | [",					"go to previous directory"},
	{"g",						"go to first item"},
	{"G",						"go to last item"},
	{"D",						"delete selected item"},
	{"h",						"close current directory"},
	{"H",						"close all directories"},
	{"l | <space>",				"toggle directory"},
	{"<right> | l | <enter>",	"connect to selected host / toggle directory"},
	{"a | i",					"create a new host"},
	{"y",						"copy selected host"},
	{"d",						"cut selected host"},
	{"p",						"paste host in clipboard"},
	{"<f7> | m",				"create a new directory"},
	{"e",						"edit selected host"},
	{"c | C | A",				"rename selected item"},
	{"<c-r>",					"reload data and configuration"},
	{"/ | f | <c-f>",			"fuzzy search"},
	{"?",						"display this help menu"},
	{"<c-c> | q",				"quit"},
}

var (
	HOST_ICONS = [4]string{" ", " ", " ", " "}
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
	"",
}

var DEFAULT_STYLE = HardStyle{
	DefColor:	COLORS[COLOR_DEFAULT],
	DirColor:	COLORS[COLOR_BOLD_BLUE],
	BoxColor:	COLORS[COLOR_DEFAULT],
	HeadColor:	COLORS[COLOR_DEFAULT],
	ErrColor:	COLORS[COLOR_RED],
	TitleColor:	COLORS[COLOR_BOLD_BLUE],
	BotColor:	COLORS[COLOR_BLUE],
	YankColor:	COLORS[COLOR_BOLD_YELLOW],
	MoveColor:	COLORS[COLOR_BOLD_RED],
}

const (
	COLOR_DEFAULT = iota
    COLOR_BLACK
    COLOR_RED
    COLOR_GREEN
    COLOR_YELLOW
    COLOR_BLUE
    COLOR_MAGENTA
    COLOR_CYAN
    COLOR_WHITE
    COLOR_GRAY
    COLOR_BOLD_BLACK
    COLOR_BOLD_RED
    COLOR_BOLD_GREEN
    COLOR_BOLD_YELLOW
    COLOR_BOLD_BLUE
    COLOR_BOLD_MAGENTA
    COLOR_BOLD_CYAN
    COLOR_BOLD_WHITE
    COLOR_BOLD_GRAY
	COLORS_MAX = COLOR_BOLD_GRAY
)

var COLORS = [COLORS_MAX + 1]string{
	"default",
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
	"gray",
	"boldblack",
	"boldred",
	"boldgreen",
	"boldyellow",
	"boldblue",
	"boldmagenta",
	"boldcyan",
	"boldwhite",
	"boldgray",
}
