package tealayout

// child pairs a widget with a layout constraint. The widget is stored as any;
// during rendering, the engine checks whether it implements SizeHinter,
// SetSizer, or Viewer.
type child struct {
	Widget     any
	Constraint constraint
}
