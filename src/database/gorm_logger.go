package database

import (
	"database/sql/driver"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"time"
)

type GormLogger struct {
	logger    *logrus.Logger
	sqlRegexp *regexp.Regexp
}

func NewGormLogger(logger *logrus.Logger) *GormLogger {
	re := regexp.MustCompile(`(\$\d+)|\?`)
	return &GormLogger{logger: logger, sqlRegexp: re}
}

func (g *GormLogger) Print(v ...interface{}) {
	if len(v) == 0 || len(v) == 1 {
		g.logger.WithFields(logrus.Fields{
			"Module": "gorm",
		}).Error(fmt.Sprintf("[Gorm] Unknown message: %v", v))
		return
	}

	level := v[0]
	if level == "info" {
		info := v[1]
		g.logger.WithFields(logrus.Fields{
			"Module": "gorm",
			"Type":   level,
			"Info":   info,
		}).Info(fmt.Sprintf("[Gorm] info: %s", info))
		return
	}
	if level != "sql" {
		g.logger.WithFields(logrus.Fields{
			"Module": "gorm",
			"Type":   level,
		}).Info(fmt.Sprintf("[Gorm] unknown level %s: %v", level, v))
		return
	}

	source := v[1]
	duration := v[2]
	sql := g.render(v[3].(string), v[4])
	rows := v[5]
	g.logger.WithFields(logrus.Fields{
		"Module":   "gorm",
		"Type":     level,
		"Source":   source,
		"Duration": duration,
		"SQL":      sql,
		"Rows":     rows,
	}).Info(fmt.Sprintf("[Gorm] rows: %3d | %10s | %s", rows, duration, sql))
}

func (g *GormLogger) render(sql string, param interface{}) string {
	values := make([]interface{}, 0)
	for _, value := range param.([]interface{}) {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() { // valid
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok { // time
				values = append(values, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
			} else if b, ok := value.([]byte); ok { // bytes
				values = append(values, fmt.Sprintf("'%v'", string(b)))
			} else if r, ok := value.(driver.Valuer); ok { // driver
				if value, err := r.Value(); err == nil && value != nil {
					values = append(values, fmt.Sprintf("'%v'", value))
				} else {
					values = append(values, "NULL")
				}
			} else { // other value
				values = append(values, fmt.Sprintf("'%v'", value))
			}
		} else { // invalid
			values = append(values, fmt.Sprintf("'%v'", value))
		}
	}

	return fmt.Sprintf(g.sqlRegexp.ReplaceAllString(sql, "%v"), values...)
}
