# Request Code
printf "\033[32m*** /auth/request-code ***\033[0m\n"
curl -s -D /tmp/req_hdr.txt -o /tmp/req_body.json \
  -X POST http://localhost:8080/api/auth/request-code \
  -H "Content-Type: application/json" \
  -d '{"email":"patrick.jiang@plu.edu","scene":"signup"}'

cat /tmp/req_hdr.txt
cat /tmp/req_body.json

codeID=$(jq -er '.code_id' /tmp/req_body.json) || {
  echo "failed to parse .code_id from body"; exit 1;
}

# Verify Code
printf "\n\033[32m*** /auth/verify-code ***\033[0m\n"
echo "Enter the code: "
read code
curl -s -D /tmp/req_hdr.txt -o /tmp/req_body.json \
  -X POST http://localhost:8080/api/auth/verify-code \
  -H "Content-Type:application/json" \
  -d "{\"code_id\":\"$codeID\",\"code\":\"$code\",\"email\":\"patrick.jiang@plu.edu\",\"scene\":\"signup\"}"

cat /tmp/req_hdr.txt
cat /tmp/req_body.json

otp=$(jq -er '.token' /tmp/req_body.json) || {
  echo "failed to parse .token from body"; exit 1;
}

## Create Account
#printf "\n\033[32m*** /auth/create-account ***\033[0m\n"
#curl -s -D /tmp/req_hdr.txt -o /tmp/req_body.json \
#  -X POST http://localhost:8080/api/auth/create-account \
#  -H "Content-Type:application/json" \
#  -H "Authorization: Bearer $otp"\
#  -d '{"password":"ThisIsFake123!"}'
#
#cat /tmp/req_hdr.txt
#cat /tmp/req_body.json
#
#ATK=$(jq -er '.access_token' /tmp/req_body.json) || {
#      echo "failed to parse .access_token from body"; exit 1;
#}
#
## Set Username
#printf "\n\033[32m*** /auth/me/set-username ***\033[0m\n"
#curl -i -X PATCH http://localhost:8080/api/auth/me/set-username \
#  -H "Content-Type:application/json" \
#  -H "Authorization: Bearer $ATK"  \
#  -d '{"username":"patrick jiang"}'