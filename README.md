<div dir="rtl">

# مسیر یادگیری Concurrency در Go

این نقشه راه از مقدمات تا مباحث پیشرفته‌ی Concurrency در Go را پوشش می‌دهد. مباحث از سطح مفهومی شروع شده و تا طراحی سیستم‌های Production‑Grade ادامه پیدا می‌کنند.

---

## بخش 0 — مقدمه و پیش‌نیازها

پیش از ورود به concurrency باید بدانیم **چرا** و **کِی** باید از آن استفاده کنیم.

### موضوعات:
- اهمیت Concurrency در سیستم‌های مدرن
- تفاوت Concurrency و Parallelism
- مدل CSP در Go (Communicating Sequential Processes)
- مروری بر Runtime زبان Go
- آشنایی با Scheduler و مدل G-M-P
- مواردی که نباید از Concurrency استفاده کرد
- اصل سادگی در مقابل پیچیدگی

---

## بخش 1 — Goroutine و Channel

این بخش پایه‌ی Concurrency در Go است.

### موضوعات:
- تعریف Goroutine
- هزینه، رفتار و زمان‌بندی گوروتین‌ها
- تعریف Channel و مدل ارتباطی
- تفاوت Channelهای Buffered و Unbuffered
- ارسال و دریافت داده
- بستن Channel
- استفاده از `range` روی Channel
- Channelهای جهت‌دار:
    - `chan<-` برای ارسال
    - `<-chan` برای دریافت
- الگوی Done Channel برای Cancellation ساده

### تمرین‌ها:
- پیاده‌سازی Fan-out ساده با چند Worker
- ساخت یک Pipeline سه‌مرحله‌ای
- تشخیص و رفع Deadlockهای متداول

---

## بخش 2 — Select و الگوهای کنترلی

دستور `select` یکی از قدرتمندترین ابزارها برای کنترل جریان Concurrency در Go است.

### موضوعات:
- select غیرمسدودکننده
- استفاده از `default`
- ایجاد Timeout با `time.After`
- استفاده از select همراه با Context
- الگوهای رایج Concurrency:
    - Fan-in
    - Fan-out
    - Tee
    - Bridge
    - Multiplex
- مدیریت Backpressure با Buffer و Select

### تمرین‌ها:
- ترکیب چند Stream در یک خروجی
- پیاده‌سازی Timeout و Retry

---

## بخش 3 — همزمانی با حافظه مشترک (sync)

زمان‌هایی که Channels کافی نیستند و به همزمانی مبتنی بر داده نیاز داریم.

### موضوعات:
- `sync.Mutex`
- `sync.RWMutex` و زمان استفاده
- `sync.WaitGroup`
- `sync.Once`
- `sync.Cond`
- `sync.Map`
- عملیات اتمیک در `sync/atomic`
- عملگر CAS
- مدل حافظه در Go

### تمرین‌ها:
- ساخت Counter ایمن در مقابل concurrent access با:
    - Mutex
    - RWMutex
    - Atomic
    - مقایسه Performance
- ساخت Cache همزمان با RWMutex و sync.Map

---

## بخش 4 — Context و Cancellation

مدیریت lifecycle گوروتین‌ها در برنامه‌های واقعی.

### موضوعات:
- `context.Background`
- `context.TODO`
- `context.WithCancel`
- `context.WithTimeout`
- `context.WithDeadline`
- `context.WithValue`
- انتقال context بین لایه‌های مختلف
- Cancellation در Pipelineها

### تمرین‌ها:
- طراحی یک Web Crawler concurrent با قابلیت Cancel
- انجام عملیات I/O با Timeout

---

## بخش 5 — الگوهای طراحی Concurrency

الگوهایی که در برنامه‌های واقعی اهمیت بالایی دارند.

### موضوعات:
- Worker Pool پیشرفته
- تنظیم پویا یا ثابت تعداد Workerها
- صف کار (Work Queue)
- Backpressure
- Bounded Parallelism (الگوی Semaphore)
- Rate Limiting:
    - Token Bucket
- Pipelineهای قابل ترکیب
- الگوی Futures/Promises
- ساخت سیستم Pub/Sub ساده

### تمرین‌ها:
- اجرای N Job با حداکثر M Worker
- پیاده‌سازی Token Bucket Rate Limiter

---

## بخش 6 — مدیریت خطا و خاتمه امن

### موضوعات:
- Propagate کردن خطا در Pipelines
- Error wrapping همراه با `errors.Is` و `errors.As`
- جلوگیری از Goroutine Leak
- آزادسازی منابع با defer
- استفاده از `errgroup`

### تمرین:
- اجرای چند گوروتین که با اولین خطا، کل سیستم cancel شود.

---

## بخش 7 — Scheduler و Performance

برای نوشتن سیستم‌های performant باید Runtime را بشناسیم.

### موضوعات:
- GOMAXPROCS
- Scheduler داخلی Go
- Work Stealing
- Preemption
- Sysmon
- Netpoller
- تأثیر GC بر Concurrency
- Allocation Pressure
- False Sharing
- Cache Line Alignment

### تمرین‌ها:
- بنچمارک با GOMAXPROCSهای مختلف
- اندازه‌گیری Mutex Contention

---

## بخش 8 — Profiling و Observability

ابزارهای دیباگ Concurrency.

### موضوعات:
- CPU Profiling
- Heap Profiling
- Goroutine Profiling
- Block Profiling
- Mutex Profiling
- Go Trace
- runtime/metrics
- expvar
- Structured Logging
- Correlation ID

### ابزارها:
- `go tool pprof`
- `go tool trace`
- `net/http/pprof`

### تمرین‌ها:
- تشخیص Goroutine Leak
- شناسایی Mutex Contention

---

## بخش 9 — I/O همزمان و شبکه

### موضوعات:
- HTTP Server Concurrent
- استفاده از Context در Handlerها
- تنظیمات HTTP Client
- Connection Pooling
- WebSocket
- Backpressure در سیستم‌های I/O
- کار با فایل و دیسک

### تمرین‌ها:
- ساخت Reverse Proxy ساده
- استریم فایل‌های بزرگ با کنترل سرعت

---

## بخش 10 — ساختار داده و Concurrency

### موضوعات:
- Channel در مقابل Mutex
- داده‌های Immutable
- Copy-on-write
- Batched updates
- Ring Buffer بدون Lock
- جلوگیری از False Sharing با Padding

### تمرین‌ها:
- ساخت Queue بدون Lock
- ساخت Logger با Batching

---

## بخش 11 — تست‌نویسی برای Concurrency

Testing در سیستم‌های concurrent چالش بسیار مهمی است.

### موضوعات:
- Race Detector
    - `go test -race`
    - محدودیت‌های Race Detector
- تست‌های Flaky
- تست Deterministic
- Fake Clock
- Fuzz Testing

### تمرین‌ها:
- تست Worker Pool
- پیدا کردن Data Race

---

## بخش 12 — الگوهای Production

الگوهای ضروری برای سیستم‌های واقعی.

### موضوعات:
- Graceful Shutdown
- مدیریت Signalها
- Drain کردن Workerها
- Health Check
- Readiness Probe
- Circuit Breaker
- مدیریت Queue
- محدودیت منابع در Kubernetes
- متریک‌های سیستم

---

## بخش 13 — کتابخانه‌های مفید

کتابخانه‌های پرکاربرد در پروژه‌های concurrent.

- `golang.org/x/sync`
    - errgroup
    - semaphore
    - singleflight
- `golang.org/x/time/rate`
- OpenTelemetry
- zap
- zerolog
- fasthttp
- تنظیمات Connection Pool در `database/sql`

---

## بخش 14 — Anti‑Patterns

اشتباهات رایج که باعث باگ‌های سخت می‌شوند.

- Goroutine Leak
- استفاده اشتباه از Buffered Channel
- Deadlockهای پنهان
- Shared Mutable State بدون مراقبت
- استفاده بی‌رویه از `default` در Select
- Premature Optimization
- Panic در Goroutine بدون Recovery

</div>
