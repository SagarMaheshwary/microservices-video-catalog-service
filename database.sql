CREATE TABLE IF NOT EXISTS "videos" (
  "id" BIGSERIAL PRIMARY KEY,
  "title" VARCHAR(250) NOT NULL,
  "description" TEXT NOT NULL,
  "original_id" VARCHAR(50) NOT NULL,
  "thumbnail_url" VARCHAR(50) NOT NULL,
  "user_id" BIGINT NOT NULL,
  "published_at" TIMESTAMP NOT NULL,
  "resolution" VARCHAR(15) NOT NULL,
  "duration" INT NOT NULL,
  "created_at" TIMESTAMP DEFAULT NOW(),
  "updated_at" TIMESTAMP NULl
);

CREATE TABLE IF NOT EXISTS "video_chunks" (
  "id" BIGSERIAL PRIMARY KEY,
  "video_id" BIGINT NOT NULL,
  "order" smallint NOT NULL,
  "resolution" VARCHAR(15) NOT NULL,
  "encoding" VARCHAR(15) NOT NULL,
  "url" VARCHAR(50) NOT NULL,
  "created_at" TIMESTAMP DEFAULT NOW(),
  "updated_at" TIMESTAMP NULL,
  FOREIGN KEY ("video_id") REFERENCES "videos"("id")
);
