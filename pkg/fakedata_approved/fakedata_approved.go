/******
 This file is a static data representation of the Fake names and street addresses 2019-10-15 file
https://docs.google.com/spreadsheets/d/1u1NO_ZWvKJc2ylOSF5-4mcm6Eg5X2zu7c_P-X4lDrE4/edit#gid=521176896

 The fake data from this file is the only approved name, address, phone, and email data allowed in our system
 for testing purposes. Mostly likely to be used in experimental (exp) or staging (stg).
 ******/

package fakedata_approved

import (
	"regexp"
	"strings"
)
type fakeName struct {
	first string
	last string
}

var fakeNames = []fakeName{
	{
		first: "Jason",
		last: "Ash",
	},
	{
		first: "Riley",
		last: "Baker",
	},
	{
		first: "Aaliyah",
		last: "Banks",
	},
	{
		first: "Ashley",
		last: "Banks",
	},
	{
		first: "Angel",
		last: "Bauer",
	},
	{
		first: "Jaime",
		last: "Childers",
	},
	{
		first: "Sofía",
		last: "Clark-Nuñez",
	},
	{
		first: "Justice",
		last: "Connelly",
	},
	{
		first: "Zoya",
		last: "Darvish",
	},
	{
		first: "Reese",
		last: "Embry",
	},
	{
		first: "Robin",
		last: "Fenstermacher",
	},
	{
		first: "Grace",
		last: "Griffin",
	},
	{
		first: "Laura Jane",
		last: "Henderson",
	},
	{
		first: "Skyler",
		last: "Hunt",
	},
	{
		first: "Jayden",
		last: "Jackson Jr.",
	},
	{
		first: "Dorothy",
		last: "Lagomarsino",
	},
	{
		first: "John",
		last: "Lee",
	},
	{
		first: "Jonathan",
		last: "Lee",
	},
	{
		first: "Lisa",
		last: "Lee",
	},
	{
		first: "Susan",
		last: "Lee",
	},
	{
		first: "W. Nathan",
		last: "Millering",
	},
	{
		first: "Owen",
		last: "Nance",
	},
	{
		first: "Avery",
		last: "O'Keefe",
	},
	{
		first: "Quinn",
		last: "Ocampo",
	},
	{
		first: "Josh",
		last: "Perez",
	},
	{
		first: "Jody",
		last: "Pitkin",
	},
	{
		first: "Saqib",
		last: "Rahman",
	},
	{
		first: "Carol",
		last: "Romilly",
	},
	{
		first: "James",
		last: "Rye",
	},
	{
		first: "Gabriela",
		last: "Sáenz Perez",
	},
	{
		first: "Jessica",
		last: "Smith",
	},
	{
		first: "Kerry",
		last: "Smith",
	},
	{
		first: "Ted",
		last: "Smith",
	},
	{
		first: "Barbara",
		last: "St. Juste",
	},
	{
		first: "Christopher",
		last: "Swinglehurst-Walters",
	},
	{
		first: "Melissa",
		last: "Taylor",
	},
	{
		first: "Edgar",
		last: "Taylor III",
	},
	{
		first: "Casey",
		last: "Thompson",
	},
	{
		first: "Gregory",
		last: "Van der Heide",
	},
	{
		first: "Catalina",
		last: "Washington",
	},
	{
		first: "Rosalie",
		last: "Wexler",
	},
	{
		first: "Nevaeh",
		last: "Wilson",
	},
	{
		first: "Peyton",
		last: "Wing",
	},
	{
		first: "Jo",
		last: "Xi",
	},
	{
		first: "Earl",
		last: "Yazzie",
	},
}

var fakeAddress = []string{
	"7 Q St",
	"17 8th St",
	"9 W 2nd Ave",
	"148 S East St",
	"412 Avenue M #3E",
	"10642 N Second Ave",
	"812 S 129th Street",
	"448 Washington Blvd NE",
	"4124 Apache Dr, Apt 18C",
	"6622 Airport Way S #1430",
	"235 Prospect Valley Road SE",
	"142 E Barrel Hoop Circle #4A",
	"441 SW Río de la Plata Drive",
	"3400 E Del Ray Place, 2nd Floor",
	"3373 NW Martin Luther King Jr Blvd",
	"1292 Orchard Terrace, Building C, Unit 10",
}

/*
 IsValidFakeDataFullNameStrict - checks the first name and last name can be found in the
 fake data. If the name is found true is returned, if not found, false is returned.
 */
func IsValidFakeDataFullNameStrict (firstName string, lastName string) (bool, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	for _, fake := range fakeNames {
		if strings.EqualFold(fake.first, strings.TrimSpace(firstName)) {
			if strings.EqualFold(fake.last, strings.TrimSpace(lastName)) {
				return true, nil
			}
		}
	}
	return false, nil
}

/*
 IsValidFakeDataFullName - checks the first name and last name can be found in the
 fake data. If the name is found true is returned, if not found, false is returned.
*/
func IsValidFakeDataFullName (firstName string, lastName string) (bool, error) {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return false, err
	}

	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	processedFirst := reg.ReplaceAllString(firstName, "")
	processedLast := reg.ReplaceAllString(lastName, "")

	for _, fake := range fakeNames {
		processedFakeFirst := reg.ReplaceAllString(fake.first, "")
		if strings.EqualFold(processedFakeFirst, processedFirst) {
			processedFakeLast := reg.ReplaceAllString(fake.last, "")
			if strings.EqualFold(processedFakeLast, processedLast) {
				return true, nil
			}
		}
	}
	return false, nil
}

/*
 IsValidFakeDataAddress - checks the that the address can be found in the
 fake data. If the address is found true is returned, if not found, false is returned.
*/
func IsValidFakeDataAddress (address string) (bool, error) {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return false, err
	}

	address = strings.TrimSpace(address)
	processed := reg.ReplaceAllString(address, "")
	for _, fake := range fakeAddress {
		processedFake := reg.ReplaceAllString(fake, "")
		if strings.EqualFold(processedFake, processed) {
			return true, nil
		}
	}
	return false, nil
}

/*
 IsValidFakeDataAddressStrict - checks the that the address can be found in the
 fake data. If the address is found true is returned, if not found, false is returned.
*/
func IsValidFakeDataAddressStrict (address string) (bool, error) {
	address = strings.TrimSpace(address)
	for _, fake := range fakeAddress {
		if strings.EqualFold(fake, address) {
			return true, nil
		}
	}
	return false, nil
}

/*
IsValidFakeDataPhone - checks for the format
 "999-999-999" or
 "###-555-####"
 */
func IsValidFakeDataPhone (phone string) (bool, error) {
	// Make a Regex to say we only want numbers
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return false, err
	}

	phone = strings.TrimSpace(phone)

	processed := reg.ReplaceAllString(phone, "")
	if processed == "9999999999" {
		return true, nil
	}
	if processed[3:6] == "555" {
		return true, nil
	}
	return false, nil
}

/*
IsValidFakeDataEmail - checks for the format
@example.com or
@email.com
 */
func IsValidFakeDataEmail (email string) (bool, error) {
	email = strings.TrimSpace(email)
	lowerPhone := strings.ToLower(email)

	if strings.HasSuffix(lowerPhone, "@example.com") {
		return true, nil
	}

	if strings.HasSuffix(lowerPhone, "@email.com") {
		return true, nil
	}

	return false, nil
}
