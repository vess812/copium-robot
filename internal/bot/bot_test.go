package bot

import (
	"reflect"
	"testing"

	"copium-bot/internal/domain"

	"github.com/stretchr/testify/mock"
)

type mockTranscriber struct {
	mock.Mock
}

func (m *mockTranscriber) Process(r domain.Request) (domain.Response, error) {
	args := m.Called(r)
	return args.Get(0).(domain.Response), args.Error(1)
}

func TestBot_Process(t *testing.T) {
	type fields struct {
		transcriber domain.Processor
	}
	type args struct {
		r domain.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Response
		wantErr bool
	}{
		{
			name: "empty user id",
			fields: fields{
				transcriber: func() *mockTranscriber {
					m := &mockTranscriber{}
					return m
				}(),
			},
			args: args{r: domain.Request{
				User: domain.User{},
				Message: domain.Message{
					ID:     2,
					ChatID: 3,
				},
			}},
			want:    domain.Response{},
			wantErr: true,
		},
		{
			name: "empty message id",
			fields: fields{
				transcriber: func() *mockTranscriber {
					m := &mockTranscriber{}
					return m
				}(),
			},
			args: args{r: domain.Request{
				User: domain.User{
					ID:   1,
					Name: "test",
				},
				Message: domain.Message{
					ChatID: 3,
				},
			}},
			want:    domain.Response{},
			wantErr: true,
		},
		{
			name: "empty chat id",
			fields: fields{
				transcriber: func() *mockTranscriber {
					m := &mockTranscriber{}
					return m
				}(),
			},
			args: args{r: domain.Request{
				User: domain.User{
					ID:   1,
					Name: "test",
				},
				Message: domain.Message{
					ID: 2,
				},
			}},
			want:    domain.Response{},
			wantErr: true,
		},
		{
			name: "voice message",
			fields: fields{
				transcriber: func() *mockTranscriber {
					m := &mockTranscriber{}
					m.On("Process", mock.Anything).Return(domain.Response{
						ChatID:  3,
						ReplyTo: 2,
						Text:    "test",
					}, nil)
					return m
				}(),
			},
			args: args{r: domain.Request{
				User: domain.User{
					ID:   1,
					Name: "test",
				},
				Message: domain.Message{
					ID:     2,
					ChatID: 3,
					Voice:  make([]byte, 0),
				},
			}},
			want: domain.Response{
				ChatID:  3,
				ReplyTo: 2,
				Text:    "test",
			},
		},
		{
			name: "video note message",
			fields: fields{
				transcriber: func() *mockTranscriber {
					m := &mockTranscriber{}
					m.On("Process", mock.Anything).Return(domain.Response{
						ChatID:  3,
						ReplyTo: 2,
						Text:    "test",
					}, nil)
					return m
				}(),
			},
			args: args{r: domain.Request{
				User: domain.User{
					ID:   1,
					Name: "test",
				},
				Message: domain.Message{
					ID:        2,
					ChatID:    3,
					VideoNote: make([]byte, 0),
				},
			}},
			want: domain.Response{
				ChatID:  3,
				ReplyTo: 2,
				Text:    "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBot(Opts{Transcriber: tt.fields.transcriber})
			got, err := b.Process(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() got = %v, want %v", got, tt.want)
			}
		})
	}
}
