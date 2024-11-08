package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Product struct {
	ID            int
	Name          string
	Description   string
	Price         float64
	Image         string
	Sizes         []string
	Discounted    bool
	OriginalPrice float64
}

var products = []Product{
	{ID: 1, Name: "Pull Capuche Vert", Price: 140, Image: "/assets/img/products/19A.webp", Sizes: []string{"S", "M", "L"}, Description: "Un pull à capuche confortable.", Discounted: false},
	{ID: 2, Name: "Pull Capuche Marine", Price: 138, Image: "/assets/img/products/21A.webp", Sizes: []string{"S", "M", "L"}, Description: "Pull pour représenter la cité londonienne.", Discounted: true, OriginalPrice: 150},
	{ID: 3, Name: "Pull Crew Noir", Price: 128, Image: "/assets/img/products/22A.webp", Sizes: []string{"M", "L"}, Description: "Un pull crew passe-partout.", Discounted: true, OriginalPrice: 130},
	{ID: 4, Name: "Pull Capuche Jaune", Price: 168, Image: "/assets/img/products/16A.webp", Sizes: []string{"S", "L"}, Description: "Pull mojito jaune pastel.", Discounted: false},
	{ID: 5, Name: "Jean Stone", Price: 125, Image: "/assets/img/products/34B.webp", Sizes: []string{"32", "34", "36"}, Description: "Jean décontracté et élégant.", Discounted: false},
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func listProducts(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func showProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	for _, product := range products {
		if product.ID == id {
			err := templates.ExecuteTemplate(w, "product.html", product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	http.NotFound(w, r)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        err := templates.ExecuteTemplate(w, "add_product.html", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    if r.Method == http.MethodPost {
        name := r.FormValue("name")
        price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
        description := r.FormValue("description")
        originalPrice, _ := strconv.ParseFloat(r.FormValue("originalPrice"), 64)
        
        discounted := originalPrice > price

        newProduct := Product{
            ID:           len(products) + 1,
            Name:         name,
            Price:        price,
            Image:        "/assets/img/products/18A.webp",
            Description:  description,
            Sizes:        []string{"M"},
            Discounted:   discounted,
            OriginalPrice: originalPrice,
        }
        products = append(products, newProduct)
        http.Redirect(w, r, "/product?id="+strconv.Itoa(newProduct.ID), http.StatusSeeOther)
    }
}


func main() {
	http.HandleFunc("/", listProducts)
	http.HandleFunc("/product", showProduct)
	http.HandleFunc("/add", addProduct)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Serveur lancé http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
