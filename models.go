// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

// BindingValueProvider is the interface that a model must implement to support
// data binding with widgets like ComboBox.
type BindingValueProvider interface {
	BindingValue(index int) interface{}
}

// ListModel is the interface that a model must implement to support widgets
// like ComboBox.
type ListModel interface {
	// ItemCount returns the number of items in the model.
	ItemCount() int

	// Value returns the value that should be displayed for the given index.
	Value(index int) interface{}

	// ItemsReset returns the event that the model should publish when the 
	// number of its items changes.
	ItemsReset() *Event

	// ItemChanged returns the event that the model should publish when an item
	// was changed.
	ItemChanged() *IntEvent
}

// ListModelBase implements the ItemsReset and ItemChanged methods of the
// ListModel interface.
type ListModelBase struct {
	itemsResetPublisher  EventPublisher
	itemChangedPublisher IntEventPublisher
}

func (lmb *ListModelBase) ItemsReset() *Event {
	return lmb.itemsResetPublisher.Event()
}

func (lmb *ListModelBase) ItemChanged() *IntEvent {
	return lmb.itemChangedPublisher.Event()
}

func (lmb *ListModelBase) PublishItemsReset() {
	lmb.itemsResetPublisher.Publish()
}

func (lmb *ListModelBase) PublishItemChanged(index int) {
	lmb.itemChangedPublisher.Publish(index)
}

// TableColumn provides column information for widgets like TableView.
type TableColumn struct {
	// Name is the optional name of the column.
	Name string

	// Title is the text to display in the column header.
	Title string

	// Format is the format string for converting a value into a string.
	Format string

	// Precision is the number of decimal places for formatting float32, float64
	// or big.Rat values.
	Precision int

	// Width is the width of the column in pixels.
	Width int

	// Alignment is the alignment of the column (who would have thought).
	Alignment Alignment1D
}

// TableModel is the interface that a model must implement to support widgets
// like TableView.
type TableModel interface {
	// Columns returns information about the columns of the model.
	Columns() []TableColumn

	// RowCount returns the number of rows in the model.
	RowCount() int

	// Value returns the value that should be displayed for the given cell.
	Value(row, col int) interface{}

	// RowsReset returns the event that the model should publish when the number
	// of its rows changes.
	RowsReset() *Event

	// RowChanged returns the event that the model should publish when a row was
	// changed.
	RowChanged() *IntEvent
}

// TableModelBase implements the RowsReset and RowChanged methods of the
// TableModel interface.
type TableModelBase struct {
	rowsResetPublisher  EventPublisher
	rowChangedPublisher IntEventPublisher
}

func (tmb *TableModelBase) RowsReset() *Event {
	return tmb.rowsResetPublisher.Event()
}

func (tmb *TableModelBase) RowChanged() *IntEvent {
	return tmb.rowChangedPublisher.Event()
}

func (tmb *TableModelBase) PublishRowsReset() {
	tmb.rowsResetPublisher.Publish()
}

func (tmb *TableModelBase) PublishRowChanged(row int) {
	tmb.rowChangedPublisher.Publish(row)
}

// ImageProvider is the interface that a model must implement to support
// displaying an item image. 
type ImageProvider interface {
	// Image returns the image to display for the item at index index.
	//
	// Supported types are *walk.Bitmap, *walk.Icon and string. A string will be
	// interpreted as a file path and the icon associated with the file will be
	// used. It is not supported to use strings together with the other options
	// in the same model instance.
	Image(index int) interface{}
}

// ItemChecker is the interface that a model must implement to support check 
// boxes in a widget like TableView.
type ItemChecker interface {
	// Checked returns if the specified item is checked.
	Checked(index int) bool

	// SetChecked sets if the specified item is checked.
	SetChecked(index int, checked bool) error
}

// SortOrder specifies the order by which items are sorted.
type SortOrder int

const (
	// SortAscending specifies ascending sort order.
	SortAscending SortOrder = iota

	// SortDescending specifies descending sort order.
	SortDescending
)

// Sorter is the interface that a model must implement to support sorting with a
// widget like TableView.
type Sorter interface {
	// ColumnSortable returns whether column col is sortable.
	ColumnSortable(col int) bool

	// Sort sorts column col in order order.
	//
	// If col is -1 then no column is to be sorted. Sort must publish the event
	// returned from SortChanged() after sorting.
	Sort(col int, order SortOrder) error

	// SortChanged returns an event that is published after sorting.
	SortChanged() *Event

	// SortedColumn returns the index of the currently sorted column, or -1 if
	// no column is currently sorted.
	SortedColumn() int

	// SortOrder returns the current sort order.
	SortOrder() SortOrder
}

// SorterBase implements the Sorter interface.
//
// You still need to provide your own implementation of at least the Sort method
// to actually sort and reset the model. Your Sort method should call the
// SorterBase implementation so the SortChanged event, that e.g. a TableView
// widget depends on, is published.
type SorterBase struct {
	changedPublisher EventPublisher
	col              int
	order            SortOrder
}

func (sb *SorterBase) ColumnSortable(col int) bool {
	return true
}

func (sb *SorterBase) Sort(col int, order SortOrder) error {
	sb.col, sb.order = col, order

	sb.changedPublisher.Publish()

	return nil
}

func (sb *SorterBase) SortChanged() *Event {
	return sb.changedPublisher.Event()
}

func (sb *SorterBase) SortedColumn() int {
	return sb.col
}

func (sb *SorterBase) SortOrder() SortOrder {
	return sb.order
}
