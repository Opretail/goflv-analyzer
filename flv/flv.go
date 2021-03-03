//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package flv

import (
	"goflv-analyzer/flv/tag"
)

// Flags ...
type Flags uint8

const (
	// FlagsAudio ...
	FlagsAudio Flags = 0x01
	// FlagsVideo ...
	FlagsVideo = 0x02
)

// HeaderSignature ...
var HeaderSignature = []byte{0x46, 0x4c, 0x56} // F, L, V
// HeaderLength ...
const HeaderLength uint32 = 9

// Header ...
type Header struct {
	Version uint8
	Flags
	DataOffset uint32
}

// Body ...
type Body struct {
	Tags []*tag.FlvTag
}
