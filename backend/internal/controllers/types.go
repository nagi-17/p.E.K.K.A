package controllers

type UpgradeBuildingReq struct {
	BuildingID string `json:"building_id"`
}

type Response struct {
	Message string `json:"message"`
}
