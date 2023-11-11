import time
from tqdm import tqdm, trange

for i in trange(100, desc="Training", unit="epoch"):
    time.sleep(0.1)

for i in tqdm(range(100), desc="Training", unit="epoch"):
    time.sleep(0.1)

if __name__ == '__main__':
    pass
