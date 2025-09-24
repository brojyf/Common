package config

import (
	"fmt"
)

// RedisKeyLoginLockIPEmail login:lock:ip:email:<ip>:<email>
func RedisKeyLoginLockIPEmail(ip, email string) string {
	return fmt.Sprintf("login:lock:ip:email:%s:%s", ip, email)
}

// RedisKeyLoginLockEmail login:lock:email:<email>
func RedisKeyLoginLockEmail(email string) string {
	return fmt.Sprintf("login:lock:email:%s", email)
}

// RedisKeyLoginIPEmailCnt login:cnt:ip:email:<ip>:<email>
func RedisKeyLoginIPEmailCnt(ip, email string) string {
	return fmt.Sprintf("login:cnt:ip:email:%s:%s", ip, email)
}

// RedisKeyLoginEmailCnt login:cnt:email:<email>
func RedisKeyLoginEmailCnt(email string) string {
	return fmt.Sprintf("login:cnt:email:%s", email)
}

// RedisKeyOTTJTIUsed ott:jti:used:<email>:<scene>:<jti>
func RedisKeyOTTJTIUsed(email, scene, jti string) string {
	return fmt.Sprintf("ott:jti:used:%s:%s:%s", email, scene, jti)
}

// RedisKeyVerifyThrottle verify:throttle:<email>:<scene>
func RedisKeyVerifyThrottle(email, scene string) string {
	return fmt.Sprintf("verify:throttle:%s:%s", email, scene)
}

// RedisKeyOTP otp:<email>:<scene>:<codeID>
func RedisKeyOTP(email, scene, codeID string) string {
	return fmt.Sprintf("otp:%s:%s:%s", email, scene, codeID)
}

// RedisKeyThrottle otp:throttle:<email>:<scene>
func RedisKeyThrottle(email, scene string) string {
	return fmt.Sprintf("otp:throttle:%s:%s", email, scene)
}
