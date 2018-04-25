-- 数据表设计
DROP TABLE IF EXISTS public.User;
CREATE TABLE public.User (
  ID BIGINT PRIMARY KEY NOT NULL,
  Name VARCHAR(255) NOT NULL,
  Type VARCHAR(16) NOT NULL DEFAULT 'user'
);
COMMENT ON TABLE public.User IS '用户表';
COMMENT ON COLUMN public.User.ID IS '主键ID';
COMMENT ON COLUMN public.User.Name IS '用户姓名';
COMMENT ON COLUMN public.User.Type IS '用户类型';

CREATE SEQUENCE user_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;
-- 自增主键
ALTER TABLE public.User ALTER COLUMN ID SET DEFAULT nextval('user_id_seq');

INSERT INTO public.User VALUES
 (1,'二毛','user'),(2,'四毛','user');

DROP TABLE IF EXISTS public.Relationship;
CREATE TABLE public.Relationship (
  ID BIGINT PRIMARY KEY NOT NULL,
  UserID BIGINT NOT NULL,
  RelateUserID BIGINT NOT NULL,
  State VARCHAR(16) NOT NULL
);
COMMENT ON TABLE public.Relationship IS '用户关系表';
COMMENT ON COLUMN public.Relationship.ID IS '主键ID';
COMMENT ON COLUMN public.Relationship.UserID IS '用户关系所属用户ID';
COMMENT ON COLUMN public.Relationship.RelateUserID IS '用户关系关联用户ID';
COMMENT ON COLUMN public.Relationship.State IS '用户类型';
CREATE INDEX idx_r ON public.Relationship (UserID);  -- 外键索引

ALTER TABLE public.Relationship ALTER COLUMN ID SET DEFAULT nextval('user_id_seq');


INSERT INTO public.Relationship VALUES
(1,1,2,'liked'),(2,2,1,'disliked');
