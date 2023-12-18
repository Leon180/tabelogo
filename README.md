# Tabelogo
## 功能
1. 使用信箱註冊
2. 允許使用位置後可以點擊地圖獲得地點資訊（位置資訊不會儲存在資料庫）
3. 儲存地點
4. 如果地點為日本餐廳, 會有一個特別按鈕：tabelogo, 點擊後等待一段時間會出現該地點近似名稱的tabelog店家資訊
5. quick search: 使用google map autocomplete實現, 輸入文字後會顯示相關結果
7. advance search： 根據條件產出最多二十筆搜尋結果

## 預計開發
1. 信箱驗證功能
2. 新增快取（tabelog資訊

## 展示網頁
https://swarm.tabelogo.com/

若要本地下載後使用請參考/project/swarm.yml檔案 （docker swarm
或/project/docker-compose.yml
使用docker-compose方式請記得建立app.env檔案
```
// front-end/app.env
BROKER_URL=http://localhost:8080
WEBSITE_URL=http://localhost:8081
BROKER_URL_DEPLOYMENT=http://localhost:8080
WEBSITE_URL_DEPLOYMENT=http://localhost:8081

// google-map/cmd/api/app.env
GOOGLE_MAP_API_KEY={Your API KEY}

// authenticate/cmd/api/app.env
DSN_TEST="host=postgres port=5432 user=postgres password=password dbname=tabelogo sslmode=disable timezone=UTC connect_timeout=5"
DSN_DEPLOYMENT="host=host.minikube.internal port=5432 user=postgres password=password dbname=tabelogo sslmode=disable timezone=UTC connect_timeout=5"
RABBITMQ_CONNECT="amqp://guest:guest@rabbitmq"
REDIS_CONNECT_SESSION="redis://redis-master-session:6379"
REDIS_CONNECT_PLACE="redis://redis-master-place:6379"
TOKEN_SYMMETRIC_KEY={32 length key}
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=24h
REFRESH_DURATION=5m

// logger-service/cmd/api/app.env
MONGO_USER={Your mongo user name}
MONGO_PASSWORD={Your mongo password}
```

資料庫初始化請參考/project/Makefile migrate_up
