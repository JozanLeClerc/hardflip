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
 * Thu 04 Jan 2024 11:50:52 AM CET
 * Joe
 *
 * the directories linked list
 */

package main

type DirsNode struct {
	ID     int
	Name   string
	Parent *DirsNode
	Depth  uint16
	lhost  *HostList
	Folded bool
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
	node.ID = last.ID + 1
	last.next = node
	ldirs.last = last.next
}

// return the list node with the according id
func (ldirs *DirsList) sel(id int) *DirsNode {
	curr := ldirs.head

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

// returns a string with the full path of the dir
func (ldirs *DirsList) path(node *DirsNode) string {
	var path string

	if node == nil {
		return ""
	}
	curr := node
	for curr != nil {
		path = curr.Name + "/" + path
		curr = curr.Parent
	}
	return path
}

func (ldirs *DirsList) count() (int, int) {
	curr := ldirs.head
	var count_dirs int
	var count_hosts int

	for count_dirs = 0; curr != nil; count_dirs++ {
		count_hosts += curr.lhost.count()
		curr = curr.next
	}
	return count_dirs, count_hosts
}
