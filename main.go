package main

import (
    _ "fmt"
    "github.com/gin-gonic/gin"
    "math/rand"
    "net/http"
    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql"
    "time"
    "strconv"
    //"github.com/gin-contrib/cors"
)

var db *sql.DB

// Object - News
type News struct {
    Id          int       `json:"id"`
    Title       string    `json:"title" form:"title"`
    Content     string    `json:"content" form:"content"`
    DateCreated time.Time `json:"date_created" form:"date_created"`
    LastUpdated time.Time `json:"last_updated" form:"last_updated"`
    OwnerId     int       `json:"owner_id" form:"owner_id"`
}

func (n News) get() (news News, err error) {
    row := db.QueryRow("SELECT id, title, content, date_created, last_updated, owner_id FROM news WHERE id=?", n.Id)
    err = row.Scan(&news.Id, &news.Title, &news.Content, &news.DateCreated, &news.LastUpdated, &news.OwnerId)
    if err != nil {
        return
    }
    return
}

func (n News) getAll() (newsList []News, err error) {
    rows, err := db.Query("SELECT id, title, content, date_created, last_updated, owner_id FROM news")
    if err != nil {
        return
    }
    for rows.Next() {
        var news News
        rows.Scan(&news.Id, &news.Title, &news.Content, &news.DateCreated, &news.LastUpdated, &news.OwnerId)
        newsList = append(newsList, news)
    }
    defer rows.Close()
    return
}

func main() {
    var err error
    // Database
    db, err = sql.Open("mysql", "root:google13794628@tcp(127.0.0.1:3306)/daybook?parseTime=true")
    if err != nil {
        log.Fatalln(err)
    }
    defer db.Close()
     
    db.SetMaxIdleConns(20)
    db.SetMaxOpenConns(20)

    if err := db.Ping(); err != nil{
        log.Fatalln(err)
    }

    // Routes
    router := gin.Default()
    // CORS必須在註冊路徑之前設定，否則無效
    //router.Use(cors.Default()) 

    v1 := router.Group("/v1")
    {
        v1.GET("/hello", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "welcome to bgops,please visit https://xxbandy.github.io!",
            })
        })
        v1.GET("/hello/:name", func(c *gin.Context) {
            name := c.Param("name")
            c.String(http.StatusOK, "Hello %s", name)
        })

        v1.GET("/line", func(c *gin.Context) {
            c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
            legendData := []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
            xAxisData := []int{120, 240, rand.Intn(500), rand.Intn(500), 150, 230, 180}
            c.JSON(200, gin.H{
                "legend_data": legendData,
                "xAxis_data":  xAxisData,
            })
        })

        // Routes - News
        v1.GET("/news", func(c *gin.Context) {
            n := News{}
            newsList, err := n.getAll()
            if err != nil {
                log.Fatalln(err)
            }
            c.JSON(http.StatusOK, gin.H{
                "result": newsList,
                "count":  len(newsList),
            })
        })

        v1.GET("/news/:id", func(c *gin.Context) {
            var result gin.H
            id := c.Param("id")

            Id, err := strconv.Atoi(id)
            if err != nil {
                log.Fatalln(err)
            }

            n := News{
                Id: Id,
            }
            news, err := n.get()
            if err != nil {
                result = gin.H{
                    "result": nil,
                    "count":  0,
                }
            } else {
                result = gin.H{
                    "result": news,
                    "count":  1,
                }

            }
            c.JSON(http.StatusOK, result)
        })
    }

    // Default Route
    router.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{
            "status": 404,
            "error":  "404, page not exists!",
        })
    })

    // Port
    router.Run(":8000")
}
