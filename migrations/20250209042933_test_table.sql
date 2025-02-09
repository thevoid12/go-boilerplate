-- +goose Up
-- +goose StatementBegin
CREATE TABLE test (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test;
-- +goose StatementEnd
