# filter.awk
@include "emoji_db.awk"
{
    # Reemplazo de shortcodes
    do {
        found = 0
        # buscar todos los patrones :alias:
        while (match($0, /:([a-z0-9_+-]+):/, m)) {
            key = ":" m[1] ":"
            if (key in map) {
                $0 = substr($0, 1, RSTART-1) map[key] substr($0, RSTART+RLENGTH)
                found = 1
                break   # salir del while para reiniciar do-while y evitar parcial match
            } else {
                # saltar y seguir buscando el siguiente
                $0 = substr($0, 1, RSTART+RLENGTH-1) substr($0, RSTART+RLENGTH)
            }
        }
    } while(found)
    print
}

