const ROTOR = [ 10, 7, 5, 3 ];
const GROUND_STEP = 4;
const SPEED = 0.01;
const SPEED_Y = 0.05;
const GRAVITY_SPEED = 0.02;
const MAX_HEIGHT = 85;
const GROUND_HEIGHT_STEP = 2;

const HIT_NOTHING = 0;
const HIT_GROUND = 1;
const HIT_PAD = 2;

player := {
    "x": 40,
    "y": 100,
    "dir": 0,
    "rotor": 0,
    "switch": 0,
    "dirchange": 0,
    "move": 0,
    "moveY": 0,
    "explode": 0,
    "lives": 5,
    "gravity_enabled": true
};
ground := [];
groundIndex := 50;    
scrollStep := 0;
turnDir := 0;

def drawRotor(color) {
    fillRect(player["x"]-ROTOR[player["rotor"]], player["y"]-7, player["x"]+ROTOR[player["rotor"]], player["y"]-10, color);
}

def handleInput() {
    if(player["explode"] > 0) {
        # return must always return a value...
        return false;
    }
    if(isKeyDown(KeyLeft)) {
        if(turnDir != -1) {
            player["dirchange"] := 0;
        }
        turnDir := -1;
    } else {
        if(isKeyDown(KeyRight)) {
            if(turnDir != 1) {
                player["dirchange"] := 0;
            }
            turnDir := 1;
        } else {
            turnDir := 0;
        }
    }
    if(isKeyDown(KeyUp) && player["y"] > 10) {
        player["y"] := player["y"] - SPEED_Y;
    }
    if(isKeyDown(KeyDown)) {
        player["y"] := player["y"] + SPEED_Y;
    }

    if(getTicks() > player["dirchange"]) {
        if(turnDir = -1 && player["dir"] > -1) {
            player["dir"] := player["dir"] - 1;
        }
        if(turnDir = 1 &&  player["dir"] < 1) {
            player["dir"] := player["dir"] + 1;
        }
        player["dirchange"] := getTicks() + 0.15;
    }
}

def movePlayer() {
    if(player["explode"] = 0 && getTicks() > player["moveY"] && player["gravity_enabled"]) {
        player["moveY"] := getTicks() + GRAVITY_SPEED;

        # gravity
        player["y"] := player["y"] + 1;
    }

    if(player["dir"] != 0 && getTicks() > player["move"]) {
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
            } else {
                if(player["x"] < 160 && player["x"] > 0) {
                    player["x"] := player["x"] + player["dir"];
                }
            } 
        }
    }    
}

def testCollision() {
    if(player["dir"] = 0) {
        sx := player["x"] - 5;
        ex := player["x"] + 5;
    }
    if(player["dir"] = 1) {
        sx := player["x"] - 12;
        ex := player["x"] + 5;
    }
    if(player["dir"] = -1) {
        sx := player["x"] - 5;
        ex := player["x"] + 12;
    }
    sy := player["y"] - 10;
    ey := player["y"] + 5;
    while(sx < ex) {
        while(sy < ey) {
            gi := groundIndex + sx / GROUND_STEP;
            if(gi >= 0 && gi < len(ground)) {
                if(ground[gi]["pad"] > -1 && sy > 200 - ground[gi]["pad"]) {
                    return HIT_PAD;
                } else {
                    groundHeight := ground[gi]["height"];
                    if(sy > 200 - groundHeight) {
                        return HIT_GROUND;
                    }
                }
            }
            sy := sy + GROUND_HEIGHT_STEP;
        }
        sx := sx + GROUND_STEP;
    }
    return HIT_NOTHING;
}

def drawPlayerExplode() {
    i := 0;
    while(i < 10) {
        if(random() > 0.5) {
            color := COLOR_YELLOW;
        } else {
            color := COLOR_WHITE;
        }
        fillCircle(player["x"] - 5 + (random() * 10), player["y"] - 5 + (random() * 10), random() * 10 + 3, color);
        i := i + 1;
    }
}

def drawPlayerHealthy() {
    fillCircle(player["x"], player["y"], 5, COLOR_RED);
    if(player["dir"] = 0) {
        fillRect(player["x"]-3, player["y"]-2, player["x"]+3, player["y"], COLOR_WHITE);
    }
    if(player["dir"] = 1) {
        fillRect(player["x"], player["y"]-2, player["x"]+3, player["y"], COLOR_WHITE);
        fillRect(player["x"] - 12, player["y"] - 5, player["x"], player["y"], COLOR_RED);
        fillRect(player["x"] - 12, player["y"] - 7, player["x"]-10, player["y"]-5, COLOR_RED);
    }
    if(player["dir"] = -1) {
        fillRect(player["x"]-3, player["y"]-2, player["x"], player["y"], COLOR_WHITE);
        fillRect(player["x"], player["y"] - 5, player["x"]+12, player["y"], COLOR_RED);
        fillRect(player["x"]+10, player["y"] - 7, player["x"]+12, player["y"]-5, COLOR_RED);
    }
    fillRect(player["x"]-1, player["y"]-10, player["x"]+1, player["y"], COLOR_RED);


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

def drawPlayer() {
    if(player["explode"] > getTicks()) {
        drawPlayerExplode();
    } else {        
        if(player["explode"] > 0) {
            player["lives"] := player["lives"] - 1;
            player["y"] := 100;
            player["explode"] := 0;
        } else {
            # collision check
            collision := testCollision();
            if(collision = HIT_GROUND) {
                player["explode"] := getTicks() + 1.5;
                player["dir"] := 0;
            }
            if(collision = HIT_PAD && player["y"] >= 200 - MAX_HEIGHT - 5) {
                player["gravity_enabled"] := false;
            } else {
                player["gravity_enabled"] := true;
            }
        }
        drawPlayerHealthy();
    }
}

def initGround() {
    length := 1000;
    h := random() * MAX_HEIGHT;
    while(len(ground) < length) {
        g := { "height": h };
        putPad := len(ground) % 300;
        if(putPad >= 25 && putPad < 35) {
            padHeight := h;
            if(len(ground) > 0 && ground[len(ground) - 1]["pad"] > -1) {
                padHeight := ground[len(ground) - 1]["pad"];
            }
            g["pad"] := padHeight;            
            g["height"] := 0;
        } else {
            g["pad"] := -1;
        }
        ground[len(ground)] := g;
        if(random() > 0.5) {
            if(h < MAX_HEIGHT) {
                h := h + GROUND_HEIGHT_STEP;
            }
        } else {
            if(h > 4) {
                h := h - GROUND_HEIGHT_STEP;
            }
        }
    }
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

def drawGround() {
    x := 0; 
    if(scrollStep > 0) {
        x := -1 * GROUND_STEP;
    }
    sx := x;
    while(x < 160) {
        gi := groundIndex + x/GROUND_STEP;

        if(ground[gi]["pad"] > -1) {
            fillRect(x + scrollStep, 200 - ground[gi]["pad"], x + scrollStep + GROUND_STEP, 200, COLOR_DARK_GRAY);
        } else {        
            fillRect(x + scrollStep, 200 - ground[gi]["height"], x + scrollStep + GROUND_STEP, 200, COLOR_GREEN);
        }

        x := x + GROUND_STEP;
    }

}

def main() {
    setVideoMode(2);
    setBackground(COLOR_DARK_BLUE);
    initGround();
    
    while(1=1) {
        clearVideo();
        if(player["lives"] > 0) {        
            handleInput();
            movePlayer();
            drawGround();
            drawPlayer();
            drawText(0, 0, COLOR_LIGHT_BLUE, COLOR_DARK_BLUE, "LIFE:" + player["lives"]);
        } else {
            fillRect(40, 60, 120, 140, COLOR_YELLOW);
            drawText(45, 94, COLOR_BLACK, COLOR_YELLOW, "Game Over");
        }
        updateVideo();
    }
}
