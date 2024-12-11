import logging

def bcnt(n):
    return bin(n)[2:].count('1')

def parity(n1, n2):
    return (bcnt(n1) + bcnt(n2)) % 2

def blen(n):
    return n.bit_length()

def to_bits(v, nbits):
    mask = '{:0'+str(nbits) + 'b}'
    return mask.format(v)

def from_bits(bs):
    return int(bs, base=2)

def inverse_map(_map):
    result = [0] * len(_map)
    for i in range(len(_map)):
        result[_map[i]] = i
    return result

def step_pbox(p_box, v):
    bits = to_bits(v, len(p_box))
    reordered = ''.join([bits[p_box[i]] for i in range(len(p_box))])
    res = from_bits(reordered)
    #print("p_box({}) = {}".format(hex(v), hex(res)))
    return res


def step_sbox(s_box, v):
    l, r = (v & 0xff00 ) >> 8, v & 0xff
    l, r = s_box[l], s_box[r]
    res = l << 8 | r
    #print("s_box({}) = {}".format(hex(v), hex(res)))
    return res

def step_key(key, v):
    #print("{} ^ {} = {}".format(hex(key), hex(v), hex(v ^ key)))
    return v ^ key

def encrypt_block(plainblock, subkeys):
    v = plainblock
    assert len(subkeys) == 3

    # first round
    v = step_key(subkeys[0], v)
   
    # intermediate rounds
    for sk in subkeys[1:-1]:
        v = step_sbox(S_BOX, v)
        v = step_pbox(P_BOX, v)
        v = step_key(sk, v)
    
    # last round
    v = step_sbox(S_BOX, v)
    cipherblock = step_key(subkeys[-1], v)
    return cipherblock

def decrypt_block(cipherblock, subkeys):
    v = cipherblock
    
    # undo last round
    v = step_key(subkeys[-1], v)
    v = step_sbox(INV_S_BOX, v)
    
    # undo intermediate rounds
    for sk in reversed(subkeys[1:-1]):
        v = step_key(sk, v)
        v = step_pbox(INV_P_BOX, v)
        v = step_sbox(INV_S_BOX, v)

    # undo first round
    plainblock = step_key(subkeys[0], v)
    return plainblock

S_BOX = [62, 117, 195, 179, 20, 210, 41, 66, 116, 178, 152, 143, 75, 105, 254, 1, 158, 95, 101, 175, 191, 166, 36, 24, 50, 39, 190, 120, 52, 242, 182, 185, 61, 225, 140, 38, 150, 80, 19, 109, 246, 252, 40, 13, 65, 236, 124, 186, 214, 86, 235, 100, 97, 49, 197, 154, 176, 199, 253, 69, 88, 112, 139, 77, 184, 45, 133, 104, 15, 54, 177, 244, 160, 169, 82, 148, 73, 30, 229, 35, 79, 137, 157, 180, 248, 163, 241, 231, 81, 94, 165, 9, 162, 233, 18, 85, 217, 84, 7, 55, 63, 171, 56, 118, 237, 132, 136, 22, 90, 221, 103, 161, 205, 11, 255, 14, 122, 47, 71, 201, 99, 220, 83, 74, 173, 76, 144, 16, 155, 126, 60, 96, 44, 234, 17, 215, 107, 138, 159, 183, 251, 3, 198, 0, 89, 170, 131, 151, 219, 29, 230, 32, 187, 125, 134, 64, 12, 202, 164, 247, 25, 223, 222, 119, 174, 67, 147, 146, 206, 51, 243, 53, 121, 239, 68, 130, 70, 203, 211, 111, 108, 113, 8, 106, 57, 240, 21, 93, 142, 238, 167, 5, 128, 72, 189, 192, 193, 92, 10, 204, 87, 145, 188, 172, 224, 226, 207, 27, 218, 48, 33, 28, 123, 6, 37, 59, 4, 102, 114, 91, 23, 209, 34, 42, 2, 196, 141, 208, 181, 245, 43, 78, 213, 216, 232, 46, 98, 26, 212, 58, 115, 194, 200, 129, 227, 249, 127, 149, 135, 228, 31, 153, 250, 156, 168, 110]
S_BOX_BITS = len(S_BOX).bit_length() - 1
SBOX_SIZE = 2 ** S_BOX_BITS
INV_S_BOX = inverse_map(S_BOX)
assert inverse_map(INV_S_BOX) == S_BOX

S_BOX_BASIC = [x for x in range(0,256)]


# Mapping the reordering of the bits'
# Bit 0 becomes bit 0 
# Bit 1 becomes bit 8
# Bit 2 becomes bit 2
# ...
P_BOX = [0, 8, 2, 10, 4, 12, 6, 14, 1, 9, 3, 11, 5, 13, 7, 15]
INV_P_BOX = inverse_map(P_BOX)
assert inverse_map(INV_P_BOX) == P_BOX

for i in range(16):
    for j in range(16):
        for k in range(16):
            for l in range(16):
                c = encrypt_block(i, (j, k, l))
                p = decrypt_block(c, (j, k, l))
                assert i == p