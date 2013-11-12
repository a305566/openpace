package main

import (
	"fmt"
    "os"
    "pace"
)

func main() {

    EF_CARDACCESS := "\x31\x81\x82\x30\x0D\x06\x08\x04\x00\x7F\x00\x07\x02\x02\x02\x02\x01\x02\x30\x12\x06\x0A\x04\x00\x7F\x00\x07\x02\x02\x03\x02\x02\x02\x01\x02\x02\x01\x41\x30\x12\x06\x0A\x04\x00\x7F\x00\x07\x02\x02\x04\x02\x02\x02\x01\x02\x02\x01\x0D\x30\x1C\x06\x09\x04\x00\x7F\x00\x07\x02\x02\x03\x02\x30\x0C\x06\x07\x04\x00\x7F\x00\x07\x01\x02\x02\x01\x0D\x02\x01\x41\x30\x2B\x06\x08\x04\x00\x7F\x00\x07\x02\x02\x06\x16\x1F\x65\x50\x41\x20\x2D\x20\x42\x44\x72\x20\x47\x6D\x62\x48\x20\x2D\x20\x54\x65\x73\x74\x6B\x61\x72\x74\x65\x20\x76\x32\x2E\x30\x04\x49\x17\x15\x41\x19\x28\x80\x0A\x01\xB4\x21\xFA\x07\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x20\x10\x10\x29\x10\x10"
    PIN := "123456"

    pace.EAC_init()

    secret := pace.PACE_SEC_new(PIN, len(PIN), pace.PACE_PIN)

    buf := pace.Get_buf(EF_CARDACCESS, len(EF_CARDACCESS))
    pace.Hexdump("EF.CardAccess", buf)

    /*fmt.Println("Secret:")*/
    /*pace.PACE_SEC_print_private(secret, 4)*/

    picc_ctx := pace.EAC_CTX_new()
    pcd_ctx := pace.EAC_CTX_new()
    pace.EAC_CTX_init_ef_cardaccess(EF_CARDACCESS, len(EF_CARDACCESS), pcd_ctx)
    pace.EAC_CTX_init_ef_cardaccess(EF_CARDACCESS, len(EF_CARDACCESS), picc_ctx)

    fmt.Println("PACE step 1")
    enc_nonce := pace.PACE_STEP1_enc_nonce(picc_ctx, secret)

    fmt.Println("PACE step 2")
    pace.PACE_STEP2_dec_nonce(pcd_ctx, secret, enc_nonce)

    fmt.Println("PACE step 3A")
    pcd_mapping_data := pace.PACE_STEP3A_generate_mapping_data(pcd_ctx)
    picc_mapping_data := pace.PACE_STEP3A_generate_mapping_data(picc_ctx)

    pace.PACE_STEP3A_map_generator(pcd_ctx, picc_mapping_data)
    pace.PACE_STEP3A_map_generator(picc_ctx, pcd_mapping_data)

    fmt.Println("PACE step 3B")
    pcd_ephemeral_pubkey := pace.PACE_STEP3B_generate_ephemeral_key(pcd_ctx)
    picc_ephemeral_pubkey := pace.PACE_STEP3B_generate_ephemeral_key(picc_ctx)

    pace.PACE_STEP3B_compute_shared_secret(pcd_ctx, picc_ephemeral_pubkey)
    pace.PACE_STEP3B_compute_shared_secret(picc_ctx, pcd_ephemeral_pubkey)

    fmt.Println("PACE step 3C")
    pace.PACE_STEP3C_derive_keys(pcd_ctx)
    pace.PACE_STEP3C_derive_keys(picc_ctx)

    fmt.Println("PACE step 3D")
    pcd_token := pace.PACE_STEP3D_compute_authentication_token(pcd_ctx, picc_ephemeral_pubkey)
    picc_token := pace.PACE_STEP3D_compute_authentication_token(picc_ctx, pcd_ephemeral_pubkey)

    pace.PACE_STEP3D_verify_authentication_token(pcd_ctx, picc_token)
    r := pace.PACE_STEP3D_verify_authentication_token(picc_ctx, pcd_token)

    /*fmt.Println("PICC's EAC_CTX:")*/
    /*pace.EAC_CTX_print_private(picc_ctx, 4)*/
    /*fmt.Println("PCD's EAC_CTX:")*/
    /*pace.EAC_CTX_print_private(pcd_ctx, 4)*/

    pace.EAC_CTX_clear_free(pcd_ctx)
    pace.EAC_CTX_clear_free(picc_ctx)
    pace.PACE_SEC_clear_free(secret)

    pace.EAC_cleanup()

    if r != 1 {
        os.Exit(1)
    }
}