package project_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/21keshav/IBackendApplication/config"
	. "github.com/21keshav/IBackendApplication/resources/project"
	"github.com/21keshav/IBackendApplication/util/fakes"
)

var _ = Describe("Project", func() {
	var (
		pm              ProjectManager
		fakeMongoClient *fakes.FakeMongoClient
		dbConfig        config.DatabaseDetails
	)
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

	Describe("CreateProject", func() {
		var (
			mongoInsertResult *mongo.InsertOneResult
			primitiveIDByte   primitive.ObjectID
			projectDetails    ProjectDetails
		)
		BeforeEach(func() {
			str := "abc"
			for k, v := range []byte(str) {
				primitiveIDByte[k] = byte(v)
			}
			mongoInsertResult = &mongo.InsertOneResult{
				InsertedID: primitiveIDByte,
			}
			projectDetails = ProjectDetails{ID: "123", SellerID: "WWW"}
			fakeMongoClient.InsertDataReturns(mongoInsertResult, nil)

		})
		It("creates project sucessfull", func() {
			err := pm.CreateProject(projectDetails)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Errors", func() {
			BeforeEach(func() {
				fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("wodo"))
			})

			It("returns a error", func() {
				err := pm.CreateProject(projectDetails)
				Expect(err).To(HaveOccurred())
			})
		})

	})

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
			fakeMongoClient.InsertDataReturns(mongoInsertResult, nil)

		})
		It("creates seller sucessfull", func() {
			err := pm.CreateSeller(seller)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Errors", func() {
			BeforeEach(func() {
				fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("wodo"))
			})

			It("returns a error", func() {
				err := pm.CreateSeller(seller)
				Expect(err).To(HaveOccurred())
			})
		})

	})

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
			fakeMongoClient.InsertDataReturns(mongoInsertResult, nil)

		})
		It("creates seller sucessfull", func() {
			err := pm.CreateBuyer(buyer)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Errors", func() {
			BeforeEach(func() {
				fakeMongoClient.InsertDataReturns(&mongo.InsertOneResult{}, errors.New("wodo"))
			})

			It("returns a error", func() {
				err := pm.CreateBuyer(buyer)
				Expect(err).To(HaveOccurred())
			})
		})

	})
})
