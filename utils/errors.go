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
	InternalError       = "internalError"
	InvalidBindingModel = "invalidBindingModel"
	SensorAlreadyExists = "sensorAlreadyExists"
	EntityCreationError = "entityCreationError"
	BadRequest          = "badRequest"
	SensorsNotFound     = "sensorsNotFound"
	SensorNotFound      = "sensorNotFound"
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

}
