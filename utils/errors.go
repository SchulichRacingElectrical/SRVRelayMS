package utils

func NewHTTPError(errorCode string) map[string]interface{} {
	m := make(map[string]interface{})
	m["error"] = errorCode
	m["error_description"] = errorMessage[errorCode]
	return m
}

func NewHTTPCustomError(errorCode, errorMsg string) map[string]interface{} {
	m := make(map[string]interface{})
	m["error"] = errorCode
	m["error_description"] = errorMsg
	return m
}

// Error codes - This can be done more nicely
const (
	// Other errors
	InternalError       = "internalError"
	InvalidBindingModel = "invalidBindingModel"
	EntityCreationError = "entityCreationError"
	BadRequest          = "badRequest"
	Unauthorized        = "unauthorized"
	Forbidden           = "forbidden"

	// Sensor errors
	SensorsNotFound     = "sensorsNotFound"
	SensorNotFound      = "sensorNotFound"
	SensorAlreadyExists = "sensorAlreadyExists"
	SensorNotUnique     = "sensorNotUnique"

	// User error
	UserNotFound    = "userNotFound"
	UsersNotFound   = "usersNotFound"
	WrongPassword   = "wrongPassword"
	UserNotApproved = "userPendingApproval"
	UserConflict    = "userConflict"
	UserLastAdmin   = "userLastAdmin"

	// Thing error
	ThingsNotFound = "thingsNotFound"
	ThingNotFound  = "thingNotFound"
	ThingNotUnique = "thingNotUnique"

	// Operator error
	OperatorsNotFound = "operatorsNotFound"
	OperatorNotFound  = "operatorNotFound"
	OperatorNotUnique = "operatorNotUnique"

	// ThingOperator error
	ThingOperatorNotUnique = "thingOperatorNotUnique"

	// Run error
	RunsNotFound = "runsNotFound"
	RunNotFound  = "runNotFound"
	RunDNE       = "runDNE"

	// File error
	NoFileReceived         = "noFileReceived"
	NotCsv                 = "notCsv"
	FileNotUploaded        = "fileNotUploaded"
	RunHasAssociatedFile   = "runHasAssociatedFile"
	RunHasNoAssociatedFile = "runHasNoAssociatedFile"
	CannotRetrieveFile     = "cannotRetrieveFile"

	// Session error
	SessionsNotFound = "sessionsNotFound"
	SessionNotFound  = "sessionNotFound"

	// Comment error
	CommentsNotFound                    = "commentsNotFound"
	CommentDoesNotExist                 = "commentDoesNotExist"
	CommentCannotUpdateOtherUserComment = "commentCannotUpdateOtherUserComment"
	UserIdMissing                       = "userIdMissing"

	// Organization error
	OrganizationDuplicate = "organizationDuplicate"
	OrganizationNotFound  = "organizationNotFound"
	OrganizationsNotFound = "organizationsNotFound"
)

// Error code with description
var errorMessage = map[string]string{
	// Generic errors
	"internalError":       "An internal error occurred.",
	"invalidBindingModel": "The model could not be bound.",
	"entityCreationError": "Could not create entity.",
	"unauthorized":        "Unauthorized.",
	"forbidden":           "Forbidden.",

	// Sensor errors
	"sensorAlreadyExists": "Sensor already exists.",
	"sensorsNotFound":     "Sensors could not be found.",
	"sensorNotFound":      "Sensor could not be found.",
	"sensorNotUnique":     "Sensor name and CAN ID must be unique for a thing.",

	// User errors
	"userNotFound":  "User could not be found.",
	"usersNotFound": "Users not found.",
	"wrongPassword": "Password was incorrect.",
	"userConflict":  "Email must be globally unique and name must be organizationally unique.",
	"userLastAdmin": "The last administrator in the organization cannot be deleted or have their role changed.",

	// Thing
	"thingsNotFound": "Things could not be found.",
	"thingNotFound":  "Thing could not be found.",
	"thingNotUnique": "Thing name must be unique",

	// Operator
	"operatorsNotFound": "Operators could not be found.",
	"operatorNotFound":  "Operator could not be found.",
	"operatorNotUnique": "Operator name must be unique.",

	// ThingOperator
	"thingOperatorNotUnique": "Thing Operator association already exists.",

	// Comment
	"commentsNotFound":                    "comments could not be found",
	"commentDoesNotExist":                 "comment does not exist",
	"commentCannotUpdateOtherUserComment": "cannot update comment of another user",
	"userIdMissing":                       "userId missing",

	// Organization
	"organizationNotFound":  "organization could not be found",
	"organizationsNotFound": "organizations could not be found",
	"organizationDuplicate": "Organization name is taken.",

	// Run errors
	"runsNotFound": "Runs could not be found",
	"runNotFound":  "Run could not be found",
	"runDNE":       "Run does not exist",

	// File
	"noFileReceived":         "No file is received",
	"notCsv":                 "Not a csv",
	"runHasAssociatedFile":   "Run already has associated file",
	"runHasNoAssociatedFile": "Run does exist or not have associated file",
	"cannotRetrieveFile":     "Cannot retrieve file",

	// Session errors
	"sessionssNotFound": "Sesssions could not be found",
	"sessionNotFound":   "Session could not be found",
}
