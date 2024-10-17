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
 * hardflip: src/c_lfuzz.go
 * Thu, 17 Oct 2024 13:22:18 +0200
 * Joe
 *
 * fuzz hard
 */

package main

type FuzzNode struct {
	ptr  *ItemsNode
	name string
	prev *FuzzNode
	next *FuzzNode
}

type FuzzList struct {
	head *FuzzNode
	last *FuzzNode
	curr *FuzzNode
	draw *FuzzNode
}

// adds a fuzz node to the list
func (lfuzz *FuzzList) add_back(node *ItemsNode) {
	name := ""
	if node.is_dir() == false {
		name = node.Host.Name
	} else {
		name = node.Dirs.Name
	}
	fuzz_node := &FuzzNode{
		node,
		name,
		nil,
		nil,
	}
	if lfuzz.head == nil {
		lfuzz.head = fuzz_node
		lfuzz.last = lfuzz.head
		lfuzz.curr = lfuzz.head
		lfuzz.draw = lfuzz.head
		return
	}
	last := lfuzz.last
	fuzz_node.prev = last
	last.next = fuzz_node
	lfuzz.last = last.next
}

// removes n fuzz node from the list
func (lfuzz *FuzzList) del(item *FuzzNode) {
    if lfuzz.head == nil {
        return
    }
    if lfuzz.head == item {
        lfuzz.head = lfuzz.head.next
		if lfuzz.head == nil {
			lfuzz.last, lfuzz.curr, lfuzz.draw = nil, nil, nil
			return
		}
		lfuzz.head.prev = nil
		lfuzz.curr, lfuzz.draw = lfuzz.head, lfuzz.head
        return
    }
	if lfuzz.last == item {
		lfuzz.last = lfuzz.last.prev
		lfuzz.last.next = nil
		lfuzz.curr = lfuzz.last
		if lfuzz.draw == item {
			lfuzz.draw = lfuzz.last
		}
		return
	}
    ptr := lfuzz.head
    for ptr.next != nil && ptr.next != item {
        ptr = ptr.next
    }
    if ptr.next == item {
        ptr.next = ptr.next.next
		ptr.next.prev = ptr
    }
}
