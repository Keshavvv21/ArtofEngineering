package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/21keshav/IBackendApplication/controller"
	"github.com/21keshav/IBackendApplication/resources/project"
	"github.com/labstack/echo"
)

// --- Mock implementations for dependencies ---

type mockBidManager struct {
	doBIDCalled     bool
	doBIDErr        error
	computeCalled   bool
	computeResult   interface{}
	computeErr      error
}

func (m *mockBidManager) DoBID(projectID string, bid project.BID) error {
	m.doBIDCalled = true
	return m.doBIDErr
}

func (m *mockBidManager) ComputeBID(projectID string) (interface{}, error) {
	m.computeCalled = true
	return m.computeResult, m.computeErr
}

type mockProjectManager struct {
	createProjectCalled bool
	createSellerCalled  bool
	createBuyerCalled   bool
	getProjectsCalled   bool

	createProjectErr error
	createSellerErr  error
	createBuyerErr   error
	getProjectsRes   []project.ProjectDetails
	getProjectsErr   error
}

func (m *mockProjectManager) CreateProject(p project.ProjectDetails) error {
	m.createProjectCalled = true
	return m.createProjectErr
}

func (m *mockProjectManager) CreateSeller(s project.Seller) error {
	m.createSellerCalled = true
	return m.createSellerErr
}

func (m *mockProjectManager) CreateBuyer(b project.Buyer) error {
	m.createBuyerCalled = true
	return m.createBuyerErr
}

func (m *mockProjectManager) GetProjects() ([]project.ProjectDetails, error) {
	m.getProjectsCalled = true
	return m.getProjectsRes, m.getProjectsErr
}

// --- Test Suite ---

var _ = Describe("Controller", func() {
	var (
		e          *echo.Echo
		rec        *httptest.ResponseRecorder
		c          controller.Controller
		mockBid    *mockBidManager
		mockProj   *mockProjectManager
	)

	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
		mockBid = &mockBidManager{}
		mockProj = &mockProjectManager{}
		c = controller.NewController(mockBid, mockProj)
	})

	// --- CreateProject ---
	Describe("CreateProject", func() {
		It("should create a project successfully", func() {
			body := `{"Name":"Test Project"}`
			req := httptest.NewRequest(http.MethodPost, "/create-project", bytes.NewBufferString(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.CreateProject(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(mockProj.createProjectCalled).To(BeTrue())
		})

		It("should return 400 on invalid JSON", func() {
			req := httptest.NewRequest(http.MethodPost, "/create-project", bytes.NewBufferString("{invalid"))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.CreateProject(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return 500 when ProjectManager fails", func() {
			mockProj.createProjectErr = errors.New("db error")
			body := `{"Name":"Test Project"}`
			req := httptest.NewRequest(http.MethodPost, "/create-project", bytes.NewBufferString(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.CreateProject(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	// --- CreateSeller ---
	Describe("CreateSeller", func() {
		It("should create a seller successfully", func() {
			body := `{"Name":"Seller1"}`
			req := httptest.NewRequest(http.MethodPost, "/create-seller", bytes.NewBufferString(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.CreateSeller(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(mockProj.createSellerCalled).To(BeTrue())
		})

		It("should return 500 when ProjectManager fails", func() {
			mockProj.createSellerErr = errors.New("db error")
			body := `{"Name":"Seller1"}`
			req := httptest.NewRequest(http.MethodPost, "/create-seller", bytes.NewBufferString(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.CreateSeller(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	// --- CreateBuyer ---
	Describe("CreateBuyer", func() {
		It("should create a buyer successfully", func() {
			body := `{"Name":"Buyer1"}`
			req := httptest.NewRequest(http.MethodPost, "/create-buyer", bytes.NewBufferString(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.CreateBuyer(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(mockProj.createBuyerCalled).To(BeTrue())
		})

		It("should return 500 when ProjectManager fails", func() {
			mockProj.createBuyerErr = errors.New("db error")
			body := `{"Name":"Buyer1"}`
			req := httptest.NewRequest(http.MethodPost, "/create-buyer", bytes.NewBufferString(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.CreateBuyer(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	// --- UpdateBID ---
	Describe("UpdateBID", func() {
		It("should update a bid successfully", func() {
			bid := project.BID{Amount: 100}
			data, _ := json.Marshal(bid)
			req := httptest.NewRequest(http.MethodPut, "/update-bid?projectID=123", bytes.NewBuffer(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.UpdateBID(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(mockBid.doBIDCalled).To(BeTrue())
		})

		It("should return 500 when BidManager fails", func() {
			mockBid.doBIDErr = errors.New("update failed")
			bid := project.BID{Amount: 100}
			data, _ := json.Marshal(bid)
			req := httptest.NewRequest(http.MethodPut, "/update-bid?projectID=123", bytes.NewBuffer(data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			err := c.UpdateBID(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	// --- GetProjects ---
	Describe("GetProjects", func() {
		It("should return projects successfully", func() {
			mockProj.getProjectsRes = []project.ProjectDetails{{Name: "P1"}}
			req := httptest.NewRequest(http.MethodGet, "/get-projects", nil)
			ctx := e.NewContext(req, rec)

			err := c.GetProjects(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(mockProj.getProjectsCalled).To(BeTrue())
		})

		It("should return 500 when ProjectManager fails", func() {
			mockProj.getProjectsErr = errors.New("db fail")
			req := httptest.NewRequest(http.MethodGet, "/get-projects", nil)
			ctx := e.NewContext(req, rec)

			err := c.GetProjects(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	// --- ComputeBID ---
	Describe("ComputeBID", func() {
		It("should compute bid successfully", func() {
			mockBid.computeResult = "Buyer1"
			req := httptest.NewRequest(http.MethodPost, "/compute-bid?projectID=123", nil)
			ctx := e.NewContext(req, rec)

			err := c.ComputeBID(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(mockBid.computeCalled).To(BeTrue())
		})

		It("should return 500 when BidManager fails", func() {
			mockBid.computeErr = errors.New("calc error")
			req := httptest.NewRequest(http.MethodPost, "/compute-bid?projectID=123", nil)
			ctx := e.NewContext(req, rec)

			err := c.ComputeBID(ctx)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
