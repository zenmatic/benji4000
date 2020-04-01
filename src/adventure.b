############################################
# a small text adventure game demo

# room indexes (enum)
const SMALL_ROOM = 0;
const GRASSY_FIELD = 1;
const OCEAN_BEACH = 2;

# room definitions
const ROOMS = [ 
    "You are in a small room. There is a cardboard box on the floor and a door to the north.",
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

# global variables
current_room := 0;

# process a direction command and return the 1 if handled, 0 otherwise
def direction_command(action)

    # todo: we don't yet handle multidimensional arrays: MAP[current_room][0], etc.
    current_room_map := MAP[current_room];

    handled := 0;
    i := 0;
    while(i < len(DIRS))
        if(DIRS[i] = action)
            handled := 1;
            proposed_room := current_room_map[i];
        end
        i := i + 1;
    end

    if(handled = 1)
        if (proposed_room = -1)
            print("You cannot go that way.");
        end
        if (proposed_room > -1)
            current_room := proposed_room;
        end
    end

    return handled;
end

def main()
    print("Welcome to the bscript adventure demo");
    print("");
    
    running := 1;    
    while(running > 0)
        
        # print the current status and ask for input
        print(ROOMS[current_room]);        
        print("");
        action := input("> ");

        # command handling
        handled := direction_command(action);
        if (handled = 0)
            if (action = "exit")
                running := 0;
            end
        end
    end
    print("Goodbye.");
end
