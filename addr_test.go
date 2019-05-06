package p2p

import (
	"testing"
)

func TestAddr_ParseString(t *testing.T) {
	type fields struct {
		Identifier string
		Port       uint16
		Secret     []byte
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"", fields{"UniqueIdentifier1", 1234, []byte("Secret")}, args{"KVXGS4-LVMVEW-IZLOOR-UWM2LF-OIYS2M-JSGM2C-2U3FMN-ZGK5A"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Addr{
				Identifier: tt.fields.Identifier,
				Port:       tt.fields.Port,
				Secret:     tt.fields.Secret,
			}
			if err := i.ParseString(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ID.ParseString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddr_String(t *testing.T) {
	type fields struct {
		Identifier string
		Port       uint16
		Secret     []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{"UniqueIdentifier1", 1234, []byte("Secret")}, "KVXGS4-LVMVEW-IZLOOR-UWM2LF-OIYS2M-JSGM2C-2U3FMN-ZGK5A"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Addr{
				Identifier: tt.fields.Identifier,
				Port:       tt.fields.Port,
				Secret:     tt.fields.Secret,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("ID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
