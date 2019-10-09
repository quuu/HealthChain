package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"

	"io"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

type M map[string]interface{}

func mapsfjson() []M {
	//open the file
	jsonFile, err := os.Open("hc_data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// read bytes in file
	byteVal, _ := ioutil.ReadAll(jsonFile)

	// store the rows of json file into list of maps
	var rows []M
	json.Unmarshal([]byte(byteVal), &rows)

	return rows
}

func valtostr(v interface{}) string {
	str := fmt.Sprintf("%v", v)
	return str
}

func db_key(patient_record M) string {
	fname := valtostr(patient_record["full_name"])
	ssn := valtostr(patient_record["ssn"])
	dob := valtostr(patient_record["dob"])
	country := valtostr(patient_record["Country"])
	return fname + ssn + dob + country
}

func encrypt_json_string(jsonStr string, encrpyt_key string) []byte {
	json_bytes := []byte(jsonStr)
	return encrypt(json_bytes, encrpyt_key)

}

func decrypt_json_string(json_bytes []byte, encrpyt_key string) string {
	// json_bytes := []byte(jsonStr)
	return string(decrypt(json_bytes, encrpyt_key))

}

// func main() {
// 	rows := mapsfjson()
// 	// fmt.Println(rows[0])
// 	jstring, err := json.Marshal(rows[0])
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	jsonStr := string(jstring)
// 	fmt.Println("heres the json string ->\n", jsonStr, "\n")

// 	var record M
// 	json.Unmarshal([]byte(jsonStr), &record)
// 	fmt.Println("Here is the json string converted back to map->\n", record, "\n")
// 	// here is how you marsh/unmarsh json into string and out of strings

// 	// todo
// 	// encrypt the json string, then decrypt, turn back into map
// 	// var fname string

// 	fname := valtostr(record["full_name"])

// 	ssn := valtostr(record["ssn"])
// 	dob := valtostr(record["dob"])
// 	country := valtostr(record["Country"])

// 	patient_id := fname + ssn + dob + country
// 	fmt.Println("Here is the patient_id string -> ", patient_id)

// 	// hash := sha1.New()

// 	// hash.Write([]byte(patient_id))
// 	// sha1_hash := hex.EncodeToString(hash.Sum(nil))

// 	// fmt.Println("Hash of ", patient_id, " -> ", sha1_hash)

// 	record_jsonStr, _ := json.Marshal(record)
// 	fmt.Println("Here is  json string (again) -> ", string(record_jsonStr))
// 	encrypted_json := encrypt_json_string(string(record_jsonStr), patient_id)
// 	fmt.Println("Here is record_json, encrypted -> ", hex.EncodeToString(encrypted_json))

// 	decrypted_json := decrypt_json_string(encrypted_json, patient_id)
// 	fmt.Println("Here is record_json, dencrypted -> ", decrypted_json)

// }
