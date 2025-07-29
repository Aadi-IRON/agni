package config

import (
	"fmt"
	"strings"
	"time"
)

// ProgressBar represents a simple progress bar
type ProgressBar struct {
	total       int
	current     int
	width       int
	description string
	startTime   time.Time
	lastUpdate  time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, description string) *ProgressBar {
	width := GetTerminalWidth() - 20 // Leave space for percentage and description
	if width < 20 {
		width = 20
	}

	return &ProgressBar{
		total:       total,
		current:     0,
		width:       width,
		description: description,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
	}
}

// Update updates the progress bar
func (bar *ProgressBar) Update(current int) {
	bar.current = current
	bar.lastUpdate = time.Now()
	bar.render()
}

// Increment increments the progress by 1
func (bar *ProgressBar) Increment() {
	bar.Update(bar.current + 1)
}

// SetTotal updates the total count
func (bar *ProgressBar) SetTotal(total int) {
	bar.total = total
	bar.render()
}

// render renders the progress bar
func (bar *ProgressBar) render() {
	if bar.total <= 0 {
		return
	}

	percentage := float64(bar.current) / float64(bar.total) * 100
	barWidth := int(float64(bar.width) * percentage / 100)

	// Create the progress bar
	progressBar := strings.Repeat("█", barWidth)
	empty := strings.Repeat("░", bar.width-barWidth)

	// Calculate elapsed time
	elapsed := time.Since(bar.startTime)

	// Estimate remaining time
	var remaining time.Duration
	if bar.current > 0 {
		rate := elapsed / time.Duration(bar.current)
		remaining = rate * time.Duration(bar.total-bar.current)
	}

	// Format the output
	output := fmt.Sprintf("\r%s %s%s %s %3.1f%% (%d/%d) %s remaining",
		BoldCyan+bar.description+Reset,
		Green+progressBar+Reset,
		empty,
		BoldYellow,
		percentage,
		bar.current,
		bar.total,
		formatDuration(remaining)+Reset)

	fmt.Print(output)
}

// Finish completes the progress bar
func (bar *ProgressBar) Finish() {
	bar.Update(bar.total)
	fmt.Println() // Move to next line
}

// formatDuration formats duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "<1s"
	}
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	return fmt.Sprintf("%.0fh", d.Hours())
}

// SimpleProgressBar creates a simpler progress indicator
type SimpleProgressBar struct {
	total       int
	current     int
	description string
}

// NewSimpleProgressBar creates a simple progress bar
func NewSimpleProgressBar(total int, description string) *SimpleProgressBar {
	return &SimpleProgressBar{
		total:       total,
		current:     0,
		description: description,
	}
}

// Update updates the simple progress bar
func (spb *SimpleProgressBar) Update(current int) {
	spb.current = current
	if spb.total > 0 {
		percentage := float64(spb.current) / float64(spb.total) * 100
		fmt.Printf("\r%s %s%3.1f%% (%d/%d)%s",
			BoldCyan+spb.description+Reset,
			BoldYellow,
			percentage,
			spb.current,
			spb.total,
			Reset)
	}
}

// Increment increments the progress by 1
func (spb *SimpleProgressBar) Increment() {
	spb.Update(spb.current + 1)
}

// Finish completes the simple progress bar
func (spb *SimpleProgressBar) Finish() {
	spb.Update(spb.total)
	fmt.Println()
}

// Spinner creates a simple spinning indicator
type Spinner struct {
	description string
	stopChan    chan bool
	spinner     []string
	index       int
}

// NewSpinner creates a new spinner
func NewSpinner(description string) *Spinner {
	return &Spinner{
		description: description,
		stopChan:    make(chan bool),
		spinner:     []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:       0,
	}
}

// Start starts the spinner
func (sp *Spinner) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("\r%s %s %s",
					BoldCyan+sp.description+Reset,
					BoldYellow+sp.spinner[sp.index]+Reset,
					Reset)
				sp.index = (sp.index + 1) % len(sp.spinner)
			case <-sp.stopChan:
				return
			}
		}
	}()
}

// Stop stops the spinner
func (sp *Spinner) Stop() {
	sp.stopChan <- true
	fmt.Println()
}
