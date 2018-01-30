package dbconn

import (
	"github.com/bububa/goconfig/config"
	"github.com/bububa/mymysql/autorc"
	_ "github.com/bububa/mymysql/thrsafe"
	"sync"
)

const (
	_CONFIG_FILE = "/var/code/go/config.cfg"
)

var (
	Mdb *autorc.Conn
)

func init() {
	var once sync.Once
	once.Do(func() {
		cfg, _ := config.ReadDefault(_CONFIG_FILE)

		hostMaster, _ := cfg.String("masterdb", "host")
		userMaster, _ := cfg.String("masterdb", "user")
		passwdMaster, _ := cfg.String("masterdb", "passwd")
		dbnameMaster, _ := cfg.String("masterdb", "dbname")

		Mdb = autorc.New("tcp", "", hostMaster, userMaster, passwdMaster, dbnameMaster)
		Mdb.Register("set names utf8")
		Mdb.Query("select 1 from dual") // test connection
	})
}
