-- 教练表
CREATE TABLE IF NOT EXISTS coaches (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id),
    title VARCHAR(64) DEFAULT '',
    specialties JSONB DEFAULT '[]',
    bio TEXT DEFAULT '',
    hourly_rate DECIMAL(10,2) DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 5.0,
    student_count INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 课程排期表
CREATE TABLE IF NOT EXISTS course_schedules (
    id BIGSERIAL PRIMARY KEY,
    coach_id BIGINT NOT NULL REFERENCES coaches(id),
    venue_id BIGINT NOT NULL REFERENCES venues(id),
    student_id BIGINT REFERENCES users(id),
    product_id BIGINT REFERENCES products(id),
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    type VARCHAR(20) DEFAULT 'private',
    status VARCHAR(20) DEFAULT 'scheduled',
    recurrence_rule VARCHAR(128) DEFAULT '',
    recurrence_group_id UUID,
    substitute_coach_id BIGINT REFERENCES coaches(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_schedules_coach_date ON course_schedules(coach_id, date);
CREATE INDEX idx_schedules_venue_date ON course_schedules(venue_id, date);
CREATE INDEX idx_schedules_coach_date_time ON course_schedules(coach_id, date, start_time, end_time);
CREATE INDEX idx_schedules_venue_date_time ON course_schedules(venue_id, date, start_time, end_time);
CREATE INDEX idx_schedules_recurrence ON course_schedules(recurrence_group_id);
CREATE INDEX idx_schedules_student ON course_schedules(student_id);

-- 教练休假表
CREATE TABLE IF NOT EXISTS coach_leaves (
    id BIGSERIAL PRIMARY KEY,
    coach_id BIGINT NOT NULL REFERENCES coaches(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason VARCHAR(256) DEFAULT '',
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 学员课程记录表
CREATE TABLE IF NOT EXISTS student_course_records (
    id BIGSERIAL PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES course_schedules(id),
    student_id BIGINT NOT NULL REFERENCES users(id),
    coach_id BIGINT NOT NULL REFERENCES coaches(id),
    attendance VARCHAR(20) DEFAULT 'present',
    coach_feedback TEXT DEFAULT '',
    rating SMALLINT,
    rating_comment TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_records_student ON student_course_records(student_id);
CREATE INDEX idx_records_coach ON student_course_records(coach_id);
CREATE INDEX idx_records_schedule ON student_course_records(schedule_id);
