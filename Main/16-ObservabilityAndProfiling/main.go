package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	//MainForCpuProfile()
	//MainForMemProfile()
	//MainForGoRoutineProfile()
	MainForPprofInProduction()
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

///////////////
