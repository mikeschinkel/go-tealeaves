package teacrumbs

// Renderer generates the display string for a single breadcrumb.
// Follows the http.Handler pattern: interface + func adapter.
type Renderer interface {
	Render(index int, model BreadcrumbsModel) string
}

// RendererFunc adapts a plain function to the Renderer interface.
type RendererFunc func(index int, model BreadcrumbsModel) string

// Render implements the Renderer interface.
func (f RendererFunc) Render(index int, model BreadcrumbsModel) string {
	return f(index, model)
}
