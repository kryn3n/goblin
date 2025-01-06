package arcgis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/kryn3n/goblin/internal/misc"
)

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties map[string]any `json:"properties"`
	Type       string         `json:"type"`
	Geometry   Geometry       `json:"geometry"`
	Id         int            `json:"id"`
}

type Geometry struct {
	Type        string `json:"type"`
	Coordinates []any  `json:"coordinates"`
}

type GeoJson struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Layer struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Server struct {
	Layers           []Layer `json:"layers"`
	LayerUrl         string
	ObjectIdField    string
	URL              string
	BatchList        [][][]int
	ObjectIds        []int
	LayerIds         []int
	MyFeatures       int
	RecordLimit      int
	RecordCount      int
	ConcurrencyLimit int
	LayerId          int
	MaxRecordCount   int `json:"maxRecordCount"`
	Lock             sync.Mutex
	Wg               sync.WaitGroup
}

const queryPath = "query"

func (s *Server) GetServerInfo() {
	values := url.Values{
		"f": {"pjson"},
	}
	metadataURL := misc.FormatURL(s.URL, "", &values)
	r, err := http.Get(metadataURL)
	misc.Check(err)

	data, err := io.ReadAll(r.Body)
	misc.Check(err)

	json.Unmarshal(data, &s)
}

func (s *Server) GetLayerURL() {
	s.LayerUrl = misc.FormatURL(s.URL, strconv.Itoa(s.LayerId), &url.Values{})
}

func (s *Server) GetRecordCount() {
	var countResponse struct {
		Count int `json:"count"`
	}

	values := url.Values{
		"where":           {"1=1"},
		"returnCountOnly": {"true"},
		"f":               {"pjson"},
	}

	queryURL := misc.FormatURL(s.LayerUrl, queryPath, &values)
	r, err := http.Get(queryURL)
	misc.Check(err)

	data, err := io.ReadAll(r.Body)
	misc.Check(err)

	json.Unmarshal(data, &countResponse)
	s.RecordCount = countResponse.Count
}

func (s *Server) GetObjectIds() {
	var objectIdsResponse struct {
		ObjectIdFieldName string `json:"objectIdFieldName"`
		ObjectIds         []int  `json:"objectIds"`
	}

	values := url.Values{
		"where":         {"1=1"},
		"returnIdsOnly": {"true"},
		"f":             {"json"},
	}

	queryURL := misc.FormatURL(s.LayerUrl, queryPath, &values)

	response, error := http.Get(queryURL)
	misc.Check(error)

	defer response.Body.Close()
	body, error := io.ReadAll(response.Body)
	misc.Check(error)

	json.Unmarshal(body, &objectIdsResponse)
	s.ObjectIdField = objectIdsResponse.ObjectIdFieldName
	s.ObjectIds = objectIdsResponse.ObjectIds
	slices.Sort(s.ObjectIds)
}

func (s *Server) GetBatches() {
	requestList := createBatches(s.ObjectIds, s.RecordLimit)
	s.BatchList = createBatches(requestList, s.ConcurrencyLimit)
}

func createBatches[I any](inputList []I, limit int) [][]I {
	var outputList [][]I
	for i := 0; i < len(inputList); i += limit {
		end := i + limit
		if end > len(inputList) {
			end = len(inputList)
		}
		batch := inputList[i:end]
		outputList = append(outputList, batch)
	}
	return outputList
}

func (s *Server) GetData(fileName string) {
	for index, batch := range s.BatchList {
		s.Wg.Add(len(batch))
		f, error := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		misc.Check(error)
		fmt.Printf("\n***************************\nBatch %d, Requests %d\n", index, len(batch))

		for _, objectIds := range batch {
			go s.GetRecords(objectIds, f)
		}

		s.Wg.Wait()
		fmt.Print("\n***************************\n")
	}
}

func (s *Server) GetRecords(objectIds []int, f *os.File) {
	defer s.Wg.Done()

	var data struct {
		Features []Feature `json:"features"`
	}

	values := url.Values{
		"returnGeometry": {"true"},
		"outFields":      {"*"},
		"f":              {"geojson"},
		"where":          {"1=1"},
		"orderByFields":  {s.ObjectIdField},
		"objectIds":      {misc.SplitToString(objectIds, ",")},
	}

	queryURL := misc.FormatURL(s.LayerUrl, queryPath, &values)

	response, error := http.Get(queryURL)
	misc.Check(error)

	body, error := io.ReadAll(response.Body)
	misc.Check(error)

	json.Unmarshal(body, &data)
	response.Close = true
	response.Body.Close()

	s.Lock.Lock()
	s.MyFeatures += len(data.Features)
	fmt.Printf("\rObject Count Is: %d", s.MyFeatures)

	for _, feature := range data.Features {
		a, error := json.Marshal(feature)
		misc.Check(error)

		_, error = fmt.Fprintf(f, "%s\n", a)
		misc.Check(error)

	}
	s.Lock.Unlock()
}
