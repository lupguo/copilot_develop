-- auto-generated definition
create table already_updated_blogs
(
    id          integer primary key,
    title       text,
    path        text,
    keywords    text,
    description text,
    summary     text,
    headers     text,
    updated_at  text,
    deleted_at  text
);


select *
from already_updated_blogs;

delete
from already_updated_blogs;
