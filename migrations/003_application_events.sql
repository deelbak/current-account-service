CREATE TABLE application_events (
    id             BIGSERIAL PRIMARY KEY,
    application_id BIGINT        NOT NULL REFERENCES applications(id),
    from_state     VARCHAR(50)   NOT NULL,
    to_state       VARCHAR(50)   NOT NULL,
    event          VARCHAR(50)   NOT NULL,
    actor_id       BIGINT,
    actor_role     VARCHAR(50),
    comment        TEXT,
    occurred_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX idx_app_events_app_id ON application_events(application_id);
CREATE INDEX idx_app_events_occurred ON application_events(occurred_at DESC);

ALTER TABLE applications ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'draft';
CREATE INDEX idx_applications_status ON applications(status);