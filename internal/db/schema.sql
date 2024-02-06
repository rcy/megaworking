create table sessions(
       id integer not null primary key,
       created_at datetime not null default current_timestamp,

       status text not null default 'init',

       num_cycles integer not null,
       start_at datetime not null,
       -- zero for individual cycles, or cycles since origin for group cycles
       start_cycle_timer_id integer not null,

       -- prepare
       accomplish text not null default '',
       important text not null default '',
       complete text not null default '',
       distractions text not null default '',
       measurable text not null default '',
       noteworthy text not null default '',

       -- debrief
       target integer not null default 0,
       done text not null default '',
       compare text not null default '',
       bogged text not null default '',
       replicate text not null default '',
       takeaways text not null default '',
       nextsteps text not null default ''
);

create table cycles(
       id integer not null primary key,
       created_at datetime not null default current_timestamp,
       session_id integer references sessions(id) not null,

       -- for individual sessions, this starts at zero
       -- for group sessions, its the number of cycles since the origin time
       cycle_timer_id integer not null,

       -- plan
       accomplish text not null default '',
       started text not null default '',
       hazards text not null default '',
       energy integer not null default '',
       morale integer not null default '',

       -- review
       target integer not null default 0,
       noteworthy text not null default '',
       distractions text not null default '',
       improve text not null default ''
);
