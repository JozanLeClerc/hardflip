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
 * hardflip: src/c_utils.go
 * Mon Jan 29 08:56:55 2024
 * Joe
 *
 * core funcs
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// this function will go get the data folder and try to create it if it does
// not exist
// the first path being checked is $XDG_CONFIG_HOME then $HOME/.config
// it returns the full data directory path
func c_get_conf_dir(ui *HardUI) string {
	var ptr string
	var home string

	if home = os.Getenv("HOME"); len(home) == 0 {
		if ui == nil {
			c_die("env variable HOME not defined", nil)
		}
		c_error_mode("env variable HOME not defined", nil, ui)
		return ""
	}
	xdg_home := os.Getenv("XDG_CONFIG_HOME")

	if len(xdg_home) > 0 {
		ptr = xdg_home
	} else {
		ptr = home + "/.config"
	}
	ptr += "/" + CONF_DIR_NAME
	if _, err := os.Stat(ptr); os.IsNotExist(err) {
	    if err := os.MkdirAll(ptr, os.ModePerm); err != nil {
			if ui == nil {
				c_die("could not create path " + ptr, err)
			}
			c_error_mode("could not create path" + ptr, err, ui)
	    }
	} else if err != nil {
		if ui == nil {
			c_die("could read path " + ptr, err)
		}
		c_error_mode("could read path" + ptr, err, ui)
		return ""
	}
	return ptr
}

// this function will go get the data folder and try to create it if it does
// not exist
// the first path being checked is $XDG_DATA_HOME then $HOME/.local/share
// it returns the full data directory path
func c_get_data_dir(ui *HardUI) string {
	var ptr string
	var home string

	if home = os.Getenv("HOME"); len(home) == 0 {
		if ui == nil {
			c_die("env variable HOME not defined", nil)
		}
		c_error_mode("env variable HOME not defined", nil, ui)
		return ""
	}
	xdg_home := os.Getenv("XDG_DATA_HOME")

	if len(xdg_home) > 0 {
		ptr = xdg_home
	} else {
		ptr = home + "/.local/share"
	}
	ptr += "/" + DATA_DIR_NAME
	if _, err := os.Stat(ptr); os.IsNotExist(err) {
	    if err := os.MkdirAll(ptr, os.ModePerm); err != nil {
			if ui == nil {
				c_die("could not create path " + ptr, err)
			}
			c_error_mode("could not create path" + ptr, err, ui)
	    }
	} else if err != nil {
		if ui == nil {
			c_die("could read path " + ptr, err)
		}
		c_error_mode("could read path" + ptr, err, ui)
		return ""
	}
	return ptr
}

// c_die displays an error string to the stderr fd and exits the program
// with the return code 1
// it takes an optional err argument of the error type as a complement of
// information
func c_die(str string, err error) {
	fmt.Fprintf(os.Stderr, "error: %s", str)
	if err != nil {
		fmt.Fprintf(os.Stderr, ": %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func c_error_mode(msg string, err error, ui *HardUI) {
	ui.mode = ERROR_MODE
	err_str := ""
	if err != nil {
		err_str = fmt.Sprintf("%v", err)
	}
	ui.err[ERROR_MSG] = msg
	ui.err[ERROR_ERR] = err_str
}

// c_encrypt_str encrypts a string with the given gpgkey
func c_encrypt_str(str string, gpg string) (string, error) {
	cmd := exec.Command("gpg", "-r", gpg, "-a", "-e")
	cmd.Stdin = strings.NewReader(str)
	out, err := cmd.Output()
	return string(out), err
}

// c_decrypt_str will try to decrypt the given str
func c_decrypt_str(str string) (string, error) {
	cmd := exec.Command("gpg", "-q", "-d")
	cmd.Stdin = strings.NewReader(str)
	out, err := cmd.Output()
	return string(out), err
}
