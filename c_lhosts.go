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
 * Thu 04 Jan 2024 11:50:55 AM CET
 * Joe
 *
 * the hosts linked list
 */

package main

// 0: ssh
// 1: rdp
type HostNode struct {
	ID        int
	Protocol  int8   `yaml:"type"`
	Name      string `yaml:"name"`
	Host      string `yaml:"host"`
	Port      uint16 `yaml:"port"`
	User      string `yaml:"user"`
	Pass      string `yaml:"pass"`
	Priv      string `yaml:"priv"`
	Jump      string `yaml:"jump"`
	JumpPort uint16 `yaml:"jump_port"`
	JumpUser string `yaml:"jump_user"`
	JumpPass string `yaml:"jump_pass"`
	JumpPriv string `yaml:"jump_priv"`
	Quality   uint8  `yaml:"quality"`
	Domain    string `yaml:"domain"`
	Width     uint16 `yaml:"width"`
	Height    uint16 `yaml:"height"`
	Dynamic   bool   `yaml:"dynamic"`
	Note      string `yaml:"note"`
	Filename  string
	next      *HostNode
}

type HostList struct {
	head *HostNode
	last *HostNode
}

// adds a host node to the list
func (lhost *HostList) add_back(node *HostNode) {
	node.ID = global_id
	if lhost.head == nil {
		lhost.head = node
		lhost.last = lhost.head
		return
	}
	last := lhost.last
	last.next = node
	lhost.last = last.next
}

func (lhost *HostList) reset_id() {
	curr := lhost.head
	for i := 0; curr != nil; i++ {
		curr.ID = i
		curr = curr.next
	}
}

// removes a host node from the list
func (lhost *HostList) del(id int) {
    if lhost.head == nil {
        return
    }
    if lhost.head.ID == id {
        lhost.head = lhost.head.next
        return
    }
    curr := lhost.head
    for curr.next != nil && curr.next.ID != id {
        curr = curr.next
    }
    if curr.next.ID == id {
        curr.next = curr.next.next
    }
}

// return the list node with the according id
func (lhost *HostList) sel(id int) *HostNode {
	curr := lhost.head

	if curr == nil {
		return nil
	}
    for curr.next != nil && curr.ID != id {
        curr = curr.next
    }
	if curr.ID != id {
		return nil
	}
	return curr
}

func (lhost *HostList) count() int {
	curr := lhost.head
	var count int

	for count = 0; curr != nil; count++ {
		curr = curr.next
	}
	return count
}

func (host *HostNode) is_dir() bool {
	return false
}

func (host *HostNode) get_id() int {
	return host.ID
}

func (host *HostNode) get_self_dirs() *DirsNode {
	return nil
}

func (host *HostNode) get_self_host() *HostNode {
	return host
}

func (host *HostNode) protocol_str() string {
	switch host.Protocol {
	case 0: return "SSH"
	case 1: return "RDP"
	default: return ""
	}
}
