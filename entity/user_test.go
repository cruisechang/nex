package entity

import (
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		userID int
		connID string
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		{
			name: "user0",
			args:args{userID:0,connID:"connID"},
			want:NewUser(0,"connID"),
		},
		{
			name: "user1",
			args:args{userID:0,connID:""},
			want:NewUser(0,""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUser(tt.args.userID,tt.args.connID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_ConnID(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			if got := u.ConnID(); got != tt.want {
				t.Errorf("user.ConnID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_GameID(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			if got := u.GameID(); got != tt.want {
				t.Errorf("user.GameID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_SetGameID(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetGameID(tt.args.v)
		})
	}
}

func Test_user_Name(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			if got := u.Name(); got != tt.want {
				t.Errorf("user.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_SetName(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetName(tt.args.name)
		})
	}
}

func Test_user_SessionID(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			if got := u.SessionID(); got != tt.want {
				t.Errorf("user.SessionID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_SetSessionID(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetSessionID(tt.args.v)
		})
	}
}

func Test_user_IsTest(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			if got := u.IsTest(); got != tt.want {
				t.Errorf("user.IsTest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_SetIsTest(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		v int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetIsTest(tt.args.v)
		})
	}
}

func Test_user_Credit(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			if got := u.Credit(); got != tt.want {
				t.Errorf("user.Credit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_SetCredit(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		v float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetCredit(tt.args.v)
		})
	}
}

func Test_user_AccessID(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			if got := u.AccessID(); got != tt.want {
				t.Errorf("user.AccessID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_SetAccessID(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetAccessID(tt.args.v)
		})
	}
}

func Test_user_SetStringVariable(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetStringVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_user_GetStringVariable(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			got, err := u.GetStringVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.GetStringVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("user.GetStringVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_SetFloatVariable(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		key   string
		value float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetFloatVariable(tt.args.key, tt.args.value)
		})
	}
}

func Test_user_GetFloatVariable(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float32
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			got, err := u.GetFloatVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.GetFloatVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("user.GetFloatVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_SetIntVariable(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		key   string
		value int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "setInt0",
			fields: fields{
				connID:        "",
				sessionID:     "",
				accessID:      "",
				name:          "name",
				isTest:        0,
				credit:        100.0,
				gameID:        "slot0",
				strVarTable:   make(map[string]string),
				floatVarTable: make(map[string]float32),
				intVarTable:   make(map[string]int),
			},
			args: args{"freeSpinTimes", 0},
		},
		{
			name: "setInt0",
			fields: fields{
				connID:        "",
				sessionID:     "",
				accessID:      "",
				name:          "name",
				isTest:        0,
				credit:        100.0,
				gameID:        "slot0",
				strVarTable:   make(map[string]string),
				floatVarTable: make(map[string]float32),
				intVarTable:   make(map[string]int),
			},
			args: args{"freeSpinTimes", 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			u.SetIntVariable(tt.args.key, tt.args.value)
			f,_:=u.GetIntVariable("freeSpinTimes")
			t.Logf("Test_user_SetIntVariable freeSpinTimes:%d\n",f)
		})
	}
}

func Test_user_GetIntVariable(t *testing.T) {
	type fields struct {
		connID        string
		sessionID     string
		accessID      string
		name          string
		isTest        int
		credit        float32
		gameID        string
		strVarTable   map[string]string
		floatVarTable map[string]float32
		intVarTable   map[string]int
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				connID:        tt.fields.connID,
				sessionID:     tt.fields.sessionID,
				accessID:      tt.fields.accessID,
				name:          tt.fields.name,
				isTest:        tt.fields.isTest,
				credit:        tt.fields.credit,
				gameID:        tt.fields.gameID,
				strVarTable:   tt.fields.strVarTable,
				floatVarTable: tt.fields.floatVarTable,
				intVarTable:   tt.fields.intVarTable,
			}
			got, err := u.GetIntVariable(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.GetIntVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("user.GetIntVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}
