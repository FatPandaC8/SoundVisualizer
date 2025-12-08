import pygame
from collections import deque
import numpy as np

SCREEN_WIDTH = 800
SCREEN_HEIGHT = 400
SCALE = 120        # amplitude scale
NUM_BARS = 200     # number of vertical bars

def render_waveform(screen, bar_heights: deque):
    """
    Draw static vertical bars and a wave line that beats up and down like a sine wave.
    bar_heights: deque of floats between -1 and 1
    """
    screen.fill((30, 30, 30))  # dark background
    mid_y = SCREEN_HEIGHT // 2
    spacing = SCREEN_WIDTH / NUM_BARS

    # Ensure bar_heights is length NUM_BARS
    if len(bar_heights) != NUM_BARS:
        indices = np.linspace(0, len(bar_heights)-1, NUM_BARS).astype(int)
        wave = np.array(bar_heights)[indices]
    else:
        wave = np.array(bar_heights)

    # Optional smoothing for sine-like motion
    wave_smooth = np.convolve(wave, np.ones(5)/5, mode='same')

    # Draw static bars (thin grid lines)
    for i in range(NUM_BARS):
        x = int(i * spacing)
        pygame.draw.line(screen, (50, 50, 50), (x, mid_y), (x, mid_y), 1)

    # Draw wave line connecting tops of bars
    points = [(int(i * spacing), int(mid_y - wave_smooth[i] * SCALE)) for i in range(NUM_BARS)]
    pygame.draw.lines(screen, (0, 200, 255), False, points, 3)
