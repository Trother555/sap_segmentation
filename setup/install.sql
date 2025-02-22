CREATE TABLE IF NOT EXISTS segmentation (
    id                  BIGSERIAL       PRIMARY KEY,
    address_sap_id      VARCHAR(255)    NOT NULL UNIQUE ,
    adr_segment         VARCHAR(16)     NOT NULL,
    segment_id          BIGINT          NOT NULL
)
