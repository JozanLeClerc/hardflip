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
 * hardflip: src/c_ldirs.go
 * Thu Jan 11 18:35:31 2024
 * Joe
 *
 * the directories linked list
 */

package main

type DirsNode struct {
	Name   string
	Parent *DirsNode
	Depth  uint16
	lhost  *HostList
	next   *DirsNode
}

type DirsList struct {
	head *DirsNode
	last *DirsNode
}

// adds a directory node to the list
func (ldirs *DirsList) add_back(node *DirsNode) {
	if ldirs.head == nil {
		ldirs.head = node
		ldirs.last = ldirs.head
		return
	}
	last := ldirs.last
	last.next = node
	ldirs.last = last.next
}

// removes a dir node from the list
// func (ldirs *DirsList) del(dir *DirsNode) {
// 	if ldirs.head == nil {
// 		return
// 	}
// 	if ldirs.head == dir {
// 		ldirs.head = ldirs.head.next
// 		if ldirs.head == nil {
// 			ldirs.last = nil
// 			return
// 		}
// 		return
// 	}
// 	if ldirs.last == dir {
// 		ptr := ldirs.head
// 		for ptr.next != nil {
// 			ptr = ptr.next
// 		}
// 		ldirs.last = ptr
// 		ldirs.last.next = nil
// 		return
// 	}
// 	ptr := ldirs.head
// 	for ptr.next != nil && ptr.next != dir {
// 		ptr = ptr.next
// 	}
// 	if ptr.next == dir {
// 		ptr.next = ptr.next.next
// 	}
// }

// returns a string with the full path of the dir
func (dir *DirsNode) path() string {
	var path string

	if dir == nil {
		return "/"
	}
	curr := dir
	for curr != nil {
		path = curr.Name + "/" + path
		curr = curr.Parent
	}
	return path
}

// func (ldirs *DirsList) prev(dir *DirsNode) *DirsNode {
// 	if ldirs.head == dir {
// 		return dir
// 	}
// 	for ptr := ldirs.head; ptr != nil; ptr = ptr.next {
// 		if ptr.next == dir {
// 			return ptr
// 		}
// 	}
// 	return nil
// }
