CREATE TABLE payments (
    id             INTEGER PRIMARY KEY,
    paymentstatus  payment_status,
    amount         DOUBLE PRECISION,
    payment_date   TIMESTAMPTZ
);