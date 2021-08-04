package ipaddress

import (
	"net"
	"reflect"
	"testing"
)

func TestNewIPv4Address(t *testing.T) {
	tests := []struct {
		name    string
		want    *IPv4Address
		wantErr bool
	}{
		{
			name:    "127.0.0.1",
			want:    &IPv4Address{n: net.ParseIP("127.0.0.1").To4(), prefix: 32},
			wantErr: false,
		},
		{
			name:    "192.168.0.1",
			want:    &IPv4Address{n: net.ParseIP("192.168.0.1").To4(), prefix: 32},
			wantErr: false,
		},
		{
			name:    "127.0.0.1/24",
			want:    &IPv4Address{n: net.ParseIP("127.0.0.1").To4(), prefix: 24},
			wantErr: false,
		},
		{name: "a.b.c.d", wantErr: true},
		{name: "127.0.0.1/a", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewIPv4Address(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIPv4Address() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIPv4Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_String(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "127.0.0.1"},
		{name: "192.168.0.1"},
		{name: "8.8.8.8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4, _ := NewIPv4Address(tt.name)
			if got := ipv4.String(); got != tt.name {
				t.Errorf("IPv4Address.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func TestIPv4Address_Bits(t *testing.T) {
	tests := []struct {
		name string
		ipv4 *IPv4Address
		want string
	}{
		{
			name: "127.0.0.1",
			ipv4: &IPv4Address{n: net.ParseIP("127.0.0.1")},
			want: "01111111000000000000000000000001",
		},
		{
			name: "192.168.0.1",
			ipv4: &IPv4Address{n: net.ParseIP("192.168.0.1")},
			want: "11000000101010000000000000000001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ipv4.Bits(); got != tt.want {
				t.Errorf("IPv4Address.Bits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Uint32(t *testing.T) {
	tests := []struct {
		name string
		ipv4 *IPv4Address
		want uint32
	}{
		{
			name: "127.0.0.1",
			ipv4: &IPv4Address{n: net.ParseIP("127.0.0.1")},
			want: 2130706433,
		},
		{
			name: "192.168.0.1",
			ipv4: &IPv4Address{n: net.ParseIP("192.168.0.1")},
			want: 3232235521,
		},
		{
			name: "255.255.255.255",
			ipv4: &IPv4Address{n: net.ParseIP("255.255.255.255")},
			want: 4294967295,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ipv4.Uint32(); got != tt.want {
				t.Errorf("IPv4Address.Bits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToIPv4Address(t *testing.T) {
	type args struct {
		n uint32
	}
	tests := []struct {
		name string
		args args
		want *IPv4Address
	}{
		{
			name: "123456789",
			args: args{n: 123456789},
			want: &IPv4Address{n: net.ParseIP("7.91.205.21").To4(), prefix: 32},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToIPv4Address(tt.args.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToIPv4Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Network(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   *IPv4Address
	}{
		{
			name: "7.91.205.21/24",
			fields: fields{
				n:      net.ParseIP("7.91.205.21").To4(),
				prefix: 24,
			},
			want: &IPv4Address{n: net.ParseIP("7.91.205.0").To4(), prefix: 24},
		},
		{
			name: "7.91.205.21/15",
			fields: fields{
				n:      net.ParseIP("7.91.205.15").To4(),
				prefix: 15,
			},
			want: &IPv4Address{n: net.ParseIP("7.90.0.0").To4(), prefix: 15},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			got := ipv4.Network()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv4Address.Network() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestIPv4Address_Broadcast(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   *IPv4Address
	}{
		{
			name: "7.91.205.21/24",
			fields: fields{
				n:      net.ParseIP("7.91.205.21").To4(),
				prefix: 24,
			},
			want: &IPv4Address{n: net.ParseIP("7.91.205.255").To4(), prefix: 24},
		},
		{
			name: "7.91.205.21/12",
			fields: fields{
				n:      net.ParseIP("7.91.205.21").To4(),
				prefix: 12,
			},
			want: &IPv4Address{n: net.ParseIP("7.95.255.255").To4(), prefix: 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			got := ipv4.Broadcast()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv4Address.Broadcast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Next(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   *IPv4Address
	}{
		{
			name: "127.0.0.1/24",
			fields: fields{
				n:      net.ParseIP("127.0.0.1").To4(),
				prefix: 24,
			},
			want: &IPv4Address{n: net.ParseIP("127.0.0.2").To4(), prefix: 24},
		},
		{
			name: "192.168.0.1/32",
			fields: fields{
				n:      net.ParseIP("192.168.0.1").To4(),
				prefix: 32,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			got := ipv4.Next()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv4Address.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_IsPrivate(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "127.0.0.1/32",
			fields: fields{
				n:      net.ParseIP("127.0.0.1"),
				prefix: 32,
			},
			want: false,
		},
		{
			name: "192.168.0.1/24",
			fields: fields{
				n:      net.ParseIP("192.168.0.1"),
				prefix: 24,
			},
			want: true,
		},
		{
			name: "172.17.18.19/22",
			fields: fields{
				n:      net.ParseIP("172.17.18.19"),
				prefix: 22,
			},
			want: true,
		},
		{
			name: "10.20.30.40/10",
			fields: fields{
				n:      net.ParseIP("10.20.30.40"),
				prefix: 10,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			if got := ipv4.IsPrivate(); got != tt.want {
				t.Errorf("IPv4Address.IsPrivate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Bytes(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "127.0.0.1",
			fields: fields{
				n:      net.ParseIP("127.0.0.1").To4(),
				prefix: 32,
			},
			want: []byte{127, 0, 0, 1},
		},
		{
			name: "1.2.3.4",
			fields: fields{
				n:      net.ParseIP("1.2.3.4").To4(),
				prefix: 32,
			},
			want: []byte{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			if got := ipv4.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv4Address.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Netmask(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "/24",
			fields: fields{
				prefix: 24,
			},
			want: "255.255.255.0",
		},
		{
			name: "/10",
			fields: fields{
				prefix: 10,
			},
			want: "255.192.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			if got := ipv4.Netmask(); got != tt.want {
				t.Errorf("IPv4Address.Netmask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Prev(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   *IPv4Address
	}{
		{
			name: "192.168.0.1",
			fields: fields{
				n:      net.ParseIP("192.168.0.1").To4(),
				prefix: 24,
			},
			want: &IPv4Address{
				n:      net.ParseIP("192.168.0.0").To4(),
				prefix: 24,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			if got := ipv4.Prev(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv4Address.Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validatePrefix(t *testing.T) {
	type args struct {
		prefix int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "-1", args: args{prefix: -1}, want: false},
		{name: "0", args: args{prefix: 0}, want: true},
		{name: "32", args: args{prefix: 32}, want: true},
		{name: "33", args: args{prefix: 33}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validatePrefix(tt.args.prefix); got != tt.want {
				t.Errorf("validatePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Class(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ClassA",
			fields: fields{
				n:      net.ParseIP("10.20.30.40").To4(),
				prefix: 8,
			},
			want: "A",
		},
		{
			name: "ClassB",
			fields: fields{
				n:      net.ParseIP("172.20.30.40").To4(),
				prefix: 12,
			},
			want: "B",
		},
		{
			name: "ClassC",
			fields: fields{
				n:      net.ParseIP("192.168.100.200").To4(),
				prefix: 16,
			},
			want: "C",
		},
		{
			name: "Loopback",
			fields: fields{
				n:      net.ParseIP("127.0.0.1").To4(),
				prefix: 32,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			if got := ipv4.Class(); got != tt.want {
				t.Errorf("IPv4Address.Class() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Sample(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Sample",
			fields: fields{
				n:      net.ParseIP("192.168.0.1").To4(),
				prefix: 24,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}

			got := ipv4.Sample()
			if ipv4.Contains(got) != tt.want {
				t.Errorf("IPv4Address.Sample() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_Size(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "/24",
			fields: fields{prefix: 24},
			want:   256,
		},
		{
			name:   "/22",
			fields: fields{prefix: 22},
			want:   1024,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			if got := ipv4.Size(); got != tt.want {
				t.Errorf("IPv4Address.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4Address_ARPA(t *testing.T) {
	type fields struct {
		n      net.IP
		prefix int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "127.0.0.1",
			fields: fields{
				n:      net.ParseIP("127.0.0.1").To4(),
				prefix: 32,
			},
			want: "1.0.0.127.in-addr.arpa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipv4 := &IPv4Address{
				n:      tt.fields.n,
				prefix: tt.fields.prefix,
			}
			if got := ipv4.ARPA(); got != tt.want {
				t.Errorf("IPv4Address.ARPA() = %v, want %v", got, tt.want)
			}
		})
	}
}
