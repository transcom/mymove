# FOR THE LOVE OF ALL THAT IS HOLY, THIS THING IS JUST:

# DELETES
# MERGES
# RENAMES
# AND NEWS

# delete related entries need to be deleted
# merge reentries need to be reassigned, then they can be deleted
# renames are easy
# news are easy

import os, sys, uuid, pandas as pd

current_dir = os.getcwd()

filename = (
    f"{current_dir}/migrations/app/schema/20230810180036_update_duty_locations.up.sql"
)

f = open(filename, "w")
f2 = open("./all-names.txt", "w")

all_names = []

if len(sys.argv) < 2:
    sys.exit("Input file required.")

deletes = [
    "Adak, AK 99546",
    "Aiea, HI 96701",
    "Akiachak, AK 99551",
    "Akiak, AK 99552",
    "Akutan, AK 99553",
    "Alakanuk, AK 99554",
    "Aleknagik, AK 99555",
    "Allakaket, AK 99720",
    "Ambler, AK 99786",
    "Anahola, HI 96703",
    "Anaktuvuk Pass, AK 99721",
    "Anchorage, AK 99501",
    "Anchorage, AK 99502",
    "Anchorage, AK 99503",
    "Anchorage, AK 99504",
    "Anchorage, AK 99507",
    "Anchorage, AK 99508",
    "Anchorage, AK 99509",
    "Anchorage, AK 99510",
    "Anchorage, AK 99511",
    "Anchorage, AK 99513",
    "Anchorage, AK 99514",
    "Anchorage, AK 99515",
    "Anchorage, AK 99516",
    "Anchorage, AK 99517",
    "Anchorage, AK 99518",
    "Anchorage, AK 99519",
    "Anchorage, AK 99520",
    "Anchorage, AK 99521",
    "Anchorage, AK 99522",
    "Anchorage, AK 99523",
    "Anchorage, AK 99524",
    "Anchorage, AK 99529",
    "Anchorage, AK 99530",
    "Anchorage, AK 99599",
    "Anchorage, AK 99695",
    "Anchor Point, AK 99556",
    "Anderson, AK 99744",
    "Andrews Air Force Base, MD 20762",
    "Angoon, AK 99820",
    "Aniak, AK 99557",
    "Anvik, AK 99558",
    "Arctic Village, AK 99722",
    "Atka, AK 99547",
    "Atqasuk, AK 99791",
    "Auke Bay, AK 99821",
    "Barrow, AK 99723",
    "Beaver, AK 99724",
    "Bethel, AK 99559",
    "Brevig Mission, AK 99785",
    "Big Lake, AK 99652",
    "Bettles Field, AK 99726",
    "Camp H M Smith, HI 96861",
    "Buckland, AK 99727",
    "Chefornak, AK 99561",
    "Cantwell, AK 99729",
    "Captain Cook, HI 96704",
    "Central, AK 99730",
    "Chalkyitsik, AK 99788",
    "Chicken, AK 99732",
    "Chevak, AK 99563",
    "Chignik, AK 99564",
    "Chugiak, AK 99567",
    "Chignik Lagoon, AK 99565",
    "Chignik Lake, AK 99548",
    "Chitina, AK 99566",
    "Circle, AK 99733",
    "Clam Gulch, AK 99568",
    "Cooper Landing, AK 99572",
    "Coffman Cove, AK 99918",
    "Cold Bay, AK 99571",
    "Cordova, AK 99574",
    "Clarks Point, AK 99569",
    "Clear, AK 99704",
    "Copper Center, AK 99573",
    "Craig, AK 99921",
    "Denali National Park, AK 99755",
    "Crooked Creek, AK 99575",
    "Deering, AK 99736",
    "Eek, AK 99578",
    "Douglas, AK 99824",
    "Delta Junction, AK 99737",
    "Dillingham, AK 99576",
    "Dutch Harbor, AK 99692",
    "Eagle, AK 99738",
    "Eagle River, AK 99577",
    "Egegik, AK 99579",
    "Eielson AFB, AK 99702",
    "Elfin Cove, AK 99825",
    "Ekwok, AK 99580",
    "Eleele, HI 96705",
    "Elim, AK 99739",
    "False Pass, AK 99583",
    "Ewa Beach, HI 96706",
    "Fairbanks, AK 99701",
    "Fairbanks, AK 99706",
    "Fairbanks, AK 99707",
    "Fairbanks, AK 99708",
    "Fairbanks, AK 99709",
    "Fairbanks, AK 99710",
    "Fairbanks, AK 99711",
    "Fairbanks, AK 99712",
    "Fairbanks, AK 99775",
    "Fairbanks, AK 99790",
    "Fairchild Air Force Base, WA 99011",
    "Emmonak, AK 99581",
    "Ester, AK 99725",
    "Fort Wainwright, AK 99703",
    "Fe Warren AFB, WY 82005",
    "Fort Benning, GA 31905",
    "Fort Benning, GA 31995",
    "Fort Bragg, NC 28307",
    "Fort Bragg, NC 28310",
    "Fort Myer, VA 22211",
    "Fort Greely, AK 99731",
    "Fort Hood, TX 76544",
    "Gakona, AK 99586",
    "Fort Lee, VA 23801",
    "Fort Yukon, AK 99740",
    "Ft Mitchell, KY 41017",
    "Fort Polk, LA 71459",
    "Fort Rucker, AL 36362",
    "Fort Shafter, HI 96858",
    "Gambell, AK 99742",
    "Galena, AK 99741",
    "Glennallen, AK 99588",
    "Girdwood, AK 99587",
    "Goodnews Bay, AK 99589",
    "Grayling, AK 99590",
    "Hana, HI 96713",
    "Haiku, HI 96708",
    "Haines, AK 99827",
    "Hakalau, HI 96710",
    "Haleiwa, HI 96712",
    "Gustavus, AK 99826",
    "Hanalei, HI 96714",
    "Hanamaulu, HI 96715",
    "Hanapepe, HI 96716",
    "Healy, AK 99743",
    "Hauula, HI 96717",
    "Hawaii National Park, HI 96718",
    "Hawi, HI 96719",
    "Hilo, HI 96720",
    "Hilo, HI 96721",
    "Hoolehua, HI 96729",
    "Holloman Air Force Base, NM 88330",
    "Hoonah, AK 99829",
    "Hooper Bay, AK 99604",
    "Hope, AK 99605",
    "Honaunau, HI 96726",
    "Hughes, AK 99745",
    "Honokaa, HI 96727",
    "Honolulu, HI 96801",
    "Honolulu, HI 96802",
    "Honolulu, HI 96803",
    "Honolulu, HI 96804",
    "Honolulu, HI 96805",
    "Honolulu, HI 96806",
    "Honolulu, HI 96807",
    "Honolulu, HI 96808",
    "Honolulu, HI 96809",
    "Honolulu, HI 96810",
    "Honolulu, HI 96811",
    "Honolulu, HI 96812",
    "Honolulu, HI 96813",
    "Honolulu, HI 96814",
    "Honolulu, HI 96815",
    "Honolulu, HI 96816",
    "Honolulu, HI 96817",
    "Honolulu, HI 96818",
    "Honolulu, HI 96819",
    "Honolulu, HI 96820",
    "Honolulu, HI 96821",
    "Honolulu, HI 96822",
    "Honolulu, HI 96823",
    "Honolulu, HI 96824",
    "Honolulu, HI 96825",
    "Honolulu, HI 96826",
    "Honolulu, HI 96828",
    "Honolulu, HI 96830",
    "Honolulu, HI 96836",
    "Honolulu, HI 96837",
    "Honolulu, HI 96838",
    "Honolulu, HI 96839",
    "Honolulu, HI 96840",
    "Honolulu, HI 96841",
    "Honolulu, HI 96843",
    "Honolulu, HI 96844",
    "Honolulu, HI 96846",
    "Honolulu, HI 96847",
    "Honolulu, HI 96848",
    "Honolulu, HI 96849",
    "Honolulu, HI 96850",
    "Honomu, HI 96728",
    "Holualoa, HI 96725",
    "Holy Cross, AK 99602",
    "Homer, AK 99603",
    "Houston, AK 99694",
    "Huslia, AK 99746",
    "Hyder, AK 99923",
    "Hydaburg, AK 99922",
    "Iliamna, AK 99606",
    "Kalskag, AK 99607",
    "JBER, AK 99505",
    "JBER, AK 99506",
    "Joint Base MDL, NJ 08640",
    "Joint Base MDL, NJ 08641",
    "Indian, AK 99540",
    "JBPHH, HI 96853",
    "JBPHH, HI 96860",
    "JBSA Ft Sam Houston, TX 78234",
    "Juneau, AK 99801",
    "Juneau, AK 99802",
    "Juneau, AK 99803",
    "Juneau, AK 99811",
    "Juneau, AK 99812",
    "Juneau, AK 99850",
    "Kaaawa, HI 96730",
    "Kaltag, AK 99748",
    "Kamuela, HI 96743",
    "Kahuku, HI 96731",
    "Kahului, HI 96732",
    "Kahului, HI 96733",
    "Kapaa, HI 96746",
    "Kapaau, HI 96755",
    "Kaneohe, HI 96744",
    "Kapolei, HI 96707",
    "Kapolei, HI 96709",
    "Kailua, HI 96734",
    "Kailua Kona, HI 96740",
    "Kailua Kona, HI 96745",
    "Kake, AK 99830",
    "Ketchikan, AK 99901",
    "Kaktovik, AK 99747",
    "Kalaheo, HI 96741",
    "Kalaupapa, HI 96742",
    "Kenai, AK 99611",
    "Karluk, AK 99608",
    "Kasigluk, AK 99609",
    "Kasilof, AK 99610",
    "Kaumakani, HI 96747",
    "Kaunakakai, HI 96748",
    "Keaau, HI 96749",
    "Kealakekua, HI 96750",
    "Kealia, HI 96751",
    "Keauhou, HI 96739",
    "Kekaha, HI 96752",
    "Kiana, AK 99749",
    "Ketchikan, AK 99950",
    "Kodiak, AK 99697",
    "Kihei, HI 96753",
    "Kilauea, HI 96754",
    "Kipnuk, AK 99614",
    "King Cove, AK 99612",
    "King Salmon, AK 99613",
    "Koloa, HI 96756",
    "Kongiganak, AK 99545",
    "Kivalina, AK 99750",
    "Kotlik, AK 99620",
    "Klawock, AK 99925",
    "Kotzebue, AK 99752",
    "Koyuk, AK 99753",
    "Koyukuk, AK 99754",
    "Kualapuu, HI 96757",
    "Kula, HI 96790",
    "Kunia, HI 96759",
    "Kobuk, AK 99751",
    "Kurtistown, HI 96760",
    "Kodiak, AK 99615",
    "Kodiak, AK 99619",
    "Kwethluk, AK 99621",
    "Kwigillingok, AK 99622",
    "Lanai City, HI 96763",
    "Lahaina, HI 96761",
    "Lahaina, HI 96767",
    "Laie, HI 96762",
    "Lake Minchumina, AK 99757",
    "Laupahoehoe, HI 96764",
    "Lawai, HI 96765",
    "Larsen Bay, AK 99624",
    "Levelock, AK 99625",
    "Little Rock Air Force Base, AR 72099",
    "Lihue, HI 96766",
    "Manley Hot Springs, AK 99756",
    "Manokotak, AK 99628",
    "Lower Kalskag, AK 99626",
    "Luke Air Force Base, AZ 85309",
    "Makawao, HI 96768",
    "Makaweli, HI 96769",
    "Mcbh Kaneohe Bay, HI 96863",
    "March Air Reserve Base, CA 92518",
    "Marshall, AK 99585",
    "Maunaloa, HI 96770",
    "Mc Call Creek, MS 39647",
    "Mc Clure, PA 17841",
    "Mc Connellsburg, PA 17233",
    "Mekoryuk, AK 99630",
    "Mc Grath, AK 99627",
    "Mililani, HI 96789",
    "Metlakatla, AK 99926",
    "Meyers Chuck, AK 99903",
    "Minto, AK 99758",
    "Napakiak, AK 99634",
    "Moose Pass, AK 99631",
    "Naalehu, HI 96772",
    "Naknek, AK 99633",
    "Mountain View, HI 96771",
    "Mountain Village, AK 99632",
    "North Pole, AK 99705",
    "Nenana, AK 99760",
    "New Stuyahok, AK 99636",
    "Nightmute, AK 99690",
    "Nikiski, AK 99635",
    "Nikolai, AK 99691",
    "Nikolski, AK 99638",
    "Naval Air Station Jrb, TX 76127",
    "Ninilchik, AK 99639",
    "Ninole, HI 96773",
    "Noatak, AK 99761",
    "Nome, AK 99762",
    "Nondalton, AK 99640",
    "Noorvik, AK 99763",
    "Paauilo, HI 96776",
    "Northway, AK 99764",
    "Nuiqsut, AK 99789",
    "Nulato, AK 99765",
    "Nunam Iqua, AK 99666",
    "Nunapitchuk, AK 99641",
    "Ocean View, HI 96737",
    "Old Harbor, AK 99643",
    "Ookala, HI 96774",
    "Ouzinkie, AK 99644",
    "Pahala, HI 96777",
    "Pahoa, HI 96778",
    "Paia, HI 96779",
    "Palmer, AK 99645",
    "Papaaloa, HI 96780",
    "Papaikou, HI 96781",
    "Parcel Return Service, DC 56901",
    "Parcel Return Service, DC 56902",
    "Parcel Return Service, DC 56904",
    "Parcel Return Service, DC 56908",
    "Parcel Return Service, DC 56915",
    "Parcel Return Service, DC 56920",
    "Parcel Return Service, DC 56933",
    "Parcel Return Service, DC 56935",
    "Parcel Return Service, DC 56944",
    "Parcel Return Service, DC 56945",
    "Parcel Return Service, DC 56950",
    "Parcel Return Service, DC 56965",
    "Parcel Return Service, DC 56966",
    "Parcel Return Service, DC 56967",
    "Parcel Return Service, DC 56968",
    "Parcel Return Service, DC 56969",
    "Parcel Return Service, DC 56970",
    "Parcel Return Service, DC 56971",
    "Parcel Return Service, DC 56972",
    "Parcel Return Service, DC 56973",
    "Parcel Return Service, DC 56998",
    "Parcel Return Service, DC 56999",
    "Pearl City, HI 96782",
    "Patrick AFB, FL 32925",
    "Pedro Bay, AK 99647",
    "Pelican, AK 99832",
    "Petersburg, AK 99833",
    "Perryville, AK 99648",
    "Pepeekeo, HI 96783",
    "Point Baker, AK 99927",
    "Platinum, AK 99651",
    "Pilot Point, AK 99649",
    "Pilot Station, AK 99650",
    "Point Hope, AK 99766",
    "Point Lay, AK 99759",
    "Port Heiden, AK 99549",
    "Prudhoe Bay, AK 99734",
    "Port Lions, AK 99550",
    "Puunene, HI 96784",
    "Port Alexander, AK 99836",
    "Port Alsworth, AK 99653",
    "Princeville, HI 96722",
    "Pukalani, HI 96788",
    "Rampart, AK 99767",
    "Quinhagak, AK 99655",
    "Salcha, AK 99714",
    "Red Devil, AK 99656",
    "Rome, NY",
    "Ruby, AK 99768",
    "Saint George Island, AK 99591",
    "Saint Michael, AK 99659",
    "Saint Paul Island, AK 99660",
    "Russian Mission, AK 99657",
    "Saint Marys, AK 99658",
    "Sand Point, AK 99661",
    "Schofield Barracks, HI 96857",
    "Scammon Bay, AK 99662",
    "Scott Air Force Base, IL 62225",
    "Savoonga, AK 99769",
    "Selawik, AK 99770",
    "Shageluk, AK 99665",
    "Shaktoolik, AK 99771",
    "Seldovia, AK 99663",
    "Shungnak, AK 99773",
    "Shishmaref, AK 99772",
    "Seward, AK 99664",
    "Sleetmute, AK 99668",
    "Sitka, AK 99835",
    "Skagway, AK 99840",
    "Stebbins, AK 99671",
    "Skwentna, AK 99667",
    "Soldotna, AK 99669",
    "South Naknek, AK 99670",
    "Sterling, AK 99672",
    "Stevens Village, AK 99774",
    "Tanacross, AK 99776",
    "Tatitlek, AK 99677",
    "Sutton, AK 99674",
    "Takotna, AK 99675",
    "Talkeetna, AK 99676",
    "Tanana, AK 99777",
    "Tenakee Springs, AK 99841",
    "Teller, AK 99778",
    "Tobyhanna, PA 18466",
    "Thorne Bay, AK 99919",
    "Togiak, AK 99678",
    "Tok, AK 99780",
    "Toksook Bay, AK 99637",
    "Trapper Creek, AK 99683",
    "Tuluksak, AK 99679",
    "Tooele, UT 84074",
    "Tripler Army Medical Center, HI 96859",
    "Tuntutuliak, AK 99680",
    "Tununak, AK 99681",
    "Volcano, HI 96785",
    "Tyonek, AK 99682",
    "Unalakleet, AK 99684",
    "Unalaska, AK 99685",
    "Valdez, AK 99686",
    "Two Rivers, AK 99716",
    "Usaf Academy, CO 80840",
    "Usaf Academy, CO 80841",
    "Venetie, AK 99781",
    "Ward Cove, AK 99928",
    "Wahiawa, HI 96786",
    "Waialua, HI 96791",
    "Waianae, HI 96792",
    "Waikoloa, HI 96738",
    "Wailuku, HI 96793",
    "Waimanalo, HI 96795",
    "Waimea, HI 96796",
    "Wainwright, AK 99782",
    "Waipahu, HI 96797",
    "Wake Island, HI 96898",
    "Wales, AK 99783",
    "Wrangell, AK 99929",
    "Wasilla, AK 99623",
    "Wasilla, AK 99629",
    "Wasilla, AK 99654",
    "Wasilla, AK 99687",
    "Wheeler Army Airfield, HI 96854",
    "Whittier, AK 99693",
    "Willow, AK 99688",
    "Yakutat, AK 99689",
    "Whiteman Air Force Base, MO 65305",
    "White Mountain, AK 99784",
]

merges = [
    ("Aberdeen Proving Ground", "Aberdeen Proving Ground, MD 21005"),
    ("Altus AFB", "Altus AFB, OK 73523"),
    ("Barksdale AFB", "Barksdale AFB, LA 71110"),
    ("Beale AFB", "Beale AFB, CA 95903"),
    ("Camp Lejeune", "Camp Lejeune, NC 28542"),
    ("Camp Pendleton", "Camp Pendleton, CA 92055"),
    ("Cannon AFB", "Cannon AFB, NM 88103"),
    ("Dover AFB", "Dover AFB, DE 19902"),
    ("Dyess AFB", "Dyess AFB, TX 79607"),
    ("Eglin AFB", "Eglin AFB, FL 32542"),
    ("Ellsworth AFB", "Ellsworth AFB, SD 57706"),
    ("Fort Belvoir", "Fort Belvoir, VA 22060"),
    ("Fort Bliss", "Fort Bliss, TX 79916"),
    ("Fort Bragg", "Fort Bragg, CA 95437"),
    ("Fort Campbell", "Fort Campbell, KY 42223"),
    ("Fort Drum", "Fort Drum, NY 13602"),
    ("Fort Huachuca", "Fort Huachuca, AZ 85613"),
    ("Fort Irwin", "Fort Irwin, CA 92310"),
    ("Fort Knox", "Fort Knox, KY 40121"),
    ("Fort Leavenworth", "Fort Leavenworth, KS 66027"),
    ("Fort Lee", "Fort Lee, NJ 07024"),
    ("Fort Leonard Wood", "Fort Leonard Wood, MO 65473"),
    ("Fort Riley", "Fort Riley, KS 66442"),
    ("Fort Sill", "Fort Sill, OK 73503"),
    ("Fort Stewart", "Fort Stewart, GA 31314"),
    ("Goodfellow AFB", "Goodfellow AFB, TX 76908"),
    ("Grand Forks AFB", "Grand Forks AFB, ND 58204"),
    ("Hanscom AFB", "Hanscom AFB, MA 01731"),
    ("Hill AFB", "Hill AFB, UT 84056"),
    ("Hurlburt Field", "Hurlburt Field, FL 32544"),
    ("JB Langley-Eustis", "JB Langley-Eustis (Eustis)"),
    ("JBSA Lackland", "JBSA Lackland, TX 78236"),
    ("JBSA Randolph", "JBSA Randolph, TX 78150"),
    ("Kirtland AFB", "Kirtland AFB, NM 87117"),
    ("Laughlin AFB", "Laughlin AFB, TX 78843"),
    ("Malmstrom AFB", "Malmstrom AFB, MT 59402"),
    ("Minot AFB", "Minot AFB, ND 58705"),
    ("Moody AFB", "Moody AFB, GA 31699"),
    ("Mountain Home AFB", "Mountain Home AFB, ID 83648"),
    ("Nellis AFB", "Nellis AFB, NV 89191"),
    ("Offutt AFB", "Offutt AFB, NE 68113"),
    ("PENNSYLVANIA STATE UNIVERSITY", "PENNSYLVANIA STATE UNIVERSITY (NROTC)"),
    ("Shaw AFB", "Shaw AFB, SC 29152"),
    ("Sheppard AFB", "Sheppard AFB, TX 76311"),
    ("Tinker AFB", "Tinker AFB, OK 73145"),
    ("Travis AFB", "Travis AFB, CA 94535"),
    ("Twentynine Palms", "Twentynine Palms, CA 92277"),
    ("Washington Navy Yard", "Washington Navy Yard, DC 20388"),
    ("West Point", "West Point, NY 10996"),
    ("White Sands Missile Range", "White Sands Missile Range, NM 88002"),
]

renames = [
    ("Fort Gordon", "Fort Eisenhower"),
    ("Fort Gordon, GA 30813", "Fort Eisenhower, GA 30813"),
    ("Fort Hood", "Fort Cavazos"),
    ("Fort Hood, TX 76544", "Fort Cavazos, TX 76544"),
    ("Fort Rucker", "Fort Novosel"),
    ("Fort Rucker, AL 36362", "Fort Novosel, AL 36362"),
    ("Fort Benning", "Fort Walker"),
    ("Fort Benning, GA 31905", "Fort Walker, GA 31905"),
    ("Fort Benning, GA 31995", "Fort Walker, GA 31995"),
    ("Fort Bragg", "Fort Liberty"),
    ("Fort Bragg, NC 28307", "Fort Liberty, NC 28307"),
    ("Fort Bragg, NC 28310", "Fort Liberty, NC 28310"),
    ("Fort Lee", "Fort Gregg-Adams"),
    ("Fort Lee, VA 23801 ", "Fort Gregg-Adams, VA 23801"),
    ("Fort Polk", "Fort Johnson"),
    ("Fort Polk, LA 71459", "Fort Johnson, LA 71459"),
]

# just insert news from the list and do nothing on name conflict

for delete in deletes:
    # e.g. "Adak, AK 99546",
    dl_id = f"(SELECT id FROM duty_locations WHERE name = '{delete}')"
    f.write(
        f"DELETE from orders where origin_duty_location_id = {dl_id} or new_duty_location_id = {dl_id};\n"
    )
    f.write(f"DELETE from service_members where duty_location_id = {dl_id};\n")
    f.write(f"DELETE from duty_location_names where duty_location_id = {dl_id};\n")
    f.write(f"DELETE from duty_locations where id = {dl_id};\n\n")

for merge in merges:
    old = merge[0]
    new = merge[1]
    old_dl_id = f"(SELECT id from duty_locations WHERE name = '{old}')"
    new_dl_id = f"(SELECT id from duty_locations WHERE name = '{old}')"

for rename in renames:
    f.write(
        f"UPDATE duty_locations SET name = '{rename[1]}' WHERE name = '{rename[0]}';\n"
    )


# takes a comma-separated string of aliases and an id
# if the id is not null, it removes existing duty_location_names and adds new ones
# otherwise, it adds them
def handle_aliases(aliases, id):
    # remove aliases from existing
    if id is not None:
        f.write("-- Remove existing associated duty_location_names\n")
        f.write(f"DELETE FROM duty_location_names WHERE duty_location_id = '{id}'")
    if not pd.isna(aliases):
        aliases = aliases.split(r"\w?,\w?")
        for alias in aliases:
            dln_id = uuid.uuid4()
            f.write("-- Insert new duty_location_names\n")
            f.write(
                f"INSERT INTO duty_location_names (id, name, duty_location_id, created_at, updated_at) VALUES('{dln_id}', '{alias}', '{id}', now(), now());"
            )


# f.write("-- Remove unique index from duty_location_names\n")
# f.write("DROP INDEX IF EXISTS duty_location_names_name_idx;\n\n")

# df = pd.read_excel(
#     pd.ExcelFile(sys.argv[1]), keep_default_na=False, na_values=["", "nan"]
# )
# df = df.reset_index()

# for index, row in df.iterrows():
#     name = row["Duty Location Name"]
#     gbloc = row["gbloc"]
#     affiliation = row["Affiliation"]
#     street_address = row["Street Address"]
#     city = row["City"]
#     state = row["State"]
#     postal_code = "%05d" % row["Postal Code"]
#     to_id = "NULL"
#     provides_services_counseling = "TRUE"
#     aliases = row["Alias"]

#     address_id = uuid.uuid4()
#     id = uuid.uuid4()

#     all_names.append(name)

#     precomma_name = name.split(", ")[0]

#     if pd.isna(street_address):
#         street_address = ""

#     if pd.isna(affiliation):
#         affiliation = "NULL"
#     else:
#         affiliation = "_".join([w.upper() for w in affiliation.split()])
#         affiliation = f"'{affiliation}'"

#     address_query = f"(SELECT id from addresses where street_address_1 = '{street_address}' and city = '{city}' and state = '{state}' and postal_code = '{postal_code}' LIMIT 1)"
#     to_query = f"(SELECT t.id FROM transportation_offices AS t, addresses AS a WHERE t.address_id = a.id AND t.gbloc='{gbloc}' AND a.city='{city}' AND a.state='{state}' LIMIT 1)"

#     f.write(
#         f"INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at) VALUES ('{address_id}', '{street_address}', '{city}', '{state}', '{postal_code}', now(), now());\n"
#     )

#     f.write(
#         "INSERT INTO duty_locations (id, address_id, transportation_office_id, name, affiliation, provides_services_counseling, updated_at, created_at)\n"
#     )
#     f.write(
#         f"VALUES ('{id}', {address_query}, {to_query}, '{name}', {affiliation}, TRUE, now(), now())\n"
#     )
#     f.write(f"ON CONFLICT (name) DO NOTHING\n")

#     f.write(
#         f"UPDATE SET affiliation = {affiliation}, transportation_office_id = {to_query}, address_id = '{address_id}';\n\n"
#     )
f2.write("','".join(all_names))
f2.close()
f.close()
sys.exit()
