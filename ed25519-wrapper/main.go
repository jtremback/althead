package ed25519

import "github.com/agl/ed25519"

func Verify(publicKey []byte, message []byte, sig []byte) bool {
	var pub [ed25519.PublicKeySize]byte
	copy(pub[:], publicKey)

	var sg [ed25519.SignatureSize]byte
	copy(sg[:], sig)

	return ed25519.Verify(&pub, message, &sg)
}

func Sign(privateKey []byte, message []byte) []byte {
	var priv [ed25519.PrivateKeySize]byte
	copy(priv[:], privateKey)

	return ed25519.Sign(&priv, message)[:]
}
