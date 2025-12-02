-- +goose Up
CREATE TABLE analytics_url_daily (
 short_code TEXT NOT NULL,
 date DATE NOT NULL,
 total_clicks INT NOT NULL,
 unique_visitors INT NOT NULL,
 PRIMARY KEY (short_code, date)
);

CREATE TABLE analytics_url_geo (
short_code TEXT NOT NULL,
country TEXT NOT NULL,
city TEXT,
total_clicks INT NOT NULL,
PRIMARY KEY (short_code, country, city)
);

CREATE TABLE analytics_url_device (
  short_code TEXT NOT NULL,
  device TEXT NOT NULL,
  total_clicks INT NOT NULL,
  PRIMARY KEY (short_code, device)
);

CREATE TABLE analytics_url_browser (
   short_code TEXT NOT NULL,
   browser TEXT NOT NULL,
   total_clicks INT NOT NULL,
   PRIMARY KEY (short_code, browser)
);

CREATE TABLE analytics_url_referrer (
    short_code TEXT NOT NULL,
    referrer TEXT NOT NULL,
    total_clicks INT NOT NULL,
    PRIMARY KEY (short_code, referrer)
);

CREATE TABLE analytics_alias_daily (
   short_code TEXT NOT NULL,
   alias TEXT NOT NULL,
   date DATE NOT NULL,
   total_clicks INT NOT NULL,
   unique_visitors INT NOT NULL,
   PRIMARY KEY (short_code, alias, date)
);

CREATE TABLE analytics_alias_geo (
 short_code TEXT NOT NULL,
 alias TEXT NOT NULL,
 country TEXT NOT NULL,
 city TEXT,
 total_clicks INT NOT NULL,
 PRIMARY KEY (short_code, alias, country, city)
);

CREATE TABLE analytics_alias_device (
    short_code TEXT NOT NULL,
    alias TEXT NOT NULL,
    device TEXT NOT NULL,
    total_clicks INT NOT NULL,
    PRIMARY KEY (short_code, alias, device)
);

CREATE TABLE analytics_alias_browser (
     short_code TEXT NOT NULL,
     alias TEXT NOT NULL,
     browser TEXT NOT NULL,
     total_clicks INT NOT NULL,
     PRIMARY KEY (short_code, alias, browser)
);

CREATE TABLE analytics_alias_referrer (
      short_code TEXT NOT NULL,
      alias TEXT NOT NULL,
      referrer TEXT NOT NULL,
      total_clicks INT NOT NULL,
      PRIMARY KEY (short_code, alias, referrer)
);

CREATE TABLE analytics_user_daily (
  user_id UUID NOT NULL,
  date DATE NOT NULL,
  total_clicks INT NOT NULL,
  unique_visitors INT NOT NULL,
  PRIMARY KEY (user_id, date)
);

CREATE TABLE analytics_user_geo (
user_id UUID NOT NULL,
country TEXT NOT NULL,
total_clicks INT NOT NULL,
PRIMARY KEY (user_id, country)
);

CREATE TABLE analytics_user_browser (
    user_id UUID NOT NULL,
    browser TEXT NOT NULL,
    total_clicks INT NOT NULL,
    PRIMARY KEY (user_id, browser)
);

CREATE TABLE analytics_user_top_links (
      user_id UUID NOT NULL,
      short_code TEXT NOT NULL,
      total_clicks INT NOT NULL,
      PRIMARY KEY (user_id, short_code)
);

-- +goose Down
DROP TABLE IF EXISTS analytics_url_daily;
DROP TABLE IF EXISTS analytics_url_geo;
DROP TABLE IF EXISTS analytics_url_device;
DROP TABLE IF EXISTS analytics_url_browser;
DROP TABLE IF EXISTS analytics_url_referrer;

DROP TABLE IF EXISTS analytics_alias_daily;
DROP TABLE IF EXISTS analytics_alias_geo;
DROP TABLE IF EXISTS analytics_alias_device;
DROP TABLE IF EXISTS analytics_alias_browser;
DROP TABLE IF EXISTS analytics_alias_referrer;

DROP TABLE IF EXISTS analytics_user_daily;
DROP TABLE IF EXISTS analytics_user_geo;
DROP TABLE IF EXISTS analytics_user_browser;
DROP TABLE IF EXISTS analytics_user_top_links;
