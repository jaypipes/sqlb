//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package symbol

// Special character symbols in order of matching precedence
const (
	SymbolANSI2003SpecialCharacterStart Symbol = iota
	SymbolSpace
	SymbolDoubleQuote
	SymbolPercent
	SymbolAmpersand
	SymbolQuote
	SymbolLeftParen
	SymbolRightParen
	SymbolAsterisk
	SymbolPlusSign
	SymbolComma
	SymbolMinusSign
	SymbolPeriod
	SymbolSolidus
	SymbolColon
	SymbolSemicolon
	SymbolLessThanOperator
	SymbolEqualsOperator
	SymbolGreaterThanOperator
	SymbolQuestionMark
	SymbolLeftBracket
	SymbolRightBracket
	SymbolCircumflex
	SymbolUnderscore
	SymbolVerticalBar
	SymbolLeftBrace
	SymbolRightBrace
	SymbolANSI2003SpecialCharacterEnd = 100
)

const (
	Space               = " "
	DoubleQuote         = "\""
	Percent             = "%"
	Ampersand           = "&"
	Quote               = "'"
	LeftParen           = "("
	RightParen          = ")"
	Asterisk            = "*"
	PlusSign            = "+"
	Comma               = ","
	MinusSign           = "-"
	Period              = "."
	Solidus             = "/"
	Colon               = ":"
	Semicolon           = ";"
	LessThanOperator    = "<"
	EqualsOperator      = "="
	GreaterThanOperator = ">"
	QuestionMark        = "?"
	LeftBracket         = "["
	RightBracket        = "]"
	Circumflex          = "^"
	Underscore          = "_"
	VerticalBar         = "|"
	LeftBrace           = "{"
	RightBrace          = "}"
)

// Reserved words in lexicographical order
const (
	SymbolANSI2003ReservedStart Symbol = SymbolANSI2003SpecialCharacterEnd + 1
	SymbolAdd
	SymbolAll
	SymbolAllocation
	SymbolAlter
	SymbolAnd
	SymbolAny
	SymbolAre
	SymbolArray
	SymbolAs
	SymbolAsensitive
	SymbolAsymmetric
	SymbolAt
	SymbolAtomic
	SymbolAuthorization
	SymbolBegin
	SymbolBetween
	SymbolBigint
	SymbolBinary
	SymbolBlob
	SymbolBoolean
	SymbolBoth
	SymbolBy
	SymbolCall
	SymbolCalled
	SymbolCascaded
	SymbolCase
	SymbolCast
	SymbolChar
	SymbolCharacter
	SymbolCheck
	SymbolClob
	SymbolClose
	SymbolCollate
	SymbolColumn
	SymbolCommit
	SymbolConnect
	SymbolConstraint
	SymbolContinue
	SymbolCorresponding
	SymbolCreate
	SymbolCross
	SymbolCube
	SymbolCurrent
	SymbolCurrentDate
	SymbolCurrentDefaultTransformGroup
	SymbolCurrentPath
	SymbolCurrentRole
	SymbolCurrentTime
	SymbolCurrentTimestamp
	SymbolCurrentTransformGroupForType
	SymbolCurrentUser
	SymbolCursor
	SymbolCycle
	SymbolDate
	SymbolDay
	SymbolDeallocate
	SymbolDec
	SymbolDecimal
	SymbolDeclare
	SymbolDefault
	SymbolDelete
	SymbolDeref
	SymbolDescribe
	SymbolDeterministic
	SymbolDisconnect
	SymbolDistinct
	SymbolDouble
	SymbolDrop
	SymbolDynamic
	SymbolEach
	SymbolElement
	SymbolElse
	SymbolEnd
	SymbolEndExec
	SymbolEscape
	SymbolExcept
	SymbolExec
	SymbolExecute
	SymbolExists
	SymbolExternal
	SymbolFalse
	SymbolFetch
	SymbolFilter
	SymbolFloat
	SymbolFor
	SymbolForeign
	SymbolFree
	SymbolFrom
	SymbolFull
	SymbolFunction
	SymbolGet
	SymbolGlobal
	SymbolGrant
	SymbolGroup
	SymbolGrouping
	SymbolHaving
	SymbolHold
	SymbolHour
	SymbolIdentity
	SymbolImmediate
	SymbolIn
	SymbolIndicator
	SymbolInner
	SymbolInOut
	SymbolInput
	SymbolInsensitive
	SymbolInsert
	SymbolInt
	SymbolInteger
	SymbolIntersect
	SymbolInterval
	SymbolInto
	SymbolIs
	SymbolIsolation
	SymbolJoin
	SymbolLanguage
	SymbolLarge
	SymbolLateral
	SymbolLeading
	SymbolLeft
	SymbolLike
	SymbolLocal
	SymbolLocalTime
	SymbolLocalTimestamp
	SymbolMatch
	SymbolMember
	SymbolMerge
	SymbolMethod
	SymbolMinute
	SymbolModifies
	SymbolModule
	SymbolMonth
	SymbolMultiSet
	SymbolNational
	SymbolNatural
	SymbolNChar
	SymbolNClob
	SymbolNew
	SymbolNo
	SymbolNone
	SymbolNot
	SymbolNull
	SymbolNumeric
	SymbolOf
	SymbolOld
	SymbolOn
	SymbolOnly
	SymbolOpen
	SymbolOr
	SymbolOrder
	SymbolOut
	SymbolOuter
	SymbolOutput
	SymbolOver
	SymbolOverlaps
	SymbolParameter
	SymbolPartition
	SymbolPrecision
	SymbolPrepare
	SymbolPrimary
	SymbolProcedure
	SymbolRange
	SymbolReads
	SymbolReal
	SymbolRecursive
	SymbolRef
	SymbolReferences
	SymbolReferencing
	SymbolRegrAvgX
	SymbolRegrAvgY
	SymbolRegrCount
	SymbolRegrIntercept
	SymbolRegrR2
	SymbolRegrSlope
	SymbolRegrSXX
	SymbolRegrSXY
	SymbolRegrSYY
	SymbolRelease
	SymbolResult
	SymbolReturn
	SymbolReturns
	SymbolRevoke
	SymbolRight
	SymbolRollback
	SymbolRollup
	SymbolRow
	SymbolRows
	SymbolSavepoint
	SymbolScroll
	SymbolSearch
	SymbolSecond
	SymbolSelect
	SymbolSensitive
	SymbolSessionUser
	SymbolSet
	SymbolSimilar
	SymbolSmallint
	SymbolSom
	SymbolSpecific
	SymbolSpecificType
	SymbolSQL
	SymbolSQLException
	SymbolSQLState
	SymbolSQLWarning
	SymbolStart
	SymbolStatic
	SymbolSubMultiSet
	SymbolSymmetric
	SymbolSystem
	SymbolSystemUser
	SymbolTable
	SymbolThen
	SymbolTime
	SymbolTimestamp
	SymbolTimezoneHour
	SymbolTimezoneMinute
	SymbolTo
	SymbolTrailing
	SymbolTranslation
	SymbolTreat
	SymbolTrigger
	SymbolTrue
	SymbolUEscape
	SymbolUnion
	SymbolUnique
	SymbolUnknown
	SymbolUnnest
	SymbolUpdate
	SymbolUpper
	SymbolUser
	SymbolUsing
	SymbolValue
	SymbolValues
	SymbolVarPop
	SymbolVarSamp
	SymbolVarchar
	SymbolVarying
	SymbolWhen
	SymbolWhenever
	SymbolWhere
	SymbolWidthBucket
	SymbolWindow
	SymbolWith
	SymbolWithin
	SymbolWithout
	SymbolYear
	SymbolANSI2003ReservedEnd = 10000
)

const (
	Add                          = "ADD"
	All                          = "ALL"
	Allocation                   = "ALLOCATION"
	Alter                        = "ALTER"
	And                          = "AND"
	Any                          = "ANY"
	Are                          = "ARE"
	Array                        = "ARRAY"
	As                           = "AS"
	Asensitive                   = "ASENSITIVE"
	Asymmetric                   = "ASYMMETRIC"
	At                           = "AT"
	Atomic                       = "ATOMIC"
	Authorization                = "AUTHORIZATION"
	Begin                        = "BEGIN"
	Between                      = "BETWEEN"
	Bigint                       = "BIGINT"
	Binary                       = "BINARY"
	Blob                         = "BLOB"
	Boolean                      = "BOOLEAN"
	Both                         = "BOTH"
	By                           = "BY"
	Call                         = "CALL"
	Called                       = "CALLED"
	Cascaded                     = "CASCADED"
	Case                         = "CASE"
	Cast                         = "CAST"
	Char                         = "CHAR"
	Character                    = "CHARACTER"
	Check                        = "CHECK"
	Clob                         = "CLOB"
	Close                        = "CLOSE"
	Collate                      = "COLLATE"
	Column                       = "COLUMN"
	Commit                       = "COMMIT"
	Connect                      = "CONNECT"
	Constraint                   = "CONSTRAINT"
	Continue                     = "CONTINUE"
	Corresponding                = "CORRESPONDING"
	Create                       = "CREATE"
	Cross                        = "CROSS"
	Cube                         = "CUBE"
	Current                      = "CURRENT"
	CurrentDate                  = "CURRENT_DATE"
	CurrentDefaultTransformGroup = "CURRENT_DEFAULT_TRANSFORM_GROUP"
	CurrentPath                  = "CURRENT_PATH"
	CurrentRole                  = "CURRENT_ROLE"
	CurrentTime                  = "CURRENT_TIME"
	CurrentTimestamp             = "CURRENT_TIMESTAMP"
	CurrentTransformGroupForType = "CURRENT_TRANSFORM_GROUP_FOR_TYPE"
	CurrentUser                  = "CURRENT_USER"
	Cursor                       = "CURSOR"
	Cycle                        = "CYCLE"
	Date                         = "DATE"
	Day                          = "DAY"
	Deallocate                   = "DEALLOCATE"
	Dec                          = "DEC"
	Decimal                      = "DECIMAL"
	Declare                      = "DECLARE"
	Default                      = "DEFAULT"
	Delete                       = "DELETE"
	Deref                        = "DEREF"
	Describe                     = "DESCRIBE"
	Deterministic                = "DETERMINISTIC"
	Disconnect                   = "DISCONNECT"
	Distinct                     = "DISTINCT"
	Double                       = "DOUBLE"
	Drop                         = "DROP"
	Dynamic                      = "DYNAMIC"
	Each                         = "EACH"
	Element                      = "ELEMENT"
	Else                         = "ELSE"
	End                          = "END"
	EndExec                      = "END-EXEC"
	Escape                       = "ESCAPE"
	Except                       = "EXCEPT"
	Exec                         = "EXEC"
	Execute                      = "EXECUTE"
	Exists                       = "EXISTS"
	External                     = "EXTERNAL"
	False                        = "FALSE"
	Fetch                        = "FETCH"
	Filter                       = "FILTER"
	Float                        = "FLOAT"
	For                          = "FOR"
	Foreign                      = "FOREIGN"
	Free                         = "FREE"
	From                         = "FROM"
	Full                         = "FULL"
	Function                     = "FUNCTION"
	Get                          = "GET"
	Global                       = "GLOBAL"
	Grant                        = "GRANT"
	Group                        = "GROUP"
	Grouping                     = "GROUPING"
	Having                       = "HAVING"
	Hold                         = "HOLD"
	Hour                         = "HOUR"
	Identity                     = "IDENTITY"
	Immediate                    = "IMMEDIATE"
	In                           = "IN"
	Indicator                    = "INDICATOR"
	Inner                        = "INNER"
	InOut                        = "INOUT"
	Input                        = "INPUT"
	Insensitive                  = "INSENSITIVE"
	Insert                       = "INSERT"
	Int                          = "INT"
	Integer                      = "INTEGER"
	Intersect                    = "INTERSECT"
	Interval                     = "INTERVAL"
	Into                         = "INTO"
	Is                           = "IS"
	Isolation                    = "ISOLATION"
	Join                         = "JOIN"
	Language                     = "LANGUAGE"
	Large                        = "LARGE"
	Lateral                      = "LATERAL"
	Leading                      = "LEADING"
	Left                         = "LEFT"
	Like                         = "LIKE"
	Local                        = "LOCAL"
	LocalTime                    = "LOCALTIME"
	LocalTimestamp               = "LOCALTIMESTAMP"
	Match                        = "MATCH"
	Member                       = "MEMBER"
	Merge                        = "MERGE"
	Method                       = "METHOD"
	Minute                       = "MINUTE"
	Modifies                     = "MODIFIES"
	Module                       = "MODULE"
	Month                        = "MONTH"
	MultiSet                     = "MULTISET"
	National                     = "NATIONAL"
	Natural                      = "NATURAL"
	NChar                        = "NCHAR"
	NClob                        = "NCLOB"
	New                          = "NEW"
	No                           = "NO"
	None                         = "NONE"
	Not                          = "NOT"
	Null                         = "NULL"
	Numeric                      = "NUMERIC"
	Of                           = "OF"
	Old                          = "OLD"
	On                           = "ON"
	Only                         = "ONLY"
	Open                         = "OPEN"
	Or                           = "OR"
	Order                        = "ORDER"
	Out                          = "OUT"
	Outer                        = "OUTER"
	Output                       = "OUTPUT"
	Over                         = "OVER"
	Overlaps                     = "OVERLAPS"
	Parameter                    = "PARAMETER"
	Partition                    = "PARTITION"
	Precision                    = "PRECISION"
	Prepare                      = "PREPARE"
	Primary                      = "PRIMARY"
	Procedure                    = "PROCEDURE"
	Range                        = "RANGE"
	Reads                        = "READS"
	Real                         = "REAL"
	Recursive                    = "RECURSIVE"
	Ref                          = "REF"
	References                   = "REFERENCES"
	Referencing                  = "REFERENCING"
	RegrAvgX                     = "REGR_AVGX"
	RegrAvgY                     = "REGR_AVGY"
	RegrCount                    = "REGR_COUNT"
	RegrIntercept                = "REGR_INTERCEPT"
	RegrR2                       = "REGR_R2"
	RegrSlope                    = "REGR_SLOPE"
	RegrSXX                      = "REGR_SXX"
	RegrSXY                      = "REGR_SXY"
	RegrSYY                      = "REGR_SYY"
	Release                      = "RELEASE"
	Result                       = "RESULT"
	Return                       = "RETURN"
	Returns                      = "RETURNS"
	Revoke                       = "REVOKE"
	Right                        = "RIGHT"
	Rollback                     = "ROLLBACK"
	Rollup                       = "ROLLUP"
	Row                          = "ROW"
	Rows                         = "ROWS"
	Savepoint                    = "SAVEPOINT"
	Scroll                       = "SCROLL"
	Search                       = "SEARCH"
	Second                       = "SECOND"
	Select                       = "SELECT"
	Sensitive                    = "SENSITIVE"
	SessionUser                  = "SESSION_USER"
	Set                          = "SET"
	Similar                      = "SIMILAR"
	Smallint                     = "SMALLINT"
	Some                         = "SOME"
	Specific                     = "SPECIFIC"
	SpecificType                 = "SPECIFICTYPE"
	SQL                          = "SQL"
	SQLException                 = "SQLEXCEPTION"
	SQLState                     = "SQLSTATE"
	SQLWarning                   = "SQLWARNING"
	Start                        = "START"
	Static                       = "STATIC"
	SubMultiSet                  = "SUBMULTISET"
	Symmetric                    = "SYMMETRIC"
	System                       = "SYSTEM"
	SystemUser                   = "SYSTEM_USER"
	Table                        = "TABLE"
	Then                         = "THEN"
	Time                         = "TIME"
	Timestamp                    = "TIMESTAMP"
	TimezoneHour                 = "TIMEZONE_HOUR"
	TimezoneMinute               = "TIMEZONE_MINUTE"
	To                           = "TO"
	Trailing                     = "TRAILING"
	Translation                  = "TRANSLATION"
	Treat                        = "TREAT"
	Trigger                      = "TRIGGER"
	True                         = "TRUE"
	UEscape                      = "UESCAPE"
	Union                        = "UNION"
	Unique                       = "UNIQUE"
	Unknown                      = "UNKNOWN"
	Unnest                       = "UNNEST"
	Update                       = "UPDATE"
	Upper                        = "UPPER"
	User                         = "USER"
	Using                        = "USING"
	Value                        = "VALUE"
	Values                       = "VALUES"
	VarPop                       = "VAR_POP"
	VarSamp                      = "VAR_SAMP"
	Varchar                      = "VARCHAR"
	Varying                      = "VARYING"
	When                         = "WHEN"
	Whenever                     = "WHENEVER"
	Where                        = "WHERE"
	WidthBucket                  = "WIDTH_BUCKET"
	Window                       = "WINDOW"
	With                         = "WITH"
	Within                       = "WITHIN"
	Without                      = "WITHOUT"
	Year                         = "YEAR"
)

// Non-reserved words in lexicographical order
const (
	SymbolANSI2003NonReservedStart Symbol = SymbolANSI2003ReservedEnd + 1
	SymbolA
	SymbolAbs
	SymbolAbsolute
	SymbolAction
	SymbolADA
	SymbolAdmin
	SymbolAfter
	SymbolAlways
	SymbolAsc
	SymbolAssertion
	SymbolAssignment
	SymbolAttribute
	SymbolAttributes
	SymbolAvg
	SymbolBefore
	SymbolBernoulli
	SymbolBreadth
	SymbolC
	SymbolCardinality
	SymbolCascade
	SymbolCatalog
	SymbolCatalogName
	SymbolCeil
	SymbolCeiling
	SymbolChain
	SymbolCharacteristics
	SymbolCharacters
	SymbolCharacterLength
	SymbolCharacterSetCatalog
	SymbolCharacterSetName
	SymbolCharacterSetSchem
	SymbolCharLength
	SymbolChecked
	SymbolClassOrigin
	SymbolCoalesce
	SymbolCOBOL
	SymbolCodeUnits
	SymbolCollation
	SymbolCollationCatalog
	SymbolCollationName
	SymbolCollationSchema
	SymbolCollect
	SymbolColumnName
	SymbolCommandFunction
	SymbolCommandFunctionCode
	SymbolCommitted
	SymbolCondition
	SymbolConditionNumber
	SymbolConnectionName
	SymbolConstraints
	SymbolConstraintCatalog
	SymbolConstraintName
	SymbolConstraintSchema
	SymbolConstructors
	SymbolContains
	SymbolConvert
	SymbolCorr
	SymbolCount
	SymbolCovarPop
	SymbolCovarSamp
	SymbolCumeDist
	SymbolCurrentCollation
	SymbolCursorName
	SymbolData
	SymbolDatetimeIntervalCode
	SymbolDatetimeIntervalPrecision
	SymbolDefaults
	SymbolDeferrable
	SymbolDeferred
	SymbolDefined
	SymbolDefiner
	SymbolDegree
	SymbolDenseRank
	SymbolDepth
	SymbolDerived
	SymbolDesc
	SymbolDescriptor
	SymbolDiagnostics
	SymbolDispatch
	SymbolDomain
	SymbolDynamicFunction
	SymbolDynamicFunctionCode
	SymbolEquals
	SymbolEvery
	SymbolException
	SymbolExclude
	SymbolExcluding
	SymbolExp
	SymbolExtract
	SymbolFinal
	SymbolFirst
	SymbolFloor
	SymbolFollowing
	SymbolFORTRAN
	SymbolFound
	SymbolFusion
	SymbolG
	SymbolGeneral
	SymbolGo
	SymbolGoto
	SymbolGranted
	SymbolHierarchy
	SymbolImplementation
	SymbolIncluding
	SymbolIncrement
	SymbolInitially
	SymbolInstance
	SymbolInstantiable
	SymbolIntersection
	SymbolInvoker
	SymbolK
	SymbolKey
	SymbolKeyMember
	SymbolKeyType
	SymbolLast
	SymbolLength
	SymbolLevel
	SymbolLn
	SymbolLocator
	SymbolLower
	SymbolM
	SymbolMap
	SymbolMatched
	SymbolMax
	SymbolMaxValue
	SymbolMessageLength
	SymbolMessageOctetLength
	SymbolMessageText
	SymbolMin
	SymbolMinValue
	SymbolMod
	SymbolMore
	SymbolMUMPS
	SymbolName
	SymbolNames
	SymbolNesting
	SymbolNext
	SymbolNormalize
	SymbolNormalized
	SymbolNullable
	SymbolNullIf
	SymbolNulls
	SymbolNumber
	SymbolObject
	SymbolOctets
	SymbolOctetLength
	SymbolOption
	SymbolOptions
	SymbolOrdering
	SymbolOrdinality
	SymbolOthers
	SymbolOverlay
	SymbolOverriding
	SymbolPad
	SymbolParameterMode
	SymbolParameterName
	SymbolParameterOrdinalPosition
	SymbolParameterSpecificCatalog
	SymbolParameterSpecificName
	SymbolParameterSpecificSchema
	SymbolPartial
	SymbolPascal
	SymbolPath
	SymbolPercentileCont
	SymbolPercentileDisc
	SymbolPercentRang
	SymbolPlacing
	SymbolPLI
	SymbolPosition
	SymbolPower
	SymbolPreceding
	SymbolPreserve
	SymbolPrior
	SymbolPrivileges
	SymbolPublic
	SymbolRank
	SymbolRead
	SymbolRelative
	SymbolRepeatable
	SymbolRestart
	SymbolReturnedCardinality
	SymbolReturnedLength
	SymbolReturnedOctetLength
	SymbolReturnedSQLState
	SymbolRole
	SymbolRoutine
	SymbolRoutineCatalog
	SymbolRoutineName
	SymbolRoutineSchema
	SymbolRowCount
	SymbolRowNumber
	SymbolScale
	SymbolSchema
	SymbolSchemaName
	SymbolScopeCatalog
	SymbolScopeName
	SymbolScopeSchema
	SymbolSection
	SymbolSecurity
	SymbolSelf
	SymbolSequence
	SymbolSerializable
	SymbolServerName
	SymbolSession
	SymbolSets
	SymbolSimple
	SymbolSize
	SymbolSource
	SymbolSpaceWord
	SymbolSpecificName
	SymbolSqrt
	SymbolState
	SymbolStatement
	SymbolStddevPop
	SymbolStddevSamp
	SymbolStructure
	SymbolStyle
	SymbolSubclassOrigin
	SymbolSubstring
	SymbolSum
	SymbolTableSample
	SymbolTableName
	SymbolTemporary
	SymbolTies
	SymbolTopLevelCount
	SymbolTransaction
	SymbolTransactionsCommitted
	SymbolTransactionsRolledBack
	SymbolTransactionActive
	SymbolTranform
	SymbolTransforms
	SymbolTranslate
	SymbolTriggerCatalog
	SymbolTriggerName
	SymbolTriggerSchema
	SymbolTrim
	SymbolType
	SymbolUnbounded
	SymbolUncommitted
	SymbolUner
	SymbolUnnamed
	SymbolUsage
	SymbolUserDefinedTypeCatalog
	SymbolUserDefinedTypeCode
	SymbolUserDefinedTypeName
	SymbolUserDefinedTypeSchema
	SymbolView
	SymbolWork
	SymbolWrite
	SymbolZone
	SymbolANSI2003NonReservedEnd = 20000
)

const (
	A                         = "A"
	Abs                       = "ABS"
	Absolute                  = "ABSOLUTE"
	Action                    = "ACTION"
	ADA                       = "ADA"
	Admin                     = "ADMIN"
	After                     = "AFTER"
	Always                    = "ALWAYS"
	Asc                       = "ASC"
	Assertion                 = "ASSERTION"
	Assignment                = "ASSIGNMENT"
	Attribute                 = "ATTRIBUTE"
	Attributes                = "ATTRIBUTES"
	Avg                       = "AVG"
	Before                    = "BEFORE"
	Bernoulli                 = "BERNOULLI"
	Breadth                   = "BREADTH"
	C                         = "C"
	Cardinality               = "CARDINALITY"
	Cascade                   = "CASCADE"
	Catalog                   = "CATALOG"
	CatalogName               = "CATALOG_NAME"
	Ceil                      = "CEIL"
	Ceiling                   = "CEILING"
	Chain                     = "CHAIN"
	Characteristics           = "CHARACTERISTICS"
	Characters                = "CHARACTERS"
	CharacterLength           = "CHARACTER_LENGTH"
	CharacterSetCatalog       = "CHARACTER_SET_CATALOG"
	CharacterSetName          = "CHARACTER_SET_NAME"
	CharacterSetSchema        = "CHARACTER_SET_SCHEMA"
	CharLength                = "CHAR_LENGTH"
	Checked                   = "CHECKED"
	ClassOrigin               = "CLASS_ORIGIN"
	Coalesce                  = "COALESCE"
	COBOL                     = "COBOL"
	CodeUnits                 = "CODE_UNITS"
	Collation                 = "COLLATION"
	CollationCatalog          = "COLLATION_CATALOG"
	CollationName             = "COLLATION_NAME"
	CollationSchema           = "COLLATION_SCHEMA"
	Collect                   = "COLLECT"
	ColumnName                = "COLUMN_NAME"
	CommandFunction           = "COMMAND_FUNCTION"
	CommandFunctionCode       = "COMMAND_FUNCTION_CODE"
	Committed                 = "COMMITTED"
	Condition                 = "CONDITION"
	ConditionNumber           = "CONDITION_NUMBER"
	ConnectionName            = "CONNECTION_NAME"
	Constraints               = "CONSTRAINTS"
	ConstraintCatalog         = "CONSTRAINT_CATALOG"
	ConstraintName            = "CONSTRAINT_NAME"
	ConstraintSchema          = "CONSTRAINT_SCHEMA"
	Constructors              = "CONSTRUCTORS"
	Contains                  = "CONTAINS"
	Convert                   = "CONVERT"
	Corr                      = "CORR"
	Count                     = "COUNT"
	CovarPop                  = "COVAR_POP"
	CovarSamp                 = "COVAR_SAMP"
	CumeDist                  = "CUME_DIST"
	CurrentCollation          = "CURRENT_COLLATION"
	CursorName                = "CURSOR_NAME"
	Data                      = "DATA"
	DatetimeIntervalCode      = "DATETIME_INTERVAL_CODE"
	DatetimeIntervalPrecision = "DATETIME_INTERVAL_PRECISION"
	Defaults                  = "DEFAULTS"
	Deferrable                = "DEFERRABLE"
	Deferred                  = "DEFERRED"
	Defined                   = "DEFINED"
	Definer                   = "DEFINER"
	Degree                    = "DEGREE"
	DenseRank                 = "DENSE_RANK"
	Depth                     = "DEPTH"
	Derived                   = "DERIVED"
	Desc                      = "DESC"
	Descriptor                = "DESCRIPTOR"
	Diagnostics               = "DIAGNOSTICS"
	Dispatch                  = "DISPATCH"
	Domain                    = "DOMAIN"
	DynamicFunction           = "DYNAMIC_FUNCTION"
	DynamicFunctionCode       = "DYNAMIC_FUNCTION_CODE"
	Equals                    = "EQUALS"
	Every                     = "EVERY"
	Exception                 = "EXCEPTION"
	Exclude                   = "EXCLUDE"
	Excluding                 = "EXCLUDING"
	Exp                       = "EXP"
	Extract                   = "EXTRACT"
	Final                     = "FINAL"
	First                     = "FIRST"
	Floor                     = "FLOOR"
	Following                 = "FOLLOWING"
	FORTRAN                   = "FORTRAN"
	Found                     = "FOUND"
	Fusion                    = "FUSION"
	G                         = "G"
	General                   = "GENERAL"
	Go                        = "GO"
	Goto                      = "GOTO"
	Granted                   = "GRANTED"
	Hierarchy                 = "HIERARCHY"
	Implementation            = "IMPLEMENTATION"
	Including                 = "INCLUDING"
	Increment                 = "INCREMENT"
	Initially                 = "INITIALLY"
	Instance                  = "INSTANCE"
	Instantiable              = "INSTANTIABLE"
	Intersection              = "INTERSECTION"
	Invoker                   = "INVOKER"
	K                         = "K"
	Key                       = "KEY"
	KeyMember                 = "KEY_MEMBER"
	KeyType                   = "KEY_TYPE"
	Last                      = "LAST"
	Length                    = "LENGTH"
	Level                     = "LEVEL"
	Ln                        = "LN"
	Locator                   = "LOCATOR"
	Lower                     = "LOWER"
	M                         = "M"
	Map                       = "MAP"
	Matched                   = "MATCHED"
	Max                       = "MAX"
	MaxValue                  = "MAXVALUE"
	MessageLength             = "MESSAGE_LENGTH"
	MessageOctetLength        = "MESSAGE_OCTET_LENGTH"
	MessageText               = "MESSAGE_TEXT"
	Min                       = "MIN"
	MinValue                  = "MINVALUE"
	Mod                       = "MOD"
	More                      = "MORE"
	MUMPS                     = "MUMPS"
	Name                      = "NAME"
	Names                     = "NAMES"
	Nesting                   = "NESTING"
	Next                      = "NEXT"
	Normalize                 = "NORMALIZE"
	Normalized                = "NORMALIZED"
	Nullable                  = "NULLABLE"
	NullIf                    = "NULLIF"
	Nulls                     = "NULLS"
	Number                    = "NUMBER"
	Object                    = "OBJECT"
	Octets                    = "OCTETS"
	OctetLength               = "OCTET_LENGTH"
	Option                    = "OPTION"
	Options                   = "OPTIONS"
	Ordering                  = "ORDERING"
	Ordinality                = "ORDINALITY"
	Others                    = "OTHERS"
	Overlay                   = "OVERLAY"
	Overriding                = "OVERRIDING"
	Pad                       = "PAD"
	ParameterMode             = "PARAMETER_MODE"
	ParameterName             = "PARAMETER_NAME"
	ParameterOrdinalPosition  = "PARAMETER_ORDINAL_POSITION"
	ParameterSpecificCatalog  = "PARAMETER_SPECIFIC_CATALOG"
	ParameterSpecificName     = "PARAMETER_SPECIFIC_NAME"
	ParameterSpecificSchema   = "PARAMETER_SPECIFIC_SCHEMA"
	Partial                   = "PARTIAL"
	Pascal                    = "PASCAL"
	Path                      = "PATH"
	PercentileCont            = "PERCENTILE_CONT"
	PercentileDisc            = "PERCENTILE_DISC"
	PercentRank               = "PERCENT_RANK"
	Placing                   = "PLACING"
	PLI                       = "PLI"
	Position                  = "POSITION"
	Power                     = "POWER"
	Preceding                 = "PRECEDING"
	Preserve                  = "PRESERVE"
	Prior                     = "PRIOR"
	Privileges                = "PRIVILEGES"
	Public                    = "PUBLIC"
	Rank                      = "RANK"
	Read                      = "READ"
	Relative                  = "RELATIVE"
	Repeatable                = "REPEATABLE"
	Restart                   = "RESTART"
	ReturnedCardinality       = "RETURNED_CARDINALITY"
	ReturnedLength            = "RETURNED_LENGTH"
	ReturnedOctetLength       = "RETURNED_OCTET_LENGTH"
	ReturnedSQLState          = "RETURNED_SQLSTATE"
	Role                      = "ROLE"
	Routine                   = "ROUTINE"
	RoutineCatalog            = "ROUTINE_CATALOG"
	RoutineName               = "ROUTINE_NAME"
	RoutineSchema             = "ROUTINE_SCHEMA"
	RowCount                  = "ROW_COUNT"
	RowNumber                 = "ROW_NUMBER"
	Scale                     = "SCALE"
	Schema                    = "SCHEMA"
	SchemaName                = "SCHEMA_NAME"
	ScopeCatalog              = "SCOPE_CATALOG"
	ScopeName                 = "SCOPE_NAME"
	ScopeSchema               = "SCOPE_SCHEMA"
	Section                   = "SECTION"
	Security                  = "SECURITY"
	Self                      = "SELF"
	Sequence                  = "SEQUENCE"
	Serializable              = "SERIALIZABLE"
	ServerName                = "SERVER_NAME"
	Session                   = "SESSION"
	Sets                      = "SETS"
	Simple                    = "SIMPLE"
	Size                      = "SIZE"
	Source                    = "SOURCE"
	SpaceWord                 = "SPACE"
	SpecificName              = "SPECIFIC_NAME"
	Sqrt                      = "SQRT"
	State                     = "STATE"
	Statement                 = "STATEMENT"
	StddevPop                 = "STDDEV_POP"
	StddevSamp                = "STDDEV_SAMP"
	Structure                 = "STRUCTURE"
	Style                     = "STYLE"
	SubclassOrigin            = "SUBCLASS_ORIGIN"
	Substring                 = "SUBSTRING"
	Sum                       = "SUM"
	TableSample               = "TABLESAMPLE"
	TableName                 = "TABLE_NAME"
	Temporary                 = "TEMPORARY"
	Ties                      = "TIES"
	TopLevelCount             = "TOP_LEVEL_COUNT"
	Transaction               = "TRANSACTION"
	TransactionsCommitted     = "TRANSACTIONS_COMMITTED"
	TransactionsRolledBack    = "TRANSACTIONS_ROLLED_BACK"
	TransactionActive         = "TRANSACTION_ACTIVE"
	Tranform                  = "TRANSFORM"
	Transforms                = "TRANSFORMS"
	Translate                 = "TRANSLATE"
	TriggerCatalog            = "TRIGGER_CATALOG"
	TriggerName               = "TRIGGER_NAME"
	TriggerSchema             = "TRIGGER_SCHEMA"
	Trim                      = "TRIM"
	Type                      = "TYPE"
	Unbounded                 = "UNBOUNDED"
	Uncommitted               = "UNCOMMITTED"
	Under                     = "UNDER"
	Unnamed                   = "UNNAMED"
	Usage                     = "USAGE"
	UserDefinedTypeCatalog    = "USER_DEFINED_TYPE_CATALOG"
	UserDefinedTypeCode       = "USER_DEFINED_TYPE_CODE"
	UserDefinedTypeName       = "USER_DEFINED_TYPE_NAME"
	UserDefinedTypeSchema     = "USER_DEFINED_TYPE_SCHEMA"
	View                      = "VIEW"
	Work                      = "WORK"
	Write                     = "WRITE"
	Zone                      = "ZONE"
)
