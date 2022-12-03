package database

const QueryAuthLogin = `
with user_find as (
	select
		id as user_id, username, email, role
	from
		users u
	where 
	  ((not $1 and u.username = $2) or 
	   ($1 and u.email = $3)) and
	  u.password  = crypt($4, u.password )
), session_insert as (
	insert into 
	  users_sessions(user_id, addr, country)
	select 
	  uf.user_id, $5, $6
	from
		user_find uf
	on conflict(user_id, addr, revoked) where not revoked
		do update set updated = now()
	returning
		id as session_id
)
select
	session_insert.*,
	user_find.*
from session_insert, user_find`

const QueryAuthRefresh = `
with session_update as (
	update
		users_sessions
	set
		updated = now(),
		addr = $2
	where
		id = $1 and
		revoked = false
	returning
		id as session_id,
		user_id
)
select
	us.session_id, id as user_id, username, email, role
from
	users u, session_update us
where
	u.id = us.user_id`

const QueryUserCreate = `
insert into
	users(username, email, password)
values
	($1, $2, crypt($3, gen_salt('bf')))
returning
	id, username, email, role`

const QueryUserUpdate = `
update
	users
set
	username = coalesce($2::text, username),
	email = coalesce($3::text, email),
	password = case 
	    when $4::text is not null then crypt($4::text, gen_salt('bf'))
		else password
	end
where
	id = $1
returning
	id, username, email, role
`

const QueryUserGet = `
select
	u.id,
	u.username,
	u.role,
	(select count(1) filter (where up.user_id = u.id and not up.private) from users_pastes up) pastes,
	(select us.updated from users_sessions us where us.user_id = u.id order by us.updated desc limit 1) online
from
	users u
where 
	u.username = $1`

const QueryUserFind = `
select
	u.id,
	u.username,
	u.role,
	(select count(1) filter (where up.user_id = u.id and not up.private) from users_pastes up) pastes,
	(select us.updated from users_sessions us where us.user_id = u.id order by us.updated desc limit 1) online,
	count(1) over() as count
from
	users u
where 
	u.username ilike '%' || $1 || '%'
order by
	id asc
offset
	$2
limit
	$3`

const QuerySessionFind = `
select 
	id,
	country,
	created,
	updated,
	revoked,
	count(1) over() as count
from 
	users_sessions
where
	user_id  = $1 and 
	($2 or not revoked)
order by 
	updated desc
offset 
	$3
limit
	$4`

const QuerySessionRevoke = `
update
	users_sessions
set
	revoked = true
where 
	user_id = $1 and 
	id = $2 and
	revoked = false
returning 
	id,
	country,
	created,
	updated,
	revoked`

const QueryPasteCreate = `
insert into
	users_pastes(user_id, lang_name, title, "password", private, unlisted, content_text, content_bytes)
values (
	$1, $2, $3, $4::text is not null, $5, $6,
	(case when $4::text is null then $7 else null end),
	(case when $4::text is not null then pgp_sym_encrypt($7, $4::text) else null end)
)
returning
	id`

const QueryPasteFetch = `
select
	id,
	user_id as uid,
	title,
	private,
	unlisted,
	created,
	password,
	(case when "password" then pgp_sym_decrypt_null_on_err(content_bytes, $3::text) else content_text end) as content,
	(select to_jsonb(language_types) from language_types where name = lang_name) as language
from
	users_pastes
where
	id = $1 and 
	(not private or user_id = $2::int)`

const QueryPasteDelete = `
delete from
	users_pastes
where 
	user_id = $1 and id = $2
returning
	id`

const QueryPasteUpdate = `
update
	users_pastes 
set
	lang_name = coalesce($3::text, lang_name),
	title = coalesce($4::text, title),
	private = coalesce($5::boolean, private),
	unlisted = coalesce($6::boolean, unlisted),
	password = $7::text is not null,
	content_text = coalesce((case when $8::text is null then content_text end), (case when $7::text is null then content_text end)),
	content_bytes = coalesce((case when $8::text is null then content_bytes end), (case when $7::text is not null then pgp_sym_encrypt($8::text, $7::text) end))
where
	id  = $1 and
	user_id = $2
returning
	id`

const QueryPasteFind = `
select
	up.id,
	up.user_id as uid,
	u.username,
	up.title,
	up.password,
	up.private,
	up.unlisted,
	up.created,
	substring(up.content_text, 0, 50) as content,
	to_jsonb(lt) as language, 
	count(1) over() as count
from
	users_pastes up
inner join users u on
	u.id = up.user_id
inner join language_types lt on
	lt.name = up.lang_name
where
	(
		up.user_id = $1::int or
		u.username = $2::text or
		up.lang_name = $3::text or
		up.title ilike $4::text or
		up.password = $5::boolean or
		up.content_text ilike $6::text or
		up.created between coalesce($7::timestamp, to_timestamp(0)) and coalesce($8::timestamp, now())
	) 
		and
	(
		(not up.private or (up.user_id = $11::int and up.private = $9::boolean)) and
		(not up.unlisted or (up.user_id = $11::int and up.private = $10::boolean))
	)
order by
	up.created desc
offset 
	$12::int
limit
	$13::int`
