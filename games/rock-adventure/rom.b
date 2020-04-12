############################################
# a pick-a-path adventure: it's a rock concert; what could go wrong?

# TODO:
##
# add var like LOGO into text
##

# room indexes (enum)
const NULL = 0;
const HOME = 1;
const PICK_JACKET = 2;
const CLUB = 3;
const WHERE_TO_GO = 4;
const MOSH_PIT = 5;
const LOUNGE = 6;
const BAR = 7;
const DANCE_STYLE = 8;
const LOUNGE_DECISION = 9;
const COCKTAIL = 10;
const HARDCORE = 11;
const LIGHT = 12;
const MEDIUM = 13;
const SIT = 14;
const T_SHIRT = 15;
const ORDER_DRINK = 16;
const BEER = 17;
const COSMO = 18;
const MANHATTAN = 19;
const ANGRY_GUY = 20;
const DRUG_DEALER = 21;
const GIRL = 22;
const FIGHT_OR_FLIGHT = 23;
const HELP_GIRL = 24;
const DRUG_DEAL = 25;
const HIT_FACE = 26;
const BACK_AWAY = 27;
const HIT_GUT = 28;
const BUY = 29;
const DONT_BUY = 30;
const EXIT = 31;
const FRONT = 32;
const BATHROOM = 33;
const RED_DOOR = 34;
const WIN = 35;
const DIE = 36;

const LOGOS = ["Deicide", "Dimmu Borgir", "Gwar"];
# room definitions
const SCENARIOS = [
        {
                "index": NULL,
                "next": NULL,
                "desc": "",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": HOME,
                "next": PICK_JACKET,
                "desc": "You've been stoked for weeks about this concert.  Your favorite band, Plague Puke, is playing at the notorious underground club, The Spike.  You're dressing for the show, standing in front of your closet in your tighty whiteys.  You have a bunch of jackets in there.  Your mama pain-stakingly stitched a band's logo on each one.  Which one do you choose?",
                "choices": LOGOS,
                "dests": [],
                "answer": ""
        },
        {
                "index": PICK_JACKET,
                "next": CLUB,
                "desc": "You don the jacket with the LOGO logo and catch the bus to downtown.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": CLUB,
                "next": WHERE_TO_GO,
                "desc": "After some searching, you hear some muffled music in the distance and find The Spike in an alley off a main street.  You pay the cover to the bouncer with all the piercings and walk into the place.  The band is already playing loud and everyone is rocking.  There's a mosh pit to the left, a bar to the right, and a lounge straight ahead of you.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": WHERE_TO_GO,
                "next": NULL,
                "desc": "Where do you go?",
                "choices": ["the mosh pit", "the grungy lounge", "the bar"],
                "dests": [MOSH_PIT, LOUNGE, BAR],
                "answer": ""
        },
        {
                "index": MOSH_PIT,
                "next": DANCE_STYLE,
                "desc": "You head into the mosh pit.  It's like falling into a pot of boiling water.  It's a free-for-all where you're not sure if peopel are dancing or fighting.  You feel like you need to do something other than just stand there.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": LOUNGE,
                "next": LOUNGE_DECISION,
                "desc": "You head into the lounge, if that's what you could call it.  You are greeted by a random assortment of chairs and couches most of which look like they have been left out on the street for a while.  There's also a dumpy stand at the edge, selling Plague t-shirts and hoodies.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": BAR,
                "next": COCKTAIL,
                "desc": "Yuck, you've seen gutters cleaner than this bar.  Before you can get away a middle-aged, balding bartender approaches and asks what you'll have.  The size of the bags under his eyes is disturbing, but you try not to focus on it.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": DANCE_STYLE,
                "next": NULL,
                "desc": "How do you want to dance?",
                "choices": ["boy band moves", "the latest hip hop", "hardcore"],
                "dests": [LIGHT, MEDIUM, HARDCORE],
                "answer": ""
        },
        {
                "index": LOUNGE_DECISION,
                "next": NULL,
                "desc": "What do you want to do?",
                "choices": ["sit", "go to T-shirt stand", "order a drink"],
                "dests": [SIT, T_SHIRT, ORDER_DRINK],
                "answer": ""
        },
        {
                "index": COCKTAIL,
                "next": NULL,
                "desc": "What do you order?",
                "choices": ["a cosmo", "a beer", "a manhattan straight up"],
                "dests": [COSMO, BEER, MANHATTAN],
                "answer": ""
        },
        {
                "index": HARDCORE,
                "next": ANGRY_GUY,
                "desc": "You let out a primodial scream and punch in every direction.  Your fist accidently connects with someone beside you.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": LIGHT,
                "next": DRUG_DEALER,
                "desc": "You do your best nsync impression and even attempt some moon walking.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": MEDIUM,
                "next": LOUNGE,
                "desc": "First an MC Hammer impression, then the Robot, finally a little twerking.  You are quickly pushed from the pit by everyone around you.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": SIT,
                "next": ANGRY_GUY,
                "desc": "You sit down in one of he musty seats.  You kick up your shoes when you hear someone snarl behind you, <<Hey, man, that's my seat.  Get the fuck out of it!>>",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": T_SHIRT,
                "next": DRUG_DEALER,
                "desc": "You go to the t-shirt stand and browse the wares.  The guy selling he stuff gives you the evil eye as you handle the merchandise.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": ORDER_DRINK,
                "next": BAR,
                "desc": "You sit down and order a drink from a woman nearby.  <<Fuck you>> she says <<do I look like a fucking waitress to you?>>  Embarrassed you head to the bar.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": BEER,
                "next": ANGRY_GUY,
                "desc": "<<Gimme a bud>> you say.  The bartender serves it and walks away.  You raise the beer to your lips when someone bumps you from behind.  You spill the beer all over.  A voice snarls behind you <<Hey, watch it, poser>>",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": COSMO,
                "next": DRUG_DEALER,
                "desc": "The bartender smirks when you order a cosmo.  He fills a glass with ice and adds tequila, fruit punch, and a packet of sweet and low.  You're pretty sure that's not how a cosmo is made, but realize you don't actually know.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": MANHATTAN,
                "next": GIRL,
                "desc": "You order a manhatten with Wild Turkey, two cherries, and extra bitters.  The bartenders appears impressed and with a wry smile he delicately prepares your cocktail for you.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": ANGRY_GUY,
                "next": FIGHT_OR_FLIGHT,
                "desc": "A huge barrel-chested dude with an oily black mullet and a scruffy chin beard stares at you angrily.  <<Yeah, that's right.  Fuck you, you little bitch>> he says, <<and fuck LOGO.  LOGO sucks!>>",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": DRUG_DEALER,
                "next": DRUG_DEAL,
                "desc": "You suddenly notice someone beside you.  He draws you aside.  <<Hey man, you look like you could use some, eh, inspiration>>.  The guy is just plain bizaare looking.  He holds out a few pills and says, <<You wanna buy some of these, man?>> and holds out his hand with a few pills.  <<And I've got some even better stuff out back if you interested...>>.  He gestures to the red door at the back of the club.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": GIRL,
                "next": HELP_GIRL,
                "desc": "A stunning punk rocker stands before you.  She's got spikey blonde hair and blue lipstick and (check this out) she's wearing a jacket with a LOGO logo on it.  <<Nice jacket>> she purrs <<we've got something in common>>.  You try to say something cool back to her, but you end up just croaking <<Hey, yeah, cool>>  The girl slides closer.  <<Hey, I know we just met and all, but do you think you could help me get back stage?>>",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": FIGHT_OR_FLIGHT,
                "next": NULL,
                "desc": "What do you do?",
                "choices": ["sock him in the nose", "run away", "sucker punch to the gut"],
                "dests": [HIT_FACE, BACK_AWAY, HIT_GUT],
                "answer": ""
        },
        {
                "index": HELP_GIRL,
                "next": NULL,
                "desc": "<<Yeah>> you lie <<I've been thinking about sneaking back stage too.>> You scope the club and figure there must be three ways to get back there.",
                "choices": ["through the mosh pit and behind the stage", "go out the entrance and around the back", "go through the red door"],
                "dests": [MOSH_PIT, FRONT, RED_DOOR],
                "answer": ""
        },
        {
                "index": DRUG_DEAL,
                "next": NULL,
                "desc": "What do you do?",
                "choices": ["buy them", "don't buy them", "head to the door"],
                "dests": [BUY, DONT_BUY, RED_DOOR],
                "answer": ""
        },
        {
                "index": HIT_FACE,
                "next": DIE,
                "desc": "You hit him in the face as hard as you can.  Everyone gasps as blood spurts from his face.  But he doesn't go down.  Now he's pissed.  He pulls a switch blade from his pocket and buries it hilt deep into your chest.  The last thing you see is everyone running as the darkness closes in.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": BACK_AWAY,
                "next": EXIT,
                "desc": "You turn and run, knocking down people in your way.  Your thoughts are racing a mile a minute. <<How do I get out of this place?>>",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": HIT_GUT,
                "next": DRUG_DEALER,
                "desc": "You channel the spirit of Tyson and sucker punch him in the gut.  The guy's eyes literally bulge.  He falls back into the crowd gasping for air.  The crowd parts in awe as you walk away.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": BUY,
                "next": DIE,
                "desc": "You buy the pills and pop them into your mouth.  You see the guy's eyes widen.  <<You ain't supposed to take them all at once, dude!>>  The pain begins in your stomache and makes it way up toward you head.  You crash to the floor, foaming at the mouth.  A woman nearby lets out a scream.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": DONT_BUY,
                "next": GIRL,
                "desc": "<<Naw>> you say <<I don't got no cash, bra>>  The dealer shrugs his shoulders and disappears into the crowd.  <<Hey, you want to share those?>> a silky voice says beside you.  <<Oh wait, you didn't buy anything>>",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": EXIT,
                "next": NULL,
                "desc": "Your bloodshot eyes scan the club.  There are three ways to high tail out of here.  Where do you run?",
                "choices": ["out the entrance", "to the bathroom", "through the red door"],
                "dests": [FRONT, BATHROOM, RED_DOOR],
                "answer": ""
        },
        {
                "index": FRONT,
                "next": WIN,
                "desc": "You walk out the entrance of the club.  You realize this may be a mistake so you turn around to go back in.  But the bouncer blocks your way.  He wants you to pay the cover again.  <<I already paid the cover>> you say, <<you just saw me walk out.  I've been waiting for months to see this band>>  But he's unmoved and tells you go get lost.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": BATHROOM,
                "next": DIE,
                "desc": "You panic and run into the bathroom.  There's a slimy sink, a broken urinal, and an overflowing toilet.  If you had time, you could admire the graffitti that covers the walls.  You turn to see the mullet guy has followed you.  His right hook floors you.  You last thing you see is stars floating in a haze.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": RED_DOOR,
                "next": DIE,
                "desc": "You walk through the red door.  It looks like it would lead into a back room, but instead it opens straight into an alley way.  A smelly dumpster is on one side and a group of thugs on the other.  They rob you and beat you senseless.",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": WIN,
                "next": NULL,
                "desc": "Little do you know it, but you escaped death!",
                "choices": [],
                "dests": [],
                "answer": ""
        },
        {
                "index": DIE,
                "next": NULL,
                "desc": "You are dead.",
                "choices": [],
                "dests": [],
                "answer": ""
        }
];

# global vars
logo := "";

def get_choice(choices) {
        ret := 4;
        while(ret = 4) {
                print("1: " + choices[0]);
                print("2: " + choices[1]);
                print("3: " + choices[2]);
                ans := input("> ");
                if(ans = "1") {
                        return 0;
                }
                if(ans = "2") {
                        return 1;
                }
                if(ans = "3") {
                        return 2;
                }
                print("You must answer 1, 2, or 3");
        }
        return ret;
}

def main() {
    print("Welcome to the Night at the Spike, a pick-a-path adventure");

    newscene := HOME;
    
    while(newscene != NULL) {

        scene := newscene;
        #print("scene=" + newscene);
        
        print("");
        print(SCENARIOS[scene]["desc"]);
        print("");

        if(len(SCENARIOS[scene]["choices"]) > 0) {
                ret := get_choice(SCENARIOS[scene]["choices"]);
                SCENARIOS[scene]["answer"] := ret;

                if (len(SCENARIOS[scene]["dests"]) > 0) {
                        i := SCENARIOS[scene]["answer"];
                        newscene := SCENARIOS[scene]["dests"][i];
                } else {
                        newscene := SCENARIOS[scene]["next"];
                }
        } else {
                newscene := SCENARIOS[scene]["next"];
        }

    }
    print("Goodbye.");
}
