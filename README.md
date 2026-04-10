مسیر یادگیری (Learning Roadmap)
بخش 0 — مقدمه و پیش‌نیازها
قبل از ورود به مبحث concurrency باید بفهمیم چرا اصلاً به concurrency نیاز داریم و در چه مواقعی استفاده از آن منطقی است.

موضوعات این بخش:

چرا Concurrency مهم است؟
تفاوت Concurrency و Parallelism
مدل CSP در Go (Communicating Sequential Processes)
مروری بر Runtime زبان Go
آشنایی با Scheduler و مدل G-M-P
چه زمانی نباید از Concurrency استفاده کرد
سادگی در مقابل پیچیدگی
بخش 1 — Goroutine و Channel
این بخش پایه و اساس concurrency در Go است.

موضوعات این بخش:

Goroutine چیست
هزینه و نحوه زمان‌بندی goroutineها
Channel چیست
تفاوت Buffered و Unbuffered Channel
ارسال و دریافت داده
بستن channel
استفاده از range روی channel
Channelهای جهت‌دار
text
chan<-   ارسال
<-chan   دریافت
الگوی Done Channel برای cancel ساده
تمرین‌ها
پیاده‌سازی Fan-out ساده با چند worker
ساخت یک Pipeline سه مرحله‌ای
تشخیص و رفع Deadlockهای رایج
بخش 2 — Select و الگوهای کنترلی
یکی از قدرتمندترین ابزارهای concurrency در Go دستور select است.

موضوعات این بخش:

select غیرمسدودکننده
استفاده از default
timeout با time.After
استفاده از select همراه با context
الگوهای concurrency
الگوها:

Fan-in

Fan-out

Tee

Bridge

Multiplex

مدیریت Backpressure با buffer و select

تمرین‌ها
ترکیب چند stream داده در یک خروجی
پیاده‌سازی Timeout و Retry
بخش 3 — همزمانی با حافظه مشترک (sync)
گاهی ارتباط از طریق channel کافی نیست و نیاز به دسترسی مشترک به حافظه داریم.

موضوعات این بخش:

sync.Mutex
sync.RWMutex
زمان استفاده از هرکدام
sync.WaitGroup
sync.Once
sync.Cond
sync.Map
عملیات اتمیک (sync/atomic)
عملیات CAS
مدل حافظه در Go
تمرین‌ها
پیاده‌سازی یک Counter امن برای concurrent access

روش‌ها:

با Mutex
با RWMutex
با Atomic
سپس مقایسه performance.

تمرین دیگر:

ساخت Cache همزمان و مقایسه:

RWMutex
sync.Map
بخش 4 — Context و Cancellation
مدیریت cancellation در سیستم‌های concurrent بسیار مهم است.

موضوعات:

context.Background
context.TODO
context.WithCancel
context.WithTimeout
context.WithDeadline
context.WithValue
انتقال context بین APIها
cancellation در pipelineها
تمرین‌ها
ساخت Web Crawler concurrent با امکان لغو
انجام عملیات I/O با timeout
بخش 5 — الگوهای طراحی Concurrency
الگوهایی که در سیستم‌های واقعی استفاده می‌شوند.

موضوعات:

Worker Pool پیشرفته
تنظیم تعداد workerها
صف کار
Backpressure
Bounded Parallelism (الگوی semaphore)
Rate Limiting
الگوریتم Token Bucket
Pipelineهای قابل ترکیب
Futures/Promises در Go
سیستم Pub/Sub ساده
تمرین‌ها
اجرای N job با حداکثر M worker
پیاده‌سازی Token Bucket Rate Limiter
بخش 6 — مدیریت خطا و خاتمه امن
موضوعات:

انتقال خطا در pipelineها
error wrapping (errors.Is و errors.As)
جلوگیری از Goroutine Leak
آزادسازی منابع با defer
استفاده از errgroup
تمرین
اجرای چند goroutine که با اولین خطا همه cancel شوند.

بخش 7 — Scheduler و Performance
برای نوشتن سیستم‌های سریع باید بفهمیم runtime چگونه کار می‌کند.

موضوعات:

GOMAXPROCS
Scheduler داخلی Go
Work stealing
Preemption
Sysmon
Netpoller
تعامل GC با concurrency
Allocation pressure
False sharing
Cache line alignment
تمرین‌ها
بنچمارک با GOMAXPROCS مختلف
اندازه‌گیری Mutex contention
بخش 8 — Profiling و Observability
دیباگ برنامه‌های concurrent.

موضوعات:

CPU profiling
Heap profiling
Goroutine profiling
Block profiling
Mutex profiling
Go Trace
runtime metrics
expvar
Structured logging
Correlation ID
ابزارها:

text
go tool pprof
go tool trace
net/http/pprof
تمرین‌ها
پیدا کردن goroutine leak
تشخیص mutex contention
بخش 9 — I/O همزمان و شبکه
کار با concurrency در برنامه‌های شبکه‌ای.

موضوعات:

HTTP server concurrent
استفاده از context در handlerها
تنظیمات HTTP Client
Connection pooling
WebSocket
Backpressure
فایل و دیسک
تمرین‌ها
ساخت Reverse Proxy ساده
استریم فایل بزرگ با کنترل سرعت
بخش 10 — ساختار داده و Concurrency
موضوعات:

Channel vs Mutex
داده‌های immutable
Copy-on-write
Batched updates
Ring buffer بدون lock
Padding ساختارها برای جلوگیری از false sharing
تمرین‌ها
پیاده‌سازی Queue بدون lock
ساخت Logger با batching
بخش 11 — تست‌نویسی برای Concurrency
تست برنامه‌های concurrent سخت است.

موضوعات:

Race detector
text
go test -race
محدودیت‌های race detector
تست‌های flaky
تست deterministic
Fake clock
Fuzz testing
تمرین‌ها
تست Worker Pool
پیدا کردن Data Race
بخش 12 — الگوهای Production
الگوهایی که در سیستم‌های واقعی استفاده می‌شوند.

موضوعات:

Graceful shutdown
مدیریت signalها
تخلیه workerها
Health check
Readiness probe
Circuit breaker
مدیریت queue
محدودیت منابع در Kubernetes
متریک‌های سیستم
بخش 13 — کتابخانه‌های مفید
کتابخانه‌هایی که در سیستم‌های concurrent زیاد استفاده می‌شوند.

golang.org/x/sync

errgroup
semaphore
singleflight
golang.org/x/time/rate

OpenTelemetry

zap

zerolog

fasthttp

تنظیمات connection pool در database/sql

بخش 14 — Anti‑Patterns
اشتباهات رایج در برنامه‌های concurrent.

Goroutine leak
استفاده اشتباه از buffered channel
Deadlockهای پنهان
shared mutable state
استفاده زیاد از select default
premature optimization
panic در goroutine بدون recovery