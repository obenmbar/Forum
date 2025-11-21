// package database gère la connexion et l'initialisation de la DB
package forumino

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Le driver SQLite
)

// InitDB ouvre une connexion à la base de données et crée le schéma (les tables)
// si elles n'existent pas.
// Elle ne prend plus de paramètre et utilise "./forum.db" par défaut.
func InitDB() *sql.DB {
	// 1. Ouvrir la connexion à la base de données
	// Le nom du fichier est maintenant défini ici
	dataSourceName := "database/forum.db"
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		// Erreur fatale, l'application ne peut pas démarrer
		log.Fatal("Erreur sql.Open: ", err)
	}

	// 2. Tester la connexion (Ping)
	if err := db.Ping(); err != nil {
		log.Fatal("Erreur db.Ping: ", err)
	}

	// 3. Activer les clés étrangères (Foreign Keys)
	// C'est CRUCIAL pour que les 'ON DELETE CASCADE' fonctionnent.
	// SQLite le désactive par défaut.
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Erreur PRAGMA foreign_keys: ", err)
	}
	
	// Exécuter la grande requête
	_, err = db.Exec(TablesSQL)
	if err != nil {
		log.Fatal("Erreur création des tables: ", err)
	}
      err = INSERTcategory(db)
	  if err != nil {
		log.Fatal("insert category error",err)
	  }

	log.Println("Base de données et tables initialisées avec succès.")
	return db // On retourne la connexion à main.go
}