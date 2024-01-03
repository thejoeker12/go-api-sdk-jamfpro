package jamfpro

// Type refers to string representation of target object type. I.e buildings, policies, computergroups

const (
	// Pagination - type: string, error: any
	errMsgFailedPaginatedGet = "failed to get paginated %s, error: %v"

	// CRUD - format always type: string, id/name: any, error: any
	errMsgFailedGet            = "failed to get %s, error: %v"
	errMsgFailedGetByID        = "failed to get %s by id: %v, error: %v"
	errMsgFailedGetByName      = "failed to get %s by name: %s, error: %v"
	errMsgFailedGetByCategory  = "failed to get %s by category: %s, error: %v"
	errMsgFailedGetByType      = "failed to get %s by type: %s, error: %v"
	errMsgFailedGetByEmail     = "failed to get %s by Email: %s, error: %v"
	errMsgFailedCreate         = "failed to create %s, error: %v"
	errMsgFailedUpdate         = "failed to update %s, error: %v"
	errMsgFailedUpdateByID     = "failed to update %s by id: %v, error: %v"
	errMsgFailedUpdateByName   = "failed to update %s by name: %s, error: %v"
	errMsgFailedUpdateByEmail  = "failed to Update %s by Email: %s, error: %v"
	errMsgFailedDeleteByID     = "failed to delete %s by id: %v, error: %v"
	errMsgFailedDeleteByName   = "failed to delete %s by name: %s, error: %v"
	errMsgFailedDeleteByEmail  = "failed to delete %s by Email: %s, error: %v"
	errMsgFailedDeleteMultiple = "failed to delete multiple %s, by ids: %v, error: %v"

	// Mapstructure - type: string, error: any
	errMsgFailedMapstruct = "failed to map interfaced %s to structs, error: %v"

	// JSON Marshalling
	errMsgFailedJsonMarshal = "failed to marshal %s, error: %v"

	// Client Credentials
	errMsgFailedRefreshClientCreds = "failed to refresh client credentials at id: %s, error :%v"
)
