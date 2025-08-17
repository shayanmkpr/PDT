package main

import (
	"bufio"
	"fmt"
	// "os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type LineChange struct {
	Type       string // "added", "removed", "unchanged"
	Content    string
	OldLineNum int // Line number in old file (0 if added)
	NewLineNum int // Line number in new file (0 if removed)
}

type FileChange struct {
	Filename    string
	OldFile     string
	NewFile     string
	LineChanges []LineChange
}

type HunkHeader struct {
	OldStart int
	OldCount int
	NewStart int
	NewCount int
}

// ParseHunkHeader parses @@ -oldStart,oldCount +newStart,newCount @@
func ParseHunkHeader(line string) (*HunkHeader, error) {
	re := regexp.MustCompile(`@@\s*-(\d+)(?:,(\d+))?\s*\+(\d+)(?:,(\d+))?\s*@@`)
	matches := re.FindStringSubmatch(line)
	
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid hunk header: %s", line)
	}
	
	oldStart, _ := strconv.Atoi(matches[1])
	newStart, _ := strconv.Atoi(matches[3])
	
	oldCount := 1
	if matches[2] != "" {
		oldCount, _ = strconv.Atoi(matches[2])
	}
	
	newCount := 1
	if matches[4] != "" {
		newCount, _ = strconv.Atoi(matches[4])
	}
	
	return &HunkHeader{
		OldStart: oldStart,
		OldCount: oldCount,
		NewStart: newStart,
		NewCount: newCount,
	}, nil
}

// ParseGitDiff parses git diff output into structured changes
func ParseGitDiff(diffOutput string) ([]FileChange, error) {
	var fileChanges []FileChange
	var currentFile *FileChange
	
	scanner := bufio.NewScanner(strings.NewReader(diffOutput))
	
	var oldLineNum, newLineNum int
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// File header: diff --git a/file b/file
		if strings.HasPrefix(line, "diff --git") {
			if currentFile != nil {
				fileChanges = append(fileChanges, *currentFile)
			}
			
			// Extract filename
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				filename := strings.TrimPrefix(parts[2], "a/")
				currentFile = &FileChange{
					Filename:    filename,
					LineChanges: []LineChange{},
				}
			}
			continue
		}
		
		// Old file: --- a/file
		if strings.HasPrefix(line, "---") {
			if currentFile != nil {
				currentFile.OldFile = strings.TrimSpace(strings.TrimPrefix(line, "---"))
			}
			continue
		}
		
		// New file: +++ b/file
		if strings.HasPrefix(line, "+++") {
			if currentFile != nil {
				currentFile.NewFile = strings.TrimSpace(strings.TrimPrefix(line, "+++"))
			}
			continue
		}
		
		// Hunk header: @@ -1,4 +1,6 @@
		if strings.HasPrefix(line, "@@") {
			hunk, err := ParseHunkHeader(line)
			if err != nil {
				continue
			}
			oldLineNum = hunk.OldStart
			newLineNum = hunk.NewStart
			continue
		}
		
		// Skip if no current file
		if currentFile == nil {
			continue
		}
		
		// Parse line changes
		if len(line) == 0 {
			continue
		}
		
		switch line[0] {
		case ' ': // Unchanged line
			content := line[1:]
			currentFile.LineChanges = append(currentFile.LineChanges, LineChange{
				Type:       "unchanged",
				Content:    content,
				OldLineNum: oldLineNum,
				NewLineNum: newLineNum,
			})
			oldLineNum++
			newLineNum++
			
		case '-': // Removed line
			content := line[1:]
			currentFile.LineChanges = append(currentFile.LineChanges, LineChange{
				Type:       "removed",
				Content:    content,
				OldLineNum: oldLineNum,
				NewLineNum: 0, // No line in new file
			})
			oldLineNum++
			
		case '+': // Added line
			content := line[1:]
			currentFile.LineChanges = append(currentFile.LineChanges, LineChange{
				Type:       "added",
				Content:    content,
				OldLineNum: 0, // No line in old file
				NewLineNum: newLineNum,
			})
			newLineNum++
		}
	}
	
	// Add last file
	if currentFile != nil {
		fileChanges = append(fileChanges, *currentFile)
	}
	
	return fileChanges, nil
}

// GetGitDiff executes git diff and returns raw output
func GetGitDiff(filename string) (string, error) {
	var cmd *exec.Cmd
	if filename != "" {
		cmd = exec.Command("git", "diff", filename)
	} else {
		cmd = exec.Command("git", "diff")
	}
	
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return string(output), nil
}

// GetAddedLines returns only added lines with their line numbers
func GetAddedLines(changes []LineChange) []LineChange {
	var added []LineChange
	for _, change := range changes {
		if change.Type == "added" {
			added = append(added, change)
		}
	}
	return added
}

// GetRemovedLines returns only removed lines with their line numbers
func GetRemovedLines(changes []LineChange) []LineChange {
	var removed []LineChange
	for _, change := range changes {
		if change.Type == "removed" {
			removed = append(removed, change)
		}
	}
	return removed
}

// GetModifiedLines returns pairs of removed/added lines that represent modifications
func GetModifiedLines(changes []LineChange) [][2]LineChange {
	var modifications [][2]LineChange
	
	// Simple heuristic: consecutive remove+add = modification
	for i := 0; i < len(changes)-1; i++ {
		if changes[i].Type == "removed" && changes[i+1].Type == "added" {
			modifications = append(modifications, [2]LineChange{changes[i], changes[i+1]})
		}
	}
	
	return modifications
}

// GetLineNumberRanges returns ranges of changed line numbers
func GetLineNumberRanges(changes []LineChange) (addedRanges, removedRanges [][2]int) {
	var currentAddedStart, currentRemovedStart int
	var currentAddedEnd, currentRemovedEnd int
	
	for _, change := range changes {
		switch change.Type {
		case "added":
			if currentAddedStart == 0 {
				currentAddedStart = change.NewLineNum
				currentAddedEnd = change.NewLineNum
			} else if change.NewLineNum == currentAddedEnd+1 {
				currentAddedEnd = change.NewLineNum
			} else {
				// Gap found, save current range and start new one
				addedRanges = append(addedRanges, [2]int{currentAddedStart, currentAddedEnd})
				currentAddedStart = change.NewLineNum
				currentAddedEnd = change.NewLineNum
			}
			
		case "removed":
			if currentRemovedStart == 0 {
				currentRemovedStart = change.OldLineNum
				currentRemovedEnd = change.OldLineNum
			} else if change.OldLineNum == currentRemovedEnd+1 {
				currentRemovedEnd = change.OldLineNum
			} else {
				// Gap found, save current range and start new one
				removedRanges = append(removedRanges, [2]int{currentRemovedStart, currentRemovedEnd})
				currentRemovedStart = change.OldLineNum
				currentRemovedEnd = change.OldLineNum
			}
		}
	}
	
	// Add final ranges
	if currentAddedStart > 0 {
		addedRanges = append(addedRanges, [2]int{currentAddedStart, currentAddedEnd})
	}
	if currentRemovedStart > 0 {
		removedRanges = append(removedRanges, [2]int{currentRemovedStart, currentRemovedEnd})
	}
	
	return addedRanges, removedRanges
}

// func main() {
// 	filename := ""
// 	if len(os.Args) > 1 {
// 		filename = os.Args[1]
// 	}
//
// 	// Get git diff
// 	diffOutput, err := GetGitDiff(filename)
// 	if err != nil {
// 		fmt.Printf("Error getting diff: %v\n", err)
// 		return
// 	}
//
// 	if diffOutput == "" {
// 		fmt.Println("No changes found")
// 		return
// 	}
//
// 	// Parse diff
// 	fileChanges, err := ParseGitDiff(diffOutput)
// 	if err != nil {
// 		fmt.Printf("Error parsing diff: %v\n", err)
// 		return
// 	}
//
// 	// Process each file
// 	for _, fileChange := range fileChanges {
// 		if !strings.HasSuffix(fileChange.Filename, ".json"){
// 			continue
// 		}
//
// 		fmt.Printf("\n=== File: %s ===\n", fileChange.Filename)
//
// 		// Get added lines
// 		added := GetAddedLines(fileChange.LineChanges)
// 		fmt.Printf("\nAdded lines (%d):\n", len(added))
// 		for _, line := range added {
// 			fmt.Printf("  +%d: %s\n", line.NewLineNum, line.Content)
// 		}
//
// 		// Get removed lines
// 		removed := GetRemovedLines(fileChange.LineChanges)
// 		fmt.Printf("\nRemoved lines (%d):\n", len(removed))
// 		for _, line := range removed {
// 			fmt.Printf("  -%d: %s\n", line.OldLineNum, line.Content)
// 		}
//
// 		// Get modified lines
// 		modifications := GetModifiedLines(fileChange.LineChanges)
// 		fmt.Printf("\nModified lines (%d pairs):\n", len(modifications))
// 		for _, mod := range modifications {
// 			fmt.Printf("  -%d: %s\n", mod[0].OldLineNum, mod[0].Content)
// 			fmt.Printf("  +%d: %s\n", mod[1].NewLineNum, mod[1].Content)
// 			fmt.Println()
// 		}
//
// 		// Get line number ranges
// 		addedRanges, removedRanges := GetLineNumberRanges(fileChange.LineChanges)
// 		fmt.Printf("\nAdded line ranges: %v\n", addedRanges)
// 		fmt.Printf("Removed line ranges: %v\n", removedRanges)
//
// 		// Example: Do something with the changed lines
// 		fmt.Printf("\n=== Processing Changes ===\n")
// 		for _, change := range fileChange.LineChanges {
// 			switch change.Type {
// 			case "added":
// 				// Your logic for added lines
// 				fmt.Printf("Process added line %d: %s\n", change.NewLineNum, change.Content)
// 			case "removed":
// 				// Your logic for removed lines
// 				fmt.Printf("Process removed line %d: %s\n", change.OldLineNum, change.Content)
// 			case "unchanged":
// 				// Your logic for unchanged lines (if needed)
// 				// fmt.Printf("Unchanged line %d->%d: %s\n", change.OldLineNum, change.NewLineNum, change.Content)
// 			}
// 		}
// 	}
// }
