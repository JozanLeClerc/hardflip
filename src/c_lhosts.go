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
 * hardflip: src/c_lhosts.go
 * Thu Feb 01 16:22:33 2024
 * Joe
 *
 * the hosts linked list
 */

package main


type StackSettings struct {
  UserDomainID string `yaml:"user_domain_id,omitempty"`
  ProjectID    string `yaml:"project_id,omitempty"`
  IdentityAPI  string `yaml:"identity_api_version,omitempty"`
  ImageAPI     string `yaml:"image_api_version,omitempty"`
  NetworkAPI   string `yaml:"network_api_version,omitempty"`
  VolumeAPI    string `yaml:"volume_api_version,omitempty"`
  RegionName   string `yaml:"region_name,omitempty"`
  EndpointType string `yaml:"endpoint_type,omitempty"`
  Interface    string `yaml:"interface,omitempty"`
}

type JumpSettings struct {
	Host string `yaml:"host,omitempty"`
	Port uint16 `yaml:"port,omitempty"`
	User string `yaml:"user,omitempty"`
	Pass string `yaml:"pass,omitempty"`
	Priv string `yaml:"priv,omitempty"`
}

// 0: ssh
// 1: rdp
// 2: single cmd
// 3: openstack
type HostNode struct {
	Protocol int8     `yaml:"type"`
	Name     string   `yaml:"name,omitempty"`
	Host     string   `yaml:"host,omitempty"`
	Port     uint16   `yaml:"port,omitempty"`
	User     string   `yaml:"user,omitempty"`
	Pass     string   `yaml:"pass,omitempty"`
	Priv     string   `yaml:"priv,omitempty"`
	RDPFile  string   `yaml:"rdp_file,omitempty"`
	Jump     JumpSettings `yaml:"jump,omitempty"`
	Quality  uint8    `yaml:"quality,omitempty"`
	Domain   string   `yaml:"domain,omitempty"`
	Width    uint16   `yaml:"width,omitempty"`
	Height   uint16   `yaml:"height,omitempty"`
	Dynamic  bool     `yaml:"dynamic,omitempty"`
	Note     string   `yaml:"note,omitempty"`
	Drive    map[string]string `yaml:"drive,omitempty"`
	Silent   bool     `yaml:"silent,omitempty"`
	Shell    []string `yaml:"shell,omitempty"`
	Stack    StackSettings `yaml:"openstack,omitempty"`
	filename string
	parent   *DirsNode
	next     *HostNode
}

type HostList struct {
	head *HostNode
	last *HostNode
}

// adds a host node to the list
func (lhost *HostList) add_back(node *HostNode) {
	if lhost.head == nil {
		lhost.head = node
		lhost.last = lhost.head
		return
	}
	last := lhost.last
	last.next = node
	lhost.last = last.next
}

// removes a host node from the list
func (lhost *HostList) del(host *HostNode) {
    if lhost.head == nil {
        return
    }
    if lhost.head == host {
        lhost.head = lhost.head.next
		if lhost.head == nil {
			lhost.last = nil
			return
		}
        return
    }
	if lhost.last == host {
		ptr := lhost.head
		for ptr.next != nil {
			ptr = ptr.next
		}
		lhost.last = ptr
		lhost.last.next = nil
		return
	}
    ptr := lhost.head
    for ptr.next != nil && ptr.next != host {
        ptr = ptr.next
    }
    if ptr.next == host {
        ptr.next = ptr.next.next
    }
}

func (lhost *HostList) count() int {
	curr := lhost.head
	var count int

	for count = 0; curr != nil; count++ {
		curr = curr.next
	}
	return count
}
