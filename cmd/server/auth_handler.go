package main

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N string `json:"n"`
	E string `json:"e"`
	X5c []string `json:"x5c"`
}

func main() {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            // Verify 'aud' claim
            aud := "http://www.workshop-on-forster.de"
            checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
            if !checkAud {
                return token, errors.New("Invalid audience.")
            }
            // Verify 'iss' claim
            iss := "https://workshop-on-forster.eu.auth0.com/"
            checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
            if !checkIss {
                return token, errors.New("Invalid issuer.")
            }

            cert, err := getPemCert(token)
            if err != nil {
                panic(err.Error())
            }

            result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
            return result, nil
        },
        SigningMethod: jwt.SigningMethodRS256,
    })
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://workshop-on-forster.eu.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	x5c := jwks.Keys[0].X5c
	for k, v := range x5c {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
