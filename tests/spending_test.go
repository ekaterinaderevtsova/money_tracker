package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"moneytracker/internal/domain"
	httpdto "moneytracker/internal/transport/dto"
	"net/http"
	"time"

	"github.com/stretchr/testify/require"
)

func (s *TrackerTestSuite) authenticateUser(login, password string) string {
	credentials := httpdto.Credentials{
		Login:    login,
		Password: password,
	}
	jsonData, err := json.Marshal(credentials)
	require.NoError(s.T(), err, "failed to marshal json data")

	req, err := http.NewRequest("POST", s.loginURL, bytes.NewBuffer(jsonData))
	s.Require().NoError(err, "failed to create login request")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	s.Require().NoError(err, "failed to execute login request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "failed to read login response")

	var tokenResp httpdto.TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	s.Require().NoError(err, "failed to unmarshal token response")

	return tokenResp.AccessToken
}

func (s *TrackerTestSuite) TestAddAndGetCurrentWeekSpendings() {
	token := s.authenticateUser("login", "password")

	today := time.Now().Format("2006-01-02")

	addRequestData := domain.DailyExpense{
		Date:   today,
		Amount: 1000,
	}
	jsonAddData, err := json.Marshal(addRequestData)
	require.NoError(s.T(), err, "failed to marshal addRequestData")

	addReq, err := http.NewRequest("POST", s.expensesURL, bytes.NewBuffer(jsonAddData))
	s.Require().NoError(err, "Error creating request")
	addReq.Header.Set("Content-Type", "application/json")
	addReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	for i := 0; i < 2; i++ {
		resp, err := client.Do(addReq)
		s.Require().NoError(err, "Error to execute request")
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		s.Require().NoError(err, "Error reading response body")

		fmt.Println("StatusCode: ", resp.StatusCode)
		fmt.Printf("Response body: '%s'\n", string(body))

		var responseMessage string
		err = json.Unmarshal(body, &responseMessage)
		s.Require().NoError(err, "Error unmarshaling JSON")

		s.Require().Equal(http.StatusCreated, resp.StatusCode, "Expected 201 status")
		s.Require().Equal("spending added", responseMessage, "Expected success message")
	}

	time.Sleep(1 * time.Millisecond)

	getURL := fmt.Sprintf("%s?date=%s", s.expensesURL+"/weekly", today)

	getReq, err := http.NewRequest("GET", getURL, nil)
	s.Require().NoError(err, "Error creating GET request")

	getReq.Header.Set("Authorization", "Bearer "+token)

	getResp, err := client.Do(getReq)
	s.Require().NoError(err, "Error executing GET request")
	defer getResp.Body.Close()

	getBody, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err, "Error reading GET response body")

	fmt.Println("GET StatusCode: ", getResp.StatusCode)
	fmt.Printf("GET Response body: '%s'\n", string(getBody))

	var weekData httpdto.WeeklyExpense
	err = json.Unmarshal(getBody, &weekData)
	s.Require().NoError(err, "Error unmarshaling GET response")

	s.Require().Equal(http.StatusOK, getResp.StatusCode, "Expected 200 OK from GET")

	today = time.Now().Format("02-01")

	var found bool
	for _, d := range weekData.DailyExpenses {
		if d.Date == today {
			found = true
			s.Require().Equal(addRequestData.Amount*2, d.Amount, "Day sum mismatch")
			break
		}
	}
	s.Require().True(found, "Expected today's date in weekly expenses")
	s.Require().Equal(addRequestData.Amount*2, weekData.TotalAmount, "Total mismatch")
}
