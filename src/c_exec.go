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
 * hardflip: src/c_exec.go
 * Tue, 26 Aug 2025 19:17:04 +0200
 * Joe
 *
 * exec the command at some point
 */

package main

import (
	"bytes"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func c_exec_cmd(cmd_fmt, cmd_env []string, silent bool) (error, string) {
	var errb bytes.Buffer
	cmd := exec.Command(cmd_fmt[0], cmd_fmt[1:]...)

	if cmd_env != nil {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, cmd_env...)
	}
	if silent == false {
		cmd.Stdin  = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = &errb
	}
	if err := cmd.Run(); err != nil {
		return err, errb.String()
	}
	return nil, ""
}

func c_format_ssh_jump(host *HostNode) string {
	jump_fmt := "-oProxyCommand=ssh"
	if len(host.Jump.Priv) > 0 {
		jump_fmt += " -i " + host.Jump.Priv
	}
	if host.Jump.Port != 0 {
		jump_fmt += " -p " + strconv.Itoa(int(host.Jump.Port))
	}
	if len(host.Jump.User) == 0 {
		jump_fmt += " root"
	} else {
		jump_fmt += " " + host.Jump.User
	}
	jump_fmt += "@" + host.Jump.Host + " -W %h:%p"
	return jump_fmt
}

func c_format_ssh(host *HostNode, pass string) ([]string, []string) {
	cmd_fmt := []string{}
	if len(pass) > 0 {
		cmd_fmt = append(cmd_fmt, "sshpass", "-p", pass)
	}

	cmd_fmt = append(cmd_fmt, "ssh", "-F", "none",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=ERROR")
	if len(host.Priv) > 0 {
		cmd_fmt = append(cmd_fmt, "-i", host.Priv)
	}
	if len(host.Jump.Host) > 0 {
		cmd_fmt = append(cmd_fmt, c_format_ssh_jump(host))
	}
	if host.Port != 0 {
		cmd_fmt = append(cmd_fmt, "-p", strconv.Itoa(int(host.Port)))
	}
	if len(host.Exec) > 0 {
		cmd_fmt = append(cmd_fmt, "-t")
	}
	cmd_fmt = append(cmd_fmt, host.User + "@" + host.Host)
	if len(host.Exec) > 0 {
		cmd_fmt = append(cmd_fmt, "--", host.Exec)
	}
	return cmd_fmt, nil
}

func c_format_rdp(host *HostNode, pass string) ([]string, []string) {
	cmd_fmt := []string{"xfreerdp"}

	if len(host.RDPFile) > 0 {
		cmd_fmt = append(cmd_fmt, host.RDPFile)
	} else {
	cmd_fmt = append(cmd_fmt,
		"/v:" + host.Host,
		"/u:" + host.User)
	}
	if len(host.Domain) > 0 {
		cmd_fmt = append(cmd_fmt, "/d:" + host.Domain)
	}
	if len(pass) > 0 {
		cmd_fmt = append(cmd_fmt, "/p:" + pass)
	}
	if host.Port != 0 && len(host.RDPFile) == 0 {
		cmd_fmt = append(cmd_fmt, "/port:" + strconv.Itoa(int(host.Port)))
	}
	if host.Dynamic == true {
		cmd_fmt = append(cmd_fmt, "/dynamic-resolution")
	}
	if host.FullScr == true {
		cmd_fmt = append(cmd_fmt, "/f")
	}
	if host.MultiMon == true {
		cmd_fmt = append(cmd_fmt, "/multimon:force")
	}
	if host.Drive != nil {
		for share, path := range host.Drive {
			cmd_fmt = append(cmd_fmt, "/drive:" + share + "," + path)
		}
	}
	switch host.Quality {
	case 0:
		cmd_fmt = append(cmd_fmt,
			"-aero", "-menu-anims", "-window-drag", "-wallpaper",
			"-decorations", "-fonts", "-themes",
			"/bpp:8", "/compression-level:2")
	case 2:
		cmd_fmt = append(cmd_fmt,
			"+aero", "+menu-anims", "+window-drag",
			"+decorations", "+fonts", "+themes", "/gfx:RFX", "/rfx", "/gdi:hw",
			"/bpp:32")
	}
	cmd_fmt = append(cmd_fmt,
		"/size:" + strconv.Itoa(int(host.Width)) +
		"x" + strconv.Itoa(int(host.Height)))
	return cmd_fmt, nil
}

func c_format_openstack(host *HostNode, pass string) ([]string, []string) {
	cmd_fmt := []string{"openstack"}
	cmd_env := []string{
		"OS_USERNAME="             + host.User,
		"OS_PASSWORD="             + pass,
		"OS_AUTH_URL="             + host.Host,
		"OS_USER_DOMAIN_ID="       + host.Stack.UserDomainID,
		"OS_PROJECT_ID="           + host.Stack.ProjectID,
		"OS_IDENTITY_API_VERSION=" + host.Stack.IdentityAPI,
		"OS_IMAGE_API_VERSION="    + host.Stack.ImageAPI,
		"OS_NETWORK_API_VERSION="  + host.Stack.NetworkAPI,
		"OS_VOLUME_API_VERSION="   + host.Stack.VolumeAPI,
		"OS_REGION_NAME="          + host.Stack.RegionName,
		"OS_ENDPOINT_TYPE="        + host.Stack.EndpointType,
		"OS_INTERFACE="            + host.Stack.Interface,
	}
	return cmd_fmt, cmd_env
}

func c_format_command(host *HostNode, pass string) ([]string, []string){
	return append(host.Shell, host.Host), nil
}

func c_format_cmd(host *HostNode, opts HardOpts,
				  ui *HardUI) ([]string, []string) {
	type format_func func(*HostNode, string) ([]string, []string)
	var pass string
	gpg, term := opts.GPG, opts.Term

	if host.Protocol > PROTOCOL_MAX {
		return nil, nil
	}
	if len(gpg) > 0 && gpg != "plain" && len(host.Pass) > 0 {
		i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " GnuPG ")
		text := "decryption using gpg..."
		left, right := i_left_right(len(text), *ui)
		i_draw_text(ui.s, left, ui.dim[H] - 3, right, ui.dim[H] - 3,
			ui.style[DEF_STYLE], text)
		ui.s.Show()
		var err error
		pass, err = c_decrypt_str(host.Pass)
		if err != nil {
			c_error_mode(host.parent.path() + host.filename +
				": password decryption failed", err, ui)
			return nil, nil
		}
		pass = strings.TrimSuffix(pass, "\n")
	}
	fp := [PROTOCOL_MAX + 1]format_func{
		c_format_ssh,
		c_format_rdp,
		c_format_command,
		c_format_openstack,
	}
	cmd_fmt, cmd_env := fp[host.Protocol](host, pass)
	if len(term) > 0 {
		// TODO: setsid
		if term == "$TERMINAL" {
			term = os.Getenv("TERMINAL")
		}
		cmd_fmt = append([]string{term, "-e"}, cmd_fmt...)
	}
	return cmd_fmt, cmd_env
}

func c_redirect_ssh(host *HostNode, local_port uint16) error {
	rdr_fmt := []string{}
	rdr_fmt = append(rdr_fmt, "ssh", "-f")
	rdr_fmt = append(rdr_fmt, "-L",
		strconv.Itoa(int(local_port)) + ":" +
		host.Host + ":" +
		strconv.Itoa(int(host.Port)))
	if len(host.Jump.Priv) > 0 {
		rdr_fmt = append(rdr_fmt, "-i", host.Jump.Priv)
	}
	if host.Jump.Port != 0 {
		rdr_fmt = append(rdr_fmt, "-p", strconv.Itoa(int(host.Jump.Port)))
	}
	rdr_fmt = append(rdr_fmt, host.Jump.User + "@" + host.Jump.Host,
		"sleep", "5")
	if err := exec.Command(rdr_fmt[0], rdr_fmt[1:]...).Run(); err != nil {
		return err
	}
	return nil
}

func c_exec(host *HostNode, opts HardOpts, ui *HardUI) {
	if host == nil {
		return
	}
	save_host, save_port := host.Host, host.Port
	tmp_host := host
	if host.Protocol == PROTOCOL_RDP && len(host.Jump.Host) != 0 {
		local_host := "127.0.0.1"
		local_port := uint16(rand.IntN(40000) + 4389)
		if err := c_redirect_ssh(host, local_port); err != nil {
			c_error_mode("ssh tunneling", err, ui)
			return
		}
		tmp_host.Host = local_host
		tmp_host.Port = local_port
	}
	cmd_fmt, cmd_env := c_format_cmd(tmp_host, opts, ui)
	if tmp_host.Port != save_port {
		tmp_host.Host = save_host
		tmp_host.Port = save_port
	}
	if cmd_fmt == nil {
		return
	}
	silent := false
	if host.Protocol == PROTOCOL_CMD {
		silent = host.Silent
	}
	if silent == false {
		if err := ui.s.Suspend(); err != nil {
			c_error_mode("screen", err, ui)
			return
		}
	} else {
		i_draw_msg(ui.s, 1, ui.style[BOX_STYLE], ui.dim, " Exec ")
		text := "running command..."
		left, right := i_left_right(len(text), *ui)
		i_draw_text(ui.s, left, ui.dim[H] - 3, right, ui.dim[H] - 3,
			ui.style[DEF_STYLE], text)
		ui.s.Show()
	}
	if err, err_str := c_exec_cmd(cmd_fmt, cmd_env, silent);
	   err != nil && host.Protocol == PROTOCOL_CMD {
		c_error_mode(err_str, err, ui)
	}
	if opts.Loop == false {
		ui.s.Fini()
		os.Exit(0)
	} else if silent == false {
		c_resume_or_die(ui)
	}
}
