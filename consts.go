package main

import "time"

// Message constants
const (
	CTimeLayout = time.RFC3339
	CTimeout    = 2000 // in milliseconds

	// Network messages
	MMessageIDDataSeparator = ";;/data/;;"
	MNewUser                = "MESSAGE_NEW_USER" + MMessageIDDataSeparator
	MBoard                  = "MESSAGE_BOARD" + MMessageIDDataSeparator
)
