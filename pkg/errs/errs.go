package errs

var ErrFmt = "err: [%v]"

// common
var (
	ErrLoadingTimeZone       = _new("CMN000", "error loading timezone data")
	ErrMarshalingJson        = _new("CMN001", "error marshaling json")
	ErrUnmarshalingJson      = _new("CMN002", "error unmarshaling json")
	ErrParsingTime           = _new("CMN003", "error parsing time")
	ErrNoEntityIdProvided    = _new("CMN004", "entity ID is required but none was provided")
	ErrNoDateProvided        = _new("CMN005", "error no date provided")
	ErrNoPayloadData         = _new("CMN006", "error event contains no payload data")
	ErrRepoMockAction        = _new("CMN007", "error repo mock action")
	ErrUnknownErrorType      = _new("CMN008", "error unknown error type")
	ErrInvalidDate           = _new("CMN009", "error invalid date format")
	ErrConvertingStringToInt = _new("CMN010", "error converting string to int")
	ErrGettingParam          = _new("CMN011", "error getting param")
	ErrConvertingPayload     = _new("CMN012", "error converting payload")
	ErrInvalidEmail          = _new("CMN013", "error invalid email")
	ErrOpenFile              = _new("CMN014", "error open file")
	ErrReadingFile           = _new("CMN015", "error reading file")
	ErrActionNotImplemented  = _new("CMN016", "error action is not implemented")
	ErrParsingAtoi           = _new("CMN017", "error converting string to a int (Atoi)")
)

// pkg/api
var (
	ErrResponseWriter = _new("API000", "error writing to response writer")
)

// pkg/config
var (
	ErrCreatingParamStore    = _new("CFG002", "unable to create param store service")
	ErrUnknownConfigProvider = _new("CFG003", "error unknown config provider")
	ErrGettingEnv            = _new("CFG004", "error unknown get env variables")
	ErrCreatingJWT           = _new("CFG005", "error signing new JWT")
)

// pkg/store
var (
	ErrCursor               = _new("STR000", "error using cursor")
	ErrDecodeCursor         = _new("STR001", "unable to decode cursor value to non-pointer variable")
	ErrClosingCursor        = _new("STR002", "error closing cursor")
	ErrMongoConnect         = _new("STR003", "error connecting to mongo")
	ErrMongoInsertOne       = _new("STR004", "error inserting one mongo document")
	ErrMongoFindOne         = _new("STR005", "error finding one mongo document")
	ErrMongoFind            = _new("STR006", "error finding mongo document(s)")
	ErrMongoUpdateOne       = _new("STR007", "error updating one mongo document")
	ErrMongoDeleteOne       = _new("STR008", "error deleting one mongo document")
	ErrNotDocumentInterface = _new("STR009", "cannot insert value that doesn't implement Document interface")
	ErrDecodingInsertedId   = _new("STR010", "error decoding inserted ID")
	ErrParseRequestURI      = _new("STR011", "error parsing request uri")
	ErrMarshalingBson       = _new("STR012", "error marshaling bson")
	ErrUnmarshalingBson     = _new("STR013", "error unmarshaling bson")
)

// pkg/middleware
var (
	ErrTokenIsEmpty          = _new("MID001", "error token is empty")
	ErrTokenSignatureInvalid = _new("MID002", "token signature is invalid")
	ErrTokenInvalid          = _new("MID003", "token is invalid")
	ErrTokenWithInvalidEmail = _new("MID004", "token with invalid email")
)

// pkg/repo
var (
	ErrTeamIsNotFound       = _new("REP001", "team is not found")
	ErrPlayerIsNotFound     = _new("REP002", "player is not found")
	ErrTournamentIsNotFound = _new("REP003", "tournament is not found")
	ErrTransferIsNotFound   = _new("REP004", "transfer is not found")
	ErrMatchIsNotFound      = _new("REP005", "match is not found")
)

// pkg/model
var (
	ErrInvalidCurrencyCode = _new("MDL001", "invalid currency Code")
	ErrInvalidMoney        = _new("MDL002", "invalid money")
	ErrInvalidActionType   = _new("MDL003", "invalid action type")
)

// pkg/kafka
var (
	ErrReadingKafkaMessage        = _new("KAF001", "error reading kafka message")
	ErrWriteKafkaMessage          = _new("KAF002", "error writing kafka message")
	ErrHandlingKafkaMessage       = _new("KAF003", "error handling kafka message")
	ErrHandlingUpdateTeamPlayer   = _new("KAF004", "error handling update team player")
	ErrHandlingGameEventStarted   = _new("KAF005", "error handling game event started")
	ErrHandlingGameEventGoal      = _new("KAF006", "error handling game event goal")
	ErrHandlingGameEventHalftime  = _new("KAF007", "error handling game event halftime")
	ErrHandlingGameEventFinish    = _new("KAF008", "error handling game event finish")
	ErrHandlingGameEventExtratime = _new("KAF009", "error handling game event extratime")
)

// general jobs
var (
	ErrBadMessage = _new("JGN000", "received a message that was not an kafka notification")
)

// validations
var (
	ErrValidation = _new("VAL000", "error on validation")
)
