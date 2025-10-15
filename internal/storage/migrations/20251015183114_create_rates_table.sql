-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exchange_rates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    currency_code VARCHAR(3) NOT NULL,
    rate DECIMAL(18, 8) NOT NULL,
    effective_date DATE NOT NULL,
    INDEX idx_currency_code (currency_code),
    UNIQUE INDEX idx_date_code (currency_code, effective_date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exchange_rates;
-- +goose StatementEnd
