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
 * Wed Jan 31 16:40:52 2024
 * Joe
 *
 * constants
 */

package main

const (
	CONF_FILE_NAME = "config.yml"
	CONF_DIR_NAME  = "hf"
	DATA_DIR_NAME  = "hf"
	VERSION        = "v0.4"
)

const (
	NORMAL_KEYS_HINTS = `!a/i: insert host -
!m: mkdir -
!s: search -
[C-r]: reload
!?: help`
	DELETE_KEYS_HINTS = `q: quit -
y: yes -
n: no`
	ERROR_KEYS_HINTS = "[Enter] Ok"
)

const (
	NORMAL_MODE = 0
	DELETE_MODE = 1
	LOAD_MODE   = 2
	ERROR_MODE  = 3
)

const (
	W           = 0
	H           = 1
	ERROR_MSG   = 0
	ERROR_ERR   = 1
	STYLE_DEF   = 0
	STYLE_DIR   = 1
	STYLE_BOX   = 2
	STYLE_HEAD  = 3
	STYLE_ERR   = 4
	STYLE_TITLE = 5
	STYLE_BOT   = 6
)

var (
	HOST_ICONS = [2]string{" ", " "}
	DIRS_ICONS = [2]string{" ", " "}
)

var DEFAULT_OPTS = HardOpts{
	true,
	true,
	"",
	false,
	"",
}

