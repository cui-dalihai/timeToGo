主要是为了解释什么情况下一个goroutine写的变量值能被另一个goroutine**可靠观察到**

无论是用通道操作来守护数据实现、还是使用sync和sync/atomic中的同步原语来实现，程序中多个goroutine同时访问相同数据时这些访问一定是串行的。

在单个goroutine中, 读写的真正执行顺序必须要和代码中指定的顺序具有相同的执行效果, 这句话的意思是, 解释器和CPU可能会对程序中单个goroutine内的一些读写操作进行重新排序, 但调整顺序前后的执行结果不能跟程序中指定的顺序执行结果不一致。由于这种对执行顺序的调整，一个goroutine中的执行顺序和其它goroutine观察到实际执行顺序可能会不同，比如一个goroutine执行了a=1; b=2,另一个goroutine可能会观察到b先被复制为2，然后再是a=1;

**Happens Before:** 

为了说清楚读和写的需求，先定义一下这个happens before: 在...之前发生, 当e1在e2之前发生时, 就是在说e2在e1之后发生，当e1 not happens before e2, 且e2 not happens before e1时, 我们说e1和e2这时是并发的,

在单个goroutine内, 在...发生之前这样的顺序是由代码表达式决定的

对于变量v的一个读需求r, 如果可以观察到写需求w对v的操作, 那么r和w要满足:
1. r不能发生在w之前
2. 在w之后且在r之前没有其他的对v的写

而如果为了保证对v的读r能够观察到对v指定的一次写w，就是说要r仅观察到这一次w, 为了实现r能够可靠观察到这次w，两者要满足:
1. w发生在r之前
2. 其它任何对共享变量v的写，要么发生在w之前，要么发生在r之后

下面这对约束要强于上面那对, 因为下面这对明确要求在w和r时没有其它的w并发发生。在单个goroutine内是不可能并发的，所以单个goroutine的情况下上面两对是一个意思：对v的读能够获取最近一次的w。但是在多个goroutine共享v的情况下, 就必须使用同步原语建立可靠的happens-before来保证一次读能够取到指定的一次写。

使用v类型的零值对v进行初始化的行为和一次对v的写操作，在内存模型中是一样的

对于一个大于一个机器字的值来说，对它的读和写和多机器字大小的操作一样，都是不确定的顺序

**同步中的happens before:**

几种可靠的发生顺序

1. 如果p导入q包, 那么q的init函数是可靠发生在p中任何逻辑之前的
2. 而main包中main函数是可靠发生在所有init函数之后的
3. goroutine创建时的go声明可靠发生在这个goroutine开始执行之前
```
var a string
func f() {
        print(a)
}
func hello() {
        a = "hello, world"   # a是被先赋值, go f()后执行, 所以print(a)一定会打印"hello,world"
        go f()
}
```
4. 不实用同步机制的话, 无法可靠保证goroutine的退出相对于程序中任何事件的先或者后
```
var a string
func hello() {
        go func() { a = "hello" }()
        print(a)   # 这可以打印空字符串, 也可以打印hello, 甚至一些激进的编译器直接删除前面的go
}
```
**通道通讯中的happens before:**

通道通讯是主要的goroutine之间的同步机制, 每个通道有对应的发送方和接收方, 通常发送和接收会在不同的goroutine

5. 一次发送可靠发生在对应这次发送的接收完成之前
```
var c = make(chan int, 10)
var a string
func f() {
        a = "hello, world"
        c <- 0
}
func main() {
        go f()
        <-c
        print(a)  # 这个一定可靠打印hello, world, 因为main中<-c接收完成之前,c<-0一定可靠发生, 那么对a的写一定也可靠发生 
}
```
6. 通道关闭可靠发生在接受方收到通道类型的零值之前, 上面c<-0改为close(c)是相同的效果
7. 对于无缓冲通道的接收是可靠发生在发送完成之前
```
var c = make(chan int)
var a string
func f() {
        a = "hello, world"
        <-c
}
func main() {
        go f()
        c <- 0
        print(a)  # c的发送发生在print之前, 
}
```
8. 第k次对容量为C的缓冲通道的接收是可靠发生在第k+C次的发送完成之前

注: 这个位置需要对比5, 7, 8理解一下, 
文档原文如下:
1. A send on a channel **happens before** the corresponding receive from that channel **completes**.
2. The closing of a channel **happens before** a receive that returns a zero value because the channel is closed.
3. A receive from an unbuffered channel **happens before** the send on that channel **completes**. 
4. The kth receive on a channel with capacity C **happens before** the k+Cth send from that channel **completes**.

前两句比较好理解, 重点是3,4两句对于非缓冲通道和缓冲通道满了情况的描述比较令人费解, 另[一篇介绍通道的文档](https://golang.org/doc/effective_go.html#channels)中有这一段
>If the channel is unbuffered, the sender blocks until the receiver has received the value. If the channel has a buffer, the sender blocks only until the value has been copied to the buffer; if the buffer is full, this means waiting until some receiver has retrieved a value.

>如果是无缓冲通道, 发送者会一直阻塞到接收者接收完成这个值. 如果是缓冲通道, 发送者会一直阻塞直到值被复制到缓冲区, 如果缓冲区满了, 那就要等接收者从缓冲区中取走一个值。

这段介绍和3,4的结论是一致的, 即对于阻塞状态下的通道, 无论是无缓冲通道还是缓冲通道满了, 接收完成一定是先于发送完成的, 这里一直使用的是has received和has retrieved, 对应3,4中的completes, 所以发送这个行为或许是先发生的, 但最终完成, 一定是接收先完成, 然后发送才完成. 

另外, 这段话还提供了缓冲通道的细节: 把发送者等待的是把值复制到缓冲区, 而不是接收者完成, 接收者等待的是缓冲区的值, 所以对于缓冲未满的情况, 发送者要先完成把值复制到缓冲区, 接收者才能从缓冲区读到值, 就是1的结论, 而非缓冲通道发送者等待的是接收者完成.(这细节有卵用, 可能是知道从阻塞状态下通道解阻塞后, 接收者先走一步，但两者处于不同goroutine, 后续各自的代码执行先后还是未知的😜)

两种通道时序图简单画一下吧

<center><img src="https://github.com/cui-dalihai/timeToGo/blob/master/concurentNonBlockingCache/channel-sender-receiver.png" width="100%" height="100%"></center>

ok, 接着读这篇内存模型的文档

通过第八条结论, 可以用缓冲通道来模拟计数型的同步机制: 缓冲数代表最大允许的活跃同步量的数量, 达到数量之后, 如果还想使用同步量就要等待其它活跃的同步量被释放, 常用来限制并发, 上代码：
```
var limit = make(chan int, 3)
func main() {
    for _, w := range work {    # 虽然for为每个work创建了一个goroutine, 但这些goroutine并不是同时活跃的   
        go func(w func()) {  
            limit <- 1          # limit满了情况下, goroutine就会阻塞在这里
            w()
            <- limit            # 直到其它goroutine执行完w(), 从limit中取一个值出来, 达到限制任何时候最大活跃goroutine只有3
        }(w)
    }
}
```

**锁中的happens before:**

9. 对于sync.Mutex或者sync.RWMutex类型变量l(小写L), n, m其中n<m, n次对l.Unlock()可靠发生在m次的l.Lock()之前
```
var l sync.Mutex
var a string
func f() {
    a = "hello"
    l.Unlock()   # n = 1
}
func main() {
    l.lock()  # m = 1 
    go f()
    l.lock()  # m = 2 上面n=1可靠发生在m=2之前, 所以对a的写发生在m=2之前, m=2发生在print之前, 所以对a的写发生在print之前, 可靠打印hello
    print(a)
}
```
10. For any call to l.RLock on a sync.RWMutex variable l, there is an n such that the l.RLock happens (returns) after call n to l.Unlock and the matching l.RUnlock happens before call n+1 to l.Lock.这句意思是下图
<center><img src="https://github.com/cui-dalihai/timeToGo/blob/master/concurentNonBlockingCache/RWMutex.png" width="30%" height="30%"></center>



**Once中的happens before:**
11. Once提供了并发场景下的初始化方案, 多个goroutine调用once.Do(f), 仅会有一个真正执行了f( ), 其它的goroutine会阻塞等待执行的那个返回, 即其中一个真正执行的那个goroutine执行f( )会发生在任何一once.Do(f)返回之前
```
var a string
var once sync.Once
func setup() {
    a = "hello"
}
func doprint() {
    once.Do(setup)
    print(a)
}
func twoprint() {
    go doprint()     # 这两个goroutine中仅有一个真正执行了setup()，但是两个都会阻塞到setup()被执行完成
    go doprint()     # 所以a写入发生在once.Do(setup)之前，print(a)会可靠打印两遍hello
}
```

**不正确的同步:**
```
var a, b int
func f() {
    a = 1
    b = 2
}
func g() {
    print(b)
    print(a)
}
func main() {
    go f()
    g()    # 这个位置几乎可print任何组合, 0-0, 0-1, 2-0, 1-2, 因为f的goroutine和主goroutine没有任何同步，
}
```
```
var a string
var done bool

func setup() {
    a = "hello"
    done = true
}
func doprint() {
    if !done {          # 重点是, 这个逻辑是在暗示读到了done就能读到在done之前写的a, 实际上是，在没有同步机制下， 读到了done也不一定
        once.Do(setup)  # 能读到a      
    }
    print(a)
}
func twoprint() {
    go doprint()    # 可能两个goroutine都会阻塞在once.Do(setup)位置, 其中一个真正执行了setup, 而另一个不会执行, 这个为执行的goroutine
    go doprint()    # 就无法可靠观察到那个执行setup的goroutine对a的写, 所以会有一个空字符串
}
```
```
var a string
var done bool

func setup() {
    a = "hello"
    done = true
}
func main() {
    go setup()
    for !done {}    # 这个也是在暗示读到done就能读到a,同样这个done可能被main goroutine读到, 但不一定表示就能读到a, 还有就是这个done
    print(a)        # 也有可能永远不会被main读到,
}
```
```
type T struct {
    msg string
}
var g *T
func setup() {
    t := new(T)
    t.msg = "hello"
    g = t
}
func main() {
    go setup()
    for g == nil {}    # main gorotine和setup gorotine共享了g, 所以main可以观察到g, 但是对g.msg的写无法可靠保证。
    print(g.msg)
}
```
只要显式使用同步原语就可以解决上面的问题
