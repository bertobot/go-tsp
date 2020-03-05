
package tsp

import (
    "net"
    "fmt"
    "bufio"
    "strconv"
    "strings"
    "time"
)

type Client struct {
    host string
    port int
    conn net.Conn
}

type KV struct {
    key string
    value float64
}

func NewKV(key string, value float64) KV {
    kv := KV{key, value}
    return kv
}

type TV struct {
    timestamp int
    value float64
}

func NewTV(millis int, value float64) TV {
    tv := TV{millis, value}
    return tv
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func (t *Client) Connect() {
    //strport := strconv.Itoa(t.port)
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.host, t.port))
    check(err)

    t.conn = conn
}

func (t *Client) Close() {
    t.conn.Close()
}

func (t *Client) Put(key string, timestamp int, value float64) bool {
    msg := fmt.Sprintf("put %s %d %f\n", key, timestamp, value)
    t.conn.Write([]byte(msg))
    result, err := bufio.NewReader(t.conn).ReadString('\n')
    check(err)
    return result == "ok"
}

func (t *Client) PutMultipleKV(timestamp int, pairs []KV) bool {
// put multiple key-values on a single timestamp
    var pairstr string
    for _,kv := range pairs {
        pairstr += fmt.Sprintf(" %s %f", kv.key, kv.value)
    }

    msg := fmt.Sprintf("mkput %d%s\n", timestamp, pairstr)
    t.conn.Write([]byte(msg))

    for range pairs {
        result, err := bufio.NewReader(t.conn).ReadString('\n')
        check(err)
        if result != "ok" {
            return false
        }
    }
    return true
}

func (t *Client) PutMultipleTV(key string, pairs []TV) bool {
// put multiple timestamp-values on a single key
    var pairstr string
    for _,tv := range pairs {
        pairstr += fmt.Sprintf(" %d %f", tv.timestamp, tv.value)
    }

    msg := fmt.Sprintf("mtput %s%s\n", key, pairstr)
    t.conn.Write([]byte(msg))

    for range pairs {
        result, err := bufio.NewReader(t.conn).ReadString('\n')
        check(err)
        if result != "ok" {
            return false
        }
    }
    return true
}

func (t *Client) Get(key string, start int, end int) []TV {
    var result []TV

    msg := fmt.Sprintf("get %s %d %d\n", key, start, end)
    t.conn.Write([]byte(msg))

    scanner := bufio.NewScanner(t.conn)

    for scanner.Scan() {
        if strings.HasPrefix(scanner.Text(), "ok") {
            break
        }

        tokens := strings.SplitN(scanner.Text(), ",", -1)
        t, err := strconv.Atoi(tokens[0])
        check(err)

        v, err := strconv.ParseFloat(tokens[1], 64)
        check(err)

        result = append(result, TV{t,v})
    }

    return result
}

func NewClient(hostname string, port int) Client {
    c := Client{hostname, port, nil}
    return c
}

type Poller struct {
    client Client
    period float64
    onLate func()
    onComplete func(elapsed time.Duration) 
    polls []func() ([]KV, int, string)
}

func (p *Poller) Add(f func() ([]KV, int, string)) {
    p.polls = append(p.polls, f)
}

func (p *Poller) Run() {
    for {
        // TODO: time each function and publish as meta?
        start := time.Now()

        for _, f := range p.polls {
            kv, timestamp, namespace := f()

            var ourkv []KV
            for _, okv := range kv {
                ourkv = append(ourkv, KV{fmt.Sprintf("%s#%s", namespace, okv.key), okv.value})
            }

            p.client.PutMultipleKV(timestamp, ourkv)
        }

        elapsed := time.Now().Sub(start)

        if elapsed.Seconds() >= p.period {
            if p.onLate != nil {
                p.onLate()
            }
        } else {
            if p.onComplete != nil {
                p.onComplete(elapsed)
            } else {
                time.Sleep(time.Duration(p.period - elapsed.Seconds()))
            }
        }
    }
}

func NewPoller(client Client, period float64, late func(), complete func(elapsed time.Duration)) Poller {
    var empty []func()([]KV, int, string)

    poller := Poller{client, period, late, complete, empty}
    return poller
}

