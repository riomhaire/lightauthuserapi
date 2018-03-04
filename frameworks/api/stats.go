package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/thoas/stats"
)

// HandleStatistics - return json or prometheus acceptable metrics
// depending on the accept type - default json
func (r *RestAPI) HandleStatistics(w http.ResponseWriter, req *http.Request) {
	requestType := req.Header.Get("Accept")
	requestType = strings.ToLower(requestType)

	// Get Statistics
	metrics := r.Statistics.Data()

	if strings.Contains(requestType, "text/plain") {
		w.Header().Set("Content-Type", "text/plain")
		metricsResponse := r.metricsToPrometheus(metrics)
		w.Write([]byte(metricsResponse))
	} else {
		//Return JSON Version
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(metrics)
		w.Write(b)
	}

}

// Takes the metrics data structure and converts it to string
func (r *RestAPI) metricsToPrometheus(metrics *stats.Data) string {
	var buffer bytes.Buffer

	buffer.WriteString("# HELP lightauthuserapi_uptime_sec How many seconds app has been up.\n")
	buffer.WriteString("# TYPE lightauthuserapi_uptime_sec counter\n")
	buffer.WriteString(fmt.Sprintf("lightauthuserapi_uptime_sec %v\n", metrics.UpTimeSec))
	buffer.WriteString("\n")

	buffer.WriteString("# HELP lightauthuserapi_total_response_time_sec Total time spent in handling requests.\n")
	buffer.WriteString("# TYPE lightauthuserapi_total_response_time_sec counter\n")
	buffer.WriteString(fmt.Sprintf("lightauthuserapi_total_response_time_sec %v\n", metrics.TotalResponseTimeSec))
	buffer.WriteString("\n")

	buffer.WriteString("# HELP lightauthuserapi_average_response_time_sec Average time spent in handling requests.\n")
	buffer.WriteString("# TYPE lightauthuserapi_average_response_time_sec guage\n")
	buffer.WriteString(fmt.Sprintf("lightauthuserapi_average_response_time_sec %v\n", metrics.AverageResponseTimeSec))
	buffer.WriteString("\n")

	// Work around for bug in underlying stats library code
	calls := 0

	// Iterate through individual request counts
	if len(metrics.TotalStatusCodeCount) > 0 {
		for statuskey := range metrics.TotalStatusCodeCount {
			tally := metrics.TotalStatusCodeCount[statuskey]
			buffer.WriteString(fmt.Sprintf("# HELP lightauthuserapi_response_status_%v Total Number of Requests returning http status %v\n", statuskey, statuskey))
			buffer.WriteString(fmt.Sprintf("# TYPE lightauthuserapi_response_status_%v counter\n", statuskey))
			buffer.WriteString(fmt.Sprintf("lightauthuserapi_response_status_%v %v\n", statuskey, tally))
			buffer.WriteString("\n")

			calls = calls + tally
		}

	}
	buffer.WriteString("# HELP lightauthuserapi_response_total_count Total Number of Requests.\n")
	buffer.WriteString("# TYPE lightauthuserapi_response_total_count counter\n")
	buffer.WriteString(fmt.Sprintf("lightauthuserapi_response_total_count %v\n", calls))
	buffer.WriteString("\n")

	rates := r.MetricsRegistry.Current()
	buffer.WriteString("# HELP lightauthuserapi_response_rate_per_second Requests/Second \n")
	buffer.WriteString("# TYPE lightauthuserapi_response_rate_per_second guage\n")
	buffer.WriteString(fmt.Sprintf("lightauthuserapi_response_rate_per_second %v\n", rates.PerSecond))
	buffer.WriteString("\n")

	buffer.WriteString("# HELP lightauthuserapi_response_rate_per_minute Requests/Minute \n")
	buffer.WriteString("# TYPE lightauthuserapi_response_rate_per_minute guage\n")
	buffer.WriteString(fmt.Sprintf("lightauthuserapi_response_rate_per_minute %v\n", rates.PerMinute))
	buffer.WriteString("\n")

	buffer.WriteString("# HELP lightauthuserapi_response_rate_per_hour Requests/Hour \n")
	buffer.WriteString("# TYPE lightauthuserapi_response_rate_per_hour guage\n")
	buffer.WriteString(fmt.Sprintf("lightauthuserapi_response_rate_per_hour %v\n", rates.PerHour))
	buffer.WriteString("\n")

	buffer.WriteString(fmt.Sprintf("\n"))

	return buffer.String()
}

/*

lightauthuserapi_response_status_200 10
lightauthuserapi_response_status_401 1
lightauthuserapi_response_total_count 11

*/
