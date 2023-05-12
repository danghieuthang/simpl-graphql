// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	_ "github.com/lib/pq"

// )

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "1"
// 	dbname   = "golang"
// )

// // album represents data about a record album.
// type album struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

// // albums slice to seed record album data.
// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

// // getAlbums responds with the list of all albums as JSON.
// func getAlbums(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, albums)
// }

// func getAlbumById(c *gin.Context) {
// 	id := c.Param("id")

// 	for _, v := range albums {
// 		if v.ID == id {
// 			c.IndentedJSON(http.StatusOK, v)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
// }

// func createAlbum(c *gin.Context) {
// 	var newAlbum album

// 	if err := c.BindJSON(&newAlbum); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Request invalid"})
// 		return
// 	}

// 	albums = append(albums, newAlbum)
// 	c.IndentedJSON(http.StatusOK, newAlbum)
// }

// func deleteAlbum(c *gin.Context) {
// 	id := c.Param("id")
// 	for i, v := range albums {
// 		if v.ID == id {
// 			albums = append(albums[:i], albums[i+1:]...)
// 			c.IndentedJSON(http.StatusOK, v)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Id not exist"})
// }

// func LoggingMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		log.Println("Middeware", c.Request.Method, c.Request.URL.Path)
// 		c.Next()
// 	}
// }
// func main() {

// 	//connection string
// 	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
// 	//open database
// 	db, err := sql.Open("postgres", psqlConn)

// 	CheckError(err)

// 	err = db.Ping()

// 	CheckError(err)

// 	insertStm := `insert into "Students"("Name", "Roll") values('John', 1)`
// 	_, e := db.Exec(insertStm)
// 	CheckError(e)

// 	router := gin.Default()

// 	router.Use(LoggingMiddleware())

// 	router.GET("albums", getAlbums)
// 	router.GET("albums/:id", getAlbumById)
// 	router.DELETE("albums/:id", deleteAlbum)
// 	router.POST("albums", createAlbum)

// 	router.Run("localhost:8080")

// }

// func CheckError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func DbModel()  {
// 	db:=pg.Con
// }

package main

import (
	"context"
	"example/web-service-gin/controller"
	"example/web-service-gin/database"
	"example/web-service-gin/gql"
	"example/web-service-gin/initializers"
	"example/web-service-gin/pkg/jwt"
	"example/web-service-gin/pkg/log"
	"example/web-service-gin/repository"
	"fmt"
	"net/http"
	"os"
	"time"
)

func init() {
	initializers.LoadEnvVariables()
	log.InitializeLogger()
	// log.Logger.Info("Initialize database...")
	database.InitializeDatabase()
	// log.Logger.Info("Sync database...")
	database.SyncDatabase()
}

func main() {
	// log.Logger.Info("Initilize graphql endpoind...")
	repos := repository.InitRepositories(database.DB)
	ServiceFactory := controller.InitServices(repos)
	schema := gql.Schema(ServiceFactory)

	http.Handle("/graphql", PerformanceMiddleware(AuthenticationMiddleware(gql.GraphqlHandlfunc(schema))))

	port := os.Getenv("PORT")
	fmt.Println("Server is started at: http://localhost:/" + port + "/")
	fmt.Println("Graphql api server is started at: http://localhost:" + port + "/graphql")
	http.ListenAndServe(":"+port, nil)
}

func PerformanceMiddleware(next http.Handler) http.Handler {
	perfomanceFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI

		next.ServeHTTP(rw, r) // serve the original request

		duration := time.Since(start)
		// log request details
		if duration.Seconds() > 10 {
			log.Logger.Info(fmt.Sprintf("%s: Log Duration %s\n", uri, duration))
		}
	}
	return http.HandlerFunc(perfomanceFn)
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	authenticationFn := func(rw http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		// Allow unauthenticated users in
		if header == "" {
			next.ServeHTTP(rw, r)
			return
		}

		//validate jwt token
		tokenStr := header
		user, err := jwt.ParseToken(tokenStr)
		if err != nil {
			http.Error(rw, "Invalid token", http.StatusForbidden)
			return
		}
		// put it in context
		ctx := context.WithValue(r.Context(), "currentUser", user)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(authenticationFn)
}
