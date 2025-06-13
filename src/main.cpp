#include "visualizer.h"

int main(int argc, char* argv[]){
    Visualizer visual;
    if(!visual.init()){
        return -1;
    }

    visual.run();
    return 0;
}