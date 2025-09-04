# Trello CLI

A powerful command-line interface for interacting with Trello boards, built in Go. View your assigned tasks, get detailed card information, and extract specific fields for automation and scripting.

## Features

- 📋 **Task Overview**: List all cards assigned to you or view all cards on a board
- 🎨 **Beautiful Card Details**: View full card information with markdown rendering
- 🔍 **Field Extraction**: Extract specific fields (title, description, assignees, etc.) for scripting
- 🎯 **Smart Filtering**: Filter cards by lists with case-insensitive matching
- 🌈 **Color Support**: ANSI color output with CLICOLOR_FORCE support for piping
- 🔐 **Secure Configuration**: Store API credentials securely in user config directory

## Installation

### Prerequisites

- Go 1.21 or later
- Trello API Key and Token

### Build from Source

```bash
git clone <repository-url>
cd trello_cli
go build -o trello_cli .
```

### Get Trello API Credentials

1. Visit [Trello Developer API Keys](https://trello.com/app-key)
2. Copy your **API Key**
3. Click "Token" to generate an **API Token**
4. Keep these credentials secure - they'll be stored in your config file

## Configuration

### First-Time Setup

The application will guide you through initial configuration:

```bash
./trello_cli
```

This will prompt you to:
1. Enter your Trello API Key
2. Enter your Trello API Token
3. Select your workspace/organization
4. Select your board

### Configuration File Location

Configuration is stored in: `~/.config/trello_cli/config.json`

```json
{
  "api_key": "your-api-key-here",
  "api_token": "your-api-token-here",
  "workspace": "workspace-id",
  "board_id": "board-id"
}
```

### Manual Configuration

You can also manually create the config file:

```bash
mkdir -p ~/.config/trello_cli
cat > ~/.config/trello_cli/config.json << EOF
{
  "api_key": "your-api-key",
  "api_token": "your-api-token",
  "workspace": "your-workspace-id",
  "board_id": "your-board-id"
}
EOF
```

## Usage

### Basic Commands

```bash
# Show cards assigned to current user (default)
./trello_cli

# Show all cards on the board
./trello_cli --all

# Show cards from specific lists
./trello_cli --lists "In Progress,Review"
```

### Card Details

```bash
# View full card details with markdown rendering
./trello_cli --card 123
./trello_cli -c 123

# View card details with color preservation when piping
CLICOLOR_FORCE=1 ./trello_cli -c 123 | less
```

### Field Extraction

Extract specific fields for scripting and automation:

```bash
# Get just the title
./trello_cli -c 123 -f title

# Get just the description
./trello_cli -c 123 -f description

# Get assignees (full names)
./trello_cli -c 123 -f assignees

# Get labels
./trello_cli -c 123 -f labels

# Get list name
./trello_cli -c 123 -f list

# Get status (Open/Closed)
./trello_cli -c 123 -f status
```

## Command Line Options

### Main Options

| Flag | Short | Description |
|------|-------|-------------|
| `--assigned` | `-a` | Show only cards assigned to current user (default) |
| `--all` | `-A` | Show all cards on the board |
| `--lists <lists>` | `-l <lists>` | Filter cards by specific lists (comma-separated) |
| `--card <id>` | `-c <id>` | Show detailed information for a specific card |
| `--field <field>` | `-f <field>` | Extract specific field from card (use with -c) |

### Field Options (use with `-f`)

- `title` - Card title
- `description` - Card description
- `assignees` - Comma-separated list of assignee full names
- `labels` - Comma-separated list of label names
- `list` - Name of the list/column the card is in
- `status` - Card status (Open/Closed)

### Examples

```bash
# Basic usage
./trello_cli                    # Show assigned cards
./trello_cli --all              # Show all cards
./trello_cli -a                 # Short form for assigned

# List filtering
./trello_cli -l "In Progress"    # Cards in "In Progress" list
./trello_cli --lists "To Do,Done" # Multiple lists
./trello_cli --all -l "Backlog" # All cards in Backlog

# Card details
./trello_cli --card #123        # Card with ID 123
./trello_cli -c 456             # Same as above
./trello_cli -c #789            # With # prefix

# Field extraction
./trello_cli -c 123 -f title           # Just the title
./trello_cli -c 123 -f description     # Just the description
./trello_cli -c 123 -f assignees       # Just assignees
./trello_cli -c 123 -f labels          # Just labels
./trello_cli -c 123 -f list            # Just list name
./trello_cli -c 123 -f status          # Just status

# Color preservation
CLICOLOR_FORCE=1 ./trello_cli -c 123 | cat    # Force colors
CLICOLOR_FORCE=1 ./trello_cli -c 456 | less   # Colors in pager
```

## Output Format

### Task List Output

```
#123 Task Title Here              In Progress
#456 Another Task                 To Do
#789 Final Task                   Done
```

- Cards are sorted by ID (ascending)
- Table uses fixed-width columns for proper alignment
- ID column fixed to 8 characters with proper padding
- Title column fixed to 80 characters with truncation when necessary
- Shows list name for each card

### Card Details Output

When viewing card details (`-c` flag), the output includes:

- **Title**: Formatted as heading
- **Status**: Open/Closed badge
- **Description**: Full description text
- **Assignees**: Full names (not IDs)
- **Labels**: All card labels
- **List**: Which column the card is in
- **Comments**: Recent comments with author and timestamp
- **Links**: Direct link to card on Trello

### Field Output

When using `-f` flag, only the raw field value is returned:

```
Implement User Authentication
```

```
John Smith, Jane Doe
```

```
Open
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CLICOLOR_FORCE=1` | Force ANSI color output even when piping |

## Dependencies

- [BubbleTea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## API Integration

This application uses the Trello REST API:

- **Base URL**: `https://api.trello.com/1`
- **Authentication**: API Key + Token
- **Endpoints Used**:
  - `GET /members/me` - Get current user
  - `GET /members/me/organizations` - List organizations
  - `GET /organizations/{id}/boards` - List boards
  - `GET /boards/{id}/cards` - Get board cards
  - `GET /boards/{id}/lists` - Get board lists
  - `GET /cards/{id}` - Get card details
  - `GET /cards/{id}/actions` - Get card comments
  - `GET /members/{id}` - Get member details

## Troubleshooting

### Common Issues

**"API credentials not found"**
- Run the application without flags first to set up credentials
- Check that `~/.config/trello_cli/config.json` exists and contains valid credentials

**"Card with ID #123 not found"**
- Verify the card ID exists on your selected board
- Try using just the number without the # prefix

**No colors when piping**
- Use `CLICOLOR_FORCE=1` environment variable
- Example: `CLICOLOR_FORCE=1 ./trello_cli -c 123 | less`

**"Unknown field" error**
- Valid fields: `title`, `description`, `assignees`, `labels`, `list`, `status`
- Field names are case-insensitive

### Getting Help

```bash
./trello_cli --help  # Show available options
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source. Please check the license file for details.

## Changelog

### Latest Version
- ✅ Field-specific output with `-f` flag
- ✅ Display assignee full names instead of IDs
- ✅ CLICOLOR_FORCE support for piping
- ✅ Enhanced markdown rendering
- ✅ Fixed API integration for card details
- ✅ Improved error handling