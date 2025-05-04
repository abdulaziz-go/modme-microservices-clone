CREATE TABLE IF NOT EXISTS tariff
(
    id            serial primary key,
    name          varchar NOT NULL,
    student_count int     NOT NULL,
    sum           float   NOT NULL,
    discounts     jsonb,
    is_deleted    bool      DEFAULT FALSE,
    created_at    timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS company
(
    id            serial primary key,
    title         varchar                    NOT NULL,
    avatar        varchar                    NOT NULL,
    start_time    varchar                    NOT NULL,
    end_time      varchar                    NOT NULL,
    company_phone varchar                    NOT NULL,
    subdomain     varchar                    NOT NULL,
    valid_date    DATE                       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tariff_id     int references tariff (id) NOT NULL,
    discount_id   varchar,
    is_demo       bool                                DEFAULT FALSE,
    sms_balance  int DEFAULT 0,
    created_at    timestamp                           DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS company_payments
(
    id                serial primary key,
    company_id        int references company (id) NOT NULL,
    tariff_id         int references tariff (id)  NOT NULL,
    discount_id       varchar,
    discount_name     varchar,
    comment           varchar,
    sum               float                       NOT NULL,
    edited_valid_date date                        NOT NULL,
    created_at        timestamp DEFAULT CURRENT_TIMESTAMP
);


CREATE table IF NOT EXISTS rooms
(
    id         serial primary key,
    title      varchar NOT NULL,
    capacity   int     NOT NULL,
    company_id int references company (id)
);


CREATE TABLE IF NOT EXISTS courses
(
    id              serial primary key,
    title           varchar                                 NOT NULL,
    duration_lesson int                                     NOT NULL,
    course_duration int                                     NOT NULL,
    price           double precision check ( price > 5000 ) NOT NULL,
    description     text,
    company_id      int references company (id)
);


CREATE TABLE IF NOT EXISTS groups
(
    id          bigserial PRIMARY KEY,
    name        varchar                                                NOT NULL,
    course_id   int                                                    NOT NULL,
    teacher_id  uuid                                                   NOT NULL,
    room_id     int references rooms (id),
    date_type   varchar check (date_type in ('JUFT', 'TOQ', 'BOSHQA')) NOT NULL,
    days        TEXT[]                                                 NOT NULL,
    start_time  varchar                                                NOT NULL,
    start_date  date                                                   NOT NULL,
    end_date    date                                                   NOT NULL,
    is_archived boolean   DEFAULT FALSE                                NOT NULL,
    created_at  timestamp DEFAULT NOW(),
    CONSTRAINT valid_days CHECK (array_length(days, 1) > 0 AND days <@
                                                               ARRAY ['DUSHANBA', 'SESHANBA', 'CHORSHANBA', 'PAYSHANBA', 'JUMA', 'SHANBA', 'YAKSHANBA']),
    company_id  int references company (id)
);

CREATE OR REPLACE FUNCTION filter_groups(
    p_is_archived BOOLEAN DEFAULT NULL,
    p_teacher_id UUID DEFAULT NULL,
    p_course_id INT DEFAULT NULL,
    p_date_type VARCHAR DEFAULT NULL,
    p_start_date DATE DEFAULT NULL,
    p_end_date DATE DEFAULT NULL,
    p_company_id INT DEFAULT NULL
)
    RETURNS TABLE (
                      id            BIGINT,
                      course_id     INT,
                      course_name   VARCHAR,
                      teacher_id    UUID,
                      room_id       INT,
                      room_name     VARCHAR,
                      room_capacity INT,
                      start_date    DATE,
                      end_date      DATE,
                      is_archived   BOOLEAN,
                      name          VARCHAR,
                      student_count INT,
                      created_at    TIMESTAMP,
                      days          TEXT[],  -- ✅ Fixed: Explicitly casted to TEXT[]
                      start_time    VARCHAR,
                      date_type     VARCHAR
                  ) AS
$$
BEGIN
    RETURN QUERY
        SELECT g.id,
               g.course_id,
               c.title                                                               AS course_name,
               g.teacher_id,
               g.room_id,
               r.title                                                               AS room_name,
               r.capacity                                                            AS room_capacity,
               g.start_date,
               g.end_date,
               g.is_archived,
               g.name,
               (SELECT COUNT(gs.id) FROM group_students gs WHERE gs.group_id = g.id) AS student_count,
               g.created_at,
               g.days::TEXT[],  -- ✅ Fixed: Explicitly casted to TEXT[]
               g.start_time,
               g.date_type
        FROM groups g
                 LEFT JOIN courses c ON g.course_id = c.id
                 LEFT JOIN rooms r ON g.room_id = r.id
        WHERE (p_is_archived IS NULL OR g.is_archived = p_is_archived)
          AND (p_teacher_id IS NULL OR g.teacher_id = p_teacher_id)
          AND (p_course_id IS NULL OR g.course_id = p_course_id)
          AND (p_date_type IS NULL OR g.date_type = p_date_type)
          AND (p_start_date IS NULL OR g.start_date >= p_start_date)
          AND (p_end_date IS NULL OR g.end_date <= p_end_date)
          AND (p_company_id IS NULL OR g.company_id = p_company_id);
END;
$$ LANGUAGE plpgsql;




CREATE OR REPLACE FUNCTION sort_groups(
    p_order_by TEXT DEFAULT 'name',
    p_order_direction TEXT DEFAULT 'ASC',
    p_company_id INTEGER DEFAULT NULL
)
    RETURNS TABLE (
                      id            BIGINT,
                      course_id     INT,
                      course_name   VARCHAR,
                      teacher_id    UUID,
                      room_id       INT,
                      room_name     VARCHAR,
                      room_capacity INT,
                      start_date    DATE,
                      end_date      DATE,
                      is_archived   BOOLEAN,
                      name          VARCHAR,
                      student_count INT,
                      created_at    TIMESTAMP,
                      days          TEXT[],  -- ✅ Fixed: Explicitly casted to TEXT[]
                      start_time    VARCHAR,
                      date_type     VARCHAR
                  ) AS
$$
BEGIN
    RETURN QUERY
        SELECT
            g.id, g.course_id, c.title AS course_name,
            g.teacher_id, g.room_id, r.title AS room_name, r.capacity AS room_capacity,
            g.start_date, g.end_date, g.is_archived, g.name,
            (SELECT COUNT(gs.id) FROM group_students gs WHERE gs.group_id = g.id) AS student_count,
            g.created_at, g.days::TEXT[],  -- ✅ Fixed: Explicitly casted to TEXT[]
            g.start_time, g.date_type
        FROM groups g
                 LEFT JOIN courses c ON g.course_id = c.id
                 LEFT JOIN rooms r ON g.room_id = r.id
        WHERE g.company_id = p_company_id
        ORDER BY
            CASE WHEN p_order_by = 'name' THEN g.name END,
            CASE WHEN p_order_by = 'start_date' THEN g.start_date END,
            CASE WHEN p_order_by = 'end_date' THEN g.end_date END
            NULLS LAST;
END;
$$ LANGUAGE plpgsql;




CREATE TABLE IF NOT EXISTS attendance
(
    is_discounted  boolean                                                               DEFAULT FALSE,
    discount_owner varchar CHECK ( discount_owner in ('TEACHER', 'CENTER'))              DEFAULT 'TEACHER',
    price_type     varchar CHECK ( price_type in ('PERCENT', 'FIXED', 'DISCOUNT')),
    total_count    float                                                        NOT NULL DEFAULT 0,
    course_price   float                                                        NOT NULL DEFAULT 0,
    price          float                                                        NOT NULL,
    group_id       bigint references groups (id),
    student_id     uuid                                                         NOT NULL,
    teacher_id     uuid                                                         NOT NULL,
    attend_date    date                                                         NOT NULL,
    status         int                                                          NOT NULL,
    created_at     timestamp                                                             DEFAULT NOW(),
    created_by     uuid                                                         NOT NULL,
    creator_role   varchar CHECK ( creator_role in ('ADMIN', 'CEO', 'TEACHER')) NOT NULL,
    company_id     int references company (id),
    sms_send       bool default false,
    PRIMARY KEY (group_id, student_id, attend_date)
);

CREATE TABLE IF NOT EXISTS students
(
    id                 uuid PRIMARY KEY,
    name               varchar NOT NULL,
    phone              varchar NOT NULL,
    date_of_birth      date                                                  default '2000-12-12',
    balance            double precision                                      DEFAULT 0,
    condition          varchar CHECK ( condition IN ('ACTIVE', 'ARCHIVED') ) DEFAULT 'ACTIVE',
    additional_contact varchar,
    address            varchar,
    telegram_username  varchar,
    passport_id        varchar,
    gender             boolean,
    created_at         timestamp                                             DEFAULT now(),
    company_id         int references company (id)
);

CREATE TABLE IF NOT EXISTS student_note
(
    id         uuid primary key,
    student_id uuid references students (id) NOT NULL,
    comment    text                          NOT NULL,
    created_at timestamp DEFAULT NOW(),
    created_by uuid,
    company_id int references company (id)
);

CREATE TABLE IF NOT EXISTS group_history
(
    id            uuid primary key,
    group_id      bigint references groups (id) NOT NULL,
    field         varchar                       NOT NULL,
    old_value     varchar                       NOT NULL,
    current_value varchar                       NOT NULL,
    created_at    timestamp DEFAULT NOW(),
    company_id    int references company (id)
);

CREATE TABLE IF NOT EXISTS student_history
(
    id            uuid primary key,
    student_id    uuid references students (id) NOT NULL,
    field         varchar                       NOT NULL,
    old_value     varchar                       NOT NULL,
    current_value varchar                       NOT NULL,
    created_at    timestamp DEFAULT NOW(),
    company_id    int references company (id)
);


CREATE TABLE IF NOT EXISTS transfer_lesson
(
    id            uuid PRIMARY KEY,
    group_id      bigint references groups (id) NOT NULL,
    real_date     date                          NOT NULL,
    transfer_date date                          NOT NULL,
    company_id    int references company (id)
);

CREATE TABLE IF NOT EXISTS group_students
(
    id                 uuid PRIMARY KEY,
    group_id           bigint references groups (id) NOT NULL,
    student_id         uuid                          NOT NULL,
    condition          varchar check ( condition in ('FREEZE', 'ACTIVE', 'DELETE')) DEFAULT 'FREEZE',
    last_specific_date date                          NOT NULL                       DEFAULT NOW(),
    created_at         timestamp                                                    DEFAULT NOW(),
    created_by         uuid                          NOT NULL,
    company_id         int references company (id),
    UNIQUE (group_id, student_id)
);

CREATE TABLE IF NOT EXISTS group_student_condition_history
(
    id                  uuid primary key,
    group_student_id    uuid references group_students (id)                                  NOT NULL,
    student_id          uuid references students (id)                                        NOT NULL,
    group_id            bigint references groups (id)                                        NOT NULL,
    old_condition       varchar check ( old_condition in ('FREEZE', 'ACTIVE', 'DELETE'))     NOT NULL,
    current_condition   varchar check ( current_condition in ('FREEZE', 'ACTIVE', 'DELETE')) NOT NULL,
    is_eliminated_trial bool                                                                          DEFAULT FALSE,
    specific_date       date                                                                 NOT NULL DEFAULT NOW(),
    return_the_money    boolean                                                              NOT NULL DEFAULT FALSE,
    created_at          timestamp                                                                     DEFAULT NOW(),
    company_id          int references company (id)
);


-- SMS SERVICE TABLES
CREATE TABLE IF NOT EXISTS "sms_payments" (
                                "id" SERIAL PRIMARY KEY,
                                "company_id" int NOT NULL,
                                "comment" varchar,
                                "sum" float NOT NULL,
                                "sms_count" float NOT NULL,
                                "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);


CREATE TABLE IF NOT EXISTS "sms_template" (
                                "id" SERIAL PRIMARY KEY,
                                "company_id" int NOT NULL,
                                "texts" text[] NOT NULL,
                                "sms_count" int NOT NULL,
                                "action_type" varchar check ( action_type in ('BEFORE_PAYMENT_ALERT' , 'INSUFFICIENT_BALANCE_ALERT' , 'PAYMENT_SUCCESSFUL_ALERT' , 'JOINED_GROUP_ALERT' , 'BIRTHDAY_ALERT' , 'NOT_PARTICIPATE_ALERT')),
                                "insufficient_balance_send_count" int NOT NULL DEFAULT 1,
                                "sms_template_type" varchar check ( sms_template_type in ('ACTION' , 'TEMPLATE')),
                                "is_active" bool DEFAULT FALSE,
                                "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS "sms_used" (
                            "id" uuid PRIMARY KEY,
                            "company_id" int NOT NULL,
                            "sms_template_id" int NOT NULL,
                            "texts" text[] NOT NULL,
                            "sms_count" int NOT NULL,
                            "student_id" uuid references students(id),
                            "sms_used_type" varchar check ( sms_used_type in ('BY_SELF' , 'BY_TEMPLATE')),
                            "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

ALTER TABLE "sms_used" ADD FOREIGN KEY ("sms_template_id") REFERENCES "sms_template" ("id");

ALTER TABLE "sms_used" ADD FOREIGN KEY ("company_id") REFERENCES "company" ("id");

ALTER TABLE "sms_payments" ADD FOREIGN KEY ("company_id") REFERENCES "company" ("id");

ALTER TABLE "sms_template" ADD FOREIGN KEY ("id") REFERENCES "company" ("id");

alter table sms_used
    add created_by_id uuid;

alter table sms_used
    add created_by_name varchar;




CREATE INDEX IF NOT EXISTS idx_attendance_group_date ON attendance (group_id, attend_date);
CREATE INDEX IF NOT EXISTS idx_group_students_group ON group_students (group_id);
CREATE OR REPLACE FUNCTION log_group_update()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.name IS DISTINCT FROM OLD.name THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'name', COALESCE(OLD.name, ''), COALESCE(NEW.name, ''), NOW());
    END IF;

    IF NEW.course_id IS DISTINCT FROM OLD.course_id THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'course_id', COALESCE(OLD.course_id::text, ''),
                COALESCE(NEW.course_id::text, ''), NOW());
    END IF;

    IF NEW.teacher_id IS DISTINCT FROM OLD.teacher_id THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'teacher_id', COALESCE(OLD.teacher_id::text, ''),
                COALESCE(NEW.teacher_id::text, ''), NOW());
    END IF;

    IF NEW.room_id IS DISTINCT FROM OLD.room_id THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'room_id', COALESCE(OLD.room_id::text, ''), COALESCE(NEW.room_id::text, ''),
                NOW());
    END IF;

    IF NEW.date_type IS DISTINCT FROM OLD.date_type THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'date_type', COALESCE(OLD.date_type, ''), COALESCE(NEW.date_type, ''),
                NOW());
    END IF;

    IF NEW.start_time IS DISTINCT FROM OLD.start_time THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'start_time', COALESCE(OLD.start_time::text, ''),
                COALESCE(NEW.start_time::text, ''), NOW());
    END IF;

    IF NEW.start_date IS DISTINCT FROM OLD.start_date THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'start_date', COALESCE(OLD.start_date::text, ''),
                COALESCE(NEW.start_date::text, ''), NOW());
    END IF;

    IF NEW.end_date IS DISTINCT FROM OLD.end_date THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'end_date', COALESCE(OLD.end_date::text, ''),
                COALESCE(NEW.end_date::text, ''), NOW());
    END IF;

    IF NEW.is_archived IS DISTINCT FROM OLD.is_archived THEN
        INSERT INTO group_history (id, group_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'is_archived', COALESCE(OLD.is_archived::text, ''),
                COALESCE(NEW.is_archived::text, ''), NOW());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_group_update
    AFTER UPDATE
    ON groups
    FOR EACH ROW
EXECUTE FUNCTION log_group_update();


CREATE OR REPLACE FUNCTION log_student_update()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.name IS DISTINCT FROM OLD.name THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'name', COALESCE(OLD.name, ''), COALESCE(NEW.name, ''), NOW());
    END IF;

    IF NEW.phone IS DISTINCT FROM OLD.phone THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'phone', COALESCE(OLD.phone, ''), COALESCE(NEW.phone, ''), NOW());
    END IF;

    IF NEW.date_of_birth IS DISTINCT FROM OLD.date_of_birth THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'date_of_birth', COALESCE(OLD.date_of_birth::text, ''),
                COALESCE(NEW.date_of_birth::text, ''), NOW());
    END IF;

    IF NEW.condition IS DISTINCT FROM OLD.condition THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'condition', COALESCE(OLD.condition, ''), COALESCE(NEW.condition, ''),
                NOW());
    END IF;

    IF NEW.additional_contact IS DISTINCT FROM OLD.additional_contact THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'additional_contact', COALESCE(OLD.additional_contact, ''),
                COALESCE(NEW.additional_contact, ''), NOW());
    END IF;

    IF NEW.address IS DISTINCT FROM OLD.address THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'address', COALESCE(OLD.address, ''), COALESCE(NEW.address, ''), NOW());
    END IF;

    IF NEW.telegram_username IS DISTINCT FROM OLD.telegram_username THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'telegram_username', COALESCE(OLD.telegram_username, ''),
                COALESCE(NEW.telegram_username, ''), NOW());
    END IF;

    IF NEW.passport_id IS DISTINCT FROM OLD.passport_id THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'passport_id', COALESCE(OLD.passport_id, ''), COALESCE(NEW.passport_id, ''),
                NOW());
    END IF;

    IF NEW.gender IS DISTINCT FROM OLD.gender THEN
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'gender', COALESCE(OLD.gender::text, ''), COALESCE(NEW.gender::text, ''),
                NOW());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_student_update
    AFTER UPDATE
    ON students
    FOR EACH ROW
EXECUTE FUNCTION log_student_update();
