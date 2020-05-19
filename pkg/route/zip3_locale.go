package route

import (
	"fmt"
)

// Zip5ToZip3LatLong looks up a zip code and returns the Lat Long from the census data
func Zip5ToZip3LatLong(zip5 string) (LatLong, error) {
	var ll LatLong
	zip5only := formatZip5(zip5)
	if len(zip5only) > 5 {
		return ll, NewUnsupportedPostalCodeError(zip5only)
	} else if len(zip5only) < 5 {
		zip5only = fmt.Sprintf("%05s", zip5only)
	}
	zip3 := zip5only[0:3]
	if len(zip3) != 3 {
		return ll, NewUnsupportedPostalCodeError(zip3)
	}
	var ok bool
	ll, ok = zip3ToLatLongMap[zip3]
	if !ok {
		return ll, NewUnsupportedPostalCodeError(zip3)
	}

	return ll, nil
}

// zip3ToLatLongMap maps Zip3 (as string) to LatLong according to
// https://pe.usps.com/Archive/HTML/DMMArchive20050106/print/L002.htm and
// simplemaps.com/data/us-cities
//
// pe.usps.com is used to get the Zip3 to City mapping and
// simplemaps.com is used to get the City to LatLong mapping
// Using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes to resolve descrepencies
var zip3ToLatLongMap = map[string]LatLong{
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"005": {40.8123, -73.0447}, // Holtsville, NY
	"006": {18.4037, -66.0636}, // San Juan, PR
	"007": {18.4037, -66.0636}, // San Juan, PR
	"008": {18.4037, -66.0636}, // San Juan, PR
	"009": {18.4037, -66.0636}, // San Juan, PR
	"010": {42.1155, -72.5395}, // Springfield, MA
	"011": {42.1155, -72.5395}, // Springfield, MA
	"012": {42.4517, -73.2605}, // Pittsfield, MA
	"013": {42.1155, -72.5395}, // Springfield, MA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"014": {42.5912, -71.8156}, // Fitchburg, MA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"015": {42.2705, -71.8079}, // Worcester, MA
	"016": {42.2705, -71.8079}, // Worcester, MA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"017": {42.3085, -71.4368}, // Framingham, MA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"018": {42.4869, -71.1543}, // Woburn, MA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"019": {42.4779, -70.9663}, // Lynn, MA
	"020": {42.0821, -71.0242}, // Brockton, MA
	"021": {42.3188, -71.0846}, // Boston, MA
	"022": {42.3188, -71.0846}, // Boston, MA
	"023": {42.0821, -71.0242}, // Brockton, MA
	// Manual fix for ZIP <024> with USPS city NORTHWEST BOS and state MA
	"024": {42.4473, -71.2272}, // Lexington, MA
	// Manual fix fo ZIP <025> with USPS city CAPE COD and state MA
	"025": {41.6688, -70.2962}, // Cape Cod, MA
	// Manual fix for ZIP <026> with USPS city CAPE COD and state MA
	"026": {41.6688, -70.2962}, // Cape Cod, MA
	"027": {41.8230, -71.4187}, // Providence, RI
	"028": {41.8230, -71.4187}, // Providence, RI
	"029": {41.8230, -71.4187}, // Providence, RI
	"030": {42.9848, -71.4447}, // Manchester, NH
	"031": {42.9848, -71.4447}, // Manchester, NH
	"032": {42.9848, -71.4447}, // Manchester, NH
	"033": {43.2305, -71.5595}, // Concord, NH
	"034": {42.9848, -71.4447}, // Manchester, NH
	// Manual fix ZIP <035> with USPS city WHITE RIV JCT and state VT
	"035": {44.3062, -71.7701}, // Littleton, NH
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"036": {43.1344, -72.4550}, // Bellows Falls, VT
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"037": {43.6496, -72.3239}, // White River Junction, VT
	"038": {43.0580, -70.7826}, // Portsmouth, NH
	"039": {43.0580, -70.7826}, // Portsmouth, NH
	"040": {43.6773, -70.2715}, // Portland, ME
	"041": {43.6773, -70.2715}, // Portland, ME
	"042": {43.6773, -70.2715}, // Portland, ME
	"043": {43.6773, -70.2715}, // Portland, ME
	"044": {44.8322, -68.7906}, // Bangor, ME
	"045": {43.6773, -70.2715}, // Portland, ME
	"046": {44.8322, -68.7906}, // Bangor, ME
	"047": {44.8322, -68.7906}, // Bangor, ME
	"048": {43.6773, -70.2715}, // Portland, ME
	"049": {44.8322, -68.7906}, // Bangor, ME
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"050": {43.6496, -72.3239}, // White River Junction, VT
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"051": {43.1344, -72.4550}, // Bellows Falls, VT
	// Manual fix for ZIP <052> with USPS city WHITE RIV JCT and state VT
	"052": {42.8781, -73.1968}, //Bennington, VT
	// Manual fix for ZIP <053> with USPS city WHITE RIV JCT and state VT
	"053": {42.8509, -72.5579}, // Brattleboro, VT
	"054": {44.4877, -73.2314}, // Burlington, VT
	// Manual fix for ZIP <055> with USPS city MIDDLESEX-ESX and state MA
	"055": {42.6583, -71.1368}, // Andover, MA
	"056": {44.4877, -73.2314}, // Burlington, VT
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"057": {43.6091, -72.9781}, // Rutland, VT
	// Manual fix for ZIP <058> with USPS city WHITE RIV JCT and state VT
	"058": {44.4193, -72.0151}, // St. Johnsbury, VT
	// Manual fix for ZIP <059> with USPS city WHITE RIV JCT and state VT
	"059": {44.4193, -72.0151}, // St. Johnsbury, VT
	"060": {41.7661, -72.6834}, // Hartford, CT
	"061": {41.7661, -72.6834}, // Hartford, CT
	"062": {41.7661, -72.6834}, // Hartford, CT
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"063": {41.3502, -72.1023}, // New London, CT
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"064": {41.3112, -72.9246}, // New Haven, CT
	"065": {41.3112, -72.9246}, // New Haven, CT
	"066": {41.1918, -73.1953}, // Bridgeport, CT
	"067": {41.5583, -73.0361}, // Waterbury, CT
	"068": {41.1035, -73.5583}, // Stamford, CT
	"069": {41.1035, -73.5583}, // Stamford, CT
	"070": {40.7245, -74.1725}, // Newark, NJ
	"071": {40.7245, -74.1725}, // Newark, NJ
	"072": {40.6657, -74.1912}, // Elizabeth, NJ
	"073": {40.7161, -74.0683}, // Jersey City, NJ
	"074": {40.9147, -74.1624}, // Paterson, NJ
	"075": {40.9147, -74.1624}, // Paterson, NJ
	"076": {40.8890, -74.0461}, // Hackensack, NJ
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"077": {40.3481, -74.0672}, // Red Bank, NJ
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"078": {40.8859, -74.5597}, // Dover, NJ
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"079": {40.7154, -74.3647}, // Summit, NJ
	// Manual fix for ZIP <080> with USPS city SOUTH JERSEY and state NJ
	"080": {39.9268, -75.0246}, // Cherry Hill, NJ
	"081": {39.9362, -75.1073}, // Camden, NJ
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"082": {39.3797, -74.4527}, // Atlantic City, NJ
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"083": {39.4286, -75.2281}, // Bridgeton, NJ
	"084": {39.3797, -74.4527}, // Atlantic City, NJ
	"085": {40.2236, -74.7641}, // Trenton, NJ
	"086": {40.2236, -74.7641}, // Trenton, NJ
	// Manual fix fox ZIP <087> with USPS city MONMOUTH and state NJ
	"087": {40.0821, -74.2097}, // Lakewood, NJ
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"088": {40.4870, -74.4450}, // New Brunswick, NJ
	"089": {40.4870, -74.4450}, // New Brunswick, NJ
	/*
		TODO PLACEHOLDER TO FIX ZIP <090> with USPS city APO and state AE
		<090> Maps to wiki zip3 city Military and state AE; Military AE note: Bases in DE
		TODO PLACEHOLDER TO FIX ZIP <091> with USPS city APO and state AE
		<091> Maps to wiki zip3 city Military and state AE; Military AE note: Bases in DE
		TODO PLACEHOLDER TO FIX ZIP <092> with USPS city APO and state AE
		<092> Maps to wiki zip3 city Military and state AE; Military AE note: Bases in DE
		TODO PLACEHOLDER TO FIX ZIP <093> with USPS city APO and state AE
		<093> Maps to wiki zip3 city Military and state AE; Military AE note: Bases in IQ, AF, other
		TODO PLACEHOLDER TO FIX ZIP <094> with USPS city APO/FPO and state AE
		<094> Maps to wiki zip3 city Military and state AE; Military AE note: Bases in UK
		TODO PLACEHOLDER TO FIX ZIP <095> with USPS city FPO and state AE
		<095> Maps to wiki zip3 city Military and state AE; Military AE note: Naval/Marine
		TODO PLACEHOLDER TO FIX ZIP <096> with USPS city APO/FPO and state AE
		<096> Maps to wiki zip3 city Military and state AE; Military AE note: Bases in ES, IT
		TODO PLACEHOLDER TO FIX ZIP <097> with USPS city APO/FPO and state AE
		<097> Maps to wiki zip3 city Military and state AE; Military AE note: Bases in EU, GL, CA
		TODO PLACEHOLDER TO FIX ZIP <098> with USPS city APO/FPO and state AE
		<098> Maps to wiki zip3 city Military and state AE; Military AE note: Africa/Middle East
		TODO PLACEHOLDER TO FIX ZIP <099> with USPS city APO/FPO and state AE
		<099> Maps to wiki zip3 city Military and state AE; Military AE note: Africa/Middle East
	*/
	"100": {40.6943, -73.9249}, // New York, NY
	"101": {40.6943, -73.9249}, // New York, NY
	"102": {40.6943, -73.9249}, // New York, NY
	"103": {40.5834, -74.1496}, // Staten Island, NY
	"104": {40.8501, -73.8662}, // Bronx, NY
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"105": {41.0220, -73.7549}, // White Plains, NY
	"106": {41.0220, -73.7549}, // White Plains, NY
	"107": {40.9466, -73.8674}, // Yonkers, NY
	"108": {40.9305, -73.7836}, // New Rochelle, NY
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"109": {41.1138, -74.1421}, // Suffern, NY
	"110": {40.7498, -73.7976}, // Queens, NY
	// Manual fix for ZIP <111> with USPS city LONG ISLAND CITY and state NY
	"111": {40.7128, -74.0060}, // Long Island City, NY (but google points to New York City)
	"112": {40.6501, -73.9496}, // Brooklyn, NY
	// Manual fix for ZIP <113> with USPS city FLUSHING and state NY
	"113": {40.7675, -73.8331}, // Flushing, NY
	// Manual fix for ZIP <114> with USPS city JAMAICA and state NY
	"114": {40.7027, -73.7890}, // Jamaica, NY
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"115": {40.7469, -73.6392}, // Mineola, NY
	// Manual fix for ZIP <116> with USPS city FAR ROCKAWAY and state NY
	"116": {40.5999, -73.7448}, // Far Rockaway, NY
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"117": {40.6696, -73.4156}, // Amityville, NY
	"118": {40.7637, -73.5245}, // Hicksville, NY
	// Manual fix for ZIP <119> with USPS city MID-ISLAND and state NY
	"119": {40.9170, -72.6620}, // Riverhead, NY
	"120": {42.6664, -73.7987}, // Albany, NY
	"121": {42.6664, -73.7987}, // Albany, NY
	"122": {42.6664, -73.7987}, // Albany, NY
	"123": {42.8025, -73.9276}, // Schenectady, NY
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"124": {41.9295, -73.9968}, // Kingston, NY
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"125": {41.6949, -73.9210}, // Poughkeepsie, NY
	"126": {41.6949, -73.9210}, // Poughkeepsie, NY
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"127": {41.6523, -74.6876}, // Monticello, NY
	"128": {43.3109, -73.6459}, // Glens Falls, NY
	"129": {44.6951, -73.4563}, // Plattsburgh, NY
	"130": {43.0409, -76.1438}, // Syracuse, NY
	"131": {43.0409, -76.1438}, // Syracuse, NY
	"132": {43.0409, -76.1438}, // Syracuse, NY
	"133": {43.0961, -75.2260}, // Utica, NY
	"134": {43.0961, -75.2260}, // Utica, NY
	"135": {43.0961, -75.2260}, // Utica, NY
	"136": {43.9734, -75.9095}, // Watertown, NY
	"137": {42.1014, -75.9093}, // Binghamton, NY
	"138": {42.1014, -75.9093}, // Binghamton, NY
	"139": {42.1014, -75.9093}, // Binghamton, NY
	"140": {42.9017, -78.8487}, // Buffalo, NY
	"141": {42.9017, -78.8487}, // Buffalo, NY
	"142": {42.9017, -78.8487}, // Buffalo, NY
	"143": {43.0921, -79.0147}, // Niagara Falls, NY
	"144": {43.1680, -77.6162}, // Rochester, NY
	"145": {43.1680, -77.6162}, // Rochester, NY
	"146": {43.1680, -77.6162}, // Rochester, NY
	"147": {42.0975, -79.2366}, // Jamestown, NY
	"148": {42.0938, -76.8097}, // Elmira, NY
	"149": {42.0938, -76.8097}, // Elmira, NY
	"150": {40.4396, -79.9763}, // Pittsburgh, PA
	"151": {40.4396, -79.9763}, // Pittsburgh, PA
	"152": {40.4396, -79.9763}, // Pittsburgh, PA
	"153": {40.4396, -79.9763}, // Pittsburgh, PA
	"154": {40.4396, -79.9763}, // Pittsburgh, PA
	"155": {40.3258, -78.9194}, // Johnstown, PA
	"156": {40.3113, -79.5444}, // Greensburg, PA
	"157": {40.3258, -78.9194}, // Johnstown, PA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"158": {41.1225, -78.7564}, // DuBois, PA
	"159": {40.3258, -78.9194}, // Johnstown, PA
	"160": {40.9956, -80.3458}, // New Castle, PA
	"161": {40.9956, -80.3458}, // New Castle, PA
	"162": {40.9956, -80.3458}, // New Castle, PA
	"163": {41.4282, -79.7035}, // Oil City, PA
	"164": {42.1168, -80.0733}, // Erie, PA
	"165": {42.1168, -80.0733}, // Erie, PA
	"166": {40.5082, -78.4007}, // Altoona, PA
	"167": {41.9604, -78.6413}, // Bradford, PA
	"168": {40.5082, -78.4007}, // Altoona, PA
	"169": {41.2398, -77.0371}, // Williamsport, PA
	"170": {40.2752, -76.8843}, // Harrisburg, PA
	"171": {40.2752, -76.8843}, // Harrisburg, PA
	"172": {40.2752, -76.8843}, // Harrisburg, PA
	"173": {40.0420, -76.3012}, // Lancaster, PA
	"174": {39.9651, -76.7315}, // York, PA
	"175": {40.0420, -76.3012}, // Lancaster, PA
	"176": {40.0420, -76.3012}, // Lancaster, PA
	"177": {41.2398, -77.0371}, // Williamsport, PA
	"178": {40.2752, -76.8843}, // Harrisburg, PA
	"179": {40.3400, -75.9267}, // Reading, PA
	// Manual fix for ZIP <180> with USPS city LEHIGH VALLEY and state PA
	"180": {40.6520, -75.5742}, // Lehigh County, PA
	"181": {40.5961, -75.4755}, // Allentown, PA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"182": {40.9504, -75.9724}, // Hazleton, PA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"183": {41.0023, -75.1779}, // East Stroudsburg, PA
	"184": {41.4044, -75.6649}, // Scranton, PA
	"185": {41.4044, -75.6649}, // Scranton, PA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"186": {41.2468, -75.8759}, // Wilkes-Barre, PA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"187": {41.2468, -75.8759}, // Wilkes-Barre, PA
	"188": {41.4044, -75.6649}, // Scranton, PA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"189": {40.3139, -75.1280}, // Doylestown, PA
	"190": {40.0077, -75.1339}, // Philadelphia, PA
	"191": {40.0077, -75.1339}, // Philadelphia, PA
	"192": {40.0077, -75.1339}, // Philadelphia, PA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"193": {40.0420, -75.4912}, // Paoli, PA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"194": {40.1224, -75.3398}, // Norristown, PA
	"195": {40.3400, -75.9267}, // Reading, PA
	"196": {40.3400, -75.9267}, // Reading, PA
	"197": {39.7415, -75.5413}, // Wilmington, DE
	"198": {39.7415, -75.5413}, // Wilmington, DE
	"199": {39.7415, -75.5413}, // Wilmington, DE
	"200": {38.9047, -77.0163}, // Washington, DC
	// Manual fix for ZIP <201> with USPS city DULLES and state VA
	"201": {38.9625, -77.4380}, // Dulles, VA
	"202": {38.9047, -77.0163}, // Washington, DC
	"203": {38.9047, -77.0163}, // Washington, DC
	"204": {38.9047, -77.0163}, // Washington, DC
	"205": {38.9047, -77.0163}, // Washington, DC
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"206": {38.6085, -76.9195}, // Waldorf, MD
	// Manual fix for ZIP <207> with USPS city SOUTHERN MD and state MD
	"207": {38.9784, -76.4922}, // Annapolis, MD
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"208": {38.9866, -77.1188}, // Bethesda, MD
	"209": {39.0028, -77.0207}, // Silver Spring, MD
	"210": {39.2088, -76.6625}, // Linthicum, MD
	"211": {39.2088, -76.6625}, // Linthicum, MD
	"212": {39.3051, -76.6144}, // Baltimore, MD
	"214": {38.9706, -76.5047}, // Annapolis, MD
	"215": {39.6515, -78.7585}, // Cumberland, MD
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"216": {38.7760, -76.0702}, // Easton, MD
	"217": {39.4336, -77.4157}, // Frederick, MD
	"218": {38.3755, -75.5867}, // Salisbury, MD
	"219": {39.3051, -76.6144}, // Baltimore, MD
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"220": {38.9047, -77.0163}, // Washington, DC
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"221": {38.9047, -77.0163}, // Washington, DC
	"222": {38.8786, -77.1011}, // Arlington, VA
	"223": {38.8185, -77.0861}, // Alexandria, VA
	"224": {37.5295, -77.4756}, // Richmond, VA
	"225": {37.5295, -77.4756}, // Richmond, VA
	"226": {39.1735, -78.1746}, // Winchester, VA
	"227": {38.4705, -78.0001}, // Culpeper, VA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"228": {38.4362, -78.8735}, // Harrisonburg, VA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"229": {38.0375, -78.4855}, // Charlottesville, VA
	"230": {37.5295, -77.4756}, // Richmond, VA
	"231": {37.5295, -77.4756}, // Richmond, VA
	"232": {37.5295, -77.4756}, // Richmond, VA
	"233": {36.8945, -76.2590}, // Norfolk, VA
	"234": {36.8945, -76.2590}, // Norfolk, VA
	"235": {36.8945, -76.2590}, // Norfolk, VA
	"236": {36.8945, -76.2590}, // Norfolk, VA
	"237": {36.8468, -76.3540}, // Portsmouth, VA
	"238": {37.5295, -77.4756}, // Richmond, VA
	"239": {37.2959, -78.4002}, // Farmville, VA
	"240": {37.2785, -79.9580}, // Roanoke, VA
	"241": {37.2785, -79.9580}, // Roanoke, VA
	"242": {36.6179, -82.1607}, // Bristol, VA
	"243": {37.2785, -79.9580}, // Roanoke, VA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"244": {38.1593, -79.0611}, // Staunton, VA
	"245": {37.4003, -79.1909}, // Lynchburg, VA
	"246": {37.2608, -81.2143}, // Bluefield, WV
	"247": {37.2608, -81.2143}, // Bluefield, WV
	"248": {37.2608, -81.2143}, // Bluefield, WV
	"249": {37.8096, -80.4327}, // Lewisburg, WV
	"250": {38.3484, -81.6323}, // Charleston, WV
	"251": {38.3484, -81.6323}, // Charleston, WV
	"252": {38.3484, -81.6323}, // Charleston, WV
	"253": {38.3484, -81.6323}, // Charleston, WV
	"254": {39.4582, -77.9776}, // Martinsburg, WV
	"255": {38.4109, -82.4344}, // Huntington, WV
	"256": {38.4109, -82.4344}, // Huntington, WV
	"257": {38.4109, -82.4344}, // Huntington, WV
	"258": {37.7877, -81.1840}, // Beckley, WV
	"259": {37.7877, -81.1840}, // Beckley, WV
	"260": {40.0751, -80.6951}, // Wheeling, WV
	"261": {39.2624, -81.5419}, // Parkersburg, WV
	"262": {39.2863, -80.3230}, // Clarksburg, WV
	"263": {39.2863, -80.3230}, // Clarksburg, WV
	"264": {39.2863, -80.3230}, // Clarksburg, WV
	"265": {39.2863, -80.3230}, // Clarksburg, WV
	"266": {38.6702, -80.7717}, // Gassaway, WV
	"267": {39.6515, -78.7585}, // Cumberland, MD
	"268": {38.9957, -79.1276}, // Petersburg, WV
	"270": {36.0956, -79.8268}, // Greensboro, NC
	"271": {36.1029, -80.2610}, // Winston-Salem, NC
	"272": {36.0956, -79.8268}, // Greensboro, NC
	"273": {36.0956, -79.8268}, // Greensboro, NC
	"274": {36.0956, -79.8268}, // Greensboro, NC
	"275": {35.8324, -78.6438}, // Raleigh, NC
	"276": {35.8324, -78.6438}, // Raleigh, NC
	"277": {35.9795, -78.9032}, // Durham, NC
	"278": {35.9676, -77.8047}, // Rocky Mount, NC
	"279": {35.9676, -77.8047}, // Rocky Mount, NC
	"280": {35.2079, -80.8304}, // Charlotte, NC
	"281": {35.2079, -80.8304}, // Charlotte, NC
	"282": {35.2079, -80.8304}, // Charlotte, NC
	"283": {35.0846, -78.9776}, // Fayetteville, NC
	"284": {35.0846, -78.9776}, // Fayetteville, NC
	"285": {35.2748, -77.5937}, // Kinston, NC
	"286": {35.7426, -81.3230}, // Hickory, NC
	"287": {35.5704, -82.5537}, // Asheville, NC
	"288": {35.5704, -82.5537}, // Asheville, NC
	"289": {35.5704, -82.5537}, // Asheville, NC
	"290": {34.0376, -80.9037}, // Columbia, SC
	"291": {34.0376, -80.9037}, // Columbia, SC
	"292": {34.0376, -80.9037}, // Columbia, SC
	"293": {34.8362, -82.3649}, // Greenville, SC
	"294": {32.8151, -79.9630}, // Charleston, SC
	"295": {34.1782, -79.7872}, // Florence, SC
	"296": {34.8362, -82.3649}, // Greenville, SC
	"297": {35.2079, -80.8304}, // Charlotte, NC
	"298": {33.3645, -82.0708}, // Augusta, GA
	"299": {32.0281, -81.1785}, // Savannah, GA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"300": {33.7627, -84.4225}, // Atlanta, GA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"301": {33.7627, -84.4225}, // Atlanta, GA
	"302": {33.7627, -84.4225}, // Atlanta, GA
	"303": {33.7627, -84.4225}, // Atlanta, GA
	"304": {32.5866, -82.3345}, // Swainsboro, GA
	"305": {33.9508, -83.3689}, // Athens, GA
	"306": {33.9508, -83.3689}, // Athens, GA
	"307": {35.0657, -85.2487}, // Chattanooga, TN
	"308": {33.3645, -82.0708}, // Augusta, GA
	"309": {33.3645, -82.0708}, // Augusta, GA
	"310": {32.8065, -83.6974}, // Macon, GA
	"311": {33.7627, -84.4225}, // Atlanta, GA
	"312": {32.8065, -83.6974}, // Macon, GA
	"313": {32.0281, -81.1785}, // Savannah, GA
	"314": {32.0281, -81.1785}, // Savannah, GA
	"315": {31.2108, -82.3579}, // Waycross, GA
	"316": {30.8502, -83.2788}, // Valdosta, GA
	"317": {31.5776, -84.1762}, // Albany, GA
	"318": {32.5100, -84.8771}, // Columbus, GA
	"319": {32.5100, -84.8771}, // Columbus, GA
	"320": {30.3322, -81.6749}, // Jacksonville, FL
	"321": {29.1994, -81.0982}, // Daytona Beach, FL
	"322": {30.3322, -81.6749}, // Jacksonville, FL
	"323": {30.4551, -84.2527}, // Tallahassee, FL
	"324": {30.1995, -85.6003}, // Panama City, FL
	"325": {30.4427, -87.1886}, // Pensacola, FL
	"326": {29.6804, -82.3458}, // Gainesville, FL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"327": {28.4772, -81.3369}, // Orlando, FL
	"328": {28.4772, -81.3369}, // Orlando, FL
	"329": {28.4772, -81.3369}, // Orlando, FL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"330": {25.7839, -80.2102}, // Miami, FL
	"331": {25.7839, -80.2102}, // Miami, FL
	"332": {25.7839, -80.2102}, // Miami, FL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"333": {26.1412, -80.1464}, // Fort Lauderdale, FL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"334": {26.7469, -80.1316}, // West Palm Beach, FL
	"335": {27.9942, -82.4451}, // Tampa, FL
	"336": {27.9942, -82.4451}, // Tampa, FL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"337": {27.7930, -82.6652}, // St. Petersburg, FL
	"338": {28.0557, -81.9545}, // Lakeland, FL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"339": {26.6195, -81.8303}, // Fort Myers, FL
	/*
		TODO PLACEHOLDER TO FIX ZIP <340> with USPS city APO/FPO and state AA
		<340> Maps to wiki zip3 city Military and state AA; Military AA note: Americas
	*/
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"341": {26.1505, -81.7936}, // Naples, FL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"342": {27.4900, -82.5740}, // Bradenton, FL
	"344": {29.6804, -82.3458}, // Gainesville, FL
	"346": {27.9942, -82.4451}, // Tampa, FL
	"347": {28.4772, -81.3369}, // Orlando, FL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"349": {27.4256, -80.3430}, // Fort Pierce, FL
	"350": {33.5277, -86.7987}, // Birmingham, AL
	"351": {33.5277, -86.7987}, // Birmingham, AL
	"352": {33.5277, -86.7987}, // Birmingham, AL
	"354": {33.2348, -87.5266}, // Tuscaloosa, AL
	"355": {33.5277, -86.7987}, // Birmingham, AL
	"356": {34.6988, -86.6412}, // Huntsville, AL
	"357": {34.6988, -86.6412}, // Huntsville, AL
	"358": {34.6988, -86.6412}, // Huntsville, AL
	"359": {33.5277, -86.7987}, // Birmingham, AL
	"360": {32.3473, -86.2666}, // Montgomery, AL
	"361": {32.3473, -86.2666}, // Montgomery, AL
	"362": {33.6713, -85.8136}, // Anniston, AL
	"363": {31.2335, -85.4068}, // Dothan, AL
	"364": {31.4342, -86.9723}, // Evergreen, AL
	"365": {30.6782, -88.1163}, // Mobile, AL
	"366": {30.6782, -88.1163}, // Mobile, AL
	"367": {32.3473, -86.2666}, // Montgomery, AL
	"368": {32.3473, -86.2666}, // Montgomery, AL
	"369": {32.3846, -88.6897}, // Meridian, MS
	"370": {36.1715, -86.7843}, // Nashville, TN
	"371": {36.1715, -86.7843}, // Nashville, TN
	"372": {36.1715, -86.7843}, // Nashville, TN
	"373": {35.0657, -85.2487}, // Chattanooga, TN
	"374": {35.0657, -85.2487}, // Chattanooga, TN
	"375": {35.1046, -89.9773}, // Memphis, TN
	"376": {36.3406, -82.3803}, // Johnson City, TN
	"377": {35.9692, -83.9496}, // Knoxville, TN
	"378": {35.9692, -83.9496}, // Knoxville, TN
	"379": {35.9692, -83.9496}, // Knoxville, TN
	"380": {35.1046, -89.9773}, // Memphis, TN
	"381": {35.1046, -89.9773}, // Memphis, TN
	"382": {36.1371, -88.5077}, // McKenzie, TN
	"383": {35.6536, -88.8353}, // Jackson, TN
	"384": {35.6236, -87.0487}, // Columbia, TN
	"385": {36.1484, -85.5114}, // Cookeville, TN
	"386": {35.1046, -89.9773}, // Memphis, TN
	"387": {33.3850, -91.0514}, // Greenville, MS
	"388": {34.2691, -88.7318}, // Tupelo, MS
	"389": {33.7816, -89.8130}, // Grenada, MS
	"390": {32.3163, -90.2124}, // Jackson, MS
	"391": {32.3163, -90.2124}, // Jackson, MS
	"392": {32.3163, -90.2124}, // Jackson, MS
	"393": {32.3846, -88.6897}, // Meridian, MS
	"394": {31.3074, -89.3170}, // Hattiesburg, MS
	"395": {30.4271, -89.0703}, // Gulfport, MS
	"396": {31.2449, -90.4714}, // McComb, MS
	"397": {33.5088, -88.4097}, // Columbus, MS
	"398": {31.5776, -84.1762}, // Albany, GA
	"399": {33.7627, -84.4225}, // Atlanta, GA
	"400": {38.1663, -85.6485}, // Louisville, KY
	"401": {38.1663, -85.6485}, // Louisville, KY
	"402": {38.1663, -85.6485}, // Louisville, KY
	"403": {38.0423, -84.4587}, // Lexington, KY
	"404": {38.0423, -84.4587}, // Lexington, KY
	"405": {38.0423, -84.4587}, // Lexington, KY
	"406": {38.1924, -84.8643}, // Frankfort, KY
	"407": {37.1209, -84.0804}, // London, KY
	"408": {37.1209, -84.0804}, // London, KY
	"409": {37.1209, -84.0804}, // London, KY
	"410": {39.1412, -84.5060}, // Cincinnati, OH
	"411": {38.4593, -82.6449}, // Ashland, KY
	"412": {38.4593, -82.6449}, // Ashland, KY
	"413": {37.7353, -83.5473}, // Campton, KY
	"414": {37.7353, -83.5473}, // Campton, KY
	"415": {37.4807, -82.5262}, // Pikeville, KY
	"416": {37.4807, -82.5262}, // Pikeville, KY
	"417": {37.2583, -83.1976}, // Hazard, KY
	"418": {37.2583, -83.1976}, // Hazard, KY
	"420": {37.0711, -88.6435}, // Paducah, KY
	"421": {36.9715, -86.4375}, // Bowling Green, KY
	"422": {36.9715, -86.4375}, // Bowling Green, KY
	"423": {37.7573, -87.1174}, // Owensboro, KY
	"424": {37.9881, -87.5341}, // Evansville, IN
	"425": {37.0816, -84.6089}, // Somerset, KY
	"426": {37.0816, -84.6089}, // Somerset, KY
	"427": {37.7030, -85.8769}, // Elizabethtown, KY
	"430": {39.9860, -82.9851}, // Columbus, OH
	"431": {39.9860, -82.9851}, // Columbus, OH
	"432": {39.9860, -82.9851}, // Columbus, OH
	"433": {39.9860, -82.9851}, // Columbus, OH
	"434": {41.6639, -83.5822}, // Toledo, OH
	"435": {41.6639, -83.5822}, // Toledo, OH
	"436": {41.6639, -83.5822}, // Toledo, OH
	"437": {39.9567, -82.0133}, // Zanesville, OH
	"438": {39.9567, -82.0133}, // Zanesville, OH
	"439": {40.3653, -80.6520}, // Steubenville, OH
	"440": {41.4767, -81.6805}, // Cleveland, OH
	"441": {41.4767, -81.6805}, // Cleveland, OH
	"442": {41.0798, -81.5219}, // Akron, OH
	"443": {41.0798, -81.5219}, // Akron, OH
	"444": {41.0993, -80.6463}, // Youngstown, OH
	"445": {41.0993, -80.6463}, // Youngstown, OH
	"446": {40.8076, -81.3678}, // Canton, OH
	"447": {40.8076, -81.3678}, // Canton, OH
	"448": {40.7656, -82.5275}, // Mansfield, OH
	"449": {40.7656, -82.5275}, // Mansfield, OH
	"450": {39.1412, -84.5060}, // Cincinnati, OH
	"451": {39.1412, -84.5060}, // Cincinnati, OH
	"452": {39.1412, -84.5060}, // Cincinnati, OH
	"453": {39.7797, -84.1998}, // Dayton, OH
	"454": {39.7797, -84.1998}, // Dayton, OH
	"455": {39.9297, -83.7957}, // Springfield, OH
	"456": {39.3393, -82.9937}, // Chillicothe, OH
	"457": {39.3269, -82.0987}, // Athens, OH
	"458": {40.7410, -84.1121}, // Lima, OH
	"459": {39.1412, -84.5060}, // Cincinnati, OH
	"460": {39.7771, -86.1458}, // Indianapolis, IN
	"461": {39.7771, -86.1458}, // Indianapolis, IN
	"462": {39.7771, -86.1458}, // Indianapolis, IN
	"463": {41.5906, -87.3472}, // Gary, IN
	"464": {41.5906, -87.3472}, // Gary, IN
	"465": {41.6771, -86.2692}, // South Bend, IN
	"466": {41.6771, -86.2692}, // South Bend, IN
	"467": {41.0885, -85.1436}, // Fort Wayne, IN
	"468": {41.0885, -85.1436}, // Fort Wayne, IN
	"469": {40.4640, -86.1277}, // Kokomo, IN
	"470": {39.1412, -84.5060}, // Cincinnati, OH
	"471": {38.1663, -85.6485}, // Louisville, KY
	"472": {39.2094, -85.9183}, // Columbus, IN
	"473": {40.1989, -85.3950}, // Muncie, IN
	"474": {39.1637, -86.5257}, // Bloomington, IN
	"475": {39.4654, -87.3763}, // Terre Haute, IN
	"476": {37.9881, -87.5341}, // Evansville, IN
	"477": {37.9881, -87.5341}, // Evansville, IN
	"478": {39.4654, -87.3763}, // Terre Haute, IN
	"479": {40.3990, -86.8593}, // Lafayette, IN
	"480": {42.5084, -83.1539}, // Royal Oak, MI
	"481": {42.3834, -83.1024}, // Detroit, MI
	"482": {42.3834, -83.1024}, // Detroit, MI
	"483": {42.5084, -83.1539}, // Royal Oak, MI
	"484": {43.0235, -83.6922}, // Flint, MI
	"485": {43.0235, -83.6922}, // Flint, MI
	"486": {43.4199, -83.9501}, // Saginaw, MI
	"487": {43.4199, -83.9501}, // Saginaw, MI
	"488": {42.7142, -84.5601}, // Lansing, MI
	"489": {42.7142, -84.5601}, // Lansing, MI
	"490": {42.2749, -85.5882}, // Kalamazoo, MI
	"491": {42.2749, -85.5882}, // Kalamazoo, MI
	"492": {42.2431, -84.4037}, // Jackson, MI
	"493": {42.9615, -85.6557}, // Grand Rapids, MI
	"494": {42.9615, -85.6557}, // Grand Rapids, MI
	"495": {42.9615, -85.6557}, // Grand Rapids, MI
	"496": {44.7547, -85.6035}, // Traverse City, MI
	"497": {45.0214, -84.6803}, // Gaylord, MI
	"498": {45.8275, -88.0599}, // Iron Mountain, MI
	"499": {45.8275, -88.0599}, // Iron Mountain, MI
	"500": {41.5725, -93.6105}, // Des Moines, IA
	"501": {41.5725, -93.6105}, // Des Moines, IA
	"502": {41.5725, -93.6105}, // Des Moines, IA
	"503": {41.5725, -93.6105}, // Des Moines, IA
	"504": {42.4920, -92.3522}, // Waterloo, IA
	"505": {42.5098, -94.1751}, // Fort Dodge, IA
	"506": {42.4920, -92.3522}, // Waterloo, IA
	"507": {42.4920, -92.3522}, // Waterloo, IA
	"508": {41.0597, -94.3650}, // Creston, IA
	"509": {41.5725, -93.6105}, // Des Moines, IA
	"510": {42.4959, -96.3901}, // Sioux City, IA
	"511": {42.4959, -96.3901}, // Sioux City, IA
	"512": {42.4959, -96.3901}, // Sioux City, IA
	"513": {42.4959, -96.3901}, // Sioux City, IA
	"514": {42.0699, -94.8647}, // Carroll, IA
	"515": {41.2628, -96.0498}, // Omaha, NE
	"516": {41.2628, -96.0498}, // Omaha, NE
	"520": {42.5007, -90.7067}, // Dubuque, IA
	"521": {43.3016, -91.7846}, // Decorah, IA
	"522": {41.9667, -91.6781}, // Cedar Rapids, IA
	"523": {41.9667, -91.6781}, // Cedar Rapids, IA
	"524": {41.9667, -91.6781}, // Cedar Rapids, IA
	"525": {41.5725, -93.6105}, // Des Moines, IA
	"526": {40.8072, -91.1247}, // Burlington, IA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"527": {41.5563, -90.6052}, // Davenport, IA
	"528": {41.5563, -90.6052}, // Davenport, IA
	"530": {43.0642, -87.9673}, // Milwaukee, WI
	"531": {43.0642, -87.9673}, // Milwaukee, WI
	"532": {43.0642, -87.9673}, // Milwaukee, WI
	"534": {42.7274, -87.8135}, // Racine, WI
	"535": {43.0827, -89.3923}, // Madison, WI
	"537": {43.0827, -89.3923}, // Madison, WI
	"538": {43.0827, -89.3923}, // Madison, WI
	"539": {43.5489, -89.4658}, // Portage, WI
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"540": {44.9477, -93.1040}, // St. Paul, MN
	"541": {44.5150, -87.9896}, // Green Bay, WI
	"542": {44.5150, -87.9896}, // Green Bay, WI
	"543": {44.5150, -87.9896}, // Green Bay, WI
	"544": {44.9615, -89.6457}, // Wausau, WI
	"545": {45.6360, -89.4256}, // Rhinelander, WI
	"546": {43.8241, -91.2268}, // La Crosse, WI
	"547": {44.8200, -91.4951}, // Eau Claire, WI
	"548": {45.8271, -91.8860}, // Spooner, WI
	"549": {44.0228, -88.5617}, // Oshkosh, WI
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"550": {44.9477, -93.1040}, // St. Paul, MN
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"551": {44.9477, -93.1040}, // St. Paul, MN
	"553": {44.9635, -93.2678}, // Minneapolis, MN
	"554": {44.9635, -93.2678}, // Minneapolis, MN
	"555": {44.9635, -93.2678}, // Minneapolis, MN
	"556": {46.7757, -92.1392}, // Duluth, MN
	"557": {46.7757, -92.1392}, // Duluth, MN
	"558": {46.7757, -92.1392}, // Duluth, MN
	"559": {44.0151, -92.4778}, // Rochester, MN
	"560": {44.1711, -93.9773}, // Mankato, MN
	"561": {44.1711, -93.9773}, // Mankato, MN
	"562": {45.1220, -95.0569}, // Willmar, MN
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"563": {45.5339, -94.1718}, // St. Cloud, MN
	"564": {46.3553, -94.1983}, // Brainerd, MN
	"565": {46.8060, -95.8449}, // Detroit Lakes, MN
	"566": {47.4830, -94.8788}, // Bemidji, MN
	"567": {47.9221, -97.0887}, // Grand Forks, ND
	"570": {43.5397, -96.7321}, // Sioux Falls, SD
	"571": {43.5397, -96.7321}, // Sioux Falls, SD
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"572": {44.9094, -97.1532}, // Watertown, SD
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"573": {43.7296, -98.0337},  // Mitchell, SD
	"574": {45.4646, -98.4680},  // Aberdeen, SD
	"575": {44.3748, -100.3205}, // Pierre, SD
	"576": {45.5411, -100.4349}, // Mobridge, SD
	"577": {44.0716, -103.2205}, // Rapid City, SD
	"580": {46.8653, -96.8292},  // Fargo, ND
	"581": {46.8653, -96.8292},  // Fargo, ND
	"582": {47.9221, -97.0887},  // Grand Forks, ND
	"583": {48.1131, -98.8753},  // Devils Lake, ND
	"584": {46.9063, -98.6937},  // Jamestown, ND
	"585": {46.8140, -100.7695}, // Bismarck, ND
	"586": {46.8140, -100.7695}, // Bismarck, ND
	"587": {48.2374, -101.2780}, // Minot, ND
	"588": {48.1814, -103.6364}, // Williston, ND
	"590": {45.7889, -108.5509}, // Billings, MT
	"591": {45.7889, -108.5509}, // Billings, MT
	"592": {48.0933, -105.6413}, // Wolf Point, MT
	"593": {46.4059, -105.8385}, // Miles City, MT
	"594": {47.5022, -111.2995}, // Great Falls, MT
	"595": {48.5427, -109.6804}, // Havre, MT
	"596": {46.5965, -112.0199}, // Helena, MT
	"597": {45.9020, -112.6571}, // Butte, MT
	"598": {46.8685, -114.0095}, // Missoula, MT
	"599": {48.2156, -114.3261}, // Kalispell, MT
	"600": {42.1181, -88.0430},  // Palatine, IL
	"601": {41.9182, -88.1308},  // Carol Stream, IL
	"602": {42.0463, -87.6942},  // Evanston, IL
	"603": {41.8872, -87.7899},  // Oak Park, IL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"604": {41.8373, -87.6862}, // Chicago, IL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"605": {41.8373, -87.6862}, // Chicago, IL
	"606": {41.8373, -87.6862}, // Chicago, IL
	"607": {41.8373, -87.6862}, // Chicago, IL
	"608": {41.8373, -87.6862}, // Chicago, IL
	"609": {41.1020, -87.8643}, // Kankakee, IL
	"610": {42.2598, -89.0641}, // Rockford, IL
	"611": {42.2598, -89.0641}, // Rockford, IL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"612": {41.4699, -90.5827}, // Rock Island, IL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"613": {41.3575, -89.0718}, // LaSalle, IL
	"614": {40.9506, -90.3763}, // Galesburg, IL
	"615": {40.7521, -89.6155}, // Peoria, IL
	"616": {40.7521, -89.6155}, // Peoria, IL
	"617": {40.4757, -88.9703}, // Bloomington, IL
	"618": {40.1144, -88.2735}, // Champaign, IL
	"619": {40.1144, -88.2735}, // Champaign, IL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"620": {38.6358, -90.2451}, // St. Louis, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"622": {38.6156, -90.1304}, // East St. Louis, IL
	"623": {39.9335, -91.3798}, // Quincy, IL
	"624": {39.1207, -88.5509}, // Effingham, IL
	"625": {39.7710, -89.6537}, // Springfield, IL
	"626": {39.7710, -89.6537}, // Springfield, IL
	"627": {39.7710, -89.6537}, // Springfield, IL
	"628": {38.5224, -89.1233}, // Centralia, IL
	"629": {37.7220, -89.2238}, // Carbondale, IL
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"630": {38.6358, -90.2451}, // St. Louis, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"631": {38.6358, -90.2451}, // St. Louis, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"633": {38.7956, -90.5156}, // St. Charles, MO
	"634": {39.9335, -91.3798}, // Quincy, IL
	"635": {39.9335, -91.3798}, // Quincy, IL
	"636": {37.3108, -89.5596}, // Cape Girardeau, MO
	"637": {37.3108, -89.5596}, // Cape Girardeau, MO
	"638": {37.3108, -89.5596}, // Cape Girardeau, MO
	"639": {37.3108, -89.5596}, // Cape Girardeau, MO
	"640": {39.1239, -94.5541}, // Kansas City, MO
	"641": {39.1239, -94.5541}, // Kansas City, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"644": {39.7598, -94.8210}, // St. Joseph, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"645": {39.7598, -94.8210}, // St. Joseph, MO
	"646": {39.7953, -93.5498}, // Chillicothe, MO
	"647": {38.6530, -94.3467}, // Harrisonville, MO
	"648": {37.1943, -93.2915}, // Springfield, MO
	"649": {39.1239, -94.5541}, // Kansas City, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"650": {38.5676, -92.1759}, // Jefferson City, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"651": {38.5676, -92.1759}, // Jefferson City, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"652": {38.9477, -92.3255}, // Columbia, MO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"653": {38.7042, -93.2351}, // Sedalia, MO
	"654": {37.1943, -93.2915}, // Springfield, MO
	"655": {37.1943, -93.2915}, // Springfield, MO
	"656": {37.1943, -93.2915}, // Springfield, MO
	"657": {37.1943, -93.2915}, // Springfield, MO
	"658": {37.1943, -93.2915}, // Springfield, MO
	"660": {39.1234, -94.7443}, // Kansas City, KS
	"661": {39.1234, -94.7443}, // Kansas City, KS
	"662": {39.1234, -94.7443}, // Kansas City, KS
	"664": {39.0346, -95.6955}, // Topeka, KS
	"665": {39.0346, -95.6955}, // Topeka, KS
	"666": {39.0346, -95.6955}, // Topeka, KS
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"667": {37.8283, -94.7038},  // Fort Scott, KS
	"668": {39.0346, -95.6955},  // Topeka, KS
	"669": {38.8137, -97.6143},  // Salina, KS
	"670": {37.6897, -97.3441},  // Wichita, KS
	"671": {37.6897, -97.3441},  // Wichita, KS
	"672": {37.6897, -97.3441},  // Wichita, KS
	"673": {37.2118, -95.7328},  // Independence, KS
	"674": {38.8137, -97.6143},  // Salina, KS
	"675": {38.0671, -97.9081},  // Hutchinson, KS
	"676": {38.8816, -99.3219},  // Hays, KS
	"677": {39.3843, -101.0459}, // Colby, KS
	"678": {37.7610, -100.0182}, // Dodge City, KS
	"679": {37.0466, -100.9295}, // Liberal, KS
	"680": {41.2628, -96.0498},  // Omaha, NE
	"681": {41.2628, -96.0498},  // Omaha, NE
	"683": {40.8088, -96.6796},  // Lincoln, NE
	"684": {40.8088, -96.6796},  // Lincoln, NE
	"685": {40.8088, -96.6796},  // Lincoln, NE
	"686": {42.0328, -97.4209},  // Norfolk, NE
	"687": {42.0328, -97.4209},  // Norfolk, NE
	"688": {40.9214, -98.3584},  // Grand Island, NE
	"689": {40.9214, -98.3584},  // Grand Island, NE
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"690": {40.2046, -100.6213}, // McCook, NE
	"691": {41.1266, -100.7640}, // North Platte, NE
	"692": {42.8739, -100.5498}, // Valentine, NE
	"693": {42.1025, -102.8766}, // Alliance, NE
	"700": {30.0687, -89.9288},  // New Orleans, LA
	"701": {30.0687, -89.9288},  // New Orleans, LA
	"703": {29.5799, -90.7058},  // Houma, LA
	"704": {30.3750, -90.0906},  // Mandeville, LA
	"705": {30.2084, -92.0323},  // Lafayette, LA
	"706": {30.2022, -93.2141},  // Lake Charles, LA
	"707": {30.4419, -91.1310},  // Baton Rouge, LA
	"708": {30.4419, -91.1310},  // Baton Rouge, LA
	"710": {32.4659, -93.7959},  // Shreveport, LA
	"711": {32.4659, -93.7959},  // Shreveport, LA
	"712": {32.5183, -92.0775},  // Monroe, LA
	"713": {31.2923, -92.4702},  // Alexandria, LA
	"714": {31.2923, -92.4702},  // Alexandria, LA
	"716": {34.2116, -92.0178},  // Pine Bluff, AR
	"717": {33.5672, -92.8467},  // Camden, AR
	"718": {33.4361, -93.9960},  // Texarkana, AR
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"719": {34.4902, -93.0498}, // Hot Springs, AR
	"720": {34.7255, -92.3580}, // Little Rock, AR
	"721": {34.7255, -92.3580}, // Little Rock, AR
	"722": {34.7255, -92.3580}, // Little Rock, AR
	"723": {35.1046, -89.9773}, // Memphis, TN
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"724": {35.8211, -90.6793}, // Jonesboro, AR
	"725": {35.7687, -91.6226}, // Batesville, AR
	"726": {36.2438, -93.1198}, // Harrison, AR
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"727": {36.0713, -94.1660},  // Fayetteville, AR
	"728": {35.2763, -93.1383},  // Russellville, AR
	"729": {35.3493, -94.3695},  // Fort Smith, AR
	"730": {35.4676, -97.5137},  // Oklahoma City, OK
	"731": {35.4676, -97.5137},  // Oklahoma City, OK
	"733": {30.3006, -97.7517},  // Austin, TX
	"734": {34.1943, -97.1253},  // Ardmore, OK
	"735": {34.6176, -98.4203},  // Lawton, OK
	"736": {35.5058, -98.9724},  // Clinton, OK
	"737": {36.4061, -97.8701},  // Enid, OK
	"738": {36.4246, -99.4057},  // Woodward, OK
	"739": {37.0466, -100.9295}, // Liberal, KS
	"740": {36.1284, -95.9043},  // Tulsa, OK
	"741": {36.1284, -95.9043},  // Tulsa, OK
	"743": {36.1284, -95.9043},  // Tulsa, OK
	"744": {35.7430, -95.3566},  // Muskogee, OK
	"745": {34.9262, -95.7698},  // McAlester, OK
	"746": {36.7235, -97.0679},  // Ponca City, OK
	"747": {33.9957, -96.3938},  // Durant, OK
	"748": {35.3525, -96.9647},  // Shawnee, OK
	"749": {35.0430, -94.6357},  // Poteau, OK
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"750": {32.7936, -96.7662}, // Dallas, TX
	"751": {32.7936, -96.7662}, // Dallas, TX
	"752": {32.7936, -96.7662}, // Dallas, TX
	"753": {32.7936, -96.7662}, // Dallas, TX
	"754": {33.1116, -96.1099}, // Greenville, TX
	"755": {33.4487, -94.0815}, // Texarkana, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"756": {32.5192, -94.7622}, // Longview, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"757": {32.3184, -95.3065}, // Tyler, TX
	"758": {31.7544, -95.6471}, // Palestine, TX
	"759": {31.3217, -94.7277}, // Lufkin, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"760": {32.7812, -97.3472}, // Fort Worth, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"761": {32.7812, -97.3472}, // Fort Worth, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"762": {33.2176, -97.1419}, // Denton, TX
	"763": {33.9072, -98.5290}, // Wichita Falls, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"764": {32.2148, -98.2205},  // Stephenville, TX
	"765": {31.5597, -97.1882},  // Waco, TX
	"766": {31.5597, -97.1882},  // Waco, TX
	"767": {31.5597, -97.1882},  // Waco, TX
	"768": {32.4543, -99.7384},  // Abilene, TX
	"769": {32.0249, -102.1137}, // Midland, TX
	"770": {29.7869, -95.3905},  // Houston, TX
	"771": {29.7869, -95.3905},  // Houston, TX
	"772": {29.7869, -95.3905},  // Houston, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"773": {30.3224, -95.4820}, // Conroe, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"774": {29.5824, -95.7602}, // Richmond, TX
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"775": {29.6584, -95.1499},  // Pasadena, TX
	"776": {30.0850, -94.1451},  // Beaumont, TX
	"777": {30.0850, -94.1451},  // Beaumont, TX
	"778": {30.6657, -96.3668},  // Bryan, TX
	"779": {28.8285, -96.9850},  // Victoria, TX
	"780": {29.4658, -98.5254},  // San Antonio, TX
	"781": {29.4658, -98.5254},  // San Antonio, TX
	"782": {29.4658, -98.5254},  // San Antonio, TX
	"783": {27.7261, -97.3755},  // Corpus Christi, TX
	"784": {27.7261, -97.3755},  // Corpus Christi, TX
	"785": {26.2273, -98.2471},  // McAllen, TX
	"786": {30.3006, -97.7517},  // Austin, TX
	"787": {30.3006, -97.7517},  // Austin, TX
	"788": {29.4658, -98.5254},  // San Antonio, TX
	"789": {30.3006, -97.7517},  // Austin, TX
	"790": {35.1989, -101.8310}, // Amarillo, TX
	"791": {35.1989, -101.8310}, // Amarillo, TX
	"792": {34.4293, -100.2516}, // Childress, TX
	"793": {33.5642, -101.8871}, // Lubbock, TX
	"794": {33.5642, -101.8871}, // Lubbock, TX
	"795": {32.4543, -99.7384},  // Abilene, TX
	"796": {32.4543, -99.7384},  // Abilene, TX
	"797": {32.0249, -102.1137}, // Midland, TX
	"798": {31.8479, -106.4309}, // El Paso, TX
	"799": {31.8479, -106.4309}, // El Paso, TX
	"800": {39.7621, -104.8759}, // Denver, CO
	"801": {39.7621, -104.8759}, // Denver, CO
	"802": {39.7621, -104.8759}, // Denver, CO
	"803": {40.0249, -105.2523}, // Boulder, CO
	"804": {39.7621, -104.8759}, // Denver, CO
	"805": {40.1690, -105.0996}, // Longmont, CO
	"806": {39.7621, -104.8759}, // Denver, CO
	"807": {39.7621, -104.8759}, // Denver, CO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"808": {38.8674, -104.7606}, // Colorado Springs, CO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"809": {38.8674, -104.7606}, // Colorado Springs, CO
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"810": {38.2713, -104.6105}, // Pueblo, CO
	"811": {37.4755, -105.8770}, // Alamosa, CO
	"812": {38.5300, -105.9984}, // Salida, CO
	"813": {37.2744, -107.8703}, // Durango, CO
	"814": {39.0877, -108.5673}, // Grand Junction, CO
	"815": {39.0877, -108.5673}, // Grand Junction, CO
	"816": {39.5455, -107.3347}, // Glenwood Springs, CO
	"820": {41.1405, -104.7927}, // Cheyenne, WY
	// Manual fix for ZIP <821> with USPS city YELLOWSTONE NL PK and state WY
	"821": {44.4280, -110.5885}, // Yellowstone, WY
	"822": {42.0516, -104.9595}, // Wheatland, WY
	"823": {41.7849, -107.2265}, // Rawlins, WY
	"824": {44.0026, -107.9543}, // Worland, WY
	"825": {43.0425, -108.4142}, // Riverton, WY
	"826": {42.8420, -106.3207}, // Casper, WY
	"827": {44.2752, -105.4984}, // Gillette, WY
	"828": {44.7962, -106.9643}, // Sheridan, WY
	"829": {41.5951, -109.2237}, // Rock Springs, WY
	"830": {41.5951, -109.2237}, // Rock Springs, WY
	"831": {41.5951, -109.2237}, // Rock Springs, WY
	"832": {42.8716, -112.4652}, // Pocatello, ID
	"833": {42.5648, -114.4617}, // Twin Falls, ID
	"834": {42.8716, -112.4652}, // Pocatello, ID
	"835": {46.3934, -116.9934}, // Lewiston, ID
	"836": {43.6007, -116.2312}, // Boise, ID
	"837": {43.6007, -116.2312}, // Boise, ID
	"838": {47.6671, -117.4330}, // Spokane, WA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"840": {40.7774, -111.9300}, // Salt Lake City, UT
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"841": {40.7774, -111.9300}, // Salt Lake City, UT
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"842": {41.2280, -111.9677}, // Ogden, UT
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"843": {41.2280, -111.9677}, // Ogden, UT
	"844": {41.2280, -111.9677}, // Ogden, UT
	"845": {40.2457, -111.6457}, // Provo, UT
	"846": {40.2457, -111.6457}, // Provo, UT
	"847": {40.2457, -111.6457}, // Provo, UT
	"850": {33.5722, -112.0891}, // Phoenix, AZ
	"852": {33.5722, -112.0891}, // Phoenix, AZ
	"853": {33.5722, -112.0891}, // Phoenix, AZ
	"855": {33.3869, -110.7514}, // Globe, AZ
	"856": {32.1545, -110.8782}, // Tucson, AZ
	"857": {32.1545, -110.8782}, // Tucson, AZ
	"859": {34.2671, -110.0384}, // Show Low, AZ
	"860": {35.1872, -111.6195}, // Flagstaff, AZ
	"863": {34.5850, -112.4475}, // Prescott, AZ
	"864": {35.2170, -114.0105}, // Kingman, AZ
	"865": {35.5183, -108.7423}, // Gallup, NM
	"870": {35.1053, -106.6464}, // Albuquerque, NM
	"871": {35.1053, -106.6464}, // Albuquerque, NM
	"872": {35.1053, -106.6464}, // Albuquerque, NM
	"873": {35.5183, -108.7423}, // Gallup, NM
	"874": {36.7555, -108.1823}, // Farmington, NM
	"875": {35.1053, -106.6464}, // Albuquerque, NM
	"877": {35.6011, -105.2206}, // Las Vegas, NM
	"878": {34.0543, -106.9066}, // Socorro, NM
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"879": {33.1864, -107.2589}, // Truth or Consequences, NM
	"880": {32.3265, -106.7893}, // Las Cruces, NM
	"881": {34.4376, -103.1923}, // Clovis, NM
	"882": {33.3730, -104.5294}, // Roswell, NM
	"883": {32.8837, -105.9624}, // Alamogordo, NM
	"884": {35.1701, -103.7042}, // Tucumcari, NM
	"885": {31.8479, -106.4309}, // El Paso, TX
	"889": {36.2333, -115.2654}, // Las Vegas, NV
	"890": {36.2333, -115.2654}, // Las Vegas, NV
	"891": {36.2333, -115.2654}, // Las Vegas, NV
	"893": {39.2649, -114.8709}, // Ely, NV
	"894": {39.5497, -119.8483}, // Reno, NV
	"895": {39.5497, -119.8483}, // Reno, NV
	"897": {39.1511, -119.7474}, // Carson City, NV
	"898": {40.8387, -115.7674}, // Elko, NV
	"900": {34.1139, -118.4068}, // Los Angeles, CA
	"901": {34.1139, -118.4068}, // Los Angeles, CA
	"902": {33.9566, -118.3444}, // Inglewood, CA
	"903": {33.9566, -118.3444}, // Inglewood, CA
	"904": {34.0232, -118.4813}, // Santa Monica, CA
	"905": {33.8346, -118.3417}, // Torrance, CA
	"906": {33.7980, -118.1675}, // Long Beach, CA
	"907": {33.7980, -118.1675}, // Long Beach, CA
	"908": {33.7980, -118.1675}, // Long Beach, CA
	"910": {34.1597, -118.1390}, // Pasadena, CA
	"911": {34.1597, -118.1390}, // Pasadena, CA
	"912": {34.1818, -118.2468}, // Glendale, CA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"913": {34.1914, -118.8755}, // Thousand Oaks, CA
	// Manual fix for ZIP <914> with USPS city VAN NUYS and state CA
	"914": {34.1899, -118.4514}, // Van Nuys, CA
	"915": {34.1879, -118.3235}, // Burbank, CA
	// Manual fix for ZIP <916> with USPS city NORTH HOLLYWOOD and state CA
	"916": {34.1870, -118.3813}, // North Hollywood, CA
	"917": {34.0175, -117.9268}, // Industry, CA
	"918": {34.0175, -117.9268}, // Industry, CA
	"919": {32.8312, -117.1225}, // San Diego, CA
	"920": {32.8312, -117.1225}, // San Diego, CA
	"921": {32.8312, -117.1225}, // San Diego, CA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"922": {33.7346, -116.2346}, // Indio, CA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"923": {34.1417, -117.2945}, // San Bernardino, CA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"924": {34.1417, -117.2945}, // San Bernardino, CA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"925": {33.9381, -117.3948}, // Riverside, CA
	"926": {33.7366, -117.8819}, // Santa Ana, CA
	"927": {33.7366, -117.8819}, // Santa Ana, CA
	"928": {33.8390, -117.8573}, // Anaheim, CA
	"930": {34.1962, -119.1819}, // Oxnard, CA
	"931": {34.4285, -119.7202}, // Santa Barbara, CA
	"932": {35.3530, -119.0359}, // Bakersfield, CA
	"933": {35.3530, -119.0359}, // Bakersfield, CA
	"934": {34.4285, -119.7202}, // Santa Barbara, CA
	"935": {35.0139, -118.1895}, // Mojave, CA
	"936": {36.7831, -119.7941}, // Fresno, CA
	"937": {36.7831, -119.7941}, // Fresno, CA
	"938": {36.7831, -119.7941}, // Fresno, CA
	"939": {36.6884, -121.6317}, // Salinas, CA
	"940": {37.7562, -122.4430}, // San Francisco, CA
	"941": {37.7562, -122.4430}, // San Francisco, CA
	"942": {38.5667, -121.4683}, // Sacramento, CA
	"943": {37.3913, -122.1467}, // Palo Alto, CA
	"944": {37.5522, -122.3122}, // San Mateo, CA
	"945": {37.7903, -122.2165}, // Oakland, CA
	"946": {37.7903, -122.2165}, // Oakland, CA
	"947": {37.8723, -122.2760}, // Berkeley, CA
	"948": {37.9477, -122.3390}, // Richmond, CA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"949": {37.9904, -122.5222}, // San Rafael, CA
	"950": {37.3021, -121.8489}, // San Jose, CA
	"951": {37.3021, -121.8489}, // San Jose, CA
	"952": {37.9766, -121.3111}, // Stockton, CA
	"953": {37.9766, -121.3111}, // Stockton, CA
	// Lookup using https://en.wikipedia.org/wiki/List_of_ZIP_Code_prefixes
	"954": {38.4458, -122.7067}, // Santa Rosa, CA
	"955": {40.7941, -124.1568}, // Eureka, CA
	"956": {38.5667, -121.4683}, // Sacramento, CA
	"957": {38.5667, -121.4683}, // Sacramento, CA
	"958": {38.5667, -121.4683}, // Sacramento, CA
	"959": {39.1518, -121.5836}, // Marysville, CA
	"960": {40.5698, -122.3650}, // Redding, CA
	"961": {39.5497, -119.8483}, // Reno, NV
	/*
		TODO PLACEHOLDER TO FIX ZIP <962> with USPS city APO/FPO and state AP
		<962> Maps to wiki zip3 city Military and state AP; Military AP note: Bases in KR
		TODO PLACEHOLDER TO FIX ZIP <963> with USPS city APO/FPO and state AP
		<963> Maps to wiki zip3 city Military and state AP; Military AP note: Bases in JP
		TODO PLACEHOLDER TO FIX ZIP <964> with USPS city APO/FPO and state AP
		<964> Maps to wiki zip3 city Military and state AP; Military AP note: Bases in PH
		TODO PLACEHOLDER TO FIX ZIP <965> with USPS city APO/FPO and state AP
		<965> Maps to wiki zip3 city Military and state AP; Military AP note: Pacific & Antarctic bases
		TODO PLACEHOLDER TO FIX ZIP <966> with USPS city FPO and state AP
		<966> Maps to wiki zip3 city Military and state AP; Military AP note: Naval/Marine
	*/
	"967": {21.3294, -157.8460}, // Honolulu, HI
	"968": {21.3294, -157.8460}, // Honolulu, HI
	// Manual fix for ZIP <969> with USPS city BARRIGADA and state GU
	"969": {13.4708, -144.8181}, // Barrigada, GU
	"970": {45.5371, -122.6500}, // Portland, OR
	"971": {45.5371, -122.6500}, // Portland, OR
	"972": {45.5371, -122.6500}, // Portland, OR
	"973": {44.9232, -123.0245}, // Salem, OR
	"974": {44.0563, -123.1173}, // Eugene, OR
	"975": {42.3372, -122.8537}, // Medford, OR
	"976": {42.2191, -121.7754}, // Klamath Falls, OR
	"977": {44.0562, -121.3087}, // Bend, OR
	"978": {45.6755, -118.8209}, // Pendleton, OR
	"979": {43.6007, -116.2312}, // Boise, ID
	"980": {47.6211, -122.3244}, // Seattle, WA
	"981": {47.6211, -122.3244}, // Seattle, WA
	"982": {47.9524, -122.1670}, // Everett, WA
	"983": {47.2431, -122.4531}, // Tacoma, WA
	"984": {47.2431, -122.4531}, // Tacoma, WA
	"985": {47.0417, -122.8959}, // Olympia, WA
	"986": {45.5371, -122.6500}, // Portland, OR
	"988": {47.4338, -120.3286}, // Wenatchee, WA
	"989": {46.5923, -120.5496}, // Yakima, WA
	"990": {47.6671, -117.4330}, // Spokane, WA
	"991": {47.6671, -117.4330}, // Spokane, WA
	"992": {47.6671, -117.4330}, // Spokane, WA
	"993": {46.2506, -119.1303}, // Pasco, WA
	"994": {46.3934, -116.9934}, // Lewiston, ID
	"995": {61.1508, -149.1091}, // Anchorage, AK
	"996": {61.1508, -149.1091}, // Anchorage, AK
	"997": {64.8353, -147.6534}, // Fairbanks, AK
	"998": {58.4546, -134.1739}, // Juneau, AK
	"999": {55.3556, -131.6698}, // Ketchikan, AK
}
