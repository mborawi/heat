package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/mborawi/forest/backend/config"
	"github.com/mborawi/forest/backend/models"
)

var db *gorm.DB
var conf config.Config

func main() {
	config.ReadConfig(&conf)

	con := fmt.Sprintf("user=%s password=%s dbname=%s port=%s  sslmode=disable",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.DbName,
		conf.Database.Port)

	var err error
	db, err = gorm.Open("postgres", con)
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}
	defer db.Close()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/leaves/list", listleavesHandler)
	router.HandleFunc("/api/leaves/{year:[0-9]+}/", listleavesHandler)
	router.HandleFunc("/api/branch/list", listBranches)
	router.HandleFunc("/api/department/{id:[0-9]+}", listDepartments)
	router.HandleFunc("/api/availability/{yrs:[0-9]+}", listAvailability)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../frontEnd/dist/static/"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "../frontEnd/dist/index.html") })
	router.NotFoundHandler = http.HandlerFunc(notFound)
	log.Printf("Launching Server Now on %s...", conf.Server.Port)
	log.Fatal(http.ListenAndServe(conf.Server.Port, router))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	log.Println("Flux Reactor not ready doc!...")
}

func listDepartments(w http.ResponseWriter, r *http.Request) {
	deps := []models.Department{}
	vars := mux.Vars(r)
	br_id, _ := strconv.Atoi(vars["id"])
	db.
		Order("id").
		Where("branch_id = ?", br_id).
		Find(&deps)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deps)
}

func listBranches(w http.ResponseWriter, r *http.Request) {
	brs := []models.Branch{}
	db.Order("id").Find(&brs)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brs)
}

func listleavesHandler(w http.ResponseWriter, r *http.Request) {
	lvs := []models.LeaveType{}
	db.
		Order("id").
		Find(&lvs)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lvs)
}

func listAvailability(w http.ResponseWriter, r *http.Request) {
	res := team_result{}
	res.PDays = make(map[string]uint)
	res.UDays = make(map[string]uint)
	res.Year = time.Now().Year()
	vars := mux.Vars(r)
	yrs, _ := strconv.Atoi(vars["yrs"])
	res.Title = fmt.Sprintf("Direct Reports Absences for %d years as of %d",
		yrs, res.Year)
	res.FileTitle = strings.Replace(res.Title, " ", "_", 0)

	st := time.Date(res.Year-10, 1, 1, 0, 0, 0, 0, time.UTC)
	ft := time.Date(res.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	dows := []DayCounts{}
	db.
		Table("leaves").
		Select("EXTRACT(dow FROM leave_date) as DOW, to_char(leave_date,'day') as Day, count(leave_date) as Count").
		Where("leave_date >= ?", st).
		Where("leave_date < ?", ft).
		Where("EXTRACT(dow FROM leave_date) IN (1,2,3,4,5)").
		Group("EXTRACT(dow FROM leave_date),to_char(leave_date,'day')").
		Order("Count DESC, DOW").
		Scan(&dows)
	res.Dows = dows

	type cc struct {
		Dom    string
		Pcount uint
		Ucount uint
	}
	tcs := []cc{}
	db.Raw("SELECT * FROM team_leaves(?)", yrs).Scan(&tcs)
	for _, t := range tcs {
		if t.Dom == "29-02" && !isLeap(res.Year) {
			continue
		}
		if t.Pcount != 0 {
			res.PDays[t.Dom] = t.Pcount
			res.PTotal += t.Pcount
		}
		if t.Ucount != 0 {
			res.UDays[t.Dom] = t.Ucount
			res.UTotal += t.Ucount
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func isLeap(year int) bool {
	return year%400 == 0 || year%4 == 0 && year%100 != 0
}
