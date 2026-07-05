package ru.bonifero.crypto.exception;

public class CryptoServiceException extends RuntimeException {

    public CryptoServiceException(String message, Throwable cause) {
        super(message, cause);
    }
}
