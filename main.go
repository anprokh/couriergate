package main

import (
	Cfg "couriergate/configs"
	Models "couriergate/models"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type program struct {
	cfg *Models.Config
}

var p *program

func init() {
	conf, err := loadConfig()
	if err != nil {
		fmt.Println("Failed to read config file:", err.Error())
		os.Exit(1)
	}
	p = &program{cfg: conf}
	// определим каталог исполняемого файла
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	Cfg.ExPath = filepath.Dir(ex)
}

func loadConfig() (*Models.Config, error) {
	conf := new(Models.Config)

	c, err := ini.LoadSources(ini.LoadOptions{
		SpaceBeforeInlineComment: true,
	}, "courier.ini")
	if err != nil {
		return nil, err
	}
	err = c.MapTo(conf)
	if err != nil {
		return nil, err
	}
	// проверим указание конфигурационных параметров
	s := reflect.ValueOf(conf).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		//fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
		if len(f.String()) == 0 {
			return nil, errors.New(fmt.Sprintf("parameter '%s' is empty", typeOfT.Field(i).Name))
		}
	}

	return conf, nil
}

func main() {

	if err := p.ConnectDB(); err != nil {
		log.Println("Database connection failed:", err.Error())
		os.Exit(1)
	}

	anotherRepeat := true
	for anotherRepeat {

		if err := p.GetAuthToken(); err != nil {
			fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
			continue
		}
		fmt.Fprintf(color.Output, "%s %s completed successfully at %s\n", color.GreenString("[info]"), color.CyanString("GetAuthToken"), time.Now().Format("2006-01-02 15:04:05"))

		p.setEnv()
		p.exportSignedDocuments()
		p.exportSignedTicketReply()
		p.processingDocumentEvents()

		// обрабатываем входящие УПД если задано в настройках фирмы
		if p.cfg.IncomingApply {
			p.processingIncomingDocuments()
		}

		p.moveCompletedDocuments()

		fmt.Fprintf(color.Output, "%s %s completed successfully at %s\n", color.GreenString("[info]"), color.CyanString("All operations"), time.Now().Format("2006-01-02 15:04:05"))
		//anotherRepeat = false

		<-time.After(time.Second * 5)
	}

}
