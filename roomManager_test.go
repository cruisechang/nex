package nex

import (
	"reflect"
	"testing"

	"github.com/cruisechang/nex/entity"
)

var (
	testRoomManager RoomManager
)

func TestNewRoomManager(t *testing.T) {
	th, _ := NewRoomManager()
	testRoomManager = th
}

func Test_RoomManager_CreateRoom(t *testing.T) {

	type args struct {
		id   int
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Room
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{id: 2001, name: "龍虎1"},
			want:    entity.NewRoom(2001, "龍虎1"),
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{id: 2002, name: "龍虎2"},
			want:    entity.NewRoom(2002, "龍虎2"),
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{id: 2001, name: "龍虎1"},  //id duplicate
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testRoomManager.CreateRoom(tt.args.id, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("roomManager.CreateRoom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roomManager.CreateRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RoomManager_GetRoom(t *testing.T) {
	type args struct {
		roomID int
	}
	tests := []struct {
		name  string
		args  args
		want  entity.Room
		want1 bool
	}{
		{
			name:  "0",
			args:  args{roomID: 2001,},
			want:  entity.NewRoom(2001, "龍虎1"),
			want1: true,
		},
		{
			name:  "1",
			args:  args{roomID: 9},      //fake id
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := testRoomManager.GetRoom(tt.args.roomID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roomManager.GetRoom() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("roomManager.GetRoom() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_roomManager_GetRooms(t *testing.T) {
	tests := []struct {
		name string
		want []entity.Room
	}{
		{
			name: "0",
			want: []entity.Room{entity.NewRoom(2001, "龍虎1"),entity.NewRoom(2002, "龍虎2")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testRoomManager.GetRooms(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roomManager.GetRooms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roomManager_RemoveRoom(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{id: 0},
		},
		{
			name: "1",
			args: args{id: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRoomManager.RemoveRoom(tt.args.id)
		})
	}
}
