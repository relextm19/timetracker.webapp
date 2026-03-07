CREATE TABLE IF NOT EXISTS Sessions(
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserToken TEXT,
    FileName TEXT,
    LanguageName TEXT,
    ProjectName TEXT,
    StartTime INTEGER,
    StartDate TEXT,
    EndTime INTEGER,
    EndDate TEXT
)
