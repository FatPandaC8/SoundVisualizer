#ifndef DEFS_H
#define DEFS_H

#include <SDL2/SDL.h>
#include <SDL2/SDL_mixer.h>
#include <sndfile.h>
#include <iostream>
#include <vector>
#include <deque>
#include <math.h>

constexpr int SCREEN_WIDTH = 1280; //compile time constant
constexpr int SCREEN_HEIGHT = 720;
const int barSpacing = 5;
constexpr int BUFFER_SIZE = (SCREEN_WIDTH - SCREEN_HEIGHT / 3) / barSpacing;

#endif