{\rtf1\ansi\ansicpg1251\cocoartf2822
\cocoatextscaling0\cocoaplatform0{\fonttbl\f0\fswiss\fcharset0 Helvetica;}
{\colortbl;\red255\green255\blue255;}
{\*\expandedcolortbl;;}
\paperw11900\paperh16840\margl1440\margr1440\vieww11520\viewh8400\viewkind0
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\f0\fs24 \cf0 # \uc0\u1069 \u1090 \u1072 \u1087  2: \u1044 \u1074 \u1080 \u1078 \u1086 \u1082  \u1076 \u1077 \u1090 \u1077 \u1082 \u1094 \u1080 \u1080  (Python)\
\
## \uc0\u1062 \u1077 \u1083 \u1100 \
\uc0\u1056 \u1077 \u1072 \u1083 \u1080 \u1079 \u1086 \u1074 \u1072 \u1090 \u1100  \u1075 \u1080 \u1073 \u1088 \u1080 \u1076 \u1085 \u1091 \u1102  \u1089 \u1080 \u1089 \u1090 \u1077 \u1084 \u1091  \u1076 \u1077 \u1090 \u1077 \u1082 \u1094 \u1080 \u1080 : \u1087 \u1088 \u1072 \u1074 \u1080 \u1083 \u1072  (SQL/Python) + ML-\u1072 \u1085 \u1086 \u1084 \u1072 \u1083 \u1080 \u1080 .\
\
## \uc0\u1060 \u1091 \u1085 \u1082 \u1094 \u1080 \u1086 \u1085 \u1072 \u1083 \u1100 \u1085 \u1099 \u1077  \u1090 \u1088 \u1077 \u1073 \u1086 \u1074 \u1072 \u1085 \u1080 \u1103 \
1. **Consumer:** \uc0\u1063 \u1090 \u1077 \u1085 \u1080 \u1077  \u1080 \u1079  Redis Stream `abuse:events:raw` \u1089  consumer group.\
2. **Rule Engine:**\
   - \uc0\u1047 \u1072 \u1075 \u1088 \u1091 \u1079 \u1082 \u1072  \u1087 \u1088 \u1072 \u1074 \u1080 \u1083  \u1080 \u1079  \u1041 \u1044  (\u1090 \u1072 \u1073 \u1083 \u1080 \u1094 \u1072  `detection_rules`).\
   - \uc0\u1055 \u1088 \u1080 \u1084 \u1077 \u1085 \u1077 \u1085 \u1080 \u1077  \u1087 \u1088 \u1072 \u1074 \u1080 \u1083  \u1082  \u1073 \u1072 \u1090 \u1095 \u1091  \u1089 \u1086 \u1073 \u1099 \u1090 \u1080 \u1081 .\
   - \uc0\u1055 \u1088 \u1080 \u1084 \u1077 \u1088  \u1087 \u1088 \u1072 \u1074 \u1080 \u1083 \u1072 : "N \u1089 \u1086 \u1073 \u1099 \u1090 \u1080 \u1081  \u1086 \u1090  \u1086 \u1076 \u1085 \u1086 \u1075 \u1086  user_id \u1079 \u1072  T \u1089 \u1077 \u1082 \u1091 \u1085 \u1076 ".\
3. **ML Detector:**\
   - \uc0\u1055 \u1077 \u1088 \u1080 \u1086 \u1076 \u1080 \u1095 \u1077 \u1089 \u1082 \u1086 \u1077  (\u1088 \u1072 \u1079  \u1074  5 \u1084 \u1080 \u1085 ) \u1086 \u1073 \u1091 \u1095 \u1077 \u1085 \u1080 \u1077 /\u1080 \u1085 \u1092 \u1077 \u1088 \u1077 \u1085 \u1089  \u1084 \u1086 \u1076 \u1077 \u1083 \u1080  Isolation Forest \u1085 \u1072  \u1087 \u1086 \u1089 \u1083 \u1077 \u1076 \u1085 \u1080 \u1093  N \u1089 \u1086 \u1073 \u1099 \u1090 \u1080 \u1103 \u1093 .\
   - \uc0\u1057 \u1082 \u1086 \u1088 \u1080 \u1085 \u1075  \u1072 \u1085 \u1086 \u1084 \u1072 \u1083 \u1100 \u1085 \u1086 \u1089 \u1090 \u1080  \u1076 \u1083 \u1103  \u1089 \u1086 \u1073 \u1099 \u1090 \u1080 \u1081 , \u1085 \u1077  \u1087 \u1086 \u1087 \u1072 \u1074 \u1096 \u1080 \u1093  \u1087 \u1086 \u1076  \u1087 \u1088 \u1072 \u1074 \u1080 \u1083 \u1072 .\
4. **Enrichment:** \uc0\u1054 \u1073 \u1086 \u1075 \u1072 \u1097 \u1077 \u1085 \u1080 \u1077  \u1089 \u1086 \u1073 \u1099 \u1090 \u1080 \u1103  \u1082 \u1086 \u1085 \u1090 \u1077 \u1082 \u1089 \u1090 \u1086 \u1084  \u1080 \u1079  \u1041 \u1044  (\u1080 \u1089 \u1090 \u1086 \u1088 \u1080 \u1103  \u1078 \u1072 \u1083 \u1086 \u1073 , \u1074 \u1086 \u1079 \u1088 \u1072 \u1089 \u1090  \u1072 \u1082 \u1082 \u1072 \u1091 \u1085 \u1090 \u1072 ).\
5. **Output:** \uc0\u1047 \u1072 \u1087 \u1080 \u1089 \u1100  \u1074 \u1077 \u1088 \u1076 \u1080 \u1082 \u1090 \u1072  (`PASS`, `WARN`, `BLOCK`) \u1074  Redis Stream `abuse:verdicts`.\
\
## \uc0\u1052 \u1086 \u1076 \u1077 \u1083 \u1100  \u1076 \u1072 \u1085 \u1085 \u1099 \u1093  (PostgreSQL)\
```sql\
CREATE TABLE events (\
    id UUID PRIMARY KEY,\
    user_id TEXT NOT NULL,\
    event_type TEXT NOT NULL,\
    timestamp TIMESTAMPTZ NOT NULL,\
    verdict TEXT,\
    score FLOAT,\
    matched_rule_id INT\
);\
\
CREATE TABLE detection_rules (\
    id SERIAL PRIMARY KEY,\
    name TEXT NOT NULL,\
    condition_json JSONB NOT NULL, -- e.g. \{"window_sec": 60, "threshold": 10\}\
    action TEXT NOT NULL, -- BLOCK, WARN\
    is_active BOOLEAN DEFAULT TRUE\
);\
}