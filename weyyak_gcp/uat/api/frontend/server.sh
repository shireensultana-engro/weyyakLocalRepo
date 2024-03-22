export SERVICE_PORT=3003 
export DB_SERVER=34.100.228.116
export DB_SERVER_READER=34.100.228.116
export DB_PORT=5432 
export DB_USER=appuser 
export DB_PASSWORD=weyyakbo345123UAT 
export DB_DATABASE=wk_frontend 
export CONTENT_DB_DATABASE=wyk_content 
export FRONTEND_CONFIG_DB_DATABASE=wyk_frontend_config 
export USER_DB_DATABASE=wk_user_management 
export DEFAULT_PAGE_SIZE=20 
export CONTENT_COMMON_API_URL= 
export CMS=https://wyk2qa.weyyak.com/v1/ 
export UM=https://wyk2qa.weyyak.com/v1/ 
export IMAGES=https://content-uat.weyyak.com/ 
export GEO_LOCATION=https://geo.weyyak1.z5.com 
export VIDEO_API=https://api-weyyak.akamaized.net/get_info/ 
export AD_TAG_URL=https://s3.ap-south-1.amazonaws.com/z5xml/mobile_apps_ads_ios.xml 
export BASE_URL=https://wyk2qa.weyyak.com/
export USER_ACTIVITY_URL=https://msapiqa-events.z5.com/event/activity 
export CONTENT_TYPE_URL=https://wyk2qa.weyyak.com/v1/en/contents/contentType?contentType= 
export CONTENT_TYPE_URL_PAGINATION=&pageNo=1&OrderBy=desc&RowCountPerPage=50&IsPaging=0 
export USER_LOG_URL=https://msapiqa-events.z5.com/event/log 
export S3_BUCKET=z5content-uat 
export S3_URLFORCONFIG=https://content-uat.weyyak.com/configqamigrate.json 
export CONFIG_KEY=configqa 
export REDIS_CACHE_URL=https://msapiqa-events.z5.com/v1/cache
export REDIS_CONTENT_KEY=GOAPIQA 
export DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token=
export PATH=$(go env GOPATH)/bin:$PATH
export JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces
export JEAGER_SERVICE_NAME=Frontend
export BUCKET_NAME=wyk-content-uat
# export TYPE=service_account
# export PROJECT_ID=testing-400312
# export PRIVATE_KEY_ID=dd03c9b33e8647b787124ad3a0775e5bcf95cf3b
# export PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC4nz+BJqcU1t0V\n9EogDpucUrxDOfWiS5z0NGTYmThH3UZF13D1Ape3Bv+aE+OWfSXkL+95+HHwSCHi\nBk8qCcluIPz7BxVYVgkP/9czg9cX4OHK7Aho7wofbznZqLr32wCnUc4Vhf0Wgd2I\npzld7U71oA2V6FuL5BMlQ/Kw/Sph/oMFLYy/PuBfL/YNnhU6jdaNWAibQM1JruW/\nHBYbcKONSrTs50WPSoW1ZZgXJFeqgAz9MSGuvBhFu6/LSUINCCrnyotkJY65p4yL\n+jpEdaqzbw5xkR/Sckk2t+DFZyRowamNXAr4iGsk/5KnKokv0fcaCLRriaHpSdXs\natckZmtZAgMBAAECggEABsZpSS8G8J/V6ylU0wpWMY1jtT/aQMNvlhKaJMyyoQiA\nvK1kSsl2kdPi61+ReYNMMayEqEAyxjOPcsDyhMNpLe5t+jRPVzeJC5pC5nQPH6Q0\nBQNWZ6tl/rNRNyiW++OAiaNZ4bZSDFJls88XLtg3jpH6DadCPMb44OQ2csJHnOLv\nK5dMzUenvh2V/zWVvNS2CaTEFqzWowbzIek/Cv2hWPizo6ZwgplTKQu1iFfbrBui\nYs8WvvpsxA8X216SyECCzo9eZv2lAwFtLedvf/4MfhNL4q6/0KTHDiTDWuWkx0f1\nzVeLbK8b4LfXplZomF9N4fcq/z8SVKpFE1DjuL9jEQKBgQDfiuCHw3K/hOAZCF+c\nNyZ5M2KIyooOv1mxZrVBifvZjlO6kji6L2L9BtIYryYRG0NLrXEpmVZz7IM3NuT1\nrCfuzKQ+l07dEtg6VngoYup1JDKrWTJD4v2cqwOJJAXdQNU03b+shTmXZvW5tJ2J\nbSibIEE6Ol9BFA3b69FV08m9PQKBgQDTbbC3bhIPa5MEo8j2kvM/S9KdJaomjST3\nywCl64sbQXzSFKH3hMB+h7zX2ztEU+o1Ls7P8lvFFQEg7R78kDZy6ulPosm1GFcE\nydOriruXDMHESBjD8YWGO+PevsJIp3j0SADFdvoJg/MHcdpJBI/JDZKUA5mNP52B\n+gD3MCKATQKBgHyh1lt7Oe/TqonqZDwZd4bdglNX8S8VunExHV+kCdmbA82ilqQf\npWYDNoHyrRuegp+f3NbfmhbZx7KyFWdvi3gVeoE3JQ4W4p1r9mQ1+hhDjUiBW4gD\n93gw1LDSd76K1hQ6ihIq2RgznE7kh2zGgnwyuIs5XkBPaQazbKwYf4LVAoGAUES3\nr5whTStpIHzSAhLeOKyfpDu1cndpjo3KjDN1l35wVg2xRBhpQGBmKIk54gH9y/0e\nVUJM5vDHgemkNvzFPzHCLBLAg3lfKBk7vEeqWnlkYxGAHXvnVoQMSfegKRczy4I+\nkLlPyicHME9gMRKDSDBX8su/EyoQsVTp4u6qWo0CgYEAoAH/l+OI6TbWmvA3oISu\n6lVDSgzp28FEJWsxZTleQHTXTa3qV2ePq7IS3Nd021Kojyn5gZ3U6R5nhg114LM+\nojqUSY3MQeXBPYfGaLTIDn/fy27Iq6R6UFVKpB3P61oQxwPSV3QcbbzzsvFcSFOZ\n5WXXSnGkpsn3iFOdCGKNsE8=\n-----END PRIVATE KEY-----\n"
# export CLIENT_EMAIL=weyyak-qa@testing-400312.iam.gserviceaccount.com
# export CLIENT_ID="106140453468102648729"
# export AUTH_URI=https://accounts.google.com/o/oauth2/auth
# export TOKEN_URI=https://oauth2.googleapis.com/token
# export AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs
# export CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/weyyak-qa%40testing-400312.iam.gserviceaccount.com
# export UNIVERSE_DOMAIN=googleapis.com
export TYPE=service_account
export PROJECT_ID=testing-400312
export PRIVATE_KEY_ID=dd03c9b33e8647b787124ad3a0775e5bcf95cf3b
export PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC4nz+BJqcU1t0V\n9EogDpucUrxDOfWiS5z0NGTYmThH3UZF13D1Ape3Bv+aE+OWfSXkL+95+HHwSCHi\nBk8qCcluIPz7BxVYVgkP/9czg9cX4OHK7Aho7wofbznZqLr32wCnUc4Vhf0Wgd2I\npzld7U71oA2V6FuL5BMlQ/Kw/Sph/oMFLYy/PuBfL/YNnhU6jdaNWAibQM1JruW/\nHBYbcKONSrTs50WPSoW1ZZgXJFeqgAz9MSGuvBhFu6/LSUINCCrnyotkJY65p4yL\n+jpEdaqzbw5xkR/Sckk2t+DFZyRowamNXAr4iGsk/5KnKokv0fcaCLRriaHpSdXs\natckZmtZAgMBAAECggEABsZpSS8G8J/V6ylU0wpWMY1jtT/aQMNvlhKaJMyyoQiA\nvK1kSsl2kdPi61+ReYNMMayEqEAyxjOPcsDyhMNpLe5t+jRPVzeJC5pC5nQPH6Q0\nBQNWZ6tl/rNRNyiW++OAiaNZ4bZSDFJls88XLtg3jpH6DadCPMb44OQ2csJHnOLv\nK5dMzUenvh2V/zWVvNS2CaTEFqzWowbzIek/Cv2hWPizo6ZwgplTKQu1iFfbrBui\nYs8WvvpsxA8X216SyECCzo9eZv2lAwFtLedvf/4MfhNL4q6/0KTHDiTDWuWkx0f1\nzVeLbK8b4LfXplZomF9N4fcq/z8SVKpFE1DjuL9jEQKBgQDfiuCHw3K/hOAZCF+c\nNyZ5M2KIyooOv1mxZrVBifvZjlO6kji6L2L9BtIYryYRG0NLrXEpmVZz7IM3NuT1\nrCfuzKQ+l07dEtg6VngoYup1JDKrWTJD4v2cqwOJJAXdQNU03b+shTmXZvW5tJ2J\nbSibIEE6Ol9BFA3b69FV08m9PQKBgQDTbbC3bhIPa5MEo8j2kvM/S9KdJaomjST3\nywCl64sbQXzSFKH3hMB+h7zX2ztEU+o1Ls7P8lvFFQEg7R78kDZy6ulPosm1GFcE\nydOriruXDMHESBjD8YWGO+PevsJIp3j0SADFdvoJg/MHcdpJBI/JDZKUA5mNP52B\n+gD3MCKATQKBgHyh1lt7Oe/TqonqZDwZd4bdglNX8S8VunExHV+kCdmbA82ilqQf\npWYDNoHyrRuegp+f3NbfmhbZx7KyFWdvi3gVeoE3JQ4W4p1r9mQ1+hhDjUiBW4gD\n93gw1LDSd76K1hQ6ihIq2RgznE7kh2zGgnwyuIs5XkBPaQazbKwYf4LVAoGAUES3\nr5whTStpIHzSAhLeOKyfpDu1cndpjo3KjDN1l35wVg2xRBhpQGBmKIk54gH9y/0e\nVUJM5vDHgemkNvzFPzHCLBLAg3lfKBk7vEeqWnlkYxGAHXvnVoQMSfegKRczy4I+\nkLlPyicHME9gMRKDSDBX8su/EyoQsVTp4u6qWo0CgYEAoAH/l+OI6TbWmvA3oISu\n6lVDSgzp28FEJWsxZTleQHTXTa3qV2ePq7IS3Nd021Kojyn5gZ3U6R5nhg114LM+\nojqUSY3MQeXBPYfGaLTIDn/fy27Iq6R6UFVKpB3P61oQxwPSV3QcbbzzsvFcSFOZ\n5WXXSnGkpsn3iFOdCGKNsE8=\n-----END PRIVATE KEY-----\n",
export CLIENT_EMAIL=weyyak-qa@testing-400312.iam.gserviceaccount.com
export CLIENT_ID=106140453468102648729
export AUTH_URI=https://accounts.google.com/o/oauth2/auth
export TOKEN_URI=https://oauth2.googleapis.com/token
export AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs
export CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/weyyak-qa%40testing-400312.iam.gserviceaccount.com
export UNIVERSE_DOMAIN=googleapis.com
# swag init -g main.go
go run *.go
