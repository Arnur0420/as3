package forms

import (
	"fmt"

	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

type MealForm struct {
	Form
}

func NewMealForm(data url.Values) *MealForm {
	return &MealForm{
		Form: Form{
			Values: data,
			Errors: errors(map[string][]string{}),
		},
	}
}

func (f *MealForm) Validate() {
	f.Required("meal_name", "weekday", "quantity")
	f.MaxLength("meal_name", 100)
}

func (f *MealForm) Valid() bool {
	return len(f.Errors) == 0
}
