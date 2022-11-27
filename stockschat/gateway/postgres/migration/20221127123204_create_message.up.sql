create table message (
	author     text not null,
	created_at timestamp default now()::timestamp,
	content    text not null
);
