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
 * hardflip: src/c_litems.go
 * Mon Jan 08 11:53:22 2024
 * Joe
 *
 * the dir and hosts linked list
 */

package main

type ItemsNode struct {
	ID   int
	Dirs *DirsNode
	Host *HostNode
	prev *ItemsNode
	next *ItemsNode
}

type ItemsList struct {
	head *ItemsNode
	last *ItemsNode
	curr *ItemsNode
}

// adds an item node to the list
func (litems *ItemsList) add_back(node *ItemsNode) {
	if litems.head == nil {
		litems.head = node
		litems.last = litems.head
		return
	}
	last := litems.last
	node.ID = last.ID + 1
	node.prev = last
	last.next = node
	litems.last = last.next
}

// sets litems.curr to be used
func (litems *ItemsList) sel(id int) {
	curr := litems.head

	if curr == nil {
		litems.curr = nil
	}
    for curr.next != nil && curr.ID != id {
        curr = curr.next
    }
	if curr.ID != id {
		litems.curr = nil
	}
	litems.curr = curr
}

func (item *ItemsNode) is_dir() bool {
	if item.Dirs == nil {
		return false
	}
	return true
}

func (item *ItemsNode) inc(jump int) *ItemsNode {
	if item == nil {
		return nil
	}
	if jump == 0 {
		return item
	} else if jump == 1 {
		if item.next != nil {
			return item.next
		}
		return item
	} else if jump == -1 {
		if item.prev != nil {
			return item.prev
		}
		return item
	}
	new_item := item
	if jump > 0 {
		for i := 0; new_item.next != nil && i < jump; i++ {
			new_item = new_item.next
		}
		return new_item
	}
	for i := 0; new_item.prev != nil && i > jump; i-- {
		new_item = new_item.prev
	}
	return new_item
}
