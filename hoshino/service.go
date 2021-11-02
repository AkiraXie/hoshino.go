package hoshino

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"sync"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/kv"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// thanks to ZeroBot/examples/manager! but Maybe something wrong...
// TODO : Optimize data persistence
var (
	bucket   = kv.New("Service")
	services = map[string]*Service{}
	mu       = sync.RWMutex{}
)

type Service struct {
	sync.RWMutex
	e                *zero.Engine
	name             string
	states           map[int64]bool
	disableOnDefault bool
}

func NewService(name string, enableOnDefault bool) *Service {
	data, _ := bucket.Get([]byte(name))
	sv := &Service{e: zero.New(), name: name, states: unpack(data), disableOnDefault: !enableOnDefault}
	mu.Lock()
	defer mu.Unlock()
	services[name] = sv
	sv.e.UsePreHandler(sv.check())
	return sv
}

func (sv *Service) Enable(groupID int64) {
	sv.Lock()
	defer sv.Unlock()
	sv.states[groupID] = true
	_ = bucket.Put([]byte(sv.name), pack(sv.states))
}

func (sv *Service) Disable(groupID int64) {
	sv.Lock()
	defer sv.Unlock()
	sv.states[groupID] = false
	_ = bucket.Put([]byte(sv.name), pack(sv.states))
}

func (sv *Service) check() zero.Rule {
	return func(ctx *zero.Ctx) bool {
		sv.RLock()
		if st, ok := sv.states[ctx.Event.GroupID]; ok {
			sv.RUnlock()
			return st
		} //如果states存在则直接返回
		sv.RUnlock()
		if sv.disableOnDefault {
			sv.Disable(ctx.Event.GroupID)
		} else {
			sv.Enable(ctx.Event.GroupID)
		}
		return !sv.disableOnDefault //如果db不存在则先按设置再返回
	}
}

func pack(m map[int64]bool) []byte {
	var (
		buf bytes.Buffer
		b   = make([]byte, 8)
	)
	for k, v := range m {
		binary.LittleEndian.PutUint64(b, uint64(k))
		if v {
			b[7] |= 0x80
		}
		buf.Write(b[:8])
	}
	return buf.Bytes()
}

func unpack(v []byte) map[int64]bool {
	var (
		m      = make(map[int64]bool)
		b      = make([]byte, 8)
		reader = bytes.NewReader(v)
		k      uint64
	)
	for {
		_, err := reader.Read(b)
		if err == io.EOF {
			break
		}
		k = binary.LittleEndian.Uint64(b)
		m[int64(k&0x7fff_ffff_ffff_ffff)] = k&8000_0000_0000_0000 != 0
	}
	return m
}

func (sv *Service) On(typ string, rules ...zero.Rule) *MatcherWrapper {
	return &MatcherWrapper{Matcher: sv.e.On(typ, rules...)}
}

func (sv *Service) OnCommandGroup(commands []string, rules ...zero.Rule) *MatcherWrapper {
	return &MatcherWrapper{Matcher: sv.e.OnCommandGroup(commands, rules...)}
}

func Lookup(service string) (*Service, bool) {
	mu.RLock()
	defer mu.RUnlock()
	m, ok := services[service]
	return m, ok
}

func init() {
	zero.OnCommandGroup([]string{"启用", "enable"}, zero.AdminPermission, zero.OnlyGroup).
		Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.Send("没有找到指定服务!")
			}
			service.Enable(ctx.Event.GroupID)
			ctx.Send(message.Text("已启用服务: " + model.Args))
		})

	zero.OnCommandGroup([]string{"禁用", "disable"}, zero.AdminPermission, zero.OnlyGroup).
		Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.Send("没有找到指定服务!")
			}
			service.Disable(ctx.Event.GroupID)
			ctx.Send(message.Text("已关闭服务: " + model.Args))
		})

	zero.OnCommandGroup([]string{"服务列表", "lssv"}, zero.AdminPermission, zero.OnlyGroup).
		Handle(func(ctx *zero.Ctx) {
			msg := strconv.Itoa(int(ctx.Event.GroupID)) + ` 服务列表:`
			mu.RLock()
			defer mu.RUnlock()
			for k, sv := range services {
				enable := sv.check()(ctx)
				if !enable {
					msg += "\n|X|" + k
				} else {
					msg += "\n|O|" + k
				}
			}
			ctx.Send(message.Text(msg))
		})
}
