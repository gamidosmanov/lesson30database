DROP TABLE IF EXISTS
    devbase.tasks.tasks_labels
    , devbase.tasks.tasks
    , devbase.tasks.labels
    , devbase.tasks.users;

-- пользователи системы
CREATE TABLE devbase.tasks.users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- метки задач
CREATE TABLE devbase.tasks.labels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- задачи
CREATE TABLE devbase.tasks.tasks (
    id SERIAL PRIMARY KEY,
    opened BIGINT NOT NULL, -- время создания задачи
    closed BIGINT DEFAULT 0, -- время выполнения задачи
    author_id INTEGER REFERENCES devbase.tasks.users(id), -- автор задачи
    assigned_id INTEGER REFERENCES devbase.tasks.users(id), -- ответственный
    title TEXT, -- название задачи
    content TEXT -- текст задачи
);

-- связь многие-ко-многим между задачами и метками
CREATE TABLE devbase.tasks.tasks_labels (
    task_id INTEGER REFERENCES devbase.tasks.tasks(id),
    label_id INTEGER REFERENCES devbase.tasks.labels(id)
);
