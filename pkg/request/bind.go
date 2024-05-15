package request

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/mr55p-dev/pagemail/internal/logging"
)

var logger = logging.NewLogger("bind")

var ErrUnsupportedType = errors.New("unsupported type")
var ErrNonPointerArg = errors.New("non-pointer argument")
var ErrNonStructArg = errors.New("non-struct argument")

func BindRequest[T any](w http.ResponseWriter, r *http.Request) *T {
	out := new(T)
	err := Bind(out, r)
	if err != nil {
		http.Error(w, "Failed to bind request", http.StatusBadRequest)
		return nil
	}
	return out
}

func Bind(v any, r *http.Request) (err error) {
	log := logging.NewLogger("bind")
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		var ok bool
		err, ok = r.(error)
		if !ok {
			err = fmt.Errorf("Bind paniced with non-error: %v", r)
			return
		}
		if errors.Is(err, ErrUnsupportedType) {
			logger.WithError(err).Error("Bind tried to attach to unsupported type")
			panic(err)
		} else if errors.Is(err, ErrNonPointerArg) {
			logger.WithError(err).Error("Bind tried to bind to non-pointer argument")
			panic(err)
		} else if errors.Is(err, ErrNonStructArg) {
			logger.WithError(err).Error("Bind tried to bind to non-struct argument")
			panic(err)
		}
		err = fmt.Errorf("Bind paniced with error: %w", err)
	}()

	// check if v is a struct pointer
	valueOf := getValue(v)
	typeOf := valueOf.Type()

	err = r.ParseForm()
	if err != nil {
		return fmt.Errorf("Bind failed to parse form: %w", err)
	}

	log.Debug("Binding struct", "struct", typeOf.Name())
	for i := 0; i < typeOf.NumField(); i++ {
		fieldType := typeOf.Field(i)
		fieldValue := valueOf.Field(i)
		fieldTag := fieldType.Tag
		var rawVal string
		if tag := fieldTag.Get("form"); tag != "" {
			rawVal = handeForm(tag, r)
		} else if tag := fieldTag.Get("query"); tag != "" {
			rawVal = handleQuery(tag, r)
		}

		switch fieldType.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(rawVal)
		case reflect.Int:
			parsed, err := strconv.Atoi(rawVal)
			if err != nil {
				return fmt.Errorf("Bind failed to parse int for field %s: %w", fieldTag, err)
			}
			fieldValue.SetInt(int64(parsed))
		case reflect.Bool:
			fieldValue.SetBool(rawVal == "true")
		default:
			return fmt.Errorf("Bind failed to bind field %s (%s): %w", fieldTag, fieldType.Type.Kind(), ErrUnsupportedType)
		}
	}
	return nil
}

func getValue(v any) reflect.Value {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Ptr {
		panic(ErrNonPointerArg)
	}
	valueOf = valueOf.Elem()
	if valueOf.Kind() != reflect.Struct {
		panic(ErrNonStructArg)
	}
	return valueOf
}

func handleQuery(tag string, r *http.Request) string {
	return r.URL.Query().Get(tag)
}

func handeForm(tag string, r *http.Request) string {
	return r.PostForm.Get(tag)
}
