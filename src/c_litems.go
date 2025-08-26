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
 * Mon May 13 12:10:57 2024
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
	draw *ItemsNode
}

// adds an item node to the list
func (litems *ItemsList) add_back(node *ItemsNode) {
	if litems.head == nil {
		litems.head = node
		litems.last = litems.head
		litems.curr = litems.head
		litems.draw = litems.head
		return
	}
	last := litems.last
	node.ID = last.ID + 1
	node.prev = last
	last.next = node
	litems.last = last.next
}

// replaces an item
// func (litems *ItemsList) overwrite(node *ItemsNode) {
// 	if litems.head == nil || litems.curr == nil {
// 		litems.add_back(node)
// 		return
// 	}
// 	curr := litems.curr
// 	node.prev = curr.prev
// 	node.next = curr.next
// 	if node.next != nil {
// 		curr.next.prev = node
// 	}
// 	if litems.last == curr {
// 		litems.last = node
// 	}
// 	if curr.prev != nil {
// 		curr.prev.next = node
// 	}
// 	litems.curr = node
// }

// adds an item node to the list after the current selected item
func (litems *ItemsList) add_after(node *ItemsNode) {
	if litems.head == nil || litems.curr == nil {
		litems.add_back(node)
		return
	}
	curr := litems.curr
	node.prev = curr
	node.next = curr.next
	curr.next = node
	if node.next != nil {
		node.next.prev = node
	}
	if litems.last == curr {
		litems.last = node
	}
	litems.curr = node
}

// removes an item node from the list and resets the ids
func (litems *ItemsList) del(item *ItemsNode) {
    if litems.head == nil {
        return
    }
    if litems.head == item {
        litems.head = litems.head.next
		if litems.head == nil {
			litems.last, litems.curr, litems.draw = nil, nil, nil
			return
		}
		litems.head.prev = nil
		litems.curr, litems.draw = litems.head, litems.head
		for ptr := litems.head; ptr != nil; ptr = ptr.next {
			ptr.ID -= 1
		}
        return
    }
	if litems.last == item {
		litems.last = litems.last.prev
		litems.last.next = nil
		litems.curr = litems.last
		if litems.draw == item {
			litems.draw = litems.last
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

// returns the previous dir
func (item *ItemsNode) prev_dir() *ItemsNode {
	for ptr := item.prev; ptr != nil; ptr = ptr.prev {
		if ptr.is_dir() == true {
			return ptr
		}
	}
	return nil
}

// returns the next dir
func (item *ItemsNode) next_dir() *ItemsNode {
	for ptr := item.next; ptr != nil; ptr = ptr.next {
		if ptr.is_dir() == true {
			return ptr
		}
	}
	return nil
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
			item_ptr = item_ptr.next
			continue
		}
		if item_ptr.Dirs == ptr {
			return item_ptr
		}
		item_ptr = item_ptr.next
	}
	return nil
}

func (litems *ItemsList) reset_id() {
	if litems.head != nil {
		litems.head.ID = 1
	}
	for ptr := litems.head; ptr != nil && ptr.next != nil; ptr = ptr.next {
		ptr.next.ID = ptr.ID + 1
	}
}

func (item *ItemsNode) path_node() *DirsNode {
	if item.is_dir() == true {
		return item.Dirs
	} else {
		return item.Host.parent
	}
}

func (item *ItemsNode) path() string {
	return item.path_node().path()
}
