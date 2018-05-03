#include "libsecp256k1-config.h"
#include "secp256k1.c"
#include "src/modules/recovery/main_impl.h"
#include "ext.h"

int main() {
  secp256k1_context* ctx = secp256k1_context_create_sign_verify();
  const int MSG_SIZE = 32;
  const int KEY_SIZE = 32;
  char seckey[KEY_SIZE+1];
  char msg[MSG_SIZE+1];
  seckey[KEY_SIZE] = '\0';
  msg[MSG_SIZE] = '\0';
  read(0, seckey, KEY_SIZE);
  read(0, msg, MSG_SIZE);
  secp256k1_ec_seckey_verify(ctx, seckey);
  secp256k1_ecdsa_recoverable_signature* sigstruct = calloc(1, sizeof(secp256k1_ecdsa_recoverable_signature));
  secp256k1_ecdsa_sign_recoverable(ctx, &sigstruct, msg, seckey, secp256k1_nonce_function_rfc6979, NULL);
  return 0;
}
