// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build boringcrypto

// Package fipsonly restricts all TLS configuration to FIPS-approved settings.
//
// The effect is triggered by importing the package anywhere in a program, as in:
//
//	import _ "crypto/tls/fipsonly"
//
// This package only exists when using Go compiled with GOEXPERIMENT=boringcrypto.
package fipsonly

// This functionality is provided as a side effect of an import to make
// it trivial to add to an existing program. It requires only a single line
// added to an existing source file, or it can be done by adding a whole
// new source file and not modifying any existing source files.

import (
	"github.com/panjf2000/gnet/v2/pkg/tls/internal/boring/fipstls"
	"github.com/panjf2000/gnet/v2/pkg/tls/internal/boring/sig"
)

func init() {
	fipstls.Force()
	sig.FIPSOnly()
}
