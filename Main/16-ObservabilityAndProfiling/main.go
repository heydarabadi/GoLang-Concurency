package main

import (
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func main() {
	//MainForCpuProfile()
	//MainForMemProfile()
	//MainForGoRoutineProfile()
	//MainForPprofInProduction()
	//MainForExampleBefore()
	MainForObservability()
}

//CPU Profiling
//هدف: فهمیدن CPU دقیقاً کجا وقت می‌سوزاند
//
//این پروفایل نشان می‌دهد:
//
//کدام functionها بیشترین CPU را مصرف می‌کنند
//کدام goroutineها باعث مصرف زیاد CPU شده‌اند
//آیا contention یا busy-loop داری؟
//آیا تابعی که فکر می‌کردی سبک است، سنگین شده؟
//معمولاً برای تشخیص:
//
//حلقه‌های بی‌پایان
//locking شدید
//preemptionهای زیاد
//contention در goroutineها

func busy() {
	for i := 0; i < 1_000_000_000; i++ {
		_ = i * 2
	}
}

func MainForCpuProfile() {
	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	busy()
	time.Sleep(1 * time.Second) // تا CPU profile کامل نوشته شود
}

///////////////////////////////

////Heap Profiling
//هدف: فهمیدن حافظه کجا مصرف می‌شود
//
//Heap profile نشان می‌دهد:
//
//چه functionهایی بیشترین allocation دارند
//objectهای short-lived و long-lived چطورند
//آیا memory leak داری؟
//آیا allocation pressure باعث GC زیاد شده؟
//برای تشخیص:
//
//memory leak
//objectهای زیادی که در goroutineها ساخته می‌شود
//فشار زیاد روی GC
//رشد بیش‌ازحد heap

func leak() {
	var s [][]byte
	for i := 0; i < 10_000; i++ {
		b := make([]byte, 1024*50) // 50KB
		s = append(s, b)           // چون نگه می‌داریم → leak
	}
}

func MainForMemProfile() {
	leak()

	f, _ := os.Create("heap.pprof")
	pprof.WriteHeapProfile(f)
}

///////////////////////////////////////

////Goroutine Profiling
//هدف: دیدن وضعیت همهٔ goroutineها
//
//این پروفایل نشان می‌دهد:
//
//هر goroutine در چه stack traceی گیر کرده
//کدام goroutineها block شده‌اند؟
//آیا deadlock داری؟
//آیا تعداد goroutineها غیرعادی زیاد شده؟ (GOROUTINE LEAK)
//آیا channel/lock باعث block شده؟
//بسیار مهم در debugging concurrency.

//در flamegraph می‌بینی:
//
//چند goroutine در waitForever
//کدام‌ها block هستند
//روی چه channelهایی گیر کرده‌اند
//برای deadlock debugging بسیار حیاتی است.

func waitForever() {
	ch := make(chan struct{})
	<-ch // block
}

func MainForGoRoutineProfile() {
	go waitForever()
	go waitForever()
	go waitForever()

	time.Sleep(500 * time.Millisecond)

	// روش صحیح - فرمت باینری
	f, _ := os.Create("goroutine.pprof")
	pprof.Lookup("goroutine").WriteTo(f, 0) // ← عدد 0 به جای 2
	f.Close()
}

////////////////////////////////////////

//// استفاده از pprof روی HTTP (محیط Production)

//Go به‌صورت built‑in یک سرور pprof دارد:

// بعد می‌توانی:
//
// http://localhost:6060/debug/pprof/
// از اینجا دانلود کنی:
//
// /profile → CPU
// /heap → Heap
// /goroutine → Goroutine stack
// /block → Blocking profile (برای lock contention)
// /mutex → Mutex contention
// این بهترین روش در production است.

// import
//
//	_ "net/http/pprof"  // ← import جانبی (side-effect)
func MainForPprofInProduction() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	time.Sleep(400 * time.Second)

}

////////////////

////چه زمانی کدام را استفاده کنیم؟
//CPU Profile
//وقتی:
//
//برنامه کند است
//هسته‌ها 100٪ هستند
//goroutineها busy-loop شده‌اند
//lock contention داری
//Heap Profile
//وقتی:
//
//حافظه رشد می‌کند
//GC زیاد اجرا می‌شود
//memory leak مشکوک است
//objectها زیاد ساخته می‌شوند
//Goroutine Profile
//وقتی:
//
//deadlock داری
//goroutine leak داری
//تعداد goroutineها ناگهان زیاد شده
//برنامه hang می‌کند

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

///Block Profiling
//برای بررسی جایی که goroutineها منتظر می‌مانند (blocking) استفاده می‌شود.
//
//یعنی نشان می‌دهد goroutineها کجا در کد منتظر resource می‌مانند مثل:
//
//channel send/receive
//mutex
//select
//sync primitives
//مثال مشکل‌هایی که پیدا می‌کند:
//
//goroutineهایی که روی channel قفل شده‌اند
//deadlockهای احتمالی
//contention روی منابع

// فعال‌سازی:
//runtime.SetBlockProfileRate(1)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Mutex Profiling چیست؟
//نشان می‌دهد:
//
//کجا mutex زیاد گرفته می‌شود
//کجا goroutineها برای lock منتظر می‌مانند
//چه lockهایی performance را خراب کرده‌اند

//// فعال‌سازی:
//runtime.SetMutexProfileFraction(1)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//Go Trace چیست؟
//قوی‌ترین ابزار runtime در Go است.
//
//اطلاعات دقیق از:
//
//زمان‌بندی goroutineها
//GC
//syscalls
//blocking
//scheduler behavior

// اجرا:
//go tool trace trace.out

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//runtime/metrics چیست؟
//یک API برای خواندن متریک‌های داخلی runtime مثل:
//
//تعداد goroutine
//heap size
//GC cycles
//scheduler stats

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var mu sync.Mutex
var counter int

func MainForExampleBefore() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	traceFile, _ := os.Create("trace.out")
	defer traceFile.Close()

	trace.Start(traceFile)
	defer trace.Stop()

	ch := make(chan int)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				time.Sleep(1 * time.Second)
				mu.Unlock()
			}
		}()
	}

	go func() {
		ch <- 1
	}()

	time.Sleep(3 * time.Second)
	samples := []metrics.Sample{
		{Name: "/sched/goroutines:goroutines"},
		{Name: "/memory/classes/heap/objects:bytes"},
	}

	metrics.Read(samples)

	fmt.Println("Number of Goroutines:", samples[0].Value.Uint64())
	fmt.Println("Heap Objects (bytes):", samples[1].Value.Uint64())

	// ذخیره block profile
	blockFile, _ := os.Create("block.prof")
	defer blockFile.Close()
	pprof.Lookup("block").WriteTo(blockFile, 0)

	// ذخیره mutex profile
	mutexFile, _ := os.Create("mutex.prof")
	defer mutexFile.Close()
	pprof.Lookup("mutex").WriteTo(mutexFile, 0)

	fmt.Println("Done. Profiles generated.")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//expvar چیست؟
//expvar یک پکیج داخلی Go است برای:
//
//اکسپوز کردن متریک‌ها از طریق HTTP
//مشاهده وضعیت داخلی برنامه
//تعریف متغیرهای global که از طریق وب قابل مشاهده باشند
//به‌صورت اتوماتیک یک endpoint می‌سازد:

// به‌صورت اتوماتیک یک endpoint می‌سازد:
///debug/vars

//وقتی Go برنامه را اجرا می‌کنی، خروجی JSON از متریک‌های runtime می‌دهد مثل:
//
//تعداد goroutineها
//حافظه مصرف‌شده
//GC stats
//متریک‌های custom که خودت اضافه می‌کنی

//مثال یک counter ساده:

// var reqCount = expvar.NewInt("requests_total")
//reqCount.Add(1)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

////Structured Logging چیست؟
//در لاگ‌نویسی معمولی، فقط متن ساده می‌نویسی:
//
//text
//User logged in
//Database saved
//Error happened
//این‌ها برای انسان خوب هستند ولی برای ماشین فاجعه هستند، چون قابل سرچ دقیق، فیلتر، تجزیه، پارس نیستند.
//
//در Structured Logging لاگ‌ها به‌صورت JSON یا کلید/مقدار ثبت می‌شوند:
//
//json
//{
//  "level": "info",
//  "msg": "user logged in",
//  "user_id": 42,
//  "ip": "192.168.1.10",
//  "time": "2026-04-26T01:02:00Z"
//}
//مزایا:
//
//بهتر دیده می‌شود (ELK / Loki / Grafana / Datadog)
//قابل فیلتر و جستجو
//شامل متادیتای کلیدی (userId, requestId, ip …)
//در Go بهترین Loggerهای structured:
//
//Zerolog
//Zap
//Logrus (نسبتاً قدیمی‌تر)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//Correlation ID چیست؟
//اگر سیستم تو چند سرویس دارد:
//
//text
//API → Auth Service → Payment Service → Database
//و یک درخواست مشکل دارد، چطور تمام لاگ مسیر همان درخواست را پیدا می‌کنی؟
//
//اینجاست که Correlation ID کمک می‌کند:
//
//یک شناسه یکتا (UUID) برای هر request ایجاد می‌شود
//در header پاس داده می‌شود
//همه سرویس‌ها همان ID را در لاگ‌های خود می‌نویسند
//نتیجه:
//
//می‌توانی تمام لاگ‌های مربوط به یک درخواست خاص را با یک سرچ ساده پیدا کنی.
//
//مثال ID:
//
//text
//X-Correlation-ID: 8a3c1e84-9d12-4c77-a4ac-22dfcd3d9135
//یا اگر نبود، خودت می‌سازی.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var reqCount = expvar.NewInt("request_count")

func MainForObservability() {
	logger := zerolog.New(log.Writer()).With().Timestamp().Logger()

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		reqCount.Add(1)
		corrID := GetCorId(r)

		logger.Info().
			Str("corrID", corrID).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Msg("Request received.")

		response := map[string]string{
			"message": "Hello World!",
			"corrID":  corrID,
		}

		json.NewEncoder(w).Encode(response)
	})
	mux.Handle("/debugs/vars/", expvar.Handler())

	handler := correlationMiddleware(logger)(mux)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", handler)
}

func correlationMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return hlog.NewHandler(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			corrID := GetCorId(r)

			w.Header().Set("X-Correlation-ID", corrID)
			ctx := context.WithValue(r.Context(), "cid", corrID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}))
	}
}

func GetCorId(r *http.Request) string {
	id := r.Header.Get("corid")
	if id == "" {
		id = uuid.New().String()
	}
	return id
}
