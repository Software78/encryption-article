package middleware

import (
        "bytes"
        "crypto/aes"
        "crypto/cipher"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "os"
        "regexp"
        "github.com/gin-gonic/gin"
)

type CryptoMiddleware struct {
        key            []byte
        iv             []byte // Initialization Vector
        excludePattern *regexp.Regexp
}

func NewCryptoMiddlewareFromEnv(excludePattern string) (*CryptoMiddleware, error) {
        key := os.Getenv("AES_SECRET_KEY")
        iv := os.Getenv("AES_IV")

        if key == "" {
                return nil, fmt.Errorf("AES_SECRET_KEY environment variable is not set")
        }
        if iv == "" {
                return nil, fmt.Errorf("AES_IV environment variable is not set")
        }

        if len(key) != 32 {
                return nil, fmt.Errorf("AES_SECRET_KEY must be exactly 32 characters (got %d characters)", len(key))
        }

        if len(iv) != 16 {
                return nil, fmt.Errorf("AES_IV must be exactly 16 characters (got %d characters)", len(iv))
        }

        var pattern *regexp.Regexp
        if excludePattern != "" {
                _, err := regexp.Compile(excludePattern)
                if err != nil {
                        return nil, fmt.Errorf("invalid exclude pattern: %w", err)
                }
        }

        return &CryptoMiddleware{
                key:            []byte(key), // Convert key to byte slice
                iv:             []byte(iv),  // Convert IV to byte slice
                excludePattern: pattern,
        }, nil
}

func (m *CryptoMiddleware) encrypt(plaintext []byte) (string, error) {
        block, err := aes.NewCipher(m.key)
        if err != nil {
                return "", err
        }

        // Pad the plaintext using PKCS7 padding
    plaintext = pkcs7Pad(plaintext, aes.BlockSize)

        ciphertext := make([]byte, len(plaintext))
        mode := cipher.NewCBCEncrypter(block, m.iv)
        mode.CryptBlocks(ciphertext, plaintext)

        return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (m *CryptoMiddleware) decrypt(encrypted string) ([]byte, error) {
        ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
        if err != nil {
                return nil, err
        }

        block, err := aes.NewCipher(m.key)
        if err != nil {
                return nil, err
        }

        plaintext := make([]byte, len(ciphertext))
        mode := cipher.NewCBCDecrypter(block, m.iv)
        mode.CryptBlocks(plaintext, ciphertext)

    plaintext = pkcs7Unpad(plaintext, aes.BlockSize)

        return plaintext, nil
}

// PKCS7 padding
func pkcs7Pad(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padtext...)
}

// PKCS7 unpadding
func pkcs7Unpad(data []byte, blockSize int) []byte {
    padding := int(data[len(data)-1])
    if padding > blockSize || padding > len(data) {
        return data
    }
    return data[:len(data)-padding]
}

type responseWriter struct {
        gin.ResponseWriter
        body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
        return w.body.Write(b)
}

func (m *CryptoMiddleware) shouldSkipEncryption(path string) bool {
        if m.excludePattern == nil {
                return false
        }
        return m.excludePattern.MatchString(path)
}

// DecryptRequestMiddleware decrypts the request body if it's JSON.
func (m *CryptoMiddleware) DecryptRequestMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
                if m.shouldSkipEncryption(c.Request.URL.Path) {
                        c.Next()
                        return
                }

               

                body, err := io.ReadAll(c.Request.Body)
                if err != nil {
                        c.AbortWithStatus(http.StatusBadRequest)
                        return
                }

                if len(body) > 0 {
                        decryptedData, err := m.decryptJSON(body)
                        if err != nil {
                                // Return 400 Bad Request for any decryption failures
                                c.AbortWithError(http.StatusBadRequest, 
                                        fmt.Errorf("failed to decrypt request body: %v", err))
                                return
                        }

                        c.Set("decryptedJSON", decryptedData)
                        c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
                }

                c.Next()
        }
}

func (m *CryptoMiddleware) EncryptResponseMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
                // Capture Response Body
                writer := &responseWriter{
                        ResponseWriter: c.Writer,
                        body:          &bytes.Buffer{},
                }
                c.Writer = writer

                c.Next() // Let the handler process the request

                if writer.body.Len() > 0 {
                        // Unmarshal the data from the handler
                        responseData := make(map[string]interface{})
                        err := json.Unmarshal(writer.body.Bytes(), &responseData)
                        if err != nil {
                                c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to unmarshal response body: %w", err))
                                return
                        }

                        encryptedData, err := m.encryptJSON(writer.body.Bytes()) // Encrypt!
                        if err != nil {
                                c.AbortWithError(http.StatusInternalServerError, err)
                                return
                        }

                        encryptedJSON, err := json.Marshal(encryptedData)
                        if err != nil {
                                c.AbortWithError(http.StatusInternalServerError, err)
                                return
                        }

                        c.Header("Content-Type", "application/json")
                        c.Writer.Write(encryptedJSON) // Write the encrypted JSON
                }
        }
}

func (m *CryptoMiddleware) encryptJSON(data []byte) (map[string]interface{}, error) {
        var jsonData map[string]interface{}
        err := json.Unmarshal(data, &jsonData)
        if err != nil {
                return nil, err
        }

        for key, value := range jsonData {
                switch v := value.(type) {
                case string:
                        processedValue, err := m.encrypt([]byte(v))
                        if err != nil {
                                return nil, fmt.Errorf("error encrypting field %s: %w", key, err)
                        }
                        jsonData[key] = processedValue

                case map[string]interface{}:
                        nestedJSON, err := json.Marshal(v)
                        if err != nil {
                                return nil, err
                        }
                        processedNestedData, err := m.encryptJSON(nestedJSON)
                        if err != nil {
                                return nil, err
                        }
                        jsonData[key] = processedNestedData

                case []interface{}:
                        for i, item := range v {
                                itemJSON, err := json.Marshal(item)
                                if err != nil {
                                        return nil, err
                                }
                                processedItem, err := m.encryptJSON(itemJSON)
                                if err != nil {
                                        return nil, err
                                }
                                v[i] = processedItem
                        }
                        jsonData[key] = v

                default:
                        processedValue, err := m.encrypt([]byte(fmt.Sprintf("%v", v)))
                        if err != nil {
                                return nil, fmt.Errorf("error encrypting field %s: %w", key, err)
                        }
                        jsonData[key] = processedValue
                }
        }

        return jsonData, nil
}

func (m *CryptoMiddleware) decryptJSON(data []byte) (map[string]interface{}, error) {
        var jsonData map[string]interface{}
        err := json.Unmarshal(data, &jsonData)
        if err != nil {
                return nil, fmt.Errorf("invalid JSON format: %v", err)
        }

        for key, value := range jsonData {
                switch v := value.(type) {
                case string:
                        decryptedBytes, err := m.decrypt(v)
                        if err != nil {
                                return nil, fmt.Errorf("failed to decrypt field '%s': %v", key, err)
                        }
                        
                        // Handle JSON string values (wrapped in quotes)
                        if len(v) > 1 && v[0] == '"' && v[len(v)-1] == '"' {
                                var decryptedValue interface{}
                                err = json.Unmarshal(decryptedBytes, &decryptedValue)
                                if err != nil {
                                        return nil, fmt.Errorf("invalid JSON in decrypted value for field '%s': %v", key, err)
                                }
                                jsonData[key] = decryptedValue
                        } else {
                                jsonData[key] = string(decryptedBytes)
                        }

                case map[string]interface{}:
                        nestedJSON, err := json.Marshal(v)
                        if err != nil {
                                return nil, fmt.Errorf("failed to marshal nested object in field '%s': %v", key, err)
                        }
                        processedNestedData, err := m.decryptJSON(nestedJSON)
                        if err != nil {
                                return nil, err
                        }
                        jsonData[key] = processedNestedData

                case []interface{}:
                        for i, item := range v {
                                itemJSON, err := json.Marshal(item)
                                if err != nil {
                                        return nil, fmt.Errorf("failed to marshal array item %d in field '%s': %v", i, key, err)
                                }
                                processedItem, err := m.decryptJSON(itemJSON)
                                if err != nil {
                                        return nil, err
                                }
                                v[i] = processedItem
                        }
                        jsonData[key] = v

                default:
                        // For non-string primitive values (numbers, booleans, null), keep as is
                }
        }

        return jsonData, nil
}



func (m *CryptoMiddleware) EncryptValues(data interface{}) ([]byte, error) {
        // 1. Marshal the interface to JSON (to handle different data structures)
        jsonData, err := json.Marshal(data)
        if err != nil {
                return nil, fmt.Errorf("failed to marshal interface to JSON: %w", err)
        }

        // 2. Unmarshal the JSON into a map[string]interface{}
        var dataMap map[string]interface{}
        err = json.Unmarshal(jsonData, &dataMap)
        if err != nil {
                return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
        }


        // 3. Encrypt only the values in the map
        encryptedMap, err := m.encryptMapValues(dataMap) // Helper function (see below)
        if err != nil {
                return nil, err
        }

        // 4. Marshal the map back to JSON
        encryptedJSON, err := json.Marshal(encryptedMap)
        if err != nil {
                return nil, fmt.Errorf("failed to marshal encrypted map to JSON: %w", err)
        }
        fmt.Println(string(encryptedJSON))

        return encryptedJSON, nil
}

func (m *CryptoMiddleware) encryptMapValues(dataMap map[string]interface{}) (map[string]interface{}, error) {
    encryptedMap := make(map[string]interface{})
    for key, value := range dataMap {
        switch v := value.(type) {
        case string:
            encryptedValue, err := m.encrypt([]byte(v))
            if err != nil {
                return nil, fmt.Errorf("error encrypting field %s: %w", key, err)
            }
            encryptedMap[key] = encryptedValue // Encrypt string value

        case map[string]interface{}: // Handle nested maps
            nestedEncryptedMap, err := m.encryptMapValues(v)
            if err != nil {
                return nil, err
            }
            encryptedMap[key] = nestedEncryptedMap

        case []interface{}: // Handle arrays
            encryptedArray := make([]interface{}, len(v))
            for i, item := range v {
                itemJSON, err := json.Marshal(item)
                if err != nil {
                    return nil, err
                }
                var itemMap map[string]interface{}
                err = json.Unmarshal(itemJSON, &itemMap)
                if err != nil {
                    // If not a map (e.g., string), encrypt directly
                    itemString, ok := item.(string)
                    if ok {
                        encryptedItem, err := m.encrypt([]byte(itemString))
                        if err != nil {
                            return nil, fmt.Errorf("error encrypting array item %d: %w", i, err)
                        }
                        encryptedArray[i] = encryptedItem
                    } else {
                        encryptedArray[i] = item // Keep non-string array items as they are
                    }
                    continue
                }
                encryptedItem, err := m.encryptMapValues(itemMap)
                if err != nil {
                    return nil, err
                }
                encryptedArray[i] = encryptedItem
            }
            encryptedMap[key] = encryptedArray

        default:
            // For other types (numbers, booleans, etc.), encrypt if needed
            // Or keep them as they are if no encryption is required
            strValue := fmt.Sprintf("%v", v)
            encryptedValue, err := m.encrypt([]byte(strValue))
            if err != nil {
                return nil, fmt.Errorf("error encrypting field %s: %w", key, err)
            }
            encryptedMap[key] = encryptedValue
        }
    }
    return encryptedMap, nil
}