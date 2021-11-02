package hoshino

import (
	"sync"
	"time"
	"unsafe"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type MatcherWrapper struct {
	sync.RWMutex
	*zero.Matcher
	sv        *Service
	handlers  []zero.Handler
	argParser zero.Handler
}

func (mw *MatcherWrapper) SetPriority(Pri int) *MatcherWrapper {
	mw.Priority = Pri
	return mw
}
func (mw *MatcherWrapper) SetBlock(block bool) *MatcherWrapper {
	mw.Matcher.Block = block
	return mw
}
func (mw *MatcherWrapper) Handle(handler zero.Handler) *MatcherWrapper {
	mw.Lock()
	defer mw.Unlock()
	mw.handlers = append(mw.handlers, handler)
	return mw
}

func (mw *MatcherWrapper) SetArgParser(parser zero.Handler) *MatcherWrapper {
	mw.argParser = parser
	return mw
}

func (mw *MatcherWrapper) Got(key, prompt string, handler zero.Handler) *MatcherWrapper {
	parse := func(ctx *zero.Ctx) {
		if _, ok := ctx.State[key]; !ok {
			if prompt != "" {
				ctx.Send(prompt)
			}
			select {
			case event := <-zero.NewFutureEvent("message", -1, true, ctx.CheckSession()).Next():
				if mw.argParser == nil {
					ctx.State[key] = event.RawMessage
					return
				}
				pc := (*zero.APICaller)(unsafe.Pointer(uintptr(unsafe.Pointer(ctx)) + unsafe.Sizeof(ctx.GetMatcher()) + unsafe.Sizeof(ctx.Event) + unsafe.Sizeof(ctx.State)))
				ctx1 := &zero.Ctx{}
				p := unsafe.Pointer(ctx1)
				pma := (**zero.Matcher)(p)
				*pma = ctx.GetMatcher()
				ctx1.Event = event
				ctx1.State = ctx.State
				pcaller := (*zero.APICaller)(unsafe.Pointer(uintptr(p) + unsafe.Sizeof(ctx.GetMatcher()) + unsafe.Sizeof(ctx1.Event) + unsafe.Sizeof(ctx1.State)))
				*pcaller = *pc
				mw.argParser(ctx1)
				for k, vv := range ctx1.State {
					ctx.State[k] = vv
				}
				return
			case <-time.After(time.Duration(HsoConfig.ExpireDuration) * time.Second):
				return
			}
		} else {
			return
		}
	}
	mw.Lock()
	defer mw.Unlock()
	mw.handlers = append(mw.handlers, parse, handler)
	return mw

}

func (mw *MatcherWrapper) Finish(ctx *zero.Ctx, message interface{}) {
	if message != nil {
		ctx.Send(message)
	}
	ctx.State["__done__"] = 1
}

func (mw *MatcherWrapper) FinishChain(ctx *zero.Ctx, message ...message.MessageSegment) {
	if message != nil {
		ctx.SendChain(message...)
	}
	ctx.State["__done__"] = 1
}

func (mw *MatcherWrapper) Ok() *MatcherWrapper {
	mw.Matcher.Handle(func(ctx *zero.Ctx) {
		log.Infof("消息 %v将会被 %v服务处理", ctx.Event.RawMessage, mw.sv.name)
		for _, handler := range mw.handlers {
			handler(ctx)
			if _, ok := ctx.State["__done__"]; ok {
				break
			}
		}
		log.Infof("消息 %v处理完毕", ctx.Event.RawMessage)
	})
	return mw
}
