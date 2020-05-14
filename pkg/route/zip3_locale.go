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
	// Manual fix for zip3 <005> with city MID-ISLAND and state NY
	"005":{40.8123, -73.0447}, // Holtsville, NY
	"006":{18.4037, -66.0636}, // San Juan, PR
	"007":{18.4037, -66.0636}, // San Juan, PR
	"008":{18.4037, -66.0636}, // San Juan, PR
	"009":{18.4037, -66.0636}, // San Juan, PR
	"010":{42.1155, -72.5395}, // Springfield, MA
	"011":{42.1155, -72.5395}, // Springfield, MA
	"012":{42.4517, -73.2605}, // Pittsfield, MA
	"013":{42.1155, -72.5395}, // Springfield, MA
	// Manual fix for zip3 <014> with city CENTRAL and state MA
	"014":{42.5912, -71.8156}, // Fitchburg, MA
	// Manual fix for zip3 <015> with city CENTRAL and state MA
	"015":{42.2705, -71.8079}, // Worcester, MA
	"016":{42.2705, -71.8079}, // Worcester, MA
	// Manual fix for zip3 <017> with city CENTRAL and state MA
	"017":{42.3085, -71.4368}, // Framingham, MA
	// Manual fix for zip3 <018> with city MIDDLESEX-ESX and state MA
	"018":{42.4869, -71.1543}, // Woburn, MA
	// Manual fix for zip3 <019> with city MIDDLESEX-ESX and state MA
	"019":{42.4779, -70.9663}, // Lynn, MA
	"020":{42.0821, -71.0242}, // Brockton, MA
	"021":{42.3188, -71.0846}, // Boston, MA
	"022":{42.3188, -71.0846}, // Boston, MA
	"023":{42.0821, -71.0242}, // Brockton, MA
	// Manual fix for zip3 <024> with city NORTHWEST BOS and state MA
	"024": {42.4473, -71.2272}, // Lexington, MA
	// Manual fix for zip3 <025> with city CAPE COD and state MA
	"025":{41.6688, -70.2962}, // Cape Cod, MA
	// Manual fix for zip3 <026> with city CAPE COD and state MA
	"026":{41.6688, -70.2962}, // Cape Cod, MA
	"027":{41.8230, -71.4187}, // Providence, RI
	"028":{41.8230, -71.4187}, // Providence, RI
	"029":{41.8230, -71.4187}, // Providence, RI
	"030":{42.9848, -71.4447}, // Manchester, NH
	"031":{42.9848, -71.4447}, // Manchester, NH
	"032":{42.9848, -71.4447}, // Manchester, NH
	"033":{43.2305, -71.5595}, // Concord, NH
	"034":{42.9848, -71.4447}, // Manchester, NH
	// Manual fix for zip3 <035> with city WHITE RIV JCT and state VT
	"035":{43.6496, -72.3239}, // White River Junction, VT
	// Manual fix for zip3 <036> with city WHITE RIV JCT and state VT
	"036":{43.6496, -72.3239}, // White River Junction, VT
	// Manual fix for zip3 <037> with city WHITE RIV JCT and state VT
	"037":{43.6496, -72.3239}, // White River Junction, VT
	"038":{43.0580, -70.7826}, // Portsmouth, NH
	"039":{43.0580, -70.7826}, // Portsmouth, NH
	"040":{43.6773, -70.2715}, // Portland, ME
	"041":{43.6773, -70.2715}, // Portland, ME
	"042":{43.6773, -70.2715}, // Portland, ME
	"043":{43.6773, -70.2715}, // Portland, ME
	"044":{44.8322, -68.7906}, // Bangor, ME
	"045":{43.6773, -70.2715}, // Portland, ME
	"046":{44.8322, -68.7906}, // Bangor, ME
	"047":{44.8322, -68.7906}, // Bangor, ME
	"048":{43.6773, -70.2715}, // Portland, ME
	"049":{44.8322, -68.7906}, // Bangor, ME
	// Manual fix for zip3 <050> with city WHITE RIV JCT and state VT
	"050":{43.6496, -72.3239}, // White River Junction, VT
	// Manual fix for zip3 <051> with city WHITE RIV JCT and state VT
	"051":{43.6496, -72.3239}, // White River Junction, VT
	// Manual fix for zip3 <052> with city WHITE RIV JCT and state VT
	"052":{43.6496, -72.3239}, // White River Junction, VT
	// Manual fix for zip3 <053> with city WHITE RIV JCT and state VT
	"053":{43.6496, -72.3239}, // White River Junction, VT
	"054":{44.4877, -73.2314}, // Burlington, VT
	// Manual fix for zip <055> with city MIDDLESEX-ESX and state MA
	"055":{42.6583, -71.1368}, // Andover, MA
	// END OF OPTIONS TO PICK
	"056":{44.4877, -73.2314}, // Burlington, VT
	// Manual fix for zip3 <057> with city WHITE RIV JCT and state VT
	"057":{43.6496, -72.3239}, // White River Junction, VT
	// Manual fix for zip3  <058> with city WHITE RIV JCT and state VT
	"058":{44.4193, -72.0151}, // Saint Johnsbury, VT
	// Manual fix for zip3 <059> with city WHITE RIV JCT and state
	"059":{44.4193, -72.0151}, // Saint Johnsbury, VT
	"060":{41.7661, -72.6834}, // Hartford, CT
	"061":{41.7661, -72.6834}, // Hartford, CT
	"062":{41.7661, -72.6834}, // Hartford, CT
	TODO PLACEHOLDER TO FIX ZIP <063> with city SOUTHERN and state CT
	"063":{41.2702, -71.9878}, // Fishers Island, NY
	"063":{41.3502, -72.1023}, // New London, CT
	"063":{41.3379, -72.0175}, // Poquonock Bridge, CT
	"063":{41.7501, -71.9091}, // Wauregan, CT
	"063":{41.6140, -72.0866}, // Baltic, CT
	"063":{41.6070, -71.9806}, // Jewett City, CT
	"063":{41.3536, -72.0517}, // Long Hill, CT
	"063":{41.5495, -72.0882}, // Norwich, CT
	"063":{41.3390, -72.0727}, // Groton, CT
	"063":{41.4644, -71.9745}, // Mashantucket, CT
	"063":{41.3774, -71.8492}, // Pawcatuck, CT
	"063":{41.6766, -71.9250}, // Plainfield Village, CT
	"063":{41.3265, -72.1949}, // Niantic, CT
	"063":{41.7169, -71.8750}, // Moosup, CT
	"063":{41.3855, -72.0686}, // Conning Towers Nautilus Park, CT
	"063":{41.3574, -71.9548}, // Mystic, CT
	"063":{41.3344, -71.9033}, // Stonington, CT
	"063":{41.3145, -72.0087}, // Groton Long Point, CT
	"063":{41.4212, -72.0859}, // Gales Ferry, CT
	"063":{41.3854, -71.9850}, // Old Mystic, CT
	"063":{41.3339, -71.9978}, // Noank, CT
	"063":{41.4441, -72.1250}, // Oxoboxo River, CT
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <064> with city SOUTHERN and state CT
	"064":{41.5832, -72.8915}, // Plantsville, CT
	"064":{41.6342, -72.4692}, // Terramuggus, CT
	"064":{41.3823, -72.4386}, // Deep River Center, CT
	"064":{41.3265, -73.0833}, // Derby, CT
	"064":{41.3552, -72.3910}, // Essex Village, CT
	"064":{41.5367, -72.7944}, // Meriden, CT
	"064":{41.3060, -73.1383}, // Shelton, CT
	"064":{41.2779, -72.8148}, // Branford Center, CT
	"064":{41.4900, -72.5562}, // Higganum, CT
	"064":{41.5961, -72.5124}, // Lake Pocotopaug, CT
	"064":{41.5026, -72.8993}, // Cheshire Village, CT
	"064":{41.2282, -72.9930}, // Woodmont, CT
	"064":{41.3443, -73.0689}, // Ansonia, CT
	"064":{41.5476, -72.6549}, // Middletown, CT
	"064":{41.2817, -72.6762}, // Guilford Center, CT
	"064":{41.4014, -72.4523}, // Chester Center, CT
	"064":{41.2918, -72.3682}, // Old Saybrook Center, CT
	"064":{41.2794, -72.6003}, // Madison Center, CT
	"064":{41.5043, -72.4491}, // Moodus, CT
	"064":{41.2825, -72.4063}, // Saybrook Manor, CT
	"064":{41.2811, -72.4424}, // Westbrook Center, CT
	"064":{41.4499, -72.8189}, // Wallingford Center, CT
	"064":{41.4119, -73.3120}, // Newtown, CT
	"064":{41.4845, -73.2351}, // Heritage Village, CT
	"064":{41.2711, -72.3546}, // Fenwick, CT
	"064":{41.2255, -73.0625}, // Milford city, CT
	// END OF OPTIONS TO PICK
	"065":{41.3112, -72.9246}, // New Haven, CT
	"066":{41.1918, -73.1953}, // Bridgeport, CT
	"067":{41.5583, -73.0361}, // Waterbury, CT
	"068":{41.1035, -73.5583}, // Stamford, CT
	"069":{41.1035, -73.5583}, // Stamford, CT
	"070":{40.7245, -74.1725}, // Newark, NJ
	"071":{40.7245, -74.1725}, // Newark, NJ
	"072":{40.6657, -74.1912}, // Elizabeth, NJ
	"073":{40.7161, -74.0683}, // Jersey City, NJ
	"074":{40.9147, -74.1624}, // Paterson, NJ
	"075":{40.9147, -74.1624}, // Paterson, NJ
	"076":{40.8890, -74.0461}, // Hackensack, NJ
	TODO PLACEHOLDER TO FIX ZIP <077> with city MONMOUTH and state NJ
	"077":{40.4489, -74.2495}, // Laurence Harbor, NJ
	"077":{40.2324, -74.2943}, // West Freehold, NJ
	"077":{40.4127, -74.2365}, // Matawan, NJ
	"077":{40.2323, -74.0014}, // Loch Arbour, NJ
	"077":{40.2005, -74.0334}, // Neptune City, NJ
	"077":{40.1913, -74.0162}, // Avon-by-the-Sea, NJ
	"077":{40.2366, -74.0294}, // Wanamassa, NJ
	"077":{40.3249, -74.0600}, // Shrewsbury, NJ
	"077":{40.4261, -74.0802}, // Belford, NJ
	"077":{40.2596, -74.2755}, // Freehold, NJ
	"077":{40.3619, -74.0392}, // Fair Haven, NJ
	"077":{40.2965, -73.9915}, // Long Branch, NJ
	"077":{40.4426, -74.2178}, // Cliffwood Beach, NJ
	"077":{40.2497, -73.9976}, // Deal, NJ
	"077":{40.2362, -74.0017}, // Allenhurst, NJ
	"077":{40.1707, -74.0376}, // West Belmar, NJ
	"077":{40.1522, -74.0430}, // Spring Lake Heights, NJ
	"077":{40.2119, -74.0078}, // Ocean Grove, NJ
	"077":{40.2708, -74.0948}, // Tinton Falls, NJ
	"077":{40.1706, -74.0262}, // Lake Como, NJ
	"077":{40.4469, -74.1316}, // Keansburg, NJ
	"077":{40.3064, -74.3385}, // Yorketown, NJ
	"077":{40.4338, -74.1009}, // Port Monmouth, NJ
	"077":{40.3391, -74.1283}, // Lincroft, NJ
	"077":{40.1983, -74.1700}, // Farmingdale, NJ
	"077":{40.2018, -74.0121}, // Bradley Beach, NJ
	"077":{40.1922, -74.0464}, // Shark River Hills, NJ
	"077":{40.4018, -74.2194}, // Strathmore, NJ
	"077":{40.4021, -74.0387}, // Navesink, NJ
	"077":{40.3364, -73.9863}, // Monmouth Beach, NJ
	"077":{40.1144, -74.1492}, // Ramtown, NJ
	"077":{40.4191, -74.0600}, // Leonardo, NJ
	"077":{40.1798, -74.0255}, // Belmar, NJ
	"077":{40.1384, -74.1032}, // Allenwood, NJ
	"077":{40.4454, -74.1699}, // Union Beach, NJ
	"077":{40.3756, -74.2444}, // Morganville, NJ
	"077":{40.3653, -73.9769}, // Sea Bright, NJ
	"077":{40.1538, -74.0271}, // Spring Lake, NJ
	"077":{40.4390, -74.1184}, // North Middletown, NJ
	"077":{40.3357, -74.0346}, // Little Silver, NJ
	"077":{40.3626, -74.0046}, // Rumson, NJ
	"077":{40.3395, -74.2939}, // Robertsville, NJ
	"077":{40.4036, -73.9898}, // Highlands, NJ
	"077":{40.2226, -74.0117}, // Asbury Park, NJ
	"077":{40.2718, -74.2425}, // East Freehold, NJ
	"077":{40.3481, -74.0672}, // Red Bank, NJ
	"077":{40.2913, -74.0558}, // Eatontown, NJ
	"077":{40.2350, -74.0166}, // Interlaken, NJ
	"077":{40.2971, -74.3607}, // Englishtown, NJ
	"077":{40.2607, -74.0263}, // Oakhurst, NJ
	"077":{40.4112, -74.0296}, // Atlantic Highlands, NJ
	"077":{40.3160, -74.0205}, // Oceanport, NJ
	"077":{40.2883, -74.0185}, // West Long Branch, NJ
	"077":{40.4327, -74.2011}, // Keyport, NJ
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <078> with city WEST JERSEY and state NJ
	"078":{40.8222, -74.8502}, // Beattystown, NJ
	"078":{40.9881, -74.9061}, // Marksboro, NJ
	"078":{40.9545, -75.0627}, // Hainesburg, NJ
	"078":{40.8294, -75.0728}, // Belvidere, NJ
	"078":{40.8859, -74.5597}, // Dover, NJ
	"078":{41.0534, -74.7527}, // Newton, NJ
	"078":{41.1234, -74.8403}, // Crandon Lakes, NJ
	"078":{40.8770, -74.9034}, // Great Meadows, NJ
	"078":{40.9541, -74.6593}, // Hopatcong, NJ
	"078":{40.9240, -74.5121}, // White Meadow Lake, NJ
	"078":{40.7621, -75.0141}, // Brass Castle, NJ
	"078":{40.8969, -74.5155}, // Rockaway, NJ
	"078":{40.7617, -74.9301}, // Anderson, NJ
	"078":{40.8999, -74.5808}, // Wharton, NJ
	"078":{40.8985, -74.7019}, // Netcong, NJ
	"078":{40.8762, -74.5435}, // Victory Gardens, NJ
	"078":{40.8507, -74.6596}, // Succasunna, NJ
	"078":{40.8733, -74.7374}, // Budd Lake, NJ
	"078":{40.9078, -74.8413}, // Panther Valley, NJ
	"078":{40.8698, -74.8810}, // Vienna, NJ
	"078":{40.9604, -74.4968}, // Lake Telemark, NJ
	"078":{40.7684, -74.9540}, // Port Colden, NJ
	"078":{41.0150, -74.6639}, // Lake Mohawk, NJ
	"078":{40.8311, -75.0052}, // Buttzville, NJ
	"078":{40.8540, -74.8257}, // Hackettstown, NJ
	"078":{40.7912, -74.9138}, // Port Murray, NJ
	"078":{40.9190, -74.6390}, // Mount Arlington, NJ
	"078":{40.9142, -74.7027}, // Stanhope, NJ
	"078":{40.7356, -75.0495}, // Broadway, NJ
	"078":{40.9269, -74.9930}, // Mount Hermon, NJ
	"078":{40.7825, -74.7777}, // Long Valley, NJ
	"078":{40.8747, -74.6274}, // Kenvil, NJ
	"078":{40.8184, -75.0613}, // Brookfield, NJ
	"078":{41.1275, -74.7135}, // Ross Corner, NJ
	"078":{40.8614, -74.9858}, // Mountain Lake, NJ
	"078":{40.8371, -75.0244}, // Bridgeville, NJ
	"078":{40.7193, -74.8372}, // Califon, NJ
	"078":{40.7587, -74.9825}, // Washington, NJ
	"078":{40.9650, -74.8782}, // Johnsonburg, NJ
	"078":{40.9859, -74.7429}, // Andover, NJ
	"078":{40.9260, -75.0945}, // Columbia, NJ
	"078":{40.9399, -74.7171}, // Byram Center, NJ
	"078":{41.1471, -74.7495}, // Branchville, NJ
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <079> with city WEST JERSEY and state NJ
	"079":{40.7966, -74.4772}, // Morristown, NJ
	"079":{40.8357, -74.4786}, // Morris Plains, NJ
	"079":{40.7268, -74.5918}, // Bernardsville, NJ
	"079":{40.7586, -74.4170}, // Madison, NJ
	"079":{40.7687, -74.6003}, // Mendham, NJ
	"079":{40.7874, -74.6921}, // Chester, NJ
	"079":{40.6907, -74.6258}, // Far Hills, NJ
	"079":{40.7154, -74.3647}, // Summit, NJ
	"079":{40.7168, -74.6562}, // Peapack and Gladstone, NJ
	"079":{40.7773, -74.3953}, // Florham Park, NJ
	"079":{40.6996, -74.4034}, // New Providence, NJ
	"079":{40.7405, -74.3838}, // Chatham, NJ
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <080> with city SOUTH JERSEY and state NJ
	"080":{39.9199, -75.0093}, // Ellisburg, NJ
	"080":{39.5974, -74.3306}, // Tuckerton, NJ
	"080":{39.6503, -75.3251}, // Woodstown, NJ
	"080":{39.7915, -74.9375}, // Berlin, NJ
	"080":{39.7458, -75.3117}, // Swedesboro, NJ
	"080":{39.7530, -74.1081}, // Barnegat Light, NJ
	"080":{39.8791, -75.0645}, // Haddon Heights, NJ
	"080":{39.7275, -75.4691}, // Penns Grove, NJ
	"080":{39.5731, -74.7148}, // Elwood, NJ
	"080":{39.8454, -75.0218}, // Somerdale, NJ
	"080":{39.9016, -74.9297}, // Marlton, NJ
	"080":{39.7565, -75.3564}, // Beckett, NJ
	"080":{39.6458, -74.1832}, // Ship Bottom, NJ
	"080":{39.7307, -74.8786}, // Chesilhurst, NJ
	"080":{39.9715, -74.6855}, // Pemberton, NJ
	"080":{39.6572, -74.7678}, // Hammonton, NJ
	"080":{39.7879, -74.9857}, // Pine Hill, NJ
	"080":{39.9488, -74.5411}, // Country Lake Estates, NJ
	"080":{39.6009, -74.2119}, // North Beach Haven, NJ
	"080":{39.8332, -74.9657}, // Gibbsboro, NJ
	"080":{39.8924, -75.1172}, // Gloucester City, NJ
	"080":{40.0783, -74.8524}, // Burlington, NJ
	"080":{39.8665, -75.0941}, // Bellmawr, NJ
	"080":{39.6702, -74.2337}, // Beach Haven West, NJ
	"080":{39.9150, -74.5648}, // Presidential Lakes Estates, NJ
	"080":{39.9565, -74.6767}, // Pemberton Heights, NJ
	"080":{39.8213, -75.0053}, // Laurel Springs, NJ
	"080":{39.8689, -75.0514}, // Barrington, NJ
	"080":{39.8806, -75.0918}, // Mount Ephraim, NJ
	"080":{39.7611, -75.4076}, // Pedricktown, NJ
	"080":{39.8582, -74.8054}, // Medford Lakes, NJ
	"080":{39.8379, -75.1524}, // Woodbury, NJ
	"080":{39.6931, -74.2493}, // Manahawkin, NJ
	"080":{39.7913, -75.1486}, // Wenonah, NJ
	"080":{39.8367, -75.0220}, // Hi-Nella, NJ
	"080":{39.6627, -75.0782}, // Clayton, NJ
	"080":{39.9013, -75.0009}, // Barclay, NJ
	"080":{39.5659, -74.3831}, // Mystic Island, NJ
	"080":{39.7335, -75.1306}, // Pitman, NJ
	"080":{40.0652, -74.9221}, // Beverly, NJ
	"080":{39.5681, -75.4724}, // Salem, NJ
	"080":{39.6874, -74.9786}, // Williamstown, NJ
	"080":{39.8676, -75.1853}, // National Park, NJ
	"080":{39.8562, -75.0365}, // Magnolia, NJ
	"080":{39.7870, -74.9742}, // Pine Valley, NJ
	"080":{39.6304, -74.9661}, // Victory Lakes, NJ
	"080":{39.8955, -75.0346}, // Haddonfield, NJ
	"080":{39.8405, -75.0678}, // Glendora, NJ
	"080":{39.9384, -75.0117}, // Cherry Hill Mall, NJ
	"080":{39.8989, -74.9614}, // Greentree, NJ
	"080":{39.8788, -75.1207}, // Brooklawn, NJ
	"080":{39.9288, -75.0398}, // Golden Triangle, NJ
	"080":{39.8704, -75.1301}, // Westville, NJ
	"080":{39.7666, -75.0614}, // Turnersville, NJ
	"080":{39.5939, -74.8821}, // Collings Lakes, NJ
	"080":{39.8046, -74.9851}, // Clementon, NJ
	"080":{39.8482, -74.9957}, // Echelon, NJ
	"080":{39.7982, -75.0629}, // Blackwood, NJ
	"080":{39.6996, -74.1421}, // Harvey Cedars, NJ
	"080":{39.7430, -74.2804}, // Ocean Acres, NJ
	"080":{39.5046, -75.4611}, // Hancocks Bridge, NJ
	"080":{39.8758, -75.0272}, // Tavistock, NJ
	"080":{39.5928, -74.8424}, // Folsom, NJ
	"080":{40.0115, -75.0148}, // Riverton, NJ
	"080":{39.8233, -75.2782}, // Gibbstown, NJ
	"080":{39.8673, -75.0289}, // Lawnside, NJ
	"080":{39.8782, -75.0085}, // Ashland, NJ
	"080":{39.7014, -75.1113}, // Glassboro, NJ
	"080":{39.8769, -74.9723}, // Springdale, NJ
	"080":{39.6645, -74.1709}, // Surf City, NJ
	"080":{39.8055, -75.1589}, // Oak Valley, NJ
	"080":{39.7266, -75.2191}, // Mullica Hill, NJ
	"080":{39.9188, -74.9898}, // Kingston Estates, NJ
	"080":{39.8987, -74.7052}, // Leisuretowne, NJ
	"080":{39.8160, -75.1512}, // Woodbury Heights, NJ
	"080":{39.9737, -74.5690}, // Browns Mills, NJ
	"080":{39.7140, -75.1734}, // Richwood, NJ
	"080":{39.9322, -74.9527}, // Ramblewood, NJ
	"080":{40.0121, -74.6715}, // Juliustown, NJ
	"080":{39.8521, -75.0739}, // Runnemede, NJ
	"080":{39.5658, -74.2489}, // Beach Haven, NJ
	"080":{39.9659, -74.9643}, // Moorestown-Lenola, NJ
	"080":{39.8400, -75.2397}, // Paulsboro, NJ
	"080":{39.8290, -75.0156}, // Stratford, NJ
	"080":{39.8172, -74.9898}, // Lindenwold, NJ
	"080":{40.0025, -75.0360}, // Palmyra, NJ
	// END OF OPTIONS TO PICK
	"081":{39.9362, -75.1073}, // Camden, NJ
	TODO PLACEHOLDER TO FIX ZIP <082> with city SOUTH JERSEY and state NJ
	"082":{38.9765, -74.9516}, // North Cape May, NJ
	"082":{39.5640, -74.5961}, // Egg Harbor City, NJ
	"082":{39.3718, -74.5543}, // Northfield, NJ
	"082":{39.0489, -74.8464}, // Burleigh, NJ
	"082":{39.5731, -74.7148}, // Elwood, NJ
	"082":{38.9587, -74.8520}, // Diamond Beach, NJ
	"082":{39.0047, -74.7990}, // North Wildwood, NJ
	"082":{38.9423, -74.9377}, // West Cape May, NJ
	"082":{39.3047, -74.7199}, // Corbin City, NJ
	"082":{39.1944, -74.6618}, // Strathmere, NJ
	"082":{39.0437, -74.7687}, // Stone Harbor, NJ
	"082":{39.4228, -74.4944}, // Absecon, NJ
	"082":{38.9409, -74.9042}, // Cape May, NJ
	"082":{39.4934, -74.4782}, // Smithville, NJ
	"082":{39.2283, -74.8096}, // Woodbine, NJ
	"082":{39.0417, -74.8673}, // Whitesboro, NJ
	"082":{38.9718, -74.8376}, // Wildwood Crest, NJ
	"082":{38.9877, -74.8188}, // Wildwood, NJ
	"082":{39.0790, -74.8209}, // Cape May Court House, NJ
	"082":{39.2682, -74.6019}, // Ocean City, NJ
	"082":{39.3587, -74.7752}, // Estell Manor, NJ
	"082":{38.9971, -74.8932}, // Erma, NJ
	"082":{39.0906, -74.7357}, // Avalon, NJ
	"082":{39.2662, -74.8672}, // Belleplain, NJ
	"082":{39.5328, -74.4852}, // Port Republic, NJ
	"082":{39.0157, -74.9350}, // Villas, NJ
	"082":{39.3167, -74.6066}, // Somers Point, NJ
	"082":{39.4687, -74.5501}, // Pomona, NJ
	"082":{39.1523, -74.6976}, // Sea Isle City, NJ
	"082":{39.0005, -74.8234}, // West Wildwood, NJ
	"082":{39.0196, -74.8762}, // Rio Grande, NJ
	"082":{38.9372, -74.9651}, // Cape May Point, NJ
	"082":{39.3435, -74.5708}, // Linwood, NJ
	"082":{39.3900, -74.5169}, // Pleasantville, NJ
	"082":{39.4138, -74.3787}, // Brigantine, NJ
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <083> with city SOUTH JERSEY and state NJ
	"083":{39.3265, -75.0306}, // Laurel Lake, NJ
	"083":{39.5015, -75.2190}, // Seabrook Farms, NJ
	"083":{39.3903, -75.0561}, // Millville, NJ
	"083":{39.3758, -75.2130}, // Fairton, NJ
	"083":{39.6627, -75.0782}, // Clayton, NJ
	"083":{39.5401, -75.1727}, // Olivet, NJ
	"083":{39.4653, -74.9981}, // Vineland, NJ
	"083":{39.5484, -75.0167}, // Newfield, NJ
	"083":{39.4286, -75.2281}, // Bridgeton, NJ
	"083":{39.4787, -75.1380}, // Rosenhayn, NJ
	"083":{39.3587, -74.7752}, // Estell Manor, NJ
	"083":{39.5046, -75.4611}, // Hancocks Bridge, NJ
	"083":{39.3382, -75.2069}, // Cedarville, NJ
	"083":{39.4524, -74.7241}, // Mays Landing, NJ
	"083":{39.5282, -74.9448}, // Buena, NJ
	"083":{39.2521, -75.0412}, // Port Norris, NJ
	"083":{39.5919, -75.1740}, // Elmer, NJ
	"083":{39.4595, -75.2970}, // Shiloh, NJ
	// END OF OPTIONS TO PICK
	"084":{39.3797, -74.4527}, // Atlantic City, NJ
	"085":{40.2236, -74.7641}, // Trenton, NJ
	"086":{40.2236, -74.7641}, // Trenton, NJ
	TODO PLACEHOLDER TO FIX ZIP <087> with city MONMOUTH and state NJ
	"087":{40.0384, -74.1678}, // Leisure Village East, NJ
	"087":{39.9452, -74.0787}, // Seaside Heights, NJ
	"087":{39.7897, -74.1925}, // Waretown, NJ
	"087":{39.9268, -74.1350}, // Ocean Gate, NJ
	"087":{40.0104, -74.2806}, // Leisure Village West, NJ
	"087":{40.0928, -74.0457}, // Point Pleasant Beach, NJ
	"087":{39.9538, -74.0798}, // Dover Beaches South, NJ
	"087":{39.9527, -74.3995}, // Cedar Glen Lakes, NJ
	"087":{39.9423, -74.1454}, // Island Heights, NJ
	"087":{39.9621, -74.2359}, // Silver Ridge, NJ
	"087":{39.9568, -74.3524}, // Crestwood Village, NJ
	"087":{40.0186, -74.2908}, // Leisure Knoll, NJ
	"087":{39.9418, -74.2087}, // South Toms River, NJ
	"087":{40.0132, -74.3201}, // Lakehurst, NJ
	"087":{39.9921, -74.0713}, // Dover Beaches North, NJ
	"087":{39.9287, -74.2023}, // Beachwood, NJ
	"087":{40.1048, -74.0636}, // Brielle, NJ
	"087":{40.1384, -74.1032}, // Allenwood, NJ
	"087":{39.9585, -74.3169}, // Pine Ridge at Crestwood, NJ
	"087":{40.0445, -74.1852}, // Leisure Village, NJ
	"087":{39.9396, -74.2572}, // Holiday Heights, NJ
	"087":{40.1307, -74.0354}, // Sea Girt, NJ
	"087":{39.9697, -74.0718}, // Lavallette, NJ
	"087":{40.0772, -74.0702}, // Point Pleasant, NJ
	"087":{40.0434, -74.0512}, // Mantoloking, NJ
	"087":{40.0418, -74.2852}, // Cedar Glen West, NJ
	"087":{40.0017, -74.2595}, // Pine Lake Park, NJ
	"087":{39.9259, -74.0782}, // Seaside Park, NJ
	"087":{40.0700, -74.0479}, // Bay Head, NJ
	"087":{39.9639, -74.2787}, // Holiday City-Berkeley, NJ
	"087":{39.8202, -74.1456}, // Forked River, NJ
	"087":{39.9533, -74.2365}, // Holiday City South, NJ
	"087":{40.1183, -74.0446}, // Manasquan, NJ
	"087":{39.9359, -74.1700}, // Pine Beach, NJ
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <088> with city KILMER and state NJ
	"088":{40.6774, -75.1558}, // Upper Pohatcong, NJ
	"088":{40.5361, -74.5745}, // Zarephath, NJ
	"088":{40.4489, -74.2495}, // Laurence Harbor, NJ
	"088":{40.6159, -74.7720}, // White House Station, NJ
	"088":{40.6048, -74.6150}, // Green Knoll, NJ
	"088":{40.3949, -74.3919}, // Spotswood, NJ
	"088":{40.5415, -74.3124}, // Fords, NJ
	"088":{40.7184, -75.0774}, // New Village, NJ
	"088":{40.5535, -74.5277}, // South Bound Brook, NJ
	"088":{40.6537, -75.0840}, // Bloomsbury, NJ
	"088":{40.5903, -74.4656}, // Dunellen, NJ
	"088":{40.6684, -74.8940}, // High Bridge, NJ
	"088":{40.4439, -74.5432}, // Franklin Park, NJ
	"088":{40.3810, -74.5435}, // Monmouth Junction, NJ
	"088":{40.5344, -74.4577}, // Society Hill, NJ
	"088":{40.5697, -74.6092}, // Somerville, NJ
	"088":{40.4656, -74.3237}, // Sayreville, NJ
	"088":{40.7818, -75.1252}, // Hutchinson, NJ
	"088":{40.5203, -74.2724}, // Perth Amboy, NJ
	"088":{40.6112, -75.1706}, // Finesville, NJ
	"088":{40.5694, -75.0916}, // Milford, NJ
	"088":{40.6360, -74.9124}, // Clinton, NJ
	"088":{40.5321, -74.5415}, // Franklin Center, NJ
	"088":{40.7028, -75.1194}, // Upper Stewartsville, NJ
	"088":{40.5087, -74.8599}, // Flemington, NJ
	"088":{40.4461, -74.2959}, // Madison Park, NJ
	"088":{40.4138, -74.5626}, // Kendall Park, NJ
	"088":{40.5421, -74.5892}, // Manville, NJ
	"088":{40.6895, -75.1821}, // Phillipsburg, NJ
	"088":{40.7003, -75.0121}, // Asbury, NJ
	"088":{40.7004, -74.9395}, // Glen Gardner, NJ
	"088":{40.3096, -74.4644}, // Clearbrook Park, NJ
	"088":{40.3494, -74.4400}, // Jamesburg, NJ
	"088":{40.6433, -74.8351}, // Lebanon, NJ
	"088":{40.5711, -74.6678}, // Bradley Gardens, NJ
	"088":{40.5083, -74.5010}, // Somerset, NJ
	"088":{40.5744, -74.5011}, // Middlesex, NJ
	"088":{40.4218, -74.5890}, // Ten Mile Run, NJ
	"088":{40.4504, -74.4350}, // Milltown, NJ
	"088":{40.4933, -74.4711}, // East Franklin, NJ
	"088":{40.4852, -74.2831}, // South Amboy, NJ
	"088":{40.3297, -74.4451}, // Whittingham, NJ
	"088":{40.5732, -74.6431}, // Raritan, NJ
	"088":{40.5424, -74.3628}, // Metuchen, NJ
	"088":{40.7356, -75.0495}, // Broadway, NJ
	"088":{40.6972, -75.1452}, // Lopatcong Overlook, NJ
	"088":{40.4523, -74.5696}, // Pleasant Plains, NJ
	"088":{40.5015, -74.5349}, // Middlebush, NJ
	"088":{40.4870, -74.4450}, // New Brunswick, NJ
	"088":{40.3815, -74.5137}, // Dayton, NJ
	"088":{40.6598, -75.1571}, // Alpha, NJ
	"088":{40.5273, -75.0571}, // Frenchtown, NJ
	"088":{40.4009, -74.2959}, // Brownville, NJ
	"088":{40.6468, -74.8881}, // Annandale, NJ
	"088":{40.5002, -74.5920}, // Millstone, NJ
	"088":{40.3777, -74.4239}, // Helmetta, NJ
	"088":{40.5676, -74.5383}, // Bound Brook, NJ
	"088":{40.4814, -74.5714}, // Blackwells Mills, NJ
	"088":{40.5011, -74.5661}, // East Millstone, NJ
	"088":{40.4455, -74.3783}, // South River, NJ
	"088":{40.7069, -74.9642}, // Hampton, NJ
	"088":{40.7718, -75.1699}, // Brainards, NJ
	"088":{40.4730, -74.5334}, // Six Mile Run, NJ
	"088":{40.7051, -75.1883}, // Delaware Park, NJ
	"088":{40.3117, -74.4477}, // Concordia, NJ
	"088":{40.6030, -74.5751}, // Martinsville, NJ
	"088":{40.5702, -74.3170}, // Iselin, NJ
	"088":{40.6939, -75.1115}, // Stewartsville, NJ
	"088":{40.4874, -74.5132}, // Clyde, NJ
	"088":{40.3361, -74.4726}, // Rossmoor, NJ
	"088":{40.5224, -74.5765}, // Weston, NJ
	"088":{40.5626, -74.5743}, // Finderne, NJ
	"088":{40.3908, -74.5756}, // Heathcote, NJ
	// END OF OPTIONS TO PICK
	"089":{40.4870, -74.4450}, // New Brunswick, NJ
	TODO PLACEHOLDER TO FIX ZIP <090> with city APO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <091> with city APO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <092> with city APO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <093> with city APO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <094> with city APO/FPO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <095> with city FPO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <096> with city APO/FPO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <097> with city APO/FPO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <098> with city APO/FPO and state AE
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <099> with city APO/FPO and state AE
	// END OF OPTIONS TO PICK
	"100":{40.6943, -73.9249}, // New York, NY
	"101":{40.6943, -73.9249}, // New York, NY
	"102":{40.6943, -73.9249}, // New York, NY
	"103":{40.5834, -74.1496}, // Staten Island, NY
	"104":{40.8501, -73.8662}, // Bronx, NY
	TODO PLACEHOLDER TO FIX ZIP <105> with city WESTCHESTER and state NY
	"105":{41.4228, -73.6069}, // Brewster Hill, NY
	"105":{40.9469, -73.7316}, // Mamaroneck, NY
	"105":{41.2362, -73.6959}, // Bedford Hills, NY
	"105":{41.0541, -73.8143}, // Elmsford, NY
	"105":{41.2883, -73.9227}, // Peekskill, NY
	"105":{41.3969, -73.6150}, // Brewster, NY
	"105":{41.1601, -73.7672}, // Chappaqua, NY
	"105":{41.1889, -73.5561}, // Scotts Corners, NY
	"105":{41.4150, -73.6855}, // Carmel Hamlet, NY
	"105":{41.2697, -73.7755}, // Yorktown Heights, NY
	"105":{41.3165, -73.8475}, // Lake Mohegan, NY
	"105":{40.9902, -73.7773}, // Scarsdale, NY
	"105":{41.4612, -73.6681}, // Lake Carmel, NY
	"105":{41.0135, -73.8395}, // Ardsley, NY
	"105":{41.3682, -73.5778}, // Peach Lake, NY
	"105":{41.1609, -73.8712}, // Ossining, NY
	"105":{41.0775, -73.7780}, // Valhalla, NY
	"105":{41.0127, -73.8698}, // Dobbs Ferry, NY
	"105":{41.2005, -73.9002}, // Croton-on-Hudson, NY
	"105":{41.2908, -73.8357}, // Crompond, NY
	"105":{41.2018, -73.7282}, // Mount Kisco, NY
	"105":{41.2558, -73.9585}, // Verplanck, NY
	"105":{41.0303, -73.6865}, // Rye Brook, NY
	"105":{41.0233, -73.7192}, // Harrison, NY
	"105":{41.3398, -73.7016}, // Heritage Hills, NY
	"105":{41.0936, -73.8724}, // Sleepy Hollow, NY
	"105":{41.1035, -73.7968}, // Hawthorne, NY
	"105":{41.1320, -73.7137}, // Armonk, NY
	"105":{41.0349, -73.8661}, // Irvington, NY
	"105":{41.2455, -73.9376}, // Montrose, NY
	"105":{41.0052, -73.6680}, // Port Chester, NY
	"105":{41.3256, -73.8295}, // Shrub Oak, NY
	"105":{41.4747, -73.5483}, // Putnam Lake, NY
	"105":{41.2643, -73.9465}, // Buchanan, NY
	"105":{41.1400, -73.8440}, // Briarcliff Manor, NY
	"105":{41.2878, -73.6681}, // Golden's Bridge, NY
	"105":{40.9982, -73.8194}, // Greenville, NY
	"105":{41.2279, -73.9260}, // Crugers, NY
	"105":{41.3306, -73.7409}, // Shenorock, NY
	"105":{40.9136, -73.8291}, // Mount Vernon, NY
	"105":{41.1186, -73.7795}, // Thornwood, NY
	"105":{41.4291, -73.9466}, // Nelsonville, NY
	"105":{40.9690, -73.6878}, // Rye, NY
	"105":{40.9258, -73.7529}, // Larchmont, NY
	"105":{41.2559, -73.6856}, // Katonah, NY
	"105":{41.4191, -73.9545}, // Cold Spring, NY
	"105":{41.3179, -73.8007}, // Jefferson Valley-Yorktown, NY
	"105":{41.1378, -73.7827}, // Pleasantville, NY
	"105":{41.3684, -73.7401}, // Mahopac, NY
	"105":{41.0153, -73.8036}, // Hartsdale, NY
	"105":{41.0647, -73.8673}, // Tarrytown, NY
	"105":{41.3361, -73.7254}, // Lincolndale, NY
	// END OF OPTIONS TO PICK
	"106":{41.0220, -73.7549}, // White Plains, NY
	"107":{40.9466, -73.8674}, // Yonkers, NY
	"108":{40.9305, -73.7836}, // New Rochelle, NY
	TODO PLACEHOLDER TO FIX ZIP <109> with city WESTCHESTER and state NY
	"109":{41.1213, -74.0685}, // Kaser, NY
	"109":{41.1578, -74.0769}, // Wesley Hills, NY
	"109":{41.4759, -74.3682}, // Scotchtown, NY
	"109":{41.1934, -73.9520}, // Haverstraw, NY
	"109":{41.1143, -73.9057}, // Upper Nyack, NY
	"109":{41.3401, -73.9853}, // Fort Montgomery, NY
	"109":{41.1317, -74.1134}, // Montebello, NY
	"109":{41.1264, -74.1705}, // Hillburn, NY
	"109":{41.1484, -73.9456}, // Congers, NY
	"109":{41.4458, -74.4228}, // Middletown, NY
	"109":{41.3088, -74.1444}, // Harriman, NY
	"109":{41.3646, -74.0135}, // West Point, NY
	"109":{41.4692, -74.4179}, // Washington Heights, NY
	"109":{41.2063, -73.9883}, // West Haverstraw, NY
	"109":{41.0289, -73.9333}, // Sparkill, NY
	"109":{41.1151, -74.0486}, // Spring Valley, NY
	"109":{41.4714, -74.5397}, // Otisville, NY
	"109":{41.1488, -74.0485}, // New Hempstead, NY
	"109":{41.3644, -73.9683}, // Highland Falls, NY
	"109":{41.0829, -74.0551}, // Chestnut Ridge, NY
	"109":{41.2551, -74.3550}, // Warwick, NY
	"109":{41.1620, -74.1902}, // Sloatsburg, NY
	"109":{41.1287, -74.0855}, // Viola, NY
	"109":{41.0992, -74.0990}, // Airmont, NY
	"109":{41.3284, -74.1004}, // Woodbury, NY
	"109":{41.3119, -74.2249}, // Walton Park, NY
	"109":{41.1410, -74.0294}, // New Square, NY
	"109":{41.0798, -73.9127}, // South Nyack, NY
	"109":{41.3570, -74.2769}, // Chester, NY
	"109":{41.0615, -74.0047}, // Pearl River, NY
	"109":{41.1160, -73.9436}, // Valley Cottage, NY
	"109":{41.4296, -74.1578}, // Washingtonville, NY
	"109":{41.0907, -73.9714}, // West Nyack, NY
	"109":{41.3198, -74.1848}, // Monroe, NY
	"109":{41.2215, -74.2891}, // Greenwood Lake, NY
	"109":{41.1926, -74.0297}, // Mount Ivy, NY
	"109":{41.0423, -73.9150}, // Piermont, NY
	"109":{41.1138, -74.1421}, // Suffern, NY
	"109":{41.1543, -73.9909}, // New City, NY
	"109":{41.4472, -74.3914}, // Mechanicstown, NY
	"109":{41.0269, -73.9520}, // Tappan, NY
	"109":{41.1181, -74.0681}, // Monsey, NY
	"109":{41.1892, -74.0543}, // Pomona, NY
	"109":{41.4017, -74.3270}, // Goshen, NY
	"109":{41.3735, -74.1790}, // South Blooming Grove, NY
	"109":{41.3312, -74.3533}, // Florida, NY
	"109":{41.1298, -74.0350}, // Hillcrest, NY
	"109":{41.0919, -73.9143}, // Nyack, NY
	"109":{41.0689, -73.9545}, // Blauvelt, NY
	"109":{41.0488, -73.9407}, // Orangeburg, NY
	"109":{41.3013, -74.5619}, // Unionville, NY
	"109":{41.0627, -73.9208}, // Grand View-on-Hudson, NY
	"109":{41.2007, -74.2060}, // Tuxedo Park, NY
	"109":{41.3404, -74.1658}, // Kiryas Joel, NY
	"109":{41.0957, -74.0155}, // Nanuet, NY
	"109":{41.2067, -74.0122}, // Thiells, NY
	"109":{41.3862, -74.1389}, // Mountain Lodge Park, NY
	"109":{41.1129, -73.9823}, // Bardonia, NY
	// END OF OPTIONS TO PICK
	"110":{40.7498, -73.7976}, // Queens, NY
	TODO PLACEHOLDER TO FIX ZIP <111> with city LONG ISLAND CITY and state NY
	"111":{40.7498, -73.7976}, // Queens, NY
	"111":{40.6943, -73.9249}, // New York, NY
	// END OF OPTIONS TO PICK
	"112":{40.6501, -73.9496}, // Brooklyn, NY
	TODO PLACEHOLDER TO FIX ZIP <113> with city FLUSHING and state NY
	"113":{40.7498, -73.7976}, // Queens, NY
	"113":{40.6943, -73.9249}, // New York, NY
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <114> with city JAMAICA and state NY
	"114":{40.7498, -73.7976}, // Queens, NY
	"114":{40.6943, -73.9249}, // New York, NY
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <115> with city WESTERN NASSAU and state NY
	"115":{40.7990, -73.6491}, // Roslyn, NY
	"115":{40.6218, -73.7507}, // Inwood, NY
	"115":{40.7043, -73.6193}, // Hempstead, NY
	"115":{40.7765, -73.6778}, // North Hills, NY
	"115":{40.7500, -73.6122}, // Carle Place, NY
	"115":{40.8123, -73.5696}, // Brookville, NY
	"115":{40.6422, -73.6942}, // Hewlett, NY
	"115":{40.7866, -73.5975}, // Old Westbury, NY
	"115":{40.6871, -73.5615}, // North Merrick, NY
	"115":{40.6840, -73.7077}, // North Valley Stream, NY
	"115":{40.6959, -73.6507}, // West Hempstead, NY
	"115":{40.7121, -73.6605}, // Garden City South, NY
	"115":{40.6579, -73.6742}, // Lynbrook, NY
	"115":{40.7454, -73.5604}, // Salisbury, NY
	"115":{40.7787, -73.6396}, // Roslyn Heights, NY
	"115":{40.6302, -73.6670}, // Bay Park, NY
	"115":{40.8782, -73.5884}, // Locust Valley, NY
	"115":{40.8709, -73.6287}, // Glen Cove, NY
	"115":{40.7608, -73.6336}, // East Williston, NY
	"115":{40.5887, -73.6660}, // Long Beach, NY
	"115":{40.6247, -73.6981}, // Hewlett Neck, NY
	"115":{40.6296, -73.6025}, // Baldwin Harbor, NY
	"115":{40.6050, -73.6437}, // Barnum Island, NY
	"115":{40.8075, -73.6755}, // Flower Hill, NY
	"115":{40.8118, -73.6262}, // Greenvale, NY
	"115":{40.6647, -73.7044}, // Valley Stream, NY
	"115":{40.7937, -73.6611}, // Roslyn Estates, NY
	"115":{40.6252, -73.7278}, // Cedarhurst, NY
	"115":{40.6328, -73.6363}, // Oceanside, NY
	"115":{40.6374, -73.7219}, // Woodmere, NY
	"115":{40.6814, -73.6233}, // South Hempstead, NY
	"115":{40.7176, -73.5947}, // Uniondale, NY
	"115":{40.6634, -73.6104}, // Baldwin, NY
	"115":{40.7197, -73.5604}, // East Meadow, NY
	"115":{40.5903, -73.5795}, // Point Lookout, NY
	"115":{40.8254, -73.5363}, // Muttontown, NY
	"115":{40.6685, -73.6737}, // North Lynbrook, NY
	"115":{40.6775, -73.6493}, // Lakeview, NY
	"115":{40.7203, -73.6853}, // Stewart Manor, NY
	"115":{40.5904, -73.6121}, // Lido Beach, NY
	"115":{40.7469, -73.6392}, // Mineola, NY
	"115":{40.5894, -73.7296}, // Atlantic Beach, NY
	"115":{40.6346, -73.6953}, // Hewlett Bay Park, NY
	"115":{40.8640, -73.5817}, // Matinecock, NY
	"115":{40.6019, -73.6647}, // Harbor Isle, NY
	"115":{40.7567, -73.6635}, // Herricks, NY
	"115":{40.6643, -73.6383}, // Rockville Centre, NY
	"115":{40.6746, -73.6721}, // Malverne, NY
	"115":{40.6042, -73.7149}, // Lawrence, NY
	"115":{40.7587, -73.6465}, // Williston Park, NY
	"115":{40.7600, -73.5649}, // New Cassel, NY
	"115":{40.6327, -73.6842}, // Hewlett Harbor, NY
	"115":{40.6051, -73.6553}, // Island Park, NY
	"115":{40.6515, -73.5535}, // Merrick, NY
	"115":{40.5876, -73.7092}, // East Atlantic Beach, NY
	"115":{40.6557, -73.7186}, // South Valley Stream, NY
	"115":{40.8332, -73.6039}, // Old Brookville, NY
	"115":{40.6817, -73.6642}, // Malverne Park Oaks, NY
	"115":{40.6797, -73.5837}, // Roosevelt, NY
	"115":{40.7599, -73.5891}, // Westbury, NY
	"115":{40.8157, -73.6378}, // Roslyn Harbor, NY
	"115":{40.8450, -73.6180}, // Glen Head, NY
	"115":{40.7705, -73.6603}, // Searingtown, NY
	"115":{40.6432, -73.6672}, // East Rockaway, NY
	"115":{40.8921, -73.5966}, // Lattingtown, NY
	"115":{40.6515, -73.5850}, // Freeport, NY
	"115":{40.7266, -73.6447}, // Garden City, NY
	"115":{40.8295, -73.6378}, // Glenwood Landing, NY
	"115":{40.7958, -73.6292}, // East Hills, NY
	"115":{40.8476, -73.5627}, // Upper Brookville, NY
	"115":{40.8441, -73.6442}, // Sea Cliff, NY
	"115":{40.7715, -73.6482}, // Albertson, NY
	"115":{40.6215, -73.7068}, // Woodsburgh, NY
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <116> with city FAR ROCKAWAY and state NY
	"116":{40.7498, -73.7976}, // Queens, NY
	"116":{40.6943, -73.9249}, // New York, NY
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <117> with city MID-ISLAND and state NY
	"117":{40.8608, -73.4488}, // Cold Spring Harbor, NY
	"117":{40.9090, -73.0492}, // Terryville, NY
	"117":{40.7193, -73.2642}, // Brightwaters, NY
	"117":{40.6388, -73.1949}, // Saltaire, NY
	"117":{40.8225, -73.3921}, // South Huntington, NY
	"117":{40.7621, -73.0185}, // Patchogue, NY
	"117":{40.9571, -72.9072}, // Shoreham, NY
	"117":{40.7623, -73.3219}, // Deer Park, NY
	"117":{40.8696, -73.0808}, // Centereach, NY
	"117":{40.8123, -73.0447}, // Holtsville, NY
	"117":{40.9060, -73.1278}, // Stony Brook, NY
	"117":{40.6377, -73.3788}, // Gilgo, NY
	"117":{40.9260, -73.0651}, // Port Jefferson Station, NY
	"117":{40.7495, -73.4856}, // Bethpage, NY
	"117":{40.8389, -73.0401}, // Farmingville, NY
	"117":{40.7506, -73.1872}, // Islip Terrace, NY
	"117":{40.6858, -73.3709}, // Lindenhurst, NY
	"117":{40.7175, -73.4471}, // South Farmingdale, NY
	"117":{40.6569, -73.5285}, // Bellmore, NY
	"117":{40.7311, -73.3251}, // North Babylon, NY
	"117":{40.9060, -73.2992}, // Fort Salonga, NY
	"117":{40.7241, -73.5125}, // Levittown, NY
	"117":{40.9380, -73.3815}, // Asharoken, NY
	"117":{40.8761, -73.1521}, // St. James, NY
	"117":{40.9099, -73.1213}, // Stony Brook University, NY
	"117":{40.7373, -73.1345}, // Oakdale, NY
	"117":{40.9464, -72.8230}, // Wading River, NY
	"117":{40.9307, -73.1018}, // Setauket-East Setauket, NY
	"117":{40.8813, -73.0059}, // Coram, NY
	"117":{40.7080, -73.2719}, // West Bay Shore, NY
	"117":{40.8467, -73.1522}, // Nesconset, NY
	"117":{40.9328, -73.3951}, // Eatons Neck, NY
	"117":{40.7518, -73.0352}, // Blue Point, NY
	"117":{40.7599, -73.1678}, // North Great River, NY
	"117":{40.6728, -73.3932}, // Copiague, NY
	"117":{40.7164, -73.1603}, // Great River, NY
	"117":{40.9614, -73.1324}, // Old Field, NY
	"117":{40.8035, -73.3370}, // Dix Hills, NY
	"117":{40.8792, -73.3232}, // East Northport, NY
	"117":{40.6678, -73.4922}, // Seaford, NY
	"117":{40.8839, -73.5582}, // Mill Neck, NY
	"117":{40.8541, -73.4755}, // Laurel Hollow, NY
	"117":{40.9460, -72.8812}, // East Shoreham, NY
	"117":{40.7005, -73.4118}, // North Amityville, NY
	"117":{40.8217, -73.2119}, // Hauppauge, NY
	"117":{40.7240, -73.4770}, // Plainedge, NY
	"117":{40.9607, -73.0672}, // Belle Terre, NY
	"117":{40.7031, -73.4679}, // North Massapequa, NY
	"117":{40.8586, -73.1168}, // Lake Grove, NY
	"117":{40.6983, -73.5086}, // North Wantagh, NY
	"117":{40.8443, -73.2834}, // Commack, NY
	"117":{40.7875, -73.5416}, // Jericho, NY
	"117":{40.7072, -73.3859}, // North Lindenhurst, NY
	"117":{40.6904, -73.5390}, // North Bellmore, NY
	"117":{40.7601, -73.2618}, // North Bay Shore, NY
	"117":{40.7557, -73.4544}, // Old Bethpage, NY
	"117":{40.8840, -73.5002}, // Cove Neck, NY
	"117":{40.6817, -73.4496}, // Massapequa Park, NY
	"117":{40.8308, -73.1112}, // Lake Ronkonkoma, NY
	"117":{40.7328, -73.4465}, // Farmingdale, NY
	"117":{40.8254, -73.5363}, // Muttontown, NY
	"117":{40.6686, -73.5104}, // Wantagh, NY
	"117":{40.6676, -73.4706}, // Massapequa, NY
	"117":{40.8462, -73.3389}, // Elwood, NY
	"117":{40.6949, -73.3270}, // Babylon, NY
	"117":{40.7097, -73.2971}, // West Islip, NY
	"117":{40.9578, -72.9726}, // Sound Beach, NY
	"117":{40.7547, -72.9423}, // Bellport, NY
	"117":{40.7624, -73.3705}, // Wheatley Heights, NY
	"117":{40.8630, -73.3642}, // Greenlawn, NY
	"117":{40.8699, -73.0462}, // Selden, NY
	"117":{40.6782, -73.0709}, // Fire Island, NY
	"117":{40.9014, -73.4163}, // Huntington Bay, NY
	"117":{40.8569, -73.5038}, // Oyster Bay Cove, NY
	"117":{40.8220, -72.9859}, // Medford, NY
	"117":{40.9022, -73.1922}, // Nissequogue, NY
	"117":{40.8156, -73.5020}, // Syosset, NY
	"117":{40.7837, -73.1945}, // Central Islip, NY
	"117":{40.7944, -73.0707}, // Holbrook, NY
	"117":{40.8864, -73.4139}, // Halesite, NY
	"117":{40.8198, -73.4339}, // West Hills, NY
	"117":{40.7294, -73.1050}, // West Sayville, NY
	"117":{40.9372, -73.0180}, // Mount Sinai, NY
	"117":{40.9018, -73.5211}, // Centre Island, NY
	"117":{40.6465, -73.2721}, // Oak Beach-Captree, NY
	"117":{40.8068, -73.1711}, // Islandia, NY
	"117":{40.9529, -73.0903}, // Poquott, NY
	"117":{40.7823, -73.4088}, // Melville, NY
	"117":{40.7704, -72.9817}, // East Patchogue, NY
	"117":{40.8446, -73.4050}, // Huntington Station, NY
	"117":{40.8524, -73.1844}, // Village of the Branch, NY
	"117":{40.7275, -73.1861}, // East Islip, NY
	"117":{40.6743, -73.4358}, // East Massapequa, NY
	"117":{40.9036, -73.3446}, // Northport, NY
	"117":{40.6696, -73.4156}, // Amityville, NY
	"117":{40.8040, -73.1258}, // Ronkonkoma, NY
	"117":{40.7317, -73.2505}, // Bay Shore, NY
	"117":{40.7833, -73.0234}, // North Patchogue, NY
	"117":{40.7336, -73.4169}, // East Farmingdale, NY
	"117":{40.9139, -73.4618}, // Lloyd Harbor, NY
	"117":{40.7868, -72.9457}, // North Bellport, NY
	"117":{40.6463, -73.1565}, // Ocean Beach, NY
	"117":{40.9357, -72.9364}, // Rocky Point, NY
	"117":{40.8981, -73.1624}, // Head of the Harbor, NY
	"117":{40.8887, -73.2452}, // Kings Park, NY
	"117":{40.7478, -73.0840}, // Sayville, NY
	"117":{40.7533, -73.2900}, // Baywood, NY
	"117":{40.7717, -73.1271}, // Bohemia, NY
	"117":{40.7112, -73.3567}, // West Babylon, NY
	"117":{40.7839, -73.2522}, // Brentwood, NY
	"117":{40.9077, -73.5603}, // Bayville, NY
	"117":{40.8476, -73.5627}, // Upper Brookville, NY
	"117":{40.9374, -72.9864}, // Miller Place, NY
	"117":{40.8496, -73.5288}, // East Norwich, NY
	"117":{40.9068, -72.8816}, // Ridge, NY
	"117":{40.9465, -73.0579}, // Port Jefferson, NY
	"117":{40.8943, -73.3714}, // Centerport, NY
	"117":{40.8645, -72.9678}, // Gordon Heights, NY
	"117":{40.7467, -73.3769}, // Wyandanch, NY
	"117":{40.7460, -73.0546}, // Bayport, NY
	// END OF OPTIONS TO PICK
	"118":{40.7637, -73.5245}, // Hicksville, NY
	TODO PLACEHOLDER TO FIX ZIP <119> with city MID-ISLAND and state NY
	"119":{40.9330, -72.4047}, // North Sea, NY
	"119":{40.8070, -72.8235}, // Moriches, NY
	"119":{41.1022, -72.3759}, // Greenport West, NY
	"119":{40.9844, -72.1326}, // Amagansett, NY
	"119":{40.9302, -72.2726}, // Sagaponack, NY
	"119":{40.8169, -72.7046}, // Remsenburg-Speonk, NY
	"119":{40.8489, -72.5783}, // East Quogue, NY
	"119":{40.9222, -72.3532}, // Water Mill, NY
	"119":{41.0931, -72.3417}, // Dering Harbor, NY
	"119":{41.0471, -71.9449}, // Montauk, NY
	"119":{40.7679, -72.8375}, // Mastic Beach, NY
	"119":{40.9427, -72.3101}, // Bridgehampton, NY
	"119":{40.9546, -72.5807}, // Jamesport, NY
	"119":{41.1425, -72.2770}, // Orient, NY
	"119":{41.0230, -72.3140}, // North Haven, NY
	"119":{40.8323, -72.9233}, // Yaphank, NY
	"119":{40.9464, -72.8230}, // Wading River, NY
	"119":{40.8015, -72.7960}, // Center Moriches, NY
	"119":{40.9725, -72.1889}, // East Hampton North, NY
	"119":{41.1031, -72.3669}, // Greenport, NY
	"119":{40.9961, -72.4767}, // New Suffolk, NY
	"119":{40.8407, -72.7251}, // Eastport, NY
	"119":{41.0167, -72.4872}, // Cutchogue, NY
	"119":{40.8857, -72.9454}, // Middle Island, NY
	"119":{40.7949, -72.8743}, // Shirley, NY
	"119":{41.1291, -72.3420}, // East Marion, NY
	"119":{40.9163, -72.7645}, // Calverton, NY
	"119":{41.0381, -72.4608}, // Peconic, NY
	"119":{40.8877, -72.4554}, // Shinnecock Hills, NY
	"119":{40.6782, -73.0709}, // Fire Island, NY
	"119":{40.9591, -72.2507}, // Wainscott, NY
	"119":{40.9950, -72.0707}, // Napeague, NY
	"119":{40.8778, -72.4004}, // Southampton, NY
	"119":{40.8692, -72.5227}, // Hampton Bays, NY
	"119":{40.9728, -72.5560}, // Laurel, NY
	"119":{40.8575, -72.7915}, // Manorville, NY
	"119":{40.8925, -72.6049}, // Flanders, NY
	"119":{40.9970, -72.2892}, // Sag Harbor, NY
	"119":{40.8215, -72.5987}, // Quogue, NY
	"119":{41.0012, -72.5419}, // Mattituck, NY
	"119":{41.0745, -72.3434}, // Shelter Island Heights, NY
	"119":{41.0212, -72.1584}, // Springs, NY
	"119":{40.9527, -72.1961}, // East Hampton, NY
	"119":{40.9827, -72.3350}, // Noyack, NY
	"119":{40.9357, -72.9364}, // Rocky Point, NY
	"119":{41.0053, -72.2220}, // Northwest Harbor, NY
	"119":{40.7776, -72.7137}, // West Hampton Dunes, NY
	"119":{40.8097, -72.7581}, // East Moriches, NY
	"119":{40.8080, -72.6457}, // Westhampton Beach, NY
	"119":{40.9645, -72.7402}, // Baiting Hollow, NY
	"119":{40.8201, -72.6280}, // Quiogue, NY
	"119":{40.8098, -72.8479}, // Mastic, NY
	"119":{40.9425, -72.6149}, // Aquebogue, NY
	"119":{40.8325, -72.6617}, // Westhampton, NY
	"119":{40.9068, -72.8816}, // Ridge, NY
	"119":{40.8645, -72.9678}, // Gordon Heights, NY
	// END OF OPTIONS TO PICK
	"120":{42.6664, -73.7987}, // Albany, NY
	"121":{42.6664, -73.7987}, // Albany, NY
	"122":{42.6664, -73.7987}, // Albany, NY
	"123":{42.8025, -73.9276}, // Schenectady, NY
	TODO PLACEHOLDER TO FIX ZIP <124> with city MID-HUDSON and state NY
	"124":{42.1773, -74.0229}, // Palenville, NY
	"124":{41.6655, -74.3912}, // Cragsmoor, NY
	"124":{41.9559, -74.0016}, // Lincoln Park, NY
	"124":{41.7972, -74.2319}, // Accord, NY
	"124":{41.7009, -74.3609}, // Ellenville, NY
	"124":{41.9048, -73.9776}, // Port Ewen, NY
	"124":{41.8504, -74.0738}, // Rosendale Hamlet, NY
	"124":{42.0209, -74.0855}, // Zena, NY
	"124":{42.0083, -74.1121}, // West Hurley, NY
	"124":{42.0750, -73.9484}, // Saugerties, NY
	"124":{42.1552, -74.5335}, // Fleischmanns, NY
	"124":{41.7798, -74.2956}, // Kerhonkson, NY
	"124":{41.8316, -74.0696}, // Tillson, NY
	"124":{42.2361, -73.8822}, // Jefferson Heights, NY
	"124":{42.1307, -74.4665}, // Pine Hill, NY
	"124":{42.2109, -74.2160}, // Hunter, NY
	"124":{41.8414, -74.1539}, // Stone Ridge, NY
	"124":{41.9177, -74.0336}, // Hillside, NY
	"124":{41.9523, -73.9704}, // East Kingston, NY
	"124":{42.0933, -73.9365}, // Malden-on-Hudson, NY
	"124":{41.9809, -74.2130}, // Shokan, NY
	"124":{41.9863, -73.9895}, // Lake Katrine, NY
	"124":{42.0790, -74.3090}, // Phoenicia, NY
	"124":{41.8288, -74.0381}, // Rifton, NY
	"124":{42.2145, -73.8656}, // Catskill, NY
	"124":{42.1937, -74.1353}, // Tannersville, NY
	"124":{42.0614, -73.9506}, // Saugerties South, NY
	"124":{42.2529, -73.8940}, // Leeds, NY
	"124":{42.0461, -73.9487}, // Glasco, NY
	"124":{42.1458, -74.6503}, // Margaretville, NY
	"124":{42.4330, -74.2287}, // Preston-Potter Hollow, NY
	"124":{41.7514, -74.3729}, // Napanoch, NY
	"124":{41.9295, -73.9968}, // Kingston, NY
	"124":{41.8275, -74.1184}, // High Falls, NY
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <125> with city MID-HUDSON and state NY
	"125":{41.5037, -74.0205}, // Newburgh, NY
	"125":{41.6191, -73.7944}, // Hillside Lake, NY
	"125":{41.4308, -74.1098}, // Salisbury Mills, NY
	"125":{41.6094, -74.2966}, // Pine Bush, NY
	"125":{41.9223, -73.9441}, // Rhinecliff, NY
	"125":{41.6655, -74.3912}, // Cragsmoor, NY
	"125":{41.5281, -74.0233}, // Balmville, NY
	"125":{41.5984, -73.9181}, // Wappingers Falls, NY
	"125":{41.9294, -73.9081}, // Rhinebeck, NY
	"125":{41.6371, -74.2633}, // Watchtower, NY
	"125":{41.5947, -73.8744}, // Myers Corner, NY
	"125":{41.5216, -74.2388}, // Montgomery, NY
	"125":{41.4441, -74.1178}, // Beaver Dam Lake, NY
	"125":{41.9533, -73.5112}, // Millerton, NY
	"125":{41.5635, -73.5989}, // Pawling, NY
	"125":{42.0587, -73.9119}, // Tivoli, NY
	"125":{41.7681, -73.9007}, // Haviland, NY
	"125":{41.6927, -74.0458}, // Clintondale, NY
	"125":{41.7395, -73.5734}, // Dover Plains, NY
	"125":{41.4588, -74.0534}, // Vails Gate, NY
	"125":{42.2263, -73.7269}, // Claverack-Red Mills, NY
	"125":{41.5511, -73.8697}, // Brinckerhoff, NY
	"125":{41.7495, -74.0809}, // New Paltz, NY
	"125":{41.4409, -74.0353}, // Firthcliffe, NY
	"125":{41.8554, -73.9255}, // Staatsburg, NY
	"125":{42.2669, -73.7700}, // Lorenz Park, NY
	"125":{41.4881, -74.2131}, // Maybrook, NY
	"125":{41.4369, -74.0145}, // Cornwall-on-Hudson, NY
	"125":{42.2477, -73.6463}, // Philmont, NY
	"125":{42.1041, -73.5494}, // Copake Hamlet, NY
	"125":{41.5385, -73.8724}, // Merritt Park, NY
	"125":{41.6695, -73.7974}, // Freedom Plains, NY
	"125":{41.6028, -73.9774}, // Marlboro, NY
	"125":{41.5036, -73.9655}, // Beacon, NY
	"125":{41.5319, -74.0936}, // Orange Lake, NY
	"125":{41.8061, -73.7900}, // Salt Point, NY
	"125":{41.4747, -73.5483}, // Putnam Lake, NY
	"125":{41.9959, -73.8769}, // Red Hook, NY
	"125":{42.1202, -73.5247}, // Copake Falls, NY
	"125":{41.8288, -74.0381}, // Rifton, NY
	"125":{42.2515, -73.7859}, // Hudson, NY
	"125":{41.5787, -73.8078}, // Hopewell Junction, NY
	"125":{42.2913, -73.7532}, // Stottville, NY
	"125":{42.1192, -73.5532}, // Taconic Shores, NY
	"125":{41.5328, -74.0594}, // Gardnertown, NY
	"125":{41.5603, -74.1879}, // Walden, NY
	"125":{41.6386, -74.3775}, // Walker Valley, NY
	"125":{41.7179, -73.9646}, // Highland, NY
	"125":{42.1416, -73.5903}, // Copake Lake, NY
	"125":{41.7842, -73.6937}, // Millbrook, NY
	"125":{41.5337, -73.8942}, // Fishkill, NY
	// END OF OPTIONS TO PICK
	"126":{41.6949, -73.9210}, // Poughkeepsie, NY
	TODO PLACEHOLDER TO FIX ZIP <127> with city MID-HUDSON and state NY
	"127":{41.3782, -74.6909}, // Port Jervis, NY
	"127":{41.7962, -74.7429}, // Liberty, NY
	"127":{41.5761, -74.4856}, // Wurtsboro, NY
	"127":{41.5519, -74.4437}, // Bloomingburg, NY
	"127":{41.7218, -74.6350}, // South Fallsburg, NY
	"127":{41.9403, -74.9132}, // Roscoe, NY
	"127":{41.7733, -74.6558}, // Loch Sheldrake, NY
	"127":{41.7661, -75.0221}, // Hortonville, NY
	"127":{41.7799, -74.9298}, // Jeffersonville, NY
	"127":{41.6006, -75.0579}, // Narrowsburg, NY
	"127":{41.6523, -74.6876}, // Monticello, NY
	"127":{41.6153, -74.5821}, // Rock Hill, NY
	"127":{41.7125, -74.5742}, // Woodridge, NY
	"127":{41.6595, -74.8200}, // Smallwood, NY
	"127":{41.8935, -74.8265}, // Livingston Manor, NY
	// END OF OPTIONS TO PICK
	"128":{43.3109, -73.6459}, // Glens Falls, NY
	"129":{44.6951, -73.4563}, // Plattsburgh, NY
	"130":{43.0409, -76.1438}, // Syracuse, NY
	"131":{43.0409, -76.1438}, // Syracuse, NY
	"132":{43.0409, -76.1438}, // Syracuse, NY
	"133":{43.0961, -75.2260}, // Utica, NY
	"134":{43.0961, -75.2260}, // Utica, NY
	"135":{43.0961, -75.2260}, // Utica, NY
	"136":{43.9734, -75.9095}, // Watertown, NY
	"137":{42.1014, -75.9093}, // Binghamton, NY
	"138":{42.1014, -75.9093}, // Binghamton, NY
	"139":{42.1014, -75.9093}, // Binghamton, NY
	"140":{42.9017, -78.8487}, // Buffalo, NY
	"141":{42.9017, -78.8487}, // Buffalo, NY
	"142":{42.9017, -78.8487}, // Buffalo, NY
	"143":{43.0921, -79.0147}, // Niagara Falls, NY
	"144":{43.1680, -77.6162}, // Rochester, NY
	"145":{43.1680, -77.6162}, // Rochester, NY
	"146":{43.1680, -77.6162}, // Rochester, NY
	"147":{42.0975, -79.2366}, // Jamestown, NY
	"148":{42.0938, -76.8097}, // Elmira, NY
	"149":{42.0938, -76.8097}, // Elmira, NY
	"150":{40.4396, -79.9763}, // Pittsburgh, PA
	"151":{40.4396, -79.9763}, // Pittsburgh, PA
	"152":{40.4396, -79.9763}, // Pittsburgh, PA
	"153":{40.4396, -79.9763}, // Pittsburgh, PA
	"154":{40.4396, -79.9763}, // Pittsburgh, PA
	"155":{40.3258, -78.9194}, // Johnstown, PA
	"156":{40.3113, -79.5444}, // Greensburg, PA
	"157":{40.3258, -78.9194}, // Johnstown, PA
	TODO PLACEHOLDER TO FIX ZIP <158> with city DU BOIS and state PA
	"158":{41.1615, -79.0827}, // Brookville, PA
	"158":{41.2469, -78.7929}, // Brockway, PA
	"158":{41.3570, -78.6101}, // Kersey, PA
	"158":{41.2764, -78.4902}, // Weedville, PA
	"158":{41.3448, -78.1340}, // Driftwood, PA
	"158":{41.0946, -78.8880}, // Reynoldsville, PA
	"158":{41.1163, -79.1880}, // Summerville, PA
	"158":{41.2487, -78.7542}, // Crenshaw, PA
	"158":{41.0259, -78.7862}, // Troutville, PA
	"158":{41.5102, -78.2363}, // Emporium, PA
	"158":{41.2581, -78.5035}, // Force, PA
	"158":{41.4269, -78.7297}, // Ridgway, PA
	"158":{41.2916, -78.5023}, // Byrnedale, PA
	"158":{41.1062, -78.7749}, // Sandy, PA
	"158":{41.5727, -78.6869}, // Wilcox, PA
	"158":{41.1421, -78.8067}, // Falls Creek, PA
	"158":{41.4913, -78.6791}, // Johnsonburg, PA
	"158":{41.0475, -78.8185}, // Sykesville, PA
	"158":{41.1713, -78.7173}, // Treasure Lake, PA
	"158":{41.1225, -78.7564}, // DuBois, PA
	"158":{41.1817, -79.2027}, // Corsica, PA
	"158":{41.4574, -78.5343}, // St. Marys, PA
	// END OF OPTIONS TO PICK
	"159":{40.3258, -78.9194}, // Johnstown, PA
	"160":{40.9956, -80.3458}, // New Castle, PA
	"161":{40.9956, -80.3458}, // New Castle, PA
	"162":{40.9956, -80.3458}, // New Castle, PA
	"163":{41.4282, -79.7035}, // Oil City, PA
	"164":{42.1168, -80.0733}, // Erie, PA
	"165":{42.1168, -80.0733}, // Erie, PA
	"166":{40.5082, -78.4007}, // Altoona, PA
	"167":{41.9604, -78.6413}, // Bradford, PA
	"168":{40.5082, -78.4007}, // Altoona, PA
	"169":{41.2398, -77.0371}, // Williamsport, PA
	"170":{40.2752, -76.8843}, // Harrisburg, PA
	"171":{40.2752, -76.8843}, // Harrisburg, PA
	"172":{40.2752, -76.8843}, // Harrisburg, PA
	"173":{40.0420, -76.3012}, // Lancaster, PA
	"174":{39.9651, -76.7315}, // York, PA
	"175":{40.0420, -76.3012}, // Lancaster, PA
	"176":{40.0420, -76.3012}, // Lancaster, PA
	"177":{41.2398, -77.0371}, // Williamsport, PA
	"178":{40.2752, -76.8843}, // Harrisburg, PA
	"179":{40.3400, -75.9267}, // Reading, PA
	TODO PLACEHOLDER TO FIX ZIP <180> with city LEHIGH VALLEY and state PA
	"180":{40.5165, -75.5545}, // Macungie, PA
	"180":{40.7544, -75.6114}, // Slatington, PA
	"180":{40.5965, -75.1987}, // Riegelsville, PA
	"180":{40.4476, -75.5508}, // Hereford, PA
	"180":{40.5360, -75.5852}, // Ancient Oaks, PA
	"180":{40.8811, -75.1861}, // East Bangor, PA
	"180":{40.8675, -75.2535}, // Pen Argyl, PA
	"180":{40.3936, -75.4965}, // Pennsburg, PA
	"180":{40.6280, -75.3401}, // Freemansburg, PA
	"180":{40.5102, -75.3915}, // Coopersburg, PA
	"180":{40.4057, -75.5059}, // East Greenville, PA
	"180":{40.5523, -75.6039}, // Trexlertown, PA
	"180":{40.7816, -75.1912}, // Martins Creek, PA
	"180":{40.6890, -75.5161}, // Cementon, PA
	"180":{40.7409, -75.2561}, // Tatamy, PA
	"180":{40.7232, -75.5366}, // Laurys Station, PA
	"180":{40.8350, -75.2226}, // Ackermanville, PA
	"180":{40.5393, -75.6342}, // Breinigsville, PA
	"180":{40.6266, -75.3679}, // Bethlehem, PA
	"180":{40.5811, -75.3378}, // Hellertown, PA
	"180":{40.5617, -75.5489}, // Wescosville, PA
	"180":{40.6693, -75.6116}, // Schnecksville, PA
	"180":{40.6607, -75.2364}, // Glendon, PA
	"180":{40.7549, -75.5309}, // Cherryville, PA
	"180":{40.8484, -75.2917}, // Wind Gap, PA
	"180":{40.7439, -75.6583}, // Slatedale, PA
	"180":{40.6710, -75.4961}, // Coplay, PA
	"180":{40.7833, -75.2755}, // Belfast, PA
	"180":{40.6575, -75.2609}, // Old Orchard, PA
	"180":{40.6858, -75.2209}, // Easton, PA
	"180":{40.6866, -75.4904}, // Northampton, PA
	"180":{40.8243, -75.6701}, // Parryville, PA
	"180":{40.7547, -75.2627}, // Stockertown, PA
	"180":{40.6029, -75.3961}, // Fountain Hill, PA
	"180":{40.5961, -75.4755}, // Allentown, PA
	"180":{40.8678, -75.2085}, // Bangor, PA
	"180":{40.6585, -75.4952}, // Hokendauqua, PA
	"180":{40.8023, -75.6160}, // Palmerton, PA
	"180":{40.6308, -75.4834}, // Fullerton, PA
	"180":{40.5352, -75.4978}, // Emmaus, PA
	"180":{40.6773, -75.7498}, // New Tripoli, PA
	"180":{40.5090, -75.6001}, // Alburtis, PA
	"180":{40.8778, -75.2204}, // Roseto, PA
	"180":{40.6907, -75.2671}, // Palmer Heights, PA
	"180":{40.5388, -75.3778}, // DeSales University, PA
	"180":{40.8019, -75.6611}, // Bowmanstown, PA
	"180":{40.7279, -75.3919}, // Bath, PA
	"180":{40.7400, -75.3132}, // Nazareth, PA
	"180":{40.3771, -75.4839}, // Red Hill, PA
	"180":{40.6640, -75.4741}, // North Catasauqua, PA
	"180":{40.6666, -75.5072}, // Stiles, PA
	"180":{40.6844, -75.2407}, // Wilson, PA
	"180":{40.6781, -75.2357}, // West Easton, PA
	"180":{40.6295, -75.2024}, // Raubsville, PA
	"180":{40.7514, -75.5955}, // Walnutport, PA
	"180":{40.6858, -75.5334}, // Egypt, PA
	"180":{40.6531, -75.4643}, // Catasauqua, PA
	"180":{40.7585, -75.4013}, // Chapman, PA
	"180":{40.3364, -75.4709}, // Green Lane, PA
	"180":{40.7480, -75.2913}, // Eastlawn Gardens, PA
	// END OF OPTIONS TO PICK
	"181":{40.5961, -75.4755}, // Allentown, PA
	TODO PLACEHOLDER TO FIX ZIP <182> with city WILKES BARRE and state PA
	"182":{41.0212, -75.8963}, // Freeland, PA
	"182":{40.9504, -75.9724}, // Hazleton, PA
	"182":{40.9160, -75.9657}, // Tresckow, PA
	"182":{40.9946, -75.5828}, // Towamensing Trails, PA
	"182":{41.0347, -75.8221}, // Hickory Hills, PA
	"182":{40.8216, -75.9865}, // Hometown, PA
	"182":{40.9426, -76.1400}, // Weston, PA
	"182":{40.9862, -75.9702}, // Harleigh, PA
	"182":{40.8249, -75.8464}, // Summit Hill, PA
	"182":{40.9009, -75.9924}, // McAdoo, PA
	"182":{40.8306, -75.7166}, // Lehighton, PA
	"182":{41.0007, -75.9680}, // Pardeesville, PA
	"182":{40.8289, -75.7009}, // Weissport, PA
	"182":{40.8659, -75.8322}, // Nesquehoning, PA
	"182":{41.0117, -75.6068}, // Albrightsville, PA
	"182":{40.8033, -75.9344}, // Tamaqua, PA
	"182":{40.8192, -76.0601}, // Park Crest, PA
	"182":{40.9903, -75.8961}, // Jeddo, PA
	"182":{40.9378, -76.1685}, // Nuremberg, PA
	"182":{40.8243, -75.6701}, // Parryville, PA
	"182":{40.9930, -75.9606}, // Lattimer, PA
	"182":{40.8197, -75.9161}, // Coaldale, PA
	"182":{40.9007, -76.0045}, // Kelayres, PA
	"182":{40.8331, -75.8847}, // Lansford, PA
	"182":{40.8960, -76.1187}, // Sheppton, PA
	"182":{41.0415, -75.9323}, // Beech Mountain Lakes, PA
	"182":{40.8369, -75.6863}, // Weissport East, PA
	"182":{40.8266, -76.0560}, // Grier City, PA
	"182":{41.0292, -75.6092}, // Holiday Pocono, PA
	"182":{41.0003, -75.5058}, // Indian Mountain Lake, PA
	"182":{40.9299, -75.9131}, // Beaver Meadows, PA
	"182":{40.9420, -75.8210}, // Weatherly, PA
	"182":{40.9701, -76.0132}, // West Hazleton, PA
	"182":{40.8414, -76.0738}, // Delano, PA
	"182":{40.8712, -75.7433}, // Jim Thorpe, PA
	"182":{40.9052, -76.1224}, // Oneida, PA
	"182":{40.9911, -76.0595}, // Conyngham, PA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <183> with city LEHIGH VALLEY and state PA
	"183":{40.8998, -75.3179}, // Saylorsburg, PA
	"183":{41.3666, -74.6997}, // Matamoras, PA
	"183":{41.1149, -75.4589}, // Pocono Pines, PA
	"183":{41.1779, -75.2632}, // Mountainhome, PA
	"183":{41.3225, -74.8902}, // Pocono Woodland Lakes, PA
	"183":{41.1225, -75.3582}, // Mount Pocono, PA
	"183":{41.0023, -75.1779}, // East Stroudsburg, PA
	"183":{40.9800, -75.4688}, // Sun Valley, PA
	"183":{41.0039, -75.2116}, // Arlington Heights, PA
	"183":{41.3040, -74.9937}, // Conashaugh Lakes, PA
	"183":{40.9436, -75.4390}, // Effort, PA
	"183":{41.1355, -74.9923}, // Pine Ridge, PA
	"183":{41.3145, -74.9434}, // Gold Key Lake, PA
	"183":{40.9215, -75.0984}, // Portland, PA
	"183":{41.1567, -74.9646}, // Pocono Mountain Lake Estates, PA
	"183":{41.1197, -75.0465}, // Saw Creek, PA
	"183":{41.0346, -75.2417}, // Penn Estates, PA
	"183":{41.1779, -74.9584}, // Pocono Ranch Lands, PA
	"183":{40.9750, -75.1392}, // Delaware Water Gap, PA
	"183":{41.0844, -75.4154}, // Emerald Lakes, PA
	"183":{41.3238, -74.8011}, // Milford, PA
	"183":{41.0003, -75.5058}, // Indian Mountain Lake, PA
	"183":{41.0007, -75.4476}, // Sierra View, PA
	"183":{41.2513, -74.9103}, // Birchwood Lakes, PA
	"183":{40.9838, -75.1972}, // Stroudsburg, PA
	"183":{41.3136, -74.9638}, // Sunrise Lake, PA
	"183":{40.9267, -75.4017}, // Brodheadsville, PA
	// END OF OPTIONS TO PICK
	"184":{41.4044, -75.6649}, // Scranton, PA
	"185":{41.4044, -75.6649}, // Scranton, PA
	TODO PLACEHOLDER TO FIX ZIP <186> with city WILKES BARRE and state PA
	"186":{41.6119, -76.0458}, // Meshoppen, PA
	"186":{41.0555, -76.2492}, // Berwick, PA
	"186":{41.0775, -76.2355}, // Foundryville, PA
	"186":{41.2968, -75.8167}, // Inkerman, PA
	"186":{41.3274, -75.7885}, // Pittston, PA
	"186":{41.1535, -76.1500}, // Shickshinny, PA
	"186":{41.4070, -75.8503}, // Upper Exeter, PA
	"186":{41.3381, -75.7423}, // Avoca, PA
	"186":{41.0347, -75.8221}, // Hickory Hills, PA
	"186":{41.5412, -75.9488}, // Tunkhannock, PA
	"186":{41.3037, -75.7825}, // Yatesville, PA
	"186":{41.2202, -76.0095}, // West Nanticoke, PA
	"186":{41.3237, -75.7421}, // Dupont, PA
	"186":{41.3537, -75.7758}, // Duryea, PA
	"186":{41.2404, -75.9505}, // Plymouth, PA
	"186":{41.0566, -75.7801}, // White Haven, PA
	"186":{41.3058, -75.8416}, // Wyoming, PA
	"186":{41.3295, -75.7998}, // West Pittston, PA
	"186":{41.1406, -76.1368}, // Mocanaqua, PA
	"186":{41.0307, -76.2995}, // Mifflinville, PA
	"186":{41.4180, -76.4920}, // Laporte, PA
	"186":{41.2639, -75.9326}, // Larksville, PA
	"186":{41.1763, -76.0384}, // Wanamie, PA
	"186":{41.0519, -76.2115}, // Nescopeck, PA
	"186":{41.4896, -76.6028}, // Forksville, PA
	"186":{41.2698, -76.0749}, // Silkworth, PA
	"186":{41.3138, -75.7826}, // Browntown, PA
	"186":{41.0654, -76.2208}, // East Berwick, PA
	"186":{41.1153, -75.7733}, // Penn Lake Park, PA
	"186":{41.3217, -75.8578}, // West Wyoming, PA
	"186":{41.3625, -76.0384}, // Harveys Lake, PA
	"186":{41.1938, -76.0188}, // Sheatown, PA
	"186":{41.3306, -75.9736}, // Dallas, PA
	"186":{41.5249, -76.3988}, // Dushore, PA
	"186":{41.3338, -75.8214}, // Exeter, PA
	"186":{41.0621, -75.7595}, // East Side, PA
	"186":{41.2004, -76.0003}, // Nanticoke, PA
	"186":{41.4530, -75.8605}, // West Falls, PA
	"186":{41.0007, -75.4476}, // Sierra View, PA
	"186":{41.4225, -76.0599}, // Noxen, PA
	"186":{41.3291, -75.7698}, // Hughestown, PA
	"186":{41.1877, -75.7552}, // Bear Creek Village, PA
	"186":{41.1796, -76.0782}, // Glen Lyon, PA
	"186":{41.5105, -75.8509}, // Lake Winola, PA
	"186":{41.3086, -76.1029}, // Pikes Creek, PA
	"186":{41.3584, -75.7027}, // Moosic, PA
	"186":{41.0463, -76.2861}, // Briar Creek, PA
	"186":{41.6458, -76.1589}, // Laceyville, PA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <187> with city WILKES BARRE and state PA
	"187":{41.2867, -75.8967}, // Luzerne, PA
	"187":{41.2896, -75.7942}, // Laflin, PA
	"187":{41.2848, -75.9142}, // Courtdale, PA
	"187":{41.1580, -75.9773}, // Nuangola, PA
	"187":{41.2128, -75.8994}, // Ashley, PA
	"187":{41.2801, -75.9650}, // Chase, PA
	"187":{41.2468, -75.8759}, // Wilkes-Barre, PA
	"187":{41.1353, -75.9045}, // Mountain Top, PA
	"187":{41.2652, -75.8875}, // Kingston, PA
	"187":{41.2975, -75.8799}, // Swoyersville, PA
	"187":{41.3101, -75.9281}, // Trucksville, PA
	"187":{41.2194, -75.8445}, // Laurel Run, PA
	"187":{41.2639, -75.9326}, // Larksville, PA
	"187":{41.1935, -75.9314}, // Sugar Notch, PA
	"187":{41.2773, -75.9014}, // Pringle, PA
	"187":{41.2758, -75.8518}, // Plains, PA
	"187":{41.1876, -75.9505}, // Warrior Run, PA
	"187":{41.2869, -75.8356}, // Hilldale, PA
	"187":{41.2773, -75.8312}, // Hudson, PA
	"187":{41.2614, -75.9071}, // Edwardsville, PA
	"187":{41.3188, -75.9405}, // Shavertown, PA
	"187":{41.2843, -75.8689}, // Forty Fort, PA
	// END OF OPTIONS TO PICK
	"188":{41.4044, -75.6649}, // Scranton, PA
	TODO PLACEHOLDER TO FIX ZIP <189> with city SOUTHEASTERN and state PA
	"189":{40.1994, -74.9986}, // Churchville, PA
	"189":{40.3259, -75.3274}, // Telford, PA
	"189":{40.2690, -75.2140}, // Brittany Farms-The Highlands, PA
	"189":{40.4412, -75.4414}, // Spinnerstown, PA
	"189":{40.2502, -75.2405}, // Montgomeryville, PA
	"189":{40.4325, -75.4056}, // Milford Square, PA
	"189":{40.2290, -74.9324}, // Newtown, PA
	"189":{40.2599, -74.9561}, // Newtown Grant, PA
	"189":{40.3139, -75.1280}, // Doylestown, PA
	"189":{40.4132, -75.3801}, // Trumbauersville, PA
	"189":{40.3863, -75.1429}, // Plumsteadville, PA
	"189":{40.2262, -75.0006}, // Richboro, PA
	"189":{40.3719, -75.2920}, // Perkasie, PA
	"189":{40.3110, -75.3223}, // Souderton, PA
	"189":{40.3473, -75.2718}, // Silverdale, PA
	"189":{40.2894, -75.2096}, // Chalfont, PA
	"189":{40.3600, -75.3083}, // Sellersville, PA
	"189":{40.3616, -74.9574}, // New Hope, PA
	"189":{40.2083, -75.0733}, // Ivyland, PA
	"189":{40.3731, -75.2041}, // Dublin, PA
	"189":{40.4398, -75.3456}, // Quakertown, PA
	"189":{40.1884, -75.0841}, // Warminster Heights, PA
	"189":{40.4726, -75.3212}, // Richlandtown, PA
	"189":{40.2016, -74.9706}, // Village Shires, PA
	"189":{40.3110, -75.4512}, // Woxall, PA
	"189":{40.2981, -75.1807}, // New Britain, PA
	// END OF OPTIONS TO PICK
	"190":{40.0077, -75.1339}, // Philadelphia, PA
	"191":{40.0077, -75.1339}, // Philadelphia, PA
	"192":{40.0077, -75.1339}, // Philadelphia, PA
	TODO PLACEHOLDER TO FIX ZIP <193> with city SOUTHEASTERN and state PA
	"193":{39.9343, -75.5306}, // Cheyney University, PA
	"193":{39.8248, -75.7819}, // Avondale, PA
	"193":{39.9618, -75.8024}, // Modena, PA
	"193":{39.8926, -75.4698}, // Chester Heights, PA
	"193":{39.9849, -75.8199}, // Coatesville, PA
	"193":{40.0496, -75.4271}, // Devon, PA
	"193":{39.9991, -75.7517}, // Thorndale, PA
	"193":{39.8206, -75.8284}, // West Grove, PA
	"193":{39.8313, -75.7565}, // Toughkenamon, PA
	"193":{40.0307, -75.6303}, // Exton, PA
	"193":{39.9473, -75.9754}, // Atglen, PA
	"193":{39.9639, -75.8850}, // Pomeroy, PA
	"193":{40.0420, -75.4912}, // Paoli, PA
	"193":{40.0076, -75.7019}, // Downingtown, PA
	"193":{39.9950, -75.7844}, // Caln, PA
	"193":{40.0937, -75.9111}, // Honey Brook, PA
	"193":{39.9690, -75.8135}, // South Coatesville, PA
	"193":{39.8065, -75.9279}, // Lincoln University, PA
	"193":{39.8905, -75.9241}, // Cochranville, PA
	"193":{39.9601, -75.6058}, // West Chester, PA
	"193":{39.7859, -75.9801}, // Oxford, PA
	"193":{40.0329, -75.5146}, // Malvern, PA
	"193":{39.9702, -75.8555}, // Westwood, PA
	"193":{39.7745, -76.1121}, // Little Britain, PA
	"193":{40.0598, -75.6789}, // Eagleview, PA
	"193":{40.0524, -75.6440}, // Lionville, PA
	"193":{39.9593, -75.9177}, // Parkesburg, PA
	"193":{39.8438, -75.7113}, // Kennett Square, PA
	"193":{40.0396, -75.4439}, // Berwyn, PA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <194> with city SOUTHEASTERN and state PA
	"194":{40.1873, -75.4581}, // Collegeville, PA
	"194":{40.2771, -75.2988}, // Hatfield, PA
	"194":{40.1992, -75.4754}, // Trappe, PA
	"194":{40.0702, -75.3189}, // West Conshohocken, PA
	"194":{40.2644, -75.6093}, // Pottsgrove, PA
	"194":{40.2225, -75.3998}, // Skippack, PA
	"194":{40.2502, -75.2405}, // Montgomeryville, PA
	"194":{40.1862, -75.5382}, // Royersford, PA
	"194":{40.1099, -75.2798}, // Plymouth Meeting, PA
	"194":{40.2111, -75.2744}, // North Wales, PA
	"194":{40.1042, -75.3437}, // Bridgeport, PA
	"194":{40.1604, -75.4090}, // Eagleville, PA
	"194":{40.1304, -75.4280}, // Audubon, PA
	"194":{40.1768, -75.5466}, // Spring City, PA
	"194":{40.2440, -75.3407}, // Kulpsville, PA
	"194":{40.2506, -75.6819}, // Stowe, PA
	"194":{40.2573, -75.4662}, // Schwenksville, PA
	"194":{40.1489, -75.3995}, // Trooper, PA
	"194":{40.1474, -75.2687}, // Blue Bell, PA
	"194":{40.0772, -75.3035}, // Conshohocken, PA
	"194":{40.1899, -75.4348}, // Evansburg, PA
	"194":{40.2791, -75.3872}, // Harleysville, PA
	"194":{40.1358, -75.5201}, // Phoenixville, PA
	"194":{40.1224, -75.3398}, // Norristown, PA
	"194":{40.2795, -75.6402}, // Halfway House, PA
	"194":{40.0962, -75.3821}, // King of Prussia, PA
	"194":{40.2242, -75.6418}, // Kenilworth, PA
	"194":{40.2755, -75.4651}, // Spring Mount, PA
	"194":{40.2367, -75.6599}, // South Pottstown, PA
	"194":{40.1847, -75.2267}, // Spring House, PA
	"194":{40.3110, -75.4512}, // Woxall, PA
	"194":{40.2417, -75.2812}, // Lansdale, PA
	"194":{40.2498, -75.5886}, // Sanatoga, PA
	"194":{40.2507, -75.6444}, // Pottstown, PA
	// END OF OPTIONS TO PICK
	"195":{40.3400, -75.9267}, // Reading, PA
	"196":{40.3400, -75.9267}, // Reading, PA
	"197":{39.7415, -75.5413}, // Wilmington, DE
	"198":{39.7415, -75.5413}, // Wilmington, DE
	"199":{39.7415, -75.5413}, // Wilmington, DE
	"200":{38.9047, -77.0163}, // Washington, DC
	TODO PLACEHOLDER TO FIX ZIP <201> with city DULLES and state VA
	"201":{39.1319, -77.7681}, // Round Hill, VA
	"201":{38.9133, -77.3969}, // Franklin Farm, VA
	"201":{38.8121, -77.6363}, // Haymarket, VA
	"201":{38.7812, -77.4817}, // Loch Lomond, VA
	"201":{38.7250, -77.4472}, // Buckhall, VA
	"201":{38.9802, -77.5323}, // Brambleton, VA
	"201":{38.6540, -77.6370}, // Catlett, VA
	"201":{39.0464, -77.3874}, // Cascades, VA
	"201":{38.8661, -77.8453}, // Marshall, VA
	"201":{38.9981, -77.4971}, // Moorefield Station, VA
	"201":{39.0265, -77.4196}, // Dulles Town Center, VA
	"201":{38.9347, -77.4083}, // Floris, VA
	"201":{39.0470, -77.3524}, // Lowes Island, VA
	"201":{39.1987, -77.7239}, // Hillsboro, VA
	"201":{38.7495, -77.7151}, // New Baltimore, VA
	"201":{39.1058, -77.5544}, // Leesburg, VA
	"201":{38.7878, -77.4961}, // Sudley, VA
	"201":{38.7882, -77.4495}, // Yorkshire, VA
	"201":{39.1349, -77.6642}, // Hamilton, VA
	"201":{38.9447, -77.5306}, // Arcola, VA
	"201":{39.1373, -77.8637}, // Shenandoah Retreat, VA
	"201":{38.7479, -77.4838}, // Manassas, VA
	"201":{38.6940, -77.5757}, // Nokesville, VA
	"201":{38.7551, -77.5750}, // Linton Hall, VA
	"201":{38.9699, -77.3867}, // Herndon, VA
	"201":{39.0168, -77.5167}, // Broadlands, VA
	"201":{39.0309, -77.3762}, // Sugarland Run, VA
	"201":{38.9801, -77.5082}, // Loudoun Valley Estates, VA
	"201":{38.9715, -77.7428}, // Middleburg, VA
	"201":{38.7176, -77.7975}, // Warrenton, VA
	"201":{39.1378, -77.7110}, // Purcellville, VA
	"201":{39.0518, -77.4124}, // Countryside, VA
	"201":{39.2736, -77.6398}, // Lovettsville, VA
	"201":{38.9295, -77.5553}, // Stone Ridge, VA
	"201":{39.0052, -77.4050}, // Sterling, VA
	"201":{39.0650, -77.4965}, // Belmont, VA
	"201":{38.8868, -77.4453}, // Chantilly, VA
	"201":{39.0601, -77.4445}, // University Center, VA
	"201":{38.6355, -77.6661}, // Calverton, VA
	"201":{38.7931, -77.6347}, // Gainesville, VA
	"201":{38.9497, -77.3461}, // Reston, VA
	"201":{38.9845, -77.4174}, // Oak Grove, VA
	"201":{38.9120, -77.5132}, // South Riding, VA
	"201":{38.7802, -77.5204}, // Bull Run, VA
	"201":{38.8622, -77.7743}, // The Plains, VA
	"201":{39.0300, -77.4711}, // Ashburn, VA
	"201":{38.9513, -77.4116}, // McNair, VA
	"201":{38.9955, -77.3693}, // Dranesville, VA
	"201":{38.7718, -77.4450}, // Manassas Park, VA
	"201":{38.6404, -77.4090}, // Independent Hill, VA
	"201":{38.8391, -77.4388}, // Centreville, VA
	"201":{38.9106, -77.6638}, // Bull Run Mountain Estates, VA
	"201":{38.6201, -77.8049}, // Opal, VA
	"201":{39.0846, -77.4839}, // Lansdowne, VA
	"201":{38.7802, -77.3866}, // Clifton, VA
	// END OF OPTIONS TO PICK
	"202":{38.9047, -77.0163}, // Washington, DC
	"203":{38.9047, -77.0163}, // Washington, DC
	"204":{38.9047, -77.0163}, // Washington, DC
	"205":{38.9047, -77.0163}, // Washington, DC
	TODO PLACEHOLDER TO FIX ZIP <206> with city SOUTHERN MD and state MD
	"206":{38.2646, -76.8496}, // Cobb Island, MD
	"206":{39.0394, -76.9211}, // Beltsville, MD
	"206":{38.4728, -76.4895}, // Calvert Beach, MD
	"206":{38.5913, -76.7074}, // Aquasco, MD
	"206":{38.6117, -76.6187}, // Huntingtown, MD
	"206":{38.1155, -76.4771}, // St. George Island, MD
	"206":{38.5498, -76.8429}, // Bryantown, MD
	"206":{38.6963, -76.8846}, // Brandywine, MD
	"206":{38.1483, -76.5201}, // Piney Point, MD
	"206":{38.5104, -77.0197}, // Port Tobacco Village, MD
	"206":{38.2969, -76.4950}, // California, MD
	"206":{38.3039, -76.6396}, // Leonardtown, MD
	"206":{38.1654, -76.5366}, // Tall Timbers, MD
	"206":{38.3574, -76.4147}, // Chesapeake Ranch Estates, MD
	"206":{38.2754, -76.8428}, // Rock Point, MD
	"206":{38.6144, -77.0850}, // Bryans Road, MD
	"206":{38.6719, -76.7428}, // Baden, MD
	"206":{38.6623, -76.8190}, // Cedarville, MD
	"206":{38.5376, -76.7748}, // Hughesville, MD
	"206":{38.5699, -77.0304}, // Pomfret, MD
	"206":{38.5440, -76.5879}, // Prince Frederick, MD
	"206":{38.5116, -76.6797}, // Benedict, MD
	"206":{38.4900, -76.7019}, // Golden Beach, MD
	"206":{38.3325, -76.4295}, // Drum Point, MD
	"206":{38.4568, -76.4737}, // Long Beach, MD
	"206":{38.3372, -76.4611}, // Solomons, MD
	"206":{38.4667, -76.7847}, // Charlotte Hall, MD
	"206":{38.5987, -77.1555}, // Indian Head, MD
	"206":{38.4202, -76.5475}, // Broomes Island, MD
	"206":{38.6176, -77.0077}, // Bensville, MD
	"206":{38.6745, -77.0023}, // Accokeek, MD
	"206":{38.7458, -76.7555}, // Croom, MD
	"206":{38.5664, -76.6872}, // Eagle Harbor, MD
	"206":{38.2543, -76.4415}, // Lexington Park, MD
	"206":{38.5352, -76.9701}, // La Plata, MD
	"206":{38.6085, -76.9195}, // Waldorf, MD
	"206":{38.7672, -76.8266}, // Rosaryville, MD
	"206":{38.4659, -76.4973}, // St. Leonard, MD
	"206":{38.3620, -76.4372}, // Lusby, MD
	"206":{38.4355, -76.7426}, // Mechanicsville, MD
	"206":{38.5987, -77.1372}, // Potomac Heights, MD
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <207> with city SOUTHERN MD and state MD
	"207":{39.1457, -76.7745}, // Jessup, MD
	"207":{39.0603, -76.8456}, // South Laurel, MD
	"207":{38.9423, -76.9645}, // Mount Rainier, MD
	"207":{39.0394, -76.9211}, // Beltsville, MD
	"207":{38.8385, -76.8232}, // Westphalia, MD
	"207":{38.9293, -76.8576}, // Glenarden, MD
	"207":{38.9070, -76.8299}, // Lake Arbor, MD
	"207":{38.9621, -76.8421}, // Lanham, MD
	"207":{38.9016, -76.9153}, // Fairmount Heights, MD
	"207":{38.7339, -77.0069}, // Fort Washington, MD
	"207":{38.9643, -76.9266}, // Riverdale Park, MD
	"207":{38.7620, -76.7857}, // Marlton, MD
	"207":{38.8392, -76.9367}, // Silver Hill, MD
	"207":{38.9502, -76.9333}, // Edmonston, MD
	"207":{39.0017, -76.9649}, // Adelphi, MD
	"207":{38.9833, -76.8040}, // Glenn Dale, MD
	"207":{39.1133, -76.8924}, // West Laurel, MD
	"207":{38.7156, -76.6730}, // Dunkirk, MD
	"207":{38.7887, -76.9733}, // Oxon Hill, MD
	"207":{38.9667, -76.9790}, // Chillum, MD
	"207":{38.8172, -76.7546}, // Upper Marlboro, MD
	"207":{38.8379, -76.5514}, // Galesville, MD
	"207":{38.9953, -76.8885}, // Greenbelt, MD
	"207":{38.9612, -76.9548}, // Hyattsville, MD
	"207":{39.1016, -76.8051}, // Maryland City, MD
	"207":{38.9254, -76.9141}, // Cheverly, MD
	"207":{38.8306, -76.7699}, // Marlboro Village, MD
	"207":{38.8022, -76.8419}, // Melwood, MD
	"207":{39.1516, -76.9163}, // Fulton, MD
	"207":{38.8940, -76.8878}, // Peppermill Village, MD
	"207":{38.9565, -76.7780}, // Fairwood, MD
	"207":{38.7499, -76.9063}, // Clinton, MD
	"207":{38.9654, -76.8773}, // New Carrollton, MD
	"207":{38.9302, -76.9437}, // Colmar Manor, MD
	"207":{38.9424, -76.8946}, // Landover Hills, MD
	"207":{38.7910, -76.5469}, // Deale, MD
	"207":{38.8106, -76.9495}, // Temple Hills, MD
	"207":{38.6719, -76.7428}, // Baden, MD
	"207":{38.8588, -76.8885}, // District Heights, MD
	"207":{38.9241, -76.8875}, // Landover, MD
	"207":{38.8285, -76.5211}, // Shady Side, MD
	"207":{38.8952, -76.9016}, // Seat Pleasant, MD
	"207":{38.9424, -76.9263}, // Bladensburg, MD
	"207":{38.9222, -76.7810}, // Woodmore, MD
	"207":{38.9574, -76.7422}, // Bowie, MD
	"207":{39.0254, -76.9751}, // Hillandale, MD
	"207":{38.8617, -76.7549}, // Brock Hall, MD
	"207":{38.8021, -76.7839}, // Queensland, MD
	"207":{38.8518, -76.8708}, // Forestville, MD
	"207":{38.8800, -76.8289}, // Largo, MD
	"207":{38.9897, -76.9808}, // Langley Park, MD
	"207":{38.9960, -76.9337}, // College Park, MD
	"207":{38.7117, -76.6055}, // Owings, MD
	"207":{38.8053, -76.8744}, // Andrews AFB, MD
	"207":{38.9719, -76.9445}, // University Park, MD
	"207":{38.7359, -76.5878}, // Friendship, MD
	"207":{38.8181, -76.9836}, // Glassmanor, MD
	"207":{38.8052, -76.9198}, // Camp Springs, MD
	"207":{38.7883, -77.0106}, // National Harbor, MD
	"207":{38.9600, -76.9108}, // East Riverdale, MD
	"207":{38.8373, -76.9641}, // Hillcrest Heights, MD
	"207":{38.9385, -76.9492}, // Cottage City, MD
	"207":{38.7080, -76.5347}, // North Beach, MD
	"207":{39.1485, -76.8228}, // Savage, MD
	"207":{38.9929, -76.9131}, // Berwyn Heights, MD
	"207":{38.8105, -76.9995}, // Forest Heights, MD
	"207":{39.0578, -76.9488}, // Calverton, MD
	"207":{38.7601, -76.9642}, // Friendly, MD
	"207":{38.7458, -76.7555}, // Croom, MD
	"207":{38.8766, -76.9074}, // Capitol Heights, MD
	"207":{38.8754, -76.8862}, // Walker Mill, MD
	"207":{38.9383, -76.8423}, // Springdale, MD
	"207":{38.8888, -76.7890}, // Kettering, MD
	"207":{39.0774, -76.9023}, // Konterra, MD
	"207":{38.8265, -76.8896}, // Morningside, MD
	"207":{39.1286, -76.8476}, // North Laurel, MD
	"207":{39.1058, -76.7438}, // Fort Meade, MD
	"207":{38.8237, -76.9485}, // Marlow Heights, MD
	"207":{38.9358, -76.8146}, // Mitchellville, MD
	"207":{38.8709, -76.9234}, // Coral Hills, MD
	"207":{38.7672, -76.8266}, // Rosaryville, MD
	"207":{38.9802, -76.8502}, // Seabrook, MD
	"207":{38.9439, -76.9571}, // Brentwood, MD
	"207":{39.1814, -76.9570}, // Highland, MD
	"207":{38.9450, -76.9510}, // North Brentwood, MD
	"207":{38.8492, -76.9225}, // Suitland, MD
	"207":{38.8374, -76.7144}, // Marlboro Meadows, MD
	"207":{38.6881, -76.5448}, // Chesapeake Beach, MD
	"207":{39.1416, -76.8843}, // Scaggsville, MD
	"207":{39.0950, -76.8622}, // Laurel, MD
	"207":{38.9042, -76.8678}, // Summerfield, MD
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <208> with city SUBURBAN MD and state MD
	"208":{39.0360, -77.0934}, // Garrett Park, MD
	"208":{39.2737, -77.2006}, // Damascus, MD
	"208":{39.0834, -77.1553}, // Rockville, MD
	"208":{38.9743, -77.1635}, // Cabin John, MD
	"208":{38.9546, -77.1292}, // Brookmont, MD
	"208":{38.9943, -77.0737}, // Chevy Chase, MD
	"208":{39.0571, -77.2458}, // Travilah, MD
	"208":{39.1808, -77.0590}, // Brookeville, MD
	"208":{39.0188, -77.0785}, // South Kensington, MD
	"208":{39.1335, -77.1465}, // Redland, MD
	"208":{38.9680, -77.1410}, // Glen Echo, MD
	"208":{38.9698, -77.0793}, // Chevy Chase Village, MD
	"208":{39.0960, -77.3033}, // Darnestown, MD
	"208":{38.9795, -77.0693}, // Martin's Additions, MD
	"208":{39.0393, -77.1191}, // North Bethesda, MD
	"208":{38.9666, -77.0963}, // Somerset, MD
	"208":{39.0804, -76.9527}, // Fairland, MD
	"208":{39.2211, -77.3798}, // Barnesville, MD
	"208":{39.0265, -77.0737}, // Kensington, MD
	"208":{39.0928, -77.0822}, // Aspen Hill, MD
	"208":{38.9866, -77.1188}, // Bethesda, MD
	"208":{39.1466, -77.0715}, // Olney, MD
	"208":{39.1346, -77.2132}, // Gaithersburg, MD
	"208":{38.9633, -77.0900}, // Friendship Heights Village, MD
	"208":{39.1136, -77.1509}, // Derwood, MD
	"208":{39.0392, -77.0723}, // North Kensington, MD
	"208":{39.0192, -77.0809}, // Chevy Chase View, MD
	"208":{39.1423, -77.4102}, // Poolesville, MD
	"208":{39.1515, -77.0065}, // Ashton-Sandy Spring, MD
	"208":{39.0141, -77.1943}, // Potomac, MD
	"208":{39.0021, -77.0743}, // North Chevy Chase, MD
	"208":{39.1405, -77.1745}, // Washington Grove, MD
	"208":{39.1755, -77.2643}, // Germantown, MD
	"208":{39.1783, -77.1957}, // Montgomery Village, MD
	"208":{39.1166, -76.9356}, // Burtonsville, MD
	"208":{38.9793, -77.0742}, // Chevy Chase Section Three, MD
	"208":{38.9840, -77.0740}, // Chevy Chase Section Five, MD
	"208":{39.2094, -77.1418}, // Laytonsville, MD
	"208":{39.2314, -77.2617}, // Clarksburg, MD
	"208":{39.0955, -77.2372}, // North Potomac, MD
	"208":{39.1190, -76.9828}, // Spencerville, MD
	// END OF OPTIONS TO PICK
	"209":{39.0028, -77.0207}, // Silver Spring, MD
	"210":{39.2088, -76.6625}, // Linthicum, MD
	"211":{39.2088, -76.6625}, // Linthicum, MD
	"212":{39.3051, -76.6144}, // Baltimore, MD
	"214":{38.9706, -76.5047}, // Annapolis, MD
	"215":{39.6515, -78.7585}, // Cumberland, MD
	TODO PLACEHOLDER TO FIX ZIP <216> with city EASTERN SHORE and state MD
	"216":{38.9853, -76.1677}, // Queenstown, MD
	"216":{38.9170, -75.9419}, // Hillsboro, MD
	"216":{39.3425, -75.8787}, // Galena, MD
	"216":{38.3357, -76.2238}, // Fishing Creek, MD
	"216":{39.2195, -76.0710}, // Chestertown, MD
	"216":{38.6082, -75.9468}, // Secretary, MD
	"216":{38.5036, -76.1530}, // Church Creek, MD
	"216":{39.0316, -75.7821}, // Goldsboro, MD
	"216":{39.1393, -76.2427}, // Rock Hall, MD
	"216":{38.9559, -76.1948}, // Grasonville, MD
	"216":{38.8770, -75.8264}, // Denton, MD
	"216":{39.1541, -76.2097}, // Edesville, MD
	"216":{38.9191, -75.9535}, // Queen Anne, MD
	"216":{38.8896, -75.8389}, // West Denton, MD
	"216":{38.7884, -76.2243}, // St. Michaels, MD
	"216":{38.6635, -76.0518}, // Trappe, MD
	"216":{38.7760, -76.0702}, // Easton, MD
	"216":{39.1450, -75.8656}, // Barclay, MD
	"216":{38.9528, -75.8828}, // Ridgely, MD
	"216":{39.1450, -75.9808}, // Church Hill, MD
	"216":{38.5737, -75.7925}, // Brookview, MD
	"216":{38.7112, -75.9090}, // Preston, MD
	"216":{38.5149, -76.2155}, // Madison, MD
	"216":{38.8309, -75.8511}, // Williston, MD
	"216":{38.4730, -76.3126}, // Taylors Island, MD
	"216":{39.1362, -75.7667}, // Templeville, MD
	"216":{39.3059, -75.9947}, // Kennedyville, MD
	"216":{38.5515, -76.0787}, // Cambridge, MD
	"216":{38.6822, -75.9492}, // Choptank, MD
	"216":{39.0750, -75.7661}, // Henderson, MD
	"216":{38.6849, -76.1703}, // Oxford, MD
	"216":{39.2718, -76.0933}, // Worton, MD
	"216":{38.5634, -75.7151}, // Galestown, MD
	"216":{38.5839, -76.0977}, // Algonquin, MD
	"216":{38.9704, -76.2385}, // Kent Narrows, MD
	"216":{38.6268, -75.8640}, // Hurlock, MD
	"216":{38.5970, -75.9233}, // East New Market, MD
	"216":{38.6929, -75.7727}, // Federalsburg, MD
	"216":{39.2036, -76.0465}, // Kingstown, MD
	"216":{38.9677, -76.2823}, // Chester, MD
	"216":{39.1833, -75.8535}, // Sudlersville, MD
	"216":{38.9745, -76.3184}, // Stevensville, MD
	"216":{38.8680, -75.9989}, // Cordova, MD
	"216":{39.1130, -75.7497}, // Marydel, MD
	"216":{39.2194, -76.1971}, // Georgetown, MD
	"216":{38.7027, -76.3354}, // Tilghman Island, MD
	"216":{39.2820, -76.0992}, // Butlertown, MD
	"216":{38.5840, -75.7893}, // Eldorado, MD
	"216":{39.2268, -76.1654}, // Fairlee, MD
	"216":{39.0420, -76.0631}, // Centreville, MD
	"216":{39.3652, -76.0690}, // Betterton, MD
	"216":{38.9764, -75.8080}, // Greensboro, MD
	"216":{39.2629, -75.8360}, // Millington, MD
	"216":{39.2218, -76.2349}, // Tolchester, MD
	// END OF OPTIONS TO PICK
	"217":{39.4336, -77.4157}, // Frederick, MD
	"218":{38.3755, -75.5867}, // Salisbury, MD
	"219":{39.3051, -76.6144}, // Baltimore, MD
	TODO PLACEHOLDER TO FIX ZIP <220> with city NORTHERN VA and state VA
	"220":{38.7773, -77.2633}, // Burke, VA
	"220":{38.8531, -77.2997}, // Fairfax, VA
	"220":{38.9133, -77.3969}, // Franklin Farm, VA
	"220":{38.7466, -77.2754}, // South Run, VA
	"220":{38.7253, -77.2638}, // Crosspointe, VA
	"220":{38.8658, -77.1445}, // Seven Corners, VA
	"220":{38.8356, -77.3186}, // George Mason, VA
	"220":{38.8302, -77.2713}, // Long Branch, VA
	"220":{38.6111, -77.3400}, // Montclair, VA
	"220":{38.9436, -77.1943}, // McLean, VA
	"220":{38.8503, -77.2322}, // Woodburn, VA
	"220":{38.8151, -77.2960}, // Kings Park West, VA
	"220":{38.6984, -77.2163}, // Lorton, VA
	"220":{38.8514, -77.1579}, // Lake Barcroft, VA
	"220":{38.8847, -77.1751}, // Falls Church, VA
	"220":{38.8717, -77.3970}, // Greenbriar, VA
	"220":{38.8896, -77.2055}, // Idylwood, VA
	"220":{38.5696, -77.2895}, // Cherry Hill, VA
	"220":{38.9699, -77.3867}, // Herndon, VA
	"220":{38.5670, -77.3233}, // Dumfries, VA
	"220":{38.7903, -77.2998}, // Burke Centre, VA
	"220":{38.8731, -77.2426}, // Merrifield, VA
	"220":{38.7942, -77.3358}, // Fairfax Station, VA
	"220":{38.7358, -77.1993}, // Newington, VA
	"220":{38.7119, -77.1459}, // Fort Belvoir, VA
	"220":{38.7371, -77.2339}, // Newington Forest, VA
	"220":{38.8324, -77.1960}, // Annandale, VA
	"220":{38.8648, -77.1878}, // West Falls Church, VA
	"220":{38.7026, -77.2422}, // Laurel Hill, VA
	"220":{38.9105, -77.1991}, // Pimmit Hills, VA
	"220":{38.8477, -77.1305}, // Bailey's Crossroads, VA
	"220":{38.8887, -77.3016}, // Oakton, VA
	"220":{38.8653, -77.3586}, // Fair Oaks, VA
	"220":{38.9497, -77.3461}, // Reston, VA
	"220":{38.6556, -77.1819}, // Mason Neck, VA
	"220":{38.8945, -77.2315}, // Dunn Loring, VA
	"220":{38.8526, -77.2571}, // Mantua, VA
	"220":{39.0110, -77.3013}, // Great Falls, VA
	"220":{38.8530, -77.3885}, // Fair Lakes, VA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <221> with city NORTHERN VA and state VA
	"221":{38.6473, -77.3459}, // Dale City, VA
	"221":{38.6917, -77.3506}, // County Center, VA
	"221":{38.7773, -77.2633}, // Burke, VA
	"221":{38.7466, -77.2754}, // South Run, VA
	"221":{38.6373, -77.2620}, // Marumsco, VA
	"221":{38.8026, -77.2396}, // Kings Park, VA
	"221":{38.5228, -77.3187}, // Quantico Base, VA
	"221":{38.8356, -77.3186}, // George Mason, VA
	"221":{38.9395, -77.2842}, // Wolf Trap, VA
	"221":{38.8024, -77.2028}, // North Springfield, VA
	"221":{38.6111, -77.3400}, // Montclair, VA
	"221":{38.9436, -77.1943}, // McLean, VA
	"221":{38.7809, -77.1839}, // Springfield, VA
	"221":{38.6984, -77.2163}, // Lorton, VA
	"221":{38.6843, -77.3059}, // Lake Ridge, VA
	"221":{38.8996, -77.2597}, // Vienna, VA
	"221":{38.6547, -77.3012}, // Potomac Mills, VA
	"221":{38.5696, -77.2895}, // Cherry Hill, VA
	"221":{38.5224, -77.2902}, // Quantico, VA
	"221":{38.6569, -77.2403}, // Woodbridge, VA
	"221":{38.8731, -77.2426}, // Merrifield, VA
	"221":{38.6825, -77.2606}, // Occoquan, VA
	"221":{38.7358, -77.1993}, // Newington, VA
	"221":{38.7771, -77.2268}, // West Springfield, VA
	"221":{38.6083, -77.2847}, // Neabsco, VA
	"221":{38.7371, -77.2339}, // Newington Forest, VA
	"221":{38.8887, -77.3016}, // Oakton, VA
	"221":{38.5483, -77.3195}, // Triangle, VA
	"221":{38.7140, -77.1043}, // Mount Vernon, VA
	"221":{38.8945, -77.2315}, // Dunn Loring, VA
	"221":{38.8526, -77.2571}, // Mantua, VA
	"221":{38.8032, -77.2223}, // Ravensworth, VA
	"221":{38.6404, -77.4090}, // Independent Hill, VA
	// END OF OPTIONS TO PICK
	"222":{38.8786, -77.1011}, // Arlington, VA
	"223":{38.8185, -77.0861}, // Alexandria, VA
	"224":{37.5295, -77.4756}, // Richmond, VA
	"225":{37.5295, -77.4756}, // Richmond, VA
	"226":{39.1735, -78.1746}, // Winchester, VA
	"227":{38.4705, -78.0001}, // Culpeper, VA
	TODO PLACEHOLDER TO FIX ZIP <228> with city CHARLOTTESVLE and state VA
	"228":{38.7394, -78.6513}, // Mount Jackson, VA
	"228":{38.4465, -78.9228}, // Belmont Estates, VA
	"228":{38.6341, -78.7730}, // Timberville, VA
	"228":{38.8235, -78.5634}, // Edinburg, VA
	"228":{38.8181, -78.7659}, // Basye, VA
	"228":{38.4362, -78.8735}, // Harrisonburg, VA
	"228":{38.6459, -78.6709}, // New Market, VA
	"228":{38.3600, -78.9413}, // Mount Crawford, VA
	"228":{38.5769, -78.5027}, // Stanley, VA
	"228":{38.6082, -78.8017}, // Broadway, VA
	"228":{38.4172, -78.9412}, // Dayton, VA
	"228":{38.3863, -78.9675}, // Bridgewater, VA
	"228":{38.4875, -78.6172}, // Shenandoah, VA
	"228":{38.3899, -78.8340}, // Massanetta Springs, VA
	"228":{38.6650, -78.4547}, // Luray, VA
	"228":{38.4114, -78.7270}, // Massanutten, VA
	"228":{38.4106, -78.6161}, // Elkton, VA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <229> with city CHARLOTTESVLE and state VA
	"229":{38.1266, -78.4386}, // Hollymead, VA
	"229":{38.1618, -78.8413}, // Crimora, VA
	"229":{38.2991, -78.4367}, // Stanardsville, VA
	"229":{37.9161, -78.9382}, // Wintergreen, VA
	"229":{38.1539, -78.5595}, // Free Union, VA
	"229":{38.0674, -78.9014}, // Waynesboro, VA
	"229":{38.0221, -78.9516}, // Lyndhurst, VA
	"229":{38.0308, -78.4444}, // Pantops, VA
	"229":{37.9889, -78.9442}, // Sherando, VA
	"229":{37.8652, -78.2627}, // Palmyra, VA
	"229":{38.1591, -78.4156}, // Piney Mountain, VA
	"229":{38.0611, -78.5987}, // Ivy, VA
	"229":{38.0645, -78.6962}, // Crozet, VA
	"229":{38.2314, -78.3756}, // Ruckersville, VA
	"229":{38.1004, -78.9686}, // Fishersville, VA
	"229":{38.1360, -78.1879}, // Gordonsville, VA
	"229":{37.7605, -78.8689}, // Lovingston, VA
	"229":{37.9210, -78.3295}, // Lake Monticello, VA
	"229":{38.0375, -78.4855}, // Charlottesville, VA
	"229":{37.6744, -78.8928}, // Arrington, VA
	"229":{37.8284, -78.5959}, // Esmont, VA
	"229":{38.0405, -78.5163}, // University of Virginia, VA
	"229":{38.2505, -78.4417}, // Twin Lakes, VA
	"229":{38.1027, -78.8495}, // Dooms, VA
	"229":{37.7272, -78.8423}, // Shipman, VA
	"229":{38.2486, -78.1127}, // Orange, VA
	"229":{38.1954, -78.9044}, // New Hope, VA
	"229":{37.9130, -78.8828}, // Nellysford, VA
	"229":{37.9945, -78.3788}, // Rivanna, VA
	"229":{37.7939, -78.7002}, // Schuyler, VA
	// END OF OPTIONS TO PICK
	"230":{37.5295, -77.4756}, // Richmond, VA
	"231":{37.5295, -77.4756}, // Richmond, VA
	"232":{37.5295, -77.4756}, // Richmond, VA
	"233":{36.8945, -76.2590}, // Norfolk, VA
	"234":{36.8945, -76.2590}, // Norfolk, VA
	"235":{36.8945, -76.2590}, // Norfolk, VA
	"236":{36.8945, -76.2590}, // Norfolk, VA
	"237":{36.8468, -76.3540}, // Portsmouth, VA
	"238":{37.5295, -77.4756}, // Richmond, VA
	"239":{37.2959, -78.4002}, // Farmville, VA
	"240":{37.2785, -79.9580}, // Roanoke, VA
	"241":{37.2785, -79.9580}, // Roanoke, VA
	"242":{36.6179, -82.1607}, // Bristol, VA
	"243":{37.2785, -79.9580}, // Roanoke, VA
	TODO PLACEHOLDER TO FIX ZIP <244> with city CHARLOTTESVLE and state VA
	"244":{37.7985, -79.7903}, // Iron Gate, VA
	"244":{38.1618, -78.8413}, // Crimora, VA
	"244":{38.0245, -79.0308}, // Stuarts Draft, VA
	"244":{38.1954, -79.4105}, // Deerfield, VA
	"244":{38.2247, -79.1619}, // Churchville, VA
	"244":{38.0013, -79.8239}, // Hot Springs, VA
	"244":{38.2586, -78.9654}, // Mount Sidney, VA
	"244":{38.0463, -79.7826}, // Warm Springs, VA
	"244":{38.0221, -78.9516}, // Lyndhurst, VA
	"244":{37.8232, -79.8250}, // Clifton Forge, VA
	"244":{37.9889, -78.9442}, // Sherando, VA
	"244":{37.8112, -80.0871}, // Callaghan, VA
	"244":{38.1149, -79.0668}, // Jolivue, VA
	"244":{38.1939, -79.0087}, // Verona, VA
	"244":{37.7319, -79.3569}, // Buena Vista, VA
	"244":{38.1004, -78.9686}, // Fishersville, VA
	"244":{38.0036, -79.1522}, // Greenville, VA
	"244":{38.1593, -79.0611}, // Staunton, VA
	"244":{38.2120, -78.8283}, // Harriston, VA
	"244":{38.2692, -78.8252}, // Grottoes, VA
	"244":{38.0552, -79.2182}, // Middlebrook, VA
	"244":{37.8033, -79.8499}, // Selma, VA
	"244":{38.1079, -79.3352}, // Augusta Springs, VA
	"244":{38.4116, -79.5809}, // Monterey, VA
	"244":{37.7785, -79.9868}, // Covington, VA
	"244":{38.2844, -78.9125}, // Weyers Cave, VA
	"244":{37.9901, -79.5067}, // Goshen, VA
	"244":{38.1954, -78.9044}, // New Hope, VA
	"244":{37.7944, -79.8706}, // Low Moor, VA
	"244":{38.0838, -79.3853}, // Craigsville, VA
	"244":{37.7825, -79.4440}, // Lexington, VA
	"244":{37.8009, -79.4160}, // East Lexington, VA
	// END OF OPTIONS TO PICK
	"245":{37.4003, -79.1909}, // Lynchburg, VA
	"246":{37.2608, -81.2143}, // Bluefield, WV
	"247":{37.2608, -81.2143}, // Bluefield, WV
	"248":{37.2608, -81.2143}, // Bluefield, WV
	"249":{37.8096, -80.4327}, // Lewisburg, WV
	"250":{38.3484, -81.6323}, // Charleston, WV
	"251":{38.3484, -81.6323}, // Charleston, WV
	"252":{38.3484, -81.6323}, // Charleston, WV
	"253":{38.3484, -81.6323}, // Charleston, WV
	"254":{39.4582, -77.9776}, // Martinsburg, WV
	"255":{38.4109, -82.4344}, // Huntington, WV
	"256":{38.4109, -82.4344}, // Huntington, WV
	"257":{38.4109, -82.4344}, // Huntington, WV
	"258":{37.7877, -81.1840}, // Beckley, WV
	"259":{37.7877, -81.1840}, // Beckley, WV
	"260":{40.0751, -80.6951}, // Wheeling, WV
	"261":{39.2624, -81.5419}, // Parkersburg, WV
	"262":{39.2863, -80.3230}, // Clarksburg, WV
	"263":{39.2863, -80.3230}, // Clarksburg, WV
	"264":{39.2863, -80.3230}, // Clarksburg, WV
	"265":{39.2863, -80.3230}, // Clarksburg, WV
	"266":{38.6702, -80.7717}, // Gassaway, WV
	"267":{39.6515, -78.7585}, // Cumberland, MD
	"268":{38.9957, -79.1276}, // Petersburg, WV
	"270":{36.0956, -79.8268}, // Greensboro, NC
	"271":{36.1029, -80.2610}, // Winston-Salem, NC
	"272":{36.0956, -79.8268}, // Greensboro, NC
	"273":{36.0956, -79.8268}, // Greensboro, NC
	"274":{36.0956, -79.8268}, // Greensboro, NC
	"275":{35.8324, -78.6438}, // Raleigh, NC
	"276":{35.8324, -78.6438}, // Raleigh, NC
	"277":{35.9795, -78.9032}, // Durham, NC
	"278":{35.9676, -77.8047}, // Rocky Mount, NC
	"279":{35.9676, -77.8047}, // Rocky Mount, NC
	"280":{35.2079, -80.8304}, // Charlotte, NC
	"281":{35.2079, -80.8304}, // Charlotte, NC
	"282":{35.2079, -80.8304}, // Charlotte, NC
	"283":{35.0846, -78.9776}, // Fayetteville, NC
	"284":{35.0846, -78.9776}, // Fayetteville, NC
	"285":{35.2748, -77.5937}, // Kinston, NC
	"286":{35.7426, -81.3230}, // Hickory, NC
	"287":{35.5704, -82.5537}, // Asheville, NC
	"288":{35.5704, -82.5537}, // Asheville, NC
	"289":{35.5704, -82.5537}, // Asheville, NC
	"290":{34.0376, -80.9037}, // Columbia, SC
	"291":{34.0376, -80.9037}, // Columbia, SC
	"292":{34.0376, -80.9037}, // Columbia, SC
	"293":{34.8362, -82.3649}, // Greenville, SC
	"294":{32.8151, -79.9630}, // Charleston, SC
	"295":{34.1782, -79.7872}, // Florence, SC
	"296":{34.8362, -82.3649}, // Greenville, SC
	"297":{35.2079, -80.8304}, // Charlotte, NC
	"298":{33.3645, -82.0708}, // Augusta, GA
	"299":{32.0281, -81.1785}, // Savannah, GA
	TODO PLACEHOLDER TO FIX ZIP <300> with city NORTH METRO and state GA
	"300":{33.9670, -84.2319}, // Peachtree Corners, GA
	"300":{34.0507, -84.0687}, // Suwanee, GA
	"300":{33.5739, -83.8940}, // Porterdale, GA
	"300":{33.7488, -84.2598}, // Belvedere Park, GA
	"300":{33.8034, -84.1724}, // Stone Mountain, GA
	"300":{34.0333, -84.2026}, // Johns Creek, GA
	"300":{33.7455, -83.8502}, // Walnut Grove, GA
	"300":{33.9805, -84.1839}, // Berkeley Lake, GA
	"300":{33.7711, -84.2968}, // Decatur, GA
	"300":{33.5180, -83.7349}, // Mansfield, GA
	"300":{33.8436, -84.2021}, // Tucker, GA
	"300":{33.6843, -84.1373}, // Stonecrest, GA
	"300":{34.0391, -84.3513}, // Roswell, GA
	"300":{33.7059, -84.2764}, // Panthersville, GA
	"300":{33.9532, -84.5422}, // Marietta, GA
	"300":{33.6645, -83.9967}, // Conyers, GA
	"300":{33.7267, -84.2723}, // Candler-McAfee, GA
	"300":{33.6277, -83.8721}, // Oxford, GA
	"300":{34.0703, -84.2738}, // Alpharetta, GA
	"300":{34.0046, -83.8128}, // Carl, GA
	"300":{33.7699, -84.2648}, // Avondale Estates, GA
	"300":{33.8561, -84.0038}, // Snellville, GA
	"300":{33.8117, -84.2405}, // Clarkston, GA
	"300":{33.9525, -83.9928}, // Lawrenceville, GA
	"300":{33.7904, -84.2060}, // Pine Lake, GA
	"300":{33.7393, -84.1644}, // Redan, GA
	"300":{33.8073, -84.2889}, // North Decatur, GA
	"300":{34.0157, -83.8319}, // Auburn, GA
	"300":{33.9815, -83.8951}, // Dacula, GA
	"300":{33.5163, -83.6959}, // Newborn, GA
	"300":{34.1353, -84.3138}, // Milton, GA
	"300":{33.7046, -84.0367}, // Lakeview Estates, GA
	"300":{33.8186, -84.3254}, // North Druid Hills, GA
	"300":{34.2064, -84.1337}, // Cumming, GA
	"300":{34.0053, -84.1491}, // Duluth, GA
	"300":{33.6049, -83.8465}, // Covington, GA
	"300":{33.9193, -84.5444}, // Fair Oaks, GA
	"300":{33.8353, -83.8957}, // Loganville, GA
	"300":{33.8633, -84.5168}, // Smyrna, GA
	"300":{33.8900, -83.9574}, // Grayson, GA
	"300":{33.8887, -84.1379}, // Lilburn, GA
	"300":{33.7129, -84.1061}, // Lithonia, GA
	"300":{33.6505, -83.7116}, // Social Circle, GA
	"300":{33.9379, -84.2065}, // Norcross, GA
	"300":{34.0830, -84.4133}, // Mountain Park, GA
	"300":{33.7177, -83.8001}, // Jersey, GA
	"300":{33.7950, -84.2634}, // Scottdale, GA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <301> with city NORTH METRO and state GA
	"301":{33.7811, -84.6486}, // Lithia Springs, GA
	"301":{33.7305, -84.9170}, // Villa Rica, GA
	"301":{34.0225, -85.2480}, // Cedartown, GA
	"301":{33.7384, -84.7068}, // Douglasville, GA
	"301":{34.2662, -85.1862}, // Rome, GA
	"301":{34.1026, -84.5085}, // Woodstock, GA
	"301":{33.8660, -84.6839}, // Powder Springs, GA
	"301":{33.7335, -85.2859}, // Tallapoosa, GA
	"301":{34.3702, -84.9214}, // Adairsville, GA
	"301":{34.3426, -84.3635}, // Ball Ground, GA
	"301":{34.5279, -84.4939}, // Talking Rock, GA
	"301":{34.4379, -84.6999}, // Fairmount, GA
	"301":{33.8202, -84.6452}, // Austell, GA
	"301":{33.4567, -85.1304}, // Roopville, GA
	"301":{33.8132, -84.5656}, // Mableton, GA
	"301":{34.2824, -84.7472}, // White, GA
	"301":{34.0459, -85.0570}, // Aragon, GA
	"301":{34.2466, -84.4901}, // Canton, GA
	"301":{33.8771, -84.7710}, // Hiram, GA
	"301":{33.6409, -85.1801}, // Mount Zion, GA
	"301":{33.5818, -85.0837}, // Carrollton, GA
	"301":{34.4709, -84.4496}, // Jasper, GA
	"301":{34.1441, -84.9327}, // Euharlee, GA
	"301":{33.7086, -85.1498}, // Bremen, GA
	"301":{34.3790, -84.3708}, // Nelson, GA
	"301":{33.4928, -84.9137}, // Whitesburg, GA
	"301":{33.5378, -85.2540}, // Bowdon, GA
	"301":{34.0262, -84.6177}, // Kennesaw, GA
	"301":{34.0103, -85.0441}, // Rockmart, GA
	"301":{33.4069, -85.2568}, // Ephesus, GA
	"301":{34.1884, -85.1808}, // Lindale, GA
	"301":{33.8029, -85.1804}, // Buchanan, GA
	"301":{33.7342, -85.0289}, // Temple, GA
	"301":{34.0860, -84.9876}, // Taylorsville, GA
	"301":{33.8633, -84.5168}, // Smyrna, GA
	"301":{34.1303, -84.7483}, // Emerson, GA
	"301":{34.1632, -84.8007}, // Cartersville, GA
	"301":{33.9836, -84.9591}, // Braswell, GA
	"301":{34.3171, -84.5505}, // Waleska, GA
	"301":{33.9153, -84.8416}, // Dallas, GA
	"301":{34.1081, -85.3389}, // Cave Spring, GA
	"301":{34.0565, -84.6709}, // Acworth, GA
	"301":{33.7026, -85.1893}, // Waco, GA
	"301":{34.1686, -84.4843}, // Holly Springs, GA
	"301":{34.2320, -84.9445}, // Kingston, GA
	"301":{34.3406, -85.0854}, // Shannon, GA
	// END OF OPTIONS TO PICK
	"302":{33.7627, -84.4225}, // Atlanta, GA
	"303":{33.7627, -84.4225}, // Atlanta, GA
	"304":{32.5866, -82.3345}, // Swainsboro, GA
	"305":{33.9508, -83.3689}, // Athens, GA
	"306":{33.9508, -83.3689}, // Athens, GA
	"307":{35.0657, -85.2487}, // Chattanooga, TN
	"308":{33.3645, -82.0708}, // Augusta, GA
	"309":{33.3645, -82.0708}, // Augusta, GA
	"310":{32.8065, -83.6974}, // Macon, GA
	"311":{33.7627, -84.4225}, // Atlanta, GA
	"312":{32.8065, -83.6974}, // Macon, GA
	"313":{32.0281, -81.1785}, // Savannah, GA
	"314":{32.0281, -81.1785}, // Savannah, GA
	"315":{31.2108, -82.3579}, // Waycross, GA
	"316":{30.8502, -83.2788}, // Valdosta, GA
	"317":{31.5776, -84.1762}, // Albany, GA
	"318":{32.5100, -84.8771}, // Columbus, GA
	"319":{32.5100, -84.8771}, // Columbus, GA
	"320":{30.3322, -81.6749}, // Jacksonville, FL
	"321":{29.1994, -81.0982}, // Daytona Beach, FL
	"322":{30.3322, -81.6749}, // Jacksonville, FL
	"323":{30.4551, -84.2527}, // Tallahassee, FL
	"324":{30.1995, -85.6003}, // Panama City, FL
	"325":{30.4427, -87.1886}, // Pensacola, FL
	"326":{29.6804, -82.3458}, // Gainesville, FL
	TODO PLACEHOLDER TO FIX ZIP <327> with city MID-FLORIDA and state FL
	"327":{28.7893, -81.2760}, // Sanford, FL
	"327":{28.6928, -80.8468}, // Mims, FL
	"327":{28.7754, -81.3721}, // Heathrow, FL
	"327":{28.8815, -81.3240}, // DeBary, FL
	"327":{29.0162, -81.3343}, // West DeLand, FL
	"327":{28.6271, -81.4354}, // Lockhart, FL
	"327":{28.7377, -81.1143}, // Geneva, FL
	"327":{28.7107, -81.1801}, // Black Hammock, FL
	"327":{29.1172, -81.3517}, // De Leon Springs, FL
	"327":{28.6569, -81.5056}, // South Apopka, FL
	"327":{28.9681, -81.6482}, // Altoona, FL
	"327":{28.5727, -80.8193}, // Titusville, FL
	"327":{28.5606, -81.0201}, // Christmas, FL
	"327":{28.9828, -81.2306}, // Lake Helen, FL
	"327":{28.7926, -81.7366}, // Tavares, FL
	"327":{28.6237, -81.5439}, // Paradise Heights, FL
	"327":{28.8560, -81.6781}, // Eustis, FL
	"327":{28.5987, -81.3438}, // Winter Park, FL
	"327":{28.6615, -81.3953}, // Altamonte Springs, FL
	"327":{28.8024, -81.5351}, // Mount Plymouth, FL
	"327":{28.6619, -81.4443}, // Forest City, FL
	"327":{28.5818, -81.4693}, // Pine Hills, FL
	"327":{28.6580, -81.1872}, // Oviedo, FL
	"327":{28.6169, -81.3911}, // Eatonville, FL
	"327":{28.9349, -81.2882}, // Orange City, FL
	"327":{28.6380, -81.1158}, // Chuluota, FL
	"327":{29.0006, -81.4240}, // Lake Mack-Forest Hills, FL
	"327":{28.7013, -81.5303}, // Apopka, FL
	"327":{28.7592, -81.3360}, // Lake Mary, FL
	"327":{28.6295, -81.3718}, // Maitland, FL
	"327":{28.7301, -81.5915}, // Zellwood, FL
	"327":{28.6484, -81.3457}, // Fern Park, FL
	"327":{28.9967, -81.6462}, // Pittman, FL
	"327":{28.9050, -81.2136}, // Deltona, FL
	"327":{28.6625, -81.3218}, // Casselberry, FL
	"327":{29.0076, -81.3113}, // DeLand Southwest, FL
	"327":{28.9274, -81.6651}, // Umatilla, FL
	"327":{28.9850, -81.5404}, // Paisley, FL
	"327":{28.8087, -81.5631}, // Sorrento, FL
	"327":{28.6889, -81.2704}, // Winter Springs, FL
	"327":{28.7014, -81.3487}, // Longwood, FL
	"327":{28.6201, -81.5002}, // Clarcona, FL
	"327":{29.0484, -81.2966}, // North DeLand, FL
	"327":{28.6021, -81.3946}, // Fairview Shores, FL
	"327":{29.0224, -81.2873}, // DeLand, FL
	"327":{28.7589, -81.6341}, // Tangerine, FL
	"327":{28.6984, -81.4251}, // Wekiwa Springs, FL
	"327":{28.8784, -80.8351}, // Oak Hill, FL
	"327":{28.6114, -81.2916}, // Goldenrod, FL
	"327":{29.0072, -81.4906}, // Lake Kathryn, FL
	"327":{28.8143, -81.6344}, // Mount Dora, FL
	"327":{28.9390, -81.4308}, // Pine Lakes, FL
	// END OF OPTIONS TO PICK
	"328":{28.4772, -81.3369}, // Orlando, FL
	"329":{28.4772, -81.3369}, // Orlando, FL
	TODO PLACEHOLDER TO FIX ZIP <330> with city SOUTH FLORIDA and state FL
	"330":{25.2574, -80.3242}, // North Key Largo, FL
	"330":{26.3218, -80.2533}, // Parkland, FL
	"330":{26.0463, -80.2862}, // Cooper City, FL
	"330":{26.2466, -80.2119}, // Margate, FL
	"330":{26.2802, -80.1842}, // Coconut Creek, FL
	"330":{25.9433, -80.2426}, // Miami Gardens, FL
	"330":{25.9407, -80.3102}, // Country Club, FL
	"330":{26.0294, -80.1679}, // Hollywood, FL
	"330":{25.5396, -80.3971}, // Princeton, FL
	"330":{26.2785, -80.0891}, // Lighthouse Point, FL
	"330":{25.0188, -80.5132}, // Tavernier, FL
	"330":{25.9854, -80.1423}, // Hallandale Beach, FL
	"330":{25.5164, -80.4221}, // Naranja, FL
	"330":{24.9408, -80.6097}, // Islamorada, Village of Islands, FL
	"330":{24.6893, -81.3676}, // Big Pine Key, FL
	"330":{25.8878, -80.3569}, // Hialeah Gardens, FL
	"330":{24.7262, -81.0376}, // Marathon, FL
	"330":{25.4665, -80.4472}, // Homestead, FL
	"330":{24.5658, -81.7351}, // Stock Island, FL
	"330":{26.3253, -80.1947}, // Hillsboro Pines, FL
	"330":{25.9840, -80.1923}, // West Park, FL
	"330":{24.5915, -81.6574}, // Big Coppitt Key, FL
	"330":{24.7233, -81.0215}, // Key Colony Beach, FL
	"330":{25.9125, -80.3214}, // Miami Lakes, FL
	"330":{24.5637, -81.7768}, // Key West, FL
	"330":{25.8696, -80.3046}, // Hialeah, FL
	"330":{25.1224, -80.4120}, // Key Largo, FL
	"330":{25.4418, -80.4685}, // Florida City, FL
	"330":{26.1605, -80.2242}, // Lauderhill, FL
	"330":{25.4937, -80.4369}, // Leisure City, FL
	"330":{26.0128, -80.3382}, // Pembroke Pines, FL
	"330":{26.2113, -80.2209}, // North Lauderdale, FL
	"330":{25.9852, -80.1777}, // Pembroke Park, FL
	"330":{26.2702, -80.2591}, // Coral Springs, FL
	"330":{24.7710, -80.9132}, // Duck Key, FL
	"330":{26.1990, -80.0972}, // Lauderdale-by-the-Sea, FL
	"330":{24.6746, -81.4986}, // Cudjoe Key, FL
	"330":{25.4936, -80.3911}, // Homestead Base, FL
	"330":{26.2837, -80.0796}, // Hillsboro Beach, FL
	"330":{26.0593, -80.1638}, // Dania Beach, FL
	"330":{24.8247, -80.8129}, // Layton, FL
	"330":{25.8997, -80.2551}, // Opa-locka, FL
	"330":{25.8195, -80.2894}, // Miami Springs, FL
	"330":{25.9773, -80.3351}, // Miramar, FL
	"330":{26.2428, -80.1312}, // Pompano Beach, FL
	"330":{25.9351, -80.3339}, // Palm Springs North, FL
	"330":{26.3049, -80.1277}, // Deerfield Beach, FL
	// END OF OPTIONS TO PICK
	"331":{25.7839, -80.2102}, // Miami, FL
	"332":{25.7839, -80.2102}, // Miami, FL
	TODO PLACEHOLDER TO FIX ZIP <333> with city FT LAUDERDALE and state FL
	"333":{26.1593, -80.1395}, // Wilton Manors, FL
	"333":{26.2057, -80.2549}, // Tamarac, FL
	"333":{26.0463, -80.2862}, // Cooper City, FL
	"333":{26.1260, -80.2617}, // Plantation, FL
	"333":{26.0294, -80.1679}, // Hollywood, FL
	"333":{26.1340, -80.1762}, // Franklin Park, FL
	"333":{26.1780, -80.1528}, // Oakland Park, FL
	"333":{26.1547, -80.2997}, // Sunrise, FL
	"333":{26.1407, -80.1809}, // Roosevelt Gardens, FL
	"333":{26.1252, -80.1822}, // Boulevard Gardens, FL
	"333":{26.0477, -80.3752}, // Southwest Ranches, FL
	"333":{26.1605, -80.2242}, // Lauderhill, FL
	"333":{26.1304, -80.1801}, // Washington Park, FL
	"333":{26.0979, -80.2088}, // Broadview Park, FL
	"333":{26.1682, -80.2017}, // Lauderdale Lakes, FL
	"333":{26.1990, -80.0972}, // Lauderdale-by-the-Sea, FL
	"333":{26.0789, -80.2870}, // Davie, FL
	"333":{26.0593, -80.1638}, // Dania Beach, FL
	"333":{26.2007, -80.0981}, // Sea Ranch Lakes, FL
	"333":{26.1412, -80.1464}, // Fort Lauderdale, FL
	"333":{26.1563, -80.1452}, // Lazy Lake, FL
	"333":{26.1007, -80.4054}, // Weston, FL
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <334> with city WEST PALM BCH and state FL
	"334":{26.6815, -80.1265}, // Royal Palm Estates, FL
	"334":{26.8486, -80.1660}, // Palm Beach Gardens, FL
	"334":{26.7325, -80.9518}, // Harlem, FL
	"334":{26.6195, -80.0591}, // Lake Worth, FL
	"334":{26.5634, -80.0531}, // Hypoluxo, FL
	"334":{26.8554, -80.0864}, // Cabana Colony, FL
	"334":{26.5899, -80.0390}, // South Palm Beach, FL
	"334":{26.6909, -80.1217}, // Haverhill, FL
	"334":{26.6885, -80.8079}, // Lake Harbor, FL
	"334":{26.6784, -80.7269}, // South Bay, FL
	"334":{26.5091, -80.0543}, // Briny Breezes, FL
	"334":{26.6932, -80.0406}, // Palm Beach, FL
	"334":{26.6349, -80.0968}, // Palm Springs, FL
	"334":{26.9199, -80.1128}, // Jupiter, FL
	"334":{26.5834, -80.0564}, // Lantana, FL
	"334":{26.9618, -80.1011}, // Tequesta, FL
	"334":{26.6458, -80.1104}, // Acacia Villas, FL
	"334":{26.5651, -80.0611}, // San Castle, FL
	"334":{26.5287, -80.0499}, // Ocean Ridge, FL
	"334":{26.7774, -80.0368}, // Palm Beach Shores, FL
	"334":{26.6919, -80.6672}, // Belle Glade, FL
	"334":{26.7038, -80.2241}, // Royal Palm Beach, FL
	"334":{26.7549, -80.2984}, // Westlake, FL
	"334":{26.9433, -80.1408}, // Limestone Creek, FL
	"334":{26.6888, -80.1371}, // Lake Belvedere Estates, FL
	"334":{26.8754, -80.0589}, // Juno Beach, FL
	"334":{26.7812, -80.0741}, // Riviera Beach, FL
	"334":{26.8338, -81.0985}, // Moore Haven, FL
	"334":{26.7258, -81.2190}, // Pioneer, FL
	"334":{26.6715, -80.0753}, // Glen Ridge, FL
	"334":{26.6994, -80.0989}, // Westgate, FL
	"334":{26.6276, -80.1151}, // Kenwood Estates, FL
	"334":{26.8202, -80.6622}, // Pahokee, FL
	"334":{26.9480, -80.0755}, // Jupiter Inlet Colony, FL
	"334":{26.4550, -80.0905}, // Delray Beach, FL
	"334":{26.7152, -80.1143}, // Schall Circle, FL
	"334":{26.6464, -80.2706}, // Wellington, FL
	"334":{27.0451, -80.1101}, // Jupiter Island, FL
	"334":{26.5043, -80.1090}, // Golf, FL
	"334":{26.7998, -80.0685}, // Lake Park, FL
	"334":{26.6460, -80.0752}, // Lake Clarke Shores, FL
	"334":{26.8217, -80.0577}, // North Palm Beach, FL
	"334":{26.8491, -80.0620}, // Juno Ridge, FL
	"334":{26.8626, -80.6222}, // Canal Point, FL
	"334":{26.5962, -80.1030}, // Atlantis, FL
	"334":{26.5840, -80.1001}, // Seminole Manor, FL
	"334":{26.6429, -81.0938}, // Montura, FL
	"334":{26.7741, -80.2779}, // The Acreage, FL
	"334":{26.3359, -80.2126}, // Watergate, FL
	"334":{26.6743, -80.0731}, // Cloud Lake, FL
	"334":{26.4890, -80.0575}, // Gulf Stream, FL
	"334":{26.5281, -80.0811}, // Boynton Beach, FL
	"334":{26.7586, -80.0761}, // Mangonia Park, FL
	"334":{26.6588, -80.1073}, // Pine Air, FL
	"334":{26.6978, -80.1238}, // Stacey Street, FL
	"334":{26.7106, -80.2764}, // Loxahatchee Groves, FL
	"334":{26.7529, -80.9399}, // Clewiston, FL
	"334":{26.6272, -80.1372}, // Greenacres, FL
	"334":{26.7469, -80.1316}, // West Palm Beach, FL
	"334":{26.6753, -80.1080}, // Gun Club Estates, FL
	"334":{26.4088, -80.0661}, // Highland Beach, FL
	"334":{26.3749, -80.1077}, // Boca Raton, FL
	"334":{27.0729, -80.1425}, // Hobe Sound, FL
	"334":{26.3049, -80.1277}, // Deerfield Beach, FL
	"334":{26.5621, -80.0439}, // Manalapan, FL
	"334":{26.9224, -80.2187}, // Jupiter Farms, FL
	// END OF OPTIONS TO PICK
	"335":{27.9942, -82.4451}, // Tampa, FL
	"336":{27.9942, -82.4451}, // Tampa, FL
	TODO PLACEHOLDER TO FIX ZIP <337> with city ST PETERSBURG and state FL
	"337":{27.9087, -82.7162}, // South Highpoint, FL
	"337":{27.8155, -82.7162}, // Kenneth City, FL
	"337":{27.9088, -82.7712}, // Largo, FL
	"337":{27.9061, -82.6813}, // Feather Sound, FL
	"337":{27.7463, -82.7100}, // Gulfport, FL
	"337":{27.7930, -82.6652}, // St. Petersburg, FL
	"337":{27.7235, -82.7387}, // St. Pete Beach, FL
	"337":{27.7985, -82.7887}, // Madeira Beach, FL
	"337":{27.6684, -82.7300}, // Tierra Verde, FL
	"337":{27.9173, -82.8455}, // Belleair Shore, FL
	"337":{27.7542, -82.7285}, // Bear Creek, FL
	"337":{27.8952, -82.8064}, // Ridgecrest, FL
	"337":{27.8589, -82.7076}, // Pinellas Park, FL
	"337":{27.9242, -82.8358}, // Belleair Beach, FL
	"337":{27.9788, -82.7624}, // Clearwater, FL
	"337":{28.0112, -82.7523}, // Greenbriar, FL
	"337":{27.8536, -82.8438}, // Indian Shores, FL
	"337":{27.9365, -82.8114}, // Belleair, FL
	"337":{27.8191, -82.7385}, // West Lealman, FL
	"337":{27.7526, -82.7394}, // South Pasadena, FL
	"337":{27.8434, -82.7839}, // Seminole, FL
	"337":{27.8575, -82.7534}, // Bardmoor, FL
	"337":{27.8128, -82.8096}, // Redington Beach, FL
	"337":{27.8139, -82.7747}, // Bay Pines, FL
	"337":{27.8293, -82.8278}, // Redington Shores, FL
	"337":{27.8961, -82.8444}, // Indian Rocks Beach, FL
	"337":{27.9197, -82.8193}, // Belleair Bluffs, FL
	"337":{27.9080, -82.8270}, // Harbor Bluffs, FL
	"337":{27.7740, -82.7663}, // Treasure Island, FL
	"337":{27.8213, -82.8190}, // North Redington Beach, FL
	"337":{27.8197, -82.6847}, // Lealman, FL
	// END OF OPTIONS TO PICK
	"338":{28.0557, -81.9545}, // Lakeland, FL
	TODO PLACEHOLDER TO FIX ZIP <339> with city FT MYERS and state FL
	"339":{26.7123, -81.8684}, // Suncoast Estates, FL
	"339":{26.5507, -82.0980}, // St. James City, FL
	"339":{26.6304, -82.0719}, // Matlacha, FL
	"339":{26.6352, -82.1253}, // Pine Island Center, FL
	"339":{26.9031, -82.0486}, // Charlotte Park, FL
	"339":{26.5205, -82.1910}, // Captiva, FL
	"339":{26.6615, -81.7399}, // Buckingham, FL
	"339":{26.7135, -81.7383}, // Fort Myers Shores, FL
	"339":{26.9388, -82.0278}, // Solana, FL
	"339":{26.6195, -81.8303}, // Fort Myers, FL
	"339":{26.9928, -82.0072}, // Harbour Heights, FL
	"339":{26.4323, -81.9167}, // Fort Myers Beach, FL
	"339":{26.5391, -81.9000}, // Cypress Lake, FL
	"339":{26.4277, -81.7951}, // Estero, FL
	"339":{26.6120, -81.6388}, // Lehigh Acres, FL
	"339":{26.5505, -81.8678}, // Villas, FL
	"339":{26.4765, -81.8193}, // San Carlos Park, FL
	"339":{26.7188, -81.6269}, // Alva, FL
	"339":{26.7113, -81.6950}, // Olga, FL
	"339":{26.5611, -81.9134}, // McGregor, FL
	"339":{26.4738, -81.7960}, // Three Oaks, FL
	"339":{26.5036, -81.9998}, // Punta Rassa, FL
	"339":{26.7057, -81.5806}, // Charleston Park, FL
	"339":{26.6445, -81.9955}, // Cape Coral, FL
	"339":{26.5733, -81.8903}, // Whiskey Creek, FL
	"339":{26.7219, -81.4506}, // LaBelle, FL
	"339":{26.8934, -82.0516}, // Punta Gorda, FL
	"339":{26.7243, -81.8491}, // North Fort Myers, FL
	"339":{26.7404, -81.5241}, // Fort Denaud, FL
	"339":{26.6892, -81.8946}, // Palmona Park, FL
	"339":{26.9918, -82.1140}, // Port Charlotte, FL
	"339":{26.5160, -81.9601}, // Iona, FL
	"339":{26.4534, -82.1023}, // Sanibel, FL
	"339":{26.6643, -82.1477}, // Pineland, FL
	"339":{26.5804, -81.7453}, // Gateway, FL
	"339":{26.7645, -82.0507}, // Burnt Store Marina, FL
	"339":{26.9629, -82.0571}, // Charlotte Harbor, FL
	"339":{26.6437, -81.9098}, // Lochmoor Waterway Estates, FL
	"339":{26.5727, -81.8775}, // Pine Manor, FL
	"339":{26.7493, -81.3876}, // Port LaBelle, FL
	"339":{26.6800, -82.1409}, // Bokeelia, FL
	"339":{26.5160, -81.9293}, // Harlem Heights, FL
	"339":{26.6758, -81.8171}, // Tice, FL
	"339":{26.8845, -82.2791}, // Rotonda, FL
	"339":{26.9529, -81.9925}, // Cleveland, FL
	"339":{26.6340, -82.0589}, // Matlacha Isles-Matlacha Shores, FL
	"339":{26.3559, -81.7861}, // Bonita Springs, FL
	"339":{26.5782, -81.8615}, // Page Park, FL
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <340> with city APO/FPO and state AA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <341> with city FT MYERS and state FL
	"341":{26.0889, -81.7031}, // Lely Resort, FL
	"341":{26.2633, -81.8094}, // Naples Park, FL
	"341":{26.2510, -81.7110}, // Island Walk, FL
	"341":{26.4253, -81.4251}, // Immokalee, FL
	"341":{26.1032, -81.7297}, // Lely, FL
	"341":{26.4277, -81.7951}, // Estero, FL
	"341":{26.1844, -81.7031}, // Golden Gate, FL
	"341":{25.9330, -81.6993}, // Marco Island, FL
	"341":{25.8579, -81.3862}, // Everglades, FL
	"341":{26.2930, -81.5786}, // Orangetree, FL
	"341":{26.2279, -81.7280}, // Vineyards, FL
	"341":{26.2326, -81.8108}, // Pelican Bay, FL
	"341":{25.8146, -81.3603}, // Chokoloskee, FL
	"341":{26.1505, -81.7936}, // Naples, FL
	"341":{25.8475, -81.3754}, // Plantation Island, FL
	"341":{25.9252, -81.6480}, // Goodland, FL
	"341":{26.0892, -81.7254}, // Naples Manor, FL
	"341":{26.3559, -81.7861}, // Bonita Springs, FL
	"341":{26.0839, -81.6794}, // Verona Walk, FL
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <342> with city MANASOTA and state FL
	"342":{27.3328, -82.4616}, // Fruitville, FL
	"342":{27.2254, -82.4944}, // Vamo, FL
	"342":{27.1227, -82.4370}, // Nokomis, FL
	"342":{27.0694, -82.4054}, // Venice Gardens, FL
	"342":{27.3591, -82.4932}, // Kensington Park, FL
	"342":{27.2874, -82.5137}, // Ridge Wood Heights, FL
	"342":{27.2916, -82.4370}, // Lake Sarasota, FL
	"342":{26.9226, -82.3537}, // Manasota Key, FL
	"342":{27.2779, -82.5516}, // Siesta Key, FL
	"342":{27.5107, -82.7152}, // Holmes Beach, FL
	"342":{27.5251, -82.5751}, // Palmetto, FL
	"342":{27.1164, -82.4135}, // Venice, FL
	"342":{27.2856, -82.5333}, // South Sarasota, FL
	"342":{27.3082, -82.5096}, // Southgate, FL
	"342":{27.2214, -81.8587}, // Arcadia, FL
	"342":{27.5266, -82.5261}, // Ellenton, FL
	"342":{27.1914, -82.4800}, // Osprey, FL
	"342":{27.2587, -82.5065}, // Gulf Gate Estates, FL
	"342":{27.0435, -82.4152}, // South Venice, FL
	"342":{27.4668, -82.6688}, // Cortez, FL
	"342":{27.1862, -81.8521}, // Southeast Arcadia, FL
	"342":{27.0577, -82.1975}, // North Port, FL
	"342":{27.3926, -82.6341}, // Longboat Key, FL
	"342":{27.4612, -82.5821}, // South Bradenton, FL
	"342":{27.4116, -82.5659}, // Whitfield, FL
	"342":{26.9603, -82.3535}, // Englewood, FL
	"342":{26.9071, -82.3258}, // Grove City, FL
	"342":{27.3386, -82.5430}, // Sarasota, FL
	"342":{27.3092, -82.4788}, // Sarasota Springs, FL
	"342":{27.4703, -82.5552}, // West Samoset, FL
	"342":{27.2855, -82.4731}, // Bee Ridge, FL
	"342":{27.4752, -82.5430}, // Samoset, FL
	"342":{27.3711, -82.5177}, // North Sarasota, FL
	"342":{27.4345, -82.5794}, // Bayshore Gardens, FL
	"342":{27.1446, -82.4618}, // Laurel, FL
	"342":{27.2856, -82.4970}, // South Gate Ridge, FL
	"342":{27.5016, -82.6146}, // West Bradenton, FL
	"342":{27.3653, -82.4725}, // The Meadows, FL
	"342":{27.4900, -82.5740}, // Bradenton, FL
	"342":{27.3730, -82.4968}, // Desoto Lakes, FL
	"342":{27.5296, -82.7338}, // Anna Maria, FL
	"342":{27.0469, -82.2702}, // Warm Mineral Springs, FL
	"342":{27.4649, -82.6957}, // Bradenton Beach, FL
	"342":{27.5435, -82.5607}, // Memphis, FL
	// END OF OPTIONS TO PICK
	"344":{29.6804, -82.3458}, // Gainesville, FL
	"346":{27.9942, -82.4451}, // Tampa, FL
	"347":{28.4772, -81.3369}, // Orlando, FL
	TODO PLACEHOLDER TO FIX ZIP <349> with city WEST PALM BCH and state FL
	"349":{27.2161, -80.2400}, // Rio, FL
	"349":{27.4096, -80.3538}, // Fort Pierce South, FL
	"349":{27.1970, -80.1982}, // Sewall's Point, FL
	"349":{27.6993, -80.8871}, // Yeehaw Junction, FL
	"349":{27.2172, -80.7927}, // Taylor Creek, FL
	"349":{27.1461, -80.1895}, // Port Salerno, FL
	"349":{27.2479, -80.8115}, // Cypress Quarters, FL
	"349":{27.4256, -80.3430}, // Fort Pierce, FL
	"349":{27.3243, -80.2425}, // Hutchinson Island South, FL
	"349":{27.1330, -80.8868}, // Buckhead Ridge, FL
	"349":{27.2438, -80.2423}, // Jensen Beach, FL
	"349":{27.3722, -80.3403}, // White City, FL
	"349":{27.2223, -80.2738}, // North River Shores, FL
	"349":{27.2414, -80.8298}, // Okeechobee, FL
	"349":{27.4966, -80.3416}, // St. Lucie Village, FL
	"349":{27.0375, -80.4913}, // Indiantown, FL
	"349":{27.1958, -80.2438}, // Stuart, FL
	"349":{27.3564, -80.2984}, // Indian River Estates, FL
	"349":{27.4736, -80.3594}, // Fort Pierce North, FL
	"349":{27.2359, -80.2206}, // Ocean Breeze Park, FL
	"349":{27.2796, -80.3884}, // Port St. Lucie, FL
	"349":{27.3214, -80.3307}, // River Park, FL
	"349":{27.5390, -80.3865}, // Lakewood Park, FL
	"349":{27.1736, -80.2861}, // Palm City, FL
	// END OF OPTIONS TO PICK
	"350":{33.5277, -86.7987}, // Birmingham, AL
	"351":{33.5277, -86.7987}, // Birmingham, AL
	"352":{33.5277, -86.7987}, // Birmingham, AL
	"354":{33.2348, -87.5266}, // Tuscaloosa, AL
	"355":{33.5277, -86.7987}, // Birmingham, AL
	"356":{34.6988, -86.6412}, // Huntsville, AL
	"357":{34.6988, -86.6412}, // Huntsville, AL
	"358":{34.6988, -86.6412}, // Huntsville, AL
	"359":{33.5277, -86.7987}, // Birmingham, AL
	"360":{32.3473, -86.2666}, // Montgomery, AL
	"361":{32.3473, -86.2666}, // Montgomery, AL
	"362":{33.6713, -85.8136}, // Anniston, AL
	"363":{31.2335, -85.4068}, // Dothan, AL
	"364":{31.4342, -86.9723}, // Evergreen, AL
	"365":{30.6782, -88.1163}, // Mobile, AL
	"366":{30.6782, -88.1163}, // Mobile, AL
	"367":{32.3473, -86.2666}, // Montgomery, AL
	"368":{32.3473, -86.2666}, // Montgomery, AL
	"369":{32.3846, -88.6897}, // Meridian, MS
	"370":{36.1715, -86.7843}, // Nashville, TN
	"371":{36.1715, -86.7843}, // Nashville, TN
	"372":{36.1715, -86.7843}, // Nashville, TN
	"373":{35.0657, -85.2487}, // Chattanooga, TN
	"374":{35.0657, -85.2487}, // Chattanooga, TN
	"375":{35.1046, -89.9773}, // Memphis, TN
	"376":{36.3406, -82.3803}, // Johnson City, TN
	"377":{35.9692, -83.9496}, // Knoxville, TN
	"378":{35.9692, -83.9496}, // Knoxville, TN
	"379":{35.9692, -83.9496}, // Knoxville, TN
	"380":{35.1046, -89.9773}, // Memphis, TN
	"381":{35.1046, -89.9773}, // Memphis, TN
	"382":{36.1371, -88.5077}, // McKenzie, TN
	"383":{35.6536, -88.8353}, // Jackson, TN
	"384":{35.6236, -87.0487}, // Columbia, TN
	"385":{36.1484, -85.5114}, // Cookeville, TN
	"386":{35.1046, -89.9773}, // Memphis, TN
	"387":{33.3850, -91.0514}, // Greenville, MS
	"388":{34.2691, -88.7318}, // Tupelo, MS
	"389":{33.7816, -89.8130}, // Grenada, MS
	"390":{32.3163, -90.2124}, // Jackson, MS
	"391":{32.3163, -90.2124}, // Jackson, MS
	"392":{32.3163, -90.2124}, // Jackson, MS
	"393":{32.3846, -88.6897}, // Meridian, MS
	"394":{31.3074, -89.3170}, // Hattiesburg, MS
	"395":{30.4271, -89.0703}, // Gulfport, MS
	"396":{31.2449, -90.4714}, // McComb, MS
	"397":{33.5088, -88.4097}, // Columbus, MS
	"398":{31.5776, -84.1762}, // Albany, GA
	"399":{33.7627, -84.4225}, // Atlanta, GA
	"400":{38.1663, -85.6485}, // Louisville, KY
	"401":{38.1663, -85.6485}, // Louisville, KY
	"402":{38.1663, -85.6485}, // Louisville, KY
	"403":{38.0423, -84.4587}, // Lexington, KY
	"404":{38.0423, -84.4587}, // Lexington, KY
	"405":{38.0423, -84.4587}, // Lexington, KY
	"406":{38.1924, -84.8643}, // Frankfort, KY
	"407":{37.1209, -84.0804}, // London, KY
	"408":{37.1209, -84.0804}, // London, KY
	"409":{37.1209, -84.0804}, // London, KY
	"410":{39.1412, -84.5060}, // Cincinnati, OH
	"411":{38.4593, -82.6449}, // Ashland, KY
	"412":{38.4593, -82.6449}, // Ashland, KY
	"413":{37.7353, -83.5473}, // Campton, KY
	"414":{37.7353, -83.5473}, // Campton, KY
	"415":{37.4807, -82.5262}, // Pikeville, KY
	"416":{37.4807, -82.5262}, // Pikeville, KY
	"417":{37.2583, -83.1976}, // Hazard, KY
	"418":{37.2583, -83.1976}, // Hazard, KY
	"420":{37.0711, -88.6435}, // Paducah, KY
	"421":{36.9715, -86.4375}, // Bowling Green, KY
	"422":{36.9715, -86.4375}, // Bowling Green, KY
	"423":{37.7573, -87.1174}, // Owensboro, KY
	"424":{37.9881, -87.5341}, // Evansville, IN
	"425":{37.0816, -84.6089}, // Somerset, KY
	"426":{37.0816, -84.6089}, // Somerset, KY
	"427":{37.7030, -85.8769}, // Elizabethtown, KY
	"430":{39.9860, -82.9851}, // Columbus, OH
	"431":{39.9860, -82.9851}, // Columbus, OH
	"432":{39.9860, -82.9851}, // Columbus, OH
	"433":{39.9860, -82.9851}, // Columbus, OH
	"434":{41.6639, -83.5822}, // Toledo, OH
	"435":{41.6639, -83.5822}, // Toledo, OH
	"436":{41.6639, -83.5822}, // Toledo, OH
	"437":{39.9567, -82.0133}, // Zanesville, OH
	"438":{39.9567, -82.0133}, // Zanesville, OH
	"439":{40.3653, -80.6520}, // Steubenville, OH
	"440":{41.4767, -81.6805}, // Cleveland, OH
	"441":{41.4767, -81.6805}, // Cleveland, OH
	"442":{41.0798, -81.5219}, // Akron, OH
	"443":{41.0798, -81.5219}, // Akron, OH
	"444":{41.0993, -80.6463}, // Youngstown, OH
	"445":{41.0993, -80.6463}, // Youngstown, OH
	"446":{40.8076, -81.3678}, // Canton, OH
	"447":{40.8076, -81.3678}, // Canton, OH
	"448":{40.7656, -82.5275}, // Mansfield, OH
	"449":{40.7656, -82.5275}, // Mansfield, OH
	"450":{39.1412, -84.5060}, // Cincinnati, OH
	"451":{39.1412, -84.5060}, // Cincinnati, OH
	"452":{39.1412, -84.5060}, // Cincinnati, OH
	"453":{39.7797, -84.1998}, // Dayton, OH
	"454":{39.7797, -84.1998}, // Dayton, OH
	"455":{39.9297, -83.7957}, // Springfield, OH
	"456":{39.3393, -82.9937}, // Chillicothe, OH
	"457":{39.3269, -82.0987}, // Athens, OH
	"458":{40.7410, -84.1121}, // Lima, OH
	"459":{39.1412, -84.5060}, // Cincinnati, OH
	"460":{39.7771, -86.1458}, // Indianapolis, IN
	"461":{39.7771, -86.1458}, // Indianapolis, IN
	"462":{39.7771, -86.1458}, // Indianapolis, IN
	"463":{41.5906, -87.3472}, // Gary, IN
	"464":{41.5906, -87.3472}, // Gary, IN
	"465":{41.6771, -86.2692}, // South Bend, IN
	"466":{41.6771, -86.2692}, // South Bend, IN
	"467":{41.0885, -85.1436}, // Fort Wayne, IN
	"468":{41.0885, -85.1436}, // Fort Wayne, IN
	"469":{40.4640, -86.1277}, // Kokomo, IN
	"470":{39.1412, -84.5060}, // Cincinnati, OH
	"471":{38.1663, -85.6485}, // Louisville, KY
	"472":{39.2094, -85.9183}, // Columbus, IN
	"473":{40.1989, -85.3950}, // Muncie, IN
	"474":{39.1637, -86.5257}, // Bloomington, IN
	"475":{39.4654, -87.3763}, // Terre Haute, IN
	"476":{37.9881, -87.5341}, // Evansville, IN
	"477":{37.9881, -87.5341}, // Evansville, IN
	"478":{39.4654, -87.3763}, // Terre Haute, IN
	"479":{40.3990, -86.8593}, // Lafayette, IN
	"480":{42.5084, -83.1539}, // Royal Oak, MI
	"481":{42.3834, -83.1024}, // Detroit, MI
	"482":{42.3834, -83.1024}, // Detroit, MI
	"483":{42.5084, -83.1539}, // Royal Oak, MI
	"484":{43.0235, -83.6922}, // Flint, MI
	"485":{43.0235, -83.6922}, // Flint, MI
	"486":{43.4199, -83.9501}, // Saginaw, MI
	"487":{43.4199, -83.9501}, // Saginaw, MI
	"488":{42.7142, -84.5601}, // Lansing, MI
	"489":{42.7142, -84.5601}, // Lansing, MI
	"490":{42.2749, -85.5882}, // Kalamazoo, MI
	"491":{42.2749, -85.5882}, // Kalamazoo, MI
	"492":{42.2431, -84.4037}, // Jackson, MI
	"493":{42.9615, -85.6557}, // Grand Rapids, MI
	"494":{42.9615, -85.6557}, // Grand Rapids, MI
	"495":{42.9615, -85.6557}, // Grand Rapids, MI
	"496":{44.7547, -85.6035}, // Traverse City, MI
	"497":{45.0214, -84.6803}, // Gaylord, MI
	"498":{45.8275, -88.0599}, // Iron Mountain, MI
	"499":{45.8275, -88.0599}, // Iron Mountain, MI
	"500":{41.5725, -93.6105}, // Des Moines, IA
	"501":{41.5725, -93.6105}, // Des Moines, IA
	"502":{41.5725, -93.6105}, // Des Moines, IA
	"503":{41.5725, -93.6105}, // Des Moines, IA
	"504":{42.4920, -92.3522}, // Waterloo, IA
	"505":{42.5098, -94.1751}, // Fort Dodge, IA
	"506":{42.4920, -92.3522}, // Waterloo, IA
	"507":{42.4920, -92.3522}, // Waterloo, IA
	"508":{41.0597, -94.3650}, // Creston, IA
	"509":{41.5725, -93.6105}, // Des Moines, IA
	"510":{42.4959, -96.3901}, // Sioux City, IA
	"511":{42.4959, -96.3901}, // Sioux City, IA
	"512":{42.4959, -96.3901}, // Sioux City, IA
	"513":{42.4959, -96.3901}, // Sioux City, IA
	"514":{42.0699, -94.8647}, // Carroll, IA
	"515":{41.2628, -96.0498}, // Omaha, NE
	"516":{41.2628, -96.0498}, // Omaha, NE
	"520":{42.5007, -90.7067}, // Dubuque, IA
	"521":{43.3016, -91.7846}, // Decorah, IA
	"522":{41.9667, -91.6781}, // Cedar Rapids, IA
	"523":{41.9667, -91.6781}, // Cedar Rapids, IA
	"524":{41.9667, -91.6781}, // Cedar Rapids, IA
	"525":{41.5725, -93.6105}, // Des Moines, IA
	"526":{40.8072, -91.1247}, // Burlington, IA
	TODO PLACEHOLDER TO FIX ZIP <527> with city QUAD CITIES and state IL
	"527":{41.7399, -90.9738}, // Bennett, IA
	"527":{41.5898, -91.0271}, // Wilton, IA
	"527":{41.8227, -90.5448}, // DeWitt, IA
	"527":{41.6915, -90.6746}, // Donahue, IA
	"527":{41.5725, -91.2621}, // West Liberty, IA
	"527":{41.3796, -91.3486}, // Conesville, IA
	"527":{41.9029, -90.8628}, // Toronto, IA
	"527":{41.6744, -90.3573}, // Princeton, IA
	"527":{41.4195, -91.0680}, // Muscatine, IA
	"527":{41.6923, -90.5804}, // Long Grove, IA
	"527":{41.8265, -90.7600}, // Calamus, IA
	"527":{41.8329, -90.8379}, // Wheatland, IA
	"527":{41.6394, -90.5805}, // Eldridge, IA
	"527":{41.5371, -90.4637}, // Riverdale, IA
	"527":{41.9680, -90.3817}, // Goose Lake, IA
	"527":{41.5716, -91.1664}, // Atalissa, IA
	"527":{41.2772, -91.1884}, // Grandview, IA
	"527":{41.6922, -90.5393}, // Park View, IA
	"527":{41.6011, -90.9137}, // Durant, IA
	"527":{41.5908, -90.8566}, // Stockton, IA
	"527":{41.2940, -91.4650}, // Cotter, IA
	"527":{41.7928, -90.2766}, // Camanche, IA
	"527":{41.9076, -90.5959}, // Welton, IA
	"527":{41.5989, -90.7744}, // Walcott, IA
	"527":{41.7701, -91.1283}, // Tipton, IA
	"527":{41.8023, -90.3542}, // Low Moor, IA
	"527":{41.2787, -91.3651}, // Columbus Junction, IA
	"527":{41.2592, -91.3747}, // Columbus City, IA
	"527":{41.3475, -91.1288}, // Fruitland, IA
	"527":{41.5556, -90.4537}, // Panorama Park, IA
	"527":{41.5096, -90.7649}, // Blue Grass, IA
	"527":{41.9621, -90.4682}, // Charlotte, IA
	"527":{41.6747, -91.1508}, // Rochester, IA
	"527":{41.7159, -90.8780}, // New Liberty, IA
	"527":{41.7422, -90.7824}, // Dixon, IA
	"527":{41.8434, -90.2409}, // Clinton, IA
	"527":{41.7435, -90.4461}, // McCausland, IA
	"527":{41.4666, -90.7170}, // Buffalo, IA
	"527":{41.4797, -91.3081}, // Nichols, IA
	"527":{41.3299, -91.2356}, // Letts, IA
	"527":{41.6496, -90.7184}, // Maysville, IA
	"527":{41.5656, -90.4764}, // Bettendorf, IA
	"527":{41.5964, -90.3687}, // Le Claire, IA
	"527":{41.4859, -91.4267}, // Lone Tree, IA
	"527":{41.2845, -91.3397}, // Fredonia, IA
	"527":{41.8232, -90.6502}, // Grand Mound, IA
	"527":{41.9799, -90.2529}, // Andover, IA
	// END OF OPTIONS TO PICK
	"528":{41.5563, -90.6052}, // Davenport, IA
	"530":{43.0642, -87.9673}, // Milwaukee, WI
	"531":{43.0642, -87.9673}, // Milwaukee, WI
	"532":{43.0642, -87.9673}, // Milwaukee, WI
	"534":{42.7274, -87.8135}, // Racine, WI
	"535":{43.0827, -89.3923}, // Madison, WI
	"537":{43.0827, -89.3923}, // Madison, WI
	"538":{43.0827, -89.3923}, // Madison, WI
	"539":{43.5489, -89.4658}, // Portage, WI
	TODO PLACEHOLDER TO FIX ZIP <540> with city ST PAUL and state MN
	"540":{44.8608, -92.6244}, // River Falls, WI
	"540":{44.9691, -92.4381}, // Hammond, WI
	"540":{45.1229, -92.5339}, // New Richmond, WI
	"540":{44.6012, -92.5336}, // Hager City, WI
	"540":{45.3045, -92.3635}, // Amery, WI
	"540":{45.3194, -92.6935}, // Osceola, WI
	"540":{45.3260, -92.1710}, // Clayton, WI
	"540":{45.1269, -92.6756}, // Somerset, WI
	"540":{44.9976, -92.7559}, // North Hudson, WI
	"540":{45.1975, -92.5321}, // Star Prairie, WI
	"540":{44.7522, -92.7882}, // Prescott, WI
	"540":{44.9586, -92.1705}, // Wilson, WI
	"540":{45.0570, -92.1714}, // Glenwood City, WI
	"540":{44.9540, -92.3709}, // Baldwin, WI
	"540":{44.6489, -92.6189}, // Diamond Bluff, WI
	"540":{45.4101, -92.6268}, // St. Croix Falls, WI
	"540":{44.9639, -92.7316}, // Hudson, WI
	"540":{45.1887, -92.3879}, // Deer Park, WI
	"540":{45.0833, -92.2579}, // Emerald, WI
	"540":{44.9732, -92.5505}, // Roberts, WI
	"540":{45.2492, -92.2675}, // Clear Lake, WI
	"540":{44.7364, -92.4806}, // Ellsworth, WI
	"540":{45.0607, -92.7904}, // Houlton, WI
	"540":{44.9483, -92.2851}, // Woodville, WI
	"540":{45.3617, -92.6344}, // Dresser, WI
	// END OF OPTIONS TO PICK
	"541":{44.5150, -87.9896}, // Green Bay, WI
	"542":{44.5150, -87.9896}, // Green Bay, WI
	"543":{44.5150, -87.9896}, // Green Bay, WI
	"544":{44.9615, -89.6457}, // Wausau, WI
	"545":{45.6360, -89.4256}, // Rhinelander, WI
	"546":{43.8241, -91.2268}, // La Crosse, WI
	"547":{44.8200, -91.4951}, // Eau Claire, WI
	"548":{45.8271, -91.8860}, // Spooner, WI
	"549":{44.0228, -88.5617}, // Oshkosh, WI
	TODO PLACEHOLDER TO FIX ZIP <550> with city ST PAUL and state MN
	"550":{45.3824, -93.0886}, // Martin Lake, MN
	"550":{44.4551, -93.1697}, // Northfield, MN
	"550":{45.8367, -92.9683}, // Pine City, MN
	"550":{44.8744, -92.9975}, // Newport, MN
	"550":{45.7957, -93.1524}, // Grasston, MN
	"550":{45.2536, -92.9583}, // Forest Lake, MN
	"550":{45.6875, -92.9654}, // Rush City, MN
	"550":{44.2240, -93.4445}, // Morristown, MN
	"550":{44.9944, -92.9031}, // Lake Elmo, MN
	"550":{45.7223, -93.1717}, // Braham, MN
	"550":{44.8360, -92.9949}, // St. Paul Park, MN
	"550":{45.1381, -93.1714}, // Lexington, MN
	"550":{45.0573, -92.8313}, // Stillwater, MN
	"550":{44.8877, -93.0411}, // South St. Paul, MN
	"550":{45.2685, -93.0810}, // Columbus, MN
	"550":{44.0914, -93.2304}, // Owatonna, MN
	"550":{45.8766, -93.2916}, // Mora, MN
	"550":{45.3941, -92.8141}, // Center City, MN
	"550":{45.9484, -93.0728}, // Brook Park, MN
	"550":{45.8713, -93.1197}, // Henriette, MN
	"550":{44.4088, -93.0303}, // Dennison, MN
	"550":{44.3429, -93.0638}, // Nerstrand, MN
	"550":{45.1679, -93.0830}, // Lino Lakes, MN
	"550":{46.1292, -92.8646}, // Sandstone, MN
	"550":{45.4022, -93.2711}, // Bethel, MN
	"550":{45.3474, -92.9116}, // Chisago City, MN
	"550":{44.8161, -92.9274}, // Cottage Grove, MN
	"550":{45.4929, -93.2414}, // Isanti, MN
	"550":{44.2574, -93.3774}, // Warsaw, MN
	"550":{45.2540, -92.8278}, // Scandia, MN
	"550":{44.6776, -93.2521}, // Lakeville, MN
	"550":{45.1640, -93.0540}, // Centerville, MN
	"550":{44.2985, -93.2787}, // Faribault, MN
	"550":{45.5976, -92.9869}, // Harris, MN
	"550":{45.6650, -93.1815}, // Stanchfield, MN
	"550":{45.5612, -93.2283}, // Cambridge, MN
	"550":{44.4007, -92.6255}, // Goodhue, MN
	"550":{45.3871, -92.8477}, // Lindstrom, MN
	"550":{44.5260, -93.0198}, // Randolph, MN
	"550":{45.3901, -92.7531}, // Shafer, MN
	"550":{44.6572, -93.1687}, // Farmington, MN
	"550":{45.3365, -92.9766}, // Wyoming, MN
	"550":{45.0540, -92.9571}, // Willernie, MN
	"550":{44.5667, -93.3381}, // Elko New Market, MN
	"550":{44.6040, -92.9328}, // New Trier, MN
	"550":{45.1671, -92.9588}, // Hugo, MN
	"550":{45.3409, -93.3264}, // Oak Grove, MN
	"550":{45.3991, -93.3902}, // St. Francis, MN
	"550":{44.1685, -93.2472}, // Medford, MN
	"550":{44.9245, -92.7657}, // Lake St. Croix Beach, MN
	"550":{46.0122, -92.9256}, // Hinckley, MN
	"550":{45.1409, -93.1505}, // Circle Pines, MN
	"550":{44.4277, -93.2039}, // Dundas, MN
	"550":{44.5123, -92.9034}, // Cannon Falls, MN
	"550":{44.9478, -92.7612}, // Lakeland Shores, MN
	"550":{45.4121, -92.6644}, // Taylors Falls, MN
	"550":{44.3710, -92.5113}, // Bellechester, MN
	"550":{44.7465, -93.0662}, // Rosemount, MN
	"550":{45.0825, -92.9090}, // Grant, MN
	"550":{45.3841, -92.9938}, // Stacy, MN
	"550":{44.7318, -92.8538}, // Hastings, MN
	"550":{45.3557, -93.2038}, // East Bethel, MN
	"550":{45.7604, -92.9088}, // Rock Creek, MN
	"550":{44.4776, -93.4220}, // Lonsdale, MN
	"550":{44.6085, -93.0025}, // Hampton, MN
	"550":{44.8247, -93.0596}, // Inver Grove Heights, MN
	"550":{44.9042, -92.8173}, // Afton, MN
	"550":{44.7160, -93.0323}, // Coates, MN
	"550":{44.4453, -92.2796}, // Lake City, MN
	"550":{45.1696, -93.2078}, // Blaine, MN
	"550":{44.5816, -92.6036}, // Red Wing, MN
	"550":{44.5061, -92.3511}, // Frontenac, MN
	"550":{45.1983, -92.7783}, // Marine on St. Croix, MN
	"550":{44.5990, -92.8196}, // Miesville, MN
	"550":{45.0324, -92.8099}, // Oak Park Heights, MN
	"550":{44.6735, -92.9643}, // Vermillion, MN
	"550":{44.9129, -92.7703}, // St. Marys Point, MN
	"550":{45.5137, -92.9601}, // North Branch, MN
	"550":{45.9155, -93.1753}, // Quamba, MN
	"550":{44.9503, -92.7700}, // Lakeland, MN
	"550":{45.0152, -92.7789}, // Bayport, MN
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <551> with city ST PAUL and state MN
	"551":{45.0598, -92.9777}, // Birchwood Village, MN
	"551":{44.9840, -93.0247}, // Maplewood, MN
	"551":{44.9477, -93.1040}, // St. Paul, MN
	"551":{44.9096, -93.1301}, // Lilydale, MN
	"551":{44.8722, -93.0943}, // Sunfish Lake, MN
	"551":{45.0658, -93.2061}, // New Brighton, MN
	"551":{44.9899, -93.1770}, // Falcon Heights, MN
	"551":{44.9942, -93.2026}, // Lauderdale, MN
	"551":{44.9877, -92.9641}, // Oakdale, MN
	"551":{45.0986, -92.9672}, // Dellwood, MN
	"551":{44.9506, -92.9769}, // Landfall, MN
	"551":{45.0155, -93.1544}, // Roseville, MN
	"551":{45.0579, -93.0405}, // Gem Lake, MN
	"551":{45.0137, -92.9995}, // North St. Paul, MN
	"551":{44.8876, -93.1609}, // Mendota, MN
	"551":{44.9056, -92.9230}, // Woodbury, MN
	"551":{45.0842, -93.1358}, // Shoreview, MN
	"551":{45.1072, -93.2078}, // Mounds View, MN
	"551":{44.9018, -93.0858}, // West St. Paul, MN
	"551":{45.1002, -93.0881}, // North Oaks, MN
	"551":{45.0540, -92.9571}, // Willernie, MN
	"551":{45.1671, -92.9588}, // Hugo, MN
	"551":{45.0825, -92.9090}, // Grant, MN
	"551":{45.0570, -93.0747}, // Vadnais Heights, MN
	"551":{45.0368, -92.9543}, // Pine Springs, MN
	"551":{44.8815, -93.1400}, // Mendota Heights, MN
	"551":{45.0655, -93.0150}, // White Bear Lake, MN
	"551":{45.0244, -93.0863}, // Little Canada, MN
	"551":{44.8169, -93.1638}, // Eagan, MN
	"551":{44.7457, -93.2006}, // Apple Valley, MN
	"551":{45.0619, -92.9660}, // Mahtomedi, MN
	"551":{45.0722, -93.1671}, // Arden Hills, MN
	// END OF OPTIONS TO PICK
	"553":{44.9635, -93.2678}, // Minneapolis, MN
	"554":{44.9635, -93.2678}, // Minneapolis, MN
	"555":{44.9635, -93.2678}, // Minneapolis, MN
	"556":{46.7757, -92.1392}, // Duluth, MN
	"557":{46.7757, -92.1392}, // Duluth, MN
	"558":{46.7757, -92.1392}, // Duluth, MN
	"559":{44.0151, -92.4778}, // Rochester, MN
	"560":{44.1711, -93.9773}, // Mankato, MN
	"561":{44.1711, -93.9773}, // Mankato, MN
	"562":{45.1220, -95.0569}, // Willmar, MN
	TODO PLACEHOLDER TO FIX ZIP <563> with city ST CLOUD and state MN
	"563":{45.4054, -94.8397}, // Regal, MN
	"563":{45.9772, -94.1008}, // Pierz, MN
	"563":{45.4650, -94.3222}, // Rockville, MN
	"563":{45.9741, -95.2925}, // Carlos, MN
	"563":{46.1245, -95.5110}, // Urbank, MN
	"563":{45.7010, -94.2742}, // St. Stephen, MN
	"563":{45.7525, -94.2317}, // Rice, MN
	"563":{45.8191, -94.4072}, // Bowlus, MN
	"563":{45.6309, -94.7525}, // New Munich, MN
	"563":{45.6148, -95.7382}, // Cyrus, MN
	"563":{45.7146, -95.1680}, // Westport, MN
	"563":{46.1530, -95.3296}, // Parkers Prairie, MN
	"563":{45.8649, -95.1524}, // Osakis, MN
	"563":{45.8245, -94.7491}, // Grey Eagle, MN
	"563":{46.0063, -93.8887}, // Hillman, MN
	"563":{46.1206, -93.5223}, // Wahkon, MN
	"563":{46.1739, -95.9156}, // Dalton, MN
	"563":{45.8340, -94.4942}, // Elmdale, MN
	"563":{45.7844, -93.5528}, // Bock, MN
	"563":{45.4619, -94.7965}, // Lake Henry, MN
	"563":{45.7050, -95.5190}, // Lowry, MN
	"563":{46.0462, -95.2931}, // Miltona, MN
	"563":{45.9168, -94.6383}, // Swanville, MN
	"563":{45.5313, -94.2528}, // Waite Park, MN
	"563":{46.0931, -95.8157}, // Ashby, MN
	"563":{45.6263, -94.8693}, // Meire Grove, MN
	"563":{45.7573, -93.6522}, // Milaca, MN
	"563":{45.9230, -94.4820}, // Sobieski, MN
	"563":{46.1398, -93.4597}, // Isle, MN
	"563":{46.1832, -93.7823}, // Vineland, MN
	"563":{45.7141, -95.2692}, // Villard, MN
	"563":{45.6517, -95.3643}, // Glenwood, MN
	"563":{45.7884, -95.3567}, // Forada, MN
	"563":{46.0700, -93.6676}, // Onamia, MN
	"563":{45.8302, -94.2931}, // Royalton, MN
	"563":{45.8865, -95.2648}, // Nelson, MN
	"563":{45.8777, -95.3766}, // Alexandria, MN
	"563":{45.8341, -95.7870}, // Hoffman, MN
	"563":{45.6509, -95.4294}, // Long Beach, MN
	"563":{45.9835, -94.3599}, // Little Falls, MN
	"563":{45.6014, -94.8594}, // Greenwald, MN
	"563":{45.8100, -94.5673}, // Upsala, MN
	"563":{45.8658, -94.6874}, // Burtrum, MN
	"563":{45.7782, -95.6970}, // Kensington, MN
	"563":{45.9387, -95.4942}, // Garfield, MN
	"563":{45.6995, -93.6490}, // Pease, MN
	"563":{45.6186, -94.2205}, // Sartell, MN
	"563":{45.4497, -94.1996}, // St. Augusta, MN
	"563":{45.4570, -94.4300}, // Cold Spring, MN
	"563":{45.6284, -94.5674}, // Albany, MN
	"563":{45.7982, -95.0885}, // West Union, MN
	"563":{45.9663, -95.5944}, // Brandon, MN
	"563":{45.5028, -94.6677}, // St. Martin, MN
	"563":{45.5627, -94.9474}, // Elrosa, MN
	"563":{45.6853, -93.8601}, // Ronneby, MN
	"563":{45.6636, -93.9095}, // Foley, MN
	"563":{45.7298, -94.4717}, // Holdingford, MN
	"563":{45.4324, -94.6365}, // Roscoe, MN
	"563":{45.9775, -94.8629}, // Long Prairie, MN
	"563":{45.4556, -94.5153}, // Richmond, MN
	"563":{45.7321, -93.7095}, // Foreston, MN
	"563":{45.5003, -95.1162}, // Brooten, MN
	"563":{45.7524, -95.6189}, // Farwell, MN
	"563":{46.1191, -94.0368}, // Harding, MN
	"563":{46.2421, -93.2751}, // McGrath, MN
	"563":{45.9653, -94.1129}, // Genola, MN
	"563":{45.4507, -94.9998}, // Belgrade, MN
	"563":{45.9480, -94.5301}, // Flensburg, MN
	"563":{46.0062, -95.6869}, // Evansville, MN
	"563":{45.7286, -94.7162}, // St. Rosa, MN
	"563":{45.5782, -95.2451}, // Sedan, MN
	"563":{45.7353, -93.9483}, // Gilman, MN
	"563":{45.5233, -94.8317}, // Spring Hill, MN
	"563":{46.0398, -94.0621}, // Lastrup, MN
	"563":{45.5339, -94.1718}, // St. Cloud, MN
	"563":{45.6758, -94.8129}, // Melrose, MN
	"563":{45.6094, -94.4609}, // Avon, MN
	"563":{45.9124, -95.8943}, // Barrett, MN
	"563":{45.7360, -94.9523}, // Sauk Centre, MN
	"563":{46.0691, -95.5571}, // Millerville, MN
	"563":{45.6627, -94.6888}, // Freeport, MN
	"563":{45.5980, -94.1540}, // Sauk Rapids, MN
	"563":{45.3785, -94.7217}, // Paynesville, MN
	"563":{45.5609, -94.3084}, // St. Joseph, MN
	"563":{45.8291, -93.4221}, // Ogilvie, MN
	"563":{45.8973, -94.0940}, // Buckman, MN
	"563":{45.6121, -95.5333}, // Starbuck, MN
	// END OF OPTIONS TO PICK
	"564":{46.3553, -94.1983}, // Brainerd, MN
	"565":{46.8060, -95.8449}, // Detroit Lakes, MN
	"566":{47.4830, -94.8788}, // Bemidji, MN
	"567":{47.9221, -97.0887}, // Grand Forks, ND
	"570":{43.5397, -96.7321}, // Sioux Falls, SD
	"571":{43.5397, -96.7321}, // Sioux Falls, SD
	TODO PLACEHOLDER TO FIX ZIP <572> with city DAKOTA CENTRAL and state SD
	"572":{45.4965, -97.4932}, // Roslyn, SD
	"572":{45.8557, -96.9173}, // New Effington, SD
	"572":{45.3349, -97.2050}, // Ortley, SD
	"572":{44.9987, -96.9742}, // Waverly, SD
	"572":{44.5899, -97.4673}, // Bryant, SD
	"572":{45.8662, -96.7318}, // Rosholt, SD
	"572":{45.3346, -97.3054}, // Waubay, SD
	"572":{45.0441, -96.7604}, // Strandburg, SD
	"572":{45.1070, -97.7806}, // Crocker, SD
	"572":{44.3862, -97.5498}, // De Smet, SD
	"572":{44.8813, -97.7346}, // Clark, SD
	"572":{45.0548, -97.3263}, // Florence, SD
	"572":{45.7248, -97.4139}, // Lake City, SD
	"572":{45.2620, -97.7117}, // Butler, SD
	"572":{44.5577, -96.5465}, // Astoria, SD
	"572":{45.6169, -97.4198}, // Eden, SD
	"572":{44.4880, -97.4409}, // Erwin, SD
	"572":{44.8809, -97.4616}, // Henry, SD
	"572":{45.3366, -97.5219}, // Webster, SD
	"572":{44.3274, -96.6424}, // Bushnell, SD
	"572":{45.3352, -96.7646}, // Corona, SD
	"572":{44.4855, -97.2079}, // Badger, SD
	"572":{45.0845, -97.4781}, // Wallace, SD
	"572":{45.6823, -97.1611}, // Long Hollow, SD
	"572":{45.1018, -96.7998}, // Stockholm, SD
	"572":{44.5768, -96.9011}, // Estelline, SD
	"572":{45.0904, -97.6420}, // Bradley, SD
	"572":{44.8775, -96.8497}, // Goodwin, SD
	"572":{45.6625, -97.0453}, // Sisseton, SD
	"572":{45.1016, -96.9300}, // South Shore, SD
	"572":{44.3635, -97.1338}, // Arlington, SD
	"572":{44.4359, -96.6460}, // White, SD
	"572":{45.2194, -96.6340}, // Milbank, SD
	"572":{44.5724, -96.6414}, // Toronto, SD
	"572":{44.7646, -96.6767}, // Clear Lake, SD
	"572":{45.4091, -96.8545}, // Wilmot, SD
	"572":{45.1818, -97.6835}, // Lily, SD
	"572":{44.3617, -97.3760}, // Lake Preston, SD
	"572":{44.4378, -96.8906}, // Bruce, SD
	"572":{45.5808, -97.0831}, // Agency Village, SD
	"572":{45.8620, -97.2869}, // Veblen, SD
	"572":{45.2985, -96.4657}, // Big Stone City, SD
	"572":{45.2605, -96.9154}, // Marvin, SD
	"572":{45.0490, -96.6752}, // La Bolt, SD
	"572":{45.3466, -97.7485}, // Bristol, SD
	"572":{44.7716, -97.5131}, // Naples, SD
	"572":{45.2090, -96.7854}, // Twin Brooks, SD
	"572":{44.9592, -97.5806}, // Garden City, SD
	"572":{45.0155, -96.5713}, // Revillo, SD
	"572":{44.9105, -97.9372}, // Raymond, SD
	"572":{44.8408, -96.6901}, // Altamont, SD
	"572":{45.5422, -96.9560}, // Peever, SD
	"572":{44.8888, -96.9092}, // Kranzburg, SD
	"572":{44.7582, -97.3810}, // Hazel, SD
	"572":{44.6571, -97.2046}, // Hayti, SD
	"572":{44.5664, -97.0798}, // Lake Poinsett, SD
	"572":{45.0513, -96.5535}, // Albee, SD
	"572":{44.7941, -96.4574}, // Gary, SD
	"572":{44.6669, -96.6252}, // Brandt, SD
	"572":{44.7032, -97.5000}, // Vienna, SD
	"572":{45.4667, -97.3900}, // Grenville, SD
	"572":{44.3771, -97.2345}, // Hetland, SD
	"572":{45.8566, -97.1035}, // Claire City, SD
	"572":{45.5673, -97.0612}, // Goodwill, SD
	"572":{44.7242, -97.0309}, // Castlewood, SD
	"572":{45.9248, -96.5760}, // White Rock, SD
	"572":{44.5799, -97.2108}, // Lake Norden, SD
	"572":{45.3056, -97.0391}, // Summit, SD
	"572":{44.6279, -97.6387}, // Willow Lake, SD
	"572":{44.9094, -97.1532}, // Watertown, SD
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <573> with city DAKOTA CENTRAL and state SD
	"573":{43.1559, -98.5362}, // Lake Andes, SD
	"573":{43.5272, -98.5894}, // Aurora Center, SD
	"573":{44.4338, -97.9885}, // Yale, SD
	"573":{43.0281, -98.8891}, // Fairfax, SD
	"573":{44.5205, -98.9869}, // Miller, SD
	"573":{44.4552, -98.6967}, // Wessington, SD
	"573":{44.1686, -97.7150}, // Carthage, SD
	"573":{43.4252, -97.8081}, // Milltown, SD
	"573":{44.1829, -98.3683}, // Alpena, SD
	"573":{43.7792, -99.1843}, // Pukwana, SD
	"573":{43.6541, -97.7800}, // Alexandria, SD
	"573":{44.5158, -99.2005}, // Ree Heights, SD
	"573":{43.7288, -97.8226}, // Fulton, SD
	"573":{43.7127, -98.2609}, // Mount Vernon, SD
	"573":{43.3929, -97.9864}, // Parkston, SD
	"573":{44.3662, -98.1835}, // Morningside, SD
	"573":{43.0775, -98.9467}, // Bonesteel, SD
	"573":{43.7252, -97.6890}, // Farmer, SD
	"573":{43.8625, -98.3526}, // Storla, SD
	"573":{43.7155, -98.4837}, // Plankinton, SD
	"573":{43.2666, -98.1592}, // Delmont, SD
	"573":{43.8038, -99.3721}, // Oacoma, SD
	"573":{43.1930, -97.8434}, // Kaylor, SD
	"573":{44.5166, -98.9379}, // St. Lawrence, SD
	"573":{43.0671, -98.5295}, // Pickstown, SD
	"573":{44.0541, -98.2724}, // Woonsocket, SD
	"573":{43.4287, -98.6073}, // New Holland, SD
	"573":{43.5894, -98.4376}, // Stickney, SD
	"573":{43.8995, -98.1442}, // Letcher, SD
	"573":{43.5463, -97.9832}, // Ethan, SD
	"573":{43.7861, -99.3269}, // Chamberlain, SD
	"573":{43.3868, -98.8437}, // Platte, SD
	"573":{43.0399, -98.1856}, // Dante, SD
	"573":{44.0335, -98.9886}, // Gann Valley, SD
	"573":{43.5502, -97.4984}, // Bridgewater, SD
	"573":{43.0052, -98.0592}, // Avon, SD
	"573":{42.9965, -98.4302}, // Marty, SD
	"573":{43.1365, -98.4268}, // Ravinia, SD
	"573":{44.0079, -97.9237}, // Artesian, SD
	"573":{44.0117, -97.5247}, // Howard, SD
	"573":{44.0060, -97.6975}, // Roswell, SD
	"573":{43.2540, -98.6974}, // Geddes, SD
	"573":{44.4105, -98.4738}, // Wolsey, SD
	"573":{44.6294, -98.4082}, // Hitchcock, SD
	"573":{43.7296, -98.0337}, // Mitchell, SD
	"573":{43.7468, -98.9569}, // Kimball, SD
	"573":{43.7931, -98.1036}, // Loomis, SD
	"573":{44.0067, -97.7895}, // Fedora, SD
	"573":{43.5285, -99.1439}, // Bijou Hills, SD
	"573":{44.0083, -97.5960}, // Vilas, SD
	"573":{44.0805, -98.5712}, // Wessington Springs, SD
	"573":{44.3693, -98.0421}, // Cavour, SD
	"573":{43.7273, -97.5911}, // Spencer, SD
	"573":{43.3190, -98.3448}, // Armour, SD
	"573":{43.0768, -98.2934}, // Wagner, SD
	"573":{44.4894, -97.7505}, // Bancroft, SD
	"573":{43.4911, -97.3848}, // Dolton, SD
	"573":{43.4311, -98.5267}, // Harrison, SD
	"573":{44.5214, -99.4393}, // Highmore, SD
	"573":{43.7283, -98.7117}, // White Lake, SD
	"573":{44.3679, -97.8495}, // Iroquois, SD
	"573":{43.6005, -99.2114}, // Ola, SD
	"573":{43.8813, -97.5040}, // Canova, SD
	"573":{44.4938, -98.3466}, // Broadland, SD
	"573":{44.0527, -99.4079}, // Fort Thompson, SD
	"573":{43.4223, -98.4073}, // Corsica, SD
	"573":{44.2909, -98.4276}, // Virgil, SD
	"573":{44.0213, -98.1032}, // Forestburg, SD
	"573":{43.2255, -97.9658}, // Tripp, SD
	"573":{43.4763, -97.9884}, // Dimock, SD
	"573":{44.3622, -98.2102}, // Huron, SD
	"573":{44.0697, -98.4246}, // Lane, SD
	"573":{43.6020, -97.6195}, // Emery, SD
	// END OF OPTIONS TO PICK
	"574":{45.4646, -98.4680},  // Aberdeen, SD
	"575":{44.3748, -100.3205}, // Pierre, SD
	"576":{45.5411, -100.4349}, // Mobridge, SD
	"577":{44.0716, -103.2205}, // Rapid City, SD
	"580":{46.8653, -96.8292},  // Fargo, ND
	"581":{46.8653, -96.8292},  // Fargo, ND
	"582":{47.9221, -97.0887},  // Grand Forks, ND
	"583":{48.1131, -98.8753},  // Devils Lake, ND
	"584":{46.9063, -98.6937},  // Jamestown, ND
	"585":{46.8140, -100.7695}, // Bismarck, ND
	"586":{46.8140, -100.7695}, // Bismarck, ND
	"587":{48.2374, -101.2780}, // Minot, ND
	"588":{48.1814, -103.6364}, // Williston, ND
	"590":{45.7889, -108.5509}, // Billings, MT
	"591":{45.7889, -108.5509}, // Billings, MT
	"592":{48.0933, -105.6413}, // Wolf Point, MT
	"593":{46.4059, -105.8385}, // Miles City, MT
	"594":{47.5022, -111.2995}, // Great Falls, MT
	"595":{48.5427, -109.6804}, // Havre, MT
	"596":{46.5965, -112.0199}, // Helena, MT
	"597":{45.9020, -112.6571}, // Butte, MT
	"598":{46.8685, -114.0095}, // Missoula, MT
	"599":{48.2156, -114.3261}, // Kalispell, MT
	"600":{42.1181, -88.0430},  // Palatine, IL
	"601":{41.9182, -88.1308},  // Carol Stream, IL
	"602":{42.0463, -87.6942},  // Evanston, IL
	"603":{41.8872, -87.7899},  // Oak Park, IL
	TODO PLACEHOLDER TO FIX ZIP <604> with city S SUBURBAN and state IL
	"604":{41.5170, -87.6924}, // Olympia Fields, IL
	"604":{41.4399, -87.6231}, // Crete, IL
	"604":{41.0021, -88.5224}, // Odell, IL
	"604":{41.6431, -87.7080}, // Robbins, IL
	"604":{41.6877, -87.7916}, // Worth, IL
	"604":{41.4211, -88.2593}, // Channahon, IL
	"604":{41.7139, -87.7528}, // Oak Lawn, IL
	"604":{41.5614, -88.0587}, // Fairmont, IL
	"604":{41.4119, -88.1346}, // Elwood, IL
	"604":{41.5638, -87.7250}, // Country Club Hills, IL
	"604":{41.7248, -87.8281}, // Hickory Hills, IL
	"604":{41.5711, -87.6187}, // Thornton, IL
	"604":{41.5058, -88.1183}, // Rockdale, IL
	"604":{41.2763, -88.2805}, // Coal City, IL
	"604":{41.6637, -87.7959}, // Palos Heights, IL
	"604":{41.5670, -87.8050}, // Tinley Park, IL
	"604":{41.6291, -87.6858}, // Posen, IL
	"604":{41.3504, -87.6171}, // Beecher, IL
	"604":{41.1744, -88.2813}, // South Wilmington, IL
	"604":{41.3312, -87.7937}, // Peotone, IL
	"604":{41.2157, -88.5049}, // Verona, IL
	"604":{41.7444, -87.7686}, // Burbank, IL
	"604":{41.6454, -87.7397}, // Crestwood, IL
	"604":{41.2713, -88.1361}, // Lakewood Shores, IL
	"604":{41.3210, -88.5880}, // Seneca, IL
	"604":{41.6134, -87.5505}, // Calumet City, IL
	"604":{41.6330, -87.6672}, // Dixmoor, IL
	"604":{41.6000, -87.6905}, // Markham, IL
	"604":{41.4723, -87.6183}, // Steger, IL
	"604":{41.7402, -87.8067}, // Bridgeview, IL
	"604":{41.5391, -87.6857}, // Flossmoor, IL
	"604":{41.6903, -88.1019}, // Bolingbrook, IL
	"604":{41.7663, -87.8607}, // Hodgkins, IL
	"604":{41.2385, -88.2452}, // Godley, IL
	"604":{41.5648, -87.5462}, // Lansing, IL
	"604":{41.5097, -87.9700}, // New Lenox, IL
	"604":{41.8073, -87.7798}, // Forest View, IL
	"604":{41.6076, -87.6521}, // Harvey, IL
	"604":{41.4831, -87.6370}, // South Chicago Heights, IL
	"604":{41.6986, -87.8266}, // Palos Hills, IL
	"604":{41.5724, -88.1128}, // Crest Hill, IL
	"604":{41.4176, -87.7510}, // Monee, IL
	"604":{41.5109, -87.5814}, // Ford Heights, IL
	"604":{41.5412, -87.6119}, // Glenwood, IL
	"604":{41.5904, -88.0293}, // Lockport, IL
	"604":{41.6118, -87.6308}, // Phoenix, IL
	"604":{41.1903, -88.5699}, // Kinsman, IL
	"604":{41.4907, -87.5702}, // Sauk Village, IL
	"604":{41.5977, -87.6022}, // South Holland, IL
	"604":{41.5234, -87.5508}, // Lynwood, IL
	"604":{41.4822, -87.7352}, // Richton Park, IL
	"604":{41.2280, -88.2464}, // Braceville, IL
	"604":{41.2852, -88.2506}, // Diamond, IL
	"604":{41.5905, -87.8413}, // Orland Hills, IL
	"604":{41.6697, -87.9828}, // Lemont, IL
	"604":{41.6578, -87.6812}, // Blue Island, IL
	"604":{41.6055, -87.7527}, // Oak Forest, IL
	"604":{41.8373, -87.6862}, // Chicago, IL
	"604":{41.4507, -88.2792}, // Minooka, IL
	"604":{41.5732, -87.6899}, // Hazel Crest, IL
	"604":{41.7335, -87.8859}, // Willow Springs, IL
	"604":{41.5758, -87.6503}, // East Hazel Crest, IL
	"604":{41.7495, -87.8345}, // Justice, IL
	"604":{41.4794, -88.4597}, // Lisbon, IL
	"604":{41.3255, -88.0532}, // Symerton, IL
	"604":{41.6044, -87.9497}, // Homer Glen, IL
	"604":{41.5324, -87.8779}, // Mokena, IL
	"604":{41.5175, -88.2149}, // Shorewood, IL
	"604":{41.1727, -88.2660}, // East Brooklyn, IL
	"604":{41.5101, -87.6347}, // Chicago Heights, IL
	"604":{41.7034, -87.7795}, // Chicago Ridge, IL
	"604":{41.6284, -87.5979}, // Dolton, IL
	"604":{41.4943, -88.0756}, // Preston Heights, IL
	"604":{41.0986, -88.4240}, // Dwight, IL
	"604":{41.7312, -87.7311}, // Hometown, IL
	"604":{41.5095, -87.7468}, // Matteson, IL
	"604":{41.2970, -88.2999}, // Carbon Hill, IL
	"604":{41.4274, -87.9802}, // Manhattan, IL
	"604":{41.4460, -87.7154}, // University Park, IL
	"604":{41.6075, -87.8619}, // Orland Park, IL
	"604":{41.2417, -88.4234}, // Mazon, IL
	"604":{41.6254, -87.7243}, // Midlothian, IL
	"604":{41.1930, -88.3147}, // Gardner, IL
	"604":{41.8433, -87.7909}, // Berwyn, IL
	"604":{41.3748, -88.4301}, // Morris, IL
	"604":{41.8183, -87.7730}, // Stickney, IL
	"604":{41.5685, -88.1638}, // Crystal Lawns, IL
	"604":{41.2696, -88.2234}, // Braidwood, IL
	"604":{41.3203, -88.1633}, // Wilmington, IL
	"604":{41.5203, -88.0346}, // Ingalls Park, IL
	"604":{41.4887, -87.8360}, // Frankfort, IL
	"604":{41.5591, -87.6610}, // Homewood, IL
	"604":{41.4817, -87.6868}, // Park Forest, IL
	"604":{41.5221, -87.8033}, // Frankfort Square, IL
	"604":{41.6319, -88.0996}, // Romeoville, IL
	"604":{41.1586, -88.6542}, // Ransom, IL
	"604":{41.6682, -87.8885}, // Palos Park, IL
	"604":{41.5189, -88.1499}, // Joliet, IL
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <605> with city FOX VALLEY and state IL
	"605":{41.7637, -88.2901}, // Aurora, IL
	"605":{41.7635, -87.9456}, // Willowbrook, IL
	"605":{41.7669, -88.5264}, // Big Rock, IL
	"605":{41.8083, -88.3413}, // North Aurora, IL
	"605":{41.8005, -87.9273}, // Hinsdale, IL
	"605":{41.8493, -89.0150}, // Steward, IL
	"605":{41.5880, -88.9229}, // Earlville, IL
	"605":{41.7683, -88.7659}, // Waterman, IL
	"605":{41.6561, -88.4508}, // Yorkville, IL
	"605":{41.7483, -88.1657}, // Naperville, IL
	"605":{41.8372, -87.9512}, // Oak Brook, IL
	"605":{41.8461, -87.8263}, // North Riverside, IL
	"605":{41.5285, -88.6790}, // Sheridan, IL
	"605":{41.6834, -88.3372}, // Oswego, IL
	"605":{41.7655, -88.8848}, // Shabbona, IL
	"605":{41.7980, -87.9569}, // Clarendon Hills, IL
	"605":{41.7959, -87.8410}, // McCook, IL
	"605":{41.8121, -87.8192}, // Lyons, IL
	"605":{41.7663, -87.8607}, // Hodgkins, IL
	"605":{41.7922, -88.0883}, // Lisle, IL
	"605":{41.6420, -88.6732}, // Somonauk, IL
	"605":{41.7447, -87.9823}, // Darien, IL
	"605":{41.7847, -88.4225}, // Prestbury, IL
	"605":{41.7112, -88.3353}, // Boulder Hill, IL
	"605":{41.6156, -88.6703}, // Lake Holiday, IL
	"605":{41.6496, -88.6177}, // Sandwich, IL
	"605":{41.8210, -88.1857}, // Warrenville, IL
	"605":{41.7741, -87.8752}, // Countryside, IL
	"605":{41.7949, -88.0170}, // Downers Grove, IL
	"605":{41.6757, -88.5294}, // Plano, IL
	"605":{41.7714, -88.6397}, // Hinckley, IL
	"605":{41.4794, -88.4597}, // Lisbon, IL
	"605":{41.6160, -88.7983}, // Leland, IL
	"605":{41.8308, -87.8723}, // La Grange Park, IL
	"605":{41.6205, -88.2260}, // Plainfield, IL
	"605":{41.5368, -88.5804}, // Newark, IL
	"605":{41.8245, -87.8470}, // Brookfield, IL
	"605":{41.8072, -87.8741}, // La Grange, IL
	"605":{41.7368, -88.0408}, // Woodridge, IL
	"605":{41.7948, -87.9742}, // Westmont, IL
	"605":{41.7877, -87.8146}, // Summit, IL
	"605":{41.7954, -88.9417}, // Lee, IL
	"605":{41.8022, -87.9006}, // Western Springs, IL
	"605":{41.7485, -87.9199}, // Burr Ridge, IL
	"605":{41.7237, -88.3631}, // Montgomery, IL
	"605":{41.8479, -88.3109}, // Batavia, IL
	"605":{41.7758, -88.4480}, // Sugar Grove, IL
	"605":{41.8310, -87.8169}, // Riverside, IL
	"605":{41.5189, -88.1499}, // Joliet, IL
	"605":{41.5343, -88.3835}, // Plattville, IL
	"605":{41.6071, -88.5454}, // Millbrook, IL
	"605":{41.5623, -88.6023}, // Millington, IL
	"605":{41.7670, -87.7915}, // Bedford Park, IL
	"605":{41.7690, -87.8977}, // Indian Head Park, IL
	// END OF OPTIONS TO PICK
	"606":{41.8373, -87.6862}, // Chicago, IL
	"607":{41.8373, -87.6862}, // Chicago, IL
	"608":{41.8373, -87.6862}, // Chicago, IL
	"609":{41.1020, -87.8643}, // Kankakee, IL
	"610":{42.2598, -89.0641}, // Rockford, IL
	"611":{42.2598, -89.0641}, // Rockford, IL
	TODO PLACEHOLDER TO FIX ZIP <612> with city QUAD CITIES and state IL
	"612":{41.3515, -90.3767}, // Orion, IL
	"612":{41.7173, -89.9238}, // Lyndon, IL
	"612":{41.4416, -90.4464}, // Coal Valley, IL
	"612":{41.6108, -90.1751}, // Hillsdale, IL
	"612":{41.6306, -89.7852}, // Tampico, IL
	"612":{41.4373, -90.7164}, // Andalusia, IL
	"612":{41.8647, -90.1576}, // Fulton, IL
	"612":{41.4002, -90.5632}, // Coyne Center, IL
	"612":{41.6098, -89.6874}, // Deer Grove, IL
	"612":{41.4678, -90.3446}, // Colona, IL
	"612":{41.2984, -90.1907}, // Cambridge, IL
	"612":{41.6166, -90.3273}, // Port Byron, IL
	"612":{41.8077, -89.9618}, // Morrison, IL
	"612":{41.1697, -91.0000}, // New Boston, IL
	"612":{41.4416, -90.5595}, // Milan, IL
	"612":{41.5172, -90.5399}, // Rock Island Arsenal, IL
	"612":{41.6703, -89.9349}, // Prophetstown, IL
	"612":{41.9726, -90.1117}, // Thomson, IL
	"612":{41.4120, -90.5742}, // Oak Grove, IL
	"612":{41.2596, -90.6055}, // Matherville, IL
	"612":{41.5788, -90.3394}, // Rapids City, IL
	"612":{41.3959, -89.8887}, // Annawan, IL
	"612":{41.7861, -90.2166}, // Albany, IL
	"612":{41.4821, -90.4919}, // Moline, IL
	"612":{41.5033, -90.3181}, // Cleveland, IL
	"612":{41.5555, -90.4046}, // Hampton, IL
	"612":{41.5199, -90.3880}, // East Moline, IL
	"612":{41.3108, -90.4990}, // Sherrard, IL
	"612":{41.6773, -90.3215}, // Cordova, IL
	"612":{41.3321, -90.6722}, // Reynolds, IL
	"612":{41.4699, -90.5827}, // Rock Island, IL
	"612":{41.2949, -90.2907}, // Andover, IL
	"612":{41.1969, -90.8788}, // Joy, IL
	"612":{41.4975, -90.4101}, // Silvis, IL
	"612":{41.4507, -90.1540}, // Geneseo, IL
	"612":{41.4983, -90.3920}, // Carbon Cliff, IL
	"612":{41.1984, -90.7459}, // Aledo, IL
	"612":{41.4167, -90.0088}, // Atkinson, IL
	"612":{41.5218, -89.9140}, // Hooppole, IL
	"612":{41.6587, -90.0815}, // Erie, IL
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <613> with city LA SALLE and state IL
	"613":{41.0486, -89.0522}, // Wenona, IL
	"613":{41.3573, -89.7372}, // Sheffield, IL
	"613":{41.1242, -88.8295}, // Streator, IL
	"613":{41.7290, -89.3679}, // Amboy, IL
	"613":{41.2637, -89.2300}, // Granville, IL
	"613":{41.0044, -89.1339}, // Toluca, IL
	"613":{41.2559, -89.1814}, // Standard, IL
	"613":{41.3809, -89.4648}, // Princeton, IL
	"613":{41.4340, -89.3959}, // Dover, IL
	"613":{41.6942, -89.0863}, // Compton, IL
	"613":{40.9920, -88.7297}, // Cornell, IL
	"613":{41.4283, -89.2137}, // Cherry, IL
	"613":{41.3099, -88.6850}, // Marseilles, IL
	"613":{41.6876, -88.9810}, // Paw Paw, IL
	"613":{41.4657, -89.0815}, // Troy Grove, IL
	"613":{41.3210, -88.5880}, // Seneca, IL
	"613":{41.3532, -88.8307}, // Ottawa, IL
	"613":{41.2960, -89.0693}, // Oglesby, IL
	"613":{41.3490, -89.1368}, // Peru, IL
	"613":{41.1252, -89.0668}, // Lostant, IL
	"613":{41.3359, -89.2034}, // Spring Valley, IL
	"613":{41.5553, -89.1042}, // Mendota, IL
	"613":{41.2651, -89.1260}, // Cedar Point, IL
	"613":{41.0353, -89.2244}, // Varna, IL
	"613":{41.2589, -89.3216}, // Hennepin, IL
	"613":{41.5112, -89.7184}, // New Bedford, IL
	"613":{41.2624, -89.2550}, // Mark, IL
	"613":{41.3816, -89.2145}, // Ladd, IL
	"613":{41.3291, -89.6791}, // Buda, IL
	"613":{41.2927, -89.5085}, // Tiskilwa, IL
	"613":{41.3648, -89.2948}, // Hollowayville, IL
	"613":{41.1772, -89.2096}, // McNabb, IL
	"613":{41.1141, -89.1958}, // Magnolia, IL
	"613":{41.3647, -89.2740}, // Seatonville, IL
	"613":{41.3553, -89.1763}, // Dalzell, IL
	"613":{41.3606, -89.5831}, // Wyanet, IL
	"613":{41.5301, -89.2833}, // La Moille, IL
	"613":{40.9838, -89.0413}, // Rutland, IL
	"613":{41.4717, -89.2482}, // Arlington, IL
	"613":{41.1480, -88.8722}, // Kangley, IL
	"613":{40.9565, -88.9500}, // Dana, IL
	"613":{41.3856, -88.7959}, // Dayton, IL
	"613":{41.3959, -89.8887}, // Annawan, IL
	"613":{41.4249, -89.3698}, // Malden, IL
	"613":{41.3843, -89.8388}, // Mineral, IL
	"613":{41.3275, -89.2963}, // De Pue, IL
	"613":{41.4553, -89.6685}, // Manlius, IL
	"613":{41.6437, -89.2311}, // Sublette, IL
	"613":{41.2159, -89.0702}, // Tonica, IL
	"613":{41.3322, -88.8781}, // Naplate, IL
	"613":{41.5570, -89.5908}, // Walnut, IL
	"613":{41.6928, -89.1473}, // West Brooklyn, IL
	"613":{41.1889, -88.9823}, // Leonore, IL
	"613":{41.2968, -89.7894}, // Neponset, IL
	"613":{41.0055, -88.8931}, // Long Point, IL
	"613":{41.5560, -89.4606}, // Ohio, IL
	"613":{41.3575, -89.0718}, // LaSalle, IL
	"613":{41.2362, -88.8310}, // Grand Ridge, IL
	"613":{41.3536, -89.0108}, // North Utica, IL
	// END OF OPTIONS TO PICK
	"614":{40.9506, -90.3763}, // Galesburg, IL
	"615":{40.7521, -89.6155}, // Peoria, IL
	"616":{40.7521, -89.6155}, // Peoria, IL
	"617":{40.4757, -88.9703}, // Bloomington, IL
	"618":{40.1144, -88.2735}, // Champaign, IL
	"619":{40.1144, -88.2735}, // Champaign, IL
	TODO PLACEHOLDER TO FIX ZIP <620> with city ST LOUIS and state MO
	"620":{39.1258, -89.8173}, // Gillespie, IL
	"620":{39.0290, -89.5247}, // Panama, IL
	"620":{38.7631, -90.0816}, // Mitchell, IL
	"620":{39.1670, -89.4735}, // Hillsboro, IL
	"620":{39.0689, -89.6189}, // Walshville, IL
	"620":{39.3001, -89.2852}, // Nokomis, IL
	"620":{39.1546, -90.1642}, // Fidelity, IL
	"620":{38.8840, -90.1073}, // East Alton, IL
	"620":{39.1122, -89.2137}, // Bingham, IL
	"620":{39.0306, -89.4745}, // Donnellson, IL
	"620":{38.7295, -90.1263}, // Granite City, IL
	"620":{39.0416, -89.9512}, // Bunker Hill, IL
	"620":{39.1154, -89.2785}, // Fillmore, IL
	"620":{39.2693, -90.2066}, // Rockbridge, IL
	"620":{39.2054, -89.4056}, // Irving, IL
	"620":{38.7261, -89.9645}, // Maryville, IL
	"620":{39.0930, -89.8023}, // Benld, IL
	"620":{39.1765, -90.1414}, // Medora, IL
	"620":{39.1002, -89.8254}, // Mount Clare, IL
	"620":{39.2971, -90.6127}, // Kampsville, IL
	"620":{38.8632, -90.0774}, // Wood River, IL
	"620":{38.8120, -90.0611}, // South Roxana, IL
	"620":{39.4496, -90.5387}, // Hillview, IL
	"620":{38.6719, -90.1689}, // Venice, IL
	"620":{39.1121, -89.7832}, // Eagarville, IL
	"620":{39.5593, -90.4329}, // Alsey, IL
	"620":{38.7922, -89.9881}, // Edwardsville, IL
	"620":{39.0008, -89.5748}, // Sorento, IL
	"620":{39.1378, -89.8127}, // East Gillespie, IL
	"620":{39.1905, -90.3512}, // Kane, IL
	"620":{38.7868, -89.7789}, // Marine, IL
	"620":{38.8891, -89.8429}, // Hamel, IL
	"620":{38.9682, -89.7642}, // Livingston, IL
	"620":{38.8212, -90.0909}, // Hartford, IL
	"620":{38.7003, -90.1425}, // Madison, IL
	"620":{39.3200, -89.2879}, // Wenonah, IL
	"620":{39.3434, -89.2192}, // Ohlman, IL
	"620":{39.2948, -90.4062}, // Carrollton, IL
	"620":{38.9699, -89.6664}, // New Douglas, IL
	"620":{39.0406, -90.1404}, // Brighton, IL
	"620":{39.2545, -89.3492}, // Witt, IL
	"620":{39.1180, -90.3276}, // Jerseyville, IL
	"620":{39.2863, -90.5538}, // Eldred, IL
	"620":{39.0998, -89.7464}, // Lake Ka-Ho, IL
	"620":{38.9315, -89.8404}, // Worden, IL
	"620":{39.1088, -90.4996}, // Fieldon, IL
	"620":{39.0331, -90.6534}, // Batchtown, IL
	"620":{38.9540, -90.3554}, // Elsah, IL
	"620":{39.0716, -89.7637}, // White City, IL
	"620":{38.9014, -90.0467}, // Bethalto, IL
	"620":{39.2848, -89.3030}, // Coalton, IL
	"620":{38.6544, -90.1679}, // Brooklyn, IL
	"620":{39.1591, -90.6248}, // Hardin, IL
	"620":{38.8318, -90.0465}, // Roxana, IL
	"620":{39.0859, -89.8878}, // Dorchester, IL
	"620":{38.9226, -89.9361}, // Holiday Shores, IL
	"620":{38.9485, -90.5890}, // Brussels, IL
	"620":{39.3443, -90.2080}, // Greenfield, IL
	"620":{38.6503, -90.1026}, // Fairmont City, IL
	"620":{39.1632, -89.4609}, // Schram City, IL
	"620":{39.0117, -89.7905}, // Staunton, IL
	"620":{39.0508, -90.3985}, // Otterville, IL
	"620":{39.4840, -90.3741}, // Roodhouse, IL
	"620":{38.7208, -90.0601}, // Pontoon Beach, IL
	"620":{38.8878, -89.7360}, // Alhambra, IL
	"620":{39.1443, -89.1105}, // Ramsey, IL
	"620":{39.1963, -89.6288}, // Litchfield, IL
	"620":{39.1976, -89.5318}, // Butler, IL
	"620":{39.0694, -89.8553}, // Wilsonville, IL
	"620":{38.9581, -90.2156}, // Godfrey, IL
	"620":{39.1310, -89.4950}, // Taylor Springs, IL
	"620":{38.9033, -90.1523}, // Alton, IL
	"620":{39.0726, -89.7277}, // Mount Olive, IL
	"620":{39.5424, -90.3303}, // Manchester, IL
	"620":{39.0785, -89.8020}, // Sawyerville, IL
	"620":{38.7579, -89.9832}, // Glen Carbon, IL
	"620":{39.2328, -90.7152}, // Hamburg, IL
	"620":{38.9765, -90.4259}, // Grafton, IL
	"620":{38.8885, -90.0722}, // Rosewood Heights, IL
	"620":{39.4388, -90.4019}, // White Hall, IL
	"620":{38.9872, -89.7644}, // Williamson, IL
	"620":{39.0883, -89.3898}, // Coffeen, IL
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <622> with city ST LOUIS and state MO
	"622":{38.6070, -89.6844}, // Trenton, IL
	"622":{38.3188, -89.8743}, // New Athens, IL
	"622":{38.5955, -89.7497}, // Summerfield, IL
	"622":{38.3402, -90.1538}, // Waterloo, IL
	"622":{38.5850, -90.1139}, // Alorton, IL
	"622":{38.3064, -90.2979}, // Valmeyer, IL
	"622":{38.6218, -89.3733}, // Carlyle, IL
	"622":{38.5452, -89.6143}, // Albers, IL
	"622":{38.5144, -90.2167}, // Dupo, IL
	"622":{38.2096, -89.9998}, // Red Bud, IL
	"622":{38.6163, -89.6088}, // Aviston, IL
	"622":{38.6826, -89.5536}, // St. Rose, IL
	"622":{38.7772, -89.5645}, // Pierron, IL
	"622":{38.5649, -90.1792}, // Cahokia, IL
	"622":{38.2649, -89.5035}, // Oakdale, IL
	"622":{38.3913, -89.4865}, // Addieville, IL
	"622":{38.4918, -89.8693}, // Rentchler, IL
	"622":{38.7198, -89.7677}, // St. Jacob, IL
	"622":{38.8945, -89.3399}, // Smithboro, IL
	"622":{38.2501, -89.7746}, // Marissa, IL
	"622":{38.1644, -90.2129}, // Fults, IL
	"622":{38.5191, -89.8045}, // Mascoutah, IL
	"622":{38.6138, -89.5231}, // Breese, IL
	"622":{38.9245, -89.2662}, // Mulberry Grove, IL
	"622":{38.0085, -89.6625}, // Steeleville, IL
	"622":{38.3764, -90.0554}, // Floraville, IL
	"622":{38.7003, -90.1425}, // Madison, IL
	"622":{38.0153, -89.6188}, // Percy, IL
	"622":{38.7268, -89.8974}, // Troy, IL
	"622":{38.0093, -89.9101}, // Ellis Grove, IL
	"622":{38.3595, -90.0465}, // Paderborn, IL
	"622":{38.5412, -89.8528}, // Scott AFB, IL
	"622":{38.1838, -89.8451}, // Baldwin, IL
	"622":{38.3054, -89.9935}, // Hecker, IL
	"622":{38.6156, -90.1304}, // East St. Louis, IL
	"622":{38.6013, -89.8120}, // Lebanon, IL
	"622":{38.4380, -89.3705}, // New Minden, IL
	"622":{38.8288, -89.6664}, // Grantfork, IL
	"622":{38.6769, -90.0059}, // Collinsville, IL
	"622":{38.0851, -89.3718}, // Pinckneyville, IL
	"622":{38.8866, -89.3893}, // Greenville, IL
	"622":{38.5976, -89.9156}, // O'Fallon, IL
	"622":{38.5404, -89.2644}, // Hoffman, IL
	"622":{38.5165, -89.9900}, // Belleville, IL
	"622":{38.3190, -89.7288}, // Darmstadt, IL
	"622":{38.5845, -90.1611}, // Sauget, IL
	"622":{38.5974, -90.0052}, // Fairview Heights, IL
	"622":{38.3966, -89.6456}, // Venedy, IL
	"622":{38.6544, -90.1679}, // Brooklyn, IL
	"622":{37.9200, -89.8260}, // Chester, IL
	"622":{38.6285, -90.0928}, // Washington Park, IL
	"622":{38.5507, -89.9859}, // Swansea, IL
	"622":{38.3514, -89.3770}, // Nashville, IL
	"622":{38.7460, -89.2737}, // Keyesport, IL
	"622":{38.1514, -89.7185}, // Sparta, IL
	"622":{38.6054, -89.2898}, // Huey, IL
	"622":{38.1850, -89.6046}, // Coulterville, IL
	"622":{38.4160, -89.9902}, // Smithton, IL
	"622":{38.6503, -90.1026}, // Fairmont City, IL
	"622":{38.4397, -89.9169}, // Freeburg, IL
	"622":{38.3775, -89.7967}, // Fayetteville, IL
	"622":{38.4579, -90.0834}, // Millstadt, IL
	"622":{38.5386, -90.2400}, // East Carondelet, IL
	"622":{38.7208, -90.0601}, // Pontoon Beach, IL
	"622":{38.4344, -89.5480}, // Okawville, IL
	"622":{38.8234, -89.5394}, // Pocahontas, IL
	"622":{38.8927, -89.5721}, // Old Ripley, IL
	"622":{38.3633, -89.7135}, // St. Libory, IL
	"622":{38.2124, -89.6840}, // Tilden, IL
	"622":{38.0889, -89.9322}, // Evansville, IL
	"622":{38.5371, -89.4671}, // Bartelso, IL
	"622":{38.5552, -89.5412}, // Germantown, IL
	"622":{38.5367, -89.7072}, // New Baden, IL
	"622":{38.6059, -89.4322}, // Beckemeyer, IL
	"622":{38.7602, -89.6810}, // Highland, IL
	"622":{38.1306, -90.0003}, // Ruma, IL
	"622":{38.5087, -89.6149}, // Damiansville, IL
	"622":{38.6300, -90.0342}, // Caseyville, IL
	"622":{37.8405, -89.6968}, // Rockwood, IL
	"622":{38.4581, -90.2156}, // Columbia, IL
	"622":{37.9828, -89.5898}, // Willisville, IL
	"622":{38.5534, -89.9160}, // Shiloh, IL
	"622":{38.5798, -90.1039}, // Centreville, IL
	"622":{38.0325, -89.5676}, // Cutler, IL
	"622":{38.0817, -90.0976}, // Prairie du Rocher, IL
	"622":{38.2859, -89.8179}, // Lenzburg, IL
	"622":{38.2267, -90.2330}, // Maeystown, IL
	// END OF OPTIONS TO PICK
	"623":{39.9335, -91.3798}, // Quincy, IL
	"624":{39.1207, -88.5509}, // Effingham, IL
	"625":{39.7710, -89.6537}, // Springfield, IL
	"626":{39.7710, -89.6537}, // Springfield, IL
	"627":{39.7710, -89.6537}, // Springfield, IL
	"628":{38.5224, -89.1233}, // Centralia, IL
	"629":{37.7220, -89.2238}, // Carbondale, IL
	TODO PLACEHOLDER TO FIX ZIP <630> with city ST LOUIS and state MO
	"630":{38.3655, -90.3645}, // Kimmswick, MO
	"630":{38.3360, -90.4046}, // Barnhart, MO
	"630":{38.3361, -90.9711}, // Parkway, MO
	"630":{38.2127, -91.1637}, // Sullivan, MO
	"630":{38.4401, -90.9928}, // Union, MO
	"630":{38.1346, -90.4583}, // Olympian Village, MO
	"630":{38.5895, -90.5884}, // Ellisville, MO
	"630":{38.4795, -90.5281}, // Parkdale, MO
	"630":{38.6048, -91.2178}, // New Haven, MO
	"630":{38.6256, -90.5945}, // Clarkson Valley, MO
	"630":{38.2871, -90.3995}, // Pevely, MO
	"630":{38.6588, -90.5803}, // Chesterfield, MO
	"630":{38.2330, -90.5670}, // Hillsboro, MO
	"630":{38.7188, -90.4750}, // Maryland Heights, MO
	"630":{38.2221, -90.3808}, // Crystal City, MO
	"630":{38.3913, -90.5908}, // Scotsdale, MO
	"630":{38.4163, -90.6777}, // LaBarque Creek, MO
	"630":{38.5950, -90.5501}, // Ballwin, MO
	"630":{38.2573, -90.3934}, // Herculaneum, MO
	"630":{38.7996, -90.3269}, // Florissant, MO
	"630":{38.7157, -90.3684}, // Breckenridge Hills, MO
	"630":{38.7266, -90.3872}, // St. Ann, MO
	"630":{38.4676, -90.5426}, // Peaceful Village, MO
	"630":{38.4294, -90.3725}, // Arnold, MO
	"630":{38.7931, -90.3900}, // Hazelwood, MO
	"630":{38.3575, -90.6409}, // Cedar Hill, MO
	"630":{38.4805, -90.7541}, // Pacific, MO
	"630":{38.6718, -91.3392}, // Berger, MO
	"630":{38.1922, -91.1918}, // West Sullivan, MO
	"630":{38.5657, -90.5010}, // Twin Oaks, MO
	"630":{38.7993, -90.2640}, // Black Jack, MO
	"630":{38.8394, -90.2847}, // Old Jamestown, MO
	"630":{38.4922, -90.4856}, // Murphy, MO
	"630":{38.4396, -90.5742}, // Byrnes Mill, MO
	"630":{38.3850, -91.4024}, // Rosebud, MO
	"630":{38.5799, -90.6698}, // Wildwood, MO
	"630":{38.3993, -91.3305}, // Gerald, MO
	"630":{38.5014, -90.6491}, // Eureka, MO
	"630":{38.4179, -91.2312}, // Leslie, MO
	"630":{38.4418, -90.7147}, // Lake Tekakwitha, MO
	"630":{38.7672, -90.4277}, // Bridgeton, MO
	"630":{38.4609, -90.5340}, // High Ridge, MO
	"630":{38.2279, -91.1493}, // Oak Grove Village, MO
	"630":{38.2198, -90.4095}, // Festus, MO
	"630":{38.2811, -91.0967}, // Charmwood, MO
	"630":{38.4680, -90.8850}, // Villa Ridge, MO
	"630":{38.2389, -91.0688}, // Miramiguoa Park, MO
	"630":{38.4951, -90.8173}, // Gray Summit, MO
	"630":{38.5513, -90.4924}, // Valley Park, MO
	"630":{38.1410, -90.5609}, // De Soto, MO
	"630":{38.7448, -90.4537}, // Champ, MO
	"630":{38.1732, -91.2127}, // St. Cloud, MO
	"630":{38.2673, -90.4320}, // Horine, MO
	"630":{38.3479, -90.9934}, // St. Clair, MO
	"630":{38.3299, -90.6582}, // Cedar Hill Lakes, MO
	"630":{38.5897, -90.5261}, // Winchester, MO
	"630":{38.5514, -91.0151}, // Washington, MO
	"630":{38.3672, -90.3706}, // Imperial, MO
	"630":{38.5279, -90.4489}, // Fenton, MO
	"630":{38.6317, -90.4789}, // Town and Country, MO
	"630":{38.5830, -90.5065}, // Manchester, MO
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <631> with city ST LOUIS and state MO
	"631":{38.5242, -90.3365}, // Green Park, MO
	"631":{38.7148, -90.3461}, // St. John, MO
	"631":{38.5117, -90.3574}, // Concord, MO
	"631":{38.7208, -90.2749}, // Country Club Hills, MO
	"631":{38.7010, -90.2965}, // Glen Echo Park, MO
	"631":{38.5866, -90.3282}, // Shrewsbury, MO
	"631":{38.5436, -90.2854}, // Bella Villa, MO
	"631":{38.6857, -90.2868}, // Hillsdale, MO
	"631":{38.7069, -90.3010}, // Normandy, MO
	"631":{38.7438, -90.3361}, // Berkeley, MO
	"631":{38.5325, -90.2845}, // Lemay, MO
	"631":{38.6309, -90.3332}, // Richmond Heights, MO
	"631":{38.6621, -90.4430}, // Creve Coeur, MO
	"631":{38.7563, -90.2766}, // Dellwood, MO
	"631":{38.6888, -90.3392}, // Vinita Park, MO
	"631":{38.7017, -90.3180}, // Bel-Nor, MO
	"631":{38.7579, -90.1984}, // Glasgow Village, MO
	"631":{38.5569, -90.3782}, // Crestwood, MO
	"631":{38.6069, -90.3912}, // Warson Woods, MO
	"631":{38.6121, -90.3240}, // Maplewood, MO
	"631":{38.7443, -90.2109}, // Riverview, MO
	"631":{38.6751, -90.2941}, // Wellston, MO
	"631":{38.7035, -90.2824}, // Northwoods, MO
	"631":{38.7130, -90.3285}, // Bel-Ridge, MO
	"631":{38.6927, -90.2829}, // Uplands Park, MO
	"631":{38.6856, -90.3250}, // Hanley Hills, MO
	"631":{38.7188, -90.4750}, // Maryland Heights, MO
	"631":{38.5866, -90.3544}, // Webster Groves, MO
	"631":{38.5682, -90.3395}, // Marlborough, MO
	"631":{38.5532, -90.3086}, // Wilbur Park, MO
	"631":{38.7529, -90.2279}, // Bellefontaine Neighbors, MO
	"631":{38.7157, -90.3684}, // Breckenridge Hills, MO
	"631":{38.6444, -90.3303}, // Clayton, MO
	"631":{38.5509, -90.3532}, // Grantwood Village, MO
	"631":{38.6378, -90.3815}, // Ladue, MO
	"631":{38.7230, -90.2643}, // Jennings, MO
	"631":{38.7174, -90.2650}, // Flordell Hills, MO
	"631":{38.7287, -90.3600}, // Woodson Terrace, MO
	"631":{38.5310, -90.4088}, // Sunset Hills, MO
	"631":{38.6979, -90.2900}, // Beverly Hills, MO
	"631":{38.5400, -90.3384}, // Lakeshire, MO
	"631":{38.7885, -90.2078}, // Spanish Lake, MO
	"631":{38.4472, -90.3199}, // Oakville, MO
	"631":{38.5973, -90.4482}, // Des Peres, MO
	"631":{38.7143, -90.2896}, // Norwood Court, MO
	"631":{38.5809, -90.3164}, // Mackenzie, MO
	"631":{38.7351, -90.3658}, // Edmundson, MO
	"631":{38.6953, -90.2756}, // Pine Lawn, MO
	"631":{38.5018, -90.3149}, // Mehlville, MO
	"631":{38.6657, -90.3315}, // University City, MO
	"631":{38.6132, -90.4090}, // Huntleigh, MO
	"631":{38.7384, -90.3249}, // Kinloch, MO
	"631":{38.7013, -90.3489}, // Sycamore Hills, MO
	"631":{38.6936, -90.3127}, // Greendale, MO
	"631":{38.7655, -90.3105}, // Calverton Park, MO
	"631":{38.7124, -90.3121}, // Bellerive Acres, MO
	"631":{38.6358, -90.2451}, // St. Louis, MO
	"631":{38.7110, -90.2974}, // Pasadena Park, MO
	"631":{38.6090, -90.3673}, // Rock Hill, MO
	"631":{38.6212, -90.4319}, // Crystal Lake Park, MO
	"631":{38.6967, -90.3689}, // Overland, MO
	"631":{38.6225, -90.4551}, // Country Life Acres, MO
	"631":{38.6434, -90.4333}, // Westwood, MO
	"631":{38.6856, -90.3293}, // Vinita Terrace, MO
	"631":{38.6725, -90.3784}, // Olivette, MO
	"631":{38.7251, -90.3058}, // Cool Valley, MO
	"631":{38.6801, -90.3082}, // Pagedale, MO
	"631":{38.6922, -90.2877}, // Velda Village Hills, MO
	"631":{38.7568, -90.2485}, // Castle Point, MO
	"631":{38.5935, -90.3826}, // Glendale, MO
	"631":{38.7456, -90.2430}, // Moline Acres, MO
	"631":{38.6194, -90.3476}, // Brentwood, MO
	"631":{38.6941, -90.2934}, // Velda City, MO
	"631":{38.5769, -90.3849}, // Oakland, MO
	"631":{38.6300, -90.4189}, // Frontenac, MO
	"631":{38.7029, -90.3426}, // Charlack, MO
	"631":{38.5499, -90.3263}, // Affton, MO
	"631":{38.7490, -90.2949}, // Ferguson, MO
	"631":{38.5260, -90.3730}, // Sappington, MO
	"631":{38.7083, -90.2922}, // Pasadena Hills, MO
	"631":{38.5789, -90.4203}, // Kirkwood, MO
	"631":{38.6317, -90.4789}, // Town and Country, MO
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <633> with city ST LOUIS and state MO
	"633":{38.8664, -90.2084}, // West Alton, MO
	"633":{39.2747, -91.5763}, // Farber, MO
	"633":{39.2622, -90.8287}, // Annada, MO
	"633":{38.7848, -90.7873}, // Lake St. Louis, MO
	"633":{38.9125, -91.5323}, // Danville, MO
	"633":{38.8546, -91.3027}, // Jonesburg, MO
	"633":{39.1682, -90.7879}, // Elsberry, MO
	"633":{39.0036, -91.3521}, // Bellflower, MO
	"633":{38.9944, -90.7437}, // Winfield, MO
	"633":{38.9689, -90.8493}, // Fountain N' Lakes, MO
	"633":{38.7558, -90.7313}, // Dardenne Prairie, MO
	"633":{39.1286, -91.4138}, // Middletown, MO
	"633":{39.3693, -90.9049}, // Clarksville, MO
	"633":{38.9092, -91.4529}, // New Florence, MO
	"633":{38.8165, -90.8671}, // Wentzville, MO
	"633":{39.3458, -91.3414}, // Curryville, MO
	"633":{38.9155, -90.8020}, // Chain of Rocks, MO
	"633":{38.8484, -90.7406}, // St. Paul, MO
	"633":{38.7117, -90.6517}, // Weldon Spring, MO
	"633":{38.8262, -91.2322}, // Pendleton, MO
	"633":{38.8019, -91.4838}, // Big Spring, MO
	"633":{38.9737, -91.5024}, // Montgomery City, MO
	"633":{38.9273, -90.3434}, // Portage Des Sioux, MO
	"633":{39.0734, -91.5692}, // Wellsville, MO
	"633":{38.9715, -91.1327}, // Hawk Point, MO
	"633":{38.9708, -90.9714}, // Troy, MO
	"633":{38.7511, -90.6587}, // Cottleville, MO
	"633":{38.7048, -90.6853}, // Weldon Spring Heights, MO
	"633":{38.7182, -90.8834}, // New Melle, MO
	"633":{39.1839, -91.0165}, // Whiteside, MO
	"633":{39.0035, -91.2403}, // Truxton, MO
	"633":{38.9400, -90.9255}, // Moscow Mills, MO
	"633":{38.8188, -91.1367}, // Warrenton, MO
	"633":{38.5684, -90.8806}, // Augusta, MO
	"633":{38.7851, -90.7177}, // O'Fallon, MO
	"633":{38.7956, -90.5156}, // St. Charles, MO
	"633":{38.8291, -90.7873}, // Josephville, MO
	"633":{38.8637, -90.8685}, // Flint Hill, MO
	"633":{39.1252, -91.0577}, // Silex, MO
	"633":{39.0214, -91.0483}, // Cave, MO
	"633":{39.2829, -91.2098}, // St. Clement, MO
	"633":{38.8754, -91.3776}, // High Hill, MO
	"633":{39.3577, -91.1837}, // Tarrants, MO
	"633":{38.8124, -91.1220}, // Truesdale, MO
	"633":{39.0458, -90.7414}, // Foley, MO
	"633":{38.6309, -91.0572}, // Marthasville, MO
	"633":{38.9382, -90.7469}, // Old Monroe, MO
	"633":{38.6327, -90.7846}, // Defiance, MO
	"633":{39.2434, -91.6424}, // Laddonia, MO
	"633":{38.8341, -91.0399}, // Wright City, MO
	"633":{39.3080, -91.4892}, // Vandalia, MO
	"633":{38.5758, -90.9971}, // Three Creeks, MO
	"633":{38.7825, -90.6061}, // St. Peters, MO
	"633":{38.7631, -91.0552}, // Innsbrook, MO
	"633":{39.2389, -91.0121}, // Eolia, MO
	"633":{39.4413, -91.0626}, // Louisiana, MO
	"633":{39.3447, -91.2031}, // Bowling Green, MO
	"633":{38.8158, -90.9622}, // Foristell, MO
	"633":{39.2537, -91.2239}, // Ashley, MO
	"633":{39.2626, -90.9004}, // Paynesville, MO
	// END OF OPTIONS TO PICK
	"634":{39.9335, -91.3798}, // Quincy, IL
	"635":{39.9335, -91.3798}, // Quincy, IL
	"636":{37.3108, -89.5596}, // Cape Girardeau, MO
	"637":{37.3108, -89.5596}, // Cape Girardeau, MO
	"638":{37.3108, -89.5596}, // Cape Girardeau, MO
	"639":{37.3108, -89.5596}, // Cape Girardeau, MO
	"640":{39.1239, -94.5541}, // Kansas City, MO
	"641":{39.1239, -94.5541}, // Kansas City, MO
	TODO PLACEHOLDER TO FIX ZIP <644> with city ST JOSEPH and state MO
	"644":{40.4045, -94.4465}, // Worth, MO
	"644":{40.3392, -95.3920}, // Fairfax, MO
	"644":{40.4006, -95.5950}, // Phelps City, MO
	"644":{40.2880, -95.0795}, // Skidmore, MO
	"644":{40.3989, -94.3233}, // Denver, MO
	"644":{40.5187, -95.1170}, // Elmo, MO
	"644":{40.3735, -95.0771}, // Quitman, MO
	"644":{40.2684, -94.0281}, // Bethany, MO
	"644":{40.5350, -95.3212}, // Westboro, MO
	"644":{40.4752, -93.9280}, // Blythedale, MO
	"644":{40.4799, -95.6233}, // Watson, MO
	"644":{40.0253, -94.9731}, // Fillmore, MO
	"644":{39.5644, -94.4615}, // Plattsburg, MO
	"644":{40.5171, -94.6143}, // Sheridan, MO
	"644":{40.3428, -94.8708}, // Maryville, MO
	"644":{39.6377, -94.3210}, // Turney, MO
	"644":{40.2166, -94.5381}, // Stanberry, MO
	"644":{40.3292, -93.7971}, // Mount Moriah, MO
	"644":{39.7496, -94.3566}, // Osborn, MO
	"644":{40.4503, -94.8420}, // Pickering, MO
	"644":{39.5247, -94.7740}, // Dearborn, MO
	"644":{39.5022, -94.6290}, // Edgerton, MO
	"644":{40.1152, -94.8211}, // Bolckow, MO
	"644":{39.7226, -94.6404}, // Easton, MO
	"644":{39.9861, -95.1433}, // Oregon, MO
	"644":{39.4537, -94.6402}, // Ridgely, MO
	"644":{39.8139, -94.5507}, // Clarksdale, MO
	"644":{39.9829, -95.1883}, // Forest City, MO
	"644":{39.9794, -94.5987}, // Union Star, MO
	"644":{40.2485, -95.4544}, // Corning, MO
	"644":{40.1923, -95.3742}, // Craig, MO
	"644":{39.8888, -94.8926}, // Amazonia, MO
	"644":{40.1365, -95.2337}, // Mound City, MO
	"644":{40.5749, -95.2109}, // Blanchard, MO
	"644":{40.5403, -94.3900}, // Irena, MO
	"644":{40.1686, -94.7361}, // Guilford, MO
	"644":{40.5511, -94.8169}, // Hopkins, MO
	"644":{40.2021, -95.0780}, // Maitland, MO
	"644":{39.7469, -94.2364}, // Cameron, MO
	"644":{40.2684, -94.6914}, // Conception Junction, MO
	"644":{40.0505, -94.5251}, // King City, MO
	"644":{40.0408, -94.8230}, // Rosendale, MO
	"644":{40.4430, -95.3835}, // Tarkio, MO
	"644":{40.1755, -94.8232}, // Barnard, MO
	"644":{40.2652, -94.1952}, // New Hampton, MO
	"644":{40.1981, -94.3997}, // Darlington, MO
	"644":{40.2594, -94.8279}, // Arkoe, MO
	"644":{39.9093, -94.2419}, // Weatherby, MO
	"644":{40.4466, -95.0667}, // Burlington Junction, MO
	"644":{40.1098, -95.2894}, // Bigelow, MO
	"644":{40.4684, -93.9861}, // Eagleville, MO
	"644":{40.5080, -95.0326}, // Clearmont, MO
	"644":{40.2411, -94.6805}, // Conception, MO
	"644":{40.2479, -94.3335}, // Albany, MO
	"644":{39.5875, -95.0235}, // Rushville, MO
	"644":{39.5375, -95.0506}, // Lewis and Clark Village, MO
	"644":{39.5515, -94.3286}, // Lathrop, MO
	"644":{40.4858, -94.4135}, // Grant City, MO
	"644":{40.4853, -94.2887}, // Allendale, MO
	"644":{40.2013, -95.0404}, // Graham, MO
	"644":{39.6129, -94.5947}, // Gower, MO
	"644":{39.4773, -94.9839}, // Iatan, MO
	"644":{39.8859, -94.3634}, // Maysville, MO
	"644":{40.0519, -95.3176}, // Fortescue, MO
	"644":{40.3523, -94.6721}, // Ravenwood, MO
	"644":{39.7551, -94.4987}, // Stewartsville, MO
	"644":{39.4764, -94.5614}, // Trimble, MO
	"644":{40.0616, -94.7646}, // Rea, MO
	"644":{39.6698, -94.7587}, // Agency, MO
	"644":{39.8638, -94.6798}, // Cosby, MO
	"644":{40.4109, -95.5332}, // Rock Port, MO
	"644":{40.4384, -94.6223}, // Parnell, MO
	"644":{40.0735, -95.3518}, // Big Lake, MO
	"644":{40.3326, -94.4233}, // Gentry, MO
	"644":{40.2662, -94.6696}, // Clyde, MO
	"644":{40.3781, -93.9379}, // Ridgeway, MO
	"644":{39.9390, -94.8279}, // Savannah, MO
	"644":{39.5883, -94.9236}, // De Kalb, MO
	"644":{39.8683, -94.4354}, // Amity, MO
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <645> with city ST JOSEPH and state MO
	"645":{39.8388, -94.8205}, // Country Club, MO
	"645":{39.8888, -94.8926}, // Amazonia, MO
	"645":{39.6698, -94.7587}, // Agency, MO
	"645":{39.7598, -94.8210}, // St. Joseph, MO
	// END OF OPTIONS TO PICK
	"646":{39.7953, -93.5498}, // Chillicothe, MO
	"647":{38.6530, -94.3467}, // Harrisonville, MO
	"648":{37.1943, -93.2915}, // Springfield, MO
	"649":{39.1239, -94.5541}, // Kansas City, MO
	TODO PLACEHOLDER TO FIX ZIP <650> with city MID-MISSOURI and state MO
	"650":{38.6983, -91.4342}, // Hermann, MO
	"650":{38.3512, -92.5766}, // Eldon, MO
	"650":{38.2035, -92.6257}, // Lake Ozark, MO
	"650":{38.4414, -92.0005}, // Westphalia, MO
	"650":{38.6970, -92.3072}, // Hartsburg, MO
	"650":{38.3164, -91.9223}, // Freeburg, MO
	"650":{38.6611, -92.6661}, // Clarksburg, MO
	"650":{38.6452, -92.1151}, // Holts Summit, MO
	"650":{38.2848, -91.7217}, // Belle, MO
	"650":{38.6763, -92.1013}, // Lake Mykee Town, MO
	"650":{38.8136, -92.5906}, // Prairie Home, MO
	"650":{38.4103, -92.5300}, // Olean, MO
	"650":{38.7204, -91.5176}, // Rhineland, MO
	"650":{38.0880, -92.4845}, // Brumley, MO
	"650":{38.6180, -92.4090}, // Centertown, MO
	"650":{38.8019, -91.4838}, // Big Spring, MO
	"650":{38.6304, -92.5667}, // California, MO
	"650":{38.7356, -91.4449}, // McKittrick, MO
	"650":{38.2560, -92.2665}, // St. Elizabeth, MO
	"650":{38.6744, -91.8723}, // Mokane, MO
	"650":{38.6763, -91.7700}, // Chamois, MO
	"650":{38.5128, -92.4384}, // Russellville, MO
	"650":{38.1974, -92.7178}, // Village of Four Seasons, MO
	"650":{38.6716, -91.6334}, // Morrison, MO
	"650":{38.6549, -92.7804}, // Tipton, MO
	"650":{38.7164, -92.0916}, // New Bloomfield, MO
	"650":{38.0121, -92.7500}, // Camdenton, MO
	"650":{38.0457, -92.7049}, // Linn Creek, MO
	"650":{38.2954, -92.0256}, // Argyle, MO
	"650":{38.2294, -92.6053}, // Bagnell, MO
	"650":{38.3004, -91.6332}, // Bland, MO
	"650":{38.3670, -92.2160}, // St. Thomas, MO
	"650":{38.4418, -92.9901}, // Stover, MO
	"650":{38.3049, -92.8244}, // Gravois Mills, MO
	"650":{38.1353, -92.6479}, // Osage Beach, MO
	"650":{38.4788, -91.8450}, // Linn, MO
	"650":{38.3490, -91.4974}, // Owensville, MO
	"650":{38.2370, -92.4600}, // Tuscumbia, MO
	"650":{38.7929, -92.2478}, // Ashland, MO
	"650":{38.2080, -92.8253}, // Laurie, MO
	"650":{38.1674, -92.7794}, // Sunrise Beach, MO
	"650":{38.3121, -92.1660}, // Meta, MO
	"650":{38.7667, -92.4802}, // Jamestown, MO
	"650":{38.5430, -92.3646}, // Lohman, MO
	"650":{38.2041, -92.6223}, // Lakeside, MO
	"650":{38.8460, -92.4540}, // Lupus, MO
	"650":{38.4333, -92.8453}, // Versailles, MO
	"650":{38.3773, -92.6744}, // Barnett, MO
	"650":{38.6703, -91.5600}, // Gasconade, MO
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <651> with city MID-MISSOURI and state MO
	"651":{38.5157, -92.0647}, // Taos, MO
	"651":{38.5676, -92.1759}, // Jefferson City, MO
	"651":{38.5943, -92.3312}, // St. Martins, MO
	"651":{38.4900, -92.1791}, // Wardsville, MO
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <652> with city MID-MISSOURI and state MO
	"652":{39.4773, -92.0039}, // Paris, MO
	"652":{38.9341, -92.7025}, // Windsor Place, MO
	"652":{38.8736, -92.9126}, // Pilot Grove, MO
	"652":{39.2105, -92.1342}, // Centralia, MO
	"652":{39.5479, -91.8572}, // Stoutsville, MO
	"652":{38.8831, -92.4549}, // McBaine, MO
	"652":{39.1622, -91.8470}, // Vandiver, MO
	"652":{38.9588, -92.7471}, // Boonville, MO
	"652":{39.2345, -92.2823}, // Sturgeon, MO
	"652":{39.2279, -92.8394}, // Glasgow, MO
	"652":{39.1626, -91.8712}, // Mexico, MO
	"652":{39.4233, -92.8025}, // Salisbury, MO
	"652":{39.4983, -93.1939}, // Triplett, MO
	"652":{39.4257, -93.1267}, // Brunswick, MO
	"652":{39.3976, -92.9921}, // Dalton, MO
	"652":{39.4937, -92.1313}, // Holliday, MO
	"652":{38.7894, -92.7992}, // Bunceton, MO
	"652":{39.0112, -92.7550}, // Franklin, MO
	"652":{39.5114, -92.4411}, // Cairo, MO
	"652":{39.1009, -91.6476}, // Martinsburg, MO
	"652":{39.4735, -92.2117}, // Madison, MO
	"652":{39.0174, -91.8958}, // Auxvasse, MO
	"652":{39.0175, -92.7406}, // New Franklin, MO
	"652":{38.9120, -92.4786}, // Huntsdale, MO
	"652":{39.5147, -91.9449}, // Goss, MO
	"652":{38.9785, -92.5633}, // Rocheport, MO
	"652":{39.4313, -92.9371}, // Keytesville, MO
	"652":{38.9066, -92.5215}, // Wooldridge, MO
	"652":{38.8633, -92.3132}, // Pierpont, MO
	"652":{38.9473, -91.9389}, // Kingdom City, MO
	"652":{39.5877, -92.4729}, // Jacksonville, MO
	"652":{39.1346, -91.7649}, // Benton City, MO
	"652":{39.1196, -92.2268}, // Hallsville, MO
	"652":{39.2692, -92.7042}, // Armstrong, MO
	"652":{38.8551, -91.9510}, // Fulton, MO
	"652":{39.2763, -92.3475}, // Clark, MO
	"652":{39.1470, -92.6857}, // Fayette, MO
	"652":{39.3419, -92.4109}, // Renick, MO
	"652":{39.1396, -92.4580}, // Harrisburg, MO
	"652":{39.4387, -92.6643}, // Clifton Hill, MO
	"652":{39.4929, -91.7907}, // Florida, MO
	"652":{39.4190, -92.4365}, // Moberly, MO
	"652":{39.4364, -92.5440}, // Huntsville, MO
	"652":{39.6335, -92.4743}, // Excello, MO
	"652":{38.9477, -92.3255}, // Columbia, MO
	"652":{39.3059, -92.5128}, // Higbee, MO
	"652":{39.2104, -91.7246}, // Rush Hill, MO
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <653> with city MID-MISSOURI and state MO
	"653":{38.6211, -93.4101}, // Green Ridge, MO
	"653":{39.1940, -93.3636}, // Malta Bend, MO
	"653":{39.1145, -93.2010}, // Marshall, MO
	"653":{38.7302, -93.5590}, // Whiteman AFB, MO
	"653":{38.8994, -93.3592}, // Houstonia, MO
	"653":{38.9793, -92.9921}, // Blackwater, MO
	"653":{39.1048, -93.4855}, // Blackburn, MO
	"653":{38.5319, -93.5227}, // Windsor, MO
	"653":{38.2474, -93.3710}, // Warsaw, MO
	"653":{39.0699, -92.9473}, // Arrow Rock, MO
	"653":{38.8369, -93.2952}, // Hughesville, MO
	"653":{39.2050, -93.4427}, // Grand Pass, MO
	"653":{38.7722, -93.4237}, // La Monte, MO
	"653":{38.9946, -93.0308}, // Nelson, MO
	"653":{38.6696, -92.8763}, // Syracuse, MO
	"653":{38.7042, -93.2351}, // Sedalia, MO
	"653":{38.7029, -93.0028}, // Otterville, MO
	"653":{38.4595, -93.2025}, // Cole Camp, MO
	"653":{39.1254, -93.3945}, // Mount Leonard, MO
	"653":{39.2227, -93.0649}, // Slater, MO
	"653":{38.9751, -93.4950}, // Emma, MO
	"653":{38.6812, -93.0930}, // Smithton, MO
	"653":{38.1004, -93.0528}, // Climax Springs, MO
	"653":{38.5037, -93.3236}, // Ionia, MO
	"653":{38.9649, -93.4152}, // Sweet Springs, MO
	"653":{38.4683, -93.6251}, // Calhoun, MO
	"653":{38.3940, -93.3313}, // Lincoln, MO
	"653":{39.3224, -93.2258}, // Miami, MO
	"653":{39.2327, -93.0041}, // Gilliam, MO
	"653":{38.7675, -93.5616}, // Knob Noster, MO
	// END OF OPTIONS TO PICK
	"654":{37.1943, -93.2915}, // Springfield, MO
	"655":{37.1943, -93.2915}, // Springfield, MO
	"656":{37.1943, -93.2915}, // Springfield, MO
	"657":{37.1943, -93.2915}, // Springfield, MO
	"658":{37.1943, -93.2915}, // Springfield, MO
	"660":{39.1234, -94.7443}, // Kansas City, KS
	"661":{39.1234, -94.7443}, // Kansas City, KS
	"662":{39.1234, -94.7443}, // Kansas City, KS
	"664":{39.0346, -95.6955}, // Topeka, KS
	"665":{39.0346, -95.6955}, // Topeka, KS
	"666":{39.0346, -95.6955}, // Topeka, KS
	TODO PLACEHOLDER TO FIX ZIP <667> with city FT SCOTT and state KS
	"667":{37.9165, -95.1715}, // Moran, KS
	"667":{37.5430, -94.7024}, // Arma, KS
	"667":{37.5186, -95.1742}, // St. Paul, KS
	"667":{37.8365, -94.8820}, // Redfield, KS
	"667":{37.0752, -94.6353}, // Galena, KS
	"667":{37.8471, -94.9757}, // Uniontown, KS
	"667":{38.0248, -95.1740}, // Mildred, KS
	"667":{37.6695, -95.4621}, // Chanute, KS
	"667":{37.5250, -95.6614}, // Altoona, KS
	"667":{37.1714, -94.8442}, // Columbus, KS
	"667":{37.2780, -94.8230}, // Scammon, KS
	"667":{37.6409, -94.6239}, // Arcadia, KS
	"667":{37.5230, -94.6994}, // Franklin, KS
	"667":{37.8682, -95.7533}, // Yates Center, KS
	"667":{37.5678, -95.9363}, // New Albany, KS
	"667":{37.3539, -95.0193}, // McCune, KS
	"667":{38.0097, -94.7196}, // Fulton, KS
	"667":{37.7490, -95.1428}, // Savonburg, KS
	"667":{37.9274, -95.4006}, // Iola, KS
	"667":{37.8958, -95.0732}, // Bronson, KS
	"667":{37.4129, -94.6984}, // Pittsburg, KS
	"667":{37.7942, -95.1498}, // Elsmore, KS
	"667":{37.0196, -94.7351}, // Baxter Springs, KS
	"667":{37.6898, -95.1441}, // Stark, KS
	"667":{37.9220, -95.5375}, // Piqua, KS
	"667":{37.9062, -95.4077}, // Bassett, KS
	"667":{37.6637, -94.9697}, // Hepler, KS
	"667":{37.4242, -95.6849}, // Neodesha, KS
	"667":{37.4721, -95.3563}, // Galesburg, KS
	"667":{37.0509, -94.7032}, // Lowell, KS
	"667":{38.0155, -94.8834}, // Mapleton, KS
	"667":{37.6268, -95.7434}, // Benedict, KS
	"667":{37.8119, -95.4370}, // Humboldt, KS
	"667":{38.0634, -94.6924}, // Prescott, KS
	"667":{37.2805, -94.8438}, // Roseland, KS
	"667":{37.3451, -94.8215}, // Cherokee, KS
	"667":{37.7986, -95.9497}, // Toronto, KS
	"667":{37.9230, -95.3453}, // Gas, KS
	"667":{38.0058, -95.5556}, // Neosho Falls, KS
	"667":{37.5094, -94.8456}, // Girard, KS
	"667":{37.3083, -94.7747}, // Weir, KS
	"667":{37.5562, -94.6236}, // Mulberry, KS
	"667":{37.0733, -94.7060}, // Riverton, KS
	"667":{37.7092, -95.6968}, // Buffalo, KS
	"667":{37.3834, -94.7445}, // Chicopee, KS
	"667":{37.5328, -95.8223}, // Fredonia, KS
	"667":{37.9165, -95.3023}, // La Harpe, KS
	"667":{37.6018, -95.0742}, // Walnut, KS
	"667":{37.4583, -94.7018}, // Frontenac, KS
	"667":{37.4851, -95.4845}, // Thayer, KS
	"667":{37.5873, -95.4696}, // Earlton, KS
	"667":{37.2839, -94.9269}, // West Mineral, KS
	"667":{37.6872, -95.8959}, // Coyville, KS
	"667":{37.8283, -94.7038}, // Fort Scott, KS
	"667":{37.5718, -95.2418}, // Erie, KS
	// END OF OPTIONS TO PICK
	"668":{39.0346, -95.6955},  // Topeka, KS
	"669":{38.8137, -97.6143},  // Salina, KS
	"670":{37.6897, -97.3441},  // Wichita, KS
	"671":{37.6897, -97.3441},  // Wichita, KS
	"672":{37.6897, -97.3441},  // Wichita, KS
	"673":{37.2118, -95.7328},  // Independence, KS
	"674":{38.8137, -97.6143},  // Salina, KS
	"675":{38.0671, -97.9081},  // Hutchinson, KS
	"676":{38.8816, -99.3219},  // Hays, KS
	"677":{39.3843, -101.0459}, // Colby, KS
	"678":{37.7610, -100.0182}, // Dodge City, KS
	"679":{37.0466, -100.9295}, // Liberal, KS
	"680":{41.2628, -96.0498},  // Omaha, NE
	"681":{41.2628, -96.0498},  // Omaha, NE
	"683":{40.8088, -96.6796},  // Lincoln, NE
	"684":{40.8088, -96.6796},  // Lincoln, NE
	"685":{40.8088, -96.6796},  // Lincoln, NE
	"686":{42.0328, -97.4209},  // Norfolk, NE
	"687":{42.0328, -97.4209},  // Norfolk, NE
	"688":{40.9214, -98.3584},  // Grand Island, NE
	"689":{40.9214, -98.3584},  // Grand Island, NE
	TODO PLACEHOLDER TO FIX ZIP <690> with city MC COOK and state NE
	"690":{40.7065, -100.2155}, // Farnam, NE
	"690":{40.1143, -101.4045}, // Max, NE
	"690":{40.2046, -100.6213}, // McCook, NE
	"690":{40.0122, -101.9384}, // Haigler, NE
	"690":{40.4550, -101.5355}, // Enders, NE
	"690":{40.6335, -100.5112}, // Curtis, NE
	"690":{40.1119, -100.1067}, // Wilsonville, NE
	"690":{40.2511, -100.3104}, // Bartley, NE
	"690":{40.0437, -101.7251}, // Parks, NE
	"690":{40.4165, -101.3766}, // Wauneta, NE
	"690":{40.3845, -101.2346}, // Hamlet, NE
	"690":{40.2284, -100.8351}, // Culbertson, NE
	"690":{40.5725, -101.9793}, // Lamar, NE
	"690":{40.1503, -101.2280}, // Stratton, NE
	"690":{40.5335, -100.3845}, // Stockville, NE
	"690":{40.6899, -100.4004}, // Moorefield, NE
	"690":{40.0487, -100.2761}, // Lebanon, NE
	"690":{40.2346, -100.4198}, // Indianola, NE
	"690":{40.4706, -101.7469}, // Champion, NE
	"690":{40.1746, -101.0136}, // Trenton, NE
	"690":{40.3483, -101.1075}, // Palisade, NE
	"690":{40.2843, -100.1654}, // Cambridge, NE
	"690":{40.6644, -100.0296}, // Eustis, NE
	"690":{40.5145, -101.6374}, // Imperial, NE
	"690":{40.0376, -100.4051}, // Danbury, NE
	"690":{40.5110, -101.0203}, // Hayes Center, NE
	"690":{40.6587, -100.6222}, // Maywood, NE
	"690":{40.0517, -101.5354}, // Benkelman, NE
	// END OF OPTIONS TO PICK
	"691":{41.1266, -100.7640}, // North Platte, NE
	"692":{42.8739, -100.5498}, // Valentine, NE
	"693":{42.1025, -102.8766}, // Alliance, NE
	"700":{30.0687, -89.9288},  // New Orleans, LA
	"701":{30.0687, -89.9288},  // New Orleans, LA
	"703":{29.5799, -90.7058},  // Houma, LA
	"704":{30.3750, -90.0906},  // Mandeville, LA
	"705":{30.2084, -92.0323},  // Lafayette, LA
	"706":{30.2022, -93.2141},  // Lake Charles, LA
	"707":{30.4419, -91.1310},  // Baton Rouge, LA
	"708":{30.4419, -91.1310},  // Baton Rouge, LA
	"710":{32.4659, -93.7959},  // Shreveport, LA
	"711":{32.4659, -93.7959},  // Shreveport, LA
	"712":{32.5183, -92.0775},  // Monroe, LA
	"713":{31.2923, -92.4702},  // Alexandria, LA
	"714":{31.2923, -92.4702},  // Alexandria, LA
	"716":{34.2116, -92.0178},  // Pine Bluff, AR
	"717":{33.5672, -92.8467},  // Camden, AR
	"718":{33.4361, -93.9960},  // Texarkana, AR
	TODO PLACEHOLDER TO FIX ZIP <719> with city HOT SPRINGS NTL PK and state AR
	"719":{34.2602, -92.9610}, // Midway, AR
	"719":{34.5810, -94.2374}, // Mena, AR
	"719":{34.0005, -93.3376}, // Okolona, AR
	"719":{34.4859, -94.3793}, // Hatfield, AR
	"719":{34.0299, -93.5057}, // Delight, AR
	"719":{34.0639, -93.0958}, // Gum Springs, AR
	"719":{34.4269, -93.0887}, // Lake Hamilton, AR
	"719":{34.4902, -93.0498}, // Hot Springs, AR
	"719":{34.2988, -94.3331}, // Wickes, AR
	"719":{34.0343, -93.4226}, // Antoine, AR
	"719":{34.5693, -93.1726}, // Mountain Pine, AR
	"719":{34.4584, -93.6791}, // Norman, AR
	"719":{34.3800, -94.3642}, // Vandervoort, AR
	"719":{34.0649, -93.6902}, // Murfreesboro, AR
	"719":{34.6215, -93.7862}, // Oden, AR
	"719":{34.3278, -93.5311}, // Glenwood, AR
	"719":{34.1888, -93.0690}, // Caddo Valley, AR
	"719":{34.4366, -94.4160}, // Cove, AR
	"719":{34.2398, -94.3220}, // Grannis, AR
	"719":{34.2350, -92.9195}, // Donaldson, AR
	"719":{34.2243, -93.0037}, // Friendship, AR
	"719":{34.2349, -93.7406}, // Daisy, AR
	"719":{34.2556, -93.6522}, // Kirby, AR
	"719":{34.6568, -92.9644}, // Hot Springs Village, AR
	"719":{34.6136, -92.9200}, // Fountain Lake, AR
	"719":{34.5509, -93.6309}, // Mount Ida, AR
	"719":{34.4607, -93.7123}, // Black Springs, AR
	"719":{34.5024, -93.1420}, // Piney, AR
	"719":{34.4641, -93.1341}, // Rockwell, AR
	"719":{34.2662, -93.4634}, // Amity, AR
	"719":{34.1250, -93.0719}, // Arkadelphia, AR
	// END OF OPTIONS TO PICK
	"720":{34.7255, -92.3580}, // Little Rock, AR
	"721":{34.7255, -92.3580}, // Little Rock, AR
	"722":{34.7255, -92.3580}, // Little Rock, AR
	"723":{35.1046, -89.9773}, // Memphis, TN
	TODO PLACEHOLDER TO FIX ZIP <724> with city NE ARKANSAS and state AR
	"724":{36.4365, -90.3907}, // McDougal, AR
	"724":{35.8592, -90.0260}, // Dell, AR
	"724":{35.9116, -90.8011}, // Bono, AR
	"724":{36.2519, -91.3604}, // Williford, AR
	"724":{36.2642, -90.2936}, // Rector, AR
	"724":{35.7456, -90.5547}, // Bay, AR
	"724":{36.3062, -90.0951}, // Nimmons, AR
	"724":{36.3609, -90.7590}, // Reyno, AR
	"724":{36.0852, -90.9462}, // Walnut Ridge, AR
	"724":{35.7225, -90.2320}, // Etowah, AR
	"724":{35.8987, -90.5764}, // Brookland, AR
	"724":{35.7597, -90.3222}, // Caraway, AR
	"724":{35.7235, -91.2027}, // Tuckerman, AR
	"724":{36.2092, -90.5043}, // Lafe, AR
	"724":{36.1901, -90.3860}, // Marmaduke, AR
	"724":{36.3925, -90.7285}, // Datto, AR
	"724":{36.0556, -90.5148}, // Paragould, AR
	"724":{36.4296, -90.2670}, // Pollard, AR
	"724":{36.4109, -90.5860}, // Corning, AR
	"724":{35.9760, -91.0275}, // Minturn, AR
	"724":{35.8244, -91.1294}, // Swifton, AR
	"724":{35.8013, -90.9313}, // Cash, AR
	"724":{36.0799, -91.3027}, // Smithville, AR
	"724":{36.0360, -90.9736}, // Hoxie, AR
	"724":{36.1070, -91.1074}, // Black Rock, AR
	"724":{35.9665, -91.3211}, // Strawberry, AR
	"724":{36.2021, -91.1823}, // Imboden, AR
	"724":{35.8848, -90.1648}, // Manila, AR
	"724":{35.6678, -91.2538}, // Campbell Station, AR
	"724":{35.6540, -91.0755}, // Grubbs, AR
	"724":{36.4556, -90.1428}, // St. Francis, AR
	"724":{35.8933, -90.3442}, // Monette, AR
	"724":{36.1271, -90.5085}, // Oak Grove Heights, AR
	"724":{36.0849, -91.0711}, // Portia, AR
	"724":{35.8943, -91.0832}, // Alicia, AR
	"724":{36.3127, -91.2233}, // Ravenden Springs, AR
	"724":{35.9226, -90.2541}, // Leachville, AR
	"724":{35.5653, -90.9346}, // Waldenburg, AR
	"724":{36.3858, -90.2016}, // Piggott, AR
	"724":{35.8211, -90.6793}, // Jonesboro, AR
	"724":{35.6197, -90.9054}, // Weiner, AR
	"724":{36.3199, -90.6024}, // Knobel, AR
	"724":{36.2812, -90.6608}, // Peach Orchard, AR
	"724":{35.4915, -90.9725}, // Fisher, AR
	"724":{36.1715, -90.8183}, // O'Kean, AR
	"724":{36.3406, -90.2223}, // Greenway, AR
	"724":{35.6761, -90.5231}, // Trumann, AR
	"724":{35.8677, -90.9454}, // Egypt, AR
	"724":{35.8202, -90.4547}, // Lake City, AR
	"724":{36.2640, -90.9703}, // Pocahontas, AR
	"724":{36.0841, -91.1202}, // Powhatan, AR
	"724":{35.5634, -90.7214}, // Harrisburg, AR
	"724":{36.2304, -90.7256}, // Delaplaine, AR
	"724":{36.4217, -90.9027}, // Maynard, AR
	"724":{36.3319, -90.8052}, // Biggers, AR
	"724":{35.8364, -90.3676}, // Black Oak, AR
	"724":{35.9752, -90.8679}, // Sedgwick, AR
	"724":{36.4545, -90.7226}, // Success, AR
	"724":{36.2417, -91.2511}, // Ravenden, AR
	"724":{36.0042, -91.2520}, // Lynn, AR
	// END OF OPTIONS TO PICK
	"725":{35.7687, -91.6226}, // Batesville, AR
	"726":{36.2438, -93.1198}, // Harrison, AR
	TODO PLACEHOLDER TO FIX ZIP <727> with city NW ARKANSAS and state AR
	"727":{36.4291, -94.3713}, // Gravette, AR
	"727":{36.1642, -94.2457}, // Tontitown, AR
	"727":{36.2263, -94.1282}, // Bethel Heights, AR
	"727":{36.4667, -94.2707}, // Bella Vista, AR
	"727":{36.1458, -93.8627}, // Hindsville, AR
	"727":{36.1844, -94.5318}, // Siloam Springs, AR
	"727":{36.3172, -94.1518}, // Rogers, AR
	"727":{36.2703, -94.2225}, // Cave Springs, AR
	"727":{35.9961, -94.1903}, // Greenland, AR
	"727":{36.3402, -94.4578}, // Decatur, AR
	"727":{35.8251, -93.7652}, // St. Paul, AR
	"727":{36.1328, -94.1757}, // Johnson, AR
	"727":{36.3400, -94.0613}, // Prairie Creek, AR
	"727":{36.3898, -93.9156}, // Lost Bridge Village, AR
	"727":{36.2767, -94.3236}, // Highfill, AR
	"727":{36.4490, -94.1212}, // Pea Ridge, AR
	"727":{36.3559, -94.2972}, // Centerton, AR
	"727":{36.4837, -94.4591}, // Sulphur Springs, AR
	"727":{36.0370, -94.2537}, // Farmington, AR
	"727":{36.3980, -94.0710}, // Avoca, AR
	"727":{36.2563, -94.1518}, // Lowell, AR
	"727":{36.4554, -93.9757}, // Garfield, AR
	"727":{36.2069, -94.2366}, // Elm Springs, AR
	"727":{36.2947, -94.5769}, // Cherokee City, AR
	"727":{36.2606, -94.4239}, // Springtown, AR
	"727":{36.0713, -94.1660}, // Fayetteville, AR
	"727":{36.3855, -94.1361}, // Little Flock, AR
	"727":{36.4019, -94.5889}, // Maysville, AR
	"727":{36.4854, -93.9365}, // Gateway, AR
	"727":{36.1898, -94.1573}, // Springdale, AR
	"727":{35.9360, -94.1798}, // West Fork, AR
	"727":{36.1042, -94.0038}, // Goshen, AR
	"727":{35.9858, -94.3048}, // Prairie Grove, AR
	"727":{36.2570, -94.4907}, // Gentry, AR
	"727":{36.0985, -93.7363}, // Huntsville, AR
	"727":{36.0163, -94.0250}, // Elkins, AR
	"727":{35.9490, -94.4174}, // Lincoln, AR
	"727":{36.3550, -94.2299}, // Bentonville, AR
	// END OF OPTIONS TO PICK
	"728":{35.2763, -93.1383},  // Russellville, AR
	"729":{35.3493, -94.3695},  // Fort Smith, AR
	"730":{35.4676, -97.5137},  // Oklahoma City, OK
	"731":{35.4676, -97.5137},  // Oklahoma City, OK
	"733":{30.3006, -97.7517},  // Austin, TX
	"734":{34.1943, -97.1253},  // Ardmore, OK
	"735":{34.6176, -98.4203},  // Lawton, OK
	"736":{35.5058, -98.9724},  // Clinton, OK
	"737":{36.4061, -97.8701},  // Enid, OK
	"738":{36.4246, -99.4057},  // Woodward, OK
	"739":{37.0466, -100.9295}, // Liberal, KS
	"740":{36.1284, -95.9043},  // Tulsa, OK
	"741":{36.1284, -95.9043},  // Tulsa, OK
	"743":{36.1284, -95.9043},  // Tulsa, OK
	"744":{35.7430, -95.3566},  // Muskogee, OK
	"745":{34.9262, -95.7698},  // McAlester, OK
	"746":{36.7235, -97.0679},  // Ponca City, OK
	"747":{33.9957, -96.3938},  // Durant, OK
	"748":{35.3525, -96.9647},  // Shawnee, OK
	"749":{35.0430, -94.6357},  // Poteau, OK
	TODO PLACEHOLDER TO FIX ZIP <750> with city NORTH TEXAS and state TX
	"750":{32.9170, -96.4377}, // Rockwall, TX
	"750":{33.2016, -96.6669}, // McKinney, TX
	"750":{33.0442, -96.5499}, // St. Paul, TX
	"750":{33.1501, -96.9184}, // Hackberry, TX
	"750":{33.0950, -96.5792}, // Lucas, TX
	"750":{33.6233, -96.7286}, // Southmayd, TX
	"750":{33.8671, -96.6596}, // Preston, TX
	"750":{32.9157, -96.5488}, // Rowlett, TX
	"750":{32.8584, -96.9702}, // Irving, TX
	"750":{32.7936, -96.7662}, // Dallas, TX
	"750":{33.1709, -96.5444}, // Lowry Crossing, TX
	"750":{32.9100, -96.6305}, // Garland, TX
	"750":{32.9726, -96.5793}, // Sachse, TX
	"750":{32.6870, -97.0209}, // Grand Prairie, TX
	"750":{33.1399, -96.6117}, // Fairview, TX
	"750":{33.0927, -96.8977}, // The Colony, TX
	"750":{33.1088, -96.6735}, // Allen, TX
	"750":{33.6266, -96.6195}, // Sherman, TX
	"750":{33.1802, -96.9911}, // Oak Point, TX
	"750":{33.0452, -96.9823}, // Lewisville, TX
	"750":{32.9890, -96.8999}, // Carrollton, TX
	"750":{33.4646, -96.7645}, // Gunter, TX
	"750":{33.0186, -96.6105}, // Murphy, TX
	"750":{33.1114, -97.0313}, // Hickory Creek, TX
	"750":{32.8464, -97.1350}, // Bedford, TX
	"750":{33.3189, -96.7866}, // Celina, TX
	"750":{33.0502, -96.7487}, // Plano, TX
	"750":{32.9638, -96.9905}, // Coppell, TX
	"750":{33.0438, -96.8793}, // Hebron, TX
	"750":{33.0897, -97.0615}, // Highland Village, TX
	"750":{33.1277, -97.0234}, // Lake Dallas, TX
	"750":{33.1378, -96.9749}, // Lakewood Village, TX
	"750":{33.1554, -96.8217}, // Frisco, TX
	"750":{33.0633, -97.1117}, // Double Oak, TX
	"750":{32.9586, -96.8356}, // Addison, TX
	"750":{33.5331, -96.6960}, // Dorchester, TX
	"750":{33.6892, -96.6184}, // Knollwood, TX
	"750":{33.7672, -96.5807}, // Denison, TX
	"750":{33.0961, -97.0975}, // Copper Canyon, TX
	"750":{33.7709, -96.6710}, // Pottsboro, TX
	"750":{32.8443, -96.4679}, // Heath, TX
	"750":{33.0344, -97.1147}, // Flower Mound, TX
	"750":{33.0362, -96.5161}, // Wylie, TX
	"750":{33.2394, -96.8088}, // Prosper, TX
	"750":{32.8512, -96.3924}, // McLendon-Chisholm, TX
	"750":{33.2116, -96.5635}, // New Hope, TX
	"750":{32.9228, -96.4111}, // Mobile City, TX
	"750":{33.0570, -96.6248}, // Parker, TX
	"750":{33.1856, -96.9290}, // Little Elm, TX
	"750":{32.9432, -96.3863}, // Fate, TX
	"750":{33.6821, -96.8487}, // Sadler, TX
	"750":{33.2100, -96.9327}, // Paloma Creek South, TX
	"750":{33.3300, -96.6676}, // Weston, TX
	"750":{32.9717, -96.7092}, // Richardson, TX
	// END OF OPTIONS TO PICK
	"751":{32.7936, -96.7662}, // Dallas, TX
	"752":{32.7936, -96.7662}, // Dallas, TX
	"753":{32.7936, -96.7662}, // Dallas, TX
	"754":{33.1116, -96.1099}, // Greenville, TX
	"755":{33.4487, -94.0815}, // Texarkana, TX
	TODO PLACEHOLDER TO FIX ZIP <756> with city EAST TEXAS and state TX
	"756":{32.2686, -94.9297}, // New London, TX
	"756":{32.5426, -94.9465}, // Gladewater, TX
	"756":{32.3159, -94.5189}, // Tatum, TX
	"756":{33.0306, -94.7250}, // Daingerfield, TX
	"756":{32.3979, -94.8603}, // Kilgore, TX
	"756":{32.8981, -94.5535}, // Avinger, TX
	"756":{32.9986, -94.6310}, // Hughes Springs, TX
	"756":{32.3569, -94.6513}, // Lake Cherokee, TX
	"756":{32.1576, -94.7960}, // Henderson, TX
	"756":{32.3820, -94.5912}, // Easton, TX
	"756":{32.5192, -94.7622}, // Longview, TX
	"756":{32.6006, -94.8527}, // East Mountain, TX
	"756":{32.7634, -94.3511}, // Jefferson, TX
	"756":{32.2436, -94.4561}, // Beckville, TX
	"756":{32.9997, -94.9668}, // Pittsburg, TX
	"756":{32.5886, -94.4468}, // Nesbitt, TX
	"756":{32.7317, -94.9460}, // Gilmer, TX
	"756":{32.5370, -94.3515}, // Marshall, TX
	"756":{31.9113, -94.6826}, // Mount Enterprise, TX
	"756":{32.4766, -94.0646}, // Waskom, TX
	"756":{32.7066, -94.1264}, // Uncertain, TX
	"756":{32.0277, -94.3679}, // Gary City, TX
	"756":{32.5027, -94.5700}, // Hallsville, TX
	"756":{32.5530, -94.9056}, // Warren City, TX
	"756":{32.2760, -94.9726}, // Overton, TX
	"756":{32.5313, -94.8564}, // White Oak, TX
	"756":{32.5329, -94.8946}, // Clarksville City, TX
	"756":{32.1526, -94.3368}, // Carthage, TX
	"756":{32.7723, -94.4992}, // Pine Harbor, TX
	"756":{32.4050, -94.7102}, // Lakeport, TX
	"756":{32.5319, -94.2465}, // Scottsville, TX
	"756":{32.5798, -94.9095}, // Union Grove, TX
	"756":{32.4502, -94.9438}, // Liberty City, TX
	"756":{32.8008, -94.7182}, // Ore City, TX
	"756":{32.9392, -94.7091}, // Lone Star, TX
	"756":{33.0286, -95.0283}, // Rocky Mound, TX
	"756":{32.4360, -93.9636}, // Greenwood, LA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <757> with city EAST TEXAS and state TX
	"757":{31.8662, -94.9833}, // Reklaw, TX
	"757":{31.8968, -95.1520}, // Gallatin, TX
	"757":{32.2750, -95.7561}, // Murchison, TX
	"757":{31.8126, -94.8411}, // Cushing, TX
	"757":{32.0561, -95.5043}, // Frankston, TX
	"757":{32.2222, -95.2217}, // Whitehouse, TX
	"757":{32.4928, -95.1730}, // Winona, TX
	"757":{32.5860, -95.1127}, // Big Sandy, TX
	"757":{32.2041, -95.8321}, // Athens, TX
	"757":{32.7125, -95.2003}, // Holly Lake Ranch, TX
	"757":{32.0793, -95.5924}, // Poynor, TX
	"757":{32.3743, -95.6099}, // Edom, TX
	"757":{32.5917, -95.2027}, // Hawkins, TX
	"757":{32.2985, -95.6130}, // Brownsboro, TX
	"757":{32.1451, -95.3196}, // Bullard, TX
	"757":{31.7978, -95.1491}, // Rusk, TX
	"757":{32.6461, -95.4775}, // Mineola, TX
	"757":{32.0876, -95.4666}, // Berryville, TX
	"757":{31.9807, -95.1150}, // New Summerfield, TX
	"757":{32.4890, -95.4578}, // Hideaway, TX
	"757":{32.2279, -95.0536}, // Arp, TX
	"757":{32.1383, -95.4871}, // Coffee City, TX
	"757":{32.3065, -95.4781}, // Chandler, TX
	"757":{32.4933, -95.4069}, // Lindale, TX
	"757":{32.7951, -95.4443}, // Quitman, TX
	"757":{31.9642, -95.2617}, // Jacksonville, TX
	"757":{32.1602, -95.4380}, // Emerald Bay, TX
	"757":{32.0374, -95.4147}, // Cuney, TX
	"757":{32.3184, -95.3065}, // Tyler, TX
	"757":{32.2441, -95.3972}, // Noonday, TX
	"757":{32.5242, -95.6373}, // Van, TX
	"757":{32.1108, -95.4212}, // Shadybrook, TX
	"757":{32.3675, -95.6983}, // Callender Lake, TX
	"757":{32.3012, -95.1674}, // New Chapel Hill, TX
	"757":{32.1450, -95.1229}, // Troup, TX
	"757":{32.1905, -95.5704}, // Moore Station, TX
	// END OF OPTIONS TO PICK
	"758":{31.7544, -95.6471}, // Palestine, TX
	"759":{31.3217, -94.7277}, // Lufkin, TX
	TODO PLACEHOLDER TO FIX ZIP <760> with city FT WORTH and state TX
	"760":{32.4885, -97.8357}, // Oak Trail Shores, TX
	"760":{32.7508, -97.6999}, // Hudson Oaks, TX
	"760":{32.4311, -97.1021}, // Venus, TX
	"760":{32.4583, -97.3851}, // Joshua, TX
	"760":{32.4484, -97.7685}, // Granbury, TX
	"760":{32.3628, -97.6554}, // Pecan Plantation, TX
	"760":{32.3528, -97.2944}, // Coyote Flats, TX
	"760":{32.2685, -97.1775}, // Grandview, TX
	"760":{32.8508, -97.0799}, // Euless, TX
	"760":{32.3140, -97.0053}, // Maypearl, TX
	"760":{32.6936, -97.1565}, // Dalworthington Gardens, TX
	"760":{32.3561, -97.4145}, // Cleburne, TX
	"760":{33.0049, -97.4856}, // Newark, TX
	"760":{32.8001, -98.0113}, // Cool, TX
	"760":{33.0647, -97.4779}, // Rhome, TX
	"760":{32.9884, -97.5528}, // Briar, TX
	"760":{33.1122, -97.4489}, // New Fairview, TX
	"760":{32.5690, -97.1211}, // Mansfield, TX
	"760":{32.5752, -97.8782}, // Horseshoe Bend, TX
	"760":{32.2461, -97.7441}, // Glen Rose, TX
	"760":{33.0560, -97.5096}, // Aurora, TX
	"760":{32.9545, -97.1503}, // Southlake, TX
	"760":{32.4958, -97.3033}, // Briaroaks, TX
	"760":{32.7147, -97.1540}, // Pantego, TX
	"760":{32.4513, -97.5322}, // Godley, TX
	"760":{32.8169, -98.0776}, // Mineral Wells, TX
	"760":{32.7535, -97.7724}, // Weatherford, TX
	"760":{32.4752, -97.7611}, // Brazos Bend, TX
	"760":{32.5780, -97.3584}, // Crowley, TX
	"760":{32.4280, -97.6911}, // DeCordova, TX
	"760":{32.2351, -97.3744}, // Rio Vista, TX
	"760":{32.8353, -97.1808}, // Hurst, TX
	"760":{32.4832, -97.3260}, // Cross Timber, TX
	"760":{32.6743, -97.6483}, // Annetta South, TX
	"760":{32.7548, -97.6499}, // Willow Park, TX
	"760":{32.8464, -97.1350}, // Bedford, TX
	"760":{32.5789, -97.2350}, // Rendon, TX
	"760":{32.9610, -97.3382}, // Haslet, TX
	"760":{32.6998, -97.1251}, // Arlington, TX
	"760":{32.8913, -97.1486}, // Colleyville, TX
	"760":{32.5299, -97.6155}, // Cresson, TX
	"760":{32.6434, -97.2173}, // Kennedale, TX
	"760":{32.6215, -97.8160}, // Western Lake, TX
	"760":{32.3927, -97.7406}, // Canyon Creek, TX
	"760":{32.7493, -98.0117}, // Millsap, TX
	"760":{32.9695, -97.6805}, // Springtown, TX
	"760":{33.0843, -97.5632}, // Boyd, TX
	"760":{32.9132, -97.5877}, // Sanctuary, TX
	"760":{32.3955, -97.3225}, // Keene, TX
	"760":{32.5170, -97.3343}, // Burleson, TX
	"760":{32.4718, -96.9877}, // Midlothian, TX
	"760":{32.4068, -97.2150}, // Alvarado, TX
	"760":{33.1504, -97.6887}, // Paradise, TX
	"760":{32.6938, -97.6581}, // Annetta, TX
	"760":{32.7193, -97.6723}, // Annetta North, TX
	"760":{32.9227, -97.5189}, // Pelican Bay, TX
	"760":{32.1586, -97.1478}, // Itasca, TX
	"760":{32.9703, -97.4727}, // Pecan Acres, TX
	"760":{32.8955, -97.5379}, // Azle, TX
	"760":{32.6974, -97.6070}, // Aledo, TX
	"760":{32.7812, -97.3472}, // Fort Worth, TX
	"760":{32.9343, -97.0744}, // Grapevine, TX
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <761> with city FT WORTH and state TX
	"761":{32.8605, -97.2180}, // North Richland Hills, TX
	"761":{32.8130, -97.4306}, // Lake Worth, TX
	"761":{32.8176, -97.2707}, // Haltom City, TX
	"761":{32.7554, -97.4605}, // White Settlement, TX
	"761":{32.6787, -97.4638}, // Benbrook, TX
	"761":{32.8543, -97.3383}, // Blue Mound, TX
	"761":{32.8718, -97.2515}, // Watauga, TX
	"761":{32.6296, -97.2827}, // Everman, TX
	"761":{32.7767, -97.3984}, // River Oaks, TX
	"761":{32.6561, -97.3406}, // Edgecliff Village, TX
	"761":{32.5789, -97.2350}, // Rendon, TX
	"761":{32.9610, -97.3382}, // Haslet, TX
	"761":{32.8028, -97.4021}, // Sansom Park, TX
	"761":{32.7436, -97.4123}, // Westover Hills, TX
	"761":{32.8657, -97.3652}, // Saginaw, TX
	"761":{32.7599, -97.4239}, // Westworth Village, TX
	"761":{32.6619, -97.2662}, // Forest Hill, TX
	"761":{32.9703, -97.4727}, // Pecan Acres, TX
	"761":{32.8095, -97.2273}, // Richland Hills, TX
	"761":{32.8219, -97.4869}, // Lakeside, TX
	"761":{32.7812, -97.3472}, // Fort Worth, TX
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <762> with city FT WORTH and state TX
	"762":{33.5703, -97.0129}, // Lake Kiowa, TX
	"762":{33.7835, -97.7302}, // Nocona, TX
	"762":{33.1330, -97.3014}, // DISH, TX
	"762":{33.2253, -96.9371}, // Paloma Creek, TX
	"762":{33.6652, -97.7210}, // Montague, TX
	"762":{33.6414, -97.2186}, // Lindsay, TX
	"762":{33.0795, -97.1506}, // Bartonville, TX
	"762":{33.1776, -97.2909}, // Ponder, TX
	"762":{33.6233, -96.7286}, // Southmayd, TX
	"762":{33.6987, -97.0159}, // Callisburg, TX
	"762":{33.2176, -97.1419}, // Denton, TX
	"762":{33.2269, -97.5886}, // Decatur, TX
	"762":{33.6334, -98.0161}, // Bellevue, TX
	"762":{33.1434, -97.0681}, // Corinth, TX
	"762":{33.0920, -97.1216}, // Lantana, TX
	"762":{33.0647, -97.4779}, // Rhome, TX
	"762":{33.6586, -97.3874}, // Muenster, TX
	"762":{33.0802, -97.2545}, // Northlake, TX
	"762":{32.9810, -97.2038}, // Westlake, TX
	"762":{33.4479, -97.7663}, // Sunset, TX
	"762":{33.5566, -97.8440}, // Bowie, TX
	"762":{33.6950, -97.5231}, // St. Jo, TX
	"762":{33.4646, -96.7645}, // Gunter, TX
	"762":{33.5594, -96.9071}, // Collinsville, TX
	"762":{33.2652, -97.2257}, // Krum, TX
	"762":{33.6612, -96.9022}, // Whitesboro, TX
	"762":{33.0148, -97.2268}, // Roanoke, TX
	"762":{33.8478, -96.8147}, // Sherwood Shores, TX
	"762":{33.0039, -97.1898}, // Trophy Club, TX
	"762":{33.4892, -97.1534}, // Valley View, TX
	"762":{33.2290, -96.9985}, // Cross Roads, TX
	"762":{33.2363, -96.9611}, // Providence Village, TX
	"762":{33.1105, -97.1864}, // Argyle, TX
	"762":{33.0633, -97.1117}, // Double Oak, TX
	"762":{33.0961, -97.0975}, // Copper Canyon, TX
	"762":{33.6391, -97.1487}, // Gainesville, TX
	"762":{33.3955, -96.9501}, // Pilot Point, TX
	"762":{33.4718, -96.9187}, // Tioga, TX
	"762":{33.0344, -97.1147}, // Flower Mound, TX
	"762":{33.0985, -97.2304}, // Corral City, TX
	"762":{33.2792, -96.9880}, // Krugerville, TX
	"762":{32.9337, -97.2255}, // Keller, TX
	"762":{33.3569, -97.6960}, // Alvord, TX
	"762":{33.3712, -97.1677}, // Sanger, TX
	"762":{33.2872, -96.9623}, // Aubrey, TX
	"762":{33.1856, -96.9290}, // Little Elm, TX
	"762":{33.2257, -96.9081}, // Savannah, TX
	"762":{33.0863, -97.3009}, // Justin, TX
	"762":{32.7812, -97.3472}, // Fort Worth, TX
	"762":{33.6821, -96.8487}, // Sadler, TX
	"762":{33.1627, -97.0394}, // Shady Shores, TX
	"762":{33.8524, -97.6435}, // Nocona Hills, TX
	// END OF OPTIONS TO PICK
	"763":{33.9072, -98.5290}, // Wichita Falls, TX
	TODO PLACEHOLDER TO FIX ZIP <764> with city FT WORTH and state TX
	"764":{32.5509, -98.4979}, // Strawn, TX
	"764":{31.8249, -98.7895}, // Blanket, TX
	"764":{32.2136, -98.6721}, // Gorman, TX
	"764":{33.2099, -97.7709}, // Bridgeport, TX
	"764":{32.5472, -99.1666}, // Moran, TX
	"764":{33.1636, -98.3891}, // Bryson, TX
	"764":{33.2074, -97.8311}, // Lake Bridgeport, TX
	"764":{32.2684, -98.8272}, // Carbon, TX
	"764":{32.5364, -98.4249}, // Mingus, TX
	"764":{32.0875, -98.3391}, // Dublin, TX
	"764":{32.1271, -99.1658}, // Cross Plains, TX
	"764":{33.0348, -98.0693}, // Perrin, TX
	"764":{32.5455, -98.3672}, // Gordon, TX
	"764":{33.1773, -97.8710}, // Runaway Bay, TX
	"764":{31.9005, -98.6044}, // Comanche, TX
	"764":{32.5187, -98.0470}, // Lipan, TX
	"764":{32.3894, -97.9191}, // Tolar, TX
	"764":{33.2235, -98.1589}, // Jacksboro, TX
	"764":{32.0976, -98.9665}, // Rising Star, TX
	"764":{32.3848, -98.9805}, // Cisco, TX
	"764":{32.7566, -98.9125}, // Breckenridge, TX
	"764":{32.7692, -98.3008}, // Palo Pinto, TX
	"764":{31.9858, -98.0290}, // Hico, TX
	"764":{33.1820, -99.1798}, // Throckmorton, TX
	"764":{32.4693, -98.6751}, // Ranger, TX
	"764":{32.7273, -99.2956}, // Albany, TX
	"764":{32.3704, -99.1958}, // Putnam, TX
	"764":{32.4030, -98.8173}, // Eastland, TX
	"764":{33.2960, -97.7987}, // Chico, TX
	"764":{33.1017, -98.5778}, // Graham, TX
	"764":{32.9374, -98.2475}, // Graford, TX
	"764":{31.8457, -98.4025}, // Gustine, TX
	"764":{33.0149, -99.0534}, // Woodson, TX
	"764":{32.2148, -98.2205}, // Stephenville, TX
	"764":{32.1114, -98.5351}, // De Leon, TX
	// END OF OPTIONS TO PICK
	"765":{31.5597, -97.1882},  // Waco, TX
	"766":{31.5597, -97.1882},  // Waco, TX
	"767":{31.5597, -97.1882},  // Waco, TX
	"768":{32.4543, -99.7384},  // Abilene, TX
	"769":{32.0249, -102.1137}, // Midland, TX
	"770":{29.7869, -95.3905},  // Houston, TX
	"771":{29.7869, -95.3905},  // Houston, TX
	"772":{29.7869, -95.3905},  // Houston, TX
	TODO PLACEHOLDER TO FIX ZIP <773> with city NORTH HOUSTON and state TX
	"773":{30.1811, -95.1550}, // Roman Forest, TX
	"773":{30.1501, -95.3217}, // Porter Heights, TX
	"773":{30.1570, -95.4421}, // Oak Ridge North, TX
	"773":{30.6093, -94.9467}, // Goodrich, TX
	"773":{29.9921, -95.2655}, // Humble, TX
	"773":{30.1840, -95.4556}, // Shenandoah, TX
	"773":{30.1956, -95.1697}, // Patton Village, TX
	"773":{30.2612, -95.8298}, // Todd Mission, TX
	"773":{30.2007, -95.0958}, // Plum Grove, TX
	"773":{30.5898, -95.1307}, // Coldspring, TX
	"773":{30.0964, -95.6186}, // Tomball, TX
	"773":{30.7100, -94.9381}, // Livingston, TX
	"773":{30.7975, -95.0749}, // Cedar Point, TX
	"773":{30.0613, -95.3830}, // Spring, TX
	"773":{30.2321, -95.1617}, // Splendora, TX
	"773":{30.3404, -95.3540}, // Cut and Shoot, TX
	"773":{30.6879, -94.7449}, // Indian Springs, TX
	"773":{30.1438, -95.7079}, // Stagecoach, TX
	"773":{30.4314, -95.4832}, // Willis, TX
	"773":{30.7050, -95.5545}, // Huntsville, TX
	"773":{30.2117, -95.7419}, // Magnolia, TX
	"773":{30.1814, -95.1835}, // Woodbranch, TX
	"773":{30.2172, -95.4130}, // Woodloch, TX
	"773":{30.4910, -95.0021}, // Shepherd, TX
	"773":{30.5374, -95.4822}, // New Waverly, TX
	"773":{30.3577, -95.1003}, // North Cleveland, TX
	"773":{30.8517, -94.8580}, // Seven Oaks, TX
	"773":{30.3224, -95.4820}, // Conroe, TX
	"773":{30.8471, -95.3985}, // Riverside, TX
	"773":{30.7473, -95.2163}, // Point Blank, TX
	"773":{30.6531, -95.1265}, // Cape Royale, TX
	"773":{29.9777, -95.1953}, // Atascocita, TX
	"773":{29.7869, -95.3905}, // Houston, TX
	"773":{30.1738, -95.5134}, // The Woodlands, TX
	"773":{30.3806, -95.4944}, // Panorama Village, TX
	"773":{30.3917, -95.6965}, // Montgomery, TX
	"773":{30.4874, -94.7678}, // Big Thicket Lake Estates, TX
	"773":{30.6957, -95.0097}, // West Livingston, TX
	"773":{30.8209, -95.1110}, // Onalaska, TX
	"773":{30.7434, -95.3072}, // Oakhurst, TX
	"773":{30.3369, -95.0924}, // Cleveland, TX
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <774> with city NORTH HOUSTON and state TX
	"774":{29.0455, -95.5673}, // Brazoria, TX
	"774":{30.0996, -96.0781}, // Hempstead, TX
	"774":{29.5460, -95.8220}, // Rosenberg, TX
	"774":{28.8727, -96.2172}, // Blessing, TX
	"774":{29.8190, -95.9761}, // Pattison, TX
	"774":{29.6930, -95.8792}, // Fulshear, TX
	"774":{29.7826, -95.9545}, // Brookshire, TX
	"774":{29.1119, -96.4123}, // Louise, TX
	"774":{29.7965, -96.1051}, // San Felipe, TX
	"774":{29.6317, -96.0638}, // Wallis, TX
	"774":{29.0278, -95.8800}, // Van Vleck, TX
	"774":{29.5935, -95.6357}, // Sugar Land, TX
	"774":{29.7395, -95.7606}, // Cinco Ranch, TX
	"774":{28.9838, -95.9601}, // Bay City, TX
	"774":{29.3138, -96.1044}, // Wharton, TX
	"774":{29.4946, -95.9146}, // Beasley, TX
	"774":{29.5835, -95.7965}, // Cumings, TX
	"774":{29.9472, -96.2597}, // Bellville, TX
	"774":{28.9627, -96.0645}, // Markham, TX
	"774":{30.0629, -95.9221}, // Waller, TX
	"774":{29.5982, -95.5520}, // Fifth Street, TX
	"774":{29.1986, -96.2725}, // El Campo, TX
	"774":{29.1421, -95.6489}, // West Columbia, TX
	"774":{29.5630, -95.5365}, // Missouri City, TX
	"774":{29.4040, -96.0837}, // Hungerford, TX
	"774":{29.6824, -95.9910}, // Simonton, TX
	"774":{29.6627, -95.9353}, // Weston Lakes, TX
	"774":{29.6705, -95.6597}, // Four Corners, TX
	"774":{29.5242, -96.0622}, // East Bernard, TX
	"774":{29.7674, -96.1679}, // Sealy, TX
	"774":{28.6973, -95.9667}, // Matagorda, TX
	"774":{29.0464, -95.6986}, // Sweeny, TX
	"774":{29.4906, -95.6304}, // Thompsons, TX
	"774":{29.5879, -96.3282}, // Eagle Lake, TX
	"774":{29.6002, -95.9694}, // Orchard, TX
	"774":{29.7547, -96.0424}, // Brazos Country, TX
	"774":{29.2832, -95.7408}, // Damon, TX
	"774":{29.7040, -95.4621}, // Bellaire, TX
	"774":{29.2665, -95.9630}, // Iago, TX
	"774":{29.4834, -95.5065}, // Sienna Plantation, TX
	"774":{29.2584, -95.9436}, // Boling, TX
	"774":{29.6235, -95.7331}, // Pecan Grove, TX
	"774":{30.0554, -96.0253}, // Pine Island, TX
	"774":{29.7905, -95.8353}, // Katy, TX
	"774":{29.4852, -95.8097}, // Pleak, TX
	"774":{29.3957, -95.8386}, // Needville, TX
	"774":{29.5824, -95.7602}, // Richmond, TX
	"774":{29.6271, -95.5653}, // Stafford, TX
	"774":{29.7869, -95.3905}, // Houston, TX
	"774":{28.7198, -96.2351}, // Palacios, TX
	"774":{30.0850, -95.9897}, // Prairie View, TX
	"774":{29.4398, -95.7777}, // Fairchilds, TX
	"774":{29.6513, -95.5873}, // Meadows Place, TX
	"774":{29.4470, -96.0010}, // Kendleton, TX
	"774":{29.0811, -95.6373}, // Wild Peach Village, TX
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <775> with city NORTH HOUSTON and state TX
	"775":{29.0516, -95.4521}, // Lake Jackson, TX
	"775":{28.9978, -95.3282}, // Oyster Creek, TX
	"775":{30.0470, -94.7908}, // Liberty, TX
	"775":{30.1464, -94.6422}, // Hull, TX
	"775":{29.4793, -95.3658}, // Manvel, TX
	"775":{29.8671, -95.0537}, // Barrett, TX
	"775":{29.4128, -94.9658}, // Texas City, TX
	"775":{29.6611, -95.2285}, // South Houston, TX
	"775":{29.1905, -94.9801}, // Jamaica Beach, TX
	"775":{29.7269, -94.8549}, // Beach City, TX
	"775":{29.3500, -95.4530}, // Rosharon, TX
	"775":{30.1135, -94.6426}, // Daisetta, TX
	"775":{29.5907, -95.3180}, // Brookside Village, TX
	"775":{29.6204, -95.0188}, // Shoreacres, TX
	"775":{30.1069, -94.8502}, // Kenefick, TX
	"775":{29.1527, -95.4977}, // Bailey's Prairie, TX
	"775":{29.5585, -95.3215}, // Pearland, TX
	"775":{30.0287, -94.5860}, // Devers, TX
	"775":{29.3257, -94.9390}, // Bayou Vista, TX
	"775":{29.0257, -95.3975}, // Clute, TX
	"775":{29.2274, -95.3461}, // Danbury, TX
	"775":{29.6584, -95.1499}, // Pasadena, TX
	"775":{29.6690, -95.0483}, // La Porte, TX
	"775":{29.5457, -95.0328}, // Clear Lake Shores, TX
	"775":{29.5733, -95.0440}, // El Lago, TX
	"775":{29.5317, -95.1188}, // Webster, TX
	"775":{29.5630, -95.5365}, // Missouri City, TX
	"775":{29.2083, -95.5146}, // Holiday Lakes, TX
	"775":{29.7452, -95.2333}, // Galena Park, TX
	"775":{29.6765, -95.0027}, // Morgan's Point, TX
	"775":{29.5077, -94.9880}, // Bacliff, TX
	"775":{29.5764, -95.0562}, // Taylor Lake Village, TX
	"775":{29.4013, -95.4816}, // Sandy Point, TX
	"775":{29.8122, -94.8201}, // Cove, TX
	"775":{28.9454, -95.3602}, // Freeport, TX
	"775":{29.1716, -95.4292}, // Angleton, TX
	"775":{30.1446, -94.8209}, // Dayton Lakes, TX
	"775":{28.9560, -95.2836}, // Surfside Beach, TX
	"775":{29.8524, -94.8785}, // Mont Belvieu, TX
	"775":{29.3690, -94.9957}, // La Marque, TX
	"775":{29.4834, -95.5065}, // Sienna Plantation, TX
	"775":{29.7649, -94.6787}, // Anahuac, TX
	"775":{30.0293, -94.9039}, // Dayton, TX
	"775":{30.1493, -94.7377}, // Hardin, TX
	"775":{28.9752, -95.4714}, // Jones Creek, TX
	"775":{29.2945, -95.0250}, // Hitchcock, TX
	"775":{29.3870, -95.2936}, // Alvin, TX
	"775":{29.4405, -95.4168}, // Iowa Colony, TX
	"775":{29.4901, -94.9403}, // San Leon, TX
	"775":{29.5307, -95.0194}, // Kemah, TX
	"775":{29.3891, -95.1005}, // Santa Fe, TX
	"775":{29.4545, -95.0584}, // Dickinson, TX
	"775":{29.5112, -95.1979}, // Friendswood, TX
	"775":{29.9146, -95.0591}, // Crosby, TX
	"775":{29.2985, -94.9175}, // Tiki Island, TX
	"775":{29.8130, -95.0577}, // Highlands, TX
	"775":{29.2487, -94.8910}, // Galveston, TX
	"775":{29.7869, -95.3905}, // Houston, TX
	"775":{29.6898, -95.1151}, // Deer Park, TX
	"775":{29.7914, -95.1145}, // Channelview, TX
	"775":{29.3924, -95.2233}, // Hillcrest, TX
	"775":{29.3013, -95.4508}, // Bonney, TX
	"775":{29.6632, -94.6882}, // Oak Island, TX
	"775":{28.9280, -95.3170}, // Quintana, TX
	"775":{29.5357, -95.4693}, // Fresno, TX
	"775":{29.7586, -94.9669}, // Baytown, TX
	"775":{29.5751, -95.0236}, // Seabrook, TX
	"775":{29.0723, -95.4045}, // Richwood, TX
	"775":{30.0451, -94.7373}, // Ames, TX
	"775":{29.8745, -94.8268}, // Old River-Winfree, TX
	"775":{29.3018, -95.2744}, // Liverpool, TX
	"775":{29.4874, -95.1087}, // League City, TX
	"775":{29.5032, -95.4693}, // Arcola, TX
	// END OF OPTIONS TO PICK
	"776":{30.0850, -94.1451},  // Beaumont, TX
	"777":{30.0850, -94.1451},  // Beaumont, TX
	"778":{30.6657, -96.3668},  // Bryan, TX
	"779":{28.8285, -96.9850},  // Victoria, TX
	"780":{29.4658, -98.5254},  // San Antonio, TX
	"781":{29.4658, -98.5254},  // San Antonio, TX
	"782":{29.4658, -98.5254},  // San Antonio, TX
	"783":{27.7261, -97.3755},  // Corpus Christi, TX
	"784":{27.7261, -97.3755},  // Corpus Christi, TX
	"785":{26.2273, -98.2471},  // McAllen, TX
	"786":{30.3006, -97.7517},  // Austin, TX
	"787":{30.3006, -97.7517},  // Austin, TX
	"788":{29.4658, -98.5254},  // San Antonio, TX
	"789":{30.3006, -97.7517},  // Austin, TX
	"790":{35.1989, -101.8310}, // Amarillo, TX
	"791":{35.1989, -101.8310}, // Amarillo, TX
	"792":{34.4293, -100.2516}, // Childress, TX
	"793":{33.5642, -101.8871}, // Lubbock, TX
	"794":{33.5642, -101.8871}, // Lubbock, TX
	"795":{32.4543, -99.7384},  // Abilene, TX
	"796":{32.4543, -99.7384},  // Abilene, TX
	"797":{32.0249, -102.1137}, // Midland, TX
	"798":{31.8479, -106.4309}, // El Paso, TX
	"799":{31.8479, -106.4309}, // El Paso, TX
	"800":{39.7621, -104.8759}, // Denver, CO
	"801":{39.7621, -104.8759}, // Denver, CO
	"802":{39.7621, -104.8759}, // Denver, CO
	"803":{40.0249, -105.2523}, // Boulder, CO
	"804":{39.7621, -104.8759}, // Denver, CO
	"805":{40.1690, -105.0996}, // Longmont, CO
	"806":{39.7621, -104.8759}, // Denver, CO
	"807":{39.7621, -104.8759}, // Denver, CO
	TODO PLACEHOLDER TO FIX ZIP <808> with city COLORADO SPGS and state CO
	"808":{39.3038, -102.4234}, // Bethune, CO
	"808":{38.9450, -105.1619}, // Divide, CO
	"808":{38.7591, -105.5024}, // Guffey, CO
	"808":{38.9343, -105.0237}, // Green Mountain Falls, CO
	"808":{38.9942, -104.8639}, // Air Force Academy, CO
	"808":{38.9987, -105.0595}, // Woodland Park, CO
	"808":{39.1222, -104.1673}, // Ramah, CO
	"808":{39.3042, -102.2714}, // Burlington, CO
	"808":{38.7461, -105.1840}, // Cripple Creek, CO
	"808":{38.7629, -102.7954}, // Kit Carson, CO
	"808":{38.8576, -104.9127}, // Manitou Springs, CO
	"808":{39.1361, -103.4735}, // Hugo, CO
	"808":{39.2980, -102.8695}, // Seibert, CO
	"808":{38.9445, -105.2899}, // Florissant, CO
	"808":{38.8469, -105.1577}, // Midland, CO
	"808":{39.2652, -103.6852}, // Limon, CO
	"808":{38.8256, -104.3829}, // Ellicott, CO
	"808":{39.2783, -103.4989}, // Genoa, CO
	"808":{38.6888, -104.6829}, // Fountain, CO
	"808":{39.3023, -102.7435}, // Vona, CO
	"808":{39.2955, -103.0767}, // Flagler, CO
	"808":{39.2841, -103.2739}, // Arriba, CO
	"808":{38.9433, -105.0025}, // Cascade-Chipita Park, CO
	"808":{39.6127, -102.5920}, // Kirk, CO
	"808":{38.8192, -102.3521}, // Cheyenne Wells, CO
	"808":{39.0330, -104.4904}, // Peyton, CO
	"808":{39.6559, -102.6786}, // Joes, CO
	"808":{38.7090, -105.1419}, // Victor, CO
	"808":{39.0345, -104.2991}, // Calhan, CO
	"808":{39.3029, -102.6035}, // Stratton, CO
	"808":{39.1408, -104.0824}, // Simla, CO
	"808":{38.7177, -105.1254}, // Goldfield, CO
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <809> with city COLORADO SPGS and state CO
	"809":{39.0453, -104.8288}, // Gleneagle, CO
	"809":{38.7731, -104.7787}, // Stratmoor, CO
	"809":{38.8674, -104.7606}, // Colorado Springs, CO
	"809":{38.7011, -104.8346}, // Rock Creek Park, CO
	"809":{38.9433, -105.0025}, // Cascade-Chipita Park, CO
	"809":{39.0608, -104.6752}, // Black Forest, CO
	"809":{38.7489, -104.7142}, // Security-Widefield, CO
	"809":{38.7403, -104.7841}, // Fort Carson, CO
	"809":{38.8597, -104.6995}, // Cimarron Hills, CO
	"809":{39.0736, -104.8468}, // Monument, CO
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <810> with city COLORADO SPGS and state CO
	"810":{38.2351, -104.3434}, // Avondale, CO
	"810":{38.0014, -103.5549}, // La Junta Gardens, CO
	"810":{37.3861, -102.2801}, // Walsh, CO
	"810":{38.1063, -102.9560}, // Hasty, CO
	"810":{38.2209, -103.7567}, // Ordway, CO
	"810":{37.1225, -104.7396}, // Segundo, CO
	"810":{38.4464, -102.4412}, // Brandon, CO
	"810":{37.2349, -104.4492}, // El Moro, CO
	"810":{37.9796, -103.5473}, // La Junta, CO
	"810":{38.4705, -102.0804}, // Towner, CO
	"810":{38.0740, -102.6155}, // Lamar, CO
	"810":{37.1445, -104.6217}, // Cokedale, CO
	"810":{37.2471, -103.3533}, // Kim, CO
	"810":{37.1046, -102.5787}, // Campo, CO
	"810":{37.1605, -105.0342}, // Stonewall Gap, CO
	"810":{37.9365, -104.8459}, // Colorado City, CO
	"810":{37.4036, -104.6550}, // Aguilar, CO
	"810":{38.0554, -102.1247}, // Holly, CO
	"810":{37.6307, -104.7818}, // Walsenburg, CO
	"810":{37.1750, -104.4908}, // Trinidad, CO
	"810":{37.4049, -102.6189}, // Springfield, CO
	"810":{38.2328, -103.6634}, // Sugar City, CO
	"810":{37.1239, -104.6794}, // Valdez, CO
	"810":{38.0142, -103.6282}, // Swink, CO
	"810":{38.2469, -104.5692}, // Blende, CO
	"810":{38.1935, -103.8598}, // Crowley, CO
	"810":{38.2447, -104.4599}, // Vineland, CO
	"810":{38.2389, -104.5881}, // Salt Creek, CO
	"810":{38.4813, -102.7798}, // Eads, CO
	"810":{38.1211, -102.2216}, // Hartman, CO
	"810":{38.1088, -103.8669}, // Manzanola, CO
	"810":{38.2493, -104.2579}, // Boone, CO
	"810":{38.2713, -104.6105}, // Pueblo, CO
	"810":{37.2816, -104.3890}, // Hoehne, CO
	"810":{37.9214, -104.9322}, // Rye, CO
	"810":{38.1554, -102.7193}, // Wiley, CO
	"810":{37.9985, -103.5228}, // North La Junta, CO
	"810":{37.1167, -104.5232}, // Starkville, CO
	"810":{38.0664, -104.9785}, // Beulah Valley, CO
	"810":{37.5086, -105.0086}, // La Veta, CO
	"810":{38.0499, -103.7227}, // Rocky Ford, CO
	"810":{38.4667, -102.2941}, // Sheridan Lake, CO
	"810":{37.3700, -102.8587}, // Pritchett, CO
	"810":{38.1663, -103.9445}, // Olney Springs, CO
	"810":{37.0156, -103.8838}, // Branson, CO
	"810":{38.0630, -102.3117}, // Granada, CO
	"810":{37.1581, -104.5501}, // Jansen, CO
	"810":{38.0695, -103.2236}, // Las Animas, CO
	"810":{38.1296, -104.0255}, // Fowler, CO
	"810":{37.5607, -102.3966}, // Two Buttes, CO
	"810":{38.3551, -104.7266}, // Pueblo West, CO
	"810":{38.4525, -103.1649}, // Haswell, CO
	"810":{37.3737, -102.4474}, // Vilas, CO
	"810":{37.4219, -104.6431}, // Lynn, CO
	"810":{37.1459, -104.8684}, // Weston, CO
	"810":{38.1078, -103.5112}, // Cheraw, CO
	// END OF OPTIONS TO PICK
	"811":{37.4755, -105.8770}, // Alamosa, CO
	"812":{38.5300, -105.9984}, // Salida, CO
	"813":{37.2744, -107.8703}, // Durango, CO
	"814":{39.0877, -108.5673}, // Grand Junction, CO
	"815":{39.0877, -108.5673}, // Grand Junction, CO
	"816":{39.5455, -107.3347}, // Glenwood Springs, CO
	"820":{41.1405, -104.7927}, // Cheyenne, WY
	TODO PLACEHOLDER TO FIX ZIP <821> with city YELLOWSTONE NL PK and state WY
	"821":{44.9733, -110.6930}, // Mammoth, WY
	// END OF OPTIONS TO PICK
	"822":{42.0516, -104.9595}, // Wheatland, WY
	"823":{41.7849, -107.2265}, // Rawlins, WY
	"824":{44.0026, -107.9543}, // Worland, WY
	"825":{43.0425, -108.4142}, // Riverton, WY
	"826":{42.8420, -106.3207}, // Casper, WY
	"827":{44.2752, -105.4984}, // Gillette, WY
	"828":{44.7962, -106.9643}, // Sheridan, WY
	"829":{41.5951, -109.2237}, // Rock Springs, WY
	"830":{41.5951, -109.2237}, // Rock Springs, WY
	"831":{41.5951, -109.2237}, // Rock Springs, WY
	"832":{42.8716, -112.4652}, // Pocatello, ID
	"833":{42.5648, -114.4617}, // Twin Falls, ID
	"834":{42.8716, -112.4652}, // Pocatello, ID
	"835":{46.3934, -116.9934}, // Lewiston, ID
	"836":{43.6007, -116.2312}, // Boise, ID
	"837":{43.6007, -116.2312}, // Boise, ID
	"838":{47.6671, -117.4330}, // Spokane, WA
	TODO PLACEHOLDER TO FIX ZIP <840> with city SALT LAKE CTY and state UT
	"840":{40.4743, -111.9383}, // Bluffdale, UT
	"840":{41.1717, -112.0480}, // Roy, UT
	"840":{40.6813, -111.2829}, // Marion, UT
	"840":{40.9846, -111.9064}, // Farmington, UT
	"840":{40.3695, -109.3556}, // Jensen, UT
	"840":{40.8983, -111.9080}, // West Bountiful, UT
	"840":{40.3448, -111.9153}, // Saratoga Springs, UT
	"840":{41.0864, -112.0697}, // Syracuse, UT
	"840":{40.3414, -112.1086}, // Cedar Fort, UT
	"840":{40.4733, -111.2533}, // Timber Lakes, UT
	"840":{40.4721, -109.9413}, // Whiterocks, UT
	"840":{40.9152, -111.3942}, // Coalville, UT
	"840":{40.6318, -111.2244}, // Samak, UT
	"840":{40.5666, -111.8636}, // White City, UT
	"840":{40.4344, -110.0308}, // Neola, UT
	"840":{40.5183, -111.4745}, // Midway, UT
	"840":{40.4496, -112.3672}, // Stockton, UT
	"840":{40.4898, -112.0169}, // Herriman, UT
	"840":{40.6137, -111.8144}, // Cottonwood Heights, UT
	"840":{40.2462, -112.0843}, // Fairfield, UT
	"840":{40.3058, -111.7544}, // Vineyard, UT
	"840":{40.6024, -112.0008}, // West Jordan, UT
	"840":{40.8439, -111.9187}, // North Salt Lake, UT
	"840":{40.3543, -110.7095}, // Tabiona, UT
	"840":{40.6430, -111.4007}, // Hideout, UT
	"840":{40.2925, -110.0094}, // Roosevelt, UT
	"840":{41.0415, -111.6801}, // Morgan, UT
	"840":{40.2307, -112.7541}, // Dugway, UT
	"840":{41.1030, -112.0237}, // Clearfield, UT
	"840":{40.3137, -112.0114}, // Eagle Mountain, UT
	"840":{40.3576, -110.2232}, // Bluebell, UT
	"840":{40.6356, -112.3054}, // Stansbury Park, UT
	"840":{40.7259, -111.2770}, // Oakley, UT
	"840":{40.4629, -111.7724}, // Alpine, UT
	"840":{40.6148, -112.4777}, // Grantsville, UT
	"840":{40.7241, -114.0250}, // Wendover, UT
	"840":{40.8754, -111.3846}, // Hoytsville, UT
	"840":{41.6643, -111.1850}, // Randolph, UT
	"840":{40.5706, -111.8510}, // Sandy, UT
	"840":{41.1220, -112.0995}, // West Point, UT
	"840":{40.4668, -111.4097}, // Daniel, UT
	"840":{40.4108, -111.2958}, // Independence, UT
	"840":{40.7257, -111.3380}, // Peoa, UT
	"840":{40.2267, -109.8308}, // Randlett, UT
	"840":{40.5412, -111.4751}, // Interlaken, UT
	"840":{40.4662, -111.4599}, // Charleston, UT
	"840":{40.6148, -111.8928}, // Midvale, UT
	"840":{40.6505, -111.5020}, // Park City, UT
	"840":{40.4956, -111.8607}, // Draper, UT
	"840":{40.4137, -111.8726}, // Lehi, UT
	"840":{40.2949, -109.9493}, // Ballard, UT
	"840":{40.3613, -112.4506}, // Rush Valley, UT
	"840":{41.0771, -111.9622}, // Layton, UT
	"840":{40.0305, -109.1884}, // Bonanza, UT
	"840":{41.8276, -111.3243}, // Laketown, UT
	"840":{41.1392, -112.0285}, // Sunset, UT
	"840":{40.3414, -111.7188}, // Lindon, UT
	"840":{41.0277, -111.9081}, // Fruit Heights, UT
	"840":{40.8890, -109.4751}, // Flaming Gorge, UT
	"840":{40.7434, -111.4905}, // Silver Summit, UT
	"840":{41.0291, -111.9455}, // Kaysville, UT
	"840":{41.0138, -111.4933}, // Henefer, UT
	"840":{40.0956, -112.4458}, // Vernon, UT
	"840":{41.8866, -111.4327}, // Garden, UT
	"840":{40.4275, -111.7955}, // Highland, UT
	"840":{40.7094, -112.0828}, // Magna, UT
	"840":{40.5581, -112.0932}, // Copperton, UT
	"840":{40.5068, -111.3984}, // Heber, UT
	"840":{40.4318, -109.4913}, // Naples, UT
	"840":{41.5212, -111.1639}, // Woodruff, UT
	"840":{41.1475, -111.7901}, // Mountain Green, UT
	"840":{40.5571, -111.9783}, // South Jordan, UT
	"840":{41.9370, -111.4122}, // Garden City, UT
	"840":{40.4517, -109.5379}, // Vernal, UT
	"840":{40.8731, -111.9170}, // Woods Cross, UT
	"840":{40.5818, -111.6229}, // Alta, UT
	"840":{40.9284, -111.8849}, // Centerville, UT
	"840":{40.3714, -111.7411}, // Pleasant Grove, UT
	"840":{40.5838, -111.2356}, // Woodland, UT
	"840":{40.9922, -109.7210}, // Manila, UT
	"840":{40.6499, -111.2723}, // Kamas, UT
	"840":{40.9301, -109.4042}, // Dutch John, UT
	"840":{40.3586, -110.2887}, // Altamont, UT
	"840":{40.5393, -112.3082}, // Tooele, UT
	"840":{40.5711, -111.7959}, // Granite, UT
	"840":{40.3793, -111.7951}, // American Fork, UT
	"840":{40.1933, -110.0624}, // Myton, UT
	"840":{40.4135, -111.7530}, // Cedar Hills, UT
	"840":{40.1754, -110.3940}, // Duchesne, UT
	"840":{40.3874, -111.4206}, // Wallsburg, UT
	"840":{40.7432, -111.5814}, // Summit Park, UT
	"840":{40.2811, -109.8770}, // Fort Duchesne, UT
	"840":{40.6105, -111.2744}, // Francis, UT
	"840":{40.2983, -111.6992}, // Orem, UT
	"840":{41.1395, -112.0656}, // Clinton, UT
	"840":{40.8722, -111.8647}, // Bountiful, UT
	"840":{40.7042, -111.5438}, // Snyderville, UT
	"840":{40.6569, -111.9493}, // Taylorsville, UT
	"840":{40.5176, -111.9635}, // Riverton, UT
	"840":{40.3718, -112.2504}, // Ophir, UT
	"840":{40.6028, -112.3214}, // Erda, UT
	"840":{40.9803, -111.4393}, // Echo, UT
	"840":{40.8168, -111.4151}, // Wanship, UT
	"840":{40.4719, -109.5786}, // Maeser, UT
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <841> with city SALT LAKE CTY and state UT
	"841":{40.6137, -111.8144}, // Cottonwood Heights, UT
	"841":{40.7774, -111.9300}, // Salt Lake City, UT
	"841":{40.6599, -111.8226}, // Holladay, UT
	"841":{40.7056, -111.8986}, // South Salt Lake, UT
	"841":{40.6498, -111.8874}, // Murray, UT
	"841":{40.6889, -112.0115}, // West Valley City, UT
	"841":{40.6892, -111.8291}, // Millcreek, UT
	"841":{40.6520, -112.0093}, // Kearns, UT
	"841":{40.6569, -111.9493}, // Taylorsville, UT
	"841":{40.7894, -111.7411}, // Emigration Canyon, UT
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <842> with city SALT LAKE CTY and state UT
	"842":{41.2624, -112.0366}, // Marriott-Slaterville, UT
	"842":{41.2280, -111.9677}, // Ogden, UT
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <843> with city SALT LAKE CTY and state UT
	"843":{41.9116, -111.9356}, // Trenton, UT
	"843":{41.7033, -111.8122}, // Providence, UT
	"843":{41.7101, -111.9778}, // Mendon, UT
	"843":{41.5678, -111.8330}, // Paradise, UT
	"843":{41.6124, -112.1250}, // Bear River City, UT
	"843":{41.9230, -111.8076}, // Richmond, UT
	"843":{41.8308, -112.0052}, // Cache, UT
	"843":{41.9204, -112.0503}, // Clarkston, UT
	"843":{41.6815, -112.3201}, // Thatcher, UT
	"843":{41.7219, -111.8194}, // River Heights, UT
	"843":{41.6725, -111.8454}, // Nibley, UT
	"843":{41.3022, -111.8086}, // Eden, UT
	"843":{41.6855, -111.8213}, // Millville, UT
	"843":{41.9756, -112.2405}, // Portage, UT
	"843":{41.6930, -112.0883}, // Deweyville, UT
	"843":{41.5497, -112.1227}, // Corinne, UT
	"843":{41.7400, -111.8419}, // Logan, UT
	"843":{41.8566, -111.8974}, // Amalga, UT
	"843":{41.6759, -112.1376}, // Elwood, UT
	"843":{41.4648, -112.0401}, // Perry, UT
	"843":{41.8347, -111.8266}, // Smithfield, UT
	"843":{41.9615, -111.8797}, // Lewiston, UT
	"843":{41.6223, -111.9430}, // Wellsville, UT
	"843":{41.7759, -111.8066}, // North Logan, UT
	"843":{41.5006, -111.9344}, // Mantua, UT
	"843":{41.3426, -111.8649}, // Liberty, UT
	"843":{41.1598, -112.2870}, // Hooper, UT
	"843":{41.5035, -112.0454}, // Brigham City, UT
	"843":{41.4146, -112.0446}, // Willard, UT
	"843":{41.8009, -111.8118}, // Hyde Park, UT
	"843":{41.5354, -111.8125}, // Avon, UT
	"843":{41.7489, -111.9169}, // Benson, UT
	"843":{41.2602, -111.7741}, // Huntsville, UT
	"843":{41.6321, -111.8428}, // Hyrum, UT
	"843":{41.8126, -112.1174}, // Fielding, UT
	"843":{41.8612, -111.9906}, // Newton, UT
	"843":{41.9728, -112.7164}, // Snowville, UT
	"843":{41.8106, -112.1403}, // Riverside, UT
	"843":{41.9717, -111.9565}, // Cornish, UT
	"843":{41.7733, -112.4450}, // Howell, UT
	"843":{41.7716, -111.9859}, // Peter, UT
	"843":{41.8734, -112.1454}, // Plymouth, UT
	"843":{41.3253, -111.8288}, // Wolf Creek, UT
	"843":{41.7362, -112.1627}, // Garland, UT
	"843":{41.3583, -112.0408}, // South Willard, UT
	"843":{41.9683, -111.7793}, // Cove, UT
	"843":{41.6361, -112.0857}, // Honeyville, UT
	"843":{41.7188, -112.1891}, // Tremonton, UT
	// END OF OPTIONS TO PICK
	"844":{41.2280, -111.9677}, // Ogden, UT
	"845":{40.2457, -111.6457}, // Provo, UT
	"846":{40.2457, -111.6457}, // Provo, UT
	"847":{40.2457, -111.6457}, // Provo, UT
	"850":{33.5722, -112.0891}, // Phoenix, AZ
	"852":{33.5722, -112.0891}, // Phoenix, AZ
	"853":{33.5722, -112.0891}, // Phoenix, AZ
	"855":{33.3869, -110.7514}, // Globe, AZ
	"856":{32.1545, -110.8782}, // Tucson, AZ
	"857":{32.1545, -110.8782}, // Tucson, AZ
	"859":{34.2671, -110.0384}, // Show Low, AZ
	"860":{35.1872, -111.6195}, // Flagstaff, AZ
	"863":{34.5850, -112.4475}, // Prescott, AZ
	"864":{35.2170, -114.0105}, // Kingman, AZ
	"865":{35.5183, -108.7423}, // Gallup, NM
	"870":{35.1053, -106.6464}, // Albuquerque, NM
	"871":{35.1053, -106.6464}, // Albuquerque, NM
	"872":{35.1053, -106.6464}, // Albuquerque, NM
	"873":{35.5183, -108.7423}, // Gallup, NM
	"874":{36.7555, -108.1823}, // Farmington, NM
	"875":{35.1053, -106.6464}, // Albuquerque, NM
	"877":{35.6011, -105.2206}, // Las Vegas, NM
	"878":{34.0543, -106.9066}, // Socorro, NM
	TODO PLACEHOLDER TO FIX ZIP <879> with city TRUTH OR CONS and state NM
	"879":{33.1864, -107.2589}, // Truth or Consequences, NM
	"879":{33.3461, -107.6488}, // Winston, NM
	"879":{33.2047, -107.2100}, // Hot Springs Landing, NM
	"879":{33.1159, -107.2951}, // Williamsburg, NM
	"879":{33.0589, -107.2982}, // Las Palomas, NM
	"879":{32.9277, -107.3164}, // Oasis, NM
	"879":{33.1806, -107.2269}, // Elephant Butte, NM
	"879":{32.7568, -107.2665}, // Garfield, NM
	"879":{32.6682, -107.1641}, // Hatch, NM
	"879":{32.8472, -107.3204}, // Arrey, NM
	"879":{32.6538, -107.1358}, // Rodey, NM
	"879":{32.9806, -107.3075}, // Caballo, NM
	"879":{32.6729, -107.0764}, // Rincon, NM
	"879":{32.7122, -107.2044}, // Salem, NM
	// END OF OPTIONS TO PICK
	"880":{32.3265, -106.7893}, // Las Cruces, NM
	"881":{34.4376, -103.1923}, // Clovis, NM
	"882":{33.3730, -104.5294}, // Roswell, NM
	"883":{32.8837, -105.9624}, // Alamogordo, NM
	"884":{35.1701, -103.7042}, // Tucumcari, NM
	"885":{31.8479, -106.4309}, // El Paso, TX
	"889":{36.2333, -115.2654}, // Las Vegas, NV
	"890":{36.2333, -115.2654}, // Las Vegas, NV
	"891":{36.2333, -115.2654}, // Las Vegas, NV
	"893":{39.2649, -114.8709}, // Ely, NV
	"894":{39.5497, -119.8483}, // Reno, NV
	"895":{39.5497, -119.8483}, // Reno, NV
	"897":{39.1511, -119.7474}, // Carson City, NV
	"898":{40.8387, -115.7674}, // Elko, NV
	"900":{34.1139, -118.4068}, // Los Angeles, CA
	"901":{34.1139, -118.4068}, // Los Angeles, CA
	"902":{33.9566, -118.3444}, // Inglewood, CA
	"903":{33.9566, -118.3444}, // Inglewood, CA
	"904":{34.0232, -118.4813}, // Santa Monica, CA
	"905":{33.8346, -118.3417}, // Torrance, CA
	"906":{33.7980, -118.1675}, // Long Beach, CA
	"907":{33.7980, -118.1675}, // Long Beach, CA
	"908":{33.7980, -118.1675}, // Long Beach, CA
	"910":{34.1597, -118.1390}, // Pasadena, CA
	"911":{34.1597, -118.1390}, // Pasadena, CA
	"912":{34.1818, -118.2468}, // Glendale, CA
	TODO PLACEHOLDER TO FIX ZIP <913> with city VAN NUYS and state CA
	"913":{34.1510, -118.7608}, // Agoura Hills, CA
	"913":{34.1847, -118.9445}, // Casa Conejo, CA
	"913":{34.4818, -118.6317}, // Castaic, CA
	"913":{34.5044, -118.3160}, // Agua Dulce, CA
	"913":{34.1375, -118.6689}, // Calabasas, CA
	"913":{34.1369, -118.8221}, // Westlake Village, CA
	"913":{34.4817, -118.6665}, // Hasley Canyon, CA
	"913":{34.3894, -118.5885}, // Stevenson Ranch, CA
	"913":{34.4504, -118.6717}, // Val Verde, CA
	"913":{34.2886, -118.4363}, // San Fernando, CA
	"913":{34.2081, -118.6876}, // Bell Canyon, CA
	"913":{34.4155, -118.4992}, // Santa Clarita, CA
	"913":{34.1139, -118.4068}, // Los Angeles, CA
	"913":{34.1914, -118.8755}, // Thousand Oaks, CA
	"913":{34.1637, -118.6612}, // Hidden Hills, CA
	"913":{34.1339, -118.8804}, // Lake Sherwood, CA
	"913":{34.1849, -118.7669}, // Oak Park, CA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <914> with city VAN NUYS and state CA
	"914":{34.1139, -118.4068}, // Los Angeles, CA
	// END OF OPTIONS TO PICK
	"915":{34.1879, -118.3235}, // Burbank, CA
	TODO PLACEHOLDER TO FIX ZIP <916> with city NORTH HOLLYWOOD and state CA
	"916":{34.1139, -118.4068}, // Los Angeles, CA
	// END OF OPTIONS TO PICK
	"917":{34.0175, -117.9268}, // Industry, CA
	"918":{34.0175, -117.9268}, // Industry, CA
	"919":{32.8312, -117.1225}, // San Diego, CA
	"920":{32.8312, -117.1225}, // San Diego, CA
	"921":{32.8312, -117.1225}, // San Diego, CA
	TODO PLACEHOLDER TO FIX ZIP <922> with city SN BERNARDINO and state CA
	"922":{33.3548, -115.7306}, // Bombay Beach, CA
	"922":{33.9797, -116.9694}, // Cherry Valley, CA
	"922":{33.9356, -116.6873}, // Whitewater, CA
	"922":{33.7346, -116.2346}, // Indio, CA
	"922":{33.9874, -117.0542}, // Calimesa, CA
	"922":{33.8150, -116.3545}, // Thousand Palms, CA
	"922":{32.7431, -116.0020}, // Ocotillo, CA
	"922":{33.9461, -116.8991}, // Banning, CA
	"922":{33.6536, -116.2785}, // La Quinta, CA
	"922":{34.0724, -116.5627}, // Morongo Valley, CA
	"922":{33.7435, -116.2874}, // Bermuda Dunes, CA
	"922":{32.6849, -115.4944}, // Calexico, CA
	"922":{33.9223, -116.4401}, // Desert Edge, CA
	"922":{33.5957, -114.7317}, // Mesa Verde, CA
	"922":{33.6905, -116.1430}, // Coachella, CA
	"922":{32.7867, -115.5586}, // El Centro, CA
	"922":{32.8129, -115.3779}, // Holtville, CA
	"922":{33.7036, -116.3396}, // Indian Wells, CA
	"922":{34.1234, -116.4216}, // Yucca Valley, CA
	"922":{33.8017, -116.5382}, // Palm Springs, CA
	"922":{33.7790, -116.2930}, // Desert Palms, CA
	"922":{33.8912, -116.3551}, // Sky Valley, CA
	"922":{33.3759, -116.0113}, // Salton Sea Beach, CA
	"922":{33.7634, -116.4271}, // Rancho Mirage, CA
	"922":{33.5157, -115.9092}, // North Shore, CA
	"922":{34.2730, -116.4145}, // Homestead Valley, CA
	"922":{33.9070, -116.9762}, // Beaumont, CA
	"922":{33.5275, -116.1261}, // Oasis, CA
	"922":{33.2387, -115.5146}, // Niland, CA
	"922":{33.5238, -114.6530}, // Ripley, CA
	"922":{33.8363, -116.4642}, // Cathedral City, CA
	"922":{33.9179, -116.4796}, // Garnet, CA
	"922":{33.8411, -116.2484}, // Indio Hills, CA
	"922":{33.5767, -116.0645}, // Mecca, CA
	"922":{33.2994, -115.9609}, // Salton City, CA
	"922":{33.6262, -116.1309}, // Thermal, CA
	"922":{32.7318, -115.5204}, // Heber, CA
	"922":{33.9551, -116.5430}, // Desert Hot Springs, CA
	"922":{33.4041, -116.0394}, // Desert Shores, CA
	"922":{34.1478, -116.0659}, // Twentynine Palms, CA
	"922":{32.7899, -115.6842}, // Seeley, CA
	"922":{33.4279, -114.7273}, // Palo Verde, CA
	"922":{33.6227, -116.2126}, // Vista Santa Rosa, CA
	"922":{34.1760, -114.2640}, // Bluewater, CA
	"922":{33.9127, -116.7828}, // Cabazon, CA
	"922":{34.1393, -114.3604}, // Big River, CA
	"922":{32.7372, -114.6378}, // Winterhaven, CA
	"922":{33.1493, -115.5056}, // Calipatria, CA
	"922":{33.0389, -115.6223}, // Westmorland, CA
	"922":{33.7378, -115.3666}, // Desert Center, CA
	"922":{34.1236, -116.3128}, // Joshua Tree, CA
	"922":{33.6219, -114.6195}, // Blythe, CA
	"922":{32.8388, -115.5723}, // Imperial, CA
	"922":{32.9783, -115.5288}, // Brawley, CA
	"922":{33.7378, -116.3695}, // Palm Desert, CA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <923> with city SN BERNARDINO and state CA
	"923":{33.9874, -117.0542}, // Calimesa, CA
	"923":{36.4277, -116.8747}, // Furnace Creek, CA
	"923":{34.4427, -116.9021}, // Lucerne Valley, CA
	"923":{34.4988, -117.2684}, // Spring Valley Lake, CA
	"923":{34.2536, -116.7903}, // Big Bear City, CA
	"923":{34.8861, -117.1078}, // Lenwood, CA
	"923":{34.5277, -117.3536}, // Victorville, CA
	"923":{34.2104, -117.1147}, // Running Springs, CA
	"923":{35.9510, -116.3056}, // Shoshone, CA
	"923":{34.0601, -117.4013}, // Bloomington, CA
	"923":{34.8164, -114.6189}, // Needles, CA
	"923":{34.1175, -117.3894}, // Rialto, CA
	"923":{34.4398, -117.5248}, // Phelan, CA
	"923":{34.8661, -117.0472}, // Barstow, CA
	"923":{35.8350, -116.2074}, // Tecopa, CA
	"923":{34.2486, -117.2890}, // Crestline, CA
	"923":{34.3975, -117.3147}, // Hesperia, CA
	"923":{34.2429, -116.8950}, // Big Bear Lake, CA
	"923":{34.0336, -117.0430}, // Yucaipa, CA
	"923":{34.0106, -117.3098}, // Highgrove, CA
	"923":{34.3912, -117.4125}, // Oak Hills, CA
	"923":{34.3495, -117.6299}, // Wrightwood, CA
	"923":{34.7519, -117.3431}, // Silver Lakes, CA
	"923":{34.0312, -117.3132}, // Grand Terrace, CA
	"923":{34.0968, -117.4599}, // Fontana, CA
	"923":{34.1417, -117.2945}, // San Bernardino, CA
	"923":{34.0538, -117.3254}, // Colton, CA
	"923":{34.0448, -116.9497}, // Oak Glen, CA
	"923":{34.0511, -117.1712}, // Redlands, CA
	"923":{34.0451, -117.2498}, // Loma Linda, CA
	"923":{34.2531, -117.1945}, // Lake Arrowhead, CA
	"923":{34.4976, -117.3472}, // Mountain View Acres, CA
	"923":{35.2476, -116.6834}, // Fort Irwin, CA
	"923":{34.0609, -117.1109}, // Mentone, CA
	"923":{34.1113, -117.1654}, // Highland, CA
	"923":{34.4438, -117.6214}, // Pinon Hills, CA
	"923":{34.2499, -117.5044}, // Lytle Creek, CA
	"923":{34.5815, -117.4397}, // Adelanto, CA
	"923":{34.5328, -117.2104}, // Apple Valley, CA
	"923":{35.2769, -116.0718}, // Baker, CA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <924> with city SN BERNARDINO and state CA
	"924":{34.1548, -117.3482}, // Muscoy, CA
	"924":{34.0968, -117.4599}, // Fontana, CA
	"924":{34.1417, -117.2945}, // San Bernardino, CA
	"924":{34.1113, -117.1654}, // Highland, CA
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <925> with city SN BERNARDINO and state CA
	"925":{33.8724, -117.4624}, // El Sobrante, CA
	"925":{33.7254, -117.2851}, // Meadowbrook, CA
	"925":{33.7970, -116.9915}, // San Jacinto, CA
	"925":{33.7443, -116.7257}, // Idyllwild-Pine Cove, CA
	"925":{33.8888, -117.2777}, // March ARB, CA
	"925":{33.5679, -116.6967}, // Anza, CA
	"925":{33.7648, -117.1572}, // Romoland, CA
	"925":{33.7067, -117.3344}, // Warm Springs, CA
	"925":{33.7436, -116.8872}, // Valle Vista, CA
	"925":{33.4928, -117.1315}, // Temecula, CA
	"925":{33.7112, -116.7248}, // Mountain Center, CA
	"925":{33.5173, -116.8107}, // Lake Riverside, CA
	"925":{33.8285, -117.1233}, // Lakeview, CA
	"925":{33.8250, -117.3683}, // Lake Mathews, CA
	"925":{34.0106, -117.3098}, // Highgrove, CA
	"925":{33.7706, -117.2772}, // Good Hope, CA
	"925":{33.7899, -117.2233}, // Perris, CA
	"925":{33.5998, -117.1069}, // French Valley, CA
	"925":{33.7301, -116.9410}, // East Hemet, CA
	"925":{33.6173, -117.2583}, // Wildomar, CA
	"925":{34.0010, -117.4706}, // Jurupa Valley, CA
	"925":{33.7459, -117.1132}, // Homeland, CA
	"925":{33.9244, -117.2045}, // Moreno Valley, CA
	"925":{33.6480, -117.3706}, // Lakeland Village, CA
	"925":{33.8011, -117.1415}, // Nuevo, CA
	"925":{33.5720, -117.1909}, // Murrieta, CA
	"925":{33.6909, -117.1849}, // Menifee, CA
	"925":{33.8789, -117.3686}, // Woodcrest, CA
	"925":{33.7341, -116.9969}, // Hemet, CA
	"925":{33.6846, -117.3344}, // Lake Elsinore, CA
	"925":{33.8333, -117.2852}, // Mead Valley, CA
	"925":{33.8784, -117.5116}, // Home Gardens, CA
	"925":{33.7146, -117.0775}, // Winchester, CA
	"925":{33.9381, -117.3948}, // Riverside, CA
	"925":{33.4522, -116.8555}, // Aguanga, CA
	"925":{33.7350, -117.0783}, // Green Acres, CA
	"925":{33.6884, -117.2621}, // Canyon Lake, CA
	// END OF OPTIONS TO PICK
	"926":{33.7366, -117.8819}, // Santa Ana, CA
	"927":{33.7366, -117.8819}, // Santa Ana, CA
	"928":{33.8390, -117.8573}, // Anaheim, CA
	"930":{34.1962, -119.1819}, // Oxnard, CA
	"931":{34.4285, -119.7202}, // Santa Barbara, CA
	"932":{35.3530, -119.0359}, // Bakersfield, CA
	"933":{35.3530, -119.0359}, // Bakersfield, CA
	"934":{34.4285, -119.7202}, // Santa Barbara, CA
	"935":{35.0139, -118.1895}, // Mojave, CA
	"936":{36.7831, -119.7941}, // Fresno, CA
	"937":{36.7831, -119.7941}, // Fresno, CA
	"938":{36.7831, -119.7941}, // Fresno, CA
	"939":{36.6884, -121.6317}, // Salinas, CA
	"940":{37.7562, -122.4430}, // San Francisco, CA
	"941":{37.7562, -122.4430}, // San Francisco, CA
	"942":{38.5667, -121.4683}, // Sacramento, CA
	"943":{37.3913, -122.1467}, // Palo Alto, CA
	"944":{37.5522, -122.3122}, // San Mateo, CA
	"945":{37.7903, -122.2165}, // Oakland, CA
	"946":{37.7903, -122.2165}, // Oakland, CA
	"947":{37.8723, -122.2760}, // Berkeley, CA
	"948":{37.9477, -122.3390}, // Richmond, CA
	TODO PLACEHOLDER TO FIX ZIP <949> with city NORTH BAY and state CA
	"949":{38.2436, -122.9560}, // Dillon Beach, CA
	"949":{37.9884, -122.5950}, // Fairfax, CA
	"949":{38.1120, -122.5190}, // Black Point-Green Point, CA
	"949":{38.0052, -122.6384}, // Woodacre, CA
	"949":{38.3743, -123.0720}, // Carmet, CA
	"949":{38.3246, -122.9148}, // Valley Ford, CA
	"949":{38.0180, -122.6914}, // Lagunitas-Forest Knolls, CA
	"949":{37.8793, -122.5382}, // Tamalpais-Homestead Valley, CA
	"949":{38.2470, -122.9054}, // Tomales, CA
	"949":{38.0067, -122.6634}, // San Geronimo, CA
	"949":{38.3480, -122.6964}, // Rohnert Park, CA
	"949":{38.0821, -122.8471}, // Inverness, CA
	"949":{38.3836, -123.0732}, // Sereno del Mar, CA
	"949":{38.3005, -122.6707}, // Penngrove, CA
	"949":{37.9821, -122.5699}, // San Anselmo, CA
	"949":{37.9056, -122.5187}, // Alto, CA
	"949":{37.8580, -122.4932}, // Sausalito, CA
	"949":{37.9904, -122.5222}, // San Rafael, CA
	"949":{38.3250, -123.0308}, // Bodega Bay, CA
	"949":{38.3488, -122.9712}, // Bodega, CA
	"949":{37.9051, -122.6457}, // Stinson Beach, CA
	"949":{37.8854, -122.4637}, // Tiburon, CA
	"949":{37.9238, -122.5129}, // Corte Madera, CA
	"949":{37.9638, -122.5615}, // Ross, CA
	"949":{38.0405, -122.5765}, // Lucas Valley-Marinwood, CA
	"949":{38.0606, -122.7008}, // Nicasio, CA
	"949":{38.0120, -122.5877}, // Sleepy Hollow, CA
	"949":{37.8925, -122.5078}, // Strawberry, CA
	"949":{38.0847, -122.8093}, // Point Reyes Station, CA
	"949":{38.2423, -122.6267}, // Petaluma, CA
	"949":{38.0920, -122.5576}, // Novato, CA
	"949":{37.8711, -122.5137}, // Marin City, CA
	"949":{37.9481, -122.5498}, // Kentfield, CA
	"949":{37.9394, -122.5313}, // Larkspur, CA
	"949":{38.3279, -122.7092}, // Cotati, CA
	"949":{38.0055, -122.5033}, // Santa Venetia, CA
	"949":{38.3463, -123.0595}, // Salmon Creek, CA
	"949":{37.8654, -122.5858}, // Muir Beach, CA
	"949":{37.9177, -122.7095}, // Bolinas, CA
	"949":{37.9085, -122.5420}, // Mill Valley, CA
	"949":{37.8735, -122.4662}, // Belvedere, CA
	"949":{38.3183, -122.8502}, // Bloomfield, CA
	// END OF OPTIONS TO PICK
	"950":{37.3021, -121.8489}, // San Jose, CA
	"951":{37.3021, -121.8489}, // San Jose, CA
	"952":{37.9766, -121.3111}, // Stockton, CA
	"953":{37.9766, -121.3111}, // Stockton, CA
	TODO PLACEHOLDER TO FIX ZIP <954> with city NORTH BAY and state CA
	"954":{39.4429, -123.3963}, // Brooktrails, CA
	"954":{38.7173, -122.9034}, // Geyserville, CA
	"954":{38.8356, -122.7230}, // Cobb, CA
	"954":{38.4376, -122.8660}, // Graton, CA
	"954":{38.8003, -122.5505}, // Hidden Valley Lake, CA
	"954":{38.4825, -122.8899}, // Forestville, CA
	"954":{38.7962, -123.0154}, // Cloverdale, CA
	"954":{38.4178, -122.7298}, // Roseland, CA
	"954":{39.4400, -123.8013}, // Fort Bragg, CA
	"954":{38.3339, -122.5065}, // Eldridge, CA
	"954":{38.7519, -122.6221}, // Middletown, CA
	"954":{38.9512, -122.7207}, // Clearlake Riviera, CA
	"954":{39.1266, -122.8525}, // Nice, CA
	"954":{39.3190, -123.1123}, // Potter Valley, CA
	"954":{38.9118, -122.6084}, // Lower Lake, CA
	"954":{38.4004, -122.9349}, // Occidental, CA
	"954":{39.8025, -123.2499}, // Covelo, CA
	"954":{38.4937, -122.7734}, // Fulton, CA
	"954":{39.3110, -123.7908}, // Mendocino, CA
	"954":{39.4927, -123.7812}, // Cleone, CA
	"954":{38.4000, -122.8277}, // Sebastopol, CA
	"954":{39.4047, -123.3494}, // Willits, CA
	"954":{38.4153, -122.5387}, // Kenwood, CA
	"954":{38.3126, -122.4888}, // Boyes Hot Springs, CA
	"954":{39.0024, -122.7794}, // Soda Bay, CA
	"954":{39.2703, -123.7817}, // Little River, CA
	"954":{38.2577, -122.4982}, // Temelec, CA
	"954":{38.5265, -123.0978}, // Cazadero, CA
	"954":{38.6229, -122.8651}, // Healdsburg, CA
	"954":{38.9123, -123.6954}, // Point Arena, CA
	"954":{39.2322, -123.1970}, // Calpella, CA
	"954":{38.5418, -122.8086}, // Windsor, CA
	"954":{39.0697, -122.7760}, // Lucerne, CA
	"954":{39.6715, -123.4945}, // Laytonville, CA
	"954":{38.3274, -122.4871}, // Fetters Hot Springs-Agua Caliente, CA
	"954":{39.3631, -123.8044}, // Caspar, CA
	"954":{39.1463, -123.2105}, // Ukiah, CA
	"954":{38.9691, -123.1173}, // Hopland, CA
	"954":{38.8126, -123.5704}, // Anchor Bay, CA
	"954":{39.0115, -123.3740}, // Boonville, CA
	"954":{39.2256, -123.7564}, // Albion, CA
	"954":{38.2902, -122.4598}, // Sonoma, CA
	"954":{39.0883, -122.9054}, // North Lakeport, CA
	"954":{38.2975, -122.4915}, // El Verano, CA
	"954":{39.0657, -123.4452}, // Philo, CA
	"954":{38.5137, -122.9894}, // Guerneville, CA
	"954":{38.5130, -122.7536}, // Larkfield-Wikiup, CA
	"954":{38.4683, -123.0147}, // Monte Rio, CA
	"954":{39.0392, -122.9218}, // Lakeport, CA
	"954":{39.1653, -122.9052}, // Upper Lake, CA
	"954":{39.1314, -123.1648}, // Talmage, CA
	"954":{38.9590, -122.6331}, // Clearlake, CA
	"954":{38.9743, -123.6910}, // Manchester, CA
	"954":{38.4458, -122.7067}, // Santa Rosa, CA
	"954":{38.7166, -123.4528}, // Sea Ranch, CA
	"954":{39.2651, -123.5897}, // Comptche, CA
	"954":{38.4511, -123.1204}, // Jenner, CA
	"954":{38.9704, -122.8327}, // Kelseyville, CA
	"954":{38.5410, -123.2590}, // Timber Cove, CA
	"954":{38.3564, -122.5412}, // Glen Ellen, CA
	"954":{39.2690, -123.2023}, // Redwood Valley, CA
	"954":{39.0218, -122.6593}, // Clearlake Oaks, CA
	// END OF OPTIONS TO PICK
	"955":{40.7941, -124.1568}, // Eureka, CA
	"956":{38.5667, -121.4683}, // Sacramento, CA
	"957":{38.5667, -121.4683}, // Sacramento, CA
	"958":{38.5667, -121.4683}, // Sacramento, CA
	"959":{39.1518, -121.5836}, // Marysville, CA
	"960":{40.5698, -122.3650}, // Redding, CA
	"961":{39.5497, -119.8483}, // Reno, NV
	TODO PLACEHOLDER TO FIX ZIP <962> with city APO/FPO and state AP
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <963> with city APO/FPO and state AP
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <964> with city APO/FPO and state AP
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <965> with city APO/FPO and state AP
	// END OF OPTIONS TO PICK
	TODO PLACEHOLDER TO FIX ZIP <966> with city FPO and state AP
	// END OF OPTIONS TO PICK
	"967":{21.3294, -157.8460}, // Honolulu, HI
	"968":{21.3294, -157.8460}, // Honolulu, HI
	TODO PLACEHOLDER TO FIX ZIP <969> with city BARRIGADA and state GU
	// END OF OPTIONS TO PICK
	"970":{45.5371, -122.6500}, // Portland, OR
	"971":{45.5371, -122.6500}, // Portland, OR
	"972":{45.5371, -122.6500}, // Portland, OR
	"973":{44.9232, -123.0245}, // Salem, OR
	"974":{44.0563, -123.1173}, // Eugene, OR
	"975":{42.3372, -122.8537}, // Medford, OR
	"976":{42.2191, -121.7754}, // Klamath Falls, OR
	"977":{44.0562, -121.3087}, // Bend, OR
	"978":{45.6755, -118.8209}, // Pendleton, OR
	"979":{43.6007, -116.2312}, // Boise, ID
	"980":{47.6211, -122.3244}, // Seattle, WA
	"981":{47.6211, -122.3244}, // Seattle, WA
	"982":{47.9524, -122.1670}, // Everett, WA
	"983":{47.2431, -122.4531}, // Tacoma, WA
	"984":{47.2431, -122.4531}, // Tacoma, WA
	"985":{47.0417, -122.8959}, // Olympia, WA
	"986":{45.5371, -122.6500}, // Portland, OR
	"988":{47.4338, -120.3286}, // Wenatchee, WA
	"989":{46.5923, -120.5496}, // Yakima, WA
	"990":{47.6671, -117.4330}, // Spokane, WA
	"991":{47.6671, -117.4330}, // Spokane, WA
	"992":{47.6671, -117.4330}, // Spokane, WA
	"993":{46.2506, -119.1303}, // Pasco, WA
	"994":{46.3934, -116.9934}, // Lewiston, ID
	"995":{61.1508, -149.1091}, // Anchorage, AK
	"996":{61.1508, -149.1091}, // Anchorage, AK
	"997":{64.8353, -147.6534}, // Fairbanks, AK
	"998":{58.4546, -134.1739}, // Juneau, AK
	"999":{55.3556, -131.6698}, // Ketchikan, AK
}
