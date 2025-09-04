package bidManager

import (
	"context"
	"math"

	"github.com/21keshav/IBackendApplication/resources/project"
	"github.com/golang/glog"
)

// BidManager defines the contract for bid-related operations.
// It encapsulates the ability to place bids and compute the winning bid.
type BidManager interface {
	// ComputeBID determines the buyer with the lowest bid for a given project.
	ComputeBID(projectID string) (project.Buyer, error)

	// DoBID places a new bid for a given project.
	DoBID(projectID string, bid project.BID) error
}

// BidManagerManagerImpl is the concrete implementation of the BidManager interface.
// It relies on ProjectManager to interact with projects and buyers stored in the database.
type BidManagerManagerImpl struct {
	projectManager project.ProjectManager // Handles project & buyer persistence
	ctx            context.Context        // Context for database operations
}

// NewBidManager initializes and returns a new BidManager instance.
// It wires together the BidManager with a ProjectManager dependency.
func NewBidManager(projectManager project.ProjectManager, ctx context.Context) BidManager {
	return &BidManagerManagerImpl{
		projectManager,
		ctx,
	}
}

// DoBID inserts or updates a bid for a project by delegating
// the operation to the ProjectManager.
func (bd *BidManagerManagerImpl) DoBID(projectID string, bid project.BID) error {
	glog.Info("Do-bid-projects")
	defer glog.Info("do-bid-completed")

	// Delegate updating the project with the new bid to ProjectManager
	err := bd.projectManager.UpdateProject(projectID, bid)
	return err
}

// ComputeBID determines the winning buyer for a given project.
// The winning buyer is the one who has placed the minimum bid.
//
// Steps:
// 1. Retrieve the project details from ProjectManager.
// 2. Iterate over all bids to find the lowest bid amount.
// 3. Fetch and return the buyer associated with that minimum bid.
func (bd *BidManagerManagerImpl) ComputeBID(projectID string) (project.Buyer, error) {
	glog.Info("compute-projects")
	defer glog.Info("compute-completed")

	// Step 1: Get the project details
	currentProject, err := bd.projectManager.GetProject(projectID)
	if err != nil {
		return project.Buyer{}, err
	}

	// Step 2: Find the minimum bid
	min := math.MaxInt64
	var minBid project.BID
	for _, bid := range currentProject.BIDS {
		if bid.Amount < min {
			min = bid.Amount
			minBid = bid
		}
	}

	// Step 3: Fetch the buyer corresponding to the lowest bid
	buyer, err := bd.projectManager.GetBuyer(minBid.BuyerID)
	if err != nil {
		return project.Buyer{}, err
	}

	return buyer, nil
}
