package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func main() {
	//GOMAXPROCS()
	//MainForPreemption()
	//MainForWorkSteal()
	//MainForSysmon()
	MainForFalseSharing()
}

// GOMAXPROCS

//GOMAXPROCS مشخص می‌کند چند CPU logical به‌طور هم‌زمان اجازهٔ اجرای goroutineها را دارند.
//
//مقدار پیش‌فرض = تعداد CPUهای سیستم
//اگر GOMAXPROCS = 1 باشد، حتی اگر 16 هسته داشته باشی، فقط یک Thread از Go Runtime هم‌زمان کار می‌کند
//زیاد کردن آن = افزایش parallelism
//خیلی زیاد کردنش = همیشه کارآمد نیست (context switch زیاد)

func GOMAXPROCS() {
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	runtime.GOMAXPROCS(2)
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
}

//2) Scheduler داخلی Go
//Scheduler مسئول مدیریت goroutine → OS thread است.
//
//این مدل به اختصار G-M-P نام دارد:
//
//G = Goroutine
//M = Machine (OS thread)
//P = Processor slot (اجازهٔ اجرا)
//توضیح خیلی ساده:
//
//هر P فقط روی یک CPU logical اجرا می‌شود و تنها وقتی که P داشته باشی، goroutine اجرا می‌شود
//تعداد P = مقدار GOMAXPROCS
//هر M یک P را “اجرا” می‌کند
//Goroutineها روی صف‌های مختلف قرار می‌گیرند و scheduler آنها را بین Pها پخش می‌کند

//3) Work Stealing
//استراتژی scheduler برای افزایش کارایی.
//
//هر P یک صف محلی از goroutineها دارد.
//
//اگر یک P کارش تمام شود، می‌رود و از Pهای دیگر goroutine می‌دزدد.
//
//مزیت:
//
//کاهش contention روی صف global
//کارکرد بهتر روی CPUهای چند هسته‌ای
//load balancing طبیعی

//4) Preemption
//Preemption یعنی scheduler بتواند یک goroutine طولانی را به زور قطع کند تا goroutineهای دیگر starving نشوند.
//
//قدیم Go فقط در نقاط خاص (function call) preemption داشت.
//
//از Go 1.14 به بعد:
//
//async preemption اضافه شده
//اگر یک goroutine CPU-bound باشد (حلقهٔ بی‌نهایت)، runtime آن را قطع می‌کند و اجازه می‌دهد دیگران هم اجرا شوند

///////////////////////////////////////////////////

// Sample For preemption + work stealing

func busy() {
	for {

	}
}

func MainForPreemption() {
	runtime.GOMAXPROCS(1)

	for i := 0; i < 3; i++ {
		go busy()
	}

	for {
		fmt.Println("main alive")
		time.Sleep(time.Second * 1)
	}
}

func Work(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for i := 0; i < 50_000_000; i++ {
		sum += i

	}
	fmt.Println("done", id)
}

func MainForWorkSteal() {
	runtime.GOMAXPROCS(4)

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go Work(i, &wg)
	}

	wg.Wait()
	fmt.Println("done")
}

//Sysmon (System Monitor)
//sysmon یک goroutine خاص داخل runtime است که دائماً در پس‌زمینه اجرا می‌شود و وظیفهٔ نظارت روی سیستم runtime را دارد.
//
//کارهای اصلی:
//
//تشخیص goroutineهای طولانی CPU-bound
//کمک به preemption
//بیدار کردن goroutineهای sleep شده
//کمک به GC
//مدیریت netpoller
//جلوگیری از block شدن scheduler
//به زبان ساده:
//
//مثل یک ناظر دائمی runtime است که هر چند میلی‌ثانیه وضعیت سیستم را چک می‌کند.

//اگر sysmon و preemption نبود:
//
//goroutine busy() کل CPU را می‌گرفت
//main دیگر اجرا نمی‌شد
//اما sysmon کمک می‌کند runtime آن goroutine را قطع کند.

func MainForSysmon() {
	go busy()

	for i := 0; i < 5; i++ {
		fmt.Println("main running")
	}
}

//Netpoller
//Netpoller سیستم داخلی Go برای مدیریت I/O غیر بلاک‌کننده است.
//
//وقتی در Go این کارها را انجام می‌دهی:
//
//network
//HTTP
//socket
//file descriptor
//runtime از netpoller استفاده می‌کند.
//
//در لینوکس معمولاً با:
//
//text
//epoll
//در macOS:
//
//text
//kqueue
//در ویندوز:
//
//text
//IOCP
//هدف:
//
//یک goroutine هنگام I/O OS thread را block نکند
//وقتی داده آماده شد، goroutine دوباره schedule شود

//اگر 10k درخواست بیاید:
//
//Go برای هر request یک goroutine می‌سازد
//وقتی منتظر network است → goroutine park می‌شود
//netpoller وقتی socket آماده شد → goroutine را بیدار می‌کند
//به همین دلیل Go می‌تواند ده‌ها هزار connection را مدیریت کند.

func MainForNetpoller() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request received")
		w.Write([]byte("hello"))
	})

	http.ListenAndServe(":8080", nil)
}

//تأثیر GC بر Concurrency
//Go از Concurrent Garbage Collector استفاده می‌کند.
//
//هدف:
//
//کاهش Stop The World
//GC مراحل زیر را دارد:
//
//mark start (stop-the-world کوتاه)
//concurrent mark
//mark termination
//sweep
//در بیشتر زمان GC همزمان با برنامه اجرا می‌شود.
//
//اما GC هنوز کمی pause ایجاد می‌کند.
//
//مشکل زمانی رخ می‌دهد که:
//
//allocation خیلی زیاد باشد
//objectهای کوتاه‌عمر زیاد باشند

//این کد:
//
//مدام memory allocate می‌کند
//GC باید دائماً اجرا شود
//throughput برنامه پایین می‌آید

func MainForGc() {
	for {
		_ = make([]byte, 1024)
	}
}

//Allocation Pressure
//Allocation pressure یعنی:
//
//سرعت allocate شدن objectها بیشتر از سرعت GC برای پاکسازی آنها باشد.
//
//نتیجه:
//
//GC بیشتر اجرا می‌شود
//CPU مصرف GC بالا می‌رود
//latency زیاد می‌شود
//معمولاً در این موارد دیده می‌شود:
//
//ساخت object زیاد در loop
//slice زیاد
//map زیاد
//string concatenation زیاد

func BadPracticeForAllocation() {
	for i := 0; i < 1_000_000; i++ {
		s := []int{1, 2, 3, 4, 5}
		_ = s
	}

}

func BestPracticeForAllocation() {
	s := make([]int, 5)

	for i := 0; i < 1_000_000; i++ {
		s[0] = 1
		s[1] = 2
	}

}

/// Cache Line چیست؟
//CPU مستقیماً با RAM کار نمی‌کند.
//
//داده‌ها را در Cache نگه می‌دارد.
//
//واحد انتقال داده در Cache معمولاً:
//
//64bytes
//
//به این 64 بایت می‌گویند:
//
//Cache Line
//
//یعنی اگر فقط یک int 8 بایتی بخوانی، در واقع کل 64 بایت اطرافش وارد cache می‌شود

//1️⃣ Cache Line دقیقاً چیست؟
//CPU داده را بایت‌به‌بایت از RAM نمی‌خواند.
//
//بلکه داده را در بلوک‌های ثابت می‌خواند که معمولاً:
//
//64
//bytes
//64bytes
//
//به این بلوک می‌گویند:
//
//Cache Line
//
//یعنی اگر فقط یک int64 (۸ بایت) بخوانی، CPU در واقع کل ۶۴ بایت اطراف آن را وارد Cache می‌کند.
//
//2️⃣ چرا Cache Line وجود دارد؟
//چون:
//
//RAM خیلی کندتر از CPU است (صدها سیکل تأخیر)
//انتقال داده تکی گران است
//spatial locality معمولاً وجود دارد
//اصل locality:
//
//🔹 Spatial Locality
//اگر به یک آدرس دسترسی داشته باشی، احتمال زیاد به اطرافش هم دسترسی خواهی داشت.
//
//پس CPU می‌گوید:
//
//بیا کل 64 بایت اطرافش را بیاورم.
//
//3️⃣ ساختار Cache (خیلی خلاصه ولی کاربردی)
//معمولاً CPU چند سطح Cache دارد:
//
//سطح	اندازه	سرعت
//L1	کوچک (~32KB)	بسیار سریع
//L2	متوسط	سریع
//L3	بزرگ‌تر	کندتر
//Cache Line در همهٔ این‌ها واحد پایه است.
//
//4️⃣ مثال واقعی در حافظه
//فرض کن این struct را داریم:
//
//go
//type Data struct {
//	A int64  // 8 bytes
//	B int64  // 8 bytes
//	C int64  // 8 bytes
//	D int64  // 8 bytes
//}
//هر int64 = 8 بایت
//
//کل struct = 32 بایت
//
//پس کل این struct داخل نصف یک Cache Line جا می‌گیرد.
//
//اگر دو struct پشت هم در memory باشند:
//
//text
//[ Data1 (32B) ][ Data2 (32B) ]
//کل 64 بایت → دقیقاً یک Cache Line.
//
//5️⃣ مهم‌ترین بخش: Cache Coherency
//در سیستم چند هسته‌ای:
//
//هر CPU هسته cache خودش را دارد.
//
//اگر:
//
//CPU1 مقدار یک آدرس را تغییر دهد
//CPU2 همان Cache Line را داشته باشد
//باید یک پروتکل اجرا شود.
//
//معمولاً:
//
//MESI Protocol
//
//حالت‌های اصلی:
//
//M (Modified)
//E (Exclusive)
//S (Shared)
//I (Invalid)
//وقتی یک CPU روی یک متغیر بنویسد:
//
//✅ کل Cache Line برای بقیه CPUها invalid می‌شود.

////////////////////////

///False Sharing چیست؟
//False Sharing زمانی رخ می‌دهد که:
//
//دو goroutine روی دو CPU مختلف
//روی دو متغیر جداگانه کار می‌کنند
//اما آن دو متغیر داخل یک Cache Line هستند
//CPUها فکر می‌کنند روی یک داده مشترک کار می‌شود
//
//و شروع می‌کنند cache را invalidate کردن.
//
//نتیجه:
//
//cache thrashing
//کندی شدید
//حتی 5 تا 10 برابر افت performance
//در حالی که هیچ data race‌ای وجود ندارد!

type Counter struct {
	A int64
	B int64
}

func MainForFalseSharing() {
	var wg sync.WaitGroup
	c := Counter{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000_000; i++ {
			c.A++
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000_000; i++ {
			c.B++
		}
	}()

	wg.Wait()
	fmt.Println("All Done", c.A, c.B)
}

////Cache Line Alignment چیست؟
//راه‌حل:
//
//باید متغیرهایی که توسط CPUهای مختلف نوشته می‌شوند را
//
//در cache lineهای جدا قرار دهیم.
//
//در Go می‌توانیم padding اضافه کنیم.

type Counter2 struct {
	A int64
	_ [56]byte // padding تا کامل شدن 64 bytes
	B int64
}

type PaddedInt64 struct {
	_     [64]byte
	Value int64
}

type Counter3 struct {
	A int64
	_ [7]int64
	B int64
}
