package errs

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
)

var errorSet = map[string]struct{}{}

type AppError interface {
	Code() string
	Error() string
	Throw(applog.Logger) AppError
	Throwf(applog.Logger, string, ...interface{}) AppError
	Annotatef(applog.Logger, string, ...interface{}) AppError
	Fatalf(applog.Logger, string, ...interface{})
	Is(target AppError) bool
}

type appError struct {
	msg         string
	code        string
	annotations []string
}

func _new(code string, format string, args ...interface{}) AppError {
	if ensureUnique(code) {
		panic(fmt.Sprintf("AppError Code %+v already exists", code))
	}

	if !ensureFormat(code) {
		panic(fmt.Sprintf("AppError Code must have format ABC123 but is %+v", code))
	}

	errorSet[code] = struct{}{}

	msg := fmt.Sprintf(format, args...)
	te := &appError{
		msg:         msg,
		code:        code,
		annotations: make([]string, 0),
	}

	return te
}

func ensureUnique(id string) bool {
	_, exists := errorSet[id]

	return exists
}

func ensureFormat(id string) bool {
	r, err := regexp.Compile("^[A-Z,0-9]{3}[0-9]{3}$")
	if err != nil {
		return false
	}
	return r.MatchString(id)
}

func (te *appError) Code() string {
	return te.code
}

func (te *appError) Error() string {
	var msg string
	if len(te.annotations) > 0 {
		annotations := strings.Join(te.annotations, ", ")
		msg = fmt.Sprintf("%s: %s - %s", te.code, te.msg, annotations)
	} else {
		msg = fmt.Sprintf("%s: %s", te.code, te.msg)
	}

	return msg
}

func (te *appError) Fatalf(log applog.Logger, format string, args ...interface{}) {
	message := te.formatMessage(format, args...)
	if log == nil {
		log = applog.New(os.Stdout)
	}

	for _, s := range StackTrace() {
		log.Stackf(s)
	}

	log.Fatalf(message)
}

func (te *appError) Throw(log applog.Logger) AppError {
	newError := &appError{
		msg:         te.msg,
		code:        te.code,
		annotations: te.annotations,
	}

	if log != nil {
		log.Errorf(newError.Error())
		for _, s := range StackTrace() {
			log.Stackf(s)
		}
	}

	return newError
}

func (te *appError) Throwf(log applog.Logger, format string, args ...interface{}) AppError {
	newError := &appError{
		msg:         te.formatMessage(format, args...),
		code:        te.code,
		annotations: te.annotations,
	}

	if log != nil {
		log.Errorf(newError.Error())
		for _, s := range StackTrace() {
			log.Stackf(s)
		}
	}

	return newError
}

func (te *appError) formatMessage(format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	if message != "" {
		message = fmt.Sprintf("%s, %s", te.msg, message)
	} else {
		message = te.msg
	}

	return message
}

func (te *appError) Annotatef(log applog.Logger, format string, args ...interface{}) AppError {
	message := fmt.Sprintf(format, args...)

	if message == "" {
		return te
	}

	if log != nil {
		log.Errorf(message)
	}

	te.annotations = append(te.annotations, message)
	return te
}

func (te *appError) Is(err AppError) bool {
	return te.Code() == err.Code()
}

func StackTrace() []string {
	var lines []string
	ok := true
	var file string
	var line int
	stack := 3 // 0 here, 1 util.Context, 2 tag_error, 3 originating call

	lines = append(lines, "[STACK TRACE]")

	for ok {
		_, file, line, ok = runtime.Caller(stack)
		if ok {
			lines = append(lines, "["+strconv.Itoa(stack)+"] "+file+":"+strconv.Itoa(line))
		} else {
			lines = append(lines, "")
		}
		stack++
	}
	// remove the last new line
	lines = lines[:len(lines)-1]

	return lines
}
