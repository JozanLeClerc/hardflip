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
 * hardflip: src/c_hardflip.go
 * Thu Jan 04 15:37:01 2024
 * Joe
 *
 * the main
 */

package main

// import "fmt"

type HardPtr interface {
	is_dir() bool
	get_self_dirs() *DirsNode
	get_self_host() *HostNode
}

// the main data structure, holds up everything important
type HardData struct {
	litems *ItemsList
	ldirs  *DirsList
	ui     HardUI
	opts   HardOpts
	data_dir string
	ptr    HardPtr
}

func c_reset_ptr(data *HardData) {
	if data.ptr = data.ldirs.head.lhost.head; data.ptr != nil {
		data.ui.line = 1
		return
	}
	data.ptr = data.ldirs.head.next
	data.ui.line = 1
}

func main() {
	data_dir := c_get_data_dir()
	opts := HardOpts{true, true, false}
	litems, ldirs := c_load_data_dir(data_dir, opts)
	data := HardData{
		litems,
		ldirs,
		HardUI{},
		opts,
		data_dir,
		nil,
	}
	c_reset_ptr(&data)

	// for ptr = ldirs.head; ptr != nil ; ptr = ptr.next {
	// 	spaces := ""
	// 	for i := 0; i < int(ptr.Depth - 1) * 2; i++ {
	// 		spaces += " "
	// 	}
	// 	if ptr.is_dir() == true {
	// 		fmt.Print(spaces, "DIR ", ptr.ID, " ")
	// 	}
	// 	fmt.Println(ptr.Name)
	// 	for ptr = ptr.lhost.head; ptr != nil; ptr = ptr.next {
	// 		spaces := ""
	// 		for i := 0; i < int(ptr.Parent.Depth - 1) * 2; i++ {
	// 			spaces += " "
	// 		}
	// 		spaces += " " 
	// 		if ptr.is_dir() == false {
	// 			fmt.Print(spaces, "HOST ", ptr.ID, " ")
	// 		}
	// 		fmt.Println(ptr.Name)
	// 	}
	// }
	// for dir := ldirs.head; dir != nil ; dir = dir.next {
	// 	for host := dir.lhost.head; host != nil; host = host.next {
	// 		fmt.Println(host.ID, host.Name, "HOST")
	// 	}
	// }
	// for item := litems.head; item != nil ; item = item.next {
	// 	if item.Dirs != nil {
	// 		fmt.Println(item.ID, item.Dirs.Name)
	// 	} else {
	// 		fmt.Println(item.ID, item.Host.Name)
	// 	}
	// }

	// PERF: test performance over a large amount of hosts with litems
	i_ui(&data)
}
