package main

import "time"

// Message constants
const (
	CTimeLayout = time.RFC3339
	CTimeout    = 1500 // in milliseconds

	MMessageIDDataSeparator = ";;/data/;;"
	MNewUser                = "MESSAGE_NEW_USER" + MMessageIDDataSeparator
	MBoard                  = "MESSAGE_BOARD" + MMessageIDDataSeparator
	MBoardF                 = "MESSAGE_BOARD_FORCED" + MMessageIDDataSeparator
)
