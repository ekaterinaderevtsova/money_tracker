package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"moneytracker/internal/domain"
	"net/http"
	"time"

	"github.com/stretchr/testify/require"
)

func (s *TrackerTestSuite) TestAddAndGetCurrentWeekSpendings() {
	addUrl := s.spendingUrl + "/spendings"

	today := time.Now().Format("2006-01-02")

	addRequestData := domain.DaySpendings{
		Day: today,
		Sum: 1000,
	}
	jsonAddData, err := json.Marshal(addRequestData)
	require.NoError(s.T(), err, "failed to marshal addRequestData")

	req, err := http.NewRequest("POST", addUrl, bytes.NewBuffer(jsonAddData))
	s.Require().NoError(err, "Error creating request")
	req.Header.Set("Content-Type", "application/json")
	s.Require().NoError(err, "Error marshaling JSON")
	client := &http.Client{}

	for i := 0; i < 2; i++ {
		resp, err := client.Do(req)
		s.Require().NoError(err, "Error to execute request")
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		s.Require().NoError(err, "Error reading response body")

		fmt.Println("StatusCode: ", resp.StatusCode)
		fmt.Printf("Response body: '%s'\n", string(body))

		// Since your handler returns JSON strings, unmarshal as string
		var responseMessage string
		err = json.Unmarshal(body, &responseMessage)
		s.Require().NoError(err, "Error unmarshaling JSON")

		// Assert the expected values
		s.Require().Equal(http.StatusCreated, resp.StatusCode, "Expected 201 status")
		s.Require().Equal("spending added", responseMessage, "Expected success message")
	}

	time.Sleep(1 * time.Millisecond)

	getUrl := fmt.Sprintf("%s/spendings?date=%s", s.spendingUrl, today)

	getReq, err := http.NewRequest("GET", getUrl, nil)
	s.Require().NoError(err, "Error creating GET request")

	getResp, err := client.Do(getReq)
	s.Require().NoError(err, "Error executing GET request")
	defer getResp.Body.Close()

	getBody, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err, "Error reading GET response body")

	fmt.Println("GET StatusCode: ", getResp.StatusCode)
	fmt.Printf("GET Response body: '%s'\n", string(getBody))

	type DaySpendingsResponse struct {
		Date      string `json:"date"`
		DayOfWeek string `json:"day_of_week"`
		Sum       int32  `json:"sum"`
	}

	type WeeklySpendings struct {
		DaySpendings []DaySpendingsResponse `json:"daySpendings"`
		Total        int32                  `json:"total"`
	}

	var weekData WeeklySpendings
	err = json.Unmarshal(getBody, &weekData)
	s.Require().NoError(err, "Error unmarshaling GET response")

	s.Require().Equal(http.StatusOK, getResp.StatusCode, "Expected 200 OK from GET")

	var found bool
	for _, d := range weekData.DaySpendings {
		if d.Date == today {
			found = true
			s.Require().Equal(addRequestData.Sum*2, d.Sum, "Day sum mismatch")
			break
		}
	}
	s.Require().True(found, "Expected today's date in weekly spendings")
	s.Require().Equal(addRequestData.Sum*2, weekData.Total, "Total mismatch")
}
