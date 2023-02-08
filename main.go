/*
Problem:
Given a list of connecting flights in any order for a trip,
where a flight is a source and destination airport code,
determine the codes of the first and last airports for the trip.
Serve with route /calculate on port 8080.

Assumptions:
All flights are part of the same trip.
There are no cycles in the trip.
Trips can have 100+ flights.

Examples:
[["SFO", "EWR"]] => ["SFO", "EWR"]
[["ATL", "EWR"], ["SFO", "ATL"]] => ["SFO", "EWR"]
[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]] => ["SFO", "EWR"]

Options:
Iterate flights, track src/dst/all airports, look for missing airport in src/dst trackers
	Time: O(N) average case map lookups, N is number of airports
	Space: O(N) average case map size, N is number of airports
Iterate flights, track src airports, count airport occurrences, look for the 2 airports with count 1
	Time: O(N) average case map lookups, N is number of airports
	Space: O(N) average case map size, N is number of airports

Solution:
Iterate flights, track src airports, count airport occurrences, look for the 2 airports with count 1

Time: O(N) average case map lookups, N is number of airports
Space: O(N) average case map size, N is number of airports
*/

package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
)

const name = "trips"

func handler() http.Handler {
	r := chi.NewRouter()
	r.Use(httprate.Limit(
		100, // Future work: Determine best number
		time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			render.Status(r, http.StatusTooManyRequests)
			render.JSON(w, r, response{
				Error: map[string]any{
					"message": "too many requests",
				},
			})
		}),
	))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.ContentCharset("utf-8"))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response{
			Error: map[string]any{
				"message": "not found",
			},
		})
	}))
	r.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusMethodNotAllowed)
		render.JSON(w, r, response{
			Error: map[string]any{
				"message": "method not allowed",
			},
		})
	}))
	r.Post("/calculate", calculate) // Future work: /v1 versioning
	return r
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(name + ": ")
	s := &http.Server{
		Addr:    ":8080",
		Handler: handler(),
		// Future work: Add TLS
	}
	if err := s.ListenAndServe(); err != nil { // Future work: Add TLS
		log.Fatalf("error: %v", err)
	}
}

// airport is an airport code like SFO.
type airport string

// flight is a pair of airport codes.
// The first is the source. The second is the destination.
type flight [2]airport

type request struct {
	Data struct {
		Flights []flight `json:"flights"`
	} `json:"data"`
}

type response struct {
	Data  map[string]any `json:"data,omitempty"`
	Error map[string]any `json:"error,omitempty"`
}

func calculate(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		if err == io.EOF {
			render.JSON(w, r, response{
				Error: map[string]any{
					"message": "cannot read request body",
				},
			})
			return
		}
		render.JSON(w, r, response{
			Error: map[string]any{
				"message": "cannot decode json",
				"details": err.Error(),
			},
		})
		return
	}
	render.JSON(w, r, response{
		Data: map[string]any{
			"trip": trip(req.Data.Flights),
		},
	})
}

func trip(flights []flight) flight {
	counts := map[airport]int{}
	sources := map[airport]struct{}{}
	for _, f := range flights {
		counts[f[0]]++
		counts[f[1]]++
		sources[f[0]] = struct{}{}
	}
	ends := make([]airport, 0, 2)
	for a, c := range counts {
		if c == 1 {
			ends = append(ends, a)
		}
	}
	src, dst := ends[0], ends[1]
	if _, ok := sources[src]; !ok {
		src, dst = dst, src
	}
	return flight{src, dst}
}
