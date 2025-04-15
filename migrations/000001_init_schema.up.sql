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

CREATE TABLE student_courses (
                                 student_id INTEGER REFERENCES students(id),
                                 course_id INTEGER REFERENCES courses(id),
                                 created_at TIMESTAMP,
                                 PRIMARY KEY(student_id, course_id)
);
