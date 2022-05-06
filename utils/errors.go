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

//Error codes
const (

	// Other errors
	InternalError       = "internalError"
	InvalidBindingModel = "invalidBindingModel"
	SensorAlreadyExists = "sensorAlreadyExists"
	EntityCreationError = "entityCreationError"
	BadRequest          = "badRequest"

	// Sensor errors
	SensorsNotFound = "sensorsNotFound"
	SensorNotFound  = "sensorNotFound"

	// User error
	UserNotFound      = "userNotFound"
	WrongPassword     = "wrongPassword"
	UserAlreadyExists = "userAlreadyExists"

	// Thing error
	ThingNotFound = "thingNotFound"

	// Run error
	RunsNotFound = "runsNotFound"
	RunNotFound  = "runNotFound"

	// Comment error
	CommentsNotFound                    = "commentsNotFound"
	CommentDoesNotExist                 = "commentDoesNotExist"
	CommentCannotUpdateOtherUserComment = "commentCannotUpdateOtherUserComment"

	// Organization error
	OrganizationNotFound  = "organizationNotFound"
	OrganizationsNotFound = "organizationsNotFound"
)

// Error code with description
var errorMessage = map[string]string{

	// Generic errors
	"internalError":       "an internal error occurred",
	"invalidBindingModel": "model could not be bound",
	"entityCreationError": "could not create entity",

	// Sensor errors
	"sensorAlreadyExists": "sensor already exists",
	"sensorsNotFound":     "sensors could not be found",
	"sensorNotFound":      "sensor could not be found",

	// User errors
	"userNotFound":      "user could not be found",
	"wrongPassword":     "password was incorrect",
	"userAlreadyExists": "email is already being used",

	// Thing
	"thingNotFound": "thing could not be found",

	// Comment
	"commentsNotFound":                    "comments could not be found",
	"commentDoesNotExist":                 "comment does not exist",
	"commentCannotUpdateOtherUserComment": "cannot update comment of another user",

	// Organization
	"organizationNotFound":  "organization could not be found",
	"organizationsNotFound": "organizations could not be found",

	// Run errors
	"runsNotFound": "Runs could not be found",
	"runNotFound":  "Run could not be found",
}
