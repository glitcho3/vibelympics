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

Use the helper script to play a movie (e.g., YouTube):

```bash
# Example: play a YouTube link
./play yt <YOUTUBE_URL>
```

Content flows: `inbound → gawk → pandoc → outbound`.

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

<!-- ====================== -->
<!-- ✅ Recommended Order -->
<!-- ====================== -->

1. `docker-compose up -d --build`  
2. `./play yt <URL>`  
3. View with `curl` or browser.

<!-- ====================== -->
<!-- Optional Notes -->
<!-- ====================== -->

- Outbound container is pure Python, does not need shell or root permissions.
- Pandoc uses a multi-stage build to keep the runtime minimal.
- gawk scripts handle emoji translation and filtering efficiently.
- Volumes are shared for `/output` so that pandoc writes and outbound reads the HTML.
- No text is used in UI, everything is emojis.
- Safe to test locally using Docker Compose.

