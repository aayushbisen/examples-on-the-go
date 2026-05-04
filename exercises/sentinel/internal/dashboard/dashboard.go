package dashboard

import (
	"fmt"
	"my-go-learning/sentinel/internal/storage"
	"net/http"
	"strings"
)

type SiteStatus struct {
	ID     string
	Status bool
}

func GetStatus(s *storage.Storage) ([]SiteStatus, error) {
	q := "SELECT website_id, status FROM results GROUP BY website_id ORDER BY timestamp DESC"
	var ss []SiteStatus
	rows, err := s.DB.Query(q)

	// err check
	if err != nil {
		return nil, err
	}
	//
	// defer close of rows
	defer rows.Close()
	// iterate using next
	//
	//
	// unpack with scan on row
	for rows.Next() {
		// We are now looking at a single row of data
		//
		var id string
		var status bool

		// The order of variables must match the order of columns in your SELECT query
		err := rows.Scan(&id, &status)
		if err != nil {
			return nil, err
		}

		ss = append(ss, SiteStatus{ID: id, Status: status})
		// fmt.Println("Website:", id, "Status:", status)
	}

	return ss, nil
}

func Handler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sss, err := GetStatus(s)

		if err != nil {
			fmt.Printf("Error %s", err)
			return
		}

		var html strings.Builder
		html.WriteString("<h1>Sentinel Dashboard</h1><table>")
		for _, res := range sss {
			fmt.Fprintf(&html, "<tr><td>%s</td><td>%t</td></tr>", res.ID, res.Status)
		}
		html.WriteString("</table>")
		fmt.Fprintf(w, "%s", html.String())
	}
}
