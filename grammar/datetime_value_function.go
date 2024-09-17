//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <datetime value function>    ::=
//          <current date value function>
//      |     <current time value function>
//      |     <current timestamp value function>
//      |     <current local time value function>
//      |     <current local timestamp value function>
//
// <current date value function>    ::=   CURRENT_DATE
//
// <current time value function>    ::=   CURRENT_TIME [ <left paren> <time precision> <right paren> ]
//
// <current local time value function>    ::=   LOCALTIME [ <left paren> <time precision> <right paren> ]
//
// <current timestamp value function>    ::=   CURRENT_TIMESTAMP [ <left paren> <timestamp precision> <right paren> ]
//
// <current local timestamp value function>    ::=   LOCALTIMESTAMP [ <left paren> <timestamp precision> <right paren> ]

type DatetimeValueFunction struct {
	CurrentDate      bool
	CurrentTime      *CurrentTimeFunction
	LocalTime        *LocalTimeFunction
	CurrentTimestamp *CurrentTimestampFunction
	LocalTimstamp    *LocalTimestampFunction
}

type CurrentTimeFunction struct {
	Precision *uint
}

type LocalTimeFunction struct {
	Precision *uint
}

type CurrentTimestampFunction struct {
	Precision *uint
}

type LocalTimestampFunction struct {
	Precision *uint
}
