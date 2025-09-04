package bidManager_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/21keshav/IBackendApplication/resources/bidManager"
	"github.com/21keshav/IBackendApplication/resources/project"
)

// --- Mock ProjectManager ---

type mockProjectManager struct {
	updateCalled   bool
	updateProjectID string
	updateBid      project.BID
	updateErr      error

	getProjectCalled bool
	getProjectID     string
	getProjectRes    project.ProjectDetails
	getProjectErr    error

	getBuyerCalled bool
	getBuyerID     string
	getBuyerRes    project.Buyer
	getBuyerErr    error
}

func (m *mockProjectManager) UpdateProject(projectID string, bid project.BID) error {
	m.updateCalled = true
	m.updateProjectID = projectID
	m.updateBid = bid
	return m.updateErr
}

func (m *mockProjectManager) GetProject(projectID string) (project.ProjectDetails, error) {
	m.getProjectCalled = true
	m.getProjectID = projectID
	return m.getProjectRes, m.getProjectErr
}

func (m *mockProjectManager) GetBuyer(buyerID string) (project.Buyer, error) {
	m.getBuyerCalled = true
	m.getBuyerID = buyerID
	return m.getBuyerRes, m.getBuyerErr
}

// Unused methods to satisfy interface (not tested here)
func (m *mockProjectManager) CreateProject(project.ProjectDetails) error { return nil }
func (m *mockProjectManager) GetProjects() ([]project.ProjectDetails, error) {
	return nil, nil
}
func (m *mockProjectManager) CreateBuyer(project.Buyer) error { return nil }
func (m *mockProjectManager) CreateSeller(project.Seller) error { return nil }

// --- Test Suite ---

var _ = Describe("BidManager", func() {
	var (
		mockPM *mockProjectManager
		bm     BidManager
	)

	BeforeEach(func() {
		mockPM = &mockProjectManager{}
		bm = NewBidManager(mockPM, context.TODO())
	})

	// --- DoBID Tests ---
	Describe("DoBID", func() {
		It("should call UpdateProject successfully", func() {
			bid := project.BID{ID: "b1", BuyerID: "buyer1", Amount: 100}
			err := bm.DoBID("p1", bid)

			Expect(err).To(BeNil())
			Expect(mockPM.updateCalled).To(BeTrue())
			Expect(mockPM.updateProjectID).To(Equal("p1"))
			Expect(mockPM.updateBid).To(Equal(bid))
		})

		It("should return error if UpdateProject fails", func() {
			mockPM.updateErr = errors.New("update failed")
			bid := project.BID{ID: "b1", BuyerID: "buyer1", Amount: 100}

			err := bm.DoBID("p1", bid)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("update failed"))
		})
	})

	// --- ComputeBID Tests ---
	Describe("ComputeBID", func() {
		It("should return the buyer with the lowest bid", func() {
			// Project with two bids
			projectDetails := project.ProjectDetails{
				ID:   "p1",
				BIDS: map[string]project.BID{
					"b1": {ID: "b1", BuyerID: "buyer1", Amount: 200},
					"b2": {ID: "b2", BuyerID: "buyer2", Amount: 100}, // lowest
				},
			}
			mockPM.getProjectRes = projectDetails
			mockPM.getBuyerRes = project.Buyer{ID: "buyer2", BuyerName: "LowestBidder"}

			buyer, err := bm.ComputeBID("p1")

			Expect(err).To(BeNil())
			Expect(mockPM.getProjectCalled).To(BeTrue())
			Expect(mockPM.getBuyerCalled).To(BeTrue())
			Expect(mockPM.getBuyerID).To(Equal("buyer2"))
			Expect(buyer.BuyerName).To(Equal("LowestBidder"))
		})

		It("should return error if GetProject fails", func() {
			mockPM.getProjectErr = errors.New("db error")

			_, err := bm.ComputeBID("p1")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("db error"))
		})

		It("should return error if GetBuyer fails", func() {
			projectDetails := project.ProjectDetails{
				ID:   "p1",
				BIDS: map[string]project.BID{
					"b1": {ID: "b1", BuyerID: "buyer1", Amount: 150},
				},
			}
			mockPM.getProjectRes = projectDetails
			mockPM.getBuyerErr = errors.New("buyer lookup failed")

			_, err := bm.ComputeBID("p1")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("buyer lookup failed"))
		})

		It("should handle empty bids gracefully", func() {
			projectDetails := project.ProjectDetails{
				ID:   "p1",
				BIDS: map[string]project.BID{}, // no bids
			}
			mockPM.getProjectRes = projectDetails

			// Expect error because no buyer can be found
			_, err := bm.ComputeBID("p1")

			Expect(err).To(HaveOccurred())
		})
	})
})
