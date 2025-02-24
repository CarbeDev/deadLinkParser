package data

import (
	"reflect"
	"testing"
)

func TestAppData_hasLink(t *testing.T) {
	type fields struct {
		InitialUrl string
		FoundLinks []FoundLink
	}
	type args struct {
		link string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should return true if link exist",
			fields: fields{
				InitialUrl: "http://anyUrl.com",
				FoundLinks: []FoundLink{
					{
						Link:    "/found",
						Visited: false,
						Alive:   false,
					},
				},
			},
			args: args{
				link: "/found",
			},
			want: true,
		},
		{
			name: "should return false if link don't exist",
			fields: fields{
				InitialUrl: "http://anyUrl.com",
				FoundLinks: []FoundLink{
					{
						Link:    "/found",
						Visited: false,
						Alive:   false,
					},
				},
			},
			args: args{
				link: "/not_found",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := AppData{
				InitialUrl: tt.fields.InitialUrl,
				FoundLinks: tt.fields.FoundLinks,
			}
			if got := data.hasLink(tt.args.link); got != tt.want {
				t.Errorf("hasLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addLinkFound(t *testing.T) {
	type args struct {
		link    string
		appData *AppData
	}
	tests := []struct {
		name string
		args args
		want []FoundLink
	}{
		{
			name: "Should add links when don't exist",
			args: args{
				link: "/new_links",
				appData: &AppData{
					InitialUrl: "http://anyUrl.com",
					FoundLinks: []FoundLink{
						{
							Link:    "/found",
							Visited: false,
							Alive:   false,
						},
					},
				},
			},
			want: []FoundLink{
				{
					Link:    "/found",
					Visited: false,
					Alive:   false,
				},
				{
					Link:    "/new_links",
					Visited: false,
					Alive:   false,
				},
			},
		},
		{
			name: "Should add links when don't exist",
			args: args{
				link: "/already_exist",
				appData: &AppData{
					InitialUrl: "http://anyUrl.com",
					FoundLinks: []FoundLink{
						{
							Link:    "/already_exist",
							Visited: true,
							Alive:   false,
						},
					},
				},
			},
			want: []FoundLink{
				{
					Link:    "/already_exist",
					Visited: true,
					Alive:   false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddLinkFound(tt.args.link, tt.args.appData)
			if !reflect.DeepEqual(tt.args.appData.FoundLinks, tt.want) {
				t.Errorf("Expected : %v /n Get : %v", tt.want, tt.args.appData.FoundLinks)
			}
		})
	}
}

func Test_initialiseAppData(t *testing.T) {
	type args struct {
		baseUrl string
	}
	tests := []struct {
		name string
		args args
		want AppData
	}{
		{
			name: "should set initial state",
			args: args{baseUrl: "http://baseurl.com"},
			want: AppData{
				InitialUrl: "http://baseurl.com",
				FoundLinks: []FoundLink{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitialiseAppData(tt.args.baseUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitialiseAppData() = %v, want %v", got, tt.want)
			}
		})
	}
}
