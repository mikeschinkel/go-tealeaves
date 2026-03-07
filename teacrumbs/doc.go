// MIT License
//
// Copyright (c) 2025 Mike Schinkel <mike@newclarity.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package teacrumbs provides a Bubble Tea v2 component for rendering
// a breadcrumb trail navigation bar.
//
// The breadcrumb trail displays the current navigation path with styled
// parent crumbs, separator, and current crumb. The trail automatically
// truncates to fit within the available width.
//
// Usage:
//
//	bc := teacrumbs.NewBreadcrumbsModel().
//	    Push(teacrumbs.Crumb{Text: "Home"}).
//	    Push(teacrumbs.Crumb{Text: "Settings"}).
//	    SetSize(80)
//
//	// In View():
//	bc.View()
//
// # Stability
//
// This package is provisional as of v0.1.0. The public API may change in
// minor releases until promoted to stable.
package teacrumbs
