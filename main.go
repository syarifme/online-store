package main

import (
	"encoding/json"
	"net/http"
)

var StockAvailable int = 0

type testOrders []int

var tO testOrders

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/checkstock", CheckStock)
	http.HandleFunc("/user", GetUser)
	http.HandleFunc("/product", GetProduct)
	http.HandleFunc("/product/update", UpdateStock)
	http.HandleFunc("/order", MakeOrder)
	http.HandleFunc("/order/test", MakeMultipleOrder)

	http.ListenAndServe(":3030", nil)
}

// func (o *testOrders) String() string {
// 	return ""
// }

// func (o *testOrders) Set(value string) error {
// 	val, err := strconv.Atoi(value)
// 	if err == nil {
// 		*o = append(*o, val)
// 	}

// 	fmt.Print(*o)
// 	return nil
// }

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		NotAllowedResponse(w)

		return
	}
	SuccessResponse(w, "Welcome to order application")
}

func CheckStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		NotAllowedResponse(w)

		return
	}
	SuccessResponse(w, StockAvailable)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		NotAllowedResponse(w)

		return
	}

	var stock StockUpdate
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		BadRequestResponse(w)
		return
	}

	StockAvailable = stock.Stock
	SuccessResponse(w, StockAvailable)
}

func MakeOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		NotAllowedResponse(w)

		return
	}

	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		BadRequestResponse(w)
		return
	}

	order.Id = 1
	if StockAvailable == 0 || order.Quantity > StockAvailable {
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "item not available for the amount specified",
			Data:    nil,
		})

		return
	}
	if StockAvailable > 0 {
		StockAvailable -= order.Quantity
	}
	SuccessResponse(w, order)
}

func MakeMultipleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		NotAllowedResponse(w)

		return
	}

	var orders []Order
	err := json.NewDecoder(r.Body).Decode(&orders)

	if err != nil {
		BadRequestResponse(w)
		return
	}

	var responses []Response
	for i := 0; i < len(orders); i++ {
		orders[i].Id = 1
		if StockAvailable == 0 || orders[i].Quantity > StockAvailable {
			resp := Response{
				Code:    400,
				Message: "item not available for the amount specified",
				Data:    nil,
			}

			responses = append(responses, resp)
		} else {
			sResp := Response{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    orders[i],
			}

			responses = append(responses, sResp)
			if StockAvailable > 0 {
				StockAvailable -= orders[i].Quantity
			}
		}
	}

	json.NewEncoder(w).Encode(responses)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		NotAllowedResponse(w)

		return
	}

	user := User{
		Id:    1,
		Name:  "Eko",
		Email: "eko@mail.com",
	}

	SuccessResponse(w, user)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		NotAllowedResponse(w)

		return
	}

	product := Product{
		Id:            1,
		Name:          "Kancing baju",
		StockQuantity: 12,
		Price:         125000,
	}

	SuccessResponse(w, product)
}

func SuccessResponse(w http.ResponseWriter, item interface{}) {
	data := Response{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    item,
	}

	json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter) {
	data := Response{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Data:    nil,
	}

	json.NewEncoder(w).Encode(data)
}

func BadRequestResponse(w http.ResponseWriter) {
	data := Response{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Data:    nil,
	}

	json.NewEncoder(w).Encode(data)
}

func NotAllowedResponse(w http.ResponseWriter) {
	data := Response{
		Code:    http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusMethodNotAllowed),
		Data:    nil,
	}

	json.NewEncoder(w).Encode(data)
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Product struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	StockQuantity int    `json:"stock_quantity"`
	Price         int    `json:"price"`
}

type Order struct {
	Id         int `json:"id"`
	UserId     int `json:"user_id"`
	ProductId  int `json:"product_id"`
	Quantity   int `json:"quantity"`
	Price      int `json:"price"`
	TotalPrice int `json:"total_price"`
}

type StockUpdate struct {
	Stock int `json:"stock"`
}
