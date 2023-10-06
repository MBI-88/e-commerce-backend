package settings

import (
	"time"

	"github.com/spf13/viper"
)

type Settings struct {
	DSN                       string
	Key                       string
	Port                      int
	Host                      string
	AllowOrigins              string
	AllowHeaders              string
	AllowMethods              string
	CsrfExpi                  time.Duration
	CookieSecure              bool
	CookieSameSite            string
	MaxAge                    int
	AllowCredentials          bool
	XSSprotection             string
	CrossOriginOpenerPolicy   string
	CrossOriginResourcePolicy string
}

// Set up production enviroments variables
func (s *Settings) setEnvPro() {
	s.DSN = viper.GetString("prod.DSN")
	s.Key = viper.GetString("prod.KEY")
	s.Port = viper.GetInt("prod.PORT")
	s.Host = viper.GetString("prod.HOST")
	s.AllowOrigins = viper.GetString("prod.ALLOW_ORIGINS")
	s.AllowHeaders = viper.GetString("prod.ALLOW_HEADERS")
	s.AllowMethods = viper.GetString("prod.ALLOW_METHODS")
	s.CookieSecure = viper.GetBool("prod.COOKIE_SAFE")
	s.CrossOriginOpenerPolicy = viper.GetString("prod.CROSS_ORIGIN_OPENER_POLICY")
	s.CrossOriginResourcePolicy = viper.GetString("prod.CROSS_ORIGIN_RESOURCE_POLICY")
	s.CsrfExpi = time.Duration(viper.GetInt("prod.CSRF_EXP")) * time.Minute
	s.MaxAge = viper.GetInt("prod.MAX_AGE")
	s.AllowCredentials = viper.GetBool("prod.ALLOW_CREDENTIALS")
	s.XSSprotection = viper.GetString("prod.XSS")
	s.CookieSameSite = viper.GetString("prod.COOKIE_SAME_SITE")

}

// Set up development enviroments variables
func (s *Settings) setEnvDev() {
	s.DSN = viper.GetString("dev.DSN")
	s.Key = viper.GetString("dev.KEY")
	s.Port = viper.GetInt("dev.PORT")
	s.Host = viper.GetString("dev.HOST")
	s.AllowOrigins = viper.GetString("dev.ALLOW_ORIGINS")
	s.AllowHeaders = viper.GetString("dev.ALLOW_HEADERS")
	s.CrossOriginOpenerPolicy = viper.GetString("dev.CROSS_ORIGIN_OPENER_POLICY")
	s.CrossOriginResourcePolicy = viper.GetString("dev.CROSS_ORIGIN_RESOURCE_POLICY")
	s.AllowMethods = viper.GetString("dev.ALLOW_METHODS")
	s.CookieSecure = viper.GetBool("dev.COOKIE_SAFE")
	s.CsrfExpi = time.Duration(viper.GetInt("dev.CSRF_EXP")) * time.Minute
	s.CookieSameSite = viper.GetString("dev.COOKIE_SAME_SITE")
	s.MaxAge = viper.GetInt("dev.MAX_AGE")
	s.AllowCredentials = viper.GetBool("dev.ALLOW_CREDENTIALS")
	s.XSSprotection = viper.GetString("dev.XSS")
}

// Get enviroment variables
func (s Settings) GetEnvVar(debug bool) Settings {
	viper.AddConfigPath("./")
	viper.SetConfigName("env")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("[-] File not found!")
		}
	}
	if debug {
		s.setEnvDev()
		return s
	}
	s.setEnvPro()
	return s
}

// Get enviroment variables for test case
func (s Settings) GetEnvVarTest() Settings {
	viper.AddConfigPath("./../")
	viper.SetConfigName("env")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("[-] File not found!")
		}
	}
	s.setEnvDev()
	return s

}
