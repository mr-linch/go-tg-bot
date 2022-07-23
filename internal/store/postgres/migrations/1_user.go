package migrations

func init() {
	include(1, query(`
		create table "user" (
			"id" serial primary key,

			"telegram_id" bigint not null,
			"telegram_username" text unique,

			"first_name" text not null,
			"last_name" text,
			"language_code" text,
			"created_at" timestamp not null,
			"updated_at" timestamp
		);
	`), query(`
		drop table "user";
	`))
}
