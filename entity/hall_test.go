package entity

import (
	"reflect"
	"testing"
)

var (
	testHall Hall
)

func TestNewHall(t *testing.T) {

	type args struct {
		hallID int
		name   string
	}
	tests := []struct {
		name string
		args args
		want Hall
	}{
		{
			name: "0",
			args: args{hallID: 0, name: "name"},
			want: NewHall(0, "name"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHall(tt.args.hallID, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHall() = %v, want %v", got, tt.want)
			}
		})
	}

	testHall = NewHall(0, "name")
}

func Test_hall_HallID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "0",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testHall.ID(); got != tt.want {
				t.Errorf("hall.HallID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hall_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "0",
			want: "name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testHall.Name(); got != tt.want {
				t.Errorf("hall.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hall_SetName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{name: "nameSet"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHall.SetName(tt.args.name)

			if got := testHall.Name(); got != tt.args.name {
				t.Errorf("hall.SetName() = %v, want %v", got, tt.args.name)
			}
		})
	}
}


func Test_hall_Active(t *testing.T) {

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

			if got := testHall.Active(); got != tt.want {
				t.Errorf("hall.Active() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hall_SetActive(t *testing.T) {

	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{v: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHall.SetActive(tt.args.v)
			if got := testHall.Active(); got != tt.args.v {
				t.Errorf("hall.SetActive() = %v, want %v", got, tt.args.v)
			}
		})
	}
}

func Test_hall_SetStringVariable(t *testing.T) {

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
			args: args{
				key:   "nickname",
				value: "nickname",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHall.SetStringVariable(tt.args.key, tt.args.value)

		})
	}
}

func Test_hall_GetStringVariable(t *testing.T) {

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
			name: "0",
			args: args{
				key: "nickname",
			},
			want:    "nickname",
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				key: "xxxx",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testHall.GetStringVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("hall.GetStringVariable() error = %v, wantErr %v, name=%s", err, tt.wantErr, tt.name)
				return
			}
			if got != tt.want {
				t.Errorf("hall.GetStringVariable() = %v, want %v, name=%s", got, tt.want, tt.name)
			}
		})
	}
}

func Test_hall_SetFloatVariable(t *testing.T) {

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
			args: args{
				key:   "age",
				value: 1.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHall.SetFloatVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_hall_GetFloatVariable(t *testing.T) {

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
			name: "0",
			args: args{
				key: "age",
			},
			want:    1.0,
			wantErr: false,
		},
		{
			name: "1",
			args: args{
				key: "ages",
			},
			want:    0,
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testHall.GetFloatVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("hall.GetFloatVariable() error = %v, wantErr %v, name=%s", err, tt.wantErr, tt.name)
				return
			}
			if got != tt.want {
				t.Errorf("hall.GetFloatVariable() = %v, want %v, name=%s", got, tt.want, tt.name)
			}
		})
	}
}

func Test_hall_SetIntVariable(t *testing.T) {

	type args struct {
		key   string
		value int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",

			args: args{
				key:   "height",
				value: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testHall.SetIntVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_hall_GetIntVariable(t *testing.T) {

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
			name: "0",
			args: args{
				key: "height",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "1",
			args: args{
				key: "heightxxxx",
			},
			want:    0,
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testHall.GetIntVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("hall.GetIntVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hall.GetIntVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hall_SetStructVariable(t *testing.T) {

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

			testHall.SetInterfaceVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_hall_GetStructVariable(t *testing.T) {

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

			got, err := testHall.GetInterfaceVariable(tt.args.key)
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
