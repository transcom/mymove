package main

import (
	"bitbucket.org/pop-book/models" // Need to sort out the import for these
	"github.com/markbates/pop"
	"log"
)

// To do: figure out how to validate this example. Then validate one. Create one more, for good measure. Then: do the next model.
func main() {
	tx, err := pop.Connect("traffic_distribution_lists")
	if err != nil {
		log.Panic(err)
	}
	tdl1 := models.TrafficDistributionList{ID: "114", SourceRateArea: "Area 1", DestinationArea: "Southwest US", CodeOfService: "135"}
	_, err = tx.ValidateAndSave(&tdl1)
	if err != nil {
		log.Panic(err)
	}
	// luke := models.User{Title: "Mr.", FirstName: "Luke", LastName: "Cage", Bio: "Hero for hire."}
	// _, err = tx.ValidateAndSave(&luke)
	// if err != nil {
	//     log.Panic(err)
	// }
	// danny := models.User{Title: "Mr.", FirstName: "Danny", LastName: "Rand", Bio: "Martial artist."}
	// _, err = tx.ValidateAndSave(&danny)
	// if err != nil {
	//     log.Panic(err)
	// }
	// matt := models.User{Title: "Mr.", FirstName: "Matthew", LastName: "Murdock", Bio: "Lawyer, sees with no eyes."}
	// _, err = tx.ValidateAndSave(&matt)
	// if err != nil {
	//     log.Panic(err)
	// }
	// frank := models.User{Title: "Mr.", FirstName: "Frank", LastName: "Castle", Bio: "USMC, badass."}
	// _, err = tx.ValidateAndSave(&frank)
	// if err != nil {
	//     log.Panic(err)
	// }
}
