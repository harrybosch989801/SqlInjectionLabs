package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/srinathgs/mysqlstore"
)

type Env struct {
	db       *sql.DB
	logger   *log.Logger
	sessions *mysqlstore.MySQLStore
}

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(localhost:32000)/RETIREMENTAPP")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	//set up logging
	file, err := os.Create("console.log")
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, file)

	var sessions *mysqlstore.MySQLStore

	sessions, err = mysqlstore.NewMySQLStore("root:root@tcp(localhost:32000)/RETIREMENTAPP?parseTime=true&loc=Local",
		"sessions", "/", 3600, []byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		panic(err)
	}
	defer sessions.Close()

	env := &Env{
		db:       db,
		logger:   log.New(mw, "RETIREMENTAPP ", log.LstdFlags),
		sessions: sessions}

	http.HandleFunc("/sources", env.getSources)
	http.HandleFunc("/participantdetails", env.getParticipantDetails)
	http.HandleFunc("/submitdeferral", env.submitDeferral)
	http.HandleFunc("/auth", env.auth)
	http.ListenAndServe(":8080", nil)

}

/*
Shallow implementation.  Not going to bother inactivating old deferrals.
Demonstration purposes only :)
*/
func (env *Env) submitDeferral(w http.ResponseWriter, r *http.Request) {
	var submitDeferralRequest SubmitDeferralRequest

	err := json.NewDecoder(r.Body).Decode(&submitDeferralRequest)
	//In a real application, we would probably want to build something to
	//handle errors in a more hollistic way.
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		env.logger.Println(err)
		return
	}
	env.logger.Printf("SubmitDeferralRequest for %s for External Plan ID %s\n",
		submitDeferralRequest.Username, submitDeferralRequest.ExternalPlanId)

	enrollmentID, err := env.createEnrollments(submitDeferralRequest)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		env.logger.Println(err)
		return
	}

	err = env.createDeferrals(submitDeferralRequest, enrollmentID)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		env.logger.Println(err)
		return
	}

	env.logger.Printf("Deferral successfully submitted for %s for External Plan ID %s\n",
		submitDeferralRequest.Username, submitDeferralRequest.ExternalPlanId)
	w.Write([]byte("Submitted"))
}

func (env *Env) getCustomerId(username string) (string, error) {

	rows, err := env.db.Query("select customerid from customers where name = '" + username + "'")
	if err != nil {
		return "Error retriving customerid", err
	}
	defer rows.Close()
	var customerid string
	for rows.Next() {
		err := rows.Scan(&customerid)
		if err != nil {
			return "Error retriving customerid", err
		}
	}

	return customerid, err

}

func (env *Env) getPlanId(externalid string) (string, error) {

	rows, err := env.db.Query("select planid from plans where externalid = '" + externalid + "'")
	if err != nil {
		return "Error retriving planid", err
	}
	defer rows.Close()
	var planid string
	for rows.Next() {
		err := rows.Scan(&planid)
		if err != nil {
			return "Error retriving planid", err
		}
	}

	return planid, err

}

func (env *Env) createDeferrals(sdr SubmitDeferralRequest, enrollmentID string) error {
	var err error

	for _, deferral := range sdr.DeferralRequest {
		id := uuid.New()
		//in real life use string builders please :)
		//lazy mode for demo :)
		sql := "insert into deferrals "
		sql += "(deferralid, sourcename, deductamount, createtime, enrollmentid) "
		sql += "values ("
		sql += "'" + id.String() + "', "
		sql += "'" + deferral.Source + "', "
		sql += "'" + strconv.Itoa(deferral.DeductAmount) + "', "
		sql += "NOW(), "
		sql += "'" + enrollmentID + "'"
		sql += ")"

		result, err := env.db.Exec(sql)
		if err != nil {
			return err
		}
		updated, err := result.RowsAffected()
		if err != nil {
			return err
		}

		env.logger.Printf("Inserted %d rows into deferrals table", updated)
	}

	return err

}
func (env *Env) createEnrollments(sdr SubmitDeferralRequest) (string, error) {

	customerid, err := env.getCustomerId(sdr.Username)
	if err != nil {
		return "Error fetching customerid", err
	}

	planid, err := env.getPlanId(sdr.ExternalPlanId)
	if err != nil {
		return "Error fetching plan", err
	}

	id := uuid.New()

	//in real life use string builders please :)
	//lazy mode for demo :)
	sql := "insert into enrollments "
	sql += "(enrollmentid, deductionmethod, planid, status, createtime, customerid) "
	sql += "values ("
	sql += "'" + id.String() + "', "
	sql += "'" + strconv.Itoa(sdr.DeductMethod) + "', "
	sql += "'" + planid + "', "
	sql += "'ACTIVE', "
	sql += "NOW(), "
	sql += "'" + customerid + "'"
	sql += ")"

	result, err := env.db.Exec(sql)
	if err != nil {
		return "Error inserting enrollment", err
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return "Error fetching update counts", err
	}

	env.logger.Printf("Inserted %d rows into enrollments table", updated)

	return id.String(), err

}

func (env *Env) getParticipantDetails(w http.ResponseWriter, r *http.Request) {

	var participantDetailRequest ParticipantDetailRequest
	err := json.NewDecoder(r.Body).Decode(&participantDetailRequest)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		env.logger.Println(err)
		return
	}
	env.logger.Printf("Participant Detail Request for %s for External Plan ID %s\n", participantDetailRequest.Username, participantDetailRequest.ExternalId)

	sql := "select p.externalid, s.sourcename, s.sourcetype, p.planname, d.deductamount, e.deductmethod from enrollments e "
	sql += "join customer c on c.customerid = e.customerid "
	sql += "join deferrals d on d.enrollmentid = e.enrollmentid "
	sql += "join sources s on s.sourcename = d.sourcename "
	sql += "join plans p on p.planid = e.planid"
	sql += "where c.name = '" + participantDetailRequest.Username + "' "
	if participantDetailRequest.ExternalId != "" {
		sql += " and p.externalid = '" + participantDetailRequest.ExternalId + "'"
	}

	rows, err := env.db.Query(sql)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		env.logger.Println(err)
		return
	}
	defer rows.Close()

	response := make(map[string][]Deferral) //externalPlanId -> []Deferral
	for rows.Next() {
		var externalPlanId string
		var sourceName string
		var sourceType string
		var planName string
		var deductAmount int
		var deductMethod int
		err := rows.Scan(&externalPlanId, &sourceName, &sourceType, &planName, &deductAmount, &deductMethod)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			env.logger.Println(err)
			return
		}
		plan := Plan{ExternalPlanId: externalPlanId, PlanName: planName}
		source := Source{SourceName: sourceName, SourceType: sourceType}
		deferral := Deferral{Source: source,
			Plan:         plan,
			DeductAmount: deductAmount,
			DeductMethod: deductMethod}

		if _, ok := response[externalPlanId]; !ok {
			response[externalPlanId] = make([]Deferral, 0, 10)
		}
		response[externalPlanId] = append(response[externalPlanId], deferral)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (env *Env) auth(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		env.logger.Println(err)
	}
	env.logger.Printf("Username is %s\n", authRequest.Username)

	//Don't do this :)
	sql := "select 1 from customers where name = '" + authRequest.Username + "' and password ='" + authRequest.Password + "'"
	rows, err := env.db.Query(sql)
	if err != nil {
		env.logger.Println(err)
	}
	defer rows.Close()

	authenticated := false

	for rows.Next() {
		authenticated = true
		break
	}

	if authenticated {
		w.Write([]byte("You have authenticated successfully!"))
	} else {
		w.Write([]byte("Invalid user/password"))
		env.logger.Printf("Invalid login attempt.  User: %s  Pass: %s", authRequest.Username, authRequest.Password)

	}

}

func (env *Env) getSources(w http.ResponseWriter, r *http.Request) {

	rows, err := env.db.Query("select * from sources")
	if err != nil {
		env.logger.Fatal(err)
	}
	defer rows.Close()
	sourceList := make([]Source, 0, 10)
	for rows.Next() {
		var sourcename string
		var sourcetype string
		err := rows.Scan(&sourcename, &sourcetype)
		if err != nil {
			env.logger.Fatal(err)
		}
		env.logger.Println(sourcename)
		s := Source{SourceName: sourcename, SourceType: sourcetype}
		sourceList = append(sourceList, s)

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sourceList)
}
