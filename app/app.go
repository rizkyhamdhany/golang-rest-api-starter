package app

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"github.com/rizkyhamdhany/kelase-micro/config"
	"fmt"
	"github.com/rizkyhamdhany/kelase-micro/app/model"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"
	"github.com/rizkyhamdhany/kelase-micro/app/cc"
	"github.com/rizkyhamdhany/kelase-micro/app/handler"
	"github.com/rizkyhamdhany/kelase-micro/app/module/users"
	"github.com/rizkyhamdhany/kelase-micro/app/module/admin"
)

type App struct {
	Router *mux.Router
	AuthRouter  *mux.Router
	DB     *gorm.DB
}

func (a *App) Run(host string) {
	handler := cors.Default().Handler(a.Router)
	log.Fatal(http.ListenAndServe(host, handler))
}

func (a *App) InitializeDB(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@%s/%s?charset=%s&parseTime=True",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBName,
		config.DBCharset)

	db, err := gorm.Open(config.DBDriver, dbURI)
	if err != nil {
		log.Fatal(err)
	}
	a.DB = model.DBMigrate(db)
}

func (a *App) InitializeRouter() {
	a.Router = mux.NewRouter()
	a.AuthRouter = mux.NewRouter()
	a.setRouters()
	a.setAuthRouters()
}

func (a *App) setRouters() {
	a.Get("/version", getVersion)
	a.Post("/login", users.Login(a.DB))
	a.Post("/signup", users.Signup(a.DB))
	a.Post("/admin/login", admin.Login(a.DB))
}

func (a *App) setAuthRouters() {
	a.GetWithAuth("/me", users.GetMe(a.DB))
	//jwt middleware
	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(cc.JwtSecret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	an := negroni.New(negroni.HandlerFunc(mw.HandlerWithNext), negroni.Wrap(a.AuthRouter))
	a.Router.PathPrefix("/").Handler(an)

	n := negroni.Classic()
	n.UseHandler(a.AuthRouter)
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.AuthRouter.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) GetWithAuth(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.AuthRouter.HandleFunc(path, f).Methods("GET")
}

func (a *App) PostWithAuth(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.AuthRouter.HandleFunc(path, f).Methods("POST")
}

func (a *App) DeleteWithAuth(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.AuthRouter.HandleFunc(path, f).Methods("DELETE")
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	config := config.Config{}
	config.Read()
	res := make(map[string] string)
	res["version"] = config.Version
	handler.RespondWithJson(w, http.StatusOK, res)
}