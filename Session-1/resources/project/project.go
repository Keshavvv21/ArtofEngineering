package project_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo" // Ginkgo for BDD test structure
	. "github.com/onsi/gomega" // Gomega for expressive assertions

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/21keshav/IBackendApplication/config"
	. "github.com/21keshav/IBackendApplication/resources/project"
	"github.com/21keshav/IBackendApplication/util/fakes"
)

var _ = Describe("ProjectManager", func() {
	var (
		pm              ProjectManager         // System under test
		fakeMongoClient *fakes.FakeMongoClient // Fake Mongo client
		dbConfig        config.DatabaseDetails // Config passed to ProjectManager
	)

	// Runs before each test case
	BeforeEach(func() {
		dbConfig = config.DatabaseDetails{
			BuyersDBName:   "buyers",
			SellersDBName:  "sellers",
			ProjectDBName:  "projects",
			CollectionName: "collection",
		}
		ctx := context.TODO()
		fakeMongoClient = &fakes.FakeMongoClient{}
		pm = NewProjectManager(fakeMongoClient, ctx, dbConfig)
	})

	// --- CreateProject ---
	Describe("CreateProject", func() {
		var projectDetails ProjectDetails

		BeforeEach(func() {
			projectDetails = ProjectDetails{ID: "p1", SellerID: "s1"}
			// Simulate successful insert
			fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, nil)
		})

		It("should create project successfully", func() {
			err := pm.CreateProject(projectDetails)
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeMongoClient.InsertDataCallCount()).To(Equal(1))
		})

		It("should return error when insert fails", func() {
			fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("insert failed"))
			err := pm.CreateProject(projectDetails)
			Expect(err).To(HaveOccurred())
		})
	})

	// --- CreateSeller ---
	Describe("CreateSeller", func() {
		var seller Seller

		BeforeEach(func() {
			seller = Seller{ID: "s1", SellerID: "seller1"}
			fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, nil)
		})

		It("should create seller successfully", func() {
			err := pm.CreateSeller(seller)
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeMongoClient.InsertDataCallCount()).To(Equal(1))
		})

		It("should return error when insert fails", func() {
			fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("insert failed"))
			err := pm.CreateSeller(seller)
			Expect(err).To(HaveOccurred())
		})
	})

	// --- CreateBuyer ---
	Describe("CreateBuyer", func() {
		var buyer Buyer

		BeforeEach(func() {
			buyer = Buyer{ID: "b1", BuyerID: "buyer1"}
			fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, nil)
		})

		It("should create buyer successfully", func() {
			err := pm.CreateBuyer(buyer)
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeMongoClient.InsertDataCallCount()).To(Equal(1))
		})

		It("should return error when insert fails", func() {
			fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("insert failed"))
			err := pm.CreateBuyer(buyer)
			Expect(err).To(HaveOccurred())
		})
	})

	// --- GetBuyer ---
	Describe("GetBuyer", func() {
		It("should return buyer when found", func() {
			expectedBuyer := Buyer{ID: "b1", BuyerID: "buyer1"}
			// Fake returns no error, fills buyer into result
			fakeMongoClient.FindObjectStub = func(dbName, coll string, filter, result interface{}) error {
				res := result.(*Buyer)
				*res = expectedBuyer
				return nil
			}

			buyer, err := pm.GetBuyer("b1")
			Expect(err).ToNot(HaveOccurred())
			Exp
