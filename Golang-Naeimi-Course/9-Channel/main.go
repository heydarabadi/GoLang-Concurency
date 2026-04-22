package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	//SimpleUnbufferedChannel1()
	//SimpleUnbufferedChannel2()

	//BufferedChannel()
	//ForLoopInChannel()
	//CheckForCloseChannelExample()
	//DirectionalChannels()
	//SelectStatement()
	//FanOutFanInChannel()
	//PipeLine()
	//InventorySystemForLockWithChannels()

	TimeoutAndCancellation()
}

func SimpleUnbufferedChannel1() {
	ch := make(chan int)

	go func() {
		ch <- 22
	}()

	val := <-ch
	fmt.Println(val)
}

func SimpleUnbufferedChannel2() {
	ch := make(chan int)
	go SendDataToChannel(ch)

	for i := 0; i < 3; i++ {
		val := <-ch
		fmt.Println(val)
	}

	time.Sleep(30 * time.Second)
}

func SendDataToChannel(channel chan int) {
	println("before sending data to channel")
	channel <- 1
	println("after sending data to channel 1")
	channel <- 2
	println("before sending data to channel 2")
	channel <- 3
	println("after sending data to channel 3")
}

func BufferedChannel() {
	ch := make(chan int, 2)

	ch <- 12
	ch <- 13
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

func ForLoopInChannel() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 200; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

func CheckForCloseChannelExample() {
	ch := make(chan int)
	close(ch)

	value, ok := <-ch
	fmt.Println(value, ok)

	if value, ok := <-ch; ok {
		fmt.Println("Channel Is Open %v", value)
	} else {
		fmt.Println("channel Is Close")
	}
}

func DirectionalChannels() {
	// فقط ارسال‌کننده
	sendOnly := func(ch chan<- int, data int) {
		ch <- data
	}

	// فقط دریافت‌کننده
	receiveOnly := func(ch <-chan int) int {
		return <-ch
	}
	ch := make(chan int)
	go sendOnly(ch, 42)
	fmt.Println(receiveOnly(ch))
}

func SelectStatement() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(10 * time.Second)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(20 * time.Second)
		ch2 <- 2
	}()

	for i := 0; i < 2; i++ {
		select {
		case val := <-ch1:
			fmt.Println(val)
		case val2 := <-ch2:
			fmt.Println(val2)

		case <-time.After(40 * time.Second):
			fmt.Println("timeout")
		}
	}

}

func FanOutFanInChannel() {
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	for w := 0; w < 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 0; j < 9; j++ {
		jobs <- j
	}
	close(jobs)

	for r := 0; r < 9; r++ {
		<-results
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker With Id:%d Alredy Work In Job Number: %d \r\n", id, job)
		time.Sleep(1 * time.Second)
		results <- job * 2
	}
}

func PipeLine() {
	gen := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			for _, num := range nums {
				out <- num
			}
			close(out)
		}()
		return out
	}

	multiply := func(in <-chan int, factor int) <-chan int {
		out := make(chan int)
		go func() {
			for num := range in {
				out <- num * factor
			}
			close(out)
		}()
		return out
	}

	numbers := gen(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	multiplied := multiply(numbers, 6)
	for v := range multiplied {
		fmt.Println(v)
	}

}

func RaceConditionWithChannels() {
	counter := 0
	ch := make(chan int, 1)

	ch <- 1

	go func() {
		localCounter := counter
		localCounter++
		counter = localCounter
		fmt.Println(localCounter)
		<-ch
	}()
	time.Sleep(time.Millisecond)
}

func InventorySystemForLockWithChannels() {
	inventory := 100
	orders := []int{20, 30, 25, 35, 10}

	mutex := make(chan int, 1)
	mutex <- 1

	var wg sync.WaitGroup

	for i, order := range orders {
		wg.Add(1)
		go func(orderId int, quantity int) {
			defer wg.Done()
			<-mutex

			if inventory >= quantity {
				time.Sleep(time.Millisecond * 50)
				localInv := inventory
				localInv -= quantity
				inventory = localInv
				fmt.Printf("✅ سفارش %d: %d عدد - تایید شد. موجودی: %d\n",
					orderId, quantity, inventory)
			} else {
				fmt.Printf("❌ سفارش %d: %d عدد - موجودی کافی نیست! موجودی: %d\n",
					orderId, quantity, inventory)
			}
			mutex <- 1
		}(i+1, order)
	}
	wg.Wait()
	fmt.Printf("\n📦 موجودی نهایی انبار: %d\n", inventory)
}

func TimeoutAndCancellation() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Done After Two Seconds"
	}()

	select {
	case res := <-ch:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("TimeOut")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Println("Cancelled:", ctx.Err())
	}
}

////////////////////////

func ProducerConsumerPattern() {
	// کانال بافر شده برای کارها
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	// تولیدکننده: تولید کارها
	go func() {
		for i := 1; i <= 20; i++ {
			fmt.Printf("📦 تولید کار: %d\n", i)
			jobs <- i
			time.Sleep(100 * time.Millisecond) // شبیه‌سازی تولید تدریجی
		}
		close(jobs) // بستن کانال بعد از اتمام تولید
		fmt.Println("✅ تولید تمام شد!")
	}()

	// مصرف‌کننده‌ها: ۳ تا کارگر همزمان
	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go consumer(w, jobs, &wg)
	}

	wg.Wait()
	fmt.Println("🎯 همه کارها انجام شد!")
}

func consumer(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("👷 کارگر %d: شروع کار %d\n", id, job)
		time.Sleep(200 * time.Millisecond) // شبیه‌سازی پردازش
		fmt.Printf("✅ کارگر %d: تمام کرد کار %d\n", id, job)
	}

	fmt.Printf("😴 کارگر %d: کاری نمونده، میخوابم\n", id)
}

///////////////////////

func RateLimiterWithTicker() {
	// محدودیت: 3 درخواست در ثانیه
	rate := 3
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	requests := []string{
		"API-1", "API-2", "API-3", "API-4", "API-5",
		"API-6", "API-7", "API-8", "API-9", "API-10",
	}

	fmt.Println("⏰ Rate Limiter: 3 درخواست در ثانیه")
	fmt.Println("----------------------------------------")

	for i, req := range requests {
		<-ticker.C // منتظر مجوز از ticker

		go func(request string, id int) {
			fmt.Printf("✅ [%d] درخواست %s پردازش شد - %s\n",
				id, request, time.Now().Format("15:04:05"))
		}(req, i+1)
	}

	// صبر برای اتمام درخواست‌ها
	time.Sleep(2 * time.Second)
}

// نسخه پیشرفته: Rate Limiter با burst
func AdvancedRateLimiter() {
	// 5 درخواست در ثانیه با burst 2 تایی
	limiter := time.NewTicker(200 * time.Millisecond) // 5 در ثانیه
	burst := make(chan struct{}, 2)

	// پر کردن burst
	for i := 0; i < 2; i++ {
		burst <- struct{}{}
	}

	go func() {
		for range limiter.C {
			select {
			case burst <- struct{}{}:
			default:
				fmt.Println("⚠️ burst پر شده، درخواست بعدی منتظر میمونه")
			}
		}
	}()

	// شبیه‌سازی درخواست‌ها
	for i := 1; i <= 10; i++ {
		<-burst // منتظر مجوز
		fmt.Printf("✅ درخواست %d پردازش شد - %s\n",
			i, time.Now().Format("15:04:05.000"))
	}

	limiter.Stop()
}

// مثال عملی: محدود کردن تماس‌های API
func APIRateLimiter() {
	type APIRequest struct {
		ID     int
		Method string
		URL    string
	}

	// محدودیت: 10 درخواست در ثانیه
	rateLimiter := time.NewTicker(100 * time.Millisecond)
	defer rateLimiter.Stop()

	requests := []APIRequest{
		{1, "GET", "/users"},
		{2, "POST", "/orders"},
		{3, "GET", "/products"},
		{4, "PUT", "/users/1"},
		{5, "DELETE", "/orders/2"},
	}

	fmt.Println("🌐 API Rate Limiter فعال شد")

	for _, req := range requests {
		<-rateLimiter.C // رعایت نرخ محدودیت

		go func(r APIRequest) {
			// شبیه‌سازی تماس API
			fmt.Printf("🌍 [%d] %s %s - %s\n",
				r.ID, r.Method, r.URL, time.Now().Format("15:04:05"))
		}(req)
	}

	time.Sleep(2 * time.Second)
}

//////////////////

// Message ساختار پیام
type Message struct {
	ID        string
	Topic     string
	Body      string
	Timestamp time.Time
}

// Subscriber مشترک
type Subscriber struct {
	ID       string
	Messages chan Message
	Done     chan struct{}
}

// Broker پیام‌رسان
type Broker struct {
	mu          sync.RWMutex
	subscribers map[string]map[string]*Subscriber // topic -> subscriberID -> subscriber
	topics      map[string]bool
}

// NewBroker ایجاد بروکر جدید
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string]map[string]*Subscriber),
		topics:      make(map[string]bool),
	}
}

// Subscribe اشتراک در یک تاپیک
func (b *Broker) Subscribe(topic string, subscriberID string) *Subscriber {
	b.mu.Lock()
	defer b.mu.Unlock()

	// ایجاد تاپیک اگر وجود نداشت
	if _, exists := b.subscribers[topic]; !exists {
		b.subscribers[topic] = make(map[string]*Subscriber)
		b.topics[topic] = true
	}

	// ایجاد مشترک
	sub := &Subscriber{
		ID:       subscriberID,
		Messages: make(chan Message, 100), // بافر 100 تایی
		Done:     make(chan struct{}),
	}

	b.subscribers[topic][subscriberID] = sub
	fmt.Printf("📝 مشترک %s به تاپیک %s اضافه شد\n", subscriberID, topic)

	return sub
}

// Unsubscribe لغو اشتراک
func (b *Broker) Unsubscribe(topic string, subscriberID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subs, exists := b.subscribers[topic]; exists {
		if sub, ok := subs[subscriberID]; ok {
			close(sub.Messages)
			close(sub.Done)
			delete(subs, subscriberID)
			fmt.Printf("❌ مشترک %s از تاپیک %s حذف شد\n", subscriberID, topic)
		}
	}
}

// Publish انتشار پیام
func (b *Broker) Publish(topic string, message Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	message.Topic = topic
	message.Timestamp = time.Now()

	if subs, exists := b.subscribers[topic]; exists {
		fmt.Printf("📢 انتشار پیام به تاپیک %s: %s\n", topic, message.Body)

		for _, sub := range subs {
			select {
			case sub.Messages <- message:
			default:
				fmt.Printf("⚠️ بافر مشترک %s پر شد!\n", sub.ID)
			}
		}
	} else {
		fmt.Printf("⚠️ تاپیک %s مشترکی ندارد!\n", topic)
	}
}

// Consume مصرف پیام‌ها
func (b *Broker) Consume(sub *Subscriber) {
	fmt.Printf("🎧 مشترک %s شروع به گوش دادن کرد\n", sub.ID)

	for {
		select {
		case msg, ok := <-sub.Messages:
			if !ok {
				fmt.Printf("🔚 مشترک %s قطع شد\n", sub.ID)
				return
			}
			fmt.Printf("📨 [%s] دریافت: %s - زمان: %s\n",
				sub.ID, msg.Body, msg.Timestamp.Format("15:04:05"))

		case <-sub.Done:
			return
		}
	}
}

// مثال استفاده از Message Broker
func MessageBrokerExample() {
	broker := NewBroker()

	// ایجاد تاپیک‌ها
	//topics := []string{"news", "sports", "weather"}

	// مشترک‌ها
	alice := broker.Subscribe("news", "alice")
	bob := broker.Subscribe("news", "bob")
	charlie := broker.Subscribe("sports", "charlie")

	// شروع مصرف پیام‌ها
	go broker.Consume(alice)
	go broker.Consume(bob)
	go broker.Consume(charlie)

	// انتشار پیام‌ها
	messages := []Message{
		{ID: "1", Body: "خبر فوری: Go 1.22 منتشر شد!"},
		{ID: "2", Body: "فوتبال: تیم ملی پیروز شد 3-1"},
		{ID: "3", Body: "هوا: امروز بارانی است"},
		{ID: "4", Body: "خبر: همایش گوگل فردا برگزار می‌شود"},
	}

	for _, msg := range messages {
		var topic string
		switch {
		case msg.Body[:2] == "خبر":
			topic = "news"
		case msg.Body[:4] == "فوتبال":
			topic = "sports"
		case msg.Body[:3] == "هوا":
			topic = "weather"
		default:
			topic = "news"
		}

		broker.Publish(topic, msg)
		time.Sleep(500 * time.Millisecond)
	}

	// منتظر ماندن برای پردازش پیام‌ها
	time.Sleep(2 * time.Second)

	// لغو اشتراک
	broker.Unsubscribe("news", "alice")

	// انتشار پیام بعد از لغو اشتراک
	broker.Publish("news", Message{ID: "5", Body: "خبر: فقط باب این رو میگیره"})

	time.Sleep(1 * time.Second)

	fmt.Println("🎬 پایان برنامه")
}

// نسخه پیشرفته: Broker با الگوی Pub/Sub کامل
func AdvancedBrokerExample() {
	broker := NewBroker()
	var wg sync.WaitGroup

	// ایجاد چندین تاپیک و مشترک
	for i := 1; i <= 3; i++ {
		topic := fmt.Sprintf("topic%d", i)

		for j := 1; j <= 2; j++ {
			subID := fmt.Sprintf("sub%d_%d", i, j)
			sub := broker.Subscribe(topic, subID)

			wg.Add(1)
			go func(s *Subscriber, id string) {
				defer wg.Done()
				broker.Consume(s)
			}(sub, subID)
		}
	}

	// تولیدکننده پیام
	go func() {
		for i := 1; i <= 10; i++ {
			topic := fmt.Sprintf("topic%d", (i%3)+1)
			msg := Message{
				ID:   fmt.Sprintf("%d", i),
				Body: fmt.Sprintf("پیام %d برای %s", i, topic),
			}
			broker.Publish(topic, msg)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("✅ برنامه تمام شد")
}
