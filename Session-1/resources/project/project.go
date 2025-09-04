package project

import (
	"context"
	"math"
	"time"

	"github.com/21keshav/IBackendApplication/config"
	"github.com/21keshav/IBackendApplication/util"
	"github.com/golang/glog"
)

//
// Domain Models
//

// ProjectDetails represents a project posted by a seller.
// Each project can have multiple bids from buyers.
type ProjectDetails struct {
	ID        string         `json:"id,omitempty" bson:"id,omitempty"`
	Details   []string       `json:"details,omitempty" bson:"details,omitempty"`   // Additional project details
	SellerID  string         `json:"seller_id,omitempty" bson:"seller_id,omitempty"`
	BIDS      map[string]BID `json:"bids,omitempty" bson:"bids,omitempty"`         // Keyed by Bid ID
	startDate time.Duration  `json:"start_date,omitempty" bson:"start_date,omitempty"`
	endDate   time.Duration  `json:"end_date,omitempty" bson:"end_date,omitempty"`
}

// BID represents a buyer's offer for a project.
type BID struct {
	ID       string `json:"id,omitempty" bson:"id,omitempty"`
	SellerID string `json:"seller_id,omitempty" bson:"seller_id,omitempty"`
	BuyerID  string `json:"buyer_id,omitempty" bson:"buyer_id,omitempty"`
	Amount   int    `json:"ammount,omitempty" bson:"ammount,omitempty"`
}

// Seller represents a seller who can create projects.
type Seller struct {
	ID         string `json:"id,omitempty" bson:"id,omitempty"`
	SellerID   string `json:"seller_id,omitempty" bson:"seller_id,omitempty"`
	SellerName string `json:"seller_name,omitempty" bson:"seller_name,omitempty"`
}

// Buyer represents a buyer who can place bids on projects.
type Buyer struct {
	ID        string `json:"id,omitempty" bson:"id,omitempty"`
	BuyerID   string `json:"buyer_id,omitempty" bson:"buyer_id,omitempty"`
	BuyerName string `json:"buyer_name,omitempty" bson:"buyer_name,omitempty"`
}

//
// ProjectManager Interface
//
// Defines operations for managing projects, buyers, sellers, and bids.
// This makes it easy to swap implementations (e.g., MongoDB vs mock for testing).
//
type ProjectManager interface {
	CreateProject(projectDetails ProjectDetails) error
	GetProjects() ([]ProjectDetails, error)
	GetProject(projectID string) (ProjectDetails, error)
	UpdateProject(projectID string, bid BID) error

	CreateBuyer(buyer Buyer) error
	GetBuyer(buyerID string) (Buyer, error)

	CreateSeller(seller Seller) error
}

//
// ProjectManagerImpl
//
// Concrete implementation of ProjectManager backed by MongoDB.
// Uses a generic MongoClient interface for all persistence operations.
//
type ProjectManagerImpl struct {
	MongoClient util.MongoClient   // Mongo client wrapper
	ctx         context.Context    // Context for DB operations
	DBConfig    config.DatabaseDetails // Config (db/collection names)
}

// NewProjectManager creates a new ProjectManager backed by MongoDB.
func NewProjectManager(mongoClient util.MongoClient, ctx context.Context, dbConfig config.DatabaseDetails) ProjectManager {
	return &ProjectManagerImpl{
		MongoClient: mongoClient,
		ctx:         ctx,
		DBConfig:    dbConfig,
	}
}

//
// Seller Operations
//

// CreateSeller inserts a new seller into the Sellers collection.
func (um *ProjectManagerImpl) CreateSeller(seller Seller) error {
	glog.Info("pm-create-seller")
	defer glog.Info("pm-create-seller-completed")

	_, err := um.MongoClient.InsertData(um.DBConfig.SellersDBName,
		um.DBConfig.CollectionName, seller)
	if err != nil {
		glog.Error("mongo error inserting seller", err)
		return err
	}
	return nil
}

//
// Buyer Operations
//

// CreateBuyer inserts a new buyer into the Buyers collection.
func (um *ProjectManagerImpl) CreateBuyer(buyer Buyer) error {
	glog.Info("pm-create-buyer")
	defer glog.Info("pm-create-buyer-completed")

	_, err := um.MongoClient.InsertData(um.DBConfig.BuyersDBName,
		um.DBConfig.CollectionName, buyer)
	if err != nil {
		glog.Error("mongo error inserting buyer", err)
		return err
	}
	return nil
}

// GetBuyer fetches a buyer by ID.
func (um *ProjectManagerImpl) GetBuyer(buyerID string) (Buyer, error) {
	glog.Info("pm-get-buyer")
	defer glog.Info("pm-get-buyer-completed")

	var buyer Buyer
	err := um.MongoClient.FindObject(um.DBConfig.BuyersDBName,
		um.DBConfig.CollectionName, Buyer{ID: buyerID}, &buyer)
	if err != nil {
		glog.Error("mongo error finding buyer", err)
		return buyer, err
	}
	return buyer, nil
}

//
// Project Operations
//

// CreateProject inserts a new project into the Projects collection.
func (um *ProjectManagerImpl) CreateProject(projectDetails ProjectDetails) error {
	glog.Info("pm-create-project")
	defer glog.Info("pm-create-project-completed")

	_, err := um.MongoClient.InsertData(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, projectDetails)
	if err != nil {
		glog.Error("mongo error inserting project", err)
		return err
	}
	return nil
}

// GetProjects fetches all projects.
func (um *ProjectManagerImpl) GetProjects() ([]ProjectDetails, error) {
	glog.Info("pm-get-projects")
	defer glog.Info("pm-get-projects-completed")

	var projects []ProjectDetails
	err := um.MongoClient.FindAllObjects(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, &projects, math.MaxInt32)
	if err != nil {
		glog.Error("mongo error finding projects", err)
		return projects, err
	}
	return projects, nil
}

// GetProject fetches a single project by ID.
func (um *ProjectManagerImpl) GetProject(projectID string) (ProjectDetails, error) {
	glog.Info("pm-get-project")
	defer glog.Info("pm-get-project-completed")

	var projectDetails ProjectDetails
	err := um.MongoClient.FindObject(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, ProjectDetails{ID: projectID}, &projectDetails)
	if err != nil {
		glog.Error("mongo error finding project", err)
		return projectDetails, err
	}
	return projectDetails, nil
}

// UpdateProject adds or updates a bid inside a project.
// Loads the project, mutates its BIDS map, then updates it in DB.
func (um *ProjectManagerImpl) UpdateProject(projectID string, bid BID) error {
	glog.Info("pm-update-project")
	defer glog.Info("pm-update-project-completed")

	// Fetch existing project
	var projectDetails ProjectDetails
	err := um.MongoClient.FindObject(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, ProjectDetails{ID: projectID}, &projectDetails)
	if err != nil {
		glog.Error("mongo error finding project", err)
		return err
	}

	// Add/update the bid
	if projectDetails.BIDS == nil {
		projectDetails.BIDS = make(map[string]BID)
	}
	projectDetails.BIDS[bid.ID] = bid

	// Persist changes
	_, err = um.MongoClient.UpdateOne(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, ProjectDetails{ID: projectID}, projectDetails)
	if err != nil {
		glog.Error("mongo error updating project", err)
		return err
	}
	return nil
}
