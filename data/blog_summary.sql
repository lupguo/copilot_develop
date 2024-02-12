create table main.blog_articles
(
    id          integer not null
        primary key autoincrement,
    title       text,
    path        text,
    keywords    text,
    description text,
    summary     text,
    draft       integer,
    weight      integer,
    word_count  integer,
    tags        text,
    categories  text,
    aliases     text,
    short_mark  text,
    date        text,
    updated_at  text,
    deleted_at  text,
    created_at  text    not null
);

create index main.blog_articles_path_index
    on main.blog_articles (path);
