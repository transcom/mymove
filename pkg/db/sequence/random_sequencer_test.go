package sequence

const testMin int64 = 10
const testMax int64 = 99

func (suite *SequenceSuite) TestRandomNextVal() {
	sequencer, err := NewRandomSequencer(testMin, testMax)
	suite.FatalNoError(err)

	// Try 100 random numbers and make sure they're all between min and max, inclusive.
	for i := 0; i < 100; i++ {
		nextVal, err := sequencer.NextVal()
		suite.NoError(err, "Error getting next value of sequence")

		suite.True(nextVal >= testMin && nextVal <= testMax,
			"NextVal returned %d, but range was %d to %d inclusive", nextVal, testMin, testMax)
	}
}

func (suite *SequenceSuite) TestRandomSetVal() {
	sequencer, err := NewRandomSequencer(testMin, testMax)
	suite.FatalNoError(err)

	// This should be a no-op; just make sure it doesn't throw any errors.
	err = sequencer.SetVal(30)
	suite.NoError(err, "Error setting value of sequence")
}

func (suite *SequenceSuite) TestRandomConstructorErrors() {
	sequencer, err := NewRandomSequencer(-3, 10)
	suite.Nil(sequencer)
	suite.Error(err)

	sequencer, err = NewRandomSequencer(20, 10)
	suite.Nil(sequencer)
	suite.Error(err)
}
