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
 * Thu Jan 11 18:37:44 2024
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
	draw_start *ItemsNode
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

// removes an item node from the list and resets the ids
func (litems *ItemsList) del(item *ItemsNode) {
    if litems.head == nil {
        return
    }
    if litems.head == item {
        litems.head = litems.head.next
		if litems.head == nil {
			litems.last, litems.curr, litems.draw_start = nil, nil, nil
			return
		}
		litems.head.prev = nil
		litems.curr, litems.draw_start = litems.head, litems.head
		for ptr := litems.head; ptr != nil; ptr = ptr.next {
			ptr.ID -= 1
		}
        return
    }
	if litems.last == item {
		litems.last = litems.last.prev
		litems.last.next = nil
		litems.curr = litems.last
		if litems.draw_start == item {
			litems.draw_start = litems.last
		}
		return
	}
    ptr := litems.head
    for ptr.next != nil && ptr.next != item {
        ptr = ptr.next
    }
    if ptr.next == item {
        ptr.next = ptr.next.next
		ptr.next.prev = ptr
    }
	for ptr := ptr.next; ptr != nil; ptr = ptr.next {
		ptr.ID -= 1
	}
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

func (litems *ItemsList) inc(jump int) {
	new_item := litems.curr

	if new_item == nil || jump == 0 {
		return
	} else if jump == +1 {
		if new_item.next != nil {
			new_item = new_item.next
		}
	} else if jump == -1 {
		if new_item.prev != nil {
			new_item = new_item.prev
		}
	} else {
		for i := 0; jump > +1 && new_item.next != nil && i < jump; i++ {
			new_item = new_item.next
		}
		for i := 0; jump < -1 && new_item.prev != nil && i > jump; i-- {
			new_item = new_item.prev
		}
	}
	litems.curr = new_item
}

// returns the next directory in line with the same or lower depth
func (item *ItemsNode) get_next_level() *ItemsNode {
	if item == nil || item.Dirs == nil {
		return nil
	}
	dir := item.Dirs
	ptr := dir.next
	for ptr != nil && ptr.Depth > dir.Depth {
		ptr = ptr.next
	}
	item_ptr := item
	for item_ptr != nil {
		if item_ptr.is_dir() == false {
			continue
		}
		if item_ptr.Dirs == ptr {
			break
		}
		item_ptr = item_ptr.next
	}
	return item_ptr
}

