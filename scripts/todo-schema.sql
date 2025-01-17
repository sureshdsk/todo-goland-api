CREATE TABLE IF NOT EXISTS todo_items(
   id serial primary key,
   task varchar(255) not null,
   status varchar(50) not null,
   created_at timestamp default CURRENT_TIMESTAMP
);