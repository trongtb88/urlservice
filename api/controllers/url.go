package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/trongtb88/urlservice/api/constants"
	"github.com/trongtb88/urlservice/api/entity"
	"github.com/trongtb88/urlservice/api/middlewares"
	"github.com/trongtb88/urlservice/api/utils"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server * Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	middlewares.JsonResponse(w, http.StatusOK, "Welcome to rate system")
}

func (server * Server) Shorten(w http.ResponseWriter, r *http.Request) {

	var param entity.HttpRequestUrl
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		middlewares.ErrorResponse(w, http.StatusBadRequest, err, "Can not read body request")
		return
	}
	if err := json.Unmarshal(body, &param); err != nil {
		middlewares.ErrorResponse(w, http.StatusBadRequest, err, "Can not parse body request")
		return
	}

	if len(param.Url) == 0 {
		middlewares.ErrorResponse(w, http.StatusBadRequest, err, "url is not present")
	}

	if len(param.Shortcode) > 0 && !utils.IsMatchRegex(constants.SHORT_CODE_REGEX, param.Shortcode) {
		middlewares.ErrorResponse(w, http.StatusUnprocessableEntity, nil, "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{6}$.")
		return
	}

	// TODO Check exist short code or not, we can use Redis : SetNX Operator to check exist key to avoid idempotent problem
	if len(param.Shortcode) > 0 {
		// check exist in mysql
		uniqueShortCode, err := server.checkUniqueShortCode(param.Shortcode)
		if err != nil {
			middlewares.ErrorResponse(w, http.StatusInternalServerError, err, "Error on check short code unique")
			return
		}
		if !uniqueShortCode {
			middlewares.ErrorResponse(w, http.StatusConflict, nil, "The the desired shortcode is already in use. Shortcodes are case-sensitive")
			return
		}
	} else {
		param.Shortcode = utils.GenerateRandom("^[0-9a-zA-Z_]{6}$")
		// We may retry to get until we got unique shortCode here,
		// or We can improve by using the other micro service to generate in advance shortCode for us, and we only need to
		// retry from it, at this tests I only try generate 2 time
		//
		isUnique, err := server.checkUniqueShortCode(param.Shortcode)
		if err != nil {
			middlewares.ErrorResponse(w, http.StatusInternalServerError, err, "Error on check short code unique")
			return
		}
		if !isUnique {
			param.Shortcode = utils.GenerateRandom("^[0-9a-zA-Z_]{6}$")
		}
	}
	// Store in DB
	url := entity.Url{
		OriginUrl: param.Url,
		ShortCode: param.Shortcode,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		RedirectCount: 0,
	}
	err = server.storeDB(url)
	if err != nil {
		middlewares.ErrorResponse(w, http.StatusInternalServerError, err, "Can not store new short code")
		return
	}

	res := entity.HttpResponseUrl{
		Shortcode: param.Shortcode,
	}
	middlewares.JsonResponse(w, http.StatusCreated, res)
}

func (server *Server) RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortcode, ok := vars["shortcode"]

	if !ok {
		middlewares.ErrorResponse(w, http.StatusBadRequest, nil, "Not found parameter")
		return
	}

	isUnique, err := server.checkUniqueShortCode(shortcode)
	if err != nil {
		middlewares.ErrorResponse(w, http.StatusInternalServerError, err, "Can not check exist shortcode")
		return
	}
	if isUnique {
		middlewares.ErrorResponse(w, http.StatusNotFound, err, "The shortcode cannot be found in the system")
		return
	}

	// increase redirect count
	err = server.DB.Exec("UPDATE urls SET redirect_count = redirect_count + 1, last_seen_at = ?  WHERE short_code = ?",
		time.Now().UTC(), shortcode).Error
	if err != nil {
		middlewares.ErrorResponse(w, http.StatusInternalServerError, err, "Error on update redirect count ")
		return
	}

	// get original
	var orignalUrl string
	err = server.DB.Raw("SELECT origin_url FROM urls  WHERE short_code = ?", shortcode).Scan(&orignalUrl).Error
	if err != nil {
		middlewares.ErrorResponse(w, http.StatusInternalServerError, err, "Error on update redirect count ")
		return
	}
	middlewares.JsonRedirectResponse(w, orignalUrl)

}

func (server *Server) Stats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortcode, ok := vars["shortcode"]

	if !ok {
		middlewares.ErrorResponse(w, http.StatusBadRequest, nil, "Not found parameter")
		return
	}

	isUnique, err := server.checkUniqueShortCode(shortcode)
	if err != nil {
		middlewares.ErrorResponse(w, http.StatusInternalServerError, err, "Can not check exist shortcode")
		return
	}
	if isUnique {
		middlewares.ErrorResponse(w, http.StatusNotFound, err, "The shortcode cannot be found in the system")
		return
	}


	// get original
	var url entity.Url
	err = server.DB.Raw("SELECT redirect_count, created_at, last_seen_at FROM urls  WHERE short_code = ?", shortcode).Scan(&url).Error
	if err != nil {
		middlewares.ErrorResponse(w, http.StatusInternalServerError, err, "Error on update redirect count ")
		return
	}

	httpResponseStats := entity.HttpResponseStatsUrl{
		RedirectCount: url.RedirectCount,
		StartDate: url.CreatedAt.Format(constants.FORMAT_DATE),
	}

	if url.RedirectCount > 0 {
		httpResponseStats.LastSeenDate = url.LastSeenAt.Format(constants.FORMAT_DATE)
	}

	middlewares.JsonResponse(w, http.StatusOK, httpResponseStats)
}


func (server *Server) checkUniqueShortCode(shortCode string) (bool, error) {
	var foundShortCode string
	err := server.DB.Raw("SELECT short_code FROM urls WHERE short_code = ?", shortCode).Scan(&foundShortCode).Error
	if err != nil {
		log.Printf("Error on check short code unique %v", err)
		return false, err
	}
	if len(foundShortCode) > 0 {
		return false, nil
	}
	return true, nil
}

func (server *Server) storeDB(url entity.Url) error {
	log.Println(url)
	err := server.DB.Exec("INSERT INTO urls (origin_url, short_code, redirect_count, created_at , updated_at) " +
		"VALUES (?, ? , ? , ? , ?)", url.OriginUrl,
		url.ShortCode, 0, url.CreatedAt, url.UpdatedAt).Error
	return err
}
