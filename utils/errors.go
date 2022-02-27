package utils

func NewHTTPError(errorCode string) map[string]interface{} {
	m := make(map[string]interface{})
	m["error"] = errorCode
	m["error_description"] = errorMessage[errorCode]

	return m
}

func NewHTTCustomError(errorCode, errorMsg string) map[string]interface{} {
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
	ThingNotFound       = "thingNotFound"
)

// Error code with description
var errorMessage = map[string]string{
	"internalError":       "an internal error occured",
	"invalidBindingModel": "model could not be bound",
	"entityCreationError": "could not create entity",
	"sensorAlreadyExists": "sensor already exists",
	// Sensors
	"sensorsNotFound": "sensors could not be found",
	"sensorNotFound":  "sensor could not be found",
	// Thing
	"thingNotFound": "thing could not be found",
}
