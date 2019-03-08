import React from 'react';
import ComboButton from 'shared/ComboButton/index.jsx';

const StyleGuide = () => (
  <div style={{ 'margin-left': '20px' }}>
    <h1>Wizard Styles</h1>
    <hr />
    <h2>This is a H2</h2>
    <h3>This is a H3</h3>
    <div>
      <ComboButton buttonText={'Approve'} disabled={false} />
    </div>
  </div>
);

export default StyleGuide;
