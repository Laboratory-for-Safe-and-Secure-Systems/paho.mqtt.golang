package mqtt

//singelton
import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	asl "github.com/Laboratory-for-Safe-and-Secure-Systems/go-asl"
	asllib "github.com/Laboratory-for-Safe-and-Secure-Systems/go-asl/lib"
)

type mqtt_asl struct {
	endpoint    *asl.ASLEndpoint
	initialized bool
}

type mqtt_asl_con struct {
	conn       net.Conn
	closed     atomic.Bool
	close_once sync.Once
}

func (c *mqtt_asl_con) Read(b []byte) (n int, err error) {
	if c.closed.Load() {
		return 0, net.ErrClosed
	}
	return c.conn.Read(b)
}

func (c *mqtt_asl_con) Write(b []byte) (n int, err error) {
	if c.closed.Load() {
		return 0, net.ErrClosed
	}
	return c.conn.Write(b)
}

// Close decrements the reference counter and only closes resources when reaching 0
func (c *mqtt_asl_con) Close() error {
	var err error
	c.close_once.Do(func() {
		c.conn.SetDeadline(time.Now())
		time.Sleep(time.Millisecond * 10)
		c.closed.Store(true)
		err = c.conn.Close()
	})
	return err
}

func (c *mqtt_asl_con) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *mqtt_asl_con) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *mqtt_asl_con) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *mqtt_asl_con) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *mqtt_asl_con) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

var mqtt_asl_instance mqtt_asl = mqtt_asl{initialized: false}

func Get_custom_function(asl_config asl.EndpointConfig) OpenConnectionFunc {
	if !mqtt_asl_instance.initialized {
		mqtt_asl_instance.initialized = true
		mqtt_asl_instance.endpoint = asl.ASLsetupClientEndpoint(&asl_config)
	}
	return asl_connection_function
}

func asl_connection_function(uri *url.URL, options ClientOptions) (net.Conn, error) {
	fmt.Printf("asl con function invoked")

	asl_con, err := asllib.Dial("tcp", uri.Host, mqtt_asl_instance.endpoint)
	if err != nil {
		return nil, err
	}

	return &mqtt_asl_con{
		conn: asl_con,
	}, nil
}
