CREATE TABLE IF NOT EXISTS Sessions(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER,
    FileName TEXT,
    LanguageName TEXT,
    ProjectName TEXT,
    StartTime INTEGER,
    StartDate TEXT,
    EndTime INTEGER,
    EndDate TEXT,
    FOREIGN KEY(UserID) REFERENCES users(ID)
)
