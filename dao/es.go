package dao

import (
	"LarsWebV0/config"
	"fmt"
	"github.com/olivere/elastic/v7"
	logger "github.com/sirupsen/logrus"
)

var EsClient *elastic.Client

func EsSetup() {
	var err error
	EsClient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(fmt.Sprintf("http://%v:%v", config.Hostname, config.EsPort)))
	/*cfg := elasticsearch.Config{
		CloudID: "Laurus_D:YXNpYS1lYXN0MS5nY3AuZWxhc3RpYy1jbG91ZC5jb206NDQzJDI4ZjRiYmNlOGM4ZjQ3YWY4ZGZiODVkYTRhMTNjN2ZiJDI3MTI5MmVkOGQyMzRjN2ZhNmUxZGRiYzVkYjkxNTFm",
		APIKey:  "LarsWebv0",
	}
	es, err := elasticsearch.NewClient(cfg)*/

	//es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Errorf("es set up err: %v", err)
		return
	}

}
