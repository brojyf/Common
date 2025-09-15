# Request Code
printf "\033[32m*** /auth/request-code ***\033[0m\n"
curl -i -X POST http://localhost:8080/api/auth/request-code \
  -H "Content-Type:application/json" \
  -d '{"email":"patrick.jiang@plu.edu","scene":"signup"}'

#Verify Code
printf "\033[32m*** /auth/verify-code ***\033[0m\n"
echo "Enter the code: "
read code
curl -i -X POST http://localhost:8080/api/auth/verify-code \
  -H "Content-Type:application/json" \
  -d "{\"email\":\"patrick.jiang@plu.edu\",\"scene\":\"signup\",\"code\":\"$code\"}"
