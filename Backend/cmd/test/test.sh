# Request Code
printf "\033[32m*** /auth/request-code ***\033[0m\n"
curl -i -X POST http://localhost:8080/api/auth/request-code \
-H "Content-Type:application/json" \
-d '{"email":"patrick.jiang@plu.edu","scene":"signup"}'