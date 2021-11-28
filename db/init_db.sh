mkdir /data
# log db
sqlite3 /data/go_note.db < /init_db/init_note_db.sql
# log tracking db
sqlite3 /data/log_track.db < /init_db/init_log_note_db.sql
