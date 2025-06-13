#include "visualizer.h"

Visualizer::Visualizer() : vWindow(nullptr), vRenderer(nullptr) {

}

Visualizer::~Visualizer(){
    //close();
}

bool Visualizer::init(){
    if(SDL_Init(SDL_INIT_EVERYTHING) < 0){
        std::cerr << "SDL could not initialized! SDL Error: " << SDL_GetError() << '\n';
        return false;
    }

    vWindow = SDL_CreateWindow("Music Visualizer", SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, SCREEN_WIDTH, SCREEN_HEIGHT, SDL_WINDOW_SHOWN);

    if(!vWindow){
        std::cerr << "Window could not be created! SDL Error: " << SDL_GetError() << '\n';
        return false;
    }

    vRenderer = SDL_CreateRenderer(vWindow, -1, SDL_RENDERER_ACCELERATED);
    if(!vRenderer){
        std::cerr << "Renderer could not be created! SDL Error: " << SDL_GetError() << '\n';
        return false;
    }

    SDL_SetRenderDrawColor(vRenderer, 0xFF, 0xFF, 0xFF, 0xFF);

    bars = convert("resource/NightcoreHero.wav", 1024);

    if (Mix_OpenAudio(44100, MIX_DEFAULT_FORMAT, 2, 2048) < 0) {
        std::cout << "SDL_mixer could not initialize! SDL_mixer Error: " << Mix_GetError() << "\n";
        return false;
    }

    music = Mix_LoadMUS("resource/NightcoreHero.wav");
    if (!music) {
        std::cerr << "Failed to load music: " << Mix_GetError() << '\n';
        return false;
    }

    if (Mix_PlayMusic(music, -1) < 0) {
        std::cerr << "Mix_PlayMusic failed: " << Mix_GetError() << "\n";
        return false;
    }

    musicStartTime = SDL_GetTicks();

    return true;
}

std::vector<float> Visualizer::convert(const char* filePath, int k){
    SF_INFO sfInfo;
    SNDFILE* file = sf_open(filePath, SFM_READ, &sfInfo);
    if(!file){
        std::cerr << "Failed to open file: " << sf_strerror(file) << '\n';
        return {};
    }

    std::vector<float> samples(sfInfo.frames * sfInfo.channels);
    sf_readf_float(file, samples.data(), sfInfo.frames);
    sf_close(file);

    //convert to mono for visualization because if stereo there are 2 channels -> hard to visualize both at the same time
    if (sfInfo.channels > 1) {
        std::vector<float> mono(samples.size() / sfInfo.channels);
        for (int i = 0, j = 0; i < samples.size(); i += sfInfo.channels, ++j) {
            float sum = 0;
            for (int c = 0; c < sfInfo.channels; ++c)
                sum += samples[i + c];
            mono[j] = sum / sfInfo.channels; //take the average of left channels and right channels as they are interweave in the samples vector
        }
        samples = mono;
    }

    int chunkSize = samples.size() / k;
    std::vector<float> bars(k, 0.0f);

    for (int i = 0; i < k; ++i) { //get the max value at each chunk, can use sliding window max here because we need to find max in each chunk size ChunkSize
        float maxVal = 0.0f;
        for (int j = i * chunkSize; j < (i + 1) * chunkSize && j < samples.size(); ++j) {
            maxVal = std::max(maxVal, std::abs(samples[j]));
        }
        bars[i] = maxVal;
    }

    songDuration = static_cast<float>(sfInfo.frames) / sfInfo.samplerate;

    return bars;
}

void Visualizer::run(){
    constexpr int FRAME_DELAY = 1000 / 60; 

    Uint32 lastTime = SDL_GetTicks();
    Uint32 frameStart = SDL_GetTicks();
    Uint32 currentTime = SDL_GetTicks();
    float deltaTime = (currentTime - lastTime) / 1000.0f;

    SDL_Event event;

    while(running){
        while(SDL_PollEvent(&event)){
            if(event.type == SDL_QUIT){
                running = false;
            }
        }

        update();
        render();
        SDL_Delay(16);
    }

    controlFrame(frameStart, FRAME_DELAY);
}

void Visualizer::controlFrame(Uint32 frameStart, int frameDelay) {
    int frameTime = SDL_GetTicks() - frameStart;
    if (frameTime < frameDelay) {
        SDL_Delay(frameDelay - frameTime);
    }
}

void Visualizer::render(){
    SDL_SetRenderDrawColor(vRenderer,0,0,0,255);
    SDL_RenderClear(vRenderer);
    SDL_Rect screen = {0,0,SCREEN_WIDTH, SCREEN_HEIGHT};

    SDL_SetRenderDrawColor(vRenderer, 255, 255, 255, 255);
    SDL_RenderFillRect(vRenderer, &screen);
    
    //----ADD THE VERTICAL BARS----//
    for (int i = 0; i < draw.size(); ++i) {
        float barValue = draw[i];
        int barHeight = static_cast<int>(barValue * 100); 
        SDL_Rect bar = {
            i*5,                        
            SCREEN_HEIGHT / 2 - barHeight, 
            3,                             
            barHeight * 2                      
        };
        SDL_SetRenderDrawColor(vRenderer, 0, 0, 255, 255); // blue
        SDL_RenderFillRect(vRenderer, &bar);
    }
    
    SDL_RenderPresent(vRenderer);
}

void Visualizer::update(){
    Uint32 now = SDL_GetTicks();

    float secondsPlayed = (now - musicStartTime) / 1000.0f;
    int barIndex = static_cast<int>((secondsPlayed / songDuration) * bars.size());
    //the songDuration is divided into bars.size() chunk -> use secondsplayed to convert from percentages to index

    if (barIndex < bars.size() && barIndex != frameIndex) {
        frameIndex = barIndex;
        draw.push_back(bars[frameIndex]);
        if (draw.size() > BUFFER_SIZE)
            draw.pop_front();
    }
}