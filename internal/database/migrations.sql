CREATE TABLE IF NOT EXISTS "package" (
	"id" varchar PRIMARY KEY NOT NULL, -- ULID
	"package_id" varchar NOT NULL,
	"latest_version" boolean NOT NULL DEFAULT false,
	"visible" boolean NOT NULL DEFAULT true,
	"quality" integer NOT NULL,
	"repository_id" varchar NOT NULL,
	"price" varchar NOT NULL,
	"version" varchar NOT NULL,
	"architecture" varchar NOT NULL,
	"package_filename" varchar NOT NULL,
	"package_size" bigint NOT NULL,
	"sha256_hash" varchar,
	"name" varchar,
	"description" varchar,
	"author" varchar,
	"maintainer" varchar,
	"depiction" varchar,
	"native_depiction" varchar,
	"sileo_depiction" varchar,
	"header_url" varchar,
	"tint_color" varchar,
	"icon_url" varchar,
	"section" varchar,
	"tags" varchar[] DEFAULT ARRAY[]::varchar[],
	"installed_size" bigint
);

CREATE TABLE IF NOT EXISTS "repository" (
	"id" varchar PRIMARY KEY NOT NULL, -- Manifest ID
	"aliases" varchar[] DEFAULT ARRAY[]::varchar[],
	"visible" boolean NOT NULL DEFAULT true,
	"quality" integer NOT NULL,
	"package_count" bigint NOT NULL DEFAULT 0,
	"sections" varchar[] DEFAULT ARRAY[]::varchar[],
	"bootstrap" boolean NOT NULL DEFAULT false,
	"uri" varchar NOT NULL,
	"suite" varchar NOT NULL DEFAULT './',
	"component" varchar,
	"name" varchar,
	"version" varchar,
	"description" varchar,
	"date" timestamp,
	"payment_gateway" varchar,
	"sileo_endpoint" varchar,

	"origin_hostname" varchar NOT NULL,
	"origin_release_path" varchar NOT NULL,
	"origin_release_hash" varchar NOT NULL,
	"origin_packages_path" varchar NOT NULL,
	"origin_packages_hash" varchar NOT NULL,
	"origin_last_updated" timestamp NOT NULL,
	"origin_has_in_release" boolean NOT NULL,
	"origin_has_release_gpg" boolean NOT NULL,
	"origin_supports_payment_v1" boolean NOT NULL,
	"origin_supports_payment_v2" boolean NOT NULL,
	"origin_uses_https" boolean NOT NULL
);

CREATE INDEX IF NOT EXISTS visible_repository ON repository ("visible");
CREATE INDEX IF NOT EXISTS visible_package ON package ("visible");
CREATE INDEX IF NOT EXISTS quality_repository ON repository ("quality");
CREATE INDEX IF NOT EXISTS quality_package ON package ("quality");
CREATE INDEX IF NOT EXISTS latest_version ON package ("latest_version");
CREATE INDEX IF NOT EXISTS package_id ON package ("package_id");
CREATE UNIQUE INDEX IF NOT EXISTS repository_id ON repository ("id");
