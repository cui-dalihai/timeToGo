package memo

//type Memo struct {
//	f Func2
//	cache map[string]result
//	mu sync.Mutex
//}

// New方法仅仅是一个普通的可导出方法，接受一个Func类型的参数f，返回一个Memo类型的结果  todo 是否可以返回值本身而不是值的指针呢？
//func New(f Func2) *Memo {
//	return &Memo{f: f, cache: make(map[string]result)}
//}

// 这个Get方法是一个绑定在Memo struct上的方法
//func (memo *Memo) Get1(key string) (interface{}, error) {
//
//	// 对于每个key都是串行执行， 相同的key第二次请求开始使用cache
//	res, ok := memo.cache[key]
//	if !ok {
//		res.value, res.err = memo.f(key)   // 这步耗时操作
//		memo.cache[key] = res
//	}
//	return res.value, res.err
//}
//
//func (memo *Memo) Get2(key string) (value interface{}, err error) {
//	memo.mu.Lock()
//	res, ok := memo.cache[key]
//	if !ok {
//		res.value, res.err = memo.f(key)  		// 把耗时操作放在了锁里面, 导致其它goroutine在当前goroutine未释放锁之前一直在等待
//		memo.cache[key] = res             		// 虽然为每个url建立了goroutine, 但是每个goroutine都需要排队领取锁,
//	}                                     		// 而每个goroutine占用锁的时间又是从检测缓存到建立缓存的全部时间
//	memo.mu.Unlock()                      		// 所以, 从使用并发优化请求的角度来说, 所有goroutine效果等同于一个goroutine
//	return res.value, res.err             		// 并发执行的仅仅是goroutine领取key的过程
//}
//
//func (memo *Memo) Get3(key string) (value interface{}, err error) {
//	memo.mu.Lock()
//	res, ok := memo.cache[key]            		// 每个goroutine依次读取缓存, 这个主要是为了显式的内存同步,
//	memo.mu.Unlock()							// 防止其它goroutine正在写key的中间被读取, 比如别的goroutine写key分为两步, 创建key和空值, 写入有效值, 不同步可能会读到前面的空值而不是第二步的有效值
//												// 防止其它cpu中的goroutine已经写了key, 但未刷到内存, 导致当前goroutine看不到这个key
//
//	if !ok {
//		res.value, res.err = memo.f(key)		// 同样, 由于f是耗时操作, 前一个拿到锁并检查key不在cache中的goroutine会在这里等待,
//												// 这时这个key的响应还未回来, 所以缓存没有建立, 所以后一个goroutine取到锁之后检测还是为空
//												// 最终两个goroutine会为同一个key各自发起请求
//		memo.mu.Lock()
//		memo.cache[key] = res                   // 写的时候虽然使用了锁, 但不能避免两个相同key的goroutine各自请求之后写入相同的key
//		memo.mu.Unlock()						// 这里不会被检测出竞态
//	}
//
//	return res.value, res.err
//}

//type Memo struct {
//	f Func2
//	mu sync.Mutex
//	cache map[string]*entry
//}
//
//func New(f Func2) *Memo {
//	return &Memo{f: f, cache: make(map[string]*entry)}
//}

// 每个goroutine在检查cache key和写入key-entry使用锁来同步, 这都是在内存中完成, 而f都是并发请求的
//func (memo *Memo) Get4(key string) (value interface{}, err error) {
//	memo.mu.Lock()
//	e := memo.cache[key]
//	if e == nil {
//		e = &entry{ready: make(chan struct{})}
//		memo.cache[key] = e
//		memo.mu.Unlock()
//
//		e.res.value, e.res.err = memo.f(key)
//		close(e.ready)  		// close会通知读e.ready通道的goroutine, 该通道已经结束, 那么那个goroutine会在读位置继续向下执行
//	} else {
//		memo.mu.Unlock()
//		<- e.ready  			// 从e.ready中取值并丢弃, 仅有close向其中发送值, 这个相当于当前goroutine一直阻塞等待e.ready调用close
//	}
//	return e.res.value, e.res.err
//}

type Func func(key string) (interface{}, error)
type result struct {
	value interface{}
	err   error
}
type entry struct {
	res   result
	ready chan struct{}
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	close(e.ready)
}
func (e *entry) deliver(response chan<- result) {
	// 对于耗时请求f, 相同的url时会出现第一个写入cache并发起请求, 而第二个需要在deliver中监听e.ready
	// 只有当第一个url的call关闭e.ready, 第二个才能收到结果
	<-e.ready
	response <- e.res
}

type request struct {
	key      string
	response chan<- result // 请求结果写入通道
}
type Memo struct {
	requests chan request // request 类型通道
}

func (memo *Memo) Close() { close(memo.requests) }

// monitor goroutine, 守护cache并且非阻塞的处理请求,
// 请求结果的call和等待结果的deliver都需要独立的goroutine, 以防止阻塞Server的主goroutine
func (memo *Memo) Server(f Func) {
	cache := make(map[string]*entry)

	// 因为仅创建了一个memo, 所以只有一个memo.requests通道
	// 所以每收到一个request, 对它处理的过程都是线性的
	for req := range memo.requests {

		// 读写cache也都是在monitor goroutine内, 所以不用加锁
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e

			// 为每个url对应的entry单独创建一个goroutine, 用来请求这个url, 写entry, 关闭entry
			go e.call(f, req.key)
		}

		// 再为每个e创建一个监听ready的通道, 一旦e.ready, 就把e.res发到这个req的响应通道
		// 这样每个url通过Get创建的监听response的通道就能收到对应的响应
		go e.deliver(req.response)
	}
}

// 为每个url创建一个goroutine, 每个goroutine来调用memo的Get方法
func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)

	// 为每个key创建一个request, 包含key和一个通道
	memo.requests <- request{key, response}

	// 然后就立即监听这个request的response通道
	// 所以每个url的goroutine都会在这个response监听的位置阻塞等待
	// 当deliver结束后每个url对应的Get创建的goroutine才会收到res, 结束阻塞
	res := <-response
	return res.value, res.err
}

// 为httpGetBody创建一个Memo, 并启动memo.requests通道监听
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.Server(f)
	return memo
}
