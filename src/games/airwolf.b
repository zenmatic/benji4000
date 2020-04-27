const ROTOR = [ 10, 7, 5, 3 ];
const GROUND_STEP = 4;
const SPEED = 0.01;
#const SPEED = 0.5;
const SPEED_Y = 0.05;
const MAX_HEIGHT = 85;

player := {
    "x": 40,
    "y": 100,
    "dir": 0,
    "rotor": 0,
    "switch": 0,
    "dirchange": 0,
    "move": 0
};
ground := [];
groundIndex := 50;    
scrollStep := 0;

def drawRotor(color) {
    fillRect(player["x"]-ROTOR[player["rotor"]], player["y"]-7, player["x"]+ROTOR[player["rotor"]], player["y"]-10, color);
}

def drawPlayerBody(color) {
    fillCircle(player["x"], player["y"], 5, color);
    if(player["dir"] = 0) {
        if(color = COLOR_RED) {
            fillRect(player["x"]-3, player["y"]-2, player["x"]+3, player["y"], COLOR_WHITE);
        }
    }
    if(player["dir"] = 1) {
        if(color = COLOR_RED) {
            fillRect(player["x"], player["y"]-2, player["x"]+3, player["y"], COLOR_WHITE);
        }
        fillRect(player["x"] - 12, player["y"] - 5, player["x"], player["y"], color);
    }
    if(player["dir"] = -1) {
        if(color = COLOR_RED) {
            fillRect(player["x"]-3, player["y"]-2, player["x"], player["y"], COLOR_WHITE);
        }
        fillRect(player["x"], player["y"] - 5, player["x"]+12, player["y"], color);
    }
    fillRect(player["x"]-1, player["y"]-10, player["x"]+1, player["y"], color);
}

def clearPlayer() {
    drawPlayerBody(COLOR_DARK_BLUE);
    drawRotor(COLOR_DARK_BLUE);
}

def movePlayer() {
    if(player["dir"] != 0 &&  getTicks() > player["move"]) {
        player["move"] := getTicks() + SPEED;

        handled := false;
        if(player["dir"] = 1 && player["x"] < 80) {
            player["x"] := player["x"] + 1;
            handled := true;
        }
        if(player["dir"] = -1 && player["x"] > 80) {
            player["x"] := player["x"] - 1;
            handled := true;
        }

        if(handled = false) {
            if(canScroll()) {
                scroll(-2 * player["dir"], 0);
                scrollStep := scrollStep - player["dir"]*2;
                if(scrollStep >= GROUND_STEP) {
                    scrollStep := 0;
                    groundIndex := groundIndex - 1;
                }
                if(scrollStep <= -1) {
                    scrollStep := GROUND_STEP - 2;
                    groundIndex := groundIndex + 1;
                }
                
                if(player["dir"] = 1) {
                    xp := 160 - (GROUND_STEP - scrollStep);
                    fillRect(xp, 200 - ground[groundIndex + 160/GROUND_STEP], xp + GROUND_STEP, 200, COLOR_GREEN);
                } else {
                    xp := (scrollStep + 2) - GROUND_STEP;
                    fillRect(xp, 200 - ground[groundIndex - 1], xp + GROUND_STEP, 200, COLOR_GREEN);
                }
            } else {
                if(player["dir"] = 1) {
                    if(player["x"] < 160) {
                        player["x"] := player["x"] + 1;
                    }
                } else {
                    if(player["x"] > 0) {
                        player["x"] := player["x"] - 1;
                    }
                }
            } 
        }
    }
}

def drawPlayer() {
    drawPlayerBody(COLOR_RED);

    # animate the rotor
    if(getTicks() > player["switch"]) {
        drawRotor(COLOR_DARK_BLUE);
        player["rotor"] := player["rotor"] + 1;
        if (player["rotor"] >= len(ROTOR)) {
            player["rotor"] := 0;
        }
        player["switch"] := getTicks() + 0.025;
        
    }
    drawRotor(COLOR_RED);
}

def initGround() {
    h := random() * MAX_HEIGHT;
    while(len(ground) < 1000) {
        ground[len(ground)] := h;
        if(random() > 0.5) {
            if(h < MAX_HEIGHT) {
                h := h + 2;
            }
        } else {
            if(h > 4) {
                h := h - 2;
            }
        }
    }
}

def changeDir(dir) {
    if(dir = -1) {
        if(player["dir"] > -1) {
            if(getTicks() > player["dirchange"]) {
                player["dir"] := player["dir"] - 1;
                
            }
        }
    } else {
        if(player["dir"] < 1) {
            if(getTicks() > player["dirchange"]) {
                player["dir"] := player["dir"] + 1;
            }
        }
    }
    player["dirchange"] := getTicks() + 0.25;
}

def canScroll() {
    if(player["dir"] = 1 && groundIndex >= len(ground) - (160/GROUND_STEP) - 1) {
        return false;
    }
    if(player["dir"] = -1 && groundIndex <= 1) {
        return false;
    }
    return true;
}

def main() {
    setVideoMode(2);
    setBackground(COLOR_DARK_BLUE);
    clearVideo();
    initGround();

    # draw the ground
    x := 0; 
    while(x < 160) {
        fillRect(x, 200 - ground[groundIndex + x/GROUND_STEP], x + GROUND_STEP, 200, COLOR_GREEN);
        x := x + GROUND_STEP;
    }

    while(1=1) {
        clearPlayer();

        if(isKeyDown(KeyLeft)) {
            changeDir(-1);
        }
        if(isKeyDown(KeyRight)) {
            changeDir(1);
        }        
        if(isKeyDown(KeyUp)) {
            player["y"] := player["y"] - SPEED_Y;
        }
        if(isKeyDown(KeyDown)) {
            player["y"] := player["y"] + SPEED_Y;
        }

        movePlayer();
        drawPlayer();
        updateVideo();
    }
}