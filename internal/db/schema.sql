create table sessions(
       id integer not null primary key,
       created_at datetime not null default current_timestamp,

       state text not null default 'init',

       num_cycles integer not null,
       start_at datetime not null,

       -- prepare
       accomplish text not null default '',
       important text not null default '',
       complete text not null default '',
       distractions text not null default '',
       measurable text not null default '',
       noteworthy text not null default ''

       -- debrief
       -- target integer not null,
       -- done text not null,
       -- compare text not null,
       -- bogged text not null,
       -- replicate text not null,
       -- takeaways text not null,
       -- nextsteps text not null,
);

create table cycles(
       id integer not null primary key,
       created_at datetime not null default current_timestamp,
       session_id integer references sessions(id) not null,

       -- plan
       accomplish text not null,
       started text not null,
       hazards text not null,
       energy integer not null,
       morale integer not null

       -- review
       -- target integer not null,
       -- noteworthy text not null,
       -- distractions text not null,
       -- improve text not null
);
