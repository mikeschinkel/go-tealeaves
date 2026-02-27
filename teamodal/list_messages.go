package teamodal

// ItemSelectedMsg is sent when user presses Space to preview-select an item.
// The dialog remains open for further browsing.
type ItemSelectedMsg[T ListItem] struct {
	Item T
}

// ListAcceptedMsg is sent when user presses Enter to accept the current selection
// and close the dialog. Contains the accepted item (cursor item if no active item).
type ListAcceptedMsg[T ListItem] struct {
	Item T
}

// NewItemRequestedMsg is sent when user presses 'a' to create new item
type NewItemRequestedMsg struct{}

// DeleteItemRequestedMsg is sent when user presses 'd' to delete an item
type DeleteItemRequestedMsg[T ListItem] struct {
	Item T
}

// ListCancelledMsg is sent when user presses Esc to close without selection
type ListCancelledMsg struct{}

// EditCompletedMsg is sent when user completes inline editing with Enter.
// Contains the item being edited and the new label text.
type EditCompletedMsg[T ListItem] struct {
	Item     T
	NewLabel string
}
