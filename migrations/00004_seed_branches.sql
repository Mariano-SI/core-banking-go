-- +goose Up
INSERT INTO branches (id, code, name) VALUES
    ('01a18509-1bbf-731b-bac6-46138a22a0f2', '0001', 'Agência Central'),
    ('01a18509-1bca-7d87-9c69-228f1be528d9', '0002', 'Agência Digital');

-- +goose Down
DELETE FROM branches WHERE code IN ('0001', '0002');
