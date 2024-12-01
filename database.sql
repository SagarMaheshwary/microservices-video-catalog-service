CREATE TABLE IF NOT EXISTS "videos" (
  "id" BIGSERIAL PRIMARY KEY,
  "title" VARCHAR(250) NOT NULL,
  "description" TEXT NOT NULL,
  "original_id" VARCHAR(50) NOT NULL,
  "user_id" BIGINT NOT NULL,
  "resolution" VARCHAR(25) NOT NULL,
  "duration" INT NOT NULL,
  "path" VARCHAR(250) NOT NULL,
  "thumbnail" VARCHAR(250) NOT NULL,
  "published_at" TIMESTAMP NOT NULL,
  "created_at" TIMESTAMP DEFAULT NOW(),
  "updated_at" TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS "video_chunks" (
  "id" BIGSERIAL PRIMARY KEY,
  "video_id" BIGINT NOT NULL,
  "resolution" VARCHAR(25) NOT NULL,
  "filename" VARCHAR(250) NOT NULL,
  "created_at" TIMESTAMP DEFAULT NOW(),
  "updated_at" TIMESTAMP NULL,
  FOREIGN KEY ("video_id") REFERENCES "videos"("id")
);
