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
 * hardflip: src/c_parse.go
 * Fri, 15 Dec 2023 10:02:29 +0100
 * Joe
 *
 * parsing of the global data
 */

package main

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func c_parse_opts(file string) (HardOpts, error) {
	var opts HardOpts

	yaml_file, err := os.ReadFile(file)
	if err != nil {
		return opts, err
	}
	err = yaml.Unmarshal(yaml_file, &opts)
	return opts, err
}

func c_read_yaml_file(file string, ui *HardUI) (*HostNode, error) {
	var host HostNode
	yaml_file, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(yaml_file, &host); err != nil {
		err = errors.New(fmt.Sprintf("%s: %v", file, err))
		return nil, err
	}
	if len(host.Name) == 0 {
		return nil, nil
	}
	if len(host.Host) == 0 {
		return nil, nil
	}
	if host.Protocol == 0 {
		if host.Port == 0 {
			host.Port = 22
		}
		if len(host.User) == 0 {
			host.User = "root"
		}
		if len(host.Jump) > 0 {
			if host.JumpPort == 0 {
				host.JumpPort = 22
			}
			if len(host.JumpUser) == 0 {
				host.JumpUser = "root"
			}
		}
	} else if host.Protocol == 1 {
		if len(host.User) == 0 {
			host.User = "Administrator"
		}
		if host.Port == 0 {
			host.Port = 3389
		}
		if host.Width == 0 {
			host.Width = 1600
		}
		if host.Height == 0 {
			host.Height = 1200
		}
	} else if host.Protocol > 1 {
		return nil, errors.New(file + ": unknown protocol")
	}
	if host.Quality > 2 {
		host.Quality = 2
	}
	return &host, nil
}
