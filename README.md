# Forum B1 — Documentation complète

Projet de forum web réalisé en Go, SQLite, HTML/CSS pur, containerisé avec Docker.  
Ce document résume l'intégralité de ce qui a été construit, modifié et pourquoi.

---

## Sommaire

1. [Lancement rapide](#lancement-rapide)
2. [Architecture du projet](#architecture-du-projet)
3. [Schéma de la base de données](#schéma-de-la-base-de-données)
4. [Ce qui a été fait — étape par étape](#ce-qui-a-été-fait--étape-par-étape)
5. [Explication des mécanismes clés](#explication-des-mécanismes-clés)
6. [Structure des fichiers](#structure-des-fichiers)
7. [Commandes utiles](#commandes-utiles)

---

## Lancement rapide

### Avec Docker (recommandé)

```bash
docker compose up -d
```

Le forum est accessible sur **http://localhost:8080**

Pour arrêter :
```bash
docker compose down
```

### Sans Docker (développement local)

```bash
go run .
```

> Prérequis : Go 1.26+

---

## Architecture du projet

```
Forum (Go)
├── main.go              → Point d'entrée, définition de toutes les routes HTTP
├── database/            → Toutes les interactions avec SQLite
│   ├── database.go      → Connexion + création des tables au démarrage
│   ├── users.go         → Création et lecture des utilisateurs
│   ├── sessions.go      → Gestion des sessions (UUID)
│   ├── posts.go         → CRUD des posts
│   ├── categories.go    → Lecture des catégories
│   ├── comments.go      → CRUD des commentaires
│   └── likes.go         → Toggle like/dislike
├── handlers/            → Handlers HTTP (logique de chaque route)
│   ├── home.go          → Page d'accueil + RenderTemplate
│   ├── login.go         → Connexion / Déconnexion
│   ├── register.go      → Inscription
│   ├── postHandler.go   → Affichage d'un post + création de commentaire
│   ├── postCreator.go   → Création d'un post (avec image)
│   ├── postLikeHandler.go    → Like/Dislike des posts
│   ├── commentLikeHandler.go → Like/Dislike des commentaires
│   ├── categoriesHandler.go  → Liste des catégories + filtrage
│   ├── errorHandlers.go      → Pages 404 et 500 custom
│   └── randomPageHandler.go  → Page aléatoire
├── models/              → Structures de données Go (structs)
│   ├── user.go          → Struct User
│   ├── post.go          → Struct Post
│   ├── comment.go       → Struct Comments
│   ├── category.go      → Struct Category
│   └── template.go      → Struct TemplateData (données envoyées aux templates)
├── fake/                → Couche d'abstraction (délègue à database/)
│   ├── users.go         → GetCurrentUser, GetCurrentUserFull
│   ├── posts.go         → GetAllPosts, GetPostById, GetPostsByCategory
│   ├── categories.go    → GetAllCategories, GetCategoryById
│   └── comments.go      → GetCommentByPostID, GetCommentByID
├── templates/           → Fichiers HTML (moteur de template Go)
│   ├── index.tmpl       → Page d'accueil (liste des posts)
│   ├── post.tmpl        → Page d'un post (contenu + commentaires + likes)
│   ├── postcreate.tmpl  → Formulaire de création de post
│   ├── login.tmpl       → Formulaire de connexion
│   ├── register.tmpl    → Formulaire d'inscription
│   ├── categories.tmpl  → Liste des catégories
│   ├── postByCategory.tmpl → Posts filtrés par catégorie
│   ├── 404.tmpl         → Page d'erreur 404
│   └── 500.tmpl         → Page d'erreur 500
├── static/              → Fichiers statiques servis directement
│   ├── style.css        → Feuille de style
│   └── upload/          → Images uploadées par les utilisateurs
├── data/                → Dossier créé automatiquement au démarrage
│   └── forum.db         → Base de données SQLite (générée au 1er lancement)
├── Dockerfile           → Instructions de build de l'image Docker
├── docker-compose.yml   → Orchestration Docker avec volumes persistants
└── .dockerignore        → Fichiers exclus du build Docker
```

---

## Schéma de la base de données

```
users
├── id              INTEGER PRIMARY KEY AUTOINCREMENT
├── email           TEXT UNIQUE NOT NULL
├── username        TEXT NOT NULL
├── password_hash   TEXT NOT NULL          ← mot de passe hashé bcrypt
└── created_at      DATETIME DEFAULT CURRENT_TIMESTAMP

sessions
├── id              TEXT PRIMARY KEY       ← UUID généré à la connexion
├── user_id         INTEGER NOT NULL       ← FK → users.id
└── expires_at      DATETIME NOT NULL      ← expiration 24h après création

categories
├── id              INTEGER PRIMARY KEY AUTOINCREMENT
└── name            TEXT UNIQUE NOT NULL   ← pré-remplies au démarrage

posts
├── id              INTEGER PRIMARY KEY AUTOINCREMENT
├── title           TEXT NOT NULL
├── content         TEXT NOT NULL
├── image_path      TEXT DEFAULT ''        ← chemin web vers l'image (/static/upload/...)
├── author_id       INTEGER NOT NULL       ← FK → users.id
└── created_at      DATETIME DEFAULT CURRENT_TIMESTAMP

post_categories                            ← table de liaison many-to-many
├── post_id         INTEGER NOT NULL       ← FK → posts.id
└── category_id     INTEGER NOT NULL       ← FK → categories.id
    PRIMARY KEY (post_id, category_id)

comments
├── id              INTEGER PRIMARY KEY AUTOINCREMENT
├── content         TEXT NOT NULL
├── author_id       INTEGER NOT NULL       ← FK → users.id
├── post_id         INTEGER NOT NULL       ← FK → posts.id
└── created_at      DATETIME DEFAULT CURRENT_TIMESTAMP

likes
├── id              INTEGER PRIMARY KEY AUTOINCREMENT
├── user_id         INTEGER NOT NULL       ← FK → users.id
├── post_id         INTEGER                ← NULL si c'est un like de commentaire
├── comment_id      INTEGER                ← NULL si c'est un like de post
└── value           INTEGER NOT NULL       ← +1 (like) ou -1 (dislike)
```

**Relation many-to-many posts ↔ catégories :**  
Un post peut avoir plusieurs catégories, une catégorie peut avoir plusieurs posts.  
La table `post_categories` stocke les associations. Ex : le post 3 est dans les catégories 1 et 4.

---

## Ce qui a été fait — étape par étape

### Étape 1 — Dépendances Go

**Fichiers modifiés :** `go.mod`, `go.sum`

Trois bibliothèques ajoutées via `go get` :

| Bibliothèque | Rôle |
|---|---|
| `modernc.org/sqlite` | Driver SQLite **pur Go** — pas besoin de GCC/CGo |
| `golang.org/x/crypto/bcrypt` | Hashage sécurisé des mots de passe |
| `github.com/google/uuid` | Génération d'identifiants de session uniques |

> **Pourquoi `modernc.org/sqlite` et pas `mattn/go-sqlite3` ?**  
> `mattn/go-sqlite3` est un wrapper C qui nécessite GCC pour compiler — difficile à installer sur Windows et dans Docker. `modernc.org/sqlite` est une traduction automatique du code C de SQLite en Go pur. Même comportement, zéro dépendance externe.

---

### Étape 2 — Base de données SQLite

**Fichiers créés :** `database/database.go`  
**Fichiers modifiés :** `main.go`

`database.Init()` est appelé au démarrage du serveur. Il :
1. Crée le dossier `./data/` s'il n'existe pas
2. Ouvre (ou crée) `./data/forum.db`
3. Exécute tous les `CREATE TABLE IF NOT EXISTS` — si les tables existent déjà, rien n'est écrasé
4. Insère les catégories par défaut avec `INSERT OR IGNORE` (même principe)

---

### Étape 3 — Authentification réelle

**Fichiers créés :** `models/user.go`, `database/users.go`, `database/sessions.go`  
**Fichiers modifiés :** `fake/users.go`, `handlers/register.go`, `handlers/login.go`

Remplacement de l'authentification fake (email hardcodé, cookie à valeur fixe) par un système complet.

**Inscription (`/register` POST) :**
1. Validation des champs (email, username, password non vides)
2. `bcrypt.GenerateFromPassword` → hash du mot de passe
3. `INSERT INTO users` avec le hash (jamais le mot de passe en clair)
4. Si email déjà pris → erreur `ErrEmailTaken` (contrainte UNIQUE SQLite)
5. Création immédiate d'une session → cookie posé → utilisateur connecté

**Connexion (`/login` POST) :**
1. `SELECT` de l'utilisateur par email
2. `bcrypt.CompareHashAndPassword` → comparaison hash sans jamais décoder le mot de passe
3. Si OK → création d'une session UUID → cookie posé
4. Si KO → redirect `/login?error=1` (délibérément sans préciser si c'est l'email ou le mot de passe qui est faux — sécurité)

**Déconnexion (`/logout` POST) :**
1. Lecture du cookie de session
2. `DELETE FROM sessions` → la session est invalidée côté serveur
3. Cookie supprimé côté navigateur (`MaxAge: -1`)

**Vérification d'une session (sur chaque requête) :**
```
Cookie "session_id" (UUID) → SELECT sessions JOIN users WHERE id = ? AND expires_at > NOW()
```
Si la session est expirée ou inexistante → utilisateur considéré comme non connecté.

---

### Étape 4 — Posts

**Fichiers créés :** `database/posts.go`, `database/categories.go`  
**Fichiers modifiés :** `models/post.go`, `fake/posts.go`, `fake/categories.go`, `handlers/postCreator.go`

**Récupération des posts (requête SQL avec JOIN) :**

La requête utilise plusieurs techniques SQL avancées :

```sql
SELECT p.id, p.title, p.content, p.image_path, u.id, u.username, p.created_at,
    COALESCE(GROUP_CONCAT(DISTINCT c.name, ', '), '') AS category,
    COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
    COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
FROM posts p
JOIN users u ON u.id = p.author_id
LEFT JOIN post_categories pc ON pc.post_id = p.id
LEFT JOIN categories c ON c.id = pc.category_id
LEFT JOIN likes l ON l.post_id = p.id
GROUP BY p.id
ORDER BY p.created_at DESC
```

- **`GROUP_CONCAT(DISTINCT c.name, ', ')`** : concatène toutes les catégories d'un post en une seule chaîne (`"Général, Jeux vidéos"`) en une seule requête
- **`SUM(CASE WHEN l.value = 1 ...)`** : compte les likes et dislikes en même temps
- **`COALESCE(..., 0)`** : retourne 0 si NULL (posts sans like)
- **`LEFT JOIN`** : garde les posts même s'ils n'ont pas de catégorie ou de like
- **`GROUP BY p.id`** : regroupe les lignes par post (nécessaire avec GROUP_CONCAT et SUM)

**Création d'un post :** utilise une **transaction SQL** — si l'insertion dans `post_categories` échoue, l'insertion dans `posts` est annulée automatiquement (`tx.Rollback()`). Évite d'avoir des posts sans catégorie en BDD.

**Validation des images par magic bytes :**  
On ne se fie pas à l'extension du fichier (n'importe qui peut renommer `virus.exe` en `photo.jpg`). On lit les premiers octets du fichier et on compare aux signatures connues :

| Format | Octets de début (hex) |
|---|---|
| JPEG | `FF D8 FF` |
| PNG | `89 50 4E 47` |
| GIF | `47 49 46` |

---

### Étape 5 — Commentaires

**Fichiers créés :** `database/comments.go`  
**Fichiers modifiés :** `models/comment.go`, `fake/comments.go`, `handlers/postHandler.go`

Les commentaires sont récupérés avec le même pattern JOIN que les posts (pour avoir le username de l'auteur et les compteurs de likes). `CreateCommentHandler` utilise désormais `GetCurrentUserFull` pour récupérer l'`ID` de l'auteur (et non juste son username) afin de l'insérer en BDD.

---

### Étape 6 — Likes / Dislikes

**Fichiers créés :** `database/likes.go`  
**Fichiers modifiés :** `handlers/postLikeHandler.go`, `handlers/commentLikeHandler.go`

Système de **toggle** : un seul enregistrement par combinaison `(user, post)` ou `(user, comment)` dans la table `likes`.

Logique de `TogglePostLike(userID, postID, value int)` :

```
SELECT value FROM likes WHERE user_id = ? AND post_id = ? AND comment_id IS NULL
│
├── Pas de résultat (ErrNoRows) → INSERT (premier like/dislike)
│
├── Même value que ce qu'on clique → DELETE (annulation / toggle off)
│
└── Value différente → UPDATE (bascule like ↔ dislike)
```

La condition `AND comment_id IS NULL` garantit qu'on ne confond pas un like de post avec un like de commentaire même si les IDs numériques sont identiques.

---

### Étape 7 — Pages d'erreur HTTP

**Fichiers créés :** `templates/404.tmpl`, `templates/500.tmpl`, `handlers/errorHandlers.go`  
**Fichiers modifiés :** `handlers/home.go`, `main.go`

- `NotFoundHandler` : enregistré sur la route `/` (catch-all du ServeMux Go) — attrape toutes les URLs qui ne correspondent à aucune autre route
- `InternalErrorHandler` : appelé dans `RenderTemplate` en cas d'erreur de template, et disponible pour tout handler qui en a besoin
- `w.WriteHeader(404)` est écrit **avant** le body — en HTTP, le status code doit précéder le contenu, sinon Go envoie un 200 par défaut

---

### Étape 8 — Docker

**Fichiers créés :** `Dockerfile`, `docker-compose.yml`, `.dockerignore`  
**Fichiers modifiés :** `database/database.go`

**Dockerfile — multi-stage build :**

```dockerfile
# Stage 1 : compilation
FROM golang:1.26-alpine AS builder
# → compile le binaire Go

# Stage 2 : image finale
FROM debian:bookworm-slim
# → copie uniquement le binaire + templates + static
```

Le multi-stage build sert à réduire la taille de l'image finale. L'image de build contient Go, les headers, les outils de compilation (~500 Mo). L'image finale ne contient que le binaire compilé + les fichiers statiques (~100 Mo). L'image de build est jetée.

**Volumes Docker (persistance des données) :**

```yaml
volumes:
  - forum_db:/app/data        # persiste forum.db entre les redémarrages
  - forum_uploads:/app/static/upload  # persiste les images uploadées
```

Sans volumes, toutes les données seraient perdues à chaque `docker compose down`.

> **Problème rencontré :** monter un volume sur un fichier (`/app/forum.db`) crée un **dossier** du même nom, pas un fichier. SQLite ne peut pas ouvrir un dossier. **Solution :** la BDD est dans `./data/forum.db` et le volume monte sur le dossier `./data/`.

---

## Explication des mécanismes clés

### Comment fonctionne l'authentification par cookie/session

```
[Navigateur]                          [Serveur Go]                    [SQLite]
     │                                      │                              │
     │── POST /login (email + password) ───►│                              │
     │                                      │── SELECT user by email ─────►│
     │                                      │◄─ user{id, hash} ────────────│
     │                                      │                              │
     │                                      │ bcrypt.Compare(hash, password)
     │                                      │                              │
     │                                      │── INSERT sessions(uuid) ────►│
     │◄── Set-Cookie: session_id=<uuid> ───│                              │
     │                                      │                              │
     │── GET / (Cookie: session_id=<uuid>) ►│                              │
     │                                      │── SELECT sessions JOIN users►│
     │                                      │◄─ user{id, username} ────────│
     │◄── HTML (connecté) ─────────────────│                              │
```

### Comment fonctionne le système de like

```
Clic "👍" sur un post
        │
        ▼
togglePostReaction(userID=5, postID=12, value=+1)
        │
        ├─► SELECT value FROM likes WHERE user_id=5 AND post_id=12 AND comment_id IS NULL
        │
        │   ┌─ Aucun résultat ──► INSERT (like = +1)
        │   │
        │   ├─ value = +1 ───────► DELETE (toggle off, annulation)
        │   │
        │   └─ value = -1 ───────► UPDATE SET value = +1 (bascule dislike → like)
        │
        ▼
Redirect vers /posts/12 (la page se recharge avec les nouveaux compteurs)
```

### Comment fonctionne RenderTemplate

```go
func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
    // 1. Parse le fichier .tmpl
    t, err := template.ParseFiles("./templates/" + tmpl)
    if err != nil {
        InternalErrorHandler(w, nil)  // → affiche 500.tmpl
        return
    }
    // 2. Injecte les données Go dans le template HTML
    if err = t.Execute(w, data); err != nil {
        InternalErrorHandler(w, nil)
    }
}
```

Les templates Go utilisent `{{ .NomDuChamp }}` pour afficher les données, `{{ if .IsLogged }}` pour les conditions, `{{ range .Posts }}` pour les boucles.

---

## Structure des fichiers

```
JS-ForumB1-Front/
├── data/
│   └── forum.db              ← créé automatiquement au 1er lancement
├── database/
│   ├── database.go
│   ├── users.go
│   ├── sessions.go
│   ├── posts.go
│   ├── categories.go
│   ├── comments.go
│   └── likes.go
├── fake/
│   ├── users.go
│   ├── posts.go
│   ├── categories.go
│   └── comments.go
├── handlers/
│   ├── home.go
│   ├── login.go
│   ├── register.go
│   ├── postHandler.go
│   ├── postCreator.go
│   ├── postLikeHandler.go
│   ├── commentLikeHandler.go
│   ├── categoriesHandler.go
│   ├── errorHandlers.go
│   └── randomPageHandler.go
├── models/
│   ├── user.go
│   ├── post.go
│   ├── comment.go
│   ├── category.go
│   └── template.go
├── static/
│   ├── style.css
│   └── upload/               ← images des posts
├── templates/
│   ├── index.tmpl
│   ├── post.tmpl
│   ├── postcreate.tmpl
│   ├── login.tmpl
│   ├── register.tmpl
│   ├── categories.tmpl
│   ├── postByCategory.tmpl
│   ├── 404.tmpl
│   └── 500.tmpl
├── .dockerignore
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── main.go
```

---

## Commandes utiles

```bash
# Lancer avec Docker
docker compose up -d

# Relancer après modification du code
docker compose up -d --build

# Voir les logs en temps réel
docker compose logs -f

# Arrêter le container
docker compose down

# Arrêter ET supprimer les volumes (remet la BDD à zéro)
docker compose down -v

# Lancer en local (sans Docker)
go run .

# Compiler pour vérifier les erreurs
go build ./...
```

---

## Fonctionnalités implémentées

| Fonctionnalité | Statut |
|---|---|
| SQLite avec schéma complet | ✅ |
| Inscription avec bcrypt | ✅ |
| Connexion / Déconnexion | ✅ |
| Sessions UUID avec expiration 24h | ✅ |
| Création de posts avec catégories multiples | ✅ |
| Upload d'image (JPEG, PNG, GIF — max 20 Mo) | ✅ |
| Validation d'image par magic bytes | ✅ |
| Commentaires | ✅ |
| Like / Dislike avec toggle (posts + commentaires) | ✅ |
| Filtrage des posts par catégorie | ✅ |
| Pages d'erreur HTTP 404 et 500 custom | ✅ |
| Docker (multi-stage build + volumes persistants) | ✅ |
| Bonus — OAuth Google/GitHub | ❌ |
| Bonus — Système de modération | ❌ |
| Bonus — HTTPS + Rate Limiting | ❌ |
| Bonus — Notifications + Page d'activité | ❌ |
