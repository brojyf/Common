# Request Code
printf "\033[32m*** /auth/request-code ***\033[0m\n"
curl -i -X POST http://localhost:8080/api/auth/request-code \
  -H "Content-Type:application/json" \
  -d '{"email":"patrick.jiang@plu.edu","scene":"signup"}'

#Verify Code
printf "\n\033[32m*** /auth/verify-code ***\033[0m\n"
echo "Enter the code: "
read code
echo "Enter the code_id: "
read codeID
curl -i -X POST http://localhost:8080/api/auth/verify-code \
  -H "Content-Type:application/json" \
  -d "{\"otp_id\":\"$codeID\",\"code\":\"$code\",\"email\":\"patrick.jiang@plu.edu\",\"scene\":\"signup\"}"
