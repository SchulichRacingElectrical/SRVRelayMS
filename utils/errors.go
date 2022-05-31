package utils

import (
	"errors"

	"github.com/jackc/pgconn"
)

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

func GetPostgresError(err error) *pgconn.PgError {
	var perr *pgconn.PgError
	errors.As(err, &perr)
	return perr
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

	// Organization error
	OrganizationDuplicate = "organizationDuplicate"
	OrganizationNotFound  = "organizationNotFound"
	OrganizationsNotFound = "organizationsNotFound"

	// Raw Data Preset Error
	RawDataPresetNotUnique = "rawDataPresetNotUnique"
	RawDataPresetNotValid  = "rawDataPresetNotValid"
	RawDataPresetNotFound  = "rawDataPresetNotFound"

	// Chart Preset Error
	ChartPresetNotUnique = "chartPresetNotUnique"
	ChartPresetNotValid  = "chartPresetNotValid"
	ChartPresetNotFound  = "chartPresetNotFound"

	// Datum Error
	DatumNotFound = "datumNotFound"

	// Collection Error
	CollectionsNotFound = "collectionsNotFound"
	CollectionNotFound  = "collectionNotFound"

	// Session Error
	SessionsNotFound = "sessionsNotFound"
	SessionNotFound  = "sessionNotFound"

	// Comments Error
	CommentsNotFound = "commentsNotFound"
	CommentNotFound  = "commentNotFound"
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

	// Organization
	"organizationDuplicate": "Organization name is taken.",
	"organizationNotFound":  "Organization could not be found.",
	"organizationsNotFound": "Organizations could not be found.",

	// File
	"noFileReceived":         "No file is received",
	"notCsv":                 "Not a csv",
	"runHasAssociatedFile":   "Run already has associated file",
	"runHasNoAssociatedFile": "Run does exist or not have associated file",
	"cannotRetrieveFile":     "Cannot retrieve file",

	// Collection errors
	"collectionsNotFound": "Collections could not be found",
	"collectionNotFound":  "Collection could not be found",

	// Session errors
	"sessionsNotFound": "Sessions could not be found",
	"sessionNotFound":  "Session could not be found",

	// Comment errors
	"commentsNotFound": "Comments could not be found",
	"commentNotFound":  "Comment could not be found",

	// Raw Data Preset
	"rawDataPresetNotUnique": "Raw Data Preset name must be unique.",
	"rawDataPresetNotValid":  "Raw Data Preset is not valid.",
	"rawDataPresetNotFound":  "Raw Data Preset not found.",

	// Chart Preset
	"chartPresetNotUnique": "Chart Preset name must be unique.",
	"chartPresetNotValid":  "Chart Preset was not valid. Ensure posted Sensors exist.",
	"chartPresetNotFound":  "Chart Preset was not found.",

	// Datum
	"datumNotFound": "Datum could not be found.",
}
