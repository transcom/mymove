import React from 'react';
import Wizard from 'shared/Wizard';
//  import PropTypes from 'prop-types';
import intro from './intro.png';
import moveType from './select-move-type.png';
import dateSelection from './select-date.png';
import mover from './select-mover.png';
import review from './review-locations.png';
const WizardDemo = props => (
  <Wizard>
    <div>
      <img src={intro} alt="intro" />
    </div>
    <div>
      <img src={moveType} alt="move type" />
    </div>
    <div>
      <img src={dateSelection} alt="select a date" />
    </div>
    <div>
      <img src={mover} alt="select a mover" />
    </div>
    <div>
      <img src={review} alt="review" />
    </div>
  </Wizard>
);

export default WizardDemo;
