package bidManager

import (
	"context"
	"math"

	"github.com/21keshav/IBackendApplication/resources/project"
	"github.com/golang/glog"
)

type BidManager interface {
	ComputeBID(projectID string) (project.Buyer, error)
	DoBID(projectID string, bid project.BID) error
}

type BidManagerManagerImpl struct {
	projectManager project.ProjectManager
	ctx            context.Context
}

func NewBidManager(projectManager project.ProjectManager, ctx context.Context) BidManager {

	return &BidManagerManagerImpl{
		projectManager,
		ctx,
	}
}

func (bd *BidManagerManagerImpl) DoBID(projectID string, bid project.BID) error {
	glog.Info("Do-bid-projects")
	defer glog.Info("do-bid-completed")
	err := bd.projectManager.UpdateProject(projectID, bid)
	return err
}

func (bd *BidManagerManagerImpl) ComputeBID(projectID string) (project.Buyer, error) {
	glog.Info("compute-projects")
	defer glog.Info("compute-completed")
	currentProject, err := bd.projectManager.GetProject(projectID)
	if err != nil {
		return project.Buyer{}, err
	}
	min := math.MaxInt64
	var minBid project.BID
	for _, bid := range currentProject.BIDS {
		if bid.Amount < min {
			min = bid.Amount
			minBid = bid
		}
	}

	buyer, err := bd.projectManager.GetBuyer(minBid.BuyerID)
	if err != nil {
		return project.Buyer{}, err
	}

	return buyer, nil
}
