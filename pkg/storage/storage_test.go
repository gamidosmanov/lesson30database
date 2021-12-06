package storage

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	storage  *Storage
	connErr  error
	password string
)

func init() {
	// Я создал на локальном сервере базу devbase и на ней схему tasks
	password = "" // Здесь пароль
	connString := fmt.Sprintf("postgresql://localhost/devbase?user=postgres&password=%s", password)
	storage, connErr = New(connString)
	if connErr != nil {
		panic(connErr)
	}
}

func TestStorage_NewDeleteUpdate(t *testing.T) {
	testTask := Task{
		Opened:     1638777178,
		Closed:     0,
		AuthorID:   5,
		AssignedID: 3,
		Title:      "Добавить новую фичу на сайт",
		Content:    "Добавить заголовок",
	}
	testID, err := storage.NewTask(testTask)
	if err != nil {
		t.Error(err)
		return
	}
	testTask.ID = testID
	selectTest, err := storage.Tasks(testTask.ID, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(selectTest[0], testTask) {
		t.Error("Задача в базе некорректна")
		return
	}
	testTask.Closed = 1641462725
	updateId, err := storage.UpdateTask(testTask)
	if err != nil {
		t.Error(err)
		return
	}
	if updateId != testTask.ID {
		t.Error("Мы обновили не ту задачу")
		return
	}
	selectTest, err = storage.Tasks(testTask.ID, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(selectTest[0], testTask) {
		t.Error("Задача в базе некорректна")
		return
	}
	deleteID, err := storage.DeleteTask(testTask)
	if err != nil {
		t.Error(err)
		return
	}
	if deleteID != testTask.ID {
		t.Error("Хаха, мы удалили не ту задачу")
		return
	}
}

func TestStorage_Tasks(t *testing.T) {
	type args struct {
		taskID   int
		authorID int
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    []Task
		wantErr bool
	}{
		// Я завел на тестовый сервер две задачи и трех пользователей
		{
			name: "Task 1",
			s:    storage,
			args: args{taskID: 1, authorID: 0},
			want: []Task{
				{
					ID:         1,
					Opened:     1638777125,
					Closed:     0,
					AuthorID:   5,
					AssignedID: 3,
					Title:      "Добавить новую фичу на сайт",
					Content:    "Добавить кнопку",
				},
			},
		},
		{
			name: "Task 2",
			s:    storage,
			args: args{taskID: 2, authorID: 0},
			want: []Task{
				{
					ID:         2,
					Opened:     1638784325,
					Closed:     0,
					AuthorID:   5,
					AssignedID: 4,
					Title:      "Добавить новую фичу на сайт",
					Content:    "Добавить таблицу",
				},
			},
		},
		{
			name: "Task Author 5",
			s:    storage,
			args: args{taskID: 0, authorID: 5},
			want: []Task{
				{
					ID:         1,
					Opened:     1638777125,
					Closed:     0,
					AuthorID:   5,
					AssignedID: 3,
					Title:      "Добавить новую фичу на сайт",
					Content:    "Добавить кнопку",
				},
				{
					ID:         2,
					Opened:     1638784325,
					Closed:     0,
					AuthorID:   5,
					AssignedID: 4,
					Title:      "Добавить новую фичу на сайт",
					Content:    "Добавить таблицу",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Tasks(tt.args.taskID, tt.args.authorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Tasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Tasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_TasksByLabel(t *testing.T) {
	type args struct {
		labelID int
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    []Task
		wantErr bool
	}{
		{
			name: "Label 1",
			s:    storage,
			args: args{labelID: 1},
			want: []Task{
				{
					ID:         1,
					Opened:     1638777125,
					Closed:     0,
					AuthorID:   5,
					AssignedID: 3,
					Title:      "Добавить новую фичу на сайт",
					Content:    "Добавить кнопку",
				},
			},
		},
		{
			name: "Label 2",
			s:    storage,
			args: args{labelID: 2},
			want: []Task{
				{
					ID:         2,
					Opened:     1638784325,
					Closed:     0,
					AuthorID:   5,
					AssignedID: 4,
					Title:      "Добавить новую фичу на сайт",
					Content:    "Добавить таблицу",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.TasksByLabel(tt.args.labelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.TasksByLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.TasksByLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}
