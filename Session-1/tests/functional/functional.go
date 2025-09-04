// +build functional

package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/21keshav/IBackendApplication/config"
	"github.com/21keshav/IBackendApplication/controller"
	"github.com/21keshav/IBackendApplication/resources/bidManager"
	"github.com/21keshav/IBackendApplication/resources/project"
	"github.com/21keshav/IBackendApplication/util"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func setupServer() *echo.Echo {
	e := echo.New()

	// Fake config for testing (use test DBs if possible)
	conf := config.Config{
		Database: struct {
			Server string
			Port   string
		}{
			Server: "localhost",
			Port:   "27017",
		},
		DatabaseDetails: config.DatabaseDetails{
			BuyersDBName:   "buyersDB_test",
			SellersDBName:  "sellersDB_test",
			ProjectDBName:  "projectsDB_test",
			CollectionName: "bids_test",
		},
	}

	// Mongo client + managers
	mongoURL := conf.Database.Server + ":" + conf.Database.Port
	mongoClient := util.NewMongoClient(context.TODO(), mongoURL)
	ctx := context.TODO()

	projectManager := project.NewProjectManager(mongoClient, ctx, conf.DatabaseDetails)
	bidMgr := bidManager.NewBidManager(projectManager, ctx)

	c := controller.NewController(bidMgr, projectManager)
	c.AttachHandlers(e)

	return e
}

func TestFunctionalFlow(t *testing.T) {
	e := setupServer()
	server := httptest.NewServer(e)
	defer server.Close()

	// --- 1. Create Seller ---
	seller := map[string]string{"Name": "TestSeller"}
	body, _ := json.Marshal(seller)
	res, err := http.Post(server.URL+"/create-seller", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// --- 2. Create Buyer ---
	buyer := map[string]string{"Name": "TestBuyer"}
	body, _ = json.Marshal(buyer)
	res, err = http.Post(server.URL+"/create-buyer", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// --- 3. Create Project ---
	projectPayload := map[string]string{"Name": "Test Project", "Description": "Demo project"}
	body, _ = json.Marshal(projectPayload)
	res, err = http.Post(server.URL+"/create-project", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// --- 4. Get Projects ---
	res, err = http.Get(server.URL + "/get-projects")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	t.Logf("Projects: %s", string(data))

	// Assume projectID = "p123" and buyerID = "b101" exist (you may parse from response)
	projectID := "p123"
	buyerID := "b101"

	// --- 5. Place a Bid ---
	bidPayload := map[string]interface{}{
		"BuyerID": buyerID,
		"Amount":  1000,
	}
	body, _ = json.Marshal(bidPayload)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/update-bid?projectID="+projectID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// --- 6. Compute Winning Bid ---
	res, err = http.Post(server.URL+"/compute-bid?projectID="+projectID, "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	data, _ = ioutil.ReadAll(res.Body)
	res.Body.Close()
	t.Logf("Winning Bidder: %s", string(data))
}
