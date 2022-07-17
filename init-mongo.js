db.createUser({
  user: "root",
  pwd: "root",
  roles: [
    { role: "readAnyDatabase", db: "flashscore" },
    { role: "dbAdminAnyDatabase", db: "flashscore" },
    { role: "userAdminAnyDatabase", db: "flashscore" },
  ],
});

let db = connect("localhost:27017/flashscore");
db.auth("root", "root");

db.createCollection("apikey");
