package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

}

// /////////////////////////////////////////////
// پایه و مادر تمام contextهاست.
// هیچ deadline، cancel یا value ندارد.
// معمولاً در تابع اصلی (main) یا پایین‌ترین سطح استفاده می‌شود.
// go
func BackGround() {
	_ = context.Background()

}

///////////////////////////////////////////////

// مثل Background است، ولی یعنی به‌صورت موقت استفاده شود.
// زمانی که هنوز نمی‌دانی چه context درست باید بدهی.

// Usage:
// 🔹 استفاده: وقتی هنوز تصمیم نگرفته‌ای از Background، WithCancel یا WithTimeout استفاده کنی. (معمولاً در توسعه اولیه)
func Todo() {
	_ = context.TODO()
}

///////////////////////////////////////////////

// یک context می‌سازد که می‌توانی **به‌صورت دستی لغو (cancel)** کنی.
// تابعی cancel() برمی‌گرداند.

// Usage:
// وقتی خودت می‌خواهی لغو را کنترل کنی (مثل دکمه stop یا error در یکی از goroutine‌ها).
func Cancel() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		fmt.Println("cancelled")
	}()

	cancel()
}

///////////////////////////////////////////////

// یک context می‌سازد که خودش بعد از زمان مشخصی لغو می‌شود.
// عد از 2 ثانیه، ctx.Done() فعال می‌شود و لغو اتفاق می‌افتد.
// استفاده برای API call، query، یا عملیات‌هایی که نباید بی‌نهایت طول بکشد.
func TimeOut() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Work")
	case <-ctx.Done():
		fmt.Println("TimeOut Reached:", ctx.Err())
	}
}

///////////////////////////////////////////////

// مثل WithTimeout است، ولی به‌جای مدت زمان، زمان دقیق (absolute time) می‌دهی.

// در زمان مشخص لغو می‌شود.
// استفاده برای کارهایی که تا لحظه خاص باید تمام شوند.
func Deadline() {
	deadLine := time.Now().Add(3 * time.Second)
	_, cancel := context.WithDeadline(context.Background(), deadLine)
	defer cancel()
}

///////////////////////////////////////////////

// برای حمل داده سبک (metadata) بین goroutine‌ها.
// استفاده برای داده‌های کوچک مثل:
//
// user id
// trace id
// request metadata
// ⚠️ نکته مهم:
//
// نباید برای داده‌های سنگین یا اصلی از WithValue استفاده شود (مثل config بزرگ یا DB handler).
//
// فقط meta-data کوچک.
func ContextWithValue() {
	ctx := context.WithValue(context.Background(), "User1", 23)
	go func(ctx context.Context) {
		fmt.Println("UserId 1 Is:", ctx.Value("User1"))
	}(ctx)
}
