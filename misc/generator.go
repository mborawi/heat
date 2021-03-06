package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/icrowley/fake"
)

func createFakeEmployees(N int) {
	f, err := os.Create("employees.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	w := bufio.NewWriter(f)
	for i := 1; i <= N; i++ {
		fname := fake.FirstName()
		lname := fake.LastName()
		email := fmt.Sprintf("%s.%s@fakemail.com", fname, lname)
		phone := fake.Phone()
		mm := time.Month(rand.Intn(12) + 1)
		sd := time.Date(time.Now().Year()-rand.Intn(20)-1,
			mm, rand.Intn(28)+1, 0, 0, 0, 0, time.UTC)
		jobTitle := fake.JobTitle()
		branchId := rand.Intn(10) + 1
		managerID := rand.Intn(i) + 1
		if i == 1 {
			managerID = 0
		}
		if i != i && managerID == 0 {
			managerID = 1
		}
		fmt.Fprintf(w, "%d,%s,%s,%s,%s,%s,%d,%d,%s\n",
			i, fname, lname, email, phone, jobTitle, branchId, managerID, sd.Format("2006-01-02 15:04:05"))
	}
	w.Flush()
	f.Close()
}

func createRanges(N, stYr, fnYr int) {
	cats := []string{
		"Carers",
		"Annual",
		"Long Service",
		"Sick",
		"Leave without pay",
		"Wellbeing",
		"Purchased",
		"Maternity",
		"Witness",
		"Miscellaneous",
		"Not Assigned",
		"Workers Compensation",
		"Study",
		"Bereavement",
	}
	lveNames := []string{"Planned", "Unplanned"}

	f, err := os.Create("LeaveRanges.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	w := bufio.NewWriter(f)
	// assume 4 leaves year, with an average of 3 days per leave for each employee
	Leaves := 4
	c := 1 // serial number for records
	for i := 0; i < N; i++ {
		for y := stYr; y <= fnYr; y++ {
			for l := 0; l < Leaves; l++ {
				m := rand.Intn(13) + 1
				d := rand.Intn(28) + 1
				leaveStarts := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
				duration := time.Hour * 24 * time.Duration(rand.Intn(8)+1)
				leaveEnds := leaveStarts.Add(duration)
				rr := rand.Intn(len(cats))
				fmt.Fprintf(w, "%d,%d,%s,%s,%.2f,%s,%s,PLWOP TCAS\n",
					c, i+1, leaveStarts.Format("2006-01-02"), leaveEnds.Format("2006-01-02"), rand.Float64()*7, lveNames[rand.Intn(2)], cats[rr])
				c++
			}
		}
	}
	w.Flush()
	f.Close()

}

func main() {
	NumberOfEmployees := 20000
	StartYear := 2008
	FinishYear := 2018
	rand.Seed(time.Now().Unix())
	fmt.Println("creating fake data")
	createFakeEmployees(NumberOfEmployees)
	createRanges(NumberOfEmployees, StartYear, FinishYear)
}
