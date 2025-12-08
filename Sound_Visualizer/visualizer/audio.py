import soundfile as sf
import numpy as np

def load_wav(file_path, chunk_count):
    samples, samplerate = sf.read(file_path)

    # Convert stereo → mono
    if samples.ndim > 1:
        samples = samples.mean(axis=1)

    frames = len(samples)
    duration = frames / samplerate

    chunk_size = frames // chunk_count
    chunks = []

    for i in range(chunk_count):
        chunk = samples[i * chunk_size : (i + 1) * chunk_size]

        if len(chunk) == 0:
            chunk = np.zeros(128)

        # Normalize for consistency
        chunk = chunk / max(1e-9, np.max(np.abs(chunk)))

        chunks.append(chunk.astype(float))

    return chunks, duration