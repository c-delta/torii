package server

import (
	"errors"
	"net"

	// "fmt"
	"github.com/c-delta/torii/pack"
)

type SubHandlerManager struct {
	WorkerID []uint32
	Tasks    []*Handler
}

type HandlerManager struct {
	Versions           []uint32
	SubHandlerManagers []*SubHandlerManager
}

func NewHandlerManager() *HandlerManager {
	return &HandlerManager{}
}

func (h *HandlerManager) VersionHas(version uint32) (uint32, bool) {
	for i, v := range h.Versions {
		if v == version {
			return uint32(i), true
		}
	}
	return 0, false
}

func (h *SubHandlerManager) TaskHas(id uint32) (uint32, bool) {
	for i, v := range h.WorkerID {
		if v == id {
			return uint32(i), true
		}
	}
	return 0, false
}

func (h *HandlerManager) Clear() {
	h = &HandlerManager{}
}

func (h *HandlerManager) Add(handler *Handler) error {
	var err error
	version := handler.Version
	if handler.TaskID < 10 {
		return errors.New("this taskid is reserved")
	}
	if i, b := h.VersionHas(version); !b {
		h.Versions = append(h.Versions, version)
		h.SubHandlerManagers = append(h.SubHandlerManagers, &SubHandlerManager{})
		i = uint32(len(h.Versions)) - 1
		err = h.SubHandlerManagers[i].Add(handler)

	} else {
		err = h.SubHandlerManagers[i].Add(handler)
	}
	return err
}

func (h *SubHandlerManager) Add(handler *Handler) error {
	id := handler.TaskID
	if _, b := h.TaskHas(id); !b {
		h.WorkerID = append(h.WorkerID, id)
		h.Tasks = append(h.Tasks, handler)
		return nil
	} else {
		return errors.New("this task is already exists")
	}
}

func (h *HandlerManager) Remove(handler Handler) error {
	var err error
	version := handler.Version
	if handler.TaskID < 10 {
		return errors.New("this taskid is reserved")
	}
	if i, b := h.VersionHas(version); b {
		err = h.SubHandlerManagers[i].Remove(handler)
		if err == nil {
			if len(h.SubHandlerManagers[i].WorkerID) == 0 {
				h.Versions = append(h.Versions[:i], h.Versions[i+1:]...)
				h.SubHandlerManagers = append(h.SubHandlerManagers[:i], h.SubHandlerManagers[i+1:]...)
			}
		}
		return err

	} else {
		return errors.New("this version is not exists")
	}
}

func (h *SubHandlerManager) Remove(handler Handler) error {
	id := handler.TaskID
	if i, b := h.TaskHas(id); b {
		h.Tasks = append(h.Tasks[:i], h.Tasks[i+1:]...)
		h.WorkerID = append(h.WorkerID[:i], h.WorkerID[i+1:]...)
		return nil
	} else {
		return errors.New("this task id is not exists")
	}
}

func (h *HandlerManager) Get(version uint32, taskID uint32) (*Handler, error) {
	if i, b := h.VersionHas(version); b {
		handler, err := h.SubHandlerManagers[i].Get(taskID)
		return handler, err

	} else {
		return nil, errors.New("this version is not exists")
	}
}

func (h *SubHandlerManager) Get(taskID uint32) (*Handler, error) {
	if i, b := h.TaskHas(taskID); b {
		return h.Tasks[i], nil
	} else {
		return nil, errors.New("this task is not exists")
	}
}

func (h *HandlerManager) AcceptTasks(conn net.Conn) {
	defer conn.Close()
	pack := pack.NewPacket()
	for {
		msg, err := pack.Read(conn)
		if err != nil {
			// fmt.Println(err)
			conn.Close()
			return
		}
		handler, err := h.Get(msg.Version, msg.ID)

		if (err != nil) || (msg.ID < 10) {
			conn.Close()
			return
		}
		task := *handler.Handler
		err = task(conn)
		if err != nil {
			break
		}

	}
}
