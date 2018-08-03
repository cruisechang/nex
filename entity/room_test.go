package entity

import (
	"reflect"
	"testing"
	"strconv"
)

var (
	testRoom     Room
	testRoomName string = "name"
	testRoomID   int    = 0
	testRoomType int    = 0
)

func TestNewRoom(t *testing.T) {
	testRoom = NewRoom(testRoomID, testRoomName)
	testRoom.SetIntVariable("roomType",testRoomType)
}

func Test_room_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "0",
			want: testRoomName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := testRoom.Name(); got != tt.want {
				t.Errorf("room.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_room_RoomID(t *testing.T) {

	tests := []struct {
		name string
		want int
	}{
		{
			name: "0",
			want: testRoomID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := testRoom.ID(); got != tt.want {
				t.Errorf("room.RoomID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_room_Active(t *testing.T) {

	tests := []struct {
		name string
		want int
	}{
		{
			name: "0",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := testRoom.Active(); got != tt.want {
				t.Errorf("room.Active() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_room_SetActive(t *testing.T) {

	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{
				0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRoom.SetActive(tt.args.v)
		})
	}
}

func Test_room_SetStringVariable(t *testing.T) {

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{key: "color", value: "white"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRoom.SetStringVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_room_GetStringVariable(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{key: "color"},
			want:    "white",
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{key: "colors"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testRoom.GetStringVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("room.GetStringVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("room.GetStringVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_room_SetFloatVariable(t *testing.T) {

	type args struct {
		key   string
		value float32
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{key: "width", value: 3.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRoom.SetFloatVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_room_GetFloatVariable(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    float32
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{key: "width"},
			want:    3.0,
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{key: "xxxx"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testRoom.GetFloatVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("room.GetFloatVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("room.GetFloatVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_room_SetIntVariable(t *testing.T) {

	type args struct {
		key   string
		value int
	}
	tests := []struct {
		name string
		args args
	}{
		{

			"0",
			args{"cup", 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testRoom.SetIntVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_room_GetIntVariable(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{key: "cup"},
			want:    1,
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{key: "xxxx"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testRoom.GetIntVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("room.GetIntVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("room.GetIntVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_room_SetStructVariable(t *testing.T) {

	v0 := struct {
		Type   string
		weight float32
	}{
		Type:   "apple",
		weight: 0.3,
	}

	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{
				key:   "fruit",
				value: v0,
			},
		},
		{
			name: "1",
			args: args{
				key:   "rice",
				value: v0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testRoom.SetInterfaceVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_room_GetStructVariable(t *testing.T) {

	v0 := struct {
		Type   string
		weight float32
	}{
		Type:   "apple",
		weight: 0.3,
	}

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:"0",
			args:args{key:"fruit"},
			want:v0,
			wantErr:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testRoom.GetInterfaceVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("room.GetStructVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("room.GetStructVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_room_SetIntVariable(b *testing.B){

	tr:= NewRoom(testRoomID, testRoomName)

	for i:=0;i<b.N;i++{
		//tr.SetIntVariable("key"+strconv.Itoa(i),i)
		tr.SetIntVariable("key"+strconv.FormatInt(int64(i),10),i)
	}
}

func Benchmark_room_SetStructVariable(b *testing.B){

	tr := NewRoom(testRoomID, testRoomName)

	v0 := struct {
		Type   string
		weight float32
	}{
		Type:   "apple",
		weight: 0.3,
	}

	for i:=0;i<b.N;i++{
		tr.SetInterfaceVariable("key"+strconv.Itoa(i),v0)
	}

}