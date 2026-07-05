package ru.bonifero.crypto.crypto;

import org.springframework.stereotype.Service;
import ru.bonifero.crypto.exception.CryptoServiceException;

import javax.crypto.Cipher;
import javax.crypto.KeyGenerator;
import javax.crypto.SecretKey;
import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.GCMParameterSpec;
import javax.crypto.spec.PBEKeySpec;
import javax.crypto.spec.SecretKeySpec;
import java.security.MessageDigest;
import java.security.PrivateKey;
import java.security.Provider;
import java.security.PublicKey;
import java.security.SecureRandom;
import java.util.Base64;
import java.util.HexFormat;

@Service
public class CryptoService {

    private static final String AES_ALGORITHM = "AES";
    private static final String AES_CIPHER = "AES/GCM/NoPadding";
    private static final int AES_KEY_SIZE_BITS = 256;
    private static final int GCM_IV_LENGTH_BYTES = 12;
    private static final int GCM_TAG_LENGTH_BITS = 128;

    private static final String RSA_CIPHER = "RSA/ECB/OAEPWithSHA-256AndMGF1Padding";

    private static final int PBKDF2_ITERATIONS = 100_000;
    private static final int PBKDF2_KEY_LENGTH_BITS = 256;
    private static final int PBKDF2_SALT_LENGTH_BYTES = 16;

    private final PublicKey publicKey;
    private final PrivateKey privateKey;
    private final Provider provider;
    private final SecureRandom secureRandom = new SecureRandom();

    public CryptoService(PublicKey publicKey, PrivateKey privateKey, Provider provider) {
        this.publicKey = publicKey;
        this.privateKey = privateKey;
        this.provider = provider;
    }

    public HybridEncryptionResult encrypt(String message) {
        try {
            KeyGenerator keyGenerator = KeyGenerator.getInstance(AES_ALGORITHM, provider);
            keyGenerator.init(AES_KEY_SIZE_BITS, secureRandom);
            SecretKey aesKey = keyGenerator.generateKey();

            byte[] iv = new byte[GCM_IV_LENGTH_BYTES];
            secureRandom.nextBytes(iv);

            Cipher aesCipher = Cipher.getInstance(AES_CIPHER, provider);
            aesCipher.init(Cipher.ENCRYPT_MODE, aesKey, new GCMParameterSpec(GCM_TAG_LENGTH_BITS, iv));
            byte[] cipherText = aesCipher.doFinal(message.getBytes(java.nio.charset.StandardCharsets.UTF_8));

            Cipher rsaCipher = Cipher.getInstance(RSA_CIPHER, provider);
            rsaCipher.init(Cipher.ENCRYPT_MODE, publicKey);
            byte[] encryptedKey = rsaCipher.doFinal(aesKey.getEncoded());

            return new HybridEncryptionResult(
                    Base64.getEncoder().encodeToString(encryptedKey),
                    Base64.getEncoder().encodeToString(iv),
                    Base64.getEncoder().encodeToString(cipherText));
        } catch (Exception e) {
            throw new CryptoServiceException("Ошибка при шифровании сообщения", e);
        }
    }

    public String decrypt(HybridEncryptionResult encrypted) {
        try {
            Cipher rsaCipher = Cipher.getInstance(RSA_CIPHER, provider);
            rsaCipher.init(Cipher.DECRYPT_MODE, privateKey);
            byte[] aesKeyBytes = rsaCipher.doFinal(Base64.getDecoder().decode(encrypted.encryptedKey()));
            SecretKey aesKey = new SecretKeySpec(aesKeyBytes, AES_ALGORITHM);

            byte[] iv = Base64.getDecoder().decode(encrypted.iv());
            Cipher aesCipher = Cipher.getInstance(AES_CIPHER, provider);
            aesCipher.init(Cipher.DECRYPT_MODE, aesKey, new GCMParameterSpec(GCM_TAG_LENGTH_BITS, iv));

            byte[] cipherText = Base64.getDecoder().decode(encrypted.cipherText());
            byte[] plainText = aesCipher.doFinal(cipherText);
            return new String(plainText, java.nio.charset.StandardCharsets.UTF_8);
        } catch (Exception e) {
            throw new CryptoServiceException("Ошибка при дешифровании сообщения. Проверьте корректность данных.", e);
        }
    }

    public HashResult hash(String message, HashAlgorithm algorithm) {
        return switch (algorithm) {
            case SHA_256 -> sha256(message);
            case PBKDF2 -> pbkdf2(message);
        };
    }

    private HashResult sha256(String message) {
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256", provider);
            byte[] hash = digest.digest(message.getBytes(java.nio.charset.StandardCharsets.UTF_8));
            return new HashResult("SHA-256", HexFormat.of().formatHex(hash), null, null);
        } catch (Exception e) {
            throw new CryptoServiceException("Ошибка при расчете SHA-256 хеша", e);
        }
    }

    private HashResult pbkdf2(String message) {
        try {
            byte[] salt = new byte[PBKDF2_SALT_LENGTH_BYTES];
            secureRandom.nextBytes(salt);

            PBEKeySpec spec = new PBEKeySpec(
                    message.toCharArray(), salt, PBKDF2_ITERATIONS, PBKDF2_KEY_LENGTH_BITS);
            SecretKeyFactory factory = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256");
            byte[] hash = factory.generateSecret(spec).getEncoded();

            return new HashResult("PBKDF2WithHmacSHA256", HexFormat.of().formatHex(hash),
                    HexFormat.of().formatHex(salt), PBKDF2_ITERATIONS);
        } catch (Exception e) {
            throw new CryptoServiceException("Ошибка при расчете PBKDF2 хеша", e);
        }
    }
}
