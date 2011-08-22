/*
Copyright 2011 Mostafa Hajizdeh

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pdf

// This file contains graphics-related functions and constants for type Document.

// TODO Describe units in documentation.

import (
	"fmt"
)

const (
	LineCapBut = iota
	LineCapRound
	LineCapProjecting

	LineJoinMiter = iota
	LineJoinRound
	LineJoinBevel
)

// LineWidth changes the width of the lines to be drawn after it.
func (d *Document) LineWidth(w int) {
	d.addc(fmt.Sprint(w, " w"))
}

// LineCapStyle changes line cap style to one of the three options. Use
// LineCapBut, LineCapRound, and LineCapProjecting constants as the argument.
func (d *Document) LineCapStyle(s int) {
	d.addc(fmt.Sprint(s, " J"))
}

// LineJoinStyle changes line join style to one of the three options. Use
// LineJoinMiter, LineJoinRound, LineJoinBevel constants as the argument.
func (d *Document) LineJoinStyle(s int) {
	d.addc(fmt.Sprint(s, " j"))
}

// MoveTo starts a new path at the given point.
func (d *Document) MoveTo(x, y int) {
	d.addc(fmt.Sprint(x, y, " m"))
}

// LineTo draws a single line from current the given point.
func (d *Document) LineTo(x, y int) {
	d.addc(fmt.Sprint(x, y, " l"))
}

// Curve draws a bézier curve from current point to point (x2, y2) using
// (x0, y0) and (x1, y1) as control points.
func (d *Document) Curve(x0, y0, x1, y1, x2, y2 int) {
	d.addc(fmt.Sprint(x0, y0, x1, y1, x2, y2, " c"))
}

// CurveV draws a bézier curve from current point to point (x1, y1) using
// current point and (x0, y0) as control points.
func (d *Document) CurveV(x0, y0, x1, y1 int) {
	d.addc(fmt.Sprint(x0, y0, x1, y1, " v"))
}

// CurveY draws a bézier curve from current point to point (x1, y1) using
// (x0, y0) and current point as control points.
func (d *Document) CurveY(x0, y0, x1, y1 int) {
	d.addc(fmt.Sprint(x0, y0, x1, y1, " y"))
}

// Rectangle draws a renctangle using PDF's 're' command.
func (d *Document) Rectangle(x, y, w, h int) {
	d.addc(fmt.Sprint(x, y, w, h, " re"))
}

// ClosePath closes the current active path by drawing a straight line from
// current point to the beginning of the path.
func (d *Document) ClosePath() {
	d.addc("h")
}

// Stroke paints the current path with stroke.
func (d *Document) Stroke() {
	d.addc("S")
}

// Fill paints inside of the current path.
func (d *Document) Fill() {
	d.addc("f")
}
