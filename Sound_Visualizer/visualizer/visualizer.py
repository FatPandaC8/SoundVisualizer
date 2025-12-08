import pygame
import time
from collections import deque
from .audio import load_wav
from .render import render_waveform

SCREEN_WIDTH = 800
SCREEN_HEIGHT = 400
BUFFER_SIZE = 160
CHUNKS = 2000   # more chunks = smoother, slower scan

class Visualizer:
    def __init__(self, wav_path):
        self.wav_path = wav_path
        self.chunks = []
        self.duration = 0

        self.current_index = 0
        self.start_time = 0

    def init(self):
        pygame.init()
        pygame.mixer.init()

        self.screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
        pygame.display.set_caption("Oscilloscope Visualizer")

        # Load waveform chunks
        self.chunks, self.duration = load_wav(self.wav_path, CHUNKS)

        # Load + play audio
        pygame.mixer.music.load(self.wav_path)
        pygame.mixer.music.play()

        self.start_time = time.time()

    def update(self):
        elapsed = time.time() - self.start_time

        idx = int((elapsed / self.duration) * len(self.chunks))

        if 0 <= idx < len(self.chunks):
            self.current_index = idx
            return self.chunks[idx]

        return None

    def run(self):
        clock = pygame.time.Clock()
        running = True

        while running:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    running = False

            wave_chunk = self.update()

            if wave_chunk is not None:
                render_waveform(self.screen, wave_chunk)

            pygame.display.flip()
            clock.tick(60)

        pygame.quit()