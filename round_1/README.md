<!--
EmojiFlics 🎬💀🍹
A fully containerized movie streaming playground where everything is emojis. No text in the interface.
Built entirely with Chainguard images for minimal, secure containers.
-->

# EmojiFlics 🎬💀🍹

<!-- ====================== -->
<!-- 🚀 Setup Instructions -->
<!-- ====================== -->

1. **Build all containers**:

```bash
# Builds and starts all services
docker-compose up -d --build
```

<!-- ====================== -->
<!-- 🧩 Architecture Overview -->
<!-- ====================== -->

We have four main services, each in its own container:

| Service  | Port | Role |
|--------- |------|------|
| inbound  | 9001 | Accepts movie content and forwards to gawk |
| gawk     | 9002 | Applies emoji filters, replaces shortcodes using AWK scripts |
| pandoc   | 9003 | Converts Markdown to HTML in `/output/index.html` |
| outbound | 8080 | Serves the final HTML to browser or curl |

**Why these components?**

- **Chainguard images:** Minimal, secure, and reproducible containers.
- **Multi-stage builds:** Used for pandoc and Python to separate build tools from runtime, keeping runtime images tiny and safe.
- **gawk:** Efficiently replaces emoji shortcodes like `:mage_man:` → `🧙‍♂️` and filters input.
- **pandoc:** Converts Markdown content into HTML, preserving emojis and formatting.
- **outbound:** Pure Python server that serves `/output` safely, no shell required.

<!-- ====================== -->
<!-- 🎬 Send Movie Content -->
<!-- ====================== -->

Use the helper script to star the game:

```bash
chmod +x play.sh
./play.sh
```

<!-- ====================== -->
<!-- 🌐 View Processed Movies -->
<!-- ====================== -->

- **Browser:**  
  Open [http://localhost:8080](http://localhost:8080)

- **CLI / curl:**  

```bash
curl http://localhost:8080
```

You should see only emojis, for example:

```
💀 🚢 🍹
🧙‍♂️ 💍 🔥
```

