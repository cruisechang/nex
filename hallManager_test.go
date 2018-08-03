package nex

import (
	"reflect"
	"testing"

	"github.com/cruisechang/nex/entity"
)

var (
	testHallManager HallManager
)

func TestNewHallManager(t *testing.T) {
	th, _ := NewHallManager()
	testHallManager = th
}

func Test_hallManager_CreateHall(t *testing.T) {
	type args struct {
		id   int
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Hall
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{id: 0, name: "hall"},
			want:    entity.NewHall(0, "hall"),
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{id: 0, name: "hall"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testHallManager.CreateHall(tt.args.id, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("hallManager.CreateHall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hallManager.CreateHall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hallManager_GetHall(t *testing.T) {
	type args struct {
		hallID int
	}
	tests := []struct {
		name  string
		args  args
		want  entity.Hall
		want1 bool
	}{
		{
			name:  "0",
			args:  args{hallID: 0,},
			want:  entity.NewHall(0, "hall"),
			want1: true,
		},
		{
			name:  "1",
			args:  args{hallID: 9,},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := testHallManager.GetHall(tt.args.hallID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hallManager.GetHall() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("hallManager.GetHall() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_hallManager_GetHalls(t *testing.T) {
	tests := []struct {
		name string
		want []entity.Hall
	}{
		{
			name: "0",
			want: []entity.Hall{entity.NewHall(0, "hall")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testHallManager.GetHalls(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hallManager.GetHalls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hallManager_RemoveHall(t *testing.T) {
	type args struct {
		hallID int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{hallID: 0},
		},
		{
			name: "1",
			args: args{hallID: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHallManager.RemoveHall(tt.args.hallID)
		})
	}
}
