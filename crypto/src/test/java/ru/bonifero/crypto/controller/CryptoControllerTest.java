package ru.bonifero.crypto.controller;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest
@AutoConfigureMockMvc
class CryptoControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @Test
    void hash_returnsSha256HexDigest() throws Exception {
        mockMvc.perform(post("/api/crypto/hash")
                        .contentType("application/json")
                        .content("{\"message\":\"Hello, world!\",\"algorithm\":\"SHA_256\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.algorithm").value("SHA-256"))
                .andExpect(jsonPath("$.hashHex").isNotEmpty());
    }

    @Test
    void encryptThenDecrypt_viaHttp_roundTrips() throws Exception {
        String encryptBody = mockMvc.perform(post("/api/crypto/encrypt")
                        .contentType("application/json")
                        .content("{\"message\":\"secret payload\"}"))
                .andExpect(status().isOk())
                .andReturn().getResponse().getContentAsString();

        JsonNode encrypted = objectMapper.readTree(encryptBody);

        String decryptRequest = objectMapper.writeValueAsString(new Object() {
            public final String encryptedKey = encrypted.get("encryptedKey").asText();
            public final String iv = encrypted.get("iv").asText();
            public final String cipherText = encrypted.get("cipherText").asText();
        });

        mockMvc.perform(post("/api/crypto/decrypt")
                        .contentType("application/json")
                        .content(decryptRequest))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.message").value("secret payload"));
    }

    @Test
    void encrypt_rejectsBlankMessage() throws Exception {
        mockMvc.perform(post("/api/crypto/encrypt")
                        .contentType("application/json")
                        .content("{\"message\":\"\"}"))
                .andExpect(status().isBadRequest());
    }

    @Test
    void decrypt_rejectsInvalidData() throws Exception {
        mockMvc.perform(post("/api/crypto/decrypt")
                        .contentType("application/json")
                        .content("{\"encryptedKey\":\"not-valid\",\"iv\":\"not-valid\",\"cipherText\":\"not-valid\"}"))
                .andExpect(status().isBadRequest());
    }

    @Test
    void logs_returnsEntryAfterHashOperation() throws Exception {
        mockMvc.perform(post("/api/crypto/hash")
                        .contentType("application/json")
                        .content("{\"message\":\"logged message\",\"algorithm\":\"SHA_256\"}"))
                .andExpect(status().isOk());

        mockMvc.perform(get("/api/crypto/logs").param("type", "HASH").param("page", "0").param("size", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content").isArray())
                .andExpect(jsonPath("$.content[0].operationType").value("HASH"));
    }

    @Test
    void logById_returnsNotFoundForMissingId() throws Exception {
        mockMvc.perform(get("/api/crypto/logs/999999999"))
                .andExpect(status().isNotFound());
    }
}
