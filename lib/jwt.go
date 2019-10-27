package lib

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

// func getKey(token *jwt.Token) (interface{}, error) {
// 	// // Don't forget to validate the alg is what you expect:
// 	// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 	// 	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 	// }

// 	// // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
// 	// return hmacSampleSecret, nil
// 	return []byte("asdf"), nil
// }

func getKey(_ *jwt.Token) (interface{}, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "asdf" // TODO Smarter secret fallback -> .env or something
	}
	return []byte(s), nil
}

func ValidateJWT(token string) bool {
	t, err := jwt.Parse(token, getKey)
	if err != nil {
		// TODO Log this
		return false
	}

	// // THIS IS THE EXAMPLE:
	//t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	//	// Don't forget to validate the alg is what you expect:
	//	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	//	}
	//
	//	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	//	return hmacSampleSecret, nil
	//})

	// // If you want to check claims:
	//if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	//	fmt.Println(claims["foo"], claims["nbf"])
	//} else {
	//	fmt.Println(err)
	//}

	return t.Valid
}
