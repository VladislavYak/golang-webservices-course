package teststests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// тут писать код тестов

func TestSearchServer(t *testing.T) {
	xml_path := "/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml"
	myData, _ := readXml(xml_path)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, &myData)
	}))
	fmt.Println("uxixui")
	fmt.Print(ts.URL)
	// time.Sleep(60 * time.Second)
}

// func TestCartCheckout(t *testing.T) {
// 	cases := []TestCase{
// 		{
// 			ID: "42",
// 			Result: &CheckoutResult{
// 				Status:  200,
// 				Balance: 100500,
// 				Err:     "",
// 			},
// 			IsError: false,
// 		},
// 		{
// 			ID: "100500",
// 			Result: &CheckoutResult{
// 				Status:  400,
// 				Balance: 0,
// 				Err:     "bad_balance",
// 			},
// 			IsError: false,
// 		},
// 		{
// 			ID:      "__broken_json",
// 			Result:  nil,
// 			IsError: true,
// 		},
// 		{
// 			ID:      "__internal_error",
// 			Result:  nil,
// 			IsError: true,
// 		},
// 	}

// 	ts := httptest.NewServer(http.HandlerFunc(CheckoutDummy))

// 	fmt.Print(ts.URL)
// 	time.Sleep(60 * time.Second)

// 	for caseNum, item := range cases {
// 		c := &Cart{
// 			PaymentApiURL: ts.URL,
// 		}
// 		result, err := c.Checkout(item.ID)

// 		fmt.Println("result", result)
// 		fmt.Println("----")

// 		if err != nil && !item.IsError {
// 			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
// 		}
// 		if err == nil && item.IsError {
// 			t.Errorf("[%d] expected error, got nil", caseNum)
// 		}
// 		if !reflect.DeepEqual(item.Result, result) {
// 			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Result, result)
// 		}
// 	}
// 	ts.Close()
// }
