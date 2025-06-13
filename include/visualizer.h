#ifndef VISUALIZER_H
#define VISUALIZER_H

#include "defs.h"
class Visualizer {
    public:
    Visualizer();
    ~Visualizer();
    bool init();

    void run();
    void update();
    void render();
    void controlFrame(Uint32 frameStart, int frameDelay);
    void close();
    std::vector<float> convert(const char* filePath, int k);

    private:
    bool running = true;
    SDL_Window* vWindow = nullptr;
    SDL_Renderer* vRenderer = nullptr;
    std::vector<float> bars;
    std::deque<float> draw;
    int frameIndex = 0;
    Uint32 musicStartTime = 0;
    float songDuration = 0.0f;

    Mix_Music* music = nullptr;
};

#endif