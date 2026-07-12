# klickops starter

Ein solider Startpunkt für deine eigene Web-App: **SvelteKit-Frontend und
Go-Backend in einem einzigen Container**, mit optionaler **PostgreSQL-
Datenbank** und **S3-Speicher**, genau so, wie [klickops](https://klickops.io)
diese Dienste bereitstellt.

Du musst kein Entwickler sein. Dieses Projekt ist dafür gebaut, mit
KI-Werkzeugen wie [Claude Code](https://claude.com/claude-code) oder Cursor
weiterentwickelt zu werden: Die Datei `CLAUDE.md` und die Regeln in
`.claude/` erklären der KI, wie das Projekt aufgebaut ist und wie sauberer
Code hier aussieht.

*English version: [README.md](README.md)*

## So startest du

1. Klicke auf GitHub auf **"Use this template"** und erstelle dein eigenes
   Repository.
2. Öffne den Projektordner mit Claude Code (oder einem ähnlichen Werkzeug).
3. Beschreibe, was du bauen willst, zum Beispiel:

   > „Ersetze die Notiz-Demo durch eine Kundenliste mit Name, E-Mail und
   > Notizfeld."

   Die KI kennt das Projekt und baut die Funktion nach den vorhandenen
   Mustern, inklusive Tests.

## Lokal ausprobieren (optional)

Am einfachsten ohne Installation: Öffne dein Repository in GitHub
Codespaces; die mitgelieferte Entwicklungsumgebung bringt alles mit,
inklusive Claude Code (tippe `claude` ins Terminal und leg los).
Oder lokal (benötigt Go, Node/pnpm und Docker):

```bash
make install         # Abhängigkeiten installieren
make services-up     # lokale Datenbank + Speicher starten
cp .env.example .env
make dev             # App auf http://localhost:5173
```

## Auf klickops veröffentlichen

1. Repository zu GitHub pushen.
2. In klickops ein Projekt erstellen und **"Deploy from repo"** wählen.
   klickops erkennt das Dockerfile, baut die App und stellt sie mit einer
   URL und Zertifikat bereit.
3. Eine **PostgreSQL**-Datenbank als Service hinzufügen; klickops schlägt
   die Verbindung (`DATABASE_URL`) automatisch vor.
4. Einen **Bucket** (Dateispeicher) hinzufügen und die Zugangsdaten als
   `S3_*`-Variablen an die App binden (Details in der englischen README).

Fertig. Kein YAML, kein Kubernetes-Wissen nötig.

## Lizenz

[MIT](LICENSE). Viel Erfolg beim Bauen!
