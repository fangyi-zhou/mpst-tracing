package globaltype

import (
	"fmt"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/globaltype/protobuf"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
)

func LoadFromProtobuf(fname string) (GlobalType, error) {
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	gtype := &protobuf.GlobalType{}
	if err := proto.Unmarshal(in, gtype); err != nil {
		return nil, fmt.Errorf("failed to load gtype from protobuf: %w", err)
	}
	return convertToGlobalType(gtype)
}

func convertFromAction(action *protobuf.Action) (GlobalType, error) {
	fromRole := action.GetFromRole()
	toRole := action.GetToRole()
	if action.GetType() == protobuf.Action_SEND {
		return &Send{
			origin: fromRole,
			dest:   toRole,
			conts:  nil,
		}, nil
	} else if action.GetType() == protobuf.Action_RECV {
		if len(action.GetContinuations()) != 1 {
			return nil, fmt.Errorf("current does not support multiple receives")
		}
		return &Recv{
			origin: fromRole,
			dest:   toRole,
			label:  action.GetContinuations()[0].GetLabel(),
			cont:   nil,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid action type %d", action.GetType())
	}
}

func patchIdx(
	action *protobuf.Action,
	gtype GlobalType,
	actions map[int32]GlobalType,
) (GlobalType, error) {
	if action.GetType() == protobuf.Action_SEND {
		gtype := gtype.(*Send)
		gtype.conts = make(map[string]GlobalType)
		for _, cont := range action.GetContinuations() {
			label := cont.GetLabel()
			if cont.GetNext() == -1 {
				gtype.conts[label] = Done{}
			} else {
				var exists bool
				gtype.conts[label], exists = actions[cont.GetNext()]
				if !exists {
					return nil, fmt.Errorf("non-existent index: %d", cont.GetNext())
				}
			}
		}
		return gtype, nil
	} else if action.GetType() == protobuf.Action_RECV {
		if len(action.GetContinuations()) != 1 {
			return nil, fmt.Errorf("current does not support multiple receives")
		}
		gtype := gtype.(*Recv)
		cont := action.GetContinuations()[0]
		var exists bool
		gtype.cont, exists = actions[cont.GetNext()]
		if !exists {
			return nil, fmt.Errorf("non-existent index: %d", cont.GetNext())
		}
		return gtype, nil
	} else {
		return nil, fmt.Errorf("invalid action type %d", action.GetType())
	}
}

func convertToGlobalType(globalType *protobuf.GlobalType) (GlobalType, error) {
	actions := make(map[int32]GlobalType)
	for _, prefix := range globalType.Actions {
		gtype, err := convertFromAction(prefix)
		if err != nil {
			return nil, err
		}
		if _, exists := actions[prefix.GetIdx()]; exists {
			return nil, fmt.Errorf("duplicate index: %d", prefix.GetIdx())
		}
		actions[prefix.GetIdx()] = gtype
	}
	for _, prefix := range globalType.Actions {
		gtype := actions[prefix.GetIdx()]
		gtype, err := patchIdx(prefix, gtype, actions)
		if err != nil {
			return nil, err
		}
	}
	start, exists := actions[globalType.GetStart()]
	if !exists {
		return nil, fmt.Errorf("non existsent index for start: %d", globalType.GetStart())
	}
	return start, nil
}
