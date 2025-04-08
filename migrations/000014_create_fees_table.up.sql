CREATE TABLE fees (
    id        INTEGER PRIMARY KEY,
    orderid    INTEGER NOT NULL,
    paymentid  INTEGER,
    sum        DOUBLE PRECISION,

    CONSTRAINT fk_fees_order
        FOREIGN KEY (orderid) REFERENCES orders(id),
    CONSTRAINT fk_fees_payment
        FOREIGN KEY (paymentid) REFERENCES payments(id)
);