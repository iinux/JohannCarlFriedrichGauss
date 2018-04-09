<?php
/**
 * Created by PhpStorm.
 * User: qzhang
 * Date: 2018/4/9
 * Time: 16:54
 */

if (PHP_VERSION_ID >= 70100) {
    //$key should have been previously generated in a cryptographically safe way, like openssl_random_pseudo_bytes
    $plaintext = "message to be encrypted";
    $cipher = "aes-128-gcm";
    if (in_array($cipher, openssl_get_cipher_methods()))
    {
        $ivlen = openssl_cipher_iv_length($cipher);
        $iv = openssl_random_pseudo_bytes($ivlen);
        $key = openssl_random_pseudo_bytes(128, $cstrong);
        $ciphertext = openssl_encrypt($plaintext, $cipher, $key, $options=0, $iv, $tag);
        //store $cipher, $iv, and $tag for decryption later
        $original_plaintext = openssl_decrypt($ciphertext, $cipher, $key, $options=0, $iv, $tag);
        echo $original_plaintext."\n";
    }
} elseif (PHP_VERSION_ID >= 50600) {
    //$key previously generated safely, ie: openssl_random_pseudo_bytes
    $plaintext = "message to be encrypted";
    $ivlen = openssl_cipher_iv_length($cipher = "AES-128-CBC");
    $iv = openssl_random_pseudo_bytes($ivlen);
    $key = openssl_random_pseudo_bytes(128, $cstrong);
    $ciphertext_raw = openssl_encrypt($plaintext, $cipher, $key, $options = OPENSSL_RAW_DATA, $iv);
    $hmac = hash_hmac('sha256', $ciphertext_raw, $key, $as_binary = true);
    $ciphertext = base64_encode($iv . $hmac . $ciphertext_raw);

    //decrypt later....
    $c = base64_decode($ciphertext);
    $ivlen = openssl_cipher_iv_length($cipher = "AES-128-CBC");
    $iv = substr($c, 0, $ivlen);
    $hmac = substr($c, $ivlen, $sha2len = 32);
    $ciphertext_raw = substr($c, $ivlen + $sha2len);
    $original_plaintext = openssl_decrypt($ciphertext_raw, $cipher, $key, $options = OPENSSL_RAW_DATA, $iv);
    $calcmac = hash_hmac('sha256', $ciphertext_raw, $key, $as_binary = true);
    if (hash_equals($hmac, $calcmac))//PHP 5.6+ timing attack safe comparison
    {
        echo $original_plaintext . "\n";
    }
} else {

}