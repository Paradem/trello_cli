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
- [x] Display assignee full names instead of IDs in card details
- [x] Add -f flag for field-specific output when used with -c flag
- [x] Add list name to table output
- [x] Sort cards by list name first, then by ShortID
- [x] Add link field to -f flag for Trello website links
- [x] Add created_at field to -f flag for querying card creation dates

## Project Summary
CLI application successfully created with the following features:
- Go-based CLI using BubbleTea for user prompts
- Configuration stored in ~/.config/trello_cli/config.json
- Simple HTTP client for Trello API integration
- Interactive workspace and board selection from available options
- Filters and displays tasks assigned to the user in #123 TITLE format (using numeric IDs)
- CLI flags for filtering: --assigned (default), --all, or --lists for specific lists
- List filtering with comma-separated values (case-insensitive)
- Table output with fixed-width columns for proper alignment
- ID column fixed to 8 characters with proper padding
- Title column fixed to 80 characters with truncation when necessary
- Card detail view with --card flag showing beautifully rendered markdown using glamour
- Cards sorted by list name first, then by ShortID (ascending) for consistent ordering
- External dependencies: BubbleTea, Bubbles, Lipgloss, Glamour for markdown rendering
- Improved error handling for 401 authentication errors
- Fixed API calls to use proper card IDs instead of ShortIDs for card details
- Uses Trello's idShort field for numeric card identifiers
- Shows which list/column each card belongs to in the table output
- CLICOLOR_FORCE support for preserving ANSI colors when piping output
- Displays assignee full names instead of IDs in card details
- Field-specific output with -f flag for extracting individual card fields including Trello links and creation dates

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

# Show specific field from card
trello_cli -c 123 -f title           # Show only the card title
trello_cli -c 123 -f description     # Show only the description
trello_cli -c 123 -f assignees       # Show only assignees (full names)
trello_cli -c 123 -f labels          # Show only labels
trello_cli -c 123 -f list            # Show only the list name
trello_cli -c 123 -f status          # Show only status (Open/Closed)
trello_cli -c 123 -f link            # Show only the Trello link
trello_cli -c 123 -f created_at      # Show only the creation date

# Preserve ANSI colors when piping output
CLICOLOR_FORCE=1 trello_cli --card #123 | cat    # Force colors even when piping
CLICOLOR_FORCE=1 trello_cli -c 456 | less        # Colors preserved in pager
```

## Output Format
```
#123 Task Title              In Progress
#456 Another Task            To Do
#789 Final Task              Done
```
*Table uses fixed-width columns with ID column fixed to 8 characters with proper padding and title fixed to 80 characters for proper alignment*
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
- Assignee names are looked up and displayed as full names instead of IDs
- -f flag allows extracting specific fields from cards for scripting/automation
- Created_at field extracts card creation date from the card ID timestamp