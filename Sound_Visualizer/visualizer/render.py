import pygame
import numpy as np

SCREEN_WIDTH = 800
SCREEN_HEIGHT = 400
SCALE = 120  # line height

def render_waveform(screen, wave):
    screen.fill((30, 30, 30))  # dark background

    if len(wave) < 2:
        return

    mid_y = SCREEN_HEIGHT // 2
    TARGET_POINTS = 200  # space out points
    indices = np.linspace(0, len(wave) - 1, TARGET_POINTS).astype(int)
    wave_smooth = wave[indices]

    # Optional smoothing with moving average
    wave_smooth = np.convolve(wave_smooth, np.ones(5)/5, mode='same')

    spacing = SCREEN_WIDTH / (TARGET_POINTS - 1)

    last_x = 0
    last_y = int(mid_y - wave_smooth[0] * SCALE)

    for i in range(1, TARGET_POINTS):
        x = int(i * spacing)
        y = int(mid_y - wave_smooth[i] * SCALE)
        # Draw single wave line
        pygame.draw.line(screen, (0, 200, 255), (last_x, last_y), (x, y), 2)
        last_x, last_y = x, y