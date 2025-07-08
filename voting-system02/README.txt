1.Take token
curl -s -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d "{\"username\":\"admin\",\"password\":\"admin123\"}"

2.CC the TOKEN insert this sentence
curl -s -X POST http://localhost:8080/api/polls -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" -d "{\"title\":\"最佳编程语言？\",\"description\":\"\",\"options\":[\"Go\",\"Python\",\"JS\"],\"end_at\":\"2025-12-31T00:00:00Z\"}"

3.Show the poll list
curl -s http://localhost:8080/api/polls

4.Check the content
curl -s http://localhost:8080/api/polls/{POLL_ID}

5.Poll
curl -s -X POST http://localhost:8080/api/polls/{POLL_ID}/vote -H "Content-Type: application/json" -H "Authorization: Bearer <TOKEN>" -d "{\"option_id\":\"{OPTION_ID}\"}"

token:b2bf98a4a5568ed2
poll_id:b0bf0d09fbae8476
option_id:138123d034f4d940

curl -s -X POST http://localhost:8080/api/polls -H "Content-Type: application/json" -H "Authorization: Bearer b2bf98a4a5568ed2" -d "{\"title\":\"最佳编程语言？\",\"description\":\"\",\"options\":[\"Go\",\"Python\",\"JS\"],\"end_at\":\"2025-12-31T00:00:00Z\"}"

curl -s -X POST http://localhost:8080/api/polls/b0bf0d09fbae8476/vote -H "Content-Type: application/json" -H "Authorization: Bearer b2bf98a4a5568ed2" -d "{\"option_id\":\"138123d034f4d940\"}"