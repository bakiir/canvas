CREATE TABLE teachers (
                          id SERIAL PRIMARY KEY,
                          name TEXT NOT NULL,
                          login TEXT UNIQUE NOT NULL,
                          password TEXT NOT NULL,
                          created_at TIMESTAMP,
                          updated_at TIMESTAMP,
                          deleted_at TIMESTAMP
);

CREATE TABLE students (
                          id SERIAL PRIMARY KEY,
                          name TEXT NOT NULL,
                          login TEXT UNIQUE NOT NULL,
                          password TEXT NOT NULL,
                          grades JSONB,
                          created_at TIMESTAMP,
                          updated_at TIMESTAMP,
                          deleted_at TIMESTAMP
);

CREATE TABLE courses (
                         id SERIAL PRIMARY KEY,
                         name TEXT NOT NULL,
                         capacity INTEGER,
                         teacher_id INTEGER REFERENCES teachers(id) ON DELETE SET NULL ON UPDATE CASCADE,
                         created_at TIMESTAMP,
                         updated_at TIMESTAMP,
                         deleted_at TIMESTAMP
);

CREATE TABLE tasks (
                       id BIGSERIAL PRIMARY KEY,
                       created_at TIMESTAMP,
                       updated_at TIMESTAMP,
                       deleted_at TIMESTAMP,

                       title VARCHAR(255),
                       description TEXT,
                       deadline TIMESTAMP,

                       course_id BIGINT,
                       CONSTRAINT fk_course FOREIGN KEY (course_id)
                           REFERENCES courses(id)
                           ON UPDATE CASCADE
                           ON DELETE SET NULL
);


CREATE TABLE homeworks (
                           id BIGSERIAL PRIMARY KEY,
                           created_at TIMESTAMP,
                           updated_at TIMESTAMP,
                           deleted_at TIMESTAMP,

                           student_id BIGINT NOT NULL,
                           task_id BIGINT NOT NULL,

                           file_url VARCHAR(512),
                           uploaded_at TIMESTAMP,

                           CONSTRAINT fk_student FOREIGN KEY (student_id)
                               REFERENCES students(id)
                               ON UPDATE CASCADE
                               ON DELETE CASCADE,

                           CONSTRAINT fk_task FOREIGN KEY (task_id)
                               REFERENCES tasks(id)
                               ON UPDATE CASCADE
                               ON DELETE CASCADE
);

CREATE TABLE student_courses (
                                 student_id INTEGER REFERENCES students(id),
                                 course_id INTEGER REFERENCES courses(id),
                                 created_at TIMESTAMP,
                                 PRIMARY KEY(student_id, course_id)
);
