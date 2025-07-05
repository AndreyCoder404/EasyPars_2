package parser

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Parser represents the main parser structure for scraping fight data
// Contains HTTP client configuration and parsing logic
type Parser struct {
	// BaseURL stores the target URL for parsing fight data
	BaseURL string
	// HTTPClient configured with proper headers and timeouts
	HTTPClient *http.Client
	// ResultChannel for goroutine communication
	ResultChannel chan FightResult
	// WaitGroup for goroutine synchronization
	WaitGroup sync.WaitGroup
	// Counter for generating unique IDs
	idCounter int64
	// Mutex for thread-safe ID generation
	idMutex sync.Mutex
}

// FightResult represents a single fight result with error handling
type FightResult struct {
	// Fight data extracted from HTML
	Fight map[string]interface{}
	// Error if parsing failed for this fight
	Error error
}

// FightEvent represents a grouped fight event (location + fight data)
type FightEvent struct {
	Location string
	Date     string
	Fighter1 string
	Fighter2 string
	Result   string
	Row      *goquery.Selection
}

// NewParser creates a new parser instance with preconfigured HTTP client
// Initializes client with headers from the working HTTP client example
func NewParser(baseURL string) *Parser {
	// Create HTTP client with reasonable timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Initialize result channel for goroutine communication
	// Buffer size of 100 to handle multiple concurrent parsing operations
	resultChan := make(chan FightResult, 100)

	return &Parser{
		BaseURL:       baseURL,
		HTTPClient:    client,
		ResultChannel: resultChan,
		idCounter:     0,
	}
}

// generateUniqueID generates a unique ID for each fight using counter and timestamp
func (p *Parser) generateUniqueID() string {
	p.idMutex.Lock()
	defer p.idMutex.Unlock()

	p.idCounter++
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("fight_%d_%d", timestamp, p.idCounter)
}

// ParseFights parses fight data from the target website using concurrent goroutines
// Returns a slice of fight data maps with proper error handling
func (p *Parser) ParseFights() ([]interface{}, error) {
	log.Println("Starting parse of fight data from vringe.com")

	// Perform HTTP request to get the HTML content
	doc, err := p.fetchHTMLDocument()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HTML document: %w", err)
	}

	// Extract fight events from the HTML document using the correct structure
	fightEvents := p.extractFightElements(doc)
	if len(fightEvents) == 0 {
		log.Println("No fight events found in the HTML document")
		return []interface{}{}, nil
	}

	log.Printf("Found %d fight events, starting concurrent parsing", len(fightEvents))

	// Process fight events concurrently using goroutines
	fights := p.processFightsConcurrently(fightEvents)

	log.Printf("Successfully parsed %d fights", len(fights))
	return fights, nil
}

// fetchHTMLDocument performs HTTP GET request with configured headers
// Returns parsed HTML document using goquery
func (p *Parser) fetchHTMLDocument() (*goquery.Document, error) {
	// Create new HTTP request with proper method and URL
	req, err := http.NewRequest("GET", p.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set headers based on the working HTTP client example
	// These headers mimic a real browser request to avoid blocking
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "ru,ru-RU;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://www.google.com/")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("sec-ch-ua-mobile", "?1")
	req.Header.Set("sec-ch-ua-platform", "Android")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36")
	// Note: Removed If-Modified-Since header to avoid 304 responses
	// Note: Removed cookie header for initial implementation - can be added later if needed

	// Execute the HTTP request
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotModified {
			return nil, fmt.Errorf("resource not modified (304), no new data available")
		}
		return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Parse HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML document: %w", err)
	}

	return doc, nil
}

// extractFightElements finds and extracts fight events from HTML document
// Groups fights by location and properly handles the table structure
func (p *Parser) extractFightElements(doc *goquery.Document) []FightEvent {
	var fightEvents []FightEvent

	// Get current date for filling missing year/month information
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()

	// Extract month information from div.month if available
	// This helps us determine the correct year and month for dates
	monthText := ""
	doc.Find("div.month").First().Each(func(i int, s *goquery.Selection) {
		monthText = strings.TrimSpace(s.Text())
		log.Printf("Found month context: %s", monthText)
	})

	// Process each table within the document
	doc.Find("table").Each(func(tableIndex int, table *goquery.Selection) {
		currentLocation := ""

		// Process each row in the table
		table.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
			// Check if this row contains location information (class="place")
			if placeCell := row.Find("td.place"); placeCell.Length() > 0 {
				// Extract and clean location text
				locationText := strings.TrimSpace(placeCell.Text())
				currentLocation = p.cleanLocationText(locationText)
				log.Printf("Found location: %s", currentLocation)
				return
			}

			// Check if this row contains fight data (has date, boxer_1, vs, boxer_2)
			dateCell := row.Find("td.date")
			boxer1Cell := row.Find("td.boxer_1")
			vsCell := row.Find("td.vs")
			boxer2Cell := row.Find("td.boxer_2")

			// Validate that we have all required fight data elements
			if dateCell.Length() > 0 && boxer1Cell.Length() > 0 && vsCell.Length() > 0 && boxer2Cell.Length() > 0 {
				// Extract date and format it properly
				dayText := strings.TrimSpace(dateCell.Text())
				formattedDate := p.formatDate(dayText, currentYear, int(currentMonth))

				// Extract fighter names (clean HTML and extra text)
				fighter1Name := p.extractFighterName(boxer1Cell)
				fighter2Name := p.extractFighterName(boxer2Cell)

				// Extract result from vs cell
				resultText := p.extractResult(vsCell)

				// Create fight event with all extracted data
				fightEvent := FightEvent{
					Location: currentLocation,
					Date:     formattedDate,
					Fighter1: fighter1Name,
					Fighter2: fighter2Name,
					Result:   resultText,
					Row:      row,
				}

				// Validate fight event has minimum required data
				if fightEvent.Fighter1 != "" && fightEvent.Fighter2 != "" {
					fightEvents = append(fightEvents, fightEvent)
					log.Printf("Extracted fight: %s vs %s on %s in %s",
						fightEvent.Fighter1, fightEvent.Fighter2, fightEvent.Date, fightEvent.Location)
				}
			}
		})
	})

	return fightEvents
}

// cleanLocationText removes HTML comments and extra whitespace from location text
func (p *Parser) cleanLocationText(locationText string) string {
	// Remove HTML comments and extra whitespace
	cleaned := strings.TrimSpace(locationText)

	// Remove any HTML artifacts that might remain
	cleaned = regexp.MustCompile(`<!--.*?-->`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	// Handle empty or invalid location
	if cleaned == "" {
		cleaned = "Unknown Location"
	}

	return cleaned
}

// formatDate converts day number to full date format (YYYY-MM-DD)
func (p *Parser) formatDate(dayText string, year int, month int) string {
	// Parse day number from text
	dayNum, err := strconv.Atoi(strings.TrimSpace(dayText))
	if err != nil {
		log.Printf("Error parsing day number '%s': %v", dayText, err)
		return time.Now().Format("2006-01-02")
	}

	// Create date with current year and month
	// Future improvement: Parse month from HTML structure or page context
	date := time.Date(year, time.Month(month), dayNum, 0, 0, 0, 0, time.UTC)

	return date.Format("2006-01-02")
}

// extractFighterName extracts clean fighter name from HTML cell
func (p *Parser) extractFighterName(cell *goquery.Selection) string {
	// First try to get text from link (a tag)
	linkText := strings.TrimSpace(cell.Find("a").First().Text())
	if linkText != "" {
		return linkText
	}

	// If no link, get all text and clean it
	fullText := strings.TrimSpace(cell.Text())

	// Remove record information in parentheses (e.g., "(7-0, 6 KO)")
	// Split by line break and take the first line (fighter name)
	lines := strings.Split(fullText, "\n")
	if len(lines) > 0 {
		fighterName := strings.TrimSpace(lines[0])

		// Remove any remaining parentheses content
		if idx := strings.Index(fighterName, "("); idx != -1 {
			fighterName = strings.TrimSpace(fighterName[:idx])
		}

		return fighterName
	}

	return "Unknown Fighter"
}

// extractResult extracts fight result from vs cell
func (p *Parser) extractResult(cell *goquery.Selection) string {
	// Get full text from vs cell
	resultText := strings.TrimSpace(cell.Text())

	// Clean up the result text
	resultText = regexp.MustCompile(`\s+`).ReplaceAllString(resultText, " ")

	// Handle empty result
	if resultText == "" {
		resultText = "TBD"
	}

	return resultText
}

// processFightsConcurrently processes fight events using goroutines
// Uses WaitGroup for synchronization and channels for result collection
func (p *Parser) processFightsConcurrently(fightEvents []FightEvent) []interface{} {
	// Set up goroutine synchronization
	p.WaitGroup.Add(len(fightEvents))

	// Start goroutines to process each fight event
	for i, event := range fightEvents {
		go p.parseSingleFight(i, event)
	}

	// Start result collector goroutine
	results := make([]interface{}, 0, len(fightEvents))
	resultCollector := make(chan []interface{}, 1)

	go func() {
		defer close(resultCollector)
		for i := 0; i < len(fightEvents); i++ {
			select {
			case result := <-p.ResultChannel:
				if result.Error != nil {
					log.Printf("Error parsing fight: %v", result.Error)
					continue
				}
				if result.Fight != nil {
					results = append(results, result.Fight)
				}
			case <-time.After(10 * time.Second):
				log.Printf("Timeout waiting for fight result %d", i)
				break
			}
		}
		resultCollector <- results
	}()

	// Wait for all goroutines to complete
	p.WaitGroup.Wait()

	// Get collected results
	select {
	case finalResults := <-resultCollector:
		return finalResults
	case <-time.After(5 * time.Second):
		log.Println("Timeout waiting for result collection")
		return results
	}
}

// parseSingleFight parses a single fight event in a goroutine
// Includes panic recovery and error handling
func (p *Parser) parseSingleFight(index int, event FightEvent) {
	defer p.WaitGroup.Done()

	// Recover from potential panics in parsing logic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in parseSingleFight for index %d: %v", index, r)
			p.ResultChannel <- FightResult{
				Fight: nil,
				Error: fmt.Errorf("panic during parsing: %v", r),
			}
		}
	}()

	// Convert fight event to structured data
	fight, err := p.convertEventToFight(event)
	if err != nil {
		p.ResultChannel <- FightResult{
			Fight: nil,
			Error: fmt.Errorf("failed to convert event to fight: %w", err),
		}
		return
	}

	// Send successful result to channel
	p.ResultChannel <- FightResult{
		Fight: fight,
		Error: nil,
	}
}

// convertEventToFight converts a FightEvent to a structured fight data map
func (p *Parser) convertEventToFight(event FightEvent) (map[string]interface{}, error) {
	// Create fight data map with extracted information
	fight := map[string]interface{}{
		"id":        p.generateUniqueID(),
		"date":      event.Date,
		"fighter1":  event.Fighter1,
		"fighter2":  event.Fighter2,
		"result":    event.Result,
		"location":  event.Location,
		"parsed_at": time.Now().Format(time.RFC3339),
	}

	// Validate that we have essential fight data
	if event.Fighter1 == "" || event.Fighter2 == "" {
		return nil, fmt.Errorf("missing essential fight data: fighter1=%s, fighter2=%s", event.Fighter1, event.Fighter2)
	}

	return fight, nil
}

// Future functions to be implemented:
// - parseEventDetails(eventURL string) (map[string]interface{}, error)
// - extractFighterStats(fighterURL string) (map[string]interface{}, error)
// - validateFightData(fight map[string]interface{}) error
// - cacheResults(fights []interface{}) error
// - setupRateLimit() error
// - handleRetryLogic(maxRetries int) error
// - parseWithPagination(startPage, endPage int) ([]interface{}, error)
// - detectHTMLStructureChanges(doc *goquery.Document) error
// - parseMonthContext(doc *goquery.Document) (int, int, error) // Extract year and month from page context
