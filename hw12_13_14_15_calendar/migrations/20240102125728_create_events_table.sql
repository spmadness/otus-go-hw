-- +goose Up
-- +goose StatementBegin
CREATE TABLE events(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    date_start timestamp NOT NULL,
    date_end timestamp NOT NULL,
    description TEXT DEFAULT NULL,
    user_id UUID NOT NULL,
    date_post timestamp DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
