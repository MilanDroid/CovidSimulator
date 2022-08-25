import matplotlib.pyplot as plt
import numpy as np

def addition(n):
    return n / 300

fig = plt.figure()
# x = np.random.normal(1, 5, 500)
# y = np.random.normal(1, 5, 500)

# size = np.random.normal(1, 60, 500)

# colour = np.arctan2(x, y)

x = (2, 4, 6, 8)
y = (4, 6, 9, 14)
size = (202, 325, 449, 444)
colour = list(map(addition, (236, 132, 7, 0)))

plt.scatter(x, y, s = size, c = colour, alpha = 0.8)
plt.colorbar()
plt.show()