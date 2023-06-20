package random

import (
	"strconv"
	"strings"
)

// Text random generator for various text values like names, addresses or phone numbers.
//
//	Example:
//		value1 := RandomText.FullName();    // Possible result: "Sergio"
//		value2 := RandomText.Verb();      	// Possible result: "Run"
//		value3 := RandomText.Text(50);    	// Possible result: "Run jorge. Red high scream?"
var Text = &_TRandomText{}

type _TRandomText struct{}

var namePrefixes = []string{"Dr.", "Mr.", "Mrs"}
var nameSuffixes = []string{"Jr.", "Sr.", "II", "III"}
var firstNames = []string{
	"John", "Bill", "Andrew", "Nick", "Pamela", "Bela", "Sergio", "George", "Hurry", "Cecilia", "Vesta", "Terry", "Patrick",
}
var lastNames = []string{
	"Doe", "Smith", "Johns", "Gates", "Carmack", "Zontak", "Clinton", "Adams", "First", "Lopez", "Due", "White", "Black",
}
var colors = []string{
	"Black", "White", "Red", "Blue", "Green", "Yellow", "Purple", "Grey", "Magenta", "Cian",
}
var stuffs = []string{
	"Game", "Ball", "Home", "Board", "Car", "Plane", "Hotel", "Wine", "Pants", "Boots", "Table", "Chair",
}
var adjectives = []string{
	"Large", "Small", "High", "Low", "Certain", "Fuzzy", "Modern", "Faster", "Slower",
}
var verbs = []string{
	"Run", "Stay", "Breeze", "Fly", "Lay", "Write", "Draw", "Scream",
}

// var streetTypes = []string{
//     "Lane", "Court", "Circle", "Drive", "Way", "Loop", "Blvd", "Street"
// }
// var streetPrefix = []string{
//     "North", "South", "East", "West", "Old", "New", "N.", "S.", "E.", "W."
// }
// var streetNames = []string{
//     "1st", "2nd", "3rd", "4th", "53rd", "6th", "8th", "Acacia", "Academy", "Adams", "Addison", "Airport", "Albany", "Alderwood", "Alton", "Amerige", "Amherst", "Anderson",
//     "Ann", "Annadale", "Applegate", "Arcadia", "Arch", "Argyle", "Arlington", "Armstrong", "Arnold", "Arrowhead", "Aspen", "Augusta", "Baker", "Bald Hill", "Bank", "Bay Meadows",
//     "Bay", "Bayberry", "Bayport", "Beach", "Beaver Ridge", "Bedford", "Beech", "Beechwood", "Belmont", "Berkshire", "Big Rock Cove", "Birch Hill", "Birchpond", "Birchwood",
//     "Bishop", "Blackburn", "Blue Spring", "Bohemia", "Border", "Boston", "Bow Ridge", "Bowman", "Bradford", "Brandywine", "Brewery", "Briarwood", "Brickell", "Brickyard",
//     "Bridge", "Bridgeton", "Bridle", "Broad", "Brookside", "Brown", "Buckingham", "Buttonwood", "Cambridge", "Campfire", "Canal", "Canterbury", "Cardinal", "Carpenter",
//     "Carriage", "Carson", "Catherine", "Cedar Swamp", "Cedar", "Cedarwood", "Cemetery", "Center", "Central", "Chapel", "Charles", "Cherry Hill", "Chestnut", "Church", "Circle",
//     "Clark", "Clay", "Cleveland", "Clinton", "Cobblestone", "Coffee", "College", "Colonial", "Columbia", "Cooper", "Corona", "Cottage", "Country Club", "Country", "County", "Court",
//     "Courtland", "Creek", "Creekside", "Crescent", "Cross", "Cypress", "Deerfield", "Del Monte", "Delaware", "Depot", "Devon", "Devonshire", "Division", "Dogwood", "Dunbar",
//     "Durham", "Eagle", "East", "Edgefield", "Edgemont", "Edgewater", "Edgewood", "El Dorado", "Elizabeth", "Elm", "Essex", "Euclid", "Evergreen", "Fairfield", "Fairground", "Fairview",
//     "Fairway", "Fawn", "Fifth", "Fordham", "Forest", "Foster", "Foxrun", "Franklin", "Fremont", "Front", "Fulton", "Galvin", "Garden", "Gartner", "Gates", "George", "Glen Creek",
//     "Glen Eagles", "Glen Ridge", "Glendale", "Glenlake", "Glenridge", "Glenwood", "Golden Star", "Goldfield", "Golf", "Gonzales", "Grand", "Grandrose", "Grant", "Green Hill",
//     "Green Lake", "Green", "Greenrose", "Greenview", "Gregory", "Griffin", "Grove", "Halifax", "Hamilton", "Hanover", "Harrison", "Hartford", "Harvard", "Harvey", "Hawthorne",
//     "Heather", "Henry Smith", "Heritage", "High Noon", "High Point", "High", "Highland", "Hill Field", "Hillcrest", "Hilldale", "Hillside", "Hilltop", "Holly", "Homestead",
//     "Homewood", "Honey Creek", "Howard", "Indian Spring", "Indian Summer", "Iroquois", "Jackson", "James", "Jefferson", "Jennings", "Jockey Hollow", "John", "Johnson", "Jones",
//     "Joy Ridge", "King", "Kingston", "Kirkland", "La Sierra", "Lafayette", "Lake Forest", "Lake", "Lakeshore", "Lakeview", "Lancaster", "Lane", "Laurel", "Leatherwood", "Lees Creek",
//     "Leeton Ridge", "Lexington", "Liberty", "Lilac", "Lincoln", "Linda", "Littleton", "Livingston", "Locust", "Longbranch", "Lookout", "Lower River", "Lyme", "Madison", "Maiden",
//     "Main", "Mammoth", "Manchester", "Manhattan", "Manor Station", "Maple", "Marconi", "Market", "Marsh", "Marshall", "Marvon", "Mayfair", "Mayfield", "Mayflower", "Meadow",
//     "Meadowbrook", "Mechanic", "Middle River", "Miles", "Mill Pond", "Miller", "Monroe", "Morris", "Mountainview", "Mulberry", "Myrtle", "Newbridge", "Newcastle", "Newport",
//     "Nichols", "Nicolls", "North", "Nut Swamp", "Oak Meadow", "Oak Valley", "Oak", "Oakland", "Oakwood", "Ocean", "Ohio", "Oklahoma", "Olive", "Orange", "Orchard", "Overlook",
//     "Pacific", "Paris Hill", "Park", "Parker", "Pawnee", "Peachtree", "Pearl", "Peg Shop", "Pendergast", "Peninsula", "Penn", "Pennington", "Pennsylvania", "Pheasant", "Philmont",
//     "Pierce", "Pin Oak", "Pine", "Pineknoll", "Piper", "Plumb Branch", "Poor House", "Prairie", "Primrose", "Prince", "Princess", "Princeton", "Proctor", "Prospect", "Pulaski",
//     "Pumpkin Hill", "Purple Finch", "Queen", "Race", "Ramblewood", "Redwood", "Ridge", "Ridgewood", "River", "Riverside", "Riverview", "Roberts", "Rock Creek", "Rock Maple",
//     "Rockaway", "Rockcrest", "Rockland", "Rockledge", "Rockville", "Rockwell", "Rocky River", "Roosevelt", "Rose", "Rosewood", "Ryan", "Saddle", "Sage", "San Carlos", "San Juan",
//     "San Pablo", "Santa Clara", "Saxon", "School", "Schoolhouse", "Second", "Shadow Brook", "Shady", "Sheffield", "Sherman", "Sherwood", "Shipley", "Shub Farm", "Sierra",
//     "Silver Spear", "Sleepy Hollow", "Smith Store", "Smoky Hollow", "Snake Hill", "Southampton", "Spring", "Spruce", "Squaw Creek", "St Louis", "St Margarets", "St Paul", "State",
//     "Stillwater", "Strawberry", "Studebaker", "Sugar", "Sulphur Springs", "Summerhouse", "Summit", "Sunbeam", "Sunnyslope", "Sunset", "Surrey", "Sutor", "Swanson", "Sycamore",
//     "Tailwater", "Talbot", "Tallwood", "Tanglewood", "Tarkiln Hill", "Taylor", "Thatcher", "Third", "Thomas", "Thompson", "Thorne", "Tower", "Trenton", "Trusel", "Tunnel",
//     "University", "Vale", "Valley Farms", "Valley View", "Valley", "Van Dyke", "Vermont", "Vernon", "Victoria", "Vine", "Virginia", "Wagon", "Wall", "Walnutwood", "Warren",
//     "Washington", "Water", "Wayne", "Westminster", "Westport", "White", "Whitemarsh", "Wild Rose", "William", "Williams", "Wilson", "Winchester", "Windfall", "Winding Way",
//     "Winding", "Windsor", "Wintergreen", "Wood", "Woodland", "Woodside", "Woodsman", "Wrangler", "York",
// }

var allWords = append(append(append(append(lastNames, colors...), stuffs...), adjectives...), verbs...)

// Color return string with random color
func (c *_TRandomText) Color() string {
	return String.Pick(colors)
}

// Stuff return string with random stuff
func (c *_TRandomText) Stuff() string {
	return String.Pick(stuffs)
}

// Adjective return string with random adjective
func (c *_TRandomText) Adjective() string {
	return String.Pick(adjectives)
}

// Verb return string with random verb
func (c *_TRandomText) Verb() string {
	return String.Pick(verbs)
}

// Phrase return string with random phrase
func (c *_TRandomText) Phrase(min int, max int) string {
	size := Integer.Next(min, max)
	if size <= 0 {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(String.Pick(allWords))
	for builder.Len() < size {
		builder.WriteString(" ")
		builder.WriteString(strings.ToLower(String.Pick(allWords)))
	}

	return builder.String()
}

// FullName return string with random full name
func (c *_TRandomText) FullName() string {
	builder := strings.Builder{}

	if Boolean.Chance(3, 5) {
		builder.WriteString(String.Pick(namePrefixes))
		builder.WriteString(" ")
	}

	builder.WriteString(String.Pick(firstNames))
	builder.WriteString(" ")
	builder.WriteString(String.Pick(lastNames))

	if Boolean.Chance(5, 10) {
		builder.WriteString(" ")
		builder.WriteString(String.Pick(nameSuffixes))
	}

	return builder.String()
}

// Word return string with random word
func (c *_TRandomText) Word() string {
	return String.Pick(allWords)
}

// Words return string with random words
func (c *_TRandomText) Words(min int, max int) string {
	builder := strings.Builder{}

	count := Integer.Next(min, max)
	for i := 0; i < count; i++ {
		builder.WriteString(String.Pick(allWords))
	}

	return builder.String()
}

// Phone return string with random phone
func (c *_TRandomText) Phone() string {
	builder := strings.Builder{}
	builder.WriteString("(")
	builder.WriteString(strconv.Itoa(Integer.Next(111, 999)))
	builder.WriteString(") ")
	builder.WriteString(strconv.Itoa(Integer.Next(111, 999)))
	builder.WriteString("-")
	builder.WriteString(strconv.Itoa(Integer.Next(1111, 9999)))
	return builder.String()
}

// Email return string with random email
func (c *_TRandomText) Email() string {
	builder := strings.Builder{}
	builder.WriteString(c.Words(2, 6))
	builder.WriteString("@")
	builder.WriteString(c.Words(1, 3))
	builder.WriteString(".com")
	return builder.String()
}

// Text return string with random text
func (c *_TRandomText) Text(min int, max int) string {
	size := Integer.Next(min, max)

	builder := strings.Builder{}
	builder.WriteString(String.Pick(allWords))

	for builder.Len() < size {
		next := String.Pick(allWords)
		if Boolean.Chance(4, 6) {
			builder.WriteString(" ")
			builder.WriteString(strings.ToLower(next))
		} else if Boolean.Chance(2, 5) {
			builder.WriteByte(String.PickChar(":,-"))
			builder.WriteString(strings.ToLower(next))
		} else if Boolean.Chance(3, 5) {
			builder.WriteByte(String.PickChar(":,-"))
			builder.WriteString(" ")
			builder.WriteString(strings.ToLower(next))
		} else {
			builder.WriteByte(String.PickChar(".!?"))
			builder.WriteString(" ")
			builder.WriteString(next)
		}
	}

	return builder.String()
}
