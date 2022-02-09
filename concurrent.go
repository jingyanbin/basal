package basal

//并发控制器
type ConcurrentController struct {
	ch chan struct{}
}

//获取
func (m *ConcurrentController) Acquire(block bool) bool {
	if block {
		m.ch <- struct{}{}
		return true
	} else {
		select {
		case m.ch <- struct{}{}:
			return true
		default:
			return false
		}
	}
}

//释放
func (m *ConcurrentController) Release(block bool) {
	if block {
		<-m.ch
	} else {
		select {
		case <-m.ch:
		default:
		}
	}
}

func NewConcurrentController(max int) *ConcurrentController {
	return &ConcurrentController{ch: make(chan struct{}, max)}
}
