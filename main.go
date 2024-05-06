package main

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
)

const email = `<html>
<head>
  <meta http-equiv="Content-Type" content="text/html;charset=UTF-8">
  <style>
      /* inter-latin-400-normal */
      @font-face {
          font-family: 'Inter';
          font-style: normal;
          font-display: swap;
          font-weight: 400;
          src: url(https://cdn.jsdelivr.net/fontsource/fonts/inter@latest/latin-400-normal.woff2) format('woff2'), url(https://cdn.jsdelivr.net/fontsource/fonts/inter@latest/latin-400-normal.woff) format('woff');
          unicode-range: U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+0304,U+0308,U+0329,U+2000-206F,U+2074,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD;
      }
  </style>
</head>
<body style="background-color:#F2F5FA;min-width:600px;font-family: Inter,serif ">

  <div style="background-color: #FFFFFF;width:600px;box-sizing: border-box;padding:32px;border-radius: 30px;margin:60px auto;">
    <img src="https://d27r9m0vtnpoa0.cloudfront.net/logo.svg"/>
    <img src="https://d27r9m0vtnpoa0.cloudfront.net/bg.png" style="width:100%;border-radius: 24px;margin:32px 0;">
    <div style="font-size: 28px;font-weight: 700;line-height: 30.8px;color:#1F2226">
      Your application has been approved
    </div>
    <div style="font-size: 16px;font-weight: 400;line-height: 22.4px;letter-spacing: -0.01em;text-align: left;color:#1F2226">
      We are thrilled to welcome you to Ankr, and excited to secure your organization's technical advantage. We appreciate your decision to choose our platform and are confident our infrastructure services will amplify your enterprise.
    </div>
    <a href="https://imo-dev.neuraprotocol.io/imo/myProjects" target="_blank" style="display: block;text-decoration: none;background-color: #356DF3;padding:0 32px;height:60px;line-height:60px;border-radius: 32px;font-size: 20px;font-weight: 600;text-align: center;color:#ffffff;margin-top:32px;">Launch Your IMO</a>
    <div style="margin-top:32px;color: #82899A;font-size: 14px;font-weight: 400;line-height: 19.6px;text-align: left;">
      <div>Want to change which emails we send you? </div>
      <div>Customize your email <span style="text-decoration: underline">notification settings</span> at any time.</div>

    </div>
    <div style="margin-top:32px;border-top:1px solid #E7EBF3">
      <div style="display: flex;margin-top:13px;font-size:14px;font-weight: 400;line-height: 20px;justify-content: space-between;align-items: center;">
        <a target="_blank" style="text-decoration: none;color:#82899A;gap:5px;display: flex;align-items: center" href="https://twitter.com/ankr">
          <img style="width:14px;height:12px;" src="https://d27r9m0vtnpoa0.cloudfront.net/twitter.svg">
          <span>Twitter</span>
        </a>
        <a target="_blank" style="text-decoration: none;color:#82899A;gap:5px;display: flex;align-items: center" href="https://discord.ankr.com/">
          <img style="width:16px;height:16px;" src="https://d27r9m0vtnpoa0.cloudfront.net/discord.svg">
          <span>Discord</span>
        </a>
        <a target="_blank" style="text-decoration: none;color:#82899A;gap:5px;display: flex;align-items: center" href="javascript:void(0)">
          <img style="width:16px;height:10px;" src="https://d27r9m0vtnpoa0.cloudfront.net/medium.svg">
          <span>Medium</span>
        </a>
      </div>
    </div>
    <div style="margin-top:32px;display: flex;align-items: center;justify-content: space-between;color:#82899A;font-size:14px;font-weight: 400">
      <div>&copy; 2024 Neura All rights reserved</div>
      <div>info@neura.com</div>
    </div>


  </div>
</body>
</html>
`

func main() {
	log.Print("init handler")
	http.HandleFunc("/email", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		encrypt := r.URL.Query().Get("encrypt")
		println(address)
		println(encrypt)
		origin, err := decryptAES(encrypt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("invalid request params %s", err.Error())))
			return
		}
		if origin != address {
			println(origin)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("encrypt doesn't match"))
			return
		}

		m := gomail.NewMessage()
		m.SetHeader("From", "info@ankr.com")
		m.SetHeader("To", "kui@ankr.com")
		m.SetHeader("Subject", "IMO events from Neura")
		m.SetBody("text/html", email)

		// Dialer configuration
		dialer := gomail.NewDialer("smtp.gmail.com", 587, "shicai@ankr.com", "jynbcdyvhpxjjdvg")

		// Send email
		if err := dialer.DialAndSend(m); err != nil {
			panic(err)
		} else {
			print(w.Write([]byte("Hello, world!")))
		}
	})
	_ = http.ListenAndServe(":8080", nil)
}

// unpad removes PKCS#7 padding from the given data.
func unpad(buf []byte) []byte {
	padLen := int(buf[len(buf)-1])
	return buf[:len(buf)-padLen]
}

// decryptAES decrypts ciphertext using AES-256 in CBC mode.
func decryptAES(cipherTextBase64 string) (string, error) {
	key := os.Getenv("key")
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	data, err := hex.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(data)%blockSize != 0 {
		return "", fmt.Errorf("invalid data length")
	}

	decrypted := make([]byte, len(data))

	for i := 0; i < len(data); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], data[i:i+blockSize])
	}

	decrypted = unpad(decrypted)
	return string(decrypted), nil
}
