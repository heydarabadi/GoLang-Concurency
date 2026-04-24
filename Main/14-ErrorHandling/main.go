package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	MainForErrorWrapping()
}

// Propagate کردن خطا در Pipelines

func generate(numbers ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range numbers {
			out <- n
		}
	}()
	return out
}

func square(in <-chan int) (<-chan int, <-chan error) {
	out := make(chan int)
	errch := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errch)

		for n := range in {
			if n < 0 {
				errch <- fmt.Errorf("عدد منفی غیرمجاز: %d", n)
				return
			}
			out <- n * n
		}
	}()

	return out, errch
}

func MainForPropagate() {
	numbers := generate(2, 3, -1, 4)
	squares, errch := square(numbers)

	for sq := range squares {
		fmt.Println(sq)
	}

	if err := <-errch; err != nil {
		fmt.Println("خطا:", err)
	}
}

// Error wrapping
var (
	ErrNotFound   = errors.New("موردی یافت نشد")
	ErrPermission = errors.New("دسترسی غیرمجاز")
	ErrTimeout    = errors.New("زمان به پایان رسید")
)

type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("خطای اعتبارسنجی در فیلد '%s' با مقدار %v: %s",
		e.Field, e.Value, e.Message)
}

func processData(id int) error {
	if id <= 0 {
		return fmt.Errorf("پردازش داده %d: %w", id, ErrNotFound)
	}

	if id > 1000 {
		return fmt.Errorf("پردازش داده %d: %w", id, ErrPermission)
	}

	if id == 42 {
		validationErr := ValidationError{
			Field:   "id",
			Value:   42,
			Message: "عدد ممنوعه",
		}
		return fmt.Errorf("پردازش ناموفق: %w", validationErr)
	}

	return nil
}

func MainForErrorWrapping() {
	err := processData(-5)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("خطای not found شناسایی شد:", err)
	}

	err = processData(42)
	var valErr ValidationError
	if errors.As(err, &valErr) {
		fmt.Printf("خطای اعتبارسنجی: فیلد=%s, پیام=%s\n",
			valErr.Field, valErr.Message)
	}

	err = processData(2000)
	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Println("موردی پیدا نشد")
	case errors.Is(err, ErrPermission):
		fmt.Println("دسترسی ندارید")
	default:
		fmt.Println("خطای ناشناخته:", err)
	}
}

// ❌ مثال بد: گوروتینی که هیچوقت آزاد نمی‌شود
func badWorker(done <-chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			// کار انجام بده
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ✅ مثال خوب: استفاده از context برای کنترل
func goodWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("گوروتین خاتمه یافت:", ctx.Err())
			return
		default:
			// شبیه‌سازی کار
			time.Sleep(100 * time.Millisecond)
			fmt.Println("کار در حال انجام...")
		}
	}
}

// ❌ مثال بد: فراموش کردن membaca از channel
func leakExample() {
	ch := make(chan int)

	go func() {
		// این گوروتین تا ابد منتظر می‌ماند
		val := <-ch // هیچکس به این channel نمی‌نویسد
		fmt.Println(val)
	}()

	// تابع تمام می‌شود ولی گوروتین باقی می‌ماند (LEAK)
}

// ✅ راه حل: استفاده از buffered channel یا context
func safeExample() {
	ch := make(chan int, 1) // بافر شده

	go func() {
		select {
		case val := <-ch:
			fmt.Println(val)
		case <-time.After(1 * time.Second):
			fmt.Println("زمان تمام شد")
		}
	}()

	// ارسال داده یا timeout
	select {
	case ch <- 42:
		fmt.Println("داده ارسال شد")
	case <-time.After(500 * time.Millisecond):
		fmt.Println("ارسال ناموفق")
	}
}

func MainForGoroutineLeak() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go goodWorker(ctx)

	time.Sleep(3 * time.Second) // بعد از 2 ثانیه context timeout می‌دهد
	fmt.Println("پایان main")
}

// مثال 1: بستن فایل
func readFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close() // حتی در صورت خطاهای بعدی بسته می‌شود

	// پردازش فایل
	data := make([]byte, 100)
	_, err = f.Read(data)
	if err != nil {
		return err // defer باز هم اجرا می‌شود
	}

	return nil
}

// مثال 2: قفل با defer
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // حتی در panic هم آزاد می‌شود

	c.count++
	// اگر panic رخ دهد، defer اجرا می‌شود
}

// مثال 3: تراکنش دیتابیس
func processTransaction(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Commit یا Rollback تضمین شده
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // panic را دوباره raise می‌کند
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// عملیات دیتابیس
	if _, err = tx.Exec("INSERT INTO users VALUES (1, 'ali')"); err != nil {
		return err
	}

	if _, err = tx.Exec("INSERT INTO orders VALUES (1, 100)"); err != nil {
		return err
	}

	return nil
}

// مثال 4: چندین defer (به ترتیب معکوس اجرا می‌شوند)
func multipleDefer() {
	defer fmt.Println("اولین defer (آخر اجرا می‌شود)")
	defer fmt.Println("دومین defer (دوم اجرا می‌شود)")
	defer fmt.Println("سومین defer (اول اجرا می‌شود)")
	fmt.Println("بدنه تابع")
}

func MainForDefer() {
	multipleDefer()
	// خروجی:
	// بدنه تابع
	// سومین defer (اول اجرا می‌شود)
	// دومین defer (دوم اجرا می‌شود)
	// اولین defer (آخر اجرا می‌شود)
}

// تابعی که یک سرویس رو شبیه‌سازی می‌کنه
func callService(ctx context.Context, name string, delay time.Duration, shouldFail bool) (string, error) {
	select {
	case <-time.After(delay): // شبیه‌سازی کار شبکه
		if shouldFail {
			return "", fmt.Errorf("%s failed!", name)
		}
		return fmt.Sprintf("result from %s", name), nil
	case <-ctx.Done(): // اگه context کنسل بشه
		return "", ctx.Err()
	}
}

func MainForErrGroup() {
	// errgroup.WithContext یه context جدید برمی‌گردونه که اگر توی یکی از
	// گوروتین‌ها خطایی رخ بده، خودکار کنسل (cancel) می‌شه.
	g, ctx := errgroup.WithContext(context.Background())

	// تعریف سه تا وظیفه (task) که هرکدوم تو یه گوروتین جدا اجرا می‌شن
	var results []string

	// وظیفه ۱: با موفقیت انجام می‌شه
	g.Go(func() error {
		res, err := callService(ctx, "Service A", 100*time.Millisecond, false)
		if err == nil {
			results = append(results, res)
		}
		return err
	})

	// وظیفه ۲: با موفقیت انجام می‌شه
	g.Go(func() error {
		res, err := callService(ctx, "Service B", 150*time.Millisecond, false)
		if err == nil {
			results = append(results, res)
		}
		return err
	})

	// وظیفه ۳: این یکی قراره خطا بده!
	g.Go(func() error {
		// این سرویس بعد از ۵۰ میلی‌ثانیه خطا برمی‌گردونه
		res, err := callService(ctx, "Service C", 50*time.Millisecond, true)
		if err == nil {
			results = append(results, res)
		}
		return err
	})

	// منتظر می‌مونه تا همه گوروتین‌ها تموم بشن.
	// اگه توی هرکدوم از اون‌ها خطایی برگردونده شده باشه،
	// اولین خطا رو بهمون میده.
	if err := g.Wait(); err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		// خروجی: Error occurred: Service C failed!
	} else {
		fmt.Println("All tasks completed successfully!")
		fmt.Println("Results:", results)
	}
}
