create table if not exists campaigns (
    "id" bigserial primary key,
    "campaign_id" text not null,
    "sensors" text[] not null,
    "type" text not null,
    "begin" timestamp not null,
    "end" timestamp not null,
    "created_at" timestamp default now()
);