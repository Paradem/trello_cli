package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"trello_cli/config"
	"trello_cli/trello"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func showCardDetails(cardID int, fieldFilter string) {
	// Load existing config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Check if we have API credentials
	if cfg.APIKey == "" || cfg.APIToken == "" {
		log.Fatalf("API credentials not found. Please run without --card flag first to set up credentials.")
	}

	// Create Trello client
	client := trello.NewClient(cfg.APIKey, cfg.APIToken)

	// First, get all cards to find the one with matching ShortID
	cards, err := client.GetCards(cfg.BoardID)
	if err != nil {
		log.Fatalf("Failed to get cards: %v", err)
	}

	// Find the card with the matching ShortID
	var targetCard *trello.Card
	for _, card := range cards {
		if card.IDShort == cardID {
			targetCard = &card
			break
		}
	}

	if targetCard == nil {
		log.Fatalf("Card with ID #%d not found on this board", cardID)
	}

	// Get card details using the full card ID
	detailedCard, err := client.GetCardDetails(targetCard.ID)
	if err != nil {
		log.Fatalf("Failed to get card details: %v", err)
	}

	// Get comments using the full card ID
	comments, err := client.GetCardComments(targetCard.ID)
	if err != nil {
		log.Fatalf("Failed to get card comments: %v", err)
	}

	// Get lists for list name lookup
	lists, err := client.GetLists(cfg.BoardID)
	if err != nil {
		log.Fatalf("Failed to get lists: %v", err)
	}

	listMap := make(map[string]string)
	for _, list := range lists {
		listMap[list.ID] = list.Name
	}

	// Build markdown content
	var markdown strings.Builder

	// Title
	markdown.WriteString(fmt.Sprintf("# %s\n\n", detailedCard.Name))

	// Status
	status := "Open"
	if detailedCard.Closed {
		status = "Closed"
	}
	markdown.WriteString(fmt.Sprintf("**%s**\n\n", status))

	// Description
	if detailedCard.Desc != "" {
		markdown.WriteString(fmt.Sprintf("## Description\n\n%s\n\n", detailedCard.Desc))
	}

	// Assignees
	if len(detailedCard.IDMembers) > 0 {
		markdown.WriteString("## Assignees\n")
		for _, memberID := range detailedCard.IDMembers {
			// Look up member details to get full name
			member, err := client.GetMember(memberID)
			if err != nil {
				// If lookup fails, show the ID
				markdown.WriteString(fmt.Sprintf("- %s\n", memberID))
			} else {
				markdown.WriteString(fmt.Sprintf("- %s\n", member.FullName))
			}
		}
		markdown.WriteString("\n")
	}

	// Labels
	if len(detailedCard.Labels) > 0 {
		markdown.WriteString("## Labels\n")
		for _, label := range detailedCard.Labels {
			markdown.WriteString(fmt.Sprintf("- %s\n", label.Name))
		}
		markdown.WriteString("\n")
	}

	// List
	if listName, exists := listMap[detailedCard.IDList]; exists {
		markdown.WriteString(fmt.Sprintf("## List\n\n%s\n\n", listName))
	}

	// Comments
	if len(comments) > 0 {
		markdown.WriteString(fmt.Sprintf("## Comments (%d)\n\n", len(comments)))
		for i, comment := range comments {
			commentTime, err := time.Parse(time.RFC3339, comment.Date)
			timeStr := "Unknown time"
			if err == nil {
				timeStr = commentTime.Format("Jan 2, 2006 at 3:04 PM")
			}

			markdown.WriteString(fmt.Sprintf("### Comment %d\n\n", i+1))
			markdown.WriteString(fmt.Sprintf("**%s** commented on %s:\n\n", comment.MemberCreator.FullName, timeStr))
			markdown.WriteString(fmt.Sprintf("%s\n\n", comment.Data.Text))
			markdown.WriteString("---\n\n")
		}
	}

	// Card link
	markdown.WriteString("## Links\n\n")
	markdown.WriteString(fmt.Sprintf("- View this card on Trello: https://trello.com/c/%s\n", detailedCard.ShortLink))

	// Handle field filtering - if fieldFilter is specified, output only that field
	if fieldFilter != "" {
		switch strings.ToLower(fieldFilter) {
		case "title":
			fmt.Print(detailedCard.Name)
		case "description":
			if detailedCard.Desc != "" {
				fmt.Print(detailedCard.Desc)
			}
		case "status":
			if detailedCard.Closed {
				fmt.Print("Closed")
			} else {
				fmt.Print("Open")
			}
		case "assignees":
			if len(detailedCard.IDMembers) > 0 {
				var names []string
				for _, memberID := range detailedCard.IDMembers {
					member, err := client.GetMember(memberID)
					if err != nil {
						names = append(names, memberID)
					} else {
						names = append(names, member.FullName)
					}
				}
				fmt.Print(strings.Join(names, ", "))
			}
		case "labels":
			if len(detailedCard.Labels) > 0 {
				var labelNames []string
				for _, label := range detailedCard.Labels {
					labelNames = append(labelNames, label.Name)
				}
				fmt.Print(strings.Join(labelNames, ", "))
			}
		case "list":
			if listName, exists := listMap[detailedCard.IDList]; exists {
				fmt.Print(listName)
			} else {
				fmt.Print("Unknown")
			}
		default:
			log.Fatalf("Unknown field: %s. Available fields: title, description, status, assignees, labels, list", fieldFilter)
		}
		return
	}

	// Render markdown with glamour
	var out string

	if os.Getenv("CLICOLOR_FORCE") == "1" {
		// Force ANSI output even when piping
		os.Setenv("NO_COLOR", "")
		os.Setenv("COLORTERM", "256color")
		os.Setenv("TERM", "xterm-256color")
		// Use Render function with dark style
		var err error
		out, err = glamour.Render(markdown.String(), "dark")
		if err != nil {
			log.Fatalf("Failed to render markdown: %v", err)
		}
	} else {
		// Use auto-style for adaptive coloring
		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(0),
		)
		var err error
		out, err = r.Render(markdown.String())
		if err != nil {
			log.Fatalf("Failed to render markdown: %v", err)
		}
	}

	fmt.Print(out)
}

func main() {
	// Define CLI flags
	assignedOnly := flag.Bool("assigned", true, "Show only cards assigned to current user")
	allCards := flag.Bool("all", false, "Show all cards on the board")
	listFilter := flag.String("lists", "", "Filter cards by specific lists (comma-separated)")
	showCard := flag.String("card", "", "Show detailed information for a specific card by ID (format: #123 or 123)")
	fieldFilter := flag.String("field", "", "Show only specific field from card (use with -c): title, description, assignees, labels, list, status")
	flag.BoolVar(assignedOnly, "a", true, "Show only cards assigned to current user (short)")
	flag.BoolVar(allCards, "A", false, "Show all cards on the board (short)")
	flag.StringVar(listFilter, "l", "", "Filter cards by specific lists (comma-separated, short)")
	flag.StringVar(showCard, "c", "", "Show detailed information for a specific card by ID (format: #123 or 123, short)")
	flag.StringVar(fieldFilter, "f", "", "Show only specific field from card (use with -c): title, description, assignees, labels, list, status (short)")
	flag.Parse()

	// Validate flags - if both are set, prefer --all
	if *assignedOnly && *allCards {
		fmt.Println("Warning: Both --assigned and --all flags specified. Using --all.")
	}

	// Handle card detail view
	if *showCard != "" {
		// Parse card ID from format #123 or 123 (bare integer)
		idStr := *showCard

		// Remove # prefix if present
		if strings.HasPrefix(idStr, "#") {
			idStr = strings.TrimPrefix(idStr, "#")
		}

		// Check if we have a valid string after removing prefix
		if idStr == "" {
			log.Fatalf("Invalid card ID format. Use format: #123 or 123")
		}

		// Parse to integer
		cardID, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("Invalid card ID: %s (must be numeric)", idStr)
		}

		showCardDetails(cardID, *fieldFilter)
		return
	}

	// Parse list filter
	var allowedLists map[string]bool
	if *listFilter != "" {
		allowedLists = make(map[string]bool)
		listNames := strings.Split(*listFilter, ",")
		for _, name := range listNames {
			// Trim whitespace and convert to lowercase for case-insensitive matching
			trimmedName := strings.TrimSpace(name)
			if trimmedName != "" {
				allowedLists[strings.ToLower(trimmedName)] = true
			}
		}
	}

	// Load existing config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// If API credentials are missing, prompt for them
	if cfg.APIKey == "" || cfg.APIToken == "" {
		fmt.Println("Please provide your Trello API credentials:")
		apiKey, apiToken, _, _, err := PromptForConfig()
		if err != nil {
			log.Fatalf("Failed to get API credentials: %v", err)
		}

		cfg.APIKey = apiKey
		cfg.APIToken = apiToken

		// Test the credentials by creating a client and fetching user info
		testClient := trello.NewClient(cfg.APIKey, cfg.APIToken)
		_, err = testClient.GetMemberID()
		if err != nil {
			log.Fatalf("Invalid API credentials: %v", err)
		}
	}

	// Create Trello client
	client := trello.NewClient(cfg.APIKey, cfg.APIToken)

	// If workspace or board is missing, prompt for selection
	if cfg.Workspace == "" || cfg.BoardID == "" {
		fmt.Println("Fetching available workspaces...")

		workspaceID, err := PromptForOrganization(client)
		if err != nil {
			log.Fatalf("Failed to select workspace: %v", err)
		}

		fmt.Println("Fetching available boards...")

		boardID, err := PromptForBoard(client, workspaceID)
		if err != nil {
			log.Fatalf("Failed to select board: %v", err)
		}

		cfg.Workspace = workspaceID
		cfg.BoardID = boardID

		// Save the config
		if err := config.SaveConfig(cfg); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}
	}

	// Get current user ID
	userID, err := client.GetMemberID()
	if err != nil {
		log.Fatalf("Failed to get user ID: %v", err)
	}

	// Get cards from the board
	cards, err := client.GetCards(cfg.BoardID)
	if err != nil {
		log.Fatalf("Failed to get cards: %v", err)
	}

	// Get lists from the board for lookup
	lists, err := client.GetLists(cfg.BoardID)
	if err != nil {
		log.Fatalf("Failed to get lists: %v", err)
	}

	// Create a map of list ID to list name for quick lookup
	listMap := make(map[string]string)
	for _, list := range lists {
		listMap[list.ID] = list.Name
	}

	// Helper function to check if a list should be included
	shouldIncludeList := func(listName string) bool {
		if allowedLists == nil {
			return true // No filter specified, include all
		}
		return allowedLists[strings.ToLower(listName)]
	}

	// Collect cards to display based on filtering
	var cardsToDisplay []struct {
		id       int
		name     string
		listName string
	}

	if *allCards {
		// Show all cards
		for _, card := range cards {
			listName := listMap[card.IDList]
			if listName == "" {
				listName = "Unknown"
			}
			if shouldIncludeList(listName) {
				cardsToDisplay = append(cardsToDisplay, struct {
					id       int
					name     string
					listName string
				}{card.IDShort, card.Name, listName})
			}
		}
	} else {
		// Show only cards assigned to current user (default behavior)
		for _, card := range cards {
			for _, memberID := range card.IDMembers {
				if memberID == userID {
					listName := listMap[card.IDList]
					if listName == "" {
						listName = "Unknown"
					}
					if shouldIncludeList(listName) {
						cardsToDisplay = append(cardsToDisplay, struct {
							id       int
							name     string
							listName string
						}{card.IDShort, card.Name, listName})
					}
					break
				}
			}
		}
	}

	// Sort cards by ShortID (ascending)
	sort.Slice(cardsToDisplay, func(i, j int) bool {
		return cardsToDisplay[i].id < cardsToDisplay[j].id
	})

	// Calculate column widths
	maxIDWidth := 0
	for _, card := range cardsToDisplay {
		idStr := fmt.Sprintf("#%d", card.id)
		if len(idStr) > maxIDWidth {
			maxIDWidth = len(idStr)
		}
	}

	// Print cards with fixed column widths
	for _, card := range cardsToDisplay {
		idStr := fmt.Sprintf("#%d", card.id)
		// Format: ID (fixed width) + Title + List (right-aligned)
		styledID := lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // Yellow color
		fmt.Printf("%-*s %s\n", maxIDWidth+1, styledID.Render(idStr), card.name)
	}
}
