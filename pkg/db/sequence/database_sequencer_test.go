package sequence

const testSequence = "test_sequence"

func (suite *SequenceSuite) TestDatabaseNextVal() {
	sequencer := NewDatabaseSequencer(testSequence)
	nextVal, err := sequencer.NextVal(suite.AppContextForTest())
	suite.NoError(err, "Error getting next value of sequence")
	suite.Equal(int64(2), nextVal)
}

func (suite *SequenceSuite) TestDatabaseSetVal() {
	sequencer := NewDatabaseSequencer(testSequence)
	err := sequencer.SetVal(suite.AppContextForTest(), 30)
	suite.NoError(err, "Error setting value of sequence")

	nextVal, err := sequencer.NextVal(suite.AppContextForTest())
	suite.NoError(err, "Error getting next value of sequence")
	suite.Equal(int64(31), nextVal)
}
