# Constants and mappings based on your previously shared or hypothetical S-Box and P-Box
S_BOX = [3, 0, 2, 1]  # Simplified S-Box for illustration; please use your actual S-Box
INV_S_BOX = [0] * len(S_BOX)  # Correct size based on S-Box

# Invert the S-Box
for i in range(len(S_BOX)):
    INV_S_BOX[S_BOX[i]] = i

P_BOX = [0, 2, 1, 3]  # Simplified P-Box for illustration; replace with actual P-Box
INV_P_BOX = [0] * len(P_BOX)  # Correct size based on P-Box

# Invert the P-Box
for i in range(len(P_BOX)):
    INV_P_BOX[P_BOX[i]] = i

def to_bits(value, length):
    return format(value, '0{}b'.format(length))

def from_bits(bits):
    return int(bits, 2)

def apply_sbox(value, sbox):
    """ Apply the S-Box to a 4-bit integer (for simplicity). """
    return sbox[value]

def apply_pbox(value, pbox):
    bits = to_bits(value, len(pbox))
    permuted = ''.join(bits[i] for i in pbox)
    return from_bits(permuted)

def encrypt(plaintext, keys):
    assert len(keys) == 3
    v = plaintext
    # Initial key mixing
    v ^= keys[0]
    # Rounds (simplified for illustration)
    v = apply_sbox(v, S_BOX)
    v = apply_pbox(v, P_BOX)
    v ^= keys[1]
    # Assuming a very simple round structure for demonstration
    return v

def decrypt(ciphertext, keys):
    v = ciphertext
    v ^= keys[1]
    v = apply_pbox(v, INV_P_BOX)
    v = apply_sbox(v, INV_S_BOX)
    v ^= keys[0]
    return v

def load_data(filename):
    """ Load plaintext and ciphertext pairs from a file. """
    plaintexts = []
    ciphertexts = []
    with open(filename, 'r') as file:
        for line in file:
            if line.strip():
                plain, cipher = line.strip().split(':')
                plaintexts.append(int(plain, 16))
                ciphertexts.append(int(cipher, 16))
    return plaintexts, ciphertexts

def brute_force_key(plaintexts, ciphertexts):
    best_keys = []
    max_score = 0
    # Simplified key range for illustration
    for k1 in range(4):
        for k2 in range(4):
            keys = (k1, k2)
            score = 0
            for pt, ct in zip(plaintexts, ciphertexts):
                if decrypt(ct, keys) == pt:
                    score += 1
            if score > max_score:
                max_score = score
                best_keys = [keys]
            elif score == max_score:
                best_keys.append(keys)
    return best_keys, max_score

# Load data from file
plaintexts, ciphertexts = load_data('cipher.txt')

# Try to brute-force the key
best_keys, max_score = brute_force_key(plaintexts, ciphertexts)
print("Best keys found:", best_keys)
print("Maximum score (number of correct decryptions):", max_score)
