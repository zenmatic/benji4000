const ROTOR = [ 10, 7, 5, 3 ];
const GROUND_STEP = 4;
const SPEED = 0.01;
const SPEED_FUEL = 0.05; const SPEED_FUEL_DOWN = 0.2;
const SPEED_Y = 0.05;
const WAIVE_SPEED = 0.15;
const GRAVITY_SPEED = 0.02;
const MAX_HEIGHT = 35;
const GROUND_HEIGHT_STEP = 2;
const WAIVE = [ 3, 2, 1, 2 ];
const HIT_NOTHING = 0;
const HIT_GROUND = 1;
const HIT_PAD = 2;
const PLAYER_COLOR = COLOR_TAN;

const DROP_SPEED = 0.5;

# todo: can't place comments inside map literal
player := {
    "x": 80,
    "y": 190,
    "switch": 0,
    "move": 0,
    "moveY": 0,
    "explode": 0,
    "lives": 5,
    "gravity_enabled": true
};
title := true;
info := false;
gameOn := false;
death := false;
ground := [];
groundIndex := 0;    
scrollStep := 0;
turnDir := 0;
soldiers := [ 
    330 * GROUND_STEP, 
    331 * GROUND_STEP, 
    334 * GROUND_STEP, 
    630 * GROUND_STEP, 
    627 * GROUND_STEP, 
    930 * GROUND_STEP, 
    932 * GROUND_STEP, 
    928 * GROUND_STEP 
];
soldierMoveTimer := 0;
waiveTimer := 0;
waiveIndex := 0;

dropTimer := 0;
drops := [
    [ 0, 0, 0, 0, 0, 0],
    [ 0, 0, 0, 0, 0, 0],
    [ 0, 0, 0, 0, 0, 0],
    [ 0, 0, 0, 0, 0, 0]
];

def drawAcidDrop(x, y) {
    i := 0;
    while(i < 10) {
        if(random() > 0.5) {
            #color := COLOR_LIGHT_BLUE;
            color := COLOR_GREEN;
        } else {
            #color := COLOR_WHITE;
            #color := COLOR_DARK_BLUE;
            color := COLOR_YELLOW;
        }
        fillCircle(x - 4 + (random() * 10), y - 4 + (random() * 10), random() * 1 + 1, color);
        i := i + 1;
    }
}

def updateAcidRain() {
    maxrows := len(drops) - 1;
    i := maxrows;
    # move all drops down a row, starting from the bottom
    while(i >= 0) {
        maxcols := len(drops[i]);
        j := 0;
        while(j < maxcols) {
            if(drops[i][j] = 1) {
                nextrow := i + 1;
                drops[i][j] := 0;
                if(nextrow <= maxrows) {
                    drops[nextrow][j] := 1;
                }
            }
            j := j + 1;
        }
        i := i - 1;
    }

    j := 0;
    while(j < len(drops)) {
        if (random() > 0.8) {
            drops[0][j] := 1;
        }
        j := j + 1;
    }
}

def drawAcidRain() {
    if(getTicks() > dropTimer) {
         updateAcidRain();
         dropTimer := getTicks() + DROP_SPEED;
    }
    step := 20;
    i := 0;
    while(i < len(drops)) {
        j := 0;
        while(j < len(drops[i])) {
            if(drops[i][j] = 1) {
                x := (j + 1) * 20 + 10;
                y := (i + 1) * 20 + 60;
                drawAcidDrop(x, y);
            }
            j := j + 1;
        }
        i := i + 1;
    }
}

def drawSoldier(index, x, y) {
    if(getTicks() > waiveTimer) {
        waiveTimer := getTicks() + WAIVE_SPEED;
        waiveIndex := waiveIndex + 1;
        if(waiveIndex >= len(WAIVE)) {
            waiveIndex := 0;
        }
    }
    wi := (waiveIndex + index) % len(WAIVE);
    drawLine(x + WAIVE[wi], y - 10, x + 4, y - 4, COLOR_LIGHT_BLUE);
    drawLine(x + 8 - WAIVE[wi], y - 10, x + 4, y - 4, COLOR_LIGHT_BLUE);
    fillRect(x + 3, y - 8, x + 5, y - 5, COLOR_LIGHT_BLUE);
    fillRect(x + 2, y - 4, x + 3, y, COLOR_LIGHT_BLUE);
    fillRect(x + 5, y - 4, x + 6, y, COLOR_LIGHT_BLUE);
}

def drawPlayerHealthy() {
    drawSoldier(0, player["x"], player["y"]);
    #fillCircle(player["x"], player["y"], 5, PLAYER_COLOR);
    #fillRect(player["x"]-3, player["y"]-2, player["x"]+3, player["y"], COLOR_WHITE);
    #fillRect(player["x"]-1, player["y"]-10, player["x"]+1, player["y"], PLAYER_COLOR);
    #drawText(20,20, COLOR_LIGHT_BLUE, COLOR_DARK_BLUE, player["x"]+":"+player["y"]);
}

def drawGround() {
    x := 0;
    y := 190;
    #maxX := 160;
    #maxY := 200;
    fillRect(x, y, x+160, y+10, COLOR_GREEN);
}

def drawUI() {
    drawText(0, 1, COLOR_LIGHT_BLUE, COLOR_DARK_BLUE, "LIFE:" + player["lives"]);
}

def drawTitle() {
    drawRect(5, 5, 155, 195, COLOR_DARK_BLUE);
    drawText(40, 20, COLOR_BROWN, COLOR_BLACK, "Acid Rain!");
    drawText(14, 45, COLOR_MID_GRAY, COLOR_BLACK, "for the Benji4000");
    drawText(25, 160, COLOR_MID_GRAY, COLOR_BLACK, "SPACE to start");
    drawText(14, 175, COLOR_DARK_GRAY, COLOR_BLACK, "2020 (c) by Matt");

    drawAcidRain();
    if(isKeyDown(KeySpace)) {
        title := false;
        info := true;
        while(isKeyDown(KeySpace)) {
        }
    }
}

def drawInfo() {
    drawRect(5, 5, 155, 195, COLOR_DARK_BLUE);
    drawText(14, 25, COLOR_MID_GRAY, COLOR_BLACK, "You find yourself");
    drawText(14, 35, COLOR_MID_GRAY, COLOR_BLACK, "deep behind enemy");
    drawText(14, 45, COLOR_MID_GRAY, COLOR_BLACK, "lines in Soviet");
    drawText(14, 55, COLOR_MID_GRAY, COLOR_BLACK, "Siberia. You look");
    drawText(14, 65, COLOR_MID_GRAY, COLOR_BLACK, "above you and");
    drawText(14, 75, COLOR_MID_GRAY, COLOR_BLACK, "see... Acid Rain!");

    drawText(14, 105, COLOR_RED, COLOR_BLACK, "Stay dry or die!");
    drawText(14, 135, COLOR_MID_GRAY, COLOR_BLACK, "SPACE to begin");
    if(isKeyDown(KeySpace)) {
        info := false;
        gameOn := true;
        while(isKeyDown(KeySpace)) {
        }
    }
}

def drawDeath() {
    setBackground(COLOR_BLACK);
    drawText(45, 87, COLOR_RED, COLOR_BLACK, "You Died!");
    if(isKeyDown(KeySpace)) {
        return 0;
    }
    return 1;
}

def handleGame() {
    setBackground(COLOR_DARK_BLUE);
    #movePlayer();
    drawGround();
    drawPlayerHealthy();
    drawUI();
    if(isKeyDown(KeySpace)) {
        gameOn := false;
        death := true;
        while(isKeyDown(KeySpace)) {
        }
    }
}

def main() {
    setVideoMode(2);
    setBackground(COLOR_BLACK);
    
    on := 1;
    while(on=1) {
        clearVideo();
        if(title) {
            drawTitle();
        }
        if(info) {
            drawInfo();
        }
        if(gameOn) {
            handleGame();
        }
        if (death) {
            on := drawDeath();
        }
        updateVideo();
    }
}
