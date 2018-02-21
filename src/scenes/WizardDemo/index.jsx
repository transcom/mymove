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
      <img src={intro} />
    </div>
    <div>
      <img src={moveType} />
    </div>
    <div>
      <img src={dateSelection} />
    </div>
    <div>
      <img src={mover} />
    </div>
    <div>
      <img src={review} />
    </div>
  </Wizard>
);

export default WizardDemo;
