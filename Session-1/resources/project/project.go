package project

import (
	"context"
	"math"
	"time"

	"github.com/21keshav/IBackendApplication/config"
	"github.com/21keshav/IBackendApplication/util"
	"github.com/golang/glog"
)

type ProjectDetails struct {
	ID        string         `json:"id,omitempty" bson:"id,omitempty"`
	Details   []string       `json:"details,omitempty" bson:"details,omitempty"`
	SellerID  string         `json:"seller_id,omitempty" bson:"seller_id,omitempty"`
	BIDS      map[string]BID `json:"bids,omitempty" bson:"bids,omitempty"`
	startDate time.Duration  `json:"start_date,omitempty" bson:"start_date,omitempty"`
	endDate   time.Duration  `json:"end_date,omitempty" bson:"end_date,omitempty"`
}

type BID struct {
	ID       string `json:"id,omitempty" bson:"id,omitempty"`
	SellerID string `json:"seller_id,omitempty" bson:"seller_id,omitempty"`
	BuyerID  string `json:"buyer_id,omitempty" bson:"buyer_id,omitempty"`
	Amount   int    `json:"ammount,omitempty" bson:"ammount,omitempty"`
}

type Seller struct {
	ID         string `json:"id,omitempty" bson:"id,omitempty"`
	SellerID   string `json:"seller_id,omitempty" bson:"seller_id,omitempty"`
	SellerName string `json:"seller_name,omitempty" bson:"seller_name,omitempty"`
}

type Buyer struct {
	ID        string `json:"id,omitempty" bson:"id,omitempty"`
	BuyerID   string `json:"buyer_id,omitempty" bson:"buyer_id,omitempty"`
	BuyerName string `json:"buyer_name,omitempty" bson:"buyer_name,omitempty"`
}

type ProjectManager interface {
	CreateProject(projectDetails ProjectDetails) error
	GetProjects() ([]ProjectDetails, error)
	UpdateProject(projectID string, bid BID) error
	GetProject(projectID string) (ProjectDetails, error)
	GetBuyer(buyerID string) (Buyer, error)
	CreateBuyer(buyer Buyer) error
	CreateSeller(seller Seller) error
}

type ProjectManagerImpl struct {
	MongoClient util.MongoClient
	ctx         context.Context
	DBConfig    config.DatabaseDetails
}

func NewProjectManager(mongoClient util.MongoClient, ctx context.Context, dbConfig config.DatabaseDetails) ProjectManager {

	return &ProjectManagerImpl{
		mongoClient,
		ctx,
		dbConfig,
	}
}

func (um *ProjectManagerImpl) CreateSeller(seller Seller) error {
	glog.Info("pm-create-seller")
	defer glog.Info("pm-create-seller-completed")
	_, err := um.MongoClient.InsertData(um.DBConfig.SellersDBName,
		um.DBConfig.CollectionName, seller)
	if err != nil {
		glog.Error("mongo error inserting object", err)
		return err
	}
	return nil
}

func (um *ProjectManagerImpl) CreateBuyer(buyer Buyer) error {
	glog.Info("pm-create-buyer")
	defer glog.Info("pm-create-buyer-completed")
	_, err := um.MongoClient.InsertData(um.DBConfig.BuyersDBName,
		um.DBConfig.CollectionName, buyer)
	if err != nil {
		glog.Error("mongo error inserting object", err)
		return err
	}
	return nil
}

func (um *ProjectManagerImpl) CreateProject(projectDetails ProjectDetails) error {
	glog.Info("pm-create-project")
	defer glog.Info("pm-create-project-completed")
	_, err := um.MongoClient.InsertData(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, projectDetails)
	if err != nil {
		glog.Error("mongo error inserting object", err)
		return err
	}
	return nil
}

func (um *ProjectManagerImpl) GetProjects() ([]ProjectDetails, error) {
	glog.Info("pm-get-projects")
	defer glog.Info("pm-get-projects-completed")
	var projectDetails []ProjectDetails
	err := um.MongoClient.FindAllObjects(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, &projectDetails, math.MaxInt32)
	if err != nil {
		glog.Error("mongo error finding object", err)
		return projectDetails, err
	}
	return projectDetails, nil
}

func (um *ProjectManagerImpl) GetBuyer(buyerID string) (Buyer, error) {
	glog.Info("pm-get-buyer")
	defer glog.Info("pm-get-buyer-completed")
	var buyer Buyer
	err := um.MongoClient.FindObject(um.DBConfig.BuyersDBName,
		um.DBConfig.CollectionName, Buyer{ID: buyerID}, &buyer)
	if err != nil {
		glog.Error("mongo error finding object", err)
		return buyer, err
	}
	return buyer, nil
}

func (um *ProjectManagerImpl) GetProject(projectID string) (ProjectDetails, error) {
	glog.Info("pm-get-projects")
	defer glog.Info("pm-get-projects-completed")
	var projectDetails ProjectDetails
	err := um.MongoClient.FindObject(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, ProjectDetails{ID: projectID}, &projectDetails)
	if err != nil {
		glog.Error("mongo error finding object", err)
		return projectDetails, err
	}
	return projectDetails, nil
}

func (um *ProjectManagerImpl) UpdateProject(projectID string, bid BID) error {
	glog.Info("pm-update-project")
	defer glog.Info("pm-update-project-completed")
	var projectDetails ProjectDetails
	err := um.MongoClient.FindObject(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, ProjectDetails{ID: projectID}, &projectDetails)
	if err != nil {
		glog.Error("mongo error finding object", err)
		return err
	}

	projectDetails.BIDS[bid.ID] = bid
	_, err = um.MongoClient.UpdateOne(um.DBConfig.ProjectDBName,
		um.DBConfig.CollectionName, ProjectDetails{ID: projectID}, projectDetails)
	if err != nil {
		glog.Error("mongo error update object", err)
		return err
	}
	return nil
}
