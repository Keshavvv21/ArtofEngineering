package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/21keshav/IBackendApplication/resources/bidManager"
	"github.com/21keshav/IBackendApplication/resources/project"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

type Controller interface {
	CreateProject(c echo.Context) error
	UpdateBID(c echo.Context) error
	ComputeBID(c echo.Context) error
	AttachHandlers(lister *echo.Echo)
	CreateSeller(c echo.Context) error
	CreateBuyer(c echo.Context) error
}

type ControllerImpl struct {
	bidManager     bidManager.BidManager
	projectManager project.ProjectManager
}

func NewController(bidManager bidManager.BidManager, projectDetails project.ProjectManager) Controller {
	return &ControllerImpl{
		bidManager,
		projectDetails,
	}
}

func (co *ControllerImpl) AttachHandlers(lister *echo.Echo) {
	lister.POST("/create-project", co.CreateProject)
	lister.POST("/create-seller", co.CreateSeller)
	lister.POST("/create-buyer", co.CreateBuyer)
	lister.PUT("/update-bid", co.UpdateBID)
	lister.GET("/get-projects", co.GetProjects)
	lister.POST("/compute-bid", co.ComputeBID)
}

func (co *ControllerImpl) UpdateBID(c echo.Context) error {
	glog.Info("update-bid")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")
	projectID := c.QueryParam("projectID")
	var bid project.BID
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		glog.Error("read-error", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err = json.Unmarshal(body, &bid)
	if err != nil {
		glog.Error("unmarshal-error", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	err = co.bidManager.DoBID(projectID, bid)
	if err != nil {
		glog.Error("update-bid-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

func (co *ControllerImpl) CreateSeller(c echo.Context) error {

	glog.Info("create-seller")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")

	var seller project.Seller
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		glog.Error("read-error", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err = json.Unmarshal(body, &seller)
	if err != nil {
		glog.Error("unmarshal-error", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	err = co.projectManager.CreateSeller(seller)
	if err != nil {
		glog.Error("create-seller-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

func (co *ControllerImpl) CreateBuyer(c echo.Context) error {

	glog.Info("create-buyer")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")

	var buyer project.Buyer
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		glog.Error("read-error", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err = json.Unmarshal(body, &buyer)
	if err != nil {
		glog.Error("unmarshal-error", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	err = co.projectManager.CreateBuyer(buyer)
	if err != nil {
		glog.Error("create-buyer-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

func (co *ControllerImpl) CreateProject(c echo.Context) error {

	glog.Info("create-project")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")

	var projectDetails project.ProjectDetails
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		glog.Error("read-error", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err = json.Unmarshal(body, &projectDetails)
	if err != nil {
		glog.Error("unmarshal-error", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	err = co.projectManager.CreateProject(projectDetails)
	if err != nil {
		glog.Error("create-project-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

func (co *ControllerImpl) GetProjects(c echo.Context) error {
	glog.Info("get-project")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")
	projectDetails, err := co.projectManager.GetProjects()
	if err != nil {
		glog.Error("get-user-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, projectDetails)
}

func (co *ControllerImpl) ComputeBID(c echo.Context) error {
	glog.Info("compute-bid")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")
	projectID := c.QueryParam("projectID")
	bidWinner, err := co.bidManager.ComputeBID(projectID)
	if err != nil {
		glog.Error("get-user-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, bidWinner)
}
