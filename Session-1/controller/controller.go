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

// Controller defines the HTTP API for managing projects, sellers, buyers, and bids.
// Each method corresponds to an HTTP endpoint.
type Controller interface {
	CreateProject(c echo.Context) error // POST /create-project
	UpdateBID(c echo.Context) error     // PUT /update-bid
	ComputeBID(c echo.Context) error    // POST /compute-bid
	AttachHandlers(lister *echo.Echo)   // Attach all routes to Echo
	CreateSeller(c echo.Context) error  // POST /create-seller
	CreateBuyer(c echo.Context) error   // POST /create-buyer
}

// ControllerImpl is the concrete implementation of Controller.
// It uses BidManager for bid operations and ProjectManager for project/seller/buyer persistence.
type ControllerImpl struct {
	bidManager     bidManager.BidManager
	projectManager project.ProjectManager
}

// NewController initializes a new Controller with the required dependencies.
func NewController(bidManager bidManager.BidManager, projectDetails project.ProjectManager) Controller {
	return &ControllerImpl{
		bidManager,
		projectDetails,
	}
}

// AttachHandlers registers all HTTP endpoints with Echo.
func (co *ControllerImpl) AttachHandlers(lister *echo.Echo) {
	lister.POST("/create-project", co.CreateProject)
	lister.POST("/create-seller", co.CreateSeller)
	lister.POST("/create-buyer", co.CreateBuyer)
	lister.PUT("/update-bid", co.UpdateBID)
	lister.GET("/get-projects", co.GetProjects)
	lister.POST("/compute-bid", co.ComputeBID)
}

// UpdateBID handles PUT /update-bid.
// Reads a bid from request body and updates it for a given project.
func (co *ControllerImpl) UpdateBID(c echo.Context) error {
	glog.Info("update-bid")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")

	projectID := c.QueryParam("projectID")

	// Parse request body into a Bid object
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

	// Delegate bid update to BidManager
	err = co.bidManager.DoBID(projectID, bid)
	if err != nil {
		glog.Error("update-bid-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

// CreateSeller handles POST /create-seller.
// Reads a seller from request body and inserts it into the database.
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

	// Insert seller using ProjectManager
	err = co.projectManager.CreateSeller(seller)
	if err != nil {
		glog.Error("create-seller-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

// CreateBuyer handles POST /create-buyer.
// Reads a buyer from request body and inserts it into the database.
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

	// Insert buyer using ProjectManager
	err = co.projectManager.CreateBuyer(buyer)
	if err != nil {
		glog.Error("create-buyer-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

// CreateProject handles POST /create-project.
// Reads a project from request body and saves it in the database.
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

	// Insert project using ProjectManager
	err = co.projectManager.CreateProject(projectDetails)
	if err != nil {
		glog.Error("create-project-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

// GetProjects handles GET /get-projects.
// Fetches and returns all projects from the database.
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

// ComputeBID handles POST /compute-bid.
// Determines the winning buyer (lowest bid) for a given project.
func (co *ControllerImpl) ComputeBID(c echo.Context) error {
	glog.Info("compute-bid")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")

	projectID := c.QueryParam("projectID")

	// Delegate to BidManager to compute winning buyer
	bidWinner, err := co.bidManager.ComputeBID(projectID)
	if err != nil {
		glog.Error("get-user-error", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, bidWinner)
}
