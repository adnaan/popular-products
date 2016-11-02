package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
)

/*Server*/
type server struct {
	httpServer *http.Server
	listener   net.Listener
}

func (s *server) listenAndServe() (chan struct{}, error) {
	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return done, err
	}
	s.listener = listener
	go s.httpServer.Serve(s.listener)
	log.Printf("Server now listening at localhost%v", s.httpServer.Addr)

	go func() {
		<-sigs
		s.shutdown()
		done <- struct{}{}
	}()

	return done, nil

}

func (s *server) shutdown() error {

	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}
	log.Println("Shutting down server")
	return nil

}

// Product ..
type Product struct {
	gorm.Model
	Vote  int    `json:"vote"`
	Title string `json:"title"`
}

// Response ...
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// newServer bootstraps a server and inits data
func newServer(port string, count int, db *gorm.DB) *server {

	// init db
	if !db.HasTable(&Product{}) {
		log.Println("Initialising DB")
		// Migrate the schema
		db.AutoMigrate(&Product{})

		for i := 0; i < count; i++ {
			db.Create(&Product{Vote: i, Title: fmt.Sprintf("Product-%d", i)})
		}
	}

	queryHandler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var products []Product
		db.Find(&products)

		response := Response{Success: true, Data: products}
		b, err := json.Marshal(response)
		if err != nil {
			log.Printf("GET /api/products json marshal error %d,  %v", http.StatusInternalServerError, err)
			http.Error(w, "json marshal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(b))
		log.Println("GET /api/products 200")
	}

	voteHandler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var product Product
		id, err := strconv.Atoi(ps.ByName("id"))
		if err != nil {
			log.Printf("GET /api/products/vote/%d Malformed ID %d,  %v", id, http.StatusBadRequest, err)
			http.Error(w, "Malformed ID", http.StatusBadRequest)
			return
		}

		db.First(&product, id)
		db.Model(&product).UpdateColumn("vote", gorm.Expr("vote + ?", 1))
		db.First(&product, id)

		response := Response{Success: true, Data: product.Vote}
		b, err := json.Marshal(response)
		if err != nil {
			log.Printf("GET /api/products/vote/%d json marshal error %d,  %v", id, http.StatusInternalServerError, err)
			http.Error(w, "json marshal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(b))
		log.Printf("GET /api/products/vote/%d 200", id)
	}

	router := httprouter.New()
	router.GET("/api/products", queryHandler)
	router.GET("/api/products/vote/:id", voteHandler)

	httpServer := &http.Server{Addr: ":" + port, Handler: router}
	return &server{httpServer: httpServer}
}

func main() {

	log.Println("Connecting to sqlite DB")
	db, err := gorm.Open("sqlite3", "products.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	server := newServer("3333", 3, db)
	done, err := server.listenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
	<-done
}
