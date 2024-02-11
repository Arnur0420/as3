package main

import (
	"alexedwards.net/snippetbox/pkg/forms"
	"alexedwards.net/snippetbox/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Meal struct {
	ID       int
	MealName string
	Weekday  string
	Quantity int
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl.html", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "show.page.tmpl.html", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl.html", &templateData{Form: form})
		return
	}
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createMealForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createmeal.page.tmpl.html", &templateData{})
}

func (app *application) createMeal(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	mealForm := forms.NewMealForm(r.PostForm)
	mealForm.Validate()

	if !mealForm.Valid() {
		app.render(w, r, "createmeal.page.tmpl.html", &templateData{})
		return
	}

	quantity, err := strconv.Atoi(mealForm.Get("quantity"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	_, err = app.meals.Insert(mealForm.Get("meal_name"), mealForm.Get("weekday"), quantity)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/list", http.StatusSeeOther)
}
func (app *application) listMeals(w http.ResponseWriter, r *http.Request) {
	meals, err := app.meals.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Meals: meals,
	}

	app.render(w, r, "canteen.page.tmpl.html", data)
}
