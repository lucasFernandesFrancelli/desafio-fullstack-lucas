-- Schema inicial: usuários, categorias, aprovadores por categoria, solicitações e histórico.

CREATE TABLE users (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    role       TEXT NOT NULL CHECK (role IN ('colaborador', 'analista', 'gestor')) DEFAULT 'colaborador',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE categories (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE category_approvers (
    id             UUID PRIMARY KEY,
    category_id    UUID NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    user_id        UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    approval_order SMALLINT NOT NULL CHECK (approval_order IN (1, 2)),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (category_id, approval_order),
    UNIQUE (category_id, user_id)
);

CREATE TABLE solicitations (
    id                    UUID PRIMARY KEY,
    title                 TEXT NOT NULL,
    problem_description   TEXT NOT NULL DEFAULT '',
    proposed_improvement  TEXT NOT NULL DEFAULT '',
    category_id           UUID REFERENCES categories (id) ON DELETE RESTRICT,
    location              TEXT NOT NULL DEFAULT '',
    requester_id          UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    status                TEXT NOT NULL CHECK (status IN
                              ('rascunho', 'em_aprovacao', 'em_analise', 'finalizada', 'recusada'))
                              DEFAULT 'rascunho',
    current_approval_step SMALLINT CHECK (current_approval_step IN (1, 2)),
    severity              SMALLINT CHECK (severity BETWEEN 1 AND 5),
    urgency               SMALLINT CHECK (urgency BETWEEN 1 AND 5),
    trend                 SMALLINT CHECK (trend BETWEEN 1 AND 5),
    priority              SMALLINT GENERATED ALWAYS AS (severity * urgency * trend) STORED,
    analysis_notes        TEXT NOT NULL DEFAULT '',
    rejection_reason      TEXT NOT NULL DEFAULT '',
    last_transition_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE solicitation_history (
    id              UUID PRIMARY KEY,
    solicitation_id UUID NOT NULL REFERENCES solicitations (id) ON DELETE CASCADE,
    from_status     TEXT CHECK (from_status IN
                        ('rascunho', 'em_aprovacao', 'em_analise', 'finalizada', 'recusada')),
    to_status       TEXT NOT NULL CHECK (to_status IN
                        ('rascunho', 'em_aprovacao', 'em_analise', 'finalizada', 'recusada')),
    action          TEXT NOT NULL CHECK (action IN
                        ('criada', 'enviada', 'aprovada_etapa1', 'aprovada_etapa2', 'recusada', 'finalizada')),
    actor_id        UUID NOT NULL REFERENCES users (id),
    approval_step   SMALLINT,
    comment         TEXT NOT NULL DEFAULT '',
    severity        SMALLINT,
    urgency         SMALLINT,
    trend           SMALLINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_solicitations_status ON solicitations (status);
CREATE INDEX idx_solicitations_category ON solicitations (category_id);
CREATE INDEX idx_solicitations_requester ON solicitations (requester_id);
CREATE INDEX idx_history_solicitation ON solicitation_history (solicitation_id, created_at);
