package hocr

import (
	"fmt"
	"github.com/eslider/go-hocr/v1_2"
	"sort"
)

type Element struct {
	// Id of each element.
	Id int16 `yaml:"-"`

	// Class short named.
	Class string //`yaml:"-"`

	// BoundingBox - "bounding box" - of an element is a rectangular box around this element,
	// which is defined by the upper-left corner (x0, y0) and the lower-right corner (x1, y1).
	//
	// * the values are with reference to the top-left corner of the document image and measured in pixels
	// * the order of the values are x0 y0 x1 y1 = "left top right bottom"
	// * use x_bboxes below for character bounding boxes
	// * do not use bbox unless the bounding box of the layout component is, in fact, rectangular
	// * some non-rectangular layout components may have rectangular bounding boxes
	//   if the non-rectangularity is caused by floating elements around which text flows
	BoundingBox []int `yaml:"-"` //`yaml:"bbox,flow"`
}

// Search for excluded properties and return a list of them.
// Uses binary search
// See: https://stackoverflow.com/questions/13520111/check-whether-a-string-slice-contains-a-certain-value-in-go
func checkForNewProperties(el *v1_2.Element, list []string) []string {
	var r []string
	l := len(list)
	sort.Strings(list)
	for k, _ := range el.GetProperties() {
		i := sort.SearchStrings(list, k)
		if !(i < l && list[i] == k) {
			r = append(r, k)
		}
	}

	if r != nil {
		fmt.Println("Found new properties:", r)
	}
	return r
}
