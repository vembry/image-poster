-- CREATE DATABASE image_poster;

\c image_poster;

CREATE TABLE public.posts (
	id char(27) NOT NULL,
	"text" text NULL,
	image jsonb NULL,
	created_by varchar NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	CONSTRAINT posts_pk PRIMARY KEY (id)
);

CREATE TABLE public.post_structures (
	post_id char(27) NOT NULL,
	parent_post_id char(27) NULL,
	CONSTRAINT post_structures_post_id_posts_fk FOREIGN KEY (post_id) REFERENCES public.posts(id),
	CONSTRAINT post_structures_partner_post_id_posts_fk FOREIGN KEY (parent_post_id) REFERENCES public.posts(id)
);
