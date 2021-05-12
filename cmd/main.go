package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/users-api/cmd/server"
	user "github.com/users-api/pkg/user"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	StartApp()
}

func StartApp() {
	initLog()
	readConfiguration()
	startWebServer()
}

func initLog() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func startWebServer() {

	s := server.New(&server.Config{
		Port: viper.GetInt("server.port"),
	})

	s.AddRoute("/health", func(w http.ResponseWriter, r *http.Request) {
		server.OK(w, r, map[string]interface{}{
			"name":   viper.GetString("server.name"),
			"status": "healthy",
		})
	}, http.MethodGet)

	user.RegisterRoutes(s)

	logrus.Info("starting http listener ...")
	go func() {
		s.ListenAndServe()
	}()

	// Wait for terminate signal to shut down server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	logrus.Infof("services started and listening in port %d ...", viper.GetInt("server.port"))
	<-c

}

func readConfiguration() {

	//Util for several environments
	env := flag.String("E", "dev", "Execution environment")
	flag.Parse()
	logrus.Infof("Starting user api in %s environment ...", *env)

	viper.AddConfigPath("./cmd/config")
	viper.SetConfigName("env_" + *env)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("error reading configuration from viper: %v", err)
	}
}
