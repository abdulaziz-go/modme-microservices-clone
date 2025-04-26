CREATE TABLE IF NOT EXISTS users
(
    id           uuid primary key,
    full_name    varchar                                                                       NOT NULL,
    phone_number varchar UNIQUE                                                                NOT NULL,
    password     varchar                                                                       NOT NULL,
    role         varchar check ( role in ('SUPER_CEO', 'CEO', 'TEACHER', 'ADMIN', 'EMPLOYEE' , 'FINANCIST')) NOT NULL,
    birth_date   date                                                                                   DEFAULT '2000-12-12',
    gender       boolean                                                                       NOT NULL DEFAULT TRUE,
    is_deleted   boolean                                                                       NOT NULL DEFAULT FALSE,
    created_at   timestamp                                                                              DEFAULT NOW(),
    company_id   int,
    has_access_finance boolean DEFAULT FALSE
);

INSERT INTO users(id, full_name, phone_number, password, role)
values (gen_random_uuid(), 'Shohruh', '+998950960153', '$2a$10$1gDxC.3v73V45QXt0R3cCurQE5YL5jB5HTRKrh8L1maJx68nySEtW',
        'SUPER_CEO');

CREATE TABLE IF NOT EXISTS users_history
(
    id            uuid primary key,
    user_id       uuid references users (id),
    updated_field varchar   NOT NULL,
    old_value     varchar   NOT NULL,
    current_value varchar   NOT NULL,
    created_at    timestamp NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION log_user_updates()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Check for changes in each field and insert a record in users_history if there are changes
    IF NEW.full_name IS DISTINCT FROM OLD.full_name THEN
        INSERT INTO users_history (id, user_id, updated_field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'full_name', OLD.full_name, NEW.full_name, NOW());
    END IF;

    IF NEW.phone_number IS DISTINCT FROM OLD.phone_number THEN
        INSERT INTO users_history (id, user_id, updated_field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'phone_number', OLD.phone_number, NEW.phone_number, NOW());
    END IF;

    IF NEW.password IS DISTINCT FROM OLD.password THEN
        INSERT INTO users_history (id, user_id, updated_field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'password', OLD.password, NEW.password, NOW());
    END IF;

    IF NEW.role IS DISTINCT FROM OLD.role THEN
        INSERT INTO users_history (id, user_id, updated_field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'role', OLD.role, NEW.role, NOW());
    END IF;

    IF NEW.birth_date IS DISTINCT FROM OLD.birth_date THEN
        INSERT INTO users_history (id, user_id, updated_field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'birth_date', OLD.birth_date::text, NEW.birth_date::text, NOW());
    END IF;

    IF NEW.gender IS DISTINCT FROM OLD.gender THEN
        INSERT INTO users_history (id, user_id, updated_field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'gender', OLD.gender::text, NEW.gender::text, NOW());
    END IF;

    IF NEW.company_id IS DISTINCT FROM OLD.company_id THEN
        INSERT INTO users_history (id, user_id, updated_field, old_value, current_value, created_at)
        VALUES (gen_random_uuid(), OLD.id, 'company_id', OLD.company_id::text, NEW.company_id::text, NOW());
    END IF;

    RETURN NEW; -- Return the new record
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_update_trigger
    AFTER UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION log_user_updates();
