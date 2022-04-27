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
	InternalError       	= "internalError"
	InvalidBindingModel 	= "invalidBindingModel"
	EntityCreationError 	= "entityCreationError"
	BadRequest          	= "badRequest"
	Unauthorized					= "unauthorized"

	// Sensor errors
	SensorsNotFound 			= "sensorsNotFound"
	SensorNotFound  			= "sensorNotFound"
	SensorAlreadyExists 	= "sensorAlreadyExists"
	SensorNotUnique				= "sensorNotUnique"

	// User error
	UserNotFound      		= "userNotFound"
	UsersNotFound					= "usersNotFound"
	WrongPassword     		= "wrongPassword"
	UserNotApproved				= "userPendingApproval"
	UserAlreadyExists 		= "userAlreadyExists"

	// Thing error
	ThingsNotFound 				= "thingsNotFound"
	ThingNotFound 				= "thingNotFound"
	ThingNotUnique				=	"thingNotUnique"

	// Operator error
	OperatorsNotFound			= "operatorsNotFound"
	OperatorNotFound			= "operatorNotFound"
	OperatorNotUnique			=	"operatorNotUnique"

	// Organization error
	OrganizationDuplicate = "organizationDuplicate"
	OrganizationNotFound  = "organizationNotFound"
	OrganizationsNotFound = "organizationsNotFound"
)

// Error code with description
var errorMessage = map[string]string {
	// Generic errors
	"internalError":       		"An internal error occurred.",
	"invalidBindingModel": 		"The model could not be bound.",
	"entityCreationError": 		"Could not create entity.",
	"unauthorized": 					"Unauthorized.",

	// Sensor errors
	"sensorAlreadyExists": 		"Sensor already exists.",
	"sensorsNotFound":     		"Sensors could not be found.",
	"sensorNotFound":      		"Sensor could not be found.",
	"sensorNotUnique":				"Sensor name and CAN ID must be unique for a thing.",

	// User errors
	"userNotFound":      			"User could not be found.",
	"usersNotFound":		 			"Users not found.",	
	"wrongPassword":     			"Password was incorrect.",
	"userAlreadyExists": 			"Email is already being used.",

	// Thing
	"thingsNotFound":					"Things could not be found.",
	"thingNotFound": 					"Thing could not be found.",
	"thingNotUnique":					"Thing name must be unique",

	// Operator
	"operatorsNotFound":			"Operators could not be found.",
	"operatorNotFound":				"Operator could not be found.",
	"operatorNotUnique":			"Operator name must be unique.",

	// Organization
	"organizationDuplicate": 	"Organization name is taken.",
	"organizationNotFound":  	"Organization could not be found.",
	"organizationsNotFound": 	"Organizations could not be found.",
}
