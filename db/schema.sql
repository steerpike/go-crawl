CREATE TABLE Artists (
    ID INTEGER PRIMARY KEY,
    Response INTEGER,
    Name TEXT,
    Url TEXT UNIQUE,
    Path TEXT UNIQUE,
    LastCrawled DATETIME
);

CREATE TABLE Seeds (
    ID INTEGER PRIMARY KEY,
    SourceUrl TEXT,
    Url TEXT UNIQUE
);

CREATE TABLE Tags (
    ID INTEGER PRIMARY KEY,
    Name TEXT UNIQUE
);

CREATE TABLE Videos (
    ID INTEGER PRIMARY KEY,
    Name TEXT,
    Url TEXT UNIQUE
);

CREATE TABLE Artist_Tags (
    ArtistUrl TEXT,
    TagName TEXT,
    FOREIGN KEY(ArtistUrl) REFERENCES Artists(Url),
    FOREIGN KEY(TagName) REFERENCES Tags(Name),
    UNIQUE(ArtistUrl, TagName)
);

CREATE TABLE Artist_Videos (
    ArtistUrl TEXT,
    VideoUrl TEXT,
    FOREIGN KEY(ArtistUrl) REFERENCES Artists(Url),
    FOREIGN KEY(VideoUrl) REFERENCES Videos(Url),
    UNIQUE(ArtistUrl, VideoUrl)
);

CREATE TABLE Similar_Artists (
    ArtistUrl1 TEXT,
    ArtistUrl2 TEXT,
    FOREIGN KEY(ArtistUrl1) REFERENCES Artists(Url),
    FOREIGN KEY(ArtistUrl2) REFERENCES Artists(Url),
    UNIQUE(ArtistUrl1, ArtistUrl2)
);
