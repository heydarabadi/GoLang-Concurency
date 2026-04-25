package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

func main() {
	//TestForHandleConcurencyInHandlerRequestWeb()
	//Example1ForContext()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// HTTP Server Concurrent در Go
// آیا در Go هر request همزمان اجرا می‌شود؟
//بله
//
//وقتی از این استفاده می‌کنی:
//
//go
//http.ListenAndServe(":8080", handler)
//برای هر درخواست جدید:
//
//یک goroutine جدید ساخته می‌شود
//handler تو داخل آن goroutine اجرا می‌شود
//درخواست‌ها کاملاً concurrent هستند
//پس این کد:
//
//go
//func handler(w http.ResponseWriter, r *http.Request) {
//    time.Sleep(5 * time.Second)
//}
//اگر 10 نفر همزمان صدا بزنند →
//
//10 goroutine همزمان اجرا می‌شود.
//
//⚠️ اما نکته مهم:
//اگر shared state داری:
//
//go
//var counter int
//و داخل handler تغییرش بدهی:
//
//go
//counter++
//⚠️ این data race ایجاد می‌کند.
//
//باید از یکی از این‌ها استفاده کنی:
//
//sync.Mutex
//sync/atomic
//Channel pattern
//Worker pool

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	counter int
	mu      sync.Mutex
)

func TestForHandleConcurencyInHandlerRequestWeb() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		counter++
		mu.Unlock()
		fmt.Fprintf(w, "Counter is : %d", counter)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//Context در Handlerها
//context.Context یکی از مهم‌ترین مفاهیم در Go است.
//
//در HTTP هر request به‌طور خودکار یک context دارد:
//
//go
//r.Context()
// Context چه کاربردهایی دارد؟
//Cancellation (اگر client قطع شود)
//Timeout
//Deadline
//انتقال metadata (مثل correlation id)

func Example1ForContext(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Fprintf(w, "Context is done")
	case <-time.After(1 * time.Second):
		fmt.Fprintf(w, "Context One Second TimeOut")
	default:
		return
	}
}

func Example2ForContext(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	result := make(chan string)

	go func() {
		time.Sleep(3 * time.Second)
		result <- "Finished"
	}()

	select {
	case res := <-result:
		fmt.Fprintln(w, res)
	case <-ctx.Done():
		http.Error(w, "Request Timeout", http.StatusRequestTimeout)
	}
}

type contextKey string

const CorrelationKey contextKey = "cid"

func correlationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := r.Header.Get("X-Correlation-ID")
		if cid == "" {
			cid = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), CorrelationKey, cid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

///چرا context حیاتی است؟
//بدون context:
//
//request قطع می‌شود
//ولی goroutine همچنان کار می‌کند
//DB query ادامه پیدا می‌کند
//memory leak ایجاد می‌شود
//resource هدر می‌رود
//با context:
//
//همه چیز cascade cancel می‌شود

//نکات Production-Level
// همیشه context را به لایه پایین‌تر پاس بده
//
// هرگز context را داخل struct نگه ندار
//
// از context برای optional param استفاده نکن
//
// از context.WithValue فقط برای metadata استفاده کن

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

///HTTP Client تنظیم‌شده
//توضیح کوتاه
//در Go بهتر است به جای http.Get یک http.Client با timeout و transport بسازی تا:
//
//درخواست‌ها معطل نشوند
//connection reuse شود
//مصرف منابع کنترل شود

func TestForHttpClient() {

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	resp, err := client.Get("https://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//Connection Pooling
//توضیح کوتاه
//Connection Pool یعنی connectionهای HTTP دوباره استفاده شوند تا برای هر request اتصال جدید ساخته نشود. این کار latency و مصرف CPU را کم می‌کند.
//

func ExampleConnectionPool() {

	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     60 * time.Second,
	}

	client := &http.Client{Transport: tr}

	for i := 0; i < 5; i++ {
		resp, err := client.Get("https://httpbin.org/get")
		if err != nil {
			panic(err)
		}

		fmt.Println("Request", i, resp.Status)
		resp.Body.Close()
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//Backpressure در I/O
//توضیح کوتاه
//Backpressure یعنی اگر مصرف‌کننده کند باشد، تولیدکننده مجبور شود صبر کند تا سیستم overload نشود.
//
//در Go معمولاً با channel buffer محدود انجام می‌شود.

func worker(jobs chan int) {
	for j := range jobs {
		fmt.Println("processing", j)
		time.Sleep(time.Second)
	}
}

func MainForBackpressure() {

	jobs := make(chan int, 3)

	go worker(jobs)

	for i := 0; i < 10; i++ {
		fmt.Println("sending", i)
		jobs <- i
	}

	close(jobs)

	time.Sleep(5 * time.Second)
}

//کار با فایل و دیسک
//توضیح کوتاه
//Go برای کار با فایل از پکیج‌های os و io استفاده می‌کند.
//
//برای فایل‌های بزرگ بهتر است streaming یا buffer استفاده شود.

func MainForIOWork() {

	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}

	file.WriteString("Hello Go File\n")
	file.Close()

	data, err := os.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
