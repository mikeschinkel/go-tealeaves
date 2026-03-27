# teaguide example

Demonstrates the teaguide workflow guide overlay with a simulated
project deployment workflow.

## Running

```bash
go run .
```

## Usage

The app simulates a 3-step deployment workflow:

1. **Not Tested** — Press [T] to run tests
2. **Tested** — Press [D] to deploy
3. **Deployed** — Press [L] to release

Press **[N]** at any time to open the guide overlay, which shows
context-aware next steps based on the current workflow state.

Inside the guide:
- Press an action key (e.g., [T]) to close the guide AND execute the action
- Press [Esc] or [N] to dismiss the guide
- Press [Space] or [Enter] to expand/collapse the "Not Yet Available" section
- Use arrow keys or j/k to scroll
