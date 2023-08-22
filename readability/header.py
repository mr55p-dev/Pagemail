#! ./readability/venv/bin/python

import sys

if __name__ != "__main__":
    sys.exit(1)

data = sys.stdin.buffer.read()
l = len(data)
hdr = bytearray([(l >> i * 8) & 0xFF for i in [3, 2, 1, 0]])
sys.stdout.buffer.write(hdr)
sys.stdout.buffer.write(data)
sys.stdout.buffer.flush()
