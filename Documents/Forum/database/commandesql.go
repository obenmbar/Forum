package forumino
	// 4. Créer toutes les tables
	// On exécute une seule grande requête de création
	 const TablesSQL  string = `
	-- 1. Table des utilisateurs
	CREATE TABLE IF NOT EXISTS users (
		"id"       INTEGER PRIMARY KEY AUTOINCREMENT,
		"email"    TEXT NOT NULL UNIQUE,  -- UNIQUE: Empêche 2x le même email
		"username" TEXT NOT NULL UNIQUE,  -- UNIQUE: Empêche 2x le même pseudo
		"password" TEXT NOT NULL          -- Stockera le mot de passe HACHÉ
	);

	-- 2. Table des sessions (pour les cookies)
	CREATE TABLE IF NOT EXISTS sessions (
		"token"     TEXT PRIMARY KEY,  -- Le UUID aléatoire du cookie
		"user_id"   INTEGER NOT NULL,
		"expiry"    DATETIME NOT NULL,
		-- FOREIGN KEY: Crée un lien vers users.id
		-- ON DELETE CASCADE: Si un user est supprimé, ses sessions sont supprimées (nettoyage auto)
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- 3. Table des catégories (Sujets)
	CREATE TABLE IF NOT EXISTS categories (
		"id"   INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" TEXT NOT NULL UNIQUE
	);

	-- 4. Table des posts (Sujets du forum)
	CREATE TABLE IF NOT EXISTS posts (
		"id"         INTEGER PRIMARY KEY AUTOINCREMENT,
		"user_id"    INTEGER NOT NULL,
		"title"      TEXT NOT NULL,
		"content"    TEXT NOT NULL,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP, -- La DB met la date/heure automatiquement
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- 5. Table des commentaires (Réponses)
	CREATE TABLE IF NOT EXISTS comments (
		"id"         INTEGER PRIMARY KEY AUTOINCREMENT,
		"post_id"    INTEGER NOT NULL,
		"user_id"    INTEGER NOT NULL,
		"content"    TEXT NOT NULL,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- 6. Table de liaison (Post <-> Catégorie)
	-- C'est la solution "pro" pour un post qui a plusieurs catégories
	CREATE TABLE IF NOT EXISTS post_categories (
		"post_id"     INTEGER NOT NULL,
		"category_id" INTEGER NOT NULL,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE,
		-- PRIMARY KEY (composée): Empêche de lier 2x le même post à la même catégorie
		PRIMARY KEY (post_id, category_id)
	);

	-- 7. Table des votes sur les POSTS
	-- C'est la conception "propre" de la garder séparée de comment_votes
	CREATE TABLE IF NOT EXISTS post_votes (
		"user_id"   INTEGER NOT NULL,
		"post_id"   INTEGER NOT NULL,
		"vote_type" BOLEAN  NOT NULL, 
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
		-- PRIMARY KEY (composée): Garantit 1 seul vote par user par post
		PRIMARY KEY (user_id, post_id)
	);

	-- 8. Table des votes sur les COMMENTAIRES
	CREATE TABLE IF NOT EXISTS comment_votes (
		"user_id"    INTEGER NOT NULL,
		"comment_id" INTEGER NOT NULL,
		"vote_type"  BOLEAN NOT NULL ,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE,
		PRIMARY KEY (user_id, comment_id)
	);
	`