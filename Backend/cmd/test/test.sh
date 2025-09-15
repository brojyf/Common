# Request Code
printf "\033[32m*** /auth/request-code ***\033[0m\n"
curl -s -D /tmp/req_hdr.txt -o /tmp/req_body.json \
  -X POST http://localhost:8080/api/auth/request-code \
  -H "Content-Type: application/json" \
  -d '{"email":"patrick.jiang@plu.edu","scene":"signup"}'

cat /tmp/req_hdr.txt
cat /tmp/req_body.json

codeID=$(jq -er '.otp_id' /tmp/req_body.json) || {
  echo "failed to parse .otp_id from body"; exit 1;
}

#Verify Code
printf "\n\033[32m*** /auth/verify-code ***\033[0m\n"
echo "Enter the code: "
read code
curl -i -X POST http://localhost:8080/api/auth/verify-code \
  -H "Content-Type:application/json" \
  -d "{\"otp_id\":\"$codeID\",\"code\":\"$code\",\"email\":\"patrick.jiang@plu.edu\",\"scene\":\"signup\"}"
