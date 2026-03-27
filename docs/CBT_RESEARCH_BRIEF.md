CBT (Cursor Backward Tabulation) in Charmbracelet's stack flows from `x/ansi` sequence generation through `x/vt` emulator interpretation to `cellbuf` diff-based rendering, with default 8-column tab stops, but no public test cases specifically exercise CBT and testing infrastructure is still being built (PR #46)
105 sources

Ultraviolet's "Cursed Render" uses an ncurses-inspired cell-based diffing algorithm with double-buffered screens, dirty-cell detection, optimized cursor movement via `relativeCursorMove()` (choosing shortest escape sequences including hard tabs), and style-change minimization to efficiently update only changed terminal cells
125 sources

Ultraviolet's `relativeCursorMove()` follows the ncurses `mvcur()` algorithm, selecting the cheapest cursor movement strategy (CBT, CUB, CR+CUF, CHA, or CUP) by comparing byte costs, with CBT chosen when backward tab stops exist between current and target positions and the total cost (including residual CUB correction) is lower than alternatives — though direct source verification was blocked by tool limitations
131 sources

Bubble Tea v2 replaced line-based diffing with a cell-based "cursed renderer" backed by ultraviolet's 2D screen buffer, where View() now returns a struct (not a string), each cell is independently tracked for changes, and only dirty cells generate minimal ANSI output — dramatically improving efficiency especially over SSH
113 sources

ncurses mvcur() optimizer uses CBT (back tab) for backward cursor movement in a narrow sweet spot: when moving ~5-20 columns left with the target near an 8-column tab stop boundary, beating both repeated cursor-left and carriage-return strategies — triggered in TUI cell-based diff rendering when right-side dirty cells on one row are followed by left-side dirty cells on the next row (e.g., columnar table layouts with staggered updates)
77 sources

Charmbracelet's "Cursed Renderer" uses a double-buffer cell-based dirty tracking algorithm (inspired by ncurses but without terminfo dependencies) that evolved from cellbuf in charmbracelet/x into ultraviolet, minimizing terminal I/O by diffing old/new buffers and only emitting ANSI sequences for changed cells
120 sources

Bubble Tea v2's "Cursed Renderer" uses ultraviolet's ncurses-inspired cell-based diffing algorithm that minimizes terminal writes by choosing the most byte-efficient ANSI cursor movement sequence (CUB, CBT, CHA, or carriage return) for each changed cell during screen buffer comparison
103 sources

Bubble Tea v2 replaced its line-based renderer with a cell-based "cursed renderer" built on `cellbuf.Screen` (not "ultraviolet," which handles input parsing), enabling ncurses-style cell-level diffing for minimal terminal output, where `model.View()` returns a `View` struct whose content flows through `renderer.write()` into the cellbuf screen for optimized ANSI diff rendering at up to 120 FPS
82 sources

Ultraviolet's cell-based renderer uses ncurses-inspired double-buffer diffing to minimize terminal output, with `relativeCursorMove` likely optimizing cursor positioning via tabs, carriage returns, and ANSI sequences, though actual source code was inaccessible
124 sources


