import pygame
import time
from collections import deque
from .audio import load_wav  # your existing audio loader
from .render import render_waveform, NUM_BARS, SCREEN_WIDTH, SCREEN_HEIGHT, SCALE

FPS = 60  # frames per second

class Visualizer:
    def __init__(self, wav_path, chunks=2000):
        self.wav_path = wav_path
        self.chunks_count = chunks
        self.chunks = []
        self.duration = 0
        self.start_time = 0
        self.bar_heights = deque([0]*NUM_BARS, maxlen=NUM_BARS)

    def init(self):
        pygame.init()
        pygame.mixer.init()

        self.screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
        pygame.display.set_caption("Oscilloscope Visualizer")

        # Load waveform chunks
        self.chunks, self.duration = load_wav(self.wav_path, self.chunks_count)

        # Load + play audio
        pygame.mixer.music.load(self.wav_path)
        pygame.mixer.music.play()

        self.start_time = time.time()

    def update(self):
        """Update bar heights with small slices of audio for smooth scrolling."""
        elapsed = time.time() - self.start_time
        idx = int((elapsed / self.duration) * len(self.chunks))

        if 0 <= idx < len(self.chunks):
            # Take a small slice per frame, compute max amplitude for bar
            slice_samples = self.chunks[idx][:10]
            max_amp = max(slice_samples)
            self.bar_heights.append(max_amp)

    def render(self):
        render_waveform(self.screen, self.bar_heights)

    def run(self):
        clock = pygame.time.Clock()
        running = True

        while running:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    running = False

            self.update()
            self.render()
            pygame.display.flip()
            clock.tick(FPS)

        pygame.quit()
