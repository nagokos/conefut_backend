# public.messages

## Description

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | varchar |  | false |  |  |  |
| content | varchar(1000) |  | false |  |  |  |
| room_id | varchar |  | false |  | [public.rooms](public.rooms.md) |  |
| user_id | varchar |  | true |  | [public.users](public.users.md) |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| messages_user_id_fkey | FOREIGN KEY | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL |
| messages_room_id_fkey | FOREIGN KEY | FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE |
| messages_pkey | PRIMARY KEY | PRIMARY KEY (id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| messages_pkey | CREATE UNIQUE INDEX messages_pkey ON public.messages USING btree (id) |
| messages_room_id_idx | CREATE INDEX messages_room_id_idx ON public.messages USING btree (room_id) |
| messages_user_id_idx | CREATE INDEX messages_user_id_idx ON public.messages USING btree (user_id) |

## Relations

![er](public.messages.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)