{\rtf1\ansi\ansicpg1251\cocoartf2822
\cocoatextscaling0\cocoaplatform0{\fonttbl\f0\fswiss\fcharset0 Helvetica;}
{\colortbl;\red255\green255\blue255;}
{\*\expandedcolortbl;;}
\paperw11900\paperh16840\margl1440\margr1440\vieww11520\viewh8400\viewkind0
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\f0\fs24 \cf0 # \uc0\u1040 \u1088 \u1093 \u1080 \u1090 \u1077 \u1082 \u1090 \u1091 \u1088 \u1072  \u1080  \u1090 \u1077 \u1093 \u1085 \u1086 \u1083 \u1086 \u1075 \u1080 \u1095 \u1077 \u1089 \u1082 \u1080 \u1081  \u1089 \u1090 \u1077 \u1082  \u1087 \u1088 \u1086 \u1077 \u1082 \u1090 \u1072  ADR-P\
\
## 1. \uc0\u1054 \u1087 \u1080 \u1089 \u1072 \u1085 \u1080 \u1077  \u1089 \u1080 \u1089 \u1090 \u1077 \u1084 \u1099 \
ADR-P (Abuse Detection & Response Platform) \'97 \uc0\u1101 \u1090 \u1086  \u1084 \u1080 \u1082 \u1088 \u1086 \u1089 \u1077 \u1088 \u1074 \u1080 \u1089 \u1085 \u1072 \u1103  \u1087 \u1083 \u1072 \u1090 \u1092 \u1086 \u1088 \u1084 \u1072  \u1076 \u1083 \u1103  \u1089 \u1073 \u1086 \u1088 \u1072  \u1090 \u1077 \u1083 \u1077 \u1084 \u1077 \u1090 \u1088 \u1080 \u1080  \u1080 \u1079  \u1074 \u1085 \u1077 \u1096 \u1085 \u1080 \u1093  \u1080 \u1089 \u1090 \u1086 \u1095 \u1085 \u1080 \u1082 \u1086 \u1074  (Telegram/VK), \u1076 \u1077 \u1090 \u1077 \u1082 \u1094 \u1080 \u1080  \u1072 \u1085 \u1086 \u1084 \u1072 \u1083 \u1080 \u1081  \u1074  \u1088 \u1077 \u1072 \u1083 \u1100 \u1085 \u1086 \u1084  \u1074 \u1088 \u1077 \u1084 \u1077 \u1085 \u1080  \u1080  \u1072 \u1074 \u1090 \u1086 \u1084 \u1072 \u1090 \u1080 \u1079 \u1080 \u1088 \u1086 \u1074 \u1072 \u1085 \u1085 \u1086 \u1075 \u1086  \u1088 \u1077 \u1072 \u1075 \u1080 \u1088 \u1086 \u1074 \u1072 \u1085 \u1080 \u1103 . \u1057 \u1080 \u1089 \u1090 \u1077 \u1084 \u1072  \u1087 \u1088 \u1077 \u1076 \u1085 \u1072 \u1079 \u1085 \u1072 \u1095 \u1077 \u1085 \u1072  \u1076 \u1083 \u1103  \u1076 \u1077 \u1084 \u1086 \u1085 \u1089 \u1090 \u1088 \u1072 \u1094 \u1080 \u1080  \u1085 \u1072 \u1074 \u1099 \u1082 \u1086 \u1074  L2-\u1072 \u1085 \u1090 \u1080 \u1092 \u1088 \u1086 \u1076 -\u1072 \u1085 \u1072 \u1083 \u1080 \u1090 \u1080 \u1082 \u1072  \u1080  \u1080 \u1085 \u1078 \u1077 \u1085 \u1077 \u1088 \u1072 .\
\
## 2. \uc0\u1062 \u1077 \u1083 \u1077 \u1074 \u1086 \u1081  \u1089 \u1090 \u1077 \u1082 \
- **Ingestion & API:** Go 1.22+, Gin/Echo, Redis (Stream/List as queue).\
- **Detection Engine:** Python 3.12+, Pandas, Scikit-learn/CatBoost, SQLAlchemy.\
- **Storage:** PostgreSQL 16+ (\uc0\u1086 \u1089 \u1085 \u1086 \u1074 \u1085 \u1099 \u1077  \u1076 \u1072 \u1085 \u1085 \u1099 \u1077 ), Redis (\u1082 \u1101 \u1096 /\u1086 \u1095 \u1077 \u1088 \u1077 \u1076 \u1080 ).\
- **Infra:** Docker Compose, Nginx (reverse proxy), Prometheus + Grafana.\
- **CI/CD:** GitHub Actions (lint, test, build).\
\
## 3. \uc0\u1057 \u1093 \u1077 \u1084 \u1072  \u1074 \u1079 \u1072 \u1080 \u1084 \u1086 \u1076 \u1077 \u1081 \u1089 \u1090 \u1074 \u1080 \u1103 \
```mermaid\
graph LR\
    Ext[VK/TG API] -->|Webhook| Ing[Go Ingestion Service]\
    Ing -->|Push Event| Redis[(Redis Queue)]\
    Redis -->|Pop Event| Det[Python Detector]\
    Det -->|Read/Write Rules| PG[(PostgreSQL)]\
    Det -->|Alert/Action| Resp[Go Response Service]\
    Resp -->|Notify| TG[TG Bot / Log]\
    PG -->|Metrics| Graf[Grafana Dashboard]\
}