package cq

import (
	"encoding/json"
	"net/http"
)

func accept(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]any
	decoder.Decode(&params)
	cqtype := params["post_type"]
	switch cqtype {
	case "message":
		if params["message_type"] == "group" { //群聊消息
			headgro(params)
			return
		}
		if params["message_type"] == "private" { //私聊消息
			headpri(params)
			return
		}
	case "request":
		if params["request_type"] == "friend" {
			//好友请求
			return
		}
		if params["request_type"] == "group" {
			//邀群请求
			return
		}
	case "notice":
		if params["notice_type"] == "group_upload" {
			//群文件上传
			return
		}
		if params["notice_type"] == "group_admin" {
			//群管理员变动
			return
		}
		if params["notice_type"] == "group_decrease" {
			//群成员减少
			return
		}
		if params["notice_type"] == "group_increase" {
			//群成员增加
			return
		}
		if params["notice_type"] == "group_ban" {
			//群禁言
			return
		}
		if params["notice_type"] == "friend_add" {
			//好友添加
			return
		}
		if params["notice_type"] == "group_recall" {
			//群消息撤回
			return
		}
		if params["notice_type"] == "friend_recall" {
			//好友消息撤回
			return
		}
		if params["notice_type"] == "notify" {
			//好友戳一戳
			return
		}
		if params["notice_type"] == "offline_file" {
			//接收到离线文件
			return
		}
		if params["notice_type"] == "client_status" {
			//其他客户端在线状态变更
			return
		}
		if params["notice_type"] == "essence" {
			//精华消息
			return
		}
	}

}

//处理群聊消息
func headgro(params map[string]any) {
	sender := params["sender"].(map[string]any)
	user := newUserByGroup(int64(params["user_id"].(float64)), sender["nickname"].(string),
		sender["sex"].(string), int32(sender["age"].(float64)), sender["card"].(string), sender["area"].(string),
		sender["level"].(string), sender["role"].(string), sender["title"].(string))
	msg := newMsg(int32(params["message_id"].(float64)), params["message_type"].(string), params["sub_type"].(string),
		params["message"].(string), params["raw_message"].(string), int32(params["font"].(float64)), int64(params["time"].(float64)))
	group := newGroup(int64(params["group_id"].(float64)))
	ai := newAi(user, group, msg)
	ai.data = params
	runGroup(ai)

}

//处理私聊消息
func headpri(params map[string]any) {
	sender := params["sender"].(map[string]any)
	user := newUser(int64(params["user_id"].(float64)), sender["nickname"].(string), sender["sex"].(string), int32(sender["age"].(float64)))
	msg := newMsg(int32(params["message_id"].(float64)), params["message_type"].(string), params["sub_type"].(string),
		params["message"].(string), params["raw_message"].(string), int32(params["font"].(float64)), int64(params["time"].(float64)))
	if params["temp_source"] != nil {
		msg.tempSource = int(params["temp_source"].(float64))
	}
	ai := newAi(user, nil, msg)
	ai.data = params
	runPrivate(ai)
}

func runGroup(ai *Ai) {

	for _, k := range roximitors {
		if ai.state == 4 {
			return
		}
		k(ai)
	}
	if ai.state == 4 {
		return
	}

	for k, v := range handleByPower {
		if hasListAdmin(v, ai.User.id) {
			for k2 := range *k {
				if k2.Match(ai.Msg.rawMessage) {
					k2.run(ai)
					afters(ai)
					return
				}
			}
		}
	}
	for k, v := range handlerBySection {
		if hasListAdmin(v.getIds(), ai.User.id) && k.Match(ai.Msg.rawMessage) {
			k.run(ai)
			afters(ai)
			return
		}
	}
	for _, v := range handlerGroups {
		if v.Match(ai.Msg.rawMessage) {
			v.run(ai)
			afters(ai)
			return
		}
	}

	for _, v := range handlers {
		if v.Match(ai.Msg.rawMessage) {
			v.run(ai)
			afters(ai)
			return
		}
	}
	for _, k := range shits {
		if ai.state == 4 {
			return
		}
		k(ai)
	}
}

func runPrivate(ai *Ai) {
	for _, k := range roximitor {
		if ai.state == 4 {
			return
		}
		k(ai)
	}
	if ai.state == 4 {
		return
	}
	for k, v := range handleByPower {
		if hasListAdmin(v, ai.User.id) {
			for k2 := range *k {
				if k2.Match(ai.Msg.rawMessage) {
					k2.run(ai)
					after(ai)
					return
				}
			}
		}
	}
	for k, v := range handlerBySection {
		if hasListAdmin(v.getIds(), ai.User.id) && k.Match(ai.Msg.rawMessage) {
			k.run(ai)
			after(ai)
			return
		}
	}
	for _, v := range handlerPrivates {
		if v.Match(ai.Msg.rawMessage) {
			v.run(ai)
			after(ai)
			return
		}
	}
	for _, v := range handlerTemporarys {
		if v.Match(ai.Msg.rawMessage) {
			v.run(ai)
			after(ai)
			return
		}
	}
	for _, v := range handlers {
		if v.Match(ai.Msg.rawMessage) {
			v.run(ai)
			after(ai)
			return
		}
	}
	for _, k := range shit {
		if ai.state == 4 {
			return
		}
		k(ai)
	}
}
func after(ai *Ai) {
	if ai.state == 4 {
		return
	}
	for _, k := range postprocessor {
		if ai.state == 4 {
			return
		}
		k(ai)
	}
}

func afters(ai *Ai) {
	if ai.state == 4 {
		return
	}
	for _, k := range postprocessors {
		if ai.state == 4 {
			return
		}
		k(ai)
	}
}
