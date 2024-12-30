CREATE SEQUENCE dashboard_success_quiz_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
    
CREATE TABLE public.success_quiz_log (
    id INT4 NOT NULL DEFAULT nextval('dashboard_success_quiz_log_id_seq'::regclass) PRIMARY KEY,
    passage INT2,
    total INT2,
    success INT2,
    create_at BIGINT NOT NULL,
    date DATE NOT NULL,
    update_at BIGINT NOT NULL,
    status INT2 NOT NULL,
    question_type TEXT,
    user_id UUID NOT NULL REFERENCES public.users(id),
    skill INT2 NOT NULL,
    answer_id INT4 REFERENCES public.answers(id),
    failed INT2 DEFAULT 0,
    skipped INT2 DEFAULT 0
);

--DONE--
