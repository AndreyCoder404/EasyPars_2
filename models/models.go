package models

// Fight represents a fight record
// Future steps: Add GORM tags, validation tags, and additional fields
type Fight struct {
	// Basic fields
	ID       uint   `json:"id"`
	Date     string `json:"date"`
	Fighter1 string `json:"fighter1"`
	Fighter2 string `json:"fighter2"`
	Result   string `json:"result"`
	Location string `json:"location"`

	// Future fields to be added:
	// ID          uint      `json:"id" gorm:"primaryKey"`
	// Date        time.Time `json:"date" gorm:"not null"`
	// Fighter1ID  uint      `json:"fighter1_id" gorm:"not null"`
	// Fighter2ID  uint      `json:"fighter2_id" gorm:"not null"`
	// ResultType  string    `json:"result_type"` // KO, TKO, Decision, etc.
	// Round       int       `json:"round"`
	// Time        string    `json:"time"`
	// Weight      float64   `json:"weight"`
	// Title       string    `json:"title"`
	// CreatedAt   time.Time `json:"created_at"`
	// UpdatedAt   time.Time `json:"updated_at"`
	// DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

// Fighter represents a fighter record
// Future steps: Add comprehensive fighter information
type Fighter struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`

	// Future fields to be added:
	// ID          uint      `json:"id" gorm:"primaryKey"`
	// Name        string    `json:"name" gorm:"not null"`
	// Nickname    string    `json:"nickname"`
	// Weight      float64   `json:"weight"`
	// Height      float64   `json:"height"`
	// Reach       float64   `json:"reach"`
	// Wins        int       `json:"wins"`
	// Losses      int       `json:"losses"`
	// Draws       int       `json:"draws"`
	// Country     string    `json:"country"`
	// BirthDate   time.Time `json:"birth_date"`
	// CreatedAt   time.Time `json:"created_at"`
	// UpdatedAt   time.Time `json:"updated_at"`
	// DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

// Future models to be implemented:
// - User (for authentication)
// - Event (for fight events)
// - WeightClass
// - Organization (UFC, Bellator, etc.)
// - Venue
