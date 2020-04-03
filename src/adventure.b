############################################
# a small text adventure game demo

# room indexes (enum)
const SMALL_ROOM = 0;
const GRASSY_FIELD = 1;
const OCEAN_BEACH = 2;

# room definitions
const ROOMS = [ 
    "You are in a small room with no furniture. There is a door to the north.",
    "You are standing on a grassy field with a small house to the south and the ocean to the east.",
    "You are on a sandy beach. The ocean is calm. There is a grassy field to the west."
];

# map: 1 row per room: room index of directions: n,s,e,w
const MAP = [
    [ GRASSY_FIELD, -1, -1, -1 ],
    [ -1, SMALL_ROOM, OCEAN_BEACH, -1 ],
    [ -1, -1, -1, GRASSY_FIELD ]
];

# direction commands
const DIRS = [ "n", "s", "e", "w" ];

# object enums
const DOOR_KEY = 0;
const SEASHELL = 1;

const OBJECTS = [
    {
        "description": "a small brass key",
        "short": "key"
    },
    {
        "description": "a beautiful sea-shell",
        "short": "shell"
    }
];


# global variables
current_room := 0;
inventory := [];
door_locked := 1;

# one list of objects per room
object_locations := [
    [DOOR_KEY],
    [],
    [SEASHELL]
];

# process a direction command and return the 1 if handled, 0 otherwise
def direction_command(action) {

    # todo: we don't yet handle multidimensional arrays: MAP[current_room][0], etc.
    current_room_map := MAP[current_room];
    handled := 0;
    i := 0;
    while(i < len(DIRS)) {
        if(DIRS[i] = action) {
            handled := 1;
            proposed_room := current_room_map[i];
        }
        i := i + 1;
    }

    if(handled = 1) {
        if (proposed_room = -1) {
            print("You cannot go that way.");
        } else {
            if(room_transition(proposed_room) = 1) {
                current_room := proposed_room;
            }
        }
    }

    return handled;
}

# if the transition is allowed, return 1, else 0
def room_transition(proposed_room) {
    # todo: need AND/OR handling
    if(proposed_room = GRASSY_FIELD) {
        if(current_room = SMALL_ROOM) {
            if(door_locked = 1){
                if(has_object("key") = 1) {
                    door_locked := 0;
                    print("You use your key to unlock the door.");
                } else {
                    print("A locked door prevents you from going outside.");
                    return 0;
                }
            }
        }
    }
    return 1;
}

def get_object(name) {
    object_index := -1;
    objects_in_room := object_locations[current_room];
    i := 0;
    while(i < len(objects_in_room)) {
        idx := objects_in_room[i];
        obj := OBJECTS[idx];
        if(obj["short"] = name) {
            object_index := idx;
            # remove from the room
            del objects_in_room[i];
            inventory[len(inventory)] := object_index;
            obj := OBJECTS[object_index];
            print("You pick up " + obj["description"] + ".");
            return 0;
        } else {
            i := i + 1;
        }
    }
    print("You don't see that anywhere.");
}

def drop_object(name) {
    object_index := -1;
    objects_in_room := object_locations[current_room];
    i := 0;
    while(i < len(OBJECTS)) {
        obj := OBJECTS[i];
        if(obj["short"] = name) {
            object_index := i;
            objects_in_room[len(objects_in_room)] := object_index;

            i := 0;
            while(i < len(inventory)) {
                if(inventory[i] = object_index) {
                    del inventory[i];
                }
                i := i + 1;
            }
            print("You drop " + obj["description"] + ".");
            return 0;
        } else {
            i := i + 1;
        }
    }
    print("You are not carrying that.");
}

def inventory() {
    if(len(inventory) > 0) {
        print("You are carrying:");
        i := 0;
        while(i < len(inventory)) {
            o := OBJECTS[i];
            print("\t-" + o["description"]);
            i := i + 1;
        }
    } else {
        print("You aren't carrying anything.");
    }
}

def has_object(name) {
    i := 0;
    while(i < len(inventory)) {
        o := OBJECTS[i];
        if(o["short"] = name) {
            return 1;
        }
        i := i + 1;
    }
    return 0;
}

def main() {
    print("Welcome to the bscript adventure demo.");
    print("The object of the game is to decorate your living room.");
    
    running := 1;    
    while(running > 0) {
        
        # print the current status and ask for input
        print("");
        print(ROOMS[current_room]);
        i := 0;
        objects_in_room := object_locations[current_room];
        while(i < len(objects_in_room)) {
            obj := OBJECTS[objects_in_room[i]];
            print("There is " + obj["description"] + " here.");
            i := i + 1;
        }
        print("");
        action := input("> ");

        # command handling
        handled := direction_command(action);
        if (handled = 0) {
            # todo: need "else if"
            if (action = "exit") {
                running := 0;
                handled := 1;
            }
            if(substr(action, 0, 4) = "get ") {
                get_object(substr(action, 4));
                handled := 1;
            }
            if(substr(action, 0, 5) = "drop ") {
                name := substr(action, 5);
                drop_object(name);
                if(name = "shell") {
                    if(current_room = SMALL_ROOM) {
                        print("");
                        print("Congratulations! You have found suitable decoration for your living room.");
                        return 0;
                    }
                }
                handled := 1;
            }
            if(substr(action, 0, 3) = "inv") {
                inventory();
                handled := 1;
            }
            if(handled = 0) {
                print("I don't understand this command.");
            }
        }
    }
    print("Goodbye.");
}
