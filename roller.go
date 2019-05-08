package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func DiceRoller(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	if dir != "/roll/" {
		e(w, http.StatusNotFound)
		return
	}
	split := strings.Split(file, "d")
	if len(split) != 2 {
		e(w, http.StatusInternalServerError)
		return
	}
	number, err := strconv.Atoi(split[0])
	if err != nil {
		e(w, http.StatusInternalServerError)
		return
	}
	sides, err := strconv.Atoi(split[1])
	if err != nil {
		e(w, http.StatusInternalServerError)
		return
	}
	rolls := make([]string, number)
	var total int

	for i :=0 ; i< number; i++ {
		roll := int(rand.Int31n(int32(sides))) + 1
		rolls[i] = strconv.Itoa(roll)
		total += roll
	}
	result := fmt.Sprintf("%d %s", total, strings.Join(rolls, " "))
	_, _ = w.Write([]byte(result))
}

func e(w http.ResponseWriter, i int) {
	w.WriteHeader(i)
	_, _ = w.Write([]byte("/roll/#d#"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/roll/", DiceRoller)


	switch {
	case os.Getenv("AWS_EXECUTION_ENV") != "":
		lambda.Start(httpadapter.New(mux).ProxyWithContext)
	default:
		fmt.Println("http://localhost:3000/roll/2d6")
		log.Fatal(http.ListenAndServe(":3000", mux))
	}
}

