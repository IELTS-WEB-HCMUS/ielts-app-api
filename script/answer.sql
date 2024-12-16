CREATE SEQUENCE answer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TABLE public.answers (
    id INT4 NOT NULL DEFAULT nextval('answer_id_seq'::regclass) PRIMARY KEY,
    sort INT4,
    user_created UUID REFERENCES public.users(id),
    date_created TIMESTAMPTZ,
    user_updated UUID REFERENCES public.users(id),
    date_updated TIMESTAMPTZ,
    quiz INT4 REFERENCES public.quiz(id),
    detail JSON,
    summary JSON,
    type INT4 REFERENCES public.type(id),
    status VARCHAR(255),
    note TEXT,
    correction_statistics JSONB,
    quiz_type INT2 DEFAULT 0,
    completed_duration INT2
);

-- Indexes
CREATE UNIQUE INDEX answer_pkey ON public.answers (id);
CREATE INDEX answer_quiz_id_idx ON public.answers (quiz);


