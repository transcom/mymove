import React from 'react';
import Wizard from 'shared/Wizard';
//  import PropTypes from 'prop-types';

const WizardDemo = props => (
  <Wizard prop1="foo" funcProp={() => console.log('whee!')}>
    <div> This is page 1</div>
    <div> this is page 2</div>
    <div> this is page 3</div>
  </Wizard>
);

export default WizardDemo;
