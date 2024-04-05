import rsa
import jwt


def generate_keys():
    public_key, private_key = rsa.newkeys(512)
    with open("public.pem", "wb") as f:
        f.write(public_key.save_pkcs1())
    with open("private.pem", "wb") as f:
        f.write(private_key.save_pkcs1())

    print("Public and private keys generated and saved")

    return public_key, private_key


def load_keys(public_key_path, private_key_path):
    with open(public_key_path, "rb") as f:
        public_key = rsa.PublicKey.load_pkcs1(f.read())
    with open(private_key_path, "rb") as f:
        private_key = rsa.PrivateKey.load_pkcs1(f.read())

    print("Public and private keys loaded")

    return public_key, private_key


def encrypt(data, public_key):
    return rsa.encrypt(data.encode(), public_key)


def decrypt(data, private_key):
    return rsa.decrypt(data, private_key)


def verify_token(token, secret):
    return jwt.decode(token, secret, algorithms=["HS256"])


def generate_token(data, secret):
    return jwt.encode(data, secret, algorithm="HS256")
