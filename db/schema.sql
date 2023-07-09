CREATE TABLE Artists (
    ID INTEGER PRIMARY KEY,
    Response INTEGER,
    Name TEXT,
    Url TEXT UNIQUE,
    Path TEXT UNIQUE,
    LastCrawled DATETIME
);

CREATE TABLE Tags (
    ID INTEGER PRIMARY KEY,
    TagName TEXT UNIQUE
);

CREATE TABLE Videos (
    ID INTEGER PRIMARY KEY,
    VideoName TEXT,
    VideoUrl TEXT UNIQUE
);

CREATE TABLE Artist_Tags (
    ArtistID INTEGER,
    TagID INTEGER,
    FOREIGN KEY(ArtistID) REFERENCES Artists(ID),
    FOREIGN KEY(TagID) REFERENCES Tags(ID),
    UNIQUE(ArtistID, TagID)
);

CREATE TABLE Artist_Videos (
    ArtistID INTEGER,
    VideoID INTEGER,
    FOREIGN KEY(ArtistID) REFERENCES Artists(ID),
    FOREIGN KEY(VideoID) REFERENCES Videos(ID),
    UNIQUE(ArtistID, VideoID)
);

CREATE TABLE Similar_Artists (
    ArtistID1 INTEGER,
    ArtistID2 INTEGER,
    FOREIGN KEY(ArtistID1) REFERENCES Artists(ID),
    FOREIGN KEY(ArtistID2) REFERENCES Artists(ID),
    UNIQUE(ArtistID1, ArtistID2)
);