create table preparations(
       created_at datetime not null default current_timestamp,
       accomplish text not null,
       important text not null,
       complete text not null,
       distractions text not null,
       measurable text not null,
       noteworthy text not null
);
