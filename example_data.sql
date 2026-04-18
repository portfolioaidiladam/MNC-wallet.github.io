-- ---------------------------------------------------------------------------
-- MNC Wallet — Sample data untuk reviewer.
--
-- PIN untuk semua 3 user = 123456 (bcrypt hash di-hardcode).
-- Dipakai untuk smoke test login + transfer antar-user.
--
-- Cara load:
--   psql $DATABASE_URL -f example_data.sql
-- ---------------------------------------------------------------------------

BEGIN;

-- Bersihkan data lama (order: child -> parent) supaya seed idempotent.
DELETE FROM transactions;
DELETE FROM refresh_tokens;
DELETE FROM wallets;
DELETE FROM users;

-- Users
INSERT INTO users (id, first_name, last_name, phone_number, address, pin_hash, created_at, updated_at)
VALUES
    ('11111111-1111-1111-1111-111111111111',
     'Aidil', 'Adam', '081234567890', 'Jakarta',
     '$2a$10$1cPR/GzLym95UXFjsMA8uejgqY2PKPAseldB8JAilTt1lKVmJEpGW',
     NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222',
     'Budi', 'Santoso', '081298765432', 'Bandung',
     '$2a$10$1cPR/GzLym95UXFjsMA8uejgqY2PKPAseldB8JAilTt1lKVmJEpGW',
     NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333',
     'Citra', 'Lestari', '081355512345', 'Surabaya',
     '$2a$10$1cPR/GzLym95UXFjsMA8uejgqY2PKPAseldB8JAilTt1lKVmJEpGW',
     NOW(), NOW());

-- Wallets (satu per user)
INSERT INTO wallets (id, user_id, balance, created_at, updated_at)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
     '11111111-1111-1111-1111-111111111111',
     380000, NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
     '22222222-2222-2222-2222-222222222222',
     250000, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc',
     '33333333-3333-3333-3333-333333333333',
     100000, NOW(), NOW());

-- Sample transactions (Aidil: topup + payment + transfer out ke Budi)
INSERT INTO transactions (id, user_id, type, status, amount, balance_before, balance_after, remarks, reference_id, counterparty_user_id, created_at)
VALUES
    -- Aidil topup 500k
    ('d1111111-1111-1111-1111-111111111111',
     '11111111-1111-1111-1111-111111111111',
     'TOPUP', 'SUCCESS',
     500000, 0, 500000,
     '', NULL, NULL, NOW() - INTERVAL '2 days'),

    -- Aidil bayar 20k ke merchant
    ('d2222222-2222-2222-2222-222222222222',
     '11111111-1111-1111-1111-111111111111',
     'PAYMENT', 'SUCCESS',
     20000, 500000, 480000,
     'beli kopi', NULL, NULL, NOW() - INTERVAL '1 day'),

    -- Aidil transfer 100k ke Budi (sukses)
    ('d3333333-3333-3333-3333-333333333333',
     '11111111-1111-1111-1111-111111111111',
     'TRANSFER_OUT', 'SUCCESS',
     100000, 480000, 380000,
     'bayar hutang',
     'd3333333-3333-3333-3333-333333333333',
     '22222222-2222-2222-2222-222222222222',
     NOW() - INTERVAL '12 hours'),

    -- Sisi Budi: TRANSFER_IN berpasangan (reference_id sama)
    ('d4444444-4444-4444-4444-444444444444',
     '22222222-2222-2222-2222-222222222222',
     'TRANSFER_IN', 'SUCCESS',
     100000, 150000, 250000,
     'bayar hutang',
     'd3333333-3333-3333-3333-333333333333',
     '11111111-1111-1111-1111-111111111111',
     NOW() - INTERVAL '12 hours'),

    -- Budi topup 150k sebelumnya
    ('d5555555-5555-5555-5555-555555555555',
     '22222222-2222-2222-2222-222222222222',
     'TOPUP', 'SUCCESS',
     150000, 0, 150000,
     '', NULL, NULL, NOW() - INTERVAL '3 days'),

    -- Citra topup 100k
    ('d6666666-6666-6666-6666-666666666666',
     '33333333-3333-3333-3333-333333333333',
     'TOPUP', 'SUCCESS',
     100000, 0, 100000,
     '', NULL, NULL, NOW() - INTERVAL '4 hours');

-- Final balance konsisten dengan jumlah transaksi:
--   Aidil : 0 + 500k - 20k - 100k = 380k
--   Budi  : 0 + 150k + 100k       = 250k
--   Citra : 0 + 100k               = 100k

COMMIT;