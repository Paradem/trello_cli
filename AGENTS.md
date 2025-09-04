# Trello CLI Development Progress

## Project Overview
CLI application in Go that connects to Trello REST API to query tasks assigned to user.

## Completed Tasks
- [x] Set up Go project structure with main.go and go.mod
- [x] Implement config file handling for ~/.config/trello_cli/config.json
- [x] Create BubbleTea prompts for API key, workspace, and board input
- [x] Implement HTTP client for Trello API calls
- [x] Parse Trello API responses and filter tasks assigned to user
- [x] Format and output tasks in tabular format (#123 TITLE)
- [x] Track work progress in AGENTS.md file
- [x] Add markdown rendering with glamour package for card details
- [x] Fix card detail API calls to use proper card IDs instead of ShortIDs
- [x] Add CLICOLOR_FORCE support for preserving ANSI colors when piping output (uses glamour.Render with dark style)

## Project Summary
CLI application successfully created with the following features:
- Go-based CLI using BubbleTea for user prompts
- Configuration stored in ~/.config/trello_cli/config.json
- Simple HTTP client for Trello API integration
- Interactive workspace and board selection from available options
- Filters and displays tasks assigned to the user in #123 TITLE format (using numeric IDs)
- CLI flags for filtering: --assigned (default), --all, or --lists for specific lists
- List filtering with comma-separated values (case-insensitive)
- Table output with fixed-size columns for proper alignment
- Card detail view with --card flag showing beautifully rendered markdown using glamour
- Cards sorted by ShortID (ascending) for consistent ordering
- External dependencies: BubbleTea, Bubbles, Lipgloss, Glamour for markdown rendering
- Improved error handling for 401 authentication errors
- Fixed API calls to use proper card IDs instead of ShortIDs for card details
- Uses Trello's idShort field for numeric card identifiers
- Shows which list/column each card belongs to
- CLICOLOR_FORCE support for preserving ANSI colors when piping output

## CLI Usage
```
trello_cli                    # Show cards assigned to current user (default)
trello_cli --assigned         # Show cards assigned to current user
trello_cli --all              # Show all cards on the board
trello_cli -a                 # Short flag for --assigned
trello_cli -A                 # Short flag for --all

# Filter by specific lists
trello_cli -l "In Progress,Needs Review"    # Show cards from specific lists
trello_cli --lists "To Do,Done"             # Long form list filtering
trello_cli --all -l "Backlog"               # Combine with --all flag

# Show detailed card information
trello_cli --card #123        # Show detailed info for card #123
trello_cli --card 123         # Same as above (bare integer)
trello_cli -c #456            # Short form for card details
trello_cli -c 456             # Same as above (bare integer)

# Preserve ANSI colors when piping output
CLICOLOR_FORCE=1 trello_cli --card #123 | cat    # Force colors even when piping
CLICOLOR_FORCE=1 trello_cli -c 456 | less        # Colors preserved in pager
```

## Output Format
```
#123 Task Title
#456 Another Task
#789 Final Task
```
*Columns are automatically sized based on the longest ID and content for proper alignment*
- [ ] Create BubbleTea prompts
- [ ] Implement HTTP client for Trello API
- [ ] Parse API responses and filter tasks
- [ ] Format output in tabular format
- [ ] Track work progress

## Notes
- External packages: BubbleTea for prompts, Lipgloss for styling, Glamour for markdown rendering
- Config stored in ~/.config/trello_cli/config.json
- Output format: #123 TITLE (no headings)
- Card details rendered with beautiful markdown formatting using Glamour
- Fixed API integration to properly handle card IDs for detailed views
- CLICOLOR_FORCE=1 environment variable preserves ANSI colors when piping output