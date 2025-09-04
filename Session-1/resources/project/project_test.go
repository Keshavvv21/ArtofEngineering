package project_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo" // Ginkgo BDD testing framework
	. "github.com/onsi/gomega" // Gomega matchers for assertions
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/21keshav/IBackendApplication/config"
	. "github.com/21keshav/IBackendApplication/resources/project"
	"github.com/21keshav/IBackendApplication/util/fakes"
)

var _ = Describe("Project", func() {
	var (
		pm              ProjectManager        // The ProjectManager under test
		fakeMongoClient *fakes.FakeMongoClient // Fake Mongo client used for mocking DB calls
		dbConfig        config.DatabaseDetails // Database config injected into ProjectManager
	)

	// Runs before each test case
	BeforeEach(func() {
		dbConfig = config.DatabaseDetails{
			BuyersDBName:   "buyers",
			SellersDBName:  "sellers",
			ProjectDBName:  "projectDetails",
			CollectionName: "Collections",
		}
		ctx := context.TODO()
		fakeMongoClient = &fakes.FakeMongoClient{}
		pm = NewProjectManager(fakeMongoClient, ctx, dbConfig)
	})

	// --- Tests for CreateProject ---
	Describe("CreateProject", func() {
		var (
			mongoInsertResult *mongo.InsertOneResult // Simulated insert result
			primitiveIDByte   primitive.ObjectID     // Fake MongoDB ObjectID
			projectDetails    ProjectDetails         // Project details used as input
		)

		BeforeEach(func() {
			// Fake MongoDB ObjectID generated from string "abc"
			str := "abc"
			for k, v := range []byte(str) {
				primitiveIDByte[k] = byte(v)
			}

			// Setup mock InsertOneResult
			mongoInsertResult = &mongo.InsertOneResult{
				InsertedID: primitiveIDByte,
			}

			// Initialize a project with dummy values
			projectDetails = ProjectDetails{ID: "123", SellerID: "WWW"}

			// Configure fake MongoClient to return success
			fakeMongoClient.InsertDataReturns(mongoInsertResult, nil)
		})

		It("creates project successfully", func() {
			err := pm.CreateProject(projectDetails)
			Expect(err).ToNot(HaveOccurred()) // Expect no error
		})

		Context("Errors", func() {
			BeforeEach(func() {
				// Simulate Mongo insert failure
				fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("wodo"))
			})

			It("returns an error", func() {
				err := pm.CreateProject(projectDetails)
				Expect(err).To(HaveOccurred()) // Expect error
			})
		})
	})

	// --- Tests for CreateSeller ---
	Describe("CreateSeller", func() {
		var (
			mongoInsertResult *mongo.InsertOneResult
			primitiveIDByte   primitive.ObjectID
			seller            Seller
		)

		BeforeEach(func() {
			str := "abc"
			for k, v := range []byte(str) {
				primitiveIDByte[k] = byte(v)
			}
			mongoInsertResult = &mongo.InsertOneResult{
				InsertedID: primitiveIDByte,
			}
			seller = Seller{ID: "123", SellerID: "WWW"}

			// Configure fake MongoClient to return success
			fakeMongoClient.InsertDataReturns(mongoInsertResult, nil)
		})

		It("creates seller successfully", func() {
			err := pm.CreateSeller(seller)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Errors", func() {
			BeforeEach(func() {
				// Simulate insert failure
				fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("wodo"))
			})

			It("returns an error", func() {
				err := pm.CreateSeller(seller)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// --- Tests for CreateBuyer ---
	Describe("CreateBuyer", func() {
		var (
			mongoInsertResult *mongo.InsertOneResult
			primitiveIDByte   primitive.ObjectID
			buyer             Buyer
		)

		BeforeEach(func() {
			str := "abc"
			for k, v := range []byte(str) {
				primitiveIDByte[k] = byte(v)
			}
			mongoInsertResult = &mongo.InsertOneResult{
				InsertedID: primitiveIDByte,
			}
			buyer = Buyer{ID: "123", BuyerID: "WWW"}

			// Configure fake MongoClient to return success
			fakeMongoClient.InsertDataReturns(mongoInsertResult, nil)
		})

		It("creates buyer successfully", func() {
			err := pm.CreateBuyer(buyer)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Errors", func() {
			BeforeEach(func() {
				// Simulate insert failure
				fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("wodo"))
			})

			It("returns an error", func() {
				err := pm.CreateBuyer(buyer)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
