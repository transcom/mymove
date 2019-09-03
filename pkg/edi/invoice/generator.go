package ediinvoice

// ICNSequenceName used to query Interchange Control Numbers from DB
const ICNSequenceName = "interchange_control_number"

// ICNRandomMin is the smallest allowed random-number based ICN (we use random ICN numbers in development)
const ICNRandomMin int64 = 100000000

// ICNRandomMax is the largest allowed random-number based ICN (we use random ICN numbers in development)
const ICNRandomMax int64 = 999999999
